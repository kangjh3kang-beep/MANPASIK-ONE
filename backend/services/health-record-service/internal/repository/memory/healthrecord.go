// Package memory는 health-record-service의 인메모리 저장소입니다.
package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/health-record-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// HealthRecordRepository는 건강 기록 인메모리 저장소입니다.
type HealthRecordRepository struct {
	mu    sync.RWMutex
	store map[string]*service.HealthRecord
}

// NewHealthRecordRepository는 새 인메모리 건강 기록 저장소를 생성합니다.
func NewHealthRecordRepository() *HealthRecordRepository {
	return &HealthRecordRepository{
		store: make(map[string]*service.HealthRecord),
	}
}

// Save는 건강 기록을 저장합니다.
func (r *HealthRecordRepository) Save(_ context.Context, rec *service.HealthRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[rec.ID] = rec
	return nil
}

// FindByID는 ID로 건강 기록을 조회합니다.
func (r *HealthRecordRepository) FindByID(_ context.Context, id string) (*service.HealthRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.store[id]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	return rec, nil
}

// FindByUserID는 사용자 건강 기록 목록을 조회합니다.
func (r *HealthRecordRepository) FindByUserID(_ context.Context, userID string, typeFilter service.HealthRecordType, startDate, endDate *time.Time, limit, offset int) ([]*service.HealthRecord, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.HealthRecord
	for _, rec := range r.store {
		if rec.UserID != userID {
			continue
		}
		if typeFilter != service.RecordTypeUnknown && rec.RecordType != typeFilter {
			continue
		}
		if startDate != nil && rec.RecordedAt.Before(*startDate) {
			continue
		}
		if endDate != nil && rec.RecordedAt.After(*endDate) {
			continue
		}
		filtered = append(filtered, rec)
	}

	// 시간 역순 정렬
	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].RecordedAt.After(filtered[i].RecordedAt) {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

// Update는 건강 기록을 업데이트합니다.
func (r *HealthRecordRepository) Update(_ context.Context, rec *service.HealthRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[rec.ID]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	r.store[rec.ID] = rec
	return nil
}

// Delete는 건강 기록을 삭제합니다.
func (r *HealthRecordRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return apperrors.New(apperrors.ErrNotFound, "리소스를 찾을 수 없습니다")
	}
	delete(r.store, id)
	return nil
}
