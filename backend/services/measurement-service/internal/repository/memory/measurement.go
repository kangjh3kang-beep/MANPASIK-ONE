package memory

import (
	"context"
	"sync"
	"time"

	"github.com/manpasik/backend/services/measurement-service/internal/service"
)

// MeasurementRepository는 인메모리 측정 데이터 저장소입니다.
type MeasurementRepository struct {
	mu   sync.RWMutex
	data []*service.MeasurementData
}

// NewMeasurementRepository는 인메모리 MeasurementRepository를 생성합니다.
func NewMeasurementRepository() *MeasurementRepository {
	return &MeasurementRepository{}
}

func (r *MeasurementRepository) Store(_ context.Context, data *service.MeasurementData) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = append(r.data, data)
	return nil
}

func (r *MeasurementRepository) GetHistory(_ context.Context, userID string, start, end time.Time, limit, offset int) ([]*service.MeasurementSummary, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.MeasurementSummary
	for _, d := range r.data {
		if d.UserID != userID {
			continue
		}
		if !start.IsZero() && d.Time.Before(start) {
			continue
		}
		if !end.IsZero() && d.Time.After(end) {
			continue
		}
		filtered = append(filtered, &service.MeasurementSummary{
			SessionID:     d.SessionID,
			CartridgeType: d.CartridgeType,
			PrimaryValue:  d.PrimaryValue,
			Unit:          d.Unit,
			MeasuredAt:    d.Time,
		})
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	endIdx := offset + limit
	if endIdx > total {
		endIdx = total
	}
	return filtered[offset:endIdx], total, nil
}
