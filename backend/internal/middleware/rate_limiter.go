package middleware

import (
	"backend/internal/cache"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig configures rate limiting
type RateLimiterConfig struct {
	RequestsPerSecond int           // Maximum requests per second
	BurstSize         int           // Burst size for token bucket
	WindowSize        time.Duration // Time window for fixed window algorithm
	KeyPrefix         string        // Redis key prefix
}

// DefaultRateLimiterConfig returns default configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		WindowSize:        time.Minute,
		KeyPrefix:         "rate_limit",
	}
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(cacheManager *cache.CacheManager, config RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP or user ID if authenticated)
		clientID := getClientID(c)
		key := cache.GenerateKey("%s:%s:%s", config.KeyPrefix, clientID, c.Request.URL.Path)

		// Check rate limit using sliding window
		allowed, remaining, resetTime, err := checkRateLimit(c.Request.Context(), cacheManager, key, config)
		if err != nil {
			// On error, allow request but log
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerSecond))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": resetTime,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit implements sliding window rate limiting
func checkRateLimit(ctx context.Context, cacheManager *cache.CacheManager, key string, config RateLimiterConfig) (bool, int, int64, error) {
	now := time.Now().Unix()
	windowStart := now - int64(config.WindowSize.Seconds())

	// Get current count from Redis
	var count int
	err := cacheManager.Get(ctx, key, &count)
	if err != nil {
		// Key doesn't exist, create new window
		count = 0
	}

	// Check if window has expired
	windowKey := key + ":window"
	var windowStartTime int64
	err = cacheManager.Get(ctx, windowKey, &windowStartTime)
	if err != nil || windowStartTime < windowStart {
		// Window expired, reset counter
		count = 0
		windowStartTime = now
		cacheManager.SetWithTTL(ctx, windowKey, windowStartTime, config.WindowSize)
	}

	// Check if limit exceeded
	if count >= config.RequestsPerSecond {
		resetTime := windowStartTime + int64(config.WindowSize.Seconds())
		return false, 0, resetTime, nil
	}

	// Increment counter
	count++
	err = cacheManager.SetWithTTL(ctx, key, count, config.WindowSize)
	if err != nil {
		return true, config.RequestsPerSecond - count, windowStartTime + int64(config.WindowSize.Seconds()), err
	}

	remaining := config.RequestsPerSecond - count
	resetTime := windowStartTime + int64(config.WindowSize.Seconds())
	return true, remaining, resetTime, nil
}

// getClientID extracts client identifier from request
func getClientID(c *gin.Context) string {
	// Check if user is authenticated
	if userID, exists := c.Get("user_id"); exists {
		return strconv.FormatUint(userID.(uint64), 10)
	}

	// Fall back to IP address
	ip := c.ClientIP()
	return ip
}

// IPRateLimiter creates a rate limiter based on IP address
func IPRateLimiter(cacheManager *cache.CacheManager, requestsPerMinute int) gin.HandlerFunc {
	config := RateLimiterConfig{
		RequestsPerSecond: requestsPerMinute,
		BurstSize:         requestsPerMinute * 2,
		WindowSize:        time.Minute,
		KeyPrefix:         "ip_rate_limit",
	}

	return RateLimiter(cacheManager, config)
}

// UserRateLimiter creates a rate limiter based on user ID
func UserRateLimiter(cacheManager *cache.CacheManager, requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to authenticated users
		if _, exists := c.Get("user_id"); !exists {
			c.Next()
			return
		}

		config := RateLimiterConfig{
			RequestsPerSecond: requestsPerMinute,
			BurstSize:         requestsPerMinute * 2,
			WindowSize:        time.Minute,
			KeyPrefix:         "user_rate_limit",
		}

		RateLimiter(cacheManager, config)(c)
	}
}

// PublicAPIRateLimiter rate limiter for public APIs
func PublicAPIRateLimiter(cacheManager *cache.CacheManager) gin.HandlerFunc {
	return IPRateLimiter(cacheManager, 30) // 30 requests per minute
}

// AuthenticatedAPIRateLimiter rate limiter for authenticated APIs
func AuthenticatedAPIRateLimiter(cacheManager *cache.CacheManager) gin.HandlerFunc {
	return UserRateLimiter(cacheManager, 60) // 60 requests per minute
}

// StrictRateLimiter strict rate limiter for sensitive operations
func StrictRateLimiter(cacheManager *cache.CacheManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		config := RateLimiterConfig{
			RequestsPerSecond: 5,
			BurstSize:         5,
			WindowSize:        time.Minute,
			KeyPrefix:         "strict_rate_limit",
		}

		RateLimiter(cacheManager, config)(c)
	}
}
