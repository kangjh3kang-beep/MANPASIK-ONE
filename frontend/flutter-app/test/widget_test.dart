import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/core/utils/validators.dart';
import 'helpers/fake_repositories.dart';

void main() {
  late ProviderContainer authContainer;

  setUp(() {
    SharedPreferences.setMockInitialValues({});
    authContainer = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
      ],
    );
  });

  tearDown(() {
    authContainer.dispose();
  });

  group('AuthProvider 테스트', () {
    test('초기 상태는 비인증', () {
      final notifier = authContainer.read(authProvider.notifier);
      expect(notifier.state.isAuthenticated, false);
      expect(notifier.state.userId, null);
    });

    test('유효한 이메일/비밀번호로 로그인 성공', () async {
      final notifier = authContainer.read(authProvider.notifier);
      final result = await notifier.login('test@manpasik.com', 'password123');
      expect(result, true);
      expect(notifier.state.isAuthenticated, true);
      expect(notifier.state.email, 'test@manpasik.com');
    });

    test('빈 이메일로 로그인 실패', () async {
      final notifier = authContainer.read(authProvider.notifier);
      final result = await notifier.login('', 'password123');
      expect(result, false);
      expect(notifier.state.isAuthenticated, false);
    });

    test('짧은 비밀번호로 로그인 실패', () async {
      final notifier = authContainer.read(authProvider.notifier);
      final result = await notifier.login('test@manpasik.com', 'short');
      expect(result, false);
    });

    test('로그아웃 시 상태 초기화', () async {
      final notifier = authContainer.read(authProvider.notifier);
      await notifier.login('test@manpasik.com', 'password123');
      expect(notifier.state.isAuthenticated, true);

      notifier.logout();
      expect(notifier.state.isAuthenticated, false);
      expect(notifier.state.userId, null);
    });

    test('회원가입 성공', () async {
      final notifier = authContainer.read(authProvider.notifier);
      final result = await notifier.register('new@manpasik.com', 'password123', '테스트');
      expect(result, true);
      expect(notifier.state.isAuthenticated, true);
      expect(notifier.state.displayName, '테스트');
    });

    test('빈 이름으로 회원가입 실패', () async {
      final notifier = authContainer.read(authProvider.notifier);
      final result = await notifier.register('new@manpasik.com', 'password123', '');
      expect(result, false);
    });
  });

  group('ThemeModeProvider 테스트', () {
    test('기본 테마는 dark', () {
      final notifier = ThemeModeNotifier();
      expect(notifier.state, ThemeMode.dark);
    });

    test('라이트 모드 전환', () {
      final notifier = ThemeModeNotifier();
      notifier.setLight();
      expect(notifier.state, ThemeMode.light);
    });

    test('다크 모드 전환', () {
      final notifier = ThemeModeNotifier();
      notifier.setDark();
      expect(notifier.state, ThemeMode.dark);
    });

    test('토글 순환 (dark → system → light → dark)', () {
      final notifier = ThemeModeNotifier();
      expect(notifier.state, ThemeMode.dark);

      notifier.toggle();
      expect(notifier.state, ThemeMode.system);

      notifier.toggle();
      expect(notifier.state, ThemeMode.light);

      notifier.toggle();
      expect(notifier.state, ThemeMode.dark);
    });
  });

  group('LocaleProvider 테스트', () {
    test('기본 로케일은 한국어', () {
      final notifier = LocaleNotifier();
      expect(notifier.state.languageCode, 'ko');
    });

    test('영어로 변경', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('en');
      expect(notifier.state.languageCode, 'en');
    });

    test('일본어로 변경', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('ja');
      expect(notifier.state.languageCode, 'ja');
    });

    test('중국어로 변경', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('zh');
      expect(notifier.state.languageCode, 'zh');
    });

    test('프랑스어로 변경', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('fr');
      expect(notifier.state.languageCode, 'fr');
    });

    test('힌디어로 변경', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('hi');
      expect(notifier.state.languageCode, 'hi');
    });

    test('미지원 언어는 무시', () {
      final notifier = LocaleNotifier();
      notifier.setLocaleByCode('ar'); // 아직 미지원
      expect(notifier.state.languageCode, 'ko'); // 변경 안됨
    });

    test('SupportedLocales 6개 언어 확인', () {
      expect(SupportedLocales.all.length, 6);
      final codes = SupportedLocales.all.map((l) => l.languageCode).toList();
      expect(codes, contains('ko'));
      expect(codes, contains('en'));
      expect(codes, contains('ja'));
      expect(codes, contains('zh'));
      expect(codes, contains('fr'));
      expect(codes, contains('hi'));
    });

    test('언어 이름 맵 확인', () {
      expect(SupportedLocales.getLanguageName('ko'), '한국어');
      expect(SupportedLocales.getLanguageName('en'), 'English');
      expect(SupportedLocales.getLanguageName('ja'), '日本語');
      expect(SupportedLocales.getLanguageName('zh'), '中文简体');
      expect(SupportedLocales.getLanguageName('fr'), 'Français');
      expect(SupportedLocales.getLanguageName('hi'), 'हिन्दी');
    });
  });

  group('Validators 테스트', () {
    test('이메일 검증 - 빈 값', () {
      expect(Validators.validateEmail(''), isNotNull);
      expect(Validators.validateEmail(null), isNotNull);
    });

    test('이메일 검증 - 유효하지 않은 형식', () {
      expect(Validators.validateEmail('notanemail'), isNotNull);
      expect(Validators.validateEmail('missing@'), isNotNull);
    });

    test('이메일 검증 - 유효한 형식', () {
      expect(Validators.validateEmail('test@example.com'), null);
      expect(Validators.validateEmail('user@manpasik.com'), null);
    });

    test('비밀번호 검증 - 빈 값', () {
      expect(Validators.validatePassword(''), isNotNull);
    });

    test('비밀번호 검증 - 너무 짧음', () {
      expect(Validators.validatePassword('abc'), isNotNull);
    });

    test('비밀번호 검증 - 숫자 없음', () {
      expect(Validators.validatePassword('abcdefgh'), isNotNull);
    });

    test('비밀번호 검증 - 영문 없음', () {
      expect(Validators.validatePassword('12345678'), isNotNull);
    });

    test('비밀번호 검증 - 유효', () {
      expect(Validators.validatePassword('abc12345'), null);
    });

    test('이름 검증 - 빈 값', () {
      expect(Validators.validateDisplayName(''), isNotNull);
    });

    test('이름 검증 - 너무 짧음', () {
      expect(Validators.validateDisplayName('a'), isNotNull);
    });

    test('이름 검증 - 유효', () {
      expect(Validators.validateDisplayName('홍길동'), null);
    });
  });
}
