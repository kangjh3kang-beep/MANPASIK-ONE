// AppLocalizations 다국어 테스트
// - 한국어/영어 번역 정확성, 지원 로케일 목록, delegate, greetingWithName

import 'dart:ui';

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/l10n/app_localizations.dart';

void main() {
  group('AppLocalizations 테스트', () {
    // 한국어(기본) 번역 확인
    test('한국어 로케일에서 appName이 "만파식"이어야 한다', () {
      final l10n = AppLocalizations(const Locale('ko'));
      expect(l10n.appName, '만파식');
      expect(l10n.login, '로그인');
      expect(l10n.register, '회원가입');
      expect(l10n.logout, '로그아웃');
    });

    // 영어 번역 확인
    test('영어 로케일에서 appName이 "ManPaSik"이어야 한다', () {
      final l10n = AppLocalizations(const Locale('en'));
      expect(l10n.appName, 'ManPaSik');
      expect(l10n.login, 'Login');
      expect(l10n.register, 'Sign Up');
      expect(l10n.logout, 'Logout');
    });

    // 일본어 로케일 로드 가능 여부
    test('일본어 로케일이 정상적으로 로드된다', () {
      final l10n = AppLocalizations(const Locale('ja'));
      // ja 번역이 존재하면 appName은 null이 아니어야 함
      expect(l10n.appName, isNotNull);
      expect(l10n.appName, isNotEmpty);
    });

    // 지원 로케일 수 확인 (6개)
    test('supportedLocales는 6개 언어를 포함해야 한다', () {
      expect(AppLocalizations.supportedLocales.length, 6);
      final codes =
          AppLocalizations.supportedLocales.map((l) => l.languageCode).toList();
      expect(codes, containsAll(['ko', 'en', 'ja', 'zh', 'fr', 'hi']));
    });

    // 미지원 로케일은 한국어(기본)로 fallback
    test('미지원 로케일은 한국어로 fallback된다', () {
      final l10n = AppLocalizations(const Locale('de')); // 독일어 (미지원)
      expect(l10n.appName, '만파식'); // 기본 한국어
    });

    // greetingWithName 치환 확인
    test('greetingWithName이 이름을 올바르게 치환한다', () {
      final l10nKo = AppLocalizations(const Locale('ko'));
      expect(l10nKo.greetingWithName('홍길동'), contains('홍길동'));

      final l10nEn = AppLocalizations(const Locale('en'));
      expect(l10nEn.greetingWithName('John'), contains('John'));
    });

    // delegate isSupported 확인
    test('delegate.isSupported는 지원 언어에 대해 true를 반환한다', () {
      const delegate = AppLocalizations.delegate;
      expect(delegate.isSupported(const Locale('ko')), isTrue);
      expect(delegate.isSupported(const Locale('en')), isTrue);
      expect(delegate.isSupported(const Locale('de')), isFalse);
    });

    // 한국어 검증 관련 번역 키 확인
    test('한국어 검증 메시지 번역이 존재해야 한다', () {
      final l10n = AppLocalizations(const Locale('ko'));
      expect(l10n.validationEmailRequired, isNotEmpty);
      expect(l10n.validationEmailInvalid, isNotEmpty);
      expect(l10n.validationPasswordRequired, isNotEmpty);
      expect(l10n.validationPasswordTooShort, isNotEmpty);
      expect(l10n.validationNameRequired, isNotEmpty);
    });
  });
}
