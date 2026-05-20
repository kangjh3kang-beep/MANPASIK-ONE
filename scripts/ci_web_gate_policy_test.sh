#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-release-gate-policy.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

if [[ ! -f "$POLICY" ]]; then
  echo "CI_WEB_GATE_POLICY_FAIL missing policy document" >&2
  exit 1
fi

if ! grep -Fq "required_on_pull_request: true" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing PR-required marker" >&2
  exit 1
fi

if ! grep -Fq "wasm_dry_run_blocking: false" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing non-blocking Wasm marker" >&2
  exit 1
fi

if ! grep -Fq "release_target: js_web" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing JS web release target marker" >&2
  exit 1
fi

if ! grep -Fq "ci_execution_mode: pull_request_and_release_branch" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing CI execution mode marker" >&2
  exit 1
fi

if ! grep -Fq "release_branch_required: true" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing release branch required marker" >&2
  exit 1
fi

if ! grep -Fq "nightly_split: false" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing nightly split marker" >&2
  exit 1
fi

if ! grep -Fq "cost_review_required_before_relaxing: true" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing cost review marker" >&2
  exit 1
fi

if ! grep -Fq "blocking_build_command: flutter build web --no-pub" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing blocking build command marker" >&2
  exit 1
fi

if ! grep -Fq "timing_capture: true" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing timing capture marker" >&2
  exit 1
fi

if ! grep -Fq "timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS" "$POLICY"; then
  echo "CI_WEB_GATE_POLICY_FAIL missing timing marker" >&2
  exit 1
fi

if ! grep -Fq "Flutter web release gate" "$WORKFLOW"; then
  echo "CI_WEB_GATE_POLICY_FAIL CI workflow is not connected" >&2
  exit 1
fi

echo "CI_WEB_GATE_POLICY_PASS"
