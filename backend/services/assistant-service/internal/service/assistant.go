package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/services/assistant-service/internal/ai"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// AssistantSession은 AI 비서 세션입니다.
type AssistantSession struct {
	ID        string
	UserID    string
	Title     string
	Status    string // active, closed
	TurnCount int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AssistantTurn은 세션 내 대화 턴입니다.
type AssistantTurn struct {
	ID           string
	SessionID    string
	Role         string // user, assistant
	Content      string
	Intent       string
	ActionType   string
	ActionResult string
	CreatedAt    time.Time
}

// AssistantRepository는 비서 저장소 인터페이스입니다.
type AssistantRepository interface {
	CreateSession(ctx context.Context, s *AssistantSession) error
	GetSession(ctx context.Context, id string) (*AssistantSession, error)
	ListSessions(ctx context.Context, userID string, limit, offset int32) ([]*AssistantSession, int32, error)
	DeleteSession(ctx context.Context, id string) error
	AddTurn(ctx context.Context, t *AssistantTurn) error
	ListTurns(ctx context.Context, sessionID string, limit, offset int32) ([]*AssistantTurn, int32, error)
	UpdateSession(ctx context.Context, s *AssistantSession) error
}

// AssistantService는 AI 비서 비즈니스 로직입니다.
type AssistantService struct {
	logger      *zap.Logger
	repo        AssistantRepository
	aiResponder ai.AIResponder // optional: nil이면 키워드 기반 응답
}

func NewAssistantService(logger *zap.Logger, repo AssistantRepository) *AssistantService {
	return &AssistantService{logger: logger, repo: repo}
}

// SetAIResponder는 AI 응답 생성기를 설정합니다 (optional).
func (s *AssistantService) SetAIResponder(r ai.AIResponder) {
	s.aiResponder = r
}

// MaxTurnsPerSession은 세션당 최대 턴 수입니다.
const MaxTurnsPerSession = 200

// intent 분류 (간단한 키워드 기반)
var intentMap = map[string]string{
	"혈당":  "measurement.blood_glucose",
	"혈압":  "measurement.blood_pressure",
	"콜레스테롤": "measurement.cholesterol",
	"체중":  "measurement.weight",
	"예약":  "reservation.create",
	"진료":  "telemedicine.start",
	"약":   "prescription.list",
	"처방":  "prescription.list",
	"가족":  "family.list",
	"코칭":  "coaching.generate",
	"추천":  "coaching.recommendation",
	"알림":  "notification.list",
	"설정":  "settings.view",
	"도움":  "help.general",
	// Phase B 추가 인텐트
	"식사":  "nutrition.food_log",
	"음식":  "nutrition.food_log",
	"칼로리": "nutrition.calorie",
	"운동":  "fitness.plan",
	"수면":  "sleep.analysis",
	"잠":   "sleep.analysis",
	"응급":  "emergency.report",
	"긴급":  "emergency.report",
	"쇼핑":  "shop.list",
	"구매":  "shop.purchase",
	"커뮤니티": "community.list",
}

func classifyIntent(text string) string {
	lower := strings.ToLower(text)
	// 긴 키워드 우선 매칭 (예: "예약" > "약")
	bestKeyword := ""
	bestIntent := ""
	for keyword, intent := range intentMap {
		if strings.Contains(lower, keyword) {
			if len(keyword) > len(bestKeyword) {
				bestKeyword = keyword
				bestIntent = intent
			}
		}
	}
	if bestIntent != "" {
		return bestIntent
	}
	return "general.chat"
}

func generateResponse(intent, userMessage string) (string, string, string) {
	actionType := ""
	actionResult := ""

	switch {
	case strings.HasPrefix(intent, "measurement"):
		actionType = "query"
		actionResult = fmt.Sprintf(`{"domain":"measurement","intent":"%s"}`, intent)
		parts := strings.Split(intent, ".")
		metric := "건강 지표"
		if len(parts) > 1 {
			switch parts[1] {
			case "blood_glucose":
				metric = "혈당"
			case "blood_pressure":
				metric = "혈압"
			case "cholesterol":
				metric = "콜레스테롤"
			case "weight":
				metric = "체중"
			}
		}
		return fmt.Sprintf("%s 관련 정보를 확인해 드리겠습니다. 최근 측정 기록을 조회합니다.", metric), actionType, actionResult

	case strings.HasPrefix(intent, "reservation"):
		actionType = "action"
		actionResult = `{"domain":"reservation","action":"create"}`
		return "예약을 도와드리겠습니다. 원하시는 날짜와 시간을 알려주세요.", actionType, actionResult

	case strings.HasPrefix(intent, "telemedicine"):
		actionType = "action"
		actionResult = `{"domain":"telemedicine","action":"start"}`
		return "원격 진료를 시작하시겠습니까? 현재 대기 중인 의사를 확인합니다.", actionType, actionResult

	case strings.HasPrefix(intent, "prescription"):
		actionType = "query"
		actionResult = `{"domain":"prescription","action":"list"}`
		return "현재 복용 중인 약 목록을 조회합니다.", actionType, actionResult

	case strings.HasPrefix(intent, "coaching"):
		actionType = "action"
		actionResult = `{"domain":"coaching","action":"generate"}`
		return "맞춤형 건강 코칭 메시지를 생성합니다. 오늘의 건강 상태를 분석하고 있습니다.", actionType, actionResult

	case strings.HasPrefix(intent, "notification"):
		actionType = "query"
		actionResult = `{"domain":"notification","action":"list"}`
		return "최근 알림 목록을 조회합니다.", actionType, actionResult

	case strings.HasPrefix(intent, "nutrition"):
		actionType = "action"
		actionResult = `{"domain":"nutrition","action":"log"}`
		return "식사 기록을 도와드리겠습니다. 드신 음식을 말씀해 주세요.", actionType, actionResult

	case strings.HasPrefix(intent, "fitness"):
		actionType = "action"
		actionResult = `{"domain":"fitness","action":"plan"}`
		return "오늘 운동 계획을 세워드리겠습니다. 선호하시는 운동이 있으신가요?", actionType, actionResult

	case strings.HasPrefix(intent, "sleep"):
		actionType = "query"
		actionResult = `{"domain":"sleep","action":"analysis"}`
		return "수면 데이터를 분석합니다. 최근 7일간 수면 패턴을 확인합니다.", actionType, actionResult

	case strings.HasPrefix(intent, "emergency"):
		actionType = "action"
		actionResult = `{"domain":"emergency","action":"report"}`
		return "긴급 상황이신가요? 119 연결 또는 긴급 연락처에 알림을 보내드릴 수 있습니다.", actionType, actionResult

	case strings.HasPrefix(intent, "shop"):
		actionType = "query"
		actionResult = `{"domain":"shop","action":"list"}`
		return "건강 관련 상품을 검색합니다. 어떤 제품을 찾으시나요?", actionType, actionResult

	case strings.HasPrefix(intent, "community"):
		actionType = "query"
		actionResult = `{"domain":"community","action":"list"}`
		return "커뮤니티 최신 게시물을 조회합니다.", actionType, actionResult

	case intent == "help.general":
		return "안녕하세요! 만파식 AI 비서입니다. 건강 측정, 예약, 코칭, 약 관리 등을 도와드릴 수 있습니다. 무엇을 도와드릴까요?", "", ""

	default:
		return "네, 알겠습니다. 더 자세히 말씀해 주시면 도움을 드리겠습니다.", "", ""
	}
}

// SendCommand는 사용자 명령을 처리합니다.
func (s *AssistantService) SendCommand(ctx context.Context, userID, sessionID, text string) (*AssistantTurn, *AssistantTurn, *AssistantSession, error) {
	if userID == "" {
		return nil, nil, nil, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if text == "" {
		return nil, nil, nil, apperrors.New(apperrors.ErrInvalidInput, "text는 필수입니다")
	}

	now := time.Now().UTC()

	// 세션이 없으면 새로 생성
	var session *AssistantSession
	if sessionID == "" {
		session = &AssistantSession{
			ID:        uuid.New().String(),
			UserID:    userID,
			Title:     truncate(text, 50),
			Status:    "active",
			TurnCount: 0,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.repo.CreateSession(ctx, session); err != nil {
			return nil, nil, nil, apperrors.New(apperrors.ErrInternal, "세션 생성 실패")
		}
	} else {
		var err error
		session, err = s.repo.GetSession(ctx, sessionID)
		if err != nil || session == nil {
			return nil, nil, nil, apperrors.New(apperrors.ErrNotFound, "세션을 찾을 수 없습니다")
		}
		if session.Status == "closed" {
			return nil, nil, nil, apperrors.New(apperrors.ErrInvalidInput, "닫힌 세션에는 명령을 보낼 수 없습니다")
		}
		if session.TurnCount >= MaxTurnsPerSession {
			return nil, nil, nil, apperrors.New(apperrors.ErrInvalidInput, "세션 턴 제한을 초과했습니다")
		}
	}

	// 사용자 턴 저장
	intent := classifyIntent(text)
	userTurn := &AssistantTurn{
		ID:        uuid.New().String(),
		SessionID: session.ID,
		Role:      "user",
		Content:   text,
		Intent:    intent,
		CreatedAt: now,
	}
	if err := s.repo.AddTurn(ctx, userTurn); err != nil {
		return nil, nil, nil, apperrors.New(apperrors.ErrInternal, "턴 저장 실패")
	}

	// 응답 생성: AI 응답기 우선, 실패 시 키워드 기반 폴백
	responseText, actionType, actionResult := s.generateResponseWithAI(ctx, session.ID, intent, text)
	assistantTurn := &AssistantTurn{
		ID:           uuid.New().String(),
		SessionID:    session.ID,
		Role:         "assistant",
		Content:      responseText,
		Intent:       intent,
		ActionType:   actionType,
		ActionResult: actionResult,
		CreatedAt:    now.Add(time.Millisecond),
	}
	if err := s.repo.AddTurn(ctx, assistantTurn); err != nil {
		return nil, nil, nil, apperrors.New(apperrors.ErrInternal, "응답 저장 실패")
	}

	// 세션 업데이트
	session.TurnCount += 2
	session.UpdatedAt = now
	s.repo.UpdateSession(ctx, session)

	s.logger.Info("비서 명령 처리 완료",
		zap.String("user_id", userID),
		zap.String("intent", intent),
	)

	return userTurn, assistantTurn, session, nil
}

// GetSession은 세션을 조회합니다.
func (s *AssistantService) GetSession(ctx context.Context, sessionID string) (*AssistantSession, error) {
	if sessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "세션 조회 실패")
	}
	if sess == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "세션을 찾을 수 없습니다")
	}
	return sess, nil
}

// ListSessions는 세션 목록을 반환합니다.
func (s *AssistantService) ListSessions(ctx context.Context, userID string, limit, offset int32) ([]*AssistantSession, int32, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 20
	}
	return s.repo.ListSessions(ctx, userID, limit, offset)
}

// ListTurns는 세션의 턴 목록을 반환합니다.
func (s *AssistantService) ListTurns(ctx context.Context, sessionID string, limit, offset int32) ([]*AssistantTurn, int32, error) {
	if sessionID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}
	if limit <= 0 {
		limit = 50
	}
	return s.repo.ListTurns(ctx, sessionID, limit, offset)
}

// CloseSession은 세션을 종료합니다.
func (s *AssistantService) CloseSession(ctx context.Context, sessionID string) (*AssistantSession, error) {
	if sessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil || sess == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "세션을 찾을 수 없습니다")
	}
	if sess.Status == "closed" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "이미 닫힌 세션입니다")
	}
	sess.Status = "closed"
	sess.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateSession(ctx, sess); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "세션 업데이트 실패")
	}
	return sess, nil
}

// GetSessionSummary는 세션의 대화 요약을 생성합니다.
func (s *AssistantService) GetSessionSummary(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil || sess == nil {
		return "", apperrors.New(apperrors.ErrNotFound, "세션을 찾을 수 없습니다")
	}

	turns, _, err := s.repo.ListTurns(ctx, sessionID, 200, 0)
	if err != nil {
		return "", apperrors.New(apperrors.ErrInternal, "턴 조회 실패")
	}

	// 인텐트별 통계 수집
	intentCounts := make(map[string]int)
	userMsgCount := 0
	for _, t := range turns {
		if t.Role == "user" {
			userMsgCount++
			if t.Intent != "" {
				intentCounts[t.Intent]++
			}
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("세션 '%s' 요약:\n", sess.Title))
	sb.WriteString(fmt.Sprintf("- 총 %d개의 사용자 메시지\n", userMsgCount))
	if len(intentCounts) > 0 {
		sb.WriteString("- 주요 인텐트: ")
		first := true
		for intent, count := range intentCounts {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s(%d회)", intent, count))
			first = false
		}
		sb.WriteString("\n")
	}
	sb.WriteString(fmt.Sprintf("- 세션 상태: %s", sess.Status))

	return sb.String(), nil
}

// DeleteSession은 세션을 삭제합니다.
func (s *AssistantService) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return apperrors.New(apperrors.ErrInvalidInput, "session_id는 필수입니다")
	}
	return s.repo.DeleteSession(ctx, sessionID)
}

// generateResponseWithAI는 AI 응답기가 있으면 사용하고, 없거나 실패하면 키워드 기반 폴백합니다.
func (s *AssistantService) generateResponseWithAI(ctx context.Context, sessionID, intent, userMessage string) (string, string, string) {
	if s.aiResponder == nil {
		return generateResponse(intent, userMessage)
	}

	// 대화 이력 조회 (최근 10턴)
	var history []ai.ChatMessage
	turns, _, err := s.repo.ListTurns(ctx, sessionID, 10, 0)
	if err == nil {
		for _, t := range turns {
			history = append(history, ai.ChatMessage{Role: t.Role, Content: t.Content})
		}
	}

	aiResp, err := s.aiResponder.GenerateResponse(ctx, history, userMessage)
	if err != nil {
		s.logger.Warn("AI 응답 생성 실패, 키워드 폴백 사용",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
		return generateResponse(intent, userMessage)
	}
	if aiResp == nil {
		return generateResponse(intent, userMessage)
	}

	actionType := aiResp.ActionType
	actionResult := aiResp.ActionResult
	// AI가 액션을 생성하지 않은 경우, 인텐트 기반 액션 사용
	if actionType == "" {
		_, actionType, actionResult = generateResponse(intent, userMessage)
	}

	return aiResp.Content, actionType, actionResult
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
