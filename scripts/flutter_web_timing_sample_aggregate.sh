#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'EOF_USAGE'
usage: flutter_web_timing_sample_aggregate.sh --input-dir <dir> --output <file> [--min-samples <n>]
EOF_USAGE
}

fail() {
  echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_FAIL $1" >&2
  exit 1
}

input_dir=""
output_path=""
min_samples=5

while [[ $# -gt 0 ]]; do
  case "$1" in
    --input-dir)
      input_dir="${2:-}"
      shift 2
      ;;
    --output)
      output_path="${2:-}"
      shift 2
      ;;
    --min-samples)
      min_samples="${2:-}"
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

if [[ -z "$input_dir" || -z "$output_path" ]]; then
  usage
  fail "missing required arguments"
fi

if [[ ! "$min_samples" =~ ^[0-9]+$ || "$min_samples" -lt 1 ]]; then
  fail "min samples must be a positive integer"
fi

if [[ ! -d "$input_dir" ]]; then
  fail "input dir not found: $input_dir"
fi

read_field() {
  local file="$1"
  local key="$2"
  grep -E "^${key}=" "$file" | head -n 1 | cut -d= -f2-
}

contains_value() {
  local needle="$1"
  shift
  local value
  for value in "$@"; do
    if [[ "$value" == "$needle" ]]; then
      return 0
    fi
  done
  return 1
}

join_by_comma() {
  local IFS=,
  echo "$*"
}

durations=()
branch_types=()
runner_contexts=()

shopt -s nullglob
files=("$input_dir"/*.env)
shopt -u nullglob

if [[ "${#files[@]}" -eq 0 ]]; then
  fail "no env artifact files found"
fi

for file in "${files[@]}"; do
  report_version="$(read_field "$file" "report_version" || true)"
  if [[ -z "$report_version" ]]; then
    continue
  fi
  if [[ "$report_version" != "1" ]]; then
    fail "unsupported report version in ${file##*/}: $report_version"
  fi

  source_marker="$(read_field "$file" "source_marker" || true)"
  duration="$(read_field "$file" "latest_duration_seconds" || true)"
  branch_type="$(read_field "$file" "branch_type" || true)"
  runner_context="$(read_field "$file" "runner_context" || true)"

  if [[ "$source_marker" != "FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS" ]]; then
    fail "invalid source marker in ${file##*/}"
  fi
  if [[ ! "$duration" =~ ^[0-9]+$ ]]; then
    fail "invalid duration in ${file##*/}"
  fi
  if [[ -z "$branch_type" ]]; then
    fail "missing branch type in ${file##*/}"
  fi
  if [[ -z "$runner_context" ]]; then
    fail "missing runner context in ${file##*/}"
  fi

  durations+=("$duration")
  if ! contains_value "$branch_type" "${branch_types[@]}"; then
    branch_types+=("$branch_type")
  fi
  if ! contains_value "$runner_context" "${runner_contexts[@]}"; then
    runner_contexts+=("$runner_context")
  fi
done

sample_count="${#durations[@]}"
if [[ "$sample_count" -lt "$min_samples" ]]; then
  fail "minimum samples not met: $sample_count < $min_samples"
fi

if ! contains_value "pull_request" "${branch_types[@]}"; then
  fail "missing pull_request samples"
fi
if ! contains_value "release_branch" "${branch_types[@]}"; then
  fail "missing release_branch samples"
fi

mapfile -t sorted_durations < <(printf '%s\n' "${durations[@]}" | sort -n)
mapfile -t sorted_branch_types < <(printf '%s\n' "${branch_types[@]}" | sort)
mapfile -t sorted_runner_contexts < <(printf '%s\n' "${runner_contexts[@]}" | sort)

middle="$((sample_count / 2))"
if [[ "$((sample_count % 2))" -eq 1 ]]; then
  median_seconds="${sorted_durations[$middle]}"
else
  left="${sorted_durations[$((middle - 1))]}"
  right="${sorted_durations[$middle]}"
  median_seconds="$(((left + right) / 2))"
fi

p95_index="$(((95 * sample_count + 99) / 100 - 1))"
p95_seconds="${sorted_durations[$p95_index]}"
worst_case_seconds="${sorted_durations[$((sample_count - 1))]}"

{
  echo "aggregate_version=1"
  echo "source_artifact=flutter-web-timing-report"
  echo "sample_count=$sample_count"
  echo "duration_seconds_values=$(join_by_comma "${sorted_durations[@]}")"
  echo "median_seconds=$median_seconds"
  echo "p95_seconds=$p95_seconds"
  echo "worst_case_seconds=$worst_case_seconds"
  echo "branch_types=$(join_by_comma "${sorted_branch_types[@]}")"
  echo "runner_contexts=$(join_by_comma "${sorted_runner_contexts[@]}")"
} >"$output_path"

echo "FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_PASS $output_path"
