# Flutter Web Timing Review Fixture P33 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 실제 `flutter-web-timing-report` 5-sample measured cost review를 같은 fixture 구조와 audit template으로 재현 가능하게 기록하는 규약을 추가한다.

**Architecture:** P32 aggregator output을 공식 review artifact로 삼고, exported sample files, manifest, aggregate output, audit review 문서의 위치를 정책 문서에 고정한다. 별도 CI guard가 fixture convention, audit template, promotion policy cross-reference, workflow, matrix 연결을 검증한다.

**Tech Stack:** Bash, Markdown policy docs, GitHub Actions policy gates.

---

## File Structure

- Create: `scripts/ci_web_timing_review_fixture_policy_test.sh`
  - review fixture convention 문서, fixture root README, audit template, promotion policy, workflow, matrix 연결을 검증한다.
- Create: `docs/ci/flutter-web-timing-review-fixtures.md`
  - real-sample review fixture root, manifest, sample file naming, aggregate output, no-PHI 규칙을 문서화한다.
- Create: `docs/ci/fixtures/flutter-web-timing/README.md`
  - 실제 review directory를 추가할 때 따라야 할 구조를 기록한다.
- Create: `docs/audit/templates/flutter-web-timing-review-template.md`
  - median/p95/worst-case, branch type, runner context, decision, sign-off를 기록하는 audit template을 제공한다.
- Modify: `docs/ci/flutter-web-timing-threshold-promotion.md`
  - fixture policy와 audit template을 threshold promotion policy에 교차 참조한다.
- Modify: `.github/workflows/ci.yml`
  - `ssot-governance` job에서 review fixture policy guard를 실행한다.
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
  - matrix가 review fixture policy gate를 참조하는지 검증한다.
- Modify: `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing review fixture policy` row와 marker를 추가한다.
- Create: `docs/audit/flutter-web-timing-review-fixture-p33.md`
  - TDD 기록과 검증 증거를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P33 진행 상황과 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Fixture Policy RED Guard

- [ ] **Step 1: Write failing guard**

Create `scripts/ci_web_timing_review_fixture_policy_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POLICY="$ROOT/docs/ci/flutter-web-timing-review-fixtures.md"
FIXTURE_README="$ROOT/docs/ci/fixtures/flutter-web-timing/README.md"
TEMPLATE="$ROOT/docs/audit/templates/flutter-web-timing-review-template.md"
PROMOTION_POLICY="$ROOT/docs/ci/flutter-web-timing-threshold-promotion.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$POLICY" "$FIXTURE_README" "$TEMPLATE" "$PROMOTION_POLICY" "$MATRIX" "$WORKFLOW"; do
  [[ -f "$file" ]] || fail "missing file: ${file#$ROOT/}"
done

required_policy_markers=(
  "review_fixture_policy_version: 1"
  "fixture_root: docs/ci/fixtures/flutter-web-timing"
  "fixture_manifest: manifest.env"
  "required_exported_artifacts: 5"
  "required_sample_files: sample-01.env,sample-02.env,sample-03.env,sample-04.env,sample-05.env"
  "aggregate_output_file: aggregate.env"
  "aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh"
  "audit_template: docs/audit/templates/flutter-web-timing-review-template.md"
  "promotion_policy: docs/ci/flutter-web-timing-threshold-promotion.md"
  "review_status_allowed: proposed,approved,rejected"
  "artifact_source_required: github_actions"
  "no_phi_allowed: true"
)

for marker in "${required_policy_markers[@]}"; do
  grep -Fq "$marker" "$POLICY" || fail "policy missing marker: $marker"
done

template_markers=(
  "review_template_version: 1"
  "review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md"
  "aggregate_source: aggregate.env"
  "required_decision_fields: decision,median_seconds,p95_seconds,worst_case_seconds,sample_count,branch_types,runner_contexts"
)

for marker in "${template_markers[@]}"; do
  grep -Fq "$marker" "$TEMPLATE" || fail "template missing marker: $marker"
done

cross_doc_markers=(
  "review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md"
  "review_fixture_root: docs/ci/fixtures/flutter-web-timing"
  "review_audit_template: docs/audit/templates/flutter-web-timing-review-template.md"
)

for marker in "${cross_doc_markers[@]}"; do
  grep -Fq "$marker" "$PROMOTION_POLICY" || fail "promotion policy missing marker: $marker"
done

workflow_markers=(
  "Flutter web timing review fixture policy"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  grep -Fq "$marker" "$WORKFLOW" || fail "workflow missing marker: $marker"
done

matrix_markers=(
  "Flutter web timing review fixture policy"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "scripts/ci_web_timing_review_fixture_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  grep -Fq "$marker" "$MATRIX" || fail "matrix missing marker: $marker"
done

echo "CI_WEB_TIMING_REVIEW_FIXTURE_POLICY_PASS"
```

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: failure because `docs/ci/flutter-web-timing-review-fixtures.md` is missing.

## Task 2: Fixture Convention and Template

- [ ] **Step 1: Add fixture convention docs**

Create `docs/ci/flutter-web-timing-review-fixtures.md` with the marker block from Task 1 and instructions for:

```text
docs/ci/fixtures/flutter-web-timing/<review-id>/
  manifest.env
  samples/sample-01.env ... samples/sample-05.env
  aggregate.env
```

- [ ] **Step 2: Add fixture root README**

Create `docs/ci/fixtures/flutter-web-timing/README.md` explaining that real exported artifacts must not contain PHI and must come from GitHub Actions `flutter-web-timing-report`.

- [ ] **Step 3: Add audit template**

Create `docs/audit/templates/flutter-web-timing-review-template.md` with marker block, aggregate result fields, decision fields, sign-off, and rollback notes.

- [ ] **Step 4: Cross-link promotion policy and workflow**

Add to `docs/ci/flutter-web-timing-threshold-promotion.md`:

```yaml
review_fixture_policy: docs/ci/flutter-web-timing-review-fixtures.md
review_fixture_root: docs/ci/fixtures/flutter-web-timing
review_audit_template: docs/audit/templates/flutter-web-timing-review-template.md
```

Add to `.github/workflows/ci.yml`:

```yaml
      - name: Flutter web timing review fixture policy
        run: bash scripts/ci_web_timing_review_fixture_policy_test.sh
```

- [ ] **Step 5: Run guard to verify matrix failure remains**

Run: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`

Expected: failure for missing matrix marker.

## Task 3: Matrix Integration

- [ ] **Step 1: Extend matrix guard**

Modify `scripts/ci_gate_matrix_policy_test.sh` to require:

```bash
"Flutter web timing review fixture policy"
"review_fixture_policy_version: 1"
"fixture_root: docs/ci/fixtures/flutter-web-timing"
```

Add `docs/ci/flutter-web-timing-review-fixtures.md`, `docs/audit/templates/flutter-web-timing-review-template.md`, and `scripts/ci_web_timing_review_fixture_policy_test.sh` to `required_files`.

- [ ] **Step 2: Run matrix guard to verify it fails**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: failure for missing matrix marker.

- [ ] **Step 3: Update matrix doc**

Add marker block entries:

```yaml
review_fixture_policy_version: 1
fixture_root: docs/ci/fixtures/flutter-web-timing
```

Add row:

```markdown
| Flutter web timing review fixture policy | `ssot-governance` / `Flutter web timing review fixture policy` | reproducible timing review fixtures and audit template | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-review-fixtures.md` | `scripts/ci_web_timing_review_fixture_policy_test.sh` |
```

- [ ] **Step 4: Run guards to verify pass**

Run:

```bash
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: both commands emit `*_PASS`.

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/flutter_web_timing_sample_aggregate_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/flutter_web_timing_sample_aggregate.sh scripts/flutter_web_timing_sample_aggregate_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-review-fixture-p33.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P33 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
