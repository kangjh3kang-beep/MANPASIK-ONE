import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/devices/data/device_repository_rest.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('DeviceRepositoryRest', () {
    test('DeviceRepositoryRest는 DeviceRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      expect(repo, isA<DeviceRepository>());
    });

    test('listDevices는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.listDevices('user-1');
      expect(devices, isEmpty);
    });

    test('listDevices 빈 userId에 대해 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.listDevices('');
      expect(devices, isEmpty);
    });

    test('getConnectedDevices는 시뮬레이션 디바이스 10개를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.getConnectedDevices();
      expect(devices, isA<List<ConnectedDevice>>());
      expect(devices.length, 10);
    });

    test('getConnectedDevices 디바이스에 올바른 필드가 있다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.getConnectedDevices();
      final first = devices.first;
      expect(first.id, isNotEmpty);
      expect(first.name, isNotEmpty);
      expect(first.batteryLevel, greaterThanOrEqualTo(0));
      expect(first.batteryLevel, lessThanOrEqualTo(100));
    });

    test('getConnectedDevices 다양한 DeviceType을 포함한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.getConnectedDevices();
      final types = devices.map((d) => d.type).toSet();
      expect(types, contains(DeviceType.gasCartridge));
      expect(types, contains(DeviceType.envCartridge));
      expect(types, contains(DeviceType.bioCartridge));
    });

    test('getConnectedDevices 연결/비연결 디바이스가 혼재한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = DeviceRepositoryRest(client);
      final devices = await repo.getConnectedDevices();
      final statuses = devices.map((d) => d.status).toSet();
      expect(statuses, contains(DeviceConnectionStatus.connected));
      expect(statuses, contains(DeviceConnectionStatus.disconnected));
    });
  });

  group('DeviceItem 도메인 모델', () {
    test('DeviceItem은 올바르게 생성된다', () {
      const item = DeviceItem(
        deviceId: 'dev-1',
        name: '테스트 디바이스',
        firmwareVersion: '2.0.0',
        status: 'online',
        batteryPercent: 90,
      );
      expect(item.deviceId, 'dev-1');
      expect(item.name, '테스트 디바이스');
      expect(item.firmwareVersion, '2.0.0');
      expect(item.status, 'online');
      expect(item.batteryPercent, 90);
      expect(item.lastSeen, isNull);
    });
  });

  group('ConnectedDevice 도메인 모델', () {
    test('ConnectedDevice는 기본값으로 생성된다', () {
      const device = ConnectedDevice(
        id: 'cd-1',
        name: '센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 80,
        signalStrength: 95,
      );
      expect(device.latestReadings, isEmpty);
      expect(device.currentValues, isEmpty);
    });

    test('ConnectedDevice는 currentValues를 포함한다', () {
      const device = ConnectedDevice(
        id: 'cd-2',
        name: '환경 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 70,
        signalStrength: 85,
        currentValues: {'Temp': '24.5°C', 'Humidity': '45%'},
        latestReadings: [24.0, 24.1, 24.5],
      );
      expect(device.currentValues, hasLength(2));
      expect(device.latestReadings, hasLength(3));
    });
  });
}
