import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/shared/widgets/glass_dock_navigation.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/floating_monitor_bubble.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/hanji_background.dart'; // Added Import

import 'package:manpasik/features/auth/presentation/splash_screen.dart';
import 'package:manpasik/features/auth/presentation/login_screen.dart';
import 'package:manpasik/features/auth/presentation/register_screen.dart';
import 'package:manpasik/features/onboarding/presentation/onboarding_screen.dart';
import 'package:manpasik/features/home/presentation/home_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_result_screen.dart';
import 'package:manpasik/features/devices/presentation/device_list_screen.dart';
import 'package:manpasik/features/settings/presentation/settings_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_settings_screen.dart';
import 'package:manpasik/features/chat/presentation/chat_screen.dart';
import 'package:manpasik/features/data_hub/presentation/data_hub_screen.dart';
import 'package:manpasik/features/data_hub/presentation/monitoring_dashboard_screen.dart'; // Added Import
import 'package:manpasik/features/ai_coach/presentation/ai_coach_screen.dart';
import 'package:manpasik/features/market/presentation/market_screen.dart';
import 'package:manpasik/features/community/presentation/community_screen.dart';
import 'package:manpasik/features/medical/presentation/medical_screen.dart';
import 'package:manpasik/features/family/presentation/family_screen.dart';
import 'package:manpasik/features/notification/presentation/notification_screen.dart';
import 'package:manpasik/features/auth/presentation/forgot_password_screen.dart';
import 'package:manpasik/features/settings/presentation/support_screen.dart';
import 'package:manpasik/features/settings/presentation/legal_screen.dart';
import 'package:manpasik/features/settings/presentation/emergency_settings_screen.dart';
import 'package:manpasik/features/settings/presentation/consent_management_screen.dart';
import 'package:manpasik/features/settings/presentation/profile_edit_screen.dart';
import 'package:manpasik/features/settings/presentation/security_screen.dart';
import 'package:manpasik/features/settings/presentation/accessibility_screen.dart';
import 'package:manpasik/features/settings/presentation/notification_settings_screen.dart';
import 'package:manpasik/features/settings/presentation/notice_screen.dart';
import 'package:manpasik/features/market/presentation/cartridge_detail_screen.dart';
import 'package:manpasik/features/market/presentation/encyclopedia_screen.dart';
import 'package:manpasik/features/market/presentation/product_detail_screen.dart';
import 'package:manpasik/features/market/presentation/cart_screen.dart';
import 'package:manpasik/features/market/presentation/order_history_screen.dart';
import 'package:manpasik/features/market/presentation/subscription_screen.dart';
import 'package:manpasik/features/market/presentation/checkout_screen.dart';
import 'package:manpasik/features/market/presentation/order_complete_screen.dart';
import 'package:manpasik/features/ai_coach/presentation/food_analysis_screen.dart';
import 'package:manpasik/features/ai_coach/presentation/exercise_video_screen.dart';
import 'package:manpasik/features/community/presentation/post_detail_screen.dart';
import 'package:manpasik/features/medical/presentation/facility_search_screen.dart';
import 'package:manpasik/features/medical/presentation/prescription_detail_screen.dart';
import 'package:manpasik/features/medical/presentation/telemedicine_screen.dart';
import 'package:manpasik/features/medical/presentation/video_call_screen.dart';
import 'package:manpasik/features/family/presentation/family_report_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_dashboard_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_users_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_audit_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_monitor_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_hierarchy_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_compliance_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_revenue_screen.dart';
import 'package:manpasik/features/admin/presentation/admin_inventory_table.dart';
import 'package:manpasik/features/community/presentation/research_post_screen.dart';
import 'package:manpasik/features/devices/presentation/device_detail_screen.dart';
import 'package:manpasik/features/community/presentation/create_post_screen.dart';
import 'package:manpasik/features/community/presentation/challenge_screen.dart';
import 'package:manpasik/features/community/presentation/qna_screen.dart';
import 'package:manpasik/features/family/presentation/family_create_screen.dart';
import 'package:manpasik/features/family/presentation/member_edit_screen.dart';
import 'package:manpasik/features/family/presentation/guardian_dashboard_screen.dart';
import 'package:manpasik/features/family/presentation/alert_detail_screen.dart';
import 'package:manpasik/features/medical/presentation/consultation_result_screen.dart';
import 'package:manpasik/features/market/presentation/order_detail_screen.dart';
import 'package:manpasik/features/market/presentation/plan_comparison_screen.dart';
import 'package:manpasik/features/settings/presentation/inquiry_create_screen.dart';
import 'package:manpasik/features/settings/presentation/escalation_screen.dart';
import 'package:manpasik/features/settings/presentation/data_export_screen.dart';
import 'package:manpasik/features/market/presentation/subscription_cancel_screen.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/network/conflict_resolver_screen.dart';
import 'package:manpasik/shared/widgets/network_indicator.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 글로벌 네비게이터 키 (ShellRoute 외부 라우트용)
final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

/// ManPaSik 앱 라우터 Provider
///
/// 사이트맵(MPK-UX-SITEMAP-v1.0) 기준 전체 라우트 등록:
/// - 인증: /, /login, /register
/// - 온보딩: /onboarding
/// - 메인(BottomNav): /home, /data, /measure, /market, /settings
/// - 서브: /coach, /community, /medical, /devices, /family, /admin/*
final appRouterProvider = Provider<GoRouter>((ref) {
  final authNotifier = ref.watch(authProvider.notifier);

  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/',
    debugLogDiagnostics: true,
    refreshListenable: GoRouterRefreshStream(authNotifier.stream),
    redirect: (context, state) {
      final authState = ref.read(authProvider);
      final isLoggedIn = authState.isAuthenticated;
      final loc = state.matchedLocation;
      final isAuthRoute = loc == '/login' || loc == '/register';
      final isSplash = loc == '/';
      final isOnboarding = loc == '/onboarding';

      if (isSplash) return null;
      if (!isLoggedIn && !isAuthRoute) return '/login';
      if (isLoggedIn && isAuthRoute) return '/home';
      if (isOnboarding && !isLoggedIn) return '/login';

      // Admin RBAC Guard: /admin/* 경로는 관리자 역할 필요
      if (loc.startsWith('/admin') && !authState.isAdmin) return '/home';

      return null;
    },
    routes: [
      // ── 인증 외 라우트 (루트 네비게이터) ──
      GoRoute(
        path: '/',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: '/login',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const RegisterScreen(),
      ),

      // ── 온보딩 (회원가입 후 첫 실행) ──
      GoRoute(
        path: '/onboarding',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const OnboardingScreen(),
      ),

      // ── 메인 앱 (BottomNav ShellRoute) ──
      ShellRoute(
        navigatorKey: _shellNavigatorKey,
        builder: (context, state, child) =>
            ScaffoldWithBottomNav(child: child),
        routes: [
          // 탭 1: 홈
          GoRoute(
            path: '/home',
            builder: (context, state) => const HomeScreen(),
          ),
          // 탭 2: 데이터 허브 (Phase 2)
          GoRoute(
            path: '/data',
            builder: (context, state) => const DataHubScreen(),
            routes: [
              GoRoute(
                path: 'monitoring',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const MonitoringDashboardScreen(),
              ),
            ],
          ),
          // 탭 3: 측정 (사이트맵: /measure)
          GoRoute(
            path: '/measure',
            builder: (context, state) => const MeasurementScreen(),
            routes: [
              GoRoute(
                path: 'result',
                parentNavigatorKey: _rootNavigatorKey,
                builder: (context, state) => const MeasurementResultScreen(),
              ),
            ],
          ),
          // 탭 4: 마켓 (Phase 2)
          GoRoute(
            path: '/market',
            builder: (context, state) => const MarketScreen(),
          ),
          // 탭 5: 설정
          GoRoute(
            path: '/settings',
            builder: (context, state) => const SettingsScreen(),
          ),
        ],
      ),

      // ── 서브 라우트 (루트 네비게이터 - 풀스크린) ──

      // AI 코치 (Phase 2)
      GoRoute(
        path: '/coach',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AiCoachScreen(),
      ),
      // AI 채팅 (코치 하위 기능)
      GoRoute(
        path: '/chat',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ChatScreen(),
      ),
      // 커뮤니티 (Phase 3)
      GoRoute(
        path: '/community',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const CommunityScreen(),
      ),
      // 의료 서비스 (Phase 3)
      GoRoute(
        path: '/medical',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const MedicalScreen(),
      ),
      // 기기 관리 (Phase 1)
      GoRoute(
        path: '/devices',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const DeviceListScreen(),
      ),
      GoRoute(
        path: '/devices/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => DeviceDetailScreen(
          deviceId: state.pathParameters['id'] ?? '',
        ),
      ),
      // 가족 관리 (Phase 3)
      GoRoute(
        path: '/family',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FamilyScreen(),
      ),
      // 알림 센터
      GoRoute(
        path: '/notifications',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const NotificationScreen(),
      ),
      // 관리자 포탈 (Phase 3)
      GoRoute(
        path: '/admin/settings',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminSettingsScreen(),
      ),
      GoRoute(
        path: '/admin/dashboard',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminDashboardScreen(),
      ),
      GoRoute(
        path: '/admin/users',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminUsersScreen(),
      ),
      GoRoute(
        path: '/admin/audit',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminAuditScreen(),
      ),
      GoRoute(
        path: '/admin/monitor',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminMonitorScreen(),
      ),
      GoRoute(
        path: '/admin/emergency',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminMonitorScreen(tab: 'emergency'),
      ),
      GoRoute(
        path: '/admin/hierarchy',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminHierarchyScreen(),
      ),
      GoRoute(
        path: '/admin/compliance',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminComplianceScreen(),
      ),
      GoRoute(
        path: '/admin/revenue',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminRevenueScreen(),
      ),
      GoRoute(
        path: '/admin/inventory',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AdminInventoryTable(),
      ),

      // ── 인증 보조 ──
      GoRoute(
        path: '/forgot-password',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ForgotPasswordScreen(),
      ),

      // ── 설정 하위 라우트 ──
      GoRoute(
        path: '/support',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SupportScreen(),
      ),
      GoRoute(
        path: '/settings/terms',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const LegalScreen(type: 'terms'),
      ),
      GoRoute(
        path: '/settings/privacy',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const LegalScreen(type: 'privacy'),
      ),
      GoRoute(
        path: '/settings/consent',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ConsentManagementScreen(),
      ),
      GoRoute(
        path: '/settings/escalation',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const EscalationScreen(),
      ),
      GoRoute(
        path: '/settings/data-export',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const DataExportScreen(),
      ),
      GoRoute(
        path: '/market/subscription/cancel',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SubscriptionCancelScreen(),
      ),
      GoRoute(
        path: '/settings/emergency',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const EmergencySettingsScreen(),
      ),
      GoRoute(
        path: '/settings/profile',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ProfileEditScreen(),
      ),
      GoRoute(
        path: '/settings/security',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SecurityScreen(),
      ),
      GoRoute(
        path: '/settings/notifications',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const NotificationSettingsScreen(),
      ),
      GoRoute(
        path: '/settings/accessibility',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const AccessibilityScreen(),
      ),
      GoRoute(
        path: '/support/notices',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const NoticeScreen(),
      ),
      GoRoute(
        path: '/settings/inquiry/create',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const InquiryCreateScreen(),
      ),

      // ── 마켓 하위 라우트 (Phase 2) ──
      GoRoute(
        path: '/market/encyclopedia',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const EncyclopediaScreen(),
      ),
      GoRoute(
        path: '/market/encyclopedia/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => CartridgeDetailScreen(
          cartridgeId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/market/product/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => ProductDetailScreen(
          productId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/market/cart',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const CartScreen(),
      ),
      GoRoute(
        path: '/market/orders',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const OrderHistoryScreen(),
      ),
      GoRoute(
        path: '/market/subscription',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const SubscriptionScreen(),
      ),
      GoRoute(
        path: '/market/checkout',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const CheckoutScreen(),
      ),
      GoRoute(
        path: '/market/order-complete/:orderId',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => OrderCompleteScreen(
          orderId: state.pathParameters['orderId'] ?? '',
        ),
      ),
      GoRoute(
        path: '/market/order/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => OrderDetailScreen(
          orderId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/market/subscription/plans',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const PlanComparisonScreen(),
      ),
      GoRoute(
        path: '/market/subscription/upgrade',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const PlanComparisonScreen(mode: 'upgrade'),
      ),
      GoRoute(
        path: '/market/subscription/downgrade',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const PlanComparisonScreen(mode: 'downgrade'),
      ),

      // ── AI 코치 하위 (Phase 2) ──
      GoRoute(
        path: '/coach/food',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FoodAnalysisScreen(),
      ),
      GoRoute(
        path: '/coach/exercise-video',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ExerciseVideoScreen(),
      ),

      // ── 커뮤니티 하위 (Phase 3) ──
      GoRoute(
        path: '/community/post/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => PostDetailScreen(
          postId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/community/create',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const CreatePostScreen(),
      ),
      GoRoute(
        path: '/community/challenge',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ChallengeScreen(),
      ),
      GoRoute(
        path: '/community/challenge/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => ChallengeScreen(
          challengeId: state.pathParameters['id'],
        ),
      ),
      GoRoute(
        path: '/community/qna',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const QnaScreen(),
      ),
      GoRoute(
        path: '/community/qna/ask',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const QnaScreen(mode: 'ask'),
      ),
      GoRoute(
        path: '/community/research',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ResearchPostScreen(),
      ),

      // ── 의료 하위 (Phase 3) ──
      GoRoute(
        path: '/medical/facility-search',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FacilitySearchScreen(),
      ),
      GoRoute(
        path: '/medical/prescription/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => PrescriptionDetailScreen(
          prescriptionId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/medical/telemedicine',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const TelemedicineScreen(),
      ),
      GoRoute(
        path: '/medical/video-call/:sessionId',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => VideoCallScreen(
          sessionId: state.pathParameters['sessionId'] ?? '',
        ),
      ),
      GoRoute(
        path: '/medical/consultation/:id/result',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => ConsultationResultScreen(
          consultationId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/medical/pharmacy',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FacilitySearchScreen(),
      ),

      // ── 가족 하위 (Phase 3) ──
      GoRoute(
        path: '/family/report',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FamilyReportScreen(),
      ),
      GoRoute(
        path: '/family/create',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FamilyCreateScreen(),
      ),
      GoRoute(
        path: '/family/invite',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const FamilyCreateScreen(mode: 'invite'),
      ),
      GoRoute(
        path: '/family/member/:id/edit',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => MemberEditScreen(
          memberId: state.pathParameters['id'] ?? '',
        ),
      ),
      GoRoute(
        path: '/family/guardian',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const GuardianDashboardScreen(),
      ),
      GoRoute(
        path: '/family/alert/:id',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => AlertDetailScreen(
          alertId: state.pathParameters['id'] ?? '',
        ),
      ),

      // ── 오프라인 충돌 해결 ──
      GoRoute(
        path: '/conflict-resolve',
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) => const ConflictResolverScreen(),
      ),
    ],
  );
});

/// 하단 네비게이션 바 포함 Scaffold (Glass Dock Ver.)
class ScaffoldWithBottomNav extends StatelessWidget {
  const ScaffoldWithBottomNav({super.key, required this.child});

  final Widget child;

  static const _tabs = [
    '/home',
    '/data',
    '/measure',
    '/market',
    '/settings',
  ];

  int _currentIndex(BuildContext context) {
    final location = GoRouterState.of(context).matchedLocation;
    for (var i = 0; i < _tabs.length; i++) {
      if (location.startsWith(_tabs[i])) return i;
    }
    return 0;
  }

  @override
  Widget build(BuildContext context) {
    final currentIndex = _currentIndex(context);
    final isDark = Theme.of(context).brightness == Brightness.dark;

    // Global Theme Override for the shared scaffold
    // This allows the Glass Dock to feel integrated
    return Scaffold(
      extendBody: true, // Key for floating dock effect
      // 2. Global Premium Background (Stacked)
      // Switch between Cosmic (Dark) and Hanji (Light)
      // 2. Global Premium Background (Stacked)
      // Switch between Cosmic (Dark) and Hanji (Light)
      body: isDark
          ? RoyalCloudBackground(
              child: Stack(
                children: [
                  Column(
                    children: [
                      const NetworkIndicator(),
                      Expanded(child: child),
                    ],
                  ),
                  const FloatingMonitorBubble(),
                ],
              ),
            )
          : HanjiBackground(
              child: Stack(
                children: [
                  Column(
                    children: [
                       const NetworkIndicator(),
                       Expanded(child: child),
                    ],
                  ),
                  const FloatingMonitorBubble(),
                ],
              ),
            ),
      bottomNavigationBar: GlassDockNavigation(
        currentIndex: currentIndex,
        onTap: (index) => context.go(_tabs[index]),
      ),
    );
  }
}

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
