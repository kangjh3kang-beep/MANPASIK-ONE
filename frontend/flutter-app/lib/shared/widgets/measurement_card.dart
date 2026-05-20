import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

// ───────────────────────────────────────────────────
// MeasurementCard — SanggamContainer 기반 측정 결과 카드
//
// [Rule 4] withOpacity 3건 → withValues(alpha:)
// [Rule 4] Colors.orange/red/green → SanggamTheme 상수
// [Rule 4] theme.colorScheme.* 4건 → SanggamTheme 직접 참조
// [Rule 4] Card+KoreanEdgeBorder+JagaeContainer 3중 → SanggamContainer
// [Rule 2] bottom 12→16, SizedBox(h:4)→8, iconBorderRadius 12→16
// ───────────────────────────────────────────────────

class MeasurementCard extends StatelessWidget {
  final DateTime date;
  final double value;
  final String unit;
  final String resultType; // 'normal', 'warning', 'danger'
  final VoidCallback? onTap;

  const MeasurementCard({
    super.key,
    required this.date,
    required this.value,
    required this.unit,
    required this.resultType,
    this.onTap,
  });

  Color _getResultColor(String type) {
    switch (type) {
      case 'warning':
        return SanggamTheme.primary; // gold — 주의
      case 'danger':
        return SanggamTheme.error; // red — 위험
      default:
        return SanggamTheme.jagaeCyan; // cyan — 정상
    }
  }

  String _getResultLabel(String type) {
    switch (type) {
      case 'warning':
        return '주의';
      case 'danger':
        return '위험';
      default:
        return '정상';
    }
  }

  @override
  Widget build(BuildContext context) {
    final color = _getResultColor(resultType);
    final label = _getResultLabel(resultType);

    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: ScaleButton(
        onPressed: onTap ?? () {},
        child: SanggamContainer(
          borderRadius: 16,
          borderWidth: 0.8,
          borderColor: color.withValues(alpha: 0.3),
          blurSigma: 8,
          jagaeOpacity: 0.05,
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Icon(
                  Icons.favorite_rounded,
                  color: color,
                  size: 24,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      label,
                      style: TextStyle(
                        color: color,
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      DateFormat('MM월 dd일 HH:mm').format(date),
                      style: const TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    value.toStringAsFixed(1),
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 22,
                      fontWeight: FontWeight.w700,
                      fontFamily: 'monospace',
                    ),
                  ),
                  Text(
                    unit,
                    style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: 10,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
