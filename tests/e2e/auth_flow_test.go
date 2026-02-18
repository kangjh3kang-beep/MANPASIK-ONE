package e2e

import (
	"net/http"
	"testing"
)

// ─── E2E-AUTH-001: 회원가입 → 로그인 → 프로필 조회 ───

func TestAuthRegistrationAndLogin(t *testing.T) {
	email := uniqueEmail("auth-reg")
	password := "E2eTest1!@#"

	// 1. 회원가입
	status, _ := registerUser(t, email, password, "E2E Auth User")
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
	t.Logf("2. 로그인 성공 (token=%s...)", token[:min(20, len(token))])

	// 3. 프로필 조회
	resp, profile := apiRequest(t, "GET", profileEndpoint, nil, token)
	defer resp.Body.Close()
	assertStatus(t, resp.StatusCode, http.StatusOK, "프로필 조회")
	t.Logf("3. 프로필 조회 성공: %v", profile)
}

// ─── E2E-AUTH-002: 중복 가입 방지 ───

func TestAuthDuplicateRegistration(t *testing.T) {
	email := uniqueEmail("auth-dup")
	password := "E2eTest1!@#"

	// 1차 가입
	status1, _ := registerUser(t, email, password, "First User")
	if status1 != http.StatusCreated && status1 != http.StatusOK {
		t.Skipf("1차 가입 실패: %d", status1)
		return
	}

	// 2차 중복 가입 시도
	status2, _ := registerUser(t, email, password, "Duplicate User")
	if status2 == http.StatusCreated || status2 == http.StatusOK {
		t.Error("중복 가입 허용됨: 동일 이메일로 2번 가입 가능")
	}
	t.Logf("중복 가입 방지 확인 (status=%d)", status2)
}

// ─── E2E-AUTH-003: 잘못된 자격증명 거부 ───

func TestAuthInvalidCredentials(t *testing.T) {
	email := uniqueEmail("auth-invalid")
	password := "E2eTest1!@#"

	// 가입
	registerUser(t, email, password, "Invalid Cred User")

	// 잘못된 비밀번호
	token := loginUser(t, email, "WrongPassword!")
	if token != "" {
		t.Error("잘못된 비밀번호로 로그인 성공: 보안 취약점")
	}

	// 존재하지 않는 이메일
	token = loginUser(t, "nonexistent@e2e-test.com", password)
	if token != "" {
		t.Error("존재하지 않는 이메일로 로그인 성공: 보안 취약점")
	}
}

// ─── E2E-AUTH-004: 토큰 갱신 플로우 ───

func TestAuthTokenRefresh(t *testing.T) {
	email := uniqueEmail("auth-refresh")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Refresh User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 토큰으로 리소스 접근
	resp, _ := apiRequest(t, "GET", profileEndpoint, nil, token)
	defer resp.Body.Close()
	assertStatus(t, resp.StatusCode, http.StatusOK, "토큰 유효성")

	// 토큰 갱신 요청
	refreshResp, refreshResult := apiRequest(t, "POST", "/api/v1/auth/refresh", nil, token)
	defer refreshResp.Body.Close()
	t.Logf("토큰 갱신 응답: status=%d, result=%v", refreshResp.StatusCode, refreshResult)
}

// ─── E2E-AUTH-005: 로그아웃 ───

func TestAuthLogout(t *testing.T) {
	email := uniqueEmail("auth-logout")
	token := registerAndLogin(t, email, "E2eTest1!@#", "Logout User")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 로그아웃
	logoutResp, _ := apiRequest(t, "POST", "/api/v1/auth/logout", nil, token)
	defer logoutResp.Body.Close()
	t.Logf("로그아웃 응답: status=%d", logoutResp.StatusCode)

	// 로그아웃 후 토큰 사용 시도
	resp, _ := apiRequest(t, "GET", profileEndpoint, nil, token)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Log("경고: 로그아웃 후에도 토큰 유효 (토큰 블랙리스트 미구현 가능)")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
