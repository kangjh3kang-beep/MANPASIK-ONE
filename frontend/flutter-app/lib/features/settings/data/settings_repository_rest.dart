import 'package:dio/dio.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/settings/domain/settings_repository.dart';

/// REST Gateway를 사용하는 SettingsRepository 구현체
class SettingsRepositoryRest implements SettingsRepository {
  SettingsRepositoryRest(this._client, this._userId);

  final ManPaSikRestClient _client;
  final String _userId;

  @override
  Future<AppSettings> getSettings() async {
    try {
      final res = await _client.getProfile(_userId);
      return AppSettings(
        locale: res['language'] as String? ?? 'ko',
        themeMode: res['theme_mode'] as String? ?? 'system',
        biometricEnabled: res['biometric_enabled'] as bool? ?? false,
        notificationsEnabled:
            res['notifications_enabled'] as bool? ?? true,
        fontScale: (res['font_scale'] as num?)?.toDouble() ?? 1.0,
        highContrast: res['high_contrast'] as bool? ?? false,
      );
    } on DioException {
      return const AppSettings();
    }
  }

  @override
  Future<void> updateSettings(AppSettings settings) async {
    await _client.updateProfile(
      _userId,
      language: settings.locale,
    );
  }

  @override
  Future<List<EmergencyContact>> getEmergencyContacts() async {
    try {
      final res = await _client.getProfile(_userId);
      final contacts = res['emergency_contacts'] as List<dynamic>? ?? [];
      return contacts
          .map((c) {
            final m = c as Map<String, dynamic>;
            return EmergencyContact(
              id: m['id'] as String? ?? '',
              name: m['name'] as String? ?? '',
              phone: m['phone'] as String? ?? '',
              relationship: m['relationship'] as String? ?? '',
            );
          })
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<void> addEmergencyContact(EmergencyContact contact) async {
    await _client.saveEmergencySettings(
      userId: _userId,
      autoReport119: true,
      emergencyContacts: [contact.phone],
    );
  }

  @override
  Future<void> removeEmergencyContact(String contactId) async {
    // Remove by updating the contacts list without the given ID
    final contacts = await getEmergencyContacts();
    final remaining =
        contacts.where((c) => c.id != contactId).map((c) => c.phone).toList();
    await _client.saveEmergencySettings(
      userId: _userId,
      autoReport119: true,
      emergencyContacts: remaining,
    );
  }

  @override
  Future<void> deleteAccount(String userId, {required String reason}) async {
    // Account deletion handled via admin/support flow
    await _client.createInquiry(
      userId: userId,
      type: 'account_deletion',
      title: '계정 삭제 요청',
      content: reason,
    );
  }

  @override
  Future<Map<String, bool>> getConsentStatus() async {
    try {
      final res = await _client.listDataSharingConsents(_userId);
      final consents = res['consents'] as List<dynamic>? ?? [];
      final result = <String, bool>{};
      for (final c in consents) {
        final m = c as Map<String, dynamic>;
        result[m['type'] as String? ?? ''] = m['active'] as bool? ?? false;
      }
      return result;
    } on DioException {
      return {};
    }
  }

  @override
  Future<void> updateConsent(String consentType, bool agreed) async {
    if (agreed) {
      await _client.createDataSharingConsent(
        userId: _userId,
        providerId: 'system',
        dataTypes: [consentType],
      );
    }
  }
}
