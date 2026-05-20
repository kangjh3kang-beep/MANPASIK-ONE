// lib/core/router/app_router.dart
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../features/home/presentation/home_screen.dart';
import '../../features/measure/presentation/measurement_screen.dart';
import '../../features/home/presentation/diagnostics_screen.dart';
import '../../features/data/presentation/data_screen.dart';
import '../../features/ai_coach/presentation/ai_coach_screen.dart';
import '../../features/market/presentation/market_screen.dart';

// (Mock) 앱 상태 참조를 위해 추후 global_providers.dart 와 연동
final authStateProvider = StateProvider<bool>((ref) => false); 

/// ManPaSik 딥링크 및 150여개의 화면 구성을 관장하는 중앙 시스템 라우터
final appRouterProvider = Provider<GoRouter>((ref) {
  final isAuthenticated = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: '/splash',
    redirect: (context, state) {
      final isSplash = state.uri.path == '/splash';
      final isAuth = state.uri.path.startsWith('/auth');

      // 미인증 시 로그인 화면 강제, 인증된 상태라면 홈으로 이동
      if (!isAuthenticated && !isAuth && !isSplash) {
        return '/auth/login';
      }
      return null;
    },
    routes: [
      GoRoute(
        path: '/splash',
        builder: (context, state) => const Scaffold(body: Center(child: Text('S-000 Splash'))),
      ),
      GoRoute(
        path: '/auth/login',
        builder: (context, state) => const Scaffold(body: Center(child: Text('A-010 Login'))),
      ),
      // 5-Tab 메인 화면 네비게이션 셸 구조
      ShellRoute(
        builder: (context, state, child) {
          final loc = state.uri.path;
          int idx = 0;
          if (loc.startsWith('/measure')) idx = 1;
          else if (loc.startsWith('/data')) idx = 2;
          else if (loc.startsWith('/ai_coach')) idx = 3;
          else if (loc.startsWith('/market')) idx = 4;

          return Scaffold(
            body: child,
            bottomNavigationBar: BottomNavigationBar(
              backgroundColor: const Color(0xFF0F172A),
              selectedItemColor: const Color(0xFF2DD4BF),
              unselectedItemColor: Colors.white54,
              type: BottomNavigationBarType.fixed,
              currentIndex: idx,
              onTap: (index) {
                if (index == 0) context.go('/home');
                if (index == 1) context.go('/measure');
                if (index == 2) context.go('/data');
                if (index == 3) context.go('/ai_coach');
                if (index == 4) context.go('/market');
              },
            items: const [
              BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
              BottomNavigationBarItem(icon: Icon(Icons.science), label: 'Measure'),
              BottomNavigationBarItem(icon: Icon(Icons.data_usage), label: 'Data'),
              BottomNavigationBarItem(icon: Icon(Icons.chat), label: 'AI Coach'),
              BottomNavigationBarItem(icon: Icon(Icons.shop), label: 'Market'),
            ],
          ),
        );
      },
      routes: [
          GoRoute(
            path: '/home',
            builder: (context, state) => const HomeScreen(),
            routes: [
              GoRoute(
                path: 'diagnostics', // H-035 진단센터 (v1.3)
                builder: (context, state) => const DiagnosticsScreen(),
              ),
              GoRoute(
                path: 'healing_log', // H-038 자가치유 로그 뷰어 (v1.4)
                builder: (context, state) => const Scaffold(body: Center(child: Text('H-038 Healing Logs'))),
              ),
            ]
          ),
          GoRoute(
            path: '/measure',
            builder: (context, state) => const MeasurementScreen(),
            routes: [
              GoRoute(
                path: 'progress', // M-040 측정 진행
                builder: (context, state) => const MeasurementScreen(),
              ),
              GoRoute(
                path: 'restore', // OV-084 세션 복원 (v1.4) 딥링크
                builder: (context, state) => const Scaffold(body: Center(child: Text('OV-084 Session Restore'))),
              ),
            ]
          ),
          GoRoute(
            path: '/data',
            builder: (context, state) => const DataScreen(),
          ),
          GoRoute(
            path: '/ai_coach',
            builder: (context, state) => const AiCoachScreen(),
          ),
          GoRoute(
            path: '/market',
            builder: (context, state) => const MarketScreen(),
          )
        ],
      )
    ],
  );
});
