#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_FAIL $1" >&2
  exit 1
}

if [[ ! -f "$WORKFLOW" ]]; then
  fail "missing workflow"
fi

required_markers=(
  "Flutter web release gate"
  "Flutter web timing report"
  "scripts/flutter_web_timing_report.sh"
  "--log /tmp/manpasik_flutter_web_build.log"
  "--output /tmp/manpasik_flutter_web_timing.env"
  "github.event_name"
  "pull_request"
  "runner.os"
  "runner.arch"
  "Upload Flutter web timing report"
  "actions/upload-artifact@v4"
  "name: flutter-web-timing-report"
  "path: /tmp/manpasik_flutter_web_timing.env"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq -- "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

echo "CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_PASS"
