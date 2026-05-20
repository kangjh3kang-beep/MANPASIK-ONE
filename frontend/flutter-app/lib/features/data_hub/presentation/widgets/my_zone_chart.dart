import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// MyZoneChart — Sanggam Orbit My Zone 차트
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] Theme.of(context) isDark 분기 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold 5x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan 4x → SanggamTheme.jagaeCyan
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] Card → SanggamTheme 스타일 Container
// [Rule 4] theme.textTheme → 직접 TextStyle
// [Rule 4] ChoiceChip → 다크 테마 스타일
// [Rule 4] withOpacity → withValues(alpha:)
// ───────────────────────────────────────────────────

/// My Zone 차트 위젯.
///
/// 바이오마커 수치의 시계열 변화를 라인 차트로 표시하고,
/// 사용자의 정상 범위(My Zone)를 반투명 배경으로 오버레이합니다.
class MyZoneChart extends StatefulWidget {
  const MyZoneChart({
    super.key,
    this.title = '바이오마커 트렌드',
    this.dataPoints,
    this.normalMin = 60.0,
    this.normalMax = 100.0,
    this.unit = '',
  });

  final String title;
  final List<double>? dataPoints;
  final double normalMin;
  final double normalMax;
  final String unit;

  @override
  State<MyZoneChart> createState() =>
      _MyZoneChartState();
}

class _MyZoneChartState
    extends State<MyZoneChart> {
  String _period = '7일';
  static const _periods = [
    '7일',
    '30일',
    '90일'
  ];

  List<double> get _data {
    if (widget.dataPoints != null) {
      return widget.dataPoints!;
    }
    final rand = math.Random(42);
    final count = _period == '7일'
        ? 7
        : _period == '30일'
            ? 30
            : 90;
    return List.generate(count,
        (i) => 65.0 + rand.nextDouble() * 45);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white
            .withValues(alpha: 0.05),
        borderRadius:
            BorderRadius.circular(16),
        border: Border.all(
          color: SanggamTheme.surfaceVariant,
        ),
      ),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment:
                MainAxisAlignment
                    .spaceBetween,
            children: [
              Text(
                widget.title,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 14,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Row(
                children:
                    _periods.map((p) {
                  final isSelected =
                      _period == p;
                  return Padding(
                    padding:
                        const EdgeInsets
                            .only(left: 4),
                    child: ChoiceChip(
                      label: Text(p,
                          style: TextStyle(
                            fontSize: 11,
                            color: isSelected
                                ? SanggamTheme
                                    .background
                                : SanggamTheme
                                    .onSurfaceDim,
                          )),
                      selected: isSelected,
                      selectedColor:
                          SanggamTheme
                              .primary,
                      backgroundColor:
                          SanggamTheme
                              .surface,
                      side: BorderSide(
                        color: isSelected
                            ? SanggamTheme
                                .primary
                            : SanggamTheme
                                .surfaceVariant,
                      ),
                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius
                                .circular(
                                    16),
                      ),
                      onSelected: (_) =>
                          setState(() =>
                              _period = p),
                      visualDensity:
                          VisualDensity
                              .compact,
                      materialTapTargetSize:
                          MaterialTapTargetSize
                              .shrinkWrap,
                    ),
                  );
                }).toList(),
              ),
            ],
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 160,
            child: CustomPaint(
              size: Size.infinite,
              painter: _MyZoneChartPainter(
                data: _data,
                normalMin:
                    widget.normalMin,
                normalMax:
                    widget.normalMax,
              ),
            ),
          ),
          const SizedBox(height: 8),
          Row(
            children: [
              Container(
                width: 12,
                height: 12,
                decoration: BoxDecoration(
                  color: SanggamTheme
                      .primary
                      .withValues(
                          alpha: 0.2),
                  border: Border.all(
                      color: SanggamTheme
                          .primary,
                      width: 1),
                  borderRadius:
                      BorderRadius.circular(
                          2),
                ),
              ),
              const SizedBox(width: 4),
              Text(
                'My Zone (${widget.normalMin.toInt()}-${widget.normalMax.toInt()}${widget.unit})',
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                ),
              ),
              const SizedBox(width: 16),
              Container(
                width: 12,
                height: 2,
                color:
                    SanggamTheme.jagaeCyan,
              ),
              const SizedBox(width: 4),
              const Text(
                '측정값',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _MyZoneChartPainter
    extends CustomPainter {
  final List<double> data;
  final double normalMin;
  final double normalMax;

  _MyZoneChartPainter({
    required this.data,
    required this.normalMin,
    required this.normalMax,
  });

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;

    final allValues = [
      ...data,
      normalMin,
      normalMax
    ];
    final minVal =
        allValues.reduce(math.min) - 5;
    final maxVal =
        allValues.reduce(math.max) + 5;
    final range = maxVal - minVal;
    if (range <= 0) return;

    double toY(double val) =>
        size.height -
        ((val - minVal) / range *
            size.height);
    double toX(int i) => data.length == 1
        ? size.width / 2
        : i / (data.length - 1) * size.width;

    // My Zone 배경
    final zoneRect = Rect.fromLTRB(
        0,
        toY(normalMax),
        size.width,
        toY(normalMin));
    canvas.drawRect(
      zoneRect,
      Paint()
        ..color = SanggamTheme.primary
            .withValues(alpha: 0.12)
        ..style = PaintingStyle.fill,
    );

    // Zone 경계선
    final zoneBorderPaint = Paint()
      ..color = SanggamTheme.primary
          .withValues(alpha: 0.3)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5;
    canvas.drawLine(
        Offset(0, toY(normalMax)),
        Offset(size.width, toY(normalMax)),
        zoneBorderPaint);
    canvas.drawLine(
        Offset(0, toY(normalMin)),
        Offset(size.width, toY(normalMin)),
        zoneBorderPaint);

    // 데이터 라인
    final path = Path();
    for (var i = 0; i < data.length; i++) {
      final p =
          Offset(toX(i), toY(data[i]));
      if (i == 0) {
        path.moveTo(p.dx, p.dy);
      } else {
        path.lineTo(p.dx, p.dy);
      }
    }

    canvas.drawPath(
      path,
      Paint()
        ..color = SanggamTheme.jagaeCyan
        ..style = PaintingStyle.stroke
        ..strokeWidth = 2
        ..strokeCap = StrokeCap.round
        ..strokeJoin = StrokeJoin.round,
    );

    // 데이터 포인트
    if (data.length <= 30) {
      final dotPaint = Paint()
        ..color = SanggamTheme.jagaeCyan;
      for (var i = 0;
          i < data.length;
          i++) {
        final p =
            Offset(toX(i), toY(data[i]));
        final inZone =
            data[i] >= normalMin &&
                data[i] <= normalMax;
        dotPaint.color = inZone
            ? SanggamTheme.jagaeCyan
            : SanggamTheme.error;
        canvas.drawCircle(p, 3, dotPaint);
      }
    }
  }

  @override
  bool shouldRepaint(
          _MyZoneChartPainter old) =>
      old.data != data ||
      old.normalMin != normalMin ||
      old.normalMax != normalMax;
}
