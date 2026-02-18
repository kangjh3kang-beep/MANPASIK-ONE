import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

import 'package:manpasik/features/measurement/domain/fingerprint_analyzer.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 생체 핑거프린트 레이더 차트 (C2)
///
/// 896차원 스펙트럼 데이터를 12개 바이오마커 클러스터로 축소하여
/// RadarChart로 시각화합니다.
class FingerprintRadarChart extends StatelessWidget {
  const FingerprintRadarChart({
    super.key,
    required this.clusters,
    this.height = 280,
  });

  final List<ClusterData> clusters;
  final double height;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (clusters.isEmpty) {
      return SizedBox(
        height: height,
        child: const Center(child: Text('핑거프린트 데이터 없음')),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '생체 핑거프린트',
          style: theme.textTheme.titleSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        SizedBox(
          height: height,
          child: RadarChart(
            RadarChartData(
              radarTouchData: RadarTouchData(enabled: true),
              dataSets: [
                RadarDataSet(
                  fillColor: AppTheme.sanggamGold.withOpacity(0.2),
                  borderColor: AppTheme.sanggamGold,
                  borderWidth: 2,
                  entryRadius: 3,
                  dataEntries: clusters
                      .map((c) => RadarEntry(value: c.value * 100))
                      .toList(),
                ),
                // 이상치 오버레이
                RadarDataSet(
                  fillColor: AppTheme.dancheongRed.withOpacity(0.1),
                  borderColor: AppTheme.dancheongRed.withOpacity(0.6),
                  borderWidth: 1,
                  entryRadius: 2,
                  dataEntries: clusters
                      .map((c) => RadarEntry(value: c.anomalyScore * 100))
                      .toList(),
                ),
              ],
              radarBackgroundColor: Colors.transparent,
              borderData: FlBorderData(show: false),
              radarBorderData:
                  BorderSide(color: theme.colorScheme.outlineVariant, width: 1),
              titlePositionPercentageOffset: 0.2,
              titleTextStyle: theme.textTheme.bodySmall!.copyWith(
                fontSize: 10,
                color: theme.colorScheme.onSurfaceVariant,
              ),
              getTitle: (index, angle) =>
                  RadarChartTitle(text: clusters[index].name),
              tickCount: 4,
              ticksTextStyle: theme.textTheme.bodySmall!.copyWith(
                fontSize: 8,
                color: theme.colorScheme.onSurfaceVariant.withOpacity(0.5),
              ),
              tickBorderData: BorderSide(
                color: theme.colorScheme.outlineVariant.withOpacity(0.3),
              ),
            ),
            duration: const Duration(milliseconds: 300),
          ),
        ),
        const SizedBox(height: 8),
        // 범례
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            _LegendItem(color: AppTheme.sanggamGold, label: '측정값'),
            const SizedBox(width: 16),
            _LegendItem(color: AppTheme.dancheongRed, label: '이상치 점수'),
          ],
        ),
      ],
    );
  }
}

class _LegendItem extends StatelessWidget {
  const _LegendItem({required this.color, required this.label});
  final Color color;
  final String label;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 12,
          height: 12,
          decoration: BoxDecoration(
            color: color.withOpacity(0.3),
            border: Border.all(color: color, width: 2),
            borderRadius: BorderRadius.circular(3),
          ),
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall,
        ),
      ],
    );
  }
}
