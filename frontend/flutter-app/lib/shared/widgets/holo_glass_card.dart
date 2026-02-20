import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';

class HoloGlassCard extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry padding;
  final EdgeInsetsGeometry margin;
  final VoidCallback? onTap;
  final bool isActive; // For glowing active state

  const HoloGlassCard({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding = const EdgeInsets.all(16),
    this.margin = EdgeInsets.zero,
    this.onTap,
    this.isActive = false,
    this.glowColor,
  });

  final Color? glowColor;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Padding(
      padding: margin,
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          width: width,
          height: height,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(20),
            // Neon Glow Shadow (Outer)
            boxShadow: [
              if (isActive) ...[
                BoxShadow(
                  color: glowColor?.withOpacity(0.4) ?? 
                         (isDark ? AppTheme.waveCyan.withOpacity(0.3) : AppTheme.celadonTeal.withOpacity(0.2)),
                  blurRadius: 20,
                  spreadRadius: -5,
                ),
                BoxShadow(
                  color: AppTheme.sanggamGold.withOpacity(0.2),
                  blurRadius: 10,
                  spreadRadius: 0,
                ),
              ] else
                BoxShadow(
                  color: isDark ? Colors.black.withOpacity(0.5) : Colors.black.withOpacity(0.05),
                  blurRadius: 15,
                  spreadRadius: 0,
                  offset: const Offset(0, 8),
                ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(20),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 20, sigmaY: 20), // Increased blur for frostier look
              child: Stack(
                children: [
                  // 1. Ultra-Thin Glass Layer (Low Opacity)
                  Positioned.fill(
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: isDark ? Colors.white.withOpacity(0.08) : Colors.black.withOpacity(0.05), // Subtle border
                          width: 0.5,
                        ),
                        gradient: LinearGradient(
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                          colors: isDark 
                            ? [
                                const Color(0xFF1B2640).withOpacity(0.25), // Much clearer
                                const Color(0xFF050B14).withOpacity(0.40), // Slightly darker bottom
                              ]
                            : [
                                Colors.white.withOpacity(0.6), // White Frost
                                Colors.white.withOpacity(0.3), 
                              ],
                        ),
                      ),
                    ),
                  ),
                  
                  // 2. Reflective Shine (Mystical Gloss)
                  Positioned.fill(
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(20),
                        gradient: LinearGradient(
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                          stops: const [0.0, 0.4, 1.0],
                          colors: [
                            Colors.white.withOpacity(isDark ? 0.12 : 0.4), // Stronger shine in light mode
                            Colors.white.withOpacity(0.0),  // Transparent middle
                            Colors.white.withOpacity(isDark ? 0.05 : 0.1), // Subtle reflect bottom-right
                          ],
                        ),
                      ),
                    ),
                  ),

                  // 3. Jagae/Data Texture
                  const Positioned.fill(
                    child: JagaeContainer(
                      showLattice: true,
                      opacity: 0.08, // Increased slightly for visibility
                      child: SizedBox(),
                    ),
                  ),

                  // 4. Sanggam Border (Gradient Overlay)
                  Positioned.fill(
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(color: Colors.transparent),
                      ),
                      child: KoreanEdgeBorder(
                        borderWidth: isActive ? 1.5 : 1.0, // Thinner border
                        borderColor: isActive ? AppTheme.waveCyan : AppTheme.sanggamGold,
                        child: const SizedBox.expand(),
                      ),
                    ),
                  ),

                  // 5. Content
                  Padding(
                    padding: padding,
                    child: child,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
