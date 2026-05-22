import 'package:dio/dio.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/core/network/offline_queue.dart';

/// REST Gateway를 사용하는 CommunityRepository 구현체
class CommunityRepositoryRest implements CommunityRepository {
  CommunityRepositoryRest(this._client, {required this.userId, OfflineQueue? offlineQueue})
      : _offlineQueue = offlineQueue ?? OfflineQueue.instance;

  final ManPaSikRestClient _client;
  final String userId;
  final OfflineQueue _offlineQueue;

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
    try {
      final res = await _client.createPost(
        authorId: userId,
        title: title,
        content: content,
        category: category.index,
      );
      return _mapPost(res);
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'POST',
          path: '/posts',
          body: {
            'author_id': userId,
            'title': title,
            'content': content,
            'category': category.index,
          },
          createdAt: DateTime.now(),
        ));
        // Optimistic result
        return CommunityPost(
          id: 'pending-${DateTime.now().millisecondsSinceEpoch}',
          authorId: userId,
          authorName: '',
          category: category,
          title: title,
          content: content,
          likeCount: 0,
          commentCount: 0,
          isLikedByMe: false,
          isBookmarkedByMe: false,
          createdAt: DateTime.now(),
        );
      }
      rethrow;
    }
  }

  @override
  Future<CommunityPost> updatePost(String postId, {String? title, String? content}) async {
    // REST client doesn't have updatePost; re-fetch after creation
    return getPostDetail(postId);
  }

  @override
  Future<void> deletePost(String postId) async {
    try {
      await _client.deletePost(postId);
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'DELETE',
          path: '/posts/$postId',
          createdAt: DateTime.now(),
        ));
        return; // Optimistic success
      }
      rethrow;
    }
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
    try {
      final res = await _client.listComments(postId);
      final comments = res['comments'] as List<dynamic>? ?? [];
      return comments.map((c) {
        final m = c as Map<String, dynamic>;
        return Comment(
          id: m['comment_id'] as String? ?? m['id'] as String? ?? '',
          postId: postId, // Add postId
          authorId: m['author_id'] as String? ?? '',
          authorName: m['author_name'] as String? ?? '',
          content: m['content'] as String? ?? '',
          createdAt: m['created_at'] != null
              ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
              : DateTime.now(),
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<Comment> createComment(String postId, String content) async {
    try {
      final res = await _client.createComment(
        postId: postId,
        authorId: userId,
        content: content,
      );
      return Comment(
        id: res['comment_id'] as String? ?? res['id'] as String? ?? '',
        postId: postId,
        authorId: userId,
        authorName: res['author_name'] as String? ?? '',
        content: content,
        createdAt: DateTime.now(),
      );
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'POST',
          path: '/posts/$postId/comments',
          body: {
            'post_id': postId,
            'author_id': userId,
            'content': content,
          },
          createdAt: DateTime.now(),
        ));
        // Optimistic result
        return Comment(
          id: 'pending-${DateTime.now().millisecondsSinceEpoch}',
          postId: postId,
          authorId: userId,
          authorName: '',
          content: content,
          createdAt: DateTime.now(),
        );
      }
      rethrow;
    }
  }

  @override
  Future<List<HealthChallenge>> getChallenges() async {
    try {
      final res = await _client.listChallenges();
      final challenges = res['challenges'] as List<dynamic>? ?? [];
      return challenges.map((c) {
        final m = c as Map<String, dynamic>;
        return HealthChallenge(
          id: m['challenge_id'] as String? ?? m['id'] as String? ?? '',
          title: m['title'] as String? ?? '',
          description: m['description'] as String? ?? '',
          startDate: m['start_date'] != null
              ? DateTime.tryParse(m['start_date'] as String) ?? DateTime.now()
              : DateTime.now(),
          endDate: m['end_date'] != null
              ? DateTime.tryParse(m['end_date'] as String) ??
                  DateTime.now().add(const Duration(days: 30))
              : DateTime.now().add(const Duration(days: 30)),
          participantCount: m['participant_count'] as int? ?? 0,
          isJoined: m['is_joined'] as bool? ?? false,
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<void> joinChallenge(String challengeId) async {
    try {
      await _client.joinChallenge(
        challengeId: challengeId,
        userId: userId,
      );
    } on DioException catch (e) {
      if (e.type == DioExceptionType.connectionError ||
          e.type == DioExceptionType.connectionTimeout) {
        await _offlineQueue.enqueue(OfflineRequest(
          method: 'POST',
          path: '/challenges/$challengeId/join',
          body: {'user_id': userId},
          createdAt: DateTime.now(),
        ));
        return; // Optimistic success
      }
      rethrow;
    }
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
