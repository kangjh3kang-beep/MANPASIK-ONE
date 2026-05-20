# Flutter Web Timing Threshold Promotion P31 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web release gate timing artifact sample이 쌓인 뒤에도 PR 필수 web gate 완화가 측정 근거 없이 변경되지 않도록 승격 조건을 정책과 CI 가드로 고정한다.

**Architecture:** P30의 `flutter-web-timing-report` artifact를 promotion review의 입력으로 정의한다. 신규 정책 문서와 guard script가 최소 sample 수, sample 구성, 통계 산출 항목, 수동 review 원칙, 자동 완화 금지를 검증하고, GitHub Actions 및 CI gate matrix가 이 guard를 실행하도록 연결한다.

**Tech Stack:** Bash, GitHub Actions, CI policy docs.

---

## File Structure

- Create: `scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - threshold promotion 정책 문서, timing collection 문서, CI workflow, gate matrix의 연결 상태를 검증한다.
- Create: `docs/ci/flutter-web-timing-threshold-promotion.md`
  - 최소 5개 성공 artifact sample 이후에도 필요한 측정/검토/문서 변경 조건을 명시한다.
- Modify: `.github/workflows/ci.yml`
  - `ssot-governance` job에서 threshold promotion policy guard를 실행한다.
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
  - matrix가 threshold promotion policy gate를 참조하는지 검증한다.
- Modify: `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing threshold promotion policy` row를 추가한다.
- Modify: `docs/ci/flutter-web-gate-timing-collection.md`
  - timing collection 정책에서 threshold promotion 정책 문서를 후속 decision gate로 참조한다.
- Create: `docs/audit/flutter-web-timing-threshold-promotion-p31.md`
  - TDD 기록과 검증 증거를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - 협업 기록 프로토콜에 따라 P31 상태를 반영한다.

## Task 1: Promotion Policy RED Guard

- [ ] **Step 1: Write the failing policy guard**

Create `scripts/ci_web_timing_threshold_promotion_policy_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-timing-threshold-promotion.md"
COLLECTION_POLICY="$ROOT/docs/ci/flutter-web-gate-timing-collection.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$POLICY" "$COLLECTION_POLICY" "$MATRIX" "$WORKFLOW"; do
  [[ -f "$file" ]] || fail "missing file: ${file#$ROOT/}"
done

required_policy_markers=(
  "promotion_policy_version: 1"
  "sample_source_artifact: flutter-web-timing-report"
  "minimum_successful_artifact_samples: 5"
  "required_branch_types: pull_request,release_branch"
  "required_runner_context_field: runner_context"
  "required_duration_field: latest_duration_seconds"
  "required_distribution_fields: min_duration_seconds,max_duration_seconds,duration_seconds_values"
  "required_statistics: median_seconds,p95_seconds,worst_case_seconds"
  "threshold_change_requires: measured_cost_review"
  "automated_gate_relaxation: false"
  "relaxation_allowed_before_minimum_samples: false"
  "nightly_split_change_requires: docs/ci/flutter-web-release-gate-policy.md"
  "review_record_required: docs/audit"
)

for marker in "${required_policy_markers[@]}"; do
  grep -Fq "$marker" "$POLICY" || fail "missing promotion marker: $marker"
done

cross_doc_markers=(
  "threshold_promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md"
  "sample_source_artifact: flutter-web-timing-report"
)

for marker in "${cross_doc_markers[@]}"; do
  grep -Fq "$marker" "$COLLECTION_POLICY" || fail "collection policy missing marker: $marker"
done

workflow_markers=(
  "Flutter web timing threshold promotion policy"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  grep -Fq "$marker" "$WORKFLOW" || fail "workflow missing marker: $marker"
done

matrix_markers=(
  "Flutter web timing threshold promotion policy"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "scripts/ci_web_timing_threshold_promotion_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  grep -Fq "$marker" "$MATRIX" || fail "matrix missing marker: $marker"
done

echo "CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_PASS"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`

Expected: failure because `docs/ci/flutter-web-timing-threshold-promotion.md` is missing.

## Task 2: Promotion Policy Document and CI Workflow

- [ ] **Step 1: Add policy document**

Create `docs/ci/flutter-web-timing-threshold-promotion.md` with marker block:

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
```

- [ ] **Step 2: Link collection policy**

Add these markers to `docs/ci/flutter-web-gate-timing-collection.md`:

```yaml
threshold_promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md
sample_source_artifact: flutter-web-timing-report
```

- [ ] **Step 3: Link workflow**

Add this step to `.github/workflows/ci.yml` in `ssot-governance`:

```yaml
      - name: Flutter web timing threshold promotion policy
        run: bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
```

- [ ] **Step 4: Run policy guard to verify matrix failure remains**

Run: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`

Expected: failure for missing matrix marker.

## Task 3: Matrix Integration

- [ ] **Step 1: Extend matrix guard**

Modify `scripts/ci_gate_matrix_policy_test.sh`:

```bash
"Flutter web timing threshold promotion policy"
"promotion_policy_version: 1"
"minimum_successful_artifact_samples: 5"
"automated_gate_relaxation: false"
```

Add `docs/ci/flutter-web-timing-threshold-promotion.md` and `scripts/ci_web_timing_threshold_promotion_policy_test.sh` to `required_files`.

- [ ] **Step 2: Run matrix guard to verify it fails**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: failure for missing matrix marker.

- [ ] **Step 3: Update matrix document**

Add marker block entries:

```yaml
promotion_policy_version: 1
automated_gate_relaxation: false
```

Add matrix row:

```markdown
| Flutter web timing threshold promotion policy | `ssot-governance` / `Flutter web timing threshold promotion policy` | measured cost review rules before web gate relaxation | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-threshold-promotion.md` | `scripts/ci_web_timing_threshold_promotion_policy_test.sh` |
```

- [ ] **Step 4: Run guards to verify they pass**

Run:

```bash
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: both commands emit `*_PASS`.

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_web_gate_timing_collection_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_flutter_web_timing_artifact_workflow_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-threshold-promotion-p31.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P31 status, changed files, RED-GREEN record, verification commands, and next-stage guidance.
