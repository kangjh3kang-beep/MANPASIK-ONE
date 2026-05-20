# Flutter Web Timing Synthetic Review Example P37 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P36 release-review checklist가 실제 fixture/audit 구조에서 재현되는지 검증하는 비운영 synthetic timing review example을 추가한다.

**Architecture:** Synthetic review directory는 real fixture convention과 동일한 `manifest.env`, five sample artifacts, `aggregate.env` 구조를 사용한다. A matching synthetic audit document is checked by the existing threshold review wrapper, and the checklist policy guard runs a synthetic example test so CI validates the documented review path.

**Tech Stack:** Bash, Markdown, key-value timing fixture files.

---

## File Structure

- Create: `scripts/ci_web_timing_synthetic_review_example_test.sh`
  - Synthetic fixture와 audit example이 존재하고 validator/wrapper를 통과하는지 검증한다.
- Create: `docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review/manifest.env`
  - Synthetic review manifest with `synthetic_example=true`.
- Create: `docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review/samples/sample-01.env` through `sample-05.env`
  - Synthetic `flutter-web-timing-report`-format samples.
- Create: `docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review/aggregate.env`
  - Expected aggregate generated from the synthetic samples.
- Create: `docs/audit/flutter-web-timing-synthetic-review-example.md`
  - Audit example that passes `scripts/flutter_web_timing_threshold_change_review_check.sh`.
- Modify: `scripts/ci_web_release_review_checklist_policy_test.sh`
  - Checklist policy guard가 synthetic example test script와 fixture/audit references를 요구하고 synthetic example test를 실행하도록 확장한다.
- Modify: `docs/ci/flutter-web-release-review-checklist.md`
  - Synthetic dry-run command와 example paths를 문서화한다.
- Create: `docs/audit/flutter-web-timing-synthetic-review-example-p37.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P37 상태와 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Synthetic Example RED Guard

- [ ] **Step 1: Write failing synthetic example guard**

Create `scripts/ci_web_timing_synthetic_review_example_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REVIEW_DIR="$ROOT/docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review"
AUDIT_FILE="$ROOT/docs/audit/flutter-web-timing-synthetic-review-example.md"

fail() {
  echo "CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_FAIL $1" >&2
  exit 1
}

[[ -d "$REVIEW_DIR" ]] || fail "missing synthetic review dir"
[[ -f "$AUDIT_FILE" ]] || fail "missing synthetic audit file"
grep -Fq "synthetic_example=true" "$REVIEW_DIR/manifest.env" || fail "manifest missing synthetic marker"
grep -Fq "synthetic_review: true" "$AUDIT_FILE" || fail "audit missing synthetic marker"

bash "$ROOT/scripts/flutter_web_timing_review_fixture_validate.sh" --review-dir "$REVIEW_DIR" >/dev/null
bash "$ROOT/scripts/flutter_web_timing_threshold_change_review_check.sh" \
  --review-dir "$REVIEW_DIR" \
  --audit-file "$AUDIT_FILE" >/dev/null

echo "CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_PASS"
```

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`

Expected: failure because synthetic review directory is missing.

## Task 2: Synthetic Fixture and Audit Example

- [ ] **Step 1: Add fixture files**

Create review directory `docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review` with:

```text
manifest.env
samples/sample-01.env
samples/sample-02.env
samples/sample-03.env
samples/sample-04.env
samples/sample-05.env
aggregate.env
```

Use synthetic durations `11,13,17,21,34`, where median is `17`, p95 is `34`, and worst-case is `34`.

- [ ] **Step 2: Add audit example**

Create `docs/audit/flutter-web-timing-synthetic-review-example.md` with machine-check fields matching `aggregate.env`:

```text
review_template_version: 1
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
aggregate_source: aggregate.env
synthetic_review: true
review_id: 2026-05-19-synthetic-review
decision: rejected
sample_count: 5
median_seconds: 17
p95_seconds: 34
worst_case_seconds: 34
branch_types: pull_request,release_branch
runner_contexts: Linux-X64
validator_command: bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review
```

- [ ] **Step 3: Run guard to verify it passes**

Run: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`

Expected: `CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_PASS`

## Task 3: Checklist Policy Integration

- [ ] **Step 1: Extend checklist policy guard**

Modify `scripts/ci_web_release_review_checklist_policy_test.sh` to require:

```bash
"synthetic_review_example: docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review"
"synthetic_review_audit: docs/audit/flutter-web-timing-synthetic-review-example.md"
"bash scripts/ci_web_timing_synthetic_review_example_test.sh"
```

Also require the synthetic script and audit file to exist, and run `bash scripts/ci_web_timing_synthetic_review_example_test.sh`.

- [ ] **Step 2: Run checklist policy guard to verify it fails**

Run: `bash scripts/ci_web_release_review_checklist_policy_test.sh`

Expected: failure for missing synthetic checklist marker.

- [ ] **Step 3: Update checklist doc**

Add a `Synthetic Dry Run` section to `docs/ci/flutter-web-release-review-checklist.md` with the exact synthetic paths and command:

```bash
bash scripts/ci_web_timing_synthetic_review_example_test.sh
```

- [ ] **Step 4: Run checklist policy guard to verify it passes**

Run: `bash scripts/ci_web_release_review_checklist_policy_test.sh`

Expected: `CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_PASS`

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_web_timing_synthetic_review_example_test.sh
bash scripts/ci_web_release_review_checklist_policy_test.sh
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/ci_web_timing_synthetic_review_example_test.sh scripts/ci_web_release_review_checklist_policy_test.sh scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_threshold_change_review_check.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-synthetic-review-example-p37.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P37 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
