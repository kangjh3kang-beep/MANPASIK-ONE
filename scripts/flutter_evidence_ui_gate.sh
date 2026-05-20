#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEFAULT_FLUTTER_BIN="/mnt/d/우리집/flutter_cache/flutter/bin/flutter"

contract_tests=(
  "test/features/measurement/domain/measurement_evidence_presentation_test.dart"
  "test/features/measurement/application/measurement_golden_path_orchestrator_test.dart"
  "test/generated/measurement_result_evidence_contract_test.dart"
  "test/generated/measurement_summary_evidence_contract_test.dart"
)

ui_tests=(
  "test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart"
  "test/features/home/presentation/home_measurement_evidence_badge_test.dart"
  "test/features/data_hub/presentation/data_hub_evidence_badge_test.dart"
  "test/features/data_hub/presentation/data_hub_layout_smoke_test.dart"
)

data_tests=(
  "test/features/measurement/data/measurement_process_gateway_mapper_test.dart"
  "test/features/measurement/data/measurement_repository_impl_test.dart"
  "test/features/measurement/data/measurement_repository_rest_test.dart"
  "test/features/data_hub/domain/data_hub_domain_test.dart"
  "test/features/data_hub/data/data_hub_repository_rest_test.dart"
)

tests=(
  "${contract_tests[@]}"
  "${ui_tests[@]}"
  "${data_tests[@]}"
)

fail() {
  echo "FLUTTER_EVIDENCE_UI_GATE_FAIL $1" >&2
  exit 1
}

mode="run"
category=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --list)
      mode="list"
      shift
      ;;
    --count)
      mode="count"
      shift
      ;;
    --list-categories)
      mode="list-categories"
      shift
      ;;
    --category)
      shift
      if [[ -z "${1:-}" ]]; then
        fail "missing category"
      fi
      category="$1"
      shift
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
done

if [[ "$mode" == "list-categories" ]]; then
  printf "%s\n" "contract" "ui" "data"
  exit 0
fi

case "$category" in
  "")
    selected_tests=("${tests[@]}")
    ;;
  contract)
    selected_tests=("${contract_tests[@]}")
    ;;
  ui)
    selected_tests=("${ui_tests[@]}")
    ;;
  data)
    selected_tests=("${data_tests[@]}")
    ;;
  *)
    fail "unknown category: $category"
    ;;
esac

if [[ "$mode" == "list" ]]; then
  printf "%s\n" "${selected_tests[@]}"
  exit 0
fi

if [[ "$mode" == "count" ]]; then
  echo "${#selected_tests[@]}"
  exit 0
fi

FLUTTER_BIN="${FLUTTER_BIN:-$DEFAULT_FLUTTER_BIN}"
if [[ ! -x "$FLUTTER_BIN" ]]; then
  FLUTTER_BIN="${FLUTTER_BIN_FALLBACK:-flutter}"
fi

cd "$ROOT/frontend/flutter-app"
"$FLUTTER_BIN" test --no-pub "${selected_tests[@]}"
echo "FLUTTER_EVIDENCE_UI_GATE_PASS"
