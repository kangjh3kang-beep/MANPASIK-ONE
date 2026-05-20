import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// MedicalStatCard — SanggamContainer 기반 생체 지표 카드
//
// [Rule 4] withOpacity 3건 → withValues(alpha:)
// [Rule 4] Color(0xFF00E5FF) → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF05101A) → SanggamContainer backgroundColor
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] 중복 ClipRRect+BackdropFilter → SanggamContainer 위임
// [Rule 2] borderRadius 12→16, padding 12→16, SizedBox(w:6)→8
// [Rule 1] 기존 BackdropFilter 유지 (SanggamContainer 내장)
// ───────────────────────────────────────────────────

class MedicalStatCard extends StatelessWidget {
  final String label;
  final String value;
  final String unit;
  final IconData icon;
  final Color color;
  final bool isAlert;
  final String? tooltipMessage;

  const MedicalStatCard({
    super.key,
    required this.label,
    required this.value,
    required this.unit,
    required this.icon,
    this.color = SanggamTheme.jagaeCyan,
    this.isAlert = false,
    this.tooltipMessage,
  });

  @override
  Widget build(BuildContext context) {
    final c = isAlert ? SanggamTheme.error : color;

    final card = SanggamContainer(
      padding: const EdgeInsets.all(16),
      borderRadius: 16,
      borderWidth: 1.5,
      borderColor: c.withValues(alpha: 0.5),
      blurSigma: 8,
      jagaeOpacity: 0.03,
      child: LayoutBuilder(builder: (context, constraints) {
        final boundedWidth =
            constraints.hasBoundedWidth && constraints.maxWidth.isFinite
                ? constraints.maxWidth
                : 180.0;
        return FittedBox(
          fit: BoxFit.scaleDown,
          alignment: Alignment.topLeft,
          child: SizedBox(
            width: boundedWidth,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                // Header: icon + label
                Row(
                  children: [
                    Icon(icon, size: 14, color: c),
                    const SizedBox(width: 8),
                    Flexible(
                      child: Text(
                        label,
                        style: TextStyle(
                          color: c,
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
                      style: const TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 10,
                      ),
                    ),
                  ],
                ),

                // Mini waveform graph
                const SizedBox(height: 8),
                if (value != '--')
                  SizedBox(
                    height: 16,
                    child: CustomPaint(
                      size: const Size(double.infinity, 16),
                      painter: _MiniGraphPainter(color: c),
                    ),
                  ),
              ],
            ),
          ),
        );
      }),
    );

    if (tooltipMessage != null) {
      return Tooltip(message: tooltipMessage!, child: card);
    }
    return card;
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
      ..style = PaintingStyle.stroke;

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
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
