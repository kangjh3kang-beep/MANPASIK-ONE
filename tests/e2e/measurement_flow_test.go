package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── E2E-MEAS-001: 완전한 측정 플로우 ───

func TestMeasurementFullFlow(t *testing.T) {
	email := uniqueEmail("meas-full")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Measurement User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 디바이스 등록
	deviceBody := map[string]string{
		"device_id":        "MPK-E2E-001",
		"name":             "E2E Test Reader",
		"firmware_version": "1.0.0",
	}
	devResp, devResult := apiRequest(t, "POST", devicesEndpoint, deviceBody, token)
	defer devResp.Body.Close()
	t.Logf("1. 디바이스 등록: status=%d, result=%v", devResp.StatusCode, devResult)

	// 2. 카트리지 스캔 (시뮬레이션)
	cartBody := map[string]interface{}{
		"cartridge_id":   "CART-GLU-E2E-001",
		"cartridge_type": "glucose",
		"lot_id":         "LOT20260101",
		"expiry_date":    "2027-12-31",
	}
	cartResp, cartResult := apiRequest(t, "POST", cartridgesEndpoint+"/scan", cartBody, token)
	defer cartResp.Body.Close()
	t.Logf("2. 카트리지 스캔: status=%d, result=%v", cartResp.StatusCode, cartResult)

	// 3. 측정 세션 생성
	sessionBody := map[string]interface{}{
		"device_id":    "MPK-E2E-001",
		"cartridge_id": "CART-GLU-E2E-001",
		"channels":     make([]float64, 88),
	}
	sessResp, sessResult := apiRequest(t, "POST", measureEndpoint, sessionBody, token)
	defer sessResp.Body.Close()
	t.Logf("3. 측정 세션: status=%d, result=%v", sessResp.StatusCode, sessResult)

	// 4. 측정 이력 조회
	histResp, histResult := apiRequest(t, "GET", historyEndpoint+"?limit=5", nil, token)
	defer histResp.Body.Close()
	t.Logf("4. 측정 이력: status=%d, entries=%v", histResp.StatusCode, histResult)
}

// ─── E2E-MEAS-002: 차동측정 정확도 검증 ───

func TestMeasurementDifferentialAccuracy(t *testing.T) {
	testCases := []struct {
		name      string
		sDet      float64
		sRef      float64
		alpha     float64
		expected  float64
		tolerance float64
	}{
		{"기본 혈당", 1.234, 0.012, 0.95, 1.2226, 0.001},
		{"높은 노이즈", 2.500, 0.500, 0.95, 2.025, 0.001},
		{"알파 1.0", 1.000, 0.100, 1.00, 0.900, 0.001},
		{"제로 참조", 0.800, 0.000, 0.95, 0.800, 0.001},
		{"고감도", 0.050, 0.001, 0.95, 0.04905, 0.0001},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := differentialCorrection(tc.sDet, tc.sRef, tc.alpha)
			diff := abs(result - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("차동측정 오차: got=%f, want=%f, diff=%f", result, tc.expected, diff)
			}
		})
	}
}

// ─── E2E-MEAS-003: 다중 카트리지 타입 측정 ───

func TestMeasurementMultipleCartridgeTypes(t *testing.T) {
	email := uniqueEmail("meas-multi")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Multi Cart User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	cartridgeTypes := []struct {
		name     string
		typeCode string
		channels int
	}{
		{"혈당 (88ch)", "glucose", 88},
		{"지질패널 (88ch)", "lipid_panel", 88},
		{"비표적448 (448ch)", "non_target_448", 448},
		{"비표적896 (896ch)", "non_target_896", 896},
	}

	for _, ct := range cartridgeTypes {
		t.Run(ct.name, func(t *testing.T) {
			sessionBody := map[string]interface{}{
				"device_id":      "MPK-E2E-001",
				"cartridge_id":   fmt.Sprintf("CART-%s-001", ct.typeCode),
				"cartridge_type": ct.typeCode,
				"channels":       make([]float64, ct.channels),
			}
			resp, result := apiRequest(t, "POST", measureEndpoint, sessionBody, token)
			defer resp.Body.Close()
			t.Logf("카트리지 %s: status=%d, result=%v", ct.name, resp.StatusCode, result)
		})
	}
}

// ─── E2E-MEAS-004: 측정 이력 페이지네이션 ───

func TestMeasurementHistoryPagination(t *testing.T) {
	email := uniqueEmail("meas-page")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Pagination User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 페이지 1
	resp1, result1 := apiRequest(t, "GET", historyEndpoint+"?limit=10&offset=0", nil, token)
	defer resp1.Body.Close()
	t.Logf("페이지 1: status=%d, result=%v", resp1.StatusCode, result1)

	// 페이지 2
	resp2, result2 := apiRequest(t, "GET", historyEndpoint+"?limit=10&offset=10", nil, token)
	defer resp2.Body.Close()
	t.Logf("페이지 2: status=%d, result=%v", resp2.StatusCode, result2)
}

// ─── E2E-MEAS-005: 인증 없는 측정 요청 거부 ───

func TestMeasurementWithoutAuth(t *testing.T) {
	sessionBody := map[string]interface{}{
		"device_id":    "MPK-NOAUTH-001",
		"cartridge_id": "CART-GLU-001",
		"channels":     make([]float64, 88),
	}

	resp, _ := apiRequest(t, "POST", measureEndpoint, sessionBody, "")
	defer resp.Body.Close()
	assertNotStatus(t, resp.StatusCode, http.StatusOK, "인증 없는 측정")
	assertNotStatus(t, resp.StatusCode, http.StatusCreated, "인증 없는 측정")
	t.Logf("인증 없는 측정 거부 확인: status=%d", resp.StatusCode)
}
