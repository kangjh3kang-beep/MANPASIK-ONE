package ai_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manpasik/backend/services/assistant-service/internal/ai"
)

func TestOpenAIResponder_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Authorization 헤더 누락")
		}
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "혈당 관리가 중요합니다."}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey:  "test-key",
		Model:   "gpt-4o",
		BaseURL: server.URL,
	})

	result, err := responder.GenerateResponse(context.Background(), nil, "혈당 관리 방법 알려줘")
	if err != nil {
		t.Fatalf("GenerateResponse 실패: %v", err)
	}
	if result == nil {
		t.Fatal("result가 nil입니다")
	}
	if result.Content != "혈당 관리가 중요합니다." {
		t.Errorf("Content = %q, want %q", result.Content, "혈당 관리가 중요합니다.")
	}
}

func TestOpenAIResponder_WithHistory(t *testing.T) {
	var receivedMessages int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Messages []map[string]string `json:"messages"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		receivedMessages = len(req.Messages)
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "네, 혈압도 확인하겠습니다."}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	history := []ai.ChatMessage{
		{Role: "user", Content: "혈당 알려줘"},
		{Role: "assistant", Content: "혈당은 95mg/dL입니다."},
	}

	result, err := responder.GenerateResponse(context.Background(), history, "혈압도 알려줘")
	if err != nil {
		t.Fatalf("GenerateResponse 실패: %v", err)
	}
	if result == nil || result.Content == "" {
		t.Fatal("응답이 비어 있습니다")
	}
	// system(1) + history(2) + user(1) = 4
	if receivedMessages != 4 {
		t.Errorf("메시지 수 = %d, want 4 (system+history+user)", receivedMessages)
	}
}

func TestOpenAIResponder_EmptyAPIKey(t *testing.T) {
	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey: "",
	})

	_, err := responder.GenerateResponse(context.Background(), nil, "test")
	if err == nil {
		t.Fatal("빈 API 키에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIResponder_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := responder.GenerateResponse(context.Background(), nil, "test")
	if err == nil {
		t.Fatal("HTTP 500에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIResponder_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{"choices": []map[string]interface{}{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := responder.GenerateResponse(context.Background(), nil, "test")
	if err == nil {
		t.Fatal("빈 선택지에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIResponder_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	responder := ai.NewOpenAIResponder(ai.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := responder.GenerateResponse(context.Background(), nil, "test")
	if err == nil {
		t.Fatal("잘못된 JSON에 에러가 반환되어야 합니다")
	}
}

func TestNoopResponder(t *testing.T) {
	responder := ai.NewNoopResponder()
	result, err := responder.GenerateResponse(context.Background(), nil, "test")
	if err != nil {
		t.Fatalf("NoopResponder는 에러를 반환하지 않아야 합니다: %v", err)
	}
	if result != nil {
		t.Fatal("NoopResponder는 nil을 반환해야 합니다")
	}
}
