import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// HoloGlassCard — SanggamContainer 기반 홀로그래픽 글래스 카드
//
// [Rule 4] withOpacity 14건 → withValues(alpha:)
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] isDark 분기 제거 (항상 dark)
// [Rule 4] 하드코딩 Color(0xFF1B2640/050B14) → SanggamContainer 위임
// [Rule 4] Container+ClipRRect+BackdropFilter+JagaeContainer+KoreanEdgeBorder
//          5중 → SanggamContainer 단일
// [Rule 2] borderRadius 20→24 (8×3)
//
// isActive=true 시 외부 글로우는 DecoratedBox로 래핑.
// ───────────────────────────────────────────────────

class HoloGlassCard extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry padding;
  final EdgeInsetsGeometry margin;
  final VoidCallback? onTap;
  final bool isActive;
  final Color? glowColor;

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

  @override
  Widget build(BuildContext context) {
    final glow = glowColor ?? SanggamTheme.jagaeCyan;
    final borderC =
        isActive ? SanggamTheme.jagaeCyan : SanggamTheme.primary;

    Widget card = SanggamContainer(
      width: width,
      height: height,
      padding: padding,
      onTap: onTap,
      borderRadius: 24,
      borderWidth: isActive ? 1.5 : 0.8,
      borderColor: borderC.withValues(alpha: isActive ? 0.8 : 0.3),
      blurSigma: 20,
      jagaeOpacity: 0.08,
      isSelected: isActive,
      child: child,
    );

    // isActive 시 외부 네온 글로우 추가
    if (isActive) {
      card = DecoratedBox(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(24),
          boxShadow: [
            BoxShadow(
              color: glow.withValues(alpha: 0.3),
              blurRadius: 20,
              spreadRadius: -5,
            ),
            BoxShadow(
              color: SanggamTheme.primary.withValues(alpha: 0.2),
              blurRadius: 10,
            ),
          ],
        ),
        child: card,
      );
    }

    return Padding(
      padding: margin,
      child: card,
    );
  }
}
