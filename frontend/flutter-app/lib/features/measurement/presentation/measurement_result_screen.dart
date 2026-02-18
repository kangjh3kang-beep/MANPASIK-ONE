import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:fl_chart/fl_chart.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/measurement/domain/fingerprint_analyzer.dart';
import 'package:manpasik/shared/widgets/fingerprint_radar_chart.dart';
import 'package:manpasik/shared/widgets/fingerprint_heatmap.dart';
import 'package:manpasik/shared/widgets/untargeted_analysis_card.dart';

/// 측정 결과 화면
///
/// 최근 측정 값 요약 + GetMeasurementHistory 기반 트렌드 차트 + AI 분석.
class MeasurementResultScreen extends ConsumerStatefulWidget {
  const MeasurementResultScreen({super.key});

  @override
  ConsumerState<MeasurementResultScreen> createState() =>
      _MeasurementResultScreenState();
}

class _MeasurementResultScreenState
    extends ConsumerState<MeasurementResultScreen> {
  AiAnalysisDto? _aiAnalysis;
  bool _isAnalyzing = false;

  Future<void> _runAiAnalysis(MeasurementHistoryItem latest,
      List<MeasurementHistoryItem> allItems) async {
    if (_isAnalyzing || _aiAnalysis != null) return;
    setState(() => _isAnalyzing = true);

    final recentValues = allItems
        .take(10)
        .map((i) => i.primaryValue)
        .toList();

    final analysis = await RustBridge.analyzeResult(
      value: latest.primaryValue,
      biomarker: latest.cartridgeType.isNotEmpty
          ? latest.cartridgeType
          : 'glucose',
      unit: latest.unit.isNotEmpty ? latest.unit : 'mg/dL',
      recentValues: recentValues,
    );

    if (!mounted) return;
    setState(() {
      _aiAnalysis = analysis;
      _isAnalyzing = false;
    });
  }

  @override
  Widget build(BuildContext context) {
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
                  Icon(Icons.analytics_outlined,
                      size: 64, color: theme.colorScheme.outline),
                  const SizedBox(height: 16),
                  Text(
                    '측정 기록이 없습니다',
                    style: theme.textTheme.titleMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  const SizedBox(height: 24),
                  FilledButton.icon(
                    onPressed: () => context.go('/measure'),
                    icon: const Icon(Icons.add),
                    label: const Text('측정하기'),
                  ),
                ],
              ),
            );
          }

          final latest = result.items.first;

          // AI 분석 자동 실행
          if (_aiAnalysis == null && !_isAnalyzing) {
            WidgetsBinding.instance.addPostFrameCallback((_) {
              _runAiAnalysis(latest, result.items);
            });
          }

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
                _buildLatestCard(theme, latest),
                const SizedBox(height: 24),

                // AI 분석 카드
                _buildAiAnalysisCard(theme),
                const SizedBox(height: 24),

                // 트렌드 차트
                Text(
                  '트렌드',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                _buildChart(theme, spots, result.items),
                const SizedBox(height: 24),

                // 핑거프린트 시각화 (C2)
                _buildFingerprintSection(theme),
                const SizedBox(height: 32),

                OutlinedButton.icon(
                  onPressed: () => context.go('/measure'),
                  icon: const Icon(Icons.refresh),
                  label: const Text('다시 측정'),
                  style: OutlinedButton.styleFrom(
                    minimumSize: const Size(double.infinity, 48),
                    shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16)),
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
              style: theme.textTheme.bodyMedium
                  ?.copyWith(color: theme.colorScheme.error),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildLatestCard(ThemeData theme, MeasurementHistoryItem latest) {
    final riskColor = _aiAnalysis != null
        ? _riskColor(_aiAnalysis!.riskLevel)
        : theme.colorScheme.primary;

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            Text(
              '최근 측정',
              style: theme.textTheme.labelLarge?.copyWith(
                color: riskColor,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              '${latest.primaryValue.toStringAsFixed(1)} ${latest.unit}',
              style: theme.textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: riskColor,
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
            if (_aiAnalysis != null) ...[
              const SizedBox(height: 12),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                decoration: BoxDecoration(
                  color: riskColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Text(
                  _riskLabel(_aiAnalysis!.riskLevel),
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: riskColor,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildAiAnalysisCard(ThemeData theme) {
    if (_isAnalyzing) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Row(
            children: [
              const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2),
              ),
              const SizedBox(width: 16),
              Text('AI 건강 분석 중...',
                  style: theme.textTheme.bodyMedium),
            ],
          ),
        ),
      );
    }

    if (_aiAnalysis == null) return const SizedBox.shrink();

    final analysis = _aiAnalysis!;
    final riskColor = _riskColor(analysis.riskLevel);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.smart_toy_rounded,
                    color: AppTheme.sanggamGold, size: 20),
                const SizedBox(width: 8),
                Text(
                  'AI 건강 분석',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: AppTheme.sanggamGold,
                  ),
                ),
                const Spacer(),
                // 건강 점수
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: riskColor.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    '${analysis.healthScore.toStringAsFixed(0)}점',
                    style: theme.textTheme.bodySmall?.copyWith(
                      color: riskColor,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),

            // 요약
            Text(analysis.summary, style: theme.textTheme.bodyMedium),
            const SizedBox(height: 8),

            // 트렌드
            Row(
              children: [
                Icon(
                  analysis.trend == 'improving'
                      ? Icons.trending_down
                      : analysis.trend == 'declining'
                          ? Icons.trending_up
                          : Icons.trending_flat,
                  size: 16,
                  color: analysis.trend == 'improving'
                      ? Colors.green
                      : analysis.trend == 'declining'
                          ? Colors.orange
                          : theme.colorScheme.onSurfaceVariant,
                ),
                const SizedBox(width: 4),
                Text(
                  '추세: ${_trendLabel(analysis.trend)}',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // 추천 사항
            ...analysis.recommendations.map((rec) => Padding(
                  padding: const EdgeInsets.only(bottom: 6),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Icon(Icons.lightbulb_outline,
                          size: 16, color: AppTheme.sanggamGold),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(rec,
                            style: theme.textTheme.bodySmall),
                      ),
                    ],
                  ),
                )),
          ],
        ),
      ),
    );
  }

  Widget _buildChart(ThemeData theme, List<FlSpot> spots,
      List<MeasurementHistoryItem> items) {
    if (spots.isEmpty) {
      return const SizedBox(
          height: 220, child: Center(child: Text('데이터 없음')));
    }

    return SizedBox(
      height: 220,
      child: LineChart(
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
                  if (idx >= 0 && idx < items.length) {
                    final item = items[items.length - 1 - idx];
                    final at = item.measuredAt;
                    if (at != null) {
                      return Text(
                        '${at.month}/${at.day}',
                        style: theme.textTheme.bodySmall,
                      );
                    }
                  }
                  return const SizedBox.shrink();
                },
              ),
            ),
            topTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false)),
            rightTitles: const AxisTitles(
                sideTitles: SideTitles(showTitles: false)),
          ),
          borderData: FlBorderData(show: true),
          lineBarsData: [
            LineChartBarData(
              spots: spots,
              isCurved: true,
              color: AppTheme.sanggamGold,
              barWidth: 2,
              dotData: const FlDotData(show: true),
              belowBarData: BarAreaData(
                show: true,
                color: AppTheme.sanggamGold.withOpacity(0.1),
              ),
            ),
          ],
          minX: 0,
          maxX: spots.isEmpty ? 1 : (spots.length - 1).toDouble(),
          minY: spots.isEmpty
              ? 0
              : (spots.map((s) => s.y).reduce((a, b) => a < b ? a : b) - 5)
                  .clamp(0, double.infinity),
          maxY: spots.isEmpty
              ? 100
              : (spots.map((s) => s.y).reduce((a, b) => a > b ? a : b) + 5),
        ),
        duration: const Duration(milliseconds: 250),
      ),
    );
  }

  Color _riskColor(String risk) {
    return switch (risk) {
      'normal' => Colors.green,
      'caution' => Colors.orange,
      'warning' => AppTheme.dancheongRed,
      'critical' => const Color(0xFF8B0000),
      _ => Colors.grey,
    };
  }

  String _riskLabel(String risk) {
    return switch (risk) {
      'normal' => '정상',
      'caution' => '주의',
      'warning' => '경고',
      'critical' => '위험',
      _ => '-',
    };
  }

  String _trendLabel(String trend) {
    return switch (trend) {
      'improving' => '개선 중',
      'declining' => '악화 추세',
      'stable' => '안정',
      _ => '-',
    };
  }

  /// 핑거프린트 시각화 섹션 (C2 + C3)
  Widget _buildFingerprintSection(ThemeData theme) {
    final simulated = FingerprintAnalyzer.generateSimulatedData();
    final clusters = FingerprintAnalyzer.reduceTo12Clusters(simulated);
    final heatmapGrid = FingerprintAnalyzer.toHeatmapGrid(simulated);
    final anomalies = FingerprintAnalyzer.detectAnomalies(simulated);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // 레이더 차트
        FingerprintRadarChart(clusters: clusters),
        const SizedBox(height: 20),
        // 히트맵
        FingerprintHeatmap(grid: heatmapGrid),
        const SizedBox(height: 20),
        // 비표적 분석 (C3)
        UntargetedAnalysisCard(
          anomalies: anomalies,
          clusters: clusters,
        ),
      ],
    );
  }
}
