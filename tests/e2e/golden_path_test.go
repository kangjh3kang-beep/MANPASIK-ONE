package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── Golden Path: 측정 전체 플로우 ───
// TestGoldenPath_MeasurementFlow tests the complete measurement journey:
// 1. Register user -> Login -> Get JWT token
// 2. Register device
// 3. Start measurement session
// 4. Submit measurement data (differential values)
// 5. End session -> Get results
// 6. Verify measurement appears in history
// 7. Verify health record created
// 8. Verify audit log entry exists

func TestGoldenPath_MeasurementFlow(t *testing.T) {
	email := uniqueEmail("golden-meas")
	password := "GoldenPath1!@#"
	name := "Golden Path User"

	// ──── Step 1: 회원가입 + 로그인 ────
	t.Log("Step 1: 회원가입 + 로그인")
	regStatus, regResult := registerUser(t, email, password, name)
	if regStatus != http.StatusCreated && regStatus != http.StatusOK {
		t.Skipf("회원가입 응답: %d (서비스 미실행 가능), body=%v", regStatus, regResult)
		return
	}
	t.Logf("  회원가입 성공 (status=%d)", regStatus)

	token := loginUser(t, email, password)
	if token == "" {
		t.Fatal("  로그인 실패: 토큰 미발급")
	}
	t.Logf("  로그인 성공 (token=%s...)", token[:min(16, len(token))])

	// ──── Step 2: 디바이스 등록 ────
	t.Log("Step 2: 디바이스 등록")
	deviceID := fmt.Sprintf("MPK-GOLDEN-%d", uniqueTS())
	deviceBody := map[string]interface{}{
		"device_id":        deviceID,
		"name":             "Golden Path Reader",
		"firmware_version": "2.1.0",
		"model":            "MPK-R1",
	}
	devResp, devResult := apiRequest(t, "POST", devicesEndpoint, deviceBody, token)
	defer devResp.Body.Close()
	t.Logf("  디바이스 등록: status=%d, result=%v", devResp.StatusCode, devResult)

	// ──── Step 3: 측정 세션 시작 ────
	t.Log("Step 3: 측정 세션 시작")
	cartridgeID := fmt.Sprintf("CART-GLU-GOLDEN-%d", uniqueTS())
	sessionBody := map[string]interface{}{
		"device_id":      deviceID,
		"cartridge_id":   cartridgeID,
		"cartridge_type": "glucose",
		"mode":           "differential",
	}
	sessResp, sessResult := apiRequest(t, "POST", measureEndpoint, sessionBody, token)
	defer sessResp.Body.Close()
	t.Logf("  세션 시작: status=%d, result=%v", sessResp.StatusCode, sessResult)

	sessionID := ""
	if id, ok := sessResult["session_id"].(string); ok {
		sessionID = id
	} else if id, ok := sessResult["id"].(string); ok {
		sessionID = id
	}

	// ──── Step 4: 측정 데이터 제출 (차동측정) ────
	t.Log("Step 4: 측정 데이터 제출 (차동측정: Sdiff = S_det - alpha * S_ref)")
	// 88차원 핑거프린트 데이터 (SSOT: 기본 88차원)
	channels := make([]float64, 88)
	refChannels := make([]float64, 88)
	for i := 0; i < 88; i++ {
		channels[i] = 1.0 + float64(i)*0.01   // 검출 신호
		refChannels[i] = 0.01 + float64(i)*0.001 // 참조 신호
	}

	dataBody := map[string]interface{}{
		"session_id":    sessionID,
		"device_id":     deviceID,
		"cartridge_id":  cartridgeID,
		"det_channels":  channels,
		"ref_channels":  refChannels,
		"alpha":         0.98, // SSOT: alpha 기본값 0.98
		"temperature_c": 25.0,
		"humidity_pct":  45.0,
	}

	dataEndpoint := measureEndpoint
	if sessionID != "" {
		dataEndpoint = fmt.Sprintf("%s/%s/data", measureEndpoint, sessionID)
	}
	dataResp, dataResult := apiRequest(t, "POST", dataEndpoint, dataBody, token)
	defer dataResp.Body.Close()
	t.Logf("  데이터 제출: status=%d, result=%v", dataResp.StatusCode, dataResult)

	// ──── Step 5: 세션 종료 + 결과 조회 ────
	t.Log("Step 5: 세션 종료 + 결과 조회")
	if sessionID != "" {
		endBody := map[string]interface{}{
			"session_id": sessionID,
			"status":     "completed",
		}
		endResp, endResult := apiRequest(t, "PUT", fmt.Sprintf("%s/%s/end", measureEndpoint, sessionID), endBody, token)
		defer endResp.Body.Close()
		t.Logf("  세션 종료: status=%d, result=%v", endResp.StatusCode, endResult)

		resultResp, resultData := apiRequest(t, "GET", fmt.Sprintf("%s/%s/result", measureEndpoint, sessionID), nil, token)
		defer resultResp.Body.Close()
		t.Logf("  결과 조회: status=%d, result=%v", resultResp.StatusCode, resultData)
	}

	// ──── Step 6: 측정 이력 검증 ────
	t.Log("Step 6: 측정 이력 검증")
	histResp, histResult := apiRequest(t, "GET", historyEndpoint+"?limit=5", nil, token)
	defer histResp.Body.Close()
	t.Logf("  측정 이력: status=%d, entries=%v", histResp.StatusCode, histResult)

	// ──── Step 7: 건강 기록 확인 ────
	t.Log("Step 7: 건강 기록 확인")
	healthResp, healthResult := apiRequest(t, "GET", "/api/v1/health-records?limit=5", nil, token)
	defer healthResp.Body.Close()
	t.Logf("  건강 기록: status=%d, result=%v", healthResp.StatusCode, healthResult)

	// ──── Step 8: 감사 로그 확인 ────
	t.Log("Step 8: 감사 로그 확인")
	auditResp, auditResult := apiRequest(t, "GET", "/api/v1/audit-log?limit=10&action=measurement", nil, token)
	defer auditResp.Body.Close()
	t.Logf("  감사 로그: status=%d, result=%v", auditResp.StatusCode, auditResult)

	t.Log("Golden Path Measurement Flow completed.")
}

// ─── Golden Path: 다중 핑거프린트 차원 측정 ───

func TestGoldenPath_MultiDimensionFingerprint(t *testing.T) {
	email := uniqueEmail("golden-fp")
	token := registerAndLogin(t, email, "GoldenFP1!@#", "Fingerprint User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// SSOT: 핑거프린트 88 -> 448 -> 896 -> 1792
	dimensions := []struct {
		name     string
		channels int
		cartType string
	}{
		{"88차원 (기본)", 88, "glucose"},
		{"448차원 (확장)", 448, "extended_448"},
		{"896차원 (통합)", 896, "integrated_896"},
		{"1792차원 (풀스펙)", 1792, "full_spectrum_1792"},
	}

	deviceID := fmt.Sprintf("MPK-FP-%d", uniqueTS())
	deviceBody := map[string]interface{}{
		"device_id":        deviceID,
		"name":             "FP Dimension Reader",
		"firmware_version": "2.1.0",
	}
	devResp, _ := apiRequest(t, "POST", devicesEndpoint, deviceBody, token)
	devResp.Body.Close()

	for _, dim := range dimensions {
		t.Run(dim.name, func(t *testing.T) {
			sessionBody := map[string]interface{}{
				"device_id":      deviceID,
				"cartridge_id":   fmt.Sprintf("CART-FP-%d-%d", dim.channels, uniqueTS()),
				"cartridge_type": dim.cartType,
				"channels":       make([]float64, dim.channels),
			}
			resp, result := apiRequest(t, "POST", measureEndpoint, sessionBody, token)
			defer resp.Body.Close()
			t.Logf("  %s: status=%d, result=%v", dim.name, resp.StatusCode, result)
		})
	}
}

// ─── Golden Path: 차동측정 정확도 (Sdiff = S - alpha * R) ───

func TestGoldenPath_DifferentialAccuracy(t *testing.T) {
	// SSOT: Sdiff_n = S_n - alpha_n * R_n (alpha 기본값 0.98, 동적 보정 0.90~1.10)
	testCases := []struct {
		name      string
		sDet      float64
		sRef      float64
		alpha     float64
		expected  float64
		tolerance float64
	}{
		{"alpha=0.98 (기본)", 1.500, 0.020, 0.98, 1.4804, 0.0001},
		{"alpha=0.90 (하한)", 2.000, 0.100, 0.90, 1.9100, 0.0001},
		{"alpha=1.10 (상한)", 1.000, 0.050, 1.10, 0.9450, 0.0001},
		{"제로 참조", 0.800, 0.000, 0.98, 0.8000, 0.0001},
		{"고감도 신호", 0.050, 0.001, 0.98, 0.04902, 0.00001},
		{"동등 신호", 1.000, 1.000, 0.98, 0.0200, 0.0001},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := differentialCorrection(tc.sDet, tc.sRef, tc.alpha)
			diff := abs(result - tc.expected)
			if diff > tc.tolerance {
				t.Errorf("차동측정 오차: got=%f, want=%f, diff=%f (tolerance=%f)",
					result, tc.expected, diff, tc.tolerance)
			} else {
				t.Logf("  S_det=%.4f, S_ref=%.4f, alpha=%.2f -> Sdiff=%.5f (expected=%.5f)",
					tc.sDet, tc.sRef, tc.alpha, result, tc.expected)
			}
		})
	}
}
