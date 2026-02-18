import 'package:flutter/material.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 생체 핑거프린트 히트맵 (C2)
///
/// 896차원 스펙트럼 데이터를 32x28 그리드로 변환하여
/// CustomPainter로 히트맵을 그립니다.
class FingerprintHeatmap extends StatelessWidget {
  const FingerprintHeatmap({
    super.key,
    required this.grid,
    this.height = 200,
  });

  final List<List<double>> grid;
  final double height;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (grid.isEmpty) {
      return SizedBox(
        height: height,
        child: const Center(child: Text('히트맵 데이터 없음')),
      );
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          '스펙트럼 히트맵',
          style: theme.textTheme.titleSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 8),
        ClipRRect(
          borderRadius: BorderRadius.circular(12),
          child: SizedBox(
            height: height,
            width: double.infinity,
            child: CustomPaint(
              painter: _HeatmapPainter(grid: grid),
            ),
          ),
        ),
        const SizedBox(height: 8),
        // 색상 범례
        SizedBox(
          height: 16,
          child: Row(
            children: [
              Text('낮음', style: theme.textTheme.bodySmall),
              const SizedBox(width: 4),
              Expanded(
                child: Container(
                  height: 10,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(5),
                    gradient: const LinearGradient(
                      colors: [
                        Color(0xFF1A237E),
                        Color(0xFF0D47A1),
                        Color(0xFF00BCD4),
                        Color(0xFF4CAF50),
                        Color(0xFFFFEB3B),
                        Color(0xFFFF9800),
                        Color(0xFFF44336),
                      ],
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 4),
              Text('높음', style: theme.textTheme.bodySmall),
            ],
          ),
        ),
      ],
    );
  }
}

class _HeatmapPainter extends CustomPainter {
  _HeatmapPainter({required this.grid});
  final List<List<double>> grid;

  @override
  void paint(Canvas canvas, Size size) {
    if (grid.isEmpty || grid[0].isEmpty) return;

    final rows = grid.length;
    final cols = grid[0].length;
    final cellW = size.width / cols;
    final cellH = size.height / rows;
    final paint = Paint();

    for (var r = 0; r < rows; r++) {
      for (var c = 0; c < cols; c++) {
        final value = grid[r][c].clamp(0.0, 1.0);
        paint.color = _valueToColor(value);
        canvas.drawRect(
          Rect.fromLTWH(c * cellW, r * cellH, cellW + 0.5, cellH + 0.5),
          paint,
        );
      }
    }
  }

  Color _valueToColor(double v) {
    // 7단계 그라데이션: 남색→파랑→시안→녹색→노랑→주황→빨강
    const stops = [0.0, 0.17, 0.33, 0.5, 0.67, 0.83, 1.0];
    const colors = [
      Color(0xFF1A237E),
      Color(0xFF0D47A1),
      Color(0xFF00BCD4),
      Color(0xFF4CAF50),
      Color(0xFFFFEB3B),
      Color(0xFFFF9800),
      Color(0xFFF44336),
    ];

    for (var i = 0; i < stops.length - 1; i++) {
      if (v <= stops[i + 1]) {
        final t = (v - stops[i]) / (stops[i + 1] - stops[i]);
        return Color.lerp(colors[i], colors[i + 1], t)!;
      }
    }
    return colors.last;
  }

  @override
  bool shouldRepaint(covariant _HeatmapPainter old) => old.grid != grid;
}
