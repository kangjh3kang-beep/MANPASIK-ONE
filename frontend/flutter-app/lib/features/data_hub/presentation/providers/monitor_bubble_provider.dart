import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitoring_providers.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

/// Floating Monitor Bubble 확장/축소 상태
final monitorBubbleExpandedProvider = StateProvider<bool>((ref) => false);

/// 연결 기기 수 (connected / total) — pollingProvider 연동
final connectedCountProvider = Provider<({int connected, int total})>((ref) {
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      final connected = devices.where((d) => d.status == DeviceConnectionStatus.connected).length;
      return (connected: connected, total: devices.length);
    },
    loading: () => (connected: 0, total: 0),
    error: (_, __) => (connected: 0, total: 0),
  );
});

/// 경고 기기 수 (disconnected 또는 배터리 < 20%) — pollingProvider 연동
final deviceAlertCountProvider = Provider<int>((ref) {
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      return devices.where((d) =>
        d.status == DeviceConnectionStatus.disconnected ||
        d.batteryLevel < 20
      ).length;
    },
    loading: () => 0,
    error: (_, __) => 0,
  );
});
