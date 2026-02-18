import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 관리자 카트리지 재고 테이블 (C12)
class AdminInventoryTable extends ConsumerWidget {
  const AdminInventoryTable({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final inventoryAsync = ref.watch(inventoryStatsProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('카트리지 재고'),
      ),
      body: inventoryAsync.when(
        data: (data) => _buildTable(theme, data),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('오류: $e')),
      ),
    );
  }

  Widget _buildTable(ThemeData theme, Map<String, dynamic> data) {
    final items = (data['items'] as List?)
            ?.map((e) => e as Map<String, dynamic>)
            .toList() ??
        [];

    if (items.isEmpty) {
      return const Center(child: Text('재고 데이터가 없습니다'));
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 요약
          Row(
            children: [
              _SummaryChip(
                label: '총 품목',
                value: '${items.length}종',
                color: theme.colorScheme.primary,
              ),
              const SizedBox(width: 8),
              _SummaryChip(
                label: '재고 부족',
                value:
                    '${items.where((i) => (i['quantity'] as num? ?? 0) < (i['reorder_level'] as num? ?? 10)).length}건',
                color: AppTheme.dancheongRed,
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 테이블
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: DataTable(
              headingRowColor: WidgetStateProperty.all(
                theme.colorScheme.surfaceContainerHighest.withOpacity(0.5),
              ),
              columns: const [
                DataColumn(label: Text('카트리지')),
                DataColumn(label: Text('카테고리')),
                DataColumn(label: Text('재고'), numeric: true),
                DataColumn(label: Text('재주문 기준'), numeric: true),
                DataColumn(label: Text('상태')),
              ],
              rows: items.map((item) {
                final name = item['name'] as String? ?? '-';
                final category = item['category'] as String? ?? '-';
                final qty = item['quantity'] as num? ?? 0;
                final reorderLevel = item['reorder_level'] as num? ?? 10;
                final isLow = qty < reorderLevel;

                return DataRow(
                  color: isLow
                      ? WidgetStateProperty.all(
                          AppTheme.dancheongRed.withOpacity(0.05))
                      : null,
                  cells: [
                    DataCell(Text(name)),
                    DataCell(Text(category)),
                    DataCell(Text('$qty')),
                    DataCell(Text('$reorderLevel')),
                    DataCell(
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: isLow
                              ? AppTheme.dancheongRed.withOpacity(0.1)
                              : Colors.green.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          isLow ? '부족' : '정상',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.bold,
                            color:
                                isLow ? AppTheme.dancheongRed : Colors.green,
                          ),
                        ),
                      ),
                    ),
                  ],
                );
              }).toList(),
            ),
          ),
        ],
      ),
    );
  }
}

class _SummaryChip extends StatelessWidget {
  const _SummaryChip({
    required this.label,
    required this.value,
    required this.color,
  });
  final String label;
  final String value;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Chip(
      avatar: CircleAvatar(
        backgroundColor: color.withOpacity(0.2),
        child: Text(
          value,
          style: TextStyle(fontSize: 10, color: color, fontWeight: FontWeight.bold),
        ),
      ),
      label: Text(label),
    );
  }
}

/// 재고 통계 Provider
final inventoryStatsProvider =
    FutureProvider<Map<String, dynamic>>((ref) async {
  try {
    return await ref.read(restClientProvider).getInventoryStats();
  } catch (_) {
    return {'items': []};
  }
});
