package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newNilHandler는 모든 gRPC 클라이언트가 nil인 핸들러를 반환합니다.
func newNilHandler() *RestHandler {
	return &RestHandler{}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("상태 코드 = %d, 원하는 값 = %d", got, want)
	}
}

func assertErrorContains(t *testing.T, body []byte, substr string) {
	t.Helper()
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("JSON 파싱 실패: %v", err)
	}
	if !strings.Contains(result["error"], substr) {
		t.Errorf("에러 메시지 %q에 %q 포함 안 됨", result["error"], substr)
	}
}

// ---------------------------------------------------------------------------
// 1. Health check
// ---------------------------------------------------------------------------

func TestHealthCheck(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)
}

// ---------------------------------------------------------------------------
// 2. Auth — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestLogin_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"email":"a@b.com","password":"pw"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleLogin(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
	assertErrorContains(t, w.Body.Bytes(), "auth service unavailable")
}

func TestRegister_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"email":"a@b.com","password":"pw","display_name":"test"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleRegister(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestRefreshToken_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"refresh_token":"rt_123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleRefreshToken(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestLogout_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"user1"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleLogout(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 3. User — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestGetProfile_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/users/user1/profile", nil)
	w := httptest.NewRecorder()
	h.handleGetProfile(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestUpdateProfile_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"display_name":"test"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/user1/profile", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleUpdateProfile(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 4. Measurement — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestStartSession_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"device_id":"d1","user_id":"u1","cartridge_id":"c1"}`
	req := httptest.NewRequest("POST", "/api/v1/measurements/sessions", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleStartSession(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestGetMeasurementHistory_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/measurements/history?user_id=u1", nil)
	w := httptest.NewRecorder()
	h.handleGetMeasurementHistory(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 5. Device — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestRegisterDevice_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"device_id":"d1","user_id":"u1","serial_number":"SN123"}`
	req := httptest.NewRequest("POST", "/api/v1/devices", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleRegisterDevice(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestListDevices_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/devices?user_id=u1", nil)
	w := httptest.NewRecorder()
	h.handleListDevices(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 6. Shop — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestListProducts_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/products", nil)
	w := httptest.NewRecorder()
	h.handleListProducts(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestAddToCart_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"u1","product_id":"p1","quantity":1}`
	req := httptest.NewRequest("POST", "/api/v1/cart", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleAddToCart(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestCreateOrder_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"u1","shipping_address":"Seoul"}`
	req := httptest.NewRequest("POST", "/api/v1/orders", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleCreateOrder(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 7. Community — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestListPosts_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/posts", nil)
	w := httptest.NewRecorder()
	h.handleListPosts(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestCreatePost_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"author_id":"u1","title":"test","content":"hello"}`
	req := httptest.NewRequest("POST", "/api/v1/posts", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleCreatePost(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 8. Notification — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestSendNotification_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"u1","type":1,"title":"test","body":"msg"}`
	req := httptest.NewRequest("POST", "/api/v1/notifications", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleSendNotification(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestListNotifications_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/notifications?user_id=u1", nil)
	w := httptest.NewRecorder()
	h.handleListNotifications(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 9. Translation — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestTranslateText_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"text":"hello","source_language":"en","target_language":"ko"}`
	req := httptest.NewRequest("POST", "/api/v1/translations/text", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleTranslateText(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 10. Subscription — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestCreateSubscription_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"u1","tier":1}`
	req := httptest.NewRequest("POST", "/api/v1/subscriptions", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleCreateSubscription(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 11. Coaching — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestSetHealthGoal_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	body := `{"user_id":"u1","metric_name":"steps","target_value":10000}`
	req := httptest.NewRequest("POST", "/api/v1/coaching/goals", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.handleSetHealthGoal(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 12. Admin — ServiceUnavailable
// ---------------------------------------------------------------------------

func TestGetSystemStats_ServiceUnavailable(t *testing.T) {
	h := newNilHandler()
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	w := httptest.NewRecorder()
	h.handleGetSystemStats(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

// ---------------------------------------------------------------------------
// 유틸리티 함수 테스트
// ---------------------------------------------------------------------------

func TestWriteJSON_ContentType(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"key": "value"})

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, 원하는 값 = %q", ct, "application/json")
	}
	assertStatus(t, w.Code, http.StatusOK)
}

func TestWriteError_Format(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "test error")

	assertStatus(t, w.Code, http.StatusBadRequest)
	assertErrorContains(t, w.Body.Bytes(), "test error")
}

func TestQueryInt_DefaultValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	got := queryInt(req, "limit", 42)
	if got != 42 {
		t.Errorf("queryInt() = %d, 원하는 값 = 42", got)
	}
}

func TestQueryInt_ParseValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=100", nil)
	got := queryInt(req, "limit", 42)
	if got != 100 {
		t.Errorf("queryInt() = %d, 원하는 값 = 100", got)
	}
}

func TestQueryInt_InvalidValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?limit=abc", nil)
	got := queryInt(req, "limit", 42)
	if got != 42 {
		t.Errorf("queryInt() = %d, 원하는 값 = 42 (기본값)", got)
	}
}

func TestQueryFloat_DefaultValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	got := queryFloat(req, "val", 3.14)
	if got != 3.14 {
		t.Errorf("queryFloat() = %f, 원하는 값 = 3.14", got)
	}
}

func TestQueryFloat_ParseValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?val=2.71", nil)
	got := queryFloat(req, "val", 3.14)
	if got != 2.71 {
		t.Errorf("queryFloat() = %f, 원하는 값 = 2.71", got)
	}
}

func TestQueryBool_True(t *testing.T) {
	req := httptest.NewRequest("GET", "/test?flag=true", nil)
	if !queryBool(req, "flag") {
		t.Error("queryBool() = false, 원하는 값 = true")
	}
}

func TestQueryBool_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	if queryBool(req, "flag") {
		t.Error("queryBool(missing) = true, 원하는 값 = false")
	}
}

// ---------------------------------------------------------------------------
// SetupRoutes 통합 테스트
// ---------------------------------------------------------------------------

func TestSetupRoutes_HealthCheck(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)
}

func TestSetupRoutes_AuthLogin_ViaMux(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	body := `{"email":"a@b.com","password":"pw"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusServiceUnavailable)
}

func TestSetupRoutes_UnknownRoute_404(t *testing.T) {
	h := newNilHandler()
	mux := h.SetupRoutes()

	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("상태 코드 = %d, 원하는 값 = 404", w.Code)
	}
}
