import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

// ───────────────────────────────────────────────────
// HeroTelemedicineCard — Sanggam Orbit 화상 진료 히어로 카드
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) 2x → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF003344) → SanggamTheme.surface
// [Rule 4] Color(0xFF00838F)/Color(0xFF00ACC1) → SanggamTheme.jagaeCyan 기반
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] 미사용 import animate_fade_in_up 제거
// ───────────────────────────────────────────────────

class HeroTelemedicineCard extends StatelessWidget {
  const HeroTelemedicineCard({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      height: 200,
      decoration: BoxDecoration(
        color: SanggamTheme.surface,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: SanggamTheme.primary.withValues(alpha: 0.5)),
        boxShadow: [
          BoxShadow(
            color: SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
            blurRadius: 20,
            spreadRadius: 0,
          ),
        ],
      ),
      child: Stack(
        children: [
          // Background Gradient Overlay
          Positioned.fill(
            child: ClipRRect(
              borderRadius: BorderRadius.circular(24),
              child: DecoratedBox(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [
                      Colors.transparent,
                      SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
                    ],
                  ),
                ),
              ),
            ),
          ),

          Padding(
            padding: const EdgeInsets.all(24.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Header (Live Status)
                Row(
                  children: [
                    BreathingGlow(
                      child: Container(
                        width: 12,
                        height: 12,
                        decoration: const BoxDecoration(
                          color: SanggamTheme.jagaeCyan,
                          shape: BoxShape.circle,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    const Text(
                      '진료 가능 (대기 0명)',
                      style: TextStyle(
                        color: SanggamTheme.jagaeCyan,
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
                const Spacer(),

                // Title
                const Text(
                  '비대면 화상 진료',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    letterSpacing: -0.5,
                  ),
                ),
                const Text(
                  '증상을 말씀하시면 AI가 적합한 의사를 연결합니다.',
                  style: TextStyle(
                    color: Colors.white70,
                    fontSize: 12,
                  ),
                ),
                const SizedBox(height: 16),

                // Action Button
                SizedBox(
                  width: double.infinity,
                  child: ScaleButton(
                    onPressed: () => context.push('/medical/video-call/session123'),
                    child: Container(
                      padding: const EdgeInsets.symmetric(vertical: 14),
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          colors: [
                            SanggamTheme.jagaeCyan.withValues(alpha: 0.7),
                            SanggamTheme.jagaeCyan,
                          ],
                        ),
                        borderRadius: BorderRadius.circular(12),
                        boxShadow: [
                          BoxShadow(
                            color: SanggamTheme.jagaeCyan.withValues(alpha: 0.4),
                            blurRadius: 10,
                            offset: const Offset(0, 4),
                          ),
                        ],
                      ),
                      child: const Center(
                        child: Text(
                          '지금 바로 진료 시작하기',
                          style: TextStyle(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),

          // Doctor Image Layout (Top Right)
          Positioned(
            top: 24,
            right: 24,
            child: CircleAvatar(
              radius: 28,
              backgroundColor: Colors.white.withValues(alpha: 0.1),
              child: const Icon(Icons.medical_services_outlined, color: Colors.white, size: 28),
            ),
          ),
        ],
      ),
    );
  }
}
