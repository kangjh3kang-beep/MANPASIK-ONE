/// 설정 도메인 모델 및 리포지토리
///
/// 앱 설정, 프로필, 보안, 접근성, 긴급 연락처

/// 앱 설정
class AppSettings {
  final String locale;
  final String themeMode; // 'light', 'dark', 'system'
  final bool biometricEnabled;
  final bool notificationsEnabled;
  final double fontScale;
  final bool highContrast;

  const AppSettings({
    this.locale = 'ko',
    this.themeMode = 'system',
    this.biometricEnabled = false,
    this.notificationsEnabled = true,
    this.fontScale = 1.0,
    this.highContrast = false,
  });
}

/// 긴급 연락처
class EmergencyContact {
  final String id;
  final String name;
  final String phone;
  final String relationship;

  const EmergencyContact({
    required this.id,
    required this.name,
    required this.phone,
    required this.relationship,
  });
}

/// 설정 리포지토리 인터페이스
abstract class SettingsRepository {
  Future<AppSettings> getSettings();
  Future<void> updateSettings(AppSettings settings);
  Future<List<EmergencyContact>> getEmergencyContacts();
  Future<void> addEmergencyContact(EmergencyContact contact);
  Future<void> removeEmergencyContact(String contactId);
  Future<void> deleteAccount(String userId, {required String reason});
  Future<Map<String, bool>> getConsentStatus();
  Future<void> updateConsent(String consentType, bool agreed);
}
