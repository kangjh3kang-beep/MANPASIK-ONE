// Package service는 measurement-service의 비즈니스 로직을 구현합니다.
//
// 핵심 비즈니스: 측정 세션 관리, 차동측정 데이터 처리, 핑거프린트 저장
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// MeasurementService는 측정 서비스의 비즈니스 로직입니다.
type MeasurementService struct {
	logger         *zap.Logger
	sessionRepo    SessionRepository
	measureRepo    MeasurementRepository
	vectorRepo     VectorRepository
	eventPublisher EventPublisher
	searchIndexer  SearchIndexer // optional: nil이면 인덱싱 비활성화
}

// SessionRepository는 측정 세션 저장소 인터페이스입니다 (PostgreSQL).
type SessionRepository interface {
	CreateSession(ctx context.Context, session *MeasurementSession) error
	GetSession(ctx context.Context, sessionID string) (*MeasurementSession, error)
	EndSession(ctx context.Context, sessionID string, totalMeasurements int, endedAt time.Time) error
}

// MeasurementRepository는 측정 데이터 저장소 인터페이스입니다 (TimescaleDB).
type MeasurementRepository interface {
	Store(ctx context.Context, data *MeasurementData) error
	GetHistory(ctx context.Context, userID string, start, end time.Time, limit, offset int) ([]*MeasurementSummary, int, error)
}

// VectorRepository는 핑거프린트 벡터 저장소 인터페이스입니다 (Milvus).
type VectorRepository interface {
	StoreFingerprint(ctx context.Context, sessionID string, vector []float32) error
	SearchSimilar(ctx context.Context, vector []float32, topK int) ([]SimilarResult, error)
}

// EventPublisher는 이벤트 발행 인터페이스입니다 (Kafka).
type EventPublisher interface {
	PublishMeasurementCompleted(ctx context.Context, event *MeasurementCompletedEvent) error
}

// SearchIndexer는 측정 데이터 검색 인덱싱 인터페이스입니다 (Elasticsearch).
type SearchIndexer interface {
	IndexMeasurement(ctx context.Context, sessionID string, data *MeasurementData) error
}

// MeasurementSession은 측정 세션 엔티티입니다.
type MeasurementSession struct {
	ID                string
	DeviceID          string
	CartridgeID       string
	UserID            string
	StartedAt         time.Time
	EndedAt           *time.Time
	TotalMeasurements int
	Status            string // "active", "completed", "error"
}

// MeasurementData는 개별 측정 데이터입니다.
type MeasurementData struct {
	Time              time.Time
	SessionID         string
	DeviceID          string
	UserID            string
	CartridgeType     string
	RawChannels       []float64
	SDet              float64
	SRef              float64
	Alpha             float64
	SCorrected        float64
	PrimaryValue      float64
	Unit              string
	Confidence        float64
	FingerprintVector []float32
	TempC             float32
	HumidityPct       float32
	BatteryPct        int
}

// MeasurementSummary는 측정 결과 요약입니다.
type MeasurementSummary struct {
	SessionID     string
	CartridgeType string
	PrimaryValue  float64
	Unit          string
	MeasuredAt    time.Time
}

// SimilarResult는 유사 벡터 검색 결과입니다.
type SimilarResult struct {
	SessionID string
	Score     float32
	Distance  float32
}

// MeasurementCompletedEvent는 측정 완료 이벤트입니다.
type MeasurementCompletedEvent struct {
	SessionID    string
	UserID       string
	DeviceID     string
	PrimaryValue float64
	Unit         string
	CompletedAt  time.Time
}

// NewMeasurementService는 새 MeasurementService를 생성합니다.
func NewMeasurementService(
	logger *zap.Logger,
	sessionRepo SessionRepository,
	measureRepo MeasurementRepository,
	vectorRepo VectorRepository,
	eventPublisher EventPublisher,
) *MeasurementService {
	return &MeasurementService{
		logger:         logger,
		sessionRepo:    sessionRepo,
		measureRepo:    measureRepo,
		vectorRepo:     vectorRepo,
		eventPublisher: eventPublisher,
	}
}

// SetSearchIndexer는 검색 인덱서를 설정합니다 (optional).
func (s *MeasurementService) SetSearchIndexer(indexer SearchIndexer) {
	s.searchIndexer = indexer
}

// StartSession은 새 측정 세션을 시작합니다.
func (s *MeasurementService) StartSession(
	ctx context.Context,
	deviceID, cartridgeID, userID string,
) (*MeasurementSession, error) {
	// 입력 검증
	if deviceID == "" || cartridgeID == "" || userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "device_id, cartridge_id, user_id는 필수입니다")
	}

	session := &MeasurementSession{
		ID:          uuid.New().String(),
		DeviceID:    deviceID,
		CartridgeID: cartridgeID,
		UserID:      userID,
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("세션 생성 실패",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return nil, apperrors.New(apperrors.ErrInternal, "측정 세션 시작에 실패했습니다")
	}

	s.logger.Info("측정 세션 시작",
		zap.String("session_id", session.ID),
		zap.String("device_id", deviceID),
		zap.String("cartridge_id", cartridgeID),
	)

	return session, nil
}

// ProcessMeasurement는 스트리밍 측정 데이터를 처리합니다.
func (s *MeasurementService) ProcessMeasurement(
	ctx context.Context,
	data *MeasurementData,
) (*ProcessedResult, error) {
	if data.SessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}

	// TimescaleDB에 시계열 데이터 저장
	if err := s.measureRepo.Store(ctx, data); err != nil {
		s.logger.Error("측정 데이터 저장 실패",
			zap.String("session_id", data.SessionID),
			zap.Error(err),
		)
		return nil, apperrors.New(apperrors.ErrInternal, "측정 데이터 저장에 실패했습니다")
	}

	// Milvus에 핑거프린트 벡터 저장
	if len(data.FingerprintVector) > 0 {
		if err := s.vectorRepo.StoreFingerprint(ctx, data.SessionID, data.FingerprintVector); err != nil {
			s.logger.Warn("벡터 저장 실패 (비치명적)",
				zap.String("session_id", data.SessionID),
				zap.Error(err),
			)
			// 벡터 저장 실패는 비치명적 에러 (측정 결과에 영향 없음)
		}
	}

	return &ProcessedResult{
		SessionID:    data.SessionID,
		PrimaryValue: data.PrimaryValue,
		Unit:         data.Unit,
		Confidence:   data.Confidence,
		ProcessedAt:  time.Now().UTC(),
	}, nil
}

// EndSession은 측정 세션을 종료합니다.
func (s *MeasurementService) EndSession(
	ctx context.Context,
	sessionID string,
) (*SessionEndResult, error) {
	if sessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}

	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil || session == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "측정 세션을 찾을 수 없습니다")
	}

	endedAt := time.Now().UTC()
	if err := s.sessionRepo.EndSession(ctx, sessionID, session.TotalMeasurements, endedAt); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "세션 종료에 실패했습니다")
	}

	// 측정 완료 이벤트 발행 (Kafka)
	event := &MeasurementCompletedEvent{
		SessionID:   sessionID,
		UserID:      session.UserID,
		DeviceID:    session.DeviceID,
		CompletedAt: endedAt,
	}
	if err := s.eventPublisher.PublishMeasurementCompleted(ctx, event); err != nil {
		s.logger.Warn("이벤트 발행 실패 (비치명적)",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
	}

	s.logger.Info("측정 세션 종료",
		zap.String("session_id", sessionID),
		zap.Int("total_measurements", session.TotalMeasurements),
	)

	return &SessionEndResult{
		SessionID:         sessionID,
		TotalMeasurements: session.TotalMeasurements,
		EndedAt:           endedAt,
	}, nil
}

// GetHistory는 측정 기록을 조회합니다.
func (s *MeasurementService) GetHistory(
	ctx context.Context,
	userID string,
	start, end time.Time,
	limit, offset int,
) ([]*MeasurementSummary, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.measureRepo.GetHistory(ctx, userID, start, end, limit, offset)
}

// defaultLOINCMap은 바이오마커 이름에 대한 LOINC 코드 매핑입니다.
var defaultLOINCMap = map[string]string{
	"blood_glucose":     "15074-8",
	"blood_pressure":    "85354-9",
	"cholesterol_total": "2093-3",
	"hemoglobin_a1c":    "4548-4",
	"heart_rate":        "8867-4",
	"body_temperature":  "8310-5",
	"oxygen_saturation": "2708-6",
}

// loincCodeFor는 바이오마커 이름에 대한 LOINC 코드를 반환합니다. 매핑이 없으면 기본값을 반환합니다.
func loincCodeFor(cartridgeType string) string {
	if code, ok := defaultLOINCMap[cartridgeType]; ok {
		return code
	}
	return "29463-7" // 기본 LOINC 코드
}

// ExportSingleMeasurement는 단일 세션의 측정 결과를 FHIR R4 Observation Bundle JSON으로 내보냅니다.
func (s *MeasurementService) ExportSingleMeasurement(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}

	// 1. 세션 조회
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil || session == nil {
		return "", apperrors.New(apperrors.ErrNotFound, "측정 세션을 찾을 수 없습니다")
	}

	// 2. 해당 세션의 측정 데이터 조회 (유저 기준으로 가져온 뒤 세션 필터링)
	summaries, _, err := s.measureRepo.GetHistory(ctx, session.UserID, time.Time{}, time.Now().UTC(), 500, 0)
	if err != nil {
		return "", fmt.Errorf("측정 데이터 조회 실패: %w", err)
	}

	// 3. 세션 ID로 필터링하고 FHIR Observation으로 변환
	observations := make([]map[string]interface{}, 0)
	for _, m := range summaries {
		if m.SessionID != sessionID {
			continue
		}
		obs := map[string]interface{}{
			"resourceType": "Observation",
			"status":       "final",
			"code": map[string]interface{}{
				"coding": []map[string]interface{}{{
					"system":  "http://loinc.org",
					"code":    loincCodeFor(m.CartridgeType),
					"display": m.CartridgeType,
				}},
			},
			"subject": map[string]interface{}{
				"reference": "Patient/" + session.UserID,
			},
			"effectiveDateTime": m.MeasuredAt.Format(time.RFC3339),
			"valueQuantity": map[string]interface{}{
				"value":  m.PrimaryValue,
				"unit":   m.Unit,
				"system": "http://unitsofmeasure.org",
				"code":   m.Unit,
			},
			"identifier": []map[string]interface{}{{
				"system": "https://manpasik.com/session",
				"value":  m.SessionID,
			}},
		}
		observations = append(observations, obs)
	}

	// 4. FHIR Bundle 구성
	entries := make([]map[string]interface{}, 0, len(observations))
	for _, o := range observations {
		entries = append(entries, map[string]interface{}{"resource": o})
	}
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry":        entries,
	}

	jsonBytes, err := json.Marshal(bundle)
	if err != nil {
		return "", fmt.Errorf("FHIR Bundle 직렬화 실패: %w", err)
	}
	return string(jsonBytes), nil
}

// ExportToFHIRObservations는 측정 결과를 FHIR R4 Observation Bundle로 내보냅니다 (Agent D).
// biomarkerNames가 비어 있으면 전체 바이오마커 포함.
func (s *MeasurementService) ExportToFHIRObservations(
	ctx context.Context,
	userID string,
	fromDate, toDate *time.Time,
	biomarkerNames []string,
) (fhirBundleJSON string, observationCount int, err error) {
	if userID == "" {
		return "", 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	var start, end time.Time
	if fromDate != nil {
		start = *fromDate
	}
	if toDate != nil {
		end = *toDate
	} else {
		end = time.Now().UTC()
	}
	summaries, _, err := s.measureRepo.GetHistory(ctx, userID, start, end, 500, 0)
	if err != nil {
		return "", 0, err
	}
	biomarkerSet := make(map[string]struct{})
	for _, n := range biomarkerNames {
		biomarkerSet[n] = struct{}{}
	}
	observations := make([]map[string]interface{}, 0, len(summaries))
	for _, m := range summaries {
		if len(biomarkerSet) > 0 {
			if _, ok := biomarkerSet[m.CartridgeType]; !ok {
				continue
			}
		}
		obs := map[string]interface{}{
			"resourceType": "Observation",
			"status":       "final",
			"code": map[string]interface{}{
				"coding": []map[string]interface{}{{
					"system":  "http://loinc.org",
					"code":   "29463-7",
					"display": m.CartridgeType,
				}},
			},
			"subject": map[string]interface{}{
				"reference": "Patient/" + userID,
			},
			"effectiveDateTime": m.MeasuredAt.Format(time.RFC3339),
			"valueQuantity": map[string]interface{}{
				"value": m.PrimaryValue,
				"unit":  m.Unit,
				"system": "http://unitsofmeasure.org",
				"code":  m.Unit,
			},
			"identifier": []map[string]interface{}{{
				"system": "https://manpasik.com/session",
				"value": m.SessionID,
			}},
		}
		observations = append(observations, obs)
	}
	entries := make([]map[string]interface{}, 0, len(observations))
	for _, o := range observations {
		entries = append(entries, map[string]interface{}{"resource": o})
	}
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry":        entries,
	}
	jsonBytes, err := json.Marshal(bundle)
	if err != nil {
		return "", 0, fmt.Errorf("FHIR Bundle 직렬화 실패: %w", err)
	}
	return string(jsonBytes), len(observations), nil
}

// ProcessedResult는 처리된 측정 결과입니다.
type ProcessedResult struct {
	SessionID    string
	PrimaryValue float64
	Unit         string
	Confidence   float64
	ProcessedAt  time.Time
}

// SessionEndResult는 세션 종료 결과입니다.
type SessionEndResult struct {
	SessionID         string
	TotalMeasurements int
	EndedAt           time.Time
}

// 사용하지 않는 import 방지
var _ = fmt.Sprintf
