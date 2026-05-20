import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

import 'package:manpasik/features/measurement/domain/fingerprint_analyzer.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

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
    if (clusters.isEmpty) {
      return SizedBox(
        height: height,
        child: const Center(
          child: Text(
            '핑거프린트 데이터 없음',
            style: TextStyle(color: SanggamTheme.onSurfaceDim),
          ),
        ),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '생체 핑거프린트',
          style: TextStyle(
            color: Colors.white,
            fontSize: 14,
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
                  fillColor: SanggamTheme.primary.withValues(alpha: 0.2),
                  borderColor: SanggamTheme.primary,
                  borderWidth: 2,
                  entryRadius: 3,
                  dataEntries: clusters
                      .map((c) => RadarEntry(value: c.value * 100))
                      .toList(),
                ),
                // 이상치 오버레이
                RadarDataSet(
                  fillColor: SanggamTheme.error.withValues(alpha: 0.1),
                  borderColor: SanggamTheme.error.withValues(alpha: 0.6),
                  borderWidth: 1,
                  entryRadius: 2,
                  dataEntries: clusters
                      .map((c) => RadarEntry(value: c.anomalyScore * 100))
                      .toList(),
                ),
              ],
              radarBackgroundColor: Colors.transparent,
              borderData: FlBorderData(show: false),
              radarBorderData: BorderSide(
                color: SanggamTheme.surfaceVariant.withValues(alpha: 0.5),
                width: 1,
              ),
              titlePositionPercentageOffset: 0.2,
              titleTextStyle: const TextStyle(
                fontSize: 10,
                color: SanggamTheme.onSurfaceDim,
              ),
              getTitle: (index, angle) =>
                  RadarChartTitle(text: clusters[index].name),
              tickCount: 4,
              ticksTextStyle: TextStyle(
                fontSize: 8,
                color: SanggamTheme.onSurfaceDim.withValues(alpha: 0.5),
              ),
              tickBorderData: BorderSide(
                color: SanggamTheme.surfaceVariant.withValues(alpha: 0.3),
              ),
            ),
            duration: const Duration(milliseconds: 300),
          ),
        ),
        const SizedBox(height: 8),
        // 범례
        const Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            _LegendItem(color: SanggamTheme.primary, label: '측정값'),
            SizedBox(width: 16),
            _LegendItem(color: SanggamTheme.error, label: '이상치 점수'),
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
            color: color.withValues(alpha: 0.3),
            border: Border.all(color: color, width: 2),
            borderRadius: BorderRadius.circular(3),
          ),
        ),
        const SizedBox(width: 4),
        Text(
          label,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          ),
        ),
      ],
    );
  }
}
