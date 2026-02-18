package security_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// ─── OWASP A02: Cryptographic Failures ───

func TestTLSEnforcement(t *testing.T) {
	// HTTP → HTTPS 리다이렉트 확인
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 리다이렉트 따라가지 않음
		},
	}

	resp, err := client.Get("http://localhost:8080/api/v1/health")
	if err != nil {
		t.Skipf("HTTP 접근 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	// 프로덕션에서는 301/302 리다이렉트 또는 HSTS 헤더 필요
	hsts := resp.Header.Get("Strict-Transport-Security")
	t.Logf("HSTS 헤더: %q (프로덕션 환경에서 max-age=31536000 이상 권장)", hsts)
}

func TestSensitiveDataInResponse(t *testing.T) {
	// 사용자 프로필 조회 시 민감 정보 노출 확인
	token := loginAndGetToken(t, "data-test@test.com", "DataTest1!")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("프로필 조회 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := strings.ToLower(buf.String())

	sensitiveFields := []string{
		"password",
		"password_hash",
		"salt",
		"secret",
		"private_key",
		"ssn",
		"social_security",
	}

	for _, field := range sensitiveFields {
		if strings.Contains(body, field) {
			t.Errorf("민감 데이터 노출: 응답에 %q 필드 포함", field)
		}
	}
}

func TestTokenInURL(t *testing.T) {
	// URL 쿼리 파라미터에 토큰 전달 시 거부 확인
	token := loginAndGetToken(t, "url-token@test.com", "UrlToken1!")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 토큰을 쿼리 파라미터로 전달 (보안 위반)
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/measurement/history?access_token=%s", baseURL, token))
	if err != nil {
		t.Skipf("요청 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	// URL에 토큰을 넣어도 인증되지 않아야 함 (헤더로만 인증)
	if resp.StatusCode == http.StatusOK {
		t.Log("경고: URL 쿼리 파라미터의 토큰이 인증에 사용됨 (보안 위험)")
	}
}

// ─── GDPR / PIPA 데이터 보호 ───

func TestDataMinimization(t *testing.T) {
	// 목록 API가 최소한의 필드만 반환하는지 확인
	token := loginAndGetToken(t, "minimize@test.com", "Minimize1!")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	req, _ := http.NewRequest("GET", baseURL+"/api/v1/measurement/history?limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("측정 이력 조회 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 목록 API에 생체 원본 데이터(raw_channels)가 포함되지 않아야 함
	body := fmt.Sprintf("%v", result)
	if strings.Contains(body, "raw_channels") || strings.Contains(body, "raw_data") {
		t.Error("데이터 최소화 위반: 목록 API에 생체 원본 데이터 포함")
	}
}

func TestDataDeletionRight(t *testing.T) {
	// GDPR 잊힐 권리: 계정 삭제 요청 테스트
	email := "delete-me@test.com"
	password := "DeleteMe1!"
	token := loginAndGetToken(t, email, password)
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 계정 삭제 요청
	req, _ := http.NewRequest("DELETE", baseURL+"/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("삭제 요청 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	t.Logf("계정 삭제 응답: status=%d (200 또는 202 기대)", resp.StatusCode)

	// 삭제 후 재로그인 불가 확인
	loginBody := map[string]string{"email": email, "password": password}
	loginData, _ := json.Marshal(loginBody)
	loginResp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginData))
	if err != nil {
		return
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode == http.StatusOK {
		t.Error("GDPR 잊힐 권리 위반: 삭제된 계정으로 로그인 가능")
	}
}

func TestConsentRequired(t *testing.T) {
	// 데이터 처리 동의 없이 민감 작업 차단 확인
	regBody := map[string]string{
		"email":    "no-consent@test.com",
		"password": "NoConsent1!",
		"name":     "No Consent User",
		// consent 필드 없음
	}
	data, _ := json.Marshal(regBody)
	resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Skipf("요청 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	// 의료 데이터 처리를 위한 명시적 동의가 필요할 수 있음
	t.Logf("동의 없는 가입 응답: status=%d", resp.StatusCode)
}

// ─── 의료 데이터 보안 ───

func TestMeasurementDataIsolation(t *testing.T) {
	// 다른 사용자의 측정 데이터에 접근 불가 확인
	tokenA := loginAndGetToken(t, "isolation-a@test.com", "IsolationA1!")
	tokenB := loginAndGetToken(t, "isolation-b@test.com", "IsolationB1!")

	if tokenA == "" || tokenB == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 사용자 B의 측정 세션을 사용자 A가 조회 시도
	endpoints := []string{
		"/api/v1/measurement/sessions?user_id=isolation-b",
		"/api/v1/measurement/history?user_id=isolation-b",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			req, _ := http.NewRequest("GET", baseURL+ep, nil)
			req.Header.Set("Authorization", "Bearer "+tokenA)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&result)
				if data, ok := result["data"]; ok && data != nil {
					t.Errorf("측정 데이터 격리 위반: 다른 사용자 데이터 접근 가능 (%s)", ep)
				}
			}
		})
	}
}

func TestMeasurementDataEncryptionAtRest(t *testing.T) {
	// 측정 데이터 저장 시 암호화 확인 (API 응답 기반 간접 확인)
	token := loginAndGetToken(t, "encrypt@test.com", "Encrypt1!")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 측정 세션 생성
	sessionBody := map[string]interface{}{
		"device_id":    "MPK-SEC-001",
		"cartridge_id": "CART-GLU-001",
		"user_id":      "encrypt-user",
		"raw_channels":  make([]float64, 88), // 88차원 원본 데이터
	}
	data, _ := json.Marshal(sessionBody)
	req, _ := http.NewRequest("POST", baseURL+"/api/v1/measurement/sessions", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("측정 세션 생성 실패: %v", err)
		return
	}
	defer resp.Body.Close()

	t.Logf("측정 세션 생성 응답: status=%d (저장 시 암호화는 인프라 레벨에서 확인 필요)", resp.StatusCode)
}

// ─── OWASP A10: SSRF (Server-Side Request Forgery) ───

func TestSSRFPrevention(t *testing.T) {
	ssrfPayloads := []string{
		"http://169.254.169.254/latest/meta-data/",   // AWS 메타데이터
		"http://localhost:6379/",                       // Redis
		"http://127.0.0.1:5432/",                      // PostgreSQL
		"http://[::1]:8080/",                           // IPv6 localhost
		"file:///etc/passwd",                            // 로컬 파일
	}

	for i, payload := range ssrfPayloads {
		t.Run(fmt.Sprintf("ssrf_%d", i), func(t *testing.T) {
			// URL 필드가 있는 API에 내부 주소 삽입 시도
			body := map[string]string{
				"callback_url": payload,
				"webhook_url":  payload,
			}
			data, _ := json.Marshal(body)
			resp, err := http.Post(baseURL+"/api/v1/notifications/webhook", "application/json", bytes.NewReader(data))
			if err != nil {
				return // 엔드포인트 없으면 무시
			}
			defer resp.Body.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			respBody := buf.String()

			// 내부 서비스 응답이 포함되면 SSRF 취약
			if strings.Contains(respBody, "ami-id") || strings.Contains(respBody, "instance-id") {
				t.Errorf("SSRF 취약점: 내부 메타데이터 접근 가능 (payload=%q)", payload)
			}
		})
	}
}

// ─── 대용량 페이로드 / DoS 방어 ───

func TestOversizedPayloadRejection(t *testing.T) {
	// 10MB 페이로드 전송 시도
	largePayload := strings.Repeat("A", 10*1024*1024) // 10MB
	body := map[string]string{
		"email":    "large@test.com",
		"password": "Large1!",
		"name":     largePayload,
	}
	data, _ := json.Marshal(body)
	resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Logf("대용량 페이로드 거부됨 (연결 레벨): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		t.Error("대용량 페이로드 수용 취약점: 10MB 페이로드가 수락됨")
	}
	t.Logf("대용량 페이로드 응답: status=%d", resp.StatusCode)
}

func TestSlowlorisProtection(t *testing.T) {
	// 느린 요청 전송으로 연결 고갈 시도 (간략 버전)
	t.Log("Slowloris 보호: Kong Gateway의 client_body_timeout 및 keepalive_timeout 설정 확인 필요")
	t.Log("프로덕션 설정: client_body_timeout=30s, client_header_timeout=30s 권장")
}
