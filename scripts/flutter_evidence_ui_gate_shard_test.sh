#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GATE="$ROOT/scripts/flutter_evidence_ui_gate.sh"

fail() {
  echo "FLUTTER_EVIDENCE_UI_GATE_SHARD_TEST_FAIL $1" >&2
  exit 1
}

require_contains() {
  local haystack="$1"
  local needle="$2"
  local label="$3"
  if [[ "$haystack" != *"$needle"* ]]; then
    fail "$label missing: $needle"
  fi
}

require_count() {
  local expected="$1"
  shift
  local actual
  actual="$(bash "$GATE" "$@")"
  if [[ "$actual" != "$expected" ]]; then
    fail "expected count $expected for '$*', got '$actual'"
  fi
}

categories="$(bash "$GATE" --list-categories)"
require_contains "$categories" "contract" "categories"
require_contains "$categories" "ui" "categories"
require_contains "$categories" "data" "categories"

require_count "13" --count
require_count "4" --count --category contract
require_count "4" --count --category ui
require_count "5" --count --category data

full_count="$(bash "$GATE" --count)"
contract_count="$(bash "$GATE" --count --category contract)"
ui_count="$(bash "$GATE" --count --category ui)"
data_count="$(bash "$GATE" --count --category data)"
if (( contract_count + ui_count + data_count != full_count )); then
  fail "category counts do not sum to full count"
fi

contract_list="$(bash "$GATE" --list --category contract)"
require_contains "$contract_list" "test/features/measurement/domain/measurement_evidence_presentation_test.dart" "contract shard"
require_contains "$contract_list" "test/generated/measurement_result_evidence_contract_test.dart" "contract shard"
require_contains "$contract_list" "test/generated/measurement_summary_evidence_contract_test.dart" "contract shard"

ui_list="$(bash "$GATE" --list --category ui)"
require_contains "$ui_list" "test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart" "ui shard"
require_contains "$ui_list" "test/features/home/presentation/home_measurement_evidence_badge_test.dart" "ui shard"
require_contains "$ui_list" "test/features/data_hub/presentation/data_hub_layout_smoke_test.dart" "ui shard"

data_list="$(bash "$GATE" --list --category data)"
require_contains "$data_list" "test/features/measurement/data/measurement_repository_impl_test.dart" "data shard"
require_contains "$data_list" "test/features/measurement/data/measurement_repository_rest_test.dart" "data shard"
require_contains "$data_list" "test/features/data_hub/data/data_hub_repository_rest_test.dart" "data shard"

if bash "$GATE" --count --category unknown >/dev/null 2>&1; then
  fail "unknown category should fail"
fi

echo "FLUTTER_EVIDENCE_UI_GATE_SHARD_TEST_PASS"
