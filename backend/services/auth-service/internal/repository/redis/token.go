// Package redis는 Redis 기반 TokenRepository 구현입니다.
package redis

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/redis/go-redis/v9"
)

// TokenRepository is a Redis-backed token repository
type TokenRepository struct {
	client *redisclient.Client
	prefix string
}

// NewTokenRepository creates a Redis token repository
func NewTokenRepository(client *redisclient.Client) *TokenRepository {
	return &TokenRepository{client: client, prefix: "auth:token:"}
}

func (r *TokenRepository) key(userID, tokenID string) string {
	return fmt.Sprintf("%s%s:%s", r.prefix, userID, tokenID)
}

func (r *TokenRepository) userPattern(userID string) string {
	return fmt.Sprintf("%s%s:*", r.prefix, userID)
}

// StoreRefreshToken은 Refresh Token을 Redis에 저장합니다.
func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID, tokenID string, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(userID, tokenID), "valid", ttl).Err()
}

// ValidateRefreshToken은 토큰 유효성을 검사합니다.
func (r *TokenRepository) ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error) {
	val, err := r.client.Get(ctx, r.key(userID, tokenID)).Result()
	if err == redisclient.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "valid", nil
}

// RevokeRefreshToken은 토큰을 철회합니다.
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, userID, tokenID string) error {
	return r.client.Del(ctx, r.key(userID, tokenID)).Err()
}

// RevokeAllUserTokens는 사용자의 모든 토큰을 철회합니다.
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	keys, err := r.client.Keys(ctx, r.userPattern(userID)).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}
