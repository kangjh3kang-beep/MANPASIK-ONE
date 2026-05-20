# Flutter Web Timing Threshold Promotion Policy

```yaml
promotion_policy_version: 1
sample_source_artifact: flutter-web-timing-report
minimum_successful_artifact_samples: 5
required_branch_types: pull_request,release_branch
required_runner_context_field: runner_context
required_duration_field: latest_duration_seconds
required_distribution_fields: min_duration_seconds,max_duration_seconds,duration_seconds_values
required_statistics: median_seconds,p95_seconds,worst_case_seconds
threshold_change_requires: measured_cost_review
automated_gate_relaxation: false
relaxation_allowed_before_minimum_samples: false
nightly_split_change_requires: docs/ci/flutter-web-release-gate-policy.md
review_record_required: docs/audit
sample_aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh
sample_aggregator_test: scripts/flutter_web_timing_sample_aggregate_test.sh
aggregation_output_fields: aggregate_version,source_artifact,sample_count,duration_seconds_values,median_seconds,p95_seconds,worst_case_seconds,branch_types,runner_contexts
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
review_fixture_root: docs/ci/fixtures/flutter-web-timing
review_audit_template: docs/audit/templates/flutter-web-timing-review-template.md
```

## Policy

The `flutter-web-timing-report` artifact is the only accepted sample source for changing Flutter web gate timing thresholds or relaxing the pull request web release gate.

At least 5 successful artifact samples are required before proposing any threshold change. Samples must include both `pull_request` and `release_branch` branch types and must preserve `runner_context` so runner drift does not distort the decision.

## Required Review Data

The measured cost review must compute and record:

- `median_seconds`
- `p95_seconds`
- `worst_case_seconds`

The review must also keep the source artifact fields:

- `latest_duration_seconds`
- `min_duration_seconds`
- `max_duration_seconds`
- `duration_seconds_values`
- `branch_type`
- `runner_context`

## Offline Aggregator

Use `scripts/flutter_web_timing_sample_aggregate.sh` after exporting at least 5 `flutter-web-timing-report` artifacts into one local directory:

```bash
bash scripts/flutter_web_timing_sample_aggregate.sh \
  --input-dir exported-artifacts \
  --output flutter-web-timing-aggregate.env
```

The aggregate output must include:

```text
aggregate_version=1
source_artifact=flutter-web-timing-report
sample_count=<n>
duration_seconds_values=<sorted comma-separated durations>
median_seconds=<median>
p95_seconds=<nearest-rank-p95>
worst_case_seconds=<max>
branch_types=<comma-separated unique branch types>
runner_contexts=<comma-separated unique runner contexts>
```

## Change Rules

- Do not relax the PR-required Flutter web release gate before 5 successful artifact samples exist.
- Do not enable automated gate relaxation. A human measured cost review is required.
- Do not change `nightly_split` without updating `docs/ci/flutter-web-release-gate-policy.md`.
- Record the review in `docs/audit` before changing thresholds or execution mode.
- Use `docs/ci/flutter-web-timing-review-fixtures.md` and `docs/audit/templates/flutter-web-timing-review-template.md` for real 5-sample reviews.
