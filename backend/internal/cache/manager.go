package cache

import (
	"backend/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheManager provides high-level caching operations
type CacheManager struct {
	client     *RedisClient
	prefix     string
	defaultTTL time.Duration
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cfg config.RedisConfig, prefix string, defaultTTL time.Duration) (*CacheManager, error) {
	client, err := New(cfg)
	if err != nil {
		return nil, err
	}

	return &CacheManager{
		client:     client,
		prefix:     prefix,
		defaultTTL: defaultTTL,
	}, nil
}

// Get retrieves a value from cache
func (c *CacheManager) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := c.fullKey(key)
	data, err := c.client.GetClient().Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Set stores a value in cache with default TTL
func (c *CacheManager) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetWithTTL(ctx, key, value, c.defaultTTL)
}

// SetWithTTL stores a value with custom TTL
func (c *CacheManager) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := c.fullKey(key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.GetClient().Set(ctx, fullKey, data, ttl).Err()
}

// Delete removes a key from cache
func (c *CacheManager) Delete(ctx context.Context, key string) error {
	fullKey := c.fullKey(key)
	return c.client.GetClient().Del(ctx, fullKey).Err()
}

// DeletePattern removes keys matching a pattern
func (c *CacheManager) DeletePattern(ctx context.Context, pattern string) error {
	fullPattern := c.fullKey(pattern)
	keys, err := c.client.GetClient().Keys(ctx, fullPattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.GetClient().Del(ctx, keys...).Err()
	}
	return nil
}

// Exists checks if a key exists
func (c *CacheManager) Exists(ctx context.Context, key string) bool {
	fullKey := c.fullKey(key)
	n, err := c.client.GetClient().Exists(ctx, fullKey).Result()
	return err == nil && n > 0
}

// GetOrSet gets value from cache or sets it using the provided function
func (c *CacheManager) GetOrSet(ctx context.Context, key string, dest interface{}, fn func() (interface{}, error)) error {
	// Try to get from cache
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// Cache miss, execute function
	value, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := c.Set(ctx, key, value); err != nil {
		// Log error but don't fail the request
		log.Printf("[WARN] Failed to cache value: %v", err)
	}

	// Set the result
	data, _ := json.Marshal(value)
	return json.Unmarshal(data, dest)
}

// Increment atomically increments a counter
func (c *CacheManager) Increment(ctx context.Context, key string) (int64, error) {
	fullKey := c.fullKey(key)
	return c.client.GetClient().Incr(ctx, fullKey).Result()
}

// Expire sets expiration on a key
func (c *CacheManager) Expire(ctx context.Context, key string, ttl time.Duration) error {
	fullKey := c.fullKey(key)
	return c.client.GetClient().Expire(ctx, fullKey, ttl).Err()
}

// FlushAll clears all cache (use with caution)
func (c *CacheManager) FlushAll(ctx context.Context) error {
	return c.client.GetClient().FlushAll(ctx).Err()
}

// fullKey returns the full key with prefix
func (c *CacheManager) fullKey(key string) string {
	if c.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Close closes the cache manager
func (c *CacheManager) Close() error {
	return c.client.Close()
}

// Cache keys for common entities
const (
	CacheKeyCruise       = "cruise:%s"
	CacheKeyCruiseList   = "cruises:list:%s"
	CacheKeyVoyage       = "voyage:%s"
	CacheKeyVoyageList   = "voyages:list:%s"
	CacheKeyCabinType    = "cabin_type:%s"
	CacheKeyPrice        = "price:%s:%s"
	CacheKeyInventory    = "inventory:%s:%s"
	CacheKeyOrder        = "order:%s"
	CacheKeyUser         = "user:%d"
	CacheKeyUserProfile  = "user_profile:%d"
	CacheKeyPopular      = "popular:%s"
	CacheKeySearchResult = "search:%s"
)

// GenerateKey generates a cache key from format and args
func GenerateKey(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// TTL constants
const (
	TTLShort  = 5 * time.Minute
	TTLMedium = 30 * time.Minute
	TTLLong   = 2 * time.Hour
	TTLDay    = 24 * time.Hour
	TTLWeek   = 7 * 24 * time.Hour
)
