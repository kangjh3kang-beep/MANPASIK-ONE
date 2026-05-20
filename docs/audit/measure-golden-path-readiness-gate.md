# Measure Golden Path Readiness Gate

**Date**: 2026-05-01  
**Owner**: Codex  
**Sprint**: H1 Measure golden path hardening

## Purpose

The Measure route must not silently run a simulated Rust path in production. This gate turns the previous MR-002 mock-retirement requirement into executable behavior.

## Implemented Contract

1. `RustBridge.diagnostics` exposes:
   - initialization state
   - native/stub mode
   - native platform eligibility
   - release build state
   - engine version
   - explicit release stub allowance

2. `MeasurementGoldenPathOrchestrator` now checks readiness before:
   - NFC cartridge read
   - session start
   - Rust measurement pipeline
   - server processing

3. Release behavior:
   - native Rust available: measurement can run
   - debug/profile with stub fallback: measurement can run with diagnostic message
   - release without native Rust: measurement is blocked
   - release stub override: only allowed with `MANPASIK_RELEASE_ALLOW_STUB_MEASUREMENT=true`

## Trace Evidence

| Evidence | Path |
|---|---|
| Runtime diagnostics | `frontend/flutter-app/lib/core/services/rust_ffi_stub.dart` |
| Gate orchestration | `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart` |
| UI phase display | `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart` |
| Pass and blocked tests | `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart` |

## Verification

- `dart format`: PASS
- `flutter analyze` on changed Measure/Rust bridge files: PASS
- `flutter test` for Measure golden path, REST repository, and domain tests: PASS

