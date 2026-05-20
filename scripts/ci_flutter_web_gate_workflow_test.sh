#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

if [[ ! -f "$WORKFLOW" ]]; then
  echo "CI_FLUTTER_WEB_GATE_WORKFLOW_FAIL missing workflow" >&2
  exit 1
fi

if ! grep -Fq "Flutter web release gate" "$WORKFLOW"; then
  echo "CI_FLUTTER_WEB_GATE_WORKFLOW_FAIL missing step name" >&2
  exit 1
fi

if ! grep -Fq "scripts/flutter_web_release_gate.sh" "$WORKFLOW"; then
  echo "CI_FLUTTER_WEB_GATE_WORKFLOW_FAIL missing release gate script path" >&2
  exit 1
fi

echo "CI_FLUTTER_WEB_GATE_WORKFLOW_PASS"
