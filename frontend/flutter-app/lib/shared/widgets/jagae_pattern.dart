import 'dart:math';
import 'package:flutter/material.dart';

/// 한국 전통 문양 엣지 보더 (금속 상감 기법 느낌)
class KoreanEdgeBorder extends StatelessWidget {
  final Widget child;
  final Color? borderColor;
  final double borderWidth;
  final BorderRadius borderRadius;

  const KoreanEdgeBorder({
    super.key,
    required this.child,
    this.borderColor,
    this.borderWidth = 2.0,
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

    final Paint paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = width;
      
    final Rect rect = Offset.zero & size;
    
    // Gold Gradient Shader
    final Gradient gradient = LinearGradient(
      colors: [
        const Color(0xFFD4AF37), // Gold
        const Color(0xFFF7EF8A), // Light Gold
        const Color(0xFFB48E26), // Dark Bronze
        const Color(0xFFD4AF37), // Gold
      ],
      stops: const [0.0, 0.3, 0.7, 1.0],
      begin: Alignment.topLeft,
      end: Alignment.bottomRight,
    );
    
    paint.shader = gradient.createShader(rect);

    // 1. 기본 테두리
    final RRect rrect = RRect.fromRectAndRadius(rect, Radius.circular(radius));
    canvas.drawRRect(rrect, paint);
    
    // 2. 모서리 장식
    _drawCornerOrnament(canvas, size, paint);
  }

  void _drawCornerOrnament(Canvas canvas, Size size, Paint paint) {
    final double s = 24.0;
    
    // 모서리 장식용 얇은 붓
    final Paint ornamentPaint = Paint()
      ..shader = paint.shader
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.5;

    // TL (Top-Left)
    Path p1 = Path();
    p1.moveTo(0, s);
    p1.lineTo(0, s*0.4);
    p1.quadraticBezierTo(0, 0, s*0.4, 0);
    p1.lineTo(s, 0);
    // 내부 문양 (Greek Key / Man-ja Style)
    p1.moveTo(6, 6);
    p1.lineTo(16, 6);
    p1.lineTo(16, 16);
    p1.lineTo(6, 16);
    p1.close();
    canvas.drawPath(p1, ornamentPaint);

    // TR (Top-Right)
    Path p2 = Path();
    p2.moveTo(size.width - s, 0);
    p2.lineTo(size.width - s*0.4, 0);
    p2.quadraticBezierTo(size.width, 0, size.width, s*0.4);
    p2.lineTo(size.width, s);
    // 문양
    p2.moveTo(size.width - 6, 6);
    p2.lineTo(size.width - 16, 6);
    p2.lineTo(size.width - 16, 16);
    p2.lineTo(size.width - 6, 16);
    p2.close();
    canvas.drawPath(p2, ornamentPaint);

    // BR (Bottom-Right)
    Path p3 = Path();
    p3.moveTo(size.width, size.height - s);
    p3.lineTo(size.width, size.height - s*0.4);
    p3.quadraticBezierTo(size.width, size.height, size.width - s*0.4, size.height);
    p3.lineTo(size.width - s, size.height);
    // 문양
    p3.moveTo(size.width - 6, size.height - 6);
    p3.lineTo(size.width - 16, size.height - 6);
    p3.lineTo(size.width - 16, size.height - 16);
    p3.lineTo(size.width - 6, size.height - 16);
    p3.close();
    canvas.drawPath(p3, ornamentPaint);
    
    // BL (Bottom-Left)
    Path p4 = Path();
    p4.moveTo(s, size.height);
    p4.lineTo(s*0.4, size.height);
    p4.quadraticBezierTo(0, size.height, 0, size.height - s*0.4);
    p4.lineTo(0, size.height - s);
    // 문양
    p4.moveTo(6, size.height - 6);
    p4.lineTo(16, size.height - 6);
    p4.lineTo(16, size.height - 16);
    p4.lineTo(6, size.height - 16);
    p4.close();
    canvas.drawPath(p4, ornamentPaint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

/// 자개 패턴이 적용된 컨테이너 (사용자 요청으로 패턴 제거, 그라디언트만 유지)
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
        // 사용자 피드백 반영: 비정형 자개 패턴(overlay) 제거.
        // Container의 decoration(그라디언트/그림자)만 유지.
    return Container(
      width: width,
      height: height,
      decoration: decoration,
          padding: padding,
      child: child,
    );
  }
}

// 하위 호환성을 위해 남겨둔 빈 클래스 로직 (실제로는 사용되지 않음)
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
