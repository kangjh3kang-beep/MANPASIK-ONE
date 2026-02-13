import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// 테마 모드 Provider
///
/// light / dark / system 테마 전환 지원.
/// 사용자 선택은 SharedPreferences에 저장 (추후 구현).
class ThemeModeNotifier extends StateNotifier<ThemeMode> {
  ThemeModeNotifier() : super(ThemeMode.system);

  /// 테마 모드 변경
  void setThemeMode(ThemeMode mode) {
    state = mode;
    // TODO: SharedPreferences에 저장
  }

  /// 라이트 모드
  void setLight() => setThemeMode(ThemeMode.light);

  /// 다크 모드
  void setDark() => setThemeMode(ThemeMode.dark);

  /// 시스템 기본값
  void setSystem() => setThemeMode(ThemeMode.system);

  /// 순환 토글 (system → light → dark → system)
  void toggle() {
    switch (state) {
      case ThemeMode.system:
        setLight();
        break;
      case ThemeMode.light:
        setDark();
        break;
      case ThemeMode.dark:
        setSystem();
        break;
    }
  }
}

/// 테마 모드 StateNotifier Provider
final themeModeProvider =
    StateNotifierProvider<ThemeModeNotifier, ThemeMode>((ref) {
  return ThemeModeNotifier();
});
