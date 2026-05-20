# Flutter Evidence UI Gate Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P20에서 추가한 Flutter evidence UI gate의 운영 정책과 비용 상한을 문서/스크립트/CI guard로 고정한다.

**Architecture:** `scripts/flutter_evidence_ui_gate.sh`는 `--count` 모드로 test file 수를 노출한다. 새 policy test는 문서 marker, gate count 상한, CI 연결을 검증한다.

**Tech Stack:** Bash, GitHub Actions workflow, Markdown policy.

---

### Task 1: Policy guard

**Files:**
- Create: `scripts/flutter_evidence_ui_gate_policy_test.sh`
- Create: `docs/ci/flutter-evidence-ui-gate-policy.md`
- Modify: `scripts/flutter_evidence_ui_gate.sh`
- Modify: `.github/workflows/ci.yml`

- [x] **Step 1: Write failing policy guard**

Create `scripts/flutter_evidence_ui_gate_policy_test.sh`.

The guard must require:
- `docs/ci/flutter-evidence-ui-gate-policy.md`
- `gate_scope: evidence_ui_contract`
- `required_on_pull_request: true`
- `duplicates_full_flutter_test: true`
- `max_test_files: 13`
- `gate_script: scripts/flutter_evidence_ui_gate.sh`
- `policy_test: scripts/flutter_evidence_ui_gate_policy_test.sh`
- `Flutter evidence UI gate policy` in CI workflow
- `scripts/flutter_evidence_ui_gate_policy_test.sh` in CI workflow

The guard must run `bash scripts/flutter_evidence_ui_gate.sh --count` and fail if count is not numeric or exceeds `max_test_files`.

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: FAIL because the policy document and CI step do not exist yet.

- [x] **Step 3: Implement policy and count mode**

Add `--count` to `scripts/flutter_evidence_ui_gate.sh`.

Create `docs/ci/flutter-evidence-ui-gate-policy.md` with policy markers and rationale.

Add this step to the `ssot-governance` job:

```yaml
      - name: Flutter evidence UI gate policy
        run: bash scripts/flutter_evidence_ui_gate_policy_test.sh
```

- [x] **Step 4: Run policy guard to verify GREEN**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/flutter-evidence-ui-gate-policy-p21.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_test.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`

Run:
`bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- .github/workflows/ci.yml scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh docs/ci/flutter-evidence-ui-gate-policy.md docs/superpowers/plans/2026-05-13-flutter-evidence-ui-gate-policy-p21.md docs/audit/flutter-evidence-ui-gate-policy-p21.md CHANGELOG.md CONTEXT.md`
