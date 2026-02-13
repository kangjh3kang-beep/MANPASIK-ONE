// Package service는 calibration-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// CalibrationType은 보정 유형입니다.
type CalibrationType int32

const (
	CalibrationTypeUnknown CalibrationType = 0
	CalibrationTypeFactory CalibrationType = 1 // 팩토리 보정 (제조 시)
	CalibrationTypeField   CalibrationType = 2 // 현장 보정 (사용자)
	CalibrationTypeAuto    CalibrationType = 3 // 자동 보정 (AI 기반)
)

// CalibrationStatus는 보정 상태입니다.
type CalibrationStatus int32

const (
	CalibrationStatusUnknown  CalibrationStatus = 0
	CalibrationStatusValid    CalibrationStatus = 1 // 유효
	CalibrationStatusExpiring CalibrationStatus = 2 // 만료 임박 (7일 이내)
	CalibrationStatusExpired  CalibrationStatus = 3 // 만료됨
	CalibrationStatusNeeded   CalibrationStatus = 4 // 보정 필요
)

// CalibrationRecord는 보정 기록 엔티티입니다.
type CalibrationRecord struct {
	ID                   string
	DeviceID             string
	CartridgeCategory    int32
	CartridgeTypeIndex   int32
	CalibrationType      CalibrationType
	Alpha                float64
	ChannelOffsets       []float64
	ChannelGains         []float64
	TempCoefficient      float64
	HumidityCoefficient  float64
	AccuracyScore        float64 // 0~1
	ReferenceStandard    string
	CalibratedBy         string
	CalibratedAt         time.Time
	ExpiresAt            time.Time
	Status               CalibrationStatus
}

// CalibrationModel은 카트리지 타입별 기본 보정 모델입니다.
type CalibrationModel struct {
	ID                   string
	CartridgeCategory    int32
	CartridgeTypeIndex   int32
	Name                 string
	Version              string
	DefaultAlpha         float64
	ValidityDays         int32
	Description          string
	CreatedAt            time.Time
}

// CalibrationRepository는 보정 데이터 저장소 인터페이스입니다.
type CalibrationRepository interface {
	Save(ctx context.Context, record *CalibrationRecord) error
	GetLatest(ctx context.Context, deviceID string, category, typeIndex int32) (*CalibrationRecord, error)
	GetLatestByType(ctx context.Context, deviceID string, category, typeIndex int32, calType CalibrationType) (*CalibrationRecord, error)
	ListByDevice(ctx context.Context, deviceID string, limit, offset int32) ([]*CalibrationRecord, int32, error)
}

// CalibrationModelRepository는 보정 모델 저장소 인터페이스입니다.
type CalibrationModelRepository interface {
	GetAll(ctx context.Context) ([]*CalibrationModel, error)
	GetByCartridgeType(ctx context.Context, category, typeIndex int32) (*CalibrationModel, error)
}

// CalibrationService는 보정 비즈니스 로직입니다.
type CalibrationService struct {
	logger   *zap.Logger
	calRepo  CalibrationRepository
	modelRepo CalibrationModelRepository
}

// NewCalibrationService는 새 CalibrationService를 생성합니다.
func NewCalibrationService(
	logger *zap.Logger,
	calRepo CalibrationRepository,
	modelRepo CalibrationModelRepository,
) *CalibrationService {
	return &CalibrationService{
		logger:    logger,
		calRepo:   calRepo,
		modelRepo: modelRepo,
	}
}

// DefaultFactoryValidityDays는 팩토리 보정 기본 유효 기간 (일)입니다.
const DefaultFactoryValidityDays = 90

// DefaultFieldValidityDays는 현장 보정 기본 유효 기간 (일)입니다.
const DefaultFieldValidityDays = 30

// ExpiringThresholdDays는 만료 임박 판단 기준 (일)입니다.
const ExpiringThresholdDays = 7

// RegisterFactoryCalibration은 팩토리 보정 데이터를 등록합니다.
func (s *CalibrationService) RegisterFactoryCalibration(
	ctx context.Context,
	deviceID string,
	category, typeIndex int32,
	alpha float64,
	channelOffsets, channelGains []float64,
	tempCoefficient, humidityCoefficient float64,
	referenceStandard, calibratedBy string,
) (*CalibrationRecord, error) {
	if deviceID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}
	if category <= 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "유효하지 않은 카트리지 카테고리입니다")
	}

	// 보정 모델에서 유효 기간 조회 (없으면 기본값 사용)
	validityDays := int32(DefaultFactoryValidityDays)
	model, _ := s.modelRepo.GetByCartridgeType(ctx, category, typeIndex)
	if model != nil && model.ValidityDays > 0 {
		validityDays = model.ValidityDays
	}

	// 정확도 점수 계산: 기준 모델의 기본 alpha와 비교
	accuracyScore := calculateAccuracyScore(alpha, model)

	now := time.Now().UTC()
	record := &CalibrationRecord{
		ID:                  uuid.New().String(),
		DeviceID:            deviceID,
		CartridgeCategory:   category,
		CartridgeTypeIndex:  typeIndex,
		CalibrationType:     CalibrationTypeFactory,
		Alpha:               alpha,
		ChannelOffsets:      copyFloat64Slice(channelOffsets),
		ChannelGains:        copyFloat64Slice(channelGains),
		TempCoefficient:     tempCoefficient,
		HumidityCoefficient: humidityCoefficient,
		AccuracyScore:       accuracyScore,
		ReferenceStandard:   referenceStandard,
		CalibratedBy:        calibratedBy,
		CalibratedAt:        now,
		ExpiresAt:           now.AddDate(0, 0, int(validityDays)),
		Status:              CalibrationStatusValid,
	}

	if err := s.calRepo.Save(ctx, record); err != nil {
		s.logger.Error("팩토리 보정 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "팩토리 보정 등록에 실패했습니다")
	}

	s.logger.Info("팩토리 보정 등록 완료",
		zap.String("device_id", deviceID),
		zap.Int32("category", category),
		zap.Int32("type_index", typeIndex),
		zap.Float64("alpha", alpha),
	)
	return record, nil
}

// PerformFieldCalibration은 현장 보정을 수행합니다.
// Alpha = mean(measured / reference) (최소제곱 방식)
func (s *CalibrationService) PerformFieldCalibration(
	ctx context.Context,
	deviceID, userID string,
	category, typeIndex int32,
	referenceValues, measuredValues []float64,
	temperatureC, humidityPct float64,
) (*CalibrationRecord, error) {
	if deviceID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}
	if userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if len(referenceValues) == 0 || len(measuredValues) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "기준값과 측정값이 필요합니다")
	}
	if len(referenceValues) != len(measuredValues) {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "기준값과 측정값의 개수가 일치해야 합니다")
	}

	// Alpha 계산: mean(measured / reference) using least squares approach
	alpha := calculateAlphaFromCalibration(referenceValues, measuredValues)

	// 정확도 점수 계산: 보정 후 잔차 기반
	accuracyScore := calculateFieldAccuracyScore(referenceValues, measuredValues, alpha)

	// 유효 기간: 현장 보정은 기본 30일
	validityDays := int32(DefaultFieldValidityDays)
	model, _ := s.modelRepo.GetByCartridgeType(ctx, category, typeIndex)
	if model != nil && model.ValidityDays > 0 {
		// 현장 보정은 모델 유효 기간의 1/3
		fieldDays := model.ValidityDays / 3
		if fieldDays > 0 {
			validityDays = fieldDays
		}
	}

	now := time.Now().UTC()
	record := &CalibrationRecord{
		ID:                  uuid.New().String(),
		DeviceID:            deviceID,
		CartridgeCategory:   category,
		CartridgeTypeIndex:  typeIndex,
		CalibrationType:     CalibrationTypeField,
		Alpha:               alpha,
		ChannelOffsets:      nil,
		ChannelGains:        nil,
		TempCoefficient:     temperatureC,
		HumidityCoefficient: humidityPct,
		AccuracyScore:       accuracyScore,
		ReferenceStandard:   "",
		CalibratedBy:        userID,
		CalibratedAt:        now,
		ExpiresAt:           now.AddDate(0, 0, int(validityDays)),
		Status:              CalibrationStatusValid,
	}

	if err := s.calRepo.Save(ctx, record); err != nil {
		s.logger.Error("현장 보정 저장 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "현장 보정 수행에 실패했습니다")
	}

	s.logger.Info("현장 보정 완료",
		zap.String("device_id", deviceID),
		zap.String("user_id", userID),
		zap.Float64("alpha", alpha),
		zap.Float64("accuracy", accuracyScore),
	)
	return record, nil
}

// GetCalibration은 디바이스 + 카트리지 타입별 최신 유효 보정 데이터를 조회합니다.
// 팩토리 보정을 우선하며, 팩토리가 만료되었으면 현장 보정을 반환합니다.
func (s *CalibrationService) GetCalibration(
	ctx context.Context,
	deviceID string,
	category, typeIndex int32,
) (*CalibrationRecord, error) {
	if deviceID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}

	now := time.Now().UTC()

	// 1) 팩토리 보정 확인 (우선)
	factory, _ := s.calRepo.GetLatestByType(ctx, deviceID, category, typeIndex, CalibrationTypeFactory)
	if factory != nil && factory.ExpiresAt.After(now) {
		factory.Status = computeStatus(factory, now)
		return factory, nil
	}

	// 2) 현장 보정 확인
	field, _ := s.calRepo.GetLatestByType(ctx, deviceID, category, typeIndex, CalibrationTypeField)
	if field != nil && field.ExpiresAt.After(now) {
		field.Status = computeStatus(field, now)
		return field, nil
	}

	// 3) 만료된 것이라도 최신 반환 (상태는 EXPIRED)
	latest, err := s.calRepo.GetLatest(ctx, deviceID, category, typeIndex)
	if err != nil {
		s.logger.Error("보정 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "보정 데이터 조회에 실패했습니다")
	}
	if latest == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "보정 데이터가 없습니다")
	}

	latest.Status = computeStatus(latest, now)
	return latest, nil
}

// ListCalibrationHistory는 디바이스의 보정 이력을 조회합니다.
func (s *CalibrationService) ListCalibrationHistory(
	ctx context.Context,
	deviceID string,
	limit, offset int32,
) ([]*CalibrationRecord, int32, error) {
	if deviceID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	records, total, err := s.calRepo.ListByDevice(ctx, deviceID, limit, offset)
	if err != nil {
		s.logger.Error("보정 이력 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "보정 이력 조회에 실패했습니다")
	}

	// 현재 시각 기준 상태 업데이트
	now := time.Now().UTC()
	for _, r := range records {
		r.Status = computeStatus(r, now)
	}

	return records, total, nil
}

// CheckCalibrationStatus는 디바이스 + 카트리지 타입의 보정 상태를 확인합니다.
func (s *CalibrationService) CheckCalibrationStatus(
	ctx context.Context,
	deviceID string,
	category, typeIndex int32,
) (CalibrationStatus, string, *CalibrationRecord, error) {
	if deviceID == "" {
		return CalibrationStatusNeeded, "", nil, apperrors.New(apperrors.ErrInvalidInput, "device_id는 필수입니다")
	}

	latest, _ := s.calRepo.GetLatest(ctx, deviceID, category, typeIndex)
	if latest == nil {
		return CalibrationStatusNeeded, "보정 데이터가 없습니다. 보정이 필요합니다.", nil, nil
	}

	now := time.Now().UTC()
	st := computeStatus(latest, now)
	latest.Status = st

	var msg string
	switch st {
	case CalibrationStatusValid:
		days := int(latest.ExpiresAt.Sub(now).Hours() / 24)
		msg = "보정 유효"
		if days > 0 {
			msg = "보정 유효 (만료까지 " + intToStr(days) + "일)"
		}
	case CalibrationStatusExpiring:
		days := int(latest.ExpiresAt.Sub(now).Hours() / 24)
		msg = intToStr(days) + "일 이내 만료 예정"
	case CalibrationStatusExpired:
		msg = "보정이 만료되었습니다. 재보정이 필요합니다."
	default:
		msg = "보정 상태를 확인할 수 없습니다"
	}

	return st, msg, latest, nil
}

// ListCalibrationModels는 모든 보정 모델 목록을 반환합니다.
func (s *CalibrationService) ListCalibrationModels(ctx context.Context) ([]*CalibrationModel, error) {
	models, err := s.modelRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("보정 모델 목록 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "보정 모델 목록 조회에 실패했습니다")
	}
	return models, nil
}

// ============================================================================
// 내부 헬퍼 함수
// ============================================================================

// computeStatus는 보정 기록과 현재 시각으로 상태를 계산합니다.
func computeStatus(r *CalibrationRecord, now time.Time) CalibrationStatus {
	if now.After(r.ExpiresAt) {
		return CalibrationStatusExpired
	}
	if r.ExpiresAt.Sub(now).Hours() < float64(ExpiringThresholdDays*24) {
		return CalibrationStatusExpiring
	}
	return CalibrationStatusValid
}

// calculateAlphaFromCalibration은 기준값과 측정값에서 alpha를 계산합니다.
// Alpha = mean(measured_i / reference_i) — 최소제곱 비율 추정
func calculateAlphaFromCalibration(referenceValues, measuredValues []float64) float64 {
	if len(referenceValues) == 0 {
		return 1.0
	}

	sum := 0.0
	count := 0
	for i := range referenceValues {
		if referenceValues[i] != 0 {
			sum += measuredValues[i] / referenceValues[i]
			count++
		}
	}

	if count == 0 {
		return 1.0
	}
	return sum / float64(count)
}

// calculateAccuracyScore는 팩토리 보정의 정확도 점수를 계산합니다.
// 기본 모델 alpha와의 차이를 기반으로 0~1 점수 반환
func calculateAccuracyScore(alpha float64, model *CalibrationModel) float64 {
	if model == nil {
		// 모델 없으면 alpha가 1.0에 가까울수록 높은 점수
		diff := math.Abs(alpha - 1.0)
		return math.Max(0, 1.0-diff*5.0)
	}

	// 기본 alpha와의 차이가 적을수록 높은 점수
	diff := math.Abs(alpha - model.DefaultAlpha)
	score := 1.0 - diff*10.0
	if score < 0 {
		score = 0
	}
	if score > 1.0 {
		score = 1.0
	}
	return score
}

// calculateFieldAccuracyScore는 현장 보정의 정확도를 잔차 기반으로 계산합니다.
func calculateFieldAccuracyScore(referenceValues, measuredValues []float64, alpha float64) float64 {
	if len(referenceValues) == 0 {
		return 0
	}

	sumSquaredError := 0.0
	sumSquaredRef := 0.0
	for i := range referenceValues {
		corrected := measuredValues[i] / alpha
		err := corrected - referenceValues[i]
		sumSquaredError += err * err
		sumSquaredRef += referenceValues[i] * referenceValues[i]
	}

	if sumSquaredRef == 0 {
		return 1.0
	}

	// R² 유사 점수 (1 - SSE/SST)
	score := 1.0 - sumSquaredError/sumSquaredRef
	if score < 0 {
		score = 0
	}
	return score
}

// copyFloat64Slice는 float64 슬라이스의 깊은 복사본을 반환합니다.
func copyFloat64Slice(src []float64) []float64 {
	if src == nil {
		return nil
	}
	dst := make([]float64, len(src))
	copy(dst, src)
	return dst
}

// intToStr는 정수를 문자열로 변환합니다 (fmt 의존 제거용).
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToStr(-n)
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

// ============================================================================
// 시드 데이터: 기본 보정 모델 (CalibrationModel seed data)
// ============================================================================

// DefaultCalibrationModels는 시스템 초기화 시 사용되는 기본 보정 모델 목록입니다.
func DefaultCalibrationModels() []*CalibrationModel {
	now := time.Now().UTC()
	models := []*CalibrationModel{
		// HealthBiomarker 타입들 (카테고리 1): alpha=0.95, 유효기간 90일
		{
			ID: uuid.New().String(), CartridgeCategory: 1, CartridgeTypeIndex: 1,
			Name: "Glucose Calibration", Version: "1.0.0",
			DefaultAlpha: 0.95, ValidityDays: 90,
			Description: "혈당 측정용 보정 모델", CreatedAt: now,
		},
		{
			ID: uuid.New().String(), CartridgeCategory: 1, CartridgeTypeIndex: 2,
			Name: "LipidPanel Calibration", Version: "1.0.0",
			DefaultAlpha: 0.95, ValidityDays: 90,
			Description: "지질 패널 측정용 보정 모델", CreatedAt: now,
		},
		{
			ID: uuid.New().String(), CartridgeCategory: 1, CartridgeTypeIndex: 3,
			Name: "HbA1c Calibration", Version: "1.0.0",
			DefaultAlpha: 0.95, ValidityDays: 90,
			Description: "당화혈색소 측정용 보정 모델", CreatedAt: now,
		},

		// ElectronicSensor 타입들 (카테고리 4): alpha=0.92, 유효기간 60일
		{
			ID: uuid.New().String(), CartridgeCategory: 4, CartridgeTypeIndex: 1,
			Name: "Temperature Sensor Calibration", Version: "1.0.0",
			DefaultAlpha: 0.92, ValidityDays: 60,
			Description: "온도 센서 보정 모델", CreatedAt: now,
		},
		{
			ID: uuid.New().String(), CartridgeCategory: 4, CartridgeTypeIndex: 2,
			Name: "Humidity Sensor Calibration", Version: "1.0.0",
			DefaultAlpha: 0.92, ValidityDays: 60,
			Description: "습도 센서 보정 모델", CreatedAt: now,
		},

		// AdvancedAnalysis 타입들 (카테고리 5): alpha=0.97, 유효기간 120일
		{
			ID: uuid.New().String(), CartridgeCategory: 5, CartridgeTypeIndex: 1,
			Name: "Advanced Biomarker Calibration", Version: "1.0.0",
			DefaultAlpha: 0.97, ValidityDays: 120,
			Description: "고급 바이오마커 분석용 보정 모델", CreatedAt: now,
		},
		{
			ID: uuid.New().String(), CartridgeCategory: 5, CartridgeTypeIndex: 2,
			Name: "Advanced Genomic Calibration", Version: "1.0.0",
			DefaultAlpha: 0.97, ValidityDays: 120,
			Description: "고급 유전체 분석용 보정 모델", CreatedAt: now,
		},

		// CustomResearch (카테고리 254): alpha=0.95, 유효기간 30일
		{
			ID: uuid.New().String(), CartridgeCategory: 254, CartridgeTypeIndex: 1,
			Name: "Custom Research Calibration", Version: "1.0.0",
			DefaultAlpha: 0.95, ValidityDays: 30,
			Description: "커스텀 리서치용 보정 모델", CreatedAt: now,
		},
	}

	// ID 안정성을 위해 정렬
	sort.Slice(models, func(i, j int) bool {
		if models[i].CartridgeCategory != models[j].CartridgeCategory {
			return models[i].CartridgeCategory < models[j].CartridgeCategory
		}
		return models[i].CartridgeTypeIndex < models[j].CartridgeTypeIndex
	})

	return models
}
