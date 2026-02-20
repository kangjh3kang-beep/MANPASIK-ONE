import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';

class DeviceStatusCard extends StatelessWidget {
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
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final statusColor = device.status == DeviceConnectionStatus.connected
        ? const Color(0xFF00E676)
        : Colors.grey;

    return Padding(
      padding: const EdgeInsets.all(16),
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          decoration: BoxDecoration(
            color: isDark ? const Color(0xFF1E2832) : Colors.white,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: isSelected 
                  ? AppTheme.sanggamGold 
                  : (isDark ? Colors.white10 : Colors.black12),
              width: isSelected ? 1.5 : 0.5,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(isDark ? 0.3 : 0.05),
                blurRadius: 10,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: ClipRRect(
            borderRadius: BorderRadius.circular(16),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Header
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: statusColor.withOpacity(0.1),
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(color: statusColor.withOpacity(0.3)),
                        ),
                        child: Row(
                          children: [
                            Container(
                              width: 6,
                              height: 6,
                              decoration: BoxDecoration(
                                  color: statusColor, shape: BoxShape.circle),
                            ),
                            const SizedBox(width: 4),
                            Text(
                              device.status == DeviceConnectionStatus.connected
                                  ? 'LIVE'
                                  : 'OFF',
                              style: TextStyle(
                                  color: statusColor,
                                  fontSize: 10,
                                  fontWeight: FontWeight.bold),
                            ),
                          ],
                        ),
                      ),
                      Icon(
                        device.type == DeviceType.gasCartridge
                            ? Icons.cloud_outlined
                            : device.type == DeviceType.envCartridge
                                ? Icons.thermostat_outlined
                                : device.type == DeviceType.bioCartridge
                                    ? Icons.science_outlined
                                    : Icons.device_unknown,
                        color: AppTheme.sanggamGold,
                        size: 20,
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),

                  // Device Name
                  Text(
                    device.name,
                    style: TextStyle(
                      color: isDark ? Colors.white : const Color(0xFF1A1A1A),
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  Text(
                    device.id,
                    style: TextStyle(
                      color: isDark ? Colors.white38 : Colors.grey,
                      fontSize: 10,
                    ),
                  ),

                  const SizedBox(height: 12),

                  // Real-time Chart
                  SizedBox(
                    height: 40,
                    width: double.infinity,
                    child: CustomPaint(
                      painter: _MiniSparklinePainter(
                        data: device.latestReadings,
                        color: isDark ? AppTheme.waveCyan : const Color(0xFF00ACC1),
                      ),
                    ),
                  ),
                  const SizedBox(height: 12),

                  // Values Grid
                  Wrap(
                    spacing: 8,
                    runSpacing: 4,
                    children: device.currentValues.entries.map((e) {
                      return Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: isDark ? Colors.white10 : Colors.grey[200],
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${e.key}: ${e.value}',
                          style: TextStyle(
                            fontSize: 10,
                            color: isDark ? Colors.white70 : Colors.black87,
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _MiniSparklinePainter extends CustomPainter {
  final List<double> data;
  final Color color;

  _MiniSparklinePainter({required this.data, required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;

    final paint = Paint()
      ..color = color
      ..strokeWidth = 2
      ..style = PaintingStyle.stroke;

    final path = Path();
    final stepX = size.width / (data.length - 1);

    final minVal = data.reduce((curr, next) => curr < next ? curr : next);
    final maxVal = data.reduce((curr, next) => curr > next ? curr : next);
    final range = maxVal - minVal;

    for (int i = 0; i < data.length; i++) {
      final x = i * stepX;
      final normalized = range == 0 ? 0.5 : (data[i] - minVal) / range;
      final y = size.height - (normalized * size.height);

      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }
    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(covariant _MiniSparklinePainter oldDelegate) => true;
}
