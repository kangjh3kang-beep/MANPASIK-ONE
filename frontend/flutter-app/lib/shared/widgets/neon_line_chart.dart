import 'dart:ui';
import 'package:flutter/material.dart';

class NeonLineChart extends StatelessWidget {
  final List<double> dataPoints;
  final Color color;
  final double strokeWidth;
  final bool showPoints;
  final double height;
  final Duration animationDuration;

  const NeonLineChart({
    super.key,
    required this.dataPoints,
    this.color = const Color(0xFF00E5FF),
    this.strokeWidth = 2.0,
    this.showPoints = true,
    this.height = 100,
    this.animationDuration = const Duration(milliseconds: 1500),
  });

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0.0, end: 1.0),
      duration: animationDuration,
      curve: Curves.easeOutCubic,
      builder: (context, value, child) {
        return CustomPaint(
          size: Size(double.infinity, height),
          painter: _NeonChartPainter(
            dataPoints: dataPoints,
            color: color,
            strokeWidth: strokeWidth,
            showPoints: showPoints,
            progress: value,
          ),
        );
      },
    );
  }
}

class _NeonChartPainter extends CustomPainter {
  final List<double> dataPoints;
  final Color color;
  final double strokeWidth;
  final bool showPoints;
  final double progress;

  _NeonChartPainter({
    required this.dataPoints,
    required this.color,
    required this.strokeWidth,
    required this.showPoints,
    required this.progress,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (dataPoints.isEmpty) return;

    final path = Path();
    final stepX = size.width / (dataPoints.length - 1);
    
    // Normalize data to 0..1 range (assuming min 0, max is max value in list)
    final double maxY = dataPoints.reduce((curr, next) => curr > next ? curr : next);
    // Add some padding to top
    final double rangeY = maxY * 1.2;

    // Build Path
    for (int i = 0; i < dataPoints.length; i++) {
      final x = i * stepX;
      // Animate Y: interpolate from baseline (size.height) to actual Y
      final targetY = size.height - (dataPoints[i] / rangeY) * size.height;
      final y = lerpDouble(size.height, targetY, progress)!;

      if (i == 0) {
        path.moveTo(x, y);
      } else {
        // Cubic Bezier for smooth curves
        final prevX = (i - 1) * stepX;
        final prevTargetY = size.height - (dataPoints[i - 1] / rangeY) * size.height;
        final prevY = lerpDouble(size.height, prevTargetY, progress)!;
        
        final controlX1 = prevX + stepX / 2;
        final controlY1 = prevY;
        final controlX2 = prevX + stepX / 2;
        final controlY2 = y;
        
        path.cubicTo(controlX1, controlY1, controlX2, controlY2, x, y);
      }
    }

    // 1. Draw Glow (Blur)
    final glowPaint = Paint()
      ..color = color.withOpacity(0.5 * progress)
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth * 2
      ..strokeCap = StrokeCap.round
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8);
    
    canvas.drawPath(path, glowPaint);

    // 2. Draw Main Line
    final linePaint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = strokeWidth
      ..strokeCap = StrokeCap.round;

    canvas.drawPath(path, linePaint);

    // 3. Gradient Fill (Area under curve)
    final fillPath = Path.from(path)
      ..lineTo(size.width, size.height)
      ..lineTo(0, size.height)
      ..close();

    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          color.withOpacity(0.3 * progress),
          color.withOpacity(0.0),
        ],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));
      
    canvas.drawPath(fillPath, fillPaint);

    // 4. Draw Points (Breath effect)
    if (showPoints) {
      final pointPaint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.fill;
      
      final pointGlowPaint = Paint()
      ..color = color.withOpacity(0.6)
      ..style = PaintingStyle.fill
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 4);

      for (int i = 0; i < dataPoints.length; i++) {
        final x = i * stepX;
        final targetY = size.height - (dataPoints[i] / rangeY) * size.height;
        final y = lerpDouble(size.height, targetY, progress)!;

        // Draw point
        canvas.drawCircle(Offset(x, y), 4 * progress, pointGlowPaint);
        canvas.drawCircle(Offset(x, y), 2 * progress, pointPaint);
      }
    }
  }

  @override
  bool shouldRepaint(covariant _NeonChartPainter oldDelegate) {
    return oldDelegate.progress != progress;
  }
}
