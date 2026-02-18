/// 커뮤니티 도메인 모델 및 리포지토리
///
/// 게시글 CRUD, 좋아요/댓글, 건강 챌린지

/// 게시글 카테고리
enum PostCategory { all, healthTips, reviews, challenge }

/// 게시글
class CommunityPost {
  final String id;
  final String authorId;
  final String authorName;
  final PostCategory category;
  final String title;
  final String content;
  final int likeCount;
  final int commentCount;
  final bool isLikedByMe;
  final bool isBookmarkedByMe;
  final DateTime createdAt;
  final DateTime? updatedAt;

  const CommunityPost({
    required this.id,
    required this.authorId,
    required this.authorName,
    required this.category,
    required this.title,
    required this.content,
    required this.likeCount,
    required this.commentCount,
    required this.isLikedByMe,
    required this.isBookmarkedByMe,
    required this.createdAt,
    this.updatedAt,
  });
}

/// 댓글
class Comment {
  final String id;
  final String postId;
  final String authorId;
  final String authorName;
  final String content;
  final DateTime createdAt;

  const Comment({
    required this.id,
    required this.postId,
    required this.authorId,
    required this.authorName,
    required this.content,
    required this.createdAt,
  });
}

/// 건강 챌린지
class HealthChallenge {
  final String id;
  final String title;
  final String description;
  final DateTime startDate;
  final DateTime endDate;
  final int participantCount;
  final bool isJoined;
  final double? myProgress; // 0.0 ~ 1.0

  const HealthChallenge({
    required this.id,
    required this.title,
    required this.description,
    required this.startDate,
    required this.endDate,
    required this.participantCount,
    required this.isJoined,
    this.myProgress,
  });
}

/// 커뮤니티 리포지토리 인터페이스
abstract class CommunityRepository {
  /// 게시글 목록 조회
  Future<List<CommunityPost>> getPosts({
    PostCategory? category,
    int page = 0,
    int size = 20,
  });

  /// 게시글 상세
  Future<CommunityPost> getPostDetail(String postId);

  /// 게시글 작성
  Future<CommunityPost> createPost({
    required PostCategory category,
    required String title,
    required String content,
  });

  /// 게시글 수정
  Future<CommunityPost> updatePost(String postId, {String? title, String? content});

  /// 게시글 삭제
  Future<void> deletePost(String postId);

  /// 좋아요 토글
  Future<void> toggleLike(String postId);

  /// 북마크 토글
  Future<void> toggleBookmark(String postId);

  /// 댓글 목록
  Future<List<Comment>> getComments(String postId);

  /// 댓글 작성
  Future<Comment> createComment(String postId, String content);

  /// 챌린지 목록
  Future<List<HealthChallenge>> getChallenges();

  /// 챌린지 참가
  Future<void> joinChallenge(String challengeId);
}
