import 'dart:math' as math;

import 'package:manpasik/core/services/app_logger.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

enum MeasurementGoldenPathPhase {
  idle,
  readinessChecked,
  cartridgeRead,
  sessionStarted,
  engineProcessed,
  serverProcessed,
  sessionEnded,
  failed,
}

class MeasurementGoldenPathSnapshot {
  const MeasurementGoldenPathSnapshot({
    required this.phase,
    this.sessionId,
    this.cartridgeId,
    this.primaryValue,
    this.unit,
    this.confidence,
    this.failureReason,
    this.engineMode,
    this.diagnosticMessage,
    this.evidenceStatus,
    this.diagnosticReady = false,
    this.evidenceGaps = const [],
  });

  final MeasurementGoldenPathPhase phase;
  final String? sessionId;
  final String? cartridgeId;
  final double? primaryValue;
  final String? unit;
  final double? confidence;
  final String? failureReason;
  final String? engineMode;
  final String? diagnosticMessage;
  final String? evidenceStatus;
  final bool diagnosticReady;
  final List<String> evidenceGaps;

  bool get isTerminal =>
      phase == MeasurementGoldenPathPhase.sessionEnded ||
      phase == MeasurementGoldenPathPhase.failed;
}

class MeasurementEngineReadiness {
  const MeasurementEngineReadiness({
    required this.canRun,
    required this.modeLabel,
    required this.engineVersion,
    this.message,
  });

  final bool canRun;
  final String modeLabel;
  final String engineVersion;
  final String? message;
}

class MeasurementGoldenPathReadinessException implements Exception {
  const MeasurementGoldenPathReadinessException(this.readiness);

  final MeasurementEngineReadiness readiness;

  @override
  String toString() {
    final message = readiness.message ?? 'Measurement engine is not ready.';
    return '$message mode=${readiness.modeLabel} version=${readiness.engineVersion}';
  }
}

class MeasurementGoldenPathTraceEvent {
  const MeasurementGoldenPathTraceEvent({
    required this.phase,
    required this.elapsedMs,
    required this.occurredAt,
    this.sessionId,
    this.cartridgeId,
    this.engineMode,
    this.primaryValue,
    this.unit,
    this.confidence,
    this.failureReason,
    this.diagnosticMessage,
  });

  final MeasurementGoldenPathPhase phase;
  final int elapsedMs;
  final DateTime occurredAt;
  final String? sessionId;
  final String? cartridgeId;
  final String? engineMode;
  final double? primaryValue;
  final String? unit;
  final double? confidence;
  final String? failureReason;
  final String? diagnosticMessage;

  bool get isFailure => phase == MeasurementGoldenPathPhase.failed;

  Map<String, Object?> toJson() => {
        'phase': phase.name,
        'elapsed_ms': elapsedMs,
        'occurred_at': occurredAt.toIso8601String(),
        'session_id': sessionId,
        'cartridge_id': cartridgeId,
        'engine_mode': engineMode,
        'primary_value': primaryValue,
        'unit': unit,
        'confidence': confidence,
        'failure_reason': failureReason,
        'diagnostic_message': diagnosticMessage,
      };

  Map<String, Object?> toRemoteObservabilityJson({
    String source = 'flutter',
    String route = '/measure',
  }) =>
      {
        'schema_version': 'measure_trace.v1',
        'source': source,
        'route': route,
        'phase': phase.name,
        'elapsed_ms': elapsedMs,
        'occurred_at': occurredAt.toIso8601String(),
        'session_id': sessionId,
        'cartridge_id': cartridgeId,
        'engine_mode': engineMode,
        'unit': unit,
        'confidence': confidence,
        'has_primary_value': primaryValue != null,
        'failure_reason': failureReason,
        'diagnostic_message': diagnosticMessage,
      };

  String toLogLine() {
    final parts = <String>[
      'phase=${phase.name}',
      'elapsed_ms=$elapsedMs',
      if (sessionId != null) 'session_id=$sessionId',
      if (cartridgeId != null) 'cartridge_id=$cartridgeId',
      if (engineMode != null) 'engine_mode=$engineMode',
      if (primaryValue != null) 'primary_value=$primaryValue',
      if (unit != null) 'unit=$unit',
      if (confidence != null) 'confidence=$confidence',
      if (failureReason != null) 'failure_reason=$failureReason',
      if (diagnosticMessage != null) 'diagnostic=$diagnosticMessage',
    ];
    return parts.join(' ');
  }
}

typedef MeasurementGoldenPathTraceSink = void Function(
  MeasurementGoldenPathTraceEvent event,
);

class MeasurementGoldenPathLogger {
  const MeasurementGoldenPathLogger._();

  static void log(MeasurementGoldenPathTraceEvent event) {
    final logger = AppLogger.instance;
    if (event.isFailure) {
      logger.warning(event.toLogLine(), tag: 'MeasureGoldenPath');
      return;
    }
    logger.info(event.toLogLine(), tag: 'MeasureGoldenPath');
  }
}

class MeasurementGoldenPathCompositeTraceSink {
  const MeasurementGoldenPathCompositeTraceSink(this._sinks);

  final List<MeasurementGoldenPathTraceSink> _sinks;

  void call(MeasurementGoldenPathTraceEvent event) {
    for (final sink in _sinks) {
      sink(event);
    }
  }
}

abstract class MeasurementEngine {
  Future<MeasurementEngineReadiness> checkReadiness();

  Future<CartridgeInfoDto> readCartridge();

  Future<MeasurementPipelineResult> runPipeline({
    required String deviceId,
    required String biomarker,
    required String unit,
  });
}

class RustMeasurementEngine implements MeasurementEngine {
  const RustMeasurementEngine();

  @override
  Future<MeasurementEngineReadiness> checkReadiness() async {
    await RustBridge.init();
    final diagnostics = RustBridge.diagnostics;
    return MeasurementEngineReadiness(
      canRun: diagnostics.canRunMeasurement,
      modeLabel: diagnostics.modeLabel,
      engineVersion: diagnostics.engineVersion,
      message: diagnostics.blockingReason ??
          (diagnostics.isStubFallback
              ? 'Rust native engine is unavailable; Dart fallback is active.'
              : 'Rust native engine is active.'),
    );
  }

  @override
  Future<CartridgeInfoDto> readCartridge() => RustBridge.nfcReadCartridge();

  @override
  Future<MeasurementPipelineResult> runPipeline({
    required String deviceId,
    required String biomarker,
    required String unit,
  }) {
    return RustBridge.runMeasurementPipeline(
      deviceId: deviceId,
      biomarker: biomarker,
      unit: unit,
    );
  }
}

class MeasurementGoldenPathResult {
  const MeasurementGoldenPathResult({
    required this.sessionId,
    required this.cartridge,
    required this.pipeline,
    required this.serverResult,
    required this.endResult,
  });

  final String sessionId;
  final CartridgeInfoDto cartridge;
  final MeasurementPipelineResult pipeline;
  final ProcessMeasurementResult serverResult;
  final EndSessionResult? endResult;
}

class MeasurementGoldenPathOrchestrator {
  MeasurementGoldenPathOrchestrator({
    required MeasurementRepository repository,
    MeasurementEngine engine = const RustMeasurementEngine(),
    void Function(MeasurementGoldenPathSnapshot snapshot)? onSnapshot,
    MeasurementGoldenPathTraceSink traceSink = MeasurementGoldenPathLogger.log,
  })  : _repository = repository,
        _engine = engine,
        _onSnapshot = onSnapshot,
        _traceSink = traceSink;

  final MeasurementRepository _repository;
  final MeasurementEngine _engine;
  final void Function(MeasurementGoldenPathSnapshot snapshot)? _onSnapshot;
  final MeasurementGoldenPathTraceSink _traceSink;

  Future<MeasurementGoldenPathResult> run({
    required String userId,
    required String deviceId,
    String biomarker = 'glucose',
    String unit = 'mg/dL',
  }) async {
    String? sessionId;
    final stopwatch = Stopwatch()..start();
    try {
      final readiness = await _engine.checkReadiness();
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.readinessChecked,
            engineMode: readiness.modeLabel,
            diagnosticMessage: readiness.message,
          ),
          stopwatch);
      if (!readiness.canRun) {
        throw MeasurementGoldenPathReadinessException(readiness);
      }

      final cartridge = await _engine.readCartridge();
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.cartridgeRead,
            cartridgeId: cartridge.cartridgeId,
          ),
          stopwatch);

      final start = await _repository.startSession(
        deviceId: deviceId,
        cartridgeId: cartridge.cartridgeId,
        userId: userId,
      );
      sessionId = start.sessionId;
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.sessionStarted,
            sessionId: sessionId,
            cartridgeId: cartridge.cartridgeId,
          ),
          stopwatch);

      final pipeline = await _engine.runPipeline(
        deviceId: deviceId,
        biomarker: biomarker,
        unit: unit,
      );
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.engineProcessed,
            sessionId: sessionId,
            cartridgeId: cartridge.cartridgeId,
            primaryValue: pipeline.measurement.primaryValue,
            unit: pipeline.measurement.unit,
            confidence: pipeline.measurement.confidence,
          ),
          stopwatch);

      final serverResult = await _repository.processMeasurement(
        _buildProcessRequest(
          sessionId: sessionId,
          userId: userId,
          deviceId: deviceId,
          cartridgeType: cartridge.cartridgeType,
          pipeline: pipeline,
        ),
      );
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.serverProcessed,
            sessionId: sessionId,
            cartridgeId: cartridge.cartridgeId,
            primaryValue: serverResult.primaryValue,
            unit: serverResult.unit,
            confidence: serverResult.confidence,
            evidenceStatus: serverResult.evidenceStatus,
            diagnosticReady: serverResult.diagnosticReady,
            evidenceGaps: serverResult.evidenceGaps,
          ),
          stopwatch);

      final end = await _repository.endSession(sessionId);
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.sessionEnded,
            sessionId: sessionId,
            cartridgeId: cartridge.cartridgeId,
            primaryValue: serverResult.primaryValue,
            unit: serverResult.unit,
            confidence: serverResult.confidence,
            evidenceStatus: serverResult.evidenceStatus,
            diagnosticReady: serverResult.diagnosticReady,
            evidenceGaps: serverResult.evidenceGaps,
          ),
          stopwatch);

      return MeasurementGoldenPathResult(
        sessionId: sessionId,
        cartridge: cartridge,
        pipeline: pipeline,
        serverResult: serverResult,
        endResult: end,
      );
    } catch (error) {
      _emit(
          MeasurementGoldenPathSnapshot(
            phase: MeasurementGoldenPathPhase.failed,
            sessionId: sessionId,
            failureReason: error.toString(),
          ),
          stopwatch);
      rethrow;
    }
  }

  ProcessMeasurementRequest _buildProcessRequest({
    required String sessionId,
    required String userId,
    required String deviceId,
    required String cartridgeType,
    required MeasurementPipelineResult pipeline,
  }) {
    final measurement = pipeline.measurement;
    final rawChannels = List<double>.generate(88, (index) {
      final wave = math.sin(index * 0.15) * 0.02;
      return measurement.primaryValue + wave;
    });
    final fingerprint = List<double>.generate(88, (index) {
      final normalized =
          rawChannels[index] / (measurement.primaryValue.abs() + 1);
      return normalized.clamp(-1.0, 1.0);
    });

    return ProcessMeasurementRequest(
      sessionId: sessionId,
      deviceId: deviceId,
      userId: userId,
      cartridgeType: cartridgeType,
      rawChannels: rawChannels,
      sDet: measurement.primaryValue,
      sRef: measurement.referenceValue,
      alpha: 0.98,
      sCorrected: measurement.differentialValue,
      primaryValue: measurement.primaryValue,
      unit: measurement.unit,
      confidence: measurement.confidence,
      fingerprintVector: fingerprint,
      tempC: 24.0,
      humidityPct: 45.0,
      batteryPct: 85,
    );
  }

  void _emit(MeasurementGoldenPathSnapshot snapshot, Stopwatch stopwatch) {
    _onSnapshot?.call(snapshot);
    _traceSink(MeasurementGoldenPathTraceEvent(
      phase: snapshot.phase,
      elapsedMs: stopwatch.elapsedMilliseconds,
      occurredAt: DateTime.now(),
      sessionId: snapshot.sessionId,
      cartridgeId: snapshot.cartridgeId,
      engineMode: snapshot.engineMode,
      primaryValue: snapshot.primaryValue,
      unit: snapshot.unit,
      confidence: snapshot.confidence,
      failureReason: snapshot.failureReason,
      diagnosticMessage: snapshot.diagnosticMessage,
    ));
  }
}
