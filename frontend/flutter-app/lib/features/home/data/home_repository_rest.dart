import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/home/domain/home_repository.dart';

/// REST Gateway를 사용하는 HomeRepository 구현체
class HomeRepositoryRest implements HomeRepository {
  HomeRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<HealthSummary> getHealthSummary(String userId) async {
    try {
      final res = await _client.getHealthSummary(userId);
      return HealthSummary(
        latestBloodSugar: (res['latest_blood_sugar'] as num?)?.toDouble(),
        latestBloodPressureSystolic:
            (res['latest_blood_pressure_systolic'] as num?)?.toDouble(),
        latestBloodPressureDiastolic:
            (res['latest_blood_pressure_diastolic'] as num?)?.toDouble(),
        latestHeartRate: res['latest_heart_rate'] as int?,
        overallStatus: res['overall_status'] as String?,
        lastMeasuredAt: res['last_measured_at'] != null
            ? DateTime.tryParse(res['last_measured_at'] as String)
            : null,
      );
    } on DioException {
      return const HealthSummary();
    }
  }

  @override
  Future<int> getUnreadNotificationCount(String userId) async {
    try {
      final res = await _client.getUnreadCount(userId);
      return res['count'] as int? ?? 0;
    } on DioException {
      return 0;
    }
  }

  @override
  Future<List<Map<String, dynamic>>> getRecentMeasurements(
    String userId, {
    int limit = 5,
  }) async {
    try {
      final res =
          await _client.getMeasurementHistory(userId, limit: limit);
      final measurements =
          res['measurements'] as List<dynamic>? ?? [];
      return measurements.cast<Map<String, dynamic>>();
    } on DioException {
      return [];
    }
  }
}
