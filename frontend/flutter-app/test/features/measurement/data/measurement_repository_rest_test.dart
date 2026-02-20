import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_rest.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('MeasurementRepositoryRest', () {
    test('MeasurementRepositoryRest는 MeasurementRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      expect(repo, isA<MeasurementRepository>());
    });

    test('getHistory는 DioException 시 빈 결과를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      final result = await repo.getHistory(userId: 'user-1');
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
    });

    test('getHistory 커스텀 limit/offset', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MeasurementRepositoryRest(client);
      final result = await repo.getHistory(userId: 'user-1', limit: 5, offset: 10);
      expect(result.items, isEmpty);
    });
  });

  group('Measurement 도메인 모델', () {
    test('StartSessionResult 생성 확인', () {
      const result = StartSessionResult(sessionId: 'sess-1');
      expect(result.sessionId, 'sess-1');
      expect(result.startedAt, isNull);
    });

    test('EndSessionResult 생성 확인', () {
      const result = EndSessionResult(
        sessionId: 'sess-1',
        totalMeasurements: 5,
      );
      expect(result.totalMeasurements, 5);
      expect(result.endedAt, isNull);
    });

    test('MeasurementHistoryResult 빈 결과 생성', () {
      const result = MeasurementHistoryResult(items: [], totalCount: 0);
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
    });

    test('MeasurementHistoryItem 생성', () {
      final item = MeasurementHistoryItem(
        sessionId: 'sess-1',
        cartridgeType: 'glucose',
        primaryValue: 95.5,
        unit: 'mg/dL',
        measuredAt: DateTime(2026, 2, 19, 10, 0),
      );
      expect(item.primaryValue, 95.5);
      expect(item.cartridgeType, 'glucose');
    });
  });
}
