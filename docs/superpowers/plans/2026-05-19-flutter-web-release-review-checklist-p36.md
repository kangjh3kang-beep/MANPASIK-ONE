# Flutter Web Release Review Checklist P36 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web timing gate, artifact, fixture, threshold review commands를 release-review 순서로 안내하는 checklist page와 CI guard를 추가한다.

**Architecture:** 신규 checklist 문서는 P30-P35에서 만든 policy/test/wrapper command를 reviewer 실행 순서대로 나열한다. 신규 guard script가 checklist marker, 필수 command, workflow step, matrix row를 검증해 문서와 CI gate가 drift되지 않게 한다.

**Tech Stack:** Bash, Markdown checklist, GitHub Actions policy guard.

---

## File Structure

- Create: `scripts/ci_web_release_review_checklist_policy_test.sh`
  - checklist 문서, workflow, matrix, 관련 script/doc 존재와 marker를 검증한다.
- Create: `docs/ci/flutter-web-release-review-checklist.md`
  - release review order와 command checklist를 문서화한다.
- Modify: `.github/workflows/ci.yml`
  - `ssot-governance` job에서 checklist policy guard를 실행한다.
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
  - CI matrix가 checklist policy gate를 참조하는지 검증한다.
- Modify: `docs/ci/ci-gate-matrix.md`
  - `Flutter web release review checklist policy` row와 marker를 추가한다.
- Create: `docs/audit/flutter-web-release-review-checklist-p36.md`
  - TDD 기록과 검증 결과를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - P36 상태와 다음 단계 지침을 최신 상태로 반영한다.

## Task 1: Checklist Policy RED Guard

- [ ] **Step 1: Write failing guard**

Create `scripts/ci_web_release_review_checklist_policy_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECKLIST="$ROOT/docs/ci/flutter-web-release-review-checklist.md"
MATRIX="$ROOT/docs/ci/ci-gate-matrix.md"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_FAIL $1" >&2
  exit 1
}

for file in "$CHECKLIST" "$MATRIX" "$WORKFLOW"; do
  [[ -f "$file" ]] || fail "missing file: ${file#$ROOT/}"
done

required_markers=(
  "release_review_checklist_version: 1"
  "review_scope: flutter_web_timing_release_review"
  "review_order: policy_gates,artifact_collection,offline_aggregation,fixture_validation,threshold_review_check,decision_record"
  "required_before_threshold_change: true"
  "no_phi_review_required: true"
  "checklist_guard: scripts/ci_web_release_review_checklist_policy_test.sh"
)

for marker in "${required_markers[@]}"; do
  grep -Fq "$marker" "$CHECKLIST" || fail "checklist missing marker: $marker"
done

required_commands=(
  "bash scripts/ci_web_gate_timing_collection_policy_test.sh"
  "bash scripts/ci_web_timing_threshold_promotion_policy_test.sh"
  "bash scripts/ci_web_timing_review_fixture_policy_test.sh"
  "bash scripts/flutter_web_timing_sample_aggregate.sh --input-dir docs/ci/fixtures/flutter-web-timing/<review-id>/samples --output docs/ci/fixtures/flutter-web-timing/<review-id>/aggregate.env"
  "bash scripts/flutter_web_timing_review_fixture_validate.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id>"
  "bash scripts/flutter_web_timing_threshold_change_review_check.sh --review-dir docs/ci/fixtures/flutter-web-timing/<review-id> --audit-file docs/audit/<review-doc>.md"
)

for command in "${required_commands[@]}"; do
  grep -Fq "$command" "$CHECKLIST" || fail "checklist missing command: $command"
done

required_files=(
  "docs/ci/flutter-web-gate-timing-collection.md"
  "docs/ci/flutter-web-timing-threshold-promotion.md"
  "docs/ci/flutter-web-timing-review-fixtures.md"
  "docs/audit/templates/flutter-web-timing-review-template.md"
  "scripts/flutter_web_timing_sample_aggregate.sh"
  "scripts/flutter_web_timing_review_fixture_validate.sh"
  "scripts/flutter_web_timing_threshold_change_review_check.sh"
)

for path in "${required_files[@]}"; do
  [[ -f "$ROOT/$path" ]] || fail "missing referenced file: $path"
  grep -Fq "$path" "$CHECKLIST" || fail "checklist does not reference file: $path"
done

workflow_markers=(
  "Flutter web release review checklist policy"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
)

for marker in "${workflow_markers[@]}"; do
  grep -Fq "$marker" "$WORKFLOW" || fail "workflow missing marker: $marker"
done

matrix_markers=(
  "Flutter web release review checklist policy"
  "docs/ci/flutter-web-release-review-checklist.md"
  "scripts/ci_web_release_review_checklist_policy_test.sh"
)

for marker in "${matrix_markers[@]}"; do
  grep -Fq "$marker" "$MATRIX" || fail "matrix missing marker: $marker"
done

echo "CI_WEB_RELEASE_REVIEW_CHECKLIST_POLICY_PASS"
```

- [ ] **Step 2: Run guard to verify it fails**

Run: `bash scripts/ci_web_release_review_checklist_policy_test.sh`

Expected: failure because `docs/ci/flutter-web-release-review-checklist.md` is missing.

## Task 2: Checklist Document and Workflow

- [ ] **Step 1: Add checklist document**

Create `docs/ci/flutter-web-release-review-checklist.md` with:

```yaml
release_review_checklist_version: 1
review_scope: flutter_web_timing_release_review
review_order: policy_gates,artifact_collection,offline_aggregation,fixture_validation,threshold_review_check,decision_record
required_before_threshold_change: true
no_phi_review_required: true
checklist_guard: scripts/ci_web_release_review_checklist_policy_test.sh
```

Include the exact commands from Task 1 in review order.

- [ ] **Step 2: Link workflow**

Add to `.github/workflows/ci.yml` in `ssot-governance`:

```yaml
      - name: Flutter web release review checklist policy
        run: bash scripts/ci_web_release_review_checklist_policy_test.sh
```

- [ ] **Step 3: Run guard to verify matrix failure remains**

Run: `bash scripts/ci_web_release_review_checklist_policy_test.sh`

Expected: failure for missing matrix marker.

## Task 3: Matrix Integration

- [ ] **Step 1: Extend matrix guard**

Modify `scripts/ci_gate_matrix_policy_test.sh` to require:

```bash
"Flutter web release review checklist policy"
"release_review_checklist_version: 1"
"review_scope: flutter_web_timing_release_review"
```

Add `docs/ci/flutter-web-release-review-checklist.md` and `scripts/ci_web_release_review_checklist_policy_test.sh` to `required_files`.

- [ ] **Step 2: Run matrix guard to verify it fails**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: failure for missing matrix marker.

- [ ] **Step 3: Update matrix doc**

Add marker block entries:

```yaml
release_review_checklist_version: 1
review_scope: flutter_web_timing_release_review
```

Add row:

```markdown
| Flutter web release review checklist policy | `ssot-governance` / `Flutter web release review checklist policy` | release-review ordered command checklist for timing gate changes | `blocking: true` | `execution_mode: documentation_guard` | `docs/ci/flutter-web-release-review-checklist.md` | `scripts/ci_web_release_review_checklist_policy_test.sh` |
```

- [ ] **Step 4: Run guards to verify pass**

Run:

```bash
bash scripts/ci_web_release_review_checklist_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: both commands emit `*_PASS`.

## Task 4: Final Verification and Logs

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_web_release_review_checklist_policy_test.sh
bash scripts/ci_web_timing_review_fixture_policy_test.sh
bash scripts/ci_web_timing_threshold_promotion_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
```

Expected: all commands exit 0.

- [ ] **Step 2: Run syntax and common gates**

Run:

```bash
bash -n scripts/ci_web_release_review_checklist_policy_test.sh scripts/ci_web_timing_review_fixture_policy_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_gate_matrix_policy_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-release-review-checklist-p36.md` and update `CHANGELOG.md` plus `CONTEXT.md` with P36 status, RED-GREEN record, changed files, verification commands, and next-stage guidance.
