# Flutter Web Timing Collection Policy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P27에서 추가한 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` marker의 수집 기준과 gate 완화 검토 조건을 문서/guard/CI matrix로 고정한다.

**Architecture:** `docs/ci/flutter-web-gate-timing-collection.md`는 timing marker 수집 단위, 최소 sample 수, threshold 상태, nightly split 전환 조건을 기록한다. `scripts/ci_web_gate_timing_collection_policy_test.sh`는 이 문서와 web release policy, CI workflow, CI gate matrix가 일치하는지 검증한다.

**Tech Stack:** Bash, Markdown CI policy, GitHub Actions workflow.

---

### Task 1: Timing collection guard

**Files:**
- Create: `scripts/ci_web_gate_timing_collection_policy_test.sh`
- Create: `docs/ci/flutter-web-gate-timing-collection.md`
- Modify: `.github/workflows/ci.yml`

- [x] **Step 1: Write failing timing collection guard**

Create `scripts/ci_web_gate_timing_collection_policy_test.sh`.

The guard must require:
- `docs/ci/flutter-web-gate-timing-collection.md`
- `docs/ci/flutter-web-release-gate-policy.md`
- `docs/ci/ci-gate-matrix.md`
- markers in timing collection doc:
  - `timing_collection_version: 1`
  - `source_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS`
  - `log_query_pattern: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=`
  - `minimum_samples_before_relaxing: 5`
  - `threshold_policy_status: advisory`
  - `blocking_threshold_seconds: unset`
  - `nightly_split_requires: measured_cost_review`
  - `current_nightly_split: false`
  - `review_dataset: ci_logs`
- web release policy still contains:
  - `timing_capture: true`
  - `timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS`
  - `nightly_split: false`
- CI workflow contains:
  - `Flutter web timing collection policy`
  - `scripts/ci_web_gate_timing_collection_policy_test.sh`

- [x] **Step 2: Run guard to verify RED**

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: FAIL because timing collection policy doc and CI step do not exist yet.

- [x] **Step 3: Create document and CI governance step**

Create `docs/ci/flutter-web-gate-timing-collection.md` with the required markers and a short procedure:
- collect `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` from CI logs.
- keep at least 5 successful samples before relaxing PR gate policy.
- threshold is advisory only until enough samples exist.

Add to `.github/workflows/ci.yml` in `ssot-governance`:

```yaml
      - name: Flutter web timing collection policy
        run: bash scripts/ci_web_gate_timing_collection_policy_test.sh
```

- [x] **Step 4: Run guard to verify GREEN**

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: PASS.

### Task 2: Matrix integration

**Files:**
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
- Modify: `docs/ci/ci-gate-matrix.md`

- [x] **Step 1: Write failing matrix expectations**

Update `scripts/ci_gate_matrix_policy_test.sh` to require:
- `Flutter web timing collection policy`
- `docs/ci/flutter-web-gate-timing-collection.md`
- `scripts/ci_web_gate_timing_collection_policy_test.sh`
- `minimum_samples_before_relaxing: 5`
- `threshold_policy_status: advisory`

- [x] **Step 2: Run matrix guard to verify RED**

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: FAIL until matrix document references the timing collection policy.

- [x] **Step 3: Update matrix document**

Add a matrix row for `Flutter web timing collection policy`.

- [x] **Step 4: Run matrix guard to verify GREEN**

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Expected: PASS.

### Task 3: Quality and logs

**Files:**
- Create: `docs/audit/flutter-web-timing-collection-policy-p28.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash scripts/ci_flutter_web_gate_workflow_test.sh`

Run:
`bash -n scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- .github/workflows/ci.yml scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh docs/ci/flutter-web-gate-timing-collection.md docs/ci/ci-gate-matrix.md docs/superpowers/plans/2026-05-13-flutter-web-timing-collection-policy-p28.md docs/audit/flutter-web-timing-collection-policy-p28.md CHANGELOG.md CONTEXT.md`
