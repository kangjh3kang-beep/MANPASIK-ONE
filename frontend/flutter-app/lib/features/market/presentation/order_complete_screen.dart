import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';
import 'package:lottie/lottie.dart';

/// 주문 완료 & 배송 추적 화면
class OrderCompleteScreen extends StatelessWidget {
  const OrderCompleteScreen({super.key, required this.orderId});

  final String orderId;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('주문 완료'),
        automaticallyImplyLeading: false,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          children: [
            // 주문 완료 아이콘 (Lottie)
            const SizedBox(height: 24),
            SizedBox(
              width: 120,
              height: 120,
              child: Lottie.asset(
                'assets/lottie/check_success.json',
                repeat: false,
                errorBuilder: (_, __, ___) => Container(
                  width: 96, height: 96,
                  decoration: BoxDecoration(color: Colors.green.withOpacity(0.1), shape: BoxShape.circle),
                  child: const Icon(Icons.check_circle, size: 64, color: Colors.green),
                ),
              ),
            ),
            const SizedBox(height: 24),
            Text('주문이 완료되었습니다!', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Text('주문번호: $orderId', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.outline)),
            const SizedBox(height: 32),

            // 배송 정보
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('배송 정보', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 12),
                    _buildInfoRow(theme, '예상 도착일', '2026-02-18 (3일 이내)'),
                    _buildInfoRow(theme, '택배사', 'CJ대한통운'),
                    _buildInfoRow(theme, '수령인', '홍길동'),
                    _buildInfoRow(theme, '배송지', '서울시 강남구 역삼동 123-45'),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 배송 상태 타임라인
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('배송 추적', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 16),
                    _buildTimelineStep(theme, '주문 접수', '2026-02-15 14:30', true, true),
                    _buildTimelineStep(theme, '상품 준비', '처리 중', true, false),
                    _buildTimelineStep(theme, '발송', '대기 중', false, false),
                    _buildTimelineStep(theme, '배송 중', '', false, false),
                    _buildTimelineStep(theme, '배송 완료', '', false, false),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 32),

            // 버튼
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () => context.push('/market/orders'),
                    style: OutlinedButton.styleFrom(minimumSize: const Size.fromHeight(48)),
                    child: const Text('주문 상세 보기'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: FilledButton(
                    onPressed: () => context.go('/market'),
                    style: FilledButton.styleFrom(
                      minimumSize: const Size.fromHeight(48),
                      backgroundColor: AppTheme.sanggamGold,
                    ),
                    child: const Text('계속 쇼핑'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(label, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
          ),
          Expanded(child: Text(value, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w500))),
        ],
      ),
    );
  }

  Widget _buildTimelineStep(ThemeData theme, String title, String subtitle, bool isActive, bool isCompleted) {
    final color = isCompleted ? Colors.green : isActive ? AppTheme.sanggamGold : theme.colorScheme.outline;

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Container(
              width: 24,
              height: 24,
              decoration: BoxDecoration(
                color: isCompleted ? Colors.green : isActive ? AppTheme.sanggamGold.withOpacity(0.2) : theme.colorScheme.surfaceContainerHighest,
                shape: BoxShape.circle,
                border: Border.all(color: color, width: 2),
              ),
              child: isCompleted ? const Icon(Icons.check, size: 14, color: Colors.white) : null,
            ),
            if (title != '배송 완료')
              Container(width: 2, height: 32, color: isCompleted ? Colors.green : theme.dividerColor),
          ],
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.only(bottom: 24),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: theme.textTheme.bodyMedium?.copyWith(
                  fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
                  color: isActive ? null : theme.colorScheme.outline,
                )),
                if (subtitle.isNotEmpty)
                  Text(subtitle, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
