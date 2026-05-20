# Flutter Web Timing Review Fixtures Policy

```yaml
review_fixture_policy_version: 1
fixture_root: docs/ci/fixtures/flutter-web-timing
fixture_manifest: manifest.env
required_exported_artifacts: 5
required_sample_files: sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env
aggregate_output_file: aggregate.env
aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh
audit_template: docs/audit/templates/flutter-web-timing-review-template.md
promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md
review_status_allowed: proposed,approved,rejected
artifact_source_required: github_actions
no_phi_allowed: true
fixture_validator_script: scripts/flutter_web_timing_review_fixture_validate.sh
fixture_validator_test: scripts/flutter_web_timing_review_fixture_validate_test.sh
manifest_required_fields: review_id,artifact_source,artifact_name,sample_count,sample_files,aggregate_output,review_status,no_phi_attestation
forbidden_fixture_fields: user_id,patient_id,device_id,access_token,refresh_token,raw_channels,s_det,s_ref,primary_value
threshold_change_review_check_script: scripts/flutter_web_timing_threshold_change_review_check.sh
threshold_change_review_check_test: scripts/flutter_web_timing_threshold_change_review_check_test.sh
threshold_change_review_check_required: true
threshold_change_review_required_fields: review_id,decision,sample_count,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts,validator_command
```

## Fixture Layout

Store each real timing review under a review-specific directory:

```text
docs/ci/fixtures/flutter-web-timing/<review-id>/
  manifest.env
  samples/
    sample-01.env
    sample-02.env
    sample-03.env
    sample-04.env
    sample-05.env
  aggregate.env
```

`<review-id>` should use `YYYY-MM-DD-<short-purpose>`, for example `2026-05-14-pr-gate-cost-review`.

## Manifest

Each `manifest.env` must record:

```text
review_id=<YYYY-MM-DD-short-purpose>
artifact_source=github_actions
artifact_name=flutter-web-timing-report
sample_count=5
sample_files=sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env
aggregate_output=aggregate.env
review_status=proposed
no_phi_attestation=true
```

Allowed `review_status` values are `proposed`, `approved`, and `rejected`.

## Aggregate Generation

Generate `aggregate.env` from the exported samples with:

```bash
bash scripts/flutter_web_timing_sample_aggregate.sh \
  --input-dir docs/ci/fixtures/flutter-web-timing/<review-id>/samples \
  --output docs/ci/fixtures/flutter-web-timing/<review-id>/aggregate.env
```

## Safety Rules

- Fixtures must come from GitHub Actions `flutter-web-timing-report` artifacts.
- Fixtures must not include PHI, user identifiers, device identifiers, access tokens, or raw medical measurements.
- Do not change Flutter web gate thresholds unless the audit record uses `docs/audit/templates/flutter-web-timing-review-template.md`.

## Validator

Before using a fixture in a threshold change PR, validate the concrete review directory:

```bash
bash scripts/flutter_web_timing_review_fixture_validate.sh \
  --review-dir docs/ci/fixtures/flutter-web-timing/<review-id>
```

The validator checks `manifest.env`, all five sample files, `aggregate.env`, no-PHI attestation, forbidden fixture fields, and aggregate reproducibility by re-running `scripts/flutter_web_timing_sample_aggregate.sh`.

## Threshold Change Review Check

Before changing Flutter web timing thresholds or PR gate policy, run the wrapper command against the concrete fixture and completed audit document:

```bash
bash scripts/flutter_web_timing_threshold_change_review_check.sh \
  --review-dir docs/ci/fixtures/flutter-web-timing/<review-id> \
  --audit-file docs/audit/<review-doc>.md
```

The audit document must record `review_id`, `decision`, `sample_count`, `median_seconds`, `p95_seconds`, `worst_case_seconds`, `branch_types`, `runner_contexts`, and `validator_command` exactly as produced by the fixture manifest and `aggregate.env`.
