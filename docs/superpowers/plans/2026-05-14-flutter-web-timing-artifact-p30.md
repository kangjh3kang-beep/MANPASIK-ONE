# Flutter Web Timing Artifact P30 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter web release gate가 성공한 CI 실행에서 timing report를 key-value 아티팩트로 생성하고 업로드하도록 고정한다.

**Architecture:** 기존 `scripts/flutter_web_release_gate.sh`가 출력하는 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` 마커를 `scripts/flutter_web_timing_report.sh`로 변환한다. GitHub Actions `flutter-app` job에 report 생성 단계와 `actions/upload-artifact@v4` 업로드 단계를 추가하고, 별도 워크플로우 가드 및 정책 문서 가드가 이를 검증한다.

**Tech Stack:** Bash, GitHub Actions, Flutter CI gate policy docs.

---

## File Structure

- Create: `scripts/ci_flutter_web_timing_artifact_workflow_test.sh`
  - `.github/workflows/ci.yml`에 timing report 생성 및 업로드 단계가 연결됐는지 검증한다.
- Modify: `.github/workflows/ci.yml`
  - `Flutter web release gate` 뒤에 `Flutter web timing report`와 `Upload Flutter web timing report` 단계를 추가한다.
- Modify: `scripts/ci_web_gate_timing_collection_policy_test.sh`
  - timing collection 정책 문서와 workflow가 CI 아티팩트 업로드 계약을 포함하는지 검증한다.
- Modify: `docs/ci/flutter-web-gate-timing-collection.md`
  - 아티팩트 이름, 경로, 업로드 액션, 생성/업로드 step 이름을 정책 마커로 문서화한다.
- Modify: `scripts/ci_gate_matrix_policy_test.sh`
  - CI gate matrix가 timing artifact 워크플로우 가드를 참조하는지 검증한다.
- Modify: `docs/ci/ci-gate-matrix.md`
  - Flutter web timing artifact gate를 matrix에 추가한다.
- Create: `docs/audit/flutter-web-timing-artifact-p30.md`
  - 구현 결과와 검증 증거를 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - 3-AI 협업 기록 프로토콜에 따라 변경 내용을 최신 상태로 반영한다.

## Task 1: 워크플로우 RED 가드

- [ ] **Step 1: Write the failing workflow guard**

Create `scripts/ci_flutter_web_timing_artifact_workflow_test.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW="$ROOT/.github/workflows/ci.yml"

fail() {
  echo "CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_FAIL $1" >&2
  exit 1
}

[[ -f "$WORKFLOW" ]] || fail "missing workflow"

required_markers=(
  "Flutter web release gate"
  "Flutter web timing report"
  "scripts/flutter_web_timing_report.sh"
  "--log /tmp/manpasik_flutter_web_build.log"
  "--output /tmp/manpasik_flutter_web_timing.env"
  "github.event_name"
  "pull_request"
  "runner.os"
  "runner.arch"
  "Upload Flutter web timing report"
  "actions/upload-artifact@v4"
  "name: flutter-web-timing-report"
  "path: /tmp/manpasik_flutter_web_timing.env"
)

for marker in "${required_markers[@]}"; do
  if ! grep -Fq -- "$marker" "$WORKFLOW"; then
    fail "workflow missing marker: $marker"
  fi
done

echo "CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_PASS"
```

- [ ] **Step 2: Run test to verify it fails**

Run: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`

Expected: `CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_FAIL workflow missing marker: Flutter web timing report`

## Task 2: GitHub Actions 연결

- [ ] **Step 1: Implement report and upload steps**

Modify `.github/workflows/ci.yml` after the existing `Flutter web release gate` step:

```yaml
      - name: Flutter web timing report
        run: |
          branch_type="release_branch"
          if [ "${{ github.event_name }}" = "pull_request" ]; then
            branch_type="pull_request"
          fi
          bash ../../scripts/flutter_web_timing_report.sh \
            --log /tmp/manpasik_flutter_web_build.log \
            --branch-type "$branch_type" \
            --runner-context "${{ runner.os }}-${{ runner.arch }}" \
            --output /tmp/manpasik_flutter_web_timing.env

      - name: Upload Flutter web timing report
        uses: actions/upload-artifact@v4
        with:
          name: flutter-web-timing-report
          path: /tmp/manpasik_flutter_web_timing.env
```

- [ ] **Step 2: Run workflow guard to verify it passes**

Run: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`

Expected: `CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_PASS`

## Task 3: 정책 문서 및 정책 가드 확장

- [ ] **Step 1: Add policy guard markers**

Extend `scripts/ci_web_gate_timing_collection_policy_test.sh` so `required_policy_markers` includes:

```bash
"ci_artifact_upload: true"
"artifact_name: flutter-web-timing-report"
"artifact_path: /tmp/manpasik_flutter_web_timing.env"
"artifact_upload_action: actions/upload-artifact@v4"
"report_generation_step: Flutter web timing report"
"report_upload_step: Upload Flutter web timing report"
```

Extend `workflow_markers` with:

```bash
"Flutter web timing report"
"Upload Flutter web timing report"
"actions/upload-artifact@v4"
"flutter-web-timing-report"
"/tmp/manpasik_flutter_web_timing.env"
```

- [ ] **Step 2: Run timing collection policy guard to verify it fails**

Run: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: failure for the first missing policy marker, `ci_artifact_upload: true`.

- [ ] **Step 3: Update policy document**

Modify `docs/ci/flutter-web-gate-timing-collection.md` marker block with:

```yaml
ci_artifact_upload: true
artifact_name: flutter-web-timing-report
artifact_path: /tmp/manpasik_flutter_web_timing.env
artifact_upload_action: actions/upload-artifact@v4
report_generation_step: Flutter web timing report
report_upload_step: Upload Flutter web timing report
```

Add a `## CI Artifact Upload` section that states successful `flutter-app` CI runs upload `/tmp/manpasik_flutter_web_timing.env` as `flutter-web-timing-report`.

- [ ] **Step 4: Run timing collection policy guard to verify it passes**

Run: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: `CI_WEB_GATE_TIMING_COLLECTION_POLICY_PASS`

## Task 4: CI gate matrix 연결

- [ ] **Step 1: Extend matrix guard**

Modify `scripts/ci_gate_matrix_policy_test.sh` so `required_markers` includes:

```bash
"Flutter web timing artifact"
"artifact_name: flutter-web-timing-report"
"artifact_path: /tmp/manpasik_flutter_web_timing.env"
"actions/upload-artifact@v4"
```

Add `scripts/ci_flutter_web_timing_artifact_workflow_test.sh` to `required_files`.

- [ ] **Step 2: Run matrix guard to verify it fails**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: failure for the first missing matrix marker, `Flutter web timing artifact`.

- [ ] **Step 3: Update matrix document**

Modify `docs/ci/ci-gate-matrix.md` marker block with:

```yaml
artifact_name: flutter-web-timing-report
artifact_path: /tmp/manpasik_flutter_web_timing.env
```

Add a matrix row:

```markdown
| Flutter web timing artifact | `flutter-app` / `Upload Flutter web timing report` | web release gate timing artifact upload via `actions/upload-artifact@v4` | `blocking: true` | `execution_mode: pull_request_and_release_branch` | `docs/ci/flutter-web-gate-timing-collection.md` | `scripts/ci_flutter_web_timing_artifact_workflow_test.sh` |
```

- [ ] **Step 4: Run matrix guard to verify it passes**

Run: `bash scripts/ci_gate_matrix_policy_test.sh`

Expected: `CI_GATE_MATRIX_POLICY_PASS`

## Task 5: 최종 검증과 기록

- [ ] **Step 1: Run targeted gates**

Run:

```bash
bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh
bash scripts/ci_web_gate_timing_collection_policy_test.sh
bash scripts/ci_gate_matrix_policy_test.sh
bash scripts/flutter_web_timing_report_test.sh
bash scripts/ci_web_gate_policy_test.sh
bash scripts/ci_flutter_web_gate_workflow_test.sh
```

Expected: all commands exit 0 with their `*_PASS` markers.

- [ ] **Step 2: Run syntax and project guards**

Run:

```bash
bash -n scripts/ci_flutter_web_timing_artifact_workflow_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Write audit and project logs**

Create `docs/audit/flutter-web-timing-artifact-p30.md` with scope, changed files, verification commands, and residual notes.

Update `CHANGELOG.md` with a top entry for P30 and update `CONTEXT.md` current-state notes for Flutter web CI timing artifact upload.
