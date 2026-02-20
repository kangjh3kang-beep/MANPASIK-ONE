import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

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
  final int _pointCount = 1500; // Increased density for volumetric feel

  @override
  void initState() {
    super.initState();
    // 1. Rotation
    _rotationController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 24),
    )..repeat();

    // 2. Energy Pulse (Expansion)
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 4),
    )..repeat();

    // 3. Data Scan (Vertical Scan)
    _scanController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 5),
    )..repeat(reverse: true);

    _generatePoints();
  }

  void _generatePoints() {
    final random = math.Random();
    for (int i = 0; i < _pointCount; i++) {
      // Golden Spiral distribution for uniform sphere coverage
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

    // 1. Pearl Core (Yeouiju) - Solid glowing center
    final coreGradient = RadialGradient(
      colors: [
        Colors.white.withOpacity(0.6),      // Bright core
        accentColor.withOpacity(0.4),       // Inner glow
        color.withOpacity(0.05),            // Outer aura
        Colors.transparent,
      ],
      stops: const [0.0, 0.2, 0.5, 1.0],
    ).createShader(Rect.fromCircle(center: center, radius: radius * 0.6));
    
    // Draw Core
    canvas.drawCircle(center, radius * 0.6, Paint()..shader = coreGradient);

    // 2. Wireframe Lat/Long (Subtle Background Structure)
    final wireframePaint = Paint()
      ..color = color.withOpacity(0.05)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5;

    _drawWireframe(canvas, center, radius, rotation, wireframePaint);

    // 3. Particles (Volumetric Cloud)
    final pointPaint = Paint()
      ..color = color
      ..strokeCap = StrokeCap.round
      ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 1.0); // Soften edges
    
    final glowPaint = Paint()
      ..color = accentColor.withOpacity(0.3)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 3.0); // Outer glow

    for (var point in points) {
      // Rotate Point
      double rotatedX = point.x * math.cos(rotation) - point.z * math.sin(rotation);
      double rotatedZ = point.x * math.sin(rotation) + point.z * math.cos(rotation);
      double y = point.y;

      // Perspective Projection
      double scale = 300 / (300 - rotatedZ);
      double x2d = rotatedX * scale + center.dx;
      double y2d = y * scale + center.dy;

      // Depth alpha
      double alpha = ((rotatedZ + radius) / (2 * radius)).clamp(0.0, 1.0);
      
      // Interaction with Scan Line
      // Map scanValue (0..1) to (-radius..+radius)
      double scanY = (scanValue * 2 - 1) * radius;
      double distToScan = (y - scanY).abs();
      
      bool isScanned = distToScan < 5.0; // Particles near scan line

      if (isScanned) {
         pointPaint.color = Colors.white.withOpacity(alpha); // White hot scan
         pointPaint.strokeWidth = 2.5 * scale;
         pointPaint.maskFilter = null; // Sharp for scanned points
         
         // Extra glow for scanned points
         canvas.drawCircle(Offset(x2d, y2d), 4.0 * scale, glowPaint);
      } else {
         pointPaint.color = color.withOpacity(alpha * 0.7); // Slightly more opaque
         pointPaint.strokeWidth = 1.5 * scale; // Slightly larger
         pointPaint.maskFilter = const MaskFilter.blur(BlurStyle.solid, 1.0);
      }
      
      canvas.drawCircle(Offset(x2d, y2d), (isScanned ? 2.0 : 1.2) * scale, pointPaint);
    }

    // 4. Shockwave (Multiple Expanding Rings)
    for(int i=0; i<3; i++) {
        double waveProgress = (pulseValue + i * 0.33) % 1.0;
        double waveRadius = radius * (0.4 + waveProgress * 0.8);
        
        if (waveRadius < size.width * 0.55) {
             final waveAlpha = (1.0 - waveProgress).clamp(0.0, 1.0);
             final pulsePaint = Paint()
              ..style = PaintingStyle.stroke
              ..strokeWidth = 1.5 * (1.0 - waveProgress)
              ..color = accentColor.withOpacity(waveAlpha * 0.4)
              ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 2);
              
             // Draw elliptical ring to match perspective
             canvas.drawOval(
               Rect.fromCenter(center: center, width: waveRadius * 2, height: waveRadius * 2), 
               pulsePaint
             );
        }
    }

    // 5. Data Scan Laser (Horizontal Plane)
    double scanPlanY = center.dy + (scanValue * 2 - 1) * radius;
    // Calculate width at this Y
    double dY = (scanPlanY - center.dy).abs();
    if (dY < radius) {
        double scanWidth = math.sqrt(radius * radius - dY * dY) * 2;
        
        // Laser Line
        canvas.drawLine(
          Offset(center.dx - scanWidth/2, scanPlanY),
          Offset(center.dx + scanWidth/2, scanPlanY),
          Paint()
            ..color = Colors.white.withOpacity(0.8)
            ..strokeWidth = 1.0
            ..shader = LinearGradient(colors: [
               Colors.transparent, accentColor, Colors.white, accentColor, Colors.transparent
            ]).createShader(Rect.fromLTWH(center.dx - scanWidth/2, scanPlanY, scanWidth, 2))
        );
        
        // Laser Glow
        canvas.drawOval(
           Rect.fromCenter(center: Offset(center.dx, scanPlanY), width: scanWidth, height: scanWidth * 0.3),
           Paint()
            ..color = accentColor.withOpacity(0.1)
            ..style = PaintingStyle.fill
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 15)
        );
    }

    // 6. Complex Waves
    _drawMultiWaves(canvas, center, size.width * 0.9, rotation);
  }

  void _drawWireframe(Canvas canvas, Offset center, double radius, double rotation, Paint paint) {
     // Longitudinal Lines (Meridians)
     for(int i=0; i<6; i++) {
        double angle = (i * 30) * math.pi / 180;
        // Draw ellipse for meridian
        // This is complex in 2D. Simplified: just a few static rings rotated? 
        // Better: Dynamic calculation based on Z rotation.
        // For efficiency, skipping full 3D wireframe mesh for now, using particles for volume.
     }
     
     // Simple Equatorial Rings
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
          double nX = x / width; // 0..1
          // Envelope: Sinc-like or Bell curve to keep edges attached
          double env = math.pow(math.sin(math.pi * nX), 2).toDouble();
          
          double y = math.sin(x * freq + time * speed) * amp * env;
          // Add interference
          y += math.sin(x * freq * 2.5 - time * speed * 1.5) * (amp * 0.3) * env;

          path.lineTo(startX + x, center.dy + y);
       }
       
       final p = Paint()
         ..color = c
         ..style = PaintingStyle.stroke
         ..strokeWidth = widthStroke;
       
       canvas.drawPath(path, p);
    }

    // 1. Primary High-Frequency Data Wave
    drawWave(0.1, 40, 5, color.withOpacity(0.8), 2.0);
    
    // 2. Secondary Harmonic Wave (Accent)
    drawWave(0.06, 30, 3, accentColor.withOpacity(0.6), 1.5);
    
    // 3. Bass Wave (Low Freq)
    drawWave(0.03, 60, 2, color.withOpacity(0.3), 1.0);
  }

  @override
  bool shouldRepaint(covariant _GlobePainter oldDelegate) => true;
}
