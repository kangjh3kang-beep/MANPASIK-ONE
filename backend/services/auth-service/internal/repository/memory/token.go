// Package memory는 인메모리 TokenRepository 구현입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"
	"time"
)

// TokenRepository는 인메모리 Refresh Token 저장소입니다.
type TokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]time.Time // "userID:tokenID" -> expiresAt
}

// NewTokenRepository는 TokenRepository를 생성합니다.
func NewTokenRepository() *TokenRepository {
	return &TokenRepository{tokens: make(map[string]time.Time)}
}

// StoreRefreshToken은 Refresh Token을 저장합니다.
func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID, tokenID string, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[userID+":"+tokenID] = time.Now().UTC().Add(ttl)
	return nil
}

// ValidateRefreshToken은 토큰 유효성을 검사합니다.
func (r *TokenRepository) ValidateRefreshToken(ctx context.Context, userID, tokenID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	exp, ok := r.tokens[userID+":"+tokenID]
	if !ok {
		return false, nil
	}
	if time.Now().UTC().After(exp) {
		delete(r.tokens, userID+":"+tokenID)
		return false, nil
	}
	return true, nil
}

// RevokeRefreshToken은 토큰을 철회합니다.
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, userID, tokenID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, userID+":"+tokenID)
	return nil
}

// RevokeAllUserTokens는 사용자의 모든 토큰을 철회합니다.
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key := range r.tokens {
		if len(key) > len(userID) && key[:len(userID)] == userID {
			delete(r.tokens, key)
		}
	}
	return nil
}
