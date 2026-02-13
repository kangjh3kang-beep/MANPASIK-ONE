import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:grpc/grpc.dart';

/// gRPC DeviceService를 사용하는 DeviceRepository 구현체
class DeviceRepositoryImpl implements DeviceRepository {
  DeviceRepositoryImpl(
    this._grpcManager, {
    required String? Function() accessTokenProvider,
  }) : _authInterceptor = AuthInterceptor(accessTokenProvider);

  final GrpcClientManager _grpcManager;
  final AuthInterceptor _authInterceptor;

  DeviceServiceClient? _client;

  DeviceServiceClient get _deviceClient {
    _client ??= DeviceServiceClient(
      _grpcManager.deviceChannel,
      interceptors: [_authInterceptor],
    );
    return _client!;
  }

  static String _statusName(int status) {
    switch (status) {
      case 1:
        return 'online';
      case 2:
        return 'offline';
      case 3:
        return 'measuring';
      case 4:
        return 'updating';
      case 5:
        return 'error';
      default:
        return 'unknown';
    }
  }

  @override
  Future<List<DeviceItem>> listDevices(String userId) async {
    try {
      final res = await _deviceClient.listDevices(
        ListDevicesRequest()..userId = userId,
      );
      return res.devices
          .map(
            (d) => DeviceItem(
              deviceId: d.deviceId,
              name: d.name.isNotEmpty ? d.name : d.deviceId,
              firmwareVersion: d.firmwareVersion,
              status: _statusName(d.status),
              batteryPercent: d.batteryPercent,
            ),
          )
          .toList();
    } on GrpcError {
      rethrow;
    }
  }
}
