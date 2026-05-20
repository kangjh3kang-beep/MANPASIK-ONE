# Flutter Web Release Review Checklist

```yaml
release_review_checklist_version: 1
review_scope: flutter_web_timing_release_review
review_order: policy_gates,artifact_collection,offline_aggregation,fixture_validation,threshold_review_check,decision_record
required_before_threshold_change: true
no_phi_review_required: true
checklist_guard: scripts/ci_web_release_review_checklist_policy_test.sh
synthetic_review_example: docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review
synthetic_review_audit: docs/audit/flutter-web-timing-synthetic-review-example.md
```

## Purpose

Use this checklist before changing Flutter web timing thresholds, PR web release gate policy, or `nightly_split`.
It links the policy gates, artifact handling, offline aggregation, concrete fixture validation, and audit decision record in review order.

## Referenced Policy And Tools

- `docs/ci/flutter-web-gate-timing-collection.md`
- `docs/ci/flutter-web-timing-threshold-promotion.md`
- `docs/ci/flutter-web-timing-review-fixtures.md`
- `docs/audit/templates/flutter-web-timing-review-template.md`
- `scripts/flutter_web_timing_sample_aggregate.sh`
- `scripts/flutter_web_timing_review_fixture_validate.sh`
- `scripts/flutter_web_timing_threshold_change_review_check.sh`
- `scripts/ci_web_timing_synthetic_review_example_test.sh`
- `docs/audit/flutter-web-timing-synthetic-review-example.md`

## Review Order

1. Run policy gates:

```bash
bash scripts/ci_web_gate_timing_collection_policy_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_web_timing_review_fixture_policy_test.sh
```

2. Export CI artifacts:

Download at least 5 successful GitHub Actions `flutter-web-timing-report` artifacts into:

```text
docs/ci/fixtures/flutter-web-timing/<review-id>/samples
```

3. Generate aggregate timing data:

```bash
bash scripts/flutter_web_timing_sample_aggregate.sh --input-dir docs/ci/fixtures/flutter-web-timing/<review-id>/samples --output docs/ci/fixtures/flutter-web-timing/<review-id>/aggregate.env
```

4. Validate the concrete review fixture:

```bash
bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id>
```

5. Complete the audit record with `docs/audit/templates/flutter-web-timing-review-template.md`.

6. Check the audit record against the fixture:

```bash
bash scripts/flutter_web_timing_threshold_change_review_check.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id> --audit-file docs/audit/<review-doc>.md
```

## Synthetic Dry Run

Use the synthetic example to verify the checklist path without production artifacts:

```bash
bash scripts/ci_web_timing_synthetic_review_example_test.sh
```

The synthetic fixture lives at:

```text
docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review
```

The matching audit example is:

```text
docs/audit/flutter-web-timing-synthetic-review-example.md
```

## Decision Rules

- Do not change timing thresholds before the checklist commands pass.
- Do not relax the PR web release gate unless the audit record decision is `approved`.
- Do not include PHI, user identifiers, device identifiers, access tokens, or raw medical measurements in fixtures or audit records.
- Keep the audit record in `docs/audit` with the review fixture path and rollback notes.
