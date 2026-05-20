# Flutter Web Timing Synthetic Review Example

review_template_version: 1
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
aggregate_source: aggregate.env
synthetic_review: true
review_id: 2026-05-19-synthetic-review
decision: rejected
sample_count: 5
median_seconds: 17
p95_seconds: 34
worst_case_seconds: 34
branch_types: pull_request,release_branch
runner_contexts: Linux-X64
validator_command: bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review

## Scope

This is a synthetic, non-production timing review example used to verify the Flutter web release review checklist.
It must not be used as measured production evidence for changing thresholds or relaxing PR gates.

## Decision

- Decision: rejected
- Reason: synthetic example only; no production threshold or release policy change is allowed from this record.
- No-PHI attestation: true
