package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// ─── E2E-OFF-001: 오프라인 측정 → 재연결 → 동기화 ───

func TestOfflineMeasurementSync(t *testing.T) {
	email := uniqueEmail("off-meas")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Offline Meas User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 오프라인 측정 데이터 업로드 (일괄 동기화 시뮬레이션)
	syncBody := map[string]interface{}{
		"measurements": []map[string]interface{}{
			{
				"local_id":       "LOCAL-001",
				"device_id":      "MPK-OFF-001",
				"cartridge_id":   "CART-GLU-OFF-001",
				"channels":       make([]float64, 88),
				"result_value":   95.5,
				"result_unit":    "mg/dL",
				"measured_at":    time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
				"offline_stored": true,
			},
			{
				"local_id":       "LOCAL-002",
				"device_id":      "MPK-OFF-001",
				"cartridge_id":   "CART-GLU-OFF-002",
				"channels":       make([]float64, 88),
				"result_value":   110.2,
				"result_unit":    "mg/dL",
				"measured_at":    time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
				"offline_stored": true,
			},
		},
	}
	syncResp, syncResult := apiRequest(t, "POST", "/api/v1/sync/measurements", syncBody, token)
	defer syncResp.Body.Close()
	t.Logf("1. 오프라인 측정 동기화: status=%d, result=%v", syncResp.StatusCode, syncResult)

	// 2. 동기화 후 이력 조회
	histResp, histResult := apiRequest(t, "GET", historyEndpoint+"?limit=10", nil, token)
	defer histResp.Body.Close()
	t.Logf("2. 동기화 후 이력: status=%d, result=%v", histResp.StatusCode, histResult)
}

// ─── E2E-OFF-002: 동기화 상태 조회 ───

func TestOfflineSyncStatus(t *testing.T) {
	email := uniqueEmail("off-status")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Sync Status User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 동기화 상태 조회
	statusResp, statusResult := apiRequest(t, "GET", "/api/v1/sync/status", nil, token)
	defer statusResp.Body.Close()
	t.Logf("동기화 상태: status=%d, result=%v", statusResp.StatusCode, statusResult)
}

// ─── E2E-OFF-003: 충돌 감지 및 해결 ───

func TestOfflineConflictResolution(t *testing.T) {
	email := uniqueEmail("off-conflict")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Conflict User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 충돌 데이터 업로드 시뮬레이션 (동일 ID, 다른 타임스탬프)
	conflictBody := map[string]interface{}{
		"measurements": []map[string]interface{}{
			{
				"local_id":     "CONFLICT-001",
				"device_id":    "MPK-OFF-001",
				"result_value": 88.0,
				"measured_at":  time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			},
		},
		"resolve_strategy": "client_wins", // client_wins | server_wins | manual
	}
	conflictResp, conflictResult := apiRequest(t, "POST", "/api/v1/sync/measurements", conflictBody, token)
	defer conflictResp.Body.Close()
	t.Logf("1. 충돌 동기화: status=%d, result=%v", conflictResp.StatusCode, conflictResult)

	// 2. 충돌 목록 조회
	listResp, listResult := apiRequest(t, "GET", "/api/v1/sync/conflicts", nil, token)
	defer listResp.Body.Close()
	t.Logf("2. 충돌 목록: status=%d, result=%v", listResp.StatusCode, listResult)

	// 3. 수동 충돌 해결
	resolveBody := map[string]interface{}{
		"conflict_id": "CONFLICT-001",
		"resolution":  "keep_local",
	}
	resolveResp, resolveResult := apiRequest(t, "POST", "/api/v1/sync/conflicts/resolve", resolveBody, token)
	defer resolveResp.Body.Close()
	t.Logf("3. 충돌 해결: status=%d, result=%v", resolveResp.StatusCode, resolveResult)
}

// ─── E2E-OFF-004: 대량 오프라인 데이터 일괄 동기화 ───

func TestOfflineBulkSync(t *testing.T) {
	email := uniqueEmail("off-bulk")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Bulk Sync User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 20건 오프라인 측정 일괄 동기화
	measurements := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		measurements[i] = map[string]interface{}{
			"local_id":     fmt.Sprintf("BULK-%03d", i+1),
			"device_id":    "MPK-BULK-001",
			"result_value": 80.0 + float64(i)*2,
			"result_unit":  "mg/dL",
			"measured_at":  time.Now().Add(-time.Duration(20-i) * time.Hour).Format(time.RFC3339),
		}
	}

	bulkBody := map[string]interface{}{
		"measurements": measurements,
	}
	bulkResp, bulkResult := apiRequest(t, "POST", "/api/v1/sync/measurements", bulkBody, token)
	defer bulkResp.Body.Close()
	t.Logf("대량 동기화 (20건): status=%d, result=%v", bulkResp.StatusCode, bulkResult)
}

// ─── E2E-OFF-005: 오프라인 설정 변경 동기화 ───

func TestOfflineSettingsSync(t *testing.T) {
	email := uniqueEmail("off-settings")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Settings Sync User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 오프라인에서 변경된 설정 동기화
	settingsBody := map[string]interface{}{
		"settings": map[string]interface{}{
			"notification_enabled": false,
			"theme":                "dark",
			"locale":               "ko",
			"measurement_reminder": "08:00",
		},
		"changed_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
	}
	settingsResp, settingsResult := apiRequest(t, "POST", "/api/v1/sync/settings", settingsBody, token)
	defer settingsResp.Body.Close()
	t.Logf("설정 동기화: status=%d, result=%v", settingsResp.StatusCode, settingsResult)
}
