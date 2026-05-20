# Flutter Web Timing Threshold Change Review Check P35 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Threshold 변경 PR에서 reviewer가 concrete timing fixture validator와 audit template 값을 빠짐없이 실행/작성했는지 확인하는 wrapper command를 추가한다.

**Architecture:** Wrapper는 P34 validator를 먼저 실행한 뒤 `manifest.env`와 `aggregate.env`의 값을 읽는다. 지정된 audit Markdown 파일에 review id, decision, sample count, median/p95/worst-case, branch types, runner contexts, validator command가 aggregate와 일치하게 기록됐는지 검사한다.

**Tech Stack:** Bash, Markdown audit file, key-value fixture files.

---

## File Structure

- Create: `scripts/flutter_web_timing_threshold_change_review_check_test.sh`
  - 임시 review fixture와 audit file을 만들고 wrapper 성공 및 audit mismatch 실패를 검증한다.
- Create: `scripts/flutter_web_timing_threshold_change_review_check.sh`
  - `--review-dir`와 `--audit-file`을 받아 fixture validator와 audit field consistency를 검증한다.
- Modify: `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - fixture policy와 audit template이 wrapper script/test 및 required audit fields를 문서화했는지 검증한다.
- Modify: `docs/ci/flutter-web-timing-review-fixtures.md`
  - threshold change review check command와 required audit fields를 문서화한다.
- Modify: `docs/audit/templates/flutter-web-timing-review-template.md`
  - wrapper가 읽을 수 있는 machine-check field block을 추가한다.
- Create: `docs/audit/flutter-web-timing-threshold-change-review-check-p35.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P35 진행 상황과 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Wrapper RED Test

- [ ] **Step 1: Write failing test**

Create `scripts/flutter_web_timing_threshold_change_review_check_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECKER="$ROOT/scripts/flutter_web_timing_threshold_change_review_check.sh"
AGGREGATOR="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

fail() {
  echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_FAIL $1" >&2
  exit 1
}

write_sample() {
  local file="$1"
  local duration="$2"
  local branch_type="$3"
  cat >"$TMP_DIR/review/samples/$file" <<EOF_SAMPLE
report_version=1
source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
sample_count=1
duration_seconds_values=$duration
latest_duration_seconds=$duration
min_duration_seconds=$duration
max_duration_seconds=$duration
branch_type=$branch_type
runner_context=Linux-X64
EOF_SAMPLE
}

mkdir -p "$TMP_DIR/review/samples"
write_sample sample-01.env 12 pull_request
write_sample sample-02.env 14 pull_request
write_sample sample-03.env 15 pull_request
write_sample sample-04.env 20 release_branch
write_sample sample-05.env 30 release_branch

cat >"$TMP_DIR/review/manifest.env" <<'EOF_MANIFEST'
review_id=2026-05-14-pr-gate-cost-review
artifact_source=github_actions
artifact_name=flutter-web-timing-report
sample_count=5
sample_files=sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env
aggregate_output=aggregate.env
review_status=proposed
no_phi_attestation=true
EOF_MANIFEST

bash "$AGGREGATOR" --input-dir "$TMP_DIR/review/samples" --output "$TMP_DIR/review/aggregate.env" >/dev/null

cat >"$TMP_DIR/review-audit.md" <<EOF_AUDIT
# Flutter Web Timing Review

review_template_version: 1
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
aggregate_source: aggregate.env
review_id: 2026-05-14-pr-gate-cost-review
decision: proposed
sample_count: 5
median_seconds: 15
p95_seconds: 30
worst_case_seconds: 30
branch_types: pull_request,release_branch
runner_contexts: Linux-X64
validator_command: bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir $TMP_DIR/review
EOF_AUDIT

bash "$CHECKER" --review-dir "$TMP_DIR/review" --audit-file "$TMP_DIR/review-audit.md" | grep -Fq "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_PASS" || fail "expected wrapper pass"

cp "$TMP_DIR/review-audit.md" "$TMP_DIR/bad-audit.md"
sed -i 's/p95_seconds: 30/p95_seconds: 29/' "$TMP_DIR/bad-audit.md"
if bash "$CHECKER" --review-dir "$TMP_DIR/review" --audit-file "$TMP_DIR/bad-audit.md" >/tmp/manpasik_threshold_review_bad.out 2>&1; then
  fail "expected p95 mismatch failure"
fi
grep -Fq "audit field mismatch: p95_seconds" /tmp/manpasik_threshold_review_bad.out || fail "wrong p95 mismatch failure"

echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_PASS"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`

Expected: failure because `scripts/flutter_web_timing_threshold_change_review_check.sh` is missing.

## Task 2: Wrapper Implementation

- [ ] **Step 1: Implement wrapper**

Create `scripts/flutter_web_timing_threshold_change_review_check.sh` supporting:

```bash
bash scripts/flutter_web_timing_threshold_change_review_check.sh \
  --review-dir docs/ci/fixtures/flutter-web-timing/<review-id> \
  --audit-file docs/audit/<review-doc>.md
```

Wrapper must:

- Run `scripts/flutter_web_timing_review_fixture_validate.sh --review-dir <dir>`.
- Read `review_id` from `manifest.env`.
- Read `sample_count`, `median_seconds`, `p95_seconds`, `worst_case_seconds`, `branch_types`, `runner_contexts` from `aggregate.env`.
- Require audit file fields:
  - `review_template_version: 1`
  - `review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md`
  - `aggregate_source: aggregate.env`
  - `review_id: <manifest review_id>`
  - `decision: proposed|approved|rejected`
  - aggregate values exactly matching `aggregate.env`
  - `validator_command: bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir <review-dir>`

- [ ] **Step 2: Run test to verify it passes**

Run: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`

Expected: `FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_TEST_PASS`

## Task 3: Policy and Template Integration

- [ ] **Step 1: Extend fixture policy guard**

Add to `scripts/ci_web_timing_review_fixture_policy_test.sh` required policy markers:

```bash
"threshold_change_review_check_script: scripts/flutter_web_timing_threshold_change_review_check.sh"
"threshold_change_review_check_test: scripts/flutter_web_timing_threshold_change_review_check_test.sh"
"threshold_change_review_check_required: true"
"threshold_change_review_required_fields: review_id,decision,sample_count,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts,validator_command"
```

Also require both wrapper files to exist and template to contain `validator_command:`.

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: failure for missing `threshold_change_review_check_script`.

- [ ] **Step 3: Update policy and audit template**

Add the same markers and command usage to `docs/ci/flutter-web-timing-review-fixtures.md`.

Add a machine-check field block to `docs/audit/templates/flutter-web-timing-review-template.md`:

```text
review_id:
decision:
sample_count:
median_seconds:
p95_seconds:
worst_case_seconds:
branch_types:
runner_contexts:
validator_command:
```

- [ ] **Step 4: Run guard to verify it passes**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: `CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS`

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/flutter_web_timing_threshold_change_review_check_test.sh
bash scripts/flutter_web_timing_review_fixture_validate_test.sh
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/flutter_web_timing_threshold_change_review_check.sh scripts/flutter_web_timing_threshold_change_review_check_test.sh scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_review_fixture_validate_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-threshold-change-review-check-p35.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P35 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
