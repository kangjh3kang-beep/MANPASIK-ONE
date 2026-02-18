import 'package:flutter/material.dart';

import 'package:manpasik/shared/widgets/simple_map_view.dart';

/// 기기 위치 맵 섹션 (C9)
///
/// DeviceListScreen에서 등록된 기기들의 위치를 지도 위에 표시합니다.
class DeviceMapSection extends StatelessWidget {
  const DeviceMapSection({
    super.key,
    required this.devices,
  });

  final List<DeviceLocation> devices;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (devices.isEmpty) {
      return Card(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Center(
            child: Text(
              '위치 정보가 있는 기기가 없습니다',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ),
        ),
      );
    }

    // 중심 좌표 계산
    final avgLat = devices.map((d) => d.latitude).reduce((a, b) => a + b) /
        devices.length;
    final avgLng = devices.map((d) => d.longitude).reduce((a, b) => a + b) /
        devices.length;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Icon(Icons.map_rounded,
                size: 18, color: theme.colorScheme.primary),
            const SizedBox(width: 8),
            Text(
              '기기 위치',
              style: theme.textTheme.titleSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const Spacer(),
            Text(
              '${devices.length}대',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.primary,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        SimpleMapView(
          height: 200,
          centerLatitude: avgLat,
          centerLongitude: avgLng,
          markers: devices
              .map((d) => MapMarker(
                    label: d.deviceName,
                    latitude: d.latitude,
                    longitude: d.longitude,
                    icon: Icons.bluetooth_connected,
                    color: d.isOnline ? Colors.green : Colors.grey,
                  ))
              .toList(),
        ),
      ],
    );
  }
}

class DeviceLocation {
  final String deviceId;
  final String deviceName;
  final double latitude;
  final double longitude;
  final bool isOnline;

  const DeviceLocation({
    required this.deviceId,
    required this.deviceName,
    required this.latitude,
    required this.longitude,
    this.isOnline = false,
  });
}
