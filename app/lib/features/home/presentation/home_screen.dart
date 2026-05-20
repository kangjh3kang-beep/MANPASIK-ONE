import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';
import '../../../core/state/global_providers.dart';
import '../../../core/theme/app_theme.dart';
import '../../measure/data/measurement_local_repository.dart';

final recentMeasurementsProvider = FutureProvider((ref) async {
  final repo = await ref.watch(measurementLocalRepoProvider.future);
  final list = await repo.getRecentMeasurements(limit: 7);
  return list.reversed.toList(); // 시간순 정렬 (과거->현재)
});

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final contextCardsAsync = ref.watch(contextCardFeedProvider);

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text('Living ManPaSik'),
        flexibleSpace: ClipRect(
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
            child: Container(color: AppTheme.backgroundDark.withValues(alpha: 0.5)),
          ),
        ),
      ),
      body: Stack(
        children: [
          // 배경 은은한 네온 글로우 효과
          Positioned(
            top: -100, left: -50,
            child: Container(
              width: 300, height: 300,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: AppTheme.primaryNeonPurple.withValues(alpha: 0.2),
              ),
              child: BackdropFilter(filter: ImageFilter.blur(sigmaX: 100, sigmaY: 100), child: Container()),
            ),
          ),
          
          CustomScrollView(
            physics: const BouncingScrollPhysics(),
            slivers: [
              SliverPadding(
                padding: const EdgeInsets.fromLTRB(20, 120, 20, 20),
                sliver: SliverToBoxAdapter(
                  child: _buildHealthDashboard(context, ref),
                ),
              ),
              
              SliverPadding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                sliver: SliverToBoxAdapter(
                  child: _buildTrendChart(ref),
                ),
              ),
              
              const SliverPadding(
                padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                sliver: SliverToBoxAdapter(
                  child: Text(
                    "지능형 인사이트 피드", 
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white, letterSpacing: 0.5),
                  ),
                ),
              ),

              contextCardsAsync.when(
                data: (cards) => SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, index) {
                      return _buildPremiumContextCard(cards[index]);
                    },
                    childCount: cards.length,
                  ),
                ),
                loading: () => const SliverToBoxAdapter(
                  child: Padding(padding: EdgeInsets.all(40), child: Center(child: CircularProgressIndicator(color: AppTheme.primaryNeonTeal))),
                ),
                error: (error, _) => SliverToBoxAdapter(child: Center(child: Text('오류 발생: $error'))),
              ),
              const SliverPadding(padding: EdgeInsets.only(bottom: 100)), // FAB 가림 방지
            ],
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.go('/measure/progress'),
        icon: const Icon(Icons.blur_on), // 생체 인식/스캔 프리미엄 아이콘 느낌
        label: const Text('Measure Now', style: TextStyle(fontWeight: FontWeight.bold)),
      ),
    );
  }

  Widget _buildHealthDashboard(BuildContext context, WidgetRef ref) {
    final recentAsync = ref.watch(recentMeasurementsProvider);
    final currentScore = recentAsync.valueOrNull?.lastOrNull?.healthScore ?? 0;

    return Container(
      height: 160,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(24),
        gradient: AppTheme.premiumGradient,
        boxShadow: [
          BoxShadow(
            color: AppTheme.primaryNeonPurple.withValues(alpha: 0.4),
            blurRadius: 20,
            offset: const Offset(0, 10),
          )
        ],
      ),
      child: Stack(
        children: [
          // Glassmorphism 오버레이
          ClipRRect(
            borderRadius: BorderRadius.circular(24),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 5, sigmaY: 5),
              child: Container(
                color: Colors.white.withValues(alpha: 0.05),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(24.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text('통합 건강 스코어', style: TextStyle(color: Colors.white.withValues(alpha: 0.9), fontSize: 16)),
                    const SizedBox(height: 8),
                    Text(currentScore == 0 ? '--' : '$currentScore', style: const TextStyle(color: Colors.white, fontSize: 48, fontWeight: FontWeight.bold, height: 1.0)),
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                      decoration: BoxDecoration(color: Colors.black26, borderRadius: BorderRadius.circular(12)),
                      child: const Text('↑ 상위 5% 최상위권', style: TextStyle(color: AppTheme.primaryNeonTeal, fontSize: 12, fontWeight: FontWeight.bold)),
                    ),
                  ],
                ),
                // 자가진단 센터 진입 버튼
                GestureDetector(
                  onTap: () => context.go('/home/diagnostics'),
                  child: Container(
                    width: 80, height: 80,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: AppTheme.surfaceDark.withValues(alpha: 0.5),
                      border: Border.all(color: AppTheme.primaryNeonTeal.withValues(alpha: 0.3), width: 4),
                    ),
                    child: const Center(child: Icon(Icons.troubleshoot, color: AppTheme.primaryNeonTeal, size: 36)),
                  ),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTrendChart(WidgetRef ref) {
    final recentAsync = ref.watch(recentMeasurementsProvider);

    return recentAsync.when(
      data: (measurements) {
        if (measurements.isEmpty) {
          return Container(
            height: 150,
            margin: const EdgeInsets.only(bottom: 20),
            decoration: BoxDecoration(color: AppTheme.surfaceDark, borderRadius: BorderRadius.circular(20)),
            child: const Center(child: Text('최근 측정 기록이 없습니다.\n"Measure Now"를 눌러주세요.', textAlign: TextAlign.center, style: TextStyle(color: Colors.white54))),
          );
        }

        final spots = measurements.asMap().entries.map((e) {
          return FlSpot(e.key.toDouble(), e.value.healthScore.toDouble());
        }).toList();

        return Container(
          height: 150,
          margin: const EdgeInsets.only(bottom: 20),
          padding: const EdgeInsets.fromLTRB(16, 24, 16, 10),
          decoration: BoxDecoration(
            color: AppTheme.surfaceDark,
            borderRadius: BorderRadius.circular(20),
            border: Border.all(color: Colors.white.withValues(alpha: 0.05)),
          ),
          child: LineChart(
            LineChartData(
              gridData: const FlGridData(show: false),
              titlesData: const FlTitlesData(show: false),
              borderData: FlBorderData(show: false),
              minY: 0, maxY: 100,
              lineBarsData: [
                LineChartBarData(
                  spots: spots,
                  isCurved: true,
                  curveSmoothness: 0.35,
                  gradient: AppTheme.premiumGradient,
                  barWidth: 4,
                  isStrokeCapRound: true,
                  dotData: const FlDotData(show: false),
                  belowBarData: BarAreaData(
                    show: true,
                    gradient: LinearGradient(
                      colors: [AppTheme.primaryNeonTeal.withValues(alpha: 0.3), AppTheme.primaryNeonPurple.withValues(alpha: 0.0)],
                      begin: Alignment.topCenter, end: Alignment.bottomCenter,
                    ),
                  ),
                )
              ],
            ),
          ),
        );
      },
      loading: () => const SizedBox(height: 150, child: Center(child: CircularProgressIndicator(color: AppTheme.primaryNeonTeal))),
      error: (e, st) => SizedBox(height: 150, child: Center(child: Text('로컬 DB 로드 실패', style: TextStyle(color: AppTheme.dangerRed)))),
    );
  }

  Widget _buildPremiumContextCard(ContextCard card) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
      decoration: BoxDecoration(
        color: AppTheme.surfaceDark,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.white.withValues(alpha: 0.05)),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(20),
          onTap: () {},
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: card.type == 'nudge' 
                        ? AppTheme.warningOrange.withValues(alpha: 0.1) 
                        : AppTheme.primaryNeonTeal.withValues(alpha: 0.1),
                    shape: BoxShape.circle,
                  ),
                  child: Icon(
                    card.type == 'nudge' ? Icons.nights_stay : Icons.shopping_bag_outlined,
                    color: card.type == 'nudge' ? AppTheme.warningOrange : AppTheme.primaryNeonTeal,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        card.type == 'nudge' ? '수면 넛지' : 'AI 보충제 추천',
                        style: TextStyle(color: Colors.white.withValues(alpha: 0.5), fontSize: 12, fontWeight: FontWeight.w600),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        card.content,
                        style: const TextStyle(color: Colors.white, fontSize: 15, height: 1.4),
                      ),
                    ],
                  ),
                ),
                const Icon(Icons.chevron_right, color: Colors.white54),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
