#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-evidence-ui-gate-policy.md"
GATE="$ROOT/scripts/flutter_evidence_ui_gate.sh"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "FLUTTER_EVIDENCE_UI_GATE_POLICY_FAIL $1" >&2
  exit 1
}

if [[ ! -f "$POLICY" ]]; then
  fail "missing policy document"
fi

required_markers=(
  "gate_scope: evidence_ui_contract"
  "required_on_pull_request: true"
  "duplicates_full_flutter_test: true"
  "max_test_files: 13"
  "category_shards: contract,ui,data"
  "max_test_files_per_category: 5"
  "ci_execution_mode: category_shards"
  "ci_shard_steps: contract,ui,data"
  "aggregate_ci_step: false"
  "aggregate_local_gate_supported: true"
  "gate_script: scripts/flutter_evidence_ui_gate.sh"
  "policy_test: scripts/flutter_evidence_ui_gate_policy_test.sh"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq "$marker" "$POLICY"; then
    fail "missing policy marker: $marker"
  fi
done

max_test_files="$(sed -n 's/^max_test_files: //p' "$POLICY" | head -n 1)"
if [[ ! "$max_test_files" =~ ^[0-9]+$ ]]; then
  fail "max_test_files is not numeric"
fi

gate_count="$(bash "$GATE" --count)"
if [[ ! "$gate_count" =~ ^[0-9]+$ ]]; then
  fail "gate count is not numeric"
fi

if (( gate_count > max_test_files )); then
  fail "gate count $gate_count exceeds max_test_files $max_test_files"
fi

max_test_files_per_category="$(sed -n 's/^max_test_files_per_category: //p' "$POLICY" | head -n 1)"
if [[ ! "$max_test_files_per_category" =~ ^[0-9]+$ ]]; then
  fail "max_test_files_per_category is not numeric"
fi

IFS=',' read -r -a category_shards <<< "$(sed -n 's/^category_shards: //p' "$POLICY" | head -n 1)"
if [[ "${#category_shards[@]}" -eq 0 ]]; then
  fail "category_shards is empty"
fi

for category in "${category_shards[@]}"; do
  category_count="$(bash "$GATE" --count --category "$category")"
  if [[ ! "$category_count" =~ ^[0-9]+$ ]]; then
    fail "category count is not numeric for $category"
  fi
  if (( category_count > max_test_files_per_category )); then
    fail "category $category count $category_count exceeds max_test_files_per_category $max_test_files_per_category"
  fi
done

ci_execution_mode="$(sed -n 's/^ci_execution_mode: //p' "$POLICY" | head -n 1)"
if [[ "$ci_execution_mode" != "category_shards" ]]; then
  fail "unsupported ci_execution_mode: $ci_execution_mode"
fi

IFS=',' read -r -a ci_shard_steps <<< "$(sed -n 's/^ci_shard_steps: //p' "$POLICY" | head -n 1)"
for category in "${ci_shard_steps[@]}"; do
  if ! grep -Fq "scripts/flutter_evidence_ui_gate.sh --category $category" "$WORKFLOW"; then
    fail "CI workflow missing shard command for $category"
  fi
done

aggregate_ci_step="$(sed -n 's/^aggregate_ci_step: //p' "$POLICY" | head -n 1)"
if [[ "$aggregate_ci_step" == "false" ]] &&
  grep -Fxq "        run: bash ../../scripts/flutter_evidence_ui_gate.sh" "$WORKFLOW"; then
  fail "aggregate CI step is disabled but workflow still runs aggregate command"
fi

if ! grep -Fq "Flutter evidence UI gate policy" "$WORKFLOW"; then
  fail "CI workflow missing policy step name"
fi

if ! grep -Fq "scripts/flutter_evidence_ui_gate_policy_test.sh" "$WORKFLOW"; then
  fail "CI workflow missing policy script path"
fi

echo "FLUTTER_EVIDENCE_UI_GATE_POLICY_PASS"
