import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/constants/app_constants.dart';

void main() {
  group('GrpcClientManager', () {
    test('기본 호스트/포트 사용', () {
      final manager = GrpcClientManager();
      expect(manager.host, AppConstants.grpcHost);
      expect(manager.authPort, AppConstants.grpcAuthPort);
      expect(manager.userPort, AppConstants.grpcUserPort);
      expect(manager.devicePort, AppConstants.grpcDevicePort);
      expect(manager.measurementPort, AppConstants.grpcMeasurementPort);
    });

    test('커스텀 호스트/포트', () {
      final manager = GrpcClientManager(
        host: '192.168.1.1',
        authPort: 50061,
        userPort: 50062,
      );
      expect(manager.host, '192.168.1.1');
      expect(manager.authPort, 50061);
      expect(manager.userPort, 50062);
      expect(manager.devicePort, AppConstants.grpcDevicePort);
    });

    test('authChannel 접근 시 채널 생성', () {
      final manager = GrpcClientManager();
      expect(manager.hasAuthChannel, false);
      final channel = manager.authChannel;
      expect(channel, isNotNull);
      expect(manager.hasAuthChannel, true);
    });

    test('shutdown 후 채널 null', () async {
      final manager = GrpcClientManager();
      manager.authChannel;
      expect(manager.hasAuthChannel, true);
      await manager.shutdown();
      expect(manager.hasAuthChannel, false);
    });
  });

  group('AuthInterceptor', () {
    test('tokenProvider null이면 메타데이터에 토큰 없음', () {
      final interceptor = AuthInterceptor(() => null);
      expect(interceptor.tokenProvider(), null);
    });

    test('tokenProvider 빈 문자열', () {
      final interceptor = AuthInterceptor(() => '');
      expect(interceptor.tokenProvider(), '');
    });

    test('tokenProvider 값 반환', () {
      final interceptor = AuthInterceptor(() => 'my-jwt-token');
      expect(interceptor.tokenProvider(), 'my-jwt-token');
    });
  });
}
