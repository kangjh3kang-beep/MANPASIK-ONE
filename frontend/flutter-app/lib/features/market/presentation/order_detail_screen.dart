import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// 주문 상세 화면
class OrderDetailScreen extends ConsumerStatefulWidget {
  const OrderDetailScreen({super.key, required this.orderId});

  final String orderId;

  @override
  ConsumerState<OrderDetailScreen> createState() => _OrderDetailScreenState();
}

class _OrderDetailScreenState extends ConsumerState<OrderDetailScreen> {
  Map<String, dynamic>? _order;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadOrder();
  }

  Future<void> _loadOrder() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.getOrderDetail(widget.orderId);
      setState(() {
        _order = data;
        _isLoading = false;
      });
    } on DioException {
      setState(() {
        _isLoading = false;
        _error = '주문 정보를 불러올 수 없습니다';
      });
    }
  }

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
                      '주문 상세',
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
                              Text(_error!,
                                  style: const TextStyle(
                                      color: Colors.white)),
                              const SizedBox(height: 16),
                              FilledButton(
                                onPressed: _loadOrder,
                                child: const Text('다시 시도'),
                              ),
                            ],
                          ),
                        )
                      : _buildOrderContent(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderContent() {
    final status = (_order?['status'] ?? 'pending') as String;
    final items = _order?['items'] as List? ?? [];
    final totalAmount =
        (_order?['total_amount'] as num?)?.toInt() ?? 0;
    final shippingAddress =
        (_order?['shipping_address'] ?? '') as String;
    final paymentMethod =
        (_order?['payment_method'] ?? '') as String;
    final createdAt = (_order?['created_at'] ?? '') as String;

    final statusLabel = _statusLabel(status);
    final statusColor = _statusColor(status);

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        // 주문 상태
        SanggamContainer(
          borderRadius: 16,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: statusColor.withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Text(statusLabel,
                        style: TextStyle(
                          color: statusColor,
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                        )),
                  ),
                  const Spacer(),
                  Text('주문번호: ${widget.orderId}',
                      style: const TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 13,
                      )),
                ],
              ),
              const SizedBox(height: 16),
              _buildTimeline(_timelineFromStatus(status, createdAt)),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // 주문 상품
        const Text('주문 상품',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            )),
        const SizedBox(height: 8),
        SanggamContainer(
          borderRadius: 16,
          padding: EdgeInsets.zero,
          child: Column(
            children: [
              for (var i = 0; i < items.length; i++) ...[
                if (i > 0)
                  const Divider(
                      height: 1, color: SanggamTheme.surfaceVariant),
                _buildProductTile(items[i] as Map<String, dynamic>),
              ],
              if (items.isEmpty)
                const Padding(
                  padding: EdgeInsets.all(16),
                  child: Text('상품 정보 없음',
                      style: TextStyle(color: SanggamTheme.onSurfaceDim)),
                ),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // 결제 정보
        const Text('결제 정보',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            )),
        const SizedBox(height: 8),
        SanggamContainer(
          borderRadius: 16,
          child: Column(
            children: [
              _paymentRow('총 결제금액', '₩${_formatPrice(totalAmount)}',
                  isBold: true),
              if (paymentMethod.isNotEmpty)
                _paymentRow('결제수단', paymentMethod),
            ],
          ),
        ),
        const SizedBox(height: 16),

        // 배송지 정보
        if (shippingAddress.isNotEmpty) ...[
          const Text('배송지 정보',
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              )),
          const SizedBox(height: 8),
          SanggamContainer(
            borderRadius: 16,
            child: Text(shippingAddress,
                style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim, fontSize: 13)),
          ),
        ],
      ],
    );
  }

  Widget _buildProductTile(Map<String, dynamic> item) {
    final name = (item['name'] ?? item['product_name'] ?? '') as String;
    final subtitle = (item['subtitle'] ?? '') as String;
    final price = (item['price'] as num?)?.toInt() ?? 0;
    return ListTile(
      leading: Container(
        width: 48,
        height: 48,
        decoration: BoxDecoration(
          color: SanggamTheme.surfaceVariant,
          borderRadius: BorderRadius.circular(8),
        ),
        child: const Icon(Icons.science, color: SanggamTheme.primary),
      ),
      title: Text(name, style: const TextStyle(color: Colors.white)),
      subtitle: subtitle.isNotEmpty
          ? Text(subtitle,
              style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim, fontSize: 12))
          : null,
      trailing: Text('₩${_formatPrice(price)}',
          style: const TextStyle(color: Colors.white)),
    );
  }

  List<_TimelineItem> _timelineFromStatus(String status, String createdAt) {
    final steps = ['pending', 'paid', 'preparing', 'shipping', 'delivered'];
    final labels = ['주문 접수', '결제 완료', '상품 준비', '배송 시작', '배송 완료'];
    final currentIdx = steps.indexOf(status);
    return List.generate(steps.length, (i) {
      return _TimelineItem(
        labels[i],
        i == 0 && createdAt.isNotEmpty
            ? createdAt.substring(0, createdAt.length.clamp(0, 16))
            : i <= currentIdx
                ? '완료'
                : '대기 중',
        i <= currentIdx,
      );
    });
  }

  Widget _buildTimeline(List<_TimelineItem> items) {
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
                  width: 12,
                  height: 12,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: item.completed
                        ? SanggamTheme.primary
                        : SanggamTheme.surfaceVariant,
                  ),
                ),
                if (!isLast)
                  Container(
                      width: 2,
                      height: 28,
                      color: item.completed
                          ? SanggamTheme.primary
                          : SanggamTheme.surfaceVariant),
              ],
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.only(bottom: 16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(item.label,
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: item.completed
                              ? FontWeight.w600
                              : FontWeight.normal,
                        )),
                    Text(item.time,
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        )),
                  ],
                ),
              ),
            ),
          ],
        );
      }).toList(),
    );
  }

  Widget _paymentRow(String label, String value, {bool isBold = false}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: TextStyle(
                color: isBold ? Colors.white : SanggamTheme.onSurfaceDim,
                fontSize: isBold ? 14 : 13,
                fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
              )),
          Text(value,
              style: TextStyle(
                color: isBold ? SanggamTheme.primary : Colors.white,
                fontSize: isBold ? 16 : 14,
                fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
              )),
        ],
      ),
    );
  }

  String _statusLabel(String status) {
    const labels = {
      'pending': '주문 접수',
      'paid': '결제 완료',
      'preparing': '상품 준비',
      'shipping': '배송 중',
      'delivered': '배송 완료',
      'cancelled': '주문 취소',
    };
    return labels[status] ?? status;
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'delivered':
        return SanggamTheme.jagaeCyan;
      case 'cancelled':
        return SanggamTheme.error;
      case 'shipping':
        return SanggamTheme.primary;
      default:
        return SanggamTheme.onSurfaceDim;
    }
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (m) => '${m[1]},',
        );
  }
}

class _TimelineItem {
  final String label, time;
  final bool completed;
  const _TimelineItem(this.label, this.time, this.completed);
}
