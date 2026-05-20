import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// GlassmorphismCard — SanggamContainer 위임 래퍼
//
// [Rule 4] 중복 제거: 자체 BackdropFilter+BoxDecoration 코드를
//          SanggamContainer에 위임하여 단일 디자인 시스템 유지.
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Color(0xFF080C17) → SanggamTheme.background
// [Rule 2] borderRadius 20 → 24 (8px×3)
//
// 기존 호출부 API 완전 호환 (파라미터 시그니처 동일).
// 신규 코드에서는 SanggamContainer를 직접 사용 권장.
// ───────────────────────────────────────────────────

class GlassmorphismCard extends StatelessWidget {
  final Widget child;
  final double blur;
  final double opacity;
  final Color baseColor;
  final EdgeInsetsGeometry padding;
  final BorderRadius? borderRadius;
  final bool useKoreanBorder;
  final double? width;
  final double? height;

  const GlassmorphismCard({
    super.key,
    required this.child,
    this.blur = 12,
    this.opacity = 0.15,
    this.baseColor = SanggamTheme.background,
    this.padding = const EdgeInsets.all(16),
    this.borderRadius,
    this.useKoreanBorder = true,
    this.width,
    this.height,
  });

  @override
  Widget build(BuildContext context) {
    // borderRadius가 BorderRadius 객체로 전달되면 첫 번째 모서리 값 추출,
    // 미지정 시 SanggamContainer 기본값(24) 사용.
    final r = borderRadius?.topLeft.x ?? 24.0;

    return SanggamContainer(
      width: width,
      height: height,
      padding: padding,
      borderRadius: r,
      borderWidth: useKoreanBorder ? 1.5 : 0.8,
      borderColor: useKoreanBorder
          ? SanggamTheme.primary
          : Colors.white.withValues(alpha: 0.1),
      backgroundColor: baseColor.withValues(alpha: opacity),
      blurSigma: blur,
      jagaeOpacity: 0.05,
      child: child,
    );
  }
}
