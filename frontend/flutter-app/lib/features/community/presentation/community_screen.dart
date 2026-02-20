import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';

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
        appBar: AppBar(
          title: const Text('커뮤니티'),
          centerTitle: true,
          bottom: const TabBar(
            isScrollable: true,
            tabs: [
              Tab(text: '전체'),
              Tab(text: '건강 팁'),
              Tab(text: '측정 후기'),
              Tab(text: 'Q&A'),
              Tab(text: '챌린지'),
              Tab(text: '연구'),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            _PostListTab(category: null),
            _PostListTab(category: PostCategory.healthTips),
            _PostListTab(category: PostCategory.reviews),
            _PostListTab(category: PostCategory.all),
            const _ChallengeTab(),
            const _ResearchTab(),
          ],
        ),
        floatingActionButton: FloatingActionButton(
          onPressed: () => context.push('/community/create'),
          child: const Icon(Icons.edit),
        ),
      ),
    );
  }

  void _showCreatePostDialog(BuildContext context, WidgetRef ref) {
    final titleCtrl = TextEditingController();
    final contentCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('게시글 작성'),
        content: _CreatePostForm(titleCtrl: titleCtrl, contentCtrl: contentCtrl),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('취소'),
          ),
          FilledButton(
            onPressed: () async {
              if (titleCtrl.text.isNotEmpty && contentCtrl.text.isNotEmpty) {
                try {
                  await ref.read(communityRepositoryProvider).createPost(
                        category: PostCategory.all,
                        title: titleCtrl.text,
                        content: contentCtrl.text,
                      );
                  ref.invalidate(communityPostsProvider);
                } catch (_) {}
                if (ctx.mounted) Navigator.pop(ctx);
              }
            },
            child: const Text('작성'),
          ),
        ],
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
    final theme = Theme.of(context);

    return postsAsync.when(
      data: (posts) {
        if (posts.isEmpty) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.forum_outlined, size: 48, color: theme.colorScheme.outline),
                const SizedBox(height: 12),
                Text(
                  '아직 게시글이 없습니다.\n첫 게시글을 작성해보세요!',
                  textAlign: TextAlign.center,
                  style: theme.textTheme.bodyLarge?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ),
          );
        }
        return RefreshIndicator(
          onRefresh: () async => ref.invalidate(communityPostsProvider(category)),
          child: ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: posts.length,
            itemBuilder: (context, index) {
              final post = posts[index];
              return _PostCard(post: post, ref: ref);
            },
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
    final theme = Theme.of(context);
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () => context.push('/community/post/${post.id}'),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  CircleAvatar(
                    radius: 16,
                    child: Text(post.authorName.isNotEmpty ? post.authorName[0] : '?'),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(post.authorName, style: theme.textTheme.labelMedium),
                        Text(
                          _formatDate(post.createdAt),
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.outline,
                          ),
                        ),
                      ],
                    ),
                  ),
                  _categoryChip(post.category, theme),
                ],
              ),
              const SizedBox(height: 12),
              Text(post.title, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 4),
              Text(
                post.content,
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
                style: theme.textTheme.bodyMedium,
              ),
              const Divider(height: 24),
              Row(
                children: [
                  InkWell(
                    onTap: () {
                      ref.read(communityRepositoryProvider).toggleLike(post.id);
                    },
                    child: Row(
                      children: [
                        Icon(
                          post.isLikedByMe ? Icons.favorite : Icons.favorite_border,
                          size: 18,
                          color: post.isLikedByMe ? Colors.red : null,
                        ),
                        const SizedBox(width: 4),
                        Text('${post.likeCount}'),
                      ],
                    ),
                  ),
                  const SizedBox(width: 16),
                  Icon(Icons.comment_outlined, size: 18, color: theme.colorScheme.outline),
                  const SizedBox(width: 4),
                  Text('${post.commentCount}'),
                  const Spacer(),
                  InkWell(
                    onTap: () {
                      showDialog(
                        context: context,
                        builder: (ctx) => AlertDialog(
                          title: const Text('게시글 신고'),
                          content: const Text('이 게시글을 부적절한 콘텐츠로 신고하시겠습니까?'),
                          actions: [
                            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
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
                    },
                    child: Icon(Icons.flag_outlined, size: 18, color: theme.colorScheme.outline),
                  ),
                  const SizedBox(width: 12),
                  Icon(
                    post.isBookmarkedByMe ? Icons.bookmark : Icons.bookmark_border,
                    size: 18,
                    color: theme.colorScheme.outline,
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _categoryChip(PostCategory cat, ThemeData theme) {
    final label = switch (cat) {
      PostCategory.healthTips => '건강 팁',
      PostCategory.reviews => '후기',
      PostCategory.challenge => '챌린지',
      _ => '',
    };
    if (label.isEmpty) return const SizedBox.shrink();
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: theme.colorScheme.primaryContainer,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Text(label, style: TextStyle(fontSize: 11, color: theme.colorScheme.onPrimaryContainer)),
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

/// 게시글 작성 폼 (데이터 첨부 + 익명 공유)
class _CreatePostForm extends StatefulWidget {
  const _CreatePostForm({required this.titleCtrl, required this.contentCtrl});
  final TextEditingController titleCtrl;
  final TextEditingController contentCtrl;

  @override
  State<_CreatePostForm> createState() => _CreatePostFormState();
}

class _CreatePostFormState extends State<_CreatePostForm> {
  bool _attachData = false;
  bool _anonymous = false;

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        TextField(
          controller: widget.titleCtrl,
          decoration: const InputDecoration(labelText: '제목'),
        ),
        const SizedBox(height: 8),
        TextField(
          controller: widget.contentCtrl,
          decoration: const InputDecoration(labelText: '내용'),
          maxLines: 4,
        ),
        const SizedBox(height: 8),
        // 이미지 첨부 버튼
        OutlinedButton.icon(
          onPressed: () {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('갤러리에서 이미지를 선택합니다.')),
            );
          },
          icon: const Icon(Icons.photo_library_outlined, size: 18),
          label: const Text('이미지 첨부'),
        ),
        const SizedBox(height: 8),
        CheckboxListTile(
          dense: true,
          contentPadding: EdgeInsets.zero,
          title: const Text('측정 데이터 첨부', style: TextStyle(fontSize: 13)),
          subtitle: const Text('최근 측정 결과를 게시글에 첨부', style: TextStyle(fontSize: 11)),
          value: _attachData,
          onChanged: (v) => setState(() => _attachData = v ?? false),
        ),
        CheckboxListTile(
          dense: true,
          contentPadding: EdgeInsets.zero,
          title: const Text('익명 공유', style: TextStyle(fontSize: 13)),
          subtitle: const Text('닉네임 대신 "익명"으로 표시', style: TextStyle(fontSize: 11)),
          value: _anonymous,
          onChanged: (v) => setState(() => _anonymous = v ?? false),
        ),
      ],
    );
  }
}

class _ChallengeTab extends ConsumerWidget {
  const _ChallengeTab();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final challengesAsync = ref.watch(healthChallengesProvider);
    final theme = Theme.of(context);

    return challengesAsync.when(
      data: (challenges) {
        if (challenges.isEmpty) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.emoji_events, size: 64, color: theme.colorScheme.primary),
                const SizedBox(height: 16),
                Text('건강 챌린지', style: theme.textTheme.headlineSmall),
                const SizedBox(height: 8),
                Text('곧 새로운 챌린지가 시작됩니다', style: theme.textTheme.bodyLarge),
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
            return Card(
              margin: const EdgeInsets.only(bottom: 12),
              child: InkWell(
                onTap: () => context.push('/community/challenge/${c.id}'),
                borderRadius: BorderRadius.circular(12),
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Icon(Icons.emoji_events, color: theme.colorScheme.primary, size: 28),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(c.title, style: theme.textTheme.titleMedium),
                                const SizedBox(height: 2),
                                Text('${c.participantCount}명 참여 중',
                                    style: theme.textTheme.bodySmall?.copyWith(
                                        color: theme.colorScheme.onSurfaceVariant)),
                              ],
                            ),
                          ),
                          if (!c.isJoined)
                            OutlinedButton(
                              onPressed: () => ref.read(communityRepositoryProvider).joinChallenge(c.id),
                              child: const Text('참가'),
                            ),
                        ],
                      ),
                      if (c.isJoined) ...[
                        const SizedBox(height: 12),
                        Row(
                          children: [
                            Text('내 진행률', style: theme.textTheme.bodySmall),
                            const Spacer(),
                            Text('${(progress * 100).toInt()}%',
                                style: theme.textTheme.bodySmall?.copyWith(
                                    fontWeight: FontWeight.bold,
                                    color: theme.colorScheme.primary)),
                          ],
                        ),
                        const SizedBox(height: 6),
                        ClipRRect(
                          borderRadius: BorderRadius.circular(4),
                          child: LinearProgressIndicator(
                            value: progress.clamp(0.0, 1.0),
                            minHeight: 6,
                            backgroundColor: theme.colorScheme.surfaceContainerHighest,
                            valueColor: AlwaysStoppedAnimation(theme.colorScheme.primary),
                          ),
                        ),
                      ],
                      // 리더보드 요약
                      const SizedBox(height: 10),
                      Row(
                        children: [
                          Icon(Icons.leaderboard, size: 16, color: theme.colorScheme.outline),
                          const SizedBox(width: 4),
                          Text('리더보드',
                              style: theme.textTheme.labelSmall?.copyWith(
                                  color: theme.colorScheme.outline)),
                          const Spacer(),
                          Icon(Icons.chevron_right, size: 16, color: theme.colorScheme.outline),
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
      error: (e, _) => Center(child: Text('불러오기 실패: $e')),
    );
  }
}

class _ResearchTab extends StatelessWidget {
  const _ResearchTab();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.science_outlined, size: 64, color: theme.colorScheme.primary),
          const SizedBox(height: 16),
          Text('건강 연구', style: theme.textTheme.headlineSmall),
          const SizedBox(height: 8),
          Text(
            '익명화된 건강 데이터를 활용한\n연구 프로젝트에 참여하세요.',
            textAlign: TextAlign.center,
            style: theme.textTheme.bodyLarge,
          ),
          const SizedBox(height: 24),
          FilledButton.icon(
            onPressed: () => context.push('/community/research'),
            icon: const Icon(Icons.arrow_forward),
            label: const Text('연구 프로젝트 보기'),
          ),
        ],
      ),
    );
  }
}
