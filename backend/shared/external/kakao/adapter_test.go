package kakao_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manpasik/backend/shared/external/kakao"
)

// TestValidateMessage_Valid는 정상 메시지 검증을 확인합니다.
func TestValidateMessage_Valid(t *testing.T) {
	msg := &kakao.Message{
		TemplateCode: "WELCOME_001",
		To:           "010-1234-5678",
	}
	if err := kakao.ValidateMessage(msg); err != nil {
		t.Errorf("ValidateMessage 실패: %v", err)
	}
}

// TestValidateMessage_InvalidPhone은 잘못된 전화번호 거부를 검증합니다.
func TestValidateMessage_InvalidPhone(t *testing.T) {
	cases := []string{"123", "abcd", "+1-555-1234"}
	for _, p := range cases {
		msg := &kakao.Message{TemplateCode: "X", To: p}
		if err := kakao.ValidateMessage(msg); err == nil {
			t.Errorf("잘못된 번호 %q 통과됨", p)
		}
	}
}

// TestValidateMessage_NoTemplate은 템플릿 코드 누락 거부를 검증합니다.
func TestValidateMessage_NoTemplate(t *testing.T) {
	msg := &kakao.Message{To: "010-1234-5678"}
	if err := kakao.ValidateMessage(msg); err == nil {
		t.Error("템플릿 코드 없이 통과됨")
	}
}

// TestExtractTemplateVariables는 변수 추출을 검증합니다.
func TestExtractTemplateVariables(t *testing.T) {
	body := "안녕하세요 #{name}님, 주문번호 #{order_id} 입니다. 다시 #{name}님께"
	vars := kakao.ExtractTemplateVariables(body)
	if len(vars) != 2 {
		t.Fatalf("vars = %d, want 2 (중복 제거)", len(vars))
	}
}

// TestRenderTemplate는 변수 치환을 검증합니다.
func TestRenderTemplate(t *testing.T) {
	body := "안녕하세요 #{name}님, 주문 #{order_id}이 #{status}되었습니다."
	rendered := kakao.RenderTemplate(body, map[string]string{
		"name":     "홍길동",
		"order_id": "ORD-123",
		"status":   "배송",
	})

	expected := "안녕하세요 홍길동님, 주문 ORD-123이 배송되었습니다."
	if rendered != expected {
		t.Errorf("rendered = %q, want %q", rendered, expected)
	}
}

// TestRenderTemplate_MissingVariable는 미치환 변수 유지를 검증합니다.
func TestRenderTemplate_MissingVariable(t *testing.T) {
	body := "안녕하세요 #{name}님, #{missing}"
	rendered := kakao.RenderTemplate(body, map[string]string{"name": "홍길동"})

	if rendered != "안녕하세요 홍길동님, #{missing}" {
		t.Errorf("rendered = %q", rendered)
	}
}

// TestNoopAdapter_Send는 Noop 발송을 검증합니다.
func TestNoopAdapter_Send(t *testing.T) {
	a := kakao.NewNoopAdapter()
	msg := &kakao.Message{
		TemplateCode: "WELCOME",
		To:           "010-1234-5678",
		Variables:    map[string]string{"name": "홍길동"},
	}

	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "noop" {
		t.Errorf("Provider = %q", result.Provider)
	}
	if a.Count() != 1 {
		t.Errorf("Count = %d, want 1", a.Count())
	}
}

// TestNoopAdapter_RegisterTemplate는 템플릿 등록을 검증합니다.
func TestNoopAdapter_RegisterTemplate(t *testing.T) {
	a := kakao.NewNoopAdapter()
	body := "주문 #{order_id} 발송 완료"
	if err := a.RegisterTemplate("ORDER_SHIPPED", body, nil); err != nil {
		t.Fatalf("RegisterTemplate 실패: %v", err)
	}

	tmpl, err := a.GetTemplate("ORDER_SHIPPED")
	if err != nil {
		t.Fatalf("GetTemplate 실패: %v", err)
	}
	if tmpl.Body != body {
		t.Errorf("Body = %q", tmpl.Body)
	}
	if len(tmpl.Variables) != 1 || tmpl.Variables[0] != "order_id" {
		t.Errorf("Variables = %v", tmpl.Variables)
	}
}

// TestNoopAdapter_GetTemplate_NotFound는 미등록 템플릿 거부를 검증합니다.
func TestNoopAdapter_GetTemplate_NotFound(t *testing.T) {
	a := kakao.NewNoopAdapter()
	_, err := a.GetTemplate("MISSING")
	if err == nil {
		t.Error("미등록 템플릿이 통과됨")
	}
}

// TestNHNCloudAdapter_Send_Success는 NHN Cloud 정상 발송을 검증합니다.
func TestNHNCloudAdapter_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-Secret-Key")
		if secret != "secret-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"header": {"isSuccessful": true, "resultCode": 0, "resultMessage": "SUCCESS"},
			"message": {"requestId": "req-123"}
		}`))
	}))
	defer server.Close()

	a := kakao.NewNHNCloudAdapter("app-key", "secret-key", "sender-key")
	a.SetEndpoint(server.URL)

	msg := &kakao.Message{
		TemplateCode: "WELCOME",
		To:           "010-1234-5678",
		Variables:    map[string]string{"name": "홍길동"},
	}

	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.MessageID != "req-123" {
		t.Errorf("MessageID = %q", result.MessageID)
	}
	if result.Status != "sent" {
		t.Errorf("Status = %q", result.Status)
	}
}

// TestNHNCloudAdapter_Send_Failure는 실패 응답 처리를 검증합니다.
func TestNHNCloudAdapter_Send_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"header": {"isSuccessful": false, "resultCode": 4001, "resultMessage": "INVALID_TEMPLATE"}
		}`))
	}))
	defer server.Close()

	a := kakao.NewNHNCloudAdapter("k", "s", "sender")
	a.SetEndpoint(server.URL)

	msg := &kakao.Message{TemplateCode: "X", To: "010-1234-5678"}
	result, _ := a.Send(context.Background(), msg)
	if result == nil {
		t.Fatal("result nil")
	}
	if result.Status != "failed" {
		t.Errorf("Status = %q", result.Status)
	}
	if result.ErrorCode != "4001" {
		t.Errorf("ErrorCode = %q", result.ErrorCode)
	}
}

// TestNHNCloudAdapter_HealthCheck는 헬스체크를 검증합니다.
func TestNHNCloudAdapter_HealthCheck(t *testing.T) {
	a := kakao.NewNHNCloudAdapter("k", "s", "sender")
	if err := a.HealthCheck(context.Background()); err != nil {
		t.Errorf("HealthCheck 실패: %v", err)
	}

	a2 := kakao.NewNHNCloudAdapter("", "", "")
	if err := a2.HealthCheck(context.Background()); err == nil {
		t.Error("자격증명 없이 헬스체크 통과")
	}
}

// TestNHNCloudAdapter_RegisterTemplate는 템플릿 등록을 검증합니다.
func TestNHNCloudAdapter_RegisterTemplate(t *testing.T) {
	a := kakao.NewNHNCloudAdapter("k", "s", "sender")
	if err := a.RegisterTemplate("T1", "안녕 #{name}", nil); err != nil {
		t.Fatalf("RegisterTemplate 실패: %v", err)
	}
	tmpl, err := a.GetTemplate("T1")
	if err != nil {
		t.Fatalf("GetTemplate 실패: %v", err)
	}
	if tmpl.Status != "approved" {
		t.Errorf("Status = %q", tmpl.Status)
	}
}

// TestKakaoBizAdapter_Send는 KakaoBiz 발송을 검증합니다.
func TestKakaoBizAdapter_Send(t *testing.T) {
	a := kakao.NewKakaoBizAdapter("api-key", "sender-key")
	msg := &kakao.Message{
		TemplateCode: "T",
		To:           "010-1234-5678",
	}
	result, err := a.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("Send 실패: %v", err)
	}
	if result.Provider != "kakaobiz" {
		t.Errorf("Provider = %q", result.Provider)
	}
}

// TestNewFromEnv_Default는 기본 Noop 폴백을 검증합니다.
func TestNewFromEnv_Default(t *testing.T) {
	t.Setenv("KAKAO_PROVIDER", "")
	a := kakao.NewFromEnv()
	if a.Provider() != "noop" {
		t.Errorf("Provider = %q, want noop", a.Provider())
	}
}

// TestNewFromEnv_NHNCloud는 nhncloud 스위칭을 검증합니다.
func TestNewFromEnv_NHNCloud(t *testing.T) {
	t.Setenv("KAKAO_PROVIDER", "nhncloud")
	t.Setenv("NHNCLOUD_APP_KEY", "key")
	a := kakao.NewFromEnv()
	if a.Provider() != "nhncloud" {
		t.Errorf("Provider = %q, want nhncloud", a.Provider())
	}
}

// TestPhoneFormats는 다양한 전화번호 형식을 검증합니다.
func TestPhoneFormats(t *testing.T) {
	valid := []string{
		"010-1234-5678",
		"01012345678",
		"+821012345678",
	}
	for _, p := range valid {
		msg := &kakao.Message{TemplateCode: "X", To: p}
		if err := kakao.ValidateMessage(msg); err != nil {
			t.Errorf("정상 번호 %q 거부됨: %v", p, err)
		}
	}
}
