import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

import '../../helpers/fake_repositories.dart';

void main() {
  test('RealMode history failure surfaces stale/error result', () async {
    final container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
        measurementRepositoryProvider.overrideWithValue(
          _ThrowingMeasurementRepository(),
        ),
      ],
    );
    addTearDown(container.dispose);

    final loggedIn = await container
        .read(authProvider.notifier)
        .login('real@manpasik.com', 'password123');
    expect(loggedIn, isTrue);

    final result = await container.read(measurementHistoryProvider.future);

    expect(result.items, isEmpty);
    expect(result.totalCount, 0);
    expect(result.isStale, isTrue);
    expect(result.hasError, isTrue);
    expect(result.errorMessage, contains('backend down'));
  });

  test('DemoMode history keeps explicit mock data fresh', () async {
    final container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
        measurementRepositoryProvider.overrideWithValue(
          _ThrowingMeasurementRepository(),
        ),
      ],
    );
    addTearDown(container.dispose);

    container.read(authProvider.notifier).loginAsDemo();

    final result = await container.read(measurementHistoryProvider.future);

    expect(result.items, hasLength(10));
    expect(result.items.first.sessionId, startsWith('mock-measure-'));
    expect(result.isStale, isFalse);
    expect(result.hasError, isFalse);
  });
}

class _ThrowingMeasurementRepository implements MeasurementRepository {
  @override
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  }) {
    throw UnimplementedError();
  }

  @override
  Future<EndSessionResult?> endSession(String sessionId) {
    throw UnimplementedError();
  }

  @override
  Future<ProcessMeasurementResult> processMeasurement(
    ProcessMeasurementRequest request,
  ) {
    throw UnimplementedError();
  }

  @override
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    throw StateError('backend down');
  }
}
