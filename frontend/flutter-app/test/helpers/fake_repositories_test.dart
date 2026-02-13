// FakeRepository 테스트
// - 기존 test/helpers/fake_repositories.dart의 동작 검증
// - 다른 테스트에서 사용하는 Fake 구현체가 올바르게 동작하는지 확인

import 'package:flutter_test/flutter_test.dart';
import 'fake_repositories.dart';

void main() {
  group('FakeAuthRepository 테스트', () {
    late FakeAuthRepository repo;

    setUp(() {
      repo = FakeAuthRepository();
    });

    // 유효한 로그인
    test('유효한 이메일/비밀번호로 로그인 시 성공한다', () async {
      final result = await repo.login('user@manpasik.com', 'Password1234');
      expect(result.success, isTrue);
      expect(result.userId, 'test-user-id');
      expect(result.accessToken, isNotNull);
    });

    // 빈 이메일로 로그인 실패
    test('빈 이메일로 로그인 시 실패한다', () async {
      final result = await repo.login('', 'Password1234');
      expect(result.success, isFalse);
    });

    // 짧은 비밀번호로 로그인 실패
    test('8자 미만 비밀번호로 로그인 시 실패한다', () async {
      final result = await repo.login('user@manpasik.com', 'short');
      expect(result.success, isFalse);
    });

    // isAuthenticated 기본값 false
    test('isAuthenticated는 항상 false를 반환한다', () async {
      expect(await repo.isAuthenticated(), isFalse);
    });
  });

  group('FakeDeviceRepository 테스트', () {
    late FakeDeviceRepository repo;

    setUp(() {
      repo = FakeDeviceRepository();
    });

    // 유효한 userId로 디바이스 목록 조회
    test('유효한 userId로 디바이스 목록을 조회할 수 있다', () async {
      final devices = await repo.listDevices('user-1');
      expect(devices, isNotEmpty);
      expect(devices.first.deviceId, 'device-1');
      expect(devices.first.status, 'online');
    });

    // 빈 userId는 빈 목록 반환
    test('빈 userId로 조회 시 빈 목록을 반환한다', () async {
      final devices = await repo.listDevices('');
      expect(devices, isEmpty);
    });
  });

  group('FakeMeasurementRepository 테스트', () {
    late FakeMeasurementRepository repo;

    setUp(() {
      repo = FakeMeasurementRepository();
    });

    // startSession
    test('startSession이 session ID를 반환한다', () async {
      final result = await repo.startSession(
        deviceId: 'dev-1',
        cartridgeId: 'cart-1',
        userId: 'user-1',
      );
      expect(result.sessionId, 'session-1');
    });

    // endSession
    test('endSession이 측정 수를 반환한다', () async {
      final result = await repo.endSession('session-1');
      expect(result, isNotNull);
      expect(result!.totalMeasurements, 5);
    });

    // getHistory - 유효한 userId
    test('유효한 userId로 기록을 조회할 수 있다', () async {
      final result = await repo.getHistory(userId: 'user-1');
      expect(result.items, isNotEmpty);
      expect(result.totalCount, 1);
      expect(result.items.first.primaryValue, 98.4);
    });

    // getHistory - 빈 userId
    test('빈 userId로 조회 시 빈 결과를 반환한다', () async {
      final result = await repo.getHistory(userId: '');
      expect(result.items, isEmpty);
      expect(result.totalCount, 0);
    });
  });

  group('FakeUserRepository 테스트', () {
    late FakeUserRepository repo;

    setUp(() {
      repo = FakeUserRepository();
    });

    // getProfile - 유효한 userId
    test('유효한 userId로 프로필을 조회할 수 있다', () async {
      final profile = await repo.getProfile('user-1');
      expect(profile, isNotNull);
      expect(profile!.userId, 'user-1');
      expect(profile.email, 'test@manpasik.com');
      expect(profile.displayName, '테스트 사용자');
    });

    // getProfile - 빈 userId
    test('빈 userId로 프로필 조회 시 null을 반환한다', () async {
      final profile = await repo.getProfile('');
      expect(profile, isNull);
    });

    // getSubscription
    test('유효한 userId로 구독 정보를 조회할 수 있다', () async {
      final sub = await repo.getSubscription('user-1');
      expect(sub, isNotNull);
      expect(sub!.tier, 0);
      expect(sub.maxDevices, 1);
    });

    // getSubscription - 빈 userId
    test('빈 userId로 구독 정보 조회 시 null을 반환한다', () async {
      final sub = await repo.getSubscription('');
      expect(sub, isNull);
    });
  });
}
