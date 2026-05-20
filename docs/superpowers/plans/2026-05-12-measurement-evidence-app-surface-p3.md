# Measurement Evidence App Surface P3 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P2에서 gRPC/Dart proto에 추가한 measurement evidence fields를 REST gateway와 Flutter repository result까지 노출한다.

**Architecture:** Gateway는 이미 `writeProtoJSON(... UseProtoNames: true)`로 proto response를 전달하므로 route contract test로 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 잠근다. Flutter는 `ProcessMeasurementResult` 도메인 모델에 evidence fields를 추가하고 native gRPC/REST mapper가 같은 값을 채우게 한다.

**Tech Stack:** Go gateway handler tests, Dart domain model, Flutter repository tests, Flutter analyzer, existing SSOT/security/assay gates.

---

## Stage Gate Rules

- 각 Task는 RED 테스트를 먼저 작성하고 실패를 확인한다.
- `research_only`를 진단 가능 표현으로 변환하지 않는다.
- Flutter 기본값은 legacy server 호환을 위해 `evidenceStatus='unknown'`, `diagnosticReady=false`, `evidenceGaps=[]`로 둔다.
- UI 표시 문구는 이번 P3에서 다루지 않고, 별도 단계에서 규정 문구 검토 후 진행한다.

## File Structure

- Modify: `backend/services/gateway/internal/handler/e2e_test.go`
  - REST `/api/v1/measurements/process` response가 evidence fields를 포함하는지 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
  - `ProcessMeasurementResult`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가한다.
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - native gRPC `MeasurementResult`에서 evidence fields를 복사한다.
- Modify: `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - native gRPC repository result가 evidence fields를 보존하는지 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart`
  - REST response의 snake_case와 camelCase evidence keys를 모두 읽는다.
- Create: `frontend/flutter-app/test/features/measurement/data/measurement_process_gateway_mapper_test.dart`
  - REST mapper가 evidence fields를 decode하는지 검증한다.
- Create: `docs/audit/measurement-evidence-app-surface-p3.md`
  - P3 구현 범위, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - 작업 기록 프로토콜에 따라 최신 상태와 다음 단계 지침을 반영한다.

## Task 1: Gateway REST Evidence Contract

**Files:**
- Modify: `backend/services/gateway/internal/handler/e2e_test.go`

- [x] **Step 1: Write failing gateway route test**

Add these assertions to `TestE2E_ProcessMeasurementGoldenPath` after the existing JSON key assertions:

```go
assertJSONKey(t, body, "evidence_status")
assertJSONKey(t, body, "diagnostic_ready")
assertJSONKey(t, body, "evidence_gaps")
if !strings.Contains(body, "research_only") {
    t.Fatalf("response body does not include research_only evidence status: %s", body)
}
```

- [x] **Step 2: Verify RED**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath
```

Expected: FAIL because the mock stream response does not yet include evidence fields.

- [x] **Step 3: Update mock stream response**

In `mockMeasurementStream.Recv()`, add:

```go
EvidenceStatus:    "research_only",
DiagnosticReady:   false,
EvidenceGaps:      []string{"clinical_lock_required"},
```

- [x] **Step 4: Verify GREEN**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/gofmt -w services/gateway/internal/handler/e2e_test.go
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath
```

Expected: PASS.

## Task 2: Flutter Native Repository Evidence Surface

**Files:**
- Modify: `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
- Modify: `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`

- [x] **Step 1: Write failing native repository test**

In `measurement_repository_impl_test.dart`, update the fake `MeasurementResult` to include:

```dart
evidenceStatus: 'research_only',
diagnosticReady: false,
evidenceGaps: ['clinical_lock_required'],
```

Then assert:

```dart
expect(result.evidenceStatus, 'research_only');
expect(result.diagnosticReady, isFalse);
expect(result.evidenceGaps, contains('clinical_lock_required'));
```

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart
```

Expected: FAIL because `ProcessMeasurementResult` has no evidence fields.

- [x] **Step 3: Add domain fields and native mapping**

In `ProcessMeasurementResult`, add:

```dart
final String evidenceStatus;
final bool diagnosticReady;
final List<String> evidenceGaps;
```

Constructor defaults:

```dart
this.evidenceStatus = 'unknown',
this.diagnosticReady = false,
this.evidenceGaps = const [],
```

In `MeasurementRepositoryImpl.processMeasurement`, map:

```dart
evidenceStatus: res.evidenceStatus.isNotEmpty ? res.evidenceStatus : 'unknown',
diagnosticReady: res.diagnosticReady,
evidenceGaps: List.unmodifiable(res.evidenceGaps),
```

- [x] **Step 4: Verify GREEN**

Run the same Flutter test. Expected: PASS.

## Task 3: Flutter REST Mapper Evidence Surface

**Files:**
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart`
- Create: `frontend/flutter-app/test/features/measurement/data/measurement_process_gateway_mapper_test.dart`

- [x] **Step 1: Write failing REST mapper test**

Create a test that calls `decodeProcessMeasurementResult()` with:

```dart
{
  'session_id': 'session-rest-1',
  'primary_value': 88.1,
  'unit': 'mg/dL',
  'confidence': 0.91,
  'evidence_status': 'research_only',
  'diagnostic_ready': false,
  'evidence_gaps': ['clinical_lock_required'],
}
```

Assert the decoded `ProcessMeasurementResult` has the same evidence values.

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_process_gateway_mapper_test.dart
```

Expected: FAIL because mapper/domain evidence fields are not fully wired.

- [x] **Step 3: Implement REST mapper**

Read both snake_case and camelCase:

```dart
final evidenceStatus =
    response['evidence_status'] as String? ??
    response['evidenceStatus'] as String? ??
    'unknown';
final diagnosticReady =
    response['diagnostic_ready'] as bool? ??
    response['diagnosticReady'] as bool? ??
    false;
final evidenceGapsRaw =
    response['evidence_gaps'] as List<dynamic>? ??
    response['evidenceGaps'] as List<dynamic>? ??
    const <dynamic>[];
```

Map `evidenceGapsRaw.map((value) => value.toString()).toList(growable: false)`.

- [x] **Step 4: Verify GREEN**

Run the mapper test. Expected: PASS.

## Task 4: Audit and Final Verification

**Files:**
- Create: `docs/audit/measurement-evidence-app-surface-p3.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Record audit**

Document Task 1~3 TDD records, code review, quality gates, and residual risks.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath
cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_process_gateway_mapper_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_impl.dart lib/features/measurement/data/measurement_process_gateway_mapper.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_process_gateway_mapper_test.dart
```

Expected: PASS.
