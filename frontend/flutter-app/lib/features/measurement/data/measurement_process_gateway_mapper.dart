import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

Map<String, dynamic> encodeProcessMeasurementRequest(
  ProcessMeasurementRequest request,
) {
  return {
    'session_id': request.sessionId,
    'device_id': request.deviceId,
    'user_id': request.userId,
    'cartridge_type': request.cartridgeType,
    'raw_channels': request.rawChannels,
    'differential': {
      's_det': request.sDet,
      's_ref': request.sRef,
      'alpha': request.alpha,
      's_corrected': request.sCorrected,
    },
    'primary_value': request.primaryValue,
    'unit': request.unit,
    'confidence': request.confidence,
    'fingerprint_vector': request.fingerprintVector,
    'env_meta': {
      'temp_c': request.tempC,
      'humidity_pct': request.humidityPct,
    },
    'battery_pct': request.batteryPct,
  };
}

ProcessMeasurementResult decodeProcessMeasurementResult(
  Map<String, dynamic> response,
  ProcessMeasurementRequest fallback,
) {
  final evidenceGapsRaw =
      _field<List<dynamic>>(response, 'evidence_gaps', 'evidenceGaps') ??
          const <dynamic>[];

  return ProcessMeasurementResult(
    sessionId: _field<String>(response, 'session_id', 'sessionId') ??
        fallback.sessionId,
    primaryValue: _numberField(response, 'primary_value', 'primaryValue') ??
        fallback.primaryValue,
    unit: response['unit'] as String? ?? fallback.unit,
    confidence: _numberField(response, 'confidence', 'confidence') ??
        fallback.confidence,
    processedAt: _field<String>(response, 'processed_at', 'processedAt') != null
        ? DateTime.tryParse(
            _field<String>(response, 'processed_at', 'processedAt')!,
          )
        : null,
    evidenceStatus:
        _field<String>(response, 'evidence_status', 'evidenceStatus') ??
            'unknown',
    diagnosticReady:
        _field<bool>(response, 'diagnostic_ready', 'diagnosticReady') ?? false,
    evidenceGaps: evidenceGapsRaw
        .map((value) => value.toString())
        .toList(growable: false),
  );
}

T? _field<T>(Map<String, dynamic> response, String snake, String camel) {
  final value = response[snake] ?? response[camel];
  return value is T ? value : null;
}

double? _numberField(
  Map<String, dynamic> response,
  String snake,
  String camel,
) {
  final value = response[snake] ?? response[camel];
  return value is num ? value.toDouble() : null;
}
