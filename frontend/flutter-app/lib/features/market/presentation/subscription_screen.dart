import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// SubscriptionScreen — Sanggam Orbit 구독 관리
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData params 4x 제거
// [Rule 4] theme.textTheme ~10x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold ~6x → SanggamTheme.primary
// [Rule 4] Card ~4x → SanggamContainer
// [Rule 4] Chip → Container+BoxDecoration
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] AlertDialog 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16
// ───────────────────────────────────────────────────

/// 구독 관리 화면
class SubscriptionScreen
    extends ConsumerWidget {
  const SubscriptionScreen({super.key});

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final plansAsync =
        ref.watch(subscriptionPlansProvider);
    final subAsync =
        ref.watch(subscriptionInfoProvider);

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
                      '구독 관리',
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
                    // 현재 구독 상태
                    _buildCurrentPlan(
                        context, subAsync),
                    const SizedBox(
                        height: 24),

                    const Text('구독 플랜 비교',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 16),

                    // 플랜 비교 카드들
                    plansAsync.when(
                      data: (plans) {
                        if (plans
                            .isEmpty) {
                          return _buildFallbackPlans(
                              context);
                        }
                        return Column(
                          children: plans
                              .map((p) =>
                                  _buildPlanCard(
                                      context,
                                      p))
                              .toList(),
                        );
                      },
                      loading: () =>
                          const Center(
                              child: CircularProgressIndicator(
                                  color: SanggamTheme
                                      .primary)),
                      error: (_, __) =>
                          _buildFallbackPlans(
                              context),
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

  Widget _buildCurrentPlan(
      BuildContext context,
      AsyncValue<dynamic> subAsync) {
    return SanggamContainer(
      borderRadius: 16,
      borderColor: SanggamTheme.primary,
      jagaeOpacity: 0.15,
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          const Text('현재 구독',
              style: TextStyle(
                color:
                    SanggamTheme.primary,
                fontSize: 13,
                fontWeight: FontWeight.w600,
              )),
          const SizedBox(height: 8),
          subAsync.when(
            data: (sub) {
              if (sub == null) {
                return const Row(
                  children: [
                    Icon(
                        Icons
                            .card_membership,
                        color: SanggamTheme
                            .primary),
                    SizedBox(width: 16),
                    Text('무료 플랜',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                  ],
                );
              }
              final tierName = [
                '무료',
                '베이직',
                '프로',
                '클리닉'
              ][sub.tier.clamp(0, 3)];
              return Column(
                crossAxisAlignment:
                    CrossAxisAlignment
                        .start,
                children: [
                  Row(
                    children: [
                      const Icon(
                          Icons
                              .card_membership,
                          color: SanggamTheme
                              .primary,
                          size: 40),
                      const SizedBox(
                          width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment:
                              CrossAxisAlignment
                                  .start,
                          children: [
                            Text(
                                '$tierName 플랜',
                                style:
                                    const TextStyle(
                                  color: Colors
                                      .white,
                                  fontSize:
                                      16,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                )),
                            const Text(
                                '활성 상태',
                                style:
                                    TextStyle(
                                  color: SanggamTheme
                                      .jagaeCyan,
                                  fontSize:
                                      12,
                                )),
                          ],
                        ),
                      ),
                    ],
                  ),
                  if (sub.tier > 0) ...[
                    const SizedBox(
                        height: 16),
                    SizedBox(
                      width:
                          double.infinity,
                      child: OutlinedButton(
                        onPressed: () =>
                            context.push(
                                '/market/subscription/cancel'),
                        style:
                            OutlinedButton
                                .styleFrom(
                          foregroundColor:
                              SanggamTheme
                                  .error,
                          side: const BorderSide(
                              color:
                                  SanggamTheme
                                      .error),
                          shape:
                              RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                          ),
                        ),
                        child: const Text(
                            '구독 해지'),
                      ),
                    ),
                  ],
                ],
              );
            },
            loading: () =>
                const CircularProgressIndicator(
                    strokeWidth: 2,
                    color: SanggamTheme
                        .primary),
            error: (_, __) => const Text(
                '구독 정보를 불러올 수 없습니다.',
                style: TextStyle(
                    color: Colors.white)),
          ),
        ],
      ),
    );
  }

  Widget _buildFallbackPlans(
      BuildContext context) {
    final plans = [
      _PlanInfo(
          '무료',
          '₩0/월',
          [
            '기본 측정 1회/일',
            '리더기 1대',
            '기본 AI 분석'
          ],
          false),
      _PlanInfo(
          '베이직',
          '₩9,900/월',
          [
            '무제한 측정',
            '리더기 2대',
            'AI 건강 코칭',
            '카트리지 도감'
          ],
          false),
      _PlanInfo(
          '프로',
          '₩29,900/월',
          [
            '모든 베이직 기능',
            '리더기 5대',
            '가족 공유 (5명)',
            '비대면 진료',
            '데이터 내보내기 (FHIR)'
          ],
          true),
      _PlanInfo(
          '클리닉',
          '₩99,000/월',
          [
            '모든 프로 기능',
            '리더기 10대',
            '의료기관급 분석',
            'MFA 필수 보안',
            '전용 고객 지원'
          ],
          false),
    ];
    return Column(
      children: plans
          .map((p) =>
              _buildFallbackPlanCard(
                  context, p))
          .toList(),
    );
  }

  Widget _buildFallbackPlanCard(
      BuildContext context,
      _PlanInfo plan) {
    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        borderColor: plan.recommended
            ? SanggamTheme.primary
            : SanggamTheme.surfaceVariant,
        borderWidth:
            plan.recommended ? 2.0 : 1.5,
        child: Column(
          crossAxisAlignment:
              CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment:
                  MainAxisAlignment
                      .spaceBetween,
              children: [
                Text(plan.name,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight:
                          FontWeight.bold,
                    )),
                if (plan.recommended)
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
            ),
            Text(plan.price,
                style: const TextStyle(
                  color: SanggamTheme
                      .primary,
                  fontSize: 20,
                  fontWeight:
                      FontWeight.bold,
                )),
            const SizedBox(height: 8),
            ...plan.features.map(
                (f) => Padding(
                      padding:
                          const EdgeInsets
                              .symmetric(
                              vertical: 2),
                      child: Row(
                        children: [
                          const Icon(
                              Icons.check,
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
            if (plan.name != '무료')
              SizedBox(
                width: double.infinity,
                child: OutlinedButton(
                  onPressed: () =>
                      context.push(
                          '/market/subscription/upgrade'),
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
                  child: Text(
                      '${plan.name} 구독하기'),
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildPlanCard(
      BuildContext context,
      SubscriptionPlan plan) {
    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        child: Column(
          crossAxisAlignment:
              CrossAxisAlignment.start,
          children: [
            Text(plan.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight:
                      FontWeight.bold,
                )),
            Text(
              '₩${plan.monthlyPrice.toString().replaceAllMapped(RegExp(r'(\d)(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}/월',
              style: const TextStyle(
                color:
                    SanggamTheme.primary,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
                '${plan.cartridgesPerMonth}개 카트리지/월 | ${plan.discountPercent}% 할인',
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 13,
                )),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton(
                onPressed: () =>
                    _subscribePlan(
                        context, plan),
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
                            .circular(16),
                  ),
                ),
                child: Text(
                    '${plan.name} 구독하기'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _subscribePlan(
      BuildContext context,
      SubscriptionPlan plan) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor:
            SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: Text('${plan.name} 구독',
            style: const TextStyle(
                color: Colors.white)),
        content: Text(
          '월 ₩${plan.monthlyPrice.toString().replaceAllMapped(RegExp(r"(\d)(?=(\d{3})+(?!\d))"), (m) => "${m[1]},")}으로 '
          '${plan.name} 플랜을 구독하시겠습니까?',
          style: const TextStyle(
              color: SanggamTheme
                  .onSurfaceDim),
        ),
        actions: [
          TextButton(
            onPressed: () =>
                Navigator.pop(ctx),
            style: TextButton.styleFrom(
                foregroundColor:
                    SanggamTheme
                        .onSurfaceDim),
            child: const Text('취소'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              context.push(
                  '/market/checkout');
            },
            style:
                FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.primary,
              foregroundColor:
                  SanggamTheme.background,
              shape:
                  RoundedRectangleBorder(
                borderRadius:
                    BorderRadius.circular(
                        16),
              ),
            ),
            child:
                const Text('결제 진행'),
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
  const _PlanInfo(this.name, this.price,
      this.features, this.recommended);
}
