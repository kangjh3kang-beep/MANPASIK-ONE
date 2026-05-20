package translator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTranslate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Verify query parameters
		q := r.URL.Query()
		if q.Get("q") != "안녕하세요" {
			t.Errorf("expected q=안녕하세요, got %q", q.Get("q"))
		}
		if q.Get("target") != "en" {
			t.Errorf("expected target=en, got %q", q.Get("target"))
		}
		if q.Get("source") != "ko" {
			t.Errorf("expected source=ko, got %q", q.Get("source"))
		}
		if q.Get("format") != "text" {
			t.Errorf("expected format=text, got %q", q.Get("format"))
		}
		if q.Get("key") != "test-api-key" {
			t.Errorf("expected key=test-api-key, got %q", q.Get("key"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"translations":[{"translatedText":"Hello","detectedSourceLanguage":"ko"}]}}`))
	}))
	defer srv.Close()

	tr := NewGoogleTranslator("test-api-key")
	tr.SetAPIURL(srv.URL)

	result, err := tr.Translate(context.Background(), "안녕하세요", "ko", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello" {
		t.Errorf("expected 'Hello', got %q", result)
	}
}

func TestTranslate_EmptyText(t *testing.T) {
	tr := NewGoogleTranslator("test-api-key")
	_, err := tr.Translate(context.Background(), "", "ko", "en")
	if err == nil {
		t.Fatal("expected error for empty text, got nil")
	}
	expected := "번역할 텍스트가 비어있습니다"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestTranslate_EmptyTargetLang(t *testing.T) {
	tr := NewGoogleTranslator("test-api-key")
	_, err := tr.Translate(context.Background(), "hello", "en", "")
	if err == nil {
		t.Fatal("expected error for empty target lang, got nil")
	}
	expected := "대상 언어가 비어있습니다"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestTranslate_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"bad request"}}`))
	}))
	defer srv.Close()

	tr := NewGoogleTranslator("test-api-key")
	tr.SetAPIURL(srv.URL)

	_, err := tr.Translate(context.Background(), "hello", "en", "ko")
	if err == nil {
		t.Fatal("expected error for HTTP 400, got nil")
	}
	if got := err.Error(); got == "" {
		t.Error("expected non-empty error message")
	}
}

func TestTranslate_EmptyTranslations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"translations":[]}}`))
	}))
	defer srv.Close()

	tr := NewGoogleTranslator("test-api-key")
	tr.SetAPIURL(srv.URL)

	_, err := tr.Translate(context.Background(), "hello", "en", "ko")
	if err == nil {
		t.Fatal("expected error for empty translations, got nil")
	}
	expected := "Google Translate 응답에 번역 결과가 없습니다"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestTranslate_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{not valid json`))
	}))
	defer srv.Close()

	tr := NewGoogleTranslator("test-api-key")
	tr.SetAPIURL(srv.URL)

	_, err := tr.Translate(context.Background(), "hello", "en", "ko")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestTranslate_AutoDetectSource(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// When sourceLang is empty, the "source" param should NOT be present
		if q.Get("source") != "" {
			t.Errorf("expected no 'source' param for auto-detect, got %q", q.Get("source"))
		}
		if q.Get("q") != "Bonjour" {
			t.Errorf("expected q=Bonjour, got %q", q.Get("q"))
		}
		if q.Get("target") != "en" {
			t.Errorf("expected target=en, got %q", q.Get("target"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"translations":[{"translatedText":"Hello","detectedSourceLanguage":"fr"}]}}`))
	}))
	defer srv.Close()

	tr := NewGoogleTranslator("test-api-key")
	tr.SetAPIURL(srv.URL)

	result, err := tr.Translate(context.Background(), "Bonjour", "", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello" {
		t.Errorf("expected 'Hello', got %q", result)
	}
}

func TestNoopTranslator_ReturnsError(t *testing.T) {
	tr := NewNoopTranslator()
	result, err := tr.Translate(context.Background(), "hello", "en", "ko")
	if err == nil {
		t.Fatal("expected error from NoopTranslator, got nil")
	}
	if result != "" {
		t.Errorf("expected empty result, got %q", result)
	}
	expected := "외부 번역 API 미설정"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
