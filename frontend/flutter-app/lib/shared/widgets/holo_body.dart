import 'dart:math' as math;
import 'dart:ui';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';

import 'holo_body_web_stub.dart'
    if (dart.library.js_interop) 'holo_body_web_impl.dart';
import 'holo_body_native_stub.dart'
    if (dart.library.io) 'holo_body_native_impl.dart';

/// HoloBody: Three.js GLB hologram (web) + V22 Bezier CustomPainter (non-web)

enum HoloGender { male, female }

enum HoloBodyProfile {
  adultMale('adult_male', 'male'),
  adultFemale('adult_female', 'female'),
  juniorMale('junior_male', 'male'),
  juniorFemale('junior_female', 'female'),
  baby('baby', 'male');

  final String key;
  final String genderKey;

  const HoloBodyProfile(this.key, this.genderKey);

  HoloGender get fallbackGender =>
      genderKey == 'female' ? HoloGender.female : HoloGender.male;
}

class HoloBody extends StatefulWidget {
  final double width;
  final double height;
  final Color color;
  final Color? accentColor;
  final HoloGender gender;
  final HoloBodyProfile? profile;
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
    this.profile,
    this.bioData = const {},
    this.showEcg = false,
    this.showHud = false,
  });

  @override
  State<HoloBody> createState() => _HoloBodyState();
}

class _HoloBodyState extends State<HoloBody> with TickerProviderStateMixin {
  late AnimationController _scanController;
  late AnimationController _pulseController;
  late AnimationController _breathController;
  late AnimationController _particleController;

  @override
  void initState() {
    super.initState();
    _scanController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 5),
    )..repeat();
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 900),
    )..repeat(reverse: true);
    _breathController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 3500),
    )..repeat(reverse: true);
    _particleController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 8),
    )..repeat();
  }

  @override
  void dispose() {
    _scanController.dispose();
    _pulseController.dispose();
    _breathController.dispose();
    _particleController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // SECURITY/UI FIX: GoRouter's ShellRoute keeps screens alive in an IndexedStack.
    // However, native WebViews (CEF/Iframe) often ignore Flutter's 'Offstage' clipping
    // and can overlap/stack if multiple tabs have holograms (e.g. Home & Data).
    // By checking TickerMode, we forcefully unmount the WebView when the tab is inactive!
    if (!TickerMode.of(context)) {
      return SizedBox(width: widget.width, height: widget.height);
    }

    return SizedBox(
      width: widget.width,
      height: widget.height,
      child: ExcludeSemantics(
        child: kIsWeb ? _buildWebContent() : _buildNativeOrPainter(),
      ),
    );
  }

  /// Web: Three.js GLB is the default renderer.
  Widget _buildWebContent() {
    final colorHex =
        widget.color.toARGB32().toRadixString(16).padLeft(8, '0').substring(2);
    final effectiveProfile = widget.profile ??
        (widget.gender == HoloGender.female
            ? HoloBodyProfile.adultFemale
            : HoloBodyProfile.adultMale);
    return buildHoloBodyWeb(effectiveProfile.key, colorHex);
  }

  Widget _buildNativeOrPainter() {
    final colorHex =
        widget.color.toARGB32().toRadixString(16).padLeft(8, '0').substring(2);
    final effectiveProfile = widget.profile ??
        (widget.gender == HoloGender.female
            ? HoloBodyProfile.adultFemale
            : HoloBodyProfile.adultMale);

    final nativeGlbView = buildHoloBodyNative(effectiveProfile.key, colorHex);
    if (nativeGlbView != null) {
      return nativeGlbView;
    }

    return _buildPainter();
  }

  /// V22 Bezier CustomPainter (fallback for unsupported native WebView platforms)
  Widget _buildPainter() {
    final effectiveGender = widget.profile?.fallbackGender ?? widget.gender;
    return RepaintBoundary(
      child: AnimatedBuilder(
        animation: Listenable.merge([
          _scanController,
          _pulseController,
          _breathController,
          _particleController,
        ]),
        builder: (context, child) {
          return CustomPaint(
            size: Size(widget.width, widget.height),
            painter: _HoloV22Painter(
              gender: effectiveGender,
              scanY: _scanController.value,
              pulse: _pulseController.value,
              breath: _breathController.value,
              particlePhase: _particleController.value,
              color: widget.color,
            ),
          );
        },
      ),
    );
  }
}

// ?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧??// Body width profiles per gender (half-widths from center)
// ?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧??
class _BW {
  final double head, jaw, neck, shoulder, chest, waist, hip;
  final double thigh, knee, calf, ankle, foot;

  const _BW({
    required this.head,
    required this.jaw,
    required this.neck,
    required this.shoulder,
    required this.chest,
    required this.waist,
    required this.hip,
    required this.thigh,
    required this.knee,
    required this.calf,
    required this.ankle,
    required this.foot,
  });

  static const male = _BW(
    head: 0.115,
    jaw: 0.100,
    neck: 0.060,
    shoulder: 0.340,
    chest: 0.290,
    waist: 0.190,
    hip: 0.225,
    thigh: 0.145,
    knee: 0.088,
    calf: 0.095,
    ankle: 0.048,
    foot: 0.068,
  );

  static const female = _BW(
    head: 0.110,
    jaw: 0.096,
    neck: 0.052,
    shoulder: 0.265,
    chest: 0.255,
    waist: 0.170,
    hip: 0.295,
    thigh: 0.155,
    knee: 0.078,
    calf: 0.082,
    ankle: 0.044,
    foot: 0.060,
  );
}

// ?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧??// V22 Multi-Pass Hologram Painter (fallback for non-web)
// ?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧?먥븧??
class _HoloV22Painter extends CustomPainter {
  final HoloGender gender;
  final double scanY;
  final double pulse;
  final double breath;
  final double particlePhase;
  final Color color;

  _HoloV22Painter({
    required this.gender,
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

    // Fixed front-facing view — no rotation to avoid visual conflict
    const ws = 1.0;
    final w = gender == HoloGender.male ? _BW.male : _BW.female;

    final body = _buildBody(w, fs, cx, cy, ws);
    final arms = [
      _buildArm(w, fs, cx, cy, ws, 1),
      _buildArm(w, fs, cx, cy, ws, -1)
    ];
    final combined = Path()
      ..addPath(body, Offset.zero)
      ..addPath(arms[0], Offset.zero)
      ..addPath(arms[1], Offset.zero);

    final scanMY = (scanY * 2.4) - 1.2;

    _drawAura(canvas, cx, cy, size);
    _drawVolume(canvas, body, arms, fs, cx, cy, scanMY, w, ws);
    _drawScanFill(canvas, combined, fs, cx, cy, scanMY, w, ws);
    _drawGrid(canvas, combined, w, fs, cx, cy, ws, scanMY);
    _drawRim(canvas, body, arms);
    _drawScanBand(canvas, cx, cy, fs, size, scanMY, w, ws);
    _drawAnatomy(canvas, cx, cy, fs);
    _drawOrgans(canvas, cx, cy, fs);
    _drawParticles(canvas, body);
  }

  Path _buildBody(_BW w, double fs, double cx, double cy, double ws) {
    double sx(double x) => cx + x * ws * fs;
    double sy(double y) => cy + y * fs;
    const g = 0.022;

    final p = Path()..moveTo(sx(0), sy(-1.0));

    p.cubicTo(sx(w.head * 0.55), sy(-1.0), sx(w.head), sy(-0.95), sx(w.head),
        sy(-0.87));
    p.cubicTo(sx(w.head), sy(-0.82), sx(w.jaw), sy(-0.79), sx(w.jaw * 0.85),
        sy(-0.76));
    p.cubicTo(sx(w.jaw * 0.7), sy(-0.74), sx(w.neck * 1.1), sy(-0.73),
        sx(w.neck), sy(-0.70));
    p.cubicTo(sx(w.neck), sy(-0.67), sx(w.shoulder * 0.5), sy(-0.64),
        sx(w.shoulder), sy(-0.62));
    p.cubicTo(sx(w.shoulder * 0.98), sy(-0.58), sx(w.chest * 1.05), sy(-0.54),
        sx(w.chest), sy(-0.48));
    p.cubicTo(sx(w.chest * 0.97), sy(-0.40), sx(w.waist * 1.15), sy(-0.32),
        sx(w.waist), sy(-0.24));
    p.cubicTo(sx(w.waist * 0.97), sy(-0.17), sx(w.hip * 0.85), sy(-0.08),
        sx(w.hip), sy(-0.02));
    p.cubicTo(sx(w.hip * 0.97), sy(0.02), sx(w.thigh * 1.25), sy(0.06),
        sx(w.thigh), sy(0.12));
    p.cubicTo(sx(w.thigh * 0.97), sy(0.25), sx(w.knee * 1.15), sy(0.42),
        sx(w.knee), sy(0.50));
    p.cubicTo(sx(w.knee * 0.98), sy(0.53), sx(w.calf * 1.05), sy(0.58),
        sx(w.calf), sy(0.68));
    p.cubicTo(sx(w.calf * 0.85), sy(0.78), sx(w.ankle * 1.3), sy(0.87),
        sx(w.ankle), sy(0.92));
    p.cubicTo(sx(w.ankle * 0.95), sy(0.96), sx(w.foot), sy(1.01), sx(w.foot),
        sy(1.037));
    p.lineTo(sx(g), sy(1.037));
    p.cubicTo(sx(g), sy(1.01), sx(w.ankle * 0.55), sy(0.96), sx(w.ankle * 0.5),
        sy(0.92));
    p.cubicTo(sx(w.ankle * 0.55), sy(0.87), sx(w.calf * 0.42), sy(0.78),
        sx(w.calf * 0.40), sy(0.68));
    p.cubicTo(sx(w.calf * 0.38), sy(0.58), sx(w.knee * 0.55), sy(0.53),
        sx(w.knee * 0.50), sy(0.50));
    p.cubicTo(sx(w.knee * 0.55), sy(0.42), sx(w.thigh * 0.55), sy(0.25),
        sx(w.thigh * 0.45), sy(0.12));
    p.cubicTo(
        sx(w.thigh * 0.35), sy(0.06), sx(g * 1.5), sy(0.04), sx(g), sy(0.03));
    p.cubicTo(sx(g * 0.5), sy(0.05), sx(0), sy(0.065), sx(-g), sy(0.03));
    p.cubicTo(sx(-g * 1.5), sy(0.04), sx(-w.thigh * 0.35), sy(0.06),
        sx(-w.thigh * 0.45), sy(0.12));
    p.cubicTo(sx(-w.thigh * 0.55), sy(0.25), sx(-w.knee * 0.55), sy(0.42),
        sx(-w.knee * 0.50), sy(0.50));
    p.cubicTo(sx(-w.calf * 0.38), sy(0.58), sx(-w.calf * 0.42), sy(0.78),
        sx(-w.calf * 0.40), sy(0.68));
    p.cubicTo(
        sx(-w.ankle * 0.55), sy(0.87), sx(-g), sy(1.01), sx(-g), sy(1.037));
    p.lineTo(sx(-w.foot), sy(1.037));
    p.cubicTo(sx(-w.foot), sy(1.01), sx(-w.ankle * 0.95), sy(0.96),
        sx(-w.ankle), sy(0.92));
    p.cubicTo(sx(-w.ankle * 1.3), sy(0.87), sx(-w.calf * 0.85), sy(0.78),
        sx(-w.calf), sy(0.68));
    p.cubicTo(sx(-w.calf * 1.05), sy(0.58), sx(-w.knee * 0.98), sy(0.53),
        sx(-w.knee), sy(0.50));
    p.cubicTo(sx(-w.knee * 1.15), sy(0.42), sx(-w.thigh * 0.97), sy(0.25),
        sx(-w.thigh), sy(0.12));
    p.cubicTo(sx(-w.thigh * 1.25), sy(0.06), sx(-w.hip * 0.97), sy(0.02),
        sx(-w.hip), sy(-0.02));
    p.cubicTo(sx(-w.hip * 0.85), sy(-0.08), sx(-w.waist * 0.97), sy(-0.17),
        sx(-w.waist), sy(-0.24));
    p.cubicTo(sx(-w.waist * 1.15), sy(-0.32), sx(-w.chest * 0.97), sy(-0.40),
        sx(-w.chest), sy(-0.48));
    p.cubicTo(sx(-w.chest * 1.05), sy(-0.54), sx(-w.shoulder * 0.98), sy(-0.58),
        sx(-w.shoulder), sy(-0.62));
    p.cubicTo(sx(-w.shoulder * 0.5), sy(-0.64), sx(-w.neck), sy(-0.67),
        sx(-w.neck), sy(-0.70));
    p.cubicTo(sx(-w.neck * 1.1), sy(-0.73), sx(-w.jaw * 0.7), sy(-0.74),
        sx(-w.jaw * 0.85), sy(-0.76));
    p.cubicTo(
        sx(-w.jaw), sy(-0.79), sx(-w.head), sy(-0.82), sx(-w.head), sy(-0.87));
    p.cubicTo(
        sx(-w.head), sy(-0.95), sx(-w.head * 0.55), sy(-1.0), sx(0), sy(-1.0));
    p.close();
    return p;
  }

  Path _buildArm(
      _BW w, double fs, double cx, double cy, double ws, double side) {
    double sx(double x) => cx + x * side * ws * fs;
    double sy(double y) => cy + y * fs;
    final s = w.shoulder;

    final p = Path()..moveTo(sx(s), sy(-0.62));
    p.cubicTo(sx(s * 0.93), sy(-0.52), sx(s * 0.84), sy(-0.38), sx(s * 0.78),
        sy(-0.24));
    p.cubicTo(sx(s * 0.74), sy(-0.12), sx(s * 0.70), sy(0.00), sx(s * 0.67),
        sy(0.10));
    p.cubicTo(
        sx(s * 0.65), sy(0.14), sx(s * 0.64), sy(0.17), sx(s * 0.62), sy(0.19));
    p.cubicTo(
        sx(s * 0.60), sy(0.21), sx(s * 0.55), sy(0.21), sx(s * 0.54), sy(0.19));
    p.cubicTo(sx(s * 0.56), sy(0.14), sx(s * 0.58), sy(0.00), sx(s * 0.62),
        sy(-0.24));
    p.cubicTo(sx(s * 0.66), sy(-0.38), sx(s * 0.72), sy(-0.52), sx(s * 0.80),
        sy(-0.60));
    p.lineTo(sx(s * 0.88), sy(-0.62));
    p.close();
    return p;
  }

  static double _bodyW(double y, _BW w, double ws) {
    const ys = [
      -1.0,
      -0.87,
      -0.76,
      -0.70,
      -0.62,
      -0.48,
      -0.24,
      -0.02,
      0.12,
      0.50,
      0.68,
      0.92,
      1.037
    ];
    final ws2 = [
      0.0,
      w.head,
      w.jaw * 0.85,
      w.neck,
      w.shoulder,
      w.chest,
      w.waist,
      w.hip,
      w.thigh,
      w.knee,
      w.calf,
      w.ankle,
      w.foot
    ];
    for (int i = 0; i < ys.length - 1; i++) {
      if (y <= ys[i + 1]) {
        final t = (y - ys[i]) / (ys[i + 1] - ys[i]);
        return (ws2[i] + t * (ws2[i + 1] - ws2[i])) * ws;
      }
    }
    return ws2.last * ws;
  }

  void _drawAura(Canvas canvas, double cx, double cy, Size size) {
    final radius = size.height * 0.45;
    final paint = Paint()
      ..shader = RadialGradient(
        colors: [
          _cyan.withOpacity(0.08),
          _cyanDark.withOpacity(0.04),
          Colors.transparent,
        ],
        stops: const [0.0, 0.5, 1.0],
      ).createShader(Rect.fromCircle(center: Offset(cx, cy), radius: radius));
    canvas.drawCircle(Offset(cx, cy), radius, paint);
  }

  void _drawVolume(Canvas canvas, Path body, List<Path> arms, double fs,
      double cx, double cy, double scanMY, _BW w, double ws) {
    final combined = Path.combine(
      PathOperation.union,
      body,
      arms.fold<Path>(
          Path(), (acc, a) => Path.combine(PathOperation.union, acc, a)),
    );
    final bounds = combined.getBounds();

    canvas.save();
    canvas.clipPath(combined);

    // (a) Vertical depth gradient — chest core brightest, extremities dim
    canvas.drawRect(
      bounds,
      Paint()
        ..shader = LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            _cyan.withOpacity(0.03),
            _cyan.withOpacity(0.07),
            _cyan.withOpacity(0.14),
            _cyan.withOpacity(0.12),
            _cyan.withOpacity(0.05),
            _cyan.withOpacity(0.02),
          ],
          stops: const [0.0, 0.20, 0.42, 0.55, 0.75, 1.0],
        ).createShader(bounds)
        ..blendMode = BlendMode.plus,
    );

    // (b) Contour-adaptive Fresnel — follows body width at each Y level
    const stripCount = 20;
    final bodyTop = cy + (-1.0) * fs;
    final bodyBot = cy + 1.037 * fs;
    final totalH = bodyBot - bodyTop;
    final stripH = totalH / stripCount;

    for (int i = 0; i < stripCount; i++) {
      final stripY = bodyTop + i * stripH;
      final meshY = -1.0 + (i + 0.5) / stripCount * 2.037;
      final bw = _bodyW(meshY, w, ws) * fs;
      if (bw <= 0) continue;
      final stripRect = Rect.fromLTRB(
          cx - bw * 1.15, stripY, cx + bw * 1.15, stripY + stripH);
      canvas.drawRect(
        stripRect,
        Paint()
          ..shader = LinearGradient(
            begin: Alignment.centerLeft,
            end: Alignment.centerRight,
            colors: [
              _cyan.withOpacity(0.18),
              _cyan.withOpacity(0.01),
              _cyan.withOpacity(0.01),
              _cyan.withOpacity(0.18),
            ],
            stops: const [0.0, 0.30, 0.70, 1.0],
          ).createShader(stripRect)
          ..blendMode = BlendMode.plus,
      );
    }

    // (c) Inner glow — subtle radial from chest center
    final glowCenter = Offset(cx, cy + (-0.45) * fs);
    final glowRadius = fs * 0.65;
    canvas.drawCircle(
      glowCenter,
      glowRadius,
      Paint()
        ..shader = RadialGradient(
          colors: [
            _cyan.withOpacity(0.06),
            _cyan.withOpacity(0.02),
            Colors.transparent,
          ],
          stops: const [0.0, 0.4, 1.0],
        ).createShader(
            Rect.fromCircle(center: glowCenter, radius: glowRadius))
        ..blendMode = BlendMode.plus,
    );

    canvas.restore();
  }

  void _drawScanFill(Canvas canvas, Path combined, double fs, double cx,
      double cy, double scanMY, _BW w, double ws) {
    final bounds = combined.getBounds();
    final scanScreenY = cy + scanMY * fs;

    canvas.save();
    canvas.clipPath(combined);

    // Scan highlight band — clipped to body contour
    final bw = _bodyW(scanMY, w, ws) * fs;
    final scanW = (bw > 0 ? bw : bounds.width * 0.4) * 1.05;
    final scanRect = Rect.fromCenter(
      center: Offset(cx, scanScreenY),
      width: scanW * 2,
      height: fs * 0.07,
    );
    canvas.drawRect(
      scanRect,
      Paint()
        ..shader = LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Colors.transparent,
            _cyanBright.withOpacity(0.10),
            _white.withOpacity(0.15),
            _cyanBright.withOpacity(0.10),
            Colors.transparent,
          ],
        ).createShader(scanRect)
        ..blendMode = BlendMode.plus,
    );

    // Afterglow trail below scan line — very subtle
    final trailH = fs * 0.22;
    final trailRect = Rect.fromLTRB(
        bounds.left, scanScreenY, bounds.right, scanScreenY + trailH);
    canvas.drawRect(
      trailRect,
      Paint()
        ..shader = LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            _cyan.withOpacity(0.05),
            _cyan.withOpacity(0.02),
            Colors.transparent,
          ],
          stops: const [0.0, 0.35, 1.0],
        ).createShader(trailRect)
        ..blendMode = BlendMode.plus,
    );

    canvas.restore();
  }

  void _drawGrid(Canvas canvas, Path combined, _BW w, double fs, double cx,
      double cy, double ws, double scanMY) {
    canvas.save();
    canvas.clipPath(combined);

    // Gaussian falloff constants
    const sigma = 0.07;
    const sigSq2 = 2.0 * sigma * sigma;

    const step = 0.04;
    for (double y = -1.0; y <= 1.04; y += step) {
      final bw = _bodyW(y, w, ws);
      final screenY = cy + y * fs;
      final halfW = bw * fs;

      // Smooth distance-based scan influence
      final dist = (y - scanMY).abs();
      final scanInfluence = math.exp(-dist * dist / sigSq2);
      final alpha = 0.08 + scanInfluence * 0.25;
      final sw = 0.35 + scanInfluence * 0.45;

      final paint = Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = sw
        ..color = Color.lerp(
            _cyan, _cyanBright, scanInfluence)!
            .withOpacity(alpha)
        ..blendMode = BlendMode.plus;

      canvas.drawLine(
        Offset(cx - halfW * 1.2, screenY),
        Offset(cx + halfW * 1.2, screenY),
        paint,
      );
    }

    canvas.drawLine(
      Offset(cx, cy + (-0.98) * fs),
      Offset(cx, cy + (1.02) * fs),
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 0.3
        ..color = _cyan.withOpacity(0.08)
        ..blendMode = BlendMode.plus,
    );

    final musclePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.5
      ..color = _cyan.withOpacity(0.10)
      ..blendMode = BlendMode.plus;

    for (final side in [-1.0, 1.0]) {
      final x = cx + side * 0.035 * ws * fs;
      canvas.drawLine(Offset(x, cy + (-0.42) * fs),
          Offset(x, cy + (-0.05) * fs), musclePaint);
    }
    for (final side in [-1.0, 1.0]) {
      final path = Path()
        ..moveTo(cx + side * 0.01 * ws * fs, cy + (-0.48) * fs)
        ..quadraticBezierTo(
          cx + side * w.chest * 0.7 * ws * fs,
          cy + (-0.46) * fs,
          cx + side * w.chest * 0.85 * ws * fs,
          cy + (-0.52) * fs,
        );
      canvas.drawPath(path, musclePaint);
    }
    canvas.restore();
  }

  void _drawRim(Canvas canvas, Path body, List<Path> arms) {
    // Merge into one path to avoid body/arm junction double-stroke
    final combined = Path.combine(
      PathOperation.union,
      body,
      arms.fold<Path>(
          Path(), (acc, a) => Path.combine(PathOperation.union, acc, a)),
    );

    // Soft outer glow
    canvas.drawPath(
        combined,
        Paint()
          ..style = PaintingStyle.stroke
          ..strokeWidth = 5.0
          ..color = _cyan.withOpacity(0.18)
          ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 4.0)
          ..blendMode = BlendMode.plus);

    // Crisp edge line
    canvas.drawPath(
        combined,
        Paint()
          ..style = PaintingStyle.stroke
          ..strokeWidth = 0.8
          ..color = _cyan.withOpacity(0.70)
          ..blendMode = BlendMode.plus);
  }

  void _drawScanBand(Canvas canvas, double cx, double cy, double fs, Size size,
      double scanMY, _BW w, double ws) {
    final scanScreenY = cy + scanMY * fs;

    // Body-contour following scan width
    final bw = _bodyW(scanMY, w, ws) * fs;
    final contourW = (bw > 0 ? bw * 2.2 : size.width * 0.54);
    final breathW = contourW * (0.94 + 0.06 * math.sin(particlePhase * math.pi * 2));
    final coreH = fs * 0.026;

    // Soft halo — reduced height, follows body width
    final haloH = fs * 0.10;
    final haloRect = Rect.fromCenter(
      center: Offset(cx, scanScreenY),
      width: contourW * 1.05,
      height: haloH,
    );
    canvas.drawOval(
      haloRect,
      Paint()
        ..shader = RadialGradient(
          colors: [
            _cyanBright.withOpacity(0.18),
            _cyan.withOpacity(0.08),
            Colors.transparent,
          ],
          stops: const [0.0, 0.45, 1.0],
        ).createShader(haloRect)
        ..blendMode = BlendMode.plus,
    );

    // Core scan line — body-contour width
    final coreRect = Rect.fromCenter(
      center: Offset(cx, scanScreenY),
      width: breathW,
      height: coreH,
    );
    canvas.drawRRect(
      RRect.fromRectAndRadius(coreRect, Radius.circular(coreH)),
      Paint()
        ..shader = LinearGradient(
          begin: Alignment.centerLeft,
          end: Alignment.centerRight,
          colors: [
            Colors.transparent,
            _cyan.withOpacity(0.12),
            _cyanBright.withOpacity(0.45),
            _white.withOpacity(0.50),
            _cyanBright.withOpacity(0.45),
            _cyan.withOpacity(0.12),
            Colors.transparent,
          ],
          stops: const [0.0, 0.10, 0.35, 0.5, 0.65, 0.90, 1.0],
        ).createShader(coreRect)
        ..blendMode = BlendMode.plus,
    );

    // Soft sparks — golden ratio spacing + blur glow
    const sparkCount = 9;
    const golden = 0.618033988;
    final softSparkPaint = Paint()
      ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 2.0)
      ..blendMode = BlendMode.plus;
    final dotPaint = Paint()..blendMode = BlendMode.plus;

    for (var i = 0; i < sparkCount; i++) {
      final t = ((i * golden + particlePhase * 0.3) % 1.0);
      final x = cx - breathW * 0.46 + breathW * 0.92 * t;
      final wobble = math.sin((particlePhase * 6.2) + i * 1.1) * fs * 0.005;
      final sparkY = scanScreenY + wobble;
      final life = (math.sin(particlePhase * math.pi * 2 + i * 0.9) + 1.0) * 0.5;
      final sparkR = fs * (0.003 + 0.002 * life);
      final alpha = (0.25 + life * 0.35).clamp(0.0, 0.6);

      // Soft outer glow
      softSparkPaint.color = _cyan.withOpacity(alpha * 0.4);
      canvas.drawCircle(Offset(x, sparkY), sparkR * 3.5, softSparkPaint);
      // Bright core
      dotPaint.color = _cyanBright.withOpacity(alpha);
      canvas.drawCircle(Offset(x, sparkY), sparkR, dotPaint);
    }

    // Edge emission sparks — at body contour intersection
    if (bw > 0) {
      final edgePaint = Paint()
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 1.5)
        ..blendMode = BlendMode.plus;
      for (final side in [-1.0, 1.0]) {
        final edgeX = cx + side * bw;
        for (var j = 0; j < 3; j++) {
          final phase = (particlePhase * 4.0 + j * 1.2 + side * 0.5) % 1.0;
          final scatter = phase * fs * 0.03 * side;
          final driftY = math.sin(phase * math.pi) * fs * 0.015;
          final eAlpha = (1.0 - phase) * 0.4;
          final eR = fs * (0.002 + 0.003 * (1.0 - phase));
          edgePaint.color = _cyanBright.withOpacity(eAlpha);
          canvas.drawCircle(
            Offset(edgeX + scatter, scanScreenY + driftY),
            eR,
            edgePaint,
          );
        }
      }
    }
  }

  void _drawAnatomy(Canvas canvas, double cx, double cy, double fs) {
    // Fixed front-facing projection (cosR=1, sinR=0)
    const frontAlpha = 1.0;

    Offset proj(double x, double y, double z) {
      final d = 800.0 - z * fs;
      final p = 800.0 / d;
      return Offset(cx + x * fs * p, cy + y * fs * p);
    }

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
          Path()
            ..moveTo(s.dx, s.dy)
            ..quadraticBezierTo(m.dx, m.dy, e.dx, e.dy),
          ribPaint,
        );
      }
    }

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

    final nervePaint = Paint()
      ..color = _purple.withValues(alpha: 0.15 * frontAlpha)
      ..strokeWidth = 0.5
      ..style = PaintingStyle.stroke;
    canvas.drawLine(proj(0, -0.87, 0.03), proj(0, -0.72, -0.08), nervePaint);
  }

  void _drawOrgans(Canvas canvas, double cx, double cy, double fs) {
    // Fixed front-facing projection (cosR=1, sinR=0)
    Offset proj(double x, double y, double z) {
      final d = 800.0 - z * fs;
      final p = 800.0 / d;
      return Offset(cx + x * fs * p, cy + y * fs * p);
    }

    final brainPos = proj(0, -0.87, 0.03);
    final brainR = fs * 0.055;
    canvas.drawCircle(
        brainPos,
        brainR * 1.8,
        Paint()
          ..shader = RadialGradient(
            colors: [_blue.withValues(alpha: 0.35), Colors.transparent],
          ).createShader(
              Rect.fromCircle(center: brainPos, radius: brainR * 1.8)));

    final rng = math.Random(42);
    for (int i = 0; i < 8; i++) {
      final angle = i * math.pi / 4 + pulse * 0.3;
      final r = brainR * (0.4 + rng.nextDouble() * 0.5);
      final node = brainPos + Offset(math.cos(angle) * r, math.sin(angle) * r);
      canvas.drawCircle(
          node, 1.5, Paint()..color = _cyanBright.withValues(alpha: 0.6));
      canvas.drawLine(
          brainPos,
          node,
          Paint()
            ..color = _blue.withValues(alpha: 0.2)
            ..strokeWidth = 0.4);
    }

    final heartPos = proj(-0.04, -0.48, 0.06);
    final hScale = 1.0 + pulse * 0.15;
    final heartR = fs * 0.04 * hScale;
    canvas.drawCircle(
        heartPos,
        heartR * 2.5,
        Paint()
          ..shader = RadialGradient(
            colors: [
              _red.withValues(alpha: 0.5 * hScale),
              _red.withValues(alpha: 0.1),
              Colors.transparent
            ],
            stops: const [0.0, 0.4, 1.0],
          ).createShader(
              Rect.fromCircle(center: heartPos, radius: heartR * 2.5)));
    canvas.drawCircle(heartPos, heartR * 0.5,
        Paint()..color = _white.withValues(alpha: 0.8 * hScale));

    for (final side in [-1.0, 1.0]) {
      final lungPos = proj(side * 0.12, -0.40, 0.02);
      final lungW = fs * 0.06 * (1.0 + breath * 0.08);
      final lungH = fs * 0.1 * (1.0 + breath * 0.08);
      final lr = Rect.fromCenter(center: lungPos, width: lungW, height: lungH);
      canvas.drawOval(
          lr,
          Paint()
            ..shader = RadialGradient(
              colors: [_cyan.withValues(alpha: 0.15), Colors.transparent],
            ).createShader(lr));
      canvas.drawOval(
          lr,
          Paint()
            ..style = PaintingStyle.stroke
            ..strokeWidth = 0.5
            ..color = _cyan.withValues(alpha: 0.2));
    }

    final liverPos = proj(0.08, -0.22, 0.04);
    canvas.drawCircle(
        liverPos,
        fs * 0.035,
        Paint()
          ..shader = RadialGradient(
            colors: [
              const Color(0xFFFF6D00).withValues(alpha: 0.15),
              Colors.transparent
            ],
          ).createShader(
              Rect.fromCircle(center: liverPos, radius: fs * 0.035)));

    for (final side in [-1.0, 1.0]) {
      final kPos = proj(side * 0.1, -0.08, -0.04);
      canvas.drawCircle(
          kPos,
          fs * 0.02,
          Paint()
            ..shader = RadialGradient(
              colors: [
                const Color(0xFFFFAB00).withValues(alpha: 0.12),
                Colors.transparent
              ],
            ).createShader(Rect.fromCircle(center: kPos, radius: fs * 0.02)));
    }
  }

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
      final life =
          ((math.sin(particlePhase * math.pi * 2 + i * 0.7) + 1.0) * 0.5);
      final radius = 1.0 + life * 1.5;
      final alpha = (0.3 + life * 0.5).clamp(0.0, 0.8);

      canvas.drawCircle(
          pos, radius, Paint()..color = _cyanBright.withValues(alpha: alpha));
      canvas.drawCircle(pos, radius * 3,
          Paint()..color = _cyan.withValues(alpha: alpha * 0.15));
    }
  }

  @override
  bool shouldRepaint(covariant _HoloV22Painter old) =>
      scanY != old.scanY ||
      pulse != old.pulse ||
      breath != old.breath ||
      particlePhase != old.particlePhase ||
      gender != old.gender;
}
