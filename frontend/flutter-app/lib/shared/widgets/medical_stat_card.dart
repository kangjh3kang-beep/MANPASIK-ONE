import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

class MedicalStatCard extends StatelessWidget {
  final String label;
  final String value;
  final String unit;
  final IconData icon;
  final Color color;
  final bool isAlert;

  const MedicalStatCard({
    super.key,
    required this.label,
    required this.value,
    required this.unit,
    required this.icon,
    this.color = const Color(0xFF00E5FF),
    this.isAlert = false,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveColor = isAlert ? Colors.redAccent : color;
    
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 8, sigmaY: 8),
        child: Container(
          width: 130, // Fixed width for HUD consistency
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: const Color(0xFF05101A).withValues(alpha: 0.7),
            borderRadius: BorderRadius.circular(12),
            border: Border.all(
              color: effectiveColor.withValues(alpha: 0.5),
              width: 1.5,
            ),
            boxShadow: [
              BoxShadow(
                color: effectiveColor.withValues(alpha: 0.15),
                blurRadius: 10,
                spreadRadius: 1,
              ),
            ],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              // Header
              Row(
                children: [
                  Icon(icon, size: 14, color: effectiveColor),
                  const SizedBox(width: 6),
                  Flexible(
                    child: Text(
                      label,
                      style: TextStyle(
                        color: effectiveColor,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        letterSpacing: 1.0,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              // Value + Unit
              Row(
                crossAxisAlignment: CrossAxisAlignment.baseline,
                textBaseline: TextBaseline.alphabetic,
                children: [
                  Text(
                    value,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 22,
                      fontWeight: FontWeight.bold,
                      fontFamily: 'monospace',
                    ),
                  ),
                  const SizedBox(width: 4),
                  Text(
                    unit,
                    style: TextStyle(
                      color: Colors.white.withValues(alpha: 0.6),
                      fontSize: 10,
                    ),
                  ),
                ],
              ),
              // Mini Graph Line (Decorative)
              const SizedBox(height: 8),
              if (value != '--')
                SizedBox(
                  height: 16,
                  child: CustomPaint(
                    size: const Size(double.infinity, 16),
                    painter: _MiniGraphPainter(color: effectiveColor),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class _MiniGraphPainter extends CustomPainter {
  final Color color;
  _MiniGraphPainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    final path = Path();
    path.moveTo(0, size.height * 0.5);
    
    // Simple mock waveform
    double x = 0;
    while (x < size.width) {
       x += 5;
       double y = size.height * 0.5 + 
                  (x % 20 < 10 ? -size.height * 0.3 : size.height * 0.3) * 
                  (x % 40 > 20 ? 0.2 : 0.8); // randomness
       path.lineTo(x, y);
    }
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
