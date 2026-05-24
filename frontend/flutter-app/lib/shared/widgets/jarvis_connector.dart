import 'dart:math' as math;
import 'package:flutter/material.dart';

/// 사이버네틱 무드의 연결선 (Smart Anchor 점선 및 회로 스타일)
class JarvisConnector extends StatelessWidget {
  final Offset startOffset;
  final Offset endOffset;
  final Color color;

  const JarvisConnector({
    super.key,
    required this.startOffset,
    required this.endOffset,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Positioned.fill(
      child: IgnorePointer(
        child: CustomPaint(
          painter: _JarvisConnectorPainter(
            startOffset: startOffset,
            endOffset: endOffset,
            color: color,
          ),
        ),
      ),
    );
  }
}

class _JarvisConnectorPainter extends CustomPainter {
  final Offset startOffset;
  final Offset endOffset;
  final Color color;

  _JarvisConnectorPainter({
    required this.startOffset,
    required this.endOffset,
    required this.color,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    // 회로처럼 L 자 형태로 꺾이게 그리기 또는 단순 곡선
    final path = Path()..moveTo(startOffset.dx, startOffset.dy);
    final midX = (startOffset.dx + endOffset.dx) / 2;
    
    path.cubicTo(
      midX, startOffset.dy,
      midX, endOffset.dy,
      endOffset.dx, endOffset.dy,
    );

    // 점선(Dotted) 처리
    const double dashWidth = 4;
    const double dashSpace = 4;
    final metric = path.computeMetrics().first;
    final double length = metric.length;
    double cur = 0;
    
    final Path dashedPath = Path();
    while (cur < length) {
      final double next = cur + dashWidth;
      dashedPath.addPath(metric.extractPath(cur, next), Offset.zero);
      cur = next + dashSpace;
    }
    
    canvas.drawPath(dashedPath, paint);

    // 시작점 (신체 앵커)
    canvas.drawCircle(startOffset, 4, Paint()..color = color.withValues(alpha: 0.8));
    canvas.drawCircle(startOffset, 1.5, Paint()..color = Colors.white);
    
    // 끝점 (HUD 앵커)
    canvas.drawCircle(endOffset, 3, Paint()..color = color.withValues(alpha: 0.5));
  }

  @override
  bool shouldRepaint(covariant _JarvisConnectorPainter oldDelegate) =>
      oldDelegate.startOffset != startOffset ||
      oldDelegate.endOffset != endOffset ||
      oldDelegate.color != color;
}
