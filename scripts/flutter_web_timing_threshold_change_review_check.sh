#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR="$ROOT/scripts/flutter_web_timing_review_fixture_validate.sh"

usage() {
  cat >&2 <<'EOF_USAGE'
usage: flutter_web_timing_threshold_change_review_check.sh --review-dir <dir> --audit-file <file>
EOF_USAGE
}

fail() {
  echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_FAIL $1" >&2
  exit 1
}

review_dir=""
audit_file=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --review-dir)
      review_dir="${2:-}"
      shift 2
      ;;
    --audit-file)
      audit_file="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown argument: $1"
      ;;
  esac
done

if [[ -z "$review_dir" || -z "$audit_file" ]]; then
  usage
  fail "missing required arguments"
fi

if [[ ! -f "$audit_file" ]]; then
  fail "audit file not found: $audit_file"
fi

bash "$VALIDATOR" --review-dir "$review_dir" >/dev/null

manifest="$review_dir/manifest.env"
aggregate="$review_dir/aggregate.env"

read_field() {
  local file="$1"
  local key="$2"
  grep -E "^${key}=" "$file" | head -n 1 | cut -d= -f2-
}

require_audit_line() {
  local label="$1"
  local expected="$2"
  if ! grep -Fq "$label: $expected" "$audit_file"; then
    fail "audit field mismatch: $label"
  fi
}

require_audit_line "review_template_version" "1"
require_audit_line "review_fixture_policy" "docs/ci/flutter-web-timing-review-fixtures.md"
require_audit_line "aggregate_source" "aggregate.env"

review_id="$(read_field "$manifest" "review_id")"
sample_count="$(read_field "$aggregate" "sample_count")"
median_seconds="$(read_field "$aggregate" "median_seconds")"
p95_seconds="$(read_field "$aggregate" "p95_seconds")"
worst_case_seconds="$(read_field "$aggregate" "worst_case_seconds")"
branch_types="$(read_field "$aggregate" "branch_types")"
runner_contexts="$(read_field "$aggregate" "runner_contexts")"

require_audit_line "review_id" "$review_id"
require_audit_line "sample_count" "$sample_count"
require_audit_line "median_seconds" "$median_seconds"
require_audit_line "p95_seconds" "$p95_seconds"
require_audit_line "worst_case_seconds" "$worst_case_seconds"
require_audit_line "branch_types" "$branch_types"
require_audit_line "runner_contexts" "$runner_contexts"
require_audit_line "validator_command" "bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir $review_dir"

decision="$(grep -E '^decision: ' "$audit_file" | head -n 1 | cut -d: -f2- | sed 's/^ //')"
case "$decision" in
  proposed|approved|rejected)
    ;;
  *)
    fail "audit field mismatch: decision"
    ;;
esac

echo "FLUTTER_WEB_TIMING_THRESHOLD_CHANGE_REVIEW_CHECK_PASS $audit_file"
