package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/user-service/internal/service"
)

// SubscriptionRepository는 인메모리 구독 저장소입니다.
type SubscriptionRepository struct {
	mu   sync.RWMutex
	subs map[string]*service.Subscription // key: userID
}

// NewSubscriptionRepository는 인메모리 SubscriptionRepository를 생성합니다.
func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subs: make(map[string]*service.Subscription),
	}
}

func (r *SubscriptionRepository) GetByUserID(_ context.Context, userID string) (*service.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.subs[userID]
	if !ok {
		return nil, nil
	}
	cp := *s
	return &cp, nil
}

func (r *SubscriptionRepository) Create(_ context.Context, sub *service.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subs[sub.UserID] = sub
	return nil
}

func (r *SubscriptionRepository) Update(_ context.Context, sub *service.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subs[sub.UserID] = sub
	return nil
}
