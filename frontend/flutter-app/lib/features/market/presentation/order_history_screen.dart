import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';

/// 주문 내역 화면
class OrderHistoryScreen extends ConsumerWidget {
  const OrderHistoryScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final ordersAsync = ref.watch(ordersProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('주문 내역'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(ordersProvider),
        child: ordersAsync.when(
          data: (orders) {
            if (orders.isEmpty) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.receipt_long_outlined, size: 64, color: theme.colorScheme.onSurfaceVariant),
                    const SizedBox(height: 16),
                    Text('주문 내역이 없습니다.', style: theme.textTheme.bodyLarge),
                    const SizedBox(height: 8),
                    FilledButton(
                      onPressed: () => context.go('/market'),
                      child: const Text('쇼핑하러 가기'),
                    ),
                  ],
                ),
              );
            }
            return ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: orders.length,
              itemBuilder: (context, index) => _buildOrderCard(context, theme, orders[index]),
            );
          },
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (_, __) => Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.error_outline, size: 48),
                const SizedBox(height: 8),
                const Text('주문 내역을 불러올 수 없습니다.'),
                const SizedBox(height: 8),
                FilledButton(
                  onPressed: () => ref.invalidate(ordersProvider),
                  child: const Text('다시 시도'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildOrderCard(BuildContext context, ThemeData theme, Order order) {
    final statusColor = _statusColor(order.status);
    final statusText = _statusText(order.status);

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ExpansionTile(
        title: Row(
          children: [
            Expanded(
              child: Text(
                '주문 #${order.id.length > 8 ? order.id.substring(0, 8) : order.id}',
                style: const TextStyle(fontWeight: FontWeight.w600),
              ),
            ),
            Chip(
              label: Text(statusText, style: TextStyle(fontSize: 11, color: statusColor)),
              backgroundColor: statusColor.withOpacity(0.1),
              side: BorderSide.none,
              visualDensity: VisualDensity.compact,
            ),
          ],
        ),
        subtitle: Text(
          '${order.orderedAt.year}-${order.orderedAt.month.toString().padLeft(2, '0')}-${order.orderedAt.day.toString().padLeft(2, '0')}',
          style: theme.textTheme.bodySmall,
        ),
        children: [
          ...order.items.map((item) => ListTile(
            dense: true,
            leading: const Icon(Icons.science, size: 20, color: AppTheme.sanggamGold),
            title: Text(item.productName, style: theme.textTheme.bodySmall),
            trailing: Text('${item.quantity}개  ₩${_formatPrice(item.unitPrice * item.quantity)}', style: theme.textTheme.bodySmall),
          )),
          const Divider(indent: 16, endIndent: 16),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('합계', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
                Text(
                  '₩${_formatPrice(order.totalAmount)}',
                  style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.bold, color: AppTheme.sanggamGold),
                ),
              ],
            ),
          ),
          if (order.status == OrderStatus.confirmed || order.status == OrderStatus.shipping)
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
              child: SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: () => context.push('/market/order-complete/${order.id}'),
                  icon: const Icon(Icons.local_shipping_outlined, size: 18),
                  label: const Text('배송 추적'),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Color _statusColor(OrderStatus status) {
    switch (status) {
      case OrderStatus.pending: return Colors.orange;
      case OrderStatus.confirmed: return Colors.blue;
      case OrderStatus.shipping: return Colors.indigo;
      case OrderStatus.delivered: return Colors.green;
      case OrderStatus.cancelled: return Colors.red;
    }
  }

  String _statusText(OrderStatus status) {
    switch (status) {
      case OrderStatus.pending: return '결제 대기';
      case OrderStatus.confirmed: return '결제 완료';
      case OrderStatus.shipping: return '배송 중';
      case OrderStatus.delivered: return '배송 완료';
      case OrderStatus.cancelled: return '취소됨';
    }
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
      RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]},',
    );
  }
}
