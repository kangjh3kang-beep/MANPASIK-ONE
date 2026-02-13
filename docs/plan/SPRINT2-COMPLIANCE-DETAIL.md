# Sprint 2 — IEC 62304 규정 문서 세부 작성 기획서

> 범위: D-1 (SDP), D-2 (SRS), D-3 (SAD) | 예상 기간: 1~2일

---

## 1. 작성 순서

```
D-1 (SDP) → D-2 (SRS) → D-3 (SAD)
```

SDP가 개발 프로세스를 정의하고, SRS가 요구사항을 명세하며, SAD가 아키텍처로 요구사항을 실현합니다.

---

## 2. D-1: IEC 62304 SDP (소프트웨어 개발 계획서)

### 파일: `docs/compliance/iec62304-sdp.md`

### 목차 구조 (IEC 62304 Clause 5.1 매핑)

| 섹션 | IEC 62304 조항 | 내용 |
|------|---------------|------|
| 1. 개요 | 5.1.1 | 문서 목적, 범위, 용어 정의 |
| 2. 소프트웨어 안전 분류 | 5.1.2 | Class B 판정 근거, 27개 서브시스템 등급 할당 |
| 3. 개발 프로세스 | 5.1.3 | 애자일+V-모델 하이브리드, Sprint 주기, 단계별 활동 |
| 4. 개발 도구 및 환경 | 5.1.4 | Go 1.22+, Flutter 3.x, Rust 1.75+, PostgreSQL, Docker, K8s |
| 5. 형상 관리 전략 | 5.1.5 | Git 브랜칭(trunk-based), 버전 체계(SemVer), 태그 규칙 |
| 6. 위험 관리 참조 | 5.1.6 | ISO 14971 위험관리 계획서 참조, STRIDE 위협 모델 연계 |
| 7. 검증 및 유효성 확인 | 5.1.7 | V&V 전략(단위/통합/E2E/수동), Quality Gate 프로세스 참조 |
| 8. 문서 관리 | 5.1.8 | 문서 버전 관리, 검토 주기, 승인 프로세스 |
| 9. 소프트웨어 유지보수 | 5.1.9 | 패치 프로세스, 핫픽스 절차, 릴리스 주기 |
| 10. 문제 해결 프로세스 | 5.1.10 | 버그 리포팅, 심각도 분류, 해결 SLA |
| 11. SOUP 관리 | 5.1.11 | SOUP 식별·검증·모니터링 절차 |
| 부록 A | — | SOUP/OTS 목록 (§6 참조) |
| 부록 B | — | 도구 검증 요약 |

### 섹션별 핵심 내용

**3. 개발 프로세스:**
- Sprint 2주, 계획→설계→구현→검증→릴리스
- 각 Sprint에서 안전 관련 변경은 추가 위험 검토
- Quality Gate: L1(자동 빌드/테스트) → L2(코드 리뷰+보안) → L3(규정 검증)

**4. 개발 도구:**
| 도구 | 용도 | 버전 |
|------|------|------|
| Go | 백엔드 서비스 | 1.22+ |
| Flutter/Dart | 모바일 앱 | 3.x / 3.x |
| Rust | 코어 엔진 | 1.75+ |
| protoc | gRPC 코드 생성 | 25.x |
| Docker | 컨테이너화 | 24.x |
| GitHub Actions | CI/CD | — |
| PostgreSQL | 데이터베이스 | 16 |

**11. SOUP 관리:**
- 신규 SOUP 도입 시 라이선스·보안·품질 평가 수행
- 분기별 SOUP 업데이트 검토 (CVE 모니터링)

### 참조 문서
- `docs/compliance/software-safety-classification.md` → §2
- `docs/compliance/risk-management-plan.md` → §6
- `docs/compliance/vnv-master-plan.md` → §7
- `QUALITY_GATES.md` → §7

---

## 3. D-2: IEC 62304 SRS (소프트웨어 요구사항 명세서)

### 파일: `docs/compliance/iec62304-srs.md`

### 목차 구조 (IEC 62304 Clause 5.2 매핑)

| 섹션 | IEC 62304 조항 | 내용 |
|------|---------------|------|
| 1. 개요 | 5.2.1 | 문서 목적, 시스템 개요 |
| 2. 기능 요구사항 | 5.2.2 | 80개 REQ (추적성 매트릭스 기반) |
| 3. 비기능 요구사항 | 5.2.3 | 성능/가용성/확장성/보안/규정 |
| 4. 인터페이스 요구사항 | 5.2.4 | gRPC API, BLE, NFC, REST, Kafka |
| 5. 위험 통제 요구사항 | 5.2.5 | ISO 14971 잔여 위험 → SW 통제 |
| 6. 규제 요구사항 | 5.2.6 | IEC 62304, ISO 13485, FDA, CE, KGMP |
| 7. 데이터 요구사항 | 5.2.7 | 데이터 모델, 보존 기간, 암호화 |
| 8. 사용 환경 | 5.2.8 | 지원 OS, 네트워크, 하드웨어 |
| 부록 A | — | 추적성 매트릭스 |
| 부록 B | — | 유스케이스 다이어그램 |

### 섹션별 핵심 내용

**2. 기능 요구사항 (80개 REQ):**

| REQ ID | 카테고리 | 요구사항 | DES | IMP | V&V |
|--------|---------|---------|-----|-----|-----|
| REQ-001 | 인증 | 사용자 등록/로그인 | DES-001 | auth-service | UT-001, E2E-001 |
| REQ-002 | 측정 | 바이오마커 측정 세션 | DES-002 | measurement-service | UT-002, E2E-002 |
| ... | ... | ... | ... | ... | ... |
> 전체 80개 REQ는 `docs/plan/plan-traceability-matrix.md`에서 참조

**3. 비기능 요구사항:**

| ID | 유형 | 요구사항 | 목표값 |
|----|------|---------|-------|
| NFR-001 | 성능 | API 응답 시간 (p95) | < 200ms |
| NFR-002 | 성능 | 동시 사용자 | 10,000+ |
| NFR-003 | 가용성 | 서비스 가용률 | 99.9% |
| NFR-004 | 보안 | 데이터 암호화 (at rest) | AES-256 |
| NFR-005 | 보안 | 데이터 암호화 (in transit) | TLS 1.3 |
| NFR-006 | 확장성 | 수평 확장 (서비스별) | K8s HPA |
| NFR-007 | 규정 | GDPR/PIPA 동의 관리 | 6개 동의 유형 |
| NFR-008 | 규정 | 의료 데이터 보존 | 10년 |

**4. 인터페이스 요구사항:**

| 인터페이스 | 프로토콜 | 설명 |
|-----------|---------|------|
| 서비스 간 통신 | gRPC (protobuf) | 20+ 서비스 상호 호출 |
| 모바일 → Gateway | gRPC-Web / REST | Flutter 앱 ↔ API Gateway |
| 디바이스 연결 | BLE 5.0 | 만파식 디바이스 데이터 수신 |
| 카트리지 인증 | NFC (ISO 14443A) | 카트리지 태그 읽기/검증 |
| 이벤트 스트림 | Kafka | 서비스 간 비동기 이벤트 |
| 외부 결제 | HTTPS (REST) | Toss Payments API |
| 푸시 알림 | HTTPS | Firebase FCM |

**5. 위험 통제 요구사항:**
- STRIDE 위협 모델(31개 위협)에서 도출된 SW 통제 조치
- ISO 14971 잔여 위험 → 소프트웨어 설계/구현으로 통제
- 예: THR-001(인증 우회) → REQ-AUTH-MFA(다중 인증) + REQ-AUTH-LOCKOUT(계정 잠금)

### 참조 문서
- `docs/plan/plan-traceability-matrix.md` → §2, 부록A
- `docs/compliance/stride-threat-model.md` → §5
- `docs/compliance/regulatory-compliance-checklist.md` → §6
- `docs/compliance/data-protection-policy.md` → §7

---

## 4. D-3: IEC 62304 SAD (소프트웨어 아키텍처 설계서)

### 파일: `docs/compliance/iec62304-sad.md`

### 목차 구조 (IEC 62304 Clause 5.3 매핑)

| 섹션 | IEC 62304 조항 | 내용 |
|------|---------------|------|
| 1. 개요 | 5.3.1 | 문서 목적, 아키텍처 원칙 |
| 2. 시스템 아키텍처 | 5.3.2 | MSA 전체 구성도, 3-tier 구조 |
| 3. 소프트웨어 항목 분해 | 5.3.3 | 20+ 서비스 목록, 각 서비스 역할·안전 등급 |
| 4. 서비스 내부 구조 | 5.3.4 | handler → service → repository 패턴 |
| 5. 데이터 흐름 | 5.3.5 | 주요 시나리오별 시퀀스 다이어그램 |
| 6. 통신 프로토콜 | 5.3.6 | gRPC, Kafka, REST, BLE, NFC 상세 |
| 7. 데이터 저장 계층 | 5.3.7 | PostgreSQL, Redis, Milvus, Elasticsearch, MinIO |
| 8. 보안 아키텍처 | 5.3.8 | JWT, mTLS, RBAC, 암호화, 네트워크 보안 |
| 9. 배포 아키텍처 | 5.3.9 | K8s, Docker, CI/CD 파이프라인 |
| 10. 인터페이스 명세 | 5.3.10 | Proto 파일 참조, API 엔드포인트 목록 |
| 11. SOUP/OTS 목록 | 5.3.11 | 전체 의존성 목록 |
| 부록 A | — | 서비스 통신 매트릭스 |
| 부록 B | — | 데이터 모델 ERD |

### 섹션별 핵심 내용

**2. 시스템 아키텍처:**
```
[Flutter App] → [API Gateway] → [gRPC Services ×20+]
                                       ↕
                              [PostgreSQL / Redis / Kafka]
                                       ↕
                              [Milvus / Elasticsearch / MinIO]

[Rust Core Engine] ↔ [Flutter App] (FFI Bridge)
```

**3. 소프트웨어 항목 분해 (20+ 서비스):**

| 서비스 | 안전 등급 | 역할 |
|--------|----------|------|
| auth-service | Class B | 사용자 인증·JWT 발급 |
| user-service | Class B | 사용자 프로필 관리 |
| device-service | Class B | 디바이스 등록·관리 |
| measurement-service | Class B | 측정 세션·결과 관리 |
| ai-inference-service | Class B | AI 추론·분석 |
| calibration-service | Class B | 디바이스 보정 |
| cartridge-service | Class B | 카트리지 인증·추적 |
| coaching-service | Class A | 건강 코칭·리포트 |
| payment-service | Class A | 결제 처리 |
| subscription-service | Class A | 구독 관리 |
| shop-service | Class A | 상품·주문 |
| notification-service | Class A | 알림 전송 |
| community-service | Class A | 커뮤니티 |
| family-service | Class A | 가족 관리 |
| health-record-service | Class B | 건강 기록 |
| telemedicine-service | Class B | 화상진료 |
| reservation-service | Class A | 예약 관리 |
| prescription-service | Class B | 처방 관리 |
| admin-service | Class A | 시스템 관리 |
| vision-service | Class B | 음식 분석 |
| translation-service | Class A | 다국어 번역 |
| video-service | Class A | 비디오 관리 |
| gateway | Class A | API 라우팅·인증 |

**7. 데이터 저장 계층:**

| 저장소 | 용도 | 서비스 |
|--------|------|--------|
| PostgreSQL | 주 데이터 저장 | 전체 서비스 |
| Redis | 캐시, 세션, 레이트리밋 | device, subscription, auth, gateway |
| Kafka (Redpanda) | 이벤트 스트림 | measurement, payment, config 등 |
| Milvus | 벡터 검색 | measurement (유사 측정 검색) |
| Elasticsearch | 전문 검색 | measurement, community (검색) |
| MinIO | 오브젝트 스토리지 | gateway (파일 업로드) |

**11. SOUP/OTS 목록:**

#### Go 의존성
| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| google.golang.org/grpc | 1.62+ | gRPC 서버/클라이언트 | Apache-2.0 |
| github.com/jackc/pgx/v5 | 5.x | PostgreSQL 드라이버 | MIT |
| github.com/go-redis/redis/v9 | 9.x | Redis 클라이언트 | BSD-2 |
| github.com/segmentio/kafka-go | 0.4+ | Kafka 프로듀서/컨슈머 | MIT |
| github.com/milvus-io/milvus-sdk-go | 2.x | Milvus 벡터 DB | Apache-2.0 |
| github.com/elastic/go-elasticsearch | 8.x | Elasticsearch | Apache-2.0 |
| github.com/minio/minio-go | 7.x | MinIO/S3 | Apache-2.0 |
| github.com/golang-jwt/jwt/v5 | 5.x | JWT 처리 | MIT |
| go.uber.org/zap | 1.x | 구조화 로깅 | MIT |
| google.golang.org/protobuf | 1.x | Protobuf | BSD-3 |

#### Flutter 의존성
| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| grpc | 3.x | gRPC 클라이언트 | BSD-3 |
| protobuf | 3.x | Protobuf 직렬화 | BSD-3 |
| provider | 6.x | 상태 관리 | MIT |
| go_router | 14.x | 라우팅 | BSD-3 |
| fl_chart | 0.68+ | 차트 | MIT |
| flutter_secure_storage | 9.x | 보안 스토리지 | BSD-3 |
| flutter_markdown | 0.7+ | 마크다운 렌더링 | BSD-3 |

#### Rust 의존성
| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| flutter_rust_bridge | 2.x | FFI 브리지 | MIT |
| btleplug | 0.11+ | BLE 통신 | MIT/Apache-2.0 |
| ring | 0.17+ | 암호화 | ISC |
| tflitec | 0.3+ | TFLite 추론 | Apache-2.0 |
| tokio | 1.x | 비동기 런타임 | MIT |

---

## 5. 참조 문서 매핑

| SDP/SRS/SAD 섹션 | 참조 문서 |
|------------------|----------|
| SDP §2 안전 분류 | `docs/compliance/software-safety-classification.md` |
| SDP §6 위험 관리 | `docs/compliance/risk-management-plan.md` |
| SDP §7 V&V | `docs/compliance/vnv-master-plan.md` |
| SRS §2 기능 요구사항 | `docs/plan/plan-traceability-matrix.md` |
| SRS §5 위험 통제 | `docs/compliance/stride-threat-model.md` |
| SRS §6 규제 | `docs/compliance/regulatory-compliance-checklist.md` |
| SRS §7 데이터 | `docs/compliance/data-protection-policy.md` |
| SAD §3 서비스 분해 | `docs/compliance/software-safety-classification.md` §서브시스템 |
| SAD §11 SOUP | `backend/go.mod`, `frontend/flutter-app/pubspec.yaml`, `rust-core/Cargo.toml` |

---

## 6. 검토 기준 체크리스트

### D-1 (SDP) 완성도

- [ ] IEC 62304 Clause 5.1의 모든 하위 조항 커버
- [ ] ManPaSik 특화 도구·환경 명시
- [ ] 형상 관리 전략이 실제 Git 워크플로우와 일치
- [ ] Quality Gate 프로세스가 QUALITY_GATES.md와 일관
- [ ] SOUP 관리 절차가 구체적 (도입/검증/모니터링)

### D-2 (SRS) 완성도

- [ ] 80개 REQ 전체 포함 (추적성 매트릭스 참조)
- [ ] 비기능 요구사항 정량화 (응답시간, TPS, SLA 등)
- [ ] 인터페이스 요구사항이 실제 Proto 정의와 일치
- [ ] 위험 통제 요구사항이 STRIDE 위협과 매핑
- [ ] 규제 요구사항이 5개국 체크리스트와 일관

### D-3 (SAD) 완성도

- [ ] 20+ 서비스 전체 목록 및 안전 등급 명시
- [ ] 서비스 내부 구조(handler→service→repository) 기술
- [ ] 데이터 흐름 시퀀스 다이어그램 최소 5개 시나리오
- [ ] SOUP/OTS 목록이 go.mod/pubspec.yaml/Cargo.toml과 일치
- [ ] 보안 아키텍처가 STRIDE 대응 조치와 일관
- [ ] 배포 아키텍처가 실제 K8s/Docker 설정과 일치

---

**마지막 업데이트**: 2026-02-12
