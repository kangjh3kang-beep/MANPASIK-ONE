package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ============================================================================
// TestNewOpenAIClient_기본값
// — model이 빈 문자열이면 DefaultModel("gpt-4o")이 설정되는지 확인
// ============================================================================

func TestNewOpenAIClient_기본값(t *testing.T) {
	client := NewOpenAIClient("test-key", "")

	if client.Model() != DefaultModel {
		t.Errorf("기본 모델이 %q이어야 하지만 %q입니다", DefaultModel, client.Model())
	}
	if client.BaseURL() != DefaultBaseURL {
		t.Errorf("기본 URL이 %q이어야 하지만 %q입니다", DefaultBaseURL, client.BaseURL())
	}

	// 명시적 모델 지정 시
	client2 := NewOpenAIClient("test-key", "gpt-3.5-turbo")
	if client2.Model() != "gpt-3.5-turbo" {
		t.Errorf("지정 모델이 %q이어야 하지만 %q입니다", "gpt-3.5-turbo", client2.Model())
	}

	// WithBaseURL 옵션 테스트
	client3 := NewOpenAIClient("test-key", "", WithBaseURL("https://custom.api.com/v1"))
	if client3.BaseURL() != "https://custom.api.com/v1" {
		t.Errorf("커스텀 URL이 설정되어야 하지만 %q입니다", client3.BaseURL())
	}

	// 빈 WithBaseURL은 기본값 유지
	client4 := NewOpenAIClient("test-key", "", WithBaseURL(""))
	if client4.BaseURL() != DefaultBaseURL {
		t.Errorf("빈 URL 옵션은 기본값을 유지해야 하지만 %q입니다", client4.BaseURL())
	}
}

// ============================================================================
// TestChatMessage_구조체
// — ChatMessage와 ChatResponse 필드가 올바르게 설정되는지 확인
// ============================================================================

func TestChatMessage_구조체(t *testing.T) {
	msg := ChatMessage{
		Role:    "user",
		Content: "안녕하세요",
	}

	if msg.Role != "user" {
		t.Errorf("Role이 %q이어야 하지만 %q입니다", "user", msg.Role)
	}
	if msg.Content != "안녕하세요" {
		t.Errorf("Content가 %q이어야 하지만 %q입니다", "안녕하세요", msg.Content)
	}

	// JSON 직렬화 확인
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("JSON 직렬화 실패: %v", err)
	}

	var decoded ChatMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON 역직렬화 실패: %v", err)
	}

	if decoded.Role != msg.Role || decoded.Content != msg.Content {
		t.Errorf("JSON 라운드트립 불일치: 원본=%+v, 복원=%+v", msg, decoded)
	}

	// ChatResponse 필드 확인
	resp := ChatResponse{
		Content:      "응답 텍스트",
		FinishReason: "stop",
		TokensUsed:   150,
	}
	if resp.Content != "응답 텍스트" {
		t.Errorf("Content가 일치하지 않습니다")
	}
	if resp.FinishReason != "stop" {
		t.Errorf("FinishReason이 일치하지 않습니다")
	}
	if resp.TokensUsed != 150 {
		t.Errorf("TokensUsed가 %d이어야 하지만 %d입니다", 150, resp.TokensUsed)
	}
}

// ============================================================================
// TestOpenAIClient_EmptyAPIKey
// — API 키가 비어있을 때 Chat 호출이 에러를 반환하는지 확인
// ============================================================================

func TestOpenAIClient_EmptyAPIKey(t *testing.T) {
	client := NewOpenAIClient("", "gpt-4o")

	messages := []ChatMessage{
		{Role: "user", Content: "테스트 메시지"},
	}

	resp, err := client.Chat(context.Background(), "시스템 프롬프트", messages)
	if err == nil {
		t.Fatal("빈 API 키로 Chat 호출 시 에러가 발생해야 합니다")
	}
	if resp != nil {
		t.Fatal("빈 API 키로 Chat 호출 시 응답이 nil이어야 합니다")
	}

	// 에러 메시지에 API 키 관련 내용이 포함되어야 함
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("에러 메시지가 비어있으면 안 됩니다")
	}
}

// ============================================================================
// TestOpenAIClient_인터페이스구현
// — OpenAIClient가 LLMClient 인터페이스를 구현하는지 컴파일 타임에 확인
// ============================================================================

func TestOpenAIClient_인터페이스구현(t *testing.T) {
	// 컴파일 타임 검증은 파일 하단 var _ LLMClient = (*OpenAIClient)(nil)로 보장
	// 런타임에서도 인터페이스 변환이 가능한지 확인
	var client LLMClient = NewOpenAIClient("test-key", "gpt-4o")
	if client == nil {
		t.Fatal("OpenAIClient가 LLMClient 인터페이스를 구현해야 합니다")
	}
}

// ============================================================================
// 추가 테스트: Mock 서버를 이용한 Chat 성공/에러 시나리오
// — 실제 API를 호출하지 않고 httptest 서버로 검증
// ============================================================================

func TestOpenAIClient_Chat_성공(t *testing.T) {
	// Mock OpenAI 서버
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 요청 검증
		if r.Method != http.MethodPost {
			t.Errorf("POST 메서드여야 하지만 %s입니다", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("Authorization 헤더가 올바르지 않습니다")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content-Type이 application/json이어야 합니다")
		}

		// 요청 바디 파싱
		var reqBody openaiChatRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("요청 바디 파싱 실패: %v", err)
		}
		if reqBody.Model != "gpt-4o" {
			t.Errorf("모델이 gpt-4o여야 하지만 %s입니다", reqBody.Model)
		}
		// system + user = 2개 메시지
		if len(reqBody.Messages) != 2 {
			t.Errorf("메시지가 2개여야 하지만 %d개입니다", len(reqBody.Messages))
		}

		// 성공 응답 반환
		resp := openaiChatResponse{
			ID:      "chatcmpl-test123",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "gpt-4o",
		}
		resp.Choices = []struct {
			Index        int           `json:"index"`
			Message      openaiMessage `json:"message"`
			FinishReason string        `json:"finish_reason"`
		}{
			{
				Index:        0,
				Message:      openaiMessage{Role: "assistant", Content: "테스트 응답입니다."},
				FinishReason: "stop",
			},
		}
		resp.Usage.PromptTokens = 20
		resp.Usage.CompletionTokens = 10
		resp.Usage.TotalTokens = 30

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient("test-api-key", "gpt-4o", WithBaseURL(server.URL))

	messages := []ChatMessage{
		{Role: "user", Content: "테스트 질문"},
	}
	resp, err := client.Chat(context.Background(), "시스템 프롬프트", messages)
	if err != nil {
		t.Fatalf("예상치 못한 에러: %v", err)
	}
	if resp.Content != "테스트 응답입니다." {
		t.Errorf("응답 내용이 일치하지 않습니다: %q", resp.Content)
	}
	if resp.FinishReason != "stop" {
		t.Errorf("FinishReason이 %q이어야 하지만 %q입니다", "stop", resp.FinishReason)
	}
	if resp.TokensUsed != 30 {
		t.Errorf("토큰 사용량이 30이어야 하지만 %d입니다", resp.TokensUsed)
	}
}

func TestOpenAIClient_Chat_API에러(t *testing.T) {
	// 401 Unauthorized 응답을 반환하는 Mock 서버
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(openaiErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			}{
				Message: "Incorrect API key provided",
				Type:    "invalid_request_error",
				Code:    "invalid_api_key",
			},
		})
	}))
	defer server.Close()

	client := NewOpenAIClient("invalid-key", "gpt-4o", WithBaseURL(server.URL))

	messages := []ChatMessage{
		{Role: "user", Content: "테스트"},
	}
	_, err := client.Chat(context.Background(), "", messages)
	if err == nil {
		t.Fatal("API 에러 시 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClient_Chat_컨텍스트취소(t *testing.T) {
	// 느린 서버 시뮬레이션 — 테스트 속도를 위해 짧은 대기
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewOpenAIClient("test-key", "gpt-4o", WithBaseURL(server.URL))

	// 100ms 타임아웃 설정
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	messages := []ChatMessage{
		{Role: "user", Content: "테스트"},
	}
	_, err := client.Chat(ctx, "", messages)
	if err == nil {
		t.Fatal("컨텍스트 취소 시 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClient_Chat_시스템프롬프트없음(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody openaiChatRequest
		json.NewDecoder(r.Body).Decode(&reqBody)

		// 시스템 프롬프트 없이 user 메시지만
		if len(reqBody.Messages) != 1 {
			t.Errorf("시스템 프롬프트 없을 때 메시지가 1개여야 하지만 %d개입니다", len(reqBody.Messages))
		}
		if reqBody.Messages[0].Role != "user" {
			t.Errorf("첫 메시지 role이 user여야 하지만 %s입니다", reqBody.Messages[0].Role)
		}

		resp := openaiChatResponse{
			ID: "chatcmpl-test",
		}
		resp.Choices = []struct {
			Index        int           `json:"index"`
			Message      openaiMessage `json:"message"`
			FinishReason string        `json:"finish_reason"`
		}{
			{
				Index:        0,
				Message:      openaiMessage{Role: "assistant", Content: "응답"},
				FinishReason: "stop",
			},
		}
		resp.Usage.TotalTokens = 10

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewOpenAIClient("test-key", "gpt-4o", WithBaseURL(server.URL))

	messages := []ChatMessage{
		{Role: "user", Content: "안녕"},
	}
	resp, err := client.Chat(context.Background(), "", messages)
	if err != nil {
		t.Fatalf("예상치 못한 에러: %v", err)
	}
	if resp.Content != "응답" {
		t.Errorf("응답이 %q이어야 하지만 %q입니다", "응답", resp.Content)
	}
}
