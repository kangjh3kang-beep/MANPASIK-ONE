// ignore: unused_import
import 'package:intl/intl.dart' as intl;
import 'app_localizations.dart';

// ignore_for_file: type=lint

/// The translations for Korean (`ko`).
class AppLocalizationsKo extends AppLocalizations {
  AppLocalizationsKo([String locale = 'ko']) : super(locale);

  @override
  String get appName => 'MANPASIK';

  @override
  String get appTagline => '초정밀 차동 계측 시스템';

  @override
  String get greeting => '안녕하세요,';

  @override
  String greetingWithName(String name) {
    return '$name님';
  }

  @override
  String get login => '로그인';

  @override
  String get register => '회원가입';

  @override
  String get logout => '로그아웃';

  @override
  String get logoutConfirm => '정말 로그아웃하시겠습니까?';

  @override
  String get email => '이메일';

  @override
  String get emailHint => 'example@manpasik.com';

  @override
  String get password => '비밀번호';

  @override
  String get passwordHint => '8자 이상 (영문 + 숫자)';

  @override
  String get passwordConfirm => '비밀번호 확인';

  @override
  String get passwordConfirmHint => '비밀번호를 다시 입력해주세요';

  @override
  String get displayName => '연구원 성명';

  @override
  String get displayNameHint => '표시될 성명을 입력해주세요';

  @override
  String get noAccountYet => '계정이 없으신가요?';

  @override
  String get alreadyHaveAccount => '이미 계정이 있으신가요?';

  @override
  String get loginFailed => '로그인에 실패했습니다. 자격 증명을 확인해주세요.';

  @override
  String get registerFailed => '연구원 등록에 실패했습니다. 다시 시도해주세요.';

  @override
  String get home => '대시보드';

  @override
  String get measurement => '분석';

  @override
  String get devices => '전술 디바이스';

  @override
  String get settings => '설정';

  @override
  String get newMeasurement => '신규 파동 분석';

  @override
  String get startMeasurement => '분석 시작';

  @override
  String get startMeasurementAction => '분석 시작하기';

  @override
  String get checkHealth => '프로젝트 파동을\n분석해 보세요';

  @override
  String get recentHistory => '최근 분석 기록';

  @override
  String get viewAll => '전체 보기';

  @override
  String get preparingDevice => '디바이스를 준비해 주세요';

  @override
  String get preparingDeviceDesc => '센서를 장착하고\n측정 버튼을 눌러주세요';

  @override
  String get connectingDevice => '디바이스 동기화 중...';

  @override
  String get connectingDeviceDesc => '시스템 연결을 시도하고 있습니다';

  @override
  String get measuring => '데이터 분석 중...';

  @override
  String get measuringDesc => '잠시만 기다려주세요';

  @override
  String get measurementComplete => '분석 완료!';

  @override
  String get measurementFailed => '분석 실패';

  @override
  String get viewResult => '인사이트 확인';

  @override
  String get retryMeasurement => '재분석';

  @override
  String get bloodSugar => '차원 수치';

  @override
  String get diagnosis => '판정 결과';

  @override
  String get noDevicesRegistered => '연결된 디바이스가 없습니다';

  @override
  String get noDevicesDesc => '우측 상단의 + 버튼을 눌러\n새 디바이스를 등록해 주세요';

  @override
  String get searchDevices => '시스템 검색';

  @override
  String get addDevice => '디바이스 등록';

  @override
  String get connected => '연결됨';

  @override
  String get disconnected => '오프라인';

  @override
  String get deviceRegistrationComingSoon => '시스템 등록 기능은 곧 활성화됩니다';

  @override
  String get profile => '연구원 프로필';

  @override
  String get general => '일반';

  @override
  String get theme => '인터페이스 테마';

  @override
  String get themeSystem => '시스템 기본값';

  @override
  String get themeLight => '퓨어 라이트';

  @override
  String get themeDark => '딥 씨 다크';

  @override
  String get themeSelect => '테마 선택';

  @override
  String get language => '시스템 언어';

  @override
  String get languageSelect => '언어 선택';

  @override
  String get appInfo => '시스템 정보';

  @override
  String get version => '버전';

  @override
  String get termsOfService => '이용약관 및 보안 지침';

  @override
  String get privacyPolicy => '개인정보 및 업무 보안 정책';

  @override
  String get account => '시스템 계정';

  @override
  String get loginRequired => '보안 로그인이 필요합니다';

  @override
  String get user => '연구원';

  @override
  String get cancel => '취소';

  @override
  String get resultNormal => '안정';

  @override
  String get resultWarning => '주의';

  @override
  String get resultDanger => '불균형';

  @override
  String get validationEmailRequired => '이메일을 입력해 주세요';

  @override
  String get validationEmailInvalid => '올바른 이메일 형식이 아닙니다';

  @override
  String get validationPasswordRequired => '비밀번호를 입력해 주세요';

  @override
  String get validationPasswordTooShort => '비밀번호는 8자 이상이어야 합니다';

  @override
  String get validationPasswordNeedsLetter => '영문자를 포함해야 합니다';

  @override
  String get validationPasswordNeedsNumber => '숫자를 포함해야 합니다';

  @override
  String get validationNameRequired => '성명을 입력해 주세요';

  @override
  String get validationNameLength => '성명은 2~50자 사이여야 합니다';

  @override
  String get validationPasswordMismatch => '비밀번호가 일치하지 않습니다';
}
