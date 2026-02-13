// Package memory는 인메모리 UserRepository 구현입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/auth-service/internal/service"
)

// UserRepository는 인메모리 사용자 저장소입니다.
type UserRepository struct {
	mu      sync.RWMutex
	byID    map[string]*service.User
	byEmail map[string]*service.User
}

// NewUserRepository는 UserRepository를 생성합니다.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:    make(map[string]*service.User),
		byEmail: make(map[string]*service.User),
	}
}

// GetByID는 ID로 사용자를 조회합니다.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*service.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

// GetByEmail은 이메일로 사용자를 조회합니다.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, nil
	}
	cp := *u
	return &cp, nil
}

// Create는 사용자를 생성합니다.
func (r *UserRepository) Create(ctx context.Context, user *service.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *user
	r.byID[user.ID] = &cp
	r.byEmail[user.Email] = &cp
	return nil
}

// UpdatePassword는 비밀번호를 업데이트합니다.
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.byID[id]; ok {
		u.HashedPassword = hashedPassword
	}
	return nil
}
