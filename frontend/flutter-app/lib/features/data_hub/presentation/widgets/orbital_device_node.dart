import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

// ───────────────────────────────────────────────────
// OrbitalDeviceNode — Sanggam Orbit 궤도 디바이스 노드
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] isDark 파라미터 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan 2x → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF1A1A1A) → SanggamTheme.surface
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// ───────────────────────────────────────────────────

class OrbitalDeviceNode
    extends StatefulWidget {
  final ConnectedDevice device;
  final double angle;
  final bool isSelected;
  final VoidCallback? onSelect;

  const OrbitalDeviceNode({
    super.key,
    required this.device,
    required this.angle,
    this.isSelected = false,
    this.onSelect,
  });

  @override
  State<OrbitalDeviceNode> createState() =>
      _OrbitalDeviceNodeState();
}

class _OrbitalDeviceNodeState
    extends State<OrbitalDeviceNode>
    with SingleTickerProviderStateMixin {
  bool _isExpanded = false;

  void _toggleExpand() {
    setState(
        () => _isExpanded = !_isExpanded);
  }

  @override
  Widget build(BuildContext context) {
    final effectiveExpanded =
        _isExpanded || widget.isSelected;

    return MouseRegion(
      onEnter: (_) =>
          setState(() => _isExpanded = true),
      onExit: (_) =>
          setState(() => _isExpanded = false),
      child: GestureDetector(
        onTap: () {
          if (widget.onSelect != null) {
            widget.onSelect!();
          } else {
            _toggleExpand();
          }
        },
        child: AnimatedContainer(
          duration: const Duration(
              milliseconds: 300),
          curve: Curves.easeOutBack,
          width:
              effectiveExpanded ? 150 : 40,
          height:
              effectiveExpanded ? 90 : 40,
          clipBehavior: Clip.hardEdge,
          decoration: BoxDecoration(
            color: SanggamTheme.surface
                .withValues(alpha: 0.9),
            borderRadius:
                BorderRadius.circular(
                    effectiveExpanded
                        ? 10
                        : 20),
            border: Border.all(
              color: widget.isSelected
                  ? SanggamTheme.jagaeCyan
                  : (widget.device.status ==
                          DeviceConnectionStatus
                              .connected
                      ? SanggamTheme.primary
                      : SanggamTheme
                          .onSurfaceDim),
              width: widget.isSelected
                  ? 2.0
                  : (effectiveExpanded
                      ? 1.5
                      : 2.0),
            ),
            boxShadow: [
              BoxShadow(
                color: (widget.isSelected
                        ? SanggamTheme
                            .jagaeCyan
                        : (widget.device
                                    .status ==
                                DeviceConnectionStatus
                                    .connected
                            ? SanggamTheme
                                .primary
                            : SanggamTheme
                                .onSurfaceDim))
                    .withValues(
                        alpha:
                            widget.isSelected
                                ? 0.5
                                : 0.3),
                blurRadius:
                    widget.isSelected
                        ? 20
                        : (effectiveExpanded
                            ? 15
                            : 5),
                spreadRadius:
                    widget.isSelected
                        ? 2
                        : 1,
              ),
            ],
          ),
          child: SingleChildScrollView(
            physics:
                const NeverScrollableScrollPhysics(),
            child: SizedBox(
              height: effectiveExpanded
                  ? 90
                  : 40,
              child: effectiveExpanded
                  ? _buildExpandedContent()
                  : _buildCollapsedContent(),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildCollapsedContent() {
    return Center(
      child: Icon(
        widget.device.type ==
                DeviceType.gasCartridge
            ? Icons.cloud_outlined
            : widget.device.type ==
                    DeviceType.envCartridge
                ? Icons.thermostat
                : widget.device.type ==
                        DeviceType
                            .bioCartridge
                    ? Icons.science
                    : Icons.device_unknown,
        color: Colors.white,
        size: 18,
      ),
    );
  }

  Widget _buildExpandedContent() {
    final statusColor = widget.device
                .status ==
            DeviceConnectionStatus.connected
        ? SanggamTheme.jagaeCyan
        : SanggamTheme.onSurfaceDim;

    return Padding(
      padding: const EdgeInsets.all(10),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            mainAxisAlignment:
                MainAxisAlignment
                    .spaceBetween,
            children: [
              Flexible(
                child: Text(
                  widget.device.name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight:
                        FontWeight.bold,
                    fontSize: 11,
                  ),
                  maxLines: 1,
                  overflow:
                      TextOverflow.ellipsis,
                ),
              ),
              const Icon(
                Icons.arrow_forward_ios,
                size: 8,
                color: SanggamTheme
                    .onSurfaceDim,
              ),
            ],
          ),
          const SizedBox(height: 2),
          Text(
            widget.device.id,
            style: const TextStyle(
              color:
                  SanggamTheme.onSurfaceDim,
              fontSize: 8,
            ),
          ),
          const Spacer(),
          Row(
            children: [
              Container(
                width: 5,
                height: 5,
                decoration: BoxDecoration(
                  color: statusColor,
                  shape: BoxShape.circle,
                ),
              ),
              const SizedBox(width: 4),
              Text(
                widget.device.status ==
                        DeviceConnectionStatus
                            .connected
                    ? '활성'
                    : '비활성',
                style: TextStyle(
                  color: statusColor,
                  fontSize: 9,
                  fontWeight:
                      FontWeight.bold,
                ),
              ),
              const Spacer(),
              if (widget.device
                  .latestReadings.isNotEmpty)
                Text(
                  '${widget.device.latestReadings.last.toInt()}',
                  style: const TextStyle(
                    color: SanggamTheme
                        .primary,
                    fontSize: 10,
                    fontWeight:
                        FontWeight.bold,
                  ),
                ),
            ],
          ),
        ],
      ),
    );
  }
}

class OrbitalConnectorPainter
    extends CustomPainter {
  final Offset start;
  final Offset end;
  final Color color;
  final double animationValue;

  OrbitalConnectorPainter({
    required this.start,
    required this.end,
    required this.color,
    required this.animationValue,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final double distance =
        (end - start).distance;
    if (distance == 0) return;

    final Offset mid = (start + end) / 2;
    final Offset direction = end - start;
    final double curveDirection =
        (end.dx > start.dx) ? 1.0 : -1.0;
    final Offset normal = Offset(
            -direction.dy,
            direction.dx *
                curveDirection) /
        distance;
    final Offset control =
        mid + normal * (distance * 0.25);

    final Path path = Path()
      ..moveTo(start.dx, start.dy)
      ..quadraticBezierTo(control.dx,
          control.dy, end.dx, end.dy);

    // Dashed Base Line
    final paint = Paint()
      ..color =
          color.withValues(alpha: 0.3)
      ..strokeWidth = 1.0
      ..style = PaintingStyle.stroke;

    const double dashWidth = 3;
    const double dashSpace = 4;

    final varyMetric =
        path.computeMetrics().first;
    final double length = varyMetric.length;
    double currentDistance = 0;

    final Path dashedPath = Path();
    while (currentDistance < length) {
      final double nextDistance =
          currentDistance + dashWidth;
      dashedPath.addPath(
        varyMetric.extractPath(
            currentDistance, nextDistance),
        Offset.zero,
      );
      currentDistance =
          nextDistance + dashSpace;
    }
    canvas.drawPath(dashedPath, paint);

    // Moving Luminous Point
    if (color.a > 0) {
      final double pointMetricPos =
          animationValue * length;
      final ui.Tangent? tangent = varyMetric
          .getTangentForOffset(
              pointMetricPos);

      if (tangent != null) {
        final Offset pointPos =
            tangent.position;

        // Glow
        canvas.drawCircle(
          pointPos,
          4.0,
          Paint()
            ..color = color.withValues(
                alpha: 0.6)
            ..maskFilter =
                const MaskFilter.blur(
                    BlurStyle.normal, 3),
        );

        // Core
        canvas.drawCircle(pointPos, 2.0,
            Paint()..color = Colors.white);

        // Trail
        final double trailPos =
            (pointMetricPos - 10)
                .clamp(0, length);
        final ui.Tangent? trailTangent =
            varyMetric.getTangentForOffset(
                trailPos);
        if (trailTangent != null) {
          canvas.drawCircle(
            trailTangent.position,
            1.5,
            Paint()
              ..color = color.withValues(
                  alpha: 0.3),
          );
        }
      }
    }

    // Anchor Points
    canvas.drawCircle(
        start,
        2.0,
        Paint()
          ..color = color.withValues(
              alpha: 0.5));
    canvas.drawCircle(
        end,
        2.0,
        Paint()
          ..color = color.withValues(
              alpha: 0.5));
  }

  @override
  bool shouldRepaint(
          covariant
          OrbitalConnectorPainter
              oldDelegate) =>
      oldDelegate.animationValue !=
          animationValue ||
      oldDelegate.start != start ||
      oldDelegate.end != end;
}
