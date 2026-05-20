// Package service는 digital-twin-service의 비즈니스 로직을 구현합니다.
//
// Rust TwinEngine의 Go 대응 서비스로, 디지털 트윈 상태 관리와
// EWMA/CUSUM 기반 드리프트 감지를 수행합니다.
package service

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TwinState는 디지털 트윈 인스턴스 상태입니다.
type TwinState struct {
	ID                    string
	SessionID             string
	UserID                string
	DeviceID              string
	EWMAValue             float64
	CUSUMPos              float64
	CUSUMNeg              float64
	HealthState           string // "nominal", "warning", "drift", "recalibrate"
	DriftScore            float64
	RemainingMeasurements int
	MeasurementCount      int
	LastSyncedAt          time.Time
	CreatedAt             time.Time
}

// TwinRepository는 디지털 트윈 상태 저장소 인터페이스입니다.
type TwinRepository interface {
	Store(ctx context.Context, state *TwinState) error
	Get(ctx context.Context, sessionID string) (*TwinState, error)
	ListByDevice(ctx context.Context, deviceID string, limit int) ([]*TwinState, error)
}

// TwinService는 디지털 트윈 서비스입니다.
type TwinService struct {
	logger *zap.Logger
	repo   TwinRepository
	// EWMA/CUSUM 파라미터 (Rust TwinEngine과 동일)
	ewmaLambda float64 // EWMA 감쇠 계수 (기본 0.2)
	cusumK     float64 // CUSUM 허용 슬랙 (기본 0.5)
	cusumH     float64 // CUSUM 결정 임계값 (기본 4.0)
}

// NewTwinService는 새 TwinService를 생성합니다.
func NewTwinService(logger *zap.Logger, repo TwinRepository) *TwinService {
	return &TwinService{
		logger:     logger,
		repo:       repo,
		ewmaLambda: 0.2,
		cusumK:     0.5,
		cusumH:     4.0,
	}
}

// SyncTwin은 디지털 트윈 상태를 동기화합니다.
func (s *TwinService) SyncTwin(ctx context.Context, sessionID, userID, deviceID string,
	residuals []float64, ewmaValue, cusumPos, cusumNeg float64,
	healthState string, driftScore float64, remainingMeasurements int,
) (*TwinState, string, error) {
	// 기존 트윈 상태 조회
	existing, _ := s.repo.Get(ctx, sessionID)

	now := time.Now().UTC()
	state := &TwinState{
		SessionID:             sessionID,
		UserID:                userID,
		DeviceID:              deviceID,
		EWMAValue:             ewmaValue,
		CUSUMPos:              cusumPos,
		CUSUMNeg:              cusumNeg,
		HealthState:           healthState,
		DriftScore:            driftScore,
		RemainingMeasurements: remainingMeasurements,
		LastSyncedAt:          now,
	}

	if existing != nil {
		state.ID = existing.ID
		state.MeasurementCount = existing.MeasurementCount + 1
		state.CreatedAt = existing.CreatedAt
	} else {
		state.ID = uuid.New().String()
		state.MeasurementCount = 1
		state.CreatedAt = now
	}

	// 서버 측 EWMA/CUSUM 검증 (Rust 엔진 결과와 교차 검증)
	serverHealthState := s.evaluateHealthState(ewmaValue, cusumPos, cusumNeg, driftScore)
	if serverHealthState != healthState {
		s.logger.Warn("클라이언트/서버 건강 상태 불일치",
			zap.String("client", healthState),
			zap.String("server", serverHealthState),
			zap.String("session_id", sessionID),
		)
		// 서버 판정이 더 보수적이면 서버 결과 사용
		if s.severityOrder(serverHealthState) > s.severityOrder(healthState) {
			state.HealthState = serverHealthState
		}
	}

	// 권장 조치 결정
	recommendedAction := s.determineAction(state)

	if err := s.repo.Store(ctx, state); err != nil {
		s.logger.Error("디지털 트윈 상태 저장 실패", zap.Error(err))
		return nil, recommendedAction, err
	}

	s.logger.Info("디지털 트윈 동기화",
		zap.String("twin_id", state.ID),
		zap.String("session_id", sessionID),
		zap.String("health_state", state.HealthState),
		zap.Float64("drift_score", driftScore),
		zap.String("action", recommendedAction),
	)

	return state, recommendedAction, nil
}

// GetTwinState는 세션의 트윈 상태를 조회합니다.
func (s *TwinService) GetTwinState(ctx context.Context, sessionID string) (*TwinState, error) {
	return s.repo.Get(ctx, sessionID)
}

// ListDeviceTwins는 디바이스의 트윈 이력을 조회합니다.
func (s *TwinService) ListDeviceTwins(ctx context.Context, deviceID string, limit int) ([]*TwinState, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListByDevice(ctx, deviceID, limit)
}

// evaluateHealthState는 EWMA/CUSUM 값으로 건강 상태를 판정합니다.
func (s *TwinService) evaluateHealthState(ewma, cusumPos, cusumNeg, driftScore float64) string {
	// CUSUM 초과 → 즉시 재교정
	if cusumPos > s.cusumH || math.Abs(cusumNeg) > s.cusumH {
		return "recalibrate"
	}
	// 드리프트 점수 기반 판정
	if driftScore > 0.8 {
		return "recalibrate"
	}
	if driftScore > 0.5 {
		return "drift"
	}
	if driftScore > 0.3 || math.Abs(ewma) > 2.0 {
		return "warning"
	}
	return "nominal"
}

// determineAction은 상태 기반 권장 조치를 결정합니다.
func (s *TwinService) determineAction(state *TwinState) string {
	switch state.HealthState {
	case "recalibrate":
		if state.RemainingMeasurements < 10 {
			return "replace_cartridge"
		}
		return "recalibrate"
	case "drift":
		return "recalibrate"
	case "warning":
		return "continue"
	default:
		return "continue"
	}
}

// severityOrder는 건강 상태의 심각도 순서를 반환합니다.
func (s *TwinService) severityOrder(state string) int {
	switch state {
	case "nominal":
		return 0
	case "warning":
		return 1
	case "drift":
		return 2
	case "recalibrate":
		return 3
	default:
		return 0
	}
}
