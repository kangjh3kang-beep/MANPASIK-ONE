import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

void main() {
  // 상태 관리 트리 최상단 래핑
  runApp(const ProviderScope(child: ManpasikApp()));
}

class ManpasikApp extends ConsumerWidget {
  const ManpasikApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 생성해둔 라우터 로드
    final router = ref.watch(appRouterProvider);

    return MaterialApp.router(
      title: 'ManPaSik (萬波息)',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,  // 다크 모드 프리미엄 테마 적용
      routerConfig: router,
    );
  }
}
