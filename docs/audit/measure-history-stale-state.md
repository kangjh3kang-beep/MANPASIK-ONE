# Measure History Stale 상태 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: MR-001 RealMode 측정 기록 실패 상태 노출

## 목적

Measure mock retirement register의 MR-001은 DemoMode mock history를 조건부 허용하되, RealMode 실패를 빈 기록처럼 숨기지 않는 것이 핵심 조건이다. 이번 단계는 measurement history 조회 실패를 `stale/error` 상태로 명확히 노출해 운영자가 "기록 없음"과 "동기화 실패"를 구분할 수 있게 한다.

## 구현 결과

- `MeasurementHistoryResult`에 `isStale`, `errorMessage`, `hasError` 계약을 추가했다.
- REST history repository는 `DioException`을 빈 정상 결과로 삼키지 않고 stale/error result로 반환한다.
- `measurementHistoryProvider`는 RealMode repository 예외를 stale/error result로 변환한다.
- DemoMode mock history는 기존처럼 명시적 `authState.isDemo`에서만 fresh mock data로 반환한다.
- measurement result screen은 history가 stale/error이고 항목이 없으면 동기화 문제 상태와 재시도 버튼을 표시한다.
- home dashboard는 history stale 상태를 hero copy와 stats card에 반영한다.

## 변경 파일

- `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
  - `MeasurementHistoryResult`에 stale/error 상태 계약을 추가했다.
- `frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart`
  - history REST 실패를 stale/error result로 반환하도록 변경했다.
- `frontend/flutter-app/lib/core/providers/grpc_provider.dart`
  - RealMode history provider 실패를 stale/error result로 노출하고 DemoMode mock freshness를 유지했다.
- `frontend/flutter-app/lib/shared/providers/ecosystem_providers.dart`
  - home dashboard aggregate model에 measurement history stale/error 상태를 연결했다.
- `frontend/flutter-app/lib/features/measurement/presentation/measurement_result_screen.dart`
  - history 동기화 실패 empty state, inline warning, retry action을 추가했다.
- `frontend/flutter-app/lib/features/home/presentation/home_screen.dart`
  - dashboard hero/stat에 측정 기록 동기화 확인 필요 상태를 표시한다.
- `frontend/flutter-app/test/core/providers/measurement_history_provider_test.dart`
  - RealMode 실패가 stale/error result로 노출되는지와 DemoMode mock이 fresh 상태인지 검증했다.
- `frontend/flutter-app/test/features/measurement/data/measurement_repository_rest_test.dart`
  - REST history failure contract를 stale/error expectation으로 갱신했다.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
export PATH="/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin"
dart format \
  lib/features/measurement/domain/measurement_repository.dart \
  lib/features/measurement/data/measurement_repository_rest.dart \
  lib/core/providers/grpc_provider.dart \
  lib/shared/providers/ecosystem_providers.dart \
  lib/features/measurement/presentation/measurement_result_screen.dart \
  lib/features/home/presentation/home_screen.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/core/providers/measurement_history_provider_test.dart
flutter test --no-pub \
  test/core/providers/measurement_history_provider_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/domain/measurement_domain_test.dart \
  test/features/domain_models_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_trace_sink_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
flutter analyze --no-pub \
  lib/features/measurement/domain/measurement_repository.dart \
  lib/features/measurement/data/measurement_repository_rest.dart \
  lib/core/providers/grpc_provider.dart \
  lib/shared/providers/ecosystem_providers.dart \
  lib/features/measurement/presentation/measurement_result_screen.dart \
  lib/features/home/presentation/home_screen.dart \
  test/core/providers/measurement_history_provider_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart
```

결과: PASS

## 다음 단계

- MR-009 native process REST bridge는 `docs/audit/measure-native-grpc-stream-binding.md`에서 gRPC stream 직결로 교체했다.
- 다음 단계는 전체 Dart proto generation 공식화와 measurement-service integration smoke test다.
