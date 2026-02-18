import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';

/// 구독 플랜 비교 화면
class PlanComparisonScreen extends ConsumerWidget {
  const PlanComparisonScreen({super.key, this.mode});

  final String? mode;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final title = mode == 'upgrade' ? '플랜 업그레이드' : mode == 'downgrade' ? '플랜 다운그레이드' : '구독 플랜 비교';
    final plansAsync = ref.watch(subscriptionPlansProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(title),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: plansAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => _buildStaticPlans(context, theme),
        data: (plans) {
          if (plans.isEmpty) return _buildStaticPlans(context, theme);
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              Text('나에게 맞는 플랜을 선택하세요', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              ...plans.map((plan) => _buildDynamicPlanCard(context, theme, plan)),
            ],
          );
        },
      ),
    );
  }

  Widget _buildStaticPlans(BuildContext context, ThemeData theme) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Text('나에게 맞는 플랜을 선택하세요', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 16),
        ..._fallbackPlans.map((plan) => _buildPlanCard(context, theme, plan)),
      ],
    );
  }

  Widget _buildDynamicPlanCard(BuildContext context, ThemeData theme, SubscriptionPlan plan) {
    final isRecommended = plan.name.toLowerCase().contains('pro');
    final priceText = plan.monthlyPrice > 0 ? '${_formatPrice(plan.monthlyPrice)}원' : '무료';
    final features = [
      '카트리지 ${plan.cartridgesPerMonth}개/월',
      ...plan.includedCartridgeTypes.map((t) => '$t 포함'),
      if (plan.discountPercent > 0) '${plan.discountPercent}% 할인',
    ];

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: isRecommended
          ? RoundedRectangleBorder(borderRadius: BorderRadius.circular(12), side: const BorderSide(color: AppTheme.sanggamGold, width: 2))
          : null,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(plan.name, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                if (isRecommended) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(color: AppTheme.sanggamGold, borderRadius: BorderRadius.circular(12)),
                    child: const Text('추천', style: TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.bold)),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 12),
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(priceText, style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: AppTheme.sanggamGold)),
                if (plan.monthlyPrice > 0) Text(' / 월', style: theme.textTheme.bodyMedium),
              ],
            ),
            const SizedBox(height: 12),
            ...features.map((f) => Padding(
              padding: const EdgeInsets.only(bottom: 4),
              child: Row(
                children: [
                  const Icon(Icons.check_circle, size: 16, color: Colors.green),
                  const SizedBox(width: 8),
                  Expanded(child: Text(f, style: theme.textTheme.bodySmall)),
                ],
              ),
            )),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: plan.monthlyPrice == 0
                  ? OutlinedButton(onPressed: () {}, child: const Text('현재 플랜'))
                  : FilledButton(
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('${plan.name} 플랜이 선택되었습니다.')));
                        context.pop();
                      },
                      style: FilledButton.styleFrom(backgroundColor: isRecommended ? AppTheme.sanggamGold : null),
                      child: Text(mode == 'upgrade' ? '업그레이드' : '선택하기'),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlanCard(BuildContext context, ThemeData theme, _PlanData plan) {
    final isRecommended = plan.name == 'Pro';

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      shape: isRecommended
          ? RoundedRectangleBorder(borderRadius: BorderRadius.circular(12), side: const BorderSide(color: AppTheme.sanggamGold, width: 2))
          : null,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(plan.name, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                if (isRecommended) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(color: AppTheme.sanggamGold, borderRadius: BorderRadius.circular(12)),
                    child: const Text('추천', style: TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.bold)),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 4),
            Text(plan.description, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 12),
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(plan.price, style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: AppTheme.sanggamGold)),
                if (plan.period.isNotEmpty) Text(' / ${plan.period}', style: theme.textTheme.bodyMedium),
              ],
            ),
            const SizedBox(height: 12),
            ...plan.features.map((f) => Padding(
              padding: const EdgeInsets.only(bottom: 4),
              child: Row(
                children: [
                  Icon(f.included ? Icons.check_circle : Icons.remove_circle_outline,
                      size: 16, color: f.included ? Colors.green : theme.colorScheme.outlineVariant),
                  const SizedBox(width: 8),
                  Expanded(child: Text(f.label, style: theme.textTheme.bodySmall?.copyWith(
                    color: f.included ? null : theme.colorScheme.onSurfaceVariant,
                  ))),
                ],
              ),
            )),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: plan.price == '무료'
                  ? OutlinedButton(onPressed: () {}, child: const Text('현재 플랜'))
                  : FilledButton(
                      onPressed: () {
                        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('${plan.name} 플랜이 선택되었습니다.')));
                        context.pop();
                      },
                      style: FilledButton.styleFrom(backgroundColor: isRecommended ? AppTheme.sanggamGold : null),
                      child: Text(mode == 'upgrade' ? '업그레이드' : '선택하기'),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatPrice(int price) {
    final str = price.toString();
    final buf = StringBuffer();
    for (var i = 0; i < str.length; i++) {
      if (i > 0 && (str.length - i) % 3 == 0) buf.write(',');
      buf.write(str[i]);
    }
    return buf.toString();
  }

  static final _fallbackPlans = [
    _PlanData(name: 'Free', description: '기본 기능 체험', price: '무료', period: '', features: [
      _Feature('기본 측정 1회/일', true), _Feature('측정 기록 7일 보관', true),
      _Feature('AI 코칭', false), _Feature('가족 공유', false), _Feature('원격 진료', false),
    ]),
    _PlanData(name: 'Basic', description: '일상 건강 관리', price: '9,900원', period: '월', features: [
      _Feature('무제한 측정', true), _Feature('측정 기록 무제한 보관', true),
      _Feature('AI 건강 코칭', true), _Feature('데이터 내보내기', true), _Feature('가족 공유 (2명)', false),
    ]),
    _PlanData(name: 'Pro', description: '가족 건강 케어', price: '19,900원', period: '월', features: [
      _Feature('무제한 측정', true), _Feature('AI 고급 분석', true),
      _Feature('가족 공유 (5명)', true), _Feature('원격 진료 월 2회', true), _Feature('우선 고객 지원', true),
    ]),
    _PlanData(name: 'Clinical', description: '전문가급 분석', price: '39,900원', period: '월', features: [
      _Feature('무제한 측정 + 연구용 데이터', true), _Feature('FHIR 의료 데이터 연동', true),
      _Feature('가족 공유 (무제한)', true), _Feature('원격 진료 무제한', true), _Feature('전담 건강 매니저', true),
    ]),
  ];
}

class _PlanData {
  final String name, description, price, period;
  final List<_Feature> features;
  const _PlanData({required this.name, required this.description, required this.price, required this.period, required this.features});
}

class _Feature {
  final String label;
  final bool included;
  const _Feature(this.label, this.included);
}
