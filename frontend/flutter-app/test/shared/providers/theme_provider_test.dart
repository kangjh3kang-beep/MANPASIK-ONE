// ThemeModeNotifier Provider 테스트
// - 기본값, setLight/setDark/setSystem, toggle 순환

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';

void main() {
  late ThemeModeNotifier notifier;

  setUp(() {
    SharedPreferences.setMockInitialValues({});
    notifier = ThemeModeNotifier();
  });

  group('ThemeModeNotifier 테스트', () {
    // 초기값 확인
    test('초기 테마 모드는 ThemeMode.dark이어야 한다', () {
      expect(notifier.state, ThemeMode.dark);
    });

    // setLight 호출
    test('setLight 호출 시 ThemeMode.light로 변경된다', () {
      notifier.setLight();
      expect(notifier.state, ThemeMode.light);
    });

    // setDark 호출
    test('setDark 호출 시 ThemeMode.dark로 변경된다', () {
      notifier.setDark();
      expect(notifier.state, ThemeMode.dark);
    });

    // setSystem 호출
    test('setSystem 호출 시 ThemeMode.system으로 변경된다', () {
      notifier.setDark(); // 먼저 dark로 변경
      notifier.setSystem();
      expect(notifier.state, ThemeMode.system);
    });

    // toggle 순환: dark → system → light → dark
    test('toggle은 dark → system → light → dark 순으로 순환한다', () {
      expect(notifier.state, ThemeMode.dark);

      notifier.toggle();
      expect(notifier.state, ThemeMode.system);

      notifier.toggle();
      expect(notifier.state, ThemeMode.light);

      notifier.toggle();
      expect(notifier.state, ThemeMode.dark);
    });

    // setThemeMode 직접 호출
    test('setThemeMode로 직접 모드를 설정할 수 있다', () {
      notifier.setThemeMode(ThemeMode.dark);
      expect(notifier.state, ThemeMode.dark);

      notifier.setThemeMode(ThemeMode.light);
      expect(notifier.state, ThemeMode.light);
    });
  });
}
