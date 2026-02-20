/// 디바이스 Repository 인터페이스
abstract class DeviceRepository {
  Future<List<DeviceItem>> listDevices(String userId);

  // New: Monitoring
  Future<List<ConnectedDevice>> getConnectedDevices();
}

/// 기기 연결 상태
enum DeviceConnectionStatus { connected, disconnected, scanning }

/// 기기 타입
enum DeviceType { gasCartridge, envCartridge, bioCartridge, unknown }

class ConnectedDevice {
  final String id;
  final String name;
  final DeviceType type;
  final DeviceConnectionStatus status;
  final int batteryLevel; // 0-100
  final int signalStrength; // 0-100
  final List<double> latestReadings; // For sparkline
  final Map<String, dynamic> currentValues; // e.g. {'CO2': 450, 'Temp': 24.5}

  const ConnectedDevice({
    required this.id,
    required this.name,
    required this.type,
    required this.status,
    required this.batteryLevel,
    required this.signalStrength,
    this.latestReadings = const [],
    this.currentValues = const {},
  });
}

/// 디바이스 목록 항목
class DeviceItem {
  final String deviceId;
  final String name;
  final String firmwareVersion;
  final String status; // online, offline, measuring, updating, error
  final int batteryPercent;
  final DateTime? lastSeen;

  const DeviceItem({
    required this.deviceId,
    required this.name,
    required this.firmwareVersion,
    required this.status,
    this.batteryPercent = 0,
    this.lastSeen,
  });
}
