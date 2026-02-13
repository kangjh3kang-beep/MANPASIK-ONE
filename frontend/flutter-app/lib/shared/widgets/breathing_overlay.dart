import 'package:flutter/material.dart';

/// [BreathingOverlay] — 호흡 애니메이션 오버레이 위젯.
///
/// 자식 위젯을 감싸 전체가 호흡하듯 미세하게 확대/축소 + 투명도 변화를
/// 반복합니다. 측정 진행 중 화면에 몰입감을 부여합니다.
///
/// MANPASIK ECOSYSTEM_DESIGN_MASTER_PLAN.md §4.2:
///   "측정 진행 시 화면 전체가 호흡하듯 움직이는 Breathing Animation"
class BreathingOverlay extends StatefulWidget {
  final Widget child;
  final Duration duration;
  final double minScale;
  final double maxScale;
  final double minOpacity;
  final double maxOpacity;
  final bool enabled;

  const BreathingOverlay({
    super.key,
    required this.child,
    this.duration = const Duration(milliseconds: 3000),
    this.minScale = 0.98,
    this.maxScale = 1.02,
    this.minOpacity = 0.85,
    this.maxOpacity = 1.0,
    this.enabled = true,
  });

  @override
  State<BreathingOverlay> createState() => _BreathingOverlayState();
}

class _BreathingOverlayState extends State<BreathingOverlay>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _scaleAnimation;
  late Animation<double> _opacityAnimation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: widget.duration,
    );

    _scaleAnimation = Tween<double>(
      begin: widget.minScale,
      end: widget.maxScale,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInOut,
    ));

    _opacityAnimation = Tween<double>(
      begin: widget.minOpacity,
      end: widget.maxOpacity,
    ).animate(CurvedAnimation(
      parent: _controller,
      curve: Curves.easeInOut,
    ));

    if (widget.enabled) {
      _controller.repeat(reverse: true);
    }
  }

  @override
  void didUpdateWidget(covariant BreathingOverlay oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.enabled && !_controller.isAnimating) {
      _controller.repeat(reverse: true);
    } else if (!widget.enabled && _controller.isAnimating) {
      _controller.stop();
      _controller.reset();
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!widget.enabled) {
      return widget.child;
    }

    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        return Opacity(
          opacity: _opacityAnimation.value,
          child: Transform.scale(
            scale: _scaleAnimation.value,
            child: child,
          ),
        );
      },
      child: widget.child,
    );
  }
}
