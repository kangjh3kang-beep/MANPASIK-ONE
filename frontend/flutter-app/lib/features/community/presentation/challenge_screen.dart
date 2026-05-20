import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/leaderboard_widget.dart';
import 'package:manpasik/shared/widgets/confetti_overlay.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ChallengeScreen — Sanggam Orbit 건강 챌린지
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar 2x → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 4x → SanggamTheme.primary
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] theme.textTheme.* ~8x → 직접 TextStyle
// [Rule 4] theme.colorScheme.surfaceContainerHighest → SanggamTheme.surfaceVariant
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] bottom:12→16, borderRadius:12→16
// ───────────────────────────────────────────────────

/// 건강 챌린지 화면
class ChallengeScreen extends ConsumerWidget {
  const ChallengeScreen({super.key, this.challengeId});

  final String? challengeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final challengesAsync = ref.watch(challengesProvider);

    if (challengeId != null) {
      return _ChallengeDetailView(id: challengeId!, ref: ref);
    }

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
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
                      '건강 챌린지',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 목록
            Expanded(
              child: challengesAsync.when(
                data: (challenges) {
                  if (challenges.isEmpty) {
                    return const Center(
                      child: Text(
                        '진행 중인 챌린지가 없습니다.',
                        style: TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 16,
                        ),
                      ),
                    );
                  }
                  return ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: challenges.length,
                    itemBuilder: (context, index) {
                      final c = challenges[index];
                      final title =
                          c['title'] as String? ?? '챌린지';
                      final desc =
                          c['description'] as String? ?? '';
                      final participants =
                          c['participant_count'] as int? ?? 0;
                      final progress =
                          (c['progress'] as num?)?.toDouble() ??
                              0.0;

                      return Padding(
                        padding:
                            const EdgeInsets.only(bottom: 16),
                        child: SanggamContainer(
                          borderRadius: 16,
                          padding: const EdgeInsets.all(16),
                          onTap: () => context.push(
                              '/community/challenge/${c['id']}'),
                          child: Column(
                            crossAxisAlignment:
                                CrossAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  const Icon(
                                      Icons.emoji_events,
                                      color: SanggamTheme
                                          .primary),
                                  const SizedBox(width: 8),
                                  Expanded(
                                    child: Text(
                                      title,
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontSize: 16,
                                        fontWeight:
                                            FontWeight.bold,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              const SizedBox(height: 8),
                              Text(
                                desc,
                                style: const TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim,
                                  fontSize: 14,
                                ),
                              ),
                              const SizedBox(height: 16),
                              LinearProgressIndicator(
                                value: progress,
                                backgroundColor:
                                    SanggamTheme.surfaceVariant,
                                color: SanggamTheme.primary,
                              ),
                              const SizedBox(height: 8),
                              Row(
                                mainAxisAlignment:
                                    MainAxisAlignment
                                        .spaceBetween,
                                children: [
                                  Text(
                                    '${(progress * 100).toInt()}% 달성',
                                    style: const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize: 12,
                                    ),
                                  ),
                                  Text(
                                    '$participants명 참여 중',
                                    style: const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize: 12,
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  );
                },
                loading: () => const Center(
                    child: CircularProgressIndicator()),
                error: (_, __) => const Center(
                  child: Text(
                    '챌린지를 불러올 수 없습니다.',
                    style: TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: 16,
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _ChallengeDetailView extends StatefulWidget {
  const _ChallengeDetailView({required this.id, required this.ref});
  final String id;
  final WidgetRef ref;

  @override
  State<_ChallengeDetailView> createState() =>
      _ChallengeDetailViewState();
}

class _ChallengeDetailViewState extends State<_ChallengeDetailView> {
  final _confetti = ConfettiController();

  @override
  Widget build(BuildContext context) {
    return ConfettiOverlay(
      controller: _confetti,
      child: Scaffold(
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
                        '챌린지 상세',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              // 내용
              Expanded(
                child: ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    SanggamContainer(
                      borderRadius: 16,
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment:
                            CrossAxisAlignment.start,
                        children: [
                          const Text(
                            '7일 연속 측정 챌린지',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 20,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            '매일 1회 이상 바이오마커 측정을 완료하세요.',
                            style: TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 14,
                            ),
                          ),
                          const SizedBox(height: 16),
                          const LinearProgressIndicator(
                            value: 0.43,
                            color: SanggamTheme.primary,
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            '3/7일 완료',
                            style: TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 12,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 16),
                    LeaderboardWidget(
                      entries: List.generate(
                        5,
                        (i) => LeaderboardEntry(
                          userId: 'user-${i + 1}',
                          displayName: '사용자 ${i + 1}',
                          score: (7 - i) * 100,
                          rank: i + 1,
                          streak: 7 - i,
                        ),
                      ),
                      currentUserId: 'user-1',
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        bottomNavigationBar: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: FilledButton(
              onPressed: () {
                _confetti.fire();
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                      content: Text('챌린지에 참여했습니다!')),
                );
              },
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(48),
                backgroundColor: SanggamTheme.primary,
                foregroundColor: SanggamTheme.background,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16)),
              ),
              child: const Text('챌린지 참여하기'),
            ),
          ),
        ),
      ),
    );
  }
}

/// 챌린지 목록 Provider
final challengesProvider =
    FutureProvider<List<Map<String, dynamic>>>((ref) async {
  try {
    final client = ref.read(restClientProvider);
    final resp = await client.getChallenges();
    return (resp['challenges'] as List?)
            ?.cast<Map<String, dynamic>>() ??
        [];
  } catch (_) {
    return [
      {
        'id': '1',
        'title': '7일 연속 측정 챌린지',
        'description': '매일 1회 이상 바이오마커 측정을 완료하세요.',
        'participant_count': 156,
        'progress': 0.43,
      },
      {
        'id': '2',
        'title': '건강 식단 기록 챌린지',
        'description': '14일간 매 식사를 기록하고 AI 분석을 받으세요.',
        'participant_count': 89,
        'progress': 0.21,
      },
      {
        'id': '3',
        'title': '만보 걷기 챌린지',
        'description': '30일간 매일 10,000보 걷기를 달성하세요.',
        'participant_count': 234,
        'progress': 0.67,
      },
    ];
  }
});
