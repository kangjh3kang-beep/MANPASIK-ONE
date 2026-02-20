// LocaleNotifier Provider 테스트
// - 기본값, setLocale, setLocaleByCode, 비지원 언어 무시

import 'dart:ui';

import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';

void main() {
  late LocaleNotifier notifier;

  setUp(() {
    SharedPreferences.setMockInitialValues({});
    notifier = LocaleNotifier();
  });

  group('LocaleNotifier 테스트', () {
    // 초기 로케일은 한국어
    test('초기 로케일은 ko(한국어)이어야 한다', () {
      expect(notifier.state.languageCode, 'ko');
    });

    // setLocale로 영어 설정
    test('setLocale로 영어 로케일을 설정할 수 있다', () {
      notifier.setLocale(const Locale('en'));
      expect(notifier.state.languageCode, 'en');
    });

    // setLocaleByCode로 일본어 설정
    test('setLocaleByCode로 일본어를 설정할 수 있다', () {
      notifier.setLocaleByCode('ja');
      expect(notifier.state.languageCode, 'ja');
    });

    // 비지원 언어 코드는 무시
    test('비지원 언어 코드(de)는 무시되고 기존 로케일이 유지된다', () {
      notifier.setLocaleByCode('en'); // 먼저 영어로 변경
      notifier.setLocaleByCode('de'); // 독일어 (미지원) 시도
      expect(notifier.state.languageCode, 'en'); // 영어 유지
    });

    // 6개 지원 언어 모두 설정 가능 확인
    test('6개 지원 언어 모두 설정 가능해야 한다', () {
      for (final code in ['ko', 'en', 'ja', 'zh', 'fr', 'hi']) {
        notifier.setLocaleByCode(code);
        expect(notifier.state.languageCode, code);
      }
    });
  });

  group('SupportedLocales 유틸리티 테스트', () {
    // 지원 언어 목록 수 확인
    test('SupportedLocales.all은 6개 언어를 포함한다', () {
      expect(SupportedLocales.all.length, 6);
    });

    // 기본 로케일 확인
    test('defaultLocale은 ko이다', () {
      expect(SupportedLocales.defaultLocale.languageCode, 'ko');
    });

    // getLanguageName 확인
    test('getLanguageName은 올바른 언어 이름을 반환한다', () {
      expect(SupportedLocales.getLanguageName('ko'), '한국어');
      expect(SupportedLocales.getLanguageName('en'), 'English');
      expect(SupportedLocales.getLanguageName('ja'), '日本語');
      expect(SupportedLocales.getLanguageName('zh'), '中文简体');
      expect(SupportedLocales.getLanguageName('fr'), 'Français');
      expect(SupportedLocales.getLanguageName('hi'), 'हिन्दी');
    });

    // 미지원 코드에 대한 fallback
    test('미지원 언어 코드는 코드 자체를 반환한다', () {
      expect(SupportedLocales.getLanguageName('de'), 'de');
      expect(SupportedLocales.getLanguageName('es'), 'es');
    });

    // languageNames 맵 크기 확인
    test('languageNames 맵은 6개 항목을 포함한다', () {
      expect(SupportedLocales.languageNames.length, 6);
    });
  });
}
