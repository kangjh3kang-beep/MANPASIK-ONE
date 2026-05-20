import 'package:flutter_test/flutter_test.dart';

import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';

/// Phase M: 도메인 모델 추가 테스트
void main() {
  group('AuthResult', () {
    test('성공 결과는 isSuccess가 true', () {
      final result = AuthResult.success(
        userId: 'u1',
        email: 'test@example.com',
        displayName: 'Test',
        accessToken: 'a',
        refreshToken: 'r',
      );
      expect(result.success, true);
      expect(result.userId, 'u1');
      expect(result.email, 'test@example.com');
    });

    test('실패 결과는 isSuccess가 false', () {
      final result = AuthResult.failure('Invalid credentials');
      expect(result.success, false);
      expect(result.errorMessage, 'Invalid credentials');
    });

    test('성공 결과는 errorMessage가 null', () {
      final result = AuthResult.success(
        userId: 'u', email: 'e', displayName: 'd',
        accessToken: 'a', refreshToken: 'r',
      );
      expect(result.errorMessage, null);
    });

    test('userId가 비어있어도 객체는 생성됨', () {
      final result = AuthResult.success(
        userId: '', email: 'e', displayName: 'd',
        accessToken: 'a', refreshToken: 'r',
      );
      expect(result.success, true);
      expect(result.userId, '');
    });
  });

  group('DeviceItem', () {
    test('생성자와 필드 검증', () {
      const device = DeviceItem(
        deviceId: 'd1',
        name: '리더기',
        firmwareVersion: '1.0.0',
        status: 'online',
        batteryPercent: 85,
      );
      expect(device.deviceId, 'd1');
      expect(device.name, '리더기');
      expect(device.batteryPercent, 85);
      expect(device.status, 'online');
    });

    test('배터리 0%도 유효', () {
      const device = DeviceItem(
        deviceId: 'low',
        name: '저전력',
        firmwareVersion: '1.0',
        status: 'online',
        batteryPercent: 0,
      );
      expect(device.batteryPercent, 0);
    });

    test('상태 종류 비교', () {
      const online = DeviceItem(
        deviceId: 'o', name: 'Online',
        firmwareVersion: '1.0', status: 'online', batteryPercent: 100,
      );
      const offline = DeviceItem(
        deviceId: 'f', name: 'Offline',
        firmwareVersion: '1.0', status: 'offline', batteryPercent: 50,
      );
      expect(online.status != offline.status, true);
    });
  });

  group('AuthResult 부가', () {
    test('role 필드 설정 가능', () {
      const result = AuthResult(
        success: true, userId: 'u', email: 'e@x.com',
        displayName: 'd', accessToken: 'a', refreshToken: 'r',
        role: 'admin',
      );
      expect(result.role, 'admin');
    });

    test('errorMessage가 nullable', () {
      const result = AuthResult(success: true);
      expect(result.errorMessage, null);
    });
  });

  group('FakeAuthRepository', () {
    final repo = FakeAuthRepository();

    test('정상 로그인', () async {
      final result = await repo.login('user@example.com', 'password123');
      expect(result.success, true);
      expect(result.email, 'user@example.com');
      expect(result.displayName, 'user');
    });

    test('빈 이메일 로그인 실패', () async {
      final result = await repo.login('', 'password123');
      expect(result.success, false);
    });

    test('짧은 비밀번호 로그인 실패', () async {
      final result = await repo.login('user@example.com', 'short');
      expect(result.success, false);
    });

    test('정상 회원가입', () async {
      final result = await repo.register('new@example.com', 'password123', '신규유저');
      expect(result.success, true);
      expect(result.displayName, '신규유저');
    });

    test('소셜 로그인은 항상 성공', () async {
      final result = await repo.socialLogin('kakao', 'token-123');
      expect(result.success, true);
      expect(result.email, 'social@kakao.com');
    });

    test('logout은 예외 없이 완료', () async {
      await expectLater(repo.logout(), completes);
    });

    test('refreshToken은 false 반환', () async {
      expect(await repo.refreshToken(), false);
    });

    test('isAuthenticated는 false 반환', () async {
      expect(await repo.isAuthenticated(), false);
    });
  });

  group('FakeDeviceRepository', () {
    final repo = FakeDeviceRepository();

    test('정상 사용자에게 디바이스 1개 반환', () async {
      final devices = await repo.listDevices('test-user');
      expect(devices.length, 1);
      expect(devices.first.deviceId, 'device-1');
    });

    test('빈 사용자ID는 빈 리스트', () async {
      final devices = await repo.listDevices('');
      expect(devices, isEmpty);
    });

    test('연결된 디바이스 빈 리스트 (기본)', () async {
      final connected = await repo.getConnectedDevices();
      expect(connected, isEmpty);
    });
  });
}

/// 헬퍼 클래스 - test/helpers의 클래스 이용
class FakeAuthRepository implements AuthRepository {
  @override
  Future<AuthResult> login(String email, String password) async {
    if (email.isEmpty || password.length < 8) {
      return AuthResult.failure('Invalid credentials');
    }
    return AuthResult.success(
      userId: 'test-user-id',
      email: email,
      displayName: email.split('@').first,
      accessToken: 'fake-access-token',
      refreshToken: 'fake-refresh-token',
    );
  }

  @override
  Future<AuthResult> register(String email, String password, String displayName) async {
    if (email.isEmpty || password.length < 8 || displayName.isEmpty) {
      return AuthResult.failure('Invalid input');
    }
    return AuthResult.success(
      userId: 'test-user-id', email: email, displayName: displayName,
      accessToken: 'fake', refreshToken: 'fake',
    );
  }

  @override
  Future<void> logout() async {}

  @override
  Future<bool> refreshToken() async => false;

  @override
  Future<bool> isAuthenticated() async => false;

  @override
  Future<AuthResult> socialLogin(String provider, String token) async {
    return AuthResult.success(
      userId: 'test-social-user',
      email: 'social@$provider.com',
      displayName: '$provider 사용자',
      accessToken: 'fake', refreshToken: 'fake',
    );
  }
}

class FakeDeviceRepository implements DeviceRepository {
  @override
  Future<List<DeviceItem>> listDevices(String userId) async {
    if (userId.isEmpty) return [];
    return [
      const DeviceItem(
        deviceId: 'device-1',
        name: '테스트 리더기',
        firmwareVersion: '1.0.0',
        status: 'online',
        batteryPercent: 85,
      ),
    ];
  }

  @override
  Future<List<ConnectedDevice>> getConnectedDevices() async => [];

  @override
  Future<bool> registerDevice(String userId, DeviceItem device) async => true;

  @override
  Future<bool> unregisterDevice(String deviceId) async => true;

  @override
  Future<DeviceItem?> getDevice(String deviceId) async => null;

  @override
  Future<bool> updateFirmware(String deviceId, String version) async => true;
}
