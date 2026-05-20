# CI Gate Matrix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 릴리스/보안/evidence/Wasm 관련 CI gate의 실행 위치, blocking 여부, policy 문서, guard script를 한 matrix로 고정한다.

**Architecture:** `docs/ci/ci-gate-matrix.md`는 gate별 `scope`, `workflow_job`, `blocking`, `policy_doc`, `guard_script`, `execution_mode` marker를 포함한다. `scripts/ci_gate_matrix_policy_test.sh`는 matrix 문서가 핵심 gate를 모두 포함하고, 각 policy/guard script가 실제 파일로 존재하며, workflow에 주요 step이 연결되어 있는지 검증한다.

**Tech Stack:** Bash, Markdown CI policy, GitHub Actions workflow.

---

### Task 1: Matrix guard

**Files:**
- Create: `scripts/ci_gate_matrix_policy_test.sh`
- Create: `docs/ci/ci-gate-matrix.md`

- [x] **Step 1: Write failing matrix guard**

Create `scripts/ci_gate_matrix_policy_test.sh`.

The guard must require:
- `docs/ci/ci-gate-matrix.md`
- matrix markers:
  - `matrix_version: 1`
  - `ssot-governance`
  - `Flutter evidence UI gate: contract`
  - `Flutter evidence UI gate: ui`
  - `Flutter evidence UI gate: data`
  - `Flutter web release gate`
  - `Wasm readiness policy`
  - `Security release gate`
  - `Assay evidence gate`
  - `blocking: true`
  - `blocking: false`
  - `execution_mode: category_shards`
  - `execution_mode: pull_request_and_release_branch`
  - `release_target: js_web`
  - `wasm_release_target: false`
- existing policy docs:
  - `docs/ci/flutter-evidence-ui-gate-policy.md`
  - `docs/ci/flutter-web-release-gate-policy.md`
  - `docs/ci/flutter-wasm-readiness.md`
- existing guard scripts:
  - `scripts/flutter_evidence_ui_gate_policy_test.sh`
  - `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`
  - `scripts/ci_web_gate_policy_test.sh`
  - `scripts/ci_wasm_readiness_policy_test.sh`
  - `scripts/security_release_gate.sh`
  - `scripts/assay_evidence_gate.sh`
- workflow markers:
  - `SSOT Governance`
  - `Flutter evidence UI gate: contract`
  - `Flutter evidence UI gate: ui`
  - `Flutter evidence UI gate: data`
  - `Flutter web release gate`

- [x] **Step 2: Run guard to verify RED**

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: FAIL because `docs/ci/ci-gate-matrix.md` does not exist.

- [x] **Step 3: Create matrix document**

Create `docs/ci/ci-gate-matrix.md` with:
- YAML marker block.
- A table for governance, evidence, web release, Wasm readiness, security, assay evidence.
- A short change rule: new blocking gates need policy doc, guard script, workflow connection, and CHANGELOG/CONTEXT update.

- [x] **Step 4: Run guard to verify GREEN**

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/ci-gate-matrix-p26.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_wasm_readiness_policy_test.sh`

Run:
`bash -n scripts/ci_gate_matrix_policy_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_wasm_readiness_policy_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- scripts/ci_gate_matrix_policy_test.sh docs/ci/ci-gate-matrix.md docs/superpowers/plans/2026-05-13-ci-gate-matrix-p26.md docs/audit/ci-gate-matrix-p26.md CHANGELOG.md CONTEXT.md`
