/// 홈 대시보드 도메인 모델 및 리포지토리
///
/// 건강 요약, 최근 측정, 알림 뱃지, 빠른 액션

/// 건강 요약 데이터
class HealthSummary {
  final double? latestBloodSugar;
  final double? latestBloodPressureSystolic;
  final double? latestBloodPressureDiastolic;
  final int? latestHeartRate;
  final String? overallStatus; // 'good', 'warning', 'critical'
  final DateTime? lastMeasuredAt;

  const HealthSummary({
    this.latestBloodSugar,
    this.latestBloodPressureSystolic,
    this.latestBloodPressureDiastolic,
    this.latestHeartRate,
    this.overallStatus,
    this.lastMeasuredAt,
  });
}

/// 홈 리포지토리 인터페이스
abstract class HomeRepository {
  Future<HealthSummary> getHealthSummary(String userId);
  Future<int> getUnreadNotificationCount(String userId);
  Future<List<Map<String, dynamic>>> getRecentMeasurements(String userId, {int limit = 5});
}
