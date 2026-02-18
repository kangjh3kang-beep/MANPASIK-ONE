import 'package:dio/dio.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 MeasurementRepository 구현체
///
/// 웹 플랫폼에서 gRPC 대신 REST API를 통해 측정 데이터 처리.
class MeasurementRepositoryRest implements MeasurementRepository {
  MeasurementRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  }) async {
    final res = await _client.startSession(
      deviceId: deviceId,
      userId: userId,
      cartridgeId: cartridgeId,
    );
    return StartSessionResult(
      sessionId: res['session_id'] as String? ?? '',
      startedAt: null,
    );
  }

  @override
  Future<EndSessionResult?> endSession(String sessionId) async {
    final res = await _client.endSession(sessionId);
    return EndSessionResult(
      sessionId: res['session_id'] as String? ?? sessionId,
      totalMeasurements: res['total_measurements'] as int? ?? 0,
      endedAt: null,
    );
  }

  @override
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final res = await _client.getMeasurementHistory(
        userId,
        limit: limit,
        offset: offset,
      );
      final measurements = res['measurements'] as List<dynamic>? ?? [];
      return MeasurementHistoryResult(
        items: measurements.map((m) {
          final map = m as Map<String, dynamic>;
          return MeasurementHistoryItem(
            sessionId: map['session_id'] as String? ?? '',
            cartridgeType: map['cartridge_type'] as String? ?? '',
            primaryValue: (map['primary_value'] as num?)?.toDouble() ?? 0.0,
            unit: map['unit'] as String? ?? '',
            measuredAt: map['measured_at'] != null
                ? DateTime.tryParse(map['measured_at'] as String)
                : null,
          );
        }).toList(),
        totalCount: res['total_count'] as int? ?? 0,
      );
    } on DioException {
      return const MeasurementHistoryResult(items: [], totalCount: 0);
    }
  }
}
