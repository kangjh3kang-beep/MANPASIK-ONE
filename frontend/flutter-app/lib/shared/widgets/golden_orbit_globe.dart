import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

/// 홈 화면 중앙에 배치되는 황금빛 궤도(Orbit) 지구본 위젯.
///
/// [CustomPaint] + [AnimationController]로 구현하며,
/// 중심 황금 구체 + 3개 타원 궤도 + 궤도 위 데이터 파티클로 구성된다.
///
/// ```dart
/// SizedBox(
///   width: 300, height: 300,
///   child: GoldenOrbitGlobe(size: 280),
/// )
/// ```
class GoldenOrbitGlobe extends StatefulWidget {
  /// 위젯이 그려질 논리적 크기 (정사각형).
  final double size;

  /// 구체·궤도·파티클의 기본 색상 (Sanggam Gold).
  final Color baseColor;

  /// 파티클 하이라이트 / 핵 밝은 부분 색상.
  final Color accentColor;

  const GoldenOrbitGlobe({
    super.key,
    this.size = 280,
    this.baseColor = SanggamTheme.primary,
    this.accentColor = Colors.white,
  });

  @override
  State<GoldenOrbitGlobe> createState() => _GoldenOrbitGlobeState();
}

class _GoldenOrbitGlobeState extends State<GoldenOrbitGlobe>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final List<_OrbitParticle> _particles;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 60),
    )..repeat();
    _particles = _OrbitParticle.generate();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, _) {
        return CustomPaint(
          size: Size(widget.size, widget.size),
          painter: _GoldenOrbitPainter(
            t: _controller.value,
            baseColor: widget.baseColor,
            accentColor: widget.accentColor,
            particles: _particles,
          ),
        );
      },
    );
  }
}

// ──────────────────────────────────────────────────────────────
//  데이터 모델
// ──────────────────────────────────────────────────────────────

/// 궤도 위를 도는 파티클 하나의 초기 상태.
class _OrbitParticle {
  /// 이 파티클이 속하는 궤도 인덱스 (0, 1, 2).
  final int orbitIndex;

  /// 궤도 위 위상 오프셋 (0 ~ 2π).
  final double phase;

  /// 상대 속도 배율 (0.3 ~ 1.2).
  final double speed;

  /// 파티클 반지름 (px).
  final double radius;

  /// 밝기/불투명도 (0.4 ~ 1.0).
  final double brightness;

  /// 글로우 여부 — brightness 가 높은 파티클만 해당.
  final bool hasGlow;

  const _OrbitParticle({
    required this.orbitIndex,
    required this.phase,
    required this.speed,
    required this.radius,
    required this.brightness,
    required this.hasGlow,
  });

  /// 3개 궤도에 걸쳐 파티클 약 40개를 생성한다.
  static List<_OrbitParticle> generate() {
    final rng = math.Random(42); // 고정 시드 → 재빌드 안정
    final list = <_OrbitParticle>[];

    for (int orbit = 0; orbit < 3; orbit++) {
      // 궤도당 12 ~ 15개
      final count = 12 + rng.nextInt(4);
      for (int i = 0; i < count; i++) {
        final brightness = 0.4 + rng.nextDouble() * 0.6;
        list.add(_OrbitParticle(
          orbitIndex: orbit,
          phase: rng.nextDouble() * 2 * math.pi,
          speed: 0.3 + rng.nextDouble() * 0.9,
          radius: 1.2 + rng.nextDouble() * 2.0,
          brightness: brightness,
          hasGlow: brightness > 0.85,
        ));
      }
    }
    return list;
  }
}

// ──────────────────────────────────────────────────────────────
//  CustomPainter
// ──────────────────────────────────────────────────────────────

class _GoldenOrbitPainter extends CustomPainter {
  /// 0 → 1 순환 (60초 주기).
  final double t;
  final Color baseColor;
  final Color accentColor;
  final List<_OrbitParticle> particles;

  _GoldenOrbitPainter({
    required this.t,
    required this.baseColor,
    required this.accentColor,
    required this.particles,
  });

  /// 3개 궤도 각각의 기울기(라디안).
  /// 0°, 60°, 120° 씩 기울여 교차하도록 한다.
  static const List<double> _orbitTilts = [
    0.0,
    math.pi / 3, // 60°
    2 * math.pi / 3, // 120°
  ];

  /// 궤도 별 회전 속도 배율 (서로 다른 속도로 회전).
  static const List<double> _orbitSpeeds = [1.0, 0.7, 0.45];

  /// 궤도 타원의 단반경 / 장반경 비율 (원근감).
  static const double _orbitEccentricity = 0.35;

  @override
  void paint(Canvas canvas, Size size) {
    final cx = size.width / 2;
    final cy = size.height / 2;
    final center = Offset(cx, cy);
    final orbitRadius = size.width * 0.42; // 궤도 장반경
    final coreRadius = size.width * 0.12; // 구체 반지름

    // ── 1. 외곽 에너지 오라 (가장 먼 레이어) ──
    _paintAura(canvas, center, orbitRadius);

    // ── 2. 궤도 후면 절반 (구체 뒤를 지나는 부분) ──
    for (int i = 0; i < 3; i++) {
      _paintOrbitArc(canvas, center, orbitRadius, i, back: true);
    }

    // ── 3. 궤도 후면 파티클 ──
    _paintParticles(canvas, center, orbitRadius, coreRadius, back: true);

    // ── 4. 중심 황금 구체 ──
    _paintCore(canvas, center, coreRadius);

    // ── 5. 궤도 전면 절반 (구체 앞을 지나는 부분) ──
    for (int i = 0; i < 3; i++) {
      _paintOrbitArc(canvas, center, orbitRadius, i, back: false);
    }

    // ── 6. 궤도 전면 파티클 ──
    _paintParticles(canvas, center, orbitRadius, coreRadius, back: false);

    // ── 7. 중심 글로스 하이라이트 (구체 맨 위) ──
    _paintCoreHighlight(canvas, center, coreRadius);
  }

  // ────────────────────────────────────────────
  //  1. 에너지 오라
  // ────────────────────────────────────────────

  void _paintAura(Canvas canvas, Offset center, double orbitR) {
    // 부드러운 외곽 발광
    final auraShader = RadialGradient(
      colors: [
        baseColor.withValues(alpha: 0.07),
        baseColor.withValues(alpha: 0.03),
        Colors.transparent,
      ],
      stops: const [0.5, 0.75, 1.0],
    ).createShader(Rect.fromCircle(center: center, radius: orbitR * 1.3));
    canvas.drawCircle(
      center,
      orbitR * 1.3,
      Paint()..shader = auraShader,
    );

    // 은은한 맥동(pulse): t 기반 사인파로 반지름 미세 변화
    final pulse = 1.0 + math.sin(t * 2 * math.pi * 3) * 0.015;
    final pulseR = orbitR * 1.15 * pulse;
    canvas.drawCircle(
      center,
      pulseR,
      Paint()
        ..color = baseColor.withValues(alpha: 0.04)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.5
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 6),
    );
  }

  // ────────────────────────────────────────────
  //  2/5. 궤도 호(Arc) 그리기
  // ────────────────────────────────────────────

  /// [back]=true → 구체 뒤쪽 반원, false → 앞쪽 반원.
  void _paintOrbitArc(
    Canvas canvas,
    Offset center,
    double orbitR,
    int orbitIndex, {
    required bool back,
  }) {
    final tilt = _orbitTilts[orbitIndex];
    final speed = _orbitSpeeds[orbitIndex];
    final baseAngle = t * 2 * math.pi * speed;

    canvas.save();
    canvas.translate(center.dx, center.dy);
    // 궤도 면 자체도 미세하게 회전 (생동감)
    canvas.rotate(tilt + baseAngle * 0.05);

    final rx = orbitR; // 장반경
    final ry = orbitR * _orbitEccentricity; // 단반경

    final orbitRect = Rect.fromCenter(
      center: Offset.zero,
      width: rx * 2,
      height: ry * 2,
    );

    // 뒤/앞 절반만 그리기
    final startAngle = back ? 0.0 : math.pi;
    const sweepAngle = math.pi;

    // 궤도 선 (golden line)
    final alpha = back ? 0.15 : 0.35;
    final linePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.8
      ..color = baseColor.withValues(alpha: alpha);
    canvas.drawArc(orbitRect, startAngle, sweepAngle, false, linePaint);

    // 궤도 글로우 (앞면만)
    if (!back) {
      final glowPaint = Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 2.5
        ..color = baseColor.withValues(alpha: 0.06)
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 4);
      canvas.drawArc(orbitRect, startAngle, sweepAngle, false, glowPaint);
    }

    canvas.restore();
  }

  // ────────────────────────────────────────────
  //  3/6. 궤도 위 파티클
  // ────────────────────────────────────────────

  void _paintParticles(
    Canvas canvas,
    Offset center,
    double orbitR,
    double coreR, {
    required bool back,
  }) {
    final paint = Paint()..style = PaintingStyle.fill;
    final glowPaint = Paint()..style = PaintingStyle.fill;

    for (final p in particles) {
      final tilt = _orbitTilts[p.orbitIndex];
      final speed = _orbitSpeeds[p.orbitIndex];
      final baseAngle = t * 2 * math.pi * speed;
      final faceAngle = tilt + baseAngle * 0.05;

      // 파티클의 궤도 위 위치 (타원 매개변수)
      final theta = p.phase + baseAngle * p.speed;
      final rx = orbitR;
      final ry = orbitR * _orbitEccentricity;

      // 타원 위의 점 (궤도 로컬 좌표)
      final localX = rx * math.cos(theta);
      final localY = ry * math.sin(theta);

      // 궤도 면 회전 적용 → 글로벌 좌표
      final cosF = math.cos(faceAngle);
      final sinF = math.sin(faceAngle);
      final gx = localX * cosF - localY * sinF + center.dx;
      final gy = localX * sinF + localY * cosF + center.dy;

      // 깊이 판별: localY > 0 은 화면 "앞" (관찰자 방향)
      final isFront = localY >= 0;
      if (back && isFront) continue;
      if (!back && !isFront) continue;

      // 구체에 가려지는 파티클 감추기 (뒤쪽 + 구체 영역 내부)
      final dx = gx - center.dx;
      final dy = gy - center.dy;
      final distFromCenter = math.sqrt(dx * dx + dy * dy);
      if (back && distFromCenter < coreR * 0.85) continue;

      // 깊이 기반 알파
      final depthAlpha = isFront ? 1.0 : 0.45;
      final alpha = (p.brightness * depthAlpha).clamp(0.0, 1.0);

      // 파티클 색상: gold → white 보간
      final color = Color.lerp(baseColor, accentColor, p.brightness * 0.4)!;

      paint.color = color.withValues(alpha: alpha);
      final pr = p.radius * (isFront ? 1.0 : 0.7);
      canvas.drawCircle(Offset(gx, gy), pr, paint);

      // 밝은 파티클에 글로우 효과
      if (p.hasGlow && isFront) {
        glowPaint
          ..color = accentColor.withValues(alpha: alpha * 0.5)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 3);
        canvas.drawCircle(Offset(gx, gy), pr * 1.8, glowPaint);
      }
    }
  }

  // ────────────────────────────────────────────
  //  4. 중심 황금 구체
  // ────────────────────────────────────────────

  void _paintCore(Canvas canvas, Offset center, double r) {
    // 레이어 1: 코어 배경 — 다층 RadialGradient로 깊이감 있는 금속 구체
    final bgShader = RadialGradient(
      center: const Alignment(-0.25, -0.3), // 좌상단 광원
      radius: 1.0,
      colors: const [
        Color(0xFFE8C84A), // 밝은 금
        SanggamTheme.primary, // 상감 금
        Color(0xFF8A6B12), // 어두운 금
        Color(0xFF5A4510), // 가장자리 그림자
      ],
      stops: const [0.0, 0.35, 0.7, 1.0],
    ).createShader(Rect.fromCircle(center: center, radius: r));
    canvas.drawCircle(center, r, Paint()..shader = bgShader);

    // 레이어 2: 내부 광택
    final sheenShader = RadialGradient(
      center: const Alignment(-0.35, -0.4),
      radius: 0.6,
      colors: [
        accentColor.withValues(alpha: 0.35),
        accentColor.withValues(alpha: 0.05),
        Colors.transparent,
      ],
      stops: const [0.0, 0.5, 1.0],
    ).createShader(Rect.fromCircle(center: center, radius: r));
    canvas.drawCircle(center, r, Paint()..shader = sheenShader);

    // 레이어 3: 가장자리 림 라이트
    canvas.drawCircle(
      center,
      r,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.0
        ..color = baseColor.withValues(alpha: 0.5),
    );

    // 레이어 4: 외곽 글로우
    canvas.drawCircle(
      center,
      r + 2,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 4
        ..color = baseColor.withValues(alpha: 0.08)
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 6),
    );
  }

  // ────────────────────────────────────────────
  //  7. 구체 글로스 하이라이트
  // ────────────────────────────────────────────

  void _paintCoreHighlight(Canvas canvas, Offset center, double r) {
    // 좌상단 스페큘러 하이라이트 (작은 밝은 점)
    final hlCenter = Offset(center.dx - r * 0.3, center.dy - r * 0.3);
    final hlR = r * 0.25;
    final hlShader = RadialGradient(
      colors: [
        accentColor.withValues(alpha: 0.6),
        accentColor.withValues(alpha: 0.1),
        Colors.transparent,
      ],
      stops: const [0.0, 0.4, 1.0],
    ).createShader(Rect.fromCircle(center: hlCenter, radius: hlR));
    canvas.drawCircle(hlCenter, hlR, Paint()..shader = hlShader);

    // 하단 반사 (바닥 반사광)
    final btCenter = Offset(center.dx + r * 0.15, center.dy + r * 0.35);
    final btR = r * 0.15;
    final btShader = RadialGradient(
      colors: [
        baseColor.withValues(alpha: 0.2),
        Colors.transparent,
      ],
    ).createShader(Rect.fromCircle(center: btCenter, radius: btR));
    canvas.drawCircle(btCenter, btR, Paint()..shader = btShader);
  }

  @override
  bool shouldRepaint(covariant _GoldenOrbitPainter oldDelegate) {
    return oldDelegate.t != t;
  }
}
