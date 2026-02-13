/// 디바이스 Repository 인터페이스
abstract class DeviceRepository {
  Future<List<DeviceItem>> listDevices(String userId);
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
