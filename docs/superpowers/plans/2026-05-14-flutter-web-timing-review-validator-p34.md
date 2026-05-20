# Flutter Web Timing Review Validator P34 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Concrete `docs/ci/fixtures/flutter-web-timing/<review-id>` directory가 threshold change PR 전에 fixture 규약을 만족하는지 검증하는 validator를 추가한다.

**Architecture:** Validator는 `manifest.env`, `samples/sample-01.env` through `sample-05.env`, `aggregate.env`를 검사한다. Manifest의 artifact source/name/status/no-PHI attestation을 검증하고, forbidden sensitive fields를 sample/aggregate 파일에서 차단하며, P32 aggregator를 재실행해 `aggregate.env`가 sample files와 일치하는지 확인한다.

**Tech Stack:** Bash, key-value fixture files, existing Flutter web timing sample aggregator.

---

## File Structure

- Create: `scripts/flutter_web_timing_review_fixture_validate_test.sh`
  - 임시 review directory를 만들고 validator 성공/실패 경로를 검증한다.
- Create: `scripts/flutter_web_timing_review_fixture_validate.sh`
  - concrete review directory의 manifest, samples, aggregate, no-PHI attestation, aggregate reproducibility를 검증한다.
- Modify: `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - fixture policy가 validator script/test와 manifest/forbidden field 계약을 문서화했는지 검증한다.
- Modify: `docs/ci/flutter-web-timing-review-fixtures.md`
  - validator 사용법, manifest required fields, forbidden fixture fields를 문서화한다.
- Create: `docs/audit/flutter-web-timing-review-validator-p34.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P34 진행 상황과 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Validator RED Test

- [ ] **Step 1: Write failing test**

Create `scripts/flutter_web_timing_review_fixture_validate_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR="$ROOT/scripts/flutter_web_timing_review_fixture_validate.sh"
AGGREGATOR="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

fail() {
  echo "FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_TEST_FAIL $1" >&2
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
bash "$VALIDATOR" --review-dir "$TMP_DIR/review" | grep -Fq "FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_PASS" || fail "expected validator pass"

cp -R "$TMP_DIR/review" "$TMP_DIR/bad-review"
sed -i 's/no_phi_attestation=true/no_phi_attestation=false/' "$TMP_DIR/bad-review/manifest.env"
if bash "$VALIDATOR" --review-dir "$TMP_DIR/bad-review" >/tmp/manpasik_fixture_validate_bad.out 2>&1; then
  fail "expected no-PHI failure"
fi
grep -Fq "no_phi_attestation must be true" /tmp/manpasik_fixture_validate_bad.out || fail "wrong no-PHI failure"

echo "FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_TEST_PASS"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`

Expected: failure because `scripts/flutter_web_timing_review_fixture_validate.sh` is missing.

## Task 2: Validator Implementation

- [ ] **Step 1: Implement validator**

Create `scripts/flutter_web_timing_review_fixture_validate.sh` supporting:

```bash
bash scripts/flutter_web_timing_review_fixture_validate.sh \
  --review-dir docs/ci/fixtures/flutter-web-timing/<review-id>
```

Validator must check:

- `manifest.env` exists.
- `artifact_source=github_actions`.
- `artifact_name=flutter-web-timing-report`.
- `sample_count=5`.
- `sample_files=sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env`.
- `aggregate_output=aggregate.env`.
- `review_status` is `proposed`, `approved`, or `rejected`.
- `no_phi_attestation=true`.
- All sample files exist.
- Files do not contain forbidden keys: `user_id`, `patient_id`, `device_id`, `access_token`, `refresh_token`, `raw_channels`, `s_det`, `s_ref`, `primary_value`.
- Re-running `scripts/flutter_web_timing_sample_aggregate.sh` produces the same `aggregate.env`.

- [ ] **Step 2: Run test to verify it passes**

Run: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`

Expected: `FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_TEST_PASS`

## Task 3: Fixture Policy Integration

- [ ] **Step 1: Extend fixture policy guard**

Add markers to `scripts/ci_web_timing_review_fixture_policy_test.sh`:

```bash
"fixture_validator_script: scripts/flutter_web_timing_review_fixture_validate.sh"
"fixture_validator_test: scripts/flutter_web_timing_review_fixture_validate_test.sh"
"manifest_required_fields: review_id,artifact_source,artifact_name,sample_count,sample_files,aggregate_output,review_status,no_phi_attestation"
"forbidden_fixture_fields: user_id,patient_id,device_id,access_token,refresh_token,raw_channels,s_det,s_ref,primary_value"
```

Also check both validator files exist.

- [ ] **Step 2: Run fixture policy guard to verify it fails**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: failure for missing `fixture_validator_script`.

- [ ] **Step 3: Update fixture policy doc**

Add the same markers and a `## Validator` section to `docs/ci/flutter-web-timing-review-fixtures.md`.

- [ ] **Step 4: Run fixture policy guard to verify it passes**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: `CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS`

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/flutter_web_timing_review_fixture_validate_test.sh
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/flutter_web_timing_review_fixture_validate.sh scripts/flutter_web_timing_review_fixture_validate_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-review-validator-p34.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P34 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
