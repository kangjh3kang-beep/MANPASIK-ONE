import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';

/// AI 건강 코칭 화면
class AiCoachScreen extends ConsumerWidget {
  const AiCoachScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final insightAsync = ref.watch(todayInsightProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('AI 건강 코치'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.chat_bubble_outline),
            tooltip: 'AI 채팅',
            onPressed: () => context.push('/chat'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(todayInsightProvider);
          ref.invalidate(aiRecommendationsProvider);
        },
        child: ListView(
          padding: const EdgeInsets.all(24),
          children: [
            // 오늘의 건강 인사이트 카드
            insightAsync.when(
              data: (insight) => _InsightCard(insight: insight),
              loading: () => const Card(
                child: Padding(
                  padding: EdgeInsets.all(32),
                  child: Center(child: CircularProgressIndicator()),
                ),
              ),
              error: (_, __) => _InsightCard(
                insight: HealthInsight(
                  summary: '건강 데이터를 분석하는 중입니다.',
                  detail: '측정 데이터가 쌓이면 더 정확한 인사이트를 제공합니다.',
                  confidence: 0.0,
                  generatedAt: DateTime.now(),
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 건강 관리 영역
            Text('건강 관리 영역', style: theme.textTheme.titleMedium),
            const SizedBox(height: 12),
            _CoachCategoryTile(
              icon: Icons.restaurant,
              title: '식이 관리',
              subtitle: '맞춤형 식단 추천',
              category: 'diet',
            ),
            Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                leading: Icon(Icons.camera_alt, color: theme.colorScheme.primary),
                title: const Text('음식 사진 분석'),
                subtitle: const Text('사진 촬영으로 칼로리·영양소 분석'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => context.push('/coach/food'),
              ),
            ),
            _CoachCategoryTile(
              icon: Icons.fitness_center,
              title: '운동 관리',
              subtitle: '활동량 기반 운동 추천',
              category: 'exercise',
            ),
            Card(
              margin: const EdgeInsets.only(bottom: 8),
              child: ListTile(
                leading: Icon(Icons.play_circle_outline, color: theme.colorScheme.primary),
                title: const Text('운동 영상 가이드'),
                subtitle: const Text('전문 트레이너의 맞춤 운동 영상'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => context.push('/coach/exercise-video'),
              ),
            ),
            _CoachCategoryTile(
              icon: Icons.bedtime,
              title: '수면 관리',
              subtitle: '수면 패턴 분석',
              category: 'sleep',
            ),
            _CoachCategoryTile(
              icon: Icons.trending_up,
              title: '트렌드 분석',
              subtitle: '바이오마커 추세 분석',
              category: 'trend',
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
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.auto_awesome, color: theme.colorScheme.primary),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    '오늘의 건강 인사이트',
                    style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
                  ),
                ),
                if (insight.confidence > 0)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primaryContainer,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      '신뢰도 ${(insight.confidence * 100).toInt()}%',
                      style: TextStyle(fontSize: 11, color: theme.colorScheme.onPrimaryContainer),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 12),
            Text(insight.summary, style: theme.textTheme.bodyLarge),
            if (insight.detail.isNotEmpty) ...[
              const SizedBox(height: 8),
              Text(insight.detail, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ],
            const SizedBox(height: 8),
            Text(
              '※ 본 정보는 의료 조언이 아닙니다. 정확한 진단은 의료 전문가와 상담하세요.',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
            ),
          ],
        ),
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
    final recsAsync = ref.watch(aiRecommendationsProvider(category));
    final theme = Theme.of(context);

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ExpansionTile(
        leading: Icon(icon, color: theme.colorScheme.primary),
        title: Text(title),
        subtitle: Text(subtitle),
        children: [
          recsAsync.when(
            data: (recs) {
              if (recs.isEmpty) {
                return Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text(
                    '아직 추천 정보가 없습니다. 측정 데이터가 쌓이면 맞춤 추천을 제공합니다.',
                    style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.outline),
                  ),
                );
              }
              return Column(
                children: recs.map((r) {
                  return ListTile(
                    dense: true,
                    title: Text(r.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
                    subtitle: Text(r.description),
                    leading: _priorityIcon(r.priority, theme),
                  );
                }).toList(),
              );
            },
            loading: () => const Padding(
              padding: EdgeInsets.all(16),
              child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
            ),
            error: (_, __) => Padding(
              padding: const EdgeInsets.all(16),
              child: Text('추천을 불러올 수 없습니다', style: theme.textTheme.bodyMedium),
            ),
          ),
        ],
      ),
    );
  }

  Widget _priorityIcon(int priority, ThemeData theme) {
    if (priority >= 3) return const Icon(Icons.priority_high, color: Colors.red, size: 20);
    if (priority >= 2) return const Icon(Icons.arrow_upward, color: Colors.orange, size: 20);
    return Icon(Icons.lightbulb_outline, color: theme.colorScheme.primary, size: 20);
  }
}
