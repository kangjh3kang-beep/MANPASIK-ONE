import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/shared/widgets/holo_body.dart' show HoloGender;

/// 30초 폴링 자동갱신 ConnectedDevices StreamProvider
final pollingConnectedDevicesProvider = StreamProvider<List<ConnectedDevice>>((ref) {
  final controller = StreamController<List<ConnectedDevice>>();

  Future<void> fetch() async {
    try {
      final repo = ref.read(deviceRepositoryProvider);
      final devices = await repo.getConnectedDevices();
      if (!controller.isClosed) {
        controller.add(devices);
      }
    } catch (e) {
      if (!controller.isClosed) {
        controller.addError(e);
      }
    }
  }

  // 즉시 1회 로드
  fetch();

  // 30초마다 반복 로드
  final timer = Timer.periodic(const Duration(seconds: 30), (_) => fetch());

  ref.onDispose(() {
    timer.cancel();
    controller.close();
  });

  return controller.stream;
});

/// 선택된 기기 ID (Dashboard ↔ Bubble 공유)
final selectedDeviceIdProvider = StateProvider<String?>((ref) => null);

/// 선택된 기기 객체 (파생)
final selectedDeviceProvider = Provider<ConnectedDevice?>((ref) {
  final selectedId = ref.watch(selectedDeviceIdProvider);
  if (selectedId == null) return null;
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      try {
        return devices.firstWhere((d) => d.id == selectedId);
      } catch (_) {
        return null;
      }
    },
    loading: () => null,
    error: (_, __) => null,
  );
});

/// 모니터링 필터 탭 (0=All, 1=Gas, 2=Env, 3=Bio)
final monitoringFilterTabProvider = StateProvider<int>((ref) => 0);

/// 필터링된 기기 목록 (파생)
final filteredDevicesProvider = Provider<AsyncValue<List<ConnectedDevice>>>((ref) {
  final tabIndex = ref.watch(monitoringFilterTabProvider);
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.whenData((devices) {
    switch (tabIndex) {
      case 1:
        return devices.where((d) => d.type == DeviceType.gasCartridge).toList();
      case 2:
        return devices.where((d) => d.type == DeviceType.envCartridge).toList();
      case 3:
        return devices.where((d) => d.type == DeviceType.bioCartridge).toList();
      default:
        return devices;
    }
  });
});

/// 요약 통계 (파생)
final monitoringSummaryProvider =
    Provider<({int total, int connected, int alerts, int avgBattery})>((ref) {
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      if (devices.isEmpty) {
        return (total: 0, connected: 0, alerts: 0, avgBattery: 0);
      }
      final connected =
          devices.where((d) => d.status == DeviceConnectionStatus.connected).length;
      final alerts = devices
          .where((d) =>
              d.status == DeviceConnectionStatus.disconnected || d.batteryLevel < 20)
          .length;
      final avgBattery =
          (devices.fold<int>(0, (sum, d) => sum + d.batteryLevel) / devices.length).round();
      return (total: devices.length, connected: connected, alerts: alerts, avgBattery: avgBattery);
    },
    loading: () => (total: 0, connected: 0, alerts: 0, avgBattery: 0),
    error: (_, __) => (total: 0, connected: 0, alerts: 0, avgBattery: 0),
  );
});

/// 경고 기기 목록 (disconnected 또는 배터리 < 20%)
final alertDevicesProvider = Provider<List<ConnectedDevice>>((ref) {
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) => devices
        .where((d) =>
            d.status == DeviceConnectionStatus.disconnected || d.batteryLevel < 20)
        .toList(),
    loading: () => [],
    error: (_, __) => [],
  );
});

/// HoloBody 성별 토글 (바이오 탭 전용) — HoloGender enum은 holo_body.dart에서 import
final holoGenderProvider = StateProvider<HoloGender>((ref) => HoloGender.male);

/// 선택된 바이오 기기의 생체 데이터 (파생)
final selectedBioDataProvider = Provider<Map<String, dynamic>>((ref) {
  final device = ref.watch(selectedDeviceProvider);
  if (device != null && device.type == DeviceType.bioCartridge) {
    return device.currentValues;
  }
  final devicesAsync = ref.watch(pollingConnectedDevicesProvider);
  return devicesAsync.when(
    data: (devices) {
      final bio = devices
          .where((d) =>
              d.type == DeviceType.bioCartridge &&
              d.status == DeviceConnectionStatus.connected)
          .toList();
      return bio.isNotEmpty ? bio.first.currentValues : <String, dynamic>{};
    },
    loading: () => <String, dynamic>{},
    error: (_, __) => <String, dynamic>{},
  );
});
