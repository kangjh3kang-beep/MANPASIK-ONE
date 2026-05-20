import 'package:flutter/material.dart';
import 'dart:math' as math;
import 'package:manpasik/core/theme/sanggam_theme.dart';

/// 신비로운 심해 배경 (Mystic Deep Sea Background)
///
/// - Base: Deep Sea Blue -> Celadon Teal Gradient
/// - Effect 1: Rainbow Mist (Aurora-like moving fog)
/// - Effect 2: Bioluminescent Particles (Rising glowing dots)
/// - Effect 3: Caustics (Underwater light refraction)
class RoyalCloudBackground extends StatefulWidget {
  final Widget child;
  const RoyalCloudBackground({super.key, required this.child});

  @override
  State<RoyalCloudBackground> createState() => _RoyalCloudBackgroundState();
}

class _RoyalCloudBackgroundState extends State<RoyalCloudBackground> with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
       vsync: this,
       duration: const Duration(seconds: 15),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Background Paint
        AnimatedBuilder(
          animation: _controller,
          builder: (context, _) {
            return CustomPaint(
              painter: _MysticDeepSeaPainter(animationValue: _controller.value),
              size: Size.infinite,
            );
          },
        ),

        // Content
        widget.child,
      ],
    );
  }
}

class _MysticDeepSeaPainter extends CustomPainter {
  final double animationValue;

  _MysticDeepSeaPainter({required this.animationValue});

  // AppTheme에서 사용하던 색상 직접 참조
  static const _celadonTeal = Color(0xFF00897B);

  @override
  void paint(Canvas canvas, Size size) {
    final rect = Rect.fromLTWH(0, 0, size.width, size.height);
    final gradient = LinearGradient(
      begin: Alignment.topCenter,
      end: Alignment.bottomCenter,
      colors: [
        const Color(0xFF001020),
        SanggamTheme.background,
        _celadonTeal.withValues(alpha: 0.4),
      ],
      stops: const [0.0, 0.6, 1.0],
    );
    canvas.drawRect(rect, Paint()..shader = gradient.createShader(rect));

    _drawRainbowMist(canvas, size, animationValue);
    _drawBioluminescence(canvas, size, animationValue);
    _drawCaustics(canvas, size, animationValue);
  }

  void _drawRainbowMist(Canvas canvas, Size size, double anim) {
    final colors = [
      Colors.purpleAccent.withValues(alpha: 0.15),
      Colors.tealAccent.withValues(alpha: 0.15),
      Colors.blueAccent.withValues(alpha: 0.15),
      Colors.amberAccent.withValues(alpha: 0.1),
    ];

    for (int i = 0; i < 3; i++) {
      final path = Path();
      final yBase = size.height * 0.75 + (i * 45);
      final amplitude = 30.0 + (i * 15);
      final shift = anim * 2 * math.pi + (i * math.pi / 2);

      path.moveTo(0, size.height);
      path.lineTo(0, yBase);

      for (double x = 0; x <= size.width; x += 10) {
        final sine1 = math.sin((x / size.width * 2 * math.pi) + shift);
        final sine2 = math.sin((x / size.width * 4 * math.pi) - shift * 0.5);
        final y = yBase + (sine1 + sine2 * 0.5) * amplitude;
        path.lineTo(x, y);
      }

      path.lineTo(size.width, size.height);
      path.close();

      final shader = LinearGradient(
        begin: Alignment.bottomCenter,
        end: Alignment.topCenter,
        colors: [
          colors[i % colors.length].withValues(alpha: 0.5),
          colors[(i+1) % colors.length].withValues(alpha: 0)
        ],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));

      canvas.drawPath(path, Paint()
        ..shader = shader
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 10)
      );
    }
  }

  void _drawBioluminescence(Canvas canvas, Size size, double anim) {
    final random = math.Random(1234);
    const count = 60;

    for (int i = 0; i < count; i++) {
        final speed = 0.2 + random.nextDouble() * 0.3;
        final x = (random.nextDouble() + math.sin(anim * 2 * math.pi * speed + i) * 0.1) * size.width;
        double yPos = (random.nextDouble() - anim * speed);
        yPos = yPos - yPos.floor();

        final y = size.height - (yPos * size.height);

        final opacity = (math.sin(anim * 4 * math.pi + i) + 1) / 2 * 0.5 + 0.2;
        final radius = random.nextDouble() * 2.0 + 0.5;

        canvas.drawCircle(Offset(x, y), radius, Paint()
          ..color = SanggamTheme.jagaeCyan.withValues(alpha: opacity * 0.8)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 2)
        );
    }
  }

  void _drawCaustics(Canvas canvas, Size size, double anim) {
    final paint = Paint()
      ..color = SanggamTheme.jagaeCyan.withValues(alpha: 0.05)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 30.0
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 40);

    for(int i=0; i<2; i++) {
       final path = Path();
       final y = size.height * 0.2 + (i * 200);
       final shift = math.sin(anim * math.pi + i) * 30;

       path.moveTo(0, y + shift);
       path.quadraticBezierTo(
         size.width / 2, y - shift - 50,
         size.width, y + shift
       );
       canvas.drawPath(path, paint);
    }
  }

  @override
  bool shouldRepaint(covariant _MysticDeepSeaPainter oldDelegate) => true;
}
