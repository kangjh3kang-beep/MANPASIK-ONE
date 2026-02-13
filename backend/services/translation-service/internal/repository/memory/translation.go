// Package memory는 translation-service의 인메모리 저장소를 구현합니다.
package memory

import (
	"context"
	"sync"

	"github.com/manpasik/backend/services/translation-service/internal/service"
)

// TranslationRepository는 인메모리 번역 이력 저장소입니다.
type TranslationRepository struct {
	mu      sync.RWMutex
	records []*service.TranslationRecord
}

// NewTranslationRepository는 TranslationRepository를 생성합니다.
func NewTranslationRepository() *TranslationRepository {
	return &TranslationRepository{}
}

func (r *TranslationRepository) Save(_ context.Context, record *service.TranslationRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, record)
	return nil
}

func (r *TranslationRepository) FindByUserID(_ context.Context, userID string, limit, offset int) ([]*service.TranslationRecord, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.TranslationRecord
	for _, rec := range r.records {
		if rec.UserID == userID {
			filtered = append(filtered, rec)
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

// UsageRepository는 인메모리 사용량 저장소입니다.
type UsageRepository struct {
	mu    sync.RWMutex
	usage map[string]*service.UsageStats
}

// NewUsageRepository는 UsageRepository를 생성합니다.
func NewUsageRepository() *UsageRepository {
	return &UsageRepository{
		usage: make(map[string]*service.UsageStats),
	}
}

func (r *UsageRepository) IncrementUsage(_ context.Context, userID string, characters int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stats, ok := r.usage[userID]
	if !ok {
		stats = &service.UsageStats{MonthlyLimit: 100000} // 기본 월 10만자
		r.usage[userID] = stats
	}

	stats.TotalCharacters += int64(characters)
	stats.MonthlyCharacters += int64(characters)
	stats.TotalRequests++
	stats.MonthlyRequests++

	return nil
}

func (r *UsageRepository) GetUsage(_ context.Context, userID string) (*service.UsageStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats, ok := r.usage[userID]
	if !ok {
		return &service.UsageStats{MonthlyLimit: 100000}, nil
	}
	return stats, nil
}
