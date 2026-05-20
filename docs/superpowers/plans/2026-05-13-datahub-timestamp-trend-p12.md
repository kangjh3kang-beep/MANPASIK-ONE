# DataHub Timestamp Trend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub summary trend가 값 정렬이 아니라 측정 timestamp 순서를 기준으로 상승/하강/안정을 계산하게 한다.

**Architecture:** P11에서 `getTrendData`가 timestamp 오름차순 `TrendDataPoint`를 반환하도록 정렬했으므로, summary trend 계산은 이 순서를 그대로 사용한다. min/max/average는 기존처럼 값 집계로 계산하되, trend만 시간 순서의 앞/뒤 절반 평균 비교로 분리한다.

**Tech Stack:** Flutter, Dart unit tests, ManPaSik REST DataHub repository.

---

### Task 1: REST summary trend contract

**Files:**
- Modify: `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`

- [ ] **Step 1: Write the failing test**

Add a history response where values decrease over time: older 130, middle 100, newer 70. `getBiomarkerSummary('glucose').trend` must be `falling`.

- [ ] **Step 2: Run test to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: FAIL because current `_computeTrend` receives sorted values and reports rising/stable incorrectly for chronological declines.

- [ ] **Step 3: Implement timestamp-based trend**

Use the timestamp-sorted points from `getTrendData` to create ordered values for trend. In `getAllBiomarkerSummaries`, keep points per biomarker, sort each list by timestamp, and compute latest/trend from that ordered list while preserving min/max from sorted values.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/datahub-timestamp-trend-p12.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run focused tests**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart`

- [ ] **Step 2: Run targeted analyzer**

Run:
`cd frontend/flutter-app && flutter analyze --no-pub --no-fatal-infos lib/features/data_hub/data/data_hub_repository_rest.dart lib/features/data_hub/domain/data_hub_repository.dart test/features/data_hub/data/data_hub_repository_rest_test.dart`

- [ ] **Step 3: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 4: Run diff hygiene checks**

Run:
`git diff --check -- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart CHANGELOG.md CONTEXT.md`
