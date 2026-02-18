import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 백자 (Porcelain) 컨테이너
///
/// 'White Mode'용 프리미엄 컨테이너.
/// - 순백색 도자기의 매끄러운 질감
/// - 가장자리의 은은한 금빛 테두리 (Sanggam)
/// - 부드러운 주변광 그림자 (Ambient Shadow)
class PorcelainContainer extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry padding;
  final double? width;
  final double? height;
  final VoidCallback? onTap;

  const PorcelainContainer({
    super.key,
    required this.child,
    this.padding = EdgeInsets.zero,
    this.width,
    this.height,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final container = Container(
      width: width,
      height: height,
      padding: padding,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(
          color: AppTheme.sanggamGold.withOpacity(0.15), // 은은한 금테
          width: 0.5,
        ),
        boxShadow: [
          // 1. 부드러운 주변광 (Soft Ambient) - Increased contrast
          BoxShadow(
            color: const Color(0xFF4A4A4A).withOpacity(0.08), // Darker grey shadow
            blurRadius: 16,
            offset: const Offset(0, 6),
            spreadRadius: 2,
          ),
          // 2. 미세한 윤곽선 강조 (Rim Light)
          const BoxShadow(
            color: Colors.white,
            blurRadius: 0,
            offset: Offset(-1, -1),
            spreadRadius: 0,
          ),
        ],
      ),
      child: child,
    );

    if (onTap != null) {
      return GestureDetector(onTap: onTap, child: container);
    }

    return container;
  }
}
