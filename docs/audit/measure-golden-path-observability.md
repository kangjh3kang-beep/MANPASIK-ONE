# Measure Golden Path Observability

**Date**: 2026-05-01  
**Owner**: Codex  
**Sprint**: H1 Measure golden path hardening

## Purpose

The Measure golden path must leave a deterministic trace for every operational phase. This makes failures diagnosable by phase instead of by a generic UI or API error.

## Trace Contract

Each emitted trace event includes:

- `phase`
- `elapsed_ms`
- `occurred_at`
- `session_id`
- `cartridge_id`
- `engine_mode`
- `primary_value`
- `unit`
- `confidence`
- `failure_reason`
- `diagnostic_message`

## Phase Sequence

Expected happy path:

1. `readinessChecked`
2. `cartridgeRead`
3. `sessionStarted`
4. `engineProcessed`
5. `serverProcessed`
6. `sessionEnded`

Failure path example:

1. `readinessChecked`
2. `failed`

## Implementation Evidence

| Evidence | Path |
|---|---|
| Trace event model and logger bridge | `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart` |
| Default app logging | `frontend/flutter-app/lib/core/services/app_logger.dart` |
| Phase trace tests | `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart` |

## Verification

- `dart format`: PASS
- `flutter analyze` on changed Measure observability files: PASS
- `flutter test` for Measure golden path, REST repository, and domain tests: PASS

