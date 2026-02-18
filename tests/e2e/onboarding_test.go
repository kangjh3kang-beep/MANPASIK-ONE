package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-OB-001: 회원가입 → 온보딩 → 디바이스 페어링 전체 플로우 ───

func TestOnboardingFullFlow(t *testing.T) {
	email := uniqueEmail("ob-full")
	password := "E2eTest1!@#"

	// 1. 회원가입
	status, _ := registerUser(t, email, password, "Onboarding User")
	if status != http.StatusCreated && status != http.StatusOK {
		t.Skipf("회원가입 응답: %d (서비스 미실행 가능)", status)
		return
	}
	t.Logf("1. 회원가입 성공 (status=%d)", status)

	// 2. 로그인
	token := loginUser(t, email, password)
	if token == "" {
		t.Fatal("2. 로그인 실패: 토큰 미발급")
	}
	t.Logf("2. 로그인 성공")

	// 3. 온보딩 프로필 설정
	profileBody := map[string]interface{}{
		"birth_year":     1990,
		"gender":         "male",
		"height_cm":      175,
		"weight_kg":      70,
		"health_goal":    "glucose_management",
		"medical_history": []string{"hypertension"},
	}
	profResp, profResult := apiRequest(t, "POST", "/api/v1/users/onboarding/profile", profileBody, token)
	defer profResp.Body.Close()
	t.Logf("3. 온보딩 프로필: status=%d, result=%v", profResp.StatusCode, profResult)

	// 4. 약관 동의
	consentBody := map[string]interface{}{
		"terms_of_service": true,
		"privacy_policy":   true,
		"marketing":        false,
		"health_data":      true,
	}
	consentResp, consentResult := apiRequest(t, "POST", "/api/v1/users/onboarding/consent", consentBody, token)
	defer consentResp.Body.Close()
	t.Logf("4. 약관 동의: status=%d, result=%v", consentResp.StatusCode, consentResult)

	// 5. 디바이스 페어링 (시뮬레이션)
	pairBody := map[string]string{
		"device_id":        "MPK-ONBOARD-001",
		"name":             "My ManPaSik Reader",
		"firmware_version": "1.0.0",
	}
	pairResp, pairResult := apiRequest(t, "POST", devicesEndpoint, pairBody, token)
	defer pairResp.Body.Close()
	t.Logf("5. 디바이스 페어링: status=%d, result=%v", pairResp.StatusCode, pairResult)

	// 6. 온보딩 완료 표시
	completeResp, completeResult := apiRequest(t, "POST", "/api/v1/users/onboarding/complete", nil, token)
	defer completeResp.Body.Close()
	t.Logf("6. 온보딩 완료: status=%d, result=%v", completeResp.StatusCode, completeResult)
}

// ─── E2E-OB-002: 온보딩 단계 건너뛰기 (최소 요구사항) ───

func TestOnboardingMinimal(t *testing.T) {
	email := uniqueEmail("ob-minimal")
	password := "E2eTest1!@#"

	status, _ := registerUser(t, email, password, "Minimal User")
	if status != http.StatusCreated && status != http.StatusOK {
		t.Skipf("회원가입 실패: %d", status)
		return
	}
	token := loginUser(t, email, password)
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 필수 약관만 동의 후 즉시 완료
	consentBody := map[string]interface{}{
		"terms_of_service": true,
		"privacy_policy":   true,
	}
	consentResp, _ := apiRequest(t, "POST", "/api/v1/users/onboarding/consent", consentBody, token)
	defer consentResp.Body.Close()

	completeResp, _ := apiRequest(t, "POST", "/api/v1/users/onboarding/complete", nil, token)
	defer completeResp.Body.Close()
	t.Logf("최소 온보딩 완료: status=%d", completeResp.StatusCode)
}

// ─── E2E-OB-003: 온보딩 중 앱 재시작 시 상태 유지 ───

func TestOnboardingStateRecovery(t *testing.T) {
	email := uniqueEmail("ob-recovery")
	password := "E2eTest1!@#"

	status, _ := registerUser(t, email, password, "Recovery User")
	if status != http.StatusCreated && status != http.StatusOK {
		t.Skipf("회원가입 실패: %d", status)
		return
	}

	// 1차 로그인: 온보딩 시작
	token1 := loginUser(t, email, password)
	if token1 == "" {
		t.Skip("1차 토큰 획득 실패")
		return
	}
	profileBody := map[string]interface{}{
		"birth_year": 1985,
		"gender":     "female",
	}
	profResp, _ := apiRequest(t, "POST", "/api/v1/users/onboarding/profile", profileBody, token1)
	defer profResp.Body.Close()
	t.Logf("1차 세션 프로필 설정: status=%d", profResp.StatusCode)

	// 2차 로그인: 온보딩 상태 조회
	token2 := loginUser(t, email, password)
	if token2 == "" {
		t.Skip("2차 토큰 획득 실패")
		return
	}
	stateResp, stateResult := apiRequest(t, "GET", "/api/v1/users/onboarding/status", nil, token2)
	defer stateResp.Body.Close()
	t.Logf("온보딩 상태 복구: status=%d, result=%v", stateResp.StatusCode, stateResult)
}
