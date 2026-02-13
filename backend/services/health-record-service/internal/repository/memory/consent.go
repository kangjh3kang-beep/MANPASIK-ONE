package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/health-record-service/internal/service"
	apperrors "github.com/manpasik/backend/shared/errors"
)

// ConsentRepository는 데이터 공유 동의 인메모리 저장소입니다.
type ConsentRepository struct {
	mu    sync.RWMutex
	store map[string]*service.DataSharingConsent
}

// NewConsentRepository는 새 인메모리 동의 저장소를 생성합니다.
func NewConsentRepository() *ConsentRepository {
	return &ConsentRepository{
		store: make(map[string]*service.DataSharingConsent),
	}
}

// Create는 데이터 공유 동의를 저장합니다.
func (r *ConsentRepository) Create(_ context.Context, consent *service.DataSharingConsent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[consent.ID] = consent
	return nil
}

// GetByID는 ID로 동의를 조회합니다.
func (r *ConsentRepository) GetByID(_ context.Context, consentID string) (*service.DataSharingConsent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	consent, ok := r.store[consentID]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, "동의를 찾을 수 없습니다")
	}
	return consent, nil
}

// ListByUser는 사용자의 동의 목록을 조회합니다.
func (r *ConsentRepository) ListByUser(_ context.Context, userID string) ([]*service.DataSharingConsent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*service.DataSharingConsent
	for _, c := range r.store {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

// ListByProvider는 제공자의 동의 목록을 조회합니다.
func (r *ConsentRepository) ListByProvider(_ context.Context, providerID string) ([]*service.DataSharingConsent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*service.DataSharingConsent
	for _, c := range r.store {
		if c.ProviderID == providerID {
			result = append(result, c)
		}
	}
	return result, nil
}

// Revoke는 동의를 철회합니다.
func (r *ConsentRepository) Revoke(_ context.Context, consentID string, reason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	consent, ok := r.store[consentID]
	if !ok {
		return apperrors.New(apperrors.ErrNotFound, "동의를 찾을 수 없습니다")
	}
	consent.Status = service.ConsentRevoked
	consent.RevokedAt = time.Now()
	consent.RevokeReason = reason
	return nil
}

// CheckAccess는 특정 사용자-제공자-scope 조합에 대한 접근 권한을 확인합니다.
func (r *ConsentRepository) CheckAccess(_ context.Context, userID, providerID, scope string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.store {
		if c.UserID != userID || c.ProviderID != providerID {
			continue
		}
		if c.Status != service.ConsentActive {
			continue
		}
		if time.Now().After(c.ExpiresAt) {
			continue
		}
		for _, s := range c.Scope {
			if s == scope {
				return true, nil
			}
		}
	}
	return false, nil
}

// DataAccessLogRepository는 데이터 접근 로그 인메모리 저장소입니다.
type DataAccessLogRepository struct {
	mu    sync.RWMutex
	store []*service.DataAccessLog
}

// NewDataAccessLogRepository는 새 인메모리 접근 로그 저장소를 생성합니다.
func NewDataAccessLogRepository() *DataAccessLogRepository {
	return &DataAccessLogRepository{
		store: make([]*service.DataAccessLog, 0),
	}
}

// Log는 데이터 접근 로그를 기록합니다.
func (r *DataAccessLogRepository) Log(_ context.Context, entry *service.DataAccessLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = append(r.store, entry)
	return nil
}

// ListByUser는 사용자의 접근 로그를 조회합니다.
func (r *DataAccessLogRepository) ListByUser(_ context.Context, userID string, limit, offset int) ([]*service.DataAccessLog, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.DataAccessLog
	for _, entry := range r.store {
		if entry.UserID == userID {
			filtered = append(filtered, entry)
		}
	}

	// 시간 역순 정렬
	for i := 0; i < len(filtered); i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].AccessedAt.After(filtered[i].AccessedAt) {
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
