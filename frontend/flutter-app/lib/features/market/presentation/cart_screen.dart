import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// 장바구니 화면
class CartScreen extends ConsumerStatefulWidget {
  const CartScreen({super.key});

  @override
  ConsumerState<CartScreen> createState() => _CartScreenState();
}

class _CartScreenState extends ConsumerState<CartScreen> {
  List<_CartItem> _items = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadCart();
  }

  Future<void> _loadCart() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final data = await rest.getCart(userId);
      final rawItems = data['items'] as List? ?? [];
      setState(() {
        _items = rawItems
            .map((e) => _CartItem.fromMap(e as Map<String, dynamic>))
            .toList();
        _isLoading = false;
      });
    } on DioException {
      setState(() {
        _isLoading = false;
        _error = '장바구니를 불러올 수 없습니다';
      });
    }
  }

  Future<void> _removeItem(int index) async {
    final item = _items[index];
    final prev = List<_CartItem>.from(_items);
    setState(() => _items.removeAt(index));
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.removeFromCart(userId: userId, itemId: item.id);
    } on DioException {
      setState(() => _items = prev);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('삭제에 실패했습니다')),
        );
      }
    }
  }

  Future<void> _updateQuantity(int index, int newQty) async {
    if (newQty < 1) return;
    final item = _items[index];
    final oldQty = item.quantity;
    setState(() {
      _items[index] = item.copyWith(quantity: newQty);
    });
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.updateCartItemQuantity(
        userId: userId,
        itemId: item.id,
        quantity: newQty,
      );
    } on DioException {
      setState(() {
        _items[index] = item.copyWith(quantity: oldQty);
      });
    }
  }

  int get _totalPrice =>
      _items.fold(0, (sum, item) => sum + item.price * item.quantity);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
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
                      '장바구니',
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
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : _error != null
                      ? Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              const Icon(Icons.error_outline,
                                  size: 48,
                                  color: SanggamTheme.onSurfaceDim),
                              const SizedBox(height: 8),
                              Text(_error!,
                                  style: const TextStyle(
                                      color: Colors.white, fontSize: 14)),
                              const SizedBox(height: 16),
                              FilledButton(
                                onPressed: _loadCart,
                                style: FilledButton.styleFrom(
                                  backgroundColor: SanggamTheme.primary,
                                  foregroundColor: SanggamTheme.background,
                                ),
                                child: const Text('다시 시도'),
                              ),
                            ],
                          ),
                        )
                      : _items.isEmpty
                          ? Center(
                              child: Column(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  const Icon(Icons.shopping_cart_outlined,
                                      size: 64,
                                      color: SanggamTheme.onSurfaceDim),
                                  const SizedBox(height: 16),
                                  const Text('장바구니가 비어있습니다.',
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 16)),
                                  const SizedBox(height: 8),
                                  FilledButton(
                                    onPressed: () => context.go('/market'),
                                    style: FilledButton.styleFrom(
                                      backgroundColor: SanggamTheme.primary,
                                      foregroundColor:
                                          SanggamTheme.background,
                                      shape: RoundedRectangleBorder(
                                        borderRadius:
                                            BorderRadius.circular(16),
                                      ),
                                    ),
                                    child: const Text('쇼핑하러 가기'),
                                  ),
                                ],
                              ),
                            )
                          : ListView.separated(
                              padding: const EdgeInsets.all(16),
                              itemCount: _items.length,
                              separatorBuilder: (_, __) => const Divider(
                                  color: SanggamTheme.surfaceVariant),
                              itemBuilder: (context, index) =>
                                  _buildCartItemTile(index),
                            ),
            ),
            // 결제 영역
            if (_items.isNotEmpty)
              SanggamContainer(
                borderRadius: 0,
                child: SafeArea(
                  top: false,
                  child: Column(
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text('총 ${_items.length}개 상품',
                              style: const TextStyle(
                                  color: Colors.white, fontSize: 14)),
                          Text(
                            '₩${_formatPrice(_totalPrice)}',
                            style: const TextStyle(
                              color: SanggamTheme.primary,
                              fontSize: 20,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      SizedBox(
                        width: double.infinity,
                        child: FilledButton(
                          onPressed: () => context.push('/market/checkout'),
                          style: FilledButton.styleFrom(
                            minimumSize: const Size.fromHeight(48),
                            backgroundColor: SanggamTheme.primary,
                            foregroundColor: SanggamTheme.background,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                          ),
                          child: Text(
                              '₩${_formatPrice(_totalPrice)} 결제하기',
                              style: const TextStyle(
                                  fontWeight: FontWeight.bold)),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildCartItemTile(int index) {
    final item = _items[index];
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Container(
          width: 64,
          height: 64,
          decoration: BoxDecoration(
            color: SanggamTheme.surfaceVariant,
            borderRadius: BorderRadius.circular(8),
          ),
          child: const Icon(Icons.science, color: SanggamTheme.primary),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(item.name,
                  style: const TextStyle(
                      color: Colors.white, fontWeight: FontWeight.w600)),
              const SizedBox(height: 4),
              Text('₩${_formatPrice(item.price)}',
                  style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim, fontSize: 13)),
              const SizedBox(height: 8),
              Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.remove_circle_outline,
                        size: 20, color: SanggamTheme.onSurfaceDim),
                    tooltip: '수량 감소',
                    onPressed: () =>
                        _updateQuantity(index, item.quantity - 1),
                  ),
                  Text('${item.quantity}',
                      style: const TextStyle(
                          color: Colors.white, fontWeight: FontWeight.bold)),
                  IconButton(
                    icon: const Icon(Icons.add_circle_outline,
                        size: 20, color: SanggamTheme.onSurfaceDim),
                    tooltip: '수량 증가',
                    onPressed: () =>
                        _updateQuantity(index, item.quantity + 1),
                  ),
                  const Spacer(),
                  IconButton(
                    icon: const Icon(Icons.delete_outline,
                        size: 20, color: SanggamTheme.error),
                    tooltip: '삭제',
                    onPressed: () => _removeItem(index),
                  ),
                ],
              ),
              Row(
                children: [
                  Checkbox(
                    value: item.subscribe,
                    activeColor: SanggamTheme.primary,
                    checkColor: SanggamTheme.background,
                    side: const BorderSide(color: SanggamTheme.onSurfaceDim),
                    onChanged: (v) {
                      setState(() {
                        _items[index] =
                            item.copyWith(subscribe: v ?? false);
                      });
                    },
                  ),
                  Text(
                    '정기 구독 (10% 할인)',
                    style: TextStyle(
                      fontSize: 12,
                      color: item.subscribe
                          ? SanggamTheme.jagaeCyan
                          : SanggamTheme.onSurfaceDim,
                    ),
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
  const _CartItem({
    required this.id,
    required this.name,
    required this.price,
    required this.quantity,
    this.subscribe = false,
  });

  factory _CartItem.fromMap(Map<String, dynamic> m) {
    return _CartItem(
      id: (m['id'] ?? m['product_id'] ?? '') as String,
      name: (m['name'] ?? m['product_name'] ?? '') as String,
      price: (m['price'] as num?)?.toInt() ?? 0,
      quantity: (m['quantity'] as num?)?.toInt() ?? 1,
      subscribe: (m['subscribe'] as bool?) ?? false,
    );
  }

  _CartItem copyWith({int? quantity, bool? subscribe}) => _CartItem(
        id: id,
        name: name,
        price: price,
        quantity: quantity ?? this.quantity,
        subscribe: subscribe ?? this.subscribe,
      );
}
