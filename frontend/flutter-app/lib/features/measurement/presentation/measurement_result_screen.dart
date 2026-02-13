import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:fl_chart/fl_chart.dart';

/// 측정 결과 화면
///
/// 최근 측정 값 요약 + GetMeasurementHistory 기반 트렌드 차트.
class MeasurementResultScreen extends ConsumerWidget {
  const MeasurementResultScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final historyAsync = ref.watch(measurementHistoryProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('측정 결과'),
      ),
      body: historyAsync.when(
        data: (result) {
          if (result.items.isEmpty) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.analytics_outlined, size: 64, color: theme.colorScheme.outline),
                  const SizedBox(height: 16),
                  Text(
                    '측정 기록이 없습니다',
                    style: theme.textTheme.titleMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 24),
                  FilledButton.icon(
                    onPressed: () => context.go('/measurement'),
                    icon: const Icon(Icons.add),
                    label: const Text('측정하기'),
                  ),
                ],
              ),
            );
          }
          final latest = result.items.first;
          final spots = result.items.asMap().entries.map((e) {
            final i = result.items.length - 1 - e.key;
            return FlSpot(i.toDouble(), e.value.primaryValue);
          }).toList().reversed.toList();

          return SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // 최근 측정 요약 카드
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      children: [
                        Text(
                          '최근 측정',
                          style: theme.textTheme.labelLarge?.copyWith(
                            color: theme.colorScheme.primary,
                          ),
                        ),
                        const SizedBox(height: 12),
                        Text(
                          '${latest.primaryValue.toStringAsFixed(1)} ${latest.unit}',
                          style: theme.textTheme.headlineMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        if (latest.cartridgeType.isNotEmpty)
                          Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              latest.cartridgeType,
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: theme.colorScheme.onSurfaceVariant,
                              ),
                            ),
                          ),
                      ],
                    ),
                  ),
                ),
                const SizedBox(height: 32),
                // 트렌드 차트
                Text(
                  '트렌드',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                SizedBox(
                  height: 220,
                  child: spots.isEmpty
                      ? const Center(child: Text('데이터 없음'))
                      : LineChart(
                          LineChartData(
                            gridData: FlGridData(show: true),
                            titlesData: FlTitlesData(
                              leftTitles: AxisTitles(
                                sideTitles: SideTitles(
                                  showTitles: true,
                                  reservedSize: 36,
                                  getTitlesWidget: (value, meta) => Text(
                                    value.toInt().toString(),
                                    style: theme.textTheme.bodySmall,
                                  ),
                                ),
                              ),
                              bottomTitles: AxisTitles(
                                sideTitles: SideTitles(
                                  showTitles: true,
                                  reservedSize: 24,
                                  getTitlesWidget: (value, meta) {
                                    final idx = value.toInt();
                                    if (idx >= 0 && idx < result.items.length) {
                                      final item = result.items[result.items.length - 1 - idx];
                                      final at = item.measuredAt;
                                      if (at != null) {
                                        return Text(
                                          '${at.month}/${at.day}',
                                          style: theme.textTheme.bodySmall,
                                        );
                                      }
                                    }
                                    return Text('', style: theme.textTheme.bodySmall);
                                  },
                                ),
                              ),
                              topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                              rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                            ),
                            borderData: FlBorderData(show: true),
                            lineBarsData: [
                              LineChartBarData(
                                spots: spots,
                                isCurved: true,
                                color: theme.colorScheme.primary,
                                barWidth: 2,
                                dotData: const FlDotData(show: true),
                                belowBarData: BarAreaData(
                                  show: true,
                                  color: theme.colorScheme.primary.withValues(alpha: 0.1),
                                ),
                              ),
                            ],
                            minX: 0,
                            maxX: spots.isEmpty ? 1 : (spots.length - 1).toDouble(),
                            minY: spots.isEmpty ? 0 : (spots.map((s) => s.y).reduce((a, b) => a < b ? a : b) - 5).clamp(0, double.infinity),
                            maxY: spots.isEmpty ? 100 : (spots.map((s) => s.y).reduce((a, b) => a > b ? a : b) + 5),
                          ),
                          duration: const Duration(milliseconds: 250),
                        ),
                ),
                const SizedBox(height: 32),
                OutlinedButton.icon(
                  onPressed: () => context.go('/measurement'),
                  icon: const Icon(Icons.refresh),
                  label: const Text('다시 측정'),
                  style: OutlinedButton.styleFrom(
                    minimumSize: const Size(double.infinity, 48),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  ),
                ),
              ],
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, _) => Center(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Text(
              '기록을 불러올 수 없습니다.\n$err',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.error),
            ),
          ),
        ),
      ),
    );
  }
}
