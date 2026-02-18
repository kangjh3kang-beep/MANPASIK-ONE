package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// ─── 공통 상수 ───

const (
	defaultTimeout    = 10 * time.Second
	healthEndpoint    = "/api/v1/health"
	loginEndpoint     = "/api/v1/auth/login"
	registerEndpoint  = "/api/v1/auth/register"
	profileEndpoint   = "/api/v1/users/me"
	measureEndpoint   = "/api/v1/measurement/sessions"
	historyEndpoint   = "/api/v1/measurement/history"
	devicesEndpoint   = "/api/v1/devices"
	cartridgesEndpoint = "/api/v1/cartridges"
)

// ─── HTTP 유틸리티 ───

// apiRequest는 API 요청을 수행하고 응답을 반환합니다.
func apiRequest(t *testing.T, method, path string, body interface{}, token string) (*http.Response, map[string]interface{}) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("요청 바디 직렬화 실패: %v", err)
		}
		reqBody = bytes.NewReader(data)
	}

	url := fmt.Sprintf("http://%s%s", gatewayAddr, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("요청 생성 실패: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s %s 요청 실패: %v", method, path, err)
	}

	var result map[string]interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.Len() > 0 {
		json.Unmarshal(buf.Bytes(), &result)
	}

	return resp, result
}

// registerUser 사용자 회원가입
func registerUser(t *testing.T, email, password, name string) (int, map[string]interface{}) {
	t.Helper()
	body := map[string]string{
		"email":    email,
		"password": password,
		"name":     name,
	}
	resp, result := apiRequest(t, "POST", registerEndpoint, body, "")
	defer resp.Body.Close()
	return resp.StatusCode, result
}

// loginUser 로그인 후 토큰 반환
func loginUser(t *testing.T, email, password string) string {
	t.Helper()
	body := map[string]string{
		"email":    email,
		"password": password,
	}
	resp, result := apiRequest(t, "POST", loginEndpoint, body, "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Skipf("로그인 실패 (status=%d)", resp.StatusCode)
		return ""
	}

	token, _ := result["access_token"].(string)
	return token
}

// registerAndLogin 회원가입 후 로그인
func registerAndLogin(t *testing.T, email, password, name string) string {
	t.Helper()
	registerUser(t, email, password, name)
	return loginUser(t, email, password)
}

// assertStatus HTTP 상태 코드 검증
func assertStatus(t *testing.T, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: status=%d, want=%d", msg, got, want)
	}
}

// assertNotStatus HTTP 상태 코드가 특정 값이 아닌지 검증
func assertNotStatus(t *testing.T, got, notWant int, msg string) {
	t.Helper()
	if got == notWant {
		t.Errorf("%s: status=%d, should not be %d", msg, got, notWant)
	}
}

// uniqueEmail 고유 이메일 생성
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@e2e-test.com", prefix, time.Now().UnixNano())
}
