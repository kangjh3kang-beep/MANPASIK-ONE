import 'package:flutter/material.dart';
import 'dart:math' as math;
import 'dart:ui' as ui;
import 'package:manpasik/core/theme/app_theme.dart';

/// 신비로운 심해 배경 (Mystic Deep Sea Background)
/// 
/// User Request: "청록색 그라디언트, 심해의 깊은 신비함, 발광, 무지개 안개"
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
       duration: const Duration(seconds: 15), // Faster, visible movement
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

  @override
  void paint(Canvas canvas, Size size) {
    // 1. Base Gradient: Deep Sea Blue -> Celadon Teal (Bottom-Up)
    // "청록색의 그라디언트 심해의 깊은 신비함"
    final rect = Rect.fromLTWH(0, 0, size.width, size.height);
    final gradient = LinearGradient(
      begin: Alignment.topCenter,
      end: Alignment.bottomCenter,
      colors: [
        const Color(0xFF001020),            // Deep Void (Top)
        AppTheme.deepSeaBlue,               // Midnight Blue (Middle)
        AppTheme.celadonTeal.withOpacity(0.4), // Teal Glow (Bottom)
      ],
      stops: const [0.0, 0.6, 1.0],
    );
    canvas.drawRect(rect, Paint()..shader = gradient.createShader(rect));

    // 2. Rainbow Mist (Aurora Effect)
    // "뒤에 무지개안개처럼 움직이는 모습"
    _drawRainbowMist(canvas, size, animationValue);

    // 3. Bioluminescent Particles
    // "발광" (Glowing particles rising)
    _drawBioluminescence(canvas, size, animationValue);

    // 4. Caustics (Underwater Light Refraction)
    // "글라스느낌" (Refracted light from surface)
    _drawCaustics(canvas, size, animationValue);
  }

  void _drawRainbowMist(Canvas canvas, Size size, double anim) {
    final colors = [
      Colors.purpleAccent.withOpacity(0.15),
      Colors.tealAccent.withOpacity(0.15),
      Colors.blueAccent.withOpacity(0.15),
      Colors.amberAccent.withOpacity(0.1), // Touch of gold
    ];

    for (int i = 0; i < 3; i++) {
      final path = Path();
      final yBase = size.height * 0.4 + (i * 150);
      final amplitude = 60.0 + (i * 30); // Larger waves
      final shift = anim * 2 * math.pi + (i * math.pi / 2);
      
      path.moveTo(0, size.height);
      path.lineTo(0, yBase);

      for (double x = 0; x <= size.width; x += 10) {
        final sine1 = math.sin((x / size.width * 2 * math.pi) + shift);
        final sine2 = math.sin((x / size.width * 4 * math.pi) - shift * 0.5); // Complex wave
        final y = yBase + (sine1 + sine2 * 0.5) * amplitude;
        path.lineTo(x, y);
      }
      
      path.lineTo(size.width, size.height);
      path.close();

      final shader = LinearGradient(
        begin: Alignment.bottomCenter, // Strongest at bottom
        end: Alignment.topCenter,      // Fades out upwards
        colors: [
          colors[i % colors.length].withOpacity(0.5), // Much stronger opacity (0.5)
          colors[(i+1) % colors.length].withOpacity(0)
        ],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));

      canvas.drawPath(path, Paint()
        ..shader = shader
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 10) // Sharper mist
      );
    }
  }

  void _drawBioluminescence(Canvas canvas, Size size, double anim) {
    final random = math.Random(1234); // Seed
    final count = 60;

    for (int i = 0; i < count; i++) {
        final speed = 0.2 + random.nextDouble() * 0.3;
        final x = (random.nextDouble() + math.sin(anim * 2 * math.pi * speed + i) * 0.1) * size.width; // Slight horizontal sway
        // Rise from bottom
        double yPos = (random.nextDouble() - anim * speed); 
        // Wrap around 0..1
        yPos = yPos - yPos.floor(); 
        
        final y = size.height - (yPos * size.height);

        final opacity = (math.sin(anim * 4 * math.pi + i) + 1) / 2 * 0.5 + 0.2; // Pulse opacity
        final radius = random.nextDouble() * 2.0 + 0.5;

        canvas.drawCircle(Offset(x, y), radius, Paint()
          ..color = AppTheme.waveCyan.withOpacity(opacity * 0.8)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 2)
        );
    }
  }

  void _drawCaustics(Canvas canvas, Size size, double anim) {
    final paint = Paint()
      ..color = Colors.cyanAccent.withOpacity(0.05)
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
