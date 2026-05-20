# Measurement Evidence Persistence P5 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans for task-by-task execution. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** measurement evidence fields를 Timescale/Postgres 저장 스키마와 history response contract까지 보존한다.

**Architecture:** `MeasurementData`에 이미 존재하는 evidence fields를 `measurement_data` 컬럼으로 저장하고, `MeasurementSummary`와 `GetMeasurementHistory` proto response에 동일한 상태를 노출한다. 기존 `MeasurementResult` stream contract와 Flutter P4 UI는 그대로 유지한다.

**Tech Stack:** Go measurement-service tests, proto generated Go/Dart, Flutter proto compile gate, schema static guard, existing SSOT/security/assay gates.

---

## Stage Gate Rules

- RED 테스트를 먼저 작성하고 실패를 확인한다.
- 새 proto fields는 기존 field number를 변경하지 않고 append-only로 추가한다.
- `research_only`는 저장/조회돼도 diagnostic ready로 승격하지 않는다.
- DB migration 범위는 init SQL 컬럼과 summary view 반영으로 한정한다.

## File Structure

- Modify: `backend/shared/proto/manpasik.proto`
  - `MeasurementSummary`에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 append-only로 추가한다.
- Modify generated:
  - `backend/shared/gen/go/v1/manpasik.pb.go`
  - `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
  - `frontend/flutter-app/lib/generated/manpasik.pb.dart`
  - `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
  - `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`
  - `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`
- Modify: `backend/services/measurement-service/internal/service/measurement.go`
  - `MeasurementSummary`에 evidence fields를 추가한다.
- Modify: `backend/services/measurement-service/internal/handler/grpc.go`
  - history response summary에 evidence fields를 매핑한다.
- Modify: `backend/services/measurement-service/internal/repository/postgres/measurement.go`
  - Store/GetHistory SQL에 evidence columns를 추가한다.
- Modify: `infrastructure/database/init/04-measurement.sql`
  - `measurement_data`와 `measurement_summary` view에 evidence columns를 추가한다.
- Modify/Add tests:
  - `backend/services/measurement-service/internal/handler/grpc_stream_test.go`
  - `backend/services/measurement-service/internal/service/measurement_test.go`
  - `backend/services/measurement-service/internal/repository/postgres/measurement_schema_test.go`
  - `frontend/flutter-app/test/generated/measurement_summary_evidence_contract_test.dart`
- Create: `docs/audit/measurement-evidence-persistence-p5.md`
- Modify: `CHANGELOG.md`, `CONTEXT.md`

## Task 1: History Contract RED

- [x] **Step 1: Write failing Go handler history test**

Add a test that returns a `service.MeasurementSummary` with `research_only`, `DiagnosticReady=false`, and `clinical_lock_required`, then asserts the gRPC `MeasurementSummary` response exposes those fields.

- [x] **Step 2: Write failing service history preservation test**

Seed the mock measurement repository with evidence fields and assert `GetHistory` returns them.

- [x] **Step 3: Write failing schema guard**

Add a lightweight repository package test that ensures `04-measurement.sql` contains `evidence_status`, `diagnostic_ready`, and `evidence_gaps`.

- [x] **Step 4: Verify RED**

Run targeted Go tests and confirm failures are due to missing fields/schema.

## Task 2: Proto and Service Mapping

- [x] **Step 1: Append proto fields**

Add fields 6~8 to `MeasurementSummary`.

- [x] **Step 2: Regenerate Go and Dart proto outputs**

Use local protoc plugins already installed in WSL.

- [x] **Step 3: Add service summary fields and handler mapping**

Map evidence fields into gRPC history response.

- [x] **Step 4: Verify handler/service GREEN**

Run targeted Go tests.

## Task 3: Persistence Columns and Repository Mapping

- [x] **Step 1: Add DB columns and summary view fields**

Update `04-measurement.sql`.

- [x] **Step 2: Persist evidence fields**

Update repository `Store` insert and `GetHistory` select/scan.

- [x] **Step 3: Verify schema/repository GREEN**

Run repository package tests.

## Task 4: Dart Generated Contract and Final Gates

- [x] **Step 1: Add Dart generated contract test**

Assert `MeasurementSummary` serializes/deserializes evidence fields.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
go test -count=1 ./services/measurement-service/internal/handler ./services/measurement-service/internal/service ./services/measurement-service/internal/repository/postgres
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/generated/measurement_summary_evidence_contract_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/generated/manpasik.pbenum.dart lib/generated/manpasik.pbjson.dart test/generated/measurement_summary_evidence_contract_test.dart
cd frontend/flutter-app && bash scripts/check_proto_generation_compile_gate.sh
```

Expected: PASS.
