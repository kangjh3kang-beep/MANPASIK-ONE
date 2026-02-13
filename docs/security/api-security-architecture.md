# ManPaSik API 보안 아키텍처 명세서 (API Security Architecture)

**문서번호**: MPK-SEC-ARCH-v1.0  
**작성일**: 2026-02-14  
**목적**: 인증 흐름, API 보안 계층, 키 관리, 서비스 간 인증, PHI 접근 제어, 보안 모니터링을 상세 정의  
**적용**: 전체 시스템 (Kong Gateway, Keycloak, Go 백엔드, Rust 코어, Flutter 클라이언트)

---

## 1. 보안 아키텍처 개요

### 1.1 보안 계층 구조

```text
[클라이언트 - Flutter/Web]
    │ TLS 1.3
    ▼
[Kong API Gateway :8090]
    ├── Rate Limiting (IP/User)
    ├── CORS / CSRF
    ├── Request Size Limit (10MB)
    ├── Bot Detection
    │ JWT 검증 (Keycloak Public Key)
    ▼
[Keycloak :8080]
    ├── OAuth 2.0 / OpenID Connect
    ├── Token 발급 (Access: 15m / Refresh: 7d)
    ├── MFA (TOTP, WebAuthn)
    ├── 소셜 로그인 (Google, Apple, Kakao)
    │ JWT (RS256)
    ▼
[gRPC 서비스 계층]
    ├── JWT Interceptor (토큰 검증)
    ├── RBAC Middleware (역할 기반 접근 제어)
    ├── PHI Filter (데이터 필터링)
    ├── Audit Interceptor (감사 로그)
    │ mTLS (서비스 간)
    ▼
[데이터 계층]
    ├── PostgreSQL (TDE - Transparent Data Encryption)
    ├── Redis (ACL + TLS)
    ├── Milvus (인증 + 암호화)
    └── MinIO (Server-Side Encryption)
```

---

## 2. 인증 흐름 상세

### 2.1 로그인 인증 흐름

```text
┌──────────┐      ┌──────────┐      ┌──────────┐      ┌──────────┐
│ Flutter  │      │  Kong    │      │ Keycloak │      │  auth-   │
│   App    │      │ Gateway  │      │  (IdP)   │      │ service  │
└────┬─────┘      └────┬─────┘      └────┬─────┘      └────┬─────┘
     │                  │                  │                  │
     │ 1. POST /auth/login                │                  │
     │ {email, password}│                  │                  │
     │─────────────────>│                  │                  │
     │                  │                  │                  │
     │                  │ 2. Rate Check    │                  │
     │                  │ (5회/분/IP)      │                  │
     │                  │                  │                  │
     │                  │ 3. gRPC Login    │                  │
     │                  │─────────────────────────────────────>
     │                  │                  │                  │
     │                  │                  │  4. Verify Password
     │                  │                  │  (Argon2id)       │
     │                  │                  │                  │
     │                  │                  │  5. MFA 필요 여부  │
     │                  │                  │  확인              │
     │                  │                  │                  │
     │                  │                  │  6a. MFA 불필요 시:│
     │                  │                  │<─── Token 요청 ───│
     │                  │                  │                  │
     │                  │                  │  7. JWT 발급      │
     │                  │                  │  Access(15m) +    │
     │                  │                  │  Refresh(7d)      │
     │                  │                  │                  │
     │                  │  8. {access_token, refresh_token}   │
     │<─────────────────│<─────────────────────────────────────
     │                  │                  │                  │
     │  9. 토큰 저장     │                  │                  │
     │  (Secure Storage) │                  │                  │
```

### 2.2 JWT 토큰 구조

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT",
    "kid": "manpasik-key-2026"
  },
  "payload": {
    "iss": "https://auth.manpasik.com/realms/manpasik",
    "sub": "uuid-user-id",
    "aud": "manpasik-app",
    "exp": 1739502000,
    "iat": 1739501100,
    "jti": "uuid-token-id",
    "roles": ["user", "subscriber:pro"],
    "tier": "pro",
    "device_ids": ["uuid-device-1", "uuid-device-2"],
    "mfa_verified": true,
    "scope": "openid profile health:read health:write"
  }
}
```

### 2.3 토큰 갱신 흐름

```text
Access Token 만료 시 (15분):
  1. Flutter → Kong: POST /auth/refresh {refresh_token}
  2. Kong → auth-service: gRPC RefreshToken(refresh_token)
  3. auth-service → Redis: Refresh Token 유효성 확인
  4. auth-service → Keycloak: 새 Access Token 요청
  5. 새 Access Token 발급 (이전 Refresh Token 무효화 — Token Rotation)
  6. 응답: {new_access_token, new_refresh_token}

Refresh Token 만료 시 (7일):
  → 재로그인 필요
  → 오프라인 모드 전환 옵션 (Rust CRDT 로컬 동작)
```

---

## 3. API 보안 계층 상세

### 3.1 Kong Gateway 보안 플러그인

| 플러그인 | 설정 | 적용 범위 |
|---------|------|----------|
| **Rate Limiting** | 인증 API: 5회/분/IP, 일반 API: 100회/분/사용자 | 글로벌 |
| **Request Size Limiting** | 10MB (파일 업로드: 50MB) | 글로벌 |
| **CORS** | origin: `*.manpasik.com`, methods: `GET,POST,PUT,DELETE` | 글로벌 |
| **IP Restriction** | 관리자 API: 허용 IP 화이트리스트 | /admin/* |
| **Bot Detection** | User-Agent 필터 + reCAPTCHA (로그인) | /auth/* |
| **Request Transformer** | 민감 헤더 제거, X-Request-ID 추가 | 글로벌 |
| **Response Transformer** | 서버 정보 헤더 제거 (Server, X-Powered-By) | 글로벌 |
| **OpenTelemetry** | 모든 요청 트레이스 ID 전파 | 글로벌 |

### 3.2 gRPC 보안 Interceptor 체인

```go
// 서버 측 Interceptor 순서 (순차 실행)
grpc.ChainUnaryInterceptor(
    // 1. 메트릭 수집 (Prometheus)
    otelgrpc.UnaryServerInterceptor(),
    
    // 2. 에러 복구 (panic → gRPC error)
    recovery.UnaryServerInterceptor(),
    
    // 3. JWT 검증 (모든 요청)
    auth.UnaryServerInterceptor(jwtAuthFunc),
    
    // 4. RBAC 역할 검사 (RPC별 필요 역할 확인)
    rbac.UnaryServerInterceptor(rbacConfig),
    
    // 5. 요청 검증 (protobuf 필드 유효성)
    validator.UnaryServerInterceptor(),
    
    // 6. 감사 로그 (PHI 접근 기록)
    audit.UnaryServerInterceptor(auditConfig),
    
    // 7. Rate Limiting (사용자별)
    ratelimit.UnaryServerInterceptor(limiter),
)
```

### 3.3 RBAC 역할 체계

| 역할 | 코드 | 접근 범위 |
|------|------|----------|
| **guest** | `role:guest` | 앱 다운로드, 회원가입, 공개 콘텐츠 |
| **user** | `role:user` | 자기 데이터 CRUD, Free 티어 기능 |
| **subscriber:basic** | `role:subscriber:basic` | Basic 카트리지, AI 코칭 기본 |
| **subscriber:pro** | `role:subscriber:pro` | Pro 카트리지, 고급 분석, 화상진료 |
| **subscriber:clinical** | `role:subscriber:clinical` | Clinical 카트리지, 전문가 기능 |
| **clinician** | `role:clinician` | 환자 데이터 열람(동의), 처방 작성 |
| **admin** | `role:admin` | 시스템 설정, 사용자 관리, 감사 로그 |
| **super_admin** | `role:super_admin` | 전체 접근, 역할 부여, 시스템 제어 |

**RPC별 역할 매핑 (예시):**

| RPC | 필요 역할 | 데이터 필터 |
|-----|----------|-----------|
| `AuthService.Login` | (인증 불필요) | — |
| `UserService.GetProfile` | `user+` | `user_id == ctx.user_id` |
| `MeasurementService.GetHistory` | `user+` | `user_id == ctx.user_id` 또는 가족 공유 |
| `AdminService.ListSystemConfigs` | `admin+` | — |
| `PrescriptionService.CreatePrescription` | `clinician+` | 담당 환자만 |
| `SubscriptionService.CancelSubscription` | `user+` | `user_id == ctx.user_id` |

---

## 4. 키 관리 (Key Management)

### 4.1 암호화 키 체계

| 키 유형 | 알고리즘 | 용도 | 저장 위치 | 순환 주기 |
|---------|---------|------|----------|----------|
| **JWT 서명 키** | RSA-2048 | Keycloak JWT 서명 | Keycloak DB (암호화) | 90일 |
| **데이터 암호화 키 (DEK)** | AES-256-GCM | PHI 필드 암호화 | K8s Secret (Sealed) | 365일 |
| **키 암호화 키 (KEK)** | AES-256 | DEK 암호화 | HSM 또는 Vault | 2년 |
| **mTLS 인증서** | ECDSA P-256 | 서비스 간 인증 | K8s Secret + cert-manager | 90일 (자동) |
| **NFC 프로비저닝 키** | HMAC-SHA256 | 카트리지 인증 | Rust 코어 (하드코딩→HSM) | 1년 |
| **Rust 벡터 암호화** | AES-256-GCM | 핑거프린트 로컬 저장 | Secure Enclave / TEE | 앱 설치 시 |

### 4.2 키 순환 절차

```text
DEK 순환 (365일 주기):
  1. 새 DEK 생성 (crypto/rand 256-bit)
  2. KEK로 새 DEK 암호화 → K8s Secret 업데이트
  3. 기존 DEK → "이전 키" 슬롯 보관 (복호화용)
  4. 새 데이터 → 새 DEK로 암호화
  5. 기존 데이터 → 배치 재암호화 (야간 Job)
  6. 완료 확인 → 이전 키 폐기

JWT 키 순환 (90일 주기):
  1. Keycloak에 새 RSA 키쌍 생성
  2. 새 키로 신규 토큰 서명
  3. 이전 키 → "검증 전용"으로 전환 (72시간)
  4. 72시간 후 이전 키 비활성화
```

---

## 5. PHI(개인건강정보) 보호 상세

### 5.1 PHI 필드 분류

| 데이터 | 민감도 | 암호화 | 접근 제한 | 감사 |
|--------|--------|--------|----------|------|
| 측정 결과 (수치) | **높음** | AES-256-GCM (필드 레벨) | 본인 + 의료진(동의) | 모든 접근 기록 |
| 핑거프린트 벡터 | **높음** | AES-256-GCM | 본인 + AI 서비스 | 모든 접근 기록 |
| 건강 기록 (FHIR) | **높음** | AES-256-GCM | 본인 + 의료진(동의) | 모든 접근 기록 |
| 처방전 | **높음** | AES-256-GCM | 본인 + 처방 의료진 | 모든 접근 기록 |
| 프로필 (이름, 생년) | **중간** | DB TDE | 본인 + 관리자(사유) | 비정상 접근 기록 |
| 기기 정보 | **낮음** | DB TDE | 본인 + 서비스 | — |
| 설정/선호 | **낮음** | 평문 | 본인 | — |

### 5.2 동의 기반 데이터 공유

```text
데이터 공유 흐름:
  1. 의료진이 환자에게 공유 요청
  2. 환자 앱에 "데이터 공유 동의" 팝업
     - 공유 범위 선택 (측정 결과 / 건강 기록 / 전체)
     - 공유 기간 (1회성 / 30일 / 진료 종료까지)
  3. 동의 → data_sharing_consents 테이블 기록
  4. 의료진 → GetSharedHealthData(patient_id) → 동의 범위 내 데이터만 반환
  5. 공유 종료 → 동의 만료 / 환자 철회 → 접근 차단
```

---

## 6. 보안 이벤트 모니터링

### 6.1 보안 이벤트 유형 및 대응

| 이벤트 | 임계값 | 심각도 | 자동 대응 | 알림 |
|--------|--------|--------|----------|------|
| 로그인 실패 (동일 IP) | 5회/5분 | **높음** | IP 15분 차단 | Admin Slack |
| 로그인 실패 (동일 계정) | 3회/5분 | **높음** | 계정 30분 잠금 | 사용자 이메일 + Admin |
| 비정상 토큰 사용 | 1회 | **긴급** | 세션 즉시 무효화 | Admin Slack + PagerDuty |
| PHI 대량 접근 | 100건/1분 | **긴급** | 접근 차단 + 세션 종료 | Admin + 보안팀 |
| Admin API 비인가 접근 | 1회 | **높음** | IP 차단 | Admin Slack |
| SQL/NoSQL Injection 시도 | 1회 | **높음** | 요청 차단 + IP 기록 | Admin |
| 인증서 만료 임박 | 7일 전 | **중간** | cert-manager 자동 갱신 | Infra 이메일 |
| Rate Limit 초과 | 설정별 | **낮음** | 429 반환 | — |

### 6.2 감사 로그 구조

```json
{
  "audit_id": "uuid",
  "timestamp": "2026-02-14T05:50:00Z",
  "user_id": "uuid",
  "action": "READ_HEALTH_RECORD",
  "resource_type": "health_record",
  "resource_id": "uuid",
  "service": "health-record-service",
  "ip_address": "192.168.1.100",
  "user_agent": "ManPaSik/1.0 (Flutter)",
  "result": "SUCCESS",
  "phi_accessed": true,
  "data_fields": ["measurement_results", "diagnosis"],
  "consent_id": "uuid-of-consent",
  "retention_until": "2036-02-14T00:00:00Z"
}
```

**보존 기간**: 10년 (IEC 62304, FDA 21 CFR Part 11 준수)  
**저장소**: Elasticsearch (검색) + S3 아카이브 (장기 보존)

---

**참조**: `docs/security/stride-threat-model.md`, `docs/compliance/data-protection-policy.md`, `docs/specs/non-functional-requirements.md` §4
