// Package memory는 인메모리 구독 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/subscription-service/internal/service"
)

// SubscriptionRepository는 인메모리 구독 저장소입니다.
type SubscriptionRepository struct {
	mu       sync.RWMutex
	byID     map[string]*service.Subscription
	byUserID map[string]*service.Subscription
}

// NewSubscriptionRepository는 인메모리 SubscriptionRepository를 생성합니다.
func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		byID:     make(map[string]*service.Subscription),
		byUserID: make(map[string]*service.Subscription),
	}
}

// GetByUserID는 사용자 ID로 구독을 조회합니다.
func (r *SubscriptionRepository) GetByUserID(_ context.Context, userID string) (*service.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.byUserID[userID]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

// GetByID는 구독 ID로 구독을 조회합니다.
func (r *SubscriptionRepository) GetByID(_ context.Context, id string) (*service.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

// Create는 구독을 생성합니다.
func (r *SubscriptionRepository) Create(_ context.Context, sub *service.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *sub
	r.byID[sub.ID] = &cp
	r.byUserID[sub.UserID] = &cp
	return nil
}

// Update는 구독을 업데이트합니다.
func (r *SubscriptionRepository) Update(_ context.Context, sub *service.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *sub
	r.byID[sub.ID] = &cp
	r.byUserID[sub.UserID] = &cp
	return nil
}

// SubscriptionHistoryRepository는 인메모리 구독 이력 저장소입니다.
type SubscriptionHistoryRepository struct {
	mu      sync.RWMutex
	entries []*service.SubscriptionHistoryEntry
}

// NewSubscriptionHistoryRepository는 인메모리 이력 저장소를 생성합니다.
func NewSubscriptionHistoryRepository() *SubscriptionHistoryRepository {
	return &SubscriptionHistoryRepository{
		entries: make([]*service.SubscriptionHistoryEntry, 0),
	}
}

// Record는 이력을 기록합니다.
func (r *SubscriptionHistoryRepository) Record(_ context.Context, entry *service.SubscriptionHistoryEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *entry
	r.entries = append(r.entries, &cp)
	return nil
}

// ListByUserID는 사용자의 이력을 조회합니다.
func (r *SubscriptionHistoryRepository) ListByUserID(_ context.Context, userID string, limit, offset int32) ([]*service.SubscriptionHistoryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*service.SubscriptionHistoryEntry
	for _, e := range r.entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}

	// 페이지네이션
	start := int(offset)
	if start >= len(result) {
		return nil, nil
	}
	end := start + int(limit)
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}
