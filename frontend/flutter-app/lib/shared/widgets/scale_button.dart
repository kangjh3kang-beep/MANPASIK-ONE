import 'package:flutter/material.dart';

class ScaleButton extends StatefulWidget {
  final Widget child;
  final VoidCallback onPressed;
  final Duration duration;
  final double scale;

  const ScaleButton({
    super.key,
    required this.child,
    required this.onPressed,
    this.duration = const Duration(milliseconds: 100),
    this.scale = 0.95,
  });

  @override
  State<ScaleButton> createState() => _ScaleButtonState();
}

class _ScaleButtonState extends State<ScaleButton> with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      lowerBound: widget.scale,
      upperBound: 1.1, // Allow over-scale for hover
      duration: widget.duration,
      value: 1.0, 
    );
    // _scaleAnimation = CurvedAnimation(parent: _controller, curve: Curves.easeInOut); // REMOVED: Causes crash with values > 1.0
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _onTapDown(TapDownDetails details) {
    _controller.animateTo(widget.scale, curve: Curves.easeInOut);
  }

  void _onTapUp(TapUpDetails details) {
    _controller.animateTo(_isHovered ? 1.05 : 1.0, curve: Curves.easeOutBack); // Use elastic out for fun
    widget.onPressed();
  }

  void _onTapCancel() {
    _controller.animateTo(_isHovered ? 1.05 : 1.0, curve: Curves.easeOut);
  }

  bool _isHovered = false;

  void _onHover(bool isHovered) {
    if (widget.onPressed == null) return; // Disable hover for disabled buttons
    setState(() {
      _isHovered = isHovered;
    });
    if (isHovered) {
      _controller.animateTo(1.05, curve: Curves.easeOut); // Scale UP on hover
    } else {
      _controller.animateTo(1.0, curve: Curves.easeOut); // Return to normal
    }
  }

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onEnter: (_) => _onHover(true),
      onExit: (_) => _onHover(false),
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTapDown: _onTapDown,
        onTapUp: _onTapUp,
        onTapCancel: _onTapCancel,
        child: ScaleTransition(
          scale: _controller, // Use controller directly
          child: widget.child,
        ),
      ),
    );
  }
}
