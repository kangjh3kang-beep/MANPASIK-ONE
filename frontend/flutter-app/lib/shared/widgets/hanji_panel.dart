import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// HanjiPanel — SanggamContainer 위임 래퍼
//
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] 하드코딩 Color(0xFF1E232E), Color(0xFFF9F6F0) → SanggamTheme.surface
// [Rule 4] 중복 제거: Container+_HanjiNoisePainter → SanggamContainer(jagaeOpacity)
// [Rule 2] padding 20px → 16px (8px×2)
// [Rule 1] BackdropFilter 부재 → SanggamContainer(blurSigma:8) 적용
//
// isDarkTheme 파라미터는 하위 호환용 유지 (앱은 항상 dark 모드).
// 신규 코드에서는 SanggamContainer를 직접 사용 권장.
// ───────────────────────────────────────────────────

/// 한지(韓紙) 질감 패널 — SanggamContainer 기반.
class HanjiPanel extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry padding;
  final double? width;
  final double? height;
  final BorderRadius? borderRadius;
  final bool isDarkTheme;

  const HanjiPanel({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.all(16),
    this.width,
    this.height,
    this.borderRadius,
    this.isDarkTheme = true,
  });

  @override
  Widget build(BuildContext context) {
    final r = borderRadius?.topLeft.x ?? 16.0;

    return SanggamContainer(
      width: width,
      height: height,
      padding: padding,
      borderRadius: r,
      borderWidth: 0.8,
      borderColor: SanggamTheme.primary.withValues(alpha: 0.3),
      backgroundColor: SanggamTheme.surface.withValues(alpha: 0.6),
      blurSigma: 8,
      jagaeOpacity: 0.02,
      child: child,
    );
  }
}
