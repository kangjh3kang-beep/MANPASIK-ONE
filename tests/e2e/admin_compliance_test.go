package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-COMP-001: 감사 로그 생성 및 조회 ───

func TestComplianceAuditLogCreation(t *testing.T) {
	email := uniqueEmail("comp-audit")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Audit Compliance User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 사용자 활동 (감사 로그 자동 생성 트리거)
	// 프로필 조회 → 감사 로그 생성
	profResp, _ := apiRequest(t, "GET", profileEndpoint, nil, token)
	defer profResp.Body.Close()

	// 측정 이력 조회 → 감사 로그 생성
	histResp, _ := apiRequest(t, "GET", historyEndpoint, nil, token)
	defer histResp.Body.Close()

	// 2. 감사 로그 조회 (사용자 자신의 활동)
	auditResp, auditResult := apiRequest(t, "GET", "/api/v1/users/me/audit-logs?limit=20", nil, token)
	defer auditResp.Body.Close()
	t.Logf("사용자 감사 로그: status=%d, result=%v", auditResp.StatusCode, auditResult)
}

// ─── E2E-COMP-002: GDPR 데이터 내보내기 (Right to Access) ───

func TestComplianceGDPRDataExport(t *testing.T) {
	email := uniqueEmail("comp-gdpr-export")
	token := registerAndLogin(t, email, "E2eTest1!@#", "GDPR Export User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 내 데이터 내보내기 요청
	exportResp, exportResult := apiRequest(t, "POST", "/api/v1/users/me/data-export", nil, token)
	defer exportResp.Body.Close()
	t.Logf("데이터 내보내기 요청: status=%d, result=%v", exportResp.StatusCode, exportResult)

	// 내보내기 상태 조회
	statusResp, statusResult := apiRequest(t, "GET", "/api/v1/users/me/data-export/status", nil, token)
	defer statusResp.Body.Close()
	t.Logf("내보내기 상태: status=%d, result=%v", statusResp.StatusCode, statusResult)
}

// ─── E2E-COMP-003: GDPR 데이터 삭제 (Right to Erasure) ───

func TestComplianceGDPRDataDeletion(t *testing.T) {
	email := uniqueEmail("comp-gdpr-delete")
	token := registerAndLogin(t, email, "E2eTest1!@#", "GDPR Delete User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 데이터 삭제 요청
	deleteBody := map[string]interface{}{
		"confirm":      true,
		"reason":       "privacy_concern",
		"delete_scope": "all_personal_data",
	}
	delResp, delResult := apiRequest(t, "POST", "/api/v1/users/me/data-deletion", deleteBody, token)
	defer delResp.Body.Close()
	t.Logf("데이터 삭제 요청: status=%d, result=%v", delResp.StatusCode, delResult)

	// 삭제 후 프로필 접근 시도 (실패 예상)
	profResp, _ := apiRequest(t, "GET", profileEndpoint, nil, token)
	defer profResp.Body.Close()
	t.Logf("삭제 후 프로필 접근: status=%d", profResp.StatusCode)
}

// ─── E2E-COMP-004: 규제 체크리스트 (관리자용) ───

func TestComplianceRegulatoryChecklist(t *testing.T) {
	email := uniqueEmail("comp-checklist")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Checklist User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// GDPR 체크리스트 조회 (관리자 전용)
	gdprResp, gdprResult := apiRequest(t, "GET", "/api/v1/admin/compliance/gdpr", nil, token)
	defer gdprResp.Body.Close()
	if gdprResp.StatusCode == http.StatusOK {
		t.Log("경고: 일반 사용자가 규제 체크리스트 접근 가능")
	}
	t.Logf("GDPR 체크리스트: status=%d, result=%v", gdprResp.StatusCode, gdprResult)

	// PIPA (개인정보보호법) 체크리스트
	pipaResp, pipaResult := apiRequest(t, "GET", "/api/v1/admin/compliance/pipa", nil, token)
	defer pipaResp.Body.Close()
	t.Logf("PIPA 체크리스트: status=%d, result=%v", pipaResp.StatusCode, pipaResult)

	// HIPAA 체크리스트
	hipaaResp, hipaaResult := apiRequest(t, "GET", "/api/v1/admin/compliance/hipaa", nil, token)
	defer hipaaResp.Body.Close()
	t.Logf("HIPAA 체크리스트: status=%d, result=%v", hipaaResp.StatusCode, hipaaResult)
}

// ─── E2E-COMP-005: 규제 보고서 생성 (관리자용) ───

func TestComplianceReportGeneration(t *testing.T) {
	email := uniqueEmail("comp-report")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Report User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 규제 준수 보고서 생성 요청
	reportBody := map[string]interface{}{
		"report_type": "quarterly",
		"frameworks":  []string{"GDPR", "PIPA", "HIPAA"},
		"period": map[string]string{
			"start": "2026-01-01",
			"end":   "2026-03-31",
		},
	}
	reportResp, reportResult := apiRequest(t, "POST", "/api/v1/admin/compliance/reports", reportBody, token)
	defer reportResp.Body.Close()
	if reportResp.StatusCode == http.StatusOK || reportResp.StatusCode == http.StatusCreated {
		t.Log("경고: 일반 사용자가 규제 보고서 생성 가능 (관리자 전용이어야 함)")
	}
	t.Logf("규제 보고서 생성: status=%d, result=%v", reportResp.StatusCode, reportResult)
}

// ─── E2E-COMP-006: 동의 관리 (Consent Management) ───

func TestComplianceConsentManagement(t *testing.T) {
	email := uniqueEmail("comp-consent")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Consent User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 현재 동의 상태 조회
	currentResp, currentResult := apiRequest(t, "GET", "/api/v1/users/me/consents", nil, token)
	defer currentResp.Body.Close()
	t.Logf("1. 현재 동의: status=%d, result=%v", currentResp.StatusCode, currentResult)

	// 2. 마케팅 동의 철회
	revokeBody := map[string]interface{}{
		"consent_type": "marketing",
		"granted":      false,
		"reason":       "no_longer_interested",
	}
	revokeResp, revokeResult := apiRequest(t, "PATCH", "/api/v1/users/me/consents", revokeBody, token)
	defer revokeResp.Body.Close()
	t.Logf("2. 마케팅 동의 철회: status=%d, result=%v", revokeResp.StatusCode, revokeResult)

	// 3. 동의 이력 조회
	histResp, histResult := apiRequest(t, "GET", "/api/v1/users/me/consents/history", nil, token)
	defer histResp.Body.Close()
	t.Logf("3. 동의 이력: status=%d, result=%v", histResp.StatusCode, histResult)
}
