import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// PlanComparisonScreen — Sanggam Orbit 플랜 비교
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData params 3x 제거
// [Rule 4] theme.textTheme ~12x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~3x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold ~8x → SanggamTheme.primary
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16
// ───────────────────────────────────────────────────

/// 구독 플랜 비교 화면
class PlanComparisonScreen
    extends ConsumerWidget {
  const PlanComparisonScreen(
      {super.key, this.mode});

  final String? mode;

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final title = mode == 'upgrade'
        ? '플랜 업그레이드'
        : mode == 'downgrade'
            ? '플랜 다운그레이드'
            : '구독 플랜 비교';
    final plansAsync =
        ref.watch(subscriptionPlansProvider);

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
                  Expanded(
                    child: Text(
                      title,
                      style: const TextStyle(
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
              child: plansAsync.when(
                loading: () => const Center(
                    child:
                        CircularProgressIndicator(
                            color:
                                SanggamTheme
                                    .primary)),
                error: (e, _) =>
                    _buildStaticPlans(
                        context),
                data: (plans) {
                  if (plans.isEmpty) {
                    return _buildStaticPlans(
                        context);
                  }
                  return ListView(
                    padding:
                        const EdgeInsets.all(
                            16),
                    children: [
                      const Text(
                          '나에게 맞는 플랜을 선택하세요',
                          style: TextStyle(
                            color: Colors
                                .white,
                            fontSize: 16,
                            fontWeight:
                                FontWeight
                                    .bold,
                          )),
                      const SizedBox(
                          height: 16),
                      ...plans.map((plan) =>
                          _buildDynamicPlanCard(
                              context,
                              plan)),
                    ],
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStaticPlans(
      BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        const Text('나에게 맞는 플랜을 선택하세요',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            )),
        const SizedBox(height: 16),
        ..._fallbackPlans.map((plan) =>
            _buildPlanCard(context, plan)),
      ],
    );
  }

  Widget _buildDynamicPlanCard(
      BuildContext context,
      SubscriptionPlan plan) {
    final isRecommended = plan.name
        .toLowerCase()
        .contains('pro');
    final priceText =
        plan.monthlyPrice > 0
            ? '${_formatPrice(plan.monthlyPrice)}원'
            : '무료';
    final features = [
      '카트리지 ${plan.cartridgesPerMonth}개/월',
      ...plan.includedCartridgeTypes
          .map((t) => '$t 포함'),
      if (plan.discountPercent > 0)
        '${plan.discountPercent}% 할인',
    ];

    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        borderColor: isRecommended
            ? SanggamTheme.primary
            : SanggamTheme.surfaceVariant,
        borderWidth:
            isRecommended ? 2.0 : 1.5,
        child: Column(
          crossAxisAlignment:
              CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(plan.name,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 20,
                      fontWeight:
                          FontWeight.bold,
                    )),
                if (isRecommended) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding:
                        const EdgeInsets
                            .symmetric(
                            horizontal: 8,
                            vertical: 2),
                    decoration:
                        BoxDecoration(
                      color: SanggamTheme
                          .primary,
                      borderRadius:
                          BorderRadius
                              .circular(
                                  12),
                    ),
                    child: const Text('추천',
                        style: TextStyle(
                          color: SanggamTheme
                              .background,
                          fontSize: 11,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 16),
            Row(
              crossAxisAlignment:
                  CrossAxisAlignment.end,
              children: [
                Text(priceText,
                    style: const TextStyle(
                      color: SanggamTheme
                          .primary,
                      fontSize: 28,
                      fontWeight:
                          FontWeight.bold,
                    )),
                if (plan.monthlyPrice > 0)
                  const Text(' / 월',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 14,
                      )),
              ],
            ),
            const SizedBox(height: 16),
            ...features.map((f) => Padding(
                  padding:
                      const EdgeInsets.only(
                          bottom: 4),
                  child: Row(
                    children: [
                      const Icon(
                          Icons
                              .check_circle,
                          size: 16,
                          color: SanggamTheme
                              .jagaeCyan),
                      const SizedBox(
                          width: 8),
                      Expanded(
                          child: Text(f,
                              style:
                                  const TextStyle(
                                color: Colors
                                    .white,
                                fontSize:
                                    13,
                              ))),
                    ],
                  ),
                )),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: plan.monthlyPrice == 0
                  ? OutlinedButton(
                      onPressed: () {},
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
                      child: const Text(
                          '현재 플랜'),
                    )
                  : FilledButton(
                      onPressed: () {
                        ScaffoldMessenger.of(
                                context)
                            .showSnackBar(SnackBar(
                                content: Text(
                                    '${plan.name} 플랜이 선택되었습니다.')));
                        context.pop();
                      },
                      style: FilledButton
                          .styleFrom(
                        backgroundColor:
                            isRecommended
                                ? SanggamTheme
                                    .primary
                                : SanggamTheme
                                    .surface,
                        foregroundColor:
                            isRecommended
                                ? SanggamTheme
                                    .background
                                : Colors
                                    .white,
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                      child: Text(
                          mode == 'upgrade'
                              ? '업그레이드'
                              : '선택하기'),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlanCard(
      BuildContext context, _PlanData plan) {
    final isRecommended = plan.name == '프로';

    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        borderColor: isRecommended
            ? SanggamTheme.primary
            : SanggamTheme.surfaceVariant,
        borderWidth:
            isRecommended ? 2.0 : 1.5,
        child: Column(
          crossAxisAlignment:
              CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(plan.name,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 20,
                      fontWeight:
                          FontWeight.bold,
                    )),
                if (isRecommended) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding:
                        const EdgeInsets
                            .symmetric(
                            horizontal: 8,
                            vertical: 2),
                    decoration:
                        BoxDecoration(
                      color: SanggamTheme
                          .primary,
                      borderRadius:
                          BorderRadius
                              .circular(
                                  12),
                    ),
                    child: const Text('추천',
                        style: TextStyle(
                          color: SanggamTheme
                              .background,
                          fontSize: 11,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                  ),
                ],
              ],
            ),
            const SizedBox(height: 4),
            Text(plan.description,
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 13,
                )),
            const SizedBox(height: 16),
            Row(
              crossAxisAlignment:
                  CrossAxisAlignment.end,
              children: [
                Text(plan.price,
                    style: const TextStyle(
                      color: SanggamTheme
                          .primary,
                      fontSize: 28,
                      fontWeight:
                          FontWeight.bold,
                    )),
                if (plan
                    .period.isNotEmpty)
                  Text(' / ${plan.period}',
                      style:
                          const TextStyle(
                        color:
                            Colors.white,
                        fontSize: 14,
                      )),
              ],
            ),
            const SizedBox(height: 16),
            ...plan.features.map(
                (f) => Padding(
                      padding:
                          const EdgeInsets
                              .only(
                              bottom: 4),
                      child: Row(
                        children: [
                          Icon(
                              f.included
                                  ? Icons
                                      .check_circle
                                  : Icons
                                      .remove_circle_outline,
                              size: 16,
                              color: f
                                      .included
                                  ? SanggamTheme
                                      .jagaeCyan
                                  : SanggamTheme
                                      .onSurfaceDim),
                          const SizedBox(
                              width: 8),
                          Expanded(
                              child: Text(
                                  f.label,
                                  style:
                                      TextStyle(
                                    color: f.included
                                        ? Colors.white
                                        : SanggamTheme.onSurfaceDim,
                                    fontSize:
                                        13,
                                  ))),
                        ],
                      ),
                    )),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: plan.price == '무료'
                  ? OutlinedButton(
                      onPressed: () {},
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
                      child: const Text(
                          '현재 플랜'),
                    )
                  : FilledButton(
                      onPressed: () {
                        ScaffoldMessenger.of(
                                context)
                            .showSnackBar(SnackBar(
                                content: Text(
                                    '${plan.name} 플랜이 선택되었습니다.')));
                        context.pop();
                      },
                      style: FilledButton
                          .styleFrom(
                        backgroundColor:
                            isRecommended
                                ? SanggamTheme
                                    .primary
                                : SanggamTheme
                                    .surface,
                        foregroundColor:
                            isRecommended
                                ? SanggamTheme
                                    .background
                                : Colors
                                    .white,
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                      child: Text(
                          mode == 'upgrade'
                              ? '업그레이드'
                              : '선택하기'),
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
      if (i > 0 &&
          (str.length - i) % 3 == 0) {
        buf.write(',');
      }
      buf.write(str[i]);
    }
    return buf.toString();
  }

  static final _fallbackPlans = [
    _PlanData(
        name: '무료',
        description: '기본 기능 체험',
        price: '무료',
        period: '',
        features: [
          _Feature('기본 측정 1회/일', true),
          _Feature('측정 기록 7일 보관', true),
          _Feature('AI 코칭', false),
          _Feature('가족 공유', false),
          _Feature('원격 진료', false),
        ]),
    _PlanData(
        name: '베이직',
        description: '일상 건강 관리',
        price: '9,900원',
        period: '월',
        features: [
          _Feature('무제한 측정', true),
          _Feature('측정 기록 무제한 보관', true),
          _Feature('AI 건강 코칭', true),
          _Feature('데이터 내보내기', true),
          _Feature('가족 공유 (2명)', false),
        ]),
    _PlanData(
        name: '프로',
        description: '가족 건강 케어',
        price: '19,900원',
        period: '월',
        features: [
          _Feature('무제한 측정', true),
          _Feature('AI 고급 분석', true),
          _Feature('가족 공유 (5명)', true),
          _Feature('원격 진료 월 2회', true),
          _Feature('우선 고객 지원', true),
        ]),
    _PlanData(
        name: '클리닉',
        description: '전문가급 분석',
        price: '39,900원',
        period: '월',
        features: [
          _Feature(
              '무제한 측정 + 연구용 데이터', true),
          _Feature(
              'FHIR 의료 데이터 연동', true),
          _Feature('가족 공유 (무제한)', true),
          _Feature('원격 진료 무제한', true),
          _Feature('전담 건강 매니저', true),
        ]),
  ];
}

class _PlanData {
  final String name, description, price,
      period;
  final List<_Feature> features;
  const _PlanData({
    required this.name,
    required this.description,
    required this.price,
    required this.period,
    required this.features,
  });
}

class _Feature {
  final String label;
  final bool included;
  const _Feature(this.label, this.included);
}
