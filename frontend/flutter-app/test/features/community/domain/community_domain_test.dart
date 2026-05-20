import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';

void main() {
  group('PostCategory enum', () {
    test('값 목록 확인 (4개)', () {
      expect(PostCategory.values, hasLength(4));
      expect(PostCategory.values, contains(PostCategory.all));
      expect(PostCategory.values, contains(PostCategory.healthTips));
      expect(PostCategory.values, contains(PostCategory.reviews));
      expect(PostCategory.values, contains(PostCategory.challenge));
    });
  });

  group('CommunityPost', () {
    test('기본 생성 확인', () {
      final now = DateTime.now();
      final post = CommunityPost(
        id: 'post-1',
        authorId: 'user-1',
        authorName: '홍길동',
        category: PostCategory.healthTips,
        title: '혈당 관리 팁',
        content: '아침 운동이 효과적입니다.',
        likeCount: 10,
        commentCount: 3,
        isLikedByMe: false,
        isBookmarkedByMe: true,
        createdAt: now,
      );
      expect(post.id, 'post-1');
      expect(post.authorName, '홍길동');
      expect(post.category, PostCategory.healthTips);
      expect(post.likeCount, 10);
      expect(post.isLikedByMe, isFalse);
      expect(post.isBookmarkedByMe, isTrue);
      expect(post.updatedAt, isNull);
    });

    test('updatedAt 포함 생성', () {
      final created = DateTime(2026, 4, 1);
      final updated = DateTime(2026, 4, 19);
      final post = CommunityPost(
        id: 'post-2',
        authorId: 'user-2',
        authorName: '김건강',
        category: PostCategory.reviews,
        title: '카트리지 리뷰',
        content: '정확도가 매우 좋습니다.',
        likeCount: 25,
        commentCount: 8,
        isLikedByMe: true,
        isBookmarkedByMe: false,
        createdAt: created,
        updatedAt: updated,
      );
      expect(post.updatedAt, updated);
      expect(post.isLikedByMe, isTrue);
    });
  });

  group('Comment', () {
    test('생성 확인', () {
      final now = DateTime.now();
      final comment = Comment(
        id: 'comment-1',
        postId: 'post-1',
        authorId: 'user-2',
        authorName: '김댓글',
        content: '좋은 정보 감사합니다!',
        createdAt: now,
      );
      expect(comment.id, 'comment-1');
      expect(comment.postId, 'post-1');
      expect(comment.authorName, '김댓글');
      expect(comment.content, '좋은 정보 감사합니다!');
      expect(comment.createdAt, now);
    });
  });

  group('HealthChallenge', () {
    test('기본 생성 확인', () {
      final challenge = HealthChallenge(
        id: 'ch-1',
        title: '30일 걷기 챌린지',
        description: '매일 1만보 걷기',
        startDate: DateTime(2026, 4, 1),
        endDate: DateTime(2026, 4, 30),
        participantCount: 500,
        isJoined: false,
      );
      expect(challenge.id, 'ch-1');
      expect(challenge.title, '30일 걷기 챌린지');
      expect(challenge.participantCount, 500);
      expect(challenge.isJoined, isFalse);
      expect(challenge.myProgress, isNull);
    });

    test('진행 중인 챌린지', () {
      final challenge = HealthChallenge(
        id: 'ch-2',
        title: '혈당 관리 챌린지',
        description: '7일 연속 정상 혈당',
        startDate: DateTime(2026, 4, 10),
        endDate: DateTime(2026, 4, 17),
        participantCount: 200,
        isJoined: true,
        myProgress: 0.71,
      );
      expect(challenge.isJoined, isTrue);
      expect(challenge.myProgress, 0.71);
    });
  });
}
