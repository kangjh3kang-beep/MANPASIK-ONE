# Flutter Web Timing Report Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** CI logs의 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 표준 timing report artifact로 변환하는 스크립트와 정책을 추가한다.

**Architecture:** `scripts/flutter_web_timing_report.sh`는 build log 파일에서 duration marker를 추출해 key-value report를 생성한다. `scripts/flutter_web_timing_report_test.sh`는 fixture log로 parser/report schema를 검증한다. Timing collection policy guard는 report script와 schema marker를 문서에 요구한다.

**Tech Stack:** Bash, Markdown CI policy, text report artifact.

---

### Task 1: Timing report script

**Files:**
- Create: `scripts/flutter_web_timing_report_test.sh`
- Create: `scripts/flutter_web_timing_report.sh`

- [x] **Step 1: Write failing report script test**

Create `scripts/flutter_web_timing_report_test.sh`.

The test must:
- Create a fixture log with:
  - `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=11`
  - `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=17`
- Run:

```bash
bash scripts/flutter_web_timing_report.sh \
  --log <fixture> \
  --branch-type pull_request \
  --runner-context ubuntu-latest \
  --output <report>
```

- Require the report to contain:
  - `report_version=1`
  - `source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS`
  - `sample_count=2`
  - `duration_seconds_values=11,17`
  - `latest_duration_seconds=17`
  - `min_duration_seconds=11`
  - `max_duration_seconds=17`
  - `branch_type=pull_request`
  - `runner_context=ubuntu-latest`
- Also verify a log without duration marker fails.

- [x] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_web_timing_report_test.sh`

Expected: FAIL because `scripts/flutter_web_timing_report.sh` does not exist.

- [x] **Step 3: Implement report script**

Create `scripts/flutter_web_timing_report.sh`.

The script must:
- Support `--log`, `--branch-type`, `--runner-context`, `--output`.
- Extract numeric `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` values.
- Fail when no marker exists.
- Write key-value report format:

```text
report_version=1
source_marker=FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS
sample_count=<n>
duration_seconds_values=<comma-separated>
latest_duration_seconds=<last>
min_duration_seconds=<min>
max_duration_seconds=<max>
branch_type=<branch-type>
runner_context=<runner-context>
```

- [x] **Step 4: Run report script test to verify GREEN**

Run:
`bash scripts/flutter_web_timing_report_test.sh`

Expected: PASS.

### Task 2: Policy integration

**Files:**
- Modify: `scripts/ci_web_gate_timing_collection_policy_test.sh`
- Modify: `docs/ci/flutter-web-gate-timing-collection.md`

- [x] **Step 1: Write failing policy expectations**

Update `scripts/ci_web_gate_timing_collection_policy_test.sh` to require:
- `report_artifact_format: key_value_v1`
- `report_script: scripts/flutter_web_timing_report.sh`
- `report_test: scripts/flutter_web_timing_report_test.sh`
- `report_required_fields: report_version,source_marker,sample_count,duration_seconds_values,latest_duration_seconds,min_duration_seconds,max_duration_seconds,branch_type,runner_context`

- [x] **Step 2: Run policy guard to verify RED**

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: FAIL because the timing collection policy does not include report markers yet.

- [x] **Step 3: Update timing collection policy document**

Add report artifact format markers and a short report schema section to `docs/ci/flutter-web-gate-timing-collection.md`.

- [x] **Step 4: Run policy guard to verify GREEN**

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Expected: PASS.

### Task 3: Quality and logs

**Files:**
- Create: `docs/audit/flutter-web-timing-report-p29.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Run targeted gates**

Run:
`bash scripts/flutter_web_timing_report_test.sh`

Run:
`bash scripts/ci_web_gate_timing_collection_policy_test.sh`

Run:
`bash scripts/ci_gate_matrix_policy_test.sh`

Run:
`bash scripts/ci_web_gate_policy_test.sh`

Run:
`bash -n scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_web_gate_policy_test.sh`

- [x] **Step 2: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [x] **Step 3: Run diff hygiene checks**

Run:
`git diff --check -- scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh docs/ci/flutter-web-gate-timing-collection.md docs/superpowers/plans/2026-05-13-flutter-web-timing-report-p29.md docs/audit/flutter-web-timing-report-p29.md CHANGELOG.md CONTEXT.md`
