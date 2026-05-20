# Flutter Web Gate Cost Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web release gate의 PR 필수 정책과 nightly split 미적용 결정을 문서/guard로 고정한다.

**Architecture:** `docs/ci/flutter-web-release-gate-policy.md`는 JS web release target과 CI 비용 정책을 함께 기록한다. `scripts/ci_web_gate_policy_test.sh`는 PR-required, release-branch-required, nightly split disabled 상태가 문서와 workflow에서 유지되는지 검증한다.

**Tech Stack:** Bash, Markdown CI policy, GitHub Actions workflow.

---

### Task 1: Web gate cost policy markers

**Files:**
- Modify: `scripts/ci_web_gate_policy_test.sh`
- Modify: `docs/ci/flutter-web-release-gate-policy.md`

- [x] **Step 1: Write failing policy expectations**

Update `scripts/ci_web_gate_policy_test.sh` to require:
- `ci_execution_mode: pull_request_and_release_branch`
- `release_branch_required: true`
- `nightly_split: false`
- `cost_review_required_before_relaxing: true`
- `blocking_build_command: flutter build web --no-pub`

- [x] **Step 2: Run policy guard to verify RED**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Expected: FAIL because the new cost policy markers are not present.

- [x] **Step 3: Update policy document**

Update `docs/ci/flutter-web-release-gate-policy.md` with the markers above and a short rationale:
- JS web is the current release artifact.
- PR blocking remains enabled until a measured CI cost review justifies moving the build to nightly/release-only.
- Wasm warning remains non-blocking.

- [x] **Step 4: Run policy guard to verify GREEN**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/flutter-web-gate-cost-policy-p25.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Run:
`bash scripts/ci_wasm_readiness_policy_test.sh`

Run:
`bash -n scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh scripts/ci_wasm_readiness_policy_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- scripts/ci_web_gate_policy_test.sh docs/ci/flutter-web-release-gate-policy.md docs/superpowers/plans/2026-05-13-flutter-web-gate-cost-policy-p25.md docs/audit/flutter-web-gate-cost-policy-p25.md CHANGELOG.md CONTEXT.md`
