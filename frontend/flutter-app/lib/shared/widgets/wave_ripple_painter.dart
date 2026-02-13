import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 물결 파동(Wave Ripple) 애니메이션을 그리는 CustomPainter.
///
/// 만파식적(萬波息笛) 설화에서 차용:
/// 거친 파도(노이즈)를 잠재우고 평온(데이터)을 찾는 과정을 시각화.
/// 동심원이 중심에서 바깥으로 퍼져나가며, Wave Cyan → Sanggam Gold
/// 그라데이션으로 전환됩니다.
class WaveRipplePainter extends CustomPainter {
  final double animationValue;
  final int rippleCount;
  final Color primaryColor;
  final Color secondaryColor;

  WaveRipplePainter({
    required this.animationValue,
    this.rippleCount = 4,
    this.primaryColor = AppTheme.waveCyan,
    this.secondaryColor = AppTheme.sanggamGold,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final maxRadius = math.max(size.width, size.height) * 0.6;

    for (int i = 0; i < rippleCount; i++) {
      final phase = (animationValue + i / rippleCount) % 1.0;
      final radius = maxRadius * phase;
      final opacity = (1.0 - phase).clamp(0.0, 1.0) * 0.35;

      // Wave Cyan → Sanggam Gold로 페이드
      final color = Color.lerp(primaryColor, secondaryColor, phase)!
          .withOpacity(opacity);

      final paint = Paint()
        ..color = color
        ..style = PaintingStyle.stroke
        ..strokeWidth = 2.0 * (1.0 - phase * 0.5);

      canvas.drawCircle(center, radius, paint);
    }

    // 중심부 글로우
    final glowPaint = Paint()
      ..shader = RadialGradient(
        colors: [
          primaryColor.withOpacity(0.2 * (0.5 + 0.5 * math.sin(animationValue * math.pi * 2))),
          primaryColor.withOpacity(0.0),
        ],
      ).createShader(Rect.fromCircle(center: center, radius: maxRadius * 0.15));
    canvas.drawCircle(center, maxRadius * 0.15, glowPaint);
  }

  @override
  bool shouldRepaint(covariant WaveRipplePainter oldDelegate) {
    return oldDelegate.animationValue != animationValue;
  }
}

/// [WaveRippleBackground] — 물결 파동 배경 위젯.
///
/// Splash 화면이나 측정 화면 배경에 배치하여
/// 은은한 파동 효과를 연출합니다.
class WaveRippleBackground extends StatefulWidget {
  final Widget? child;
  final Duration duration;
  final int rippleCount;
  final Color primaryColor;
  final Color secondaryColor;

  const WaveRippleBackground({
    super.key,
    this.child,
    this.duration = const Duration(seconds: 4),
    this.rippleCount = 4,
    this.primaryColor = AppTheme.waveCyan,
    this.secondaryColor = AppTheme.sanggamGold,
  });

  @override
  State<WaveRippleBackground> createState() => _WaveRippleBackgroundState();
}

class _WaveRippleBackgroundState extends State<WaveRippleBackground>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: widget.duration,
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        return CustomPaint(
          painter: WaveRipplePainter(
            animationValue: _controller.value,
            rippleCount: widget.rippleCount,
            primaryColor: widget.primaryColor,
            secondaryColor: widget.secondaryColor,
          ),
          child: child,
        );
      },
      child: widget.child,
    );
  }
}

/// [WavePainter] — 수평 사인파 물결을 그리는 CustomPainter.
///
/// 측정 진행 중 "파동 안정화" 프로세스를 시각화합니다.
/// 진행률에 따라 진폭(amplitude)이 줄어들어 직선으로 수렴.
class WavePainter extends CustomPainter {
  final double animationValue;
  final double progress; // 0.0 ~ 1.0, 진행률
  final Color waveColor;

  WavePainter({
    required this.animationValue,
    this.progress = 0.0,
    this.waveColor = AppTheme.waveCyan,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.5
      ..strokeCap = StrokeCap.round;

    final path = Path();
    final midY = size.height / 2;
    // 진행률에 따라 진폭 감소: 파동 → 직선
    final amplitude = size.height * 0.25 * (1.0 - progress);
    final frequency = 2 * math.pi * 3 / size.width;

    path.moveTo(0, midY);
    for (double x = 0; x <= size.width; x += 1) {
      final y = midY +
          amplitude *
              math.sin(frequency * x + animationValue * math.pi * 2) *
              math.cos(frequency * x * 0.5 + animationValue * math.pi);
      path.lineTo(x, y);
    }

    // 글로우 효과 (아래 레이어)
    final glowPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 6.0
      ..strokeCap = StrokeCap.round
      ..color = waveColor.withOpacity(0.15)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8);
    canvas.drawPath(path, glowPaint);

    // 메인 라인 (위 레이어)
    paint.shader = LinearGradient(
      colors: [
        waveColor.withOpacity(0.3),
        waveColor,
        AppTheme.sanggamGold,
        waveColor,
        waveColor.withOpacity(0.3),
      ],
    ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant WavePainter oldDelegate) {
    return oldDelegate.animationValue != animationValue ||
        oldDelegate.progress != progress;
  }
}
