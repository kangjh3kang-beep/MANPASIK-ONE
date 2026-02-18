import 'package:flutter/material.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 챌린지 리더보드 위젯 (C8)
///
/// 건강 챌린지 참가자 순위를 메달 아이콘과 함께 표시합니다.
class LeaderboardWidget extends StatelessWidget {
  const LeaderboardWidget({
    super.key,
    required this.entries,
    this.currentUserId,
  });

  final List<LeaderboardEntry> entries;
  final String? currentUserId;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (entries.isEmpty) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Center(
            child: Text(
              '아직 참가자가 없습니다',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ),
        ),
      );
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.leaderboard_rounded, size: 20),
                const SizedBox(width: 8),
                Text(
                  '리더보드',
                  style: theme.textTheme.titleSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            ...entries.asMap().entries.map((e) {
              final rank = e.key + 1;
              final entry = e.value;
              final isMe = entry.userId == currentUserId;

              return Container(
                margin: const EdgeInsets.only(bottom: 8),
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                decoration: BoxDecoration(
                  color: isMe
                      ? theme.colorScheme.primaryContainer.withOpacity(0.3)
                      : theme.colorScheme.surfaceContainerHighest
                          .withOpacity(0.5),
                  borderRadius: BorderRadius.circular(12),
                  border: isMe
                      ? Border.all(
                          color: theme.colorScheme.primary.withOpacity(0.5))
                      : null,
                ),
                child: Row(
                  children: [
                    // 순위 / 메달
                    SizedBox(
                      width: 32,
                      child: _buildRankBadge(rank, theme),
                    ),
                    const SizedBox(width: 12),
                    // 프로필
                    CircleAvatar(
                      radius: 16,
                      backgroundColor:
                          theme.colorScheme.primary.withOpacity(0.2),
                      child: Text(
                        entry.displayName.isNotEmpty
                            ? entry.displayName[0]
                            : '?',
                        style: theme.textTheme.bodySmall?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: theme.colorScheme.primary,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    // 이름
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            entry.displayName + (isMe ? ' (나)' : ''),
                            style: theme.textTheme.bodyMedium?.copyWith(
                              fontWeight:
                                  isMe ? FontWeight.bold : FontWeight.normal,
                            ),
                          ),
                          if (entry.streak > 0)
                            Text(
                              '${entry.streak}일 연속',
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: theme.colorScheme.onSurfaceVariant,
                                fontSize: 11,
                              ),
                            ),
                        ],
                      ),
                    ),
                    // 점수
                    Text(
                      '${entry.score}점',
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: theme.colorScheme.primary,
                      ),
                    ),
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }

  Widget _buildRankBadge(int rank, ThemeData theme) {
    if (rank <= 3) {
      final colors = [
        const Color(0xFFFFD700),
        const Color(0xFFC0C0C0),
        const Color(0xFFCD7F32),
      ];
      return Container(
        width: 28,
        height: 28,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: colors[rank - 1].withOpacity(0.2),
          border: Border.all(color: colors[rank - 1], width: 2),
        ),
        child: Center(
          child: Text(
            '$rank',
            style: theme.textTheme.bodySmall?.copyWith(
              fontWeight: FontWeight.bold,
              color: colors[rank - 1],
            ),
          ),
        ),
      );
    }
    return Center(
      child: Text(
        '$rank',
        style: theme.textTheme.bodyMedium?.copyWith(
          color: theme.colorScheme.onSurfaceVariant,
        ),
      ),
    );
  }
}

class LeaderboardEntry {
  final String userId;
  final String displayName;
  final int score;
  final int rank;
  final int streak;

  const LeaderboardEntry({
    required this.userId,
    required this.displayName,
    required this.score,
    required this.rank,
    this.streak = 0,
  });
}
