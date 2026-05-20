# DataHub Evidence UI And Latest Selection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub 화면이 측정 evidence 상태를 안전한 배지로 표시하고, 요약의 최신 evidence가 `measured_at` 기준으로 안정적으로 선택되게 한다.

**Architecture:** P4/P7에서 만든 `MeasurementEvidenceBadge`와 `MeasurementEvidencePresentation`을 재사용해 새로운 의료 판정 문구를 만들지 않는다. REST repository는 history item을 `TrendDataPoint`로 표준화한 뒤 timestamp 오름차순으로 정렬하고, summary는 정렬된 latest point를 기준으로 evidence metadata를 전달한다.

**Tech Stack:** Flutter, Riverpod, Dart widget/domain tests, ManPaSik REST gateway mapping.

---

### Task 1: DataHub 화면 evidence badge 표시

**Files:**
- Create: `frontend/flutter-app/test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart`

- [ ] **Step 1: Write the failing widget test**

`DataHubScreen`에 `research_only` summary를 주입하고 `연구용` 배지가 보이는지 검증한다. 동시에 `정상`, `위험`, `진단`, `확정` 문구가 나오지 않아야 한다.

- [ ] **Step 2: Run test to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`

Expected: FAIL because `DataHubScreen` does not render an evidence badge yet.

- [ ] **Step 3: Implement minimal UI**

Import `MeasurementEvidenceBadge` in `data_hub_screen.dart`. Render it from `_HeroChartCard` and `_DetailPanel` using `BiomarkerSummary.latestEvidenceStatus`, `latestDiagnosticReady`, and `latestEvidenceGaps`.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`

Expected: PASS.

### Task 2: timestamp 기반 latest evidence 선택

**Files:**
- Modify: `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`

- [ ] **Step 1: Write failing REST tests**

Add tests with out-of-order history items. The newer `measured_at` item must decide `getBiomarkerSummary.latestEvidenceStatus`, `latestDiagnosticReady`, `latestEvidenceGaps`, and `getTrendData` must return timestamp-sorted points.

- [ ] **Step 2: Run tests to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: FAIL because the repository currently preserves response order in trend data and uses response/value order for latest summary evidence.

- [ ] **Step 3: Implement deterministic mapping**

Sort `getTrendData` results by `timestamp`. In `getAllBiomarkerSummaries`, track latest point per biomarker by comparing parsed `measured_at` timestamps instead of first response order. Keep missing timestamps conservative by ignoring latest replacement unless a parsed timestamp exists.

- [ ] **Step 4: Run tests to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: PASS.

### Task 3: Documentation and quality gates

**Files:**
- Create: `docs/audit/datahub-evidence-ui-and-latest-p10-p11.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run focused Flutter tests**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`

- [ ] **Step 2: Run targeted analyzer**

Run:
`cd frontend/flutter-app && flutter analyze --no-pub --no-fatal-infos lib/features/data_hub/presentation/data_hub_screen.dart lib/features/data_hub/data/data_hub_repository_rest.dart lib/features/data_hub/domain/data_hub_repository.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart`

- [ ] **Step 3: Run project gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 4: Run diff hygiene checks**

Run:
`git diff --check -- frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart frontend/flutter-app/test/features/data_hub/presentation/data_hub_evidence_badge_test.dart frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart CHANGELOG.md CONTEXT.md`

- [ ] **Step 5: Update audit and shared logs**

Document RED/GREEN evidence, self-review notes, quality gate output, and next stage recommendation.
