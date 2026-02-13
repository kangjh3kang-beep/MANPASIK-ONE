import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/features/auth/presentation/splash_screen.dart';
import 'package:manpasik/features/auth/presentation/login_screen.dart';
import 'package:manpasik/features/auth/presentation/register_screen.dart';
import 'package:manpasik/features/home/presentation/home_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_result_screen.dart';
import 'package:manpasik/features/devices/presentation/device_list_screen.dart';
import 'package:manpasik/features/settings/presentation/settings_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_settings_screen.dart';
import 'package:manpasik/features/chat/presentation/chat_screen.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// ManPaSik 앱 라우터 Provider
///
/// P0 라우트 + chat: splash, login, register, home, measure, devices, settings, chat
/// 인증 상태에 따른 리다이렉트 처리 포함
final appRouterProvider = Provider<GoRouter>((ref) {
  // routerConfig가 재생성되지 않도록 ref.watch(authProvider) 제거
  final authNotifier = ref.watch(authProvider.notifier);

  return GoRouter(
    initialLocation: '/',
    debugLogDiagnostics: true,
    // 인증 상태 변경 감지
    refreshListenable: GoRouterRefreshStream(authNotifier.stream),
    redirect: (context, state) {
      // 리다이렉트 시점에 최신 상태 읽기
      final authState = ref.read(authProvider);
      final isLoggedIn = authState.isAuthenticated;
      final isAuthRoute = state.matchedLocation == '/login' ||
          state.matchedLocation == '/register';
      final isSplash = state.matchedLocation == '/';

      // 스플래시 화면에서는 리다이렉트 없음
      if (isSplash) return null;

      // 미인증 상태에서 인증 라우트 외 접근 시 → 로그인
      if (!isLoggedIn && !isAuthRoute) return '/login';

      // 인증 완료 상태에서 인증 라우트 접근 시 → 홈
      if (isLoggedIn && isAuthRoute) return '/home';

      return null;
    },
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/measurement',
        builder: (context, state) => const MeasurementScreen(),
      ),
      GoRoute(
        path: '/measurement/result',
        builder: (context, state) => const MeasurementResultScreen(),
      ),
      GoRoute(
        path: '/devices',
        builder: (context, state) => const DeviceListScreen(),
      ),
      GoRoute(
        path: '/settings',
        builder: (context, state) => const SettingsScreen(),
      ),
      GoRoute(
        path: '/admin/settings',
        builder: (context, state) => const AdminSettingsScreen(),
      ),
      GoRoute(
        path: '/chat',
        builder: (context, state) => const ChatScreen(),
      ),
    ],
  );
});

/// Stream을 Listenable로 변환하는 클래스
class GoRouterRefreshStream extends ChangeNotifier {
  GoRouterRefreshStream(Stream<dynamic> stream) {
    notifyListeners();
    _subscription = stream.asBroadcastStream().listen(
          (dynamic _) => notifyListeners(),
        );
  }

  late final StreamSubscription<dynamic> _subscription;

  @override
  void dispose() {
    _subscription.cancel();
    super.dispose();
  }
}
