# Measurement Evidence History App Surface P6 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans and superpowers:test-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P5에서 저장/조회 proto까지 올라온 measurement evidence summary fields를 gateway REST history와 Flutter history domain/native/REST repository까지 보존한다.

**Architecture:** `MeasurementSummary` proto evidence fields는 이미 P5에서 확장되었다. P6에서는 HTTP JSON contract와 Flutter domain model에 같은 fields를 붙이고, REST snake_case 및 legacy camelCase response 모두 안전하게 decode한다.

**Tech Stack:** Go gateway route tests, Flutter repository tests, targeted Flutter analyze, existing SSOT/security/assay/proto gates.

---

## Stage Gate Rules

- RED 테스트를 먼저 작성하고 실패를 확인한다.
- history item의 legacy 기본값은 `unknown`, `false`, `[]`로 둔다.
- `research_only`는 `diagnosticReady=false`로 보존하며 UI/의료 판정 문구는 추가하지 않는다.
- P5 proto/generated field numbers는 변경하지 않는다.

## File Structure

- Modify: `backend/services/gateway/internal/handler/e2e_test.go`
  - `/measurements/history` REST response가 evidence fields를 보존하는지 검증한다.
- Modify: `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
  - `MeasurementHistoryItem`에 evidence fields를 추가한다.
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart`
  - native gRPC `MeasurementSummary` evidence fields를 도메인 history item으로 매핑한다.
- Modify: `frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart`
  - REST history item snake_case/camelCase evidence fields를 도메인으로 매핑한다.
- Modify: `frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart`
  - native history evidence mapping contract를 추가한다.
- Modify: `frontend/flutter-app/test/features/measurement/data/measurement_repository_rest_test.dart`
  - REST history evidence mapping contract를 추가한다.
- Create: `docs/audit/measurement-evidence-history-app-surface-p6.md`
- Modify: `CHANGELOG.md`, `CONTEXT.md`

## Task 1: Gateway REST History Contract

- [x] **Step 1: Write failing gateway test**

Update `TestE2E_GetMeasurementHistory` to require `evidence_status`, `diagnostic_ready`, and `evidence_gaps` in the REST response body.

- [x] **Step 2: Verify RED**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory
```

Expected: FAIL because mock history response is empty.

- [x] **Step 3: Update gateway mock history response**

Return one `MeasurementSummary` with `research_only`, `DiagnosticReady=false`, and `clinical_lock_required`.

- [x] **Step 4: Verify GREEN**

Run the same Go test. Expected: PASS.

## Task 2: Flutter History Domain and Mappers

- [x] **Step 1: Write failing Flutter history mapping tests**

Add tests requiring:

- native repository history maps generated summary evidence fields
- REST repository history maps snake_case evidence fields
- REST repository history maps legacy camelCase evidence fields

- [x] **Step 2: Verify RED**

Run targeted Flutter tests. Expected: FAIL due missing domain fields and mapper support.

- [x] **Step 3: Add domain fields and mappings**

Add fields to `MeasurementHistoryItem`, native mapping, REST mapping helper.

- [x] **Step 4: Verify GREEN**

Run the same Flutter tests. Expected: PASS.

## Task 3: Audit and Final Gates

- [x] **Step 1: Record audit and logs**

Document TDD record, self-review, quality gates, and next-stage guidance.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_rest.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart
```

Expected: PASS.
