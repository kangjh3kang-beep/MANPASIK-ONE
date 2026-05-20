import 'dart:ui';

import 'package:flutter/material.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';

/// 상감(象嵌) 기법 + Glassmorphism을 결합한 커스텀 컨테이너.
///
/// - **Glassmorphism**: [BackdropFilter]로 뒤가 비치는 유리 질감
/// - **상감 금선**: Sanggam Gold 테두리로 전통 금입사 표현
/// - **자개 빛반사**: [SweepGradient]로 미세한 무지개빛 광택
///
/// ```dart
/// SanggamContainer(
///   padding: EdgeInsets.all(20),
///   child: Text('상감 오르빗'),
/// )
/// ```
class SanggamContainer extends StatelessWidget {
  /// 내부에 렌더링할 컨텐츠.
  final Widget child;

  /// 컨테이너 너비. null이면 부모에 맞춤.
  final double? width;

  /// 컨테이너 높이. null이면 자식에 맞춤.
  final double? height;

  /// 내부 여백.
  final EdgeInsetsGeometry? padding;

  /// 외부 여백.
  final EdgeInsetsGeometry? margin;

  /// 탭 콜백.
  final VoidCallback? onTap;

  /// 테두리 둥글기. 기본값 24.
  final double borderRadius;

  /// 금선 테두리 두께. 기본값 1.5.
  final double borderWidth;

  /// 금선 테두리 색상. 기본값 Sanggam Gold.
  final Color borderColor;

  /// 유리 배경색. 기본값 DeepSea Blue (opacity 0.4).
  final Color? backgroundColor;

  /// 블러 강도. 기본값 15.
  final double blurSigma;

  /// 자개 빛반사 불투명도. 기본값 0.1. 0이면 비활성.
  final double jagaeOpacity;

  /// 선택 상태 — true이면 금선이 더 밝아진다.
  final bool isSelected;

  const SanggamContainer({
    super.key,
    required this.child,
    this.width,
    this.height,
    this.padding,
    this.margin,
    this.onTap,
    this.borderRadius = 24,
    this.borderWidth = 1.5,
    this.borderColor = SanggamTheme.primary,
    this.backgroundColor,
    this.blurSigma = 15,
    this.jagaeOpacity = 0.1,
    this.isSelected = false,
  });

  @override
  Widget build(BuildContext context) {
    final radius = BorderRadius.circular(borderRadius);
    final effectiveBg = backgroundColor ??
        SanggamTheme.background.withValues(alpha: 0.4);
    final effectiveBorderColor =
        isSelected ? borderColor : borderColor.withValues(alpha: 0.6);
    final effectiveBorderWidth = isSelected ? borderWidth + 0.5 : borderWidth;

    Widget container = ClipRRect(
      borderRadius: radius,
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: blurSigma, sigmaY: blurSigma),
        child: Container(
          width: width,
          height: height,
          decoration: BoxDecoration(
            color: effectiveBg,
            borderRadius: radius,
            border: Border.all(
              color: effectiveBorderColor,
              width: effectiveBorderWidth,
            ),
            boxShadow: [
              // 깊이감을 주는 외부 그림자
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.35),
                blurRadius: 20,
                spreadRadius: 2,
                offset: const Offset(0, 8),
              ),
              // 금선 주변 은은한 발광
              BoxShadow(
                color: borderColor.withValues(alpha: 0.08),
                blurRadius: 12,
                spreadRadius: 0,
              ),
            ],
          ),
          child: Stack(
            children: [
              // ── 자개 빛반사 레이어 ──
              if (jagaeOpacity > 0)
                Positioned.fill(
                  child: IgnorePointer(
                    child: _JagaeIridescentLayer(
                      borderRadius: radius,
                      opacity: jagaeOpacity,
                    ),
                  ),
                ),

              // ── 상단 하이라이트 (유리 반사) ──
              Positioned.fill(
                child: IgnorePointer(
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: radius,
                      gradient: LinearGradient(
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        stops: const [0.0, 0.35, 1.0],
                        colors: [
                          Colors.white.withValues(alpha: 0.06),
                          Colors.white.withValues(alpha: 0.015),
                          Colors.transparent,
                        ],
                      ),
                    ),
                  ),
                ),
              ),

              // ── 좌상단 글로스 포인트 ──
              Positioned(
                top: -10,
                left: -10,
                width: 80,
                height: 80,
                child: IgnorePointer(
                  child: Container(
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      gradient: RadialGradient(
                        colors: [
                          Colors.white.withValues(alpha: 0.05),
                          Colors.transparent,
                        ],
                      ),
                    ),
                  ),
                ),
              ),

              // ── 실제 컨텐츠 ──
              Padding(
                padding: padding ?? EdgeInsets.zero,
                child: child,
              ),
            ],
          ),
        ),
      ),
    );

    // 외부 여백 적용
    if (margin != null) {
      container = Padding(padding: margin!, child: container);
    }

    // 탭 콜백 적용
    if (onTap != null) {
      container = GestureDetector(onTap: onTap, child: container);
    }

    return container;
  }
}

/// 자개(螺鈿) 빛반사 효과 레이어.
///
/// [SweepGradient]를 사용하여 Cyan → Magenta → Gold가
/// 방사형으로 순환하는 미세한 무지개빛 광택을 만든다.
class _JagaeIridescentLayer extends StatelessWidget {
  final BorderRadius borderRadius;
  final double opacity;

  const _JagaeIridescentLayer({
    required this.borderRadius,
    required this.opacity,
  });

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: borderRadius,
      child: ShaderMask(
        shaderCallback: (Rect bounds) {
          return SweepGradient(
            center: Alignment.center,
            startAngle: 0.0,
            endAngle: 6.2832, // 2π — 완전한 한 바퀴
            colors: [
              SanggamTheme.jagaeCyan.withValues(alpha: opacity),
              SanggamTheme.jagaeMagenta.withValues(alpha: opacity * 0.8),
              SanggamTheme.primary.withValues(alpha: opacity * 0.6),
              SanggamTheme.jagaeCyan.withValues(alpha: opacity * 0.4),
              SanggamTheme.jagaeMagenta.withValues(alpha: opacity * 0.7),
              SanggamTheme.jagaeCyan.withValues(alpha: opacity),
            ],
            stops: const [0.0, 0.2, 0.4, 0.6, 0.8, 1.0],
          ).createShader(bounds);
        },
        blendMode: BlendMode.srcATop,
        child: Container(
          decoration: BoxDecoration(
            borderRadius: borderRadius,
            color: Colors.white.withValues(alpha: opacity),
          ),
        ),
      ),
    );
  }
}
