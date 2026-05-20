import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/data/measurement_process_gateway_mapper.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

void main() {
  group('decodeProcessMeasurementResult', () {
    test('snake_case evidence fields를 도메인 결과로 매핑한다', () {
      final result = decodeProcessMeasurementResult(
        {
          'session_id': 'session-rest-1',
          'primary_value': 88.1,
          'unit': 'mg/dL',
          'confidence': 0.91,
          'evidence_status': 'research_only',
          'diagnostic_ready': false,
          'evidence_gaps': ['clinical_lock_required'],
        },
        _fallback(),
      );

      expect(result.sessionId, 'session-rest-1');
      expect(result.evidenceStatus, 'research_only');
      expect(result.diagnosticReady, isFalse);
      expect(result.evidenceGaps, contains('clinical_lock_required'));
    });

    test('camelCase evidence fields도 하위 호환으로 매핑한다', () {
      final result = decodeProcessMeasurementResult(
        {
          'sessionId': 'session-rest-2',
          'primaryValue': 77.2,
          'unit': 'mg/dL',
          'confidence': 0.89,
          'evidenceStatus': 'research_only',
          'diagnosticReady': false,
          'evidenceGaps': ['clinical_lock_required'],
        },
        _fallback(),
      );

      expect(result.sessionId, 'session-rest-2');
      expect(result.primaryValue, 77.2);
      expect(result.evidenceStatus, 'research_only');
      expect(result.diagnosticReady, isFalse);
      expect(result.evidenceGaps, contains('clinical_lock_required'));
    });
  });
}

ProcessMeasurementRequest _fallback() {
  return const ProcessMeasurementRequest(
    sessionId: 'fallback-session',
    deviceId: 'device-1',
    userId: 'user-1',
    cartridgeType: 'glucose',
    rawChannels: [1.0, 2.0],
    sDet: 10.0,
    sRef: 1.0,
    alpha: 0.98,
    sCorrected: 9.02,
    primaryValue: 9.02,
    unit: 'mg/dL',
    confidence: 0.9,
    fingerprintVector: [1.0, 2.0],
    tempC: 24.0,
    humidityPct: 45.0,
    batteryPct: 90,
  );
}
