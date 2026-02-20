import 'package:flutter/material.dart';
import 'dart:math' as math;
import 'dart:ui' as ui;
import 'package:manpasik/core/theme/app_theme.dart';

/// Jarvis Style Data Connector Line
/// 
/// Connects a [startOffset] to an [endOffset] with an animated
/// "Data Packet" traveling along the line.
class JarvisConnector extends StatefulWidget {
  final Offset startOffset;
  final Offset endOffset;
  final Color color;
  final bool isActive;

  const JarvisConnector({
    super.key,
    required this.startOffset,
    required this.endOffset,
    this.color = AppTheme.sanggamGold,
    this.isActive = true,
  });

  @override
  State<JarvisConnector> createState() => _JarvisConnectorState();
}

class _JarvisConnectorState extends State<JarvisConnector> with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 6), // Slow, elegant flow
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!widget.isActive) return const SizedBox();

    return IgnorePointer(
      child: AnimatedBuilder( // Fix: Wrap in AnimatedBuilder to trigger repaints
        animation: _controller,
        builder: (context, _) {
          return CustomPaint(
            painter: _ConnectorPainter(
              start: widget.startOffset,
              end: widget.endOffset,
              progress: _controller.value,
              color: widget.color,
            ),
            size: Size.infinite,
          );
        },
      ),
    );
  }
}

class _ConnectorPainter extends CustomPainter {
  final Offset start;
  final Offset end;
  final double progress;
  final Color color;

  _ConnectorPainter({
    required this.start,
    required this.end,
    required this.progress,
    required this.color,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final path = Path();
    path.moveTo(start.dx, start.dy);

    // Neural Curve (Organic Bezier)
    final controlPoint1 = Offset(start.dx, end.dy);
    final controlPoint2 = Offset(end.dx, start.dy);
    // Use a cubic bezier that feels like a nerve fiber
    // Vary control points based on distance for "tension"
    final midX = (start.dx + end.dx) / 2;
    path.cubicTo(
      midX, start.dy, 
      midX, end.dy, 
      end.dx, end.dy
    );

    // 1. Faint Synapse Trace (Background)
    final tracePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.0
      ..shader = ui.Gradient.linear(
        start, end,
        [color.withOpacity(0.0), color.withOpacity(0.2), color.withOpacity(0.0)],
        [0.0, 0.5, 1.0], // Fix: Add stops for 3 colors
      );
    canvas.drawPath(path, tracePaint);

    // 2. Terminals (Nodes)
    _drawNode(canvas, start, color);
    _drawNode(canvas, end, color);

    // 3. Active Neural Impulses (Data Packets)
    ui.PathMetrics pathMetrics = path.computeMetrics();
    for (ui.PathMetric metric in pathMetrics) {
      final length = metric.length;
      
      // Packet 1: Main fast impulse (Clean Loop)
      final val1 = progress; // 0.0 -> 1.0 (Start -> End)
      _drawImpulse(canvas, metric, val1 * length, color, 3.0, 1.5); // Size 3.0 (Small), Intensity 1.5

      // Packet 2: Slower/Offset impulse
      final val2 = (progress + 0.5) % 1.0; 
      _drawImpulse(canvas, metric, val2 * length, color.withOpacity(0.5), 2.0, 1.0); // Size 2.0
    }
  }

  void _drawNode(Canvas canvas, Offset pos, Color color) {
    // Outer Glow
    canvas.drawCircle(pos, 4.0, Paint() // Smaller node
      ..color = color.withOpacity(0.5)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 5)
    );
    // Core
    canvas.drawCircle(pos, 1.5, Paint()..color = Colors.white);
  }

  void _drawImpulse(Canvas canvas, ui.PathMetric metric, double dist, Color color, double size, double intensity) {
    final tangent = metric.getTangentForOffset(dist);
    if (tangent == null) return;
    
    final pos = tangent.position;
    
    // 1. Strong Optical Glow (Blurry Halo)
    canvas.drawCircle(pos, size * 3, Paint()
      ..color = color.withOpacity((0.4 * intensity).clamp(0.0, 1.0))
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8));

    // 2. Core (Glowing Dot)
    // White hot core
    canvas.drawCircle(pos, size * 0.8, Paint()
      ..color = Colors.white.withOpacity((1.0 * intensity).clamp(0.0, 1.0))
      ..maskFilter = const MaskFilter.blur(BlurStyle.solid, 2) // Slight soft edge
    );
    
    // Colored rim around core
    canvas.drawCircle(pos, size * 1.2, Paint()
      ..color = color.withOpacity((0.8 * intensity).clamp(0.0, 1.0))
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.5
    );

    // 3. Trail (Comet Tail)
    final trailLength = 60.0 * intensity;
    final trailPath = metric.extractPath(math.max(0, dist - trailLength), dist);
    canvas.drawPath(trailPath, Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = size * 0.6
      ..strokeCap = StrokeCap.round
      ..shader = ui.Gradient.linear(
         pos - (tangent.vector * trailLength), pos,
         [color.withOpacity(0), color.withOpacity((0.6 * intensity).clamp(0.0, 1.0))],
         [0.0, 1.0]
      )
    );
  }

  @override
  bool shouldRepaint(covariant _ConnectorPainter oldDelegate) => true;
}
