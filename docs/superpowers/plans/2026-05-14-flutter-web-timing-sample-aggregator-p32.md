# Flutter Web Timing Sample Aggregator P32 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Exported `flutter-web-timing-report` artifact files 5개 이상에서 median/p95/worst-case timing 값을 산출하는 offline aggregator를 추가한다.

**Architecture:** P30 report 파일은 `key=value` 형식이므로 Bash parser로 `latest_duration_seconds`, `branch_type`, `runner_context`를 읽는다. Aggregator는 최소 sample 수와 branch type 구성을 검증한 뒤 정렬된 duration 배열로 median, nearest-rank p95, worst-case 값을 계산해 review용 key-value report를 출력한다.

**Tech Stack:** Bash, key-value artifact files, CI policy docs.

---

## File Structure

- Create: `scripts/flutter_web_timing_sample_aggregate_test.sh`
  - fixture artifact 5개를 만들고 aggregate output의 sample count, branch types, median, p95, worst-case를 검증한다.
- Create: `scripts/flutter_web_timing_sample_aggregate.sh`
  - exported `.env` artifact directory를 읽어 aggregate report를 생성한다.
- Modify: `scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - promotion policy가 aggregator script/test와 output schema를 문서화했는지 검증한다.
- Modify: `docs/ci/flutter-web-timing-threshold-promotion.md`
  - offline aggregator 사용법과 output fields를 기록한다.
- Create: `docs/audit/flutter-web-timing-sample-aggregator-p32.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P32 진행 상황과 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Aggregator RED Test

- [ ] **Step 1: Write failing test**

Create `scripts/flutter_web_timing_sample_aggregate_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

fail() {
  echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_FAIL $1" >&2
  exit 1
}

write_sample() {
  local file="$1"
  local duration="$2"
  local branch_type="$3"
  cat >"$TMP_DIR/$file" <<EOF_SAMPLE
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

write_sample pr-1.env 12 pull_request
write_sample pr-2.env 14 pull_request
write_sample pr-3.env 15 pull_request
write_sample release-1.env 20 release_branch
write_sample release-2.env 30 release_branch

OUTPUT="$TMP_DIR/aggregate.env"
bash "$SCRIPT" --input-dir "$TMP_DIR" --output "$OUTPUT"

grep -Fq "aggregate_version=1" "$OUTPUT" || fail "missing aggregate_version"
grep -Fq "sample_count=5" "$OUTPUT" || fail "missing sample_count"
grep -Fq "branch_types=pull_request,release_branch" "$OUTPUT" || fail "missing branch types"
grep -Fq "median_seconds=15" "$OUTPUT" || fail "missing median"
grep -Fq "p95_seconds=30" "$OUTPUT" || fail "missing p95"
grep -Fq "worst_case_seconds=30" "$OUTPUT" || fail "missing worst case"

if bash "$SCRIPT" --input-dir "$TMP_DIR" --min-samples 6 --output "$TMP_DIR/too-few.env" >/tmp/manpasik_aggregate_too_few.out 2>&1; then
  fail "expected min sample failure"
fi
grep -Fq "minimum samples not met" /tmp/manpasik_aggregate_too_few.out || fail "wrong min sample failure"

echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_PASS"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`

Expected: failure because `scripts/flutter_web_timing_sample_aggregate.sh` is missing.

## Task 2: Aggregator Implementation

- [ ] **Step 1: Implement parser and statistics**

Create `scripts/flutter_web_timing_sample_aggregate.sh` that supports:

```bash
bash scripts/flutter_web_timing_sample_aggregate.sh \
  --input-dir exported-artifacts \
  --output flutter-web-timing-aggregate.env
```

Required output:

```text
aggregate_version=1
source_artifact=flutter-web-timing-report
sample_count=<n>
duration_seconds_values=<sorted comma-separated durations>
median_seconds=<median>
p95_seconds=<nearest-rank-p95>
worst_case_seconds=<max>
branch_types=<comma-separated unique branch types>
runner_contexts=<comma-separated unique runner contexts>
```

- [ ] **Step 2: Run test to verify it passes**

Run: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`

Expected: `FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_PASS`

## Task 3: Promotion Policy Integration

- [ ] **Step 1: Extend promotion policy guard**

Add these markers to `scripts/ci_web_timing_threshold_promotion_policy_test.sh`:

```bash
"sample_aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh"
"sample_aggregator_test: scripts/flutter_web_timing_sample_aggregate_test.sh"
"aggregation_output_fields: aggregate_version,source_artifact,sample_count,duration_seconds_values,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts"
```

- [ ] **Step 2: Run promotion policy guard to verify it fails**

Run: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`

Expected: failure for missing `sample_aggregator_script`.

- [ ] **Step 3: Update promotion policy doc**

Add the same markers and an `Offline Aggregator` section to `docs/ci/flutter-web-timing-threshold-promotion.md`.

- [ ] **Step 4: Run promotion policy guard to verify it passes**

Run: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`

Expected: `CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_PASS`

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/flutter_web_timing_sample_aggregate_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_web_gate_timing_collection_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/flutter_web_timing_sample_aggregate.sh scripts/flutter_web_timing_sample_aggregate_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-sample-aggregator-p32.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P32 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
