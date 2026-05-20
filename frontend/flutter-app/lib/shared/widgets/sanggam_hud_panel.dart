import 'dart:ui';

import 'package:flutter/material.dart';

/// Sanggam Orbit hybrid HUD panel.
///
/// Layer 1: photorealistic frame asset.
/// Layer 2: glass blur + content overlay.
/// Bottom action: photorealistic jagae button asset.
class SanggamHudPanel extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final String? title;
  final EdgeInsetsGeometry padding;
  final VoidCallback? onActionTap;

  const SanggamHudPanel({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.title,
    this.padding = const EdgeInsets.all(20),
    this.onActionTap,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: width,
      height: height,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(24),
        child: Stack(
          children: [
            // Layer 2: glass blur overlay.
            Positioned.fill(
              child: BackdropFilter(
                filter: ImageFilter.blur(sigmaX: 12, sigmaY: 12),
                child: Container(
                  color: Colors.black.withValues(alpha: 0.28),
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(24),
                    border: Border.all(
                      color: const Color(0xFFD4AF37).withValues(alpha: 0.15),
                      width: 0.5,
                    ),
                  ),
                ),
              ),
            ),

            // Foreground content.
            Padding(
              padding: padding,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (title != null && title!.isNotEmpty) ...[
                    Text(
                      title!,
                      style: const TextStyle(
                        color: Color(0xFFEAD9A0),
                        fontSize: 16,
                        fontWeight: FontWeight.w800,
                        letterSpacing: 1.2,
                      ),
                    ),
                    const SizedBox(height: 12),
                  ],
                  Expanded(child: child),
                  const SizedBox(height: 12),
                  Align(
                    alignment: Alignment.bottomRight,
                    child: _JagaeActionButton(
                      onTap: onActionTap,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _JagaeActionButton extends StatelessWidget {
  final VoidCallback? onTap;

  const _JagaeActionButton({
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: SizedBox(
        width: 64,
        height: 64,
        child: Stack(
          alignment: Alignment.center,
          children: [
            Image.asset(
              'assets/images/btn_3d_action_slim.png',
              width: 72,
              height: 72,
              fit: BoxFit.contain,
              colorBlendMode: BlendMode.screen,
              color: Colors.white,
              errorBuilder: (context, error, stackTrace) {
                return Container(
                  width: 58,
                  height: 58,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: const LinearGradient(
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                      colors: [Color(0xFF7CE9F5), Color(0xFF1F8DA8)],
                    ),
                    border:
                        Border.all(color: const Color(0xFFD4AF37), width: 2),
                    boxShadow: [
                      BoxShadow(
                        color: const Color(0xFF00E5FF).withValues(alpha: 0.35),
                        blurRadius: 14,
                        offset: const Offset(0, 6),
                      ),
                    ],
                  ),
                );
              },
            ),
            const Icon(
              Icons.chevron_right_rounded,
              color: Colors.white,
              size: 28,
            ),
          ],
        ),
      ),
    );
  }
}
