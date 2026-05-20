import 'package:dio/dio.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/measurement/data/measurement_process_gateway_mapper.dart';

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
  Future<ProcessMeasurementResult> processMeasurement(
    ProcessMeasurementRequest request,
  ) async {
    final res = await _client.processMeasurement(
      data: encodeProcessMeasurementRequest(request),
    );
    return decodeProcessMeasurementResult(res, request);
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
            sessionId: _stringField(map, 'session_id', 'sessionId'),
            cartridgeType: _stringField(map, 'cartridge_type', 'cartridgeType'),
            primaryValue:
                _numberField(map, 'primary_value', 'primaryValue') ?? 0.0,
            unit: _stringField(map, 'unit', 'unit'),
            evidenceStatus:
                _stringField(map, 'evidence_status', 'evidenceStatus',
                    fallback: 'unknown'),
            diagnosticReady:
                _boolField(map, 'diagnostic_ready', 'diagnosticReady'),
            evidenceGaps: _stringListField(
              map,
              'evidence_gaps',
              'evidenceGaps',
            ),
            measuredAt:
                _stringField(map, 'measured_at', 'measuredAt').isNotEmpty
                    ? DateTime.tryParse(
                        _stringField(map, 'measured_at', 'measuredAt'),
                      )
                : null,
          );
        }).toList(),
        totalCount: _intField(res, 'total_count', 'totalCount'),
      );
    } on DioException catch (error) {
      return MeasurementHistoryResult(
        items: const [],
        totalCount: 0,
        isStale: true,
        errorMessage: '측정 기록을 새로고침할 수 없습니다: ${error.message}',
      );
    }
  }
}

String _stringField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey, {
  String fallback = '',
}) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is String ? value : fallback;
}

double? _numberField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is num ? value.toDouble() : null;
}

int _intField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is int ? value : 0;
}

bool _boolField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  return value is bool ? value : false;
}

List<String> _stringListField(
  Map<String, dynamic> map,
  String snakeKey,
  String camelKey,
) {
  final value = map[snakeKey] ?? map[camelKey];
  if (value is! List) {
    return const [];
  }
  return List.unmodifiable(value.whereType<String>());
}
