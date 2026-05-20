#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REVIEW_DIR_REL="docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review"
REVIEW_DIR="$ROOT/$REVIEW_DIR_REL"
AUDIT_FILE="$ROOT/docs/audit/flutter-web-timing-synthetic-review-example.md"

fail() {
  echo "CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_FAIL $1" >&2
  exit 1
}

if [[ ! -d "$REVIEW_DIR" ]]; then
  fail "missing synthetic review dir"
fi
if [[ ! -f "$AUDIT_FILE" ]]; then
  fail "missing synthetic audit file"
fi
if ! grep -Fq "synthetic_example=true" "$REVIEW_DIR/manifest.env"; then
  fail "manifest missing synthetic marker"
fi
if ! grep -Fq "synthetic_review: true" "$AUDIT_FILE"; then
  fail "audit missing synthetic marker"
fi

bash "$ROOT/scripts/flutter_web_timing_review_fixture_validate.sh" --review-dir "$REVIEW_DIR" >/dev/null
bash "$ROOT/scripts/flutter_web_timing_threshold_change_review_check.sh" \
  --review-dir "$REVIEW_DIR_REL" \
  --audit-file "$AUDIT_FILE" >/dev/null

echo "CI_WEB_TIMING_SYNTHETIC_REVIEW_EXAMPLE_PASS"
