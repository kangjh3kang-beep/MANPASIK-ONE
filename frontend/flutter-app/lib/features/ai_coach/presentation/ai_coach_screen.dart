import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AiCoachScreen — Sanggam Orbit AI 건강 코칭
//
// [Rule 4] +sanggam_theme.dart, +sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 3x 제거
// [Rule 4] ThemeData 파라미터 1x 제거
// [Rule 4] theme.textTheme.* ~10x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~10x → SanggamTheme 상수
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Card 5x → SanggamContainer / ListTile(tileColor)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] padding:20→24, h:12→16, borderRadius:12→16
// ───────────────────────────────────────────────────

/// AI 건강 코칭 화면
class AiCoachScreen extends ConsumerWidget {
  const AiCoachScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final insightAsync =
        ref.watch(todayInsightProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      'AI 건강 코치',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(
                        Icons.chat_bubble_outline,
                        color: Colors.white),
                    tooltip: 'AI 채팅',
                    onPressed: () =>
                        context.push('/chat'),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: RefreshIndicator(
                onRefresh: () async {
                  ref.invalidate(todayInsightProvider);
                  ref.invalidate(
                      aiRecommendationsProvider);
                },
                child: ListView(
                  padding: const EdgeInsets.all(24),
                  children: [
                    // 오늘의 건강 인사이트 카드
                    insightAsync.when(
                      data: (insight) => _InsightCard(
                          insight: insight),
                      loading: () =>
                          const SanggamContainer(
                        borderRadius: 16,
                        padding: EdgeInsets.all(32),
                        child: Center(
                            child:
                                CircularProgressIndicator()),
                      ),
                      error: (_, __) => _InsightCard(
                        insight: HealthInsight(
                          summary:
                              '건강 데이터를 분석하는 중입니다.',
                          detail:
                              '측정 데이터가 쌓이면 더 정확한 인사이트를 제공합니다.',
                          confidence: 0.0,
                          generatedAt: DateTime.now(),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),

                    // 건강 관리 영역
                    const Text(
                      '건강 관리 영역',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    const _CoachCategoryTile(
                      icon: Icons.restaurant,
                      title: '식이 관리',
                      subtitle: '맞춤형 식단 추천',
                      category: 'diet',
                    ),
                    Padding(
                      padding: const EdgeInsets.only(
                          bottom: 8),
                      child: ListTile(
                        tileColor: SanggamTheme.surface,
                        shape: RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius.circular(16),
                        ),
                        leading: const Icon(
                            Icons.camera_alt,
                            color:
                                SanggamTheme.primary),
                        title:
                            const Text('음식 사진 분석'),
                        subtitle: const Text(
                            '사진 촬영으로 칼로리·영양소 분석'),
                        trailing: const Icon(
                            Icons.chevron_right,
                            color: SanggamTheme
                                .onSurfaceDim),
                        onTap: () => context
                            .push('/coach/food'),
                      ),
                    ),
                    const _CoachCategoryTile(
                      icon: Icons.fitness_center,
                      title: '운동 관리',
                      subtitle: '활동량 기반 운동 추천',
                      category: 'exercise',
                    ),
                    Padding(
                      padding: const EdgeInsets.only(
                          bottom: 8),
                      child: ListTile(
                        tileColor: SanggamTheme.surface,
                        shape: RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius.circular(16),
                        ),
                        leading: const Icon(
                            Icons.play_circle_outline,
                            color:
                                SanggamTheme.primary),
                        title:
                            const Text('운동 영상 가이드'),
                        subtitle: const Text(
                            '전문 트레이너의 맞춤 운동 영상'),
                        trailing: const Icon(
                            Icons.chevron_right,
                            color: SanggamTheme
                                .onSurfaceDim),
                        onTap: () => context.push(
                            '/coach/exercise-video'),
                      ),
                    ),
                    const _CoachCategoryTile(
                      icon: Icons.bedtime,
                      title: '수면 관리',
                      subtitle: '수면 패턴 분석',
                      category: 'sleep',
                    ),
                    const _CoachCategoryTile(
                      icon: Icons.trending_up,
                      title: '트렌드 분석',
                      subtitle: '바이오마커 추세 분석',
                      category: 'trend',
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
}

class _InsightCard extends StatelessWidget {
  const _InsightCard({required this.insight});
  final HealthInsight insight;

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.auto_awesome,
                  color: SanggamTheme.primary),
              const SizedBox(width: 8),
              const Expanded(
                child: Text(
                  '오늘의 건강 인사이트',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              if (insight.confidence > 0)
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: SanggamTheme.primary
                        .withValues(alpha: 0.15),
                    borderRadius:
                        BorderRadius.circular(16),
                  ),
                  child: Text(
                    '신뢰도 ${(insight.confidence * 100).toInt()}%',
                    style: const TextStyle(
                      fontSize: 11,
                      color: SanggamTheme.primary,
                    ),
                  ),
                ),
            ],
          ),
          const SizedBox(height: 16),
          Text(
            insight.summary,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
            ),
          ),
          if (insight.detail.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text(
              insight.detail,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
              ),
            ),
          ],
          const SizedBox(height: 8),
          const Text(
            '※ 본 정보는 의료 조언이 아닙니다. 정확한 진단은 의료 전문가와 상담하세요.',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }
}

class _CoachCategoryTile extends ConsumerWidget {
  const _CoachCategoryTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.category,
  });

  final IconData icon;
  final String title;
  final String subtitle;
  final String category;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final recsAsync =
        ref.watch(aiRecommendationsProvider(category));

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: ExpansionTile(
          leading: Icon(icon,
              color: SanggamTheme.primary),
          title: Text(
            title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
            ),
          ),
          subtitle: Text(
            subtitle,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
          children: [
            recsAsync.when(
              data: (recs) {
                if (recs.isEmpty) {
                  return const Padding(
                    padding: EdgeInsets.all(16),
                    child: Text(
                      '아직 추천 정보가 없습니다. 측정 데이터가 쌓이면 맞춤 추천을 제공합니다.',
                      style: TextStyle(
                        color:
                            SanggamTheme.onSurfaceDim,
                        fontSize: 14,
                      ),
                    ),
                  );
                }
                return Column(
                  children: recs.map((r) {
                    return ListTile(
                      dense: true,
                      title: Text(
                        r.title,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      subtitle: Text(r.description),
                      leading:
                          _priorityIcon(r.priority),
                    );
                  }).toList(),
                );
              },
              loading: () => const Padding(
                padding: EdgeInsets.all(16),
                child: Center(
                    child:
                        CircularProgressIndicator(
                            strokeWidth: 2)),
              ),
              error: (_, __) => const Padding(
                padding: EdgeInsets.all(16),
                child: Text(
                  '추천을 불러올 수 없습니다',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _priorityIcon(int priority) {
    if (priority >= 3) {
      return const Icon(Icons.priority_high,
          color: SanggamTheme.error, size: 20);
    }
    if (priority >= 2) {
      return const Icon(Icons.arrow_upward,
          color: SanggamTheme.primary, size: 20);
    }
    return const Icon(Icons.lightbulb_outline,
        color: SanggamTheme.primary, size: 20);
  }
}
