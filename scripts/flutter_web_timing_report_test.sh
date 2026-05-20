#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORTER="$ROOT/scripts/flutter_web_timing_report.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

pass_log="$TMP_DIR/pass.log"
missing_log="$TMP_DIR/missing.log"
report="$TMP_DIR/report.env"

cat >"$pass_log" <<'LOG'
Compiling lib/main.dart for the Web...
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=11
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=17
FLUTTER_WEB_RELEASE_GATE_PASS
LOG

cat >"$missing_log" <<'LOG'
Compiling lib/main.dart for the Web...
FLUTTER_WEB_RELEASE_GATE_PASS
LOG

bash "$REPORTER" \
  --log "$pass_log" \
  --branch-type pull_request \
  --runner-context ubuntu-latest \
  --output "$report"

required_lines=(
  "report_version=1"
  "source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"
  "sample_count=2"
  "duration_seconds_values=11,17"
  "latest_duration_seconds=17"
  "min_duration_seconds=11"
  "max_duration_seconds=17"
  "branch_type=pull_request"
  "runner_context=ubuntu-latest"
)

for line in "${required_lines[@]}"; do
  if ! grep -Fxq "$line" "$report"; then
    echo "FLUTTER_WEB_TIMING_REPORT_TEST_FAIL missing report line: $line" >&2
    exit 1
  fi
done

if bash "$REPORTER" \
  --log "$missing_log" \
  --branch-type pull_request \
  --runner-context ubuntu-latest \
  --output "$TMP_DIR/missing.env" >/dev/null 2>&1; then
  echo "FLUTTER_WEB_TIMING_REPORT_TEST_FAIL missing marker log should fail" >&2
  exit 1
fi

echo "FLUTTER_WEB_TIMING_REPORT_TEST_PASS"
