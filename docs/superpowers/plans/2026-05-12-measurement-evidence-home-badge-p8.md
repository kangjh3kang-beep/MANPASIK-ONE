# Measurement Evidence Home Badge P8 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans and superpowers:test-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Home dashboard의 최근 측정 요약에도 P7 `MeasurementEvidenceBadge`를 표시해, 결과 화면 밖에서도 `research_only` 상태가 일관되게 보이도록 한다.

**Architecture:** `homeDashboardProvider.latestMeasurement`는 P6 `MeasurementHistoryItem`을 그대로 전달한다. Home hero card는 해당 evidence fields를 `MeasurementEvidenceBadge`에 넘겨 compact badge만 렌더링한다.

**Tech Stack:** Flutter screen widget test, targeted Flutter analyze, existing SSOT/security/assay gates.

---

## Stage Gate Rules

- RED screen widget test를 먼저 작성하고 실패를 확인한다.
- Home 화면에서도 `research_only`는 `연구용`으로만 표시한다.
- 진단/정상/위험/확정 표현은 추가하지 않는다.
- DataHub 확장은 이번 P8 범위가 아니다.

## File Structure

- Modify: `frontend/flutter-app/lib/features/home/presentation/home_screen.dart`
  - `MeasurementEvidenceBadge`를 import하고 hero card 최근 측정 아래에 표시한다.
- Create: `frontend/flutter-app/test/features/home/presentation/home_measurement_evidence_badge_test.dart`
  - HomeScreen이 latest measurement evidence badge를 표시하는지 검증한다.
- Create: `docs/audit/measurement-evidence-home-badge-p8.md`
- Modify: `CHANGELOG.md`, `CONTEXT.md`

## Task 1: Home Badge Screen Test

- [x] **Step 1: Write failing screen test**

Pump `HomeScreen` with `homeDashboardProvider` override containing a `research_only` latest measurement, then assert `연구용` renders and forbidden copy does not.

- [x] **Step 2: Verify RED**

Run:

```bash
cd frontend/flutter-app
/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart
```

Expected: FAIL because HomeScreen does not render evidence badge yet.

- [x] **Step 3: Wire badge into HomeScreen**

Render `MeasurementEvidenceBadge` in `_HeroBentoCard` when latest measurement exists.

- [x] **Step 4: Verify GREEN**

Run the same widget test. Expected: PASS.

## Task 2: Audit and Final Gates

- [x] **Step 1: Record audit and logs**

Document TDD record, self-review, quality gates, and next-stage guidance.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/home/presentation/home_screen.dart lib/features/measurement/presentation/widgets/measurement_evidence_badge.dart test/features/home/presentation/home_measurement_evidence_badge_test.dart
```

Expected: PASS.
