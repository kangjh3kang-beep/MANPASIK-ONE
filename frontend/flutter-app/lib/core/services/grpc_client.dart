/// gRPC 채널 관리자
///
/// 호스트/포트 설정 및 Auth/User/Device/Measurement 서비스별 채널 제공.
/// AppConstants.grpcHost, grpcAuthPort 등과 연동.
import 'package:grpc/grpc.dart';
import 'package:manpasik/core/constants/app_constants.dart';

/// gRPC 채널 및 연결 상태 관리
class GrpcClientManager {
  GrpcClientManager({
    String? host,
    int? authPort,
    int? userPort,
    int? devicePort,
    int? measurementPort,
    int? adminPort,
    int? aiInferencePort,
    bool useInsecure = true,
  })  : _host = host ?? AppConstants.grpcHost,
        _authPort = authPort ?? AppConstants.grpcAuthPort,
        _userPort = userPort ?? AppConstants.grpcUserPort,
        _devicePort = devicePort ?? AppConstants.grpcDevicePort,
        _measurementPort = measurementPort ?? AppConstants.grpcMeasurementPort,
        _adminPort = adminPort ?? AppConstants.grpcAdminPort,
        _aiInferencePort = aiInferencePort ?? AppConstants.grpcAiInferencePort,
        _useInsecure = useInsecure;

  final String _host;
  final int _authPort;
  final int _userPort;
  final int _devicePort;
  final int _measurementPort;
  final int _adminPort;
  final int _aiInferencePort;
  final bool _useInsecure;

  ClientChannel? _authChannel;
  ClientChannel? _userChannel;
  ClientChannel? _deviceChannel;
  ClientChannel? _measurementChannel;
  ClientChannel? _adminChannel;
  ClientChannel? _aiInferenceChannel;

  String get host => _host;
  int get authPort => _authPort;
  int get userPort => _userPort;
  int get devicePort => _devicePort;
  int get measurementPort => _measurementPort;
  int get adminPort => _adminPort;
  int get aiInferencePort => _aiInferencePort;

  ChannelOptions get _channelOptions => ChannelOptions(
        credentials: ChannelCredentials.insecure(),
      );

  /// Auth 서비스 채널 (50051). 인증 전용이라 인터셉터 없음.
  ClientChannel get authChannel {
    _authChannel ??= ClientChannel(
      _host,
      port: _authPort,
      options: _channelOptions,
    );
    return _authChannel!;
  }

  /// User 서비스 채널 (50052). JWT 인터셉터와 함께 사용.
  ClientChannel get userChannel {
    _userChannel ??= ClientChannel(
      _host,
      port: _userPort,
      options: _channelOptions,
    );
    return _userChannel!;
  }

  /// Device 서비스 채널 (50053)
  ClientChannel get deviceChannel {
    _deviceChannel ??= ClientChannel(
      _host,
      port: _devicePort,
      options: _channelOptions,
    );
    return _deviceChannel!;
  }

  /// Measurement 서비스 채널 (50054)
  ClientChannel get measurementChannel {
    _measurementChannel ??= ClientChannel(
      _host,
      port: _measurementPort,
      options: _channelOptions,
    );
    return _measurementChannel!;
  }

  /// Admin 서비스 채널 (50055)
  ClientChannel get adminChannel {
    _adminChannel ??= ClientChannel(
      _host,
      port: _adminPort,
      options: _channelOptions,
    );
    return _adminChannel!;
  }

  /// AI Inference 서비스 채널 (50058)
  ClientChannel get aiInferenceChannel {
    _aiInferenceChannel ??= ClientChannel(
      _host,
      port: _aiInferencePort,
      options: _channelOptions,
    );
    return _aiInferenceChannel!;
  }

  /// 연결 상태 확인 (채널이 생성되었는지)
  bool get hasAuthChannel => _authChannel != null;
  bool get hasUserChannel => _userChannel != null;
  bool get hasDeviceChannel => _deviceChannel != null;
  bool get hasMeasurementChannel => _measurementChannel != null;
  bool get hasAdminChannel => _adminChannel != null;
  bool get hasAiInferenceChannel => _aiInferenceChannel != null;

  /// 채널 종료 (앱 종료 시 호출)
  Future<void> shutdown() async {
    await _authChannel?.shutdown();
    await _userChannel?.shutdown();
    await _deviceChannel?.shutdown();
    await _measurementChannel?.shutdown();
    await _adminChannel?.shutdown();
    await _aiInferenceChannel?.shutdown();
    _authChannel = null;
    _userChannel = null;
    _deviceChannel = null;
    _measurementChannel = null;
    _adminChannel = null;
    _aiInferenceChannel = null;
  }
}
