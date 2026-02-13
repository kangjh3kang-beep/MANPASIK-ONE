import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:grpc/grpc.dart';

/// gRPC MeasurementService를 사용하는 MeasurementRepository 구현체
class MeasurementRepositoryImpl implements MeasurementRepository {
  MeasurementRepositoryImpl(
    this._grpcManager, {
    required String? Function() accessTokenProvider,
  }) : _authInterceptor = AuthInterceptor(accessTokenProvider);

  final GrpcClientManager _grpcManager;
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
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final res = await _measurementClient.getMeasurementHistory(
        GetHistoryRequest()
          ..userId = userId
          ..limit = limit
          ..offset = offset,
      );
      return MeasurementHistoryResult(
        items: res.measurements
            .map(
              (m) => MeasurementHistoryItem(
                sessionId: m.sessionId,
                cartridgeType: m.cartridgeType,
                primaryValue: m.primaryValue,
                unit: m.unit,
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
}
