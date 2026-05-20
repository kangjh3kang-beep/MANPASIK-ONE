import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'dart:ui';

class JagaeCard extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry padding;
  final BorderRadius? borderRadius;
  final double opacity;

  const JagaeCard({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding = const EdgeInsets.all(16.0),
    this.borderRadius,
    this.opacity = 0.08,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveRadius = borderRadius ?? BorderRadius.circular(20);

    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        borderRadius: effectiveRadius,
        border: Border.all(color: SanggamTheme.primary.withValues(alpha: 0.3), width: 1.0),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.4),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: effectiveRadius,
        child: Stack(
          children: [
            // Frosted Glass Effect
            BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 4, sigmaY: 4),
              child: Container(
                color: SanggamTheme.surfaceVariant.withValues(alpha: 0.3),
              ),
            ),
            // High-Res Jagae Texture Overlay
            Positioned.fill(
              child: Opacity(
                opacity: opacity,
                child: Image.asset(
                  'assets/images/premium/premium_jagae_iridescent_button_texture_1771749140880.png',
                  fit: BoxFit.cover,
                ),
              ),
            ),
            // Content
            Padding(
              padding: padding,
              child: child,
            ),
          ],
        ),
      ),
    );
  }
}
