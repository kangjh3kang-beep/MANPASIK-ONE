import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/wave_ripple_painter.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:lottie/lottie.dart';

/// 스플래시 화면
///
/// 앱 초기화 + 인증 상태 확인 후 자동 리다이렉트.
/// - 인증됨 → /home
/// - 비인증 → /login
///
/// Wave Ripple 배경 효과로 만파식적 정체성 강화.
class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen>
    with TickerProviderStateMixin {
  late AnimationController _fadeController;
  late Animation<double> _fadeAnimation;
  late AnimationController _waveController;

  @override
  void initState() {
    super.initState();
    _fadeController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    );
    _fadeAnimation = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _fadeController, curve: Curves.easeIn),
    );
    _fadeController.forward();

    _waveController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 4),
    )..repeat();

    _initializeApp();
  }

  Future<void> _initializeApp() async {
    // 인증 상태 확인
    await ref.read(authProvider.notifier).checkAuthStatus();

    // 최소 스플래시 표시 시간
    await Future.delayed(const Duration(seconds: 2));

    if (!mounted) return;

    final isAuthenticated = ref.read(authProvider).isAuthenticated;
    if (isAuthenticated) {
      context.go('/home');
    } else {
      context.go('/login');
    }
  }

  @override
  void dispose() {
    _fadeController.dispose();
    _waveController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      body: CosmicBackground(
        child: SizedBox(
          width: double.infinity,
          height: double.infinity,
          child: Stack(
            children: [
              // Wave Ripple 배경
              Positioned.fill(
                child: AnimatedBuilder(
                  animation: _waveController,
                  builder: (context, _) {
                    return CustomPaint(
                      painter: WaveRipplePainter(
                        animationValue: _waveController.value,
                        rippleCount: 6,
                        primaryColor: AppTheme.waveCyan.withOpacity(0.5),
                        secondaryColor: AppTheme.sanggamGold.withOpacity(0.3),
                      ),
                    );
                  },
                ),
              ),
              // 전경 콘텐츠
              Center(
                child: FadeTransition(
                  opacity: _fadeAnimation,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      // 앱 아이콘
                      Container(
                        width: 120,
                        height: 120,
                        decoration: BoxDecoration(
                          color: Colors.white.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(32),
                          border: Border.all(
                            color: AppTheme.sanggamGold.withOpacity(0.3),
                            width: 1,
                          ),
                        ),
                        child: const Icon(
                          Icons.biotech_rounded,
                          size: 64,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(height: 32),
                      // 브랜드명 (BRAND_GUIDELINE 준수)
                      Text(
                        'MANPASIK',
                        style: theme.textTheme.headlineLarge?.copyWith(
                          color: AppTheme.sanggamGold,
                          fontWeight: FontWeight.bold,
                          letterSpacing: 4,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '초정밀 차동 계측 시스템',
                        style: theme.textTheme.bodyLarge?.copyWith(
                          color: Colors.white.withOpacity(0.7),
                        ),
                      ),
                      const SizedBox(height: 64),
                      // 로딩 인디케이터 (Lottie)
                      SizedBox(
                        width: 64,
                        height: 64,
                        child: Lottie.asset(
                          'assets/lottie/logo_intro.json',
                          repeat: true,
                          errorBuilder: (_, __, ___) => CircularProgressIndicator(
                            strokeWidth: 2.5,
                            color: AppTheme.sanggamGold.withOpacity(0.6),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

