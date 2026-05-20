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
    // REST API를 통해 등록된 디바이스 조회 시도
    try {
      final devices = await listDevices('');
      if (devices.isNotEmpty) {
        return devices
            .map((d) => ConnectedDevice(
                  id: d.deviceId,
                  name: d.name,
                  type: DeviceType.bioCartridge,
                  status: d.status == 'online'
                      ? DeviceConnectionStatus.connected
                      : DeviceConnectionStatus.disconnected,
                  batteryLevel: d.batteryPercent,
                  signalStrength: d.status == 'online' ? 85 : 0,
                ))
            .toList();
      }
    } on DioException {
      // API 실패 시 DEMO 데이터 사용 (H6 실패격리)
    }

    // DEMO: 서버 미연결 시 시뮬레이션 디바이스 반환
    return _demoDevices();
  }

  /// DEMO 모드용 시뮬레이션 디바이스 목록 (서버 미연결 시 사용)
  static List<ConnectedDevice> _demoDevices() {
    return [
      ConnectedDevice(
        id: 'demo-bio-001',
        name: '[DEMO] 바이오 카트리지',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 85,
        signalStrength: 92,
        currentValues: {'Glucose': '95 mg/dL'},
      ),
      ConnectedDevice(
        id: 'demo-env-001',
        name: '[DEMO] 환경 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 90,
        signalStrength: 88,
        currentValues: {'Temp': '24.5°C', 'Humidity': '45%'},
      ),
      ConnectedDevice(
        id: 'demo-gas-001',
        name: '[DEMO] 가스 감지기',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.disconnected,
        batteryLevel: 0,
        signalStrength: 0,
      ),
    ];
  }
}
