import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/onboarding/domain/onboarding_repository.dart';

/// REST Gateway를 사용하는 OnboardingRepository 구현체
class OnboardingRepositoryRest implements OnboardingRepository {
  OnboardingRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<bool> isOnboardingCompleted(String userId) async {
    try {
      final res = await _client.getProfile(userId);
      // Onboarding is complete when the user has a display_name set
      final displayName = res['display_name'] as String?;
      return displayName != null && displayName.isNotEmpty;
    } on DioException {
      return false;
    }
  }

  @override
  Future<void> saveHealthProfile(
    String userId,
    HealthProfile profile,
  ) async {
    await _client.createHealthRecord(
      userId: userId,
      recordType: 0, // onboarding health profile type
      title: 'Health Profile',
      description: 'Onboarding health profile',
      metadata: {
        if (profile.birthYear != null)
          'birth_year': profile.birthYear.toString(),
        if (profile.gender != null) 'gender': profile.gender!,
        if (profile.heightCm != null)
          'height_cm': profile.heightCm.toString(),
        if (profile.weightKg != null)
          'weight_kg': profile.weightKg.toString(),
        'health_conditions': profile.healthConditions.join(','),
        'medications': profile.medications.join(','),
      },
    );
  }

  @override
  Future<void> completeOnboarding(String userId) async {
    await _client.updateProfile(userId, language: 'ko');
  }
}
