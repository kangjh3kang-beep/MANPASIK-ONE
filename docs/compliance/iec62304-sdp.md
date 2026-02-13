# IEC 62304:2015 소프트웨어 개발 계획서 (SDP)

> **문서 ID**: DOC-SDP-001  
> **버전**: v1.1  
> **안전 등급**: Class B (IEC 62304)  
> **작성일**: 2026-02-13  
> **상태**: 초안 (검토 대기)  
> **적용 표준**: IEC 62304:2006+AMD1:2015 (의료기기 소프트웨어 — 소프트웨어 수명주기 프로세스)

---

## 1. 목적 및 범위 (Clause 5.1.1)

### 1.1 목적

본 문서는 ManPaSik (만파식) 의료기기 소프트웨어의 **소프트웨어 개발 계획(Software Development Plan)**을 정의합니다. IEC 62304:2015 §5.1 요구사항을 충족하며, 소프트웨어 개발 전 과정에서 일관된 품질과 안전성을 보장하기 위한 프로세스, 도구, 역할, 형상 관리, 유지보수 계획을 명시합니다.

### 1.2 범위

| 항목 | 내용 |
|------|------|
| **제품명** | ManPaSik (만파식, 萬波息) — 차동측정 기반 범용 분석 헬스케어 AI 생태계 |
| **의료기기 등급** | Class II (체외진단의료기기, IVD) |
| **소프트웨어 안전 등급** | IEC 62304 Class B |
| **대상 시장** | KR(MFDS), US(FDA 510(k)), EU(CE-IVDR) |
| **핵심 기능** | 차동측정(S_det − α×S_ref), 896차원 핑거프린트, 오프라인 AI 추론, 원격진료, 건강 관리 |

### 1.3 소프트웨어 시스템 식별

| 구성요소 | 기술 스택 | 역할 |
|---------|----------|------|
| **모바일 앱** | Flutter 3.x + Dart | 사용자 인터페이스, BLE/NFC 연동, 오프라인 측정 |
| **코어 엔진** | Rust (Ed. 2021) | 차동측정, 핑거프린트, 엣지 AI(TFLite), 암호화 |
| **백엔드** | Go 1.22+ (gRPC MSA) | 30+ 마이크로서비스 (인증, 측정, AI, 결제 등) |
| **데이터 계층** | PostgreSQL 16, TimescaleDB, Milvus 2.4, Redis 7, Kafka (Redpanda), Elasticsearch 8.14 | 관계형·시계열·벡터·캐시·이벤트·검색 |
| **API Gateway** | Kong 3.7 + Keycloak 25.0 | 인증, 라우팅, RBAC |

### 1.4 용어 정의

| 용어 | 정의 |
|------|------|
| SDP | 소프트웨어 개발 계획서 (Software Development Plan) |
| SRS | 소프트웨어 요구사항 명세서 (Software Requirements Specification) |
| SAD | 소프트웨어 아키텍처 설계서 (Software Architecture Design) |
| SOUP | 알려진 출처의 소프트웨어 (Software of Unknown Provenance) |
| OTS | 기성품 소프트웨어 (Off-The-Shelf Software) |
| V&V | 검증 및 유효성 확인 (Verification and Validation) |
| PHI | 보호 대상 건강 정보 (Protected Health Information) |
| Sprint | 2주 단위 반복 개발 주기 |
| Quality Gate | 단계별 품질 검증 통과 조건 |

---

## 2. 참조 문서 (Clause 5.1.1)

### 2.1 규정·표준

| 문서 | 버전 | 적용 |
|------|------|------|
| **IEC 62304** | 2006+AMD1:2015 | 의료기기 소프트웨어 — 소프트웨어 수명주기 프로세스 |
| **ISO 14971** | 2019 | 의료기기 — 위험관리의 의료기기 적용 |
| **ISO 13485** | 2016 | 의료기기 — 품질경영시스템 |
| **IEC 62366-1** | 2015+Amd.1 | 의료기기 — 사용적합성 공학(Usability Engineering) |
| **FDA SW Guidance** | 2002 | General Principles of Software Validation |
| **FDA Cyber** | 2023 | Cybersecurity in Medical Devices |

### 2.2 프로젝트 내부 문서

| 문서 | 경로 | 용도 |
|------|------|------|
| SRS | `docs/compliance/iec62304-srs.md` | 소프트웨어 요구사항 |
| SAD | `docs/compliance/iec62304-sad.md` | 아키텍처 설계 |
| 안전 등급 | `docs/compliance/software-safety-classification.md` | Class B 판정 근거 |
| 위험관리 계획 | `docs/compliance/iso14971-risk-management-plan.md` | ISO 14971 위험관리 |
| V&V 마스터 플랜 | `docs/compliance/vnv-master-plan.md` | 검증·확인 전략 |
| Quality Gates | `QUALITY_GATES.md` | 단계별 품질 검증 |
| 추적성 매트릭스 | `docs/plan/plan-traceability-matrix.md` | REQ↔DES↔IMP↔V&V |

---

## 3. 소프트웨어 안전 등급 (Clause 5.1.2)

### 3.1 Class B 판정

ManPaSik 소프트웨어는 **IEC 62304 Class B**로 분류됩니다.

### 3.2 판정 근거

| 기준 | 설명 |
|------|------|
| **IEC 62304 §4.3** | Class B: 위해를 유발하거나 기여할 수 있으나 **심각하지 않음** |
| **위해 시나리오** | 소프트웨어 고장 시 중상 가능성은 낮으나, 건강 측정 데이터의 오류가 잘못된 건강 판단으로 이어질 수 있어 **경상(non-serious injury)** 가능성 존재 |
| **의도된 사용** | 진단 기기가 아닌 **정보 제공 기기** — 의료 전문가 확인 없이 자가 관리에 사용될 수 있음 |
| **직접 치료 결정** | 소프트웨어 결과가 직접적 치료 결정이 아닌 참고 정보 제공 목적 |

### 3.3 서브시스템별 안전 등급

| 서브시스템 | 안전 등급 | 근거 |
|-----------|----------|------|
| measurement-service | Class B | 측정 데이터 정확성 |
| ai-inference-service | Class B | 분석 결과 신뢰성 |
| calibration-service | Class B | 보정 정확성 |
| device-service | Class B | 디바이스 제어 |
| health-record-service | Class B | 건강 기록 무결성 |
| prescription-service | Class B | 처방 데이터 정확성 |
| telemedicine-service | Class B | 원격진료 안전성 |
| auth-service | Class B | 의료 데이터 접근 통제 |
| Rust Core Engine | Class B | 차동측정·핑거프린트 계산 |
| Flutter 앱 (측정/결과 화면) | Class B | 건강 데이터 표시 정확성 |
| 기타 서비스 (14개+) | Class A | 직접적 의료 기능 없음 |

**상세**: `docs/compliance/software-safety-classification.md` 참조

---

## 4. 개발 생명주기 모델 (Clause 5.1.3)

### 4.1 Agile + V-Model 하이브리드

**Sprint + Gate** 조합 모델을 채택합니다.

```
┌───────────────────────────────────────────────────────────────────────────┐
│                     V-Model (왼쪽: 개발, 오른쪽: 검증)                        │
├───────────────────────────────────────────────────────────────────────────┤
│                                                                           │
│  요구사항 분석 ──────────────────────────── 시스템 테스트·확인              │
│       │                                              ▲                    │
│       ▼                                              │                    │
│  아키텍처 설계 ─────────────────────────── 통합 테스트                      │
│       │                                              ▲                    │
│       ▼                                              │                    │
│  상세 설계 ─────────────────────────────── 단위 테스트                    │
│       │                                              ▲                    │
│       ▼                                              │                    │
│  코딩·구현 ───────────────────────────────────────────┘                   │
│                                                                           │
│  ※ Sprint(2주) 단위로 위 단계가 반복·점진적 진행                              │
│  ※ Quality Gate(L1/L2/L3)로 각 단계 완료 조건 검증                           │
└───────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Sprint 주기 (2주)

| Day | 활동 |
|-----|------|
| Day 1-2 | Sprint 계획 + 요구사항 분석 + 백로그 정제 |
| Day 3-8 | 설계 + 구현 + 단위 테스트 (TDD) |
| Day 9-10 | 통합 테스트 + 코드 리뷰 + 보안 검토 |
| Day 11-12 | 검증 + Sprint 리뷰 + 회고 + Gate 판정 |

### 4.3 Quality Gate 프로세스

| Level | 시점 | 활동 | 통과 조건 |
|-------|------|------|----------|
| **L1** | 매 작업 완료 | 린트 + 단위 테스트 + 빌드 | 에러 0, 테스트 통과, 빌드 성공 |
| **L2** | Stage 완료 | 코드 리뷰 + 보안 검토 + 통합 테스트 | Stage별 Gate 체크리스트 충족 |
| **L3** | Phase 완료 | 성능 테스트 + 규정 문서 정합성 + 릴리스 판정 | E2E·보안·문서 검증 완료 |

**참조**: `QUALITY_GATES.md`

---

## 5. 개발 프로세스 (Clause 5.1.3)

### 5.1 프로세스 흐름

```
요구사항(SRS) → 아키텍처(SAD) → 상세설계 → 코딩 → 단위테스트 → 통합 → 시스템
     │              │               │          │         │          │        │
     ▼              ▼               ▼          ▼         ▼          ▼        ▼
  REQ 분석      구조 분해       모듈 설계    TDD 구현    UT 80%+   IT·E2E   V&V
  추적성 매트릭스  서비스 목록    인터페이스   코드 리뷰   회귀 테스트  배포 전 검증
```

### 5.2 단계별 활동

| 단계 | IEC 62304 절 | 활동 | 산출물 |
|------|-------------|------|--------|
| 1. 요구사항 | §5.2 | 사용자/시스템 요구사항 분석, SRS 작성, 위험 통제 요구사항 도출 | SRS, 추적성 매트릭스 |
| 2. 아키텍처 | §5.3 | 소프트웨어 항목 분해, 인터페이스 정의, SAD 작성 | SAD, 서비스 목록 |
| 3. 상세 설계 | §5.4 (Class C) | Class B에서는 복잡 안전 모듈만 권고 적용 | 상세 설계서 (선택) |
| 4. 코딩 | §5.5 | TDD, 코드 리뷰, 정적 분석, 코딩 표준 준수 | 소스 코드, 테스트 코드 |
| 5. 단위 테스트 | §5.5 | 단위 검증, 커버리지 ≥80% | 단위 테스트 보고서 |
| 6. 통합 | §5.6 | 모듈/서비스 통합, 통합 테스트 | 통합 테스트 보고서 |
| 7. 시스템 | §5.7 | E2E 테스트, 사용적합성, 임상 확인 | 시스템 테스트 보고서 |

### 5.3 TDD 원칙

- **"No Code Without Tests"**: 기능 구현 전 실패하는 테스트 먼저 작성
- **언어별 프레임워크**: Rust `#[cfg(test)]`, Go `testing`, Dart `flutter_test`
- **명명 규칙**: `test_기능_시나리오_기대결과` (한국어 허용)

### 5.4 안전 관련 변경 프로세스

Class B 이상 서브시스템 변경 시:

1. **변경 영향 분석** (위험관리 계획서 참조)
2. **추적성 매트릭스 갱신** (REQ → DES → IMP → V&V)
3. **회귀 테스트** 실행
4. **위험 평가** 갱신

---

## 6. 유지보수 계획 (Clause 5.1.9, §6)

### 6.1 릴리스 주기

| 유형 | 주기 | 조건 |
|------|------|------|
| **정기 릴리스** | 월 1회 (MINOR) | 신규 기능, 개선 |
| **패치 릴리스** | 필요 시 (PATCH) | 버그 수정 |
| **긴급 핫픽스** | 24시간 내 | 심각도 Critical (환자 안전, 데이터 손실) |

### 6.2 변경 관리 프로세스

```
[이슈 등록] → [영향 분석] → [구현+테스트] → [코드 리뷰] → [Quality Gate] → [릴리스] → [문서 갱신]
     │              │              │              │              │              │
  GitHub Issues   안전 등급 확인   TDD, CI 통과   PR 승인       L1/L2 통과    CHANGELOG, SRS
```

### 6.3 시판 후 유지보수

| 활동 | 주기 | 담당 |
|------|------|------|
| 시판 후 감시 | 상시 | PMS 담당 |
| 버그·이슈 추적 | 상시 | GitHub Issues, KNOWN_ISSUES.md |
| 의존성 CVE 모니터링 | 분기 | Dependabot, Renovate |
| SOUP 업데이트 검토 | 분기 | 품질팀 |
| 위험관리 파일 갱신 | 변경 시 | RMR |

---

## 7. 형상 관리 (Clause 5.1.5, §8)

### 7.1 버전 관리 도구

| 도구 | 용도 |
|------|------|
| **Git** | 소스 코드, 문서, 스크립트 버전 관리 |
| **GitHub** | 원격 저장소, PR, Issues, Actions |

### 7.2 브랜치 전략

| 브랜치 | 용도 | 수명 |
|--------|------|------|
| `main` | 프로덕션 릴리스, 항상 배포 가능 상태 | 영구 |
| `feature/*` | 기능 개발 (예: `feature/MPK-123-measurement-ui`) | 최대 1주 |
| `hotfix/*` | 긴급 버그 수정 | 최소화 |
| `release/*` | 릴리스 준비, 최종 검증 | 릴리스 시 |

**전략**: Trunk-Based Development — 단기 기능 브랜치, 자주 main 병합

### 7.3 버전 체계

**Semantic Versioning (SemVer)**: `MAJOR.MINOR.PATCH`

| 구분 | 의미 | 예 |
|------|------|-----|
| MAJOR | 호환성 깨지는 변경 | 1.0.0 → 2.0.0 |
| MINOR | 하위 호환 기능 추가 | 1.1.0 → 1.2.0 |
| PATCH | 하위 호환 버그 수정 | 1.1.1 → 1.1.2 |

**의료기기 빌드 번호**: `v1.2.3+build.456`

### 7.4 형상 항목 (Configuration Items)

| 항목 | 경로/관리 |
|------|----------|
| 소스 코드 | `backend/`, `frontend/`, `rust-core/` |
| Proto 정의 | `backend/shared/proto/manpasik.proto` |
| DB 마이그레이션 | `backend/migrations/`, `infrastructure/database/init/` |
| 규정 문서 | `docs/compliance/` |
| 테스트 코드 | 각 패키지 `*_test.go`, `**/test/` |
| CI/CD 설정 | `.github/workflows/` |
| 인프라 매니페스트 | `infrastructure/kubernetes/`, `infrastructure/docker/` |

### 7.5 릴리스 태그

- 형식: `v{MAJOR}.{MINOR}.{PATCH}`
- 예: `v1.2.3`, `v1.2.3+build.456`

---

## 8. 도구 및 환경 (Clause 5.1.4)

### 8.1 개발 언어 및 프레임워크

| 도구 | 버전 | 용도 | 검증 수준 |
|------|------|------|----------|
| Go | 1.22+ | 백엔드 마이크로서비스 (30+ 서비스) | 컴파일러, `go vet` |
| Flutter/Dart | 3.x | 모바일 앱 (iOS/Android) | `flutter analyze` |
| Rust | 1.75+ | 코어 엔진 (차동측정, AI) | `cargo clippy` |
| Protocol Buffers | 3.x | gRPC 인터페이스 정의 | protoc 컴파일 |

### 8.2 데이터베이스 및 인프라

| 도구 | 버전 | 용도 |
|------|------|------|
| PostgreSQL | 16 | 메인 관계형 DB |
| TimescaleDB | (pg16) | 시계열 데이터 (측정 데이터) |
| Milvus | 2.4 | 벡터 DB (핑거프린트 검색) |
| Redis | 7 | 캐시, 세션 |
| Kafka (Redpanda) | 24.2 | 이벤트 스트리밍 |
| Elasticsearch | 8.14 | 전문 검색, 로그 |
| MinIO | latest | 오브젝트 스토리지 |
| Docker | 24.x | 컨테이너화 |
| Kubernetes | 1.28+ | 오케스트레이션 |

### 8.3 CI/CD

| 도구 | 용도 |
|------|------|
| GitHub Actions | CI (빌드, 테스트, 린트) |
| GitHub Actions | CD (Docker 빌드, K8s 배포) |
| Docker Hub / GHCR | 컨테이너 이미지 레지스트리 |

### 8.4 검증 명령어

| 언어 | 린트 | 테스트 | 빌드 |
|------|------|--------|------|
| Rust | `cargo clippy --all-targets` | `cargo test` | `cargo build` |
| Go | `golangci-lint run` | `go test ./...` | `go build ./...` |
| Dart | `dart analyze` | `flutter test` | `flutter build` |

---

## 9. 위험 관리 참조 (Clause 5.1.6, ISO 14971)

### 9.1 참조 문서

| 문서 | 경로 | 용도 |
|------|------|------|
| 위험관리 계획서 | `docs/compliance/iso14971-risk-management-plan.md` | ISO 14971 프로세스 |
| STRIDE 위협 모델 | `docs/security/stride-threat-model.md` (참조) | 사이버보안 위협 |
| SRS 위험 통제 요구사항 | `docs/compliance/iec62304-srs.md` §5 | 소프트웨어 통제 조치 |

### 9.2 위험 관리 통합

| 개발 단계 | 위험 관리 활동 |
|----------|---------------|
| 요구사항 | 위험 통제 요구사항 도출 (SRS §5) |
| 설계 | 아키텍처 위험 분석 (SAD §8) |
| 구현 | 보안 코딩 지침 준수 (OWASP, 입력 검증) |
| 검증 | 위험 통제 조치 검증 (테스트 케이스 연결) |

---

## 10. 문서화 계획 (Clause 5.1.8)

### 10.1 문서 체계

| 유형 | 경로 | 형식 |
|------|------|------|
| 규정 문서 | `docs/compliance/` | 마크다운 (.md) |
| 기술 문서 | `docs/`, `docs/plan/` | 마크다운 |
| 작업 기록 | 루트 | CHANGELOG.md, CONTEXT.md, KNOWN_ISSUES.md |

### 10.2 문서 검토 주기

| 문서 유형 | 검토 시점 |
|----------|----------|
| 규정 문서 (SDP, SRS, SAD) | Phase 완료 시 (최소 분기 1회) |
| 기술 문서 | Sprint 완료 시 |
| 작업 기록 | 매 작업 세션 (CHANGELOG.md) |

### 10.3 승인 프로세스

```
[작성자 자기 검토] → [기술 리더 기술 검토] → [품질 관리자 규정 검토] → [승인 + 버전 태그]
```

### 10.4 문서 변경 추적

- 문서 변경 이력은 **Git 커밋 히스토리**로 추적
- 규정 문서 수정 시 PR + 검토 필수

---

## 11. 문제 해결 프로세스 (Clause 5.1.10)

### 11.1 심각도 분류

| 심각도 | 설명 | 해결 SLA |
|--------|------|---------|
| **Critical** | 환자 안전 영향, 데이터 손실 | 4시간 |
| **Major** | 핵심 기능 중단 | 24시간 |
| **Minor** | 비핵심 기능 이상 | 1주 |
| **Trivial** | UI 개선, 문구 수정 | 다음 Sprint |

### 11.2 문제 추적

| 도구 | 용도 |
|------|------|
| GitHub Issues | 이슈 등록, 할당, 마일스톤 |
| KNOWN_ISSUES.md | 미해결 이슈 (삭제 금지, 해결 시 상태 변경) |

### 11.3 필수 정보

- 재현 단계, 환경, 예상/실제 결과, 스크린샷

---

## 12. SOUP 관리 (Clause 5.1.11)

### 12.1 SOUP 식별

| 언어 | 의존성 파일 |
|------|------------|
| Go | `go.mod` |
| Flutter | `pubspec.yaml` |
| Rust | `Cargo.toml` |

### 12.2 SOUP 검증 기준

| 항목 | 기준 |
|------|------|
| 라이선스 | 상용 사용 허용 (Apache-2.0, MIT, BSD 등) |
| 보안 | CVE 확인, 최신 패치 적용 |
| 품질 | 유지보수 활성도, 테스트 커버리지 |
| Class B 서브시스템 SOUP | 추가 위험 평가 |

### 12.3 SOUP 모니터링

- **분기별** 의존성 업데이트 검토
- **CVE 모니터링**: Dependabot, Renovate
- **주요 SOUP 업데이트** 시 회귀 테스트 실행

### 12.4 SOUP 목록 (요약)

**부록 A** 참조 — Go 10+, Flutter 7+, Rust 5+ 주요 라이브러리

---

## 부록 A: SOUP/OTS 주요 목록

### A.1 Go 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 | 안전 영향 |
|-----------|------|------|---------|----------|
| google.golang.org/grpc | 1.62+ | gRPC 통신 | Apache-2.0 | Class B |
| github.com/jackc/pgx/v5 | 5.x | PostgreSQL | MIT | Class B |
| github.com/go-redis/redis/v9 | 9.x | Redis 캐시 | BSD-2 | Class A |
| github.com/twmb/franz-go | 1.x | Kafka 클라이언트 | BSD-3 | Class A |
| github.com/golang-jwt/jwt/v5 | 5.x | JWT 처리 | MIT | Class B |
| go.uber.org/zap | 1.x | 구조화 로깅 | MIT | Class A |
| google.golang.org/protobuf | 1.x | Protobuf | BSD-3 | Class B |

### A.2 Flutter 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| grpc | 3.x | gRPC 클라이언트 | BSD-3 |
| protobuf | 3.x | Protobuf | BSD-3 |
| provider/riverpod | 6.x/2.x | 상태 관리 | MIT |
| go_router | 14.x | 라우팅 | BSD-3 |
| flutter_secure_storage | 9.x | 보안 저장 | BSD-3 |

### A.3 Rust 의존성

| 라이브러리 | 버전 | 용도 | 라이선스 |
|-----------|------|------|---------|
| flutter_rust_bridge | 2.x | FFI 브리지 | MIT |
| ring | 0.17+ | 암호화 | ISC |
| tokio | 1.x | 비동기 런타임 | MIT |
| serde | 1.x | 직렬화 | MIT/Apache-2.0 |

---

## 부록 B: 도구 검증 요약

| 도구 | 검증 방법 | 상태 |
|------|----------|------|
| Go 컴파일러 | 공식 릴리스, `go vet` | 검증됨 |
| Rust 컴파일러 | 공식 릴리스, `clippy` | 검증됨 |
| Dart/Flutter | 공식 릴리스, `flutter analyze` | 검증됨 |
| protoc | 공식 릴리스, 생성 코드 빌드 검증 | 검증됨 |
| Docker | 공식 이미지, 재현 가능한 빌드 | 검증됨 |
| GitHub Actions | 결정적 워크플로우, 로그 보존 | 검증됨 |

---

## 부록 C: IEC 62304 §5.1 체크리스트

| §5.1 요구사항 | 대응 섹션 | 상태 |
|--------------|----------|------|
| 5.1.1 목적·범위·참조 | §1, §2 | ✅ |
| 5.1.2 소프트웨어 안전 등급 | §3 | ✅ |
| 5.1.3 개발 프로세스·생명주기 | §4, §5 | ✅ |
| 5.1.4 개발 도구 | §8 | ✅ |
| 5.1.5 형상 관리 | §7 | ✅ |
| 5.1.6 위험 관리 참조 | §9 | ✅ |
| 5.1.7 검증 및 유효성 확인 | V&V 마스터 플랜 참조 | ✅ |
| 5.1.8 문서화 계획 | §10 | ✅ |
| 5.1.9 유지보수 | §6 | ✅ |
| 5.1.10 문제 해결 | §11 | ✅ |
| 5.1.11 SOUP | §12 | ✅ |

---

**마지막 업데이트**: 2026-02-13 (v1.1)  
**다음 검토**: Phase 완료 시 또는 분기 1회
