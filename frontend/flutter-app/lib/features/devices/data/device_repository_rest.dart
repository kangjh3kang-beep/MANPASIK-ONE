import 'package:dio/dio.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 DeviceRepository 구현체
///
/// 웹 플랫폼에서 gRPC 대신 REST API를 통해 디바이스 관리.
class DeviceRepositoryRest implements DeviceRepository {
  DeviceRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<List<DeviceItem>> listDevices(String userId) async {
    try {
      final res = await _client.listDevices(userId);
      final devices = res['devices'] as List<dynamic>? ?? [];
      return devices.map((d) {
        final map = d as Map<String, dynamic>;
        return DeviceItem(
          deviceId: map['device_id'] as String? ?? '',
          name: (map['name'] as String?)?.isNotEmpty == true
              ? map['name'] as String
              : map['device_id'] as String? ?? '',
          firmwareVersion: map['firmware_version'] as String? ?? '',
          status: _statusName(map['status'] as int? ?? 0),
          batteryPercent: map['battery_percent'] as int? ?? 0,
        );
      }).toList();
    } on DioException {
      return [];
    }
  }

  static String _statusName(int status) {
    switch (status) {
      case 1:
        return 'online';
      case 2:
        return 'offline';
      case 3:
        return 'measuring';
      case 4:
        return 'updating';
      case 5:
        return 'error';
      default:
        return 'unknown';
    }
  }

  @override
  Future<List<ConnectedDevice>> getConnectedDevices() async {
    // Mock implementation for REST - 10 Simulated Devices
    await Future.delayed(const Duration(milliseconds: 500));
    return [
      ConnectedDevice(
        id: 'gas-001',
        name: '거실 공기질 측정기',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 85,
        signalStrength: 92,
        currentValues: {'CO2': '450 ppm', 'VOC': '0.05 ppm', 'Radon': 'Safe'},
        latestReadings: [420, 430, 450, 440, 460, 450, 455, 450, 448, 452],
      ),
      ConnectedDevice(
        id: 'env-002',
        name: '안방 환경 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 90,
        signalStrength: 88,
        currentValues: {'Temp': '24.5°C', 'Humidity': '45%', 'Light': '300 lux'},
        latestReadings: [24.0, 24.1, 24.2, 24.5, 24.5, 24.4, 24.5, 24.6, 24.5, 24.5],
      ),
      ConnectedDevice(
        id: 'gas-002',
        name: '주방 가스 감지기',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 72,
        signalStrength: 75,
        currentValues: {'CO': '0 ppm', 'LNG': '0%', 'Smoke': 'None'},
        latestReadings: [0, 0, 0, 0, 0, 0, 1, 0, 0, 0],
      ),
      ConnectedDevice(
        id: 'bio-001',
        name: '바이오 카트리지 #1',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.disconnected,
        batteryLevel: 0,
        signalStrength: 0,
      ),
      // --- Additional 6 Devices ---
      ConnectedDevice(
        id: 'env-003',
        name: '아이방 온습도계',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 95,
        signalStrength: 98,
        currentValues: {'Temp': '23.0°C', 'Humidity': '50%'},
        latestReadings: [23, 23, 23, 23, 23],
      ),
      ConnectedDevice(
        id: 'gas-003',
        name: '베란다 환기 센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 60,
        signalStrength: 82,
        currentValues: {'Dust': '15 ug/m3'},
        latestReadings: [10, 12, 15, 14, 15],
      ),
      ConnectedDevice(
        id: 'bio-002',
        name: '웨어러블 밴드 Left',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 45,
        signalStrength: 90,
        currentValues: {'Pulse': '72 bpm', 'O2': '98%'},
        latestReadings: [70, 72, 71, 72, 75],
      ),
      ConnectedDevice(
        id: 'bio-003',
        name: '웨어러블 밴드 Right',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 42,
        signalStrength: 88,
        currentValues: {'Pulse': '73 bpm', 'Stress': 'Low'},
        latestReadings: [72, 73, 73, 74, 73],
      ),
      ConnectedDevice(
        id: 'env-004',
        name: '서재 조명 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.disconnected,
        batteryLevel: 10,
        signalStrength: 20,
      ),
      ConnectedDevice(
        id: 'gas-004',
        name: '차고 배기 센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 88,
        signalStrength: 65,
        currentValues: {'CO': '2 ppm'},
        latestReadings: [1, 1, 2, 2, 1],
      ),
    ];
  }
}
