import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/shared/widgets/glass_morphism_card.dart';

// ───────────────────────────────────────────────────
// DeviceStatusCard — Sanggam Orbit 디바이스 상태 카드
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] Theme.of(context) isDark 분기 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

class DeviceStatusCard
    extends StatelessWidget {
  final ConnectedDevice device;
  final bool isSelected;
  final VoidCallback? onTap;

  const DeviceStatusCard({
    super.key,
    required this.device,
    this.isSelected = false,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final statusColor = device.status ==
            DeviceConnectionStatus.connected
        ? SanggamTheme.jagaeCyan
        : SanggamTheme.onSurfaceDim;

    return Padding(
      padding: const EdgeInsets.all(16),
      child: GestureDetector(
        onTap: onTap,
        child: GlassmorphismCard(
          useKoreanBorder: isSelected,
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment:
                CrossAxisAlignment.start,
            children: [
              // Header
              Row(
                mainAxisAlignment:
                    MainAxisAlignment
                        .spaceBetween,
                children: [
                  Container(
                    padding:
                        const EdgeInsets
                            .symmetric(
                            horizontal: 8,
                            vertical: 4),
                    decoration: BoxDecoration(
                      color: statusColor
                          .withValues(
                              alpha: 0.1),
                      borderRadius:
                          BorderRadius.circular(
                              8),
                      border: Border.all(
                          color: statusColor
                              .withValues(
                                  alpha:
                                      0.3)),
                    ),
                    child: Row(
                      children: [
                        Container(
                          width: 6,
                          height: 6,
                          decoration:
                              BoxDecoration(
                            color: statusColor,
                            shape:
                                BoxShape.circle,
                          ),
                        ),
                        const SizedBox(
                            width: 4),
                        Text(
                          device.status ==
                                  DeviceConnectionStatus
                                      .connected
                              ? '활성'
                              : '비활성',
                          style: TextStyle(
                            color: statusColor,
                            fontSize: 10,
                            fontWeight:
                                FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Icon(
                    device.type ==
                            DeviceType
                                .gasCartridge
                        ? Icons.cloud_outlined
                        : device.type ==
                                DeviceType
                                    .envCartridge
                            ? Icons
                                .thermostat_outlined
                            : device.type ==
                                    DeviceType
                                        .bioCartridge
                                ? Icons
                                    .science_outlined
                                : Icons
                                    .device_unknown,
                    color:
                        SanggamTheme.primary,
                    size: 20,
                  ),
                ],
              ),
              const SizedBox(height: 12),

              // Device Name
              Text(
                device.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
                maxLines: 1,
                overflow:
                    TextOverflow.ellipsis,
              ),
              Text(
                device.id,
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 10,
                ),
              ),

              const SizedBox(height: 12),

              // Real-time Chart
              SizedBox(
                height: 40,
                width: double.infinity,
                child: CustomPaint(
                  painter:
                      _MiniSparklinePainter(
                    data: device
                        .latestReadings,
                    color: SanggamTheme
                        .jagaeCyan,
                  ),
                ),
              ),
              const SizedBox(height: 12),

              // Values Grid
              Wrap(
                spacing: 8,
                runSpacing: 4,
                children: device
                    .currentValues.entries
                    .map((e) {
                  return Container(
                    padding: const EdgeInsets
                        .symmetric(
                        horizontal: 6,
                        vertical: 2),
                    decoration: BoxDecoration(
                      color: Colors.white
                          .withValues(
                              alpha: 0.1),
                      borderRadius:
                          BorderRadius.circular(
                              4),
                    ),
                    child: Text(
                      '${e.key}: ${e.value}',
                      style: const TextStyle(
                        fontSize: 10,
                        color: SanggamTheme
                            .onSurfaceDim,
                      ),
                    ),
                  );
                }).toList(),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _MiniSparklinePainter
    extends CustomPainter {
  final List<double> data;
  final Color color;

  _MiniSparklinePainter(
      {required this.data,
      required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;

    final paint = Paint()
      ..color = color
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke;

    final path = Path();
    final stepX =
        size.width / (data.length - 1);

    final minVal = data.reduce(
        (curr, next) =>
            curr < next ? curr : next);
    final maxVal = data.reduce(
        (curr, next) =>
            curr > next ? curr : next);
    final range = maxVal - minVal;

    for (int i = 0; i < data.length; i++) {
      final x = i * stepX;
      final normalized = range == 0
          ? 0.5
          : (data[i] - minVal) / range;
      final y = size.height -
          (normalized * size.height);

      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(
          covariant
          _MiniSparklinePainter
              oldDelegate) =>
      true;
}
