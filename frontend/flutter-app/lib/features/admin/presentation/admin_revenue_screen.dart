import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 관리자 매출 통계 화면 (C12)
class AdminRevenueScreen extends ConsumerWidget {
  const AdminRevenueScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final statsAsync = ref.watch(revenueStatsProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('매출 통계'),
      ),
      body: statsAsync.when(
        data: (stats) => _buildContent(theme, stats),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('오류: $e')),
      ),
    );
  }

  Widget _buildContent(ThemeData theme, Map<String, dynamic> stats) {
    final periods = stats['periods'] as List? ?? [];
    final totalRevenue = stats['total_revenue'] as num? ?? 0;
    final subscriptionRevenue = stats['subscription_revenue'] as num? ?? 0;
    final cartridgeRevenue = stats['cartridge_revenue'] as num? ?? 0;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 매출 요약 카드
          Row(
            children: [
              Expanded(
                child: _StatCard(
                  title: '총 매출',
                  value: '${_formatKrw(totalRevenue)}원',
                  icon: Icons.account_balance_wallet_rounded,
                  color: AppTheme.sanggamGold,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _StatCard(
                  title: '구독',
                  value: '${_formatKrw(subscriptionRevenue)}원',
                  icon: Icons.card_membership_rounded,
                  color: theme.colorScheme.primary,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _StatCard(
            title: '카트리지 판매',
            value: '${_formatKrw(cartridgeRevenue)}원',
            icon: Icons.science_rounded,
            color: theme.colorScheme.tertiary,
          ),
          const SizedBox(height: 24),

          // 월별 차트
          Text(
            '월별 매출 추이',
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 220,
            child: periods.isEmpty
                ? const Center(child: Text('데이터 없음'))
                : BarChart(
                    BarChartData(
                      alignment: BarChartAlignment.spaceAround,
                      barTouchData: BarTouchData(enabled: true),
                      titlesData: FlTitlesData(
                        bottomTitles: AxisTitles(
                          sideTitles: SideTitles(
                            showTitles: true,
                            getTitlesWidget: (value, meta) {
                              final idx = value.toInt();
                              if (idx >= 0 && idx < periods.length) {
                                final label =
                                    periods[idx]['label'] as String? ?? '';
                                return Text(label,
                                    style: theme.textTheme.bodySmall
                                        ?.copyWith(fontSize: 9));
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
                      gridData: FlGridData(show: false),
                      barGroups: periods.asMap().entries.map((e) {
                        final amount =
                            (e.value['amount'] as num?)?.toDouble() ?? 0;
                        return BarChartGroupData(
                          x: e.key,
                          barRods: [
                            BarChartRodData(
                              toY: amount,
                              color: AppTheme.sanggamGold,
                              width: 16,
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
        ],
      ),
    );
  }

  String _formatKrw(num value) {
    if (value >= 100000000) return '${(value / 100000000).toStringAsFixed(1)}억';
    if (value >= 10000) return '${(value / 10000).toStringAsFixed(0)}만';
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
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Icon(icon, color: color, size: 24),
            const SizedBox(height: 8),
            Text(title,
                style: theme.textTheme.bodySmall
                    ?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 4),
            Text(value,
                style: theme.textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
          ],
        ),
      ),
    );
  }
}

/// 매출 통계 Provider
final revenueStatsProvider =
    FutureProvider<Map<String, dynamic>>((ref) async {
  try {
    return await ref.read(restClientProvider).getRevenueStats();
  } catch (_) {
    return {
      'total_revenue': 0,
      'subscription_revenue': 0,
      'cartridge_revenue': 0,
      'periods': [],
    };
  }
});
