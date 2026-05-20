import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// CheckoutScreen — Sanggam Orbit 결제
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData param 1x 제거
// [Rule 4] theme.textTheme ~8x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Card 3x → SanggamContainer
// [Rule 4] RadioListTile/CheckboxListTile 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 결제 화면
class CheckoutScreen
    extends ConsumerStatefulWidget {
  const CheckoutScreen({super.key});

  @override
  ConsumerState<CheckoutScreen>
      createState() =>
          _CheckoutScreenState();
}

class _CheckoutScreenState
    extends ConsumerState<CheckoutScreen> {
  String _paymentMethod = 'card';
  bool _agreeTerms = false;
  bool _agreePrivacy = false;
  bool _processing = false;

  @override
  Widget build(BuildContext context) {
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
                      '결제',
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
              child: SingleChildScrollView(
                padding:
                    const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .stretch,
                  children: [
                    // 배송지 정보
                    const Text('배송지 정보',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 8),
                    SanggamContainer(
                      borderRadius: 16,
                      child: Row(
                        children: [
                          const Icon(
                              Icons
                                  .location_on_outlined,
                              size: 20,
                              color: SanggamTheme
                                  .onSurfaceDim),
                          const SizedBox(
                              width: 8),
                          const Expanded(
                            child: Column(
                              crossAxisAlignment:
                                  CrossAxisAlignment
                                      .start,
                              children: [
                                Text('홍길동',
                                    style: TextStyle(
                                        color: Colors
                                            .white,
                                        fontWeight:
                                            FontWeight.w600)),
                                Text(
                                    '010-1234-5678',
                                    style: TextStyle(
                                        color: SanggamTheme
                                            .onSurfaceDim,
                                        fontSize:
                                            13)),
                                Text(
                                    '서울시 강남구 역삼동 123-45, 101동 1001호',
                                    style: TextStyle(
                                        color: SanggamTheme
                                            .onSurfaceDim,
                                        fontSize:
                                            13)),
                              ],
                            ),
                          ),
                          TextButton(
                            onPressed: () {},
                            style: TextButton
                                .styleFrom(
                                    foregroundColor:
                                        SanggamTheme
                                            .primary),
                            child: const Text(
                                '변경'),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 결제 수단
                    const Text('결제 수단',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 8),
                    SanggamContainer(
                      borderRadius: 16,
                      padding:
                          EdgeInsets.zero,
                      child: Column(
                        children: [
                          RadioListTile<
                              String>(
                            title: const Text(
                                '신용/체크카드',
                                style: TextStyle(
                                    color: Colors
                                        .white)),
                            secondary: const Icon(
                                Icons
                                    .credit_card,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            activeColor:
                                SanggamTheme
                                    .primary,
                            value: 'card',
                            groupValue:
                                _paymentMethod,
                            onChanged: (v) =>
                                setState(() =>
                                    _paymentMethod =
                                        v!),
                          ),
                          RadioListTile<
                              String>(
                            title: const Text(
                                '계좌이체',
                                style: TextStyle(
                                    color: Colors
                                        .white)),
                            secondary: const Icon(
                                Icons
                                    .account_balance,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            activeColor:
                                SanggamTheme
                                    .primary,
                            value: 'bank',
                            groupValue:
                                _paymentMethod,
                            onChanged: (v) =>
                                setState(() =>
                                    _paymentMethod =
                                        v!),
                          ),
                          RadioListTile<
                              String>(
                            title: const Text(
                                '휴대폰 결제',
                                style: TextStyle(
                                    color: Colors
                                        .white)),
                            secondary: const Icon(
                                Icons
                                    .phone_android,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            activeColor:
                                SanggamTheme
                                    .primary,
                            value: 'phone',
                            groupValue:
                                _paymentMethod,
                            onChanged: (v) =>
                                setState(() =>
                                    _paymentMethod =
                                        v!),
                          ),
                          RadioListTile<
                              String>(
                            title: const Text(
                                '간편결제 (TossPay)',
                                style: TextStyle(
                                    color: Colors
                                        .white)),
                            secondary: const Icon(
                                Icons.wallet,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            activeColor:
                                SanggamTheme
                                    .primary,
                            value: 'toss',
                            groupValue:
                                _paymentMethod,
                            onChanged: (v) =>
                                setState(() =>
                                    _paymentMethod =
                                        v!),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 주문 요약
                    const Text('주문 요약',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 8),
                    SanggamContainer(
                      borderRadius: 16,
                      child: Column(
                        children: [
                          _buildSummaryRow(
                              '혈당 측정 카트리지 (10개입) x1',
                              '₩29,900'),
                          _buildSummaryRow(
                              '콜레스테롤 카트리지 (5개입) x1',
                              '₩39,900'),
                          const Divider(
                              height: 24,
                              color: SanggamTheme
                                  .surfaceVariant),
                          _buildSummaryRow(
                              '상품 금액',
                              '₩69,800'),
                          _buildSummaryRow(
                              '배송비', '무료'),
                          _buildSummaryRow(
                              '할인', '-₩0'),
                          const Divider(
                              height: 24,
                              color: SanggamTheme
                                  .surfaceVariant),
                          const Row(
                            mainAxisAlignment:
                                MainAxisAlignment
                                    .spaceBetween,
                            children: [
                              Text(
                                  '총 결제 금액',
                                  style:
                                      TextStyle(
                                    color: Colors
                                        .white,
                                    fontSize:
                                        16,
                                    fontWeight:
                                        FontWeight
                                            .bold,
                                  )),
                              Text(
                                  '₩69,800',
                                  style:
                                      TextStyle(
                                    color: SanggamTheme
                                        .primary,
                                    fontSize:
                                        16,
                                    fontWeight:
                                        FontWeight
                                            .bold,
                                  )),
                            ],
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 16),

                    // 동의 체크박스
                    CheckboxListTile(
                      contentPadding:
                          EdgeInsets.zero,
                      dense: true,
                      activeColor:
                          SanggamTheme
                              .primary,
                      checkColor:
                          SanggamTheme
                              .background,
                      title: const Text(
                          '주문 내용을 확인하였으며 결제에 동의합니다',
                          style: TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 13)),
                      value: _agreeTerms,
                      onChanged: (v) =>
                          setState(() =>
                              _agreeTerms =
                                  v ??
                                      false),
                    ),
                    CheckboxListTile(
                      contentPadding:
                          EdgeInsets.zero,
                      dense: true,
                      activeColor:
                          SanggamTheme
                              .primary,
                      checkColor:
                          SanggamTheme
                              .background,
                      title: const Text(
                          '개인정보 수집 및 이용 동의',
                          style: TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 13)),
                      value: _agreePrivacy,
                      onChanged: (v) =>
                          setState(() =>
                              _agreePrivacy =
                                  v ??
                                      false),
                    ),
                    const SizedBox(
                        height: 16),

                    // 결제 버튼
                    FilledButton(
                      onPressed:
                          _agreeTerms &&
                                  _agreePrivacy &&
                                  !_processing
                              ? _processPayment
                              : null,
                      style: FilledButton
                          .styleFrom(
                        minimumSize:
                            const Size
                                .fromHeight(
                                56),
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
                      child: _processing
                          ? const SizedBox(
                              width: 24,
                              height: 24,
                              child: CircularProgressIndicator(
                                  strokeWidth:
                                      2,
                                  color: Colors
                                      .white))
                          : const Text(
                              '₩69,800 결제하기',
                              style:
                                  TextStyle(
                                fontSize: 16,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                    ),
                    const SizedBox(
                        height: 16),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryRow(
      String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 4),
      child: Row(
        mainAxisAlignment:
            MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 13,
              )),
          Text(value,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 13,
                fontWeight: FontWeight.w500,
              )),
        ],
      ),
    );
  }

  Future<void> _processPayment() async {
    setState(() => _processing = true);
    try {
      final client =
          ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ??
              '';

      final orderResp =
          await client.createOrder(
        userId: userId,
        shippingAddress:
            '서울시 강남구 역삼동 123-45, 101동 1001호',
        paymentMethod: _paymentMethod,
      );
      final orderId = orderResp['order_id']
              as String? ??
          orderResp['id'] as String? ??
          '';

      final paymentResp =
          await client.createPayment(
        userId: userId,
        orderId: orderId,
        paymentType: 0,
        amountKrw: 69800,
        paymentMethod: _paymentMethod,
      );
      final paymentId =
          paymentResp['payment_id']
                  as String? ??
              paymentResp['id']
                  as String? ??
              '';

      await client.confirmPayment(
        paymentId,
        pgTransactionId:
            'TOSS-${DateTime.now().millisecondsSinceEpoch}',
        pgProvider: 'tosspayments',
      );

      if (mounted) {
        context.go(
            '/market/order-complete/$orderId');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(
          SnackBar(
              content: Text(
                  '결제 처리 중 오류가 발생했습니다: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(
            () => _processing = false);
      }
    }
  }
}
