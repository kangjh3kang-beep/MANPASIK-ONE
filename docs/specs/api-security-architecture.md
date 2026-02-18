# API 보안 아키텍처 (API Security Architecture)

> **문서 ID**: MPS-SEC-API-001
> **버전**: 1.0
> **작성일**: 2026-02-14
> **적용 표준**: IEC 81001-5-1, OWASP API Security Top 10, NIST SP 800-53

---

## 1. 아키텍처 개요

```
클라이언트 (Flutter/Web)
    │  TLS 1.3
    ▼
┌─────────────────────────┐
│  Kong API Gateway       │  Rate Limiting, IP Whitelist, WAF
│  (Port 8000/8443)       │  JWT 검증, CORS, Security Headers
└────────────┬────────────┘
             │  mTLS (내부)
    ┌────────┼────────┐
    ▼        ▼        ▼
┌───────┐┌───────┐┌───────┐
│Auth   ││User   ││Meas.  │  22개 gRPC 마이크로서비스
│Service││Service││Service│  각 서비스별 RBAC 검증
└───┬───┘└───┬───┘└───┬───┘
    │        │        │
    ▼        ▼        ▼
┌─────────────────────────┐
│  Keycloak OIDC          │  ID/Access/Refresh Token 발급
│  (SSO + MFA)            │  5개 역할, Realm 관리
└─────────────────────────┘
```

---

## 2. 인증 체계 (Authentication)

### 2.1 토큰 아키텍처

| 토큰 | 형식 | 유효기간 | 용도 |
|------|------|---------|------|
| Access Token | JWT (RS256) | 15분 | API 호출 인증 |
| Refresh Token | Opaque + DB 저장 | 7일 | Access Token 갱신 |
| ID Token | JWT (RS256) | 15분 | 사용자 정보 (OIDC) |
| Device Token | JWT (ES256) | 30일 | 리더기-서버 인증 |
| Offline Token | Encrypted (AES-256) | 72시간 | 오프라인 모드 |

### 2.2 JWT 구조

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT",
    "kid": "manpasik-key-2026-02"
  },
  "payload": {
    "sub": "user-uuid",
    "iss": "https://auth.manpasik.com",
    "aud": "manpasik-api",
    "exp": 1708000000,
    "iat": 1707999100,
    "roles": ["user"],
    "tier": "standard",
    "device_id": "device-uuid",
    "scope": "measurement:read measurement:write profile:read"
  }
}
```

### 2.3 다중 인증 (MFA)

| 조건 | MFA 요구 | 방식 |
|------|---------|------|
| 일반 로그인 | 선택 | TOTP / SMS |
| Clinical 티어 | **필수** | TOTP + 생체인증 |
| 관리자 로그인 | **필수** | TOTP + HW 키 (FIDO2) |
| 민감 데이터 접근 | 재인증 | Step-up Auth |
| 비대면 진료 | **필수** | 본인인증 + 생체 |

---

## 3. 인가 체계 (Authorization)

### 3.1 RBAC 역할 매트릭스

| 역할 | 측정 | 결과 조회 | 가족 | 진료 | 관리 | 감사 로그 |
|------|------|---------|------|------|------|---------|
| user | RW | Own | Join | Book | - | - |
| family_admin | RW | Family | CRUD | Book | - | - |
| medical | R | Patient | R | RW | - | R |
| admin | R | All | R | R | RW | R |
| system | RW | All | RW | RW | RW | RW |

### 3.2 리소스 수준 접근 제어

```go
// middleware/rbac.go 패턴
func AuthorizeMeasurement(ctx context.Context, measurementID string) error {
    userID := auth.GetUserID(ctx)
    roles := auth.GetRoles(ctx)

    // 1. 본인 데이터 → 항상 허용
    if measurement.OwnerID == userID { return nil }

    // 2. 가족 데이터 → 공유 동의 확인
    if hasRole(roles, "family_admin") {
        if hasSharingConsent(measurement.OwnerID, userID) { return nil }
    }

    // 3. 의료진 → 진료 관계 확인
    if hasRole(roles, "medical") {
        if hasActiveTreatment(measurement.OwnerID, userID) { return nil }
    }

    return ErrForbidden
}
```

### 3.3 Scope 기반 API 제어

| Scope | 허용 API | 설명 |
|-------|---------|------|
| `measurement:read` | GET /measurements/* | 측정 데이터 조회 |
| `measurement:write` | POST /measurements/* | 측정 실행/저장 |
| `profile:read` | GET /users/me | 프로필 조회 |
| `profile:write` | PUT /users/me | 프로필 수정 |
| `family:manage` | /families/* | 가족 그룹 관리 |
| `medical:access` | /telemedicine/*, /prescriptions/* | 의료 서비스 |
| `admin:full` | /admin/* | 관리자 전체 |

---

## 4. API Gateway 보안 (Kong)

### 4.1 보안 플러그인 구성

| 플러그인 | 설정 | 목적 |
|---------|------|------|
| jwt | RS256, Keycloak JWKS | 토큰 검증 |
| rate-limiting | 100 req/min (user), 1000 req/min (service) | DDoS 방어 |
| ip-restriction | 관리자 IP 화이트리스트 | 접근 제한 |
| cors | 허용 도메인 목록 | XSS 방어 |
| request-size-limiting | max 5MB | 대용량 공격 방어 |
| bot-detection | User-Agent 검증 | 봇 차단 |
| response-transformer | Server 헤더 제거 | 정보 노출 방지 |

### 4.2 보안 헤더

```
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 0
Content-Security-Policy: default-src 'self'; script-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(self), geolocation=()
```

---

## 5. 키 관리 시스템 (KMS)

### 5.1 키 계층 구조

```
Master Key (AWS KMS / HashiCorp Vault)
    ├── JWT Signing Key (RSA-2048, 90일 회전)
    ├── Data Encryption Key (AES-256, 연 1회 회전)
    │   ├── PHI Encryption (L4 데이터)
    │   ├── PII Encryption (L3 데이터)
    │   └── Backup Encryption
    ├── BLE Session Key (ECDH P-256, 세션별)
    ├── NFC Auth Key (AES-128, 카트리지별)
    └── HMAC Signing Key (SHA-256, API 무결성)
```

### 5.2 키 회전 정책

| 키 유형 | 회전 주기 | 자동화 | 무중단 |
|---------|---------|--------|--------|
| JWT Signing | 90일 | Vault auto-rotate | kid 기반 다중 키 |
| Data Encryption | 365일 | Vault + re-encryption job | Envelope 암호화 |
| BLE Session | 세션별 | ECDH 핸드셰이크 | N/A |
| TLS Certificate | 90일 | cert-manager (Let's Encrypt) | 자동 갱신 |

---

## 6. PHI 보호 아키텍처

### 6.1 데이터 흐름 보안

```
[리더기] --AES-CCM/BLE--> [앱]
   │                        │
   │  측정 원시 데이터         │  AES-256 로컬 암호화
   │                        │
   └─── NFC 인증 ───────────┘
                             │
                    TLS 1.3  │  gRPC (Protobuf)
                             ▼
                    [API Gateway]
                             │
                    mTLS     │  내부 서비스 간
                             ▼
                    [Measurement Service]
                             │
                    AES-256  │  Column-level encryption
                             ▼
                    [TimescaleDB / PostgreSQL]
```

### 6.2 PHI 접근 감사

| 이벤트 | 로그 항목 | 보존 |
|--------|---------|------|
| PHI 읽기 | user_id, resource_id, timestamp, IP, role | 10년 |
| PHI 쓰기 | user_id, resource_id, before/after hash | 10년 |
| PHI 내보내기 | user_id, format, record_count | 10년 |
| PHI 삭제 | user_id, resource_id, deletion_type | 10년 |
| 비정상 접근 시도 | user_id, resource_id, denial_reason | 10년 |

---

## 7. 동시 세션 관리 정책

### 7.1 세션 제한

| 역할 | 최대 동시 세션 | 초과 시 동작 |
|------|-------------|------------|
| user | 3 (앱 2 + 웹 1) | 가장 오래된 세션 종료 |
| family_admin | 3 | 가장 오래된 세션 종료 |
| medical | 2 | 새 로그인 차단 (확인 필요) |
| admin | 1 | 기존 세션 강제 종료 + 알림 |

### 7.2 세션 무효화

| 트리거 | 동작 |
|--------|------|
| 로그아웃 | 즉시 Refresh Token 폐기 |
| 비밀번호 변경 | 전 세션 무효화 |
| 역할 변경 | 전 세션 무효화, 재인증 필요 |
| 30분 미활동 | Access Token 미갱신 (자동 만료) |
| 이상 행위 탐지 | 전 세션 강제 종료 + 계정 잠금 |
| 디바이스 변경 | Step-up 인증 요구 |

---

## 8. API 공격 방어

### 8.1 OWASP API Security Top 10 대응

| # | 위협 | 방어 | 구현 위치 |
|---|------|------|---------|
| API1 | 객체 수준 권한 취약 | 리소스 소유자 검증 | middleware/rbac.go |
| API2 | 인증 취약 | JWT + MFA + Rate Limit | Kong + Keycloak |
| API3 | 객체 속성 수준 권한 | DTO 필터링 (역할별 응답 필드) | handler/grpc.go |
| API4 | 무제한 리소스 소비 | Rate Limiting + Request Size | Kong plugins |
| API5 | 기능 수준 권한 취약 | RBAC + Scope 검증 | middleware/rbac.go |
| API6 | 민감 데이터 노출 | TLS + 응답 필터링 + 마스킹 | Kong + service layer |
| API7 | 보안 설정 오류 | 보안 헤더 + CORS | Kong response-transformer |
| API8 | 자동화 위협 | Bot Detection + CAPTCHA | Kong bot-detection |
| API9 | 부적절한 자산 관리 | API 버전 관리 + 문서화 | OpenAPI spec |
| API10 | 서버측 요청 위조 | SSRF 필터 (내부 IP 차단) | validation/sanitizer.go |

---

## 9. 침해 대응 자동화

### 9.1 실시간 탐지 규칙

| 규칙 | 조건 | 자동 대응 |
|------|------|---------|
| 무차별 대입 | 5회 연속 인증 실패 | 계정 15분 잠금 |
| 토큰 재사용 | 폐기된 Refresh Token 사용 | 전 세션 무효화 |
| 비정상 지역 | 새로운 국가에서 로그인 | MFA 강제 + 알림 |
| 대량 데이터 접근 | 1시간 내 100건 이상 PHI 조회 | 접근 차단 + 보안팀 알림 |
| API 스캐닝 | 존재하지 않는 엔드포인트 연속 호출 | IP 차단 (1시간) |
| SQL 인젝션 시도 | WAF 패턴 매칭 | 요청 차단 + 로깅 |

### 9.2 알림 체계

```
탐지 → Prometheus Alert → AlertManager
                              │
                 ┌────────────┼────────────┐
                 ▼            ▼            ▼
            Slack #보안    PagerDuty     이메일
            (정보)        (긴급)       (보고서)
```

---

## 10. 문서 이력

| 버전 | 날짜 | 변경 | 작성자 |
|------|------|------|--------|
| 1.0 | 2026-02-14 | 초안 작성 | Claude |

---

*본 문서는 IEC 81001-5-1:2021 (Health software cybersecurity) 및 OWASP API Security Top 10 (2023)을 준수합니다.*
