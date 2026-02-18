import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 한지 (Korean Paper) 배경
///
/// 백자(Baekja)와 창호지의 미학을 담은 은은한 배경 위젯.
/// - 닥나무 섬유질(Fiber) 텍스처
/// - 은은하게 일렁이는 햇살(Sunlight) 애니메이션
/// - 따뜻한 미색(Off-White) 베이스
class HanjiBackground extends StatefulWidget {
  final Widget child;

  const HanjiBackground({super.key, required this.child});

  @override
  State<HanjiBackground> createState() => _HanjiBackgroundState();
}

class _HanjiBackgroundState extends State<HanjiBackground> with SingleTickerProviderStateMixin {
  late AnimationController _sunlightController;

  @override
  void initState() {
    super.initState();
    _sunlightController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 10),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _sunlightController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // 1. 따뜻한 백자색 베이스
        Container(
          color: const Color(0xFFF9F9F7), // Baekja White
        ),

        // 2. 한지 섬유질 텍스처 (Custom Painter)
        Positioned.fill(
          child: CustomPaint(
            painter: _HanjiFiberPainter(),
          ),
        ),

        // 3. 은은한 햇살 (Breathing Gradient) - Toned down for readability
        Positioned.fill(
          child: AnimatedBuilder(
            animation: _sunlightController,
            builder: (context, _) {
              return Container(
                decoration: BoxDecoration(
                  gradient: RadialGradient(
                    center: const Alignment(0.0, -0.3), // 약간 위쪽
                    radius: 1.5 + (_sunlightController.value * 0.1), // Reduced movement
                    colors: [
                      Colors.white.withOpacity(0.2), // Reduced opacity from 0.4
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

class _HanjiFiberPainter extends CustomPainter {
  final _random = math.Random(42); // 고정된 시드

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5
      ..color = const Color(0xFFDCDCDC).withOpacity(0.15); // 아주 연한 회색 (가독성 위해 투명도 조절)

    // 무작위 섬유질 그리기
    for (int i = 0; i < 500; i++) {
      final x = _random.nextDouble() * size.width;
      final y = _random.nextDouble() * size.height;
      final length = _random.nextDouble() * 10 + 5;
      final angle = _random.nextDouble() * 2 * math.pi;

      final dx = math.cos(angle) * length;
      final dy = math.sin(angle) * length;

      // 닥나무 섬유의 불규칙한 느낌 (곡선)
      final path = Path();
      path.moveTo(x, y);
      path.quadraticBezierTo(
        x + dx * 0.5 + _random.nextDouble() * 2, 
        y + dy * 0.5 + _random.nextDouble() * 2, 
        x + dx, 
        y + dy
      );
      
      canvas.drawPath(path, paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
