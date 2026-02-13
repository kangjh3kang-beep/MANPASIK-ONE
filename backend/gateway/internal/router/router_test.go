//go:build integration

package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestRouter creates a gateway router with dummy backend addresses.
// The gRPC backends are not expected to be running; these tests only verify
// that the HTTP layer (routing, JSON, status codes) works correctly.
func newTestRouter() *Router {
	return NewRouter(Config{
		AuthAddr:         "localhost:50051",
		MeasurementAddr:  "localhost:50054",
		UserAddr:         "localhost:50052",
		DeviceAddr:       "localhost:50053",
		ReservationAddr:  "localhost:50055",
		PrescriptionAddr: "localhost:50062",
		SubscriptionAddr: "localhost:50055",
		ShopAddr:         "localhost:50056",
		PaymentAddr:      "localhost:50057",
		HealthRecordAddr: "localhost:50064",
		NotificationAddr: "localhost:50068",
		CommunityAddr:    "localhost:50065",
		AdminAddr:        "localhost:50067",
	})
}

// --------------------------------------------------------------------------
// Health / Version endpoints  (no gRPC backend required)
// --------------------------------------------------------------------------

func TestGatewayHealthEndpoint(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}

	if body["status"] != "healthy" {
		t.Errorf("status = %v, want healthy", body["status"])
	}
	if body["service"] != "gateway" {
		t.Errorf("service = %v, want gateway", body["service"])
	}
	if _, ok := body["time"]; !ok {
		t.Error("response missing 'time' field")
	}

	t.Logf("✅ /health → %v", body)
}

func TestGatewayHealthLiveEndpoint(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health/live status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}

	if body["alive"] != true {
		t.Errorf("alive = %v, want true", body["alive"])
	}
}

func TestGatewayHealthReadyEndpoint(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /health/ready status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}

	if body["ready"] != true {
		t.Errorf("ready = %v, want true", body["ready"])
	}
}

func TestGatewayVersionEndpoint(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/version status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("JSON decode error: %v", err)
	}

	if body["service"] != "gateway" {
		t.Errorf("service = %v, want gateway", body["service"])
	}
	if body["version"] != "1.0.0" {
		t.Errorf("version = %v, want 1.0.0", body["version"])
	}
	if body["api_version"] != "v1" {
		t.Errorf("api_version = %v, want v1", body["api_version"])
	}
	if body["services"] != float64(13) {
		t.Errorf("services = %v, want 13", body["services"])
	}

	t.Logf("✅ /api/v1/version → %v", body)
}

// --------------------------------------------------------------------------
// CORS preflight
// --------------------------------------------------------------------------

func TestGatewayCORSPreflight(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS status = %d, want %d", w.Code, http.StatusNoContent)
	}

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("CORS Allow-Origin = %q, want *", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("CORS Allow-Methods header is empty")
	}
	if got := w.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Error("CORS Allow-Headers header is empty")
	}

	t.Log("✅ OPTIONS → CORS headers present")
}

// --------------------------------------------------------------------------
// Route existence tests — verify endpoints return the right error
// (503 Service Unavailable because the gRPC backends are not running,
//  NOT 404 which would mean the route is not registered)
// --------------------------------------------------------------------------

func TestGatewayRouteRegistration(t *testing.T) {
	r := newTestRouter()

	tests := []struct {
		method string
		path   string
		body   interface{} // nil for GET/DELETE, struct for POST/PUT
	}{
		// Auth
		{"POST", "/api/v1/auth/register", map[string]string{"email": "a@b.c", "password": "p", "display_name": "n"}},
		{"POST", "/api/v1/auth/login", map[string]string{"email": "a@b.c", "password": "p"}},
		{"POST", "/api/v1/auth/refresh", map[string]string{"refresh_token": "tok"}},
		{"POST", "/api/v1/auth/logout", map[string]string{"user_id": "uid"}},

		// User
		{"GET", "/api/v1/users/user123/profile", nil},
		{"PUT", "/api/v1/users/user123/profile", map[string]string{"display_name": "test"}},

		// Measurement
		{"POST", "/api/v1/measurements/sessions", map[string]string{"device_id": "d", "user_id": "u"}},
		{"POST", "/api/v1/measurements/sessions/sess1/end", nil},
		{"GET", "/api/v1/measurements/history?user_id=u1", nil},

		// Device
		{"POST", "/api/v1/devices", map[string]string{"device_id": "d", "user_id": "u"}},
		{"GET", "/api/v1/devices?user_id=u1", nil},

		// Reservation
		{"GET", "/api/v1/facilities?query=test", nil},
		{"GET", "/api/v1/facilities/fac1", nil},
		{"POST", "/api/v1/reservations", map[string]string{"user_id": "u", "facility_id": "f"}},
		{"GET", "/api/v1/reservations?user_id=u1", nil},
		{"GET", "/api/v1/reservations/res1", nil},

		// Prescription
		{"POST", "/api/v1/prescriptions/rx1/pharmacy", map[string]string{"pharmacy_id": "ph1"}},
		{"POST", "/api/v1/prescriptions/rx1/send", nil},
		{"GET", "/api/v1/prescriptions/token/tok123", nil},

		// Subscription
		{"GET", "/api/v1/subscriptions/plans", nil},
		{"GET", "/api/v1/subscriptions/user123", nil},
		{"POST", "/api/v1/subscriptions", map[string]string{"user_id": "u1"}},
		{"DELETE", "/api/v1/subscriptions/sub1", nil},

		// Shop
		{"GET", "/api/v1/products", nil},
		{"GET", "/api/v1/products/prod1", nil},
		{"POST", "/api/v1/cart", map[string]string{"user_id": "u", "product_id": "p"}},
		{"GET", "/api/v1/cart/user1", nil},
		{"POST", "/api/v1/orders", map[string]string{"user_id": "u"}},
		{"GET", "/api/v1/orders?user_id=u1", nil},

		// Payment
		{"POST", "/api/v1/payments", map[string]string{"user_id": "u"}},
		{"POST", "/api/v1/payments/pay1/confirm", map[string]string{"pg_transaction_id": "tx"}},
		{"GET", "/api/v1/payments/pay1", nil},

		// Health Record
		{"POST", "/api/v1/health-records", map[string]string{"user_id": "u"}},
		{"GET", "/api/v1/health-records?user_id=u1", nil},
		{"GET", "/api/v1/health-records/rec1", nil},
		{"POST", "/api/v1/health-records/export/fhir", map[string]string{"user_id": "u"}},

		// Notification
		{"GET", "/api/v1/notifications?user_id=u1", nil},
		{"POST", "/api/v1/notifications/notif1/read", nil},
		{"GET", "/api/v1/notifications/unread-count?user_id=u1", nil},

		// Community
		{"GET", "/api/v1/posts", nil},
		{"POST", "/api/v1/posts", map[string]string{"author_id": "a", "title": "t", "content": "c"}},
		{"GET", "/api/v1/posts/post1", nil},
		{"POST", "/api/v1/posts/post1/like", map[string]string{"user_id": "u"}},

		// Admin
		{"GET", "/api/v1/admin/stats", nil},
		{"GET", "/api/v1/admin/users", nil},
		{"GET", "/api/v1/admin/audit-log", nil},
	}

	for _, tc := range tests {
		name := tc.method + " " + tc.path
		t.Run(name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				b, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewReader(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// The route must be registered. Our handlers always return JSON
			// (via writeJSON / writeError). The mux's built-in 404 handler
			// returns plain text "404 page not found". So a valid JSON response
			// means the handler ran, i.e. the route is registered.
			// Note: GET handlers for single resources return 404 with JSON when
			// the gRPC backend is unreachable — that's still a registered route.
			var body map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
				t.Errorf("%s %s → %d with non-JSON response (route not registered): %v",
					tc.method, tc.path, w.Code, err)
			}

			// 405 Method Not Allowed means no handler for this method
			if w.Code == http.StatusMethodNotAllowed {
				t.Errorf("%s %s → %d (method not allowed)", tc.method, tc.path, w.Code)
			}

			t.Logf("✅ %s %s → %d", tc.method, tc.path, w.Code)
		})
	}
}

// --------------------------------------------------------------------------
// Content-Type header
// --------------------------------------------------------------------------

func TestGatewayJSONContentType(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}
