import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

class BreathingGlow extends StatefulWidget {
  final Widget child;
  final Color? glowColor;

  const BreathingGlow({
    super.key,
    required this.child,
    this.glowColor,
  });

  @override
  State<BreathingGlow> createState() => _BreathingGlowState();
}

class _BreathingGlowState extends State<BreathingGlow> with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _scale;
  late Animation<double> _opacity;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat(reverse: true);

    _scale = Tween<double>(begin: 1.0, end: 1.02).animate(
      CurvedAnimation(parent: _controller, curve: Curves.easeInOut),
    );

    _opacity = Tween<double>(begin: 0.2, end: 0.6).animate(
      CurvedAnimation(parent: _controller, curve: Curves.easeInOut),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final color = widget.glowColor ?? AppTheme.sanggamGold;

    return Stack(
      alignment: Alignment.center,
      children: [
        // 1. The Glowing Shadow (Behind)
        Positioned.fill(
          child: AnimatedBuilder(
            animation: _controller,
            builder: (context, child) {
              return Transform.scale(
                scale: _scale.value,
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(20), // Match HoloGlassCard
                    color: Colors.transparent, // Important: No solid color
                    boxShadow: [
                      BoxShadow(
                        color: color.withOpacity(_opacity.value),
                        blurRadius: 30, // Softer, wider diffusion
                        spreadRadius: 5,
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
        // 2. The Actual Content
        widget.child,
      ],
    );
  }
}
