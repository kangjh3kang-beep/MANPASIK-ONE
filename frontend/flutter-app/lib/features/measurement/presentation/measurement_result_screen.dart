import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:fl_chart/fl_chart.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/shared/providers/ecosystem_providers.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/measurement/domain/fingerprint_analyzer.dart';
import 'package:manpasik/shared/widgets/fingerprint_radar_chart.dart';
import 'package:manpasik/shared/widgets/fingerprint_heatmap.dart';
import 'package:manpasik/shared/widgets/untargeted_analysis_card.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'package:manpasik/features/measurement/presentation/widgets/measurement_evidence_badge.dart';

// ───────────────────────────────────────────────────
// MeasurementResultScreen — Sanggam Orbit 측정 결과
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 4x 제거
// [Rule 4] theme.textTheme.* ~15x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~7x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 5x → SanggamTheme.primary
// [Rule 4] AppTheme.dancheongRed → SanggamTheme.error
// [Rule 4] Colors.green 2x → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange 2x → SanggamTheme.primary
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] withOpacity 3x → withValues(alpha:)
// [Rule 4] Card 3x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] h:12→16, h:20→24, padding:20→24, horizontal:12→16,
//          vertical:4→8, horizontal:10→16, borderRadius:12/20→16,
//          bottom:6→8
// ───────────────────────────────────────────────────

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
  bool _completionFired = false;

  void _fireMeasurementCompletion() {
    if (_completionFired) return;
    _completionFired = true;
    ref.read(measurementCompletionProvider.notifier).update((s) => s + 1);
    ref.invalidate(measurementHistoryProvider);
  }

  Future<void> _runAiAnalysis(MeasurementHistoryItem latest,
      List<MeasurementHistoryItem> allItems) async {
    if (_isAnalyzing || _aiAnalysis != null) return;
    setState(() => _isAnalyzing = true);

    final recentValues = allItems.take(10).map((i) => i.primaryValue).toList();

    final analysis = await RustBridge.analyzeResult(
      value: latest.primaryValue,
      biomarker:
          latest.cartridgeType.isNotEmpty ? latest.cartridgeType : 'glucose',
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
    final historyAsync = ref.watch(measurementHistoryProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back, color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '측정 결과',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: historyAsync.when(
                data: (result) {
                  if (result.hasError && result.items.isEmpty) {
                    return Center(
                      child: Padding(
                        padding: const EdgeInsets.all(24),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(
                              Icons.sync_problem_rounded,
                              size: 64,
                              color: SanggamTheme.error,
                            ),
                            const SizedBox(height: 16),
                            const Text(
                              '측정 기록을 확인할 수 없습니다',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                color: SanggamTheme.error,
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Text(
                              result.errorMessage!,
                              textAlign: TextAlign.center,
                              style: const TextStyle(
                                color: SanggamTheme.onSurfaceDim,
                                fontSize: 13,
                              ),
                            ),
                            const SizedBox(height: 24),
                            FilledButton.icon(
                              onPressed: () =>
                                  ref.invalidate(measurementHistoryProvider),
                              icon: const Icon(Icons.refresh),
                              label: const Text('다시 시도'),
                            ),
                          ],
                        ),
                      ),
                    );
                  }
                  if (result.items.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Icon(Icons.analytics_outlined,
                              size: 64, color: SanggamTheme.onSurfaceDim),
                          const SizedBox(height: 16),
                          const Text(
                            '측정 기록이 없습니다',
                            style: TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 24),
                          FilledButton.icon(
                            onPressed: () => context.go('/measure'),
                            icon: const Icon(Icons.add),
                            label: const Text('측정하기'),
                            style: FilledButton.styleFrom(
                              backgroundColor: SanggamTheme.primary,
                              foregroundColor: SanggamTheme.background,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(16),
                              ),
                            ),
                          ),
                        ],
                      ),
                    );
                  }

                  final latest = result.items.first;

                  // 측정 완료 이벤트 → 홈/데이터허브 자동 갱신
                  WidgetsBinding.instance.addPostFrameCallback((_) {
                    _fireMeasurementCompletion();
                  });

                  // AI 분석 자동 실행
                  if (_aiAnalysis == null && !_isAnalyzing) {
                    WidgetsBinding.instance.addPostFrameCallback((_) {
                      _runAiAnalysis(latest, result.items);
                    });
                  }

                  final spots = result.items
                      .asMap()
                      .entries
                      .map((e) {
                        final i = result.items.length - 1 - e.key;
                        return FlSpot(i.toDouble(), e.value.primaryValue);
                      })
                      .toList()
                      .reversed
                      .toList();

                  return SingleChildScrollView(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        if (result.hasError) ...[
                          _buildHistoryWarning(result.errorMessage!),
                          const SizedBox(height: 16),
                        ],
                        // 최근 측정 요약 카드
                        _buildLatestCard(latest),
                        const SizedBox(height: 24),

                        // AI 분석 카드
                        _buildAiAnalysisCard(),
                        const SizedBox(height: 24),

                        // 트렌드 차트
                        const Text(
                          '트렌드',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 16),
                        _buildChart(spots, result.items),
                        const SizedBox(height: 24),

                        // 핑거프린트 시각화 (C2)
                        _buildFingerprintSection(),
                        const SizedBox(height: 32),

                        // 크로스 도메인 CTA 버튼
                        FilledButton.icon(
                          onPressed: () => context.go('/data'),
                          icon: const Icon(Icons.hub_outlined),
                          label: const Text('데이터 허브에서 상세 보기'),
                          style: FilledButton.styleFrom(
                            minimumSize: const Size(double.infinity, 48),
                            backgroundColor: SanggamTheme.primary,
                            foregroundColor: SanggamTheme.background,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                          ),
                        ),
                        const SizedBox(height: 8),
                        OutlinedButton.icon(
                          onPressed: () => context.go('/coach'),
                          icon: const Icon(Icons.psychology_outlined),
                          label: const Text('AI 코치 분석 받기'),
                          style: OutlinedButton.styleFrom(
                            minimumSize: const Size(double.infinity, 48),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                          ),
                        ),
                        const SizedBox(height: 8),
                        OutlinedButton.icon(
                          onPressed: () => context.go('/measure'),
                          icon: const Icon(Icons.refresh),
                          label: const Text('다시 측정'),
                          style: OutlinedButton.styleFrom(
                            minimumSize: const Size(double.infinity, 48),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
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
                      style: const TextStyle(
                        color: SanggamTheme.error,
                        fontSize: 14,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLatestCard(MeasurementHistoryItem latest) {
    final riskColor = _aiAnalysis != null
        ? _riskColor(_aiAnalysis!.riskLevel)
        : SanggamTheme.primary;

    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          Text(
            '최근 측정',
            style: TextStyle(
              color: riskColor,
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 16),
          Text(
            '${latest.primaryValue.toStringAsFixed(1)} ${latest.unit}',
            style: TextStyle(
              color: riskColor,
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          if (latest.cartridgeType.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(top: 8),
              child: Text(
                latest.cartridgeType,
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                ),
              ),
            ),
          const SizedBox(height: 12),
          MeasurementEvidenceBadge(
            evidenceStatus: latest.evidenceStatus,
            diagnosticReady: latest.diagnosticReady,
            evidenceGaps: latest.evidenceGaps,
          ),
          if (_aiAnalysis != null) ...[
            const SizedBox(height: 16),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: riskColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Text(
                _riskLabel(_aiAnalysis!.riskLevel),
                style: TextStyle(
                  color: riskColor,
                  fontSize: 12,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildHistoryWarning(String message) {
    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(16),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(
            Icons.sync_problem_rounded,
            color: SanggamTheme.error,
            size: 20,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              message,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 13,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAiAnalysisCard() {
    if (_isAnalyzing) {
      return const SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.all(24),
        child: Row(
          children: [
            SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(strokeWidth: 2),
            ),
            SizedBox(width: 16),
            Text(
              'AI 건강 분석 중...',
              style: TextStyle(
                color: Colors.white,
                fontSize: 14,
              ),
            ),
          ],
        ),
      );
    }

    if (_aiAnalysis == null) {
      return const SizedBox.shrink();
    }

    final analysis = _aiAnalysis!;
    final riskColor = _riskColor(analysis.riskLevel);

    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.smart_toy_rounded,
                  color: SanggamTheme.primary, size: 20),
              const SizedBox(width: 8),
              const Text(
                'AI 건강 분석',
                style: TextStyle(
                  color: SanggamTheme.primary,
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const Spacer(),
              Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                decoration: BoxDecoration(
                  color: riskColor.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(16),
                ),
                child: Text(
                  '${analysis.healthScore.toStringAsFixed(0)}점',
                  style: TextStyle(
                    color: riskColor,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 요약
          Text(
            analysis.summary,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
            ),
          ),
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
                    ? SanggamTheme.jagaeCyan
                    : analysis.trend == 'declining'
                        ? SanggamTheme.primary
                        : SanggamTheme.onSurfaceDim,
              ),
              const SizedBox(width: 8),
              Text(
                '추세: ${_trendLabel(analysis.trend)}',
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 추천 사항
          ...analysis.recommendations.map((rec) => Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Icon(Icons.lightbulb_outline,
                        size: 16, color: SanggamTheme.primary),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        rec,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 12,
                        ),
                      ),
                    ),
                  ],
                ),
              )),
        ],
      ),
    );
  }

  Widget _buildChart(List<FlSpot> spots, List<MeasurementHistoryItem> items) {
    if (spots.isEmpty) {
      return const SizedBox(
        height: 220,
        child: Center(
          child: Text(
            '데이터 없음',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
        ),
      );
    }

    return SizedBox(
      height: 220,
      child: LineChart(
        LineChartData(
          gridData: const FlGridData(show: true),
          titlesData: FlTitlesData(
            leftTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 36,
                getTitlesWidget: (value, meta) => Text(
                  value.toInt().toString(),
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
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
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      );
                    }
                  }
                  return const SizedBox.shrink();
                },
              ),
            ),
            topTitles:
                const AxisTitles(sideTitles: SideTitles(showTitles: false)),
            rightTitles:
                const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          ),
          borderData: FlBorderData(show: true),
          lineBarsData: [
            LineChartBarData(
              spots: spots,
              isCurved: true,
              color: SanggamTheme.primary,
              barWidth: 2,
              dotData: const FlDotData(show: true),
              belowBarData: BarAreaData(
                show: true,
                color: SanggamTheme.primary.withValues(alpha: 0.1),
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
      'normal' => SanggamTheme.jagaeCyan,
      'caution' => SanggamTheme.primary,
      'warning' => SanggamTheme.error,
      'critical' => const Color(0xFF8B0000),
      _ => SanggamTheme.onSurfaceDim,
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
  Widget _buildFingerprintSection() {
    // 실측정 차동 데이터가 있으면 fromDifferentialData 사용,
    // 없으면 DEMO 모드로 시뮬레이션 데이터 사용
    final fingerprintData = FingerprintAnalyzer.fromDifferentialData([]);
    final clusters = FingerprintAnalyzer.reduceTo12Clusters(fingerprintData);
    final heatmapGrid = FingerprintAnalyzer.toHeatmapGrid(fingerprintData);
    final anomalies = FingerprintAnalyzer.detectAnomalies(fingerprintData);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        FingerprintRadarChart(clusters: clusters),
        const SizedBox(height: 24),
        FingerprintHeatmap(grid: heatmapGrid),
        const SizedBox(height: 24),
        UntargetedAnalysisCard(
          anomalies: anomalies,
          clusters: clusters,
        ),
      ],
    );
  }
}
