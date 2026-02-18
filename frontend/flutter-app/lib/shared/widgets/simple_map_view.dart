import 'dart:math';
import 'package:flutter/material.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 간이 지도 뷰 (C4/C9)
///
/// 외부 지도 SDK 없이 CustomPainter로 좌표 마커를 표시합니다.
/// Google Maps / Naver Maps SDK 연동 시 이 위젯을 교체합니다.
class SimpleMapView extends StatelessWidget {
  const SimpleMapView({
    super.key,
    required this.markers,
    this.height = 250,
    this.centerLatitude = 37.5665,
    this.centerLongitude = 126.9780,
    this.zoomLevel = 14,
  });

  final List<MapMarker> markers;
  final double height;
  final double centerLatitude;
  final double centerLongitude;
  final int zoomLevel;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return ClipRRect(
      borderRadius: BorderRadius.circular(16),
      child: Container(
        height: height,
        width: double.infinity,
        color: theme.colorScheme.surfaceContainerHighest.withOpacity(0.3),
        child: Stack(
          children: [
            // 배경 그리드
            CustomPaint(
              size: Size.infinite,
              painter: _GridPainter(
                lineColor: theme.colorScheme.outlineVariant.withOpacity(0.3),
              ),
            ),

            // 중심 좌표 표시
            Center(
              child: Icon(
                Icons.add,
                color: theme.colorScheme.outlineVariant.withOpacity(0.5),
                size: 20,
              ),
            ),

            // 마커들
            ...markers.map((m) {
              final dx = (m.longitude - centerLongitude) * 5000;
              final dy = -(m.latitude - centerLatitude) * 5000;
              return Positioned(
                left: height / 2 + dx - 12, // 근사 위치
                top: height / 2 + dy - 24,
                child: _MarkerWidget(marker: m),
              );
            }),

            // 좌표 정보
            Positioned(
              bottom: 8,
              right: 8,
              child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: theme.colorScheme.surface.withOpacity(0.8),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '${centerLatitude.toStringAsFixed(4)}, ${centerLongitude.toStringAsFixed(4)}',
                  style: theme.textTheme.bodySmall?.copyWith(fontSize: 10),
                ),
              ),
            ),

            // SDK 안내
            Positioned(
              bottom: 8,
              left: 8,
              child: Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: theme.colorScheme.surface.withOpacity(0.8),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  '간이 지도 (SDK 연동 시 교체)',
                  style: theme.textTheme.bodySmall?.copyWith(
                    fontSize: 9,
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _MarkerWidget extends StatelessWidget {
  const _MarkerWidget({required this.marker});
  final MapMarker marker;

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: marker.label,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            marker.icon ?? Icons.location_on,
            color: marker.color ?? AppTheme.dancheongRed,
            size: 24,
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
            decoration: BoxDecoration(
              color: Colors.black54,
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              marker.label,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 8,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _GridPainter extends CustomPainter {
  _GridPainter({required this.lineColor});
  final Color lineColor;

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = lineColor
      ..strokeWidth = 0.5;

    const spacing = 30.0;
    for (var x = 0.0; x < size.width; x += spacing) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (var y = 0.0; y < size.height; y += spacing) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

class MapMarker {
  final String label;
  final double latitude;
  final double longitude;
  final IconData? icon;
  final Color? color;
  final Map<String, dynamic>? data;

  const MapMarker({
    required this.label,
    required this.latitude,
    required this.longitude,
    this.icon,
    this.color,
    this.data,
  });
}
