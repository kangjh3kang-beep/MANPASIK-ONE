package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/manpasik/backend/services/assistant-service/internal/ai"
	"github.com/manpasik/backend/services/assistant-service/internal/repository/memory"
	"github.com/manpasik/backend/services/assistant-service/internal/service"
	"go.uber.org/zap"
)

func newTestService() *service.AssistantService {
	repo := memory.NewAssistantRepository()
	logger := zap.NewNop()
	return service.NewAssistantService(logger, repo)
}

func TestSendCommand(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, asstTurn, session, err := svc.SendCommand(ctx, "user1", "", "혈당 측정 결과 알려줘")
	if err != nil {
		t.Fatalf("SendCommand failed: %v", err)
	}
	if session == nil || session.ID == "" {
		t.Fatal("session should be created")
	}
	if userTurn.Role != "user" {
		t.Errorf("expected user role, got %s", userTurn.Role)
	}
	if asstTurn.Role != "assistant" {
		t.Errorf("expected assistant role, got %s", asstTurn.Role)
	}
	if userTurn.Intent != "measurement.blood_glucose" {
		t.Errorf("expected blood_glucose intent, got %s", userTurn.Intent)
	}
}

func TestSendCommandEmptyUser(t *testing.T) {
	svc := newTestService()
	_, _, _, err := svc.SendCommand(context.Background(), "", "", "hello")
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestSendCommandExistingSession(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "안녕")
	_, _, session2, err := svc.SendCommand(ctx, "user1", session.ID, "혈압 알려줘")
	if err != nil {
		t.Fatalf("SendCommand with session failed: %v", err)
	}
	if session2.ID != session.ID {
		t.Error("should reuse same session")
	}
}

func TestListSessions(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	svc.SendCommand(ctx, "user1", "", "첫번째")
	svc.SendCommand(ctx, "user1", "", "두번째")

	sessions, total, err := svc.ListSessions(ctx, "user1", 10, 0)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 sessions, got %d", total)
	}
	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestListTurns(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "혈당")
	turns, total, err := svc.ListTurns(ctx, session.ID, 50, 0)
	if err != nil {
		t.Fatalf("ListTurns failed: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2 turns, got %d", total)
	}
	if len(turns) != 2 {
		t.Errorf("expected 2 turns, got %d", len(turns))
	}
}

func TestDeleteSession(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "test")
	err := svc.DeleteSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}
	_, err = svc.GetSession(ctx, session.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestIntentClassificationViaCommand(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	// 혈당 의도
	_, asstTurn, _, err := svc.SendCommand(ctx, "user1", "", "혈당 알려줘")
	if err != nil {
		t.Fatal(err)
	}
	if asstTurn.Intent != "measurement.blood_glucose" {
		t.Errorf("expected blood_glucose intent, got %s", asstTurn.Intent)
	}

	// 예약 의도
	_, asstTurn2, _, err := svc.SendCommand(ctx, "user1", "", "예약 잡아줘")
	if err != nil {
		t.Fatal(err)
	}
	if asstTurn2.Intent != "reservation.create" {
		t.Errorf("expected reservation intent, got %s", asstTurn2.Intent)
	}
}

// ============================================================================
// Phase B 신규 테스트
// ============================================================================

func TestSendCommand_ClosedSession(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "안녕")
	svc.CloseSession(ctx, session.ID)

	_, _, _, err := svc.SendCommand(ctx, "user1", session.ID, "추가 메시지")
	if err == nil {
		t.Fatal("닫힌 세션에 명령 전송 시 에러가 반환되어야 합니다")
	}
}

func TestCloseSession_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "안녕")

	closed, err := svc.CloseSession(ctx, session.ID)
	if err != nil {
		t.Fatalf("CloseSession 실패: %v", err)
	}
	if closed.Status != "closed" {
		t.Errorf("Status = %s, want closed", closed.Status)
	}
}

func TestCloseSession_AlreadyClosed(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "안녕")
	svc.CloseSession(ctx, session.ID)

	_, err := svc.CloseSession(ctx, session.ID)
	if err == nil {
		t.Fatal("이미 닫힌 세션 재닫기에 에러가 반환되어야 합니다")
	}
}

func TestGetSessionSummary(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "혈당 측정 결과")
	svc.SendCommand(ctx, "user1", session.ID, "혈압도 알려줘")

	summary, err := svc.GetSessionSummary(ctx, session.ID)
	if err != nil {
		t.Fatalf("GetSessionSummary 실패: %v", err)
	}
	if summary == "" {
		t.Fatal("요약이 비어 있습니다")
	}
}

func TestIntentClassification_Food(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, _, _, err := svc.SendCommand(ctx, "user1", "", "오늘 식사 기록해줘")
	if err != nil {
		t.Fatal(err)
	}
	if userTurn.Intent != "nutrition.food_log" {
		t.Errorf("Intent = %s, want nutrition.food_log", userTurn.Intent)
	}
}

func TestSendCommand_EmptyText(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, _, err := svc.SendCommand(ctx, "user1", "", "")
	if err == nil {
		t.Fatal("빈 텍스트에 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase F 테스트 보강
// ============================================================================

func TestGetSession_EmptyID(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetSession(context.Background(), "")
	if err == nil {
		t.Fatal("빈 session_id에 에러가 반환되어야 합니다")
	}
}

func TestGetSession_NotFound(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetSession(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 세션에 에러가 반환되어야 합니다")
	}
}

func TestListSessions_EmptyUserID(t *testing.T) {
	svc := newTestService()
	_, _, err := svc.ListSessions(context.Background(), "", 10, 0)
	if err == nil {
		t.Fatal("빈 user_id에 에러가 반환되어야 합니다")
	}
}

func TestListTurns_EmptySessionID(t *testing.T) {
	svc := newTestService()
	_, _, err := svc.ListTurns(context.Background(), "", 50, 0)
	if err == nil {
		t.Fatal("빈 session_id에 에러가 반환되어야 합니다")
	}
}

func TestDeleteSession_EmptyID(t *testing.T) {
	svc := newTestService()
	err := svc.DeleteSession(context.Background(), "")
	if err == nil {
		t.Fatal("빈 session_id에 에러가 반환되어야 합니다")
	}
}

func TestCloseSession_EmptyID(t *testing.T) {
	svc := newTestService()
	_, err := svc.CloseSession(context.Background(), "")
	if err == nil {
		t.Fatal("빈 session_id에 에러가 반환되어야 합니다")
	}
}

func TestGetSessionSummary_EmptyID(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetSessionSummary(context.Background(), "")
	if err == nil {
		t.Fatal("빈 session_id에 에러가 반환되어야 합니다")
	}
}

func TestGetSessionSummary_NotFound(t *testing.T) {
	svc := newTestService()
	_, err := svc.GetSessionSummary(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("존재하지 않는 세션에 에러가 반환되어야 합니다")
	}
}

func TestIntentClassification_Exercise(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, _, _, err := svc.SendCommand(ctx, "user1", "", "운동 계획 세워줘")
	if err != nil {
		t.Fatal(err)
	}
	if userTurn.Intent != "fitness.plan" {
		t.Errorf("Intent = %s, want fitness.plan", userTurn.Intent)
	}
}

func TestIntentClassification_Emergency(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, _, _, err := svc.SendCommand(ctx, "user1", "", "응급 상황이에요")
	if err != nil {
		t.Fatal(err)
	}
	if userTurn.Intent != "emergency.report" {
		t.Errorf("Intent = %s, want emergency.report", userTurn.Intent)
	}
}

func TestIntentClassification_Sleep(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, _, _, err := svc.SendCommand(ctx, "user1", "", "수면 패턴 분석해줘")
	if err != nil {
		t.Fatal(err)
	}
	if userTurn.Intent != "sleep.analysis" {
		t.Errorf("Intent = %s, want sleep.analysis", userTurn.Intent)
	}
}

func TestIntentClassification_General(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	userTurn, _, _, err := svc.SendCommand(ctx, "user1", "", "오늘 날씨 어때")
	if err != nil {
		t.Fatal(err)
	}
	if userTurn.Intent != "general.chat" {
		t.Errorf("Intent = %s, want general.chat", userTurn.Intent)
	}
}

func TestListSessions_DefaultLimit(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	svc.SendCommand(ctx, "user-dl", "", "첫번째")

	sessions, _, err := svc.ListSessions(ctx, "user-dl", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Errorf("sessions: got %d, want 1", len(sessions))
	}
}

func TestListTurns_DefaultLimit(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, session, _ := svc.SendCommand(ctx, "user1", "", "테스트")

	turns, _, err := svc.ListTurns(ctx, session.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(turns) != 2 {
		t.Errorf("turns: got %d, want 2", len(turns))
	}
}

func TestSendCommand_NonexistentSession(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, _, _, err := svc.SendCommand(ctx, "user1", "nonexistent-session", "메시지")
	if err == nil {
		t.Fatal("존재하지 않는 세션에 에러가 반환되어야 합니다")
	}
}

// ============================================================================
// Phase C-4 AI 통합 테스트
// ============================================================================

// mockAIResponder는 테스트용 AI 응답 생성기입니다.
type mockAIResponder struct {
	response *ai.AIResponse
	err      error
}

func (m *mockAIResponder) GenerateResponse(_ context.Context, _ []ai.ChatMessage, _ string) (*ai.AIResponse, error) {
	return m.response, m.err
}

func TestSendCommand_WithAIResponder_Success(t *testing.T) {
	svc := newTestService()
	svc.SetAIResponder(&mockAIResponder{
		response: &ai.AIResponse{
			Content: "AI가 생성한 건강 관리 조언입니다.",
		},
	})

	_, asstTurn, _, err := svc.SendCommand(context.Background(), "user1", "", "혈당 관리 방법")
	if err != nil {
		t.Fatalf("SendCommand 실패: %v", err)
	}
	if asstTurn.Content != "AI가 생성한 건강 관리 조언입니다." {
		t.Errorf("Content = %q, want AI 응답", asstTurn.Content)
	}
}

func TestSendCommand_WithAIResponder_Fallback(t *testing.T) {
	svc := newTestService()
	svc.SetAIResponder(&mockAIResponder{
		err: fmt.Errorf("API 오류"),
	})

	_, asstTurn, _, err := svc.SendCommand(context.Background(), "user1", "", "혈당 알려줘")
	if err != nil {
		t.Fatalf("SendCommand 실패: %v", err)
	}
	// AI 실패 시 키워드 기반 응답으로 폴백
	if asstTurn.Content == "" {
		t.Fatal("폴백 응답이 비어 있습니다")
	}
}

func TestSendCommand_WithNoopResponder(t *testing.T) {
	svc := newTestService()
	svc.SetAIResponder(ai.NewNoopResponder())

	_, asstTurn, _, err := svc.SendCommand(context.Background(), "user1", "", "예약 잡아줘")
	if err != nil {
		t.Fatalf("SendCommand 실패: %v", err)
	}
	// Noop은 nil 반환 → 키워드 기반 폴백
	if asstTurn.Content == "" {
		t.Fatal("폴백 응답이 비어 있습니다")
	}
	if asstTurn.Intent != "reservation.create" {
		t.Errorf("Intent = %q, want reservation.create", asstTurn.Intent)
	}
}

func TestSendCommand_WithAIResponder_ActionFallback(t *testing.T) {
	svc := newTestService()
	svc.SetAIResponder(&mockAIResponder{
		response: &ai.AIResponse{
			Content: "AI 응답: 혈당을 확인하겠습니다.",
			// ActionType/ActionResult 비어있음 → 인텐트 기반 액션 사용
		},
	})

	_, asstTurn, _, err := svc.SendCommand(context.Background(), "user1", "", "혈당 알려줘")
	if err != nil {
		t.Fatalf("SendCommand 실패: %v", err)
	}
	if asstTurn.Content != "AI 응답: 혈당을 확인하겠습니다." {
		t.Errorf("Content = %q, want AI 응답", asstTurn.Content)
	}
	// 인텐트 기반 액션이 사용되어야 함
	if asstTurn.ActionType != "query" {
		t.Errorf("ActionType = %q, want query (인텐트 기반 폴백)", asstTurn.ActionType)
	}
}
