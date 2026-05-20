# CI Web Gate Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web release gate를 PR 필수 CI gate로 유지한다는 정책을 문서화하고 guard test로 고정한다.

**Architecture:** 정책 문서는 `docs/ci/flutter-web-release-gate-policy.md`에 둔다. guard script는 정책 문서의 machine-readable marker와 CI workflow 연결을 동시에 확인한다.

**Tech Stack:** Markdown policy document, Bash guard script, GitHub Actions workflow.

---

### Task 1: Policy guard

**Files:**
- Create: `scripts/ci_web_gate_policy_test.sh`
- Create: `docs/ci/flutter-web-release-gate-policy.md`

- [ ] **Step 1: Write the failing policy guard**

Create a shell test that requires:
- `docs/ci/flutter-web-release-gate-policy.md` exists.
- It contains `required_on_pull_request: true`.
- It contains `wasm_dry_run_blocking: false`.
- CI workflow still contains `Flutter web release gate`.

- [ ] **Step 2: Run test to verify RED**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Expected: FAIL because the policy document does not exist yet.

- [ ] **Step 3: Add policy document**

Create `docs/ci/flutter-web-release-gate-policy.md` with:

```yaml
required_on_pull_request: true
wasm_dry_run_blocking: false
release_target: js_web
```

Explain that JS web build failures block PR/release, while Wasm dry-run warnings are tracked but non-blocking until Wasm becomes a product target.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/ci-web-gate-policy-p18.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run script tests and syntax checks**

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Run:
`bash -n scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`

- [ ] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- docs/ci/flutter-web-release-gate-policy.md scripts/ci_web_gate_policy_test.sh CHANGELOG.md CONTEXT.md`
