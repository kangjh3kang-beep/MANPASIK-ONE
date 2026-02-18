// Package service는 notification-service의 에스컬레이션 로직을 구현합니다.
package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// EscalationStage는 에스컬레이션 단계를 나타냅니다.
type EscalationStage int

const (
	StageNone EscalationStage = iota
	// Stage1: AI 분석 + 인앱 알림 (즉시)
	StageInAppAlert
	// Stage2: 보호자 푸시 알림 (3분 후 미확인 시)
	StageGuardianPush
	// Stage3: AI 음성 확인 전화 (6분 후 미확인 시)
	StageAIVoiceCall
	// Stage4: 119 자동 신고 (10분 후 미확인 시)
	StageEmergencyCall
	// 종료 상태
	StageResolved
	StageCancelled
)

// EscalationEvent는 에스컬레이션 이벤트입니다.
type EscalationEvent struct {
	ID             string
	UserID         string
	AlertType      string // "health_critical", "fall_detected", "no_response"
	Stage          EscalationStage
	MeasurementID  string
	Value          string
	CreatedAt      time.Time
	LastEscalation time.Time
	ResolvedAt     *time.Time
	ResolvedBy     string // "user", "guardian", "operator", ""
}

// EscalationConfig는 에스컬레이션 타이밍 설정입니다.
type EscalationConfig struct {
	Stage2Delay time.Duration // Stage1→Stage2 대기 시간 (기본 3분)
	Stage3Delay time.Duration // Stage2→Stage3 대기 시간 (기본 3분)
	Stage4Delay time.Duration // Stage3→Stage4 대기 시간 (기본 4분)
}

// DefaultEscalationConfig는 기본 에스컬레이션 설정입니다.
var DefaultEscalationConfig = EscalationConfig{
	Stage2Delay: 3 * time.Minute,
	Stage3Delay: 3 * time.Minute,
	Stage4Delay: 4 * time.Minute,
}

// EscalationService는 에스컬레이션 상태 머신을 관리합니다.
type EscalationService struct {
	log       *zap.Logger
	notiSvc   *NotificationService
	config    EscalationConfig
	active    map[string]*EscalationEvent // eventID → event
	guardians map[string][]string         // userID → guardian userIDs
}

// NewEscalationService는 EscalationService를 생성합니다.
func NewEscalationService(log *zap.Logger, notiSvc *NotificationService) *EscalationService {
	return &EscalationService{
		log:       log,
		notiSvc:   notiSvc,
		config:    DefaultEscalationConfig,
		active:    make(map[string]*EscalationEvent),
		guardians: make(map[string][]string),
	}
}

// SetConfig는 에스컬레이션 타이밍 설정을 변경합니다.
func (es *EscalationService) SetConfig(cfg EscalationConfig) {
	es.config = cfg
}

// TriggerEscalation은 건강 이상 감지 시 에스컬레이션 체인을 시작합니다.
func (es *EscalationService) TriggerEscalation(ctx context.Context, userID, alertType, measurementID, value string) (*EscalationEvent, error) {
	if userID == "" || alertType == "" {
		return nil, fmt.Errorf("userID와 alertType은 필수입니다")
	}

	event := &EscalationEvent{
		ID:             fmt.Sprintf("esc-%d", time.Now().UnixNano()),
		UserID:         userID,
		AlertType:      alertType,
		Stage:          StageInAppAlert,
		MeasurementID:  measurementID,
		Value:          value,
		CreatedAt:      time.Now(),
		LastEscalation: time.Now(),
	}

	es.active[event.ID] = event

	// Stage 1: 인앱 알림 즉시 발송
	es.executeStage1(ctx, event)

	es.log.Info("에스컬레이션 시작",
		zap.String("event_id", event.ID),
		zap.String("user_id", userID),
		zap.String("alert_type", alertType),
		zap.String("value", value),
	)

	// 비동기 에스컬레이션 체인 시작
	go es.runEscalationChain(event)

	return event, nil
}

// AcknowledgeEscalation은 사용자/보호자가 에스컬레이션을 확인하여 중단합니다.
func (es *EscalationService) AcknowledgeEscalation(ctx context.Context, eventID, resolvedBy string) error {
	event, exists := es.active[eventID]
	if !exists {
		return fmt.Errorf("에스컬레이션 이벤트를 찾을 수 없습니다: %s", eventID)
	}

	now := time.Now()
	event.ResolvedAt = &now
	event.ResolvedBy = resolvedBy
	event.Stage = StageResolved

	delete(es.active, eventID)

	es.log.Info("에스컬레이션 해제",
		zap.String("event_id", eventID),
		zap.String("resolved_by", resolvedBy),
	)

	return nil
}

// GetActiveEscalations은 사용자의 활성 에스컬레이션 목록을 반환합니다.
func (es *EscalationService) GetActiveEscalations(userID string) []*EscalationEvent {
	var result []*EscalationEvent
	for _, event := range es.active {
		if event.UserID == userID {
			result = append(result, event)
		}
	}
	return result
}

// SetGuardians는 사용자의 보호자 목록을 설정합니다.
func (es *EscalationService) SetGuardians(userID string, guardianIDs []string) {
	es.guardians[userID] = guardianIDs
}

// --- 내부 메서드 ---

// executeStage1: AI 분석 + 인앱 알림
func (es *EscalationService) executeStage1(ctx context.Context, event *EscalationEvent) {
	title := "건강 이상 감지"
	body := fmt.Sprintf("[긴급] %s 수치가 위험 범위입니다: %s. 즉시 확인해주세요.", event.AlertType, event.Value)
	_, _ = es.notiSvc.SendNotification(ctx, event.UserID, TypeHealthAlert, ChannelInApp, PriorityUrgent, title, body, event.MeasurementID)
}

// executeStage2: 보호자 푸시 알림
func (es *EscalationService) executeStage2(ctx context.Context, event *EscalationEvent) {
	guardianIDs := es.guardians[event.UserID]
	for _, gid := range guardianIDs {
		title := "가족 건강 긴급 알림"
		body := fmt.Sprintf("가족 구성원의 %s 수치가 위험합니다: %s. 확인이 필요합니다.", event.AlertType, event.Value)
		_, _ = es.notiSvc.SendNotification(ctx, gid, TypeHealthAlert, ChannelPush, PriorityUrgent, title, body, event.MeasurementID)
	}
	event.Stage = StageGuardianPush
	event.LastEscalation = time.Now()
}

// executeStage3: AI 음성 확인 전화
func (es *EscalationService) executeStage3(_ context.Context, event *EscalationEvent) {
	// 실제 구현: TTS/STT 기반 AI 음성 통화 서비스 연동
	// 사용자에게 전화하여 "괜찮으세요? 1번을 눌러 확인해주세요" 시나리오 실행
	es.log.Info("Stage3: AI 음성 확인 전화 발신",
		zap.String("event_id", event.ID),
		zap.String("user_id", event.UserID),
	)
	event.Stage = StageAIVoiceCall
	event.LastEscalation = time.Now()
}

// executeStage4: 119 자동 신고
func (es *EscalationService) executeStage4(_ context.Context, event *EscalationEvent) {
	// 실제 구현: 119 API 또는 긴급 연락 시스템 연동
	// 사용자 위치 정보 + 증상 정보를 포함하여 신고
	es.log.Warn("Stage4: 119 자동 신고 발동",
		zap.String("event_id", event.ID),
		zap.String("user_id", event.UserID),
		zap.String("alert_type", event.AlertType),
		zap.String("value", event.Value),
	)
	event.Stage = StageEmergencyCall
	event.LastEscalation = time.Now()
}

// runEscalationChain은 비동기로 에스컬레이션 체인을 실행합니다.
func (es *EscalationService) runEscalationChain(event *EscalationEvent) {
	ctx := context.Background()

	// Stage1 → Stage2 대기
	time.Sleep(es.config.Stage2Delay)
	if _, ok := es.active[event.ID]; !ok {
		return // 이미 해제됨
	}
	es.executeStage2(ctx, event)

	// Stage2 → Stage3 대기
	time.Sleep(es.config.Stage3Delay)
	if _, ok := es.active[event.ID]; !ok {
		return
	}
	es.executeStage3(ctx, event)

	// Stage3 → Stage4 대기
	time.Sleep(es.config.Stage4Delay)
	if _, ok := es.active[event.ID]; !ok {
		return
	}
	es.executeStage4(ctx, event)
}
