#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GATE="$ROOT/scripts/flutter_web_release_gate.sh"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

pass_log="$TMP_DIR/pass_with_wasm_warning.log"
fail_log="$TMP_DIR/missing_build_marker.log"

cat >"$pass_log" <<'LOG'
Compiling lib/main.dart for the Web...
✓ Built build/web
Wasm dry run findings:
Found incompatibilities with WebAssembly.
LOG

cat >"$fail_log" <<'LOG'
Compiling lib/main.dart for the Web...
Build failed.
LOG

output="$(FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=42 bash "$GATE" --policy-check "$pass_log")"
if [[ "$output" != *"FLUTTER_WEB_RELEASE_GATE_PASS"* ]]; then
  echo "expected pass marker in policy-check output" >&2
  exit 1
fi

if [[ "$output" != *"FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=42"* ]]; then
  echo "expected duration seconds marker in policy-check output" >&2
  exit 1
fi

if [[ "$output" != *"FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN"* ]]; then
  echo "expected wasm warning marker in policy-check output" >&2
  exit 1
fi

if bash "$GATE" --policy-check "$fail_log" >/tmp/flutter_web_gate_fail.out 2>&1; then
  echo "expected missing build marker log to fail" >&2
  exit 1
fi

echo "FLUTTER_WEB_RELEASE_GATE_TEST_PASS"
