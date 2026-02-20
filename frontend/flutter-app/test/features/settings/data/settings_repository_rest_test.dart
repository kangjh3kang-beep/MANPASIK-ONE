import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/settings/data/settings_repository_rest.dart';
import 'package:manpasik/features/settings/domain/settings_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('SettingsRepositoryRest', () {
    test('SettingsRepositoryRest는 SettingsRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = SettingsRepositoryRest(client, 'user-1');
      expect(repo, isA<SettingsRepository>());
    });

    test('getSettings는 DioException 시 기본 설정을 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = SettingsRepositoryRest(client, 'user-1');
      final settings = await repo.getSettings();
      expect(settings, isA<AppSettings>());
    });

    test('getEmergencyContacts는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = SettingsRepositoryRest(client, 'user-1');
      final contacts = await repo.getEmergencyContacts();
      expect(contacts, isEmpty);
    });
  });

  group('Settings 도메인 모델', () {
    test('AppSettings 기본값 확인', () {
      const settings = AppSettings();
      expect(settings.locale, 'ko');
      expect(settings.themeMode, 'system');
      expect(settings.biometricEnabled, isFalse);
      expect(settings.notificationsEnabled, isTrue);
      expect(settings.fontScale, 1.0);
    });

    test('AppSettings 커스텀 값 확인', () {
      const settings = AppSettings(
        locale: 'en',
        themeMode: 'dark',
        biometricEnabled: true,
        fontScale: 1.5,
        highContrast: true,
      );
      expect(settings.locale, 'en');
      expect(settings.themeMode, 'dark');
      expect(settings.biometricEnabled, isTrue);
      expect(settings.fontScale, 1.5);
      expect(settings.highContrast, isTrue);
    });
  });
}
