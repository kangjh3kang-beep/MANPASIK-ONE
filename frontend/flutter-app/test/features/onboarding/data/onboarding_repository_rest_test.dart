import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/onboarding/data/onboarding_repository_rest.dart';
import 'package:manpasik/features/onboarding/domain/onboarding_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('OnboardingRepositoryRest', () {
    test('OnboardingRepositoryRest는 OnboardingRepository를 구현한다', () {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = OnboardingRepositoryRest(client);
      expect(repo, isA<OnboardingRepository>());
    });

    test('isOnboardingCompleted는 DioException 시 false를 반환한다', () async {
      final client =
          ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = OnboardingRepositoryRest(client);
      final completed = await repo.isOnboardingCompleted('user-1');
      expect(completed, isFalse);
    });
  });

  group('OnboardingStep enum', () {
    test('OnboardingStep 값 확인', () {
      expect(OnboardingStep.values, hasLength(5));
      expect(OnboardingStep.values, contains(OnboardingStep.welcome));
      expect(OnboardingStep.values, contains(OnboardingStep.healthProfile));
      expect(OnboardingStep.values, contains(OnboardingStep.devicePairing));
      expect(OnboardingStep.values, contains(OnboardingStep.permissions));
      expect(OnboardingStep.values, contains(OnboardingStep.complete));
    });
  });

  group('HealthProfile 도메인 모델', () {
    test('기본 생성자로 빈 HealthProfile 생성', () {
      const profile = HealthProfile();
      expect(profile.birthYear, isNull);
      expect(profile.gender, isNull);
      expect(profile.heightCm, isNull);
      expect(profile.weightKg, isNull);
      expect(profile.healthConditions, isEmpty);
      expect(profile.medications, isEmpty);
    });

    test('모든 필드가 있는 HealthProfile 생성', () {
      const profile = HealthProfile(
        birthYear: 1990,
        gender: 'male',
        heightCm: 175.0,
        weightKg: 70.0,
        healthConditions: ['diabetes', 'hypertension'],
        medications: ['metformin'],
      );
      expect(profile.birthYear, 1990);
      expect(profile.gender, 'male');
      expect(profile.heightCm, 175.0);
      expect(profile.weightKg, 70.0);
      expect(profile.healthConditions, hasLength(2));
      expect(profile.medications, hasLength(1));
    });

    test('healthConditions 기본값은 빈 리스트', () {
      const profile = HealthProfile(birthYear: 2000);
      expect(profile.healthConditions, isEmpty);
      expect(profile.medications, isEmpty);
    });
  });
}
