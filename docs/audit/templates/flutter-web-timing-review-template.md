# Flutter Web Timing Review Template

```yaml
review_template_version: 1
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
aggregate_source: aggregate.env
required_decision_fields: decision,median_seconds,p95_seconds,worst_case_seconds,sample_count,branch_types,runner_contexts
```

## Review Metadata

- Review ID:
- Review date:
- Reviewer:
- Fixture path:
- Source artifact name: `flutter-web-timing-report`
- No-PHI attestation: `true`

## Fixture

- Manifest: `manifest.env`
- Samples: `samples/sample-01.env` through `samples/sample-05.env`
- Aggregate: `aggregate.env`
- Aggregator command:

```bash
bash scripts/flutter_web_timing_sample_aggregate.sh \
  --input-dir docs/ci/fixtures/flutter-web-timing/<review-id>/samples \
  --output docs/ci/fixtures/flutter-web-timing/<review-id>/aggregate.env
```

## Aggregate Results

- `review_id`:
- `decision`:
- `sample_count`:
- `duration_seconds_values`:
- `median_seconds`:
- `p95_seconds`:
- `worst_case_seconds`:
- `branch_types`:
- `runner_contexts`:
- `validator_command`:

## Machine Check Fields

```text
review_id:
decision:
sample_count:
median_seconds:
p95_seconds:
worst_case_seconds:
branch_types:
runner_contexts:
validator_command:
```

## Decision

- `decision`: `proposed | approved | rejected`
- Threshold change:
- PR web gate relaxation:
- Nightly split change:
- Rationale:

## Sign-Off

- Engineering reviewer:
- Release owner:
- Security/compliance reviewer:

## Rollback Notes

- Revert target:
- Expected restored policy:
- Follow-up issue:
