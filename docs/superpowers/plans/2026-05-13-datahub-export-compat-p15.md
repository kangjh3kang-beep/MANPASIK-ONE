# DataHub Export Compatibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub export result가 REST response의 snake_case와 camelCase 파일/카운트 필드를 모두 보존하게 한다.

**Architecture:** P14에서 추가한 `_intField`와 P9의 `_stringField` fallback 패턴을 export에도 적용한다. `file_path`, `filePath`, `fhir_json`, `fhirJson` 순서로 파일 경로를 해석하고 `record_count`, `recordCount`를 모두 record count로 처리한다.

**Tech Stack:** Flutter, Dart unit tests, ManPaSik REST DataHub repository.

---

### Task 1: Export response compatibility

**Files:**
- Modify: `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
- Modify: `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`

- [ ] **Step 1: Write the failing test**

Add a local HTTP server test for `/api/v1/health-records/export/fhir` returning:

```json
{"filePath": "/tmp/export.fhir.json", "recordCount": 3}
```

`exportData(format: ExportFormat.json)` must return the same file path and record count.

- [ ] **Step 2: Run test to verify RED**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: FAIL because `exportData` currently reads only `file_path`, `fhir_json`, and `record_count`.

- [ ] **Step 3: Implement fallback**

Use `_stringField` and `_intField` fallbacks to read `filePath`, `fhirJson`, and `recordCount` without changing the public `ExportResult` API.

- [ ] **Step 4: Run test to verify GREEN**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`

Expected: PASS.

### Task 2: Quality and logs

**Files:**
- Create: `docs/audit/datahub-export-compat-p15.md`
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
