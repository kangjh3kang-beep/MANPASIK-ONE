import 'dart:math';
import 'package:flutter/material.dart';

/// 간단한 Confetti 오버레이 위젯.
///
/// [ConfettiController.fire]를 호출하면 파티클이 분사됩니다.
class ConfettiOverlay extends StatefulWidget {
  const ConfettiOverlay({super.key, required this.controller, required this.child});

  final ConfettiController controller;
  final Widget child;

  @override
  State<ConfettiOverlay> createState() => _ConfettiOverlayState();
}

class _ConfettiOverlayState extends State<ConfettiOverlay>
    with SingleTickerProviderStateMixin {
  late final AnimationController _anim;
  List<_Particle> _particles = [];
  final _rand = Random();

  @override
  void initState() {
    super.initState();
    _anim = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 2500),
    )..addListener(() => setState(() {}));

    widget.controller._attach(this);
  }

  @override
  void dispose() {
    _anim.dispose();
    super.dispose();
  }

  void _fire() {
    _particles = List.generate(40, (_) => _Particle(_rand));
    _anim.forward(from: 0);
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        widget.child,
        if (_anim.isAnimating)
          Positioned.fill(
            child: IgnorePointer(
              child: CustomPaint(
                painter: _ConfettiPainter(
                  particles: _particles,
                  progress: _anim.value,
                ),
              ),
            ),
          ),
      ],
    );
  }
}

class ConfettiController {
  _ConfettiOverlayState? _state;

  void _attach(_ConfettiOverlayState state) => _state = state;

  void fire() => _state?._fire();
}

class _Particle {
  final double x;
  final double angle;
  final double speed;
  final double size;
  final Color color;
  final double rotSpeed;

  static const _colors = [
    Color(0xFFD4AF37), // Gold
    Color(0xFF00E5FF), // Cyan
    Color(0xFFFF1744), // Red
    Color(0xFF2962FF), // Blue
    Color(0xFF00C853), // Green
    Color(0xFFFF6D00), // Orange
    Color(0xFFAA00FF), // Purple
  ];

  _Particle(Random rand)
      : x = rand.nextDouble(),
        angle = -pi / 2 + (rand.nextDouble() - 0.5) * pi * 0.8,
        speed = 300 + rand.nextDouble() * 400,
        size = 4 + rand.nextDouble() * 6,
        color = _colors[rand.nextInt(_colors.length)],
        rotSpeed = rand.nextDouble() * 6 - 3;
}

class _ConfettiPainter extends CustomPainter {
  final List<_Particle> particles;
  final double progress;

  _ConfettiPainter({required this.particles, required this.progress});

  @override
  void paint(Canvas canvas, Size size) {
    final gravity = size.height * 1.2;
    final t = progress;
    final opacity = (1 - t).clamp(0.0, 1.0);

    for (final p in particles) {
      final px = p.x * size.width + cos(p.angle) * p.speed * t;
      final py = size.height * 0.5 + sin(p.angle) * p.speed * t + 0.5 * gravity * t * t;

      if (py > size.height || py < -20) continue;

      canvas.save();
      canvas.translate(px, py);
      canvas.rotate(p.rotSpeed * t * pi);

      final paint = Paint()
        ..color = p.color.withOpacity(opacity)
        ..style = PaintingStyle.fill;

      canvas.drawRect(
        Rect.fromCenter(center: Offset.zero, width: p.size, height: p.size * 0.6),
        paint,
      );
      canvas.restore();
    }
  }

  @override
  bool shouldRepaint(_ConfettiPainter oldDelegate) => oldDelegate.progress != progress;
}
