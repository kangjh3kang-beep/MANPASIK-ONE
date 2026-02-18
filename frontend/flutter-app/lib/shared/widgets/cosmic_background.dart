import 'package:flutter/material.dart';
import 'dart:math' as math;
import 'dart:ui' as ui;
import 'package:manpasik/core/theme/app_theme.dart';

class CosmicBackground extends StatefulWidget {
  final Widget child;
  const CosmicBackground({super.key, required this.child});

  @override
  State<CosmicBackground> createState() => _CosmicBackgroundState();
}

class _CosmicBackgroundState extends State<CosmicBackground> with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
       vsync: this, 
       duration: const Duration(seconds: 20),
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
        // 1. Deep Space Void (Base)
        Container(
          decoration: const BoxDecoration(
            gradient: RadialGradient(
              center: Alignment.center,
              radius: 1.5,
              colors: [
                Color(0xFF001020), // Deep Blue Center
                Color(0xFF000510), // Darker Edge
                Colors.black,      // Absolute Black Corners
              ],
              stops: [0.2, 0.6, 1.0],
            ),
          ),
        ),

        // 2. Animated Star Field
        AnimatedBuilder(
          animation: _controller,
          builder: (context, _) {
            return CustomPaint(
              painter: _StarFieldPainter(animationValue: _controller.value),
              size: Size.infinite,
            );
          },
        ),

        // 3. Nebula/Aurora Clouds (Blurred Gradients)
        Positioned(
          top: -200, left: -200,
          child: _buildNebulaCloud(Colors.purple.withOpacity(0.15), 600),
        ),
        Positioned(
          bottom: -200, right: -200,
          child: _buildNebulaCloud(AppTheme.waveCyan.withOpacity(0.1), 700),
        ),
        
        // 4. Content Overlay
        widget.child,
        
        // 5. Vignette (Cinematic Edge Darkening)
        IgnorePointer(
          child: Container(
            decoration: BoxDecoration(
              gradient: RadialGradient(
                center: Alignment.center,
                radius: 1.2,
                colors: [
                  Colors.transparent, 
                  Colors.black.withOpacity(0.6)
                ],
                stops: const [0.6, 1.0],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildNebulaCloud(Color color, double size) {
    return Container(
      width: size, height: size,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: RadialGradient(
          colors: [color, Colors.transparent],
        ),
      ),
      child: BackdropFilter(
        filter: ui.ImageFilter.blur(sigmaX: 50, sigmaY: 50),
        child: Container(color: Colors.transparent),
      ),
    );
  }
}

class _StarFieldPainter extends CustomPainter {
  final double animationValue;
  final List<_Star> _stars = [];

  _StarFieldPainter({required this.animationValue}) {
    // Deterministic random stars
    final random = math.Random(42); 
    for(int i=0; i<300; i++) { // Increased from 100 to 300
      _stars.add(_Star(
        x: random.nextDouble(),
        y: random.nextDouble(),
        size: random.nextDouble() * 2 + 0.5,
        brightness: random.nextDouble() * 1.5 + 0.5, // Brighter
        speed: random.nextDouble() * 0.2 + 0.05, // Slower, deeper
      ));
    }
  }

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()..color = Colors.white;
    
    for (var star in _stars) {
      // Parallax movement
      double y = (star.y + animationValue * star.speed) % 1.0;
      
      // Twinkle Effect (Sine wave based on time + position)
      double twinkle = (math.sin(animationValue * 20 * math.pi * star.speed + star.x * 50) + 1) / 2;
      double opacity = (twinkle * 0.5 + 0.5) * star.brightness;
      opacity = opacity.clamp(0.0, 1.0); // Ensure valid range

      paint.color = Colors.white.withOpacity(opacity);
      canvas.drawCircle(Offset(star.x * size.width, y * size.height), star.size, paint);
    }
  }

  @override
  bool shouldRepaint(covariant _StarFieldPainter oldDelegate) => true;
}

class _Star {
  final double x, y, size, brightness, speed;
  _Star({required this.x, required this.y, required this.size, required this.brightness, required this.speed});
}
