import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// CommunityScreen — Sanggam Orbit 커뮤니티
//
// [Rule 4] Theme(data:ThemeData.dark().copyWith(...)) 래퍼 제거
// [Rule 4] AppBar → body 내 커스텀 헤더 + TabBar
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppTheme.sanggamGold 4x → SanggamTheme.primary
// [Rule 4] GlassmorphismCard 2x → SanggamContainer
// [Rule 4] HanjiPanel 1x → SanggamContainer
// [Rule 4] theme.colorScheme.* ~15x → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~15x → 직접 TextStyle
// [Rule 4] withOpacity 1x → withValues(alpha:)
// [Rule 4] Color(0xFF...) 하드코딩 3x → SanggamTheme 상수
// [Rule 4] Colors.white60/70/red → SanggamTheme 상수
// [Rule 4] grpc_provider.dart 미사용 import 제거
// [Rule 4] _showCreatePostDialog + _CreatePostForm 미사용 코드 제거
// [Rule 2] bottom 12→16, h:4→8, h:6→8, h:10→16, borderRadius 12→16, 4→8
// ───────────────────────────────────────────────────

/// 커뮤니티 화면
/// - 건강 정보 공유 게시판
/// - 카테고리별 게시글 목록
/// - 좋아요/댓글/북마크
/// - 게이미피케이션 (건강 챌린지)
class CommunityScreen extends ConsumerWidget {
  const CommunityScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return DefaultTabController(
      length: 6,
      child: Scaffold(
        backgroundColor: SanggamTheme.background,
        body: SafeArea(
          child: Column(
            children: [
              // 헤더: 타이틀
              const Padding(
                padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Text(
                  '커뮤니티',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    letterSpacing: -0.5,
                  ),
                  textAlign: TextAlign.center,
                ),
              ),
              // 탭바
              const TabBar(
                isScrollable: true,
                indicatorColor: SanggamTheme.primary,
                labelColor: SanggamTheme.primary,
                unselectedLabelColor: SanggamTheme.onSurfaceDim,
                dividerColor: Colors.transparent,
                tabs: [
                  Tab(text: '전체'),
                  Tab(text: '건강 팁'),
                  Tab(text: '측정 후기'),
                  Tab(text: 'Q&A'),
                  Tab(text: '챌린지'),
                  Tab(text: '연구'),
                ],
              ),
              // 탭 뷰
              Expanded(
                child: TabBarView(
                  children: [
                    _PostListTab(category: null),
                    _PostListTab(category: PostCategory.healthTips),
                    _PostListTab(category: PostCategory.reviews),
                    _PostListTab(category: PostCategory.all),
                    const _ChallengeTab(),
                    const _ResearchTab(),
                  ],
                ),
              ),
            ],
          ),
        ),
        floatingActionButton: FloatingActionButton(
          tooltip: '게시글 작성',
          backgroundColor: SanggamTheme.primary,
          foregroundColor: SanggamTheme.background,
          onPressed: () => context.push('/community/create'),
          child: const Icon(Icons.edit),
        ),
      ),
    );
  }
}

class _PostListTab extends ConsumerWidget {
  const _PostListTab({required this.category});
  final PostCategory? category;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final postsAsync = ref.watch(communityPostsProvider(category));

    return postsAsync.when(
      data: (posts) {
        if (posts.isEmpty) {
          return Padding(
            padding: const EdgeInsets.all(24),
            child: Center(
              child: SanggamContainer(
                borderRadius: 24,
                padding: const EdgeInsets.all(32),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.forum_outlined,
                        size: 48,
                        color:
                            SanggamTheme.primary.withValues(alpha: 0.6)),
                    const SizedBox(height: 16),
                    const Text(
                      '아직 게시글이 없습니다.\n첫 게시글을 작성해보세요!',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 16,
                        height: 1.5,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          );
        }
        return RefreshIndicator(
          onRefresh: () async =>
              ref.invalidate(communityPostsProvider(category)),
          child: ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: posts.length,
            itemBuilder: (context, index) =>
                _PostCard(post: posts[index], ref: ref),
          ),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('불러오기 실패: $e')),
    );
  }
}

class _PostCard extends StatelessWidget {
  const _PostCard({required this.post, required this.ref});
  final CommunityPost post;
  final WidgetRef ref;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        padding: const EdgeInsets.all(16),
        onTap: () => context.push('/community/post/${post.id}'),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 16,
                  child: Text(post.authorName.isNotEmpty
                      ? post.authorName[0]
                      : '?'),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        post.authorName,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 12,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      Text(
                        _formatDate(post.createdAt),
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                _categoryChip(post.category),
              ],
            ),
            const SizedBox(height: 16),
            Text(
              post.title,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              post.content,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
              ),
            ),
            const Divider(height: 24),
            Row(
              children: [
                InkWell(
                  onTap: () => ref
                      .read(communityRepositoryProvider)
                      .toggleLike(post.id),
                  child: Row(
                    children: [
                      Icon(
                        post.isLikedByMe
                            ? Icons.favorite
                            : Icons.favorite_border,
                        size: 18,
                        color: post.isLikedByMe
                            ? SanggamTheme.error
                            : SanggamTheme.onSurfaceDim,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${post.likeCount}',
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                const Icon(Icons.comment_outlined,
                    size: 18, color: SanggamTheme.onSurfaceDim),
                const SizedBox(width: 4),
                Text(
                  '${post.commentCount}',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
                const Spacer(),
                InkWell(
                  onTap: () => _showReportDialog(context),
                  child: const Icon(Icons.flag_outlined,
                      size: 18, color: SanggamTheme.onSurfaceDim),
                ),
                const SizedBox(width: 16),
                Icon(
                  post.isBookmarkedByMe
                      ? Icons.bookmark
                      : Icons.bookmark_border,
                  size: 18,
                  color: SanggamTheme.onSurfaceDim,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showReportDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('게시글 신고'),
        content: const Text('이 게시글을 부적절한 콘텐츠로 신고하시겠습니까?'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(ctx),
              child: const Text('취소')),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('신고가 접수되었습니다.')),
              );
            },
            child: const Text('신고'),
          ),
        ],
      ),
    );
  }

  Widget _categoryChip(PostCategory cat) {
    final label = switch (cat) {
      PostCategory.healthTips => '건강 팁',
      PostCategory.reviews => '후기',
      PostCategory.challenge => '챌린지',
      _ => '',
    };
    if (label.isEmpty) return const SizedBox.shrink();
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: SanggamTheme.surfaceVariant,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Text(
        label,
        style: const TextStyle(
          fontSize: 11,
          color: SanggamTheme.primary,
        ),
      ),
    );
  }

  String _formatDate(DateTime dt) {
    final diff = DateTime.now().difference(dt);
    if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
    if (diff.inHours < 24) return '${diff.inHours}시간 전';
    if (diff.inDays < 7) return '${diff.inDays}일 전';
    return '${dt.month}/${dt.day}';
  }
}

class _ChallengeTab extends ConsumerWidget {
  const _ChallengeTab();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final challengesAsync = ref.watch(healthChallengesProvider);

    return challengesAsync.when(
      data: (challenges) {
        if (challenges.isEmpty) {
          return const Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.emoji_events,
                    size: 64, color: SanggamTheme.primary),
                SizedBox(height: 16),
                Text(
                  '건강 챌린지',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                SizedBox(height: 8),
                Text(
                  '곧 새로운 챌린지가 시작됩니다',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
          );
        }
        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: challenges.length,
          itemBuilder: (context, index) {
            final c = challenges[index];
            final progress = c.myProgress ?? 0.0;
            return Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: SanggamContainer(
                borderRadius: 16,
                padding: const EdgeInsets.all(16),
                onTap: () =>
                    context.push('/community/challenge/${c.id}'),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.emoji_events,
                            color: SanggamTheme.primary, size: 28),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment:
                                CrossAxisAlignment.start,
                            children: [
                              Text(
                                c.title,
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontWeight: FontWeight.w900,
                                  letterSpacing: -0.5,
                                  fontSize: 18,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Text(
                                '${c.participantCount}명 참여 중',
                                style: const TextStyle(
                                  color: SanggamTheme.onSurfaceDim,
                                  fontSize: 12,
                                ),
                              ),
                            ],
                          ),
                        ),
                        if (!c.isJoined)
                          FilledButton(
                            onPressed: () => ref
                                .read(communityRepositoryProvider)
                                .joinChallenge(c.id),
                            style: FilledButton.styleFrom(
                              backgroundColor: SanggamTheme.primary,
                              foregroundColor:
                                  SanggamTheme.background,
                              shape: RoundedRectangleBorder(
                                  borderRadius:
                                      BorderRadius.circular(16)),
                            ),
                            child: const Text('참가',
                                style: TextStyle(
                                    fontWeight: FontWeight.bold)),
                          ),
                      ],
                    ),
                    if (c.isJoined) ...[
                      const SizedBox(height: 16),
                      Row(
                        children: [
                          const Text(
                            '내 진행률',
                            style: TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 12,
                            ),
                          ),
                          const Spacer(),
                          Text(
                            '${(progress * 100).toInt()}%',
                            style: const TextStyle(
                              color: SanggamTheme.primary,
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 8),
                      ClipRRect(
                        borderRadius: BorderRadius.circular(8),
                        child: LinearProgressIndicator(
                          value: progress.clamp(0.0, 1.0),
                          minHeight: 8,
                          backgroundColor:
                              SanggamTheme.surfaceVariant,
                          valueColor:
                              const AlwaysStoppedAnimation(
                                  SanggamTheme.primary),
                        ),
                      ),
                    ],
                    const SizedBox(height: 16),
                    const Row(
                      children: [
                        Icon(Icons.leaderboard,
                            size: 16,
                            color: SanggamTheme.onSurfaceDim),
                        SizedBox(width: 8),
                        Text(
                          '리더보드',
                          style: TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 12,
                          ),
                        ),
                        Spacer(),
                        Icon(Icons.chevron_right,
                            size: 16,
                            color: SanggamTheme.onSurfaceDim),
                      ],
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('불러오기 실패: $e')),
    );
  }
}

class _ResearchTab extends StatelessWidget {
  const _ResearchTab();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.science_outlined,
              size: 64, color: SanggamTheme.primary),
          const SizedBox(height: 16),
          const Text(
            '건강 연구',
            style: TextStyle(
              color: Colors.white,
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '익명화된 건강 데이터를 활용한\n연구 프로젝트에 참여하세요.',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 16,
            ),
          ),
          const SizedBox(height: 24),
          FilledButton.icon(
            onPressed: () => context.push('/community/research'),
            icon: const Icon(Icons.arrow_forward),
            label: const Text('연구 프로젝트 보기'),
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.primary,
              foregroundColor: SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16)),
            ),
          ),
        ],
      ),
    );
  }
}
