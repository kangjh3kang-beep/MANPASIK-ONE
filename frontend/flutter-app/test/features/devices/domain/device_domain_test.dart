import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

void main() {
  group('DeviceConnectionStatus enum', () {
    test('값 목록 확인', () {
      expect(DeviceConnectionStatus.values, hasLength(3));
      expect(DeviceConnectionStatus.values,
          contains(DeviceConnectionStatus.connected));
      expect(DeviceConnectionStatus.values,
          contains(DeviceConnectionStatus.disconnected));
      expect(DeviceConnectionStatus.values,
          contains(DeviceConnectionStatus.scanning));
    });
  });

  group('DeviceType enum', () {
    test('값 목록 확인', () {
      expect(DeviceType.values, hasLength(4));
      expect(DeviceType.values, contains(DeviceType.gasCartridge));
      expect(DeviceType.values, contains(DeviceType.envCartridge));
      expect(DeviceType.values, contains(DeviceType.bioCartridge));
      expect(DeviceType.values, contains(DeviceType.unknown));
    });
  });

  group('ConnectedDevice', () {
    test('기본 생성 확인', () {
      const device = ConnectedDevice(
        id: 'dev-1',
        name: 'ManPaSik Reader',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 85,
        signalStrength: 92,
      );
      expect(device.id, 'dev-1');
      expect(device.name, 'ManPaSik Reader');
      expect(device.type, DeviceType.bioCartridge);
      expect(device.status, DeviceConnectionStatus.connected);
      expect(device.batteryLevel, 85);
      expect(device.signalStrength, 92);
      expect(device.latestReadings, isEmpty);
      expect(device.currentValues, isEmpty);
    });

    test('latestReadings 포함 생성', () {
      const device = ConnectedDevice(
        id: 'dev-2',
        name: '가스 센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 50,
        signalStrength: 70,
        latestReadings: [1.0, 2.0, 3.0, 4.0],
        currentValues: {'CO2': 450, 'Temp': 24.5},
      );
      expect(device.latestReadings, hasLength(4));
      expect(device.currentValues['CO2'], 450);
      expect(device.currentValues['Temp'], 24.5);
    });
  });

  group('DeviceItem', () {
    test('기본 생성 확인', () {
      const item = DeviceItem(
        deviceId: 'dev-1',
        name: 'Reader A',
        firmwareVersion: '2.1.0',
        status: 'online',
      );
      expect(item.deviceId, 'dev-1');
      expect(item.name, 'Reader A');
      expect(item.firmwareVersion, '2.1.0');
      expect(item.status, 'online');
      expect(item.batteryPercent, 0); // default
      expect(item.lastSeen, isNull);
    });

    test('모든 필드 포함 생성', () {
      final now = DateTime.now();
      final item = DeviceItem(
        deviceId: 'dev-2',
        name: 'Reader B',
        firmwareVersion: '3.0.1',
        status: 'measuring',
        batteryPercent: 75,
        lastSeen: now,
      );
      expect(item.batteryPercent, 75);
      expect(item.lastSeen, now);
    });

    test('다양한 status 값', () {
      for (final status in ['online', 'offline', 'measuring', 'updating', 'error']) {
        final item = DeviceItem(
          deviceId: 'dev-test',
          name: 'Test',
          firmwareVersion: '1.0.0',
          status: status,
        );
        expect(item.status, status);
      }
    });
  });
}
