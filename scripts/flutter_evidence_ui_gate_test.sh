#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GATE="$ROOT/scripts/flutter_evidence_ui_gate.sh"

output="$(bash "$GATE" --list)"

required_tests=(
  "test/features/measurement/domain/measurement_evidence_presentation_test.dart"
  "test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart"
  "test/features/home/presentation/home_measurement_evidence_badge_test.dart"
  "test/features/data_hub/presentation/data_hub_evidence_badge_test.dart"
  "test/features/data_hub/presentation/data_hub_layout_smoke_test.dart"
  "test/features/measurement/data/measurement_repository_rest_test.dart"
  "test/features/data_hub/data/data_hub_repository_rest_test.dart"
  "test/generated/measurement_result_evidence_contract_test.dart"
  "test/generated/measurement_summary_evidence_contract_test.dart"
)

for test_path in "${required_tests[@]}"; do
  if [[ "$output" != *"$test_path"* ]]; then
    echo "FLUTTER_EVIDENCE_UI_GATE_TEST_FAIL missing test: $test_path" >&2
    exit 1
  fi
done

echo "FLUTTER_EVIDENCE_UI_GATE_TEST_PASS"
