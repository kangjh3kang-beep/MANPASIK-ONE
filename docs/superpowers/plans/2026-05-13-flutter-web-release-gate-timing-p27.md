# Flutter Web Release Gate Timing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web release gate가 build elapsed seconds marker를 출력하게 해, P25의 CI cost review를 수치 기반으로 진행할 수 있게 한다.

**Architecture:** `scripts/flutter_web_release_gate.sh`는 실제 build 경로에서 시작/종료 epoch seconds를 측정하고 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 출력한다. `--policy-check` 경로는 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` env가 주어졌을 때 marker를 출력해 빠른 parser test가 timing 출력 형식을 검증할 수 있게 한다.

**Tech Stack:** Bash, Flutter web build gate, Markdown CI policy.

---

### Task 1: Timing marker

**Files:**
- Modify: `scripts/flutter_web_release_gate_test.sh`
- Modify: `scripts/flutter_web_release_gate.sh`

- [x] **Step 1: Write failing timing test**

Update `scripts/flutter_web_release_gate_test.sh` so the successful policy-check run sets:

```bash
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=42
```

and requires output to contain:

```text
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=42
```

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Expected: FAIL because the gate does not emit timing marker yet.

- [x] **Step 3: Implement timing marker**

Update `scripts/flutter_web_release_gate.sh` so:
- `policy_check` prints `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` when env var `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` is set to a numeric value.
- actual build path captures `start_time="$(date +%s)"`, `end_time="$(date +%s)"`, computes elapsed seconds, and passes it to `policy_check`.
- invalid timing env values fail with `FLUTTER_WEB_RELEASE_GATE_FAIL invalid duration seconds`.

- [x] **Step 4: Run timing test to verify GREEN**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Expected: PASS.

### Task 2: Policy and matrix markers

**Files:**
- Modify: `scripts/ci_web_gate_policy_test.sh`
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
- Modify: `docs/ci/flutter-web-release-gate-policy.md`
- Modify: `docs/ci/ci-gate-matrix.md`

- [x] **Step 1: Write failing policy expectations**

Update guards to require:
- `timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS`
- `timing_capture: true`

- [x] **Step 2: Run policy guards to verify RED**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: FAIL until docs are updated.

- [x] **Step 3: Update policy and matrix docs**

Add timing markers to:
- `docs/ci/flutter-web-release-gate-policy.md`
- `docs/ci/ci-gate-matrix.md`

- [x] **Step 4: Run policy guards to verify GREEN**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: PASS.

### Task 3: Quality and logs

**Files:**
- Create: `docs/audit/flutter-web-release-gate-timing-p27.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Run:
`bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_gate_matrix_policy_test.sh docs/ci/flutter-web-release-gate-policy.md docs/ci/ci-gate-matrix.md docs/superpowers/plans/2026-05-13-flutter-web-release-gate-timing-p27.md docs/audit/flutter-web-release-gate-timing-p27.md CHANGELOG.md CONTEXT.md`
