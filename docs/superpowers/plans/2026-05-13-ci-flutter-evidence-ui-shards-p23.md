# CI Flutter Evidence UI Shards Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P22에서 만든 Flutter evidence UI category shard를 GitHub Actions CI step으로 분리해 실패 위치와 CI 비용 관리를 더 명확히 한다.

**Architecture:** 기존 `flutter-app` job은 `Run tests` 이후 aggregate evidence gate를 한 번 실행했다. 이번 단계는 이를 `contract`, `ui`, `data` 세 step으로 바꾸고, workflow guard가 category별 `--category` 호출과 aggregate 단독 호출 제거를 검증하게 한다.

**Tech Stack:** Bash, GitHub Actions workflow, Flutter test gate scripts.

---

### Task 1: CI shard workflow guard

**Files:**
- Create: `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`
- Modify: `.github/workflows/ci.yml`

- [x] **Step 1: Write failing workflow guard**

Create `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`.

The guard must require:
- `Flutter evidence UI gate: contract`
- `Flutter evidence UI gate: ui`
- `Flutter evidence UI gate: data`
- `scripts/flutter_evidence_ui_gate.sh --category contract`
- `scripts/flutter_evidence_ui_gate.sh --category ui`
- `scripts/flutter_evidence_ui_gate.sh --category data`

The guard must fail if the workflow still contains the aggregate-only command:

```yaml
run: bash ../../scripts/flutter_evidence_ui_gate.sh
```

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`

Expected: FAIL because CI still runs the aggregate gate.

- [x] **Step 3: Split CI steps**

Replace the aggregate evidence gate step in `.github/workflows/ci.yml`:

```yaml
      - name: Flutter evidence UI gate
        run: bash ../../scripts/flutter_evidence_ui_gate.sh
```

with:

```yaml
      - name: Flutter evidence UI gate: contract
        run: bash ../../scripts/flutter_evidence_ui_gate.sh --category contract

      - name: Flutter evidence UI gate: ui
        run: bash ../../scripts/flutter_evidence_ui_gate.sh --category ui

      - name: Flutter evidence UI gate: data
        run: bash ../../scripts/flutter_evidence_ui_gate.sh --category data
```

- [x] **Step 4: Run workflow guard to verify GREEN**

Run:
`bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/ci-flutter-evidence-ui-shards-p23.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted workflow and policy guards**

Run:
`bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_shard_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Run:
`bash -n scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_shard_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh`

- [x] **Step 2: Run actual Flutter shard gates**

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category contract`

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category ui`

Run:
`bash scripts/flutter_evidence_ui_gate.sh --category data`

- [x] **Step 3: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 4: Run diff hygiene checks**

Run:
`git diff --check -- .github/workflows/ci.yml scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh docs/superpowers/plans/2026-05-13-ci-flutter-evidence-ui-shards-p23.md docs/audit/ci-flutter-evidence-ui-shards-p23.md CHANGELOG.md CONTEXT.md`
