import 'package:flutter/material.dart';
import 'dart:ui';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart'; // Add import

/// MANPASIK R&D Lab - Sanggam Inlay Decoration
/// Provides traditional gold inlay effect and glassmorphism for Flutter widgets.
class SanggamContainer extends StatelessWidget {
  final Widget child;
  final double borderRadius;
  final bool showGlow;
  final double borderWidth;
  final Color? backgroundColor;

  const SanggamContainer({
    super.key,
    required this.child,
    this.borderRadius = 16.0,
    this.showGlow = true,
    this.borderWidth = 1.5, // Thicker default border
    this.backgroundColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(borderRadius),
        boxShadow: showGlow ? [
          BoxShadow(
            color: AppTheme.sanggamGold.withOpacity(0.2), // Increased glow opacity
            blurRadius: 15,
            spreadRadius: 2,
          )
        ] : null,
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(borderRadius),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
          child: CustomPaint(
            painter: SanggamPainter(
              borderRadius: borderRadius,
              borderWidth: borderWidth,
            ),
            foregroundPainter: JagaePatternPainter( // Use the public painter if possible, or duplicate logic
               color: Colors.white.withOpacity(0.15), // Increased from 0.05
            ), 
            child: Container(
              padding: const EdgeInsets.all(1), // Border alignment
              color: backgroundColor ?? AppTheme.deepSeaBlue.withOpacity(0.8), // Slightly darker for contrast
              child: child,
            ),
          ),
        ),
      ),
    );
  }
}

class SanggamPainter extends CustomPainter {
  final double borderRadius;
  final double borderWidth;

  SanggamPainter({
    required this.borderRadius,
    required this.borderWidth,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final rect = Offset.zero & size;
    final rrect = RRect.fromRectAndRadius(rect, Radius.circular(borderRadius));

    // 1. Foundation Border (Fine Gold Line)
    final paint = Paint()
      ..color = AppTheme.sanggamGold.withOpacity(0.3)
      ..style = PaintingStyle.stroke
      ..strokeWidth = borderWidth;

    canvas.drawRRect(rrect, paint);

    // 2. Corner Highlights (The 'Sanggam' touch)
    final highlightPaint = Paint()
      ..color = AppTheme.sanggamGold
      ..style = PaintingStyle.stroke
      ..strokeWidth = borderWidth * 1.5
      ..strokeCap = StrokeCap.round;

    // Drawing L-shaped accents at corners
    final cornerLength = borderRadius * 1.5;
    
    // Top Left
    canvas.drawPath(
      Path()
        ..moveTo(0, cornerLength)
        ..lineTo(0, borderRadius)
        ..arcToPoint(Offset(borderRadius, 0), radius: Radius.circular(borderRadius))
        ..lineTo(cornerLength, 0),
      highlightPaint,
    );

    // Bottom Right
    canvas.drawPath(
      Path()
        ..moveTo(size.width - cornerLength, size.height)
        ..lineTo(size.width - borderRadius, size.height)
        ..arcToPoint(Offset(size.width, size.height - borderRadius), radius: Radius.circular(borderRadius), clockwise: false)
        ..lineTo(size.width, size.height - cornerLength),
      highlightPaint,
    );
  }

  @override
  bool shouldRepaint(SanggamPainter oldDelegate) => false;
}
