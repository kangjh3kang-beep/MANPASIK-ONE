import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/features/measurement/application/measurement_golden_path_orchestrator.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

void main() {
  test(
      'MeasurementGoldenPathOrchestrator runs cartridge-session-engine-server-end flow',
      () async {
    final repo = _FakeMeasurementRepository();
    final engine = _FakeMeasurementEngine();
    final phases = <MeasurementGoldenPathPhase>[];
    final snapshots = <MeasurementGoldenPathSnapshot>[];
    final traceEvents = <MeasurementGoldenPathTraceEvent>[];
    final orchestrator = MeasurementGoldenPathOrchestrator(
      repository: repo,
      engine: engine,
      onSnapshot: (snapshot) {
        phases.add(snapshot.phase);
        snapshots.add(snapshot);
      },
      traceSink: traceEvents.add,
    );

    final result = await orchestrator.run(
      userId: 'user-1',
      deviceId: 'device-1',
    );

    expect(result.sessionId, 'session-1');
    expect(repo.started, isTrue);
    expect(repo.processed, isTrue);
    expect(repo.ended, isTrue);
    expect(repo.lastProcessRequest?.rawChannels, hasLength(88));
    expect(repo.lastProcessRequest?.fingerprintVector, hasLength(88));
    expect(phases, [
      MeasurementGoldenPathPhase.readinessChecked,
      MeasurementGoldenPathPhase.cartridgeRead,
      MeasurementGoldenPathPhase.sessionStarted,
      MeasurementGoldenPathPhase.engineProcessed,
      MeasurementGoldenPathPhase.serverProcessed,
      MeasurementGoldenPathPhase.sessionEnded,
    ]);
    expect(traceEvents.map((event) => event.phase), phases);
    expect(traceEvents.first.engineMode, 'native');
    expect(traceEvents.last.sessionId, 'session-1');
    expect(traceEvents.last.toJson()['phase'], 'sessionEnded');
    expect(traceEvents.last.toLogLine(), contains('phase=sessionEnded'));
    final remotePayload = traceEvents.last.toRemoteObservabilityJson();
    expect(remotePayload['has_primary_value'], isTrue);
    expect(remotePayload.containsKey('primary_value'), isFalse);
    final serverSnapshot = snapshots.firstWhere(
      (snapshot) => snapshot.phase == MeasurementGoldenPathPhase.serverProcessed,
    );
    expect(serverSnapshot.evidenceStatus, 'research_only');
    expect(serverSnapshot.diagnosticReady, isFalse);
    expect(serverSnapshot.evidenceGaps, contains('clinical_lock_required'));
  });

  test('MeasurementGoldenPathOrchestrator blocks when engine is not ready',
      () async {
    final repo = _FakeMeasurementRepository();
    final engine = _FakeMeasurementEngine(canRun: false);
    final phases = <MeasurementGoldenPathPhase>[];
    final traceEvents = <MeasurementGoldenPathTraceEvent>[];
    final orchestrator = MeasurementGoldenPathOrchestrator(
      repository: repo,
      engine: engine,
      onSnapshot: (snapshot) => phases.add(snapshot.phase),
      traceSink: traceEvents.add,
    );

    await expectLater(
      orchestrator.run(userId: 'user-1', deviceId: 'device-1'),
      throwsA(isA<MeasurementGoldenPathReadinessException>()),
    );

    expect(repo.started, isFalse);
    expect(repo.processed, isFalse);
    expect(repo.ended, isFalse);
    expect(phases, [
      MeasurementGoldenPathPhase.readinessChecked,
      MeasurementGoldenPathPhase.failed,
    ]);
    expect(traceEvents.map((event) => event.phase), phases);
    expect(traceEvents.first.engineMode, 'stub-fallback');
    expect(traceEvents.last.isFailure, isTrue);
    expect(traceEvents.last.failureReason, contains('blocked'));
  });
}

class _FakeMeasurementEngine implements MeasurementEngine {
  _FakeMeasurementEngine({this.canRun = true});

  final bool canRun;

  @override
  Future<MeasurementEngineReadiness> checkReadiness() async {
    return MeasurementEngineReadiness(
      canRun: canRun,
      modeLabel: canRun ? 'native' : 'stub-fallback',
      engineVersion: canRun ? '0.1.0-test' : '0.1.0-stub',
      message: canRun ? 'ready' : 'blocked',
    );
  }

  @override
  Future<CartridgeInfoDto> readCartridge() async {
    return const CartridgeInfoDto(
      cartridgeId: 'cart-1',
      cartridgeType: 'Glucose',
      lotId: 'LOT-A',
      expiryDate: '20270101',
      remainingUses: 7,
    );
  }

  @override
  Future<MeasurementPipelineResult> runPipeline({
    required String deviceId,
    required String biomarker,
    required String unit,
  }) async {
    return MeasurementPipelineResult(
      measurement: MeasurementResultDto(
        primaryValue: 95,
        referenceValue: 10,
        differentialValue: 85.2,
        snr: 35,
        confidence: 0.96,
        unit: unit,
        biomarker: biomarker,
        timestamp: DateTime(2026, 5, 1),
      ),
      analysis: const AiAnalysisDto(
        riskLevel: 'normal',
        healthScore: 90,
        summary: 'ok',
        recommendations: ['continue'],
        trend: 'stable',
      ),
      pipelineDurationMs: 120,
    );
  }
}

class _FakeMeasurementRepository implements MeasurementRepository {
  bool started = false;
  bool processed = false;
  bool ended = false;
  ProcessMeasurementRequest? lastProcessRequest;

  @override
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  }) async {
    started = true;
    return const StartSessionResult(sessionId: 'session-1');
  }

  @override
  Future<ProcessMeasurementResult> processMeasurement(
    ProcessMeasurementRequest request,
  ) async {
    processed = true;
    lastProcessRequest = request;
    return ProcessMeasurementResult(
      sessionId: request.sessionId,
      primaryValue: request.primaryValue,
      unit: request.unit,
      confidence: request.confidence,
      evidenceStatus: 'research_only',
      diagnosticReady: false,
      evidenceGaps: const ['clinical_lock_required'],
    );
  }

  @override
  Future<EndSessionResult?> endSession(String sessionId) async {
    ended = true;
    return EndSessionResult(sessionId: sessionId, totalMeasurements: 1);
  }

  @override
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
}
