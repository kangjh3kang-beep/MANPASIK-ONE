import 'package:flutter/material.dart';
import 'package:manpasik/shared/widgets/medical_glass_card.dart';

// ───────────────────────────────────────────────────
// MedicalStatCard — Living Digital Twin 디자인이 적용된 생체 지표 카드
// [PDCA Design] High Contrast 다크 모드, 44x44px 터치 보장, 사이버네틱 무드.
// ───────────────────────────────────────────────────

class MedicalStatCard extends StatelessWidget {
  final String label;
  final String value;
  final String unit;
  final IconData icon;
  final Color? color;
  final bool isAlert;
  final String? tooltipMessage;
  final VoidCallback? onTap;

  const MedicalStatCard({
    super.key,
    required this.label,
    required this.value,
    required this.unit,
    required this.icon,
    this.color,
    this.isAlert = false,
    this.tooltipMessage,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    const Color cyberCyan = Color(0xFF22D3EE);
    const Color warningRed = Color(0xFFEF4444);
    final Color accentColor = isAlert ? warningRed : (color ?? cyberCyan);

    Widget content = Semantics(
      label: '$label, $value $unit, ${isAlert ? '경고 상태' : '정상 상태'}',
      child: MedicalGlassCard(
        isAlert: isAlert,
        customAccentColor: accentColor,
        onTap: onTap,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            // Header: icon + label
            Row(
              children: [
                Icon(icon, size: 14, color: accentColor),
                const SizedBox(width: 8),
                Flexible(
                  child: Text(
                    label,
                    style: TextStyle(
                      color: accentColor,
                      fontSize: 12, // 디자인 명세: 12px
                      fontWeight: FontWeight.bold,
                      letterSpacing: 1.0,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),

            // Value + Unit
            Row(
              crossAxisAlignment: CrossAxisAlignment.baseline,
              textBaseline: TextBaseline.alphabetic,
              children: [
                Text(
                  value,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 24, // 디자인 명세: 24px
                    fontWeight: FontWeight.w700,
                    fontFamily: 'monospace',
                  ),
                ),
                const SizedBox(width: 6),
                Text(
                  unit,
                  style: const TextStyle(
                    color: Color(0xFFCBD5E1), // Slate 300
                    fontSize: 10,
                  ),
                ),
              ],
            ),

            // Mini waveform graph (생체 리듬 시각화)
            const SizedBox(height: 12),
            if (value != '--')
              SizedBox(
                height: 16,
                child: RepaintBoundary(
                  child: CustomPaint(
                    size: const Size(double.infinity, 16),
                    painter: _MiniGraphPainter(color: accentColor),
                  ),
                ),
              ),
          ],
        ),
      ),
    );

    if (tooltipMessage != null) {
      return Tooltip(
        message: tooltipMessage!,
        child: content,
      );
    }
    
    return content;
  }
}

// ═══════════════════════════════════════════════════
//  Mini waveform graph painter
// ═══════════════════════════════════════════════════

class _MiniGraphPainter extends CustomPainter {
  final Color color;
  _MiniGraphPainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke
      ..strokeJoin = StrokeJoin.round;

    final path = Path()..moveTo(0, size.height * 0.5);

    double x = 0;
    while (x < size.width) {
      x += 5;
      final y = size.height * 0.5 +
          (x % 20 < 10 ? -size.height * 0.3 : size.height * 0.3) *
              (x % 40 > 20 ? 0.2 : 0.8);
      path.lineTo(x, y);
    }
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant _MiniGraphPainter oldDelegate) => 
    oldDelegate.color != color;
}
