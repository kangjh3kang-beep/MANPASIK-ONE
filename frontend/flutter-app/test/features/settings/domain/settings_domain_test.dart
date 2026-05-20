import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/settings/domain/settings_repository.dart';

void main() {
  group('AppSettings', () {
    test('기본값 확인', () {
      const settings = AppSettings();
      expect(settings.locale, 'ko');
      expect(settings.themeMode, 'system');
      expect(settings.biometricEnabled, isFalse);
      expect(settings.notificationsEnabled, isTrue);
      expect(settings.fontScale, 1.0);
      expect(settings.highContrast, isFalse);
    });

    test('커스텀 설정 생성', () {
      const settings = AppSettings(
        locale: 'en',
        themeMode: 'dark',
        biometricEnabled: true,
        notificationsEnabled: false,
        fontScale: 1.5,
        highContrast: true,
      );
      expect(settings.locale, 'en');
      expect(settings.themeMode, 'dark');
      expect(settings.biometricEnabled, isTrue);
      expect(settings.notificationsEnabled, isFalse);
      expect(settings.fontScale, 1.5);
      expect(settings.highContrast, isTrue);
    });

    test('fontScale 범위 확인', () {
      const smallFont = AppSettings(fontScale: 0.8);
      const largeFont = AppSettings(fontScale: 2.0);
      expect(smallFont.fontScale, 0.8);
      expect(largeFont.fontScale, 2.0);
    });
  });

  group('EmergencyContact', () {
    test('생성 확인', () {
      const contact = EmergencyContact(
        id: 'ec-1',
        name: '홍길동',
        phone: '010-1234-5678',
        relationship: '배우자',
      );
      expect(contact.id, 'ec-1');
      expect(contact.name, '홍길동');
      expect(contact.phone, '010-1234-5678');
      expect(contact.relationship, '배우자');
    });

    test('다양한 관계 확인', () {
      for (final rel in ['배우자', '부모', '자녀', '보호자', '기타']) {
        final contact = EmergencyContact(
          id: 'test',
          name: '테스트',
          phone: '000-0000-0000',
          relationship: rel,
        );
        expect(contact.relationship, rel);
      }
    });
  });
}
