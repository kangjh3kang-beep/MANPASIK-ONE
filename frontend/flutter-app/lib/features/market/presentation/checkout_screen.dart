import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 결제 화면
class CheckoutScreen extends ConsumerStatefulWidget {
  const CheckoutScreen({super.key});

  @override
  ConsumerState<CheckoutScreen> createState() => _CheckoutScreenState();
}

class _CheckoutScreenState extends ConsumerState<CheckoutScreen> {
  String _paymentMethod = 'card';
  bool _agreeTerms = false;
  bool _agreePrivacy = false;
  bool _processing = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('결제'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 배송지 정보
            Text('배송지 정보', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.location_on_outlined, size: 20),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('홍길동', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                              Text('010-1234-5678', style: theme.textTheme.bodySmall),
                              Text('서울시 강남구 역삼동 123-45, 101동 1001호', style: theme.textTheme.bodySmall),
                            ],
                          ),
                        ),
                        TextButton(onPressed: () {}, child: const Text('변경')),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            // 결제 수단
            Text('결제 수단', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Card(
              child: Column(
                children: [
                  RadioListTile<String>(
                    title: const Text('신용/체크카드'),
                    secondary: const Icon(Icons.credit_card),
                    value: 'card',
                    groupValue: _paymentMethod,
                    onChanged: (v) => setState(() => _paymentMethod = v!),
                  ),
                  RadioListTile<String>(
                    title: const Text('계좌이체'),
                    secondary: const Icon(Icons.account_balance),
                    value: 'bank',
                    groupValue: _paymentMethod,
                    onChanged: (v) => setState(() => _paymentMethod = v!),
                  ),
                  RadioListTile<String>(
                    title: const Text('휴대폰 결제'),
                    secondary: const Icon(Icons.phone_android),
                    value: 'phone',
                    groupValue: _paymentMethod,
                    onChanged: (v) => setState(() => _paymentMethod = v!),
                  ),
                  RadioListTile<String>(
                    title: const Text('간편결제 (TossPay)'),
                    secondary: const Icon(Icons.wallet),
                    value: 'toss',
                    groupValue: _paymentMethod,
                    onChanged: (v) => setState(() => _paymentMethod = v!),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // 주문 요약
            Text('주문 요약', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    _buildSummaryRow(theme, '혈당 측정 카트리지 (10개입) x1', '₩29,900'),
                    _buildSummaryRow(theme, '콜레스테롤 카트리지 (5개입) x1', '₩39,900'),
                    const Divider(height: 24),
                    _buildSummaryRow(theme, '상품 금액', '₩69,800'),
                    _buildSummaryRow(theme, '배송비', '무료'),
                    _buildSummaryRow(theme, '할인', '-₩0'),
                    const Divider(height: 24),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('총 결제 금액', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                        Text('₩69,800', style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: AppTheme.sanggamGold,
                        )),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 동의 체크박스
            CheckboxListTile(
              contentPadding: EdgeInsets.zero,
              dense: true,
              title: const Text('주문 내용을 확인하였으며 결제에 동의합니다', style: TextStyle(fontSize: 13)),
              value: _agreeTerms,
              onChanged: (v) => setState(() => _agreeTerms = v ?? false),
            ),
            CheckboxListTile(
              contentPadding: EdgeInsets.zero,
              dense: true,
              title: const Text('개인정보 수집 및 이용 동의', style: TextStyle(fontSize: 13)),
              value: _agreePrivacy,
              onChanged: (v) => setState(() => _agreePrivacy = v ?? false),
            ),
            const SizedBox(height: 16),

            // 결제 버튼
            FilledButton(
              onPressed: _agreeTerms && _agreePrivacy && !_processing ? _processPayment : null,
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(56),
                backgroundColor: AppTheme.sanggamGold,
              ),
              child: _processing
                  ? const SizedBox(width: 24, height: 24, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                  : const Text('₩69,800 결제하기', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: theme.textTheme.bodySmall),
          Text(value, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }

  Future<void> _processPayment() async {
    setState(() => _processing = true);
    try {
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';

      // 1. 주문 생성
      final orderResp = await client.createOrder(
        userId: userId,
        shippingAddress: '서울시 강남구 역삼동 123-45, 101동 1001호',
        paymentMethod: _paymentMethod,
      );
      final orderId = orderResp['order_id'] as String? ?? orderResp['id'] as String? ?? '';

      // 2. 결제 요청 생성
      final paymentResp = await client.createPayment(
        userId: userId,
        orderId: orderId,
        paymentType: 0, // PRODUCT_PURCHASE
        amountKrw: 69800,
        paymentMethod: _paymentMethod,
      );
      final paymentId = paymentResp['payment_id'] as String? ?? paymentResp['id'] as String? ?? '';

      // 3. PG 결제 승인 (Toss Payments SDK 연동 시 여기에 위젯 호출)
      // 현재는 서버 사이드 승인으로 처리
      await client.confirmPayment(
        paymentId,
        pgTransactionId: 'TOSS-${DateTime.now().millisecondsSinceEpoch}',
        pgProvider: 'tosspayments',
      );

      if (mounted) {
        context.go('/market/order-complete/$orderId');
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('결제 처리 중 오류가 발생했습니다: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _processing = false);
    }
  }
}
