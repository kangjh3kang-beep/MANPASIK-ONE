// 도메인 모델 단위 테스트
// - DeviceItem, MeasurementRepository 모델들, UserProfileInfo, SubscriptionInfoDto

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';

void main() {
  group('DeviceItem 모델 테스트', () {
    // 기본 생성
    test('DeviceItem을 필수 필드로 생성할 수 있다', () {
      const device = DeviceItem(
        deviceId: 'dev-1',
        name: '만파식 리더기 v1',
        firmwareVersion: '2.0.0',
        status: 'online',
      );

      expect(device.deviceId, 'dev-1');
      expect(device.name, '만파식 리더기 v1');
      expect(device.firmwareVersion, '2.0.0');
      expect(device.status, 'online');
      expect(device.batteryPercent, 0); // 기본값
      expect(device.lastSeen, isNull);
    });

    // 모든 필드 지정
    test('DeviceItem의 모든 필드를 지정하여 생성할 수 있다', () {
      final now = DateTime.now();
      final device = DeviceItem(
        deviceId: 'dev-2',
        name: '테스트 디바이스',
        firmwareVersion: '1.5.0',
        status: 'offline',
        batteryPercent: 75,
        lastSeen: now,
      );

      expect(device.batteryPercent, 75);
      expect(device.lastSeen, now);
      expect(device.status, 'offline');
    });
  });

  group('Measurement 모델 테스트', () {
    // StartSessionResult
    test('StartSessionResult를 생성할 수 있다', () {
      final now = DateTime.now();
      final result = StartSessionResult(
        sessionId: 'session-abc',
        startedAt: now,
      );

      expect(result.sessionId, 'session-abc');
      expect(result.startedAt, now);
    });

    // EndSessionResult
    test('EndSessionResult를 생성할 수 있다', () {
      const result = EndSessionResult(
        sessionId: 's1',
        totalMeasurements: 10,
      );

      expect(result.sessionId, 's1');
      expect(result.totalMeasurements, 10);
      expect(result.endedAt, isNull);
    });

    // MeasurementHistoryItem
    test('MeasurementHistoryItem을 생성할 수 있다', () {
      const item = MeasurementHistoryItem(
        sessionId: 's1',
        cartridgeType: 'glucose',
        primaryValue: 105.3,
        unit: 'mg/dL',
      );

      expect(item.sessionId, 's1');
      expect(item.cartridgeType, 'glucose');
      expect(item.primaryValue, 105.3);
      expect(item.unit, 'mg/dL');
      expect(item.measuredAt, isNull);
    });

    // MeasurementHistoryResult
    test('MeasurementHistoryResult를 생성할 수 있다', () {
      const result = MeasurementHistoryResult(
        items: [
          MeasurementHistoryItem(
            sessionId: 's1',
            cartridgeType: 'glucose',
            primaryValue: 98.4,
            unit: 'mg/dL',
          ),
          MeasurementHistoryItem(
            sessionId: 's2',
            cartridgeType: 'cholesterol',
            primaryValue: 180.0,
            unit: 'mg/dL',
          ),
        ],
        totalCount: 2,
      );

      expect(result.items.length, 2);
      expect(result.totalCount, 2);
      expect(result.items.first.cartridgeType, 'glucose');
    });
  });

  group('User 모델 테스트', () {
    // UserProfileInfo 기본 생성
    test('UserProfileInfo를 필수 필드로 생성할 수 있다', () {
      const profile = UserProfileInfo(
        userId: 'user-1',
        email: 'test@manpasik.com',
        displayName: '테스트 사용자',
      );

      expect(profile.userId, 'user-1');
      expect(profile.email, 'test@manpasik.com');
      expect(profile.displayName, '테스트 사용자');
      expect(profile.avatarUrl, isNull);
      expect(profile.language, isNull);
      expect(profile.timezone, isNull);
      expect(profile.subscriptionTier, 0); // 기본 Free
    });

    // SubscriptionInfoDto
    test('SubscriptionInfoDto를 생성할 수 있다', () {
      const sub = SubscriptionInfoDto(
        userId: 'user-1',
        tier: 2,
        maxDevices: 5,
        maxFamilyMembers: 4,
        aiCoachingEnabled: true,
        telemedicineEnabled: true,
      );

      expect(sub.tier, 2);
      expect(sub.maxDevices, 5);
      expect(sub.maxFamilyMembers, 4);
      expect(sub.aiCoachingEnabled, isTrue);
      expect(sub.telemedicineEnabled, isTrue);
    });

    // SubscriptionInfoDto 기본값
    test('SubscriptionInfoDto 기본값은 Free 티어이다', () {
      const sub = SubscriptionInfoDto(userId: 'user-1');

      expect(sub.tier, 0);
      expect(sub.maxDevices, 1);
      expect(sub.maxFamilyMembers, 1);
      expect(sub.aiCoachingEnabled, isFalse);
      expect(sub.telemedicineEnabled, isFalse);
    });
  });
}
