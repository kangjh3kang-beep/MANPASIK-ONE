import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/leaderboard_widget.dart';

/// 건강 챌린지 화면
class ChallengeScreen extends ConsumerWidget {
  const ChallengeScreen({super.key, this.challengeId});

  final String? challengeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final challengesAsync = ref.watch(challengesProvider);

    if (challengeId != null) {
      return _ChallengeDetailView(id: challengeId!, ref: ref);
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('건강 챌린지'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: challengesAsync.when(
        data: (challenges) {
          if (challenges.isEmpty) {
            return const Center(child: Text('진행 중인 챌린지가 없습니다.'));
          }
          return ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: challenges.length,
            itemBuilder: (context, index) {
              final c = challenges[index];
              final title = c['title'] as String? ?? '챌린지';
              final desc = c['description'] as String? ?? '';
              final participants = c['participant_count'] as int? ?? 0;
              final progress = (c['progress'] as num?)?.toDouble() ?? 0.0;

              return Card(
                margin: const EdgeInsets.only(bottom: 12),
                child: InkWell(
                  onTap: () => context.push('/community/challenge/${c['id']}'),
                  borderRadius: BorderRadius.circular(12),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            const Icon(Icons.emoji_events, color: AppTheme.sanggamGold),
                            const SizedBox(width: 8),
                            Expanded(child: Text(title, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold))),
                          ],
                        ),
                        const SizedBox(height: 8),
                        Text(desc, style: theme.textTheme.bodyMedium),
                        const SizedBox(height: 12),
                        LinearProgressIndicator(value: progress, backgroundColor: theme.colorScheme.surfaceContainerHighest, color: AppTheme.sanggamGold),
                        const SizedBox(height: 8),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text('${(progress * 100).toInt()}% 달성', style: theme.textTheme.bodySmall),
                            Text('$participants명 참여 중', style: theme.textTheme.bodySmall),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              );
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (_, __) => const Center(child: Text('챌린지를 불러올 수 없습니다.')),
      ),
    );
  }
}

class _ChallengeDetailView extends StatelessWidget {
  const _ChallengeDetailView({required this.id, required this.ref});
  final String id;
  final WidgetRef ref;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('챌린지 상세'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('7일 연속 측정 챌린지', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  Text('매일 1회 이상 바이오마커 측정을 완료하세요.', style: theme.textTheme.bodyMedium),
                  const SizedBox(height: 16),
                  const LinearProgressIndicator(value: 0.43, color: AppTheme.sanggamGold),
                  const SizedBox(height: 8),
                  Text('3/7일 완료', style: theme.textTheme.bodySmall),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          const SizedBox(height: 8),
          LeaderboardWidget(
            entries: List.generate(5, (i) => LeaderboardEntry(
              userId: 'user-${i + 1}',
              displayName: '사용자 ${i + 1}',
              score: (7 - i) * 100,
              rank: i + 1,
              streak: 7 - i,
            )),
            currentUserId: 'user-1',
          ),
        ],
      ),
      bottomNavigationBar: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: FilledButton(
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('챌린지에 참여했습니다!')),
              );
            },
            style: FilledButton.styleFrom(
              minimumSize: const Size.fromHeight(48),
              backgroundColor: AppTheme.sanggamGold,
            ),
            child: const Text('챌린지 참여하기'),
          ),
        ),
      ),
    );
  }
}

/// 챌린지 목록 Provider
final challengesProvider = FutureProvider<List<Map<String, dynamic>>>((ref) async {
  try {
    final client = ref.read(restClientProvider);
    final resp = await client.getChallenges();
    return (resp['challenges'] as List?)?.cast<Map<String, dynamic>>() ?? [];
  } catch (_) {
    return [
      {'id': '1', 'title': '7일 연속 측정 챌린지', 'description': '매일 1회 이상 바이오마커 측정을 완료하세요.', 'participant_count': 156, 'progress': 0.43},
      {'id': '2', 'title': '건강 식단 기록 챌린지', 'description': '14일간 매 식사를 기록하고 AI 분석을 받으세요.', 'participant_count': 89, 'progress': 0.21},
      {'id': '3', 'title': '만보 걷기 챌린지', 'description': '30일간 매일 10,000보 걷기를 달성하세요.', 'participant_count': 234, 'progress': 0.67},
    ];
  }
});
