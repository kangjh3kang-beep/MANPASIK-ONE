import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

// ───────────────────────────────────────────────────
// MiniDeviceCard — Sanggam Orbit 미니 디바이스 카드
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] Theme.of(context) isDark 분기 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// Mini Device Card — DeviceStatusCard 소형 버전 (120x160px)
class MiniDeviceCard extends StatelessWidget {
  final ConnectedDevice device;
  final VoidCallback? onTap;

  const MiniDeviceCard(
      {super.key,
      required this.device,
      this.onTap});

  @override
  Widget build(BuildContext context) {
    final isConnected = device.status ==
        DeviceConnectionStatus.connected;
    final statusColor = isConnected
        ? SanggamTheme.jagaeCyan
        : SanggamTheme.onSurfaceDim;

    return Tooltip(
      message: '${device.name} 상세 보기',
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        child: GestureDetector(
          onTap: onTap,
          child: Container(
            width: 120,
            height: 160,
            margin: const EdgeInsets.only(
                right: 8),
            padding:
                const EdgeInsets.all(10),
            decoration: BoxDecoration(
              color: Colors.white
                  .withValues(
                      alpha: 0.06),
              borderRadius:
                  BorderRadius.circular(
                      14),
              border: Border.all(
                color: Colors.white
                    .withValues(
                        alpha: 0.08),
              ),
            ),
            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,
              children: [
                // Status + Type Icon
                Row(
                  mainAxisAlignment:
                      MainAxisAlignment
                          .spaceBetween,
                  children: [
                    Container(
                      padding:
                          const EdgeInsets
                              .symmetric(
                              horizontal: 5,
                              vertical: 2),
                      decoration:
                          BoxDecoration(
                        color: statusColor
                            .withValues(
                                alpha:
                                    0.15),
                        borderRadius:
                            BorderRadius
                                .circular(
                                    6),
                      ),
                      child: Row(
                        mainAxisSize:
                            MainAxisSize
                                .min,
                        children: [
                          Container(
                            width: 5,
                            height: 5,
                            decoration:
                                BoxDecoration(
                              color:
                                  statusColor,
                              shape: BoxShape
                                  .circle,
                            ),
                          ),
                          const SizedBox(
                              width: 3),
                          Text(
                            isConnected
                                ? '활성'
                                : '비활성',
                            style:
                                TextStyle(
                              color:
                                  statusColor,
                              fontSize: 8,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            ),
                          ),
                        ],
                      ),
                    ),
                    Icon(
                      _deviceIcon(
                          device.type),
                      color: SanggamTheme
                          .primary,
                      size: 14,
                    ),
                  ],
                ),
                const SizedBox(height: 6),

                // Device Name
                Text(
                  device.name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight:
                        FontWeight.w600,
                    fontSize: 11,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow
                      .ellipsis,
                ),
                const SizedBox(height: 4),

                // Mini Sparkline
                Expanded(
                  child: SizedBox(
                    width: double.infinity,
                    child: CustomPaint(
                      painter:
                          _MiniSparkPainter(
                        data: device
                            .latestReadings,
                        color: SanggamTheme
                            .jagaeCyan,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 4),

                // Current Values (first 2)
                if (device.currentValues
                    .isNotEmpty)
                  ...device
                      .currentValues
                      .entries
                      .take(2)
                      .map((e) => Text(
                            '${e.key}: ${e.value}',
                            style:
                                const TextStyle(
                              fontSize: 8,
                              color: SanggamTheme
                                  .onSurfaceDim,
                            ),
                            maxLines: 1,
                            overflow:
                                TextOverflow
                                    .ellipsis,
                          )),

                // Battery
                const SizedBox(height: 2),
                Row(
                  children: [
                    Icon(
                      device.batteryLevel <
                              20
                          ? Icons
                              .battery_alert
                          : Icons
                              .battery_full,
                      size: 10,
                      color: device
                                  .batteryLevel <
                              20
                          ? SanggamTheme
                              .error
                          : SanggamTheme
                              .jagaeCyan,
                    ),
                    const SizedBox(
                        width: 2),
                    Text(
                      '${device.batteryLevel}%',
                      style: TextStyle(
                        fontSize: 8,
                        color: device
                                    .batteryLevel <
                                20
                            ? SanggamTheme
                                .error
                            : SanggamTheme
                                .onSurfaceDim,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  IconData _deviceIcon(DeviceType type) {
    switch (type) {
      case DeviceType.gasCartridge:
        return Icons.cloud_outlined;
      case DeviceType.envCartridge:
        return Icons.thermostat_outlined;
      case DeviceType.bioCartridge:
        return Icons.science_outlined;
      case DeviceType.unknown:
        return Icons.device_unknown;
    }
  }
}

class _MiniSparkPainter
    extends CustomPainter {
  final List<double> data;
  final Color color;

  _MiniSparkPainter(
      {required this.data,
      required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    if (data.length < 2) return;

    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    final path = Path();
    final stepX =
        size.width / (data.length - 1);
    final minVal = data
        .reduce((a, b) => a < b ? a : b);
    final maxVal = data
        .reduce((a, b) => a > b ? a : b);
    final range = maxVal - minVal;

    for (int i = 0; i < data.length; i++) {
      final x = i * stepX;
      final norm = range == 0
          ? 0.5
          : (data[i] - minVal) / range;
      final y = size.height -
          (norm * size.height);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    canvas.drawPath(path, paint);

    // Fill gradient
    final fillPath = Path.from(path)
      ..lineTo(size.width, size.height)
      ..lineTo(0, size.height)
      ..close();
    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          color.withValues(alpha: 0.2),
          color.withValues(alpha: 0.0),
        ],
      ).createShader(Rect.fromLTWH(
          0, 0, size.width, size.height));
    canvas.drawPath(fillPath, fillPaint);
  }

  @override
  bool shouldRepaint(
          covariant _MiniSparkPainter
              oldDelegate) =>
      true;
}
