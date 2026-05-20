import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

/// 관리자 카트리지 재고 테이블 (C12)
class AdminInventoryTable extends ConsumerWidget {
  const AdminInventoryTable({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final inventoryAsync = ref.watch(inventoryStatsProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          tooltip: '뒤로 가기',
          onPressed: () => context.pop(),
        ),
        title: const Text('카트리지 재고', style: TextStyle(
          color: SanggamTheme.primary,
          fontWeight: FontWeight.bold,
        )),
      ),
      body: inventoryAsync.when(
        data: (data) => _buildTable(data),
        loading: () => const Center(child: CircularProgressIndicator(color: SanggamTheme.primary)),
        error: (e, _) => Center(child: Text('오류: $e', style: const TextStyle(color: SanggamTheme.error))),
      ),
    );
  }

  Widget _buildTable(Map<String, dynamic> data) {
    final items = (data['items'] as List?)
            ?.map((e) => e as Map<String, dynamic>)
            .toList() ??
        [];

    if (items.isEmpty) {
      return const Center(child: Text('재고 데이터가 없습니다', style: TextStyle(color: SanggamTheme.onSurfaceDim)));
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 요약
          Row(
            children: [
              const _SummaryChip(
                label: '총 품목',
                value: '-',
                color: SanggamTheme.primary,
              ),
              const SizedBox(width: 8),
              _SummaryChip(
                label: '재고 부족',
                value:
                    '${items.where((i) => (i['quantity'] as num? ?? 0) < (i['reorder_level'] as num? ?? 10)).length}건',
                color: SanggamTheme.error,
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 테이블
          SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: DataTable(
              headingRowColor: WidgetStateProperty.all(
                SanggamTheme.surfaceVariant.withValues(alpha: 0.5),
              ),
              columns: const [
                DataColumn(label: Text('카트리지', style: TextStyle(color: Colors.white))),
                DataColumn(label: Text('카테고리', style: TextStyle(color: Colors.white))),
                DataColumn(label: Text('재고', style: TextStyle(color: Colors.white)), numeric: true),
                DataColumn(label: Text('재주문 기준', style: TextStyle(color: Colors.white)), numeric: true),
                DataColumn(label: Text('상태', style: TextStyle(color: Colors.white))),
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
                          SanggamTheme.error.withValues(alpha: 0.05))
                      : null,
                  cells: [
                    DataCell(Text(name, style: const TextStyle(color: Colors.white))),
                    DataCell(Text(category, style: const TextStyle(color: Colors.white))),
                    DataCell(Text('$qty', style: const TextStyle(color: Colors.white))),
                    DataCell(Text('$reorderLevel', style: const TextStyle(color: Colors.white))),
                    DataCell(
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: isLow
                              ? SanggamTheme.error.withValues(alpha: 0.1)
                              : SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          isLow ? '부족' : '정상',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.bold,
                            color:
                                isLow ? SanggamTheme.error : SanggamTheme.jagaeCyan,
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
        backgroundColor: color.withValues(alpha: 0.2),
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
