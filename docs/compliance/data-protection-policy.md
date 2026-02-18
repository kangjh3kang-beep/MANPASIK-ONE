# 만파식 데이터 보호 정책

> **문서 ID**: MPS-SEC-DPP-001
> **적용 규정**: HIPAA, GDPR, PIPA, PIPL, APPI
> **작성일**: 2026-02-09
> **작성자**: Claude (Security & Architecture Agent)
> **검토 상태**: 초안 (Legal Review 필요)
> **적용 범위**: ManPaSik 플랫폼 전체 (앱, 서버, 리더기)

---

## 1. 목적 및 범위

### 1.1 목적

본 정책은 만파식 헬스케어 AI 생태계에서 처리되는 모든 개인정보 및 건강정보의 수집, 이용, 저장, 전송, 삭제에 관한 기준을 수립하여 5개국 규정을 동시에 준수합니다.

### 1.2 적용 범위

- **대상 시스템**: ManPaSik 모바일 앱, 웹 대시보드, Go 마이크로서비스, 데이터베이스, 리더기 펌웨어
- **대상 데이터**: 개인식별정보(PII), 보호대상 건강정보(PHI), 측정 데이터, 디바이스 데이터
- **대상 주체**: 앱 사용자, 가족 구성원, 의료진, 관리자

---

## 2. 데이터 분류 체계

### 2.1 4단계 분류

| 등급 | 분류 | 설명 | 예시 | 보호 수준 |
|------|------|------|------|---------|
| **L4** | 극비 (Critical) | 유출 시 법적 제재 + 심각한 피해 | 건강 측정 원시 데이터, 진단 결과, 바이오마커 수치, 해시체인 | 암호화 필수, 접근 감사, 최소 접근 |
| **L3** | 기밀 (Confidential) | 유출 시 프라이버시 침해 | 이름, 이메일, 생년월일, 의료 기록 ID, 가족 관계 | 암호화 필수, 접근 제어 |
| **L2** | 내부 (Internal) | 내부 운영 데이터 | 디바이스 ID, 펌웨어 버전, 앱 사용 로그, 구독 상태 | 접근 제어, 로깅 |
| **L1** | 공개 (Public) | 공개 가능 정보 | 앱 버전, API 문서, 카트리지 유형 목록 | 기본 보호 |

### 2.2 데이터 인벤토리

| 데이터 항목 | 분류 | 저장소 | 보존 기간 | 암호화 | HIPAA PHI | GDPR 민감 |
|-----------|------|--------|---------|--------|----------|---------|
| 측정 원시 데이터 (채널값) | L4 | TimescaleDB | 10년 | AES-256-GCM | ✅ | ✅ |
| 차동측정 결과 (S_corrected) | L4 | TimescaleDB | 10년 | AES-256-GCM | ✅ | ✅ |
| 핑거프린트 벡터 (88-896차원) | L4 | Milvus | 10년 | AES-256-GCM | ✅ | ✅ |
| AI 추론 결과 (바이오마커 수치) | L4 | PostgreSQL | 10년 | AES-256-GCM | ✅ | ✅ |
| 이름, 이메일, 생년월일 | L3 | PostgreSQL | 계정 삭제+30일 | AES-256-GCM | ✅ | ✅ |
| 가족 관계 정보 | L3 | PostgreSQL | 탈퇴 시 삭제 | AES-256-GCM | - | ✅ |
| 디바이스 ID/시리얼 | L2 | PostgreSQL | 해제+90일 | 해시 | - | - |
| 앱 사용 로그 | L2 | Elasticsearch | 1년 | - | - | - |
| 구독 결제 정보 | L3 | 외부 PG | PG 정책 | PG 관리 | - | - |
| BLE 통신 로그 | L2 | 로컬/서버 | 90일 | - | - | - |
| 감사 추적 로그 | L3 | Elasticsearch | 10년 | 서명 | ✅ | - |

---

## 3. 데이터 수집 원칙

### 3.1 최소 수집 원칙 (Data Minimization)

- 서비스 제공에 **필수적인 데이터만** 수집
- 수집 시 **구체적 목적** 명시
- **목적 외 사용 금지**

### 3.2 수집 항목 및 법적 근거

| 수집 항목 | 수집 목적 | 법적 근거 (GDPR) | 법적 근거 (PIPA) | 필수/선택 |
|----------|---------|----------------|----------------|---------|
| 이메일, 비밀번호 | 계정 생성 및 인증 | 계약 이행 (Art.6(1)(b)) | 동의 (§15) | 필수 |
| 이름, 프로필 사진 | 사용자 식별 | 동의 (Art.6(1)(a)) | 동의 (§15) | 선택 |
| 생년월일, 성별 | 건강 지표 기준값 보정 | 동의 (Art.9(2)(a)) | 민감정보 동의 (§23) | 선택 |
| 건강 측정 데이터 | 핵심 서비스 제공 | 명시적 동의 (Art.9(2)(a)) | 민감정보 동의 (§23) | 필수 |
| 디바이스 정보 | 기기 관리 및 호환성 | 정당한 이익 (Art.6(1)(f)) | 동의 (§15) | 필수 |
| 위치 정보 | 환경 분석 보정 | 동의 (Art.6(1)(a)) | 동의 (§15) + 위치정보법 | 선택 |
| 앱 사용 패턴 | 서비스 개선 | 정당한 이익 (Art.6(1)(f)) | 동의 (§15) | 선택 |

### 3.3 동의 관리 (Consent Management)

```
┌─────────────────────────────────────────────────────┐
│ 동의 관리 UI 요구사항                                   │
├─────────────────────────────────────────────────────┤
│                                                       │
│ 1. 계층적 동의 (Layered Consent)                      │
│    ├── 1차: 필수 동의 (서비스 이용약관, 개인정보 처리)   │
│    ├── 2차: 건강정보 수집 동의 (GDPR Art.9 명시적)      │
│    ├── 3차: 위치정보 수집 동의 (선택)                   │
│    ├── 4차: 마케팅/분석 동의 (선택)                     │
│    └── 5차: AI 학습 데이터 활용 동의 (선택)             │
│                                                       │
│ 2. 동의 속성                                          │
│    - 자유로운 (Freely given): 동의 거부 시 불이익 없음  │
│    - 구체적 (Specific): 각 목적별 개별 동의             │
│    - 정보에 입각한 (Informed): 이해 가능한 언어         │
│    - 명확한 (Unambiguous): 적극적 동의 행위 필요        │
│                                                       │
│ 3. 동의 기록                                          │
│    - 동의 일시, 동의 버전, 동의 내용, 동의 방식          │
│    - 철회 일시 (동의 철회 시)                           │
│    - 최소 5년 보존 (규정 증빙)                         │
│                                                       │
│ 4. 동의 철회                                          │
│    - 설정 > 개인정보 > 동의 관리에서 언제든 철회 가능    │
│    - 철회 즉시 해당 데이터 처리 중단                    │
│    - 필수 동의 철회 = 서비스 탈퇴와 동일                │
│                                                       │
└─────────────────────────────────────────────────────┘
```

---

## 4. 데이터 보호 조치

### 4.1 암호화

| 상태 | 방식 | 대상 | 키 관리 |
|------|------|------|--------|
| 저장 (At Rest) | AES-256-GCM | L3, L4 데이터 | AWS KMS / Vault |
| 전송 (In Transit) | TLS 1.3 | 모든 외부 통신 | cert-manager |
| 클라이언트 (Local) | AES-256 (SQLCipher) | 오프라인 측정 데이터 | 디바이스 키체인 |
| BLE | AES-CCM | 측정 패킷 | ECDH 세션 키 |
| 백업 | AES-256-GCM | 전체 백업 | 별도 백업 키 |

### 4.2 접근 제어

- **인증**: Keycloak OIDC + JWT (RS256) + MFA (Clinical 티어 필수)
- **인가**: RBAC (5개 역할: 사용자, 가족관리자, 의료진, 관리자, 시스템)
- **최소 권한 원칙**: 각 역할에 필요한 최소한의 접근 권한만 부여
- **API 수준 제어**: Kong Gateway에서 엔드포인트별 역할 검증
- **동시 세션 제한**: 사용자 3개 (앱 2+웹 1), 의료진 2개, 관리자 1개 (상세: `api-security-architecture.md §7`)
- **세션 무효화**: 로그아웃 시 즉시 Refresh Token 폐기, 비밀번호 변경 시 전 세션 무효화, 30분 미활동 시 자동 만료

### 4.3 감사 추적 (Audit Trail)

| 기록 항목 | 상세 |
|----------|------|
| 누가 (Who) | 사용자 ID, 역할, IP 주소 |
| 무엇을 (What) | 접근한 리소스, 수행한 작업 (CRUD) |
| 언제 (When) | 타임스탬프 (UTC, 밀리초 정밀도) |
| 어디서 (Where) | 서비스명, 서버 ID |
| 결과 (Result) | 성공/실패, 에러 코드 |

- **보존 기간**: 10년 (IEC 62304 + HIPAA)
- **불변성**: 로그 서명 (SHA-256), 별도 스토리지
- **접근 제한**: 감사 로그 읽기는 관리자/CISO만 가능

### 4.4 데이터 무결성 (해시체인)

```
측정 데이터 패킷:
  hash[n] = SHA-256(data[n] || hash[n-1] || timestamp[n])

검증:
  체인의 어떤 단일 항목이라도 변조되면
  이후 모든 해시가 불일치 → 변조 탐지
```

---

## 5. 정보주체 권리 보장

### 5.1 권리 매트릭스

| 권리 | GDPR | HIPAA | PIPA | PIPL | APPI | 구현 방안 |
|------|------|-------|------|------|------|---------|
| **접근권** (Right of Access) | Art.15 | §164.524 | §35 | Art.44 | Art.33 | 설정 > 내 데이터 > 데이터 조회 |
| **정정권** (Right to Rectification) | Art.16 | §164.526 | §36 | Art.46 | Art.34 | 설정 > 프로필 편집 |
| **삭제권** (Right to Erasure) | Art.17 | 제한적 | §36 | Art.47 | Art.33 | 설정 > 계정 삭제 |
| **처리 제한권** | Art.18 | - | §37 | Art.44 | - | 설정 > 데이터 처리 일시중지 |
| **데이터 이동권** (Portability) | Art.20 | - | - | Art.45 | - | 설정 > 내 데이터 > 내보내기 |
| **이의 제기권** | Art.21 | - | §37 | Art.44 | Art.30 | 설정 > 개인정보 > 이의 제기 |
| **자동화된 결정 거부권** | Art.22 | - | - | Art.24 | - | AI 코치 끄기 옵션 |

### 5.2 권리 행사 처리 절차

```
1. 요청 접수 (앱 내 + 이메일 + 고객센터)
   ↓ (본인 확인: MFA 재인증)
2. 요청 검증 (유효성, 본인 여부)
   ↓ (1영업일 내)
3. 요청 처리
   ↓
4. 결과 통지
   ↓
5. 기록 보관 (3년)

처리 기한:
- GDPR: 1개월 (최대 3개월 연장)
- PIPA: 10일
- HIPAA: 30일 (최대 60일 연장)
- PIPL: 15일
- APPI: 지체 없이
```

### 5.3 데이터 내보내기 형식

- **건강 데이터**: FHIR R4 JSON (국제 의료 표준)
- **측정 이력**: CSV + JSON (선택)
- **개인정보**: JSON (GDPR Art.20 호환)
- **전체 내보내기**: ZIP 압축 (암호화된 다운로드 링크, 48시간 유효)

---

## 6. 데이터 보존 및 삭제

### 6.1 보존 기간

| 데이터 유형 | 보존 기간 | 근거 | 삭제 방법 |
|-----------|---------|------|---------|
| 건강 측정 데이터 | 10년 | IEC 62304, 의료기기법 | 암호화 키 파기 |
| 감사 추적 로그 | 10년 | HIPAA §164.530(j) | 보존 기간 후 자동 삭제 |
| 개인식별정보 | 계정 삭제 후 30일 | GDPR Art.17, PIPA §21 | 완전 삭제 (물리) |
| 동의 기록 | 5년 | GDPR 증빙 의무 | 보존 기간 후 자동 삭제 |
| 앱 사용 로그 | 1년 | 서비스 개선 | 자동 삭제 |
| BLE/NFC 통신 로그 | 90일 | 디버깅 | 자동 삭제 |
| 결제 정보 | 외부 PG 정책 | 전자상거래법 | PG사 관리 |

### 6.2 삭제 절차

```
사용자 삭제 요청 시:

즉시 (0일):
  - 계정 비활성화 (로그인 차단)
  - 개인정보 익명화 처리 시작

30일 이내:
  - PII (이름, 이메일, 프로필) 완전 삭제
  - 가족 그룹 관계 삭제
  - 디바이스 연결 해제

건강 데이터 보존 (법적 요구):
  - 측정 데이터: user_id → 익명 해시로 대체
  - 핑거프린트: user_id 연결 해제
  - 10년 보존 후 물리적 삭제
  
※ 삭제 불가 사유 (법적 보존 의무):
  - 의료기기 측정 기록: 10년 (의료기기법)
  - 감사 추적: 10년 (HIPAA)
  - 법적 분쟁 관련 데이터: 분쟁 종료 시까지
```

---

## 7. 국외 데이터 이전

### 7.1 이전 매트릭스

| 출발 국가 | 도착 국가 | 이전 근거 | 추가 조치 |
|----------|---------|---------|---------|
| 한국 → AWS 서울 | 국내 | N/A | 없음 |
| 한국 → AWS 미국 | 국외 | PIPA §17 동의 | 동의 획득 + 보호 조치 |
| EU → AWS 아일랜드 | EEA 내 | N/A | 없음 |
| EU → 한국 | 적정성 결정국 | GDPR Art.45 | 한국은 적정성 인정 |
| EU → 미국 | EU-US DPF | GDPR Art.45 | DPF 인증 필요 |
| 중국 → 해외 | 금지 (원칙) | PIPL Art.38 | CAC 안전 평가 필수 |
| 일본 → 해외 | APPI Art.28 | 동의 + 적정 보호 | 이전 기록 보관 |

### 7.2 중국 데이터 현지화 전략

```
⚠️ 중국 PIPL에 따른 필수 조치:

1. 중국 내 데이터센터 구축 (AWS China/Alibaba Cloud)
2. 중국 사용자 데이터는 중국 내에만 저장
3. 해외 이전 필요 시:
   a. CAC 보안 평가 통과 (100만명 이상 또는 중요 데이터)
   b. 표준 계약 체결
   c. 개인정보 보호 인증 취득
4. 중국 내 데이터 보호 책임자 지정
```

---

## 8. AI/ML 데이터 활용 정책

### 8.1 연합학습 (Federated Learning)

```
원칙: 데이터는 디바이스를 떠나지 않음

1. 로컬 학습: 각 디바이스에서 모델 업데이트 계산
2. 파라미터만 전송: 모델 가중치 차이(gradient)만 서버로 전송
3. Secure Aggregation: 서버도 개별 gradient 볼 수 없음
4. 차분 프라이버시: 노이즈 추가로 개인 데이터 추론 방지
```

### 8.2 AI 투명성

- **AI 코치 결과에 대한 고지**: "이 결과는 AI 분석에 기반한 참고 정보이며, 의료 진단이 아닙니다"
- **모델 버전 추적**: 어떤 모델 버전이 어떤 결과를 생성했는지 기록
- **설명 가능성**: 주요 기여 특성(feature importance) 사용자에게 제공

---

## 9. 침해 통지 절차

### 9.1 탐지 후 액션

```
침해 발생 → 탐지 (모니터링/사용자 보고)
  ↓
[0~2시간] 초기 평가
  - 영향 범위 (몇 명, 어떤 데이터)
  - 침해 유형 (유출/변조/삭제/접근)
  ↓
[2~24시간] 상세 조사
  - 근본 원인 분석
  - 영향받은 정보주체 목록
  ↓
[24~72시간] 규제당국 보고
  - GDPR: 72시간 내 감독기관
  - PIPA: 72시간 내 개인정보보호위원회
  - HIPAA: 60일 내 HHS OCR
  ↓
[지체없이] 정보주체 통지
  - 통지 내용: 사고 개요, 영향 데이터, 완화 조치, 연락처
  - 통지 방식: 이메일 + 앱 내 알림 + 웹사이트 게시
```

---

## 10. 정책 이행 체크리스트

### 10.1 개발 시 필수 사항

| # | 항목 | 검증 방법 | 담당 |
|---|------|---------|------|
| DP01 | 모든 L3/L4 데이터 필드에 암호화 적용 | 코드 리뷰 + 자동 스캔 | 보안 에이전트 |
| DP02 | API 응답에 불필요한 PII 미포함 | API 리뷰 | 백엔드 에이전트 |
| DP03 | 로그에 PII/PHI 미포함 | 로그 정책 검증 | 인프라 에이전트 |
| DP04 | 동의 없는 데이터 수집 코드 없음 | 코드 리뷰 | 보안 에이전트 |
| DP05 | 삭제/익명화 API 구현 | 기능 테스트 | 백엔드 에이전트 |
| DP06 | 감사 로그 기록 | 통합 테스트 | 인프라 에이전트 |
| DP07 | 데이터 내보내기 기능 | 기능 테스트 | 프론트엔드 에이전트 |
| DP08 | 동의 관리 UI | UI 테스트 | 프론트엔드 에이전트 |
| DP09 | 오프라인 데이터 암호화 | 보안 테스트 | Rust Core 에이전트 |
| DP10 | BLE/NFC 통신 암호화 | 프로토콜 테스트 | Rust Core 에이전트 |

---

---

## 11. 구현 현황 및 검증 매핑

### 11.1 보안 테스트 구현 현황

| 정책 항목 | 검증 테스트 | 파일 위치 |
|----------|-----------|---------|
| DP01 암호화 | TestTLSEnforcement, TestMeasurementDataEncryptionAtRest | tests/security/data_security_test.go |
| DP02 API 응답 PII | TestSensitiveDataInResponse, TestDataMinimization | tests/security/data_security_test.go |
| DP03 로그 PII | TestServerInformationLeakage, TestErrorDetailExposure | tests/security/api_security_test.go |
| DP04 동의 검증 | TestConsentRequired | tests/security/data_security_test.go |
| DP05 삭제/익명화 | TestDataDeletionRight | tests/security/data_security_test.go |
| DP06 감사 로그 | TestAuditTrailOnSensitiveActions | tests/security/api_security_test.go |
| DP09 오프라인 암호화 | AES-256 SQLCipher (Rust crypto 모듈) | rust-core/manpasik-engine/src/crypto/ |
| DP10 BLE/NFC 암호화 | AES-CCM + ECDH (Rust ble/nfc 모듈) | rust-core/manpasik-engine/src/ble/, src/nfc/ |

### 11.2 접근 제어 검증

| 역할 | 검증 시나리오 | 테스트 파일 |
|------|------------|-----------|
| 일반 사용자 | 관리자 API 차단 (TestAdminUserManagement) | tests/e2e/admin_test.go |
| 타 사용자 | 수평적 권한 상승 차단 (TestHorizontalPrivilegeEscalation) | tests/security/auth_security_test.go |
| 관리자 | 수직적 권한 상승 차단 (TestVerticalPrivilegeEscalation) | tests/security/auth_security_test.go |
| 무인증 | 측정 API 차단 (TestMeasurementWithoutAuth) | tests/e2e/measurement_flow_test.go |
| 가족 | 데이터 공유 권한 (TestFamilyMemberPermissions) | tests/e2e/community_family_test.go |

### 11.3 OWASP Top 10 방어 검증

| OWASP | 위협 | 방어 테스트 | 상태 |
|-------|------|-----------|------|
| A01 | 취약한 접근 제어 | TestHorizontal/VerticalPrivilegeEscalation | ✅ 구현 |
| A02 | 암호화 실패 | TestTLSEnforcement, TestTokenInURL | ✅ 구현 |
| A03 | 인젝션 | TestSQLInjection, TestXSSInjection, TestCommandInjection | ✅ 구현 |
| A04 | 불안전한 설계 | TestRateLimiting | ✅ 구현 |
| A05 | 보안 구성 오류 | TestSecurityHeaders, TestServerInformationLeakage | ✅ 구현 |
| A06 | 취약 컴포넌트 | TestDirectoryTraversal, dependency_scan.sh | ✅ 구현 |
| A07 | 인증 실패 | TestExpired/MalformedTokenRejection, TestBruteForceProtection | ✅ 구현 |
| A08 | 무결성 실패 | TestJWTAlgorithmConfusion | ✅ 구현 |
| A09 | 로깅 실패 | TestAuditTrailOnSensitiveActions | ✅ 구현 |
| A10 | SSRF | TestSSRFPrevention | ✅ 구현 |

---

**Document Version**: 1.1.0
**Next Review**: 2026-02-23 (2주 후)
**Approval Required**: DPO, Legal, CISO
