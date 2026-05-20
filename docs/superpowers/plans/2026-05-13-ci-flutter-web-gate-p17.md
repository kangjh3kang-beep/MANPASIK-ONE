# CI Flutter Web Gate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P16의 Flutter web release gate를 GitHub Actions CI의 Flutter job에 연결한다.

**Architecture:** 기존 `flutter-app` job의 analyze/test 이후에 `bash ../../scripts/flutter_web_release_gate.sh`를 실행한다. workflow guard script로 CI 파일이 release gate step과 script path를 포함하는지 검증한다.

**Tech Stack:** GitHub Actions YAML, Bash workflow guard, Flutter web release gate.

---

### Task 1: CI workflow guard

**Files:**
- Create: `scripts/ci_flutter_web_gate_workflow_test.sh`
- Modify: `.github/workflows/ci.yml`

- [ ] **Step 1: Write the failing guard test**

Create a shell test that verifies `.github/workflows/ci.yml` contains:
- `Flutter web release gate`
- `scripts/flutter_web_release_gate.sh`

- [ ] **Step 2: Run test to verify RED**

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Expected: FAIL because CI has not connected the new gate yet.

- [ ] **Step 3: Add CI step**

In the `flutter-app` job, after `Run tests`, add:

```yaml
      - name: Flutter web release gate
        run: bash ../../scripts/flutter_web_release_gate.sh
```

- [ ] **Step 4: Run guard test to verify GREEN**

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/ci-flutter-web-gate-p17.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run script tests and syntax checks**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Run:
`bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`

- [ ] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- .github/workflows/ci.yml scripts/ci_flutter_web_gate_workflow_test.sh CHANGELOG.md CONTEXT.md`
