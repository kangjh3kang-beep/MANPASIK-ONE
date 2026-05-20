import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:lottie/lottie.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// 주문 완료 & 배송 추적 화면
class OrderCompleteScreen extends ConsumerStatefulWidget {
  const OrderCompleteScreen({super.key, required this.orderId});

  final String orderId;

  @override
  ConsumerState<OrderCompleteScreen> createState() =>
      _OrderCompleteScreenState();
}

class _OrderCompleteScreenState extends ConsumerState<OrderCompleteScreen> {
  Map<String, dynamic>? _order;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadOrder();
  }

  Future<void> _loadOrder() async {
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.getOrderDetail(widget.orderId);
      setState(() {
        _order = data;
        _isLoading = false;
      });
    } on DioException {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final totalAmount = (_order?['total_amount'] as num?)?.toInt() ?? 0;
    final shippingAddress =
        (_order?['shipping_address'] ?? '') as String;
    final status = (_order?['status'] ?? 'pending') as String;

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
                  const SizedBox(width: 48),
                  const Expanded(
                    child: Text(
                      '주문 완료',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  const SizedBox(width: 48),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : SingleChildScrollView(
                      padding: const EdgeInsets.all(24),
                      child: Column(
                        children: [
                          const SizedBox(height: 24),
                          SizedBox(
                            width: 120,
                            height: 120,
                            child: Lottie.asset(
                              'assets/lottie/check_success.json',
                              repeat: false,
                              errorBuilder: (_, __, ___) => Container(
                                width: 96,
                                height: 96,
                                decoration: BoxDecoration(
                                  color: SanggamTheme.jagaeCyan
                                      .withValues(alpha: 0.1),
                                  shape: BoxShape.circle,
                                ),
                                child: const Icon(Icons.check_circle,
                                    size: 64,
                                    color: SanggamTheme.jagaeCyan),
                              ),
                            ),
                          ),
                          const SizedBox(height: 24),
                          const Text('주문이 완료되었습니다!',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 22,
                                fontWeight: FontWeight.bold,
                              )),
                          const SizedBox(height: 8),
                          Text('주문번호: ${widget.orderId}',
                              style: const TextStyle(
                                color: SanggamTheme.onSurfaceDim,
                                fontSize: 14,
                              )),
                          if (totalAmount > 0) ...[
                            const SizedBox(height: 4),
                            Text(
                                '결제 금액: ₩${_formatPrice(totalAmount)}',
                                style: const TextStyle(
                                  color: SanggamTheme.primary,
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold,
                                )),
                          ],
                          const SizedBox(height: 32),

                          // 배송 정보
                          SanggamContainer(
                            borderRadius: 16,
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment:
                                  CrossAxisAlignment.start,
                              children: [
                                const Text('배송 정보',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontSize: 14,
                                      fontWeight: FontWeight.bold,
                                    )),
                                const SizedBox(height: 16),
                                _buildInfoRow('주문 상태',
                                    _statusLabel(status)),
                                if (shippingAddress.isNotEmpty)
                                  _buildInfoRow(
                                      '배송지', shippingAddress),
                              ],
                            ),
                          ),
                          const SizedBox(height: 16),

                          // 배송 상태 타임라인
                          SanggamContainer(
                            borderRadius: 16,
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment:
                                  CrossAxisAlignment.start,
                              children: [
                                const Text('배송 추적',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontSize: 14,
                                      fontWeight: FontWeight.bold,
                                    )),
                                const SizedBox(height: 16),
                                _buildTimelineStep('주문 접수', '완료',
                                    true, true),
                                _buildTimelineStep(
                                    '상품 준비',
                                    status == 'paid' ||
                                            status == 'preparing'
                                        ? '처리 중'
                                        : '대기 중',
                                    status != 'pending',
                                    false),
                                _buildTimelineStep('발송', '대기 중',
                                    false, false),
                                _buildTimelineStep(
                                    '배송 중', '', false, false),
                                _buildTimelineStep(
                                    '배송 완료', '', false, false),
                              ],
                            ),
                          ),
                          const SizedBox(height: 32),

                          // 버튼
                          Row(
                            children: [
                              Expanded(
                                child: OutlinedButton(
                                  onPressed: () => context.push(
                                      '/market/orders/${widget.orderId}'),
                                  style: OutlinedButton.styleFrom(
                                    minimumSize:
                                        const Size.fromHeight(48),
                                    foregroundColor:
                                        SanggamTheme.onSurfaceDim,
                                    side: const BorderSide(
                                        color: SanggamTheme
                                            .surfaceVariant),
                                    shape: RoundedRectangleBorder(
                                      borderRadius:
                                          BorderRadius.circular(16),
                                    ),
                                  ),
                                  child: const Text('주문 상세 보기'),
                                ),
                              ),
                              const SizedBox(width: 16),
                              Expanded(
                                child: FilledButton(
                                  onPressed: () =>
                                      context.go('/market'),
                                  style: FilledButton.styleFrom(
                                    minimumSize:
                                        const Size.fromHeight(48),
                                    backgroundColor:
                                        SanggamTheme.primary,
                                    foregroundColor:
                                        SanggamTheme.background,
                                    shape: RoundedRectangleBorder(
                                      borderRadius:
                                          BorderRadius.circular(16),
                                    ),
                                  ),
                                  child: const Text('계속 쇼핑'),
                                ),
                              ),
                            ],
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

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(label,
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 13,
                )),
          ),
          Expanded(
            child: Text(value,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                )),
          ),
        ],
      ),
    );
  }

  Widget _buildTimelineStep(
      String title, String subtitle, bool isActive, bool isCompleted) {
    final color = isCompleted
        ? SanggamTheme.jagaeCyan
        : isActive
            ? SanggamTheme.primary
            : SanggamTheme.onSurfaceDim;

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Container(
              width: 24,
              height: 24,
              decoration: BoxDecoration(
                color: isCompleted
                    ? SanggamTheme.jagaeCyan
                    : isActive
                        ? SanggamTheme.primary.withValues(alpha: 0.2)
                        : SanggamTheme.surfaceVariant,
                shape: BoxShape.circle,
                border: Border.all(color: color, width: 2),
              ),
              child: isCompleted
                  ? const Icon(Icons.check, size: 14, color: Colors.white)
                  : null,
            ),
            if (title != '배송 완료')
              Container(
                  width: 2,
                  height: 32,
                  color: isCompleted
                      ? SanggamTheme.jagaeCyan
                      : SanggamTheme.surfaceVariant),
          ],
        ),
        const SizedBox(width: 16),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.only(bottom: 24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title,
                    style: TextStyle(
                      color: isActive
                          ? Colors.white
                          : SanggamTheme.onSurfaceDim,
                      fontSize: 14,
                      fontWeight:
                          isActive ? FontWeight.bold : FontWeight.normal,
                    )),
                if (subtitle.isNotEmpty)
                  Text(subtitle,
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
  }

  String _statusLabel(String status) {
    const labels = {
      'pending': '주문 접수',
      'paid': '결제 완료',
      'preparing': '상품 준비 중',
      'shipping': '배송 중',
      'delivered': '배송 완료',
    };
    return labels[status] ?? status;
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (m) => '${m[1]},',
        );
  }
}
