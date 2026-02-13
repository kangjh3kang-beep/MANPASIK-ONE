// Package memory는 user-service의 인메모리 저장소 구현입니다.
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/user-service/internal/service"
)

// ProfileRepository는 인메모리 프로필 저장소입니다.
type ProfileRepository struct {
	mu       sync.RWMutex
	profiles map[string]*service.UserProfile // key: userID
}

// NewProfileRepository는 인메모리 ProfileRepository를 생성합니다.
func NewProfileRepository() *ProfileRepository {
	return &ProfileRepository{
		profiles: make(map[string]*service.UserProfile),
	}
}

func (r *ProfileRepository) GetByID(_ context.Context, userID string) (*service.UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.profiles[userID]
	if !ok {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

func (r *ProfileRepository) Update(_ context.Context, profile *service.UserProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[profile.UserID] = profile
	return nil
}

// Seed는 테스트/개발용으로 프로필을 추가합니다.
func (r *ProfileRepository) Seed(profile *service.UserProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[profile.UserID] = profile
}
