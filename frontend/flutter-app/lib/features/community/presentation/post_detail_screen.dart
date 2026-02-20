import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/shared/widgets/cached_image.dart';

/// 커뮤니티 게시글 상세 화면
class PostDetailScreen extends ConsumerStatefulWidget {
  const PostDetailScreen({super.key, required this.postId});

  final String postId;

  @override
  ConsumerState<PostDetailScreen> createState() => _PostDetailScreenState();
}

class _PostDetailScreenState extends ConsumerState<PostDetailScreen> {
  final _commentController = TextEditingController();
  bool _isLiked = false;
  bool _isBookmarked = false;
  bool _isSubmitting = false;
  late Future<Map<String, dynamic>> _postFuture;
  List<Comment> _comments = [];
  bool _commentsLoading = true;

  @override
  void initState() {
    super.initState();
    _postFuture = ref.read(restClientProvider).getPost(widget.postId);
    _loadComments();
  }

  Future<void> _loadComments() async {
    try {
      final comments = await ref.read(communityRepositoryProvider).getComments(widget.postId);
      if (!mounted) return;
      setState(() {
        _comments = comments;
        _commentsLoading = false;
      });
    } catch (_) {
      if (!mounted) return;
      setState(() => _commentsLoading = false);
    }
  }

  Future<void> _submitComment() async {
    final text = _commentController.text.trim();
    if (text.isEmpty || _isSubmitting) return;
    setState(() => _isSubmitting = true);
    try {
      final comment = await ref.read(communityRepositoryProvider).createComment(widget.postId, text);
      if (!mounted) return;
      setState(() {
        _comments.add(comment);
        _isSubmitting = false;
      });
      _commentController.clear();
    } catch (_) {
      if (!mounted) return;
      setState(() => _isSubmitting = false);
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('댓글 등록에 실패했습니다.')),
      );
    }
  }

  @override
  void dispose() {
    _commentController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('게시글'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(
            icon: Icon(_isBookmarked ? Icons.bookmark : Icons.bookmark_outline),
            onPressed: () => setState(() => _isBookmarked = !_isBookmarked),
          ),
        ],
      ),
      body: FutureBuilder<Map<String, dynamic>>(
        future: _postFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  const Icon(Icons.error_outline, size: 48),
                  const SizedBox(height: 8),
                  const Text('게시글을 불러올 수 없습니다.'),
                  const SizedBox(height: 8),
                  FilledButton(
                    onPressed: () => setState(() {
                      _postFuture = ref.read(restClientProvider).getPost(widget.postId);
                    }),
                    child: const Text('다시 시도'),
                  ),
                ],
              ),
            );
          }
          final data = snapshot.data;
          return _buildContent(theme, data);
        },
      ),
      bottomNavigationBar: _buildCommentInput(theme),
    );
  }

  Widget _buildContent(ThemeData theme, Map<String, dynamic>? data) {
    final title = data?['title'] as String? ?? '게시글을 불러올 수 없습니다.';
    final content = data?['content'] as String? ?? '서버 연결을 확인해주세요.';
    final author = data?['author_name'] as String? ?? '사용자';
    final authorRole = data?['author_role'] as String? ?? 'user';
    final likeCount = data?['like_count'] as int? ?? 0;
    final commentCount = data?['comment_count'] as int? ?? 0;
    final createdAt = data?['created_at'] as String? ?? '';

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 작성자 정보
          Row(
            children: [
              CircleAvatar(
                radius: 20,
                backgroundColor: theme.colorScheme.primaryContainer,
                child: Text(author.isNotEmpty ? author[0] : '?'),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(author, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                        if (authorRole == 'clinician' || authorRole == 'expert') ...[
                          const SizedBox(width: 4),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                            decoration: BoxDecoration(
                              color: Colors.blue.withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(Icons.verified, size: 12, color: Colors.blue[700]),
                                const SizedBox(width: 2),
                                Text('전문가', style: TextStyle(fontSize: 10, color: Colors.blue[700], fontWeight: FontWeight.w600)),
                              ],
                            ),
                          ),
                        ],
                      ],
                    ),
                    Text(createdAt.length > 10 ? createdAt.substring(0, 10) : createdAt,
                        style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 제목
          Text(title, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 12),

          // 첨부 이미지 (있는 경우)
          if (data?['image_url'] != null && (data!['image_url'] as String).isNotEmpty) ...[
            ClipRRect(
              borderRadius: BorderRadius.circular(12),
              child: ManpasikCachedImage(
                imageUrl: data['image_url'] as String,
                width: double.infinity,
                height: 200,
                fit: BoxFit.cover,
                borderRadius: BorderRadius.circular(12),
              ),
            ),
            const SizedBox(height: 16),
          ],

          // 본문
          Text(content, style: theme.textTheme.bodyMedium?.copyWith(height: 1.8)),
          const SizedBox(height: 16),
          const Divider(),

          // 좋아요/댓글
          Row(
            children: [
              InkWell(
                onTap: () {
                  setState(() => _isLiked = !_isLiked);
                  ref.read(communityRepositoryProvider).toggleLike(widget.postId);
                },
                child: Row(
                  children: [
                    Icon(
                      _isLiked ? Icons.favorite : Icons.favorite_outline,
                      size: 20,
                      color: _isLiked ? Colors.red : null,
                    ),
                    const SizedBox(width: 4),
                    Text('${likeCount + (_isLiked ? 1 : 0)}'),
                  ],
                ),
              ),
              const SizedBox(width: 24),
              Row(
                children: [
                  const Icon(Icons.chat_bubble_outline, size: 20),
                  const SizedBox(width: 4),
                  Text('$commentCount'),
                ],
              ),
            ],
          ),
          const Divider(),

          // 댓글 섹션
          Text('댓글 (${_comments.length})', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          if (_commentsLoading)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 24),
              child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
            )
          else if (_comments.isEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 24),
              child: Center(
                child: Text('첫 댓글을 남겨보세요!', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ),
            )
          else
            ..._comments.map((c) => _buildCommentTile(theme, c)),
        ],
      ),
    );
  }

  Widget _buildCommentTile(ThemeData theme, Comment comment) {
    final timeAgo = DateTime.now().difference(comment.createdAt);
    final timeLabel = timeAgo.inMinutes < 60
        ? '${timeAgo.inMinutes}분 전'
        : timeAgo.inHours < 24
            ? '${timeAgo.inHours}시간 전'
            : '${timeAgo.inDays}일 전';

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: 14,
            backgroundColor: theme.colorScheme.secondaryContainer,
            child: Text(
              comment.authorName.isNotEmpty ? comment.authorName[0] : '?',
              style: const TextStyle(fontSize: 11),
            ),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(comment.authorName, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
                    const SizedBox(width: 8),
                    Text(timeLabel, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant, fontSize: 11)),
                  ],
                ),
                const SizedBox(height: 2),
                Text(comment.content, style: theme.textTheme.bodySmall),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCommentInput(ThemeData theme) {
    return SafeArea(
      child: Container(
        padding: const EdgeInsets.fromLTRB(16, 8, 8, 8),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          boxShadow: [BoxShadow(color: Colors.black12, blurRadius: 4, offset: const Offset(0, -2))],
        ),
        child: Row(
          children: [
            Expanded(
              child: TextField(
                controller: _commentController,
                decoration: const InputDecoration(
                  hintText: '댓글을 입력하세요...',
                  border: InputBorder.none,
                  contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                ),
              ),
            ),
            _isSubmitting
                ? const Padding(
                    padding: EdgeInsets.all(12),
                    child: SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2)),
                  )
                : IconButton(
                    icon: const Icon(Icons.send, color: AppTheme.sanggamGold),
                    onPressed: _submitComment,
                  ),
          ],
        ),
      ),
    );
  }
}
