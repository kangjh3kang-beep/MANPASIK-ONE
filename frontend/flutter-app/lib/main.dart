import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:manpasik/core/router/app_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/core/services/crash_reporter.dart';
import 'package:manpasik/core/services/app_logger.dart';
import 'package:manpasik/core/services/app_lifecycle_observer.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/l10n/app_localizations.dart';

/// 전역 크래시 리포터 인스턴스
final CrashReporter crashReporter = ConsoleCrashReporter();

void main() async {
  // Zone.runGuarded로 비동기 에러 포착
  runZonedGuarded<Future<void>>(() async {
    WidgetsFlutterBinding.ensureInitialized();

    final log = AppLogger.instance;
    log.info('ManPaSik 앱 시작', tag: 'Bootstrap');

    // Flutter 프레임워크 에러 핸들러
    FlutterError.onError = (FlutterErrorDetails details) {
      log.error(
        'Flutter framework error: ${details.exceptionAsString()}',
        tag: 'FlutterError',
        error: details.exception,
        stackTrace: details.stack,
      );
      crashReporter.reportError(
        details.exception,
        details.stack,
        context: details.context?.toString(),
      );
    };

    // 플랫폼 디스패처 에러 핸들러 (Flutter 3.10+)
    PlatformDispatcher.instance.onError = (error, stack) {
      log.error(
        'Platform dispatcher error',
        tag: 'PlatformError',
        error: error,
        stackTrace: stack,
      );
      crashReporter.reportFatalError(error, stack);
      return true;
    };

    // Hive 초기화 (로컬 저장소)
    // await Hive.initFlutter();

    // Rust Core 엔진 초기화 (스텁/네이티브 자동 선택)
    await RustBridge.init();
    log.info('Rust Bridge 초기화 완료', tag: 'Bootstrap');

    // Firebase/FCM 초기화 (firebase_core, firebase_messaging 패키지 설치 후 활성화)
    // await Firebase.initializeApp();
    // final messaging = FirebaseMessaging.instance;
    // await messaging.requestPermission();
    // FirebaseMessaging.onBackgroundMessage(_firebaseMessagingBackgroundHandler);

    runApp(
      const ProviderScope(
        child: ManpasikApp(),
      ),
    );

    log.info('ManPaSik 앱 부트스트랩 완료', tag: 'Bootstrap');
  }, (error, stackTrace) {
    // Zone 바깥 비동기 에러 포착
    AppLogger.instance.error(
      'Uncaught async error',
      tag: 'ZoneError',
      error: error,
      stackTrace: stackTrace,
    );
    crashReporter.reportFatalError(error, stackTrace);
  });
}

/// ManPaSik 앱 루트 위젯
///
/// Material Design 3 기반, 다국어 6개 언어 지원, 다크모드 지원
class ManpasikApp extends ConsumerStatefulWidget {
  const ManpasikApp({super.key});

  @override
  ConsumerState<ManpasikApp> createState() => _ManpasikAppState();
}

class _ManpasikAppState extends ConsumerState<ManpasikApp> {
  late final AppLifecycleObserver _lifecycleObserver;

  @override
  void initState() {
    super.initState();
    _lifecycleObserver = AppLifecycleObserver(
      onResumed: _onAppResumed,
      onPaused: _onAppPaused,
    );
    _lifecycleObserver.register();
  }

  @override
  void dispose() {
    _lifecycleObserver.unregister();
    super.dispose();
  }

  void _onAppResumed() {
    AppLogger.instance.info('Foreground 복귀 — 데이터 갱신', tag: 'Lifecycle');
    // 토큰 유효성 검사, 동기화 트리거 등
  }

  void _onAppPaused() {
    AppLogger.instance.info('Background 진입 — 리소스 정리', tag: 'Lifecycle');
    // 불필요한 스트림/타이머 일시정지 등
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(appRouterProvider);
    final themeMode = ref.watch(themeModeProvider);
    final locale = ref.watch(localeProvider);

    return MaterialApp.router(
      title: 'MANPASIK Measurement System',
      debugShowCheckedModeBanner: false,

      // 테마
      theme: AppTheme.koreanWhite, // Korean White Mode (Baekja)
      darkTheme: AppTheme.dark,
      themeMode: themeMode, // Dynamic Theme Switching

      // 라우터
      routerConfig: router,

      // 다국어 (6개 언어: ko, en, ja, zh, fr, hi)
      locale: locale,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: AppLocalizations.supportedLocales,
    );
  }
}
