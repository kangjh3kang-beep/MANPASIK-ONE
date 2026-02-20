import 'dart:math' as math;
import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// v6.0 바이오 데이터 패널 — Flutter 위젯 기반 (120×70px).
///
/// CustomPaint 내부 TextPainter(55×40px, 7px 라벨) 대비 2.2배 확대,
/// BackdropFilter + 스파크라인 + 상태 색상 인디케이터.
class BioDataPanel extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color statusColor;
  final double scanValue;

  const BioDataPanel({
    super.key,
    required this.label,
    required this.value,
    required this.icon,
    required this.statusColor,
    this.scanValue = 0.0,
  });

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(8),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
        child: Container(
          width: 120,
          height: 70,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            color: const Color(0xFF0A1628).withValues(alpha: 0.75),
            border: Border.all(
              color: AppTheme.waveCyan.withValues(alpha: 0.30),
              width: 0.5,
            ),
          ),
          child: Row(
            children: [
              // 좌측 상태 인디케이터 바
              Container(
                width: 3,
                decoration: BoxDecoration(
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.circular(8),
                    bottomLeft: Radius.circular(8),
                  ),
                  color: statusColor,
                ),
              ),
              // 콘텐츠
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(6, 5, 6, 4),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // 라벨 행 (아이콘 + 라벨)
                      Row(
                        children: [
                          Icon(icon, size: 10, color: AppTheme.waveCyan.withValues(alpha: 0.7)),
                          const SizedBox(width: 3),
                          Text(
                            label,
                            style: TextStyle(
                              color: AppTheme.waveCyan.withValues(alpha: 0.7),
                              fontSize: 11,
                              fontWeight: FontWeight.w500,
                              letterSpacing: 0.5,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 2),
                      // 값
                      Text(
                        value,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          height: 1.1,
                        ),
                      ),
                      const Spacer(),
                      // 스파크라인
                      SizedBox(
                        height: 16,
                        child: CustomPaint(
                          size: const Size(double.infinity, 16),
                          painter: _SparklinePainter(
                            color: statusColor,
                            scanValue: scanValue,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _SparklinePainter extends CustomPainter {
  final Color color;
  final double scanValue;

  _SparklinePainter({required this.color, required this.scanValue});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color.withValues(alpha: 0.5)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.0;

    final path = Path();
    path.moveTo(0, size.height / 2);
    for (double x = 0; x <= size.width; x += 2) {
      final nX = x / size.width;
      final y = size.height / 2 +
          math.sin(nX * math.pi * 4 + scanValue * math.pi * 6) * size.height * 0.35;
      path.lineTo(x, y);
    }
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant _SparklinePainter old) =>
      old.scanValue != scanValue || old.color != color;
}
