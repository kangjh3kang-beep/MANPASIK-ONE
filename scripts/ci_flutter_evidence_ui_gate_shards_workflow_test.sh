#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_FLUTTER_EVIDENCE_UI_GATE_SHARDS_WORKFLOW_FAIL $1" >&2
  exit 1
}

if [[ ! -f "$WORKFLOW" ]]; then
  fail "missing workflow"
fi

required_markers=(
  "Flutter evidence UI gate: contract"
  "Flutter evidence UI gate: ui"
  "Flutter evidence UI gate: data"
  "scripts/flutter_evidence_ui_gate.sh --category contract"
  "scripts/flutter_evidence_ui_gate.sh --category ui"
  "scripts/flutter_evidence_ui_gate.sh --category data"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$WORKFLOW"; then
    fail "missing workflow marker: $marker"
  fi
done

if grep -Fxq "        run: bash ../../scripts/flutter_evidence_ui_gate.sh" "$WORKFLOW"; then
  fail "aggregate-only evidence gate command still present"
fi

echo "CI_FLUTTER_EVIDENCE_UI_GATE_SHARDS_WORKFLOW_PASS"
