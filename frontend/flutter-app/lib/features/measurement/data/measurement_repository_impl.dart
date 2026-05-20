import 'dart:async' as async;

import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:grpc/grpc.dart';

typedef MeasurementStreamCall = async.Stream<MeasurementResult> Function(
  async.Stream<MeasurementData> request,
);

typedef MeasurementHistoryCall = Future<GetHistoryResponse> Function(
  GetHistoryRequest request,
);

/// gRPC MeasurementService를 사용하는 MeasurementRepository 구현체
class MeasurementRepositoryImpl implements MeasurementRepository {
  MeasurementRepositoryImpl(
    this._grpcManager, {
    required String? Function() accessTokenProvider,
    MeasurementStreamCall? streamMeasurement,
    MeasurementHistoryCall? getMeasurementHistory,
  })  : _streamMeasurementOverride = streamMeasurement,
        _getMeasurementHistoryOverride = getMeasurementHistory,
        _authInterceptor = AuthInterceptor(accessTokenProvider);

  final GrpcClientManager _grpcManager;
  final MeasurementStreamCall? _streamMeasurementOverride;
  final MeasurementHistoryCall? _getMeasurementHistoryOverride;
  final AuthInterceptor _authInterceptor;

  MeasurementServiceClient? _client;

  MeasurementServiceClient get _measurementClient {
    _client ??= MeasurementServiceClient(
      _grpcManager.measurementChannel,
      interceptors: [_authInterceptor],
    );
    return _client!;
  }

  @override
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  }) async {
    try {
      final res = await _measurementClient.startSession(
        StartSessionRequest()
          ..deviceId = deviceId
          ..cartridgeId = cartridgeId
          ..userId = userId,
      );
      return StartSessionResult(
        sessionId: res.sessionId,
        startedAt: null,
      );
    } on GrpcError {
      rethrow;
    }
  }

  @override
  Future<EndSessionResult?> endSession(String sessionId) async {
    try {
      final res = await _measurementClient.endSession(
        EndSessionRequest()..sessionId = sessionId,
      );
      return EndSessionResult(
        sessionId: res.sessionId,
        totalMeasurements: res.totalMeasurements,
        endedAt: null,
      );
    } on GrpcError {
      rethrow;
    }
  }

  @override
  Future<ProcessMeasurementResult> processMeasurement(
    ProcessMeasurementRequest request,
  ) async {
    try {
      final res = await _streamMeasurements(
        async.Stream.value(_toMeasurementData(request)),
      ).first;
      return ProcessMeasurementResult(
        sessionId: res.sessionId.isNotEmpty ? res.sessionId : request.sessionId,
        primaryValue: res.primaryValue,
        unit: res.unit.isNotEmpty ? res.unit : request.unit,
        confidence: res.confidence,
        processedAt: null,
        evidenceStatus:
            res.evidenceStatus.isNotEmpty ? res.evidenceStatus : 'unknown',
        diagnosticReady: res.diagnosticReady,
        evidenceGaps: List.unmodifiable(res.evidenceGaps),
      );
    } on GrpcError {
      rethrow;
    }
  }

  @override
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final request = GetHistoryRequest()
        ..userId = userId
        ..limit = limit
        ..offset = offset;
      final override = _getMeasurementHistoryOverride;
      final res = override != null
          ? await override(request)
          : await _measurementClient.getMeasurementHistory(request);
      return MeasurementHistoryResult(
        items: res.measurements
            .map(
              (m) => MeasurementHistoryItem(
                sessionId: m.sessionId,
                cartridgeType: m.cartridgeType,
                primaryValue: m.primaryValue,
                unit: m.unit,
                evidenceStatus:
                    m.evidenceStatus.isNotEmpty ? m.evidenceStatus : 'unknown',
                diagnosticReady: m.diagnosticReady,
                evidenceGaps: List.unmodifiable(m.evidenceGaps),
                measuredAt: null,
              ),
            )
            .toList(),
        totalCount: res.totalCount,
      );
    } on GrpcError {
      rethrow;
    }
  }

  async.Stream<MeasurementResult> _streamMeasurements(
    async.Stream<MeasurementData> frames,
  ) {
    final override = _streamMeasurementOverride;
    if (override != null) {
      return override(frames);
    }
    return _measurementClient.streamMeasurement(frames);
  }

  MeasurementData _toMeasurementData(ProcessMeasurementRequest request) {
    return MeasurementData(
      sessionId: request.sessionId,
      rawChannels: request.rawChannels,
      differential: DifferentialCorrection(
        sDet: request.sDet,
        sRef: request.sRef,
        alpha: request.alpha,
        sCorrected: request.sCorrected,
      ),
      envMeta: EnvironmentMeta(
        tempC: request.tempC,
        humidityPct: request.humidityPct,
      ),
    );
  }
}
