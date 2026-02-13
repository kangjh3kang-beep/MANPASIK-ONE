# Flutter UI 구현 스펙 (ChatGPT용)

## 프로젝트 정보

- **프로젝트명**: ManPaSik (만파식)
- **앱 유형**: 크로스플랫폼 (iOS, Android, Desktop)
- **프레임워크**: Flutter 3.x
- **상태관리**: Riverpod 2.x
- **로컬 DB**: SQLite (sqflite) + Hive

---

## 요청 작업

### 1. 프로젝트 초기 설정

**pubspec.yaml 생성:**
```yaml
name: manpasik
description: 만파식 헬스케어 AI 생태계
version: 1.0.0+1

environment:
  sdk: '>=3.2.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter
  flutter_localizations:
    sdk: flutter
    
  # 상태관리
  flutter_riverpod: ^2.4.9
  riverpod_annotation: ^2.3.3
  
  # 라우팅
  go_router: ^13.0.0
  
  # 네트워크
  dio: ^5.4.0
  
  # 로컬 저장소
  sqflite: ^2.3.0
  hive_flutter: ^1.1.0
  
  # UI
  flutter_svg: ^2.0.9
  cached_network_image: ^3.3.1
  shimmer: ^3.0.0
  
  # 유틸리티
  freezed_annotation: ^2.4.1
  json_annotation: ^4.8.1
  intl: ^0.19.0
  
  # BLE/NFC (Rust FFI로 대체 예정)
  flutter_blue_plus: ^1.31.0
  
dev_dependencies:
  flutter_test:
    sdk: flutter
  build_runner: ^2.4.8
  freezed: ^2.4.6
  json_serializable: ^6.7.1
  riverpod_generator: ^2.3.9
  flutter_lints: ^3.0.1

flutter:
  uses-material-design: true
  generate: true
  
  assets:
    - assets/images/
    - assets/icons/
    - assets/fonts/
```

---

### 2. 앱 진입점 (main.dart)

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/router/app_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // TODO: Rust FFI 초기화
  // await ManpasikEngine.initialize();
  
  runApp(
    const ProviderScope(
      child: ManpasikApp(),
    ),
  );
}

class ManpasikApp extends ConsumerWidget {
  const ManpasikApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    
    return MaterialApp.router(
      title: '만파식',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.system,
      routerConfig: router,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [
        Locale('ko', 'KR'),
        Locale('en', 'US'),
        Locale('ja', 'JP'),
        Locale('zh', 'CN'),
      ],
    );
  }
}
```

---

### 3. 필요한 화면 목록

| 화면 | 경로 | 우선순위 |
|------|------|----------|
| 스플래시 | `/splash` | P0 |
| 로그인 | `/login` | P0 |
| 회원가입 | `/register` | P0 |
| 홈 대시보드 | `/home` | P0 |
| 측정 | `/measure` | P0 |
| 측정 결과 | `/measure/result` | P0 |
| 데이터 허브 | `/data-hub` | P1 |
| AI 코치 | `/coach` | P1 |
| 마켓 | `/market` | P2 |
| 커뮤니티 | `/community` | P2 |
| 설정 | `/settings` | P1 |
| 디바이스 관리 | `/devices` | P0 |

---

### 4. 디자인 시스템

**색상 팔레트:**
```dart
// Primary: 만파식 블루그린 (치유, 생명)
static const primary = Color(0xFF00897B);
static const primaryLight = Color(0xFF4EBAAA);
static const primaryDark = Color(0xFF005B4F);

// Secondary: 골드 (프리미엄, 신뢰)
static const secondary = Color(0xFFFFB300);

// Background
static const backgroundLight = Color(0xFFF5F7FA);
static const backgroundDark = Color(0xFF121212);

// 상태 색상
static const success = Color(0xFF4CAF50);
static const warning = Color(0xFFFF9800);
static const error = Color(0xFFE53935);
static const info = Color(0xFF2196F3);
```

---

### 5. 폴더 구조

```
lib/
├── core/
│   ├── router/
│   │   └── app_router.dart
│   ├── theme/
│   │   └── app_theme.dart
│   ├── constants/
│   ├── utils/
│   └── services/
├── features/
│   ├── auth/
│   │   ├── data/
│   │   ├── domain/
│   │   └── presentation/
│   ├── home/
│   ├── measurement/
│   ├── data_hub/
│   ├── ai_coach/
│   ├── market/
│   ├── community/
│   ├── devices/
│   └── settings/
├── shared/
│   ├── widgets/
│   ├── models/
│   └── providers/
└── l10n/
    ├── app_ko.arb
    └── app_en.arb
```

---

## 산출물

1. `pubspec.yaml`
2. `lib/main.dart`
3. `lib/core/router/app_router.dart`
4. `lib/core/theme/app_theme.dart`
5. `lib/features/home/presentation/home_screen.dart`
6. `lib/features/measurement/presentation/measurement_screen.dart`

---

**작성자**: Antigravity (Gemini)
**작성일**: 2026-02-09
