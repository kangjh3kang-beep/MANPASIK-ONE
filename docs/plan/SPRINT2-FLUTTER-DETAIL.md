# Sprint 2 Flutter 프론트엔드 세부 구현 기획서

> **문서 ID**: PLAN-SPRINT2-FLUTTER-001  
> **버전**: v1.0  
> **작성일**: 2026-02-12  
> **상태**: 기획 확정 (구현 대기)  
> **범위**: F-1 (단위 테스트 60개), AS-6 (Admin 설정 UI), AS-8 Flutter 부분 (LLM 채팅 UI)

---

## 목차

1. [현행 아키텍처 분석](#1-현행-아키텍처-분석)
2. [F-1: Flutter 단위 테스트 60개](#2-f-1-flutter-단위-테스트-60개)
3. [AS-6: Admin 설정 UI](#3-as-6-admin-설정-ui)
4. [AS-8: LLM 채팅 UI](#4-as-8-llm-채팅-ui)
5. [패키지 의존성](#5-패키지-의존성)
6. [다국어 (l10n) 추가 키](#6-다국어-l10n-추가-키)
7. [gRPC 연동 명세](#7-grpc-연동-명세)
8. [검증 기준](#8-검증-기준)

---

## 1. 현행 아키텍처 분석

### 1.1 디렉토리 구조 (현재)

```
frontend/flutter-app/
├── lib/
│   ├── core/
│   │   ├── constants/app_constants.dart        # 상수 정의
│   │   ├── providers/grpc_provider.dart         # gRPC Repository Provider 모음
│   │   ├── router/app_router.dart               # GoRouter (7개 라우트)
│   │   ├── services/
│   │   │   ├── auth_interceptor.dart            # JWT 인터셉터
│   │   │   ├── grpc_client.dart                 # GrpcClientManager
│   │   │   ├── rest_client.dart                 # REST 클라이언트
│   │   │   └── rust_ffi_stub.dart               # Rust FFI 스텁
│   │   ├── theme/app_theme.dart                 # 라이트/다크 테마
│   │   └── utils/validators.dart                # 입력 검증 유틸리티
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/auth_repository_impl.dart
│   │   │   ├── domain/auth_repository.dart
│   │   │   └── presentation/ (login, register, splash)
│   │   ├── devices/
│   │   │   ├── data/device_repository_impl.dart
│   │   │   ├── domain/device_repository.dart
│   │   │   └── presentation/ (device_list, ble_scan_dialog)
│   │   ├── home/presentation/home_screen.dart
│   │   ├── measurement/
│   │   │   ├── data/measurement_repository_impl.dart
│   │   │   ├── domain/measurement_repository.dart
│   │   │   └── presentation/ (measurement, result)
│   │   ├── settings/presentation/settings_screen.dart
│   │   └── user/
│   │       ├── data/user_repository_impl.dart
│   │       └── domain/user_repository.dart
│   ├── generated/ (manpasik.pb.dart, manpasik.pbgrpc.dart)
│   ├── l10n/
│   │   ├── app_localizations.dart
│   │   └── translations/ (ko, en, ja, zh, fr, hi)
│   ├── shared/
│   │   ├── providers/ (auth_provider, theme_provider, locale_provider)
│   │   └── widgets/ (measurement_card, jagae_pattern, primary_button)
│   └── main.dart
├── test/
│   ├── helpers/fake_repositories.dart
│   ├── widget_test.dart            # 34개 테스트 (auth, theme, locale, validator)
│   ├── repository_test.dart        # 16개 테스트 (fake repository)
│   ├── screen_widget_test.dart     # 5개 테스트 (화면 위젯)
│   └── grpc_client_test.dart       # 6개 테스트 (gRPC 클라이언트)
└── pubspec.yaml
```

### 1.2 상태 관리 패턴

| Provider | 타입 | 역할 |
|----------|------|------|
| `authProvider` | `StateNotifierProvider<AuthNotifier, AuthState>` | 인증 상태 (로그인/로그아웃/게스트) |
| `themeModeProvider` | `StateNotifierProvider<ThemeModeNotifier, ThemeMode>` | 테마 모드 (system/light/dark) |
| `localeProvider` | `StateNotifierProvider<LocaleNotifier, Locale>` | 언어 설정 (6개 언어) |
| `measurementHistoryProvider` | `FutureProvider<MeasurementHistoryResult>` | 최근 측정 기록 |
| `deviceListProvider` | `FutureProvider<List<DeviceItem>>` | 디바이스 목록 |
| `userProfileProvider` | `FutureProvider<UserProfileInfo?>` | 사용자 프로필 |
| `subscriptionInfoProvider` | `FutureProvider<SubscriptionInfoDto?>` | 구독 정보 |

### 1.3 기존 라우트 (7개)

| 경로 | 화면 | 인증 필요 |
|------|------|-----------|
| `/` | SplashScreen | N |
| `/login` | LoginScreen | N |
| `/register` | RegisterScreen | N |
| `/home` | HomeScreen | Y |
| `/measurement` | MeasurementScreen | Y |
| `/measurement/result` | MeasurementResultScreen | Y |
| `/devices` | DeviceListScreen | Y |
| `/settings` | SettingsScreen | Y |

### 1.4 기존 테스트 현황

| 파일 | 테스트 수 | 카테고리 |
|------|-----------|----------|
| `widget_test.dart` | 34 | auth provider(7), theme(4), locale(10), validator(11), misc(2) |
| `repository_test.dart` | 16 | fake auth(7), device(2), measurement(4), user(4) |
| `screen_widget_test.dart` | 5 | home(3), device(2), measurement_result(1) |
| `grpc_client_test.dart` | 6 | GrpcClientManager(4), AuthInterceptor(3) |
| **합계** | **~61** | |

> **참고**: 기존 테스트가 이미 약 61개 존재하나, 일부 테스트가 누락/불완전하고 AS-6/AS-8 관련 테스트가 없으므로 F-1에서는 기존 테스트 보강 + 신규 Admin/LLM UI 테스트를 추가하여 총 60개 이상의 **신규 또는 개선** 테스트를 확보한다.

---

## 2. F-1: Flutter 단위 테스트 60개

### 2.1 테스트 전략

```
test/
├── helpers/
│   ├── fake_repositories.dart          # [기존] + admin/llm fake 추가
│   ├── test_app_wrapper.dart           # [신규] ProviderScope + MaterialApp 래퍼
│   └── golden_comparator.dart          # [신규] 골든 테스트 비교기 (향후 확장용)
├── unit/
│   ├── auth_provider_test.dart         # [신규] auth provider 전용 단위 테스트
│   ├── theme_provider_test.dart        # [신규] theme provider 전용 단위 테스트
│   ├── locale_provider_test.dart       # [신규] locale provider 전용 단위 테스트
│   ├── validators_test.dart            # [신규] validators 전용 단위 테스트
│   ├── app_constants_test.dart         # [신규] 상수 검증 테스트
│   ├── auth_interceptor_test.dart      # [신규] AuthInterceptor 단위 테스트
│   └── admin_settings_provider_test.dart  # [신규] AS-6 AdminSettings provider 테스트
├── widget/
│   ├── measurement_card_test.dart      # [신규] MeasurementCard 위젯 테스트
│   ├── primary_button_test.dart        # [신규] PrimaryButton 위젯 테스트
│   ├── home_screen_test.dart           # [신규] HomeScreen 위젯 테스트
│   ├── login_screen_test.dart          # [신규] LoginScreen 위젯 테스트
│   ├── device_list_screen_test.dart    # [신규] DeviceListScreen 위젯 테스트
│   ├── settings_screen_test.dart       # [신규] SettingsScreen 위젯 테스트
│   ├── admin_settings_screen_test.dart # [신규] AS-6 AdminSettingsScreen 테스트
│   ├── config_edit_dialog_test.dart    # [신규] AS-6 ConfigEditDialog 테스트
│   └── llm_chat_screen_test.dart       # [신규] AS-8 LLM 채팅 화면 테스트
├── integration/
│   └── fake_repository_integration_test.dart  # [신규] Fake Repository 통합 테스트
├── widget_test.dart                    # [기존 유지]
├── repository_test.dart                # [기존 유지]
├── screen_widget_test.dart             # [기존 유지]
└── grpc_client_test.dart               # [기존 유지]
```

### 2.2 테스트 목록 (60개 신규 테스트)

#### 2.2.1 `test/unit/auth_provider_test.dart` (8개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 1 | `초기 상태는 비인증` | `AuthState.isAuthenticated == false`, 모든 필드 null |
| 2 | `유효한 자격증명으로 로그인 성공` | login 후 `isAuthenticated == true`, email/userId/token 설정 |
| 3 | `빈 이메일로 로그인 실패` | `login('', 'password123')` → false, 상태 변경 없음 |
| 4 | `짧은 비밀번호로 로그인 실패` | `login('test@test.com', 'short')` → false |
| 5 | `로그아웃 시 모든 상태 초기화` | 로그인 후 logout → 모든 필드 초기화 확인 |
| 6 | `회원가입 성공 시 displayName 설정` | register 후 displayName 확인 |
| 7 | `게스트 로그인 시 guest-user 설정` | loginAsGuest → userId == 'guest-user' |
| 8 | `checkAuthStatus 미인증 시 상태 초기화` | isAuthenticated가 false → state 초기화 확인 |

#### 2.2.2 `test/unit/theme_provider_test.dart` (6개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 9 | `기본 테마 모드는 system` | 초기 state == ThemeMode.system |
| 10 | `setLight 호출 시 light 모드` | state == ThemeMode.light |
| 11 | `setDark 호출 시 dark 모드` | state == ThemeMode.dark |
| 12 | `setSystem 호출 시 system 모드` | setDark → setSystem → state == ThemeMode.system |
| 13 | `toggle 순환 검증` | system → light → dark → system 순환 3회 |
| 14 | `setThemeMode 직접 호출` | setThemeMode(ThemeMode.dark) → dark 확인 |

#### 2.2.3 `test/unit/locale_provider_test.dart` (8개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 15 | `기본 로케일은 ko` | state.languageCode == 'ko' |
| 16 | `6개 지원 언어 전환 (ko→en→ja→zh→fr→hi)` | 순차 전환 후 각각 languageCode 확인 |
| 17 | `미지원 언어 코드 무시` | setLocaleByCode('ar') → 변경 없음 (ko 유지) |
| 18 | `setLocale(Locale) 직접 호출` | setLocale(Locale('en')) → 'en' 확인 |
| 19 | `SupportedLocales.all 6개 확인` | length == 6, 모든 코드 포함 |
| 20 | `SupportedLocales.defaultLocale == ko` | defaultLocale.languageCode == 'ko' |
| 21 | `getLanguageName 6개 언어 이름 확인` | 각 코드별 올바른 이름 반환 |
| 22 | `getLanguageName 미지원 코드 시 코드 자체 반환` | getLanguageName('xx') == 'xx' |

#### 2.2.4 `test/unit/validators_test.dart` (10개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 23 | `이메일 - null 입력` | validateEmail(null) → 에러 메시지 |
| 24 | `이메일 - 빈 문자열` | validateEmail('') → 에러 메시지 |
| 25 | `이메일 - @ 없음` | validateEmail('notanemail') → 에러 |
| 26 | `이메일 - 도메인 없음` | validateEmail('test@') → 에러 |
| 27 | `이메일 - 유효한 형식` | validateEmail('user@example.com') → null |
| 28 | `비밀번호 - 빈 입력` | validatePassword('') → 에러 |
| 29 | `비밀번호 - 7자 이하` | validatePassword('abc1234') → 에러 |
| 30 | `비밀번호 - 숫자 없음` | validatePassword('abcdefgh') → 에러 |
| 31 | `비밀번호 - 영문 없음` | validatePassword('12345678') → 에러 |
| 32 | `비밀번호 - 유효 (영문+숫자 8자 이상)` | validatePassword('abc12345') → null |

#### 2.2.5 `test/unit/app_constants_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 33 | `앱 이름 상수 확인` | appName == '만파식', appNameEn == 'ManPaSik' |
| 34 | `gRPC 포트 기본값 확인` | auth=50051, user=50052, device=50053, measurement=50054 |
| 35 | `구독 티어 문자열 확인` | FREE/BASIC/PRO/CLINICAL 존재 |
| 36 | `UI 상수 범위 확인` | borderRadius, padding 값 > 0 |

#### 2.2.6 `test/unit/auth_interceptor_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 37 | `tokenProvider null 시 빈 메타데이터` | 토큰 없을 때 authorization 헤더 미포함 |
| 38 | `tokenProvider 빈 문자열 시 메타데이터 미포함` | '' 반환 시 authorization 미포함 |
| 39 | `tokenProvider 유효 토큰 시 Bearer 토큰 포함` | 'my-token' → 'Bearer my-token' 확인 |
| 40 | `interceptStreaming에서도 토큰 전달` | 스트리밍 호출 시도 동일 동작 확인 |

#### 2.2.7 `test/unit/admin_settings_provider_test.dart` (6개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 41 | `초기 카테고리 목록 로드` | FakeAdminRepository → 10개 카테고리 반환 |
| 42 | `카테고리별 설정 필터링` | category='payment' → payment 설정만 반환 |
| 43 | `설정 검색 (키워드)` | 'toss' 검색 → toss 관련 설정만 필터 |
| 44 | `설정 값 변경 시 상태 업데이트` | setValue → state 변경 확인 |
| 45 | `설정 유효성 검증 호출` | validateConfigValue → valid/invalid 결과 |
| 46 | `secret 타입 설정 마스킹 확인` | security_level=secret → 값 '****' 표시 |

#### 2.2.8 `test/widget/measurement_card_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 47 | `정상 결과 - 초록색 아이콘 표시` | resultType='normal' → 초록색 확인 |
| 48 | `주의 결과 - 주황색 아이콘 표시` | resultType='warning' → 주황색 확인 |
| 49 | `위험 결과 - 빨간색 아이콘 표시` | resultType='danger' → 빨간색 확인 |
| 50 | `날짜 포맷 및 값 표시` | '01월 15일 14:30', '98.4', 'mg/dL' 텍스트 확인 |

#### 2.2.9 `test/widget/primary_button_test.dart` (3개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 51 | `버튼 텍스트 표시` | text='측정 시작' → 텍스트 finder 확인 |
| 52 | `로딩 상태에서 CircularProgressIndicator 표시` | isLoading=true → 인디케이터 존재 |
| 53 | `탭 콜백 호출` | onPressed mock → 호출 횟수 1 확인 |

#### 2.2.10 `test/widget/home_screen_test.dart` (3개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 54 | `인증 시 사용자명 표시` | authState.displayName 텍스트 존재 |
| 55 | `최근 기록 섹션 표시` | '최근 기록' 텍스트 finder |
| 56 | `측정 시작하기 버튼 존재` | '측정 시작하기' 텍스트 finder |

#### 2.2.11 `test/widget/login_screen_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 57 | `이메일/비밀번호 입력 필드 존재` | TextFormField 2개 finder |
| 58 | `로그인 버튼 탭 시 폼 유효성 검증` | 빈 상태에서 탭 → 에러 메시지 표시 |
| 59 | `비밀번호 가시성 토글` | 아이콘 버튼 탭 → obscureText 전환 |
| 60 | `회원가입 링크 존재` | '계정이 없으신가요?' 텍스트 확인 |

#### 2.2.12 `test/widget/device_list_screen_test.dart` (2개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 61 | `앱바에 디바이스 타이틀 표시` | '디바이스' 텍스트 finder |
| 62 | `빈 목록 시 빈 상태 UI 표시` | '등록된 디바이스가 없습니다' 텍스트 확인 |

#### 2.2.13 `test/widget/settings_screen_test.dart` (3개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 63 | `프로필/일반/앱 정보 섹션 헤더 존재` | '프로필', '일반', '앱 정보' 텍스트 finder |
| 64 | `테마 선택 다이얼로그 표시` | 테마 ListTile 탭 → 다이얼로그 표시 확인 |
| 65 | `언어 선택 다이얼로그에 6개 언어 표시` | 언어 ListTile 탭 → 6개 옵션 확인 |

#### 2.2.14 `test/widget/admin_settings_screen_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 66 | `카테고리 탭 10개 표시` | general~integration 탭 존재 |
| 67 | `설정 카드 목록 표시` | 카테고리 선택 시 해당 설정 카드 표시 |
| 68 | `설정 검색 필드 동작` | 검색어 입력 → 필터링 결과 표시 |
| 69 | `카테고리별 설정 수 배지 표시` | 각 탭에 Badge 위젯으로 수 표시 |

#### 2.2.15 `test/widget/config_edit_dialog_test.dart` (5개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 70 | `string 유형 - TextField 표시` | valueType='string' → TextField 존재 |
| 71 | `boolean 유형 - Switch 표시` | valueType='boolean' → Switch 위젯 존재 |
| 72 | `select 유형 - DropdownButton 표시` | valueType='select' → DropdownButton 존재 |
| 73 | `secret 유형 - 마스킹 토글 표시` | valueType='secret' → 가시성 토글 버튼 존재 |
| 74 | `변경 확인 다이얼로그 - 사유 입력 필드` | 저장 버튼 탭 → 변경 사유 TextField 존재 |

#### 2.2.16 `test/widget/llm_chat_screen_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 75 | `채팅 입력 필드 및 전송 버튼 존재` | TextField + 전송 IconButton finder |
| 76 | `메시지 버블 표시` | 메시지 전송 시 ChatBubble 위젯 생성 |
| 77 | `제안 패널에 ConfigSuggestion 카드 표시` | 제안 도착 시 카드 위젯 존재 |
| 78 | `적용/거부 버튼 동작` | 적용 버튼 탭 → onApply 콜백 호출 확인 |

#### 2.2.17 `test/integration/fake_repository_integration_test.dart` (4개)

| # | 테스트명 | 검증 내용 |
|---|---------|-----------|
| 79 | `AuthProvider + FakeAuth 로그인→기록조회 흐름` | 로그인 → measurementHistory 로드 성공 |
| 80 | `AuthProvider + FakeDevice 디바이스 목록 연동` | 로그인 → deviceList 로드 → 1개 디바이스 |
| 81 | `AuthProvider + FakeUser 프로필 조회` | 로그인 → userProfile 로드 → displayName 확인 |
| 82 | `로그아웃 후 모든 Provider 초기화 확인` | 로그아웃 → history/device/profile 모두 빈 상태 |

### 2.3 테스트 헬퍼 파일

#### `test/helpers/test_app_wrapper.dart` (신규)

```dart
/// 테스트용 ProviderScope + MaterialApp 래퍼
///
/// Fake Repository 오버라이드를 일괄 적용하여
/// 모든 위젯 테스트에서 일관된 환경을 제공한다.
class TestAppWrapper extends StatelessWidget {
  final Widget child;
  final List<Override> overrides;

  const TestAppWrapper({
    super.key,
    required this.child,
    this.overrides = const [],
  });

  @override
  Widget build(BuildContext context) {
    return ProviderScope(
      overrides: [
        authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
        deviceRepositoryProvider.overrideWithValue(FakeDeviceRepository()),
        measurementRepositoryProvider.overrideWithValue(FakeMeasurementRepository()),
        userRepositoryProvider.overrideWithValue(FakeUserRepository()),
        ...overrides,
      ],
      child: MaterialApp(
        localizationsDelegates: const [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: AppLocalizations.supportedLocales,
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        home: child,
      ),
    );
  }
}
```

#### `test/helpers/fake_repositories.dart` 확장 (기존 + 추가)

기존 `FakeAuthRepository`, `FakeDeviceRepository`, `FakeMeasurementRepository`, `FakeUserRepository`에 아래를 추가:

```dart
/// 테스트용 Fake AdminSettingsRepository
class FakeAdminSettingsRepository implements AdminSettingsRepository {
  @override
  Future<ListConfigsResult> listConfigs({
    String languageCode = 'ko',
    String? category,
    bool includeSecrets = false,
  }) async {
    // 10개 카테고리 × 2~3개 설정 = 더미 데이터 반환
  }

  @override
  Future<ValidateResult> validateConfigValue(String key, String value) async {
    return ValidateResult(valid: value.isNotEmpty);
  }

  @override
  Future<void> setConfig(String key, String value, String reason) async {}
}

/// 테스트용 Fake LlmAssistantRepository
class FakeLlmAssistantRepository implements LlmAssistantRepository {
  @override
  Future<ConfigSessionInfo> startSession({...}) async { ... }

  @override
  Future<AssistantResponse> sendMessage(String sessionId, String message) async {
    return AssistantResponse(
      message: '테스트 응답입니다.',
      suggestions: [],
    );
  }

  @override
  Future<ApplyResult> applySuggestion(String sessionId, String suggestionId) async { ... }
}
```

### 2.4 테스트 실행 기준

```bash
# 전체 테스트 실행
flutter test

# 특정 파일 실행
flutter test test/unit/auth_provider_test.dart

# 커버리지 리포트
flutter test --coverage
genhtml coverage/lcov.info -o coverage/html
```

**목표 커버리지**: 라인 커버리지 70% 이상 (`lib/` 기준, `generated/` 제외)

---

## 3. AS-6: Admin 설정 UI

### 3.1 신규 파일 목록

```
lib/features/admin/
├── data/
│   └── admin_settings_repository_impl.dart     # gRPC 연동 구현체
├── domain/
│   ├── admin_settings_repository.dart           # Repository 인터페이스
│   └── models/
│       ├── config_item.dart                     # ConfigItem 도메인 모델
│       ├── config_category.dart                 # ConfigCategory enum
│       └── validate_result.dart                 # ValidateResult 모델
├── presentation/
│   ├── admin_settings_screen.dart               # 메인 설정 화면
│   ├── widgets/
│   │   ├── category_tab_bar.dart                # 카테고리 탭 바
│   │   ├── config_card.dart                     # 개별 설정 카드
│   │   ├── config_edit_dialog.dart              # 설정 편집 다이얼로그
│   │   ├── config_search_bar.dart               # 설정 검색 바
│   │   ├── change_reason_dialog.dart            # 변경 사유 입력 다이얼로그
│   │   └── markdown_help_viewer.dart            # 마크다운 도움말 뷰어
│   └── providers/
│       └── admin_settings_provider.dart         # 설정 상태 관리
└── README.md (불필요 - 생성하지 않음)
```

### 3.2 라우트 추가

`lib/core/router/app_router.dart` 에 아래 2개 라우트 추가:

```dart
GoRoute(
  path: '/admin/settings',
  builder: (context, state) => const AdminSettingsScreen(),
),
GoRoute(
  path: '/admin/settings/assistant',
  builder: (context, state) => const LlmChatScreen(),
),
```

> 인증 리다이렉트 로직에서 `/admin/*` 경로는 인증 필수 + admin 역할 검증 추가.

### 3.3 도메인 모델

#### `lib/features/admin/domain/models/config_category.dart`

```dart
/// 설정 카테고리 열거형
///
/// 백엔드 config_category enum과 1:1 매핑.
enum ConfigCategory {
  general('general', '일반', Icons.settings),
  payment('payment', '결제', Icons.payment),
  auth('auth', '인증', Icons.lock),
  storage('storage', '스토리지', Icons.cloud),
  messaging('messaging', '메시징', Icons.message),
  database('database', '데이터베이스', Icons.storage),
  ai('ai', 'AI/ML', Icons.smart_toy),
  notification('notification', '알림', Icons.notifications),
  security('security', '보안', Icons.security),
  integration('integration', '외부 연동', Icons.extension);

  const ConfigCategory(this.value, this.displayName, this.icon);

  final String value;
  final String displayName;
  final IconData icon;

  /// 문자열에서 ConfigCategory로 변환
  static ConfigCategory fromString(String value) {
    return ConfigCategory.values.firstWhere(
      (c) => c.value == value,
      orElse: () => ConfigCategory.general,
    );
  }
}
```

#### `lib/features/admin/domain/models/config_item.dart`

```dart
/// 시스템 설정 항목 도메인 모델
///
/// gRPC ConfigWithMeta 메시지와 매핑.
class ConfigItem {
  final String key;
  final String value;
  final String rawValue;          // secret일 때 실제 값 (권한 있을 때만)

  // 메타데이터
  final ConfigCategory category;
  final String valueType;         // string, number, boolean, secret, url, select, json, multiline
  final String securityLevel;     // public, internal, confidential, secret
  final bool isRequired;
  final String defaultValue;
  final List<String> allowedValues;
  final String? validationRegex;
  final double? validationMin;
  final double? validationMax;
  final String? dependsOn;
  final String? dependsValue;
  final String? envVarName;
  final String? serviceName;
  final bool restartRequired;

  // 다국어 번역
  final String displayName;
  final String description;
  final String? placeholder;
  final String? helpText;         // 마크다운
  final String? validationMessage;

  // 변경 정보
  final String? updatedBy;
  final DateTime? updatedAt;

  const ConfigItem({...});

  /// gRPC ConfigWithMeta에서 변환
  factory ConfigItem.fromProto(ConfigWithMeta proto) { ... }
}
```

#### `lib/features/admin/domain/admin_settings_repository.dart`

```dart
/// Admin 설정 Repository 인터페이스
abstract class AdminSettingsRepository {
  /// 설정 목록 조회 (카테고리 필터, 언어별 번역 포함)
  Future<ListConfigsResult> listConfigs({
    String languageCode = 'ko',
    String? category,
    bool includeSecrets = false,
  });

  /// 단일 설정 조회 (메타데이터 포함)
  Future<ConfigItem> getConfigWithMeta(String key, {String languageCode = 'ko'});

  /// 설정 값 변경
  Future<void> setConfig(String key, String value, String reason);

  /// 설정 값 유효성 검증
  Future<ValidateResult> validateConfigValue(String key, String value);

  /// 일괄 설정 변경
  Future<BulkSetResult> bulkSetConfigs(List<ConfigChange> changes, String reason);
}

class ListConfigsResult {
  final List<ConfigItem> configs;
  final Map<String, int> categoryCounts;

  const ListConfigsResult({required this.configs, required this.categoryCounts});
}

class ValidateResult {
  final bool valid;
  final String? errorMessage;
  final List<String> suggestions;

  const ValidateResult({required this.valid, this.errorMessage, this.suggestions = const []});
}

class ConfigChange {
  final String key;
  final String value;

  const ConfigChange({required this.key, required this.value});
}

class BulkSetResult {
  final int successCount;
  final int failureCount;
  final List<ConfigChangeResult> results;

  const BulkSetResult({...});
}
```

### 3.4 Provider (상태 관리)

#### `lib/features/admin/presentation/providers/admin_settings_provider.dart`

```dart
/// Admin 설정 상태 모델
class AdminSettingsState {
  final List<ConfigItem> allConfigs;
  final List<ConfigItem> filteredConfigs;
  final Map<String, int> categoryCounts;
  final ConfigCategory selectedCategory;
  final String searchQuery;
  final bool isLoading;
  final String? error;

  const AdminSettingsState({
    this.allConfigs = const [],
    this.filteredConfigs = const [],
    this.categoryCounts = const {},
    this.selectedCategory = ConfigCategory.general,
    this.searchQuery = '',
    this.isLoading = false,
    this.error,
  });

  AdminSettingsState copyWith({...});
}

/// Admin 설정 StateNotifier
class AdminSettingsNotifier extends StateNotifier<AdminSettingsState> {
  AdminSettingsNotifier(this._repository) : super(const AdminSettingsState());

  final AdminSettingsRepository _repository;

  /// 설정 목록 로드
  Future<void> loadConfigs({String languageCode = 'ko'}) async { ... }

  /// 카테고리 변경
  void selectCategory(ConfigCategory category) { ... }

  /// 검색어 변경
  void search(String query) { ... }

  /// 설정 값 변경
  Future<bool> updateConfig(String key, String value, String reason) async { ... }

  /// 설정 유효성 검증
  Future<ValidateResult> validate(String key, String value) async { ... }

  /// 내부: 필터링 적용
  void _applyFilters() { ... }
}

/// Provider 정의
final adminSettingsProvider =
    StateNotifierProvider<AdminSettingsNotifier, AdminSettingsState>((ref) {
  final repository = ref.watch(adminSettingsRepositoryProvider);
  return AdminSettingsNotifier(repository);
});

final adminSettingsRepositoryProvider = Provider<AdminSettingsRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return AdminSettingsRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});
```

### 3.5 화면 구조 (위젯 트리)

#### AdminSettingsScreen

```
Scaffold
├── AppBar
│   ├── BackButton
│   ├── Title: "시스템 설정"
│   └── Actions
│       └── IconButton (LLM 어시스턴트 → /admin/settings/assistant)
├── Body: Column
│   ├── ConfigSearchBar                        ← 검색 바
│   │   └── TextField (검색 아이콘, 실시간 필터)
│   ├── CategoryTabBar                         ← 카테고리 탭 바
│   │   └── SingleChildScrollView (horizontal)
│   │       └── Row
│   │           └── [10개] ChoiceChip / FilterChip
│   │               ├── Icon (카테고리 아이콘)
│   │               ├── Text (카테고리명)
│   │               └── Badge (설정 수)
│   └── Expanded
│       └── ListView.builder                   ← 설정 카드 목록
│           └── [N개] ConfigCard
│               ├── ListTile
│               │   ├── leading: Icon (valueType 아이콘)
│               │   ├── title: Text (displayName)
│               │   ├── subtitle: Column
│               │   │   ├── Text (description, maxLines: 2)
│               │   │   └── Row (chips: security_level, restart_required)
│               │   └── trailing: Column
│               │       ├── Text (현재 값, secret이면 '****')
│               │       └── Icon (chevron_right)
│               └── onTap → ConfigEditDialog
```

#### ConfigEditDialog

```
AlertDialog
├── title: Row
│   ├── Icon (valueType 아이콘)
│   └── Text (displayName)
├── content: SingleChildScrollView
│   └── Column
│       ├── Text (description)
│       ├── SizedBox(height: 8)
│       ├── if (helpText != null)
│       │   └── ExpansionTile ("도움말")
│       │       └── MarkdownBody (flutter_markdown)
│       ├── SizedBox(height: 16)
│       ├── _buildInputWidget(valueType)       ← 유형별 입력 위젯
│       │   ├── string    → TextField
│       │   ├── number    → TextField (keyboardType: number, inputFormatters)
│       │   ├── boolean   → SwitchListTile
│       │   ├── secret    → TextField + IconButton (visibility toggle)
│       │   ├── url       → TextField + URL 검증 (validateUrl)
│       │   ├── select    → DropdownButtonFormField (allowedValues)
│       │   ├── json      → TextField (maxLines: 10, monospace font) + Format 버튼
│       │   └── multiline → TextField (maxLines: 5)
│       ├── SizedBox(height: 8)
│       ├── if (validationMessage != null)
│       │   └── Text (에러 메시지, 빨간색)
│       ├── SizedBox(height: 8)
│       └── Row (메타 정보)
│           ├── Chip ("카테고리: ${category}")
│           ├── Chip ("보안: ${securityLevel}")
│           └── if (restartRequired) Chip ("재시작 필요", 주황)
└── actions
    ├── TextButton ("취소")
    └── FilledButton ("저장") → ChangeReasonDialog
```

#### ChangeReasonDialog

```
AlertDialog
├── title: Text ("설정 변경 확인")
├── content: Column
│   ├── Text ("다음 설정을 변경하시겠습니까?")
│   ├── Card
│   │   ├── Text ("키: ${key}")
│   │   ├── Text ("현재 값: ${oldValue}")
│   │   └── Text ("변경 값: ${newValue}")
│   ├── SizedBox(height: 16)
│   ├── TextField (변경 사유 입력, required)
│   └── if (restartRequired)
│       └── Banner ("이 변경은 서비스 재시작이 필요합니다", 주황)
└── actions
    ├── TextButton ("취소")
    └── FilledButton ("변경 적용")
```

### 3.6 주요 클래스/메서드 명세

#### `AdminSettingsScreen` (`ConsumerStatefulWidget`)

| 메서드 | 설명 |
|--------|------|
| `initState()` | `adminSettingsProvider.loadConfigs()` 호출 |
| `build()` | Scaffold + SearchBar + TabBar + ListView 구성 |
| `_onCategorySelected(ConfigCategory)` | 카테고리 변경 → provider.selectCategory() |
| `_onSearchChanged(String)` | 검색어 변경 → provider.search() |
| `_openConfigEditor(ConfigItem)` | ConfigEditDialog 표시 |
| `_navigateToAssistant()` | context.push('/admin/settings/assistant') |

#### `ConfigEditDialog` (`ConsumerStatefulWidget`)

| 메서드 | 설명 |
|--------|------|
| `_buildInputWidget(String valueType)` | 유형별 입력 위젯 반환 |
| `_validateInput()` | 로컬 유효성 검증 + gRPC validateConfigValue |
| `_formatJson()` | JSON 포맷팅 (indent 2) |
| `_toggleSecretVisibility()` | secret 마스킹/표시 토글 |
| `_onSave()` | ChangeReasonDialog 표시 → confirm → setConfig 호출 |

#### `ConfigCard` (`StatelessWidget`)

| 프로퍼티 | 타입 | 설명 |
|----------|------|------|
| `config` | `ConfigItem` | 설정 항목 |
| `onTap` | `VoidCallback` | 탭 시 편집 다이얼로그 열기 |

#### `CategoryTabBar` (`StatelessWidget`)

| 프로퍼티 | 타입 | 설명 |
|----------|------|------|
| `selectedCategory` | `ConfigCategory` | 현재 선택 카테고리 |
| `categoryCounts` | `Map<String, int>` | 카테고리별 설정 수 |
| `onCategorySelected` | `ValueChanged<ConfigCategory>` | 카테고리 선택 콜백 |

---

## 4. AS-8: LLM 채팅 UI

### 4.1 신규 파일 목록

```
lib/features/admin/
├── data/
│   └── llm_assistant_repository_impl.dart      # gRPC 연동 구현체
├── domain/
│   ├── llm_assistant_repository.dart            # Repository 인터페이스
│   └── models/
│       ├── chat_message.dart                    # 채팅 메시지 모델
│       ├── config_suggestion.dart               # 설정 제안 모델
│       └── assistant_session.dart               # 세션 모델
└── presentation/
    ├── llm_chat_screen.dart                     # LLM 채팅 메인 화면
    ├── widgets/
    │   ├── chat_message_bubble.dart             # 메시지 버블
    │   ├── chat_input_bar.dart                  # 채팅 입력 바
    │   ├── suggestion_panel.dart                # 우측 제안 패널
    │   ├── suggestion_card.dart                 # 개별 제안 카드
    │   └── language_category_bar.dart           # 언어 + 카테고리 필터 바
    └── providers/
        └── llm_chat_provider.dart               # 채팅 상태 관리
```

### 4.2 도메인 모델

#### `lib/features/admin/domain/models/chat_message.dart`

```dart
/// 채팅 메시지 모델
class ChatMessage {
  final String id;
  final ChatRole role;               // user, assistant, system
  final String content;              // 마크다운 지원
  final List<ConfigSuggestion> suggestions;
  final DateTime createdAt;

  const ChatMessage({...});
}

enum ChatRole { user, assistant, system }
```

#### `lib/features/admin/domain/models/config_suggestion.dart`

```dart
/// LLM 설정 변경 제안 모델
///
/// gRPC ConfigSuggestion 메시지와 매핑.
class ConfigSuggestion {
  final String suggestionId;
  final String configKey;
  final String currentValue;         // secret이면 '****'
  final String suggestedValue;       // secret이면 직접 입력 필요
  final String reason;
  final bool isSecret;
  final SuggestionStatus status;     // pending, applied, rejected

  const ConfigSuggestion({...});
}

enum SuggestionStatus { pending, applied, rejected }
```

#### `lib/features/admin/domain/llm_assistant_repository.dart`

```dart
/// LLM 설정 어시스턴트 Repository 인터페이스
abstract class LlmAssistantRepository {
  /// 새 대화 세션 시작
  Future<ConfigSessionInfo> startSession({
    required String adminId,
    String languageCode = 'ko',
    String? category,
  });

  /// 메시지 전송 → 어시스턴트 응답
  Future<AssistantResponse> sendMessage(String sessionId, String message);

  /// 설정 제안 적용
  Future<ApplyResult> applySuggestion(
    String sessionId,
    String suggestionId, {
    String? overrideValue,
  });

  /// 세션 종료
  Future<void> endSession(String sessionId);
}

class ConfigSessionInfo {
  final String sessionId;
  final String status;
  final String welcomeMessage;

  const ConfigSessionInfo({...});
}

class AssistantResponse {
  final String message;               // 마크다운 응답
  final List<ConfigSuggestion> suggestions;
  final bool requiresConfirmation;

  const AssistantResponse({...});
}

class ApplyResult {
  final bool success;
  final String message;

  const ApplyResult({...});
}
```

### 4.3 Provider (상태 관리)

#### `lib/features/admin/presentation/providers/llm_chat_provider.dart`

```dart
/// LLM 채팅 상태 모델
class LlmChatState {
  final String? sessionId;
  final List<ChatMessage> messages;
  final List<ConfigSuggestion> activeSuggestions;   // 현재 적용 대기 제안
  final String selectedLanguage;
  final ConfigCategory? selectedCategory;
  final bool isLoading;                              // 어시스턴트 응답 대기
  final String? error;

  const LlmChatState({
    this.sessionId,
    this.messages = const [],
    this.activeSuggestions = const [],
    this.selectedLanguage = 'ko',
    this.selectedCategory,
    this.isLoading = false,
    this.error,
  });

  LlmChatState copyWith({...});
}

/// LLM 채팅 StateNotifier
class LlmChatNotifier extends StateNotifier<LlmChatState> {
  LlmChatNotifier(this._repository) : super(const LlmChatState());

  final LlmAssistantRepository _repository;

  /// 세션 시작
  Future<void> startSession({String? category}) async { ... }

  /// 메시지 전송
  Future<void> sendMessage(String text) async {
    // 1. 사용자 메시지 추가
    // 2. isLoading = true
    // 3. _repository.sendMessage → 어시스턴트 응답
    // 4. 응답 메시지 + 제안 추가
    // 5. isLoading = false
  }

  /// 제안 적용
  Future<void> applySuggestion(String suggestionId, {String? overrideValue}) async { ... }

  /// 제안 거부
  void rejectSuggestion(String suggestionId) { ... }

  /// 언어 변경
  void setLanguage(String languageCode) { ... }

  /// 카테고리 필터 변경
  void setCategory(ConfigCategory? category) { ... }

  /// 세션 종료
  Future<void> endSession() async { ... }
}

/// Provider 정의
final llmChatProvider =
    StateNotifierProvider<LlmChatNotifier, LlmChatState>((ref) {
  final repository = ref.watch(llmAssistantRepositoryProvider);
  return LlmChatNotifier(repository);
});

final llmAssistantRepositoryProvider = Provider<LlmAssistantRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return LlmAssistantRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});
```

### 4.4 화면 구조 (위젯 트리)

#### LlmChatScreen (반응형: 좁으면 탭 전환, 넓으면 좌우 분할)

```
Scaffold
├── AppBar
│   ├── BackButton
│   ├── Title: "설정 어시스턴트"
│   └── Actions
│       └── PopupMenuButton (세션 종료, 초기화)
├── Body: LayoutBuilder
│   ├── if (width >= 768)                    ← 태블릿/데스크톱: 좌우 분할
│   │   └── Row
│   │       ├── Expanded(flex: 3)            ← 좌: 채팅 영역
│   │       │   └── _ChatPanel
│   │       ├── VerticalDivider
│   │       └── Expanded(flex: 2)            ← 우: 제안 패널
│   │           └── _SuggestionPanel
│   └── else                                 ← 모바일: 탭 전환
│       └── Column
│           ├── LanguageCategoryBar
│           └── Expanded
│               └── _ChatPanel (제안은 인라인 표시)
```

#### _ChatPanel

```
Column
├── LanguageCategoryBar                       ← 언어 + 카테고리 필터
│   └── Row
│       ├── DropdownButton<String>             (6개 언어)
│       └── DropdownButton<ConfigCategory?>    (카테고리 필터)
├── Expanded
│   └── ListView.builder                      ← 메시지 목록
│       └── [N개] ChatMessageBubble
│           ├── if (role == user)
│           │   └── Align(right) → Container (primary color)
│           │       └── Text (content)
│           └── if (role == assistant)
│               └── Align(left) → Container (surface color)
│                   ├── MarkdownBody (content)  ← 마크다운 렌더링
│                   └── if (suggestions.isNotEmpty && 모바일)
│                       └── Column
│                           └── [N개] SuggestionCard (인라인)
├── if (isLoading)
│   └── LinearProgressIndicator
└── ChatInputBar
    └── Row
        ├── Expanded → TextField (hint: '설정에 대해 질문해보세요...')
        └── IconButton (send) → sendMessage()
```

#### _SuggestionPanel (데스크톱 우측)

```
Column
├── Padding
│   └── Text ("설정 변경 제안", style: titleMedium)
├── Divider
└── Expanded
    └── ListView.builder
        └── [N개] SuggestionCard
            └── Card
                ├── ListTile
                │   ├── title: Text (configKey)
                │   ├── subtitle: Column
                │   │   ├── Row ("현재: ${currentValue}")
                │   │   ├── Row ("제안: ${suggestedValue}", bold)
                │   │   └── Text (reason, maxLines: 3)
                │   └── trailing: Column (상태 아이콘)
                ├── if (isSecret)
                │   └── Padding → TextField ("직접 입력")
                └── ButtonBar
                    ├── OutlinedButton ("거부") → rejectSuggestion()
                    └── FilledButton ("적용") → applySuggestion()
```

### 4.5 주요 클래스/메서드 명세

#### `LlmChatScreen` (`ConsumerStatefulWidget`)

| 메서드 | 설명 |
|--------|------|
| `initState()` | `llmChatProvider.startSession()` 호출 |
| `dispose()` | `llmChatProvider.endSession()` 호출 |
| `build()` | LayoutBuilder → 반응형 레이아웃 |
| `_buildChatPanel()` | 채팅 메시지 목록 + 입력 바 |
| `_buildSuggestionPanel()` | 제안 카드 목록 |

#### `ChatMessageBubble` (`StatelessWidget`)

| 프로퍼티 | 타입 | 설명 |
|----------|------|------|
| `message` | `ChatMessage` | 메시지 데이터 |
| `onSuggestionApply` | `Function(String)` | 인라인 제안 적용 콜백 |
| `onSuggestionReject` | `Function(String)` | 인라인 제안 거부 콜백 |

#### `SuggestionCard` (`StatelessWidget`)

| 프로퍼티 | 타입 | 설명 |
|----------|------|------|
| `suggestion` | `ConfigSuggestion` | 제안 데이터 |
| `onApply` | `VoidCallback` | 적용 버튼 콜백 |
| `onReject` | `VoidCallback` | 거부 버튼 콜백 |
| `overrideController` | `TextEditingController?` | secret용 직접 입력 컨트롤러 |

#### `ChatInputBar` (`StatelessWidget`)

| 프로퍼티 | 타입 | 설명 |
|----------|------|------|
| `onSend` | `ValueChanged<String>` | 메시지 전송 콜백 |
| `isLoading` | `bool` | 전송 중 비활성화 |

---

## 5. 패키지 의존성

### 5.1 추가 필요 패키지 (dependencies)

| 패키지 | 버전 | 용도 | 사용 위치 |
|--------|------|------|-----------|
| `flutter_markdown` | `^0.7.4` | 마크다운 렌더링 | ConfigEditDialog (helpText), ChatMessageBubble |
| `url_launcher` | `^6.2.3` | 마크다운 내 링크 오픈 | MarkdownHelpViewer |
| `badges` | `^3.1.2` | 카테고리 설정 수 배지 | CategoryTabBar |

### 5.2 추가 필요 패키지 (dev_dependencies)

| 패키지 | 버전 | 용도 |
|--------|------|------|
| `network_image_mock` | `^2.1.1` | CachedNetworkImage 테스트 목킹 |

### 5.3 pubspec.yaml 수정 diff

```yaml
dependencies:
  # 기존 패키지 유지...

  # [신규] Admin 설정 UI
  flutter_markdown: ^0.7.4
  url_launcher: ^6.2.3
  badges: ^3.1.2

dev_dependencies:
  # 기존 패키지 유지...

  # [신규] 테스트 보조
  network_image_mock: ^2.1.1
```

> **주의**: `mocktail: ^1.0.4`는 이미 dev_dependencies에 포함되어 있으므로 추가 불필요.

---

## 6. 다국어 (l10n) 추가 키

### 6.1 Admin 설정 UI 관련 (AS-6)

6개 언어 파일 (`ko.dart`, `en.dart`, `ja.dart`, `zh.dart`, `fr.dart`, `hi.dart`) 모두에 추가:

| 키 | 한국어 (ko) | 영어 (en) |
|----|-------------|-----------|
| `adminSettings` | `시스템 설정` | `System Settings` |
| `adminSettingsSearch` | `설정 검색...` | `Search settings...` |
| `adminSettingsNoResult` | `일치하는 설정이 없습니다` | `No matching settings found` |
| `adminSettingsCount` | `{count}개 설정` | `{count} settings` |
| `configCategoryGeneral` | `일반` | `General` |
| `configCategoryPayment` | `결제` | `Payment` |
| `configCategoryAuth` | `인증` | `Authentication` |
| `configCategoryStorage` | `스토리지` | `Storage` |
| `configCategoryMessaging` | `메시징` | `Messaging` |
| `configCategoryDatabase` | `데이터베이스` | `Database` |
| `configCategoryAi` | `AI/ML` | `AI/ML` |
| `configCategoryNotification` | `알림` | `Notification` |
| `configCategorySecurity` | `보안` | `Security` |
| `configCategoryIntegration` | `외부 연동` | `Integration` |
| `configEditTitle` | `설정 편집` | `Edit Setting` |
| `configEditSave` | `저장` | `Save` |
| `configEditCancel` | `취소` | `Cancel` |
| `configEditHelpToggle` | `도움말` | `Help` |
| `configEditCurrentValue` | `현재 값` | `Current Value` |
| `configEditDefaultValue` | `기본값: {value}` | `Default: {value}` |
| `configEditRequired` | `필수 항목` | `Required` |
| `configEditRestartRequired` | `변경 시 서비스 재시작 필요` | `Service restart required` |
| `configEditSecurityLevel` | `보안 등급: {level}` | `Security: {level}` |
| `configChangeConfirm` | `설정 변경 확인` | `Confirm Setting Change` |
| `configChangeReason` | `변경 사유` | `Change Reason` |
| `configChangeReasonHint` | `변경 사유를 입력해주세요` | `Enter reason for change` |
| `configChangeApply` | `변경 적용` | `Apply Change` |
| `configChangeCancelConfirm` | `변경을 취소하시겠습니까?` | `Cancel the change?` |
| `configValidationFailed` | `유효성 검증 실패` | `Validation failed` |
| `configSecretMasked` | `(암호화됨)` | `(encrypted)` |
| `configSecretToggle` | `값 보기/숨기기` | `Show/hide value` |
| `configJsonFormat` | `JSON 포맷팅` | `Format JSON` |
| `configJsonInvalid` | `올바른 JSON 형식이 아닙니다` | `Invalid JSON format` |

### 6.2 LLM 채팅 UI 관련 (AS-8)

| 키 | 한국어 (ko) | 영어 (en) |
|----|-------------|-----------|
| `llmAssistantTitle` | `설정 어시스턴트` | `Settings Assistant` |
| `llmAssistantHint` | `설정에 대해 질문해보세요...` | `Ask about settings...` |
| `llmAssistantSend` | `전송` | `Send` |
| `llmAssistantSuggestions` | `설정 변경 제안` | `Setting Change Suggestions` |
| `llmAssistantNoSuggestions` | `아직 제안이 없습니다` | `No suggestions yet` |
| `llmSuggestionApply` | `적용` | `Apply` |
| `llmSuggestionReject` | `거부` | `Reject` |
| `llmSuggestionApplied` | `적용됨` | `Applied` |
| `llmSuggestionRejected` | `거부됨` | `Rejected` |
| `llmSuggestionCurrent` | `현재` | `Current` |
| `llmSuggestionProposed` | `제안` | `Proposed` |
| `llmSuggestionReason` | `사유` | `Reason` |
| `llmSuggestionSecretInput` | `값을 직접 입력해주세요` | `Please enter the value directly` |
| `llmLanguageSelect` | `응답 언어` | `Response Language` |
| `llmCategoryFilter` | `카테고리 필터` | `Category Filter` |
| `llmCategoryAll` | `전체` | `All` |
| `llmSessionEnd` | `세션 종료` | `End Session` |
| `llmSessionReset` | `대화 초기화` | `Reset Chat` |
| `llmTyping` | `입력 중...` | `Typing...` |

### 6.3 키 추가 방법

각 언어 파일(`lib/l10n/translations/{lang}.dart`)의 Map에 위 키-값을 추가하고, `lib/l10n/app_localizations.dart`에 접근자 getter를 추가:

```dart
// app_localizations.dart에 추가할 접근자 예시
String get adminSettings => _translations['adminSettings']!;
String get adminSettingsSearch => _translations['adminSettingsSearch']!;
String get llmAssistantTitle => _translations['llmAssistantTitle']!;
// ... 등
```

---

## 7. gRPC 연동 명세

### 7.1 기존 gRPC 클라이언트

현재 `lib/generated/manpasik.pbgrpc.dart`에 정의된 4개 서비스 클라이언트:

| 클라이언트 | 포트 | 사용 중인 메서드 |
|-----------|------|-----------------|
| `AuthServiceClient` | 50051 | Register, Login, RefreshToken, Logout |
| `UserServiceClient` | 50052 | GetProfile, GetSubscription |
| `DeviceServiceClient` | 50053 | ListDevices |
| `MeasurementServiceClient` | 50054 | StartSession, EndSession, GetMeasurementHistory |

### 7.2 AS-6 추가 필요 gRPC 메서드 (AdminService 확장)

| RPC | 요청 메시지 | 응답 메시지 | 사용 위치 |
|-----|-----------|-----------|-----------|
| `ListSystemConfigs` | `ListSystemConfigsRequest` | `ListSystemConfigsResponse` | AdminSettingsScreen 초기 로드 |
| `GetConfigWithMeta` | `GetConfigWithMetaRequest` | `ConfigWithMeta` | ConfigEditDialog 상세 조회 |
| `ListConfigsByCategory` | `ListConfigsByCategoryRequest` | `ListSystemConfigsResponse` | 카테고리 탭 전환 시 |
| `SetSystemConfig` | `SetSystemConfigRequest` | `SystemConfig` | 설정 값 변경 |
| `ValidateConfigValue` | `ValidateConfigValueRequest` | `ValidateConfigValueResponse` | 입력 실시간 검증 |
| `BulkSetConfigs` | `BulkSetConfigsRequest` | `BulkSetConfigsResponse` | 일괄 변경 (향후) |

### 7.3 AS-8 추가 필요 gRPC 메서드 (AiInferenceService 확장)

| RPC | 요청 메시지 | 응답 메시지 | 사용 위치 |
|-----|-----------|-----------|-----------|
| `StartConfigSession` | `StartConfigSessionRequest` | `ConfigSessionResponse` | 채팅 화면 진입 시 |
| `SendConfigMessage` | `SendConfigMessageRequest` | `ConfigAssistantResponse` | 메시지 전송 |
| `ApplyConfigSuggestion` | `ApplyConfigSuggestionRequest` | `ApplyConfigSuggestionResponse` | 제안 적용 버튼 |
| `EndConfigSession` | `EndConfigSessionRequest` | `ConfigSessionResponse` | 세션 종료 |

### 7.4 AdminServiceClient 추가 (신규)

`lib/generated/manpasik.pbgrpc.dart` 또는 별도 파일에 AdminServiceClient 정의:

```dart
class AdminServiceClient extends Client {
  AdminServiceClient(super.channel, {super.interceptors, super.options});

  // 기존
  ResponseFuture<SystemConfig> setSystemConfig(SetSystemConfigRequest req, {...});
  ResponseFuture<SystemConfig> getSystemConfig(GetSystemConfigRequest req, {...});

  // 신규 (AS-6)
  ResponseFuture<ListSystemConfigsResponse> listSystemConfigs(ListSystemConfigsRequest req, {...});
  ResponseFuture<ConfigWithMeta> getConfigWithMeta(GetConfigWithMetaRequest req, {...});
  ResponseFuture<ValidateConfigValueResponse> validateConfigValue(ValidateConfigValueRequest req, {...});
  ResponseFuture<BulkSetConfigsResponse> bulkSetConfigs(BulkSetConfigsRequest req, {...});
}
```

### 7.5 GrpcClientManager 확장

`lib/core/services/grpc_client.dart`에 admin 채널 추가:

```dart
// 추가 포트 상수 (AppConstants에 추가)
static const int grpcAdminPort = 50055;
static const int grpcAiInferencePort = 50056;

// GrpcClientManager에 추가
ClientChannel get adminChannel { ... }
ClientChannel get aiInferenceChannel { ... }
```

### 7.6 grpc_provider.dart 확장

```dart
// Admin Settings Repository Provider
final adminSettingsRepositoryProvider = Provider<AdminSettingsRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return AdminSettingsRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

// LLM Assistant Repository Provider
final llmAssistantRepositoryProvider = Provider<LlmAssistantRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return LlmAssistantRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});
```

---

## 8. 검증 기준

### 8.1 flutter analyze 통과 기준

```bash
flutter analyze --no-fatal-infos
```

| 기준 | 조건 |
|------|------|
| Error | 0개 (zero tolerance) |
| Warning | 0개 (zero tolerance) |
| Info | 허용 (but should minimize) |

**주요 lint 규칙** (기존 `flutter_lints: ^3.0.1` 기반):
- `prefer_const_constructors`: 가능한 const 사용
- `avoid_print`: logger 패키지 사용
- `prefer_final_locals`: final 변수 우선
- `always_declare_return_types`: 반환 타입 명시
- `sort_constructors_first`: 생성자 우선 배치

### 8.2 flutter test 통과 기준

```bash
flutter test --reporter=expanded
```

| 기준 | 조건 |
|------|------|
| 전체 테스트 통과 | 100% (0 failures) |
| 신규 테스트 수 | 60개 이상 |
| 테스트 실행 시간 | 60초 이내 |
| 커버리지 (lib/ - generated/) | 70% 이상 |

### 8.3 빌드 검증

```bash
# Android 빌드 (APK)
flutter build apk --debug --no-tree-shake-icons

# Web 빌드
flutter build web --no-tree-shake-icons

# iOS 빌드 (macOS에서만)
flutter build ios --no-codesign --no-tree-shake-icons
```

| 기준 | 조건 |
|------|------|
| Android Debug APK | 빌드 성공 |
| Web 빌드 | 빌드 성공 |
| Dart 포맷 | `dart format --set-exit-if-changed lib/ test/` 통과 |

### 8.4 CI 파이프라인 검증 체크리스트

```yaml
# .github/workflows/ci.yml 에 추가될 검증 단계
- flutter pub get
- flutter analyze --no-fatal-infos
- flutter test --coverage
- dart format --set-exit-if-changed lib/ test/
- flutter build apk --debug --no-tree-shake-icons
```

### 8.5 수동 검증 항목

| # | 검증 항목 | 방법 |
|---|----------|------|
| 1 | Admin 설정 화면 진입 | `/admin/settings` 라우트 이동 확인 |
| 2 | 카테고리 탭 전환 | 10개 카테고리 탭 클릭 → 필터 동작 |
| 3 | 설정 검색 | 'toss' 입력 → payment 카테고리 설정만 표시 |
| 4 | 설정 편집 다이얼로그 | 각 valueType별 입력 위젯 표시 확인 |
| 5 | 마크다운 도움말 렌더링 | helpText가 있는 설정의 도움말 펼치기 |
| 6 | 변경 확인 다이얼로그 | 저장 시 변경 사유 입력 다이얼로그 표시 |
| 7 | LLM 채팅 화면 진입 | `/admin/settings/assistant` 라우트 이동 확인 |
| 8 | 채팅 메시지 송수신 | 메시지 입력 → 전송 → 응답 버블 표시 |
| 9 | 제안 적용/거부 | 적용 버튼 → 상태 변경, 거부 버튼 → 상태 변경 |
| 10 | 반응형 레이아웃 | 창 너비 768px 이상/이하에서 레이아웃 전환 |
| 11 | 다크 모드 | 설정 → 다크 모드 → Admin UI 테마 적용 |
| 12 | 다국어 전환 | 설정 → 언어 변경 → Admin UI 텍스트 전환 |

---

## 부록 A: 전체 신규/수정 파일 목록

### 신규 파일 (21개)

| # | 파일 경로 | 설명 |
|---|----------|------|
| 1 | `lib/features/admin/domain/admin_settings_repository.dart` | Admin 설정 Repository 인터페이스 |
| 2 | `lib/features/admin/domain/llm_assistant_repository.dart` | LLM 어시스턴트 Repository 인터페이스 |
| 3 | `lib/features/admin/domain/models/config_item.dart` | ConfigItem 도메인 모델 |
| 4 | `lib/features/admin/domain/models/config_category.dart` | ConfigCategory enum |
| 5 | `lib/features/admin/domain/models/validate_result.dart` | ValidateResult 모델 |
| 6 | `lib/features/admin/domain/models/chat_message.dart` | ChatMessage 모델 |
| 7 | `lib/features/admin/domain/models/config_suggestion.dart` | ConfigSuggestion 모델 |
| 8 | `lib/features/admin/domain/models/assistant_session.dart` | ConfigSessionInfo/AssistantResponse 모델 |
| 9 | `lib/features/admin/data/admin_settings_repository_impl.dart` | gRPC AdminService 연동 구현체 |
| 10 | `lib/features/admin/data/llm_assistant_repository_impl.dart` | gRPC AiInferenceService 연동 구현체 |
| 11 | `lib/features/admin/presentation/admin_settings_screen.dart` | Admin 설정 메인 화면 |
| 12 | `lib/features/admin/presentation/llm_chat_screen.dart` | LLM 채팅 화면 |
| 13 | `lib/features/admin/presentation/providers/admin_settings_provider.dart` | 설정 상태 관리 |
| 14 | `lib/features/admin/presentation/providers/llm_chat_provider.dart` | 채팅 상태 관리 |
| 15 | `lib/features/admin/presentation/widgets/category_tab_bar.dart` | 카테고리 탭 바 |
| 16 | `lib/features/admin/presentation/widgets/config_card.dart` | 설정 카드 |
| 17 | `lib/features/admin/presentation/widgets/config_edit_dialog.dart` | 설정 편집 다이얼로그 |
| 18 | `lib/features/admin/presentation/widgets/config_search_bar.dart` | 검색 바 |
| 19 | `lib/features/admin/presentation/widgets/change_reason_dialog.dart` | 변경 사유 다이얼로그 |
| 20 | `lib/features/admin/presentation/widgets/markdown_help_viewer.dart` | 마크다운 도움말 |
| 21 | `lib/features/admin/presentation/widgets/chat_message_bubble.dart` | 채팅 버블 |

추가 채팅 위젯:

| # | 파일 경로 | 설명 |
|---|----------|------|
| 22 | `lib/features/admin/presentation/widgets/chat_input_bar.dart` | 채팅 입력 바 |
| 23 | `lib/features/admin/presentation/widgets/suggestion_panel.dart` | 제안 패널 |
| 24 | `lib/features/admin/presentation/widgets/suggestion_card.dart` | 제안 카드 |
| 25 | `lib/features/admin/presentation/widgets/language_category_bar.dart` | 언어+카테고리 바 |

### 신규 테스트 파일 (14개)

| # | 파일 경로 | 테스트 수 |
|---|----------|----------|
| 1 | `test/helpers/test_app_wrapper.dart` | (헬퍼) |
| 2 | `test/unit/auth_provider_test.dart` | 8 |
| 3 | `test/unit/theme_provider_test.dart` | 6 |
| 4 | `test/unit/locale_provider_test.dart` | 8 |
| 5 | `test/unit/validators_test.dart` | 10 |
| 6 | `test/unit/app_constants_test.dart` | 4 |
| 7 | `test/unit/auth_interceptor_test.dart` | 4 |
| 8 | `test/unit/admin_settings_provider_test.dart` | 6 |
| 9 | `test/widget/measurement_card_test.dart` | 4 |
| 10 | `test/widget/primary_button_test.dart` | 3 |
| 11 | `test/widget/home_screen_test.dart` | 3 |
| 12 | `test/widget/login_screen_test.dart` | 4 |
| 13 | `test/widget/device_list_screen_test.dart` | 2 |
| 14 | `test/widget/settings_screen_test.dart` | 3 |
| 15 | `test/widget/admin_settings_screen_test.dart` | 4 |
| 16 | `test/widget/config_edit_dialog_test.dart` | 5 |
| 17 | `test/widget/llm_chat_screen_test.dart` | 4 |
| 18 | `test/integration/fake_repository_integration_test.dart` | 4 |
| | **합계** | **82** |

### 수정 파일 (6개)

| # | 파일 경로 | 수정 내용 |
|---|----------|-----------|
| 1 | `lib/core/router/app_router.dart` | +2 라우트 (`/admin/settings`, `/admin/settings/assistant`) |
| 2 | `lib/core/providers/grpc_provider.dart` | +2 Provider (adminSettings, llmAssistant) |
| 3 | `lib/core/services/grpc_client.dart` | +2 채널 (admin, aiInference) |
| 4 | `lib/core/constants/app_constants.dart` | +2 포트 상수 (grpcAdminPort, grpcAiInferencePort) |
| 5 | `lib/l10n/app_localizations.dart` | +50 접근자 getter |
| 6 | `lib/l10n/translations/*.dart` (6파일) | +50 번역 키-값 |
| 7 | `pubspec.yaml` | +3 패키지 (flutter_markdown, url_launcher, badges) |
| 8 | `test/helpers/fake_repositories.dart` | +2 Fake Repository (Admin, LLM) |

---

## 부록 B: 구현 순서 권장

```
Phase 1: 테스트 인프라 + 단위 테스트 (F-1)
  1. test/helpers/test_app_wrapper.dart 생성
  2. test/unit/ 단위 테스트 파일 8개 생성 (auth, theme, locale, validators, constants, interceptor)
  3. test/widget/ 위젯 테스트 파일 6개 생성 (measurement_card, primary_button, home, login, device, settings)
  4. test/integration/ 통합 테스트 1개 생성
  5. flutter test 전체 통과 확인

Phase 2: Admin 설정 UI (AS-6)
  1. 도메인 모델 정의 (config_item, config_category, validate_result)
  2. Repository 인터페이스 정의 (admin_settings_repository)
  3. Fake Repository 추가 (test/helpers/)
  4. Provider 구현 (admin_settings_provider)
  5. 위젯 구현 (config_card → config_edit_dialog → admin_settings_screen)
  6. gRPC 연동 구현체 (admin_settings_repository_impl)
  7. 라우트 추가
  8. 다국어 키 추가
  9. AS-6 테스트 추가 (admin_settings_screen_test, config_edit_dialog_test)

Phase 3: LLM 채팅 UI (AS-8)
  1. 도메인 모델 정의 (chat_message, config_suggestion, assistant_session)
  2. Repository 인터페이스 정의 (llm_assistant_repository)
  3. Fake Repository 추가
  4. Provider 구현 (llm_chat_provider)
  5. 위젯 구현 (chat_message_bubble → chat_input_bar → suggestion_card → llm_chat_screen)
  6. gRPC 연동 구현체 (llm_assistant_repository_impl)
  7. 라우트 추가
  8. 다국어 키 추가
  9. AS-8 테스트 추가 (llm_chat_screen_test)

Phase 4: 통합 검증
  1. flutter analyze 통과
  2. flutter test 전체 통과 (82개+ 테스트)
  3. dart format 통과
  4. flutter build apk --debug 성공
  5. 수동 검증 12항목 확인
```

---

> **문서 끝**  
> 구현 시 이 기획서를 참조하되, 실제 코드 작성 중 발견되는 이슈는 KNOWN_ISSUES.md에 기록하고 기획서를 업데이트한다.
