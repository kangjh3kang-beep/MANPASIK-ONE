import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 주문 상세 화면
class OrderDetailScreen extends ConsumerWidget {
  const OrderDetailScreen({super.key, required this.orderId});

  final String orderId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('주문 상세'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 주문 상태
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(color: Colors.blue.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
                        child: const Text('배송 중', style: TextStyle(color: Colors.blue, fontSize: 12, fontWeight: FontWeight.w600)),
                      ),
                      const Spacer(),
                      Text('주문번호: $orderId', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // 배송 타임라인
                  _buildTimeline(theme, [
                    _TimelineItem('주문 접수', '2024-02-13 10:00', true),
                    _TimelineItem('결제 완료', '2024-02-13 10:01', true),
                    _TimelineItem('상품 준비', '2024-02-13 14:00', true),
                    _TimelineItem('배송 시작', '2024-02-14 09:30', true),
                    _TimelineItem('배송 완료', '예상: 2024-02-15', false),
                  ]),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),

          // 주문 상품
          Text('주문 상품', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Column(
              children: [
                ListTile(
                  leading: Container(
                    width: 48, height: 48,
                    decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerHighest, borderRadius: BorderRadius.circular(8)),
                    child: const Icon(Icons.science),
                  ),
                  title: const Text('바이오마커 카트리지 PRO'),
                  subtitle: const Text('5개입 패키지'),
                  trailing: const Text('49,000원'),
                ),
                const Divider(height: 1),
                ListTile(
                  leading: Container(
                    width: 48, height: 48,
                    decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerHighest, borderRadius: BorderRadius.circular(8)),
                    child: const Icon(Icons.bluetooth),
                  ),
                  title: const Text('리더기 보호 케이스'),
                  subtitle: const Text('프리미엄 가죽'),
                  trailing: const Text('25,000원'),
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),

          // 결제 정보
          Text('결제 정보', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  _paymentRow(theme, '상품 합계', '74,000원'),
                  _paymentRow(theme, '배송비', '무료'),
                  _paymentRow(theme, '할인', '-5,000원'),
                  const Divider(),
                  _paymentRow(theme, '총 결제금액', '69,000원', isBold: true),
                  const SizedBox(height: 8),
                  _paymentRow(theme, '결제수단', '신용카드 (****1234)'),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),

          // 배송지 정보
          Text('배송지 정보', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('홍길동', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                  Text('010-1234-5678', style: theme.textTheme.bodySmall),
                  Text('서울특별시 강남구 테헤란로 123 만파식빌딩 5F', style: theme.textTheme.bodySmall),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTimeline(ThemeData theme, List<_TimelineItem> items) {
    return Column(
      children: items.asMap().entries.map((entry) {
        final i = entry.key;
        final item = entry.value;
        final isLast = i == items.length - 1;

        return Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Column(
              children: [
                Container(
                  width: 12, height: 12,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: item.completed ? AppTheme.sanggamGold : theme.colorScheme.outlineVariant,
                  ),
                ),
                if (!isLast) Container(width: 2, height: 28, color: item.completed ? AppTheme.sanggamGold : theme.colorScheme.outlineVariant),
              ],
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.only(bottom: 16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(item.label, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: item.completed ? FontWeight.w600 : null)),
                    Text(item.time, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ),
            ),
          ],
        );
      }).toList(),
    );
  }

  Widget _paymentRow(ThemeData theme, String label, String value, {bool isBold = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: isBold ? theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold) : theme.textTheme.bodySmall),
          Text(value, style: isBold ? theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.bold, color: AppTheme.sanggamGold) : theme.textTheme.bodyMedium),
        ],
      ),
    );
  }
}

class _TimelineItem {
  final String label, time;
  final bool completed;
  const _TimelineItem(this.label, this.time, this.completed);
}
