import 'dart:math' as math;
import 'package:flutter/material.dart';

class HoloGlobe extends StatefulWidget {
  final double size;
  final Color color;
  final Color? accentColor;

  const HoloGlobe({
    super.key,
    this.size = 300,
    this.color = const Color(0xFF00E5FF),
    this.accentColor,
  });

  @override
  State<HoloGlobe> createState() => _HoloGlobeState();
}

class _HoloGlobeState extends State<HoloGlobe> with TickerProviderStateMixin {
  late AnimationController _rotationController;
  late AnimationController _pulseController;
  late AnimationController _scanController;

  final List<_Point3D> _points = [];
  final int _pointCount = 1500;

  @override
  void initState() {
    super.initState();
    _rotationController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 24),
    )..repeat();

    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 4),
    )..repeat();

    _scanController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 5),
    )..repeat(reverse: true);

    _generatePoints();
  }

  void _generatePoints() {
    for (int i = 0; i < _pointCount; i++) {
      final phi = math.acos(1 - 2 * (i + 0.5) / _pointCount);
      final theta = math.pi * (1 + math.sqrt(5)) * i;

      final r = widget.size * 0.4;
      final x = r * math.sin(phi) * math.cos(theta);
      final y = r * math.sin(phi) * math.sin(theta);
      final z = r * math.cos(phi);

      _points.add(_Point3D(x, y, z));
    }
  }

  @override
  void dispose() {
    _rotationController.dispose();
    _pulseController.dispose();
    _scanController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ExcludeSemantics(
      child: AnimatedBuilder(
        animation: Listenable.merge([_rotationController, _pulseController, _scanController]),
        builder: (context, child) {
          return CustomPaint(
            size: Size(widget.size, widget.size),
            painter: _GlobePainter(
              points: _points,
              rotation: _rotationController.value * 2 * math.pi,
              pulseValue: _pulseController.value,
              scanValue: _scanController.value,
              color: widget.color,
              accentColor: widget.accentColor ?? widget.color,
            ),
          );
        },
      ),
    );
  }
}

class _Point3D {
  double x, y, z;
  _Point3D(this.x, this.y, this.z);
}

class _GlobePainter extends CustomPainter {
  final List<_Point3D> points;
  final double rotation;
  final double pulseValue;
  final double scanValue;
  final Color color;
  final Color accentColor;

  _GlobePainter({
    required this.points,
    required this.rotation,
    required this.pulseValue,
    required this.scanValue,
    required this.color,
    required this.accentColor,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width * 0.4;

    // 1. Pearl Core (Yeouiju)
    final coreGradient = RadialGradient(
      colors: [
        Colors.white.withValues(alpha: 0.6),
        accentColor.withValues(alpha: 0.4),
        color.withValues(alpha: 0.05),
        Colors.transparent,
      ],
      stops: const [0.0, 0.2, 0.5, 1.0],
    ).createShader(Rect.fromCircle(center: center, radius: radius * 0.6));

    canvas.drawCircle(center, radius * 0.6, Paint()..shader = coreGradient);

    // 2. Wireframe Lat/Long
    final wireframePaint = Paint()
      ..color = color.withValues(alpha: 0.05)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5;

    _drawWireframe(canvas, center, radius, rotation, wireframePaint);

    // 3. Particles (Volumetric Cloud)
    final pointPaint = Paint()
      ..color = color
      ..strokeCap = StrokeCap.round
      ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 1.0);

    final glowPaint = Paint()
      ..color = accentColor.withValues(alpha: 0.3)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 3.0);

    for (var point in points) {
      final rotatedX = point.x * math.cos(rotation) - point.z * math.sin(rotation);
      final rotatedZ = point.x * math.sin(rotation) + point.z * math.cos(rotation);
      final y = point.y;

      final scale = 300 / (300 - rotatedZ);
      final x2d = rotatedX * scale + center.dx;
      final y2d = y * scale + center.dy;

      final alpha = ((rotatedZ + radius) / (2 * radius)).clamp(0.0, 1.0);

      final scanY = (scanValue * 2 - 1) * radius;
      final distToScan = (y - scanY).abs();

      final isScanned = distToScan < 5.0;

      if (isScanned) {
         pointPaint.color = Colors.white.withValues(alpha: alpha);
         pointPaint.strokeWidth = 2.5 * scale;
         pointPaint.maskFilter = null;

         canvas.drawCircle(Offset(x2d, y2d), 4.0 * scale, glowPaint);
      } else {
         pointPaint.color = color.withValues(alpha: alpha * 0.7);
         pointPaint.strokeWidth = 1.5 * scale;
         pointPaint.maskFilter = const MaskFilter.blur(BlurStyle.solid, 1.0);
      }

      canvas.drawCircle(Offset(x2d, y2d), (isScanned ? 2.0 : 1.2) * scale, pointPaint);
    }

    // 4. Shockwave
    for(int i=0; i<3; i++) {
        final waveProgress = (pulseValue + i * 0.33) % 1.0;
        final waveRadius = radius * (0.4 + waveProgress * 0.8);

        if (waveRadius < size.width * 0.55) {
             final waveAlpha = (1.0 - waveProgress).clamp(0.0, 1.0);
             final pulsePaint = Paint()
              ..style = PaintingStyle.stroke
              ..strokeWidth = 1.5 * (1.0 - waveProgress)
              ..color = accentColor.withValues(alpha: waveAlpha * 0.4)
              ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 2);

             canvas.drawOval(
               Rect.fromCenter(center: center, width: waveRadius * 2, height: waveRadius * 2),
               pulsePaint
             );
        }
    }

    // 5. Data Scan Laser
    final scanPlanY = center.dy + (scanValue * 2 - 1) * radius;
    final dY = (scanPlanY - center.dy).abs();
    if (dY < radius) {
        final scanWidth = math.sqrt(radius * radius - dY * dY) * 2;

        canvas.drawLine(
          Offset(center.dx - scanWidth/2, scanPlanY),
          Offset(center.dx + scanWidth/2, scanPlanY),
          Paint()
            ..color = Colors.white.withValues(alpha: 0.8)
            ..strokeWidth = 1.0
            ..shader = LinearGradient(colors: [
               Colors.transparent, accentColor, Colors.white, accentColor, Colors.transparent
            ]).createShader(Rect.fromLTWH(center.dx - scanWidth/2, scanPlanY, scanWidth, 2))
        );

        canvas.drawOval(
           Rect.fromCenter(center: Offset(center.dx, scanPlanY), width: scanWidth, height: scanWidth * 0.3),
           Paint()
            ..color = accentColor.withValues(alpha: 0.1)
            ..style = PaintingStyle.fill
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 15)
        );
    }

    // 6. Complex Waves
    _drawMultiWaves(canvas, center, size.width * 0.9, rotation);
  }

  void _drawWireframe(Canvas canvas, Offset center, double radius, double rotation, Paint paint) {
     for(int i=0; i<6; i++) {
        // Meridian drawing placeholder
     }

     canvas.drawCircle(center, radius, paint);
     canvas.drawOval(Rect.fromCenter(center: center, width: radius * 2, height: radius * 0.6), paint);
     canvas.drawOval(Rect.fromCenter(center: center, width: radius * 1.5, height: radius * 2), paint);
  }

  void _drawMultiWaves(Canvas canvas, Offset center, double width, double time) {
    void drawWave(double freq, double amp, double speed, Color c, double widthStroke) {
       final path = Path();
       final startX = center.dx - width / 2;
       path.moveTo(startX, center.dy);

       for(double x = 0; x <= width; x += 5) {
          final nX = x / width;
          final env = math.pow(math.sin(math.pi * nX), 2).toDouble();

          double y = math.sin(x * freq + time * speed) * amp * env;
          y += math.sin(x * freq * 2.5 - time * speed * 1.5) * (amp * 0.3) * env;

          path.lineTo(startX + x, center.dy + y);
       }

       final p = Paint()
         ..color = c
         ..style = PaintingStyle.stroke
         ..strokeWidth = widthStroke;

       canvas.drawPath(path, p);
    }

    drawWave(0.1, 40, 5, color.withValues(alpha: 0.8), 2.0);
    drawWave(0.06, 30, 3, accentColor.withValues(alpha: 0.6), 1.5);
    drawWave(0.03, 60, 2, color.withValues(alpha: 0.3), 1.0);
  }

  @override
  bool shouldRepaint(covariant _GlobePainter oldDelegate) => true;
}
