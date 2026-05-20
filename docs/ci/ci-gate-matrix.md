# CI Gate Matrix

```yaml
matrix_version: 1
release_target: js_web
wasm_release_target: false
timing_capture: true
timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
minimum_samples_before_relaxing: 5
threshold_policy_status: advisory
artifact_name: flutter-web-timing-report
artifact_path: /tmp/manpasik_flutter_web_timing.env
promotion_policy_version: 1
automated_gate_relaxation: false
review_fixture_policy_version: 1
fixture_root: docs/ci/fixtures/flutter-web-timing
release_review_checklist_version: 1
review_scope: flutter_web_timing_release_review
governance_index_version: 1
governance_scope: flutter_web_timing_ci
```

## Matrix

| Gate | Workflow job/step | Scope | Blocking | Execution mode | Policy doc | Guard script |
|------|-------------------|-------|----------|----------------|------------|--------------|
| SSOT constants | `ssot-governance` / `Validate SSOT constants` | 상수/문서 일치성 | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `AGENTS.md` | `scripts/validate_ssot_constants.py` |
| Security release gate | `ssot-governance` / `Security release gate` | 릴리스 보안 baseline | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `docs/audit/security-release-gates.md` | `scripts/security_release_gate.sh` |
| Assay evidence gate | `ssot-governance` / `Assay evidence gate` | assay evidence contract | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `docs/audit/assay-evidence-factory-p1.md` | `scripts/assay_evidence_gate.sh` |
| Flutter evidence policy | `ssot-governance` / `Flutter evidence UI gate policy` | evidence gate policy drift | `blocking: true` | `execution_mode: category_shards` | `docs/ci/flutter-evidence-ui-gate-policy.md` | `scripts/flutter_evidence_ui_gate_policy_test.sh` |
| Flutter evidence workflow | `flutter-app` / `Flutter evidence UI gate: contract` | evidence contract shard | `blocking: true` | `execution_mode: category_shards` | `docs/ci/flutter-evidence-ui-gate-policy.md` | `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh` |
| Flutter evidence workflow | `flutter-app` / `Flutter evidence UI gate: ui` | evidence UI shard | `blocking: true` | `execution_mode: category_shards` | `docs/ci/flutter-evidence-ui-gate-policy.md` | `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh` |
| Flutter evidence workflow | `flutter-app` / `Flutter evidence UI gate: data` | evidence mapper/data shard | `blocking: true` | `execution_mode: category_shards` | `docs/ci/flutter-evidence-ui-gate-policy.md` | `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh` |
| Flutter web release gate | `flutter-app` / `Flutter web release gate` | JS web release artifact | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `docs/ci/flutter-web-release-gate-policy.md` | `scripts/ci_web_gate_policy_test.sh` |
| Flutter web timing collection policy | `ssot-governance` / `Flutter web timing collection policy` | web release gate cost review samples | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-gate-timing-collection.md` | `scripts/ci_web_gate_timing_collection_policy_test.sh` |
| Flutter web timing artifact | `flutter-app` / `Upload Flutter web timing report` | web release gate timing artifact upload via `actions/upload-artifact@v4` | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `docs/ci/flutter-web-gate-timing-collection.md` | `scripts/ci_flutter_web_timing_artifact_workflow_test.sh` |
| Flutter web timing threshold promotion policy | `ssot-governance` / `Flutter web timing threshold promotion policy` | measured cost review rules before web gate relaxation | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-threshold-promotion.md` | `scripts/ci_web_timing_threshold_promotion_policy_test.sh` |
| Flutter web timing review fixture policy | `ssot-governance` / `Flutter web timing review fixture policy` | reproducible timing review fixtures and audit template | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-review-fixtures.md` | `scripts/ci_web_timing_review_fixture_policy_test.sh` |
| Flutter web release review checklist policy | `ssot-governance` / `Flutter web release review checklist policy` | release-review ordered command checklist for timing gate changes | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-release-review-checklist.md` | `scripts/ci_web_release_review_checklist_policy_test.sh` |
| Flutter web timing governance index policy | `ssot-governance` / `Flutter web timing governance index policy` | canonical index for Flutter web timing CI governance | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-governance-index.md` | `scripts/ci_web_timing_governance_index_policy_test.sh` |
| Wasm readiness policy | `ssot-governance` / policy guard only | Wasm non-release readiness | `blocking: false` | `execution_mode: documentation_guard` | `docs/ci/flutter-wasm-readiness.md` | `scripts/ci_wasm_readiness_policy_test.sh` |

## Change Rules

- New blocking gates need a policy document, a guard script, a workflow connection, and CHANGELOG/CONTEXT entries.
- New non-blocking readiness gates still need a policy document and guard script if they influence release decisions.
- `blocking: true` gates must be runnable from a fresh checkout without hidden local state beyond normal project dependencies.
- `blocking: false` gates must clearly state promotion criteria before becoming release blocking.
- Evidence gate CI execution remains category-sharded while `docs/ci/flutter-evidence-ui-gate-policy.md` declares `ci_execution_mode: category_shards`.
- Flutter web release gate remains PR/release branch blocking while `docs/ci/flutter-web-release-gate-policy.md` declares `release_target: js_web` and `nightly_split: false`.
- Flutter web release gate timing policy remains advisory until `docs/ci/flutter-web-gate-timing-collection.md` records enough samples for a measured cost review.
- Flutter web timing artifact upload remains tied to successful `flutter-app` web release gate runs and stores `flutter-web-timing-report`.
- Flutter web timing threshold promotion cannot be automated while `automated_gate_relaxation: false`.
- Flutter web timing reviews use fixture root `docs/ci/fixtures/flutter-web-timing` and the audit template before policy changes.
- Flutter web timing reviews must use `docs/audit/templates/flutter-web-timing-review-template.md`.
- Flutter web timing release reviews must follow `docs/ci/flutter-web-release-review-checklist.md`.
- Flutter web timing governance starts at `docs/ci/flutter-web-timing-governance-index.md`.
