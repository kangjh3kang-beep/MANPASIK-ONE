// Package cache는 subscription-service의 Redis 캐시 데코레이터입니다.
//
// Cache-Aside 패턴:
//   - 읽기: 캐시 hit → 반환, miss → DB 조회 후 캐시 저장
//   - 쓰기: DB 처리 후 관련 캐시 키 삭제
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/manpasik/backend/services/subscription-service/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	subByUserPrefix = "sub:user:"
	subByIDPrefix   = "sub:id:"
	subTTL          = 10 * time.Minute
)

// SubscriptionRepository는 Redis 캐시를 적용한 SubscriptionRepository 래퍼입니다.
type SubscriptionRepository struct {
	inner service.SubscriptionRepository
	rdb   *redis.Client
}

// NewSubscriptionRepository는 캐시 SubscriptionRepository를 생성합니다.
func NewSubscriptionRepository(inner service.SubscriptionRepository, rdb *redis.Client) *SubscriptionRepository {
	return &SubscriptionRepository{inner: inner, rdb: rdb}
}

// GetByUserID는 캐시를 확인한 후 DB에서 조회합니다.
func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*service.Subscription, error) {
	key := subByUserPrefix + userID

	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var sub service.Subscription
		if json.Unmarshal(data, &sub) == nil {
			return &sub, nil
		}
	}

	sub, err := r.inner.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}

	if b, e := json.Marshal(sub); e == nil {
		r.rdb.Set(ctx, key, b, subTTL)
	}
	return sub, nil
}

// GetByID는 캐시를 확인한 후 DB에서 조회합니다.
func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*service.Subscription, error) {
	key := subByIDPrefix + id

	data, err := r.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var sub service.Subscription
		if json.Unmarshal(data, &sub) == nil {
			return &sub, nil
		}
	}

	sub, err := r.inner.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}

	if b, e := json.Marshal(sub); e == nil {
		r.rdb.Set(ctx, key, b, subTTL)
	}
	return sub, nil
}

// Create는 구독을 생성하고 관련 캐시를 무효화합니다.
func (r *SubscriptionRepository) Create(ctx context.Context, sub *service.Subscription) error {
	if err := r.inner.Create(ctx, sub); err != nil {
		return err
	}
	r.rdb.Del(ctx, subByUserPrefix+sub.UserID)
	return nil
}

// Update는 구독을 업데이트하고 관련 캐시를 무효화합니다.
func (r *SubscriptionRepository) Update(ctx context.Context, sub *service.Subscription) error {
	if err := r.inner.Update(ctx, sub); err != nil {
		return err
	}
	r.rdb.Del(ctx,
		subByUserPrefix+sub.UserID,
		subByIDPrefix+sub.ID,
	)
	return nil
}
