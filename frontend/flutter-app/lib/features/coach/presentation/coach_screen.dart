import 'dart:math' as math;
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ═══════════════════════════════════════════════════════════════════
// 만파식(萬波息) AI 주치의 — "24시간 나를 지키는 주치의"
//
// 5대 생명주기 스토리라인 最終 종착지: Prevent & Care
//
// ┌─ 5. 적극적 예방 ────────────────────────────────────────────┐
// │  이 화면은 Measure → Analyze → Accumulate → Predict의      │
// │  모든 데이터가 하나의 '영혼 오브(Spirit Orb)'에 집약되어     │
// │  사용자를 24시간 보호하는 만파식 철학의 최종 구현이다.        │
// │                                                              │
// │  - 위험 예측 시 오브에 붉은 심박동 애니메이션 (Breathing)     │
// │  - 평온 시 자개 시안/마젠타 빛의 명상적 펄스                  │
// │  - 즉시 원격진료 연결 플로팅 3D 액션 버튼                    │
// │  - AI 학습 진행률 인디케이터 (축적 데이터 시각화)             │
// │  - 예측 트렌드 경고 배너 (Sanggam Gold)                     │
// └──────────────────────────────────────────────────────────────┘
//
// [Rule 1] Z-축 격리 — 텍스트 반드시 SanggamContainer/BackdropFilter 내부
// [Rule 2] 8px 그리드 — 간격 8/16/24/32px, 벤토 박스
// [Rule 3] Slim-fit — SlimHeader + Expanded + InputPanel
// [Rule 4] SanggamTheme 상수 + withValues(alpha:) 통일
// ═══════════════════════════════════════════════════════════════════

// ── Layout constants (8px grid) ──
const double _kHeaderH = 36;
const double _kOrbDiameter = 120.0;
const double _kOrbZoneH = 200;
const double _kBubbleMaxRatio = 0.75;
const double _kBubbleVPad = 16.0;
const double _kBubbleHPad = 24.0;
const double _kBubbleGap = 16.0;

// ── Risk levels ──
enum _RiskLevel { safe, caution, danger }

// ═══════════════════════════════════════════════════════════════════
//  CoachScreen — 5단계 생명주기 최종 종착지
// ═══════════════════════════════════════════════════════════════════

class CoachScreen extends StatefulWidget {
  const CoachScreen({super.key});

  @override
  State<CoachScreen> createState() => _CoachScreenState();
}

class _CoachScreenState extends State<CoachScreen>
    with TickerProviderStateMixin {
  // ── Animations ──
  late AnimationController _orbPulseCtrl;   // 오브 호흡 펄스
  late AnimationController _orbRotateCtrl;  // 오브 궤도 회전
  late AnimationController _riskGlowCtrl;   // 위험 경고 글로우
  late AnimationController _particleCtrl;   // 파티클 시스템
  late AnimationController _predictionCtrl; // 예측 배너 등장

  late Animation<double> _orbPulse;
  late Animation<double> _riskGlow;
  late Animation<double> _predictionFade;

  // ── State ──
  final _msgCtrl = TextEditingController();
  final _scrollCtrl = ScrollController();
  final _RiskLevel _riskLevel = _RiskLevel.caution; // 시뮬레이션: 주의 상태
  double _aiLearningProgress = 0.72;          // AI 학습 진행률 (72%)
  final int _accumulatedDays = 147;                 // 축적 데이터 일수

  final List<_ChatMsg> _messages = [
    const _ChatMsg(
      text: '홍길동님, 현재 건강 위험도 "주의" 단계입니다.\n'
          '최근 3일간 수면 불규칙 패턴이 감지되었고,\n'
          '현 추세라면 2주 뒤 스트레스 지수가 위험 수치에\n'
          '도달할 것으로 예측됩니다.',
      isAi: true,
      time: '오후 10:20',
      msgType: _MsgType.prediction,
    ),
    const _ChatMsg(
      text: '어떻게 하면 좋을까요?',
      isAi: false,
      time: '오후 10:21',
    ),
    const _ChatMsg(
      text: '세 가지 예방 플랜을 제안합니다:\n\n'
          '① 오늘 밤 22시 수면 명상 모드 가동\n'
          '② 내일 아침 10분 호흡 운동 (AI 가이드)\n'
          '③ 3일 내 원격 진료로 전문 상담\n\n'
          '바로 실행할 항목을 선택해주세요.',
      isAi: true,
      time: '오후 10:22',
      msgType: _MsgType.actionPlan,
    ),
  ];

  @override
  void initState() {
    super.initState();

    // 1) 오브 호흡 펄스 — 위험도에 따라 속도 변경
    _orbPulseCtrl = AnimationController(
      vsync: this,
      duration: _pulseDuration,
    )..repeat(reverse: true);
    _orbPulse = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _orbPulseCtrl, curve: Curves.easeInOut),
    );

    // 2) 오브 궤도 회전 (60초 1바퀴)
    _orbRotateCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 60),
    )..repeat();

    // 3) 위험 경고 글로우
    _riskGlowCtrl = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    )..repeat(reverse: true);
    _riskGlow = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _riskGlowCtrl, curve: Curves.easeInOut),
    );

    // 4) 파티클 시스템 (30초 순환)
    _particleCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 30),
    )..repeat();

    // 5) 예측 배너 페이드인
    _predictionCtrl = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    )..forward();
    _predictionFade = CurvedAnimation(
      parent: _predictionCtrl,
      curve: Curves.easeOut,
    );
  }

  Duration get _pulseDuration {
    switch (_riskLevel) {
      case _RiskLevel.danger:
        return const Duration(milliseconds: 800);  // 빠른 심박동
      case _RiskLevel.caution:
        return const Duration(milliseconds: 1500); // 경계 호흡
      case _RiskLevel.safe:
        return const Duration(seconds: 4);          // 명상적 펄스
    }
  }

  @override
  void dispose() {
    _orbPulseCtrl.dispose();
    _orbRotateCtrl.dispose();
    _riskGlowCtrl.dispose();
    _particleCtrl.dispose();
    _predictionCtrl.dispose();
    _msgCtrl.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  void _onSend() {
    final txt = _msgCtrl.text.trim();
    if (txt.isEmpty) return;
    setState(() {
      _messages.add(_ChatMsg(text: txt, isAi: false, time: '방금'));
    });
    _msgCtrl.clear();

    // AI 응답 시뮬레이션 (800ms 후)
    Future.delayed(const Duration(milliseconds: 800), () {
      if (!mounted) return;
      setState(() {
        _messages.add(const _ChatMsg(
          text: '명상 모드를 가동합니다. 잠시 후 심해(DeepSea) 파동이\n'
              '안정될 것입니다. 만파식적이 파도를 잠재우듯이요.',
          isAi: true,
          time: '방금',
          msgType: _MsgType.action,
        ));
        // 학습 진행률 미세 증가
        _aiLearningProgress = (_aiLearningProgress + 0.003).clamp(0.0, 1.0);
      });
    });
  }

  void _goToMedical() => context.push('/medical/telemedicine');

  @override
  Widget build(BuildContext context) {
    final mq = MediaQuery.of(context);
    final topPad = mq.padding.top;
    final bottomPad = mq.padding.bottom;
    final screenW = mq.size.width;

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: Column(
        children: [
          // ── SlimHeader ──
          _SlimHeader(topPad: topPad),

          // ── Main content: Orb + Prediction + Chat ──
          Expanded(
            child: ClipRect(
              child: Stack(
                children: [
                  // Layer 0: Deep sea background
                  const Positioned.fill(child: _DeepSeaBackground()),

                  // Layer 1: Particle system (축적 데이터 시각화)
                  Positioned.fill(
                    child: AnimatedBuilder(
                      animation: _particleCtrl,
                      builder: (_, __) => CustomPaint(
                        painter: _DataParticlePainter(
                          progress: _particleCtrl.value,
                          particleCount: (_accumulatedDays * 0.2).round(),
                          riskLevel: _riskLevel,
                        ),
                      ),
                    ),
                  ),

                  // Layer 2: Chat list
                  Positioned.fill(
                    child: CustomScrollView(
                      controller: _scrollCtrl,
                      physics: const BouncingScrollPhysics(),
                      slivers: [
                        // Orb zone spacer
                        const SliverToBoxAdapter(
                          child: SizedBox(height: _kOrbZoneH + 56),
                        ),

                        // AI 예측 경고 배너 (Predict 스토리라인)
                        if (_riskLevel != _RiskLevel.safe)
                          SliverToBoxAdapter(
                            child: FadeTransition(
                              opacity: _predictionFade,
                              child: _PredictionBanner(
                                riskLevel: _riskLevel,
                                onMedicalTap: _goToMedical,
                              ),
                            ),
                          ),

                        // AI 학습 진행률 인디케이터 (Analyze & Learn)
                        SliverToBoxAdapter(
                          child: _AiLearningIndicator(
                            progress: _aiLearningProgress,
                            accumulatedDays: _accumulatedDays,
                          ),
                        ),

                        // Chat bubbles
                        SliverPadding(
                          padding: const EdgeInsets.symmetric(horizontal: 16),
                          sliver: SliverList.builder(
                            itemCount: _messages.length,
                            itemBuilder: (_, i) => _ChatBubbleRow(
                              msg: _messages[i],
                              maxW: screenW * _kBubbleMaxRatio,
                            ),
                          ),
                        ),

                        // Bottom spacer
                        const SliverToBoxAdapter(
                          child: SizedBox(height: 80),
                        ),
                      ],
                    ),
                  ),

                  // Layer 3: Orb zone overlay (fixed)
                  Positioned(
                    top: 0,
                    left: 0,
                    right: 0,
                    height: _kOrbZoneH + 40,
                    child: _OrbZone(
                      orbPulse: _orbPulse,
                      orbRotateCtrl: _orbRotateCtrl,
                      riskGlow: _riskGlow,
                      riskLevel: _riskLevel,
                      accumulatedDays: _accumulatedDays,
                    ),
                  ),

                  // Layer 4: 플로팅 원격진료 3D 액션 버튼
                  if (_riskLevel != _RiskLevel.safe)
                    Positioned(
                      right: 16,
                      bottom: 16,
                      child: _FloatingMedicalButton(
                        riskGlow: _riskGlow,
                        riskLevel: _riskLevel,
                        onTap: _goToMedical,
                      ),
                    ),
                ],
              ),
            ),
          ),

          // ── Bottom: Glassmorphism input panel ──
          _InputPanel(
            controller: _msgCtrl,
            bottomPad: bottomPad,
            onSend: _onSend,
          ),
        ],
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  SlimHeader — "AI SPIRIT DOCTOR" sub-page with back
// ═══════════════════════════════════════════════════════════════════

class _SlimHeader extends StatelessWidget {
  final double topPad;
  const _SlimHeader({required this.topPad});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: topPad + _kHeaderH,
      child: Stack(
        fit: StackFit.expand,
        children: [
          Image.asset(
            'assets/images/header_3d_frame_slim.png',
            fit: BoxFit.fill,
            errorBuilder: (_, __, ___) =>
                const ColoredBox(color: Color(0xFF080E1E)),
          ),
          const CustomPaint(painter: _HeaderBezelPainter()),

          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            height: _kHeaderH,
            child: Row(
              children: [
                const SizedBox(width: 4),
                IconButton(
                  icon: const Icon(Icons.arrow_back_ios_new,
                      color: Colors.white70, size: 18),
                  onPressed: () => context.pop(),
                  padding: EdgeInsets.zero,
                  constraints:
                      const BoxConstraints(minWidth: 32, minHeight: 32),
                ),
                const Expanded(
                  child: Text(
                    'AI SPIRIT DOCTOR',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      color: SanggamTheme.primary,
                      fontSize: 10,
                      letterSpacing: 5,
                      fontWeight: FontWeight.w900,
                    ),
                  ),
                ),
                const SizedBox(width: 36),
              ],
            ),
          ),

          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            height: 1,
            child: ColoredBox(
              color: SanggamTheme.primary.withValues(alpha: 0.4),
            ),
          ),
        ],
      ),
    );
  }
}

class _HeaderBezelPainter extends CustomPainter {
  const _HeaderBezelPainter();

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          Colors.black.withValues(alpha: 0.6),
          Colors.black.withValues(alpha: 0.1),
        ],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));
    canvas.drawRect(Rect.fromLTWH(0, 0, size.width, size.height), paint);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

// ═══════════════════════════════════════════════════════════════════
//  DeepSea Background — 심해 배경 (만파식 세계관)
// ═══════════════════════════════════════════════════════════════════

class _DeepSeaBackground extends StatelessWidget {
  const _DeepSeaBackground();

  @override
  Widget build(BuildContext context) {
    return const DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          stops: [0.0, 0.3, 0.7, 1.0],
          colors: [
            Color(0xFF040810),  // 심해 최상단 — 칠흑
            Color(0xFF060C1A),  // 오브 뒤 — 깊은 남색
            SanggamTheme.background,
            Color(0xFF0A1326),  // 하단 — 미세한 푸른빛
          ],
        ),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  Data Particle Painter — 축적 데이터를 떠다니는 빛 입자로 표현
//  (3. 무한한 축적 스토리라인: 앱이 나와 함께 성장)
// ═══════════════════════════════════════════════════════════════════

class _DataParticlePainter extends CustomPainter {
  final double progress;
  final int particleCount;
  final _RiskLevel riskLevel;

  _DataParticlePainter({
    required this.progress,
    required this.particleCount,
    required this.riskLevel,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final rng = math.Random(42);

    for (int i = 0; i < particleCount; i++) {
      final baseX = rng.nextDouble() * size.width;
      final baseY = rng.nextDouble() * size.height;
      final speed = 0.3 + rng.nextDouble() * 0.7;
      final phase = (progress * speed + i * 0.037) % 1.0;

      // 파티클이 천천히 위로 부유
      final x = baseX + math.sin(phase * math.pi * 2 + i) * 20;
      final y = (baseY - phase * size.height * 0.3) % size.height;

      final depth = rng.nextDouble();
      final radius = 0.5 + depth * 1.5;
      final alpha = (0.1 + depth * 0.25) * (1.0 - (phase * 0.5));

      // 위험도에 따라 파티클 색상 변화
      Color baseColor;
      switch (riskLevel) {
        case _RiskLevel.safe:
          baseColor = SanggamTheme.jagaeCyan;
        case _RiskLevel.caution:
          baseColor = Color.lerp(
              SanggamTheme.jagaeCyan, SanggamTheme.primary, depth)!;
        case _RiskLevel.danger:
          baseColor = Color.lerp(
              SanggamTheme.primary, SanggamTheme.error, depth)!;
      }

      final paint = Paint()
        ..color = baseColor.withValues(alpha: alpha)
        ..maskFilter = MaskFilter.blur(BlurStyle.normal, radius * 2);

      canvas.drawCircle(Offset(x, y), radius, paint);
    }
  }

  @override
  bool shouldRepaint(covariant _DataParticlePainter old) =>
      old.progress != progress || old.riskLevel != riskLevel;
}

// ═══════════════════════════════════════════════════════════════════
//  Orb Zone — Spirit Orb + 궤도 링 + 상태 레이블
//  (5. 적극적 예방: 위험 시 붉은 심박동, 평온 시 자개 펄스)
// ═══════════════════════════════════════════════════════════════════

class _OrbZone extends StatelessWidget {
  final Animation<double> orbPulse;
  final AnimationController orbRotateCtrl;
  final Animation<double> riskGlow;
  final _RiskLevel riskLevel;
  final int accumulatedDays;

  const _OrbZone({
    required this.orbPulse,
    required this.orbRotateCtrl,
    required this.riskGlow,
    required this.riskLevel,
    required this.accumulatedDays,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        // Fade-out mask (hides chat bubbles scrolling behind)
        const Positioned.fill(
          child: DecoratedBox(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                stops: [0.0, 0.6, 1.0],
                colors: [
                  Color(0xFF040810),
                  Color(0xDD0B1021),
                  Color(0x000B1021),
                ],
              ),
            ),
          ),
        ),

        // 오브 + 궤도 시스템
        Center(
          child: Padding(
            padding: const EdgeInsets.only(bottom: 32),
            child: AnimatedBuilder(
              animation: Listenable.merge([orbPulse, orbRotateCtrl, riskGlow]),
              builder: (_, __) => _SpiritOrb(
                pulseValue: orbPulse.value,
                rotateValue: orbRotateCtrl.value,
                riskGlowValue: riskGlow.value,
                riskLevel: riskLevel,
                accumulatedDays: accumulatedDays,
              ),
            ),
          ),
        ),

        // 상태 레이블 (Rule 1: Z-axis isolation)
        Positioned(
          left: 0,
          right: 0,
          bottom: 24,
          child: Center(
            child: SanggamContainer(
              blurSigma: 12,
              borderRadius: 16,
              borderWidth: 0.8,
              borderColor: _riskBorderColor(riskLevel),
              jagaeOpacity: 0.04,
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  _RiskDot(riskLevel: riskLevel),
                  const SizedBox(width: 8),
                  Text(
                    _statusText(riskLevel),
                    style: TextStyle(
                      color: _riskTextColor(riskLevel),
                      fontSize: 10,
                      letterSpacing: 1.5,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  static Color _riskBorderColor(_RiskLevel level) {
    switch (level) {
      case _RiskLevel.safe: return SanggamTheme.jagaeCyan;
      case _RiskLevel.caution: return SanggamTheme.primary;
      case _RiskLevel.danger: return SanggamTheme.error;
    }
  }

  static Color _riskTextColor(_RiskLevel level) {
    switch (level) {
      case _RiskLevel.safe: return SanggamTheme.jagaeCyan;
      case _RiskLevel.caution: return SanggamTheme.primary;
      case _RiskLevel.danger: return SanggamTheme.error;
    }
  }

  static String _statusText(_RiskLevel level) {
    switch (level) {
      case _RiskLevel.safe: return '안정 · 24시간 모니터링 중';
      case _RiskLevel.caution: return '주의 · 예방 조치 권고';
      case _RiskLevel.danger: return '경고 · 즉시 진료 필요';
    }
  }
}

// ── Risk dot (animated) ──
class _RiskDot extends StatelessWidget {
  final _RiskLevel riskLevel;
  const _RiskDot({required this.riskLevel});

  @override
  Widget build(BuildContext context) {
    final color = _OrbZone._riskTextColor(riskLevel);
    return Container(
      width: 6,
      height: 6,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        color: color,
        boxShadow: [
          BoxShadow(color: color.withValues(alpha: 0.6), blurRadius: 6),
        ],
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  Spirit Orb — 5단계 생명주기의 모든 데이터가 집약된 영혼 구체
//
//  위험 시: 붉은 심박동 (Breathing) + 빠른 펄스
//  주의 시: 금색 경계 글로우 + 중간 펄스
//  안정 시: 자개 시안/마젠타 명상적 빛 + 느린 펄스
//  축적 데이터량에 비례하여 궤도 입자가 풍성해짐
// ═══════════════════════════════════════════════════════════════════

class _SpiritOrb extends StatelessWidget {
  final double pulseValue;
  final double rotateValue;
  final double riskGlowValue;
  final _RiskLevel riskLevel;
  final int accumulatedDays;

  const _SpiritOrb({
    required this.pulseValue,
    required this.rotateValue,
    required this.riskGlowValue,
    required this.riskLevel,
    required this.accumulatedDays,
  });

  @override
  Widget build(BuildContext context) {
    final size = _kOrbDiameter + (pulseValue * 16);

    // 위험도에 따른 색상 매핑
    final Color innerColor;
    final Color outerColor;
    final Color glowColor;

    switch (riskLevel) {
      case _RiskLevel.safe:
        innerColor = SanggamTheme.jagaeCyan.withValues(alpha: 0.5);
        outerColor = SanggamTheme.jagaeMagenta.withValues(alpha: 0.2);
        glowColor = SanggamTheme.jagaeCyan;
      case _RiskLevel.caution:
        innerColor = SanggamTheme.primary.withValues(alpha: 0.5);
        outerColor = SanggamTheme.error.withValues(alpha: 0.15);
        glowColor = SanggamTheme.primary;
      case _RiskLevel.danger:
        // 위험 시 붉은색이 미세하게 감도는 심박동
        final redPulse = 0.3 + riskGlowValue * 0.4;
        innerColor = Color.lerp(
          SanggamTheme.primary.withValues(alpha: 0.4),
          SanggamTheme.error.withValues(alpha: 0.6),
          riskGlowValue,
        )!;
        outerColor = SanggamTheme.error.withValues(alpha: redPulse);
        glowColor = SanggamTheme.error;
    }

    return SizedBox(
      width: size + 80,
      height: size + 80,
      child: Stack(
        alignment: Alignment.center,
        children: [
          // 궤도 링 3개 (축적 데이터에 비례하여 밝아짐)
          CustomPaint(
            size: Size(size + 60, size + 60),
            painter: _OrbitRingPainter(
              rotation: rotateValue,
              orbitAlpha: (accumulatedDays / 365.0).clamp(0.15, 0.6),
              riskLevel: riskLevel,
            ),
          ),

          // 외부 글로우
          Container(
            width: size + 40,
            height: size + 40,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: glowColor.withValues(
                      alpha: 0.15 + pulseValue * 0.2),
                  blurRadius: 60,
                  spreadRadius: 20,
                ),
              ],
            ),
          ),

          // 주 오브 — 4-layer RadialGradient
          Container(
            width: size,
            height: size,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: RadialGradient(
                center: const Alignment(-0.3, -0.3),
                colors: [
                  Colors.white.withValues(alpha: 0.6),
                  innerColor,
                  outerColor,
                  Colors.transparent,
                ],
                stops: const [0.0, 0.3, 0.7, 1.0],
              ),
              boxShadow: [
                BoxShadow(
                  color: glowColor.withValues(
                      alpha: 0.25 + pulseValue * 0.15),
                  blurRadius: 30,
                  spreadRadius: 5,
                ),
              ],
            ),
          ),

          // 핵심 — 밝은 중심점
          Container(
            width: 20,
            height: 20,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: Colors.white.withValues(alpha: 0.8 + pulseValue * 0.2),
              boxShadow: [
                BoxShadow(
                  color: Colors.white.withValues(alpha: 0.6),
                  blurRadius: 20 * (0.5 + pulseValue * 0.5),
                  spreadRadius: 3,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ── Orbit ring painter ──
class _OrbitRingPainter extends CustomPainter {
  final double rotation;
  final double orbitAlpha;
  final _RiskLevel riskLevel;

  _OrbitRingPainter({
    required this.rotation,
    required this.orbitAlpha,
    required this.riskLevel,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final cx = size.width / 2;
    final cy = size.height / 2;
    final rx = size.width * 0.45;
    final ry = size.height * 0.18;

    final Color ringColor;
    switch (riskLevel) {
      case _RiskLevel.safe: ringColor = SanggamTheme.jagaeCyan;
      case _RiskLevel.caution: ringColor = SanggamTheme.primary;
      case _RiskLevel.danger: ringColor = SanggamTheme.error;
    }

    // 3개 궤도 (0°, 60°, 120°)
    for (int i = 0; i < 3; i++) {
      final tilt = i * math.pi / 3;
      final paint = Paint()
        ..color = ringColor.withValues(alpha: orbitAlpha * 0.6)
        ..style = PaintingStyle.stroke
        ..strokeWidth = 0.8;

      canvas.save();
      canvas.translate(cx, cy);
      canvas.rotate(rotation * math.pi * 2 + tilt);

      final path = Path();
      for (int j = 0; j <= 72; j++) {
        final angle = j * math.pi * 2 / 72;
        final x = rx * math.cos(angle);
        final y = ry * math.sin(angle);
        if (j == 0) {
          path.moveTo(x, y);
        } else {
          path.lineTo(x, y);
        }
      }
      path.close();
      canvas.drawPath(path, paint);

      // 궤도 위 파티클 (축적 데이터량에 비례)
      final particleCount = (orbitAlpha * 8).round();
      for (int p = 0; p < particleCount; p++) {
        final angle = (rotation * math.pi * 2 + p * math.pi * 2 / particleCount + i * 1.2) % (math.pi * 2);
        final px = rx * math.cos(angle);
        final py = ry * math.sin(angle);
        final pSize = 1.5 + (math.sin(angle) + 1) * 0.8;

        final pPaint = Paint()
          ..color = ringColor.withValues(alpha: orbitAlpha)
          ..maskFilter = MaskFilter.blur(BlurStyle.normal, pSize);
        canvas.drawCircle(Offset(px, py), pSize, pPaint);
      }

      canvas.restore();
    }
  }

  @override
  bool shouldRepaint(covariant _OrbitRingPainter old) =>
      old.rotation != rotation || old.riskLevel != riskLevel;
}

// ═══════════════════════════════════════════════════════════════════
//  Prediction Banner — AI 예측 경고
//  (4. 선제적 예측: "다가올 파도를 미리 보다")
// ═══════════════════════════════════════════════════════════════════

class _PredictionBanner extends StatelessWidget {
  final _RiskLevel riskLevel;
  final VoidCallback onMedicalTap;

  const _PredictionBanner({
    required this.riskLevel,
    required this.onMedicalTap,
  });

  @override
  Widget build(BuildContext context) {
    final isUrgent = riskLevel == _RiskLevel.danger;
    final borderColor = isUrgent ? SanggamTheme.error : SanggamTheme.primary;

    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
      child: SanggamContainer(
        borderRadius: 16,
        borderWidth: 1.2,
        borderColor: borderColor,
        blurSigma: 12,
        jagaeOpacity: 0.05,
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 상단: 아이콘 + 제목
            Row(
              children: [
                Icon(
                  isUrgent ? Icons.warning_rounded : Icons.trending_up_rounded,
                  color: borderColor,
                  size: 18,
                ),
                const SizedBox(width: 8),
                Text(
                  isUrgent ? 'AI 예측 — 즉시 조치 필요' : 'AI 예측 — 주의 권고',
                  style: TextStyle(
                    color: borderColor,
                    fontSize: 13,
                    fontWeight: FontWeight.w700,
                    letterSpacing: 0.5,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),

            // 예측 텍스트
            Text(
              isUrgent
                  ? '현재 추세라면 72시간 내 위험 수치 도달이 예상됩니다.\n즉시 원격 진료를 권장합니다.'
                  : '현재 추세라면 2주 뒤 스트레스 지수가\n위험 수치에 도달할 것으로 예측됩니다.',
              style: TextStyle(
                color: Colors.white.withValues(alpha: 0.85),
                fontSize: 13,
                height: 1.5,
              ),
            ),
            const SizedBox(height: 12),

            // 예측 트렌드 바 시각화
            _MiniTrendBar(riskLevel: riskLevel),
            const SizedBox(height: 12),

            // 원격 진료 바로가기
            GestureDetector(
              onTap: onMedicalTap,
              child: Container(
                padding: const EdgeInsets.symmetric(
                    vertical: 10, horizontal: 16),
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(12),
                  gradient: isUrgent
                      ? const LinearGradient(colors: [
                          Color(0xFF8B1A1A),
                          Color(0xFFCC3333),
                        ])
                      : SanggamTheme.goldGradient,
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.videocam_rounded,
                      color: isUrgent ? Colors.white : SanggamTheme.background,
                      size: 16,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      isUrgent ? '지금 원격 진료 시작' : '전문 상담 예약하기',
                      style: TextStyle(
                        color: isUrgent
                            ? Colors.white
                            : SanggamTheme.background,
                        fontSize: 12,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// ── Mini trend bar — 예측 트렌드 시각화 ──
class _MiniTrendBar extends StatelessWidget {
  final _RiskLevel riskLevel;
  const _MiniTrendBar({required this.riskLevel});

  @override
  Widget build(BuildContext context) {
    final dangerRatio = riskLevel == _RiskLevel.danger ? 0.85 : 0.55;

    return Column(
      children: [
        // 라벨 행
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              '현재',
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim.withValues(alpha: 0.7),
                fontSize: 9,
              ),
            ),
            Text(
              '2주 후 예측',
              style: TextStyle(
                color: SanggamTheme.primary.withValues(alpha: 0.9),
                fontSize: 9,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        const SizedBox(height: 4),

        // 그래디언트 바
        Container(
          height: 6,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(3),
            color: SanggamTheme.surfaceVariant.withValues(alpha: 0.3),
          ),
          child: Stack(
            children: [
              // 현재 수준 (실선)
              FractionallySizedBox(
                widthFactor: dangerRatio,
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(3),
                    gradient: LinearGradient(
                      colors: [
                        SanggamTheme.jagaeCyan,
                        SanggamTheme.primary,
                        SanggamTheme.error.withValues(alpha: dangerRatio),
                      ],
                    ),
                  ),
                ),
              ),
              // 위험 임계선
              Positioned(
                left: MediaQuery.of(context).size.width * 0.75 * 0.7,
                top: 0,
                bottom: 0,
                child: Container(
                  width: 1.5,
                  color: SanggamTheme.error.withValues(alpha: 0.8),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  AI Learning Indicator — 학습 진행률
//  (2. 심층 분석 및 학습: "데이터 이면의 진실을 찾다")
// ═══════════════════════════════════════════════════════════════════

class _AiLearningIndicator extends StatelessWidget {
  final double progress;
  final int accumulatedDays;

  const _AiLearningIndicator({
    required this.progress,
    required this.accumulatedDays,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
      child: SanggamContainer(
        borderRadius: 16,
        borderWidth: 0.6,
        blurSigma: 10,
        jagaeOpacity: 0.03,
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            // 원형 진행률
            SizedBox(
              width: 44,
              height: 44,
              child: CustomPaint(
                painter: _CircularProgressPainter(progress: progress),
                child: Center(
                  child: Text(
                    '${(progress * 100).round()}%',
                    style: const TextStyle(
                      color: SanggamTheme.primary,
                      fontSize: 10,
                      fontWeight: FontWeight.w800,
                    ),
                  ),
                ),
              ),
            ),
            const SizedBox(width: 16),

            // 텍스트 정보
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'AI 개인화 학습 진행률',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '$accumulatedDays일간 축적 · 1,792차원 분석 중',
                    style: TextStyle(
                      color: SanggamTheme.onSurfaceDim.withValues(alpha: 0.7),
                      fontSize: 10,
                    ),
                  ),
                ],
              ),
            ),

            // 성장 아이콘
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: SanggamTheme.primary.withValues(alpha: 0.1),
              ),
              child: const Icon(
                Icons.auto_graph_rounded,
                color: SanggamTheme.primary,
                size: 16,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _CircularProgressPainter extends CustomPainter {
  final double progress;
  _CircularProgressPainter({required this.progress});

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width / 2 - 3;

    // 배경 원
    final bgPaint = Paint()
      ..color = SanggamTheme.surfaceVariant.withValues(alpha: 0.3)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 3;
    canvas.drawCircle(center, radius, bgPaint);

    // 진행률 호
    final progressPaint = Paint()
      ..shader = SweepGradient(
        startAngle: -math.pi / 2,
        endAngle: math.pi * 2 * progress - math.pi / 2,
        colors: const [
          SanggamTheme.jagaeCyan,
          SanggamTheme.primary,
        ],
      ).createShader(Rect.fromCircle(center: center, radius: radius))
      ..style = PaintingStyle.stroke
      ..strokeWidth = 3
      ..strokeCap = StrokeCap.round;

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      -math.pi / 2,
      math.pi * 2 * progress,
      false,
      progressPaint,
    );
  }

  @override
  bool shouldRepaint(covariant _CircularProgressPainter old) =>
      old.progress != progress;
}

// ═══════════════════════════════════════════════════════════════════
//  Floating Medical Button — 즉시 원격진료 연결
//  (5. 예방: 플로팅 3D 액션 버튼을 눈앞에 대령)
// ═══════════════════════════════════════════════════════════════════

class _FloatingMedicalButton extends StatelessWidget {
  final Animation<double> riskGlow;
  final _RiskLevel riskLevel;
  final VoidCallback onTap;

  const _FloatingMedicalButton({
    required this.riskGlow,
    required this.riskLevel,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isUrgent = riskLevel == _RiskLevel.danger;
    final glowColor = isUrgent ? SanggamTheme.error : SanggamTheme.primary;

    return AnimatedBuilder(
      animation: riskGlow,
      builder: (_, __) => GestureDetector(
        onTap: onTap,
        child: Container(
          width: 56,
          height: 56,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            image: const DecorationImage(
              image: AssetImage('assets/images/btn_3d_action_slim.png'),
              fit: BoxFit.cover,
              onError: null,
            ),
            gradient: isUrgent
                ? LinearGradient(colors: [
                    SanggamTheme.error.withValues(alpha: 0.9),
                    const Color(0xFFCC3333),
                  ])
                : SanggamTheme.goldGradient,
            boxShadow: [
              BoxShadow(
                color: glowColor.withValues(
                    alpha: 0.3 + riskGlow.value * 0.3),
                blurRadius: 20 + riskGlow.value * 10,
                spreadRadius: 4,
              ),
            ],
          ),
          child: Center(
            child: Icon(
              isUrgent ? Icons.emergency_rounded : Icons.videocam_rounded,
              color: isUrgent ? Colors.white : SanggamTheme.background,
              size: 24,
            ),
          ),
        ),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  Chat Bubble Row — Rule 1: text inside SanggamContainer
// ═══════════════════════════════════════════════════════════════════

class _ChatBubbleRow extends StatelessWidget {
  final _ChatMsg msg;
  final double maxW;
  const _ChatBubbleRow({required this.msg, required this.maxW});

  @override
  Widget build(BuildContext context) {
    final isAi = msg.isAi;

    // 메시지 타입별 테두리 색상
    Color bubbleBorder;
    if (!isAi) {
      bubbleBorder = SanggamTheme.primary;
    } else {
      switch (msg.msgType) {
        case _MsgType.prediction:
          bubbleBorder = SanggamTheme.primary.withValues(alpha: 0.5);
        case _MsgType.actionPlan:
          bubbleBorder = SanggamTheme.jagaeCyan.withValues(alpha: 0.4);
        case _MsgType.action:
          bubbleBorder = SanggamTheme.jagaeMagenta.withValues(alpha: 0.3);
        default:
          bubbleBorder = Colors.white.withValues(alpha: 0.12);
      }
    }

    return Padding(
      padding: const EdgeInsets.only(bottom: _kBubbleGap),
      child: Column(
        crossAxisAlignment:
            isAi ? CrossAxisAlignment.start : CrossAxisAlignment.end,
        children: [
          Row(
            mainAxisAlignment:
                isAi ? MainAxisAlignment.start : MainAxisAlignment.end,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              if (isAi) ...[
                const _AiBadge(),
                const SizedBox(width: 8),
              ],
              ConstrainedBox(
                constraints: BoxConstraints(maxWidth: maxW),
                child: SanggamContainer(
                  borderRadius: 20,
                  borderWidth: isAi ? 0.8 : 1.2,
                  borderColor: bubbleBorder,
                  backgroundColor: isAi
                      ? SanggamTheme.surface.withValues(alpha: 0.5)
                      : SanggamTheme.surfaceVariant.withValues(alpha: 0.6),
                  blurSigma: 12,
                  jagaeOpacity: isAi ? 0.03 : 0.06,
                  padding: const EdgeInsets.symmetric(
                    vertical: _kBubbleVPad,
                    horizontal: _kBubbleHPad,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // 예측 메시지 태그
                      if (isAi && msg.msgType == _MsgType.prediction) ...[
                        const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.trending_up_rounded,
                                color: SanggamTheme.primary, size: 12),
                            SizedBox(width: 4),
                            Text(
                              'AI 예측 분석',
                              style: TextStyle(
                                color: SanggamTheme.primary,
                                fontSize: 10,
                                fontWeight: FontWeight.w700,
                                letterSpacing: 0.5,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                      ],

                      // 행동 계획 태그
                      if (isAi && msg.msgType == _MsgType.actionPlan) ...[
                        const Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.shield_rounded,
                                color: SanggamTheme.jagaeCyan, size: 12),
                            SizedBox(width: 4),
                            Text(
                              '예방 행동 플랜',
                              style: TextStyle(
                                color: SanggamTheme.jagaeCyan,
                                fontSize: 10,
                                fontWeight: FontWeight.w700,
                                letterSpacing: 0.5,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                      ],

                      Text(
                        msg.text,
                        style: TextStyle(
                          color: isAi
                              ? Colors.white
                              : Colors.white.withValues(alpha: 0.92),
                          fontSize: 16,
                          height: 1.5,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),

          const SizedBox(height: 4),
          Padding(
            padding: EdgeInsets.only(left: isAi ? 32 : 0),
            child: Text(
              msg.time,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 10,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  AI badge — golden circle with sparkle
// ═══════════════════════════════════════════════════════════════════

class _AiBadge extends StatelessWidget {
  const _AiBadge();

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 24,
      height: 24,
      decoration: const BoxDecoration(
        shape: BoxShape.circle,
        gradient: SanggamTheme.goldGradient,
      ),
      child: const Center(
        child: Icon(Icons.auto_awesome,
            size: 14, color: SanggamTheme.background),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  Input panel — glassmorphism bottom bar
// ═══════════════════════════════════════════════════════════════════

class _InputPanel extends StatelessWidget {
  final TextEditingController controller;
  final double bottomPad;
  final VoidCallback onSend;

  const _InputPanel({
    required this.controller,
    required this.bottomPad,
    required this.onSend,
  });

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 15, sigmaY: 15),
        child: Container(
          padding: EdgeInsets.fromLTRB(16, 16, 16, bottomPad + 16),
          decoration: BoxDecoration(
            color: SanggamTheme.background.withValues(alpha: 0.85),
            border: Border(
              top: BorderSide(
                color: SanggamTheme.primary.withValues(alpha: 0.2),
              ),
            ),
          ),
          child: Row(
            children: [
              // Voice button
              GestureDetector(
                onTap: () {},
                child: Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: SanggamTheme.surface.withValues(alpha: 0.5),
                    border: Border.all(
                      color: Colors.white.withValues(alpha: 0.08),
                    ),
                  ),
                  child: const Center(
                    child: Icon(Icons.mic_rounded,
                        color: SanggamTheme.onSurfaceDim, size: 18),
                  ),
                ),
              ),
              const SizedBox(width: 8),

              // Text field
              Expanded(
                child: Container(
                  height: 48,
                  decoration: BoxDecoration(
                    color: SanggamTheme.surface.withValues(alpha: 0.5),
                    borderRadius: BorderRadius.circular(24),
                    border: Border.all(
                      color: Colors.white.withValues(alpha: 0.08),
                    ),
                  ),
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  child: TextField(
                    controller: controller,
                    style: const TextStyle(color: Colors.white, fontSize: 16),
                    decoration: const InputDecoration(
                      hintText: '주치의에게 질문하기...',
                      hintStyle:
                          TextStyle(color: SanggamTheme.onSurfaceDim),
                      border: InputBorder.none,
                    ),
                    onSubmitted: (_) => onSend(),
                  ),
                ),
              ),
              const SizedBox(width: 8),

              // Send button — golden orb
              GestureDetector(
                onTap: onSend,
                child: Container(
                  width: 48,
                  height: 48,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: SanggamTheme.goldGradient,
                    boxShadow: [
                      BoxShadow(
                        color: SanggamTheme.primary.withValues(alpha: 0.4),
                        blurRadius: 16,
                        spreadRadius: 2,
                      ),
                    ],
                  ),
                  child: const Center(
                    child: Icon(Icons.bolt,
                        color: SanggamTheme.background, size: 24),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════════════════
//  Data models
// ═══════════════════════════════════════════════════════════════════

enum _MsgType { normal, prediction, actionPlan, action }

class _ChatMsg {
  final String text;
  final bool isAi;
  final String time;
  final _MsgType msgType;

  const _ChatMsg({
    required this.text,
    required this.isAi,
    required this.time,
    this.msgType = _MsgType.normal,
  });
}
