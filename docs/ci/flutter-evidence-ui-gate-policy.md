# Flutter Evidence UI Gate Policy

```yaml
gate_scope: evidence_ui_contract
required_on_pull_request: true
duplicates_full_flutter_test: true
max_test_files: 13
category_shards: contract,ui,data
max_test_files_per_category: 5
ci_execution_mode: category_shards
ci_shard_steps: contract,ui,data
aggregate_ci_step: false
aggregate_local_gate_supported: true
gate_script: scripts/flutter_evidence_ui_gate.sh
policy_test: scripts/flutter_evidence_ui_gate_policy_test.sh
```

## Policy

The Flutter evidence UI gate is a named regression gate for:

- Measurement evidence contract fields.
- Regulatory-safe evidence copy.
- Measurement result and home badge display.
- DataHub evidence badge display and layout smoke coverage.
- Native and REST mapper compatibility for snake_case and legacy camelCase evidence fields.

This gate is required on pull requests. It currently duplicates a small subset of the full Flutter test job so evidence regressions are visible as a named CI failure rather than being buried inside the general test output.

The gate must stay bounded. If more than `max_test_files` test files are needed, first split the gate into narrower categories or update this policy with an explicit rationale.

The current categories are:

- `contract`: generated protobuf evidence contracts and measurement evidence presentation/snapshot behavior.
- `ui`: safe evidence badge rendering on measurement, home, and DataHub surfaces.
- `data`: native/REST repository and mapper compatibility for evidence fields.

New evidence tests must join the closest category. If a category would exceed `max_test_files_per_category`, split that category or update this policy with the reason.

CI runs this gate as category shard steps. Aggregate local execution remains supported through:

```bash
bash scripts/flutter_evidence_ui_gate.sh
```

The aggregate local gate is useful for developer verification, but CI must keep category steps separate while `ci_execution_mode` is `category_shards`.

## When This Changes

If the full Flutter test job is later sharded by category, these category shards can become separate CI steps instead of one aggregate evidence gate. Keep `scripts/flutter_evidence_ui_gate.sh --count`, `--count --category <name>`, and this policy guard in place so CI cost changes are intentional.
