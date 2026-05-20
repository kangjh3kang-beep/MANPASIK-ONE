package classifier_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manpasik/backend/services/nlp-service/internal/classifier"
)

func TestOpenAIClassifier_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Authorization 헤더 누락")
		}
		result := `{"intent":"health_inquiry.blood_sugar","entities":["blood_sugar"],"confidence":0.95,"urgency":"routine"}`
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": result}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey:  "test-key",
		Model:   "gpt-4o-mini",
		BaseURL: server.URL,
	})

	result, err := cls.Classify(context.Background(), "혈당 수치 알려줘")
	if err != nil {
		t.Fatalf("Classify 실패: %v", err)
	}
	if result == nil {
		t.Fatal("result가 nil입니다")
	}
	if result.Intent != "health_inquiry.blood_sugar" {
		t.Errorf("Intent = %q, want %q", result.Intent, "health_inquiry.blood_sugar")
	}
	if result.Confidence != 0.95 {
		t.Errorf("Confidence = %f, want 0.95", result.Confidence)
	}
	if result.Urgency != "routine" {
		t.Errorf("Urgency = %q, want %q", result.Urgency, "routine")
	}
}

func TestOpenAIClassifier_EmptyAPIKey(t *testing.T) {
	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey: "",
	})

	_, err := cls.Classify(context.Background(), "test")
	if err == nil {
		t.Fatal("빈 API 키에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClassifier_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := cls.Classify(context.Background(), "test")
	if err == nil {
		t.Fatal("HTTP 503에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClassifier_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{"choices": []map[string]interface{}{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := cls.Classify(context.Background(), "test")
	if err == nil {
		t.Fatal("빈 선택지에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClassifier_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": "not json at all"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	_, err := cls.Classify(context.Background(), "test")
	if err == nil {
		t.Fatal("잘못된 JSON 응답에 에러가 반환되어야 합니다")
	}
}

func TestOpenAIClassifier_EmergencyUrgency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := `{"intent":"health_inquiry","entities":["흉통","호흡곤란"],"confidence":0.98,"urgency":"emergency"}`
		resp := map[string]interface{}{
			"choices": []map[string]interface{}{
				{"message": map[string]string{"content": result}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cls := classifier.NewOpenAIClassifier(classifier.OpenAIConfig{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	result, err := cls.Classify(context.Background(), "흉통이 심하고 호흡이 곤란합니다")
	if err != nil {
		t.Fatalf("Classify 실패: %v", err)
	}
	if result.Urgency != "emergency" {
		t.Errorf("Urgency = %q, want %q", result.Urgency, "emergency")
	}
}

func TestNoopClassifier(t *testing.T) {
	cls := classifier.NewNoopClassifier()
	result, err := cls.Classify(context.Background(), "test")
	if err != nil {
		t.Fatalf("NoopClassifier는 에러를 반환하지 않아야 합니다: %v", err)
	}
	if result != nil {
		t.Fatal("NoopClassifier는 nil을 반환해야 합니다")
	}
}
