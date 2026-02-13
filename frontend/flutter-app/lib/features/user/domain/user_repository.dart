/// 사용자 프로필·구독 Repository 인터페이스
abstract class UserRepository {
  Future<UserProfileInfo?> getProfile(String userId);
  Future<SubscriptionInfoDto?> getSubscription(String userId);
}

class UserProfileInfo {
  final String userId;
  final String email;
  final String displayName;
  final String? avatarUrl;
  final String? language;
  final String? timezone;
  final int subscriptionTier;

  const UserProfileInfo({
    required this.userId,
    required this.email,
    required this.displayName,
    this.avatarUrl,
    this.language,
    this.timezone,
    this.subscriptionTier = 0,
  });
}

class SubscriptionInfoDto {
  final String userId;
  final int tier;
  final int maxDevices;
  final int maxFamilyMembers;
  final bool aiCoachingEnabled;
  final bool telemedicineEnabled;

  const SubscriptionInfoDto({
    required this.userId,
    this.tier = 0,
    this.maxDevices = 1,
    this.maxFamilyMembers = 1,
    this.aiCoachingEnabled = false,
    this.telemedicineEnabled = false,
  });
}
