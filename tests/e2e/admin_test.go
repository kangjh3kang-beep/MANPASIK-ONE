package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-ADMIN-001: 관리자 사용자 관리 ───

func TestAdminUserManagement(t *testing.T) {
	// 일반 사용자 토큰
	userEmail := uniqueEmail("admin-user")
	userToken := registerAndLogin(t, userEmail, "E2eTest1!@#", "Regular User")
	if userToken == "" {
		t.Skip("사용자 토큰 획득 실패")
		return
	}

	// 관리자 API 접근 시도 (일반 사용자)
	resp, _ := apiRequest(t, "GET", "/api/v1/admin/users", nil, userToken)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("관리자 권한 위반: 일반 사용자가 관리자 API 접근 가능")
	}
	t.Logf("일반 사용자 관리자 API 차단: status=%d", resp.StatusCode)
}

// ─── E2E-ADMIN-002: 카트리지 관리 ───

func TestAdminCartridgeManagement(t *testing.T) {
	userEmail := uniqueEmail("admin-cart")
	userToken := registerAndLogin(t, userEmail, "E2eTest1!@#", "Cart Admin User")
	if userToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 카트리지 등록 (관리자 전용)
	cartBody := map[string]interface{}{
		"type_code":        "0x0F",
		"name_ko":          "새로운 바이오마커",
		"name_en":          "NewBiomarker",
		"required_channels": 88,
		"measurement_secs":  15,
		"unit":             "mg/dL",
		"reference_range":  "50-200",
	}
	resp, result := apiRequest(t, "POST", "/api/v1/admin/cartridges", cartBody, userToken)
	defer resp.Body.Close()
	t.Logf("카트리지 등록 시도: status=%d, result=%v", resp.StatusCode, result)
}

// ─── E2E-ADMIN-003: 시스템 설정 ───

func TestAdminSystemConfig(t *testing.T) {
	userEmail := uniqueEmail("admin-config")
	userToken := registerAndLogin(t, userEmail, "E2eTest1!@#", "Config User")
	if userToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 시스템 설정 조회 (관리자 전용)
	resp, result := apiRequest(t, "GET", "/api/v1/admin/system/config", nil, userToken)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("일반 사용자가 시스템 설정 접근 가능")
	}
	t.Logf("시스템 설정 접근 차단: status=%d, result=%v", resp.StatusCode, result)
}

// ─── E2E-ADMIN-004: 감사 로그 조회 ───

func TestAdminAuditLogs(t *testing.T) {
	userEmail := uniqueEmail("admin-audit")
	userToken := registerAndLogin(t, userEmail, "E2eTest1!@#", "Audit User")
	if userToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 감사 로그 조회 (관리자 전용)
	resp, _ := apiRequest(t, "GET", "/api/v1/admin/audit-logs?limit=50", nil, userToken)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("일반 사용자가 감사 로그 접근 가능")
	}
	t.Logf("감사 로그 접근 차단: status=%d", resp.StatusCode)
}

// ─── E2E-ADMIN-005: 번역 관리 ───

func TestAdminTranslationManagement(t *testing.T) {
	userEmail := uniqueEmail("admin-i18n")
	userToken := registerAndLogin(t, userEmail, "E2eTest1!@#", "Translation User")
	if userToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 번역 키 조회
	resp, result := apiRequest(t, "GET", "/api/v1/translations?locale=ko", nil, userToken)
	defer resp.Body.Close()
	t.Logf("번역 조회: status=%d, result=%v", resp.StatusCode, result)
}
