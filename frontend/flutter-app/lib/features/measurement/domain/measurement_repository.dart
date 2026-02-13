/// 측정 Repository 인터페이스
abstract class MeasurementRepository {
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  });
  Future<EndSessionResult?> endSession(String sessionId);
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  });
}

class StartSessionResult {
  final String sessionId;
  final DateTime? startedAt;

  const StartSessionResult({
    required this.sessionId,
    this.startedAt,
  });
}

class EndSessionResult {
  final String sessionId;
  final int totalMeasurements;
  final DateTime? endedAt;

  const EndSessionResult({
    required this.sessionId,
    required this.totalMeasurements,
    this.endedAt,
  });
}

class MeasurementHistoryResult {
  final List<MeasurementHistoryItem> items;
  final int totalCount;

  const MeasurementHistoryResult({
    required this.items,
    required this.totalCount,
  });
}

class MeasurementHistoryItem {
  final String sessionId;
  final String cartridgeType;
  final double primaryValue;
  final String unit;
  final DateTime? measuredAt;

  const MeasurementHistoryItem({
    required this.sessionId,
    required this.cartridgeType,
    required this.primaryValue,
    required this.unit,
    this.measuredAt,
  });
}
