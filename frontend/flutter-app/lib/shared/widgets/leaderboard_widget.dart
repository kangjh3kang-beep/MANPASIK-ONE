import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// LeaderboardWidget — SanggamContainer 기반 챌린지 리더보드
//
// [Rule 4] app_theme.dart 미사용 import 제거
// [Rule 4] Card 2건 → SanggamContainer
// [Rule 4] theme.colorScheme.* 8건 → SanggamTheme 직접 참조
// [Rule 4] withOpacity 5건 → withValues(alpha:)
// [Rule 4] theme.textTheme.* → 직접 TextStyle
// [Rule 2] 12px 간격 → 16, vertical 10→8, borderRadius 12→16
// ───────────────────────────────────────────────────

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
    if (entries.isEmpty) {
      return SanggamContainer(
        borderRadius: 16,
        borderWidth: 0.8,
        borderColor: SanggamTheme.primary.withValues(alpha: 0.3),
        blurSigma: 8,
        jagaeOpacity: 0.03,
        padding: const EdgeInsets.all(24),
        child: const Center(
          child: Text(
            '아직 참가자가 없습니다',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
        ),
      );
    }

    return SanggamContainer(
      borderRadius: 16,
      borderWidth: 0.8,
      borderColor: SanggamTheme.primary.withValues(alpha: 0.3),
      blurSigma: 8,
      jagaeOpacity: 0.03,
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Header
          const Row(
            children: [
              Icon(Icons.leaderboard_rounded,
                  size: 20, color: SanggamTheme.primary),
              SizedBox(width: 8),
              Text(
                '리더보드',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // Entries
          ...entries.asMap().entries.map((e) {
            final rank = e.key + 1;
            final entry = e.value;
            final isMe = entry.userId == currentUserId;

            return Container(
              margin: const EdgeInsets.only(bottom: 8),
              padding:
                  const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: isMe
                    ? SanggamTheme.primary.withValues(alpha: 0.1)
                    : SanggamTheme.surface.withValues(alpha: 0.5),
                borderRadius: BorderRadius.circular(16),
                border: isMe
                    ? Border.all(
                        color:
                            SanggamTheme.primary.withValues(alpha: 0.5))
                    : null,
              ),
              child: Row(
                children: [
                  // 순위 / 메달
                  SizedBox(
                    width: 32,
                    child: _buildRankBadge(rank),
                  ),
                  const SizedBox(width: 16),
                  // 프로필
                  CircleAvatar(
                    radius: 16,
                    backgroundColor:
                        SanggamTheme.primary.withValues(alpha: 0.2),
                    child: Text(
                      entry.displayName.isNotEmpty
                          ? entry.displayName[0]
                          : '?',
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 12,
                        color: SanggamTheme.primary,
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  // 이름
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          entry.displayName + (isMe ? ' (나)' : ''),
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 14,
                            fontWeight:
                                isMe ? FontWeight.bold : FontWeight.normal,
                          ),
                        ),
                        if (entry.streak > 0)
                          Text(
                            '${entry.streak}일 연속',
                            style: const TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 11,
                            ),
                          ),
                      ],
                    ),
                  ),
                  // 점수
                  Text(
                    '${entry.score}점',
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                      color: SanggamTheme.primary,
                    ),
                  ),
                ],
              ),
            );
          }),
        ],
      ),
    );
  }

  Widget _buildRankBadge(int rank) {
    if (rank <= 3) {
      const colors = [
        SanggamTheme.primary, // 금
        Color(0xFFC0C0C0), // 은
        Color(0xFFCD7F32), // 동
      ];
      final c = colors[rank - 1];
      return Container(
        width: 28,
        height: 28,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: c.withValues(alpha: 0.2),
          border: Border.all(color: c, width: 2),
        ),
        child: Center(
          child: Text(
            '$rank',
            style: TextStyle(
              fontWeight: FontWeight.bold,
              fontSize: 12,
              color: c,
            ),
          ),
        ),
      );
    }
    return Center(
      child: Text(
        '$rank',
        style: const TextStyle(
          color: SanggamTheme.onSurfaceDim,
          fontSize: 14,
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
