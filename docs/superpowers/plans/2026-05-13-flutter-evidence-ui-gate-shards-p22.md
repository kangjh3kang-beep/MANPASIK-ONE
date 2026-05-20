# Flutter Evidence UI Gate Shards Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P20/P21의 Flutter evidence UI gate를 `contract`, `ui`, `data` category shard로 분리해 test scope와 CI 비용을 더 작게 관리한다.

**Architecture:** `scripts/flutter_evidence_ui_gate.sh`는 기본 실행 시 기존 13개 테스트를 모두 실행한다. `--category <contract|ui|data>`와 `--list-categories`, category-aware `--list`/`--count`를 추가해 shard별 목록과 실행을 지원한다. Policy guard는 전체 상한 13개와 shard별 상한 5개를 함께 검증한다.

**Tech Stack:** Bash, Flutter tests, Markdown CI policy, GitHub Actions workflow.

---

### Task 1: Shard guard

**Files:**
- Create: `scripts/flutter_evidence_ui_gate_shard_test.sh`
- Modify: `scripts/flutter_evidence_ui_gate.sh`

- [x] **Step 1: Write failing shard guard**

Create `scripts/flutter_evidence_ui_gate_shard_test.sh`.

The guard must verify:
- `bash scripts/flutter_evidence_ui_gate.sh --list-categories` contains `contract`, `ui`, `data`.
- `bash scripts/flutter_evidence_ui_gate.sh --count` returns `13`.
- `bash scripts/flutter_evidence_ui_gate.sh --count --category contract` returns `4`.
- `bash scripts/flutter_evidence_ui_gate.sh --count --category ui` returns `4`.
- `bash scripts/flutter_evidence_ui_gate.sh --count --category data` returns `5`.
- Shard counts sum to the full count.
- `--list --category contract` includes:
  - `test/features/measurement/domain/measurement_evidence_presentation_test.dart`
  - `test/generated/measurement_result_evidence_contract_test.dart`
  - `test/generated/measurement_summary_evidence_contract_test.dart`
- `--list --category ui` includes:
  - `test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`
  - `test/features/home/presentation/home_measurement_evidence_badge_test.dart`
  - `test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`
- `--list --category data` includes:
  - `test/features/measurement/data/measurement_repository_impl_test.dart`
  - `test/features/measurement/data/measurement_repository_rest_test.dart`
  - `test/features/data_hub/data/data_hub_repository_rest_test.dart`
- Unknown category fails.

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_evidence_ui_gate_shard_test.sh`

Expected: FAIL because `--category` and `--list-categories` are not implemented.

- [x] **Step 3: Implement category shard support**

Update `scripts/flutter_evidence_ui_gate.sh` so:
- `--list-categories` prints `contract`, `ui`, `data`.
- `--list --category <name>` prints only that shard.
- `--count --category <name>` prints only that shard count.
- `--category <name>` runs only that shard.
- Unknown category exits non-zero with `FLUTTER_EVIDENCE_UI_GATE_FAIL unknown category`.
- Default behavior still runs all 13 tests and prints `FLUTTER_EVIDENCE_UI_GATE_PASS`.

- [x] **Step 4: Run shard guard to verify GREEN**

Run:
`bash scripts/flutter_evidence_ui_gate_shard_test.sh`

Expected: PASS.

### Task 2: Policy guard expansion

**Files:**
- Modify: `docs/ci/flutter-evidence-ui-gate-policy.md`
- Modify: `scripts/flutter_evidence_ui_gate_policy_test.sh`

- [x] **Step 1: Write failing policy expectations**

Update `scripts/flutter_evidence_ui_gate_policy_test.sh` to require:
- `category_shards: contract,ui,data`
- `max_test_files_per_category: 5`
- Each category count is numeric and does not exceed `max_test_files_per_category`.

- [x] **Step 2: Run policy guard to verify RED**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: FAIL because policy markers are not present yet.

- [x] **Step 3: Update policy document**

Update `docs/ci/flutter-evidence-ui-gate-policy.md` with:
- `category_shards: contract,ui,data`
- `max_test_files_per_category: 5`
- A short policy note that new tests must join the closest category or trigger shard split review.

- [x] **Step 4: Run policy guard to verify GREEN**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: PASS.

### Task 3: Quality and logs

**Files:**
- Create: `docs/audit/flutter-evidence-ui-gate-shards-p22.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/flutter_evidence_ui_gate_shard_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_test.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Run:
`bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh scripts/flutter_evidence_ui_gate_shard_test.sh`

- [x] **Step 2: Run actual Flutter gate paths**

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category contract`

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category ui`

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category data`

Run:
`bash scripts/flutter_evidence_ui_gate.sh`

- [x] **Step 3: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 4: Run diff hygiene checks**

Run:
`git diff --check -- scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_shard_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh docs/ci/flutter-evidence-ui-gate-policy.md docs/superpowers/plans/2026-05-13-flutter-evidence-ui-gate-shards-p22.md docs/audit/flutter-evidence-ui-gate-shards-p22.md CHANGELOG.md CONTEXT.md`
