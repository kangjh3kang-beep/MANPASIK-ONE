import 'dart:async';

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_impl.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/generated/manpasik.pb.dart';

void main() {
  test('native processMeasurement streams frame to MeasurementService gRPC',
      () async {
    final received = Completer<MeasurementData>();
    final repo = MeasurementRepositoryImpl(
      GrpcClientManager(),
      accessTokenProvider: () => 'native-token',
      streamMeasurement: (frames) async* {
        final frame = await frames.single;
        received.complete(frame);
        yield MeasurementResult(
          sessionId: 'session-1',
          primaryValue: 101.2,
          unit: 'mg/dL',
          confidence: 0.91,
          fingerprintVector: [0.1, 0.2, 0.3],
          evidenceStatus: 'research_only',
          diagnosticReady: false,
          evidenceGaps: ['clinical_lock_required'],
        );
      },
      getMeasurementHistory: (request) async {
        return GetHistoryResponse(
          measurements: [
            MeasurementSummary(
              sessionId: 'session-history',
              cartridgeType: 'glucose',
              primaryValue: 99.5,
              unit: 'mg/dL',
              evidenceStatus: 'research_only',
              diagnosticReady: false,
              evidenceGaps: ['clinical_lock_required'],
            ),
          ],
          totalCount: 1,
        );
      },
    );

    final result = await repo.processMeasurement(_request());
    final frame = await received.future.timeout(const Duration(seconds: 2));

    expect(frame.sessionId, 'session-1');
    expect(frame.rawChannels, [1.0, 2.0, 3.0]);
    expect(frame.differential.sDet, 100.0);
    expect(frame.differential.sRef, 5.0);
    expect(frame.differential.alpha, 0.95);
    expect(frame.differential.sCorrected, 95.25);
    expect(frame.envMeta.tempC, 24.5);
    expect(frame.envMeta.humidityPct, 45.0);

    expect(result.sessionId, 'session-1');
    expect(result.primaryValue, 101.2);
    expect(result.unit, 'mg/dL');
    expect(result.confidence, 0.91);
    expect(result.processedAt, isNull);
    expect(result.evidenceStatus, 'research_only');
    expect(result.diagnosticReady, isFalse);
    expect(result.evidenceGaps, contains('clinical_lock_required'));
  });

  test('native getHistory maps MeasurementSummary evidence fields', () async {
    final repo = MeasurementRepositoryImpl(
      GrpcClientManager(),
      accessTokenProvider: () => 'native-token',
      getMeasurementHistory: (request) async {
        return GetHistoryResponse(
          measurements: [
            MeasurementSummary(
              sessionId: 'session-history',
              cartridgeType: 'glucose',
              primaryValue: 99.5,
              unit: 'mg/dL',
              evidenceStatus: 'research_only',
              diagnosticReady: false,
              evidenceGaps: ['clinical_lock_required'],
            ),
          ],
          totalCount: 1,
        );
      },
    );

    final result = await repo.getHistory(userId: 'user-1');

    expect(result.totalCount, 1);
    expect(result.items.single.evidenceStatus, 'research_only');
    expect(result.items.single.diagnosticReady, isFalse);
    expect(
      result.items.single.evidenceGaps,
      contains('clinical_lock_required'),
    );
  });
}

ProcessMeasurementRequest _request() {
  return const ProcessMeasurementRequest(
    sessionId: 'session-1',
    deviceId: 'device-1',
    userId: 'user-1',
    cartridgeType: 'Glucose',
    rawChannels: [1.0, 2.0, 3.0],
    sDet: 100.0,
    sRef: 5.0,
    alpha: 0.95,
    sCorrected: 95.25,
    primaryValue: 95.0,
    unit: 'mg/dL',
    confidence: 0.96,
    fingerprintVector: [0.1, 0.2, 0.3],
    tempC: 24.5,
    humidityPct: 45.0,
    batteryPct: 88,
  );
}
