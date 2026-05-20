# Flutter Evidence UI Gate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** measurement evidence contract, safe copy, UI badge, history/app/DataHub mapping tests를 하나의 실행 가능한 Flutter evidence UI gate로 묶고 CI에서 실행한다.

**Architecture:** 기존 테스트를 재사용한다. 새 gate script는 `frontend/flutter-app`에서 evidence 관련 test file 목록만 실행한다. `--list` 모드는 CI/문서 guard가 테스트 목록을 빠르게 확인할 수 있게 한다.

**Tech Stack:** Bash, Flutter tests, GitHub Actions workflow.

---

### Task 1: Evidence gate script

**Files:**
- Create: `scripts/flutter_evidence_ui_gate_test.sh`
- Create: `scripts/flutter_evidence_ui_gate.sh`

- [x] **Step 1: Write the failing script test**

Create a shell test that runs `scripts/flutter_evidence_ui_gate.sh --list` and requires these paths:
- `test/features/measurement/domain/measurement_evidence_presentation_test.dart`
- `test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`
- `test/features/home/presentation/home_measurement_evidence_badge_test.dart`
- `test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`
- `test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`
- `test/features/measurement/data/measurement_repository_rest_test.dart`
- `test/features/data_hub/data/data_hub_repository_rest_test.dart`
- `test/generated/measurement_result_evidence_contract_test.dart`
- `test/generated/measurement_summary_evidence_contract_test.dart`

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_evidence_ui_gate_test.sh`

Expected: FAIL because `scripts/flutter_evidence_ui_gate.sh` does not exist.

- [x] **Step 3: Implement gate script**

The script must:
- Support `--list`.
- Use `FLUTTER_BIN` override or fallback to `/mnt/d/우리집/flutter_cache/flutter/bin/flutter`.
- Run `flutter test --no-pub` with the evidence-related test file list.
- Print `FLUTTER_EVIDENCE_UI_GATE_PASS` after tests pass.

- [x] **Step 4: Run script test to verify GREEN**

Run:
`bash scripts/flutter_evidence_ui_gate_test.sh`

Expected: PASS.

### Task 2: CI connection

**Files:**
- Create: `scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`
- Modify: `.github/workflows/ci.yml`

- [x] **Step 1: Write failing CI guard**

Create a shell test requiring `.github/workflows/ci.yml` to contain:
- `Flutter evidence UI gate`
- `scripts/flutter_evidence_ui_gate.sh`

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Expected: FAIL because CI does not run the evidence gate yet.

- [x] **Step 3: Add CI step**

Add this step in the `flutter-app` job after `Run tests` and before `Flutter web release gate`:

```yaml
      - name: Flutter evidence UI gate
        run: bash ../../scripts/flutter_evidence_ui_gate.sh
```

- [x] **Step 4: Run CI guard to verify GREEN**

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Expected: PASS.

### Task 3: Quality and logs

**Files:**
- Create: `docs/audit/flutter-evidence-ui-gate-p20.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run gate and script checks**

Run:
`bash scripts/flutter_evidence_ui_gate_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Run:
`bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- .github/workflows/ci.yml scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh CHANGELOG.md CONTEXT.md`
