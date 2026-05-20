# Measurement Evidence History Badge P7 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans and superpowers:test-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P6에서 Flutter history item까지 전달된 evidence status를 measurement result 화면의 최신 측정 카드에 규정 친화적인 compact badge로 표시한다.

**Architecture:** P4의 `MeasurementEvidencePresentation`을 재사용하는 presentation widget을 추가하고, `MeasurementResultScreen._buildLatestCard()`에서 `MeasurementHistoryItem` evidence fields로 badge를 렌더링한다. UI copy는 badge label만 표시하며 상세 의료 문구나 판정 문구는 추가하지 않는다.

**Tech Stack:** Flutter widget test, targeted Flutter analyze, existing SSOT/security/assay gates.

---

## Stage Gate Rules

- RED widget test를 먼저 작성하고 실패를 확인한다.
- `research_only` badge는 "정상", "위험", "진단", "확정" 표현을 포함하지 않는다.
- 표시 문구는 P4 helper에서만 생성한다.
- DataHub 전체 표시 확장은 이번 P7 범위가 아니다.

## File Structure

- Create: `frontend/flutter-app/lib/features/measurement/presentation/widgets/measurement_evidence_badge.dart`
  - `MeasurementEvidencePresentation`을 사용해 compact badge를 렌더링한다.
- Create: `frontend/flutter-app/test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`
  - `research_only` badge copy와 금지 표현 부재를 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/presentation/measurement_result_screen.dart`
  - 최신 측정 카드에 badge를 표시한다.
- Create: `docs/audit/measurement-evidence-history-badge-p7.md`
- Modify: `CHANGELOG.md`, `CONTEXT.md`

## Task 1: Evidence Badge Widget

- [x] **Step 1: Write failing widget test**

Test that a `research_only` badge renders `연구용` and does not render diagnostic/normal/risk/confirmed language.

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart
```

Expected: FAIL because widget file/class is missing.

- [x] **Step 3: Implement badge widget**

Add a compact Material-friendly badge using P4 presentation helper.

- [x] **Step 4: Verify GREEN**

Run the same widget test. Expected: PASS.

## Task 2: Result Screen Integration

- [x] **Step 1: Wire badge into latest card**

Import the badge widget and render it below cartridge type or value in `_buildLatestCard()`.

- [x] **Step 2: Verify analyze**

Run targeted Flutter analyze.

## Task 3: Audit and Final Gates

- [x] **Step 1: Record audit and logs**

Document TDD record, self-review, quality gates, and next-stage guidance.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/measurement/domain/measurement_evidence_presentation.dart lib/features/measurement/presentation/widgets/measurement_evidence_badge.dart lib/features/measurement/presentation/measurement_result_screen.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart
```

Expected: PASS.
