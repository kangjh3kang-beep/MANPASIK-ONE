import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

enum SanggamButtonStyle { gold, jagae }

class SanggamButton extends StatelessWidget {
  final String text;
  final VoidCallback? onPressed;
  final SanggamButtonStyle style;
  final double width;
  final double height;
  final Widget? icon;

  static const _goldAsset = 'assets/images/premium/premium_glossy_gold_button_texture_1771749127373.png';
  static const _jagaeAsset = 'assets/images/premium/premium_jagae_iridescent_button_texture_1771749140880.png';

  const SanggamButton({
    super.key,
    required this.text,
    this.onPressed,
    this.style = SanggamButtonStyle.gold,
    this.width = double.infinity,
    this.height = 56,
    this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final assetPath = style == SanggamButtonStyle.gold
        ? _goldAsset
        : _jagaeAsset;

    return Container(
      width: width,
      height: height,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: (style == SanggamButtonStyle.gold ? SanggamTheme.primary : Colors.white).withValues(alpha: 0.2),
            blurRadius: 12,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: Stack(
          children: [
            // Background Image Asset (The Texture)
            Positioned.fill(
              child: Image.asset(
                assetPath,
                fit: BoxFit.cover,
              ),
            ),
            // Material Ripple & Content
            Material(
              color: Colors.transparent,
              child: InkWell(
                onTap: onPressed,
                splashColor: Colors.white.withValues(alpha: 0.2),
                child: Center(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      if (icon != null) ...[
                        icon!,
                        const SizedBox(width: 8),
                      ],
                      Text(
                        text,
                        style: TextStyle(
                          color: style == SanggamButtonStyle.gold ? Colors.black : Colors.white,
                          fontFamily: 'GowunBatang',
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                          letterSpacing: 0.5,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
