# DataHub Layout Smoke Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub evidence badge와 trend header가 좁은 모바일 화면에서도 overflow 없이 표시되게 한다.

**Architecture:** 화면 문구와 evidence semantics는 P10-P12 구현을 그대로 유지한다. 레이아웃만 보강하며, 긴 metric label과 badge/period chip이 같은 행에서 충돌하지 않도록 title 영역을 flexible하게 만들고 summary pill row도 작은 화면에서 줄바꿈 가능하게 한다.

**Tech Stack:** Flutter widget tests, Riverpod provider overrides, ManPaSik DataHub UI.

---

### Task 1: Narrow mobile layout smoke

**Files:**
- Create: `frontend/flutter-app/test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart`

- [ ] **Step 1: Write the failing test**

Create a widget test that renders `DataHubScreen` at 320px width with a long selected metric name and `research_only` evidence. The test must assert:
- `연구용` badge is visible.
- `tester.takeException()` is `null` after pump.

- [ ] **Step 2: Run test to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`

Expected: FAIL with a Flutter layout exception or RenderFlex overflow caused by long title + badge + period chip.

- [ ] **Step 3: Implement minimal layout fix**

In `_HeroChartCard`, place metric title and evidence badge inside a flexible `Wrap`/`Expanded` region before the period chip. Keep the period chip visible. In chart summary, use `Wrap` instead of a fixed-width `Row` so three summary pills can wrap on narrow widths.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`

Expected: PASS.

### Task 2: Regression gates

**Files:**
- Create: `docs/audit/datahub-layout-smoke-p13.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run focused Flutter tests**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`

- [ ] **Step 2: Run targeted analyzer**

Run:
`cd frontend/flutter-app && flutter analyze --no-pub --no-fatal-infos lib/features/data_hub/presentation/data_hub_screen.dart test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`

- [ ] **Step 3: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 4: Run build**

Run:
`cd frontend/flutter-app && flutter build web --no-pub`

- [ ] **Step 5: Run diff hygiene checks**

Run:
`git diff --check -- frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart frontend/flutter-app/test/features/data_hub/presentation/data_hub_layout_smoke_test.dart CHANGELOG.md CONTEXT.md`
