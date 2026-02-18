import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart' as intl;

import 'app_localizations_en.dart';
import 'app_localizations_fr.dart';
import 'app_localizations_hi.dart';
import 'app_localizations_ja.dart';
import 'app_localizations_ko.dart';
import 'app_localizations_zh.dart';

// ignore_for_file: type=lint

/// Callers can lookup localized strings with an instance of AppLocalizations
/// returned by `AppLocalizations.of(context)`.
///
/// Applications need to include `AppLocalizations.delegate()` in their app's
/// `localizationDelegates` list, and the locales they support in the app's
/// `supportedLocales` list. For example:
///
/// ```dart
/// import 'l10n/app_localizations.dart';
///
/// return MaterialApp(
///   localizationsDelegates: AppLocalizations.localizationsDelegates,
///   supportedLocales: AppLocalizations.supportedLocales,
///   home: MyApplicationHome(),
/// );
/// ```
///
/// ## Update pubspec.yaml
///
/// Please make sure to update your pubspec.yaml to include the following
/// packages:
///
/// ```yaml
/// dependencies:
///   # Internationalization support.
///   flutter_localizations:
///     sdk: flutter
///   intl: any # Use the pinned version from flutter_localizations
///
///   # Rest of dependencies
/// ```
///
/// ## iOS Applications
///
/// iOS applications define key application metadata, including supported
/// locales, in an Info.plist file that is built into the application bundle.
/// To configure the locales supported by your app, you’ll need to edit this
/// file.
///
/// First, open your project’s ios/Runner.xcworkspace Xcode workspace file.
/// Then, in the Project Navigator, open the Info.plist file under the Runner
/// project’s Runner folder.
///
/// Next, select the Information Property List item, select Add Item from the
/// Editor menu, then select Localizations from the pop-up menu.
///
/// Select and expand the newly-created Localizations item then, for each
/// locale your application supports, add a new item and select the locale
/// you wish to add from the pop-up menu in the Value field. This list should
/// be consistent with the languages listed in the AppLocalizations.supportedLocales
/// property.
abstract class AppLocalizations {
  AppLocalizations(String locale)
      : localeName = intl.Intl.canonicalizedLocale(locale.toString());

  final String localeName;

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  /// A list of this localizations delegate along with the default localizations
  /// delegates.
  ///
  /// Returns a list of localizations delegates containing this delegate along with
  /// GlobalMaterialLocalizations.delegate, GlobalCupertinoLocalizations.delegate,
  /// and GlobalWidgetsLocalizations.delegate.
  ///
  /// Additional delegates can be added by appending to this list in
  /// MaterialApp. This list does not have to be used at all if a custom list
  /// of delegates is preferred or required.
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates =
      <LocalizationsDelegate<dynamic>>[
    delegate,
    GlobalMaterialLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
  ];

  /// A list of this localizations delegate's supported locales.
  static const List<Locale> supportedLocales = <Locale>[
    Locale('en'),
    Locale('fr'),
    Locale('hi'),
    Locale('ja'),
    Locale('ko'),
    Locale('zh')
  ];

  /// MANPASIK 차동측정시스템 브랜드명
  ///
  /// In ko, this message translates to:
  /// **'MANPASIK'**
  String get appName;

  /// MANPASIK Device 슬로건
  ///
  /// In ko, this message translates to:
  /// **'초정밀 차동 계측 시스템'**
  String get appTagline;

  /// 홈 화면 인사말
  ///
  /// In ko, this message translates to:
  /// **'안녕하세요,'**
  String get greeting;

  /// 이름이 포함된 인사말
  ///
  /// In ko, this message translates to:
  /// **'{name}님'**
  String greetingWithName(String name);

  /// No description provided for @login.
  ///
  /// In ko, this message translates to:
  /// **'로그인'**
  String get login;

  /// No description provided for @register.
  ///
  /// In ko, this message translates to:
  /// **'회원가입'**
  String get register;

  /// No description provided for @logout.
  ///
  /// In ko, this message translates to:
  /// **'로그아웃'**
  String get logout;

  /// No description provided for @logoutConfirm.
  ///
  /// In ko, this message translates to:
  /// **'정말 로그아웃하시겠습니까?'**
  String get logoutConfirm;

  /// No description provided for @email.
  ///
  /// In ko, this message translates to:
  /// **'이메일'**
  String get email;

  /// No description provided for @emailHint.
  ///
  /// In ko, this message translates to:
  /// **'example@manpasik.com'**
  String get emailHint;

  /// No description provided for @password.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호'**
  String get password;

  /// No description provided for @passwordHint.
  ///
  /// In ko, this message translates to:
  /// **'8자 이상 (영문 + 숫자)'**
  String get passwordHint;

  /// No description provided for @passwordConfirm.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호 확인'**
  String get passwordConfirm;

  /// No description provided for @passwordConfirmHint.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호를 다시 입력해주세요'**
  String get passwordConfirmHint;

  /// No description provided for @displayName.
  ///
  /// In ko, this message translates to:
  /// **'연구원 성명'**
  String get displayName;

  /// No description provided for @displayNameHint.
  ///
  /// In ko, this message translates to:
  /// **'표시될 성명을 입력해주세요'**
  String get displayNameHint;

  /// No description provided for @noAccountYet.
  ///
  /// In ko, this message translates to:
  /// **'계정이 없으신가요?'**
  String get noAccountYet;

  /// No description provided for @alreadyHaveAccount.
  ///
  /// In ko, this message translates to:
  /// **'이미 계정이 있으신가요?'**
  String get alreadyHaveAccount;

  /// No description provided for @loginFailed.
  ///
  /// In ko, this message translates to:
  /// **'로그인에 실패했습니다. 자격 증명을 확인해주세요.'**
  String get loginFailed;

  /// No description provided for @registerFailed.
  ///
  /// In ko, this message translates to:
  /// **'연구원 등록에 실패했습니다. 다시 시도해주세요.'**
  String get registerFailed;

  /// No description provided for @home.
  ///
  /// In ko, this message translates to:
  /// **'대시보드'**
  String get home;

  /// No description provided for @measurement.
  ///
  /// In ko, this message translates to:
  /// **'분석'**
  String get measurement;

  /// No description provided for @devices.
  ///
  /// In ko, this message translates to:
  /// **'전술 디바이스'**
  String get devices;

  /// No description provided for @settings.
  ///
  /// In ko, this message translates to:
  /// **'설정'**
  String get settings;

  /// No description provided for @newMeasurement.
  ///
  /// In ko, this message translates to:
  /// **'신규 파동 분석'**
  String get newMeasurement;

  /// No description provided for @startMeasurement.
  ///
  /// In ko, this message translates to:
  /// **'분석 시작'**
  String get startMeasurement;

  /// No description provided for @startMeasurementAction.
  ///
  /// In ko, this message translates to:
  /// **'분석 시작하기'**
  String get startMeasurementAction;

  /// No description provided for @checkHealth.
  ///
  /// In ko, this message translates to:
  /// **'프로젝트 파동을\n분석해 보세요'**
  String get checkHealth;

  /// No description provided for @recentHistory.
  ///
  /// In ko, this message translates to:
  /// **'최근 분석 기록'**
  String get recentHistory;

  /// No description provided for @viewAll.
  ///
  /// In ko, this message translates to:
  /// **'전체 보기'**
  String get viewAll;

  /// No description provided for @preparingDevice.
  ///
  /// In ko, this message translates to:
  /// **'디바이스를 준비해 주세요'**
  String get preparingDevice;

  /// No description provided for @preparingDeviceDesc.
  ///
  /// In ko, this message translates to:
  /// **'센서를 장착하고\n측정 버튼을 눌러주세요'**
  String get preparingDeviceDesc;

  /// No description provided for @connectingDevice.
  ///
  /// In ko, this message translates to:
  /// **'디바이스 동기화 중...'**
  String get connectingDevice;

  /// No description provided for @connectingDeviceDesc.
  ///
  /// In ko, this message translates to:
  /// **'시스템 연결을 시도하고 있습니다'**
  String get connectingDeviceDesc;

  /// No description provided for @measuring.
  ///
  /// In ko, this message translates to:
  /// **'데이터 분석 중...'**
  String get measuring;

  /// No description provided for @measuringDesc.
  ///
  /// In ko, this message translates to:
  /// **'잠시만 기다려주세요'**
  String get measuringDesc;

  /// No description provided for @measurementComplete.
  ///
  /// In ko, this message translates to:
  /// **'분석 완료!'**
  String get measurementComplete;

  /// No description provided for @measurementFailed.
  ///
  /// In ko, this message translates to:
  /// **'분석 실패'**
  String get measurementFailed;

  /// No description provided for @viewResult.
  ///
  /// In ko, this message translates to:
  /// **'인사이트 확인'**
  String get viewResult;

  /// No description provided for @retryMeasurement.
  ///
  /// In ko, this message translates to:
  /// **'재분석'**
  String get retryMeasurement;

  /// No description provided for @bloodSugar.
  ///
  /// In ko, this message translates to:
  /// **'차원 수치'**
  String get bloodSugar;

  /// No description provided for @diagnosis.
  ///
  /// In ko, this message translates to:
  /// **'판정 결과'**
  String get diagnosis;

  /// No description provided for @noDevicesRegistered.
  ///
  /// In ko, this message translates to:
  /// **'연결된 디바이스가 없습니다'**
  String get noDevicesRegistered;

  /// No description provided for @noDevicesDesc.
  ///
  /// In ko, this message translates to:
  /// **'우측 상단의 + 버튼을 눌러\n새 디바이스를 등록해 주세요'**
  String get noDevicesDesc;

  /// No description provided for @searchDevices.
  ///
  /// In ko, this message translates to:
  /// **'시스템 검색'**
  String get searchDevices;

  /// No description provided for @addDevice.
  ///
  /// In ko, this message translates to:
  /// **'디바이스 등록'**
  String get addDevice;

  /// No description provided for @connected.
  ///
  /// In ko, this message translates to:
  /// **'연결됨'**
  String get connected;

  /// No description provided for @disconnected.
  ///
  /// In ko, this message translates to:
  /// **'오프라인'**
  String get disconnected;

  /// No description provided for @deviceRegistrationComingSoon.
  ///
  /// In ko, this message translates to:
  /// **'시스템 등록 기능은 곧 활성화됩니다'**
  String get deviceRegistrationComingSoon;

  /// No description provided for @profile.
  ///
  /// In ko, this message translates to:
  /// **'연구원 프로필'**
  String get profile;

  /// No description provided for @general.
  ///
  /// In ko, this message translates to:
  /// **'일반'**
  String get general;

  /// No description provided for @theme.
  ///
  /// In ko, this message translates to:
  /// **'인터페이스 테마'**
  String get theme;

  /// No description provided for @themeSystem.
  ///
  /// In ko, this message translates to:
  /// **'시스템 기본값'**
  String get themeSystem;

  /// No description provided for @themeLight.
  ///
  /// In ko, this message translates to:
  /// **'퓨어 라이트'**
  String get themeLight;

  /// No description provided for @themeDark.
  ///
  /// In ko, this message translates to:
  /// **'딥 씨 다크'**
  String get themeDark;

  /// No description provided for @themeSelect.
  ///
  /// In ko, this message translates to:
  /// **'테마 선택'**
  String get themeSelect;

  /// No description provided for @language.
  ///
  /// In ko, this message translates to:
  /// **'시스템 언어'**
  String get language;

  /// No description provided for @languageSelect.
  ///
  /// In ko, this message translates to:
  /// **'언어 선택'**
  String get languageSelect;

  /// No description provided for @appInfo.
  ///
  /// In ko, this message translates to:
  /// **'시스템 정보'**
  String get appInfo;

  /// No description provided for @version.
  ///
  /// In ko, this message translates to:
  /// **'버전'**
  String get version;

  /// No description provided for @termsOfService.
  ///
  /// In ko, this message translates to:
  /// **'이용약관 및 보안 지침'**
  String get termsOfService;

  /// No description provided for @privacyPolicy.
  ///
  /// In ko, this message translates to:
  /// **'개인정보 및 업무 보안 정책'**
  String get privacyPolicy;

  /// No description provided for @account.
  ///
  /// In ko, this message translates to:
  /// **'시스템 계정'**
  String get account;

  /// No description provided for @loginRequired.
  ///
  /// In ko, this message translates to:
  /// **'보안 로그인이 필요합니다'**
  String get loginRequired;

  /// No description provided for @user.
  ///
  /// In ko, this message translates to:
  /// **'연구원'**
  String get user;

  /// No description provided for @cancel.
  ///
  /// In ko, this message translates to:
  /// **'취소'**
  String get cancel;

  /// No description provided for @resultNormal.
  ///
  /// In ko, this message translates to:
  /// **'안정'**
  String get resultNormal;

  /// No description provided for @resultWarning.
  ///
  /// In ko, this message translates to:
  /// **'주의'**
  String get resultWarning;

  /// No description provided for @resultDanger.
  ///
  /// In ko, this message translates to:
  /// **'불균형'**
  String get resultDanger;

  /// No description provided for @validationEmailRequired.
  ///
  /// In ko, this message translates to:
  /// **'이메일을 입력해 주세요'**
  String get validationEmailRequired;

  /// No description provided for @validationEmailInvalid.
  ///
  /// In ko, this message translates to:
  /// **'올바른 이메일 형식이 아닙니다'**
  String get validationEmailInvalid;

  /// No description provided for @validationPasswordRequired.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호를 입력해 주세요'**
  String get validationPasswordRequired;

  /// No description provided for @validationPasswordTooShort.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호는 8자 이상이어야 합니다'**
  String get validationPasswordTooShort;

  /// No description provided for @validationPasswordNeedsLetter.
  ///
  /// In ko, this message translates to:
  /// **'영문자를 포함해야 합니다'**
  String get validationPasswordNeedsLetter;

  /// No description provided for @validationPasswordNeedsNumber.
  ///
  /// In ko, this message translates to:
  /// **'숫자를 포함해야 합니다'**
  String get validationPasswordNeedsNumber;

  /// No description provided for @validationNameRequired.
  ///
  /// In ko, this message translates to:
  /// **'성명을 입력해 주세요'**
  String get validationNameRequired;

  /// No description provided for @validationNameLength.
  ///
  /// In ko, this message translates to:
  /// **'성명은 2~50자 사이여야 합니다'**
  String get validationNameLength;

  /// No description provided for @validationPasswordMismatch.
  ///
  /// In ko, this message translates to:
  /// **'비밀번호가 일치하지 않습니다'**
  String get validationPasswordMismatch;
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  Future<AppLocalizations> load(Locale locale) {
    return SynchronousFuture<AppLocalizations>(lookupAppLocalizations(locale));
  }

  @override
  bool isSupported(Locale locale) => <String>[
        'en',
        'fr',
        'hi',
        'ja',
        'ko',
        'zh'
      ].contains(locale.languageCode);

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}

AppLocalizations lookupAppLocalizations(Locale locale) {
  // Lookup logic when only language code is specified.
  switch (locale.languageCode) {
    case 'en':
      return AppLocalizationsEn();
    case 'fr':
      return AppLocalizationsFr();
    case 'hi':
      return AppLocalizationsHi();
    case 'ja':
      return AppLocalizationsJa();
    case 'ko':
      return AppLocalizationsKo();
    case 'zh':
      return AppLocalizationsZh();
  }

  throw FlutterError(
      'AppLocalizations.delegate failed to load unsupported locale "$locale". This is likely '
      'an issue with the localizations generation tool. Please file an issue '
      'on GitHub with a reproducible sample app and the gen-l10n configuration '
      'that was used.');
}
