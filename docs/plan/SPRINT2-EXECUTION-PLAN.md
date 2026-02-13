# Sprint 2 실행 계획

> **기준일**: 2026-02-12  
> **범위**: Sprint 1 잔여 + 관리자 설정 시스템 + Phase 3 준비  
> **예상 기간**: 5~7일

---

## 1. 현재 완료 현황

### 완료됨 (Sprint 0~1)

| 영역 | 완료 항목 |
|------|----------|
| **Go Backend** | B-1 Kafka 확장 ✅, B-2 vision-service ✅, B-3 Toss PG ✅, B-4 FCM ✅, B-5 예약 ✅, B-6 처방 ✅ |
| **인프라** | DB init 22~24 추가, implementation-patterns 보완, 종합 검증(build/vet/test) 통과 |
| **기획** | 관리자 설정+LLM 기획서(`admin-settings-llm-assistant-spec.md`), Vision Proto 제안서 |

### 미완료 (Sprint 1 잔여)

| ID | 항목 | 에이전트 | 우선순위 |
|----|------|----------|----------|
| R-1 | Rust FFI 브리지 활성화 | Agent 1 | P1 |
| R-2 | BLE 실제 구현 (btleplug) | Agent 1 | P2 |
| R-3 | NFC 실제 구현 (ISO 14443A) | Agent 1 | P2 |
| R-4 | AI TFLite 추론 실제 구현 | Agent 1 | P2 |
| R-5 | OTA 펌웨어 업데이트 | Agent 1 | P3 |
| F-1 | Flutter 단위 테스트 60개 | Agent 2 | **P0** |
| F-2 | market Feature (상품/장바구니/결제) | Agent 2 | P1 |
| F-3 | medical Feature (화상진료/예약) | Agent 2 | P2 |
| F-4 | community Feature | Agent 2 | P2 |
| F-5 | family Feature | Agent 2 | P3 |
| D-1 | IEC 62304 SDP | Agent 4 | P1 |
| D-2 | IEC 62304 SRS | Agent 4 | P1 |
| D-3 | IEC 62304 SAD | Agent 4 | P1 |
| D-4 | DPIA 템플릿 | Agent 4 | P2 |
| D-5 | Predicate Device 조사 | Agent 4 | P2 |
| I-1 | Docker Compose 갱신 | Agent 5 | P1 |
| I-2 | Grafana 대시보드 | Agent 5 | P2 |
| I-3 | E2E 테스트 확장 | Agent 5 | P1 |
| I-4 | CD 파이프라인 갱신 | Agent 5 | P2 |
| I-5 | K8s Overlay 갱신 | Agent 5 | P2 |

### 신규 (관리자 설정 시스템)

| ID | 항목 | 우선순위 |
|----|------|----------|
| AS-1 | DB 스키마 확장 (25-admin-settings-ext.sql) + 시드 데이터 | P0 |
| AS-2 | admin-service 설정 관리 RPC 확장 (ListConfigs, GetConfigWithMeta, Validate, Bulk) | P0 |
| AS-3 | 암호화 저장/복호화 (AES-256-GCM) | P0 |
| AS-4 | Kafka config.changed 이벤트 + ConfigWatcher | P1 |
| AS-5 | payment-service DB config 로드 (env fallback) | P1 |
| AS-6 | Flutter Admin 설정 UI (목록/편집/이력) | P1 |
| AS-7 | LLM 클라이언트 (OpenAI/Anthropic) 통합 | P2 |
| AS-8 | LLM 어시스턴트 RPC + Flutter 채팅 UI | P2 |
| AS-9 | 다국어 자동 번역 + 변경 대기열 | P3 |

---

## 2. Sprint 2 실행 순서

### Day 1: 관리자 설정 기반 구축 (AS-1 ~ AS-3)

**목표**: DB 스키마·메타데이터·다국어·암호화 완성

```
[AS-1] 25-admin-settings-ext.sql 생성
       - config_metadata, config_translations, llm_config_sessions,
         llm_config_messages, config_change_queue 테이블
       - 설정 메타데이터 시드 (25+ 설정 항목)
       - 다국어 설명 시드 (ko/en/ja 주요 설정)

[AS-2] admin-service 확장
       - ConfigMetadataRepository 인터페이스 + PG/Memory 구현
       - ConfigTranslationRepository 인터페이스 + PG/Memory 구현
       - ListSystemConfigs(category, language) RPC
       - GetConfigWithMeta(key, language) RPC
       - ValidateConfigValue(key, value) RPC
       - BulkSetConfigs RPC

[AS-3] 암호화 모듈
       - backend/services/admin-service/internal/crypto/aes.go
       - AES-256-GCM Encrypt/Decrypt
       - CONFIG_ENCRYPTION_KEY 환경변수 로드
       - SetSystemConfig에서 security_level=secret이면 암호화 저장
```

**검증**: `go build ./...`, `go test ./services/admin-service/...`

### Day 2: 동적 반영 + 서비스 연동 (AS-4 ~ AS-5, F-1 시작)

**목표**: 설정 변경이 서비스에 즉시 반영되는 파이프라인

```
[AS-4] 설정 변경 이벤트
       - admin-service: SetSystemConfig 시 Kafka "manpasik.config.changed" 발행
       - shared/events/config_watcher.go: ConfigWatcher 인터페이스 + Kafka 구현
       - payment-service: Toss 설정 변경 시 pgGateway 재초기화
       - notification-service: FCM 설정 변경 시 FCMClient 재초기화

[AS-5] DB config 우선 로드
       - payment-service cmd/main.go: DB system_configs에서 toss.* 로드 → 없으면 env fallback
       - 패턴: loadConfigFromDB(pool, "toss.secret_key") || os.Getenv("TOSS_SECRET_KEY")

[F-1 시작] Flutter 단위 테스트 작성
       - auth, theme, locale, validator 테스트 (20개 목표)
       - repository fake + widget 테스트 시작
```

**검증**: 전체 빌드·테스트, Kafka 이벤트 발행 확인

### Day 3: Flutter Admin 설정 UI (AS-6, F-1 계속)

**목표**: 관리자가 브라우저/앱에서 설정을 조회·변경할 수 있는 UI

```
[AS-6] Flutter Admin 화면
       - /admin/settings 라우트 추가 (app_router.dart)
       - AdminSettingsScreen: 카테고리 탭 + 설정 카드 목록
       - ConfigEditDialog: 유형별 입력 (string, number, boolean, secret, select)
       - 마크다운 help_text 렌더링 (flutter_markdown)
       - 설정 검색·필터
       - gRPC 클라이언트: ListSystemConfigs, SetSystemConfig, ValidateConfigValue

[F-1 계속] Flutter 테스트 40개 도달
       - grpc_client, measurement_card, device_list 위젯 테스트
```

### Day 4: IEC 62304 규정 문서 + E2E (D-1~D-3, I-3)

**목표**: 의료기기 필수 규정 문서 3종 초안 + E2E 테스트 확장

```
[D-1] IEC 62304 SDP (소프트웨어 개발 계획서)
       - docs/compliance/iec62304-sdp.md
       - 개발 프로세스, 도구, 형상관리, 위험관리 참조

[D-2] IEC 62304 SRS (소프트웨어 요구사항 명세서)
       - docs/compliance/iec62304-srs.md
       - 기능/비기능 요구사항, 추적성 매트릭스 참조

[D-3]roke IEC 62304 SAD (소프트웨어 아키텍처 설계서)
       - docs/compliance/iec62304-sad.md
       - MSA 아키텍처, 서비스 간 통신, 데이터 흐름

[I-3] E2E 테스트 확장
       - payment flow (CreatePayment → ConfirmPayment → RefundPayment)
       - subscription flow (Create → Upgrade → Cancel)
       - admin config flow (Set → Get → List)
```

### Day 5: LLM 어시스턴트 (AS-7 ~ AS-8)

**목표**: LLM 클라이언트 + 설정 어시스턴트 채팅

```
[AS-7] LLM 클라이언트
       - backend/services/ai-inference-service/internal/llm/client.go
       - OpenAI Chat Completion API (gpt-4o)
       - 인터페이스: LLMClient.Chat(ctx, systemPrompt, messages) → response
       - 설정: llm.provider, llm.api_key, llm.model (system_configs에서 로드)

[AS-8] 설정 어시스턴트 RPC
       - StartConfigSession: 세션 생성 + 시스템 프롬프트 구성
       - SendConfigMessage: 대화 → LLM 호출 → ConfigSuggestion 추출
       - ApplyConfigSuggestion: 관리자 승인 → SetSystemConfig 호출
       - 안전 규칙: secret 값 미전달, 제안→승인 2단계

       Flutter 채팅 UI:
       - /admin/settings/assistant 라우트
       - 채팅 인터페이스 + 제안 패널
       - [적용] [거부] 버튼
```

### Day 6~7: 인프라 갱신 + 마무리 (I-1, I-5, F-1 완료)

```
[I-1] Docker Compose 갱신
       - 신규 서비스(vision, 22+ 서비스) 반영
       - 환경변수 정리 (CONFIG_ENCRYPTION_KEY 등)

[I-5] K8s Overlay 갱신
       - 신규 서비스 Deployment/Service YAML
       - ConfigMap에 신규 환경변수 추가

[F-1 완료] Flutter 단위 테스트 60개 도달
       - admin settings 위젯 테스트 포함

최종 검증:
       - go build ./... + go vet ./... + go test ./...
       - flutter analyze + flutter test
       - CONTEXT.md, CHANGELOG.md, KNOWN_ISSUES.md 갱신
```

---

## 3. 의존 관계

```
AS-1 (DB) ──→ AS-2 (RPC) ──→ AS-6 (Flutter UI)
              │                │
              └──→ AS-3 (암호화) ──→ AS-4 (Kafka) ──→ AS-5 (서비스 연동)
                                                       │
AS-7 (LLM 클라이언트) ──→ AS-8 (어시스턴트 RPC+UI) ───┘

D-1/D-2/D-3 (규정 문서) — 독립, 병렬 가능
F-1 (Flutter 테스트) — 독립, 병렬 가능
I-1/I-3/I-5 (인프라) — AS-1 이후 병렬 가능
```

---

## 4. Sprint 2 완료 기준

### 필수 (P0~P1)

- [ ] `25-admin-settings-ext.sql` + 시드 데이터 적용
- [ ] admin-service: ListSystemConfigs, GetConfigWithMeta, ValidateConfigValue, BulkSetConfigs RPC
- [ ] AES-256-GCM 암호화 저장/복호화
- [ ] Kafka config.changed 이벤트 + ConfigWatcher (payment, notification)
- [ ] Flutter Admin 설정 UI (/admin/settings)
- [ ] Flutter 단위 테스트 60개
- [ ] IEC 62304 SDP/SRS/SAD 초안
- [ ] E2E 테스트 확장 (payment + subscription + admin)
- [ ] Docker Compose 갱신

### 선택 (P2~P3)

- [ ] LLM 클라이언트 (OpenAI) + 어시스턴트 RPC
- [ ] Flutter LLM 채팅 UI
- [ ] 다국어 자동 번역 (translation-service 연동)
- [ ] 설정 변경 대기열 + 일괄 승인
- [ ] Grafana 대시보드
- [ ] CD 파이프라인 갱신

### 검증

- `go build ./...` + `go vet ./...` + `go test ./...` 전체 통과
- `flutter analyze` + `flutter test` 통과
- CONTEXT.md, CHANGELOG.md, KNOWN_ISSUES.md, QUALITY_GATES.md 갱신

---

## 5. 이후 전망 (Sprint 3~)

| Sprint | 주요 내용 |
|--------|----------|
| **Sprint 3** | F-2 market Feature, R-1 Rust FFI 활성화, LLM 고도화(AS-9), Phase 3 준비 |
| **Sprint 4** | F-3 medical Feature, R-2/R-3 BLE/NFC, I-2 Grafana, Phase 3 Gate |
| **Sprint 5** | F-4/F-5 community/family, R-4/R-5 AI/OTA, Phase 4 준비 |

---

**참조 문서**:
- `docs/specs/admin-settings-llm-assistant-spec.md` — 관리자 설정 + LLM 세부 기획서
- `docs/plan/B3-toss-pg-integration.md` — Toss PG 연동 (완료)
- `docs/plan/NEXT-STEPS-DETAILED-PLAN.md` — B-3 세부 계획 (완료)
- `docs/plan/proto-extension-vision-service.md` — Vision Proto 확장 제안

**마지막 업데이트**: 2026-02-12
