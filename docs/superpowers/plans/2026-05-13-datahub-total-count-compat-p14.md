# DataHub Total Count Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub REST repository가 measurement history response의 `total_count`와 legacy `totalCount`를 모두 총 측정 수로 해석하게 한다.

**Architecture:** P9-P12에서 history item field를 snake_case/camelCase 양쪽으로 처리한 방식과 맞춰 count field도 작은 helper로 통일한다. 기존 `total_count` 동작은 유지하고, `totalCount` fallback만 추가한다.

**Tech Stack:** Flutter, Dart unit tests, ManPaSik REST DataHub repository.

---

### Task 1: Count field compatibility

**Files:**
- Modify: `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`

- [ ] **Step 1: Write the failing test**

Add a test where the local history server returns:

```json
{"measurements": [], "totalCount": 7}
```

`getTotalMeasurementCount()` must return `7`.

- [ ] **Step 2: Run test to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: FAIL because `getTotalMeasurementCount` currently reads only `total_count`.

- [ ] **Step 3: Implement fallback**

Add `_intField(map, 'total_count', 'totalCount')` helper and use it in `getTotalMeasurementCount`.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/datahub-total-count-compat-p14.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run focused tests**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart`

- [ ] **Step 2: Run targeted analyzer**

Run:
`cd frontend/flutter-app && flutter analyze --no-pub --no-fatal-infos lib/features/data_hub/data/data_hub_repository_rest.dart test/features/data_hub/data/data_hub_repository_rest_test.dart`

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
