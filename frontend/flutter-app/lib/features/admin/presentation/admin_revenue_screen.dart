import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminRevenueScreen — Sanggam Orbit 매출 통계
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] +sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] ThemeData 파라미터 1x 제거
// [Rule 4] theme.textTheme.* ~4x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~3x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16, borderRadius 4→8
// ───────────────────────────────────────────────────

/// 관리자 매출 통계 화면 (C12)
class AdminRevenueScreen extends ConsumerWidget {
  const AdminRevenueScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final statsAsync =
        ref.watch(revenueStatsProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '매출 통계',
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
              child: statsAsync.when(
                data: (stats) =>
                    _buildContent(stats),
                loading: () => const Center(
                    child:
                        CircularProgressIndicator(
                            color: SanggamTheme
                                .primary)),
                error: (e, _) => Center(
                    child: Text('오류: $e',
                        style: const TextStyle(
                            color:
                                Colors.white))),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildContent(
      Map<String, dynamic> stats) {
    final periods =
        stats['periods'] as List? ?? [];
    final totalRevenue =
        stats['total_revenue'] as num? ?? 0;
    final subscriptionRevenue =
        stats['subscription_revenue'] as num? ??
            0;
    final cartridgeRevenue =
        stats['cartridge_revenue'] as num? ?? 0;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.stretch,
        children: [
          // 매출 요약 카드
          Row(
            children: [
              Expanded(
                child: _StatCard(
                  title: '총 매출',
                  value:
                      '${_formatKrw(totalRevenue)}원',
                  icon: Icons
                      .account_balance_wallet_rounded,
                  color: SanggamTheme.primary,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: _StatCard(
                  title: '구독',
                  value:
                      '${_formatKrw(subscriptionRevenue)}원',
                  icon: Icons
                      .card_membership_rounded,
                  color: SanggamTheme.jagaeCyan,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          _StatCard(
            title: '카트리지 판매',
            value:
                '${_formatKrw(cartridgeRevenue)}원',
            icon: Icons.science_rounded,
            color: SanggamTheme.jagaeMagenta,
          ),
          const SizedBox(height: 24),

          // 월별 차트
          const Text(
            '월별 매출 추이',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 220,
            child: periods.isEmpty
                ? const Center(
                    child: Text('데이터 없음',
                        style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim)))
                : BarChart(
                    BarChartData(
                      alignment:
                          BarChartAlignment
                              .spaceAround,
                      barTouchData:
                          BarTouchData(
                              enabled: true),
                      titlesData: FlTitlesData(
                        bottomTitles: AxisTitles(
                          sideTitles: SideTitles(
                            showTitles: true,
                            getTitlesWidget:
                                (value, meta) {
                              final idx =
                                  value.toInt();
                              if (idx >= 0 &&
                                  idx <
                                      periods
                                          .length) {
                                final label =
                                    periods[idx][
                                            'label']
                                        as String? ??
                                        '';
                                return Text(
                                    label,
                                    style:
                                        const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize: 9,
                                    ));
                              }
                              return const SizedBox
                                  .shrink();
                            },
                          ),
                        ),
                        leftTitles:
                            const AxisTitles(
                                sideTitles:
                                    SideTitles(
                                        showTitles:
                                            false)),
                        topTitles:
                            const AxisTitles(
                                sideTitles:
                                    SideTitles(
                                        showTitles:
                                            false)),
                        rightTitles:
                            const AxisTitles(
                                sideTitles:
                                    SideTitles(
                                        showTitles:
                                            false)),
                      ),
                      borderData: FlBorderData(
                          show: false),
                      gridData: FlGridData(
                          show: false),
                      barGroups: periods
                          .asMap()
                          .entries
                          .map((e) {
                        final amount = (e.value[
                                    'amount']
                                as num?)
                            ?.toDouble() ?? 0;
                        return BarChartGroupData(
                          x: e.key,
                          barRods: [
                            BarChartRodData(
                              toY: amount,
                              color:
                                  SanggamTheme
                                      .primary,
                              width: 16,
                              borderRadius:
                                  const BorderRadius
                                      .vertical(
                                      top: Radius
                                          .circular(
                                              8)),
                            ),
                          ],
                        );
                      }).toList(),
                    ),
                    duration: const Duration(
                        milliseconds: 250),
                  ),
          ),
        ],
      ),
    );
  }

  String _formatKrw(num value) {
    if (value >= 100000000) {
      return '${(value / 100000000).toStringAsFixed(1)}억';
    }
    if (value >= 10000) {
      return '${(value / 10000).toStringAsFixed(0)}만';
    }
    return value.toString();
  }
}

class _StatCard extends StatelessWidget {
  const _StatCard({
    required this.title,
    required this.value,
    required this.icon,
    required this.color,
  });
  final String title;
  final String value;
  final IconData icon;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(height: 8),
          Text(title,
              style: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim,
                fontSize: 12,
              )),
          const SizedBox(height: 4),
          Text(value,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              )),
        ],
      ),
    );
  }
}

/// 매출 통계 Provider
final revenueStatsProvider =
    FutureProvider<Map<String, dynamic>>(
        (ref) async {
  try {
    return await ref
        .read(restClientProvider)
        .getRevenueStats();
  } catch (_) {
    return {
      'total_revenue': 0,
      'subscription_revenue': 0,
      'cartridge_revenue': 0,
      'periods': [],
    };
  }
});
