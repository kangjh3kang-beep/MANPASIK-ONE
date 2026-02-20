import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/data_hub/data/data_hub_repository_rest.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('DataHubRepositoryRest', () {
    test('DataHubRepositoryRest는 DataHubRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DataHubRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<DataHubRepository>());
    });

    test('getTrendData는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DataHubRepositoryRest(client, userId: 'user-1');
      final data = await repo.getTrendData(
        biomarkerType: 'glucose',
        from: DateTime(2026, 1, 1),
        to: DateTime(2026, 2, 19),
      );
      expect(data, isEmpty);
    });
  });

  group('DataHub 도메인 모델', () {
    test('TrendDataPoint 생성 확인', () {
      final point = TrendDataPoint(
        timestamp: DateTime(2026, 2, 19, 10, 0),
        value: 95.0,
        unit: 'mg/dL',
        biomarkerType: 'glucose',
        isWithinRange: true,
      );
      expect(point.value, 95.0);
      expect(point.isWithinRange, isTrue);
    });
  });
}
