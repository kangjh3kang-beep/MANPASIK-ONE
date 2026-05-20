import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

import 'package:manpasik/features/measurement/domain/fingerprint_analyzer.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

/// 비표적 분석 결과 카드 (C3)
///
/// 12개 바이오마커 클러스터의 이상치 점수를 BarChart로 시각화하고,
/// 감지된 이상 항목에 대한 설명을 표시합니다.
class UntargetedAnalysisCard extends StatelessWidget {
  const UntargetedAnalysisCard({
    super.key,
    required this.anomalies,
    required this.clusters,
    this.chartHeight = 180,
  });

  final List<AnomalyResult> anomalies;
  final List<ClusterData> clusters;
  final double chartHeight;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: SanggamTheme.surfaceVariant.withValues(alpha: 0.3),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.biotech_rounded,
                    color: SanggamTheme.primary, size: 20),
                const SizedBox(width: 8),
                const Text(
                  '비표적 분석',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Spacer(),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: anomalies.isEmpty
                        ? SanggamTheme.jagaeCyan.withValues(alpha: 0.1)
                        : SanggamTheme.error.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    anomalies.isEmpty
                        ? '정상'
                        : '${anomalies.length}건 감지',
                    style: TextStyle(
                      color: anomalies.isEmpty
                          ? SanggamTheme.jagaeCyan
                          : SanggamTheme.error,
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // 바 차트: 클러스터별 이상치 점수
            if (clusters.isNotEmpty)
              SizedBox(
                height: chartHeight,
                child: BarChart(
                  BarChartData(
                    alignment: BarChartAlignment.spaceAround,
                    maxY: 1.0,
                    barTouchData: BarTouchData(
                      touchTooltipData: BarTouchTooltipData(
                        getTooltipItem: (group, groupIndex, rod, rodIndex) {
                          return BarTooltipItem(
                            '${clusters[group.x].name}\n${(rod.toY * 100).toStringAsFixed(0)}%',
                            const TextStyle(
                                color: Colors.white, fontSize: 12),
                          );
                        },
                      ),
                    ),
                    titlesData: FlTitlesData(
                      show: true,
                      bottomTitles: AxisTitles(
                        sideTitles: SideTitles(
                          showTitles: true,
                          reservedSize: 32,
                          getTitlesWidget: (value, meta) {
                            final idx = value.toInt();
                            if (idx >= 0 && idx < clusters.length) {
                              return Padding(
                                padding: const EdgeInsets.only(top: 4),
                                child: Text(
                                  clusters[idx].name.length > 3
                                      ? clusters[idx].name.substring(0, 3)
                                      : clusters[idx].name,
                                  style: const TextStyle(
                                    color: SanggamTheme.onSurfaceDim,
                                    fontSize: 8,
                                  ),
                                ),
                              );
                            }
                            return const SizedBox.shrink();
                          },
                        ),
                      ),
                      leftTitles: const AxisTitles(
                          sideTitles: SideTitles(showTitles: false)),
                      topTitles: const AxisTitles(
                          sideTitles: SideTitles(showTitles: false)),
                      rightTitles: const AxisTitles(
                          sideTitles: SideTitles(showTitles: false)),
                    ),
                    borderData: FlBorderData(show: false),
                    gridData: const FlGridData(show: false),
                    barGroups: clusters.asMap().entries.map((e) {
                      final color = e.value.anomalyScore > 0.25
                          ? SanggamTheme.error
                          : e.value.anomalyScore > 0.15
                              ? SanggamTheme.primary
                              : SanggamTheme.jagaeCyan;
                      return BarChartGroupData(
                        x: e.key,
                        barRods: [
                          BarChartRodData(
                            toY: e.value.anomalyScore,
                            color: color,
                            width: 12,
                            borderRadius: const BorderRadius.vertical(
                                top: Radius.circular(4)),
                          ),
                        ],
                      );
                    }).toList(),
                  ),
                  duration: const Duration(milliseconds: 250),
                ),
              ),

            // 이상 항목 설명
            if (anomalies.isNotEmpty) ...[
              const SizedBox(height: 16),
              const Divider(color: SanggamTheme.surfaceVariant),
              const SizedBox(height: 8),
              ...anomalies.take(3).map((a) => Padding(
                    padding: const EdgeInsets.only(bottom: 8),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Icon(
                          Icons.warning_amber_rounded,
                          size: 16,
                          color: a.severity == AnomalySeverity.high
                              ? SanggamTheme.error
                              : SanggamTheme.primary,
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            a.description,
                            style: const TextStyle(
                              color: Colors.white70,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ],
                    ),
                  )),
            ],

            if (anomalies.isEmpty) ...[
              const SizedBox(height: 12),
              const Text(
                '모든 바이오마커가 정상 범위 내에 있습니다.',
                style: TextStyle(
                  color: SanggamTheme.jagaeCyan,
                  fontSize: 12,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
