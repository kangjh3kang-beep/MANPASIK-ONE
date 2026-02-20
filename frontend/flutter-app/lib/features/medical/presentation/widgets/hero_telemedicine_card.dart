import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

class HeroTelemedicineCard extends StatelessWidget {
  const HeroTelemedicineCard({super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      height: 200,
      decoration: BoxDecoration(
        color: const Color(0xFF003344), // Deep Teal
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.5)),
        boxShadow: [
          BoxShadow(
            color: AppTheme.waveCyan.withOpacity(0.1),
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
                      AppTheme.waveCyan.withOpacity(0.1),
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
                          color: Color(0xFF00E676), // Bright Green
                          shape: BoxShape.circle,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    const Text(
                      '진료 가능 (대기 0명)',
                      style: TextStyle(
                        color: Color(0xFF00E676),
                        fontWeight: FontWeight.bold,
                        fontSize: 14,
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                
                // Icon Illustration (Placeholder for now, drawing simple shapes)
                // ...
                
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
                        gradient: const LinearGradient(
                          colors: [Color(0xFF00838F), Color(0xFF00ACC1)],
                        ),
                        borderRadius: BorderRadius.circular(12),
                        boxShadow: [
                          BoxShadow(
                            color: const Color(0xFF00ACC1).withOpacity(0.4),
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
              backgroundColor: Colors.white.withOpacity(0.1),
              child: const Icon(Icons.medical_services_outlined, color: Colors.white, size: 28),
            ),
          )
        ],
      ),
    );
  }
}
