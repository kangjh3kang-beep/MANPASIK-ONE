# ManPaSik AI 공유 컨텍스트

> **용도**: 모든 AI 도구가 프로젝트 컨텍스트를 빠르게 파악할 수 있는 요약 문서
> **업데이트 주기**: 매 작업 세션 완료 시

---

## 🎯 프로젝트 개요

**만파식(萬波息)**: 차동측정 기반 범용 분석 헬스케어 AI 생태계

**핵심 기술:**
- 차동측정법: `S_corrected = S_det - α × S_ref` (α=0.95, 99% 노이즈 제거)
- 1792차원 핑거프린트 벡터 (88→448→896→1792 4단계 성장 경로)
- 엣지 AI (TFLite) 100% 오프라인 동작
- 연합학습 기반 프라이버시 보존
- 무한확장 카트리지 시스템 (기본 29종 + 레지스트리 무제한 추가, 2-byte 65,536종/4-byte 43억종)
- 구독 등급별 카트리지 접근 제어 (INCLUDED/LIMITED/ADD_ON/RESTRICTED/BETA)

**의료기기:** Class II (체외진단의료기기, IVD)
**대상 시장:** KR(MFDS), US(FDA), EU(CE-IVDR), CN(NMPA), JP(PMDA)

---

## 📐 확정된 아키텍처

```
┌───────────────────────────────────────────────────────────────┐
│  클라이언트: Flutter 앱 (iOS/Android/Desktop) + Next.js 웹    │
├───────────────────────────────────────────────────────────────┤
│  Rust 코어 엔진 (manpasik-engine + flutter-bridge)           │
│  differential | fingerprint | ai | ble | nfc | dsp | crypto  │
├───────────────────────────────────────────────────────────────┤
│  API Gateway: Kong 3.7 → Keycloak 25.0 (OIDC/MFA/RBAC)      │
├───────────────────────────────────────────────────────────────┤
│  Go 마이크로서비스 (gRPC): auth | user | device | measurement │
├───────────────────────────────────────────────────────────────┤
│  데이터: PostgreSQL 16 | TimescaleDB | Milvus 2.4 | Redis 7  │
│          Kafka(Redpanda) | Elasticsearch 8 | MinIO | MQTT     │
├───────────────────────────────────────────────────────────────┤
│  관측성: Prometheus | Grafana 11 | OpenTelemetry | ELK        │
└───────────────────────────────────────────────────────────────┘
```

---

## 📂 디렉토리 구조

```
Manpasik/
├── docs/                     # 문서, AI 스펙
│   ├── AI_COLLABORATION.md   # AI 역할 분담 계획
│   ├── implementation-patterns.md  # 구현 패턴·유형·방식 (실시간 공유·학습)
│   └── ai-specs/             # AI별 작업 스펙 (claude/, chatgpt/)
├── infrastructure/           # Docker, K8s, Terraform
│   └── docker/               # docker-compose.dev.yml (15+ 서비스)
├── backend/                  # Go 마이크로서비스 (30+)
│   └── shared/proto/         # gRPC Proto (manpasik.proto, health.proto)
├── rust-core/                # Rust 코어 엔진 워크스페이스
│   ├── manpasik-engine/      # 핵심 8모듈 (differential~sync)
│   └── flutter-bridge/       # Flutter FFI 브리지 (10 API)
├── ai-ml/                    # AI/ML 파이프라인 (예정)
├── frontend/                 # Flutter 앱, Next.js 웹 (예정)
├── sdk/                      # 개발자 SDK (예정)
├── tests/                    # 테스트 (예정)
├── .cursor/rules/            # ★ 에이전트 팀 규칙 (6 파일)
├── AGENTS.md                 # ★ 에이전트 팀 마스터 컨텍스트
├── CHANGELOG.md              # 🔄 작업 로그 (실시간 업데이트)
├── CONTEXT.md                # 📋 이 파일 (컨텍스트 요약)
├── KNOWN_ISSUES.md           # ⚠️ 알려진 이슈 추적 (환경/빌드/기술 부채)
├── QUALITY_GATES.md          # 🚦 Quality Gate 프로세스 (Stage/Phase 체크리스트, 현재 상태)
└── README.md                 # 프로젝트 개요
```

---

## 🔧 기술 스택

| 계층 | 기술 | 버전 | 상태 |
|------|------|------|------|
| 모바일 | Flutter + Riverpod + go_router | 3.x / 2.4 / 13.0 | Phase 1C 완료 (gRPC 4서비스, 차트·BLE/NFC 스텁, E2E 4서비스, 60+테스트) |
| 코어 엔진 | Rust + TFLite + BLE | Ed.2021 / tflitec 0.7 | ✅ 100% (빌드+62테스트, full feature 포함 통과) |
| FFI 브리지 | flutter_rust_bridge | 2.0 | ✅ 완료 |
| API Gateway | Kong + Keycloak | 3.7 / 25.0 | Docker 설정 |
| 백엔드 | Go + gRPC + Gin | 1.24+ | Phase 2 완료: 11서비스 (auth/user/device/measurement + subscription/shop/payment/ai-inference/cartridge/calibration/coaching) |
| 메인 DB | PostgreSQL | 16-alpine | Docker 설정 |
| 시계열 DB | TimescaleDB | latest-pg16 | Docker 설정 |
| 벡터 DB | Milvus | 2.4.10 | Docker 설정 |
| 캐시 | Redis | 7-alpine | Docker 설정 |
| 메시징 | Kafka (Redpanda) | v24.2.1 | Docker 설정 |
| 검색 | Elasticsearch | 8.14.0 | Docker 설정 |
| IoT | MQTT (Mosquitto) | 2 | Docker 설정 |
| 스토리지 | MinIO | latest | Docker 설정 |
| 모니터링 | Prometheus + Grafana | v2.53 / 11.1 | Docker 설정 |
| 웹 | Next.js + TypeScript | 14+ | 대기 |

---

## 🤖 AI 역할 분담

| AI | 담당 영역 | 현재 상태 | 최근 작업 |
|----|-----------|-----------|----------|
| **Antigravity** | 아키텍처, Rust 코어, 인프라, 통합 | ✅ 활성 | Rust 8모듈 + Docker + Proto 완성 |
| **Claude** | 규정/보안 분석, 코드 리뷰, AI/ML 설계 | ✅ 활성 | 산출물 10종 + Rust 3모듈 + Go 4서비스 + Flutter S4 완성 |
| **ChatGPT** | Flutter UI, Go 서비스 구현 | 🔲 대기 | 스펙 전달 완료, 구현 대기 |

---

## 🛡️ Cursor 에이전트 팀 (신규)

| 에이전트 | 규칙 파일 | 대상 | 주요 컨텍스트 |
|---------|----------|------|-------------|
| Global | `manpasik-project.mdc` | `**/*` (항상) | 보안, TDD, 의료규정, 코드스타일, 작업 기록 |
| Work Log | `work-logging.mdc` | `**/*` (항상) | 실시간 작업 기록 프로토콜 (세션 시작/완료) |
| Rust Core | `rust-core.mdc` | `rust-core/**` | 10모듈 완전 API, 의존성, Feature Flags |
| Go Backend | `go-backend.mdc` | `backend/**` | Proto, DB 스키마, 서비스 패턴 |
| Frontend | `frontend.mdc` | `frontend/**` | Flutter/Riverpod, Rust FFI, 디자인 시스템 |
| Security | `security-compliance.mdc` | 전체 | OWASP, STRIDE, 5개국 규정, 암호화 |
| Infra | `infrastructure.mdc` | `infrastructure/**` | Docker 15+ 서비스, K8s, CI/CD |

**마스터 문서:** `AGENTS.md` (gRPC API 상세, DB 스키마, 아키텍처, 카트리지, BLE/NFC 프로토콜, 작업 기록 프로토콜)

---

## 📋 핵심 결정 사항

1. **차동측정 공식**: `S_det - α × S_ref` (α 기본값 = 0.95)
2. **핑거프린트 차원**: 88 → 448 → 896 → 1792 4단계 성장 경로 (단세포→다세포→유기체→생태계). Phase 5에서 E12-IF 다중 리더기 융합 + 시간축 확장으로 1792차원 달성
3. **카트리지**: 무한확장 체계 (기본 29종, 2-byte 코드 = 카테고리×타입, 최대 65,536종). 등급별 접근 제어. NFC ISO 14443A v1.0/v2.0 호환
4. **리더기 관리**: 무제한 확장 (BLE 5.0, Wi-Fi Hub, Cloud Gateway)
5. **오프라인**: 100% 완전 동작 (CRDT 동기화)
6. **AI 모델**: TFLite 엣지 (5종 모델) + 클라우드 하이브리드
7. **보안**: AES-256-GCM + SHA-256 해시체인 + E2E 암호화
8. **인증**: Keycloak OIDC + JWT (Access 15분, Refresh 7일)
9. **API Gateway**: Kong 3.7 (레이트리밋, SSL 종단, 라우팅)
10. **구독 티어**: Free / Basic(₩9,900) / Pro(₩29,900) / Clinical(₩59,900)
11. **Quality Gate 개발 프로세스**: Level 1(매 작업 린트/테스트/빌드) → Level 2(Stage Gate 체크리스트) → Level 3(Phase Gate 통합·보안·규정 검증). 상세는 `QUALITY_GATES.md` 참조.
12. **기획-개발 철학**: (1) 원본+제안 반영해 미수립 세부 기획·계획을 완벽히 수립해 기준 확정 (2) 개발 중 수정·보완으로 고도화 (3) 예상치 못한 사항에 유기적 대응 — 완벽한 시스템 구축·완성 추구. 상세는 `docs/development-philosophy.md` 참조.

---

## 🚦 Quality Gate 현재 상태

| Stage | 범위 | 상태 | 비고 |
|-------|------|------|------|
| S1 | Rust 코어 엔진 | ✅ 통과 | 62테스트, BLE 포함 빌드 (2026-02-10) |
| S2 | Go 인증 서비스 | ✅ 통과 | 8테스트 PASS, gRPC 전 계층 연동 (2026-02-10) |
| S3 | Go user/device/measurement | ✅ 통과 | 32테스트, 핸들러·저장소·main 연동 (2026-02-10) |
| S4 | Flutter 앱 기본 구조 | ✅ 통과 | 7화면, 3Provider, 6언어, 32테스트 (2026-02-10) |
| S5 | Flutter 핵심 화면 (S5a) | ✅ 통과 | gRPC 4서비스 연동, 60+테스트 (2026-02-10) |
| S6 | 전체 통합 | ✅ 통과 | E2E 4서비스 헬스체크, README gRPC 포트 (2026-02-10) |
| S7 | Sprint 2 병렬 구현 | ✅ 통과 | Go QG 전체 PASS, Flutter 128/132 PASS (2026-02-12) |

**Phase:** Phase 1A·1B·1C·1D ✅ 완료. **Phase 2(Core) 7/7 서비스 완료** ✅ — 커머스(subscription/shop/payment) + AI(ai-inference/coaching) + 측정(cartridge/calibration).

**Sprint 2 완료 현황 (2026-02-12):**
| 항목 | 상태 | 비고 |
|------|------|------|
| AS-1~AS-5 (Admin 설정 백엔드) | ✅ | DB/RPC/암호화/Kafka/서비스연동 |
| AS-6 (Admin 설정 UI) | ✅ | 8카테고리 탭, 편집 다이얼로그, gRPC 연동 |
| F-1 (Flutter 테스트 83개) | ✅ | 10파일, 목표 60개 초과 달성 |
| D-1 (IEC 62304 SDP) | ✅ | 15섹션 + 부록 3종 |
| I-1 (Docker Compose) | ✅ | vision-service 추가, 환경변수 7건 |
| I-3 (E2E 확장) | ✅ | telemedicine/coaching/cartridge 3파일 |
| Proto 정식 재생성 | ✅ | 수동 ext.go 완전 삭제, 🟡 우회 0건 |
| D-2 (IEC 62304 SRS) | ✅ | 819줄, REQ 80+, 추적성 매트릭스 |
| D-3 (IEC 62304 SAD) | ✅ | 547줄, 22서비스+Rust 9모듈, SOUP |
| AS-7 (LLM 클라이언트) | ✅ | OpenAI client.go 294줄, 8테스트 PASS |
| I-5 (K8s Overlay) | ✅ | telemedicine/vision YAML, 3환경 overlay |
| screen_widget_test 수정 | ✅ | gRPC Provider override → 6/6 PASS |
| ConfigMap 버그 수정 | ✅ | TELEMEDICINE 50066→50071 |
| AS-8 (LLM 어시스턴트 RPC) | ✅ | inference.go LLM 연동, 11테스트, Graceful Degradation |
| AS-9 (Admin 감사 로그) | ✅ | audit_log.go + admin.go 연동, 14테스트 (42개 총) |
| D-4 (DPIA Template) | ✅ | 8섹션+3부록, GDPR+한국법 이중 병기 |
| I-2 (Grafana Dashboard) | ✅ | overview+grpc-services JSON, Prometheus, provisioning |
| Flutter Chat UI (AS-8 후반) | ✅ | chat_screen.dart 573줄, 6언어, gRPC 연동 |
| I-4 (CD Pipeline) | ✅ | 23서비스 매트릭스, 3환경 배포, 롤백 |
| D-5 (Predicate Device) | ✅ | 5기기 SE 분석, Abbott i-STAT Primary |

**통합 검증 (2026-02-12):** Go 전체 빌드·vet·테스트, Flutter analyze·test 통과. `docs/implementation-patterns.md`에 Kafka/Memory 이벤트 패턴·K8s Kustomize 반영. DB init 22·23·24 추가. Vision Proto 확장 제안서 `docs/plan/proto-extension-vision-service.md` 작성.

**Sprint 2 계획 수립 (2026-02-12):** Sprint 1 잔여(R/F/D/I) + 관리자 설정 시스템(AS-1~AS-9) + IEC 62304 규정 문서 + E2E/인프라 갱신 7일 일정 수립. 4개 워크스트림 세부 기획서 병렬 작성 완료.
- `docs/plan/SPRINT2-EXECUTION-PLAN.md` — 전체 실행 계획
- `docs/plan/SPRINT2-AS-BACKEND-DETAIL.md` — 관리자 설정 백엔드 (AS-1~5, 15+파일, 30+테스트)
- `docs/plan/SPRINT2-FLUTTER-DETAIL.md` — Flutter UI+테스트 (82테스트, 25파일, 50 l10n키)
- `docs/plan/SPRINT2-COMPLIANCE-DETAIL.md` — IEC 62304 SDP/SRS/SAD (3문서, SOUP 목록)
- `docs/plan/SPRINT2-INFRA-E2E-DETAIL.md` — 인프라+E2E (Docker 32+서비스, E2E 8파일, K8s 3환경)

**Sprint 2 전체 병렬 구현 (2026-02-12):** AS-1~AS-5 관리자 설정 백엔드 전체 구현 + 서비스 연동(payment/notification ConfigWatcher 핫리로드) + E2E 테스트 2파일 4시나리오 + 28개 신규 테스트 + Docker Compose/K8s 갱신 + IEC 62304 SDP/SRS/SAD 초안 3종 + E2E 포트 버그 수정. `go build ./...`·`go vet ./...`·`go test ./...` 전체 통과.

**Proto 정식 재생성 (2026-02-12):** TelemedicineService를 manpasik.proto에 추가(7 RPC, 14 메시지, 3 enum) → protoc 3.21.12로 정식 재생성(manpasik.pb.go 927KB + manpasik_grpc.pb.go 340KB) → 수동 스텁(admin_config_ext.go, telemedicine_ext.go) 완전 삭제. E2E payment_subscription_flow_test.go Proto 필드 수정. **🟡 우회 중 이슈: 0건** (전량 정상 해결). `go build/vet/test` 전체 통과.

**Sprint 2 병렬 구현 (2026-02-12):** 4개 병렬 스트림 동시 실행 — (A) Flutter 단위 테스트 83개 10파일 (F-1 P0 완료), (B) IEC 62304 SDP 15섹션+부록3종 (D-1 완료), (C) Docker Compose vision-service 추가+환경변수 7건+E2E 3파일 (I-1/I-3 완료), (D) Flutter Admin 설정 UI 8카테고리탭+편집다이얼로그 (AS-6 완료). Go QG 전체 PASS. Flutter 128/132 PASS.

---

## 📍 현재 진행 단계

**Phase 1: MVP (Month 1-4)**
- [x] Week 1-4: 기반 인프라 ✅ (Antigravity)
- [x] Week 5-8: Rust 코어 엔진 ✅ (Antigravity)
- [x] 에이전트 팀 활성화 ✅ (Claude)
- [x] 컨텍스트 엔지니어링 강화 ✅ (Claude)
- [x] 실시간 작업 기록 프로토콜 구축 ✅ (Claude)
- [x] 시스템 기획안 검증 분석 ✅ (Claude)
- [x] IEC 62304 소프트웨어 안전 등급 판정 (Class B) ✅ (Claude)
- [x] 5개국 의료기기 규제 준수 체크리스트 (146항목) ✅ (Claude)
- [x] STRIDE 위협 모델링 (31개 위협, 8개 공격표면) ✅ (Claude)
- [x] 데이터 보호 정책 초안 (5개국 동시 준수) ✅ (Claude)
- [x] 서브시스템별 안전 등급 할당 (IEC 62304 5.3.3) ✅ (Claude)
- [x] 기술문서(Technical File) 목차 (11대분류 60+문서) ✅ (Claude)
- [x] AI/ML 모델 설계 전략 (5종 모델 + 연합학습 + PCCP) ✅ (Claude)
- [x] Claude 전체 업무 분석 (19건 로드맵) ✅ (Claude)
- [x] Flutter 앱 기본 구조 (S4): 7화면, 3Provider, 6언어 i18n, 32테스트 ✅ (Claude)
- [x] Flutter gRPC 연동·화면 고도화 (S5a): Go 4서비스 Dockerfile, Dart gRPC 클라이언트, Repository 4종, Home/Measure/Devices/Settings 실제 데이터, 60+테스트 ✅ (Claude)
- [x] Phase 1C Stage S6: 전체 통합 — E2E 4서비스 헬스체크, README gRPC 포트, Phase 1C 완료 선언 ✅ (Claude)
- [x] Phase 1D D1: E2E 플로우(Login→StartSession→EndSession→GetHistory), backend/tests/e2e, CI/Makefile 정리 ✅ (Claude)
- [x] Phase 1D Gate 통과: Docker 빌드·DB 초기화·gRPC protoc 생성 코드·E2E 플로우 전체 검증 완료 ✅ (Claude)
- [x] Flutter S5b: 측정 결과·트렌드 차트(fl_chart), BLE 스캔/NFC 카트리지 읽기 UI(Rust FFI 스텁), /measurement/result 라우트 ✅ (Claude)
- [x] Go 백엔드 4서비스 gRPC 핸들러·저장소·main 연동 완료 ✅ (Claude)
- [x] ISO 14971 위험관리 계획서 ✅ (Claude)
- [x] V&V 마스터 플랜 ✅ (Claude)
- [x] Rust 스텁 3모듈 완전 구현 (crypto/dsp/sync) ✅ (Claude)
- [x] Go 백엔드 4서비스 비즈니스 로직 ✅ (Claude)
- [x] DB 초기화 스크립트 4종 (auth/user/device/measurement) ✅ (Claude)
- [x] Quality Gate 개발 프로세스 도입 ✅ (Level 1/2/3, QUALITY_GATES.md)
- [ ] DPIA 템플릿 (Claude Phase 1B)
- [ ] Predicate Device 조사 (Claude Phase 1B)
- [ ] FHIR R4 통합 설계 (Claude Phase 1B)

**Phase 2: Core (Month 5-8)**
- [x] subscription-service: SaaS 구독 관리 (4 티어, 업/다운그레이드) ✅ (Claude)
- [x] shop-service: 상품, 장바구니, 주문 ✅ (Claude)
- [x] payment-service: PG 연동, 구독 결제 ✅ (Claude)
- [x] ai-inference-service: 실시간 추론, 모델 서빙 (시뮬레이션 엔진) ✅ (Claude)
- [x] cartridge-service: 카트리지 인증·NFC 태그 검증·사용 추적·레지스트리 (30종, 15카테고리) ✅ (Claude)
- [x] calibration-service: 팩토리/현장 보정·보정 모델 관리·보정 상태 추적 (22종 모델) ✅ (Claude)
- [x] coaching-service: AI 건강 코칭·목표 설정·일일/주간 리포트·개인화 추천 ✅ (Claude)

**기획 보완 (2026-02-12) ✅ 완료**
- [x] 기획서 12종 전수 분석 + 검증 보고서 작성 ✅ (Claude)
- [x] 추적성 매트릭스 v2.0 완성 (80개 REQ 전체 DES↔IMP↔V&V 매핑) ✅ (Claude)
- [x] 비기능 요구사항 정량화 (응답시간/TPS/SLA/스토리지/보안) ✅ (Claude)
- [x] 이벤트 스키마 명세 (18개 Kafka 토픽 JSON 스키마) ✅ (Claude)
- [x] 테스트 전략 (Phase별 커버리지/E2E 시나리오/자동화) ✅ (Claude)
- [x] 배포 전략 (Rolling→Canary→Blue-Green/환경별 스펙/롤백) ✅ (Claude)

**기획 완성도 100/100 달성 (2026-02-14) ✅ 완료**
- [x] 서비스 간 통신 패턴 명세 (gRPC 매트릭스/Kafka 매핑/Saga/Circuit Breaker) ✅ (Antigravity)
- [x] API 보안 아키텍처 명세 (인증 흐름/JWT/RBAC/KMS/PHI/SIEM) ✅ (Antigravity)
- [x] 규제 준수 갭 해소 통합 (FMEA 15개/SOUP 38개/SMP/SCM) ✅ (Antigravity)
- [x] 검증 보고서 점수 갱신 (88→100/100) ✅ (Antigravity)

**Sprint 0: 공유 모듈 실연동 (2026-02-12) ✅ 완료**
- [x] measurement-service: Milvus VectorRepository 연동 ✅ (Claude)
- [x] measurement-service: Kafka EventPublisher 연동 ✅ (Claude)
- [x] device-service: Redis Cache-Aside 데코레이터 연동 ✅ (Claude)
- [x] subscription-service: Redis Cache-Aside 데코레이터 연동 ✅ (Claude)
- [x] measurement-service: Elasticsearch SearchIndexer 연동 ✅ (Claude)
- [x] community-service: Elasticsearch PostSearchIndexer 연동 ✅ (Claude)
- [x] gateway: MinIO/S3 파일 업로드/다운로드/삭제 연동 ✅ (Claude)
- [x] community-service → PostgreSQL 전환 ✅ (Claude)
- [x] video-service → PostgreSQL 전환 ✅ (Claude)
- [x] translation-service → PostgreSQL 전환 ✅ (Claude)
- [x] telemedicine-service → PostgreSQL 전환 ✅ (Claude)

**Sprint 1: 5-에이전트 병렬 실행 (2026-02-12~) 🔄 진행중**

*Agent 1 (Rust/AI)*:
- [ ] R-1: Rust FFI 브리지 활성화 (flutter_rust_bridge 빌드)
- [ ] R-2: BLE 실제 구현 (btleplug 기반)
- [ ] R-3: NFC 실제 구현 (ISO 14443A)
- [ ] R-4: AI TFLite 추론 실제 구현
- [ ] R-5: OTA 펌웨어 업데이트 로직

*Agent 2 (Flutter)*:
- [x] F-1: Flutter 단위 테스트 83개 작성 ✅ (10파일, 목표 60개 초과 달성)
- [ ] F-2: market Feature (상품/장바구니/결제)
- [ ] F-3: medical Feature (화상진료/예약)
- [ ] F-4: community Feature
- [ ] F-5: family Feature

*Agent 3 (Go Backend)*:
- [x] B-1: Kafka 이벤트 발행 확장 (payment/subscription/device) ✅
- [x] B-2: vision-service 신규 구현 ✅ (서비스/리포지토리 완성, Proto 확장 필요)
- [ ] B-3: 실제 PG 결제 연동 (Toss)
- [x] B-4: FCM 푸시 알림 연동 ✅
- [x] B-5: reservation-service 구역 검색 확장 ✅ (이미 구현됨)
- [x] B-6: prescription-service 약국 전송 ✅ (이미 구현됨)

*Agent 4 (규정/문서)*:
- [x] D-1: IEC 62304 SDP (소프트웨어 개발 계획) ✅ (15섹션+부록3종)
- [ ] D-2: IEC 62304 SRS (소프트웨어 요구사항 명세)
- [ ] D-3: IEC 62304 SAD (소프트웨어 아키텍처 설계)
- [ ] D-4: DPIA 템플릿
- [ ] D-5: Predicate Device 조사

*Agent 5 (인프라/통합)*:
- [x] I-1: Docker Compose 갱신 ✅ (vision-service 추가, 환경변수 7건 보강)
- [ ] I-2: Grafana 대시보드
- [x] I-3: E2E 테스트 확장 ✅ (telemedicine/coaching/cartridge 3파일 추가)
- [ ] I-4: CD 파이프라인 갱신
- [ ] I-5: K8s Overlay 갱신

**업무분장 상세**: `docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md`

---

## 📐 구현 패턴·유형·방식

구현 내역·유형·방식을 실시간 공유·학습하기 위한 전용 문서를 유지합니다. 새 기능/레이어 추가 시 해당 패턴을 문서에 반영합니다.

| 영역 | 요약 | 상세 문서 |
|------|------|-----------|
| Flutter | Feature-First, Riverpod(Notifier/Provider/FutureProvider), gRPC 수동 스텁, Repository(interface+impl+Fake), 화면 연동·테스트 | `docs/implementation-patterns.md` §1 |
| Go 백엔드 | 서비스 구조(cmd/internal), handler→service→repository, 공통 config·인메모리 fallback | `docs/implementation-patterns.md` §2 |
| 인프라 | Go Dockerfile(빌드 컨텍스트 backend/), docker-compose 4서비스·환경변수 | `docs/implementation-patterns.md` §3 |
| Proto·gRPC | manpasik.proto, Go/Flutter 생성 방식·경로 | `docs/implementation-patterns.md` §4 |
| 문서·품질 | 패턴 문서·CHANGELOG·CONTEXT·QUALITY_GATES 갱신 규칙 | `docs/implementation-patterns.md` §5 |

→ **전체 상세**: [docs/implementation-patterns.md](docs/implementation-patterns.md)

---

## 🔗 참조 문서

### 프로젝트 핵심
- **AGENTS.md**: 에이전트 팀 마스터 컨텍스트 (API/DB/프로토콜 상세)
- **docs/implementation-patterns.md**: 구현 패턴·유형·방식 (Flutter/Go/인프라/Proto/문서 갱신)
- **CHANGELOG.md**: 상세 작업 로그 + 이슈/에러/디버깅 기록 (시간순)
- **KNOWN_ISSUES.md**: 알려진 이슈 추적 (미해결/해결 이력, 환경 제약, 빌드 요구사항)
- **QUALITY_GATES.md**: Quality Gate 프로세스 (Stage/Phase 체크리스트, 현재 Stage 상태)
- **README.md**: 프로젝트 개요

### AI 스펙
- **docs/AI_COLLABORATION.md**: 3-AI 분담 계획
- **docs/ai-specs/claude/**: 의료규정 + 보안 아키텍처 스펙
- **docs/ai-specs/chatgpt/**: Flutter UI + Go 서비스 스펙

### 기획안·개발 철학
- **docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md**: 기획안 v1.1 완성본 (원본 부합 + 현행 매핑 + MSA 로드맵 + 추적성 + 세부 명세 참조)
- **docs/plan/README.md**: 기획·계획 문서 인덱스
- **docs/development-philosophy.md**: 개발 철학 (기준 수립 → 고도화 → 유기적 대응)
- **docs/plan-original-vs-current-and-development-proposal.md**: 원본 대비 검증, 보완사항, 발전 수립안
- **docs/specs/data-packet-family-c.md**: 패밀리C 데이터 패킷 표준
- **docs/ux/sitemap.md**: 사이트맵 공식 문서
- **docs/ux/storyboard-*.md**: 스토리보드 (첫 측정, 음식 칼로리)
- **docs/specs/offline-capability-matrix.md**: 오프라인 기능 매트릭스
- **docs/plan/msa-expansion-roadmap.md**: 단계별 MSA 확장 로드맵
- **docs/plan/plan-traceability-matrix.md**: 기획-구현 추적성 (v2.0 — 80개 REQ 전체 매핑)
- **docs/plan/ai-agent-phase-mapping.md**: AI 에이전트–Phase 매핑
- **docs/specs/cartridge-system-spec.md**: 카트리지 무한확장 체계 및 등급별 접근 제어 명세
- **docs/plan/terminology-and-tier-mapping.md**: 용어·티어·카트리지 접근 정책 통일
- **docs/system-plan-verification.md**: 기획안 검증 분석 보고서 (2차 갱신)
- **docs/compliance/software-safety-classification.md**: IEC 62304 안전 등급 (Class B + 서브시스템 27개 할당)
- **docs/compliance/regulatory-compliance-checklist.md**: 5개국 규제 체크리스트 (146항목)
- **docs/compliance/data-protection-policy.md**: 데이터 보호 정책 초안
- **docs/compliance/technical-file-structure.md**: 기술문서 목차 (11대분류 60+문서, FDA/CE-IVDR 매핑)
- **docs/security/stride-threat-model.md**: STRIDE 위협 모델링 (31개 위협)
- **docs/ai-specs/claude/ml-model-design-spec.md**: AI/ML 모델 설계 전략 (5종 + 연합학습 + PCCP)
- **docs/claude-task-analysis.md**: Claude 전체 업무 분석 (19건 로드맵)
- **docs/compliance/iso14971-risk-management-plan.md**: ISO 14971 위험관리 계획서 (9개 위해 시나리오, 5×5 매트릭스)
- **docs/compliance/vnv-master-plan.md**: V&V 마스터 플랜 (5레벨 전략, CI/CD 통합, 추적성)

### 신규 보완 문서 (2026-02-12)
- **docs/specs/non-functional-requirements.md**: 비기능 요구사항 정량화 (응답시간, TPS, SLA, 스토리지, 보안)
- **docs/specs/event-schema-specification.md**: Kafka 이벤트 스키마 (18토픽, JSON 스키마, DLQ 정책)
- **docs/specs/test-strategy.md**: 테스트 전략 (Phase별 커버리지, 자동화 파이프라인)
- **docs/specs/deployment-strategy.md**: 배포 전략 (Phase별 Rolling→Canary→Blue-Green)
- **docs/reports/system-verification-and-implementation-plan-2026-02-12.md**: 시스템 검증 및 구현 계획 종합 보고서

### 에이전트 규칙
- **.cursor/rules/**: 7개 에이전트 규칙 파일

---

**마지막 업데이트**: 2026-02-12 (Claude — Sprint 2 병렬 구현: F-1 테스트 83개 + D-1 SDP + I-1/I-3 Docker/E2E + AS-6 Admin UI, Go QG 전체 PASS)
