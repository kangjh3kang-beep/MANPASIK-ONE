#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEFAULT_FLUTTER_BIN="/mnt/d/우리집/flutter_cache/flutter/bin/flutter"
LOG_PATH="${FLUTTER_WEB_RELEASE_GATE_LOG:-/tmp/manpasik_flutter_web_build.log}"

policy_check() {
  local log_file="$1"
  if [[ ! -f "$log_file" ]]; then
    echo "FLUTTER_WEB_RELEASE_GATE_FAIL missing log file: $log_file" >&2
    return 1
  fi

  if ! grep -Fq "Built build/web" "$log_file"; then
    echo "FLUTTER_WEB_RELEASE_GATE_FAIL build/web success marker missing" >&2
    return 1
  fi

  echo "FLUTTER_WEB_RELEASE_GATE_PASS"

  if [[ -n "${FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS:-}" ]]; then
    if [[ ! "$FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS" =~ ^[0-9]+$ ]]; then
      echo "FLUTTER_WEB_RELEASE_GATE_FAIL invalid duration seconds" >&2
      return 1
    fi
    echo "FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=$FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS"
  fi

  if grep -Fq "Wasm dry run findings:" "$log_file"; then
    echo "FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN"
  fi
}

if [[ "${1:-}" == "--policy-check" ]]; then
  if [[ $# -ne 2 ]]; then
    echo "usage: $0 --policy-check <build-log>" >&2
    exit 2
  fi
  policy_check "$2"
  exit $?
fi

FLUTTER_BIN="${FLUTTER_BIN:-$DEFAULT_FLUTTER_BIN}"
if [[ ! -x "$FLUTTER_BIN" ]]; then
  FLUTTER_BIN="${FLUTTER_BIN_FALLBACK:-flutter}"
fi

cd "$ROOT/frontend/flutter-app"
start_time="$(date +%s)"
"$FLUTTER_BIN" build web --no-pub >"$LOG_PATH" 2>&1
end_time="$(date +%s)"
duration_seconds="$((end_time - start_time))"
cat "$LOG_PATH"
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS="$duration_seconds" policy_check "$LOG_PATH"
