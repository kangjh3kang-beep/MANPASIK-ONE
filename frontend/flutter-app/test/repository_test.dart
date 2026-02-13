import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/test/helpers/fake_repositories.dart';

void main() {
  group('FakeAuthRepository', () {
    late FakeAuthRepository repo;

    setUp(() => repo = FakeAuthRepository());

    test('login 성공 - 유효한 이메일/비밀번호', () async {
      final result = await repo.login('a@b.com', 'password1');
      expect(result.success, true);
      expect(result.accessToken, 'fake-access-token');
      expect(result.userId, 'test-user-id');
    });

    test('login 실패 - 빈 이메일', () async {
      final result = await repo.login('', 'password12');
      expect(result.success, false);
      expect(result.errorMessage, isNotNull);
    });

    test('login 실패 - 짧은 비밀번호', () async {
      final result = await repo.login('a@b.com', 'short');
      expect(result.success, false);
    });

    test('register 성공', () async {
      final result = await repo.register('x@y.com', 'pass1234', 'Display');
      expect(result.success, true);
      expect(result.displayName, 'Display');
    });

    test('register 실패 - 빈 displayName', () async {
      final result = await repo.register('x@y.com', 'pass1234', '');
      expect(result.success, false);
    });

    test('logout 완료', () async {
      await expectLater(repo.logout(), completes);
    });

    test('isAuthenticated 항상 false', () async {
      expect(await repo.isAuthenticated(), false);
    });
  });

  group('FakeDeviceRepository', () {
    late FakeDeviceRepository repo;

    setUp(() => repo = FakeDeviceRepository());

    test('listDevices - 빈 userId면 빈 목록', () async {
      final list = await repo.listDevices('');
      expect(list, isEmpty);
    });

    test('listDevices - userId 있으면 1개 반환', () async {
      final list = await repo.listDevices('user-1');
      expect(list.length, 1);
      expect(list.first.deviceId, 'device-1');
      expect(list.first.name, '테스트 리더기');
      expect(list.first.status, 'online');
    });
  });

  group('FakeMeasurementRepository', () {
    late FakeMeasurementRepository repo;

    setUp(() => repo = FakeMeasurementRepository());

    test('startSession 성공', () async {
      final result = await repo.startSession(
        deviceId: 'd1',
        cartridgeId: 'c1',
        userId: 'u1',
      );
      expect(result.sessionId, 'session-1');
    });

    test('endSession 성공', () async {
      final result = await repo.endSession('session-1');
      expect(result, isNotNull);
      expect(result!.totalMeasurements, 5);
    });

    test('getHistory - 빈 userId면 빈 결과', () async {
      final result = await repo.getHistory(userId: '');
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
    });

    test('getHistory - userId 있으면 1개 항목', () async {
      final result = await repo.getHistory(userId: 'u1');
      expect(result.totalCount, 1);
      expect(result.items.first.primaryValue, 98.4);
      expect(result.items.first.unit, 'mg/dL');
    });
  });

  group('FakeUserRepository', () {
    late FakeUserRepository repo;

    setUp(() => repo = FakeUserRepository());

    test('getProfile - 빈 userId면 null', () async {
      expect(await repo.getProfile(''), null);
    });

    test('getProfile - userId 있으면 프로필 반환', () async {
      final p = await repo.getProfile('u1');
      expect(p, isNotNull);
      expect(p!.email, 'test@manpasik.com');
      expect(p.displayName, '테스트 사용자');
    });

    test('getSubscription - 빈 userId면 null', () async {
      expect(await repo.getSubscription(''), null);
    });

    test('getSubscription - userId 있으면 구독 정보 반환', () async {
      final s = await repo.getSubscription('u1');
      expect(s, isNotNull);
      expect(s!.maxDevices, 1);
      expect(s.tier, 0);
    });
  });
}
