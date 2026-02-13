//go:build integration

package e2e

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestGatewayREST_AllEndpoints verifies that all major REST gateway
// endpoints return well-formed JSON responses.
func TestGatewayREST_AllEndpoints(t *testing.T) {
	gatewayURL := "http://" + GatewayAddr()
	client := &http.Client{Timeout: 5 * time.Second}

	// ── Health ──────────────────────────────────────────────────────
	t.Run("Health", func(t *testing.T) {
		resp, err := client.Get(gatewayURL + "/health")
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("health: got %d, want 200", resp.StatusCode)
		}
	})

	// ── Auth endpoints ─────────────────────────────────────────────
	t.Run("Auth_Register", func(t *testing.T) {
		body := `{"email":"test@e2e.com","password":"test1234","display_name":"E2E User"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/auth/register",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	t.Run("Auth_Login", func(t *testing.T) {
		body := `{"email":"test@e2e.com","password":"test1234"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/auth/login",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Measurement endpoints ──────────────────────────────────────
	t.Run("Measurement_StartSession", func(t *testing.T) {
		body := `{"device_id":"dev-001","cartridge_id":"cart-001"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/measurements/sessions",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Coaching endpoints ─────────────────────────────────────────
	t.Run("Coaching_SetGoal", func(t *testing.T) {
		body := `{"user_id":"user-001","category":"WEIGHT","target_value":70.0}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/coaching/goals",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── AI Inference ───────────────────────────────────────────────
	t.Run("AI_Analyze", func(t *testing.T) {
		body := `{"session_id":"sess-001","user_id":"user-001"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/ai/analyze",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Subscription endpoints ─────────────────────────────────────
	t.Run("Subscription_Create", func(t *testing.T) {
		body := `{"user_id":"user-001","plan":"premium"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/subscriptions",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Shop endpoints ─────────────────────────────────────────────
	t.Run("Shop_ListProducts", func(t *testing.T) {
		resp, err := client.Get(gatewayURL + "/api/v1/shop/products")
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Community endpoints ────────────────────────────────────────
	t.Run("Community_ListPosts", func(t *testing.T) {
		resp, err := client.Get(gatewayURL + "/api/v1/community/posts")
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Reservation endpoints ──────────────────────────────────────
	t.Run("Reservation_Create", func(t *testing.T) {
		body := `{"user_id":"user-001","facility_name":"서울대병원","date":"2026-03-01"}`
		resp, err := client.Post(
			gatewayURL+"/api/v1/reservations",
			"application/json",
			strings.NewReader(body),
		)
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})

	// ── Family endpoints ───────────────────────────────────────────
	t.Run("Family_ListMembers", func(t *testing.T) {
		resp, err := client.Get(gatewayURL + "/api/v1/family/members?user_id=user-001")
		if err != nil {
			t.Skip("Gateway not running:", err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Errorf("expected JSON response: %v", err)
		}
	})
}
