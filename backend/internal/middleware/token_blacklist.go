package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklist provides token revocation support via Redis
type TokenBlacklist interface {
	// Revoke adds a token to the blacklist with its remaining TTL
	Revoke(ctx context.Context, tokenID string, expiry time.Duration) error
	// IsRevoked checks if a token has been revoked
	IsRevoked(ctx context.Context, tokenID string) bool
}

// redisTokenBlacklist implements TokenBlacklist using Redis
type redisTokenBlacklist struct {
	client *redis.Client
	prefix string
}

// NewTokenBlacklist creates a new Redis-based token blacklist
func NewTokenBlacklist(client *redis.Client) TokenBlacklist {
	return &redisTokenBlacklist{
		client: client,
		prefix: "token:blacklist:",
	}
}

// Revoke adds a token to the blacklist
func (b *redisTokenBlacklist) Revoke(ctx context.Context, tokenID string, expiry time.Duration) error {
	key := fmt.Sprintf("%s%s", b.prefix, tokenID)
	return b.client.Set(ctx, key, "revoked", expiry).Err()
}

// IsRevoked checks if a token has been revoked
func (b *redisTokenBlacklist) IsRevoked(ctx context.Context, tokenID string) bool {
	key := fmt.Sprintf("%s%s", b.prefix, tokenID)
	result, err := b.client.Exists(ctx, key).Result()
	return err == nil && result > 0
}
