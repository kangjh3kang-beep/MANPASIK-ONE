# Flutter Evidence CI Shard Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P23에서 CI를 category shard step으로 분리한 운영 상태를 policy marker와 guard로 고정한다.

**Architecture:** `docs/ci/flutter-evidence-ui-gate-policy.md`는 evidence gate 범위와 CI 실행 모드를 함께 기록한다. `scripts/flutter_evidence_ui_gate_policy_test.sh`는 policy marker와 workflow 상태가 일치하는지 확인한다.

**Tech Stack:** Bash, Markdown CI policy, GitHub Actions workflow.

---

### Task 1: CI shard policy markers

**Files:**
- Modify: `scripts/flutter_evidence_ui_gate_policy_test.sh`
- Modify: `docs/ci/flutter-evidence-ui-gate-policy.md`

- [x] **Step 1: Write failing policy expectations**

Update `scripts/flutter_evidence_ui_gate_policy_test.sh` to require:
- `ci_execution_mode: category_shards`
- `ci_shard_steps: contract,ui,data`
- `aggregate_ci_step: false`
- `aggregate_local_gate_supported: true`

Also verify:
- each `ci_shard_steps` category has a matching CI workflow `--category <name>` command.
- when `aggregate_ci_step: false`, the workflow must not contain `run: bash ../../scripts/flutter_evidence_ui_gate.sh`.

- [x] **Step 2: Run policy guard to verify RED**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: FAIL because policy markers are not present yet.

- [x] **Step 3: Update policy document**

Update `docs/ci/flutter-evidence-ui-gate-policy.md` with:

```yaml
ci_execution_mode: category_shards
ci_shard_steps: contract,ui,data
aggregate_ci_step: false
aggregate_local_gate_supported: true
```

Add a short note that aggregate local execution remains supported for developer verification, but CI uses category shards.

- [x] **Step 4: Run policy guard to verify GREEN**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/flutter-evidence-ci-shard-policy-p24.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/flutter_evidence_ui_gate_policy_test.sh`

Run:
`bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`

Run:
`bash scripts/flutter_evidence_ui_gate_shard_test.sh`

Run:
`bash -n scripts/flutter_evidence_ui_gate_policy_test.sh scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/flutter_evidence_ui_gate_shard_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- scripts/flutter_evidence_ui_gate_policy_test.sh docs/ci/flutter-evidence-ui-gate-policy.md docs/superpowers/plans/2026-05-13-flutter-evidence-ci-shard-policy-p24.md docs/audit/flutter-evidence-ci-shard-policy-p24.md CHANGELOG.md CONTEXT.md`
