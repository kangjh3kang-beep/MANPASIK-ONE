#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AGGREGATOR="$ROOT/scripts/flutter_web_timing_sample_aggregate.sh"

usage() {
  cat >&2 <<'EOF_USAGE'
usage: flutter_web_timing_review_fixture_validate.sh --review-dir <dir>
EOF_USAGE
}

fail() {
  echo "FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_FAIL $1" >&2
  exit 1
}

review_dir=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --review-dir)
      review_dir="${2:-}"
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

if [[ -z "$review_dir" ]]; then
  usage
  fail "missing review dir"
fi

if [[ ! -d "$review_dir" ]]; then
  fail "review dir not found: $review_dir"
fi

manifest="$review_dir/manifest.env"
samples_dir="$review_dir/samples"

if [[ ! -f "$manifest" ]]; then
  fail "missing manifest.env"
fi
if [[ ! -d "$samples_dir" ]]; then
  fail "missing samples dir"
fi

read_field() {
  local file="$1"
  local key="$2"
  grep -E "^${key}=" "$file" | head -n 1 | cut -d= -f2-
}

require_manifest_field() {
  local key="$1"
  local expected="$2"
  local actual
  actual="$(read_field "$manifest" "$key" || true)"
  if [[ "$actual" != "$expected" ]]; then
    fail "$key must be $expected"
  fi
}

require_manifest_field "artifact_source" "github_actions"
require_manifest_field "artifact_name" "flutter-web-timing-report"
require_manifest_field "sample_count" "5"
require_manifest_field "sample_files" "sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env"
require_manifest_field "aggregate_output" "aggregate.env"
require_manifest_field "no_phi_attestation" "true"

review_id="$(read_field "$manifest" "review_id" || true)"
if [[ -z "$review_id" ]]; then
  fail "review_id is required"
fi

review_status="$(read_field "$manifest" "review_status" || true)"
case "$review_status" in
  proposed|approved|rejected)
    ;;
  *)
    fail "review_status must be proposed, approved, or rejected"
    ;;
esac

IFS=',' read -r -a sample_files <<<"$(read_field "$manifest" "sample_files")"
for sample_file in "${sample_files[@]}"; do
  if [[ ! -f "$samples_dir/$sample_file" ]]; then
    fail "missing sample file: $sample_file"
  fi
done

aggregate_file="$review_dir/$(read_field "$manifest" "aggregate_output")"
if [[ ! -f "$aggregate_file" ]]; then
  fail "missing aggregate file: ${aggregate_file#$review_dir/}"
fi

forbidden_pattern='^(user_id|patient_id|device_id|access_token|refresh_token|raw_channels|s_det|s_ref|primary_value)='
if grep -REq "$forbidden_pattern" "$manifest" "$samples_dir" "$aggregate_file"; then
  fail "forbidden fixture field found"
fi

tmp_aggregate="$(mktemp)"
trap 'rm -f "$tmp_aggregate"' EXIT
bash "$AGGREGATOR" --input-dir "$samples_dir" --output "$tmp_aggregate" >/dev/null

if ! cmp -s "$tmp_aggregate" "$aggregate_file"; then
  fail "aggregate.env is stale or not reproducible"
fi

echo "FLUTTER_WEB_TIMING_REVIEW_FIXTURE_VALIDATE_PASS $review_dir"
