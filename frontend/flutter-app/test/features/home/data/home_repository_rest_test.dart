import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/home/data/home_repository_rest.dart';
import 'package:manpasik/features/home/domain/home_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('HomeRepositoryRest', () {
    test('HomeRepositoryRest는 HomeRepository를 구현한다', () {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = HomeRepositoryRest(client);
      expect(repo, isA<HomeRepository>());
    });

    test('getHealthSummary는 DioException 시 빈 HealthSummary를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = HomeRepositoryRest(client);
      final summary = await repo.getHealthSummary('user-1');
      expect(summary, isA<HealthSummary>());
      expect(summary.latestBloodSugar, isNull);
      expect(summary.overallStatus, isNull);
      expect(summary.lastMeasuredAt, isNull);
    });

    test('getUnreadNotificationCount는 DioException 시 0을 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = HomeRepositoryRest(client);
      final count = await repo.getUnreadNotificationCount('user-1');
      expect(count, 0);
    });

    test('getRecentMeasurements는 DioException 시 빈 리스트를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = HomeRepositoryRest(client);
      final measurements = await repo.getRecentMeasurements('user-1');
      expect(measurements, isEmpty);
    });

    test('getRecentMeasurements limit 파라미터 확인', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = HomeRepositoryRest(client);
      final measurements =
          await repo.getRecentMeasurements('user-1', limit: 10);
      expect(measurements, isEmpty);
    });
  });

  group('HealthSummary 도메인 모델', () {
    test('기본 생성자로 빈 HealthSummary 생성', () {
      const summary = HealthSummary();
      expect(summary.latestBloodSugar, isNull);
      expect(summary.latestBloodPressureSystolic, isNull);
      expect(summary.latestBloodPressureDiastolic, isNull);
      expect(summary.latestHeartRate, isNull);
      expect(summary.overallStatus, isNull);
      expect(summary.lastMeasuredAt, isNull);
    });

    test('모든 필드가 있는 HealthSummary 생성', () {
      final now = DateTime.now();
      final summary = HealthSummary(
        latestBloodSugar: 95.0,
        latestBloodPressureSystolic: 120.0,
        latestBloodPressureDiastolic: 80.0,
        latestHeartRate: 72,
        overallStatus: 'good',
        lastMeasuredAt: now,
      );
      expect(summary.latestBloodSugar, 95.0);
      expect(summary.latestBloodPressureSystolic, 120.0);
      expect(summary.latestBloodPressureDiastolic, 80.0);
      expect(summary.latestHeartRate, 72);
      expect(summary.overallStatus, 'good');
      expect(summary.lastMeasuredAt, now);
    });
  });
}
