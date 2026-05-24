import 'dart:ui';
import 'package:flutter/material.dart';

/// Living Digital Twin 대시보드의 'Foreground UI Layer'에 배치될 Glassmorphism 카드
/// [PDCA Design] 3D 홀로그램과 겹치지 않으면서 투명하고 직관적인 정보 전달.
class MedicalGlassCard extends StatefulWidget {
  final Widget child;
  final VoidCallback? onTap;
  final bool isAlert;
  final Color? customAccentColor;
  final EdgeInsetsGeometry padding;
  final double width;

  const MedicalGlassCard({
    super.key,
    required this.child,
    this.onTap,
    this.isAlert = false,
    this.customAccentColor,
    this.padding = const EdgeInsets.all(16),
    this.width = 180.0,
  });

  @override
  State<MedicalGlassCard> createState() => _MedicalGlassCardState();
}

class _MedicalGlassCardState extends State<MedicalGlassCard>
    with SingleTickerProviderStateMixin {
  bool _isHovered = false;
  late AnimationController _pulseController;
  late Animation<double> _pulseAnimation;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1500),
    );
    _pulseAnimation = Tween<double>(begin: 0.4, end: 0.8).animate(
      CurvedAnimation(parent: _pulseController, curve: Curves.easeInOut),
    );

    if (widget.isAlert) {
      _pulseController.repeat(reverse: true);
    }
  }

  @override
  void didUpdateWidget(covariant MedicalGlassCard oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.isAlert && !oldWidget.isAlert) {
      _pulseController.repeat(reverse: true);
    } else if (!widget.isAlert && oldWidget.isAlert) {
      _pulseController.stop();
      _pulseController.value = 0.0;
    }
  }

  @override
  void dispose() {
    _pulseController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // Design Intelligence - Color Palette
    const Color darkSlate = Color(0xFF0F172A);
    const Color warningRed = Color(0xFFEF4444);
    final Color accentColor = widget.customAccentColor ?? const Color(0xFF22D3EE);
    
    final Color currentAccent = widget.isAlert ? warningRed : accentColor;
    
    return Semantics(
      container: true,
      button: widget.onTap != null,
      child: RepaintBoundary(
        child: MouseRegion(
          cursor: widget.onTap != null
              ? SystemMouseCursors.click
              : SystemMouseCursors.basic,
          onEnter: (_) => setState(() => _isHovered = true),
          onExit: (_) => setState(() => _isHovered = false),
          child: AnimatedScale(
            scale: _isHovered ? 1.02 : 1.0,
            duration: const Duration(milliseconds: 200),
            curve: Curves.easeOutQuad,
            child: AnimatedBuilder(
              animation: _pulseAnimation,
              builder: (context, child) {
                // Pulse effect on alert
                final double opacity = widget.isAlert 
                    ? _pulseAnimation.value 
                    : (_isHovered ? 0.8 : 0.6);
                    
                return ClipRRect(
                  borderRadius: BorderRadius.circular(16),
                  child: BackdropFilter(
                    filter: ImageFilter.blur(sigmaX: 16.0, sigmaY: 16.0),
                    child: Container(
                      width: widget.width,
                      decoration: BoxDecoration(
                        color: darkSlate.withValues(alpha: widget.isAlert ? 0.7 : 0.5),
                        borderRadius: BorderRadius.circular(16),
                        border: Border.all(
                          color: currentAccent.withValues(alpha: opacity),
                          width: _isHovered || widget.isAlert ? 2.0 : 1.0,
                        ),
                        boxShadow: [
                          if (_isHovered || widget.isAlert)
                            BoxShadow(
                              color: currentAccent.withValues(alpha: 0.3),
                              blurRadius: 12,
                              spreadRadius: 2,
                            )
                        ],
                      ),
                      child: Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: widget.onTap,
                          splashColor: currentAccent.withValues(alpha: 0.2),
                          highlightColor: currentAccent.withValues(alpha: 0.1),
                          child: Padding(
                            padding: widget.padding,
                            child: widget.child,
                          ),
                        ),
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
        ),
      ),
    );
  }
}
