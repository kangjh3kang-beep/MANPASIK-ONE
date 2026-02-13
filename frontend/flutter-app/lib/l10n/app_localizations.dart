import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';

import 'package:manpasik/l10n/translations/ko.dart';
import 'package:manpasik/l10n/translations/en.dart';
import 'package:manpasik/l10n/translations/ja.dart';
import 'package:manpasik/l10n/translations/zh.dart';
import 'package:manpasik/l10n/translations/fr.dart';
import 'package:manpasik/l10n/translations/hi.dart';

/// ManPaSik 다국어 지원 클래스
///
/// 6개 언어 지원: ko, en, ja, zh, fr, hi
/// 확장 시 translations/ 디렉토리에 새 파일 추가 + supportedLocales 등록
class AppLocalizations {
  final Locale locale;
  late final Map<String, String> _translations;

  AppLocalizations(this.locale) {
    _translations = _loadTranslations(locale.languageCode);
  }

  static Map<String, String> _loadTranslations(String code) {
    switch (code) {
      case 'en':
        return enTranslations;
      case 'ja':
        return jaTranslations;
      case 'zh':
        return zhTranslations;
      case 'fr':
        return frTranslations;
      case 'hi':
        return hiTranslations;
      case 'ko':
      default:
        return koTranslations;
    }
  }

  /// BuildContext에서 현재 로케일의 AppLocalizations 가져오기
  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  /// LocalizationsDelegate
  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// 지원 로케일 목록
  static const List<Locale> supportedLocales = [
    Locale('ko'),
    Locale('en'),
    Locale('ja'),
    Locale('zh'),
    Locale('fr'),
    Locale('hi'),
  ];

  // ── 번역 키 접근자 ──

  String get appName => _translations['appName']!;
  String get appTagline => _translations['appTagline']!;
  String get greeting => _translations['greeting']!;
  String greetingWithName(String name) =>
      (_translations['greetingWithName'] ?? '{name}').replaceAll('{name}', name);

  String get login => _translations['login']!;
  String get register => _translations['register']!;
  String get logout => _translations['logout']!;
  String get logoutConfirm => _translations['logoutConfirm']!;
  String get email => _translations['email']!;
  String get emailHint => _translations['emailHint']!;
  String get password => _translations['password']!;
  String get passwordHint => _translations['passwordHint']!;
  String get passwordConfirm => _translations['passwordConfirm']!;
  String get passwordConfirmHint => _translations['passwordConfirmHint']!;
  String get displayName => _translations['displayName']!;
  String get displayNameHint => _translations['displayNameHint']!;
  String get noAccountYet => _translations['noAccountYet']!;
  String get alreadyHaveAccount => _translations['alreadyHaveAccount']!;
  String get loginFailed => _translations['loginFailed']!;
  String get registerFailed => _translations['registerFailed']!;

  String get home => _translations['home']!;
  String get measurement => _translations['measurement']!;
  String get devices => _translations['devices']!;
  String get settings => _translations['settings']!;

  String get newMeasurement => _translations['newMeasurement']!;
  String get startMeasurement => _translations['startMeasurement']!;
  String get startMeasurementAction => _translations['startMeasurementAction']!;
  String get checkHealth => _translations['checkHealth']!;
  String get recentHistory => _translations['recentHistory']!;
  String get viewAll => _translations['viewAll']!;

  String get preparingDevice => _translations['preparingDevice']!;
  String get preparingDeviceDesc => _translations['preparingDeviceDesc']!;
  String get connectingDevice => _translations['connectingDevice']!;
  String get connectingDeviceDesc => _translations['connectingDeviceDesc']!;
  String get measuring => _translations['measuring']!;
  String get measuringDesc => _translations['measuringDesc']!;
  String get measurementComplete => _translations['measurementComplete']!;
  String get measurementFailed => _translations['measurementFailed']!;
  String get viewResult => _translations['viewResult']!;
  String get retryMeasurement => _translations['retryMeasurement']!;
  String get bloodSugar => _translations['bloodSugar']!;
  String get diagnosis => _translations['diagnosis']!;

  String get noDevicesRegistered => _translations['noDevicesRegistered']!;
  String get noDevicesDesc => _translations['noDevicesDesc']!;
  String get searchDevices => _translations['searchDevices']!;
  String get addDevice => _translations['addDevice']!;
  String get connected => _translations['connected']!;
  String get disconnected => _translations['disconnected']!;
  String get deviceRegistrationComingSoon =>
      _translations['deviceRegistrationComingSoon']!;

  String get profile => _translations['profile']!;
  String get general => _translations['general']!;
  String get theme => _translations['theme']!;
  String get themeSystem => _translations['themeSystem']!;
  String get themeLight => _translations['themeLight']!;
  String get themeDark => _translations['themeDark']!;
  String get themeSelect => _translations['themeSelect']!;
  String get language => _translations['language']!;
  String get languageSelect => _translations['languageSelect']!;
  String get appInfo => _translations['appInfo']!;
  String get version => _translations['version']!;
  String get termsOfService => _translations['termsOfService']!;
  String get privacyPolicy => _translations['privacyPolicy']!;
  String get account => _translations['account']!;
  String get loginRequired => _translations['loginRequired']!;
  String get user => _translations['user']!;
  String get cancel => _translations['cancel']!;

  String get resultNormal => _translations['resultNormal']!;
  String get resultWarning => _translations['resultWarning']!;
  String get resultDanger => _translations['resultDanger']!;

  String get validationEmailRequired =>
      _translations['validationEmailRequired']!;
  String get validationEmailInvalid =>
      _translations['validationEmailInvalid']!;
  String get validationPasswordRequired =>
      _translations['validationPasswordRequired']!;
  String get validationPasswordTooShort =>
      _translations['validationPasswordTooShort']!;
  String get validationPasswordNeedsLetter =>
      _translations['validationPasswordNeedsLetter']!;
  String get validationPasswordNeedsNumber =>
      _translations['validationPasswordNeedsNumber']!;
  String get validationNameRequired =>
      _translations['validationNameRequired']!;
  String get validationNameLength => _translations['validationNameLength']!;
  String get validationPasswordMismatch =>
      _translations['validationPasswordMismatch']!;

  // ── Chat (AI 건강 어시스턴트) ──
  String get chatTitle => _translations['chatTitle']!;
  String get chatWelcome => _translations['chatWelcome']!;
  String get chatInputHint => _translations['chatInputHint']!;
  String get chatSend => _translations['chatSend']!;
  String get chatExampleBloodSugar => _translations['chatExampleBloodSugar']!;
  String get chatExampleBloodPressure =>
      _translations['chatExampleBloodPressure']!;
  String get chatExampleExercise => _translations['chatExampleExercise']!;
  String get chatExampleDiet => _translations['chatExampleDiet']!;
  String get chatExampleSleep => _translations['chatExampleSleep']!;
  String get chatDisclaimer => _translations['chatDisclaimer']!;
  String get chatClearHistory => _translations['chatClearHistory']!;
  String get chatClearConfirm => _translations['chatClearConfirm']!;
  String get chatErrorGeneric => _translations['chatErrorGeneric']!;
  String get chatTyping => _translations['chatTyping']!;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return ['ko', 'en', 'ja', 'zh', 'fr', 'hi']
        .contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(AppLocalizations(locale));
  }

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}
