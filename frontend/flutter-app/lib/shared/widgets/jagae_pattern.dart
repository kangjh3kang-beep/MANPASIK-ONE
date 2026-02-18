import 'dart:math' as math;
import 'package:flutter/material.dart';

/// 한국 전통 문양 엣지 보더 (Update: Premium Sanggam Gold Inlay)
/// 금속 상감 기법의 정교함을 표현하기 위해 이중 테두리와 메탈릭 쉐이더를 적용합니다.
class KoreanEdgeBorder extends StatelessWidget {
  final Widget child;
  final Color? borderColor;
  final double borderWidth;
  final BorderRadius borderRadius;

  const KoreanEdgeBorder({
    super.key,
    required this.child,
    this.borderColor,
    this.borderWidth = 1.0, // Thinner base width for elegance
    this.borderRadius = const BorderRadius.all(Radius.circular(16)),
  });

  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      foregroundPainter: _KoreanEdgePainter(
        borderColor: borderColor ?? const Color(0xFFD4AF37), // Metallic Gold
        width: borderWidth,
        radius: borderRadius.topLeft.x,
      ),
      child: child,
    );
  }
}

class _KoreanEdgePainter extends CustomPainter {
  final Color borderColor;
  final double width;
  final double radius;

  _KoreanEdgePainter({
    required this.borderColor,
    required this.width,
    required this.radius,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (size.isEmpty) return;

    final Rect rect = Offset.zero & size;

    // Premium Gold Metallic Gradient
    final Gradient gradient = LinearGradient(
      colors: [
        const Color(0xFFD4AF37), // Classic Gold
        const Color(0xFFFFF8C5), // High-light (Pale Gold)
        const Color(0xFFC59D2E), // Shadow Gold
        const Color(0xFFF7EF8A), // Reflected Light
        const Color(0xFFD4AF37), // Classic Gold
      ],
      stops: const [0.0, 0.2, 0.5, 0.8, 1.0],
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
    );
    
    final Paint paint = Paint()
      ..shader = gradient.createShader(rect)
      ..style = PaintingStyle.stroke
      ..strokeWidth = width;

    // 1. Double Border (Inner & Outer) to simulate inlay depth
    // Outer Line
    final RRect outerRRect = RRect.fromRectAndRadius(rect, Radius.circular(radius));
    canvas.drawRRect(outerRRect, paint);

    // Inner Line (Thinner, subtle)
    final Paint innerPaint = Paint()
      ..shader = gradient.createShader(rect)
      ..style = PaintingStyle.stroke
      ..strokeWidth = width * 0.5;
    
    final RRect innerRRect = outerRRect.deflate(3.0);
    canvas.drawRRect(innerRRect, innerPaint);
    
    // 2. Premium Corner Ornaments (Stylized Guigap/Turtle Shell)
    _drawPremiumCorners(canvas, size, paint);
  }

  void _drawPremiumCorners(Canvas canvas, Size size, Paint paint) {
    final double length = 28.0;
    final double offset = 4.0; 

    // Corner Paint (Slightly thicker for emphasis)
    final Paint cornerPaint = Paint()
      ..shader = paint.shader
      ..style = PaintingStyle.stroke
      ..strokeWidth = width * 1.5
      ..strokeCap = StrokeCap.round;

    // Top-Left Corner
    Path tl = Path();
    tl.moveTo(offset, length);
    tl.lineTo(offset, radius);
    tl.arcToPoint(Offset(radius, offset), radius: Radius.circular(radius - offset));
    tl.lineTo(length, offset);
    // Decorative Loop
    tl.moveTo(offset + 4, length - 4);
    tl.quadraticBezierTo(offset + 4, offset + 4, length - 4, offset + 4);
    canvas.drawPath(tl, cornerPaint);

    // Top-Right Corner
    Path tr = Path();
    tr.moveTo(size.width - length, offset);
    tr.lineTo(size.width - radius, offset);
    tr.arcToPoint(Offset(size.width - offset, radius), radius: Radius.circular(radius - offset));
    tr.lineTo(size.width - offset, length);
    // Decorative Loop
    tr.moveTo(size.width - length + 4, offset + 4);
    tr.quadraticBezierTo(size.width - offset - 4, offset + 4, size.width - offset - 4, length - 4);
    canvas.drawPath(tr, cornerPaint);

    // Bottom-Right Corner
    Path br = Path();
    br.moveTo(size.width - offset, size.height - length);
    br.lineTo(size.width - offset, size.height - radius);
    br.arcToPoint(Offset(size.width - radius, size.height - offset), radius: Radius.circular(radius - offset));
    br.lineTo(size.width - length, size.height - offset);
    // Decorative Loop
    br.moveTo(size.width - offset - 4, size.height - length + 4);
    br.quadraticBezierTo(size.width - offset - 4, size.height - offset - 4, size.width - length + 4, size.height - offset - 4);
    canvas.drawPath(br, cornerPaint);

    // Bottom-Left Corner
    Path bl = Path();
    bl.moveTo(length, size.height - offset);
    bl.lineTo(radius, size.height - offset);
    bl.arcToPoint(Offset(offset, size.height - radius), radius: Radius.circular(radius - offset));
    bl.lineTo(offset, size.height - length);
    // Decorative Loop
    bl.moveTo(length - 4, size.height - offset - 4);
    bl.quadraticBezierTo(offset + 4, size.height - offset - 4, offset + 4, size.height - length + 4);
    canvas.drawPath(bl, cornerPaint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

/// 자개 패턴 컨테이너 (Update: Premium Lattice Pattern)
class JagaeContainer extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry? padding;
  final double? width;
  final double? height;
  final Decoration? decoration;
  final double opacity;
  final bool showLattice;
  final Color? baseColor;

  const JagaeContainer({
    super.key,
    required this.child,
    this.padding,
    this.width,
    this.height,
    this.decoration,
    this.opacity = 0.15,
    this.showLattice = false,
    this.baseColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      height: height,
      decoration: decoration,
      padding: padding,
      child: showLattice 
          ? CustomPaint(
              painter: JagaePatternPainter(
                color: baseColor ?? (Theme.of(context).brightness == Brightness.dark 
                    ? Colors.white.withOpacity(opacity) 
                    : Colors.black.withOpacity(opacity * 0.6)), 
              ),
              child: child,
            )
          : child,
    );
  }
}

/// Premium Lattice Painter (세살문/Seosal-mun Style)
class JagaePatternPainter extends CustomPainter {
  final Color color;

  JagaePatternPainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    if (size.isEmpty) return;

    final paint = Paint()
      ..color = color.withOpacity(color.opacity * 0.8)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5; 

    // Seosal-mun (Fine Lattice) Spacing
    const double step = 20.0;
    const double innerStep = 4.0; // Double line gap

    // 1. Vertical Double Lines
    for (double x = 0; x <= size.width; x += step) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
      canvas.drawLine(Offset(x + innerStep, 0), Offset(x + innerStep, size.height), paint);
    }

    // 2. Horizontal Double Lines
    for (double y = 0; y <= size.height; y += step) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
      canvas.drawLine(Offset(0, y + innerStep), Offset(size.width, y + innerStep), paint);
    }

    // 3. Flower/Star Accents at Intersections (Octagonal feel)
    final accentPaint = Paint()
      ..color = color.withOpacity(color.opacity) // Brighter
      ..style = PaintingStyle.fill;

    for (double x = 0; x <= size.width; x += step) {
      for (double y = 0; y <= size.height; y += step) {
        // Draw small diamond at intersection
        double cx = x + innerStep / 2;
        double cy = y + innerStep / 2;
        canvas.drawCircle(Offset(cx, cy), 1.5, accentPaint);
      }
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

// Legacy support
class JagaePattern extends StatelessWidget {
  final Widget? child;
  final double opacity;
  final Color? baseColor;
  final bool showLattice;

  const JagaePattern({
    super.key,
    this.child,
    this.opacity = 0.25,
    this.baseColor,
    this.showLattice = false,
  });

  @override
  Widget build(BuildContext context) {
    return child ?? const SizedBox();
  }
}
