# Wasm Readiness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter Wasm을 릴리스 타깃으로 승격하기 전 필요한 blocker와 정책을 명확히 문서화한다.

**Architecture:** 현재 `release_target`은 P18 정책에 따라 `js_web`이다. P19는 Wasm을 즉시 지원하지 않고, dry-run warning에서 반복 확인된 dependency blocker를 readiness 문서와 guard script로 추적한다.

**Tech Stack:** Markdown readiness document, Bash guard script, Flutter web build warning evidence.

---

### Task 1: Wasm readiness guard

**Files:**
- Create: `scripts/ci_wasm_readiness_policy_test.sh`
- Create: `docs/ci/flutter-wasm-readiness.md`

- [ ] **Step 1: Write the failing guard**

Create a shell test that requires `docs/ci/flutter-wasm-readiness.md` to contain:
- `wasm_release_target: false`
- `flutter_secure_storage_web`
- `share_plus`
- `connectivity_plus`
- `package:js`

- [ ] **Step 2: Run test to verify RED**

Run:
`bash scripts/ci_wasm_readiness_policy_test.sh`

Expected: FAIL because readiness document does not exist yet.

- [ ] **Step 3: Add readiness document**

Create a document describing:
- Current Wasm release target status: false
- Current blocker packages from dry-run output
- Promotion criteria before `wasm_dry_run_blocking` can become true

- [ ] **Step 4: Run test to verify GREEN**

Run:
`bash scripts/ci_wasm_readiness_policy_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/wasm-readiness-p19.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run script tests and syntax checks**

Run:
`bash scripts/ci_wasm_readiness_policy_test.sh`

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash -n scripts/ci_wasm_readiness_policy_test.sh scripts/ci_web_gate_policy_test.sh`

- [ ] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- docs/ci/flutter-wasm-readiness.md scripts/ci_wasm_readiness_policy_test.sh CHANGELOG.md CONTEXT.md`
