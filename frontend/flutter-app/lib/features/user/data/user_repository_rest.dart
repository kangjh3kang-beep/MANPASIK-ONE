import 'package:dio/dio.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 UserRepository 구현체
///
/// 웹 플랫폼에서 gRPC 대신 REST API를 통해 사용자 프로필/구독 관리.
class UserRepositoryRest implements UserRepository {
  UserRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<UserProfileInfo?> getProfile(String userId) async {
    try {
      final res = await _client.getProfile(userId);
      return UserProfileInfo(
        userId: res['user_id'] as String? ?? userId,
        email: res['email'] as String? ?? '',
        displayName: res['display_name'] as String? ?? '',
        avatarUrl: res['avatar_url'] as String?,
        language: res['language'] as String?,
        timezone: res['timezone'] as String?,
        subscriptionTier: res['subscription_tier'] as int? ?? 0,
      );
    } on DioException {
      return null;
    }
  }

  @override
  Future<SubscriptionInfoDto?> getSubscription(String userId) async {
    try {
      final res = await _client.getSubscription(userId);
      return SubscriptionInfoDto(
        userId: res['user_id'] as String? ?? userId,
        tier: res['tier'] as int? ?? 0,
        maxDevices: res['max_devices'] as int? ?? 1,
        maxFamilyMembers: res['max_family_members'] as int? ?? 1,
        aiCoachingEnabled: res['ai_coaching_enabled'] as bool? ?? false,
        telemedicineEnabled: res['telemedicine_enabled'] as bool? ?? false,
      );
    } on DioException {
      return null;
    }
  }
}
