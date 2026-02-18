# ManPaSik 보안 테스트 (Security Tests)

## 개요
OWASP Top 10 (2021) 및 의료 데이터 보호 규정(GDPR, PIPA, HIPAA)에 기반한 보안 테스트 스위트.

## 테스트 구성

### auth_security_test.go — 인증/권한 보안
| 테스트 | OWASP | 설명 |
|--------|-------|------|
| TestHorizontalPrivilegeEscalation | A01 | 다른 사용자 리소스 접근 차단 |
| TestVerticalPrivilegeEscalation | A01 | 일반→관리자 권한 상승 차단 |
| TestExpiredTokenRejection | A07 | 만료 토큰 거부 |
| TestMalformedTokenRejection | A07 | 비정상 토큰 거부 (alg:none 포함) |
| TestBruteForceProtection | A07 | 무차별 대입 공격 차단 |
| TestPasswordComplexityEnforcement | A07 | 약한 비밀번호 거부 |
| TestSessionFixation | A07 | 세션 고정 공격 방지 |
| TestRateLimiting | A04 | API Rate Limiting 적용 확인 |

### api_security_test.go — API 보안
| 테스트 | OWASP | 설명 |
|--------|-------|------|
| TestSQLInjection | A03 | SQL 인젝션 방어 |
| TestXSSInjection | A03 | XSS 인젝션 방어 |
| TestCommandInjection | A03 | 명령 인젝션 방어 |
| TestSecurityHeaders | A05 | 보안 헤더 설정 확인 |
| TestServerInformationLeakage | A05 | 서버 정보 노출 방지 |
| TestErrorDetailExposure | A05 | 에러 상세 정보 노출 방지 |
| TestDirectoryTraversal | A06 | 경로 탐색 공격 방어 |
| TestJWTAlgorithmConfusion | A08 | JWT 알고리즘 혼동 공격 방어 |
| TestAuditTrailOnSensitiveActions | A09 | 감사 추적 생성 확인 |

### data_security_test.go — 데이터 보안
| 테스트 | 규정 | 설명 |
|--------|------|------|
| TestTLSEnforcement | A02 | TLS/HSTS 적용 확인 |
| TestSensitiveDataInResponse | A02 | 민감 데이터 응답 노출 방지 |
| TestTokenInURL | A02 | URL 내 토큰 전달 방지 |
| TestDataMinimization | GDPR | 데이터 최소화 원칙 |
| TestDataDeletionRight | GDPR | 잊힐 권리 (계정 삭제) |
| TestConsentRequired | PIPA | 데이터 처리 동의 확인 |
| TestMeasurementDataIsolation | HIPAA | 측정 데이터 사용자 격리 |
| TestMeasurementDataEncryptionAtRest | HIPAA | 저장 시 암호화 확인 |
| TestSSRFPrevention | A10 | SSRF 공격 방어 |
| TestOversizedPayloadRejection | - | 대용량 페이로드 거부 |

### dependency_scan.sh — 의존성 취약점 스캔
- Go 모듈: `govulncheck` 기반 취약점 탐지
- Flutter/Dart: `dart pub outdated` 기반 버전 확인
- Rust crate: `cargo-audit` 기반 보안 감사
- Docker: Dockerfile 보안 체크리스트 + `trivy` 스캔

## 실행 방법

```bash
# 전체 보안 테스트
cd tests/security
go test -v -timeout 5m ./...

# 개별 테스트 실행
go test -v -run TestSQLInjection ./...
go test -v -run TestBruteForceProtection ./...

# 의존성 스캔
chmod +x dependency_scan.sh
./dependency_scan.sh

# 특정 OWASP 카테고리 테스트
go test -v -run "Injection" ./...        # A03
go test -v -run "SecurityHeaders" ./...  # A05
go test -v -run "Token" ./...            # A07
```

## 사전 요구사항
- Go 1.21+
- ManPaSik 서비스 실행 중 (`localhost:8080`)
- (선택) `govulncheck`, `cargo-audit`, `trivy`

## 보안 목표 (Phase 1)

| 항목 | 목표 |
|------|------|
| OWASP Top 10 커버리지 | 10/10 (100%) |
| 인증 취약점 테스트 | 8개 시나리오 |
| 인젝션 방어 테스트 | SQL, XSS, CMD, SSRF |
| 데이터 보호 테스트 | GDPR + PIPA + HIPAA |
| 의존성 스캔 | Go + Dart + Rust + Docker |
