import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// OrderHistoryScreen — Sanggam Orbit 주문 내역
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData param 1x 제거
// [Rule 4] theme.textTheme ~6x → 직접 TextStyle
// [Rule 4] theme.colorScheme 1x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Colors mapped (orange→primary, blue→jagaeCyan, green→jagaeCyan, red→error)
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16
// ───────────────────────────────────────────────────

/// 주문 내역 화면
class OrderHistoryScreen
    extends ConsumerWidget {
  const OrderHistoryScreen({super.key});

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final ordersAsync =
        ref.watch(ordersProvider);

    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '주문 내역',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: RefreshIndicator(
                color: SanggamTheme.primary,
                onRefresh: () async => ref
                    .invalidate(
                        ordersProvider),
                child: ordersAsync.when(
                  data: (orders) {
                    if (orders.isEmpty) {
                      return Center(
                        child: Column(
                          mainAxisSize:
                              MainAxisSize
                                  .min,
                          children: [
                            const Icon(
                                Icons
                                    .receipt_long_outlined,
                                size: 64,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            const SizedBox(
                                height: 16),
                            const Text(
                                '주문 내역이 없습니다.',
                                style: TextStyle(
                                    color: Colors
                                        .white,
                                    fontSize:
                                        16)),
                            const SizedBox(
                                height: 8),
                            FilledButton(
                              onPressed: () =>
                                  context.go(
                                      '/market'),
                              style: FilledButton
                                  .styleFrom(
                                backgroundColor:
                                    SanggamTheme
                                        .primary,
                                foregroundColor:
                                    SanggamTheme
                                        .background,
                                shape:
                                    RoundedRectangleBorder(
                                  borderRadius:
                                      BorderRadius
                                          .circular(
                                              16),
                                ),
                              ),
                              child: const Text(
                                  '쇼핑하러 가기'),
                            ),
                          ],
                        ),
                      );
                    }
                    return ListView.builder(
                      padding:
                          const EdgeInsets
                              .all(16),
                      itemCount:
                          orders.length,
                      itemBuilder:
                          (context, index) =>
                              _buildOrderCard(
                                  context,
                                  orders[
                                      index]),
                    );
                  },
                  loading: () => const Center(
                      child:
                          CircularProgressIndicator(
                              color:
                                  SanggamTheme
                                      .primary)),
                  error: (_, __) => Center(
                    child: Column(
                      mainAxisSize:
                          MainAxisSize.min,
                      children: [
                        const Icon(
                            Icons
                                .error_outline,
                            size: 48,
                            color:
                                SanggamTheme
                                    .error),
                        const SizedBox(
                            height: 8),
                        const Text(
                            '주문 내역을 불러올 수 없습니다.',
                            style: TextStyle(
                                color: Colors
                                    .white)),
                        const SizedBox(
                            height: 8),
                        FilledButton(
                          onPressed: () => ref
                              .invalidate(
                                  ordersProvider),
                          style: FilledButton
                              .styleFrom(
                            backgroundColor:
                                SanggamTheme
                                    .primary,
                            foregroundColor:
                                SanggamTheme
                                    .background,
                            shape:
                                RoundedRectangleBorder(
                              borderRadius:
                                  BorderRadius
                                      .circular(
                                          16),
                            ),
                          ),
                          child: const Text(
                              '다시 시도'),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOrderCard(
      BuildContext context, Order order) {
    final statusColor =
        _statusColor(order.status);
    final statusText =
        _statusText(order.status);

    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: Theme(
          data: ThemeData.dark().copyWith(
            dividerColor:
                SanggamTheme.surfaceVariant,
          ),
          child: ExpansionTile(
            title: Row(
              children: [
                Expanded(
                  child: Text(
                    '주문 #${order.id.length > 8 ? order.id.substring(0, 8) : order.id}',
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight:
                          FontWeight.w600,
                    ),
                  ),
                ),
                Container(
                  padding: const EdgeInsets
                      .symmetric(
                      horizontal: 8,
                      vertical: 2),
                  decoration: BoxDecoration(
                    color: statusColor
                        .withValues(
                            alpha: 0.1),
                    borderRadius:
                        BorderRadius.circular(
                            8),
                  ),
                  child: Text(statusText,
                      style: TextStyle(
                        fontSize: 11,
                        color: statusColor,
                      )),
                ),
              ],
            ),
            subtitle: Text(
              '${order.orderedAt.year}-${order.orderedAt.month.toString().padLeft(2, '0')}-${order.orderedAt.day.toString().padLeft(2, '0')}',
              style: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 13,
              ),
            ),
            children: [
              ...order.items.map((item) =>
                  ListTile(
                    dense: true,
                    leading: const Icon(
                        Icons.science,
                        size: 20,
                        color: SanggamTheme
                            .primary),
                    title: Text(
                        item.productName,
                        style:
                            const TextStyle(
                          color:
                              Colors.white,
                          fontSize: 13,
                        )),
                    trailing: Text(
                        '${item.quantity}개  ₩${_formatPrice(item.unitPrice * item.quantity)}',
                        style:
                            const TextStyle(
                          color:
                              Colors.white,
                          fontSize: 13,
                        )),
                  )),
              const Divider(
                  indent: 16,
                  endIndent: 16,
                  color: SanggamTheme
                      .surfaceVariant),
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 8),
                child: Row(
                  mainAxisAlignment:
                      MainAxisAlignment
                          .spaceBetween,
                  children: [
                    const Text('합계',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 14,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    Text(
                      '₩${_formatPrice(order.totalAmount)}',
                      style: const TextStyle(
                        color: SanggamTheme
                            .primary,
                        fontSize: 16,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 8),
                child: SizedBox(
                  width: double.infinity,
                  child:
                      OutlinedButton.icon(
                    onPressed: () =>
                        context.push(
                            '/market/order/${order.id}'),
                    icon: const Icon(
                        Icons
                            .receipt_outlined,
                        size: 18),
                    label: const Text(
                        '주문 상세'),
                    style: OutlinedButton
                        .styleFrom(
                      foregroundColor:
                          SanggamTheme
                              .onSurfaceDim,
                      side: const BorderSide(
                          color: SanggamTheme
                              .surfaceVariant),
                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius
                                .circular(
                                    16),
                      ),
                    ),
                  ),
                ),
              ),
              if (order.status ==
                      OrderStatus
                          .confirmed ||
                  order.status ==
                      OrderStatus.shipping)
                Padding(
                  padding:
                      const EdgeInsets
                          .fromLTRB(
                          16, 0, 16, 16),
                  child: SizedBox(
                    width: double.infinity,
                    child:
                        OutlinedButton.icon(
                      onPressed: () =>
                          context.push(
                              '/market/order-complete/${order.id}'),
                      icon: const Icon(
                          Icons
                              .local_shipping_outlined,
                          size: 18),
                      label: const Text(
                          '배송 추적'),
                      style: OutlinedButton
                          .styleFrom(
                        foregroundColor:
                            SanggamTheme
                                .onSurfaceDim,
                        side: const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  Color _statusColor(OrderStatus status) {
    switch (status) {
      case OrderStatus.pending:
        return SanggamTheme.primary;
      case OrderStatus.confirmed:
        return SanggamTheme.jagaeCyan;
      case OrderStatus.shipping:
        return SanggamTheme.jagaeMagenta;
      case OrderStatus.delivered:
        return SanggamTheme.jagaeCyan;
      case OrderStatus.cancelled:
        return SanggamTheme.error;
    }
  }

  String _statusText(OrderStatus status) {
    switch (status) {
      case OrderStatus.pending:
        return '결제 대기';
      case OrderStatus.confirmed:
        return '결제 완료';
      case OrderStatus.shipping:
        return '배송 중';
      case OrderStatus.delivered:
        return '배송 완료';
      case OrderStatus.cancelled:
        return '취소됨';
    }
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
          RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
          (m) => '${m[1]},',
        );
  }
}
