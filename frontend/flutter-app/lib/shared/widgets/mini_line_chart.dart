import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

/// 48x24 크기의 미니 스파크라인 차트.
/// 건강 요약 카드, 가디언 대시보드 등에서 추세를 표시할 때 사용.
class MiniLineChart extends StatelessWidget {
  const MiniLineChart({
    super.key,
    required this.data,
    this.width = 48,
    this.height = 24,
    this.lineColor,
    this.fillColor,
    this.strokeWidth = 1.5,
  });

  final List<double> data;
  final double width;
  final double height;
  final Color? lineColor;
  final Color? fillColor;
  final double strokeWidth;

  @override
  Widget build(BuildContext context) {
    if (data.isEmpty) {
      return SizedBox(width: width, height: height);
    }

    final color = lineColor ?? Theme.of(context).colorScheme.primary;

    final spots = <FlSpot>[];
    for (var i = 0; i < data.length; i++) {
      spots.add(FlSpot(i.toDouble(), data[i]));
    }

    return SizedBox(
      width: width,
      height: height,
      child: LineChart(
        LineChartData(
          gridData: const FlGridData(show: false),
          titlesData: const FlTitlesData(show: false),
          borderData: FlBorderData(show: false),
          lineTouchData: const LineTouchData(enabled: false),
          clipData: const FlClipData.all(),
          lineBarsData: [
            LineChartBarData(
              spots: spots,
              isCurved: true,
              curveSmoothness: 0.3,
              color: color,
              barWidth: strokeWidth,
              isStrokeCapRound: true,
              dotData: const FlDotData(show: false),
              belowBarData: BarAreaData(
                show: true,
                color: (fillColor ?? color).withOpacity(0.15),
              ),
            ),
          ],
        ),
        duration: Duration.zero,
      ),
    );
  }
}
