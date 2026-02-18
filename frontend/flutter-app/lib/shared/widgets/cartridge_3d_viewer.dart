import 'dart:math';
import 'package:flutter/material.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 카트리지 360도 뷰어 (C10)
///
/// CustomPainter + GestureDetector로 카트리지를 3D-like로 렌더링합니다.
/// 드래그로 회전, 핀치로 확대/축소를 지원합니다.
class Cartridge3DViewer extends StatefulWidget {
  const Cartridge3DViewer({
    super.key,
    this.height = 300,
    this.primaryColor,
    this.label = 'ManPaSik Cartridge',
  });

  final double height;
  final Color? primaryColor;
  final String label;

  @override
  State<Cartridge3DViewer> createState() => _Cartridge3DViewerState();
}

class _Cartridge3DViewerState extends State<Cartridge3DViewer>
    with SingleTickerProviderStateMixin {
  double _rotationY = 0;
  double _rotationX = 0.3;
  double _scale = 1.0;
  late AnimationController _autoRotate;

  @override
  void initState() {
    super.initState();
    _autoRotate = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 20),
    )..repeat();
  }

  @override
  void dispose() {
    _autoRotate.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = widget.primaryColor ?? AppTheme.sanggamGold;

    return Column(
      children: [
        GestureDetector(
          onPanUpdate: (d) {
            setState(() {
              _rotationY += d.delta.dx * 0.01;
              _rotationX = (_rotationX + d.delta.dy * 0.01).clamp(-0.8, 0.8);
            });
          },
          onScaleUpdate: (d) {
            setState(() {
              _scale = (_scale * d.scale).clamp(0.5, 2.0);
            });
          },
          child: AnimatedBuilder(
            animation: _autoRotate,
            builder: (context, child) {
              final angle = _rotationY + _autoRotate.value * 2 * pi;
              return SizedBox(
                height: widget.height,
                width: double.infinity,
                child: CustomPaint(
                  painter: _CartridgePainter(
                    rotationY: angle,
                    rotationX: _rotationX,
                    scale: _scale,
                    primaryColor: color,
                    label: widget.label,
                    textStyle: theme.textTheme.bodySmall ?? const TextStyle(),
                  ),
                ),
              );
            },
          ),
        ),
        const SizedBox(height: 8),
        Text(
          '드래그하여 회전 | 핀치하여 확대',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant.withOpacity(0.6),
            fontSize: 11,
          ),
        ),
      ],
    );
  }
}

class _CartridgePainter extends CustomPainter {
  _CartridgePainter({
    required this.rotationY,
    required this.rotationX,
    required this.scale,
    required this.primaryColor,
    required this.label,
    required this.textStyle,
  });

  final double rotationY;
  final double rotationX;
  final double scale;
  final Color primaryColor;
  final String label;
  final TextStyle textStyle;

  @override
  void paint(Canvas canvas, Size size) {
    final cx = size.width / 2;
    final cy = size.height / 2;
    final baseW = 60.0 * scale;
    final baseH = 120.0 * scale;

    canvas.save();
    canvas.translate(cx, cy);

    // 3D 투영 시뮬레이션
    final cosY = cos(rotationY);
    final sinY = sin(rotationY);
    final cosX = cos(rotationX);

    // 카트리지 본체 (둥근 사각형)
    final depth = baseW * 0.4 * sinY.abs();
    final perspectiveW = baseW * cosY.abs().clamp(0.3, 1.0);

    // 그림자
    final shadowPaint = Paint()
      ..color = Colors.black.withOpacity(0.15)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 12);
    canvas.drawRRect(
      RRect.fromRectAndRadius(
        Rect.fromCenter(
          center: Offset(4, 8),
          width: perspectiveW * 2 + 8,
          height: baseH * cosX.abs().clamp(0.5, 1.0) + 8,
        ),
        const Radius.circular(12),
      ),
      shadowPaint,
    );

    // 본체
    final bodyPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [
          primaryColor.withOpacity(0.9),
          primaryColor.withOpacity(0.6),
          primaryColor.withOpacity(0.4),
        ],
      ).createShader(Rect.fromCenter(
        center: Offset.zero,
        width: perspectiveW * 2,
        height: baseH,
      ));

    final bodyRect = RRect.fromRectAndRadius(
      Rect.fromCenter(
        center: Offset.zero,
        width: perspectiveW * 2,
        height: baseH * cosX.abs().clamp(0.5, 1.0),
      ),
      const Radius.circular(10),
    );
    canvas.drawRRect(bodyRect, bodyPaint);

    // 테두리
    final borderPaint = Paint()
      ..color = primaryColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2;
    canvas.drawRRect(bodyRect, borderPaint);

    // 측면 깊이 (회전 시 보이는 면)
    if (sinY.abs() > 0.1) {
      final sidePaint = Paint()
        ..color = primaryColor.withOpacity(0.3);
      final sideOffset = sinY > 0 ? perspectiveW : -perspectiveW;
      canvas.drawRRect(
        RRect.fromRectAndRadius(
          Rect.fromLTWH(
            sideOffset - depth / 2,
            -baseH * cosX.abs().clamp(0.5, 1.0) / 2,
            depth,
            baseH * cosX.abs().clamp(0.5, 1.0),
          ),
          const Radius.circular(4),
        ),
        sidePaint,
      );
    }

    // NFC 칩 표시
    final chipPaint = Paint()
      ..color = Colors.white.withOpacity(0.6)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.5;
    canvas.drawCircle(
      Offset(0, -baseH * 0.2 * cosX.abs().clamp(0.5, 1.0)),
      12 * scale,
      chipPaint,
    );
    canvas.drawCircle(
      Offset(0, -baseH * 0.2 * cosX.abs().clamp(0.5, 1.0)),
      6 * scale,
      Paint()..color = Colors.white.withOpacity(0.4),
    );

    // 라벨
    final textPainter = TextPainter(
      text: TextSpan(
        text: label,
        style: textStyle.copyWith(
          color: Colors.white.withOpacity(cosY.abs().clamp(0.0, 0.8)),
          fontSize: 9 * scale,
        ),
      ),
      textDirection: TextDirection.ltr,
    );
    textPainter.layout();
    textPainter.paint(
      canvas,
      Offset(-textPainter.width / 2,
          baseH * 0.15 * cosX.abs().clamp(0.5, 1.0)),
    );

    canvas.restore();
  }

  @override
  bool shouldRepaint(covariant _CartridgePainter old) =>
      old.rotationY != rotationY ||
      old.rotationX != rotationX ||
      old.scale != scale;
}
