import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 장바구니 화면
class CartScreen extends ConsumerStatefulWidget {
  const CartScreen({super.key});

  @override
  ConsumerState<CartScreen> createState() => _CartScreenState();
}

class _CartScreenState extends ConsumerState<CartScreen> {
  // 로컬 장바구니 (서버 미연결 시 fallback)
  final List<_CartItem> _items = [
    _CartItem(id: 'BIO-001', name: '혈당 측정 카트리지 (10개입)', price: 29900, quantity: 1),
    _CartItem(id: 'BIO-002', name: '콜레스테롤 카트리지 (5개입)', price: 39900, quantity: 1),
  ];

  int get _totalPrice => _items.fold(0, (sum, item) => sum + item.price * item.quantity);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('장바구니'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: _items.isEmpty
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.shopping_cart_outlined, size: 64, color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(height: 16),
                  Text('장바구니가 비어있습니다.', style: theme.textTheme.bodyLarge),
                  const SizedBox(height: 8),
                  FilledButton(
                    onPressed: () => context.go('/market'),
                    child: const Text('쇼핑하러 가기'),
                  ),
                ],
              ),
            )
          : Column(
              children: [
                Expanded(
                  child: ListView.separated(
                    padding: const EdgeInsets.all(16),
                    itemCount: _items.length,
                    separatorBuilder: (_, __) => const Divider(),
                    itemBuilder: (context, index) => _buildCartItemTile(theme, index),
                  ),
                ),
                // 결제 영역
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.surface,
                    boxShadow: [BoxShadow(color: Colors.black12, blurRadius: 4, offset: const Offset(0, -2))],
                  ),
                  child: SafeArea(
                    child: Column(
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text('총 ${_items.length}개 상품', style: theme.textTheme.bodyMedium),
                            Text(
                              '₩${_formatPrice(_totalPrice)}',
                              style: theme.textTheme.titleLarge?.copyWith(
                                fontWeight: FontWeight.bold,
                                color: AppTheme.sanggamGold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),
                        SizedBox(
                          width: double.infinity,
                          child: FilledButton(
                            onPressed: () => context.push('/market/checkout'),
                            style: FilledButton.styleFrom(
                              minimumSize: const Size.fromHeight(48),
                              backgroundColor: AppTheme.sanggamGold,
                            ),
                            child: Text('₩${_formatPrice(_totalPrice)} 결제하기'),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
    );
  }

  Widget _buildCartItemTile(ThemeData theme, int index) {
    final item = _items[index];
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // 상품 아이콘
        Container(
          width: 64, height: 64,
          decoration: BoxDecoration(
            color: theme.colorScheme.surfaceContainerHighest,
            borderRadius: BorderRadius.circular(8),
          ),
          child: const Icon(Icons.science, color: AppTheme.sanggamGold),
        ),
        const SizedBox(width: 12),

        // 상품 정보
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(item.name, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
              const SizedBox(height: 4),
              Text('₩${_formatPrice(item.price)}', style: theme.textTheme.bodySmall),
              const SizedBox(height: 8),
              // 수량 조절
              Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.remove_circle_outline, size: 20),
                    onPressed: () {
                      setState(() {
                        if (item.quantity > 1) {
                          _items[index] = item.copyWith(quantity: item.quantity - 1);
                        }
                      });
                    },
                  ),
                  Text('${item.quantity}', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
                  IconButton(
                    icon: const Icon(Icons.add_circle_outline, size: 20),
                    onPressed: () {
                      setState(() {
                        _items[index] = item.copyWith(quantity: item.quantity + 1);
                      });
                    },
                  ),
                  const Spacer(),
                  IconButton(
                    icon: const Icon(Icons.delete_outline, size: 20, color: Colors.red),
                    onPressed: () => setState(() => _items.removeAt(index)),
                  ),
                ],
              ),
              // 정기 구독 옵션
              Row(
                children: [
                  Checkbox(
                    value: item.subscribe,
                    onChanged: (v) {
                      setState(() {
                        _items[index] = item.copyWith(subscribe: v ?? false);
                      });
                    },
                  ),
                  Text(
                    '정기 구독 (10% 할인)',
                    style: TextStyle(fontSize: 12, color: item.subscribe ? Colors.green : Colors.grey),
                  ),
                ],
              ),
            ],
          ),
        ),
      ],
    );
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
      RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]},',
    );
  }
}

class _CartItem {
  final String id, name;
  final int price, quantity;
  final bool subscribe;
  const _CartItem({required this.id, required this.name, required this.price, required this.quantity, this.subscribe = false});
  _CartItem copyWith({int? quantity, bool? subscribe}) =>
      _CartItem(id: id, name: name, price: price, quantity: quantity ?? this.quantity, subscribe: subscribe ?? this.subscribe);
}
