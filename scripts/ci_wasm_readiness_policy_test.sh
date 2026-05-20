#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
READINESS="$ROOT/docs/ci/flutter-wasm-readiness.md"

if [[ ! -f "$READINESS" ]]; then
  echo "CI_WASM_READINESS_POLICY_FAIL missing readiness document" >&2
  exit 1
fi

required_markers=(
  "wasm_release_target: false"
  "flutter_secure_storage_web"
  "share_plus"
  "connectivity_plus"
  "package:js"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$READINESS"; then
    echo "CI_WASM_READINESS_POLICY_FAIL missing marker: $marker" >&2
    exit 1
  fi
done

echo "CI_WASM_READINESS_POLICY_PASS"
