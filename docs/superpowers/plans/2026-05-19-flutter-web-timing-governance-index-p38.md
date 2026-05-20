# Flutter Web Timing Governance Index P38 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web timing CI governance 문서와 도구의 canonical entrypoint를 추가한다.

**Architecture:** 신규 governance index 문서는 release gate, timing collection, artifact upload, promotion policy, fixture convention, review checklist, synthetic example을 하나의 policy chain으로 연결한다. 신규 guard script가 index marker, 필수 문서/스크립트 링크, workflow step, matrix row를 검증한다.

**Tech Stack:** Bash, Markdown, GitHub Actions policy guard.

---

## File Structure

- Create: `scripts/ci_web_timing_governance_index_policy_test.sh`
  - governance index 문서, 참조 문서/스크립트, 기존 정책 가드 체인, 감사 기록 체인, workflow, matrix 연결을 검증한다.
- Create: `docs/ci/flutter-web-timing-governance-index.md`
  - Flutter web timing governance의 canonical entrypoint와 P27-P38 통합 추적 경로를 문서화한다.
- Modify: `.github/workflows/ci.yml`
  - `ssot-governance` job에서 governance index policy guard를 실행한다.
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
  - CI matrix가 governance index policy gate를 참조하는지 검증한다.
- Modify: `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing governance index policy` row와 marker를 추가한다.
- Create: `docs/audit/flutter-web-timing-governance-index-p38.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P38 상태와 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Governance Index RED Guard

- [ ] **Step 1: Write failing guard**

Create `scripts/ci_web_timing_governance_index_policy_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INDEX="$ROOT/docs/ci/flutter-web-timing-governance-index.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_TIMING_GOVERNANCE_INDEX_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$INDEX" "$MATRIX" "$WORKFLOW"; do
  [[ -f "$file" ]] || fail "missing file: ${file#$ROOT/}"
done

required_markers=(
  "governance_index_version: 1"
  "governance_scope: flutter_web_timing_ci"
  "entrypoint_status: canonical"
  "policy_chain: release_gate,timing_collection,artifact_upload,promotion_policy,review_fixture,release_review_checklist,synthetic_example"
  "index_guard: scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${required_markers[@]}"; do
  grep -Fq "$marker" "$INDEX" || fail "index missing marker: $marker"
done

required_refs=(
  "docs/ci/flutter-web-release-gate-policy.md"
  "docs/ci/flutter-web-gate-timing-collection.md"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "docs/ci/flutter-web-release-review-checklist.md"
  "docs/ci/ci-gate-matrix.md"
  "docs/audit/templates/flutter-web-timing-review-template.md"
  "scripts/flutter_web_release_gate.sh"
  "scripts/flutter_web_timing_report.sh"
  "scripts/flutter_web_timing_sample_aggregate.sh"
  "scripts/flutter_web_timing_review_fixture_validate.sh"
  "scripts/flutter_web_timing_threshold_change_review_check.sh"
  "scripts/ci_web_timing_synthetic_review_example_test.sh"
)

for ref in "${required_refs[@]}"; do
  [[ -f "$ROOT/$ref" ]] || fail "missing referenced file: $ref"
  grep -Fq "$ref" "$INDEX" || fail "index does not reference file: $ref"
done

workflow_markers=(
  "Flutter web timing governance index policy"
  "scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  grep -Fq "$marker" "$WORKFLOW" || fail "workflow missing marker: $marker"
done

matrix_markers=(
  "Flutter web timing governance index policy"
  "docs/ci/flutter-web-timing-governance-index.md"
  "scripts/ci_web_timing_governance_index_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  grep -Fq "$marker" "$MATRIX" || fail "matrix missing marker: $marker"
done

echo "CI_WEB_TIMING_GOVERNANCE_INDEX_POLICY_PASS"
```

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_timing_governance_index_policy_test.sh`

Expected: failure because `docs/ci/flutter-web-timing-governance-index.md` is missing.

## Task 2: Index Document and Workflow

- [ ] **Step 1: Add governance index**

Create `docs/ci/flutter-web-timing-governance-index.md` with:

```yaml
governance_index_version: 1
governance_scope: flutter_web_timing_ci
entrypoint_status: canonical
policy_chain: release_gate,timing_collection,artifact_upload,promotion_policy,review_fixture,release_review_checklist,synthetic_example
index_guard: scripts/ci_web_timing_governance_index_policy_test.sh
```

Include every referenced doc/script from Task 1.

- [ ] **Step 2: Link workflow**

Add to `.github/workflows/ci.yml` in `ssot-governance`:

```yaml
      - name: Flutter web timing governance index policy
        run: bash scripts/ci_web_timing_governance_index_policy_test.sh
```

- [ ] **Step 3: Run guard to verify matrix failure remains**

Run: `bash scripts/ci_web_timing_governance_index_policy_test.sh`

Expected: failure for missing matrix marker.

## Task 3: Matrix Integration

- [ ] **Step 1: Extend matrix guard**

Modify `scripts/ci_gate_matrix_policy_test.sh` to require:

```bash
"Flutter web timing governance index policy"
"governance_index_version: 1"
"governance_scope: flutter_web_timing_ci"
```

Add `docs/ci/flutter-web-timing-governance-index.md` and `scripts/ci_web_timing_governance_index_policy_test.sh` to `required_files`.

- [ ] **Step 2: Run matrix guard to verify it fails**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: failure for missing matrix marker.

- [ ] **Step 3: Update matrix doc**

Add marker block entries:

```yaml
governance_index_version: 1
governance_scope: flutter_web_timing_ci
```

Add row:

```markdown
| Flutter web timing governance index policy | `ssot-governance` / `Flutter web timing governance index policy` | canonical index for Flutter web timing CI governance | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-timing-governance-index.md` | `scripts/ci_web_timing_governance_index_policy_test.sh` |
```

- [ ] **Step 4: Run guards to verify pass**

Run:

```bash
bash scripts/ci_web_timing_governance_index_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: both commands emit `*_PASS`.

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_web_timing_governance_index_policy_test.sh
bash scripts/ci_web_release_review_checklist_policy_test.sh
bash scripts/ci_web_timing_synthetic_review_example_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

## Task 5: Existing Work Integration Supplement

- [ ] **Step 1: Extend guard for existing policy and audit chains**

Modify `scripts/ci_web_timing_governance_index_policy_test.sh` to require:

```bash
"guard_chain: release_gate_policy,release_gate_workflow,timing_collection_policy,timing_artifact_workflow,threshold_promotion_policy,review_fixture_policy,release_review_checklist_policy,synthetic_review_example,governance_index_policy,gate_matrix_policy"
```

Also require references to the existing guard scripts:

```bash
"scripts/ci_web_gate_policy_test.sh"
"scripts/ci_flutter_web_gate_workflow_test.sh"
"scripts/ci_web_gate_timing_collection_policy_test.sh"
"scripts/ci_flutter_web_timing_artifact_workflow_test.sh"
"scripts/ci_web_timing_threshold_promotion_policy_test.sh"
"scripts/ci_web_timing_review_fixture_policy_test.sh"
"scripts/ci_web_release_review_checklist_policy_test.sh"
"scripts/ci_gate_matrix_policy_test.sh"
```

And require P27-P38 audit records:

```bash
"docs/audit/flutter-web-release-gate-timing-p27.md"
"docs/audit/flutter-web-timing-collection-policy-p28.md"
"docs/audit/flutter-web-timing-report-p29.md"
"docs/audit/flutter-web-timing-artifact-p30.md"
"docs/audit/flutter-web-timing-threshold-promotion-p31.md"
"docs/audit/flutter-web-timing-sample-aggregator-p32.md"
"docs/audit/flutter-web-timing-review-fixture-p33.md"
"docs/audit/flutter-web-timing-review-validator-p34.md"
"docs/audit/flutter-web-timing-threshold-change-review-check-p35.md"
"docs/audit/flutter-web-release-review-checklist-p36.md"
"docs/audit/flutter-web-timing-synthetic-review-example-p37.md"
"docs/audit/flutter-web-timing-governance-index-p38.md"
```

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_timing_governance_index_policy_test.sh`

Expected: failure for missing `guard_chain`.

- [ ] **Step 3: Update canonical index**

Add `guard_chain`, `Guard Chain`, and `Audit Trail` sections to `docs/ci/flutter-web-timing-governance-index.md`.

- [ ] **Step 4: Run guard to verify pass**

Run: `bash scripts/ci_web_timing_governance_index_policy_test.sh`

Expected: `CI_WEB_TIMING_GOVERNANCE_INDEX_POLICY_PASS`.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/ci_web_timing_governance_index_policy_test.sh scripts/ci_web_release_review_checklist_policy_test.sh scripts/ci_web_timing_synthetic_review_example_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-governance-index-p38.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P38 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
