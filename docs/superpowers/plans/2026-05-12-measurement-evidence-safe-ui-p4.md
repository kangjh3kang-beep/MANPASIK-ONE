# Measurement Evidence Safe UI P4 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** measurement evidence status를 규정 친화적 UI copy로 변환하고, 측정 골든패스 snapshot까지 전달해 화면이 안전하게 표시할 수 있게 한다.

**Architecture:** Flutter domain에 evidence presentation helper를 추가해 `research_only`를 진단/정상/위험 판정과 분리된 문구로 고정한다. Orchestrator는 repository result의 evidence fields를 snapshot으로 전달하고, MeasureScreen은 compact label만 표시한다.

**Tech Stack:** Flutter/Dart domain tests, measurement orchestrator tests, targeted Flutter analyze, existing Go/SSOT/security/assay gates.

---

## Stage Gate Rules

- 각 Task는 RED 테스트를 먼저 작성하고 실패를 확인한다.
- `research_only` UI copy는 "정상", "위험", "진단", "확정" 표현을 포함하지 않는다.
- UI 표시는 짧은 badge label만 연결하고, 상세 의료 문구/상담 유도는 별도 단계에서 검토한다.
- 저장 계층 영속화는 이번 P4 범위가 아니다.

## File Structure

- Create: `frontend/flutter-app/lib/features/measurement/domain/measurement_evidence_presentation.dart`
  - evidence status를 안전한 label/detail copy로 변환한다.
- Create: `frontend/flutter-app/test/features/measurement/domain/measurement_evidence_presentation_test.dart`
  - `research_only` copy가 진단 claim을 포함하지 않는지 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart`
  - `MeasurementGoldenPathSnapshot`에 evidence fields를 추가한다.
  - serverProcessed/sessionEnded snapshot에 server result evidence를 전달한다.
- Modify: `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`
  - orchestrator snapshot이 evidence fields를 전달하는지 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
  - 완료/서버 처리 상태의 compact badge에 안전한 evidence label을 표시한다.
- Create: `docs/audit/measurement-evidence-safe-ui-p4.md`
  - P4 구현 범위, TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - 작업 기록 프로토콜을 갱신한다.

## Task 1: Evidence Presentation Copy

**Files:**
- Create: `frontend/flutter-app/test/features/measurement/domain/measurement_evidence_presentation_test.dart`
- Create: `frontend/flutter-app/lib/features/measurement/domain/measurement_evidence_presentation.dart`

- [x] **Step 1: Write failing copy tests**

Create tests that require:

```dart
final copy = MeasurementEvidencePresentation.from(
  evidenceStatus: 'research_only',
  diagnosticReady: false,
  evidenceGaps: const ['clinical_lock_required'],
);
expect(copy.badgeLabel, '연구용');
expect(copy.detailText, contains('참고용'));
expect(copy.detailText, isNot(contains('정상')));
expect(copy.detailText, isNot(contains('위험')));
expect(copy.detailText, isNot(contains('확정')));
```

Also test missing status:

```dart
final copy = MeasurementEvidencePresentation.from(
  evidenceStatus: 'unknown',
  diagnosticReady: false,
  evidenceGaps: const [],
);
expect(copy.badgeLabel, '검증 확인 중');
```

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart
```

Expected: FAIL because `MeasurementEvidencePresentation` is not defined.

- [x] **Step 3: Implement helper**

Add immutable class with `badgeLabel`, `detailText`, and factory `from(...)`.

- [x] **Step 4: Verify GREEN**

Run the same Flutter test. Expected: PASS.

## Task 2: Orchestrator Snapshot Evidence Propagation

**Files:**
- Modify: `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart`
- Modify: `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`

- [x] **Step 1: Write failing snapshot test**

Update the orchestrator test to collect snapshots and assert:

```dart
final serverSnapshot = snapshots.firstWhere(
  (snapshot) => snapshot.phase == MeasurementGoldenPathPhase.serverProcessed,
);
expect(serverSnapshot.evidenceStatus, 'research_only');
expect(serverSnapshot.diagnosticReady, isFalse);
expect(serverSnapshot.evidenceGaps, contains('clinical_lock_required'));
```

Update fake repository result with evidence fields.

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
```

Expected: FAIL because snapshot has no evidence fields.

- [x] **Step 3: Add snapshot fields and emit mapping**

Add `evidenceStatus`, `diagnosticReady`, `evidenceGaps` to `MeasurementGoldenPathSnapshot`, and fill them from `serverResult` for `serverProcessed` and `sessionEnded`.

- [x] **Step 4: Verify GREEN**

Run the same orchestrator test. Expected: PASS.

## Task 3: Measure Screen Badge Copy

**Files:**
- Modify: `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`

- [x] **Step 1: Wire compact label**

Import `measurement_evidence_presentation.dart` and update `_measureStatusText`:

```dart
case MeasurementGoldenPathPhase.serverProcessed:
case MeasurementGoldenPathPhase.sessionEnded:
  if (snapshot.evidenceStatus != null) {
    return MeasurementEvidencePresentation.from(
      evidenceStatus: snapshot.evidenceStatus!,
      diagnosticReady: snapshot.diagnosticReady,
      evidenceGaps: snapshot.evidenceGaps,
    ).badgeLabel;
  }
  return snapshot.phase == MeasurementGoldenPathPhase.sessionEnded ? 'DONE' : 'SERVER';
```

- [x] **Step 2: Verify analyze**

Run targeted analyzer. Expected: PASS.

## Task 4: Audit and Final Verification

**Files:**
- Create: `docs/audit/measurement-evidence-safe-ui-p4.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Record audit and logs**

Document P4 scope, TDD records, code review, quality gates, and next-stage guidance.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/measurement/domain/measurement_evidence_presentation.dart lib/features/measurement/application/measurement_golden_path_orchestrator.dart lib/features/measurement/presentation/measure_screen.dart test/features/measurement/domain/measurement_evidence_presentation_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath
```

Expected: PASS.
