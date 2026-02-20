import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

class PorcelainContainer extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry? padding;
  final EdgeInsetsGeometry? margin;
  final VoidCallback? onTap;
  final bool isSelected;

  const PorcelainContainer({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding,
    this.margin,
    this.onTap,
    this.isSelected = false,
    this.color,
  });

  final Color? color;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Padding(
      padding: margin ?? EdgeInsets.zero,
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          width: width,
          height: height,
          padding: padding,
          decoration: BoxDecoration(
            color: color ?? (isDark ? const Color(0xFF1A1A1A) : const Color(0xFFFAFAFA)), // Off-white/Dark Grey
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: isSelected 
                  ? AppTheme.sanggamGold 
                  : (isDark ? Colors.white10 : const Color(0xFFE0E0E0)),
              width: isSelected ? 1.5 : 0.5,
            ),
            boxShadow: [
              // Soft Ambient Shadow (Porcelain Glaze feel)
              BoxShadow(
                color: isDark ? Colors.black.withOpacity(0.3) : const Color(0xFF8D6E63).withOpacity(0.05),
                blurRadius: 15,
                spreadRadius: 2,
                offset: const Offset(0, 8),
              ),
              // Rim Highlight (Gloss)
              BoxShadow(
                color: isDark ? Colors.white.withOpacity(0.05) : Colors.white.withOpacity(0.8),
                blurRadius: 1,
                spreadRadius: 0,
                offset: const Offset(-1, -1),
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: Stack(
              children: [
                // Subtle Noise Texture (Optional, using simple gradient for now)
                Positioned.fill(
                  child: Container(
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                        colors: isDark
                            ? [Colors.white.withOpacity(0.02), Colors.transparent]
                            : [Colors.white.withOpacity(0.4), Colors.transparent],
                      ),
                    ),
                  ),
                ),
                child,
              ],
            ),
          ),
        ),
      ),
    );
  }
}
