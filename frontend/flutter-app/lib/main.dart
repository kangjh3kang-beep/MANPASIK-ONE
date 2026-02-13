import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:manpasik/core/router/app_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/l10n/app_localizations.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Hive 초기화 (로컬 저장소)
  // await Hive.initFlutter();

  // Rust Core 초기화 (S5에서 활성화)
  // await RustBridge.init();

  runApp(
    const ProviderScope(
      child: ManpasikApp(),
    ),
  );
}

/// ManPaSik 앱 루트 위젯
///
/// Material Design 3 기반, 다국어 6개 언어 지원, 다크모드 지원
class ManpasikApp extends ConsumerWidget {
  const ManpasikApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    final themeMode = ref.watch(themeModeProvider);
    final locale = ref.watch(localeProvider);

    return MaterialApp.router(
      title: 'MANPASIK Measurement System',
      debugShowCheckedModeBanner: false,

      // 테마
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: themeMode,

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
