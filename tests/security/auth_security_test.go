package security_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

// ─── OWASP A01: Broken Access Control ───

func TestHorizontalPrivilegeEscalation(t *testing.T) {
	// 사용자 A의 토큰으로 사용자 B의 리소스에 접근 시도
	tokenA := loginAndGetToken(t, "userA@test.com", "PasswordA1!")
	tokenB := loginAndGetToken(t, "userB@test.com", "PasswordB1!")

	// 사용자 A가 사용자 B의 측정 기록에 접근 시도
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/measurement/history?user_id=userB", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusNotFound {
		t.Errorf("수평 권한 상승 취약점: 다른 사용자 리소스에 접근 가능 (status=%d)", resp.StatusCode)
	}
	_ = tokenB // tokenB는 비교용으로만 생성
}

func TestVerticalPrivilegeEscalation(t *testing.T) {
	// 일반 사용자 토큰으로 관리자 API 접근 시도
	token := loginAndGetToken(t, "regular@test.com", "Regular1!")

	adminEndpoints := []string{
		"/api/v1/admin/users",
		"/api/v1/admin/cartridges",
		"/api/v1/admin/system/config",
		"/api/v1/admin/audit-logs",
	}

	for _, endpoint := range adminEndpoints {
		t.Run(endpoint, func(t *testing.T) {
			req, _ := http.NewRequest("GET", baseURL+endpoint, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Skipf("엔드포인트 접근 불가: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Errorf("수직 권한 상승 취약점: 일반 사용자가 관리자 API 접근 가능 (%s)", endpoint)
			}
		})
	}
}

func TestExpiredTokenRejection(t *testing.T) {
	// 만료된 토큰으로 접근 시도
	expiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IiwiZXhwIjoxNjAwMDAwMDAwfQ.invalid"

	req, _ := http.NewRequest("GET", baseURL+"/api/v1/measurement/history", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("요청 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("만료 토큰이 거부되지 않음: status=%d", resp.StatusCode)
	}
}

func TestMalformedTokenRejection(t *testing.T) {
	malformedTokens := []string{
		"",
		"Bearer",
		"Bearer ",
		"Bearer invalid-token",
		"Bearer eyJhbGciOiJub25lIn0.eyJzdWIiOiJhZG1pbiJ9.", // alg:none 공격
		"Basic dXNlcjpwYXNz",                                  // 잘못된 인증 방식
	}

	for i, token := range malformedTokens {
		t.Run(fmt.Sprintf("malformed_%d", i), func(t *testing.T) {
			req, _ := http.NewRequest("GET", baseURL+"/api/v1/measurement/history", nil)
			if token != "" {
				req.Header.Set("Authorization", token)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Errorf("비정상 토큰이 승인됨: token=%q, status=%d", token, resp.StatusCode)
			}
		})
	}
}

// ─── OWASP A07: Authentication Failures ───

func TestBruteForceProtection(t *testing.T) {
	// 동일 계정에 연속 로그인 실패 시 차단 확인
	loginBody := map[string]string{
		"email":    "bruteforce-target@test.com",
		"password": "WrongPassword!",
	}
	body, _ := json.Marshal(loginBody)

	var lastStatus int
	blocked := false
	for i := 0; i < 15; i++ {
		resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
		if err != nil {
			continue
		}
		lastStatus = resp.StatusCode
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
			blocked = true
			t.Logf("무차별 대입 보호 활성: %d회 시도 후 차단 (status=%d)", i+1, resp.StatusCode)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !blocked {
		t.Errorf("무차별 대입 공격 보호 미흡: 15회 실패 후에도 차단 안됨 (last_status=%d)", lastStatus)
	}
}

func TestPasswordComplexityEnforcement(t *testing.T) {
	weakPasswords := []string{
		"123456",
		"password",
		"abc",
		"aaaaaaaaaa",
		"12345678",
	}

	for _, pw := range weakPasswords {
		t.Run(pw, func(t *testing.T) {
			regBody := map[string]string{
				"email":    fmt.Sprintf("weak-%s@test.com", pw),
				"password": pw,
				"name":     "WeakPW Test",
			}
			body, _ := json.Marshal(regBody)
			resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(body))
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				t.Errorf("약한 비밀번호(%q)가 허용됨", pw)
			}
		})
	}
}

func TestSessionFixation(t *testing.T) {
	// 로그인 전후 세션/토큰이 변경되는지 확인
	loginBody := map[string]string{
		"email":    "session-fix@test.com",
		"password": "SessionFix1!",
	}
	body, _ := json.Marshal(loginBody)

	// 첫 번째 로그인
	resp1, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Skipf("로그인 실패: %v", err)
		return
	}
	defer resp1.Body.Close()

	var result1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&result1)
	token1, _ := result1["access_token"].(string)

	time.Sleep(1 * time.Second)

	// 두 번째 로그인 (같은 자격증명)
	resp2, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Skipf("재로그인 실패: %v", err)
		return
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&result2)
	token2, _ := result2["access_token"].(string)

	if token1 != "" && token2 != "" && token1 == token2 {
		t.Error("세션 고정 취약점: 재로그인 시 동일 토큰 발급")
	}
}

// ─── OWASP A04: Insecure Design (Rate Limiting) ───

func TestRateLimiting(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/register"},
		{"GET", "/api/v1/measurement/history"},
	}

	for _, ep := range endpoints {
		t.Run(ep.path, func(t *testing.T) {
			rateLimited := false
			for i := 0; i < 200; i++ {
				var resp *http.Response
				var err error
				if ep.method == "POST" {
					resp, err = http.Post(baseURL+ep.path, "application/json",
						strings.NewReader(`{"email":"rate@test.com","password":"Rate1!"}`))
				} else {
					resp, err = http.Get(baseURL + ep.path)
				}
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == http.StatusTooManyRequests {
					rateLimited = true
					t.Logf("Rate limit 활성: %d회 요청 후 차단", i+1)
					break
				}
			}
			if !rateLimited {
				t.Logf("경고: %s %s에 rate limiting 미적용 (200회 요청 후에도 미차단)", ep.method, ep.path)
			}
		})
	}
}

// ─── Helper Functions ───

func loginAndGetToken(t *testing.T, email, password string) string {
	t.Helper()

	// 먼저 회원가입 시도
	regBody := map[string]string{"email": email, "password": password, "name": "Security Test"}
	regData, _ := json.Marshal(regBody)
	http.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(regData))

	// 로그인
	loginBody := map[string]string{"email": email, "password": password}
	loginData, _ := json.Marshal(loginBody)
	resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginData))
	if err != nil {
		t.Skipf("로그인 실패: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	token, _ := result["access_token"].(string)
	return token
}
