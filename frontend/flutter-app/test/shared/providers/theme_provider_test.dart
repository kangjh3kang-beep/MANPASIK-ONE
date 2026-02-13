// ThemeModeNotifier Provider 테스트
// - 기본값, setLight/setDark/setSystem, toggle 순환

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';

void main() {
  late ThemeModeNotifier notifier;

  setUp(() {
    notifier = ThemeModeNotifier();
  });

  group('ThemeModeNotifier 테스트', () {
    // 초기값 확인
    test('초기 테마 모드는 ThemeMode.system이어야 한다', () {
      expect(notifier.state, ThemeMode.system);
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

    // toggle 순환: system → light → dark → system
    test('toggle은 system → light → dark → system 순으로 순환한다', () {
      expect(notifier.state, ThemeMode.system);

      notifier.toggle();
      expect(notifier.state, ThemeMode.light);

      notifier.toggle();
      expect(notifier.state, ThemeMode.dark);

      notifier.toggle();
      expect(notifier.state, ThemeMode.system);
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
