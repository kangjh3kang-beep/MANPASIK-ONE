# 만파식 보안 아키텍처 STRIDE 위협 모델링

> **문서 ID**: MPS-SEC-STM-001
> **표준 근거**: IEC 81001-5-1:2021, FDA Cybersecurity Guidance (2023), STRIDE (Microsoft)
> **작성일**: 2026-02-09
> **작성자**: Claude (Security & Architecture Agent)
> **검토 상태**: 초안 (Security Review 필요)

---

## 1. 위협 모델링 범위

### 1.1 대상 시스템

```
┌─────────────┐    BLE 5.0     ┌──────────────┐   HTTPS/gRPC   ┌──────────────┐
│  ManPaSik   │◄──────────────►│  Flutter App  │◄─────────────►│  API Gateway │
│  리더기     │                │  (모바일)     │                │  (Kong)      │
│  + NFC      │                │              │                │              │
│  카트리지   │                │  Rust Core   │                │  Keycloak    │
└─────────────┘                │  (FFI)       │                └──────┬───────┘
                               └──────────────┘                       │ gRPC+mTLS
                                                                ┌─────┼─────────────┐
                                                                │     ▼             │
                                                                │  Go 마이크로서비스  │
                                                                │  (auth/user/      │
                                                                │   device/measure) │
                                                                │     │             │
                                                                │     ▼             │
                                                                │  데이터 계층      │
                                                                │  (PG/TS/Milvus/  │
                                                                │   Redis/Kafka)    │
                                                                └───────────────────┘
```

### 1.2 공격 표면 (Attack Surface)

| # | 공격 표면 | 프로토콜 | 데이터 유형 | 노출도 |
|---|----------|---------|-----------|--------|
| AS1 | BLE 통신 (리더기 ↔ 앱) | BLE 5.0 GATT | 측정 원시 데이터, 명령 | 물리적 근거리 |
| AS2 | NFC 통신 (카트리지 ↔ 리더기) | ISO 14443A | 보정 데이터, 카트리지 ID | 물리적 접촉 |
| AS3 | HTTPS API (앱 ↔ Kong Gateway) | TLS 1.3 | 인증 토큰, 건강 데이터 | 인터넷 |
| AS4 | gRPC 내부 (Gateway ↔ Services) | gRPC + mTLS | 서비스 간 데이터 | 내부 네트워크 |
| AS5 | MQTT (IoT 허브 → 서버) | MQTT 3.1.1 | 디바이스 상태, 텔레메트리 | 인터넷 |
| AS6 | 로컬 저장소 (앱 내부) | SQLite/Hive | 오프라인 측정 데이터 | 디바이스 물리 접근 |
| AS7 | 데이터베이스 (내부) | TCP | PHI, PII, 벡터 데이터 | 내부 네트워크 |
| AS8 | 관리자 웹 (Next.js) | HTTPS | 관리 기능, 대시보드 | 인터넷 (제한) |

---

## 2. STRIDE 분석

### STRIDE 카테고리 설명

| 카테고리 | 설명 | 보안 속성 |
|---------|------|---------|
| **S**poofing | 신원 위조/위장 | 인증 (Authentication) |
| **T**ampering | 데이터 변조 | 무결성 (Integrity) |
| **R**epudiation | 부인 | 부인 방지 (Non-repudiation) |
| **I**nformation Disclosure | 정보 유출 | 기밀성 (Confidentiality) |
| **D**enial of Service | 서비스 거부 | 가용성 (Availability) |
| **E**levation of Privilege | 권한 상승 | 인가 (Authorization) |

---

### 2.1 AS1: BLE 통신 (리더기 ↔ 앱)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-BLE-01 | **위조 리더기 연결**: 공격자가 가짜 BLE 기기로 위장하여 앱에 연결 | S | 높음 | 중간 | 🔴 높음 | BLE Secure Connections (ECDH 키교환), 디바이스 인증서 검증, MAC 주소 화이트리스트 | 🟡 설계됨, 미구현 |
| T-BLE-02 | **MITM 공격**: BLE 통신 중간에서 데이터 가로채기/변조 | T,I | 높음 | 낮음 | 🟡 중간 | BLE 5.0 LE Secure Connections (AES-CCM), 페어링 시 OOB/Numeric Comparison | 🟡 설계됨, 미구현 |
| T-BLE-03 | **측정 데이터 변조**: 전송 중 측정값 조작 | T | 높음 | 낮음 | 🟡 중간 | AES-CCM 암호화 + HMAC 무결성 검증, 해시체인 | ❌ crypto 모듈 미구현 |
| T-BLE-04 | **BLE 재전송 공격**: 이전 측정 패킷 재전송 | S,T | 중간 | 낮음 | 🟡 중간 | 시퀀스 번호 (패킷 헤더 u16), 타임스탬프 검증, Nonce 사용 | 🟡 시퀀스 번호 구현 |
| T-BLE-05 | **BLE 서비스 거부**: 다수 BLE 연결 시도로 리더기 마비 | D | 중간 | 중간 | 🟡 중간 | 동시 연결 수 제한, 연결 타임아웃, 레이트 리밋 | ❌ 미구현 |
| T-BLE-06 | **BLE 스니핑**: 페어링 과정 도청으로 키 탈취 | I | 높음 | 낮음 | 🟡 중간 | LE Secure Connections (P-256 ECDH), Just Works 금지 | 🟡 설계됨 |

### 2.2 AS2: NFC 카트리지 통신

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-NFC-01 | **위조 카트리지**: 가짜/복제 NFC 태그 사용 | S | 높음 | 중간 | 🔴 높음 | MIFARE DESFire EV3 인증, 카트리지 ID 서버 검증, 사용 횟수 차감 서버 동기화 | 🟡 검증 로직 구현, 서버 동기화 미구현 |
| T-NFC-02 | **보정 데이터 변조**: NFC 태그 내 보정 계수(α, offset, gain) 조작 | T | 높음 | 낮음 | 🟡 중간 | 보정 데이터 HMAC 서명, 서버 원본 대조, DESFire 쓰기 보호 | ❌ HMAC 미구현 |
| T-NFC-03 | **만료 카트리지 재사용**: 유효기간/사용횟수 조작 | T | 중간 | 중간 | 🟡 중간 | 서버 기반 사용 이력 관리, NFC 카운터 감소 원자적 연산 | 🟡 validate_cartridge() 구현 |
| T-NFC-04 | **NFC 데이터 도청**: 태그 통신 스니핑 | I | 낮음 | 낮음 | 🟢 낮음 | DESFire 암호화된 통신, 민감정보 미포함 (보정값만) | 🟡 설계됨 |

### 2.3 AS3: HTTPS API (앱 ↔ Kong Gateway)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-API-01 | **JWT 토큰 탈취**: XSS/피싱으로 Access Token 탈취 | S,I | 높음 | 중간 | 🔴 높음 | Access Token TTL 15분, HttpOnly 쿠키, Secure 플래그, SameSite=Strict | 🟡 Keycloak 설계 |
| T-API-02 | **Refresh Token 탈취**: 기기 분실/악성코드로 Refresh Token 탈취 | S | 높음 | 중간 | 🔴 높음 | 디바이스 바인딩, Refresh Token Rotation, Redis 블랙리스트, MFA | 🟡 설계됨 |
| T-API-03 | **API 인젝션 공격**: SQL Injection, NoSQL Injection, Command Injection | T | 높음 | 중간 | 🔴 높음 | ORM만 사용 (sqlc/GORM), 입력 검증 (protobuf 스키마), WAF (Kong 플러그인) | 🟡 설계됨 (ORM 규칙) |
| T-API-04 | **DDoS/Rate Limit 우회**: API 과부하 공격 | D | 중간 | 높음 | 🔴 높음 | Kong Rate Limiting 플러그인, IP별/사용자별/전역 제한, Cloudflare WAF | 🟡 Kong 설정 존재 |
| T-API-05 | **IDOR (Insecure Direct Object Reference)**: 다른 사용자 데이터 접근 | I,E | 높음 | 중간 | 🔴 높음 | JWT claims에서 user_id 추출 (요청 파라미터 무시), RBAC 검증 | ❌ 미구현 |
| T-API-06 | **SSRF (Server-Side Request Forgery)**: 내부 서비스 접근 | I | 높음 | 낮음 | 🟡 중간 | 내부 네트워크 분리, 허용 URL 화이트리스트 | 🟡 Docker 네트워크 분리 |
| T-API-07 | **Man-in-the-Middle**: TLS 우회/다운그레이드 | I,T | 높음 | 낮음 | 🟡 중간 | TLS 1.3 강제, Certificate Pinning (앱), HSTS | 🟡 Kong TLS 설정 |
| T-API-08 | **부적절한 에러 노출**: 스택 트레이스/내부 정보 응답 포함 | I | 중간 | 중간 | 🟡 중간 | 표준 에러 응답 형식, 프로덕션 디버그 모드 비활성화, 에러 코드 체계 | 🟡 규칙에 정의됨 |

### 2.4 AS4: gRPC 내부 통신 (Gateway ↔ Services)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-GRPC-01 | **서비스 위장**: 악성 서비스가 정상 서비스로 위장 | S | 높음 | 낮음 | 🟡 중간 | mTLS (상호 인증), 서비스 메시 (Istio/Linkerd), K8s NetworkPolicy | 🟡 설계됨 |
| T-GRPC-02 | **내부 통신 도청**: 서비스 간 평문 통신 가로채기 | I | 높음 | 낮음 | 🟡 중간 | mTLS 필수, K8s Pod Security Policy | 🟡 설계됨 |
| T-GRPC-03 | **Protobuf 메시지 변조**: 직렬화/역직렬화 과정에서 조작 | T | 중간 | 낮음 | 🟢 낮음 | 스키마 검증 (protobuf), 서명 검증, TLS 무결성 | ✅ protobuf 사용 |
| T-GRPC-04 | **내부 서비스 DoS**: 대량 gRPC 호출로 서비스 마비 | D | 중간 | 중간 | 🟡 중간 | 서비스별 연결 제한, Circuit Breaker, K8s 리소스 리밋 | ❌ 미구현 |

### 2.5 AS5: MQTT (IoT)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-MQTT-01 | **위조 디바이스 메시지**: 가짜 디바이스 데이터 발행 | S,T | 중간 | 중간 | 🟡 중간 | MQTT 인증 (사용자/비밀번호 + 인증서), Topic ACL, 디바이스 인증서 | ❌ Mosquitto 기본 설정 |
| T-MQTT-02 | **MQTT 도청**: 텔레메트리 데이터 도청 | I | 중간 | 중간 | 🟡 중간 | MQTT over TLS, 토픽별 접근 제어 | ❌ 미구현 |
| T-MQTT-03 | **MQTT Flood**: 대량 메시지로 브로커 마비 | D | 중간 | 중간 | 🟡 중간 | 메시지 크기 제한, 클라이언트별 레이트 리밋 | ❌ 미구현 |

### 2.6 AS6: 로컬 저장소 (모바일 앱)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-LOCAL-01 | **루팅/탈옥 기기에서 데이터 추출** | I | 높음 | 중간 | 🔴 높음 | SQLCipher 암호화 DB, 키체인 저장, 루팅 탐지 | ❌ 미구현 |
| T-LOCAL-02 | **앱 백업에서 PHI 노출** | I | 높음 | 중간 | 🟡 중간 | android:allowBackup=false, iOS Keychain exclusion | ❌ 미구현 |
| T-LOCAL-03 | **오프라인 데이터 변조** | T | 중간 | 낮음 | 🟡 중간 | SHA-256 해시체인 + HMAC, 동기화 시 서버 검증 | 🟡 해시체인 설계됨 |

### 2.7 AS7: 데이터베이스

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-DB-01 | **DB 직접 접근**: 네트워크 침투 후 DB 직접 연결 | I | 높음 | 낮음 | 🟡 중간 | K8s NetworkPolicy, DB 포트 외부 비노출, 강력한 비밀번호, SSL 연결 | 🟡 Docker 네트워크 분리 |
| T-DB-02 | **PHI 평문 저장**: 암호화 없이 건강 데이터 저장 | I | 높음 | 확정 | 🔴 높음 | 컬럼 레벨 암호화 (pgcrypto), TDE, 앱 레벨 암호화 | ❌ 미구현 |
| T-DB-03 | **백업 데이터 유출**: 백업 파일에서 데이터 추출 | I | 높음 | 낮음 | 🟡 중간 | 암호화된 백업, 접근 제어, 안전한 폐기 | ❌ 미구현 |
| T-DB-04 | **감사 로그 변조**: 접근 기록 삭제/수정 | R | 높음 | 낮음 | 🟡 중간 | 감사 로그 별도 저장, 불변 스토리지, 로그 서명 | ❌ 감사 로그 미구현 |

### 2.8 AS8: 관리자 웹 (Next.js)

| ID | 위협 | STRIDE | 심각도 | 가능성 | 위험도 | 완화 방안 | 현재 상태 |
|----|------|--------|--------|--------|--------|---------|---------|
| T-WEB-01 | **XSS (Cross-Site Scripting)** | I,S | 높음 | 중간 | 🔴 높음 | React 자동 이스케이프, CSP 헤더, DOMPurify | ❌ 웹 미구현 |
| T-WEB-02 | **CSRF (Cross-Site Request Forgery)** | T | 중간 | 중간 | 🟡 중간 | CSRF 토큰, SameSite 쿠키, Origin 검증 | ❌ 웹 미구현 |
| T-WEB-03 | **관리자 계정 탈취** | S,E | 높음 | 중간 | 🔴 높음 | MFA 필수, IP 화이트리스트, 세션 타임아웃 | 🟡 Keycloak MFA 설계 |

---

## 3. 위협 우선순위 매트릭스

### 3.1 위험도별 분류

#### 🔴 높음 (즉시 조치 필요) — 12건

| ID | 위협 | 공격 표면 | 핵심 조치 |
|----|------|---------|---------|
| T-BLE-01 | 위조 리더기 연결 | BLE | BLE Secure Connections 구현 |
| T-NFC-01 | 위조 카트리지 | NFC | 서버 기반 카트리지 인증 |
| T-API-01 | JWT 토큰 탈취 | HTTPS | Token TTL 단축, HttpOnly |
| T-API-02 | Refresh Token 탈취 | HTTPS | Token Rotation, 디바이스 바인딩 |
| T-API-03 | 인젝션 공격 | HTTPS | ORM + 입력 검증 |
| T-API-04 | DDoS | HTTPS | Kong Rate Limiting |
| T-API-05 | IDOR | HTTPS | JWT 기반 사용자 식별 |
| T-LOCAL-01 | 로컬 데이터 추출 | 모바일 | SQLCipher 암호화 |
| T-DB-02 | PHI 평문 저장 | DB | 컬럼 레벨 암호화 |
| T-WEB-01 | XSS | 웹 | CSP + 자동 이스케이프 |
| T-WEB-03 | 관리자 계정 탈취 | 웹 | MFA + IP 제한 |
| T-BLE-03 | 측정 데이터 변조 | BLE | AES-CCM + HMAC |

#### 🟡 중간 (Phase 1 내 조치) — 17건

- T-BLE-02, T-BLE-04, T-BLE-05, T-BLE-06
- T-NFC-02, T-NFC-03
- T-API-06, T-API-07, T-API-08
- T-GRPC-01, T-GRPC-02, T-GRPC-04
- T-MQTT-01, T-MQTT-02, T-MQTT-03
- T-LOCAL-02, T-LOCAL-03
- T-DB-01, T-DB-03, T-DB-04
- T-WEB-02

#### 🟢 낮음 (Phase 2 조치 가능) — 2건

- T-NFC-04, T-GRPC-03

---

## 4. 보안 아키텍처 권고 사항

### 4.1 인증/인가 아키텍처

```
┌──────────────────────────────────────────────────────────────┐
│ 인증 흐름 (Authentication Flow)                                │
│                                                                │
│  사용자 → Keycloak (OIDC) → Access Token (JWT, 15분)          │
│                            → Refresh Token (Opaque, 7일)       │
│                            → ID Token (사용자 정보)             │
│                                                                │
│  MFA 정책:                                                     │
│  - Free 티어: 선택적 MFA                                       │
│  - Basic 이상: 필수 MFA (TOTP/SMS)                             │
│  - Clinical 티어: 필수 MFA + 생체인식 권고                      │
│                                                                │
│  세션 정책:                                                     │
│  - 동시 세션 제한: 최대 3 디바이스                              │
│  - 유휴 타임아웃: 30분                                         │
│  - 절대 타임아웃: 12시간                                       │
│  - 민감 작업 재인증: 비밀번호 변경, 결제, 데이터 삭제          │
└──────────────────────────────────────────────────────────────┘
```

### 4.2 RBAC 접근제어 매트릭스

| 리소스 \ 역할 | 일반사용자 | 가족관리자 | 의료진 | 관리자 | 시스템 |
|-------------|----------|----------|--------|--------|--------|
| 자기 측정 데이터 | CRUD | R (가족) | R (환자) | R | CRUD |
| 타인 측정 데이터 | - | R (가족) | R (환자) | R | CRUD |
| 사용자 프로필 | CRU(자기) | CRU(자기) | R(환자) | CRUD | CRUD |
| 디바이스 관리 | CRU(자기) | CRU(자기) | R | CRUD | CRUD |
| 구독 관리 | RU(자기) | RU(자기) | - | CRUD | CRUD |
| 관리자 대시보드 | - | - | R(제한) | CRUD | CRUD |
| 감사 로그 | - | - | - | R | CRUD |
| 시스템 설정 | - | - | - | RU | CRUD |

### 4.3 암호화 전략

| 데이터 상태 | 암호화 방식 | 키 관리 | 구현 상태 |
|-----------|-----------|--------|---------|
| **저장 (At Rest)** | AES-256-GCM | AWS KMS / Vault | ❌ 미구현 |
| **전송 (In Transit)** | TLS 1.3 (HTTPS/gRPC) | 자동 인증서 (cert-manager) | 🟡 설계 |
| **클라이언트 (Local)** | AES-256 (SQLCipher) | 디바이스 키체인/키스토어 | ❌ 미구현 |
| **BLE 통신** | AES-CCM (BLE SC) | ECDH P-256 키교환 | ❌ 미구현 |
| **NFC 태그** | DESFire EV3 암호화 | 카트리지 제조 시 프로비저닝 | ❌ 미구현 |
| **백업** | AES-256-GCM | 별도 백업 키 | ❌ 미구현 |

### 4.4 키 관리 정책

| 키 유형 | 생성 | 로테이션 | 폐기 | 저장 |
|--------|------|---------|------|------|
| JWT 서명 키 (RSA 2048) | Keycloak 자동 | 90일 | 폐기 후 1년 보관 | Keycloak 내부 |
| DB 암호화 키 (AES-256) | KMS/Vault | 365일 | 데이터 재암호화 후 | KMS/Vault |
| BLE 세션 키 (AES-CCM) | ECDH 키교환 | 매 연결 | 연결 종료 시 | 메모리 (비영속) |
| NFC 인증 키 (DESFire) | 제조 프로비저닝 | 불변 | 카트리지 수명 | Secure Element |
| TLS 인증서 | cert-manager | 90일 (Let's Encrypt) | 자동 | K8s Secret |
| API 키 | 관리자 생성 | 365일 | 즉시 무효화 | Vault |

---

## 5. 침해사고 대응 계획 (Incident Response Plan)

### 5.1 사고 분류

| 등급 | 정의 | 응답 시간 | 보고 대상 |
|------|------|---------|---------|
| **P1 (Critical)** | PHI 유출, 시스템 전면 장애, 랜섬웨어 | 15분 내 탐지, 1시간 내 대응 | CEO, CISO, 규제당국 |
| **P2 (High)** | 대규모 서비스 장애, 인증 시스템 침해 | 30분 내 탐지, 4시간 내 대응 | CTO, CISO |
| **P3 (Medium)** | 개별 서비스 장애, 의심스러운 접근 | 2시간 내 탐지, 24시간 내 대응 | 보안팀 |
| **P4 (Low)** | 정책 위반, 취약점 발견 | 24시간 내 탐지, 1주 내 대응 | 보안팀 |

### 5.2 대응 절차 (7단계)

```
1. 탐지 (Detection)
   → Prometheus 알림, ELK 이상 탐지, WAF 로그, 사용자 보고

2. 분류 (Triage)
   → 사고 등급 판정, 영향 범위 평가, 에스컬레이션 결정

3. 격리 (Containment)
   → 영향받은 서비스 격리, 네트워크 분리, 계정 비활성화

4. 근절 (Eradication)
   → 공격 벡터 제거, 취약점 패치, 악성코드 제거

5. 복구 (Recovery)
   → 서비스 복원, 데이터 무결성 검증, 모니터링 강화

6. 보고 (Notification)
   → 규제당국 보고 (GDPR 72시간, HIPAA 60일, PIPA 72시간)
   → 영향받은 사용자 통지

7. 사후 분석 (Post-Incident)
   → RCA (Root Cause Analysis), 재발 방지 대책, 프로세스 개선
```

### 5.3 규제 보고 요건

| 규정 | 보고 기한 | 보고 대상 | 조건 |
|------|---------|---------|------|
| GDPR | 72시간 | 감독기관 + 정보주체 | 개인정보 유출 |
| HIPAA | 60일 | HHS OCR + 환자 | 500명 이상 PHI 유출 |
| PIPA | 72시간 | 개인정보보호위원회 + 정보주체 | 개인정보 유출 |
| PIPL | 즉시 | CAC | 개인정보 유출 |
| MFDS | 15일 | 식약처 | 의료기기 안전성 문제 |
| FDA | 30일 (MDR) | FDA MAUDE | 심각한 부작용 |

---

## 6. SBOM (Software Bill of Materials) 요구사항

### 6.1 형식

- **표준**: CycloneDX 1.5 (JSON/XML)
- **대안**: SPDX 2.3
- **생성 도구**: `cargo-cyclonedx` (Rust), `cyclonedx-gomod` (Go), `flutter pub deps` + 변환

### 6.2 포함 항목

| 항목 | 설명 |
|------|------|
| 컴포넌트 이름 | 라이브러리/패키지 이름 |
| 버전 | 정확한 버전 |
| 라이선스 | SPDX 식별자 |
| 공급자 | 패키지 저장소 |
| 해시 | SHA-256 |
| 취약점 | 알려진 CVE |

### 6.3 취약점 스캔 도구

| 언어 | 도구 | 주기 |
|------|------|------|
| Rust | `cargo-audit`, `cargo-deny` | 매 빌드 |
| Go | `govulncheck`, `nancy` | 매 빌드 |
| Flutter/Dart | `dart pub outdated`, OSV Scanner | 주 1회 |
| Docker | Trivy, Snyk Container | 매 빌드 |
| 전체 | Dependabot / Renovate | 자동 PR |

---

## 7. 다음 단계

### 즉시 (이번 스프린트)
1. Rust `crypto` 모듈 AES-256-GCM 구현 → T-BLE-03, T-DB-02 해소
2. SBOM 첫 생성 (cargo-cyclonedx + cyclonedx-gomod)
3. Kong Rate Limiting 플러그인 활성화 → T-API-04 해소

### Phase 1 내
4. BLE Secure Connections 구현 → T-BLE-01, T-BLE-02 해소
5. Go 서비스 JWT 인터셉터 구현 → T-API-01, T-API-05 해소
6. MQTT TLS + 인증 설정 → T-MQTT-01, T-MQTT-02 해소
7. DB 컬럼 레벨 암호화 설계 → T-DB-02 해소

### Phase 2
8. 침입 탐지 시스템 (IDS) 구축
9. 보안 자동화 (SAST/DAST) 파이프라인
10. 침투 테스트 (외부 업체)

---

**Document Version**: 1.0.0
**Next Review**: 2026-02-16
**Approval Required**: CISO, CTO
