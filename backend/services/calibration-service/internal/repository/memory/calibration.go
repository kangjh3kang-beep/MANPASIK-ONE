// Package memory는 인메모리 보정 저장소입니다 (개발/테스트용).
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/manpasik/backend/services/calibration-service/internal/service"
)

// ============================================================================
// CalibrationRepository (보정 기록 저장소)
// ============================================================================

// CalibrationRepository는 인메모리 보정 기록 저장소입니다.
type CalibrationRepository struct {
	mu      sync.RWMutex
	records []*service.CalibrationRecord
}

// NewCalibrationRepository는 인메모리 CalibrationRepository를 생성합니다.
func NewCalibrationRepository() *CalibrationRepository {
	return &CalibrationRepository{
		records: make([]*service.CalibrationRecord, 0),
	}
}

// Save는 보정 기록을 저장합니다.
func (r *CalibrationRepository) Save(_ context.Context, record *service.CalibrationRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := copyRecord(record)
	r.records = append(r.records, cp)
	return nil
}

// GetLatest는 디바이스 + 카트리지 타입별 최신 보정 기록을 반환합니다.
func (r *CalibrationRepository) GetLatest(_ context.Context, deviceID string, category, typeIndex int32) (*service.CalibrationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *service.CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID && rec.CartridgeCategory == category && rec.CartridgeTypeIndex == typeIndex {
			if latest == nil || rec.CalibratedAt.After(latest.CalibratedAt) {
				latest = rec
			}
		}
	}

	if latest == nil {
		return nil, nil
	}
	cp := copyRecord(latest)
	return cp, nil
}

// GetLatestByType은 디바이스 + 카트리지 타입 + 보정 유형별 최신 기록을 반환합니다.
func (r *CalibrationRepository) GetLatestByType(_ context.Context, deviceID string, category, typeIndex int32, calType service.CalibrationType) (*service.CalibrationRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *service.CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID &&
			rec.CartridgeCategory == category &&
			rec.CartridgeTypeIndex == typeIndex &&
			rec.CalibrationType == calType {
			if latest == nil || rec.CalibratedAt.After(latest.CalibratedAt) {
				latest = rec
			}
		}
	}

	if latest == nil {
		return nil, nil
	}
	cp := copyRecord(latest)
	return cp, nil
}

// ListByDevice는 디바이스의 보정 이력을 반환합니다 (최신순 정렬).
func (r *CalibrationRepository) ListByDevice(_ context.Context, deviceID string, limit, offset int32) ([]*service.CalibrationRecord, int32, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*service.CalibrationRecord
	for _, rec := range r.records {
		if rec.DeviceID == deviceID {
			filtered = append(filtered, rec)
		}
	}

	// 최신순 정렬
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CalibratedAt.After(filtered[j].CalibratedAt)
	})

	total := int32(len(filtered))

	// 페이지네이션
	start := int(offset)
	if start >= len(filtered) {
		return nil, total, nil
	}
	end := start + int(limit)
	if end > len(filtered) {
		end = len(filtered)
	}

	// 결과 복사
	result := make([]*service.CalibrationRecord, 0, end-start)
	for _, rec := range filtered[start:end] {
		result = append(result, copyRecord(rec))
	}

	return result, total, nil
}

// copyRecord는 CalibrationRecord의 깊은 복사본을 반환합니다.
func copyRecord(src *service.CalibrationRecord) *service.CalibrationRecord {
	cp := *src
	if src.ChannelOffsets != nil {
		cp.ChannelOffsets = make([]float64, len(src.ChannelOffsets))
		copy(cp.ChannelOffsets, src.ChannelOffsets)
	}
	if src.ChannelGains != nil {
		cp.ChannelGains = make([]float64, len(src.ChannelGains))
		copy(cp.ChannelGains, src.ChannelGains)
	}
	return &cp
}

// ============================================================================
// CalibrationModelRepository (보정 모델 저장소)
// ============================================================================

// CalibrationModelRepository는 인메모리 보정 모델 저장소입니다.
type CalibrationModelRepository struct {
	mu     sync.RWMutex
	models []*service.CalibrationModel
}

// NewCalibrationModelRepository는 인메모리 CalibrationModelRepository를 생성합니다.
// seedModels를 전달하면 초기 데이터로 등록합니다.
func NewCalibrationModelRepository(seedModels []*service.CalibrationModel) *CalibrationModelRepository {
	models := make([]*service.CalibrationModel, 0, len(seedModels))
	for _, m := range seedModels {
		cp := *m
		models = append(models, &cp)
	}
	return &CalibrationModelRepository{
		models: models,
	}
}

// GetAll은 전체 보정 모델 목록을 반환합니다.
func (r *CalibrationModelRepository) GetAll(_ context.Context) ([]*service.CalibrationModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*service.CalibrationModel, 0, len(r.models))
	for _, m := range r.models {
		cp := *m
		result = append(result, &cp)
	}
	return result, nil
}

// GetByCartridgeType은 카트리지 타입별 보정 모델을 반환합니다.
func (r *CalibrationModelRepository) GetByCartridgeType(_ context.Context, category, typeIndex int32) (*service.CalibrationModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, m := range r.models {
		if m.CartridgeCategory == category && m.CartridgeTypeIndex == typeIndex {
			cp := *m
			return &cp, nil
		}
	}
	return nil, nil
}
