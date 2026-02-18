/// 온보딩 도메인 모델 및 리포지토리
///
/// 초기 설정, 건강 프로필, 기기 페어링 안내

/// 온보딩 단계
enum OnboardingStep { welcome, healthProfile, devicePairing, permissions, complete }

/// 건강 프로필 (온보딩 시 수집)
class HealthProfile {
  final int? birthYear;
  final String? gender; // 'male', 'female', 'other'
  final double? heightCm;
  final double? weightKg;
  final List<String> healthConditions; // 'diabetes', 'hypertension', etc.
  final List<String> medications;

  const HealthProfile({
    this.birthYear,
    this.gender,
    this.heightCm,
    this.weightKg,
    this.healthConditions = const [],
    this.medications = const [],
  });
}

/// 온보딩 리포지토리 인터페이스
abstract class OnboardingRepository {
  Future<bool> isOnboardingCompleted(String userId);
  Future<void> saveHealthProfile(String userId, HealthProfile profile);
  Future<void> completeOnboarding(String userId);
}
