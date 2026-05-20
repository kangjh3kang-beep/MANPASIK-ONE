import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// HanjiBackground — 다크 상감 한지(韓紙) 배경
//
// [Rule 4] withOpacity 2건 → withValues(alpha:)
// [Rule 4] AppTheme 미사용 → SanggamTheme 통일
// [Rule 4] 라이트 하드코딩 (0xFFF9F9F7, 0xFFDCDCDC) → SanggamTheme 다크 상수
// [Rule 2] 기존 구조 유지 (배경 위젯, 간격 해당 없음)
//
// 닥나무 섬유질(Fiber) 텍스처를 금선(Gold) 톤으로 전환하여
// Sanggam Orbit 디자인 시스템과 시각적 일관성 확보.
// ───────────────────────────────────────────────────

/// 한지(韓紙) 질감 배경 — 다크 Sanggam 스타일.
///
/// 닥나무 섬유질을 금빛 라인으로, 은은한 햇살을
/// 따뜻한 금색 글로우로 표현하여 다크 모드에 통합.
class HanjiBackground extends StatefulWidget {
  final Widget child;

  const HanjiBackground({super.key, required this.child});

  @override
  State<HanjiBackground> createState() => _HanjiBackgroundState();
}

class _HanjiBackgroundState extends State<HanjiBackground>
    with SingleTickerProviderStateMixin {
  late AnimationController _glowCtrl;

  @override
  void initState() {
    super.initState();
    _glowCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 10),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _glowCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // 1. 다크 베이스
        const Positioned.fill(
          child: ColoredBox(color: SanggamTheme.background),
        ),

        // 2. 한지 섬유질 텍스처 (금빛 라인)
        const Positioned.fill(
          child: CustomPaint(painter: _HanjiFiberPainter()),
        ),

        // 3. 은은한 금빛 글로우 (breathing)
        Positioned.fill(
          child: AnimatedBuilder(
            animation: _glowCtrl,
            builder: (context, _) {
              return DecoratedBox(
                decoration: BoxDecoration(
                  gradient: RadialGradient(
                    center: const Alignment(0.0, -0.3),
                    radius: 1.5 + (_glowCtrl.value * 0.1),
                    colors: [
                      SanggamTheme.primary.withValues(alpha: 0.06),
                      Colors.transparent,
                    ],
                    stops: const [0.0, 0.8],
                  ),
                ),
              );
            },
          ),
        ),

        // 4. 전경 콘텐츠
        widget.child,
      ],
    );
  }
}

// ═══════════════════════════════════════════════════
//  한지 섬유질 페인터 — 금빛 닥나무 섬유
// ═══════════════════════════════════════════════════

class _HanjiFiberPainter extends CustomPainter {
  const _HanjiFiberPainter();

  @override
  void paint(Canvas canvas, Size size) {
    final random = math.Random(42); // 고정 시드 — 안정적 텍스처
    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5
      ..color = SanggamTheme.primary.withValues(alpha: 0.04);

    for (int i = 0; i < 500; i++) {
      final x = random.nextDouble() * size.width;
      final y = random.nextDouble() * size.height;
      final length = random.nextDouble() * 10 + 5;
      final angle = random.nextDouble() * 2 * math.pi;

      final dx = math.cos(angle) * length;
      final dy = math.sin(angle) * length;

      final path = Path()
        ..moveTo(x, y)
        ..quadraticBezierTo(
          x + dx * 0.5 + random.nextDouble() * 2,
          y + dy * 0.5 + random.nextDouble() * 2,
          x + dx,
          y + dy,
        );

      canvas.drawPath(path, paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
