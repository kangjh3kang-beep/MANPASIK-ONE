#!/usr/bin/env bash
set -euo pipefail

log_path=""
branch_type=""
runner_context=""
output_path=""
source_marker="FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"

fail() {
  echo "FLUTTER_WEB_TIMING_REPORT_FAIL $1" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --log)
      shift
      log_path="${1:-}"
      ;;
    --branch-type)
      shift
      branch_type="${1:-}"
      ;;
    --runner-context)
      shift
      runner_context="${1:-}"
      ;;
    --output)
      shift
      output_path="${1:-}"
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

[[ -n "$log_path" ]] || fail "missing --log"
[[ -n "$branch_type" ]] || fail "missing --branch-type"
[[ -n "$runner_context" ]] || fail "missing --runner-context"
[[ -n "$output_path" ]] || fail "missing --output"
[[ -f "$log_path" ]] || fail "log file not found: $log_path"

mapfile -t durations < <(
  grep -Eo "${source_marker}=[0-9]+" "$log_path" |
    sed "s/^${source_marker}=//"
)

if [[ "${#durations[@]}" -eq 0 ]]; then
  fail "no duration marker found"
fi

min_duration="${durations[0]}"
max_duration="${durations[0]}"
latest_duration="${durations[${#durations[@]} - 1]}"
for duration in "${durations[@]}"; do
  if (( duration < min_duration )); then
    min_duration="$duration"
  fi
  if (( duration > max_duration )); then
    max_duration="$duration"
  fi
done

duration_values="$(IFS=,; echo "${durations[*]}")"
mkdir -p "$(dirname "$output_path")"
{
  echo "report_version=1"
  echo "source_marker=$source_marker"
  echo "sample_count=${#durations[@]}"
  echo "duration_seconds_values=$duration_values"
  echo "latest_duration_seconds=$latest_duration"
  echo "min_duration_seconds=$min_duration"
  echo "max_duration_seconds=$max_duration"
  echo "branch_type=$branch_type"
  echo "runner_context=$runner_context"
} >"$output_path"

echo "FLUTTER_WEB_TIMING_REPORT_PASS $output_path"
