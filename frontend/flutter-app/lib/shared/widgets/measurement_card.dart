import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';

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
        return Colors.orange;
      case 'danger':
        return Colors.red;
      default:
        return Colors.green;
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
    final theme = Theme.of(context);
    final color = _getResultColor(resultType);
    final label = _getResultLabel(resultType);
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: KoreanEdgeBorder(
        borderRadius: BorderRadius.circular(16),
        child: Card(
          elevation: 0,
          color: theme.colorScheme.surface,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          margin: EdgeInsets.zero,
          child: JagaeContainer(
            opacity: 0.05,
            showLattice: false,
            decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(16),
            ),
            child: InkWell(
              onTap: onTap,
              borderRadius: BorderRadius.circular(16),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    Container(
                      width: 48,
                      height: 48,
                      decoration: BoxDecoration(
                        color: color.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(12),
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
                            style: theme.textTheme.labelMedium?.copyWith(
                              color: color,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            DateFormat('MM월 dd일 HH:mm').format(date),
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
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
                          style: theme.textTheme.headlineSmall?.copyWith(
                            fontWeight: FontWeight.w700,
                            color: theme.colorScheme.onSurface,
                          ),
                        ),
                        Text(
                          unit,
                          style: theme.textTheme.labelSmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
