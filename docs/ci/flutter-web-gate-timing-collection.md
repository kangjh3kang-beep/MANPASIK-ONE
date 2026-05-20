# Flutter Web Gate Timing Collection Policy

```yaml
timing_collection_version: 1
source_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
log_query_pattern: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=
minimum_samples_before_relaxing: 5
threshold_policy_status: advisory
blocking_threshold_seconds: unset
nightly_split_requires: measured_cost_review
current_nightly_split: false
review_dataset: ci_logs
report_artifact_format: key_value_v1
report_script: scripts/flutter_web_timing_report.sh
report_test: scripts/flutter_web_timing_report_test.sh
report_required_fields: report_version,source_marker,sample_count,duration_seconds_values,latest_duration_seconds,min_duration_seconds,max_duration_seconds,branch_type,runner_context
ci_artifact_upload: true
artifact_name: flutter-web-timing-report
artifact_path: /tmp/manpasik_flutter_web_timing.env
artifact_upload_action: actions/upload-artifact@v4
report_generation_step: Flutter web timing report
report_upload_step: Upload Flutter web timing report
threshold_promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md
sample_source_artifact: flutter-web-timing-report
```

## Policy

Flutter web release gate timing is collected from CI logs by searching for:

```text
FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=
```

The value is emitted by `scripts/flutter_web_release_gate.sh` after a successful `flutter build web --no-pub` run.

At least 5 successful CI samples are required before relaxing the PR-required web release gate policy. Until then, timing thresholds are advisory only and must not automatically fail the build.

Threshold promotion rules are governed by `docs/ci/flutter-web-timing-threshold-promotion.md`.

## Report Artifact

Use `scripts/flutter_web_timing_report.sh` to convert a CI log into a key-value report artifact:

```bash
bash scripts/flutter_web_timing_report.sh \
  --log ci.log \
  --branch-type pull_request \
  --runner-context ubuntu-latest \
  --output flutter-web-timing.env
```

The report must include:

```text
report_version=1
source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
sample_count=<n>
duration_seconds_values=<comma-separated>
latest_duration_seconds=<last>
min_duration_seconds=<min>
max_duration_seconds=<max>
branch_type=<pull_request|release_branch>
runner_context=<runner-label>
```

## CI Artifact Upload

Successful `flutter-app` CI runs generate `/tmp/manpasik_flutter_web_timing.env` in the `Flutter web timing report` step.
The following `Upload Flutter web timing report` step uploads that file with `actions/upload-artifact@v4` as `flutter-web-timing-report`.

## Review Procedure

1. Collect duration markers from successful pull request and release branch CI runs.
2. Keep the raw duration seconds with date, branch type, runner context, and the key-value report artifact.
3. Review median and worst-case duration before proposing `nightly_split: true`.
4. Keep JS web release risk blocking while `current_nightly_split: false`.

## Change Rules

- Do not set `blocking_threshold_seconds` to a number without a measured cost review.
- Do not change `current_nightly_split` without updating `docs/ci/flutter-web-release-gate-policy.md`.
- If timing collection moves from logs to artifacts, update `review_dataset` and this guard in the same change.
