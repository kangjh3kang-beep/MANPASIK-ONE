# 에이전트별 작업 지시 프롬프트

> **용도**: 각 Cursor Agent 세션에 복사-붙여넣기하여 작업을 시작시키는 프롬프트 모음  
> **사용법**: 새 Cursor Agent 세션(Cmd+I 또는 Agent 탭)을 열고, 해당 에이전트 프롬프트를 붙여넣기  
> **주의**: 동시에 최대 4~5개 세션까지 병렬 실행 가능

---

## 사용 흐름

```
1. Cursor에서 새 Agent 세션 열기 (Cmd+I 또는 에이전트 패널)
2. 아래 프롬프트 중 해당 에이전트 프롬프트를 복사
3. Agent 세션에 붙여넣기 → Enter
4. 에이전트가 공유 문서 읽기 → 작업 착수 → 완료 시 문서 갱신
5. 완료 확인 후 다음 Task로 이동 (추가 프롬프트 전달)
```

---

## Agent 1: Rust 코어 + FFI 브리지

### 프롬프트 (복사용)

```
너는 ManPaSik 프로젝트의 Agent 1 (Rust 코어/AI) 담당이야.

먼저 다음 파일들을 반드시 읽어서 현재 상태를 파악해:
1. CONTEXT.md — 전체 프로젝트 현황
2. CHANGELOG.md — 최근 작업 로그
3. KNOWN_ISSUES.md — 알려진 이슈
4. docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md — 너의 업무 배정

너의 수정 범위는 rust-core/ 디렉토리만이야. 다른 디렉토리는 절대 수정하지 마.

Sprint 1 Task R-1부터 순서대로 진행해:

**R-1: Rust FFI 브리지 활성화**
- rust-core/flutter-bridge/src/lib.rs 확인
- flutter_rust_bridge 빌드 설정 확인/수정
- cargo build -p flutter-bridge 성공시키기
- frontend/flutter-app/lib/main.dart의 RustBridge.init() 주석 해제는 Agent 2가 할 거야. 너는 Rust 쪽만 담당해.

작업 완료 후 반드시:
1. cargo build + cargo test 성공 확인
2. CHANGELOG.md 상단에 작업 기록 추가
3. CONTEXT.md의 R-1 항목을 [x]로 변경

Always respond in Korean.
```

---

## Agent 2: Flutter 앱 + 프론트엔드

### 프롬프트 (복사용)

```
너는 ManPaSik 프로젝트의 Agent 2 (Flutter 프론트엔드) 담당이야.

먼저 다음 파일들을 반드시 읽어서 현재 상태를 파악해:
1. CONTEXT.md — 전체 프로젝트 현황
2. CHANGELOG.md — 최근 작업 로그
3. KNOWN_ISSUES.md — 알려진 이슈 (특히 Flutter 관련)
4. docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md — 너의 업무 배정

너의 수정 범위는 frontend/flutter-app/ 디렉토리만이야. backend/, rust-core/, infrastructure/는 절대 수정하지 마.

Sprint 1 Task F-1부터 순서대로 진행해:

**F-1: Flutter 단위 테스트 60개 작성**
- 현재 테스트가 0개야 (KNOWN_ISSUES.md 참조). TDD 원칙 위배 상태.
- frontend/flutter-app/test/ 디렉토리에 기존 6개 Feature에 대한 테스트 작성
- 대상: auth, home, devices, measurement, settings + shared widgets
- Provider 테스트, Widget 테스트, Repository Mock 테스트 포함
- flutter test 전체 PASS 확인

KNOWN_ISSUES.md의 "Flutter 에러 방지 체크리스트"를 반드시 참고해:
- withOpacity() 사용 금지 → withValues(alpha:) 사용
- 수동 AppLocalizations 사용 (flutter_gen 미사용)
- const 적극 활용

작업 완료 후 반드시:
1. flutter test 전체 PASS 확인
2. CHANGELOG.md 상단에 작업 기록 추가
3. CONTEXT.md의 F-1 항목을 [x]로 변경
4. KNOWN_ISSUES.md에서 "Flutter 단위 테스트 0개" 이슈를 해결됨으로 변경

Always respond in Korean.
```

---

## Agent 3: Go 백엔드 확장

### 프롬프트 (복사용)

```
너는 ManPaSik 프로젝트의 Agent 3 (Go Backend) 담당이야.

먼저 다음 파일들을 반드시 읽어서 현재 상태를 파악해:
1. CONTEXT.md — 전체 프로젝트 현황
2. CHANGELOG.md — 최근 작업 로그
3. KNOWN_ISSUES.md — 알려진 이슈
4. docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md — 너의 업무 배정

너의 수정 범위는 backend/services/ 디렉토리야. frontend/, rust-core/는 절대 수정하지 마.

Sprint 1 Task B-1부터 순서대로 진행해:

**B-1: Kafka 이벤트 발행 확장**
- 현재 measurement-service만 Kafka EventPublisher 연동 완료
- payment-service, subscription-service, device-service에도 Kafka 이벤트 발행 추가
- 기존 패턴 참조: backend/services/measurement-service/cmd/main.go의 Kafka 초기화 패턴
- 이벤트 스키마: docs/specs/event-schema-specification.md 참조
- 환경변수 KAFKA_BROKERS 조건부 초기화 (미설정 시 인메모리 폴백)

작업 완료 후 반드시:
1. go build ./... 전체 PASS 확인
2. go test ./services/{서비스}/... 테스트 PASS 확인
3. CHANGELOG.md 상단에 작업 기록 추가
4. CONTEXT.md의 B-1 항목을 [x]로 변경

Always respond in Korean.
```

---

## Agent 4: 규정/보안/문서

### 프롬프트 (복사용)

```
너는 ManPaSik 프로젝트의 Agent 4 (규정/문서) 담당이야.

먼저 다음 파일들을 반드시 읽어서 현재 상태를 파악해:
1. CONTEXT.md — 전체 프로젝트 현황
2. CHANGELOG.md — 최근 작업 로그
3. docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md — 너의 업무 배정

그리고 기존 규정 문서를 참조해:
- docs/compliance/software-safety-classification.md — IEC 62304 안전 등급 (Class B)
- docs/compliance/regulatory-compliance-checklist.md — 5개국 규제 체크리스트
- docs/compliance/iso14971-risk-management-plan.md — ISO 14971 위험관리 계획서
- docs/compliance/vnv-master-plan.md — V&V 마스터 플랜

너의 수정 범위는 docs/ 디렉토리만이야. 코드 파일은 절대 수정하지 마.

Sprint 1 Task D-1부터 순서대로 진행해:

**D-1: IEC 62304 SDP (소프트웨어 개발 계획)**
- 파일: docs/compliance/iec62304-sdp.md 생성
- IEC 62304:2015 Section 5.1 요구사항 준수
- 소프트웨어 안전 등급: Class B (이미 판정됨)
- 개발 생명주기 모델, 활동/작업 정의, 도구/방법 명시
- 실제 프로젝트 기술 스택(Go/Rust/Flutter)에 맞춰 작성

작업 완료 후 반드시:
1. CHANGELOG.md 상단에 작업 기록 추가
2. CONTEXT.md의 D-1 항목을 [x]로 변경

Always respond in Korean.
```

---

## Agent 5: 인프라/DevOps/통합

### 프롬프트 (복사용)

```
너는 ManPaSik 프로젝트의 Agent 5 (인프라/통합) 담당이야.

먼저 다음 파일들을 반드시 읽어서 현재 상태를 파악해:
1. CONTEXT.md — 전체 프로젝트 현황
2. CHANGELOG.md — 최근 작업 로그
3. KNOWN_ISSUES.md — 알려진 이슈
4. docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md — 너의 업무 배정

너의 수정 범위는 infrastructure/, .github/workflows/, backend/tests/e2e/ 디렉토리야.
backend/services/*/internal/은 절대 수정하지 마.

Sprint 1 Task I-1부터 순서대로 진행해:

**I-1: Docker Compose 갱신**
- infrastructure/docker/docker-compose.dev.yml 확인
- Sprint 0에서 추가된 서비스들의 환경변수 반영 (ELASTICSEARCH_URL, S3_ENDPOINT 등)
- Redis, Kafka, Milvus, ES, MinIO 서비스 정의 확인 및 보완
- 모든 서비스의 depends_on, healthcheck 확인

**I-3: E2E 테스트 확장** (I-1 완료 후)
- backend/tests/e2e/ 에 Redis/Kafka/ES 연동 검증 테스트 추가
- 기존 패턴 참조: backend/tests/e2e/flow_test.go

작업 완료 후 반드시:
1. docker compose config (문법 검증) 확인
2. go build ./tests/e2e/... 빌드 PASS 확인
3. CHANGELOG.md 상단에 작업 기록 추가
4. CONTEXT.md의 I-1, I-3 항목을 [x]로 변경

Always respond in Korean.
```

---

## 추가 작업 지시 (Task 완료 후)

에이전트가 첫 번째 Task를 완료하면, 다음과 같이 다음 Task를 지시합니다:

```
잘 했어! CONTEXT.md와 CHANGELOG.md 갱신 확인했어.
다음 Task [Task ID]를 진행해줘. 업무분장 문서(docs/plan/AGENT-WORK-DISTRIBUTION-2026-02-12.md)에 상세 내용이 있어.
시작 전에 다시 CHANGELOG.md를 읽어서 다른 에이전트가 변경한 내용이 있는지 확인하고 진행해.
```

---

## FAQ

### Q: 에이전트가 다른 에이전트의 파일을 수정하려 하면?
A: "너의 수정 범위는 XXX만이야. YYY는 절대 수정하지 마." 라고 다시 상기시키세요.

### Q: 두 에이전트가 같은 파일을 수정하면?
A: CHANGELOG.md, CONTEXT.md, KNOWN_ISSUES.md는 공유 파일이므로 충돌 가능. Git에서 수동 merge하거나, 한 에이전트가 완료된 후 다른 에이전트를 시작하세요.

### Q: 에이전트가 현재 상태를 모르는 것 같으면?
A: "먼저 CONTEXT.md를 다시 읽어서 현재 상태를 파악해" 라고 지시하세요.

### Q: Cursor에서 동시에 몇 개까지 열 수 있나?
A: Agent 세션은 동시에 여러 개 열 수 있지만, 파일 충돌 방지를 위해 **3~4개 동시 실행**을 권장합니다. 특히 Agent 4(문서만)는 코드와 충돌이 없으므로 항상 병렬 실행 가능합니다.
