package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestRecordRequest(t *testing.T) {
	m := NewMetrics()

	// Record multiple requests
	m.RecordRequest("GET", 100*time.Millisecond, 200)
	m.RecordRequest("GET", 200*time.Millisecond, 200)
	m.RecordRequest("POST", 150*time.Millisecond, 201)
	m.RecordRequest("GET", 50*time.Millisecond, 500)

	stats := m.GetStats()

	totalReqs, ok := stats["total_requests"].(int64)
	if !ok || totalReqs != 4 {
		t.Errorf("expected total_requests=4, got %v", stats["total_requests"])
	}

	totalErrs, ok := stats["total_errors"].(int64)
	if !ok || totalErrs != 1 {
		t.Errorf("expected total_errors=1, got %v", stats["total_errors"])
	}

	methods, ok := stats["methods"].(map[string]int64)
	if !ok {
		t.Fatal("expected methods to be map[string]int64")
	}
	if methods["GET"] != 3 {
		t.Errorf("expected GET count=3, got %d", methods["GET"])
	}
	if methods["POST"] != 1 {
		t.Errorf("expected POST count=1, got %d", methods["POST"])
	}

	uptime, ok := stats["uptime_seconds"].(float64)
	if !ok || uptime < 0 {
		t.Errorf("expected positive uptime, got %v", stats["uptime_seconds"])
	}
}

func TestPrometheusHandler(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("GET", 100*time.Millisecond, 200)
	m.RecordRequest("POST", 200*time.Millisecond, 500)

	handler := m.PrometheusHandler()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/plain; version=0.0.4" {
		t.Errorf("expected Prometheus content type, got %s", contentType)
	}

	body := w.Body.String()

	// Verify Prometheus format elements
	checks := []string{
		"# HELP manpasik_uptime_seconds",
		"# TYPE manpasik_uptime_seconds gauge",
		"manpasik_uptime_seconds",
		"# HELP manpasik_requests_total",
		"# TYPE manpasik_requests_total counter",
		"manpasik_requests_total{method=\"GET\"}",
		"manpasik_requests_total{method=\"POST\"}",
		"# HELP manpasik_errors_total",
		"# TYPE manpasik_errors_total counter",
		"manpasik_errors_total{method=\"POST\"}",
		"# HELP manpasik_request_duration_seconds",
		"# TYPE manpasik_request_duration_seconds gauge",
		"manpasik_request_duration_seconds{method=\"GET\"}",
		"manpasik_request_duration_seconds{method=\"POST\"}",
	}

	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Errorf("expected body to contain %q\nbody:\n%s", check, body)
		}
	}
}

func TestHealthCheck(t *testing.T) {
	hc := NewHealthCheck("test-service", "1.0.0")

	handler := hc.Handler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected application/json, got %s", contentType)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("expected status=healthy, got %v", result["status"])
	}
	if result["service"] != "test-service" {
		t.Errorf("expected service=test-service, got %v", result["service"])
	}
	if result["version"] != "1.0.0" {
		t.Errorf("expected version=1.0.0, got %v", result["version"])
	}
	if _, ok := result["goroutines"]; !ok {
		t.Error("expected goroutines field in response")
	}
	if _, ok := result["memory_mb"]; !ok {
		t.Error("expected memory_mb field in response")
	}
	if _, ok := result["uptime"]; !ok {
		t.Error("expected uptime field in response")
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	m := NewMetrics()
	interceptor := UnaryServerInterceptor(m)

	// Test successful handler
	successHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/manpasik.v1.AuthService/Login",
	}

	resp, err := interceptor(context.Background(), nil, info, successHandler)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Errorf("expected resp=ok, got %v", resp)
	}

	// Test error handler
	errHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, context.DeadlineExceeded
	}

	infoErr := &grpc.UnaryServerInfo{
		FullMethod: "/manpasik.v1.AuthService/Verify",
	}

	resp2, err2 := interceptor(context.Background(), nil, infoErr, errHandler)
	if err2 == nil {
		t.Error("expected error, got nil")
	}
	if resp2 != nil {
		t.Errorf("expected nil resp, got %v", resp2)
	}

	// Verify metrics were recorded
	stats := m.GetStats()
	totalReqs, ok := stats["total_requests"].(int64)
	if !ok || totalReqs != 2 {
		t.Errorf("expected total_requests=2, got %v", stats["total_requests"])
	}
	totalErrs, ok := stats["total_errors"].(int64)
	if !ok || totalErrs != 1 {
		t.Errorf("expected total_errors=1 (from error handler), got %v", stats["total_errors"])
	}

	methods, ok := stats["methods"].(map[string]int64)
	if !ok {
		t.Fatal("expected methods to be map[string]int64")
	}
	if methods["/manpasik.v1.AuthService/Login"] != 1 {
		t.Errorf("expected Login count=1, got %d", methods["/manpasik.v1.AuthService/Login"])
	}
	if methods["/manpasik.v1.AuthService/Verify"] != 1 {
		t.Errorf("expected Verify count=1, got %d", methods["/manpasik.v1.AuthService/Verify"])
	}
}
