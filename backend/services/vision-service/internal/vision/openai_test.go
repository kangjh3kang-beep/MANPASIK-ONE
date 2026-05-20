package vision

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// makeVisionResponse builds a mock OpenAI Vision API JSON response body.
func makeVisionResponse(content string) string {
	resp := visionResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{Message: struct {
				Content string `json:"content"`
			}{Content: content}},
		},
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

func TestAnalyzeFood_Success(t *testing.T) {
	foodJSON := `[{"name":"김치찌개","confidence":0.92,"calorie_kcal":180,"portion_g":300,"nutrients":[{"name":"단백질","amount":15,"unit":"g","dv":0.3}]}]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request basics
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			t.Errorf("expected 'Bearer test-api-key', got %q", auth)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(makeVisionResponse(foodJSON)))
	}))
	defer srv.Close()

	analyzer := NewOpenAIVisionAnalyzer("test-api-key", "gpt-4o", srv.URL)
	items, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Name != "김치찌개" {
		t.Errorf("expected name 김치찌개, got %s", items[0].Name)
	}
	if items[0].Confidence != 0.92 {
		t.Errorf("expected confidence 0.92, got %f", items[0].Confidence)
	}
	if items[0].CalorieKcal != 180 {
		t.Errorf("expected calorie_kcal 180, got %f", items[0].CalorieKcal)
	}
	if items[0].PortionG != 300 {
		t.Errorf("expected portion_g 300, got %f", items[0].PortionG)
	}
	if len(items[0].Nutrients) != 1 {
		t.Fatalf("expected 1 nutrient, got %d", len(items[0].Nutrients))
	}
	n := items[0].Nutrients[0]
	if n.Name != "단백질" {
		t.Errorf("expected nutrient name 단백질, got %s", n.Name)
	}
	if n.Amount != 15 {
		t.Errorf("expected nutrient amount 15, got %f", n.Amount)
	}
	if n.Unit != "g" {
		t.Errorf("expected nutrient unit g, got %s", n.Unit)
	}
	if n.DV != 0.3 {
		t.Errorf("expected nutrient dv 0.3, got %f", n.DV)
	}
}

func TestAnalyzeFood_EmptyAPIKey(t *testing.T) {
	analyzer := NewOpenAIVisionAnalyzer("", "gpt-4o", "http://localhost")
	_, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error for empty API key, got nil")
	}
	expected := "Vision API 키가 설정되지 않았습니다"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestAnalyzeFood_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer srv.Close()

	analyzer := NewOpenAIVisionAnalyzer("test-key", "gpt-4o", srv.URL)
	_, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}

func TestAnalyzeFood_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	analyzer := NewOpenAIVisionAnalyzer("test-key", "gpt-4o", srv.URL)
	_, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
	expected := "Vision API 응답에 결과가 없습니다"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestAnalyzeFood_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{not valid json`))
	}))
	defer srv.Close()

	analyzer := NewOpenAIVisionAnalyzer("test-key", "gpt-4o", srv.URL)
	_, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestAnalyzeFood_InvalidFoodJSON(t *testing.T) {
	// Valid OpenAI response structure, but the content is not valid food JSON
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(makeVisionResponse("this is not json at all")))
	}))
	defer srv.Close()

	analyzer := NewOpenAIVisionAnalyzer("test-key", "gpt-4o", srv.URL)
	_, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error for invalid food JSON in content, got nil")
	}
}

func TestLogVisionAnalyzer_ReturnsError(t *testing.T) {
	analyzer := NewLogVisionAnalyzer()
	items, err := analyzer.AnalyzeFood(context.Background(), "https://example.com/food.jpg")
	if err == nil {
		t.Fatal("expected error from LogVisionAnalyzer, got nil")
	}
	if items != nil {
		t.Errorf("expected nil items, got %v", items)
	}
	expected := "Vision API 미설정"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
