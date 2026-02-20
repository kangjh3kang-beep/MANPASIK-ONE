import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/user/data/user_repository_rest.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('UserRepositoryRest', () {
    test('UserRepositoryRest는 UserRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = UserRepositoryRest(client);
      expect(repo, isA<UserRepository>());
    });

    test('getProfile은 DioException 시 null을 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = UserRepositoryRest(client);
      final profile = await repo.getProfile('user-1');
      expect(profile, isNull);
    });

    test('getSubscription은 DioException 시 null을 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = UserRepositoryRest(client);
      final sub = await repo.getSubscription('user-1');
      expect(sub, isNull);
    });
  });

  group('User 도메인 모델', () {
    test('UserProfileInfo 생성 확인', () {
      const profile = UserProfileInfo(
        userId: 'user-1',
        email: 'test@manpasik.com',
        displayName: '테스트',
        avatarUrl: null,
        language: 'ko',
        timezone: 'Asia/Seoul',
        subscriptionTier: 2,
      );
      expect(profile.userId, 'user-1');
      expect(profile.subscriptionTier, 2);
      expect(profile.avatarUrl, isNull);
    });

    test('SubscriptionInfoDto 생성 확인', () {
      const sub = SubscriptionInfoDto(
        userId: 'user-1',
        tier: 3,
        maxDevices: 10,
        maxFamilyMembers: 5,
        aiCoachingEnabled: true,
        telemedicineEnabled: true,
      );
      expect(sub.tier, 3);
      expect(sub.aiCoachingEnabled, isTrue);
    });
  });
}
