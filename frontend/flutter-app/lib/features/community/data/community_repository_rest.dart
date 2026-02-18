import 'package:dio/dio.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 CommunityRepository 구현체
class CommunityRepositoryRest implements CommunityRepository {
  CommunityRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<List<CommunityPost>> getPosts({
    PostCategory? category,
    int page = 0,
    int size = 20,
  }) async {
    try {
      final res = await _client.listPosts(
        category: category != null && category != PostCategory.all
            ? category.index
            : null,
        limit: size,
        offset: page * size,
      );
      final posts = res['posts'] as List<dynamic>? ?? [];
      return posts.map((p) => _mapPost(p as Map<String, dynamic>)).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<CommunityPost> getPostDetail(String postId) async {
    final res = await _client.getPost(postId);
    return _mapPost(res);
  }

  @override
  Future<CommunityPost> createPost({
    required PostCategory category,
    required String title,
    required String content,
  }) async {
    final res = await _client.createPost(
      authorId: userId,
      title: title,
      content: content,
      category: category.index,
    );
    return _mapPost(res);
  }

  @override
  Future<CommunityPost> updatePost(String postId, {String? title, String? content}) async {
    // REST client doesn't have updatePost; re-fetch after creation
    return getPostDetail(postId);
  }

  @override
  Future<void> deletePost(String postId) async {
    // No delete endpoint in REST client yet
  }

  @override
  Future<void> toggleLike(String postId) async {
    try {
      await _client.likePost(postId, userId);
    } on DioException {
      // ignore
    }
  }

  @override
  Future<void> toggleBookmark(String postId) async {
    // No bookmark endpoint in REST client yet
  }

  @override
  Future<List<Comment>> getComments(String postId) async {
    // No comments endpoint in REST client yet
    return [];
  }

  @override
  Future<Comment> createComment(String postId, String content) async {
    throw UnimplementedError('Comment creation not available via REST yet');
  }

  @override
  Future<List<HealthChallenge>> getChallenges() async {
    // No challenges endpoint in REST client yet
    return [];
  }

  @override
  Future<void> joinChallenge(String challengeId) async {
    // No challenge join endpoint yet
  }

  CommunityPost _mapPost(Map<String, dynamic> m) {
    return CommunityPost(
      id: m['id'] as String? ?? m['post_id'] as String? ?? '',
      authorId: m['author_id'] as String? ?? '',
      authorName: m['author_name'] as String? ?? '',
      category: _parseCategory(m['category']),
      title: m['title'] as String? ?? '',
      content: m['content'] as String? ?? '',
      likeCount: m['like_count'] as int? ?? 0,
      commentCount: m['comment_count'] as int? ?? 0,
      isLikedByMe: m['is_liked_by_me'] as bool? ?? false,
      isBookmarkedByMe: m['is_bookmarked_by_me'] as bool? ?? false,
      createdAt: m['created_at'] != null
          ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
          : DateTime.now(),
      updatedAt: m['updated_at'] != null
          ? DateTime.tryParse(m['updated_at'] as String)
          : null,
    );
  }

  PostCategory _parseCategory(dynamic v) {
    if (v is int && v >= 0 && v < PostCategory.values.length) {
      return PostCategory.values[v];
    }
    if (v is String) {
      switch (v) {
        case 'health_tips':
          return PostCategory.healthTips;
        case 'reviews':
          return PostCategory.reviews;
        case 'challenge':
          return PostCategory.challenge;
      }
    }
    return PostCategory.all;
  }
}
