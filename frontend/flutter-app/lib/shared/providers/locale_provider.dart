import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 지원 언어 정의
///
/// 새 언어 추가 시 이 목록에 추가하고 해당 ARB 파일만 생성하면 됨.
/// 추후 확장 대상: 아랍어(ar), 스페인어(es), 독일어(de), 포르투갈어(pt), 태국어(th)
class SupportedLocales {
  SupportedLocales._();

  /// 기본 지원 6개 언어
  static const List<Locale> all = [
    Locale('ko'), // 한국어
    Locale('en'), // 영어
    Locale('ja'), // 일본어
    Locale('zh'), // 중국어 간체
    Locale('fr'), // 프랑스어
    Locale('hi'), // 힌디어
  ];

  /// 기본 로케일
  static const Locale defaultLocale = Locale('ko');

  /// 언어 이름 맵 (UI 표시용)
  static const Map<String, String> languageNames = {
    'ko': '한국어',
    'en': 'English',
    'ja': '日本語',
    'zh': '中文简体',
    'fr': 'Français',
    'hi': 'हिन्दी',
  };

  /// 언어 코드로 이름 가져오기
  static String getLanguageName(String code) {
    return languageNames[code] ?? code;
  }
}

/// 로케일 Notifier
///
/// 사용자 언어 설정 관리. SharedPreferences에 자동 저장.
class LocaleNotifier extends StateNotifier<Locale> {
  LocaleNotifier() : super(SupportedLocales.defaultLocale) {
    _loadSavedLocale();
  }

  static const _prefKey = 'app_locale';

  Future<void> _loadSavedLocale() async {
    final prefs = await SharedPreferences.getInstance();
    final code = prefs.getString(_prefKey);
    if (code != null && SupportedLocales.languageNames.containsKey(code)) {
      state = Locale(code);
    }
  }

  /// 로케일 변경
  void setLocale(Locale locale) {
    if (SupportedLocales.all.contains(locale)) {
      state = locale;
      SharedPreferences.getInstance().then((prefs) {
        prefs.setString(_prefKey, locale.languageCode);
      });
    }
  }

  /// 언어 코드로 로케일 변경
  void setLocaleByCode(String languageCode) {
    setLocale(Locale(languageCode));
  }
}

/// 현재 로케일 Provider
final localeProvider = StateNotifierProvider<LocaleNotifier, Locale>((ref) {
  return LocaleNotifier();
});
