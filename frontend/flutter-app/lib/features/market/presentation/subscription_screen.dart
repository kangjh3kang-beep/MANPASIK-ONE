import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';

/// 구독 관리 화면
///
/// 현재 플랜 표시 + 4개 플랜 비교표 + 변경/취소
class SubscriptionScreen extends ConsumerWidget {
  const SubscriptionScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final plansAsync = ref.watch(subscriptionPlansProvider);
    final subAsync = ref.watch(subscriptionInfoProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('구독 관리'),
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
            // 현재 구독 상태
            _buildCurrentPlan(theme, subAsync),
            const SizedBox(height: 24),

            Text('구독 플랜 비교', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 12),

            // 플랜 비교 카드들
            plansAsync.when(
              data: (plans) {
                if (plans.isEmpty) return _buildFallbackPlans(context, theme);
                return Column(
                  children: plans.map((p) => _buildPlanCard(context, theme, p)).toList(),
                );
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (_, __) => _buildFallbackPlans(context, theme),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCurrentPlan(ThemeData theme, AsyncValue<dynamic> subAsync) {
    return Card(
      color: AppTheme.sanggamGold.withOpacity(0.1),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('현재 구독', style: theme.textTheme.labelLarge?.copyWith(color: AppTheme.sanggamGold)),
            const SizedBox(height: 8),
            subAsync.when(
              data: (sub) {
                if (sub == null) {
                  return Row(
                    children: [
                      const Icon(Icons.card_membership, color: AppTheme.sanggamGold),
                      const SizedBox(width: 12),
                      Text('Free 플랜', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    ],
                  );
                }
                final tierName = ['Free', 'Basic', 'Pro', 'Clinical'][sub.tier.clamp(0, 3)];
                return Row(
                  children: [
                    const Icon(Icons.card_membership, color: AppTheme.sanggamGold),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('$tierName 플랜', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                          Text('활성 상태', style: theme.textTheme.bodySmall?.copyWith(color: Colors.green)),
                        ],
                      ),
                    ),
                  ],
                );
              },
              loading: () => const CircularProgressIndicator(strokeWidth: 2),
              error: (_, __) => const Text('구독 정보를 불러올 수 없습니다.'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFallbackPlans(BuildContext context, ThemeData theme) {
    final plans = [
      _PlanInfo('Free', '₩0/월', ['기본 측정 1회/일', '리더기 1대', '기본 AI 분석'], false),
      _PlanInfo('Basic', '₩9,900/월', ['무제한 측정', '리더기 2대', 'AI 건강 코칭', '카트리지 도감'], false),
      _PlanInfo('Pro', '₩29,900/월', ['모든 Basic 기능', '리더기 5대', '가족 공유 (5명)', '비대면 진료', '데이터 내보내기 (FHIR)'], true),
      _PlanInfo('Clinical', '₩99,000/월', ['모든 Pro 기능', '리더기 10대', '의료기관급 분석', 'MFA 필수 보안', '전용 고객 지원'], false),
    ];
    return Column(
      children: plans.map((p) => _buildFallbackPlanCard(context, theme, p)).toList(),
    );
  }

  Widget _buildFallbackPlanCard(BuildContext context, ThemeData theme, _PlanInfo plan) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: plan.recommended
          ? RoundedRectangleBorder(
              side: const BorderSide(color: AppTheme.sanggamGold, width: 2),
              borderRadius: BorderRadius.circular(12),
            )
          : null,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(plan.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                if (plan.recommended)
                  Chip(
                    label: const Text('추천', style: TextStyle(fontSize: 11, color: Colors.white)),
                    backgroundColor: AppTheme.sanggamGold,
                    side: BorderSide.none,
                    visualDensity: VisualDensity.compact,
                  ),
              ],
            ),
            Text(plan.price, style: theme.textTheme.titleLarge?.copyWith(color: AppTheme.sanggamGold, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            ...plan.features.map((f) => Padding(
              padding: const EdgeInsets.symmetric(vertical: 2),
              child: Row(
                children: [
                  const Icon(Icons.check, size: 16, color: Colors.green),
                  const SizedBox(width: 8),
                  Expanded(child: Text(f, style: theme.textTheme.bodySmall)),
                ],
              ),
            )),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: plan.name == 'Free'
                  ? const SizedBox.shrink()
                  : OutlinedButton(
                      onPressed: () => context.push('/market/subscription/upgrade'),
                      child: Text('${plan.name} 구독하기'),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlanCard(BuildContext context, ThemeData theme, SubscriptionPlan plan) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(plan.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(
              '₩${plan.monthlyPrice.toString().replaceAllMapped(RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}/월',
              style: theme.textTheme.titleLarge?.copyWith(color: AppTheme.sanggamGold, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Text('${plan.cartridgesPerMonth}개 카트리지/월 | ${plan.discountPercent}% 할인', style: theme.textTheme.bodySmall),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton(
                onPressed: () => _subscribePlan(context, plan),
                child: Text('${plan.name} 구독하기'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _subscribePlan(BuildContext context, SubscriptionPlan plan) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('${plan.name} 구독'),
        content: Text(
          '월 ₩${plan.monthlyPrice.toString().replaceAllMapped(RegExp(r"(\d)(?=(\d{3})+(?!\d))"), (m) => "${m[1]},")}으로 '
          '${plan.name} 플랜을 구독하시겠습니까?',
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              context.push('/market/checkout');
            },
            child: const Text('결제 진행'),
          ),
        ],
      ),
    );
  }
}

class _PlanInfo {
  final String name, price;
  final List<String> features;
  final bool recommended;
  const _PlanInfo(this.name, this.price, this.features, this.recommended);
}
