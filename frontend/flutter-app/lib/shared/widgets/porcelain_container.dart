import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// PorcelainContainer — SanggamContainer 위임 래퍼
//
// [Rule 4] withOpacity 5건 → withValues(alpha:)
// [Rule 4] AppTheme → SanggamTheme 통일
// [Rule 4] 하드코딩 Color(0xFF1A1A1A) 등 4종 → SanggamTheme 상수
// [Rule 4] 중복 제거: Container+BoxDecoration → SanggamContainer 위임
//
// 기존 호출부 API 완전 호환 (파라미터 시그니처 동일).
// 신규 코드에서는 SanggamContainer를 직접 사용 권장.
// ───────────────────────────────────────────────────

class PorcelainContainer extends StatelessWidget {
  final Widget child;
  final double? width;
  final double? height;
  final EdgeInsetsGeometry? padding;
  final EdgeInsetsGeometry? margin;
  final VoidCallback? onTap;
  final bool isSelected;
  final Color? color;

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

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      width: width,
      height: height,
      padding: padding,
      margin: margin,
      onTap: onTap,
      borderRadius: 16,
      borderWidth: isSelected ? 1.5 : 0.8,
      borderColor: isSelected
          ? SanggamTheme.primary
          : Colors.white.withValues(alpha: 0.08),
      backgroundColor: color ?? SanggamTheme.surface.withValues(alpha: 0.5),
      blurSigma: 10,
      jagaeOpacity: 0.03,
      isSelected: isSelected,
      child: child,
    );
  }
}
