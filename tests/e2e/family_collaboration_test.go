package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-FCOL-001: 가족 그룹 생성 → 초대 → 수락 → 데이터 공유 전체 플로우 ───

func TestFamilyCollaborationFullFlow(t *testing.T) {
	parentEmail := uniqueEmail("fcol-parent")
	childEmail := uniqueEmail("fcol-child")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "부모 (협업)")
	childToken := registerAndLogin(t, childEmail, "E2eTest1!@#", "자녀 (협업)")
	if parentToken == "" || childToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 1. 가족 그룹 생성
	familyBody := map[string]interface{}{
		"name":        "만파식 가족",
		"description": "건강 관리 협업 가족 그룹",
	}
	createResp, createResult := apiRequest(t, "POST", "/api/v1/family/groups", familyBody, parentToken)
	defer createResp.Body.Close()
	t.Logf("1. 가족 그룹 생성: status=%d, result=%v", createResp.StatusCode, createResult)

	// 2. 초대 링크 생성
	inviteBody := map[string]interface{}{
		"email": childEmail,
		"role":  "child",
		"mode":  "monitor", // 모니터링 모드
	}
	inviteResp, inviteResult := apiRequest(t, "POST", "/api/v1/family/groups/1/invite", inviteBody, parentToken)
	defer inviteResp.Body.Close()
	t.Logf("2. 초대 생성: status=%d, result=%v", inviteResp.StatusCode, inviteResult)

	// 3. 자녀가 초대 수락
	acceptBody := map[string]interface{}{
		"invite_code": "INVITE-001",
	}
	acceptResp, acceptResult := apiRequest(t, "POST", "/api/v1/family/invites/accept", acceptBody, childToken)
	defer acceptResp.Body.Close()
	t.Logf("3. 초대 수락: status=%d, result=%v", acceptResp.StatusCode, acceptResult)

	// 4. 부모가 가족 대시보드 조회 (자녀 데이터 포함)
	dashResp, dashResult := apiRequest(t, "GET", "/api/v1/family/dashboard", nil, parentToken)
	defer dashResp.Body.Close()
	t.Logf("4. 가족 대시보드: status=%d, result=%v", dashResp.StatusCode, dashResult)

	// 5. 자녀 측정 데이터를 부모가 조회
	memberDataResp, memberDataResult := apiRequest(t, "GET", "/api/v1/family/members/child-1/health-data?period=7d", nil, parentToken)
	defer memberDataResp.Body.Close()
	t.Logf("5. 자녀 건강 데이터: status=%d, result=%v", memberDataResp.StatusCode, memberDataResult)
}

// ─── E2E-FCOL-002: 보호자 대시보드 7일 트렌드 ───

func TestGuardianDashboardTrend(t *testing.T) {
	guardianEmail := uniqueEmail("fcol-guardian")
	guardianToken := registerAndLogin(t, guardianEmail, "E2eTest1!@#", "보호자")
	if guardianToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 보호자 대시보드 (7일 트렌드)
	resp, result := apiRequest(t, "GET", "/api/v1/family/guardian/dashboard?period=7d", nil, guardianToken)
	defer resp.Body.Close()
	t.Logf("보호자 대시보드 7일: status=%d, result=%v", resp.StatusCode, result)

	// 30일 트렌드
	resp30, result30 := apiRequest(t, "GET", "/api/v1/family/guardian/dashboard?period=30d", nil, guardianToken)
	defer resp30.Body.Close()
	t.Logf("보호자 대시보드 30일: status=%d, result=%v", resp30.StatusCode, result30)
}

// ─── E2E-FCOL-003: 긴급 알림 전파 ───

func TestFamilyEmergencyAlert(t *testing.T) {
	parentEmail := uniqueEmail("fcol-emg-p")
	childEmail := uniqueEmail("fcol-emg-c")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "부모 (긴급)")
	_ = registerAndLogin(t, childEmail, "E2eTest1!@#", "자녀 (긴급)")
	if parentToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 긴급 알림 수동 발송
	alertBody := map[string]interface{}{
		"type":    "emergency",
		"message": "자녀 저혈당 위험 수치 감지",
		"severity": "critical",
		"member_id": "child-1",
	}
	alertResp, alertResult := apiRequest(t, "POST", "/api/v1/family/alerts", alertBody, parentToken)
	defer alertResp.Body.Close()
	t.Logf("긴급 알림 발송: status=%d, result=%v", alertResp.StatusCode, alertResult)

	// 긴급 알림 목록 조회
	listResp, listResult := apiRequest(t, "GET", "/api/v1/family/alerts?type=emergency", nil, parentToken)
	defer listResp.Body.Close()
	t.Logf("긴급 알림 목록: status=%d, result=%v", listResp.StatusCode, listResult)
}

// ─── E2E-FCOL-004: 멤버 역할/모드 변경 ───

func TestFamilyMemberRoleChange(t *testing.T) {
	parentEmail := uniqueEmail("fcol-role-p")
	parentToken := registerAndLogin(t, parentEmail, "E2eTest1!@#", "역할변경 부모")
	if parentToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 멤버 역할 변경 (child → guardian)
	roleBody := map[string]interface{}{
		"role": "guardian",
		"mode": "full_access",
	}
	roleResp, roleResult := apiRequest(t, "PATCH", "/api/v1/family/members/member-1", roleBody, parentToken)
	defer roleResp.Body.Close()
	t.Logf("역할 변경: status=%d, result=%v", roleResp.StatusCode, roleResult)

	// 멤버 목록으로 변경 확인
	listResp, listResult := apiRequest(t, "GET", "/api/v1/family/groups/1/members", nil, parentToken)
	defer listResp.Body.Close()
	t.Logf("멤버 목록: status=%d, result=%v", listResp.StatusCode, listResult)
}

// ─── E2E-FCOL-005: 비가족원 데이터 접근 차단 ───

func TestFamilyNonMemberAccessDenied(t *testing.T) {
	strangerEmail := uniqueEmail("fcol-stranger")
	strangerToken := registerAndLogin(t, strangerEmail, "E2eTest1!@#", "비가족원")
	if strangerToken == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 타인의 가족 대시보드 접근 시도
	resp, _ := apiRequest(t, "GET", "/api/v1/family/members/child-1/health-data", nil, strangerToken)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("비가족원이 가족 건강 데이터 접근 가능: 프라이버시 위반")
	}
	t.Logf("비가족원 접근 차단: status=%d", resp.StatusCode)
}
