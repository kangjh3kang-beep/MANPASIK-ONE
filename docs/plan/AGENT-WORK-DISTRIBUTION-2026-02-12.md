# ManPaSik 에이전트 업무분장 계획 v2.0

> **발행일**: 2026-02-12  
> **기준**: Sprint 0 완료 후 전체 시스템 분석 기반  
> **목적**: 다중 에이전트 병렬 작업 시 충돌 방지 및 효율 극대화

---

## 1. 현재 구축 현황 종합 분석

### 1.1 완료된 작업

| 영역 | 완료 항목 | 완료 시점 |
|------|----------|----------|
| **Rust 코어** | 8모듈 100% 빌드 + 62테스트 (differential/fingerprint/ai/ble/nfc/dsp/crypto/sync) | Phase 1 |
| **Go 백엔드** | 20+ 마이크로서비스 (auth/user/device/measurement + subscription/shop/payment/ai-inference/cartridge/calibration/coaching + community/video/translation/telemedicine/family/notification/admin/health-record/prescription/reservation) | Phase 1-3 |
| **공유 모듈 5개 연동** | Redis(3서비스), Kafka(1), Milvus(1), Elasticsearch(2), MinIO(gateway) | Sprint 0 ✅ |
| **PostgreSQL** | 20/20 서비스 조건부 DB 지원 (4개 서비스 Sprint 0에서 추가) | Sprint 0 ✅ |
| **Flutter 앱** | 7화면, gRPC 4서비스 연동, BLE/NFC 스텁, 60+ 테스트 | Phase 1C |
| **인프라** | Docker Compose 21서비스, K8s 39 YAML, Prometheus+Grafana, CI/CD | Phase 9-10 |
| **기획/문서** | 기획안 12종, 추적성 80REQ, NFR/이벤트스키마/테스트/배포전략, 규정문서 10종 | 기획 보완 ✅ |
| **Proto** | manpasik.proto 2650줄, 20서비스 인터페이스, Phase 3 포함 | Phase 3 Proto |

### 1.2 미완료 항목 (우선순위별)

#### P0 — 즉시 해결 필요
| # | 항목 | 현재 상태 | 담당 추천 |
|---|------|----------|----------|
| 1 | **Flutter 단위 테스트** | 0개 (TDD 원칙 위배) | Agent 2 (Flutter) |
| 2 | **Rust FFI 브리지 활성화** | main.dart에서 주석 처리 | Agent 1 (Rust) |
| 3 | **Kafka 이벤트 발행 확장** | measurement만 연동, 다른 서비스 미연동 | Agent 3 (Backend) |
| 4 | **연합학습 기반선 학습** | ai-inference 시뮬레이션 엔진만 존재 | Agent 1 (Rust/AI) |
| 5 | **vision-service 신규 구현** | 서비스 미존재 (음식 인식/칼로리) | Agent 3 (Backend) |

#### P1 — Sprint 1-2 대상
| # | 항목 | 현재 상태 | 담당 추천 |
|---|------|----------|----------|
| 6 | **IEC 62304 문서 작성** (SDP/SRS/SAD) | 미작성 | Agent 4 (규정) |
| 7 | **Flutter market Feature** (상품/장바구니/결제) | 미구현 | Agent 2 (Flutter) |
| 8 | **Flutter medical Feature** (화상진료/예약) | 미구현 | Agent 2 (Flutter) |
| 9 | **실제 PG 결제 연동** (Toss/NicePay) | payment-service 시뮬레이션만 | Agent 3 (Backend) |
| 10 | **실제 푸시 알림 (FCM)** | notification-service 인메모리만 | Agent 3 (Backend) |
| 11 | **위험 예측/경고** | ai-inference 확장 필요 | Agent 1 (AI) |
| 12 | **OTA 펌웨어 업데이트 실제 구현** | 인터페이스만 존재 | Agent 1 (Rust) |

#### P2 — Sprint 3-4 대상
| # | 항목 | 현재 상태 | 담당 추천 |
|---|------|----------|----------|
| 13 | Flutter community Feature | 미구현 | Agent 2 |
| 14 | Flutter family Feature | 미구현 | Agent 2 |
| 15 | WebRTC 실제 연동 | video-service 인메모리 시그널링만 | Agent 3 |
| 16 | 실제 번역 API 연동 | translation-service 시뮬레이션만 | Agent 3 |
| 17 | Grafana 대시보드 | 미구성 | Agent 5 (Infra) |
| 18 | Next.js 웹 (관리자/의료진) | 미구현 | Agent 2 |

---

## 2. 에이전트 업무 분장

### Agent 1: Rust 코어 + AI/ML (Antigravity / 전문 에이전트)

**역할**: Rust 엔진 실제 구현, FFI 브리지 활성화, AI 모델 통합

**Sprint 1 작업 (우선순위순):**

| Task ID | 작업 | 파일 범위 | 의존성 | 예상 공수 |
|---------|------|----------|--------|----------|
| R-1 | **Rust FFI 브리지 활성화** — flutter_rust_bridge 빌드 설정, main.dart 주석 해제 | `rust-core/flutter-bridge/`, `frontend/flutter-app/lib/main.dart` | 없음 | 2일 |
| R-2 | **BLE 실제 구현** — btleplug 기반 리더기 스캔·연결·데이터 수신 | `rust-core/manpasik-engine/src/ble/` | R-1 | 5일 |
| R-3 | **NFC 실제 구현** — 카트리지 태그 읽기/검증, ISO 14443A | `rust-core/manpasik-engine/src/nfc/` | R-1 | 3일 |
| R-4 | **AI TFLite 추론 실제 구현** — 모델 로드·추론·결과 반환 | `rust-core/manpasik-engine/src/ai/` | R-1 | 5일 |
| R-5 | **OTA 펌웨어 업데이트 로직** | `rust-core/manpasik-engine/`, `backend/services/device-service/` | R-2 | 3일 |

**수정 금지 영역**: `backend/shared/`, `backend/gateway/`, `frontend/flutter-app/lib/features/`

---

### Agent 2: Flutter 앱 + 프론트엔드 (ChatGPT / 전문 에이전트)

**역할**: Flutter UI/UX, 단위 테스트, 새 Feature 구현

**Sprint 1 작업 (우선순위순):**

| Task ID | 작업 | 파일 범위 | 의존성 | 예상 공수 |
|---------|------|----------|--------|----------|
| F-1 | **Flutter 단위 테스트 60개 작성** (기존 6 Feature 커버) | `frontend/flutter-app/test/` | 없음 | 3일 |
| F-2 | **market Feature** — 상품 목록, 장바구니, 주문, 결제 UI | `frontend/flutter-app/lib/features/market/` | F-1 | 5일 |
| F-3 | **medical Feature** — 화상진료 예약, 의사 프로필, 비디오 통화 UI | `frontend/flutter-app/lib/features/medical/` | F-1 | 7일 |
| F-4 | **community Feature** — 포럼, 댓글, 챌린지, 좋아요 UI | `frontend/flutter-app/lib/features/community/` | F-2 | 5일 |
| F-5 | **family Feature** — 가족 구성원, 모니터링, 알림 UI | `frontend/flutter-app/lib/features/family/` | F-2 | 5일 |

**수정 금지 영역**: `backend/`, `rust-core/`, `infrastructure/`

---

### Agent 3: Go 백엔드 확장 (Claude / 현재 에이전트)

**역할**: 기존 서비스 고도화, 신규 서비스 구현, Kafka 이벤트 확장

**Sprint 1 작업 (우선순위순):**

| Task ID | 작업 | 파일 범위 | 의존성 | 예상 공수 |
|---------|------|----------|--------|----------|
| B-1 | **Kafka 이벤트 발행 확장** — payment, subscription, device 서비스에 이벤트 발행 추가 | `backend/services/{payment,subscription,device}-service/` | 없음 | 3일 |
| B-2 | **vision-service 신규 구현** — 음식 인식, 칼로리 분석 | `backend/services/vision-service/` | 없음 | 7일 |
| B-3 | **실제 PG 결제 연동** — Toss Payments API 래퍼 | `backend/services/payment-service/` | 없음 | 5일 |
| B-4 | **FCM 푸시 알림 연동** — Firebase Cloud Messaging 클라이언트 | `backend/services/notification-service/` | 없음 | 3일 |
| B-5 | **Agent A 작업** — reservation-service 구역 검색 + 의사 프로필 확장 | `backend/services/reservation-service/` | 없음 | 3일 |
| B-6 | **Agent B 작업** — prescription-service 약국 전송 + 배송 추적 | `backend/services/prescription-service/` | 없음 | 3일 |

**수정 금지 영역**: `frontend/`, `rust-core/`

---

### Agent 4: 규정/보안/문서 (Claude 규정 모드)

**역할**: IEC 62304, ISO 14971, 보안 감사, 기술 문서

**Sprint 1 작업 (우선순위순):**

| Task ID | 작업 | 파일 범위 | 의존성 | 예상 공수 |
|---------|------|----------|--------|----------|
| D-1 | **IEC 62304 SDP** (소프트웨어 개발 계획) | `docs/compliance/iec62304-sdp.md` | 없음 | 3일 |
| D-2 | **IEC 62304 SRS** (소프트웨어 요구사항 명세) | `docs/compliance/iec62304-srs.md` | D-1 | 3일 |
| D-3 | **IEC 62304 SAD** (소프트웨어 아키텍처 설계) | `docs/compliance/iec62304-sad.md` | D-1 | 3일 |
| D-4 | **DPIA 템플릿** (데이터 보호 영향평가) | `docs/compliance/dpia-template.md` | 없음 | 2일 |
| D-5 | **Predicate Device 조사** (유사 의료기기 분석) | `docs/compliance/predicate-device-analysis.md` | 없음 | 2일 |
| D-6 | **보안 감사 보고서** (OWASP Top 10 기준) | `docs/security/security-audit-report.md` | 없음 | 2일 |

**수정 금지 영역**: 모든 코드 파일 (문서만 생성/수정)

---

### Agent 5: 인프라/DevOps/통합 (통합 에이전트)

**역할**: Docker/K8s, CI/CD, 모니터링, 통합 테스트

**Sprint 1 작업 (우선순위순):**

| Task ID | 작업 | 파일 범위 | 의존성 | 예상 공수 |
|---------|------|----------|--------|----------|
| I-1 | **Docker Compose 갱신** — Sprint 0 신규 서비스 반영, ES/MinIO 환경변수 | `infrastructure/docker/` | 없음 | 1일 |
| I-2 | **Grafana 대시보드** — 서비스별 메트릭, 알림 규칙 | `infrastructure/monitoring/` | I-1 | 3일 |
| I-3 | **E2E 테스트 확장** — 신규 연동 서비스 검증 (ES/Redis/Kafka) | `backend/tests/e2e/` | I-1 | 3일 |
| I-4 | **CD 파이프라인 갱신** — Canary 배포 설정, 롤백 자동화 | `.github/workflows/cd.yml` | I-1 | 2일 |
| I-5 | **K8s Overlay 갱신** — 신규 서비스 매니페스트, HPA 설정 | `infrastructure/kubernetes/` | I-1 | 2일 |

**수정 금지 영역**: `backend/services/*/internal/` (서비스 내부 로직)

---

## 3. 병렬 작업 충돌 방지 규칙

### 3.1 파일 소유권 매트릭스

| 디렉토리 | 소유 에이전트 | 다른 에이전트 수정 시 |
|----------|-------------|-------------------|
| `rust-core/` | Agent 1 | 금지 |
| `frontend/flutter-app/` | Agent 2 | 금지 |
| `backend/services/*/internal/` | Agent 3 | Agent 1(device-service OTA만 허용) |
| `backend/services/*/cmd/main.go` | Agent 3/5 | 조율 필요 |
| `backend/shared/` | Agent 3/5 | 조율 필요 (PR 기반) |
| `backend/gateway/` | Agent 3/5 | 조율 필요 |
| `docs/compliance/` | Agent 4 | 금지 |
| `docs/plan/` | 공용 | 자유 |
| `infrastructure/` | Agent 5 | 금지 |
| `.github/workflows/` | Agent 5 | 금지 |
| `CHANGELOG.md` | 공용 | 작업 완료 시 상단 추가 |
| `CONTEXT.md` | 공용 | 작업 완료 시 갱신 |
| `KNOWN_ISSUES.md` | 공용 | 이슈 발견 시 추가 |

### 3.2 동기화 프로토콜

1. **작업 시작 전**: `CONTEXT.md`, `CHANGELOG.md`, `KNOWN_ISSUES.md` 최신 내용 확인
2. **작업 중**: 다른 에이전트 파일 수정 필요 시 `CHANGELOG.md`에 "요청" 기록
3. **작업 완료 시**: `CHANGELOG.md` 상단에 변경 사항 기록, `CONTEXT.md` 갱신
4. **충돌 시**: 가장 최근 빌드 성공 상태를 기준으로 병합

### 3.3 빌드 검증 의무

| 시점 | 검증 명령 | 담당 |
|------|----------|------|
| 코드 수정 후 | `go build ./services/{서비스}/...` | 해당 에이전트 |
| 서비스 완성 후 | `go test ./services/{서비스}/...` | 해당 에이전트 |
| Sprint 완료 시 | `go build ./...` + `go test ./...` | Agent 5 |
| 주간 | 전체 E2E 테스트 | Agent 5 |

---

## 4. Sprint 1 타임라인

```
Week 1 (즉시):
  Agent 1: R-1 (FFI 활성화) + R-2 시작 (BLE)
  Agent 2: F-1 (Flutter 테스트 60개)
  Agent 3: B-1 (Kafka 확장) + B-5 (reservation)
  Agent 4: D-1 (IEC 62304 SDP)
  Agent 5: I-1 (Docker 갱신) + I-3 (E2E 확장)

Week 2:
  Agent 1: R-2 (BLE) + R-3 (NFC)
  Agent 2: F-2 (market Feature)
  Agent 3: B-2 (vision-service) + B-6 (prescription)
  Agent 4: D-2 (SRS) + D-3 (SAD)
  Agent 5: I-2 (Grafana) + I-4 (CD)

Week 3:
  Agent 1: R-4 (AI TFLite)
  Agent 2: F-3 (medical Feature)
  Agent 3: B-3 (PG 결제) + B-4 (FCM)
  Agent 4: D-4 (DPIA) + D-5 (Predicate)
  Agent 5: I-5 (K8s) + 통합 검증

Week 4:
  Agent 1: R-5 (OTA)
  Agent 2: F-4 (community) + F-5 (family)
  Agent 3: 미진 항목 마무리 + 코드 리뷰
  Agent 4: D-6 (보안 감사)
  Agent 5: 전체 통합 테스트 + 품질 게이트 검증
```

---

## 5. 성공 기준 (Sprint 1 Gate)

| 항목 | 기준 |
|------|------|
| Go 빌드 | `go build ./...` 전체 PASS |
| Go 테스트 | `go test ./...` 전체 PASS |
| Flutter 테스트 | 60+ 단위 테스트, `flutter test` PASS |
| Rust 빌드 | `cargo build --features full` PASS |
| Rust FFI | Flutter에서 Rust 함수 호출 성공 |
| E2E | 10+ 시나리오 PASS |
| 커버리지 | Go 60%+, Flutter 40%+ |
| 문서 | IEC 62304 SDP/SRS/SAD 3종 완성 |
| Docker | 전체 서비스 Compose 기동 성공 |

---

**마지막 업데이트**: 2026-02-12 (Claude — Sprint 0 완료 후 전체 현황 분석 기반)
