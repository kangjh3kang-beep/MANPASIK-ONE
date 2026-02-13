// AuthNotifier (Provider) 단위 테스트
// - FakeAuthRepository를 이용한 로그인/회원가입/게스트/로그아웃 시나리오

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import '../../helpers/fake_repositories.dart';

void main() {
  late FakeAuthRepository fakeRepo;
  late AuthNotifier notifier;

  setUp(() {
    fakeRepo = FakeAuthRepository();
    notifier = AuthNotifier(fakeRepo);
  });

  group('AuthNotifier 테스트', () {
    // 초기 상태 확인
    test('초기 상태는 미인증 상태여야 한다', () {
      expect(notifier.state.isAuthenticated, isFalse);
      expect(notifier.state.userId, isNull);
      expect(notifier.state.accessToken, isNull);
    });

    // 정상 로그인 성공 시나리오
    test('유효한 이메일/비밀번호로 로그인하면 인증 상태가 된다', () async {
      final result = await notifier.login('user@manpasik.com', 'Password1234');

      expect(result, isTrue);
      expect(notifier.state.isAuthenticated, isTrue);
      expect(notifier.state.userId, 'test-user-id');
      expect(notifier.state.email, 'user@manpasik.com');
      expect(notifier.state.accessToken, isNotNull);
      expect(notifier.state.refreshToken, isNotNull);
    });

    // 잘못된 비밀번호로 로그인 실패
    test('짧은 비밀번호로 로그인하면 실패해야 한다', () async {
      final result = await notifier.login('user@manpasik.com', 'short');

      expect(result, isFalse);
      expect(notifier.state.isAuthenticated, isFalse);
    });

    // 빈 이메일로 로그인 실패
    test('빈 이메일로 로그인하면 실패해야 한다', () async {
      final result = await notifier.login('', 'Password1234');

      expect(result, isFalse);
      expect(notifier.state.isAuthenticated, isFalse);
    });

    // 회원가입 성공 시나리오
    test('유효한 정보로 회원가입하면 인증 상태가 된다', () async {
      final result = await notifier.register(
        'new@manpasik.com',
        'Password1234',
        '새 사용자',
      );

      expect(result, isTrue);
      expect(notifier.state.isAuthenticated, isTrue);
      expect(notifier.state.email, 'new@manpasik.com');
      expect(notifier.state.displayName, '새 사용자');
    });

    // 회원가입 실패 시나리오 (빈 displayName)
    test('빈 이름으로 회원가입하면 실패해야 한다', () async {
      final result = await notifier.register(
        'new@manpasik.com',
        'Password1234',
        '', // 빈 이름
      );

      expect(result, isFalse);
      expect(notifier.state.isAuthenticated, isFalse);
    });

    // 게스트 로그인
    test('게스트 로그인 시 guest 정보로 인증 상태가 된다', () {
      notifier.loginAsGuest();

      expect(notifier.state.isAuthenticated, isTrue);
      expect(notifier.state.userId, 'guest-user');
      expect(notifier.state.email, 'guest@example.com');
      expect(notifier.state.displayName, 'Guest');
      expect(notifier.state.accessToken, 'guest-token');
    });

    // 로그아웃
    test('로그아웃 시 미인증 상태로 초기화된다', () async {
      // 먼저 로그인
      await notifier.login('user@manpasik.com', 'Password1234');
      expect(notifier.state.isAuthenticated, isTrue);

      // 로그아웃
      notifier.logout();
      expect(notifier.state.isAuthenticated, isFalse);
      expect(notifier.state.userId, isNull);
      expect(notifier.state.accessToken, isNull);
    });

    // 로그인 → 게스트 → 로그아웃 순차 시나리오
    test('로그인 후 로그아웃하고 게스트 로그인이 정상 동작해야 한다', () async {
      await notifier.login('user@manpasik.com', 'Password1234');
      expect(notifier.state.email, 'user@manpasik.com');

      notifier.logout();
      expect(notifier.state.isAuthenticated, isFalse);

      notifier.loginAsGuest();
      expect(notifier.state.userId, 'guest-user');
    });

    // checkAuthStatus 호출 시 미인증 상태 확인
    // (FakeAuthRepository.isAuthenticated()는 항상 false 반환)
    test('checkAuthStatus 호출 시 미인증으로 초기화된다', () async {
      // 먼저 게스트로 로그인
      notifier.loginAsGuest();
      expect(notifier.state.isAuthenticated, isTrue);

      await notifier.checkAuthStatus();
      // FakeAuthRepository.isAuthenticated()가 false이므로 초기화됨
      expect(notifier.state.isAuthenticated, isFalse);
    });
  });
}
