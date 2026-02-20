import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/community/data/community_repository_rest.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('CommunityRepositoryRest', () {
    test('CommunityRepositoryRest는 CommunityRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = CommunityRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<CommunityRepository>());
    });

    test('getPosts는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = CommunityRepositoryRest(client, userId: 'user-1');
      final posts = await repo.getPosts();
      expect(posts, isEmpty);
    });

    test('getPosts category 필터 적용', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = CommunityRepositoryRest(client, userId: 'user-1');
      final posts = await repo.getPosts(category: PostCategory.healthTips);
      expect(posts, isEmpty);
    });

    test('getPosts 페이지네이션 파라미터', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = CommunityRepositoryRest(client, userId: 'user-1');
      final posts = await repo.getPosts(page: 1, size: 10);
      expect(posts, isEmpty);
    });
  });

  group('Community 도메인 모델', () {
    test('PostCategory enum 값 확인', () {
      expect(PostCategory.values.length, greaterThanOrEqualTo(2));
      expect(PostCategory.values, contains(PostCategory.all));
      expect(PostCategory.values, contains(PostCategory.healthTips));
    });

    test('CommunityPost 생성 확인', () {
      final post = CommunityPost(
        id: 'post-1',
        authorId: 'user-1',
        authorName: '홍길동',
        title: '건강 관리 팁',
        content: '매일 운동하세요',
        category: PostCategory.healthTips,
        likeCount: 10,
        commentCount: 3,
        isLikedByMe: false,
        isBookmarkedByMe: false,
        createdAt: DateTime(2026, 2, 19),
      );
      expect(post.id, 'post-1');
      expect(post.likeCount, 10);
      expect(post.isLikedByMe, isFalse);
    });

    test('HealthChallenge 생성 확인', () {
      final challenge = HealthChallenge(
        id: 'ch-1',
        title: '만보 걷기 챌린지',
        description: '매일 10000보 걷기',
        startDate: DateTime(2026, 2, 1),
        endDate: DateTime(2026, 2, 28),
        participantCount: 150,
        isJoined: true,
        myProgress: 0.75,
      );
      expect(challenge.participantCount, 150);
      expect(challenge.myProgress, 0.75);
    });
  });
}
