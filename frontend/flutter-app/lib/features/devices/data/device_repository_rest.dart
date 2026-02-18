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
}
