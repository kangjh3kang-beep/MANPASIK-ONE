import 'dart:math' as math;
import 'dart:ui';
import 'package:flutter/material.dart';

/// HoloBody V22: Path-based Anatomical Hologram Engine
///
/// Bézier curve paths for anatomically correct human body silhouette.
/// 8-pass rendering: aura, volume, grid, rim, scan, anatomy, organs, particles.

enum HoloGender { male, female }

class HoloBody extends StatefulWidget {
  final double width;
  final double height;
  final Color color;
  final Color? accentColor;
  final HoloGender gender;
  final Map<String, dynamic> bioData;
  final bool showEcg;
  final bool showHud;

  const HoloBody({
    super.key,
    this.width = 300,
    this.height = 500,
    this.color = const Color(0xFF00E5FF),
    this.accentColor,
    this.gender = HoloGender.male,
    this.bioData = const {},
    this.showEcg = false,
    this.showHud = false,
  });

  @override
  State<HoloBody> createState() => _HoloBodyState();
}

class _HoloBodyState extends State<HoloBody> with TickerProviderStateMixin {
  late AnimationController _rotateController;
  late AnimationController _scanController;
  late AnimationController _pulseController;
  late AnimationController _breathController;
  late AnimationController _particleController;

  @override
  void initState() {
    super.initState();
    _rotateController = AnimationController(
      vsync: this, duration: const Duration(seconds: 40),
    )..repeat();
    _scanController = AnimationController(
      vsync: this, duration: const Duration(seconds: 5),
    )..repeat();
    _pulseController = AnimationController(
      vsync: this, duration: const Duration(milliseconds: 900),
    )..repeat(reverse: true);
    _breathController = AnimationController(
      vsync: this, duration: const Duration(milliseconds: 3500),
    )..repeat(reverse: true);
    _particleController = AnimationController(
      vsync: this, duration: const Duration(seconds: 8),
    )..repeat();
  }

  @override
  void dispose() {
    _rotateController.dispose();
    _scanController.dispose();
    _pulseController.dispose();
    _breathController.dispose();
    _particleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: widget.width,
      height: widget.height,
      child: ExcludeSemantics(
        child: RepaintBoundary(
          child: AnimatedBuilder(
            animation: Listenable.merge([
              _rotateController, _scanController, _pulseController,
              _breathController, _particleController,
            ]),
            builder: (context, child) {
              return CustomPaint(
                size: Size(widget.width, widget.height),
                painter: _HoloV22Painter(
                  gender: widget.gender,
                  rotation: _rotateController.value * 2 * math.pi,
                  scanY: _scanController.value,
                  pulse: _pulseController.value,
                  breath: _breathController.value,
                  particlePhase: _particleController.value,
                  color: widget.color,
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════════════════
// Body width profiles per gender (half-widths from center)
// ═══════════════════════════════════════════════════════════════════════════════

class _BW {
  final double head, jaw, neck, shoulder, chest, waist, hip;
  final double thigh, knee, calf, ankle, foot;

  const _BW({
    required this.head, required this.jaw, required this.neck,
    required this.shoulder, required this.chest, required this.waist,
    required this.hip, required this.thigh, required this.knee,
    required this.calf, required this.ankle, required this.foot,
  });

  static const male = _BW(
    head: 0.095, jaw: 0.088, neck: 0.052,
    shoulder: 0.200, chest: 0.170, waist: 0.130,
    hip: 0.155, thigh: 0.090, knee: 0.055,
    calf: 0.062, ankle: 0.038, foot: 0.052,
  );

  static const female = _BW(
    head: 0.090, jaw: 0.083, neck: 0.046,
    shoulder: 0.175, chest: 0.165, waist: 0.110,
    hip: 0.185, thigh: 0.100, knee: 0.052,
    calf: 0.057, ankle: 0.035, foot: 0.048,
  );
}

// ═══════════════════════════════════════════════════════════════════════════════
// V22 Multi-Pass Hologram Painter
// ═══════════════════════════════════════════════════════════════════════════════

class _HoloV22Painter extends CustomPainter {
  final HoloGender gender;
  final double rotation;
  final double scanY;
  final double pulse;
  final double breath;
  final double particlePhase;
  final Color color;

  _HoloV22Painter({
    required this.gender,
    required this.rotation,
    required this.scanY,
    required this.pulse,
    required this.breath,
    required this.particlePhase,
    required this.color,
  });

  static const _cyan = Color(0xFF00E5FF);
  static const _cyanDark = Color(0xFF006B7A);
  static const _cyanBright = Color(0xFF80F0FF);
  static const _blue = Color(0xFF2962FF);
  static const _red = Color(0xFFFF1744);
  static const _purple = Color(0xFFAA00FF);
  static const _white = Color(0xFFFFFFFF);

  @override
  void paint(Canvas canvas, Size size) {
    final cx = size.width / 2;
    const meshH = 2.037;
    const padding = 0.05;
    final fs = size.height * (1.0 - 2 * padding) / meshH;
    const meshCenterY = (-1.0 + 1.037) / 2;
    final cy = size.height / 2 - meshCenterY * fs;

    final cosR = math.cos(rotation);
    final sinR = math.sin(rotation);
    final ws = 0.7 + 0.3 * cosR.abs();
    final w = gender == HoloGender.male ? _BW.male : _BW.female;

    final body = _buildBody(w, fs, cx, cy, ws);
    final arms = [_buildArm(w, fs, cx, cy, ws, 1), _buildArm(w, fs, cx, cy, ws, -1)];
    final combined = Path()
      ..addPath(body, Offset.zero)
      ..addPath(arms[0], Offset.zero)
      ..addPath(arms[1], Offset.zero);

    final scanMY = (scanY * 2.4) - 1.2;

    _drawAura(canvas, cx, cy, size);
    _drawVolume(canvas, body, arms, fs, cx, cy, scanMY);
    _drawGrid(canvas, combined, w, fs, cx, cy, ws, scanMY);
    _drawRim(canvas, body, arms);
    _drawScanBand(canvas, cx, cy, fs, size, scanMY);
    _drawAnatomy(canvas, cx, cy, fs, cosR, sinR);
    _drawOrgans(canvas, cx, cy, fs, cosR, sinR);
    _drawParticles(canvas, body);
  }

  // ══════════════════════════════════════════════════════════════════════
  // Body Outline Path (torso + head + legs, no arms)
  // ══════════════════════════════════════════════════════════════════════

  Path _buildBody(_BW w, double fs, double cx, double cy, double ws) {
    double sx(double x) => cx + x * ws * fs;
    double sy(double y) => cy + y * fs;
    const g = 0.022; // half-gap between legs

    final p = Path()..moveTo(sx(0), sy(-1.0));

    // ── Right head: top → temple → jaw ──
    p.cubicTo(sx(w.head * 0.55), sy(-1.0), sx(w.head), sy(-0.95), sx(w.head), sy(-0.87));
    p.cubicTo(sx(w.head), sy(-0.82), sx(w.jaw), sy(-0.79), sx(w.jaw * 0.85), sy(-0.76));

    // ── Jaw → neck → shoulder ──
    p.cubicTo(sx(w.jaw * 0.7), sy(-0.74), sx(w.neck * 1.1), sy(-0.73), sx(w.neck), sy(-0.70));
    p.cubicTo(sx(w.neck), sy(-0.67), sx(w.shoulder * 0.5), sy(-0.64), sx(w.shoulder), sy(-0.62));

    // ── Shoulder → chest → waist → hip ──
    p.cubicTo(sx(w.shoulder * 0.98), sy(-0.58), sx(w.chest * 1.05), sy(-0.54), sx(w.chest), sy(-0.48));
    p.cubicTo(sx(w.chest * 0.97), sy(-0.40), sx(w.waist * 1.15), sy(-0.32), sx(w.waist), sy(-0.24));
    p.cubicTo(sx(w.waist * 0.97), sy(-0.17), sx(w.hip * 0.85), sy(-0.08), sx(w.hip), sy(-0.02));

    // ── Hip → thigh → knee → calf → ankle → foot ──
    p.cubicTo(sx(w.hip * 0.97), sy(0.02), sx(w.thigh * 1.25), sy(0.06), sx(w.thigh), sy(0.12));
    p.cubicTo(sx(w.thigh * 0.97), sy(0.25), sx(w.knee * 1.15), sy(0.42), sx(w.knee), sy(0.50));
    p.cubicTo(sx(w.knee * 0.98), sy(0.53), sx(w.calf * 1.05), sy(0.58), sx(w.calf), sy(0.68));
    p.cubicTo(sx(w.calf * 0.85), sy(0.78), sx(w.ankle * 1.3), sy(0.87), sx(w.ankle), sy(0.92));
    p.cubicTo(sx(w.ankle * 0.95), sy(0.96), sx(w.foot), sy(1.01), sx(w.foot), sy(1.037));

    // ── Right foot bottom ──
    p.lineTo(sx(g), sy(1.037));

    // ── Right inner leg (up) ──
    p.cubicTo(sx(g), sy(1.01), sx(w.ankle * 0.55), sy(0.96), sx(w.ankle * 0.5), sy(0.92));
    p.cubicTo(sx(w.ankle * 0.55), sy(0.87), sx(w.calf * 0.42), sy(0.78), sx(w.calf * 0.40), sy(0.68));
    p.cubicTo(sx(w.calf * 0.38), sy(0.58), sx(w.knee * 0.55), sy(0.53), sx(w.knee * 0.50), sy(0.50));
    p.cubicTo(sx(w.knee * 0.55), sy(0.42), sx(w.thigh * 0.55), sy(0.25), sx(w.thigh * 0.45), sy(0.12));
    p.cubicTo(sx(w.thigh * 0.35), sy(0.06), sx(g * 1.5), sy(0.04), sx(g), sy(0.03));

    // ── Crotch V ──
    p.cubicTo(sx(g * 0.5), sy(0.05), sx(0), sy(0.065), sx(-g), sy(0.03));

    // ── Left inner leg (down) ──
    p.cubicTo(sx(-g * 1.5), sy(0.04), sx(-w.thigh * 0.35), sy(0.06), sx(-w.thigh * 0.45), sy(0.12));
    p.cubicTo(sx(-w.thigh * 0.55), sy(0.25), sx(-w.knee * 0.55), sy(0.42), sx(-w.knee * 0.50), sy(0.50));
    p.cubicTo(sx(-w.calf * 0.38), sy(0.58), sx(-w.calf * 0.42), sy(0.78), sx(-w.calf * 0.40), sy(0.68));
    p.cubicTo(sx(-w.ankle * 0.55), sy(0.87), sx(-g), sy(1.01), sx(-g), sy(1.037));

    // ── Left foot bottom ──
    p.lineTo(sx(-w.foot), sy(1.037));

    // ── Left outer leg (up) ──
    p.cubicTo(sx(-w.foot), sy(1.01), sx(-w.ankle * 0.95), sy(0.96), sx(-w.ankle), sy(0.92));
    p.cubicTo(sx(-w.ankle * 1.3), sy(0.87), sx(-w.calf * 0.85), sy(0.78), sx(-w.calf), sy(0.68));
    p.cubicTo(sx(-w.calf * 1.05), sy(0.58), sx(-w.knee * 0.98), sy(0.53), sx(-w.knee), sy(0.50));
    p.cubicTo(sx(-w.knee * 1.15), sy(0.42), sx(-w.thigh * 0.97), sy(0.25), sx(-w.thigh), sy(0.12));
    p.cubicTo(sx(-w.thigh * 1.25), sy(0.06), sx(-w.hip * 0.97), sy(0.02), sx(-w.hip), sy(-0.02));

    // ── Left hip → waist → chest → shoulder → neck → jaw → head ──
    p.cubicTo(sx(-w.hip * 0.85), sy(-0.08), sx(-w.waist * 0.97), sy(-0.17), sx(-w.waist), sy(-0.24));
    p.cubicTo(sx(-w.waist * 1.15), sy(-0.32), sx(-w.chest * 0.97), sy(-0.40), sx(-w.chest), sy(-0.48));
    p.cubicTo(sx(-w.chest * 1.05), sy(-0.54), sx(-w.shoulder * 0.98), sy(-0.58), sx(-w.shoulder), sy(-0.62));
    p.cubicTo(sx(-w.shoulder * 0.5), sy(-0.64), sx(-w.neck), sy(-0.67), sx(-w.neck), sy(-0.70));
    p.cubicTo(sx(-w.neck * 1.1), sy(-0.73), sx(-w.jaw * 0.7), sy(-0.74), sx(-w.jaw * 0.85), sy(-0.76));
    p.cubicTo(sx(-w.jaw), sy(-0.79), sx(-w.head), sy(-0.82), sx(-w.head), sy(-0.87));
    p.cubicTo(sx(-w.head), sy(-0.95), sx(-w.head * 0.55), sy(-1.0), sx(0), sy(-1.0));

    p.close();
    return p;
  }

  // ══════════════════════════════════════════════════════════════════════
  // Arm Path (simple tapered shape)
  // ══════════════════════════════════════════════════════════════════════

  Path _buildArm(_BW w, double fs, double cx, double cy, double ws, double side) {
    double sx(double x) => cx + x * side * ws * fs;
    double sy(double y) => cy + y * fs;
    final s = w.shoulder;

    final p = Path()..moveTo(sx(s), sy(-0.62));

    // Outer edge (down)
    p.cubicTo(sx(s * 0.93), sy(-0.52), sx(s * 0.84), sy(-0.38), sx(s * 0.78), sy(-0.24));
    p.cubicTo(sx(s * 0.74), sy(-0.12), sx(s * 0.70), sy(0.00), sx(s * 0.67), sy(0.10));
    p.cubicTo(sx(s * 0.65), sy(0.14), sx(s * 0.64), sy(0.17), sx(s * 0.62), sy(0.19));

    // Fingertips
    p.cubicTo(sx(s * 0.60), sy(0.21), sx(s * 0.55), sy(0.21), sx(s * 0.54), sy(0.19));

    // Inner edge (up)
    p.cubicTo(sx(s * 0.56), sy(0.14), sx(s * 0.58), sy(0.00), sx(s * 0.62), sy(-0.24));
    p.cubicTo(sx(s * 0.66), sy(-0.38), sx(s * 0.72), sy(-0.52), sx(s * 0.80), sy(-0.60));
    p.lineTo(sx(s * 0.88), sy(-0.62));

    p.close();
    return p;
  }

  // ══════════════════════════════════════════════════════════════════════
  // Body width at Y (for grid lines)
  // ══════════════════════════════════════════════════════════════════════

  static double _bodyW(double y, _BW w, double ws) {
    const ys =  [-1.0, -0.87, -0.76, -0.70, -0.62, -0.48, -0.24, -0.02, 0.12, 0.50, 0.68, 0.92, 1.037];
    final ws2 = [0.0, w.head, w.jaw * 0.85, w.neck, w.shoulder, w.chest, w.waist, w.hip, w.thigh, w.knee, w.calf, w.ankle, w.foot];
    for (int i = 0; i < ys.length - 1; i++) {
      if (y <= ys[i + 1]) {
        final t = (y - ys[i]) / (ys[i + 1] - ys[i]);
        return (ws2[i] + t * (ws2[i + 1] - ws2[i])) * ws;
      }
    }
    return ws2.last * ws;
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 1: Ambient Aura
  // ══════════════════════════════════════════════════════════════════════

  void _drawAura(Canvas canvas, double cx, double cy, Size size) {
    final radius = size.height * 0.45;
    final paint = Paint()
      ..shader = RadialGradient(
        colors: [
          _cyan.withValues(alpha: 0.08),
          _cyanDark.withValues(alpha: 0.04),
          Colors.transparent,
        ],
        stops: const [0.0, 0.5, 1.0],
      ).createShader(Rect.fromCircle(center: Offset(cx, cy), radius: radius));
    canvas.drawCircle(Offset(cx, cy), radius, paint);
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 2: Body Volume (filled path + Fresnel gradient)
  // ══════════════════════════════════════════════════════════════════════

  void _drawVolume(Canvas canvas, Path body, List<Path> arms,
      double fs, double cx, double cy, double scanMY) {
    // Base fill
    final basePaint = Paint()
      ..color = _cyan.withValues(alpha: 0.07)
      ..blendMode = BlendMode.plus;
    canvas.drawPath(body, basePaint);
    for (final arm in arms) {
      canvas.drawPath(arm, Paint()
        ..color = _cyan.withValues(alpha: 0.05)
        ..blendMode = BlendMode.plus);
    }

    // Horizontal gradient for Fresnel-like edge brightening
    final bounds = body.getBounds();
    final fresnelPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.centerLeft,
        end: Alignment.centerRight,
        colors: [
          _cyan.withValues(alpha: 0.18),
          _cyan.withValues(alpha: 0.02),
          _cyan.withValues(alpha: 0.02),
          _cyan.withValues(alpha: 0.18),
        ],
        stops: const [0.0, 0.3, 0.7, 1.0],
      ).createShader(bounds)
      ..blendMode = BlendMode.plus;
    canvas.drawPath(body, fresnelPaint);

    // Scan highlight band (clipped to body)
    final scanScreenY = cy + scanMY * fs;
    final scanRect = Rect.fromCenter(
      center: Offset(cx, scanScreenY),
      width: bounds.width * 1.2,
      height: fs * 0.16,
    );
    canvas.save();
    canvas.clipPath(body);
    canvas.drawRect(scanRect, Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          Colors.transparent,
          _cyanBright.withValues(alpha: 0.35),
          _white.withValues(alpha: 0.5),
          _cyanBright.withValues(alpha: 0.35),
          Colors.transparent,
        ],
      ).createShader(scanRect)
      ..blendMode = BlendMode.plus);
    canvas.restore();
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 3: Holographic Grid (internal lines clipped to body)
  // ══════════════════════════════════════════════════════════════════════

  void _drawGrid(Canvas canvas, Path combined, _BW w,
      double fs, double cx, double cy, double ws, double scanMY) {
    canvas.save();
    canvas.clipPath(combined);

    final gridPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.4
      ..color = _cyan.withValues(alpha: 0.12)
      ..blendMode = BlendMode.plus;

    final scanGridPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.8
      ..color = _cyanBright.withValues(alpha: 0.5)
      ..blendMode = BlendMode.plus;

    // Horizontal scan lines
    const step = 0.04;
    for (double y = -1.0; y <= 1.04; y += step) {
      final bw = _bodyW(y, w, ws);
      final screenY = cy + y * fs;
      final halfW = bw * fs;

      final isScanNear = (y - scanMY).abs() < 0.06;
      canvas.drawLine(
        Offset(cx - halfW * 1.2, screenY),
        Offset(cx + halfW * 1.2, screenY),
        isScanNear ? scanGridPaint : gridPaint,
      );
    }

    // Center line (body axis)
    canvas.drawLine(
      Offset(cx, cy + (-0.98) * fs),
      Offset(cx, cy + (1.02) * fs),
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 0.3
        ..color = _cyan.withValues(alpha: 0.08)
        ..blendMode = BlendMode.plus,
    );

    // Muscle contour hints
    final musclePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5
      ..color = _cyan.withValues(alpha: 0.10)
      ..blendMode = BlendMode.plus;

    // Rectus abdominis (abs midlines)
    for (final side in [-1.0, 1.0]) {
      final x = cx + side * 0.035 * ws * fs;
      canvas.drawLine(Offset(x, cy + (-0.42) * fs), Offset(x, cy + (-0.05) * fs), musclePaint);
    }
    // Pectoral lines
    for (final side in [-1.0, 1.0]) {
      final path = Path()
        ..moveTo(cx + side * 0.01 * ws * fs, cy + (-0.48) * fs)
        ..quadraticBezierTo(
          cx + side * w.chest * 0.7 * ws * fs, cy + (-0.46) * fs,
          cx + side * w.chest * 0.85 * ws * fs, cy + (-0.52) * fs,
        );
      canvas.drawPath(path, musclePaint);
    }

    canvas.restore();
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 4: Rim Glow (3-layer edge bloom)
  // ══════════════════════════════════════════════════════════════════════

  void _drawRim(Canvas canvas, Path body, List<Path> arms) {
    // Outer glow (wide, dim)
    final outerPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 6.0
      ..color = _cyan.withValues(alpha: 0.12)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 5.0)
      ..blendMode = BlendMode.plus;
    canvas.drawPath(body, outerPaint);
    for (final arm in arms) canvas.drawPath(arm, outerPaint);

    // Mid glow
    final midPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.5
      ..color = _cyan.withValues(alpha: 0.25)
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 2.0)
      ..blendMode = BlendMode.plus;
    canvas.drawPath(body, midPaint);
    for (final arm in arms) canvas.drawPath(arm, midPaint);

    // Inner sharp edge
    final innerPaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.8
      ..color = _cyanBright.withValues(alpha: 0.5)
      ..blendMode = BlendMode.plus;
    canvas.drawPath(body, innerPaint);
    for (final arm in arms) canvas.drawPath(arm, innerPaint);
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 5: Scan Band (unchanged from V20)
  // ══════════════════════════════════════════════════════════════════════

  void _drawScanBand(Canvas canvas, double cx, double cy, double fs,
      Size size, double scanMY) {
    final scanScreenY = cy + scanMY * fs;
    final bandH = size.height * 0.04;
    final bandW = size.width * 0.6;

    final rect = Rect.fromCenter(
      center: Offset(cx, scanScreenY), width: bandW, height: bandH,
    );
    canvas.drawRect(rect, Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          Colors.transparent,
          _cyan.withValues(alpha: 0.2),
          _cyanBright.withValues(alpha: 0.35),
          _cyan.withValues(alpha: 0.2),
          Colors.transparent,
        ],
      ).createShader(rect));

    canvas.drawLine(
      Offset(cx - bandW * 0.5, scanScreenY),
      Offset(cx + bandW * 0.5, scanScreenY),
      Paint()
        ..color = _cyanBright.withValues(alpha: 0.6)
        ..strokeWidth = 0.5,
    );

    for (int i = -5; i <= 5; i++) {
      final tickX = cx + i * (bandW / 10);
      final tickLen = i % 2 == 0 ? 4.0 : 2.0;
      canvas.drawLine(
        Offset(tickX, scanScreenY - tickLen),
        Offset(tickX, scanScreenY + tickLen),
        Paint()
          ..color = _cyan.withValues(alpha: 0.3)
          ..strokeWidth = 0.5,
      );
    }
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 6: Anatomy Overlay (unchanged from V21)
  // ══════════════════════════════════════════════════════════════════════

  void _drawAnatomy(Canvas canvas, double cx, double cy,
      double fs, double cosR, double sinR) {
    if (sinR.abs() > 0.7) return;
    final frontAlpha = (1.0 - sinR.abs() * 1.4).clamp(0.0, 1.0);

    Offset proj(double x, double y, double z) {
      final rx = x * cosR - z * sinR;
      final rz = x * sinR + z * cosR;
      final d = 800.0 - rz * fs;
      final p = 800.0 / d;
      return Offset(cx + rx * fs * p, cy + y * fs * p);
    }

    // Spine: 24 vertebrae
    final spinePaint = Paint()
      ..color = _cyan.withValues(alpha: 0.3 * frontAlpha)
      ..strokeWidth = 1.0
      ..style = PaintingStyle.stroke;
    final vertPaint = Paint()
      ..color = _cyanBright.withValues(alpha: 0.4 * frontAlpha);

    Offset? prev;
    for (int i = 0; i < 24; i++) {
      final t = i / 23.0;
      final y = -0.72 + t * 0.74;
      final z = -0.08 - math.sin(t * math.pi) * 0.03;
      final pt = proj(0, y, z);
      if (prev != null) canvas.drawLine(prev, pt, spinePaint);
      canvas.drawCircle(pt, 1.5, vertPaint);
      prev = pt;
    }

    // Ribs: 6 pairs
    final ribPaint = Paint()
      ..color = _cyan.withValues(alpha: 0.15 * frontAlpha)
      ..strokeWidth = 0.6
      ..style = PaintingStyle.stroke;

    for (int i = 0; i < 6; i++) {
      final y = -0.50 + i * 0.05;
      final width = 0.18 - (i - 3).abs() * 0.02;
      for (final side in [-1.0, 1.0]) {
        final s = proj(0, y, -0.08);
        final e = proj(side * width, y + 0.02, 0.02);
        final m = proj(side * width * 0.6, y - 0.01, -0.02);
        canvas.drawPath(
          Path()..moveTo(s.dx, s.dy)..quadraticBezierTo(m.dx, m.dy, e.dx, e.dy),
          ribPaint,
        );
      }
    }

    // Arteries
    final artPaint = Paint()
      ..color = _red.withValues(alpha: 0.25 * frontAlpha)
      ..strokeWidth = 0.8
      ..style = PaintingStyle.stroke;

    final heart = proj(-0.04, -0.48, 0.06);
    final aortaTop = proj(0, -0.70, 0.04);
    canvas.drawLine(heart, aortaTop, artPaint);
    final aortaDown = proj(0, 0.02, -0.02);
    canvas.drawLine(heart, aortaDown, artPaint);
    for (final side in [-1.0, 1.0]) {
      canvas.drawLine(aortaTop, proj(side * 0.2, -0.55, 0.02), artPaint);
      canvas.drawLine(aortaDown, proj(side * 0.1, 0.5, -0.01), artPaint);
    }

    // Nervous system
    final nervePaint = Paint()
      ..color = _purple.withValues(alpha: 0.15 * frontAlpha)
      ..strokeWidth = 0.5
      ..style = PaintingStyle.stroke;
    canvas.drawLine(proj(0, -0.87, 0.03), proj(0, -0.72, -0.08), nervePaint);
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 7: Organ Systems (unchanged from V21)
  // ══════════════════════════════════════════════════════════════════════

  void _drawOrgans(Canvas canvas, double cx, double cy,
      double fs, double cosR, double sinR) {
    Offset proj(double x, double y, double z) {
      final rx = x * cosR - z * sinR;
      final rz = x * sinR + z * cosR;
      final d = 800.0 - rz * fs;
      final p = 800.0 / d;
      return Offset(cx + rx * fs * p, cy + y * fs * p);
    }

    // Brain
    final brainPos = proj(0, -0.87, 0.03);
    final brainR = fs * 0.055;
    canvas.drawCircle(brainPos, brainR * 1.8, Paint()
      ..shader = RadialGradient(
        colors: [_blue.withValues(alpha: 0.35), Colors.transparent],
      ).createShader(Rect.fromCircle(center: brainPos, radius: brainR * 1.8)));

    final rng = math.Random(42);
    for (int i = 0; i < 8; i++) {
      final angle = i * math.pi / 4 + pulse * 0.3;
      final r = brainR * (0.4 + rng.nextDouble() * 0.5);
      final node = brainPos + Offset(math.cos(angle) * r, math.sin(angle) * r);
      canvas.drawCircle(node, 1.5, Paint()..color = _cyanBright.withValues(alpha: 0.6));
      canvas.drawLine(brainPos, node, Paint()
        ..color = _blue.withValues(alpha: 0.2)
        ..strokeWidth = 0.4);
    }

    // Heart
    final heartPos = proj(-0.04, -0.48, 0.06);
    final hScale = 1.0 + pulse * 0.15;
    final heartR = fs * 0.04 * hScale;
    canvas.drawCircle(heartPos, heartR * 2.5, Paint()
      ..shader = RadialGradient(
        colors: [_red.withValues(alpha: 0.5 * hScale), _red.withValues(alpha: 0.1), Colors.transparent],
        stops: const [0.0, 0.4, 1.0],
      ).createShader(Rect.fromCircle(center: heartPos, radius: heartR * 2.5)));
    canvas.drawCircle(heartPos, heartR * 0.5, Paint()..color = _white.withValues(alpha: 0.8 * hScale));

    // Lungs
    for (final side in [-1.0, 1.0]) {
      final lungPos = proj(side * 0.12, -0.40, 0.02);
      final lungW = fs * 0.06 * (1.0 + breath * 0.08);
      final lungH = fs * 0.1 * (1.0 + breath * 0.08);
      final lr = Rect.fromCenter(center: lungPos, width: lungW, height: lungH);
      canvas.drawOval(lr, Paint()
        ..shader = RadialGradient(
          colors: [_cyan.withValues(alpha: 0.15), Colors.transparent],
        ).createShader(lr));
      canvas.drawOval(lr, Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 0.5
        ..color = _cyan.withValues(alpha: 0.2));
    }

    // Liver
    final liverPos = proj(0.08, -0.22, 0.04);
    canvas.drawCircle(liverPos, fs * 0.035, Paint()
      ..shader = RadialGradient(
        colors: [const Color(0xFFFF6D00).withValues(alpha: 0.15), Colors.transparent],
      ).createShader(Rect.fromCircle(center: liverPos, radius: fs * 0.035)));

    // Kidneys
    for (final side in [-1.0, 1.0]) {
      final kPos = proj(side * 0.1, -0.08, -0.04);
      canvas.drawCircle(kPos, fs * 0.02, Paint()
        ..shader = RadialGradient(
          colors: [const Color(0xFFFFAB00).withValues(alpha: 0.12), Colors.transparent],
        ).createShader(Rect.fromCircle(center: kPos, radius: fs * 0.02)));
    }
  }

  // ══════════════════════════════════════════════════════════════════════
  // Pass 8: Energy Particles (PathMetrics along outline)
  // ══════════════════════════════════════════════════════════════════════

  void _drawParticles(Canvas canvas, Path body) {
    const count = 30;
    final metrics = body.computeMetrics().toList();
    if (metrics.isEmpty) return;

    final totalLen = metrics.fold<double>(0, (sum, m) => sum + m.length);
    if (totalLen <= 0) return;

    for (int i = 0; i < count; i++) {
      final baseT = (i / count + particlePhase) % 1.0;
      var dist = baseT * totalLen;

      Tangent? tangent;
      for (final m in metrics) {
        if (dist <= m.length) {
          tangent = m.getTangentForOffset(dist);
          break;
        }
        dist -= m.length;
      }
      if (tangent == null) continue;

      final pos = tangent.position;
      final life = ((math.sin(particlePhase * math.pi * 2 + i * 0.7) + 1.0) * 0.5);
      final radius = 1.0 + life * 1.5;
      final alpha = (0.3 + life * 0.5).clamp(0.0, 0.8);

      canvas.drawCircle(pos, radius, Paint()..color = _cyanBright.withValues(alpha: alpha));
      canvas.drawCircle(pos, radius * 3, Paint()..color = _cyan.withValues(alpha: alpha * 0.15));
    }
  }

  @override
  bool shouldRepaint(covariant _HoloV22Painter old) =>
      rotation != old.rotation ||
      scanY != old.scanY ||
      pulse != old.pulse ||
      breath != old.breath ||
      particlePhase != old.particlePhase ||
      gender != old.gender;
}
