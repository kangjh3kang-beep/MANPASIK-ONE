/// 측정 Repository 인터페이스
abstract class MeasurementRepository {
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  });
  Future<EndSessionResult?> endSession(String sessionId);
  Future<ProcessMeasurementResult> processMeasurement(
    ProcessMeasurementRequest request,
  );
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

class ProcessMeasurementRequest {
  final String sessionId;
  final String deviceId;
  final String userId;
  final String cartridgeType;
  final List<double> rawChannels;
  final double sDet;
  final double sRef;
  final double alpha;
  final double sCorrected;
  final double primaryValue;
  final String unit;
  final double confidence;
  final List<double> fingerprintVector;
  final double tempC;
  final double humidityPct;
  final int batteryPct;

  const ProcessMeasurementRequest({
    required this.sessionId,
    required this.deviceId,
    required this.userId,
    required this.cartridgeType,
    required this.rawChannels,
    required this.sDet,
    required this.sRef,
    required this.alpha,
    required this.sCorrected,
    required this.primaryValue,
    required this.unit,
    required this.confidence,
    required this.fingerprintVector,
    required this.tempC,
    required this.humidityPct,
    required this.batteryPct,
  });
}

class ProcessMeasurementResult {
  final String sessionId;
  final double primaryValue;
  final String unit;
  final double confidence;
  final DateTime? processedAt;
  final String evidenceStatus;
  final bool diagnosticReady;
  final List<String> evidenceGaps;

  const ProcessMeasurementResult({
    required this.sessionId,
    required this.primaryValue,
    required this.unit,
    required this.confidence,
    this.processedAt,
    this.evidenceStatus = 'unknown',
    this.diagnosticReady = false,
    this.evidenceGaps = const [],
  });
}

class MeasurementHistoryResult {
  final List<MeasurementHistoryItem> items;
  final int totalCount;
  final bool isStale;
  final String? errorMessage;

  const MeasurementHistoryResult({
    required this.items,
    required this.totalCount,
    this.isStale = false,
    this.errorMessage,
  });

  bool get hasError => errorMessage != null && errorMessage!.isNotEmpty;
}

class MeasurementHistoryItem {
  final String sessionId;
  final String cartridgeType;
  final double primaryValue;
  final String unit;
  final String evidenceStatus;
  final bool diagnosticReady;
  final List<String> evidenceGaps;
  final DateTime? measuredAt;

  const MeasurementHistoryItem({
    required this.sessionId,
    required this.cartridgeType,
    required this.primaryValue,
    required this.unit,
    this.evidenceStatus = 'unknown',
    this.diagnosticReady = false,
    this.evidenceGaps = const [],
    this.measuredAt,
  });
}
