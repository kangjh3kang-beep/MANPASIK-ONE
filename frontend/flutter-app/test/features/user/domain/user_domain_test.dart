import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';

void main() {
  group('UserProfileInfo', () {
    test('기본 생성 확인', () {
      const profile = UserProfileInfo(
        userId: 'user-1',
        email: 'test@manpasik.com',
        displayName: '홍길동',
      );
      expect(profile.userId, 'user-1');
      expect(profile.email, 'test@manpasik.com');
      expect(profile.displayName, '홍길동');
      expect(profile.avatarUrl, isNull);
      expect(profile.language, isNull);
      expect(profile.timezone, isNull);
      expect(profile.subscriptionTier, 0);
    });

    test('모든 필드 포함 생성', () {
      const profile = UserProfileInfo(
        userId: 'user-2',
        email: 'premium@manpasik.com',
        displayName: '김프리미엄',
        avatarUrl: 'https://example.com/avatar.png',
        language: 'ko',
        timezone: 'Asia/Seoul',
        subscriptionTier: 3,
      );
      expect(profile.avatarUrl, isNotNull);
      expect(profile.language, 'ko');
      expect(profile.timezone, 'Asia/Seoul');
      expect(profile.subscriptionTier, 3);
    });
  });

  group('SubscriptionInfoDto', () {
    test('기본값 확인 (Free tier)', () {
      const sub = SubscriptionInfoDto(userId: 'user-1');
      expect(sub.userId, 'user-1');
      expect(sub.tier, 0);
      expect(sub.maxDevices, 1);
      expect(sub.maxFamilyMembers, 1);
      expect(sub.aiCoachingEnabled, isFalse);
      expect(sub.telemedicineEnabled, isFalse);
    });

    test('Premium tier 생성', () {
      const sub = SubscriptionInfoDto(
        userId: 'user-2',
        tier: 3,
        maxDevices: 5,
        maxFamilyMembers: 10,
        aiCoachingEnabled: true,
        telemedicineEnabled: true,
      );
      expect(sub.tier, 3);
      expect(sub.maxDevices, 5);
      expect(sub.maxFamilyMembers, 10);
      expect(sub.aiCoachingEnabled, isTrue);
      expect(sub.telemedicineEnabled, isTrue);
    });

    test('tier별 기능 활성화 패턴', () {
      // Free
      const free = SubscriptionInfoDto(userId: 'u1', tier: 0);
      expect(free.aiCoachingEnabled, isFalse);

      // Premium
      const premium = SubscriptionInfoDto(
        userId: 'u2',
        tier: 2,
        aiCoachingEnabled: true,
        telemedicineEnabled: true,
      );
      expect(premium.aiCoachingEnabled, isTrue);
      expect(premium.telemedicineEnabled, isTrue);
    });
  });
}
