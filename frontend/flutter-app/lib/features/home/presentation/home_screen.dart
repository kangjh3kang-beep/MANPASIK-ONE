import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/providers/ecosystem_providers.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/layout/responsive_layout.dart';
import 'package:manpasik/shared/widgets/golden_orbit_globe.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'package:manpasik/features/measurement/presentation/widgets/measurement_evidence_badge.dart';
import 'package:manpasik/features/community/presentation/challenge_screen.dart';

// ════════════════════════════════════════════════════════════════════════════
//  HomeScreen — Bento-box Slim-fit Hybrid Architecture
//
//  4대 통합 최적화 규칙 적용:
//  [규칙 1] Z축 레이어 분리 — 텍스트는 반드시 SanggamContainer 내부
//  [규칙 2] 8px 그리드 + 벤토 박스 — 모든 간격 8/16/24/32px
//  [규칙 3] 슬림핏 SafeArea — header_3d_frame_slim.png + bottom_3d_dock_slim.png
//  [규칙 4] 중복 제거 + 통합 버튼 — btn_3d_action_slim.png 어포던스
//
//  [Top Layer]    8~10 %  — 브러시드 메탈 베젤 + 프로필/상태 뱃지
//  [Center Layer] Expanded (75%+) — 벤토 그리드 (히어로 + 통계 + 인사이트)
//  [Bottom Layer] 12~15 % — 자개 독 + 플로팅 측정 버튼
// ════════════════════════════════════════════════════════════════════════════

/// 상단 헤더 컨텐츠 높이 (SafeArea 제외).
const double _kHeaderRatio = 0.09;
const double _kHeaderMinH = 48.0;
const double _kHeaderMaxH = 76.0;

/// 하단 독 컨텐츠 높이 (SafeArea 제외).
const double _kDockRatio = 0.14;
const double _kDockMinH = 96.0;
const double _kDockMaxH = 140.0;

/// 황금비 — title:body 폰트 크기 비율.
const double _kGoldenRatio = 1.618;

/// 본문 기준 폰트 크기.
const double _kBodyFont = 16.0;

/// 히어로 스코어 폰트 크기.
const double _kScoreFont = 48.0;

/// 소형 캡션 폰트 크기.
const double _kCaptionFont = 10.0;

const List<({IconData icon, String label, String path})> _kPrimaryNavActions = [
  (icon: Icons.home_outlined, label: '홈', path: '/home'),
  (icon: Icons.bar_chart_outlined, label: '데이터', path: '/data'),
  (icon: Icons.hexagon_outlined, label: '측정', path: '/measure'),
  (icon: Icons.local_hospital_outlined, label: '의료', path: '/medical'),
  (icon: Icons.shopping_cart_outlined, label: '마켓', path: '/market'),
  (icon: Icons.forum_outlined, label: '커뮤니티', path: '/community'),
  (icon: Icons.family_restroom_outlined, label: '가족', path: '/family'),
  (icon: Icons.settings_outlined, label: '설정', path: '/settings'),
];

Future<void> _showAppNavigationSheet(BuildContext context) async {
  await showModalBottomSheet<void>(
    context: context,
    backgroundColor: const Color(0xFF0A1428),
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
    ),
    builder: (sheetContext) {
      return SafeArea(
        top: false,
        child: ListView(
          shrinkWrap: true,
          children: [
            const ListTile(
              title: Text(
                '전체 메뉴',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
            const Divider(height: 1, color: Color(0x33D4AF37)),
            ..._kPrimaryNavActions.map((item) {
              return ListTile(
                leading: Icon(item.icon, color: const Color(0xFFD4AF37)),
                title: Text(
                  item.label,
                  style: const TextStyle(color: Colors.white),
                ),
                onTap: () {
                  Navigator.of(sheetContext).pop();
                  context.go(item.path);
                },
              );
            }),
          ],
        ),
      );
    },
  );
}

// ── Main Screen ──────────────────────────────────────────────────────────

class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen>
    with SingleTickerProviderStateMixin {
  late final AnimationController _pulseCtrl;

  @override
  void initState() {
    super.initState();
    _pulseCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _pulseCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final mq = MediaQuery.of(context);
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: LayoutBuilder(
        builder: (context, constraints) {
          final safeTop = mq.padding.top;
          final safeBottom = mq.padding.bottom;
          final usableHeight = math.max(
            0.0,
            constraints.maxHeight - safeTop - safeBottom,
          );
          final headerBodyH =
              (usableHeight * _kHeaderRatio).clamp(_kHeaderMinH, _kHeaderMaxH);
          final dockBodyH =
              (usableHeight * _kDockRatio).clamp(_kDockMinH, _kDockMaxH);

          return Column(
            children: [
              // ═══ Top Layer (8~10 %) ═════════════════════════════════════
              _SlimHeader(
                totalHeight: safeTop + headerBodyH,
                safeTop: safeTop,
                contentHeight: headerBodyH,
              ),

              // ═══ Center Layer (Expanded ≥ 75 %) ════════════════════════
              Expanded(
                child: Stack(
                  children: [
                    // 배경: 심해 + 발광 파티클 + 오로라
                    const Positioned.fill(
                      child: RoyalCloudBackground(child: SizedBox.shrink()),
                    ),
                    // 벤토 그리드 컨텐츠
                    Positioned.fill(
                      child: _BentoGrid(pulseCtrl: _pulseCtrl),
                    ),
                  ],
                ),
              ),

              // ═══ Bottom Layer (12~15 %) ════════════════════════════════
              _SlimDock(
                totalHeight: safeBottom + dockBodyH,
                safeBottom: safeBottom,
                contentHeight: dockBodyH,
                pulse: _pulseCtrl,
                onMeasure: () => context.push('/measure'),
              ),
            ],
          );
        },
      ),
    );
  }
}

// ════════════════════════════════════════════════════════════════════════════
//  TOP LAYER — 브러시드 메탈 베젤 + 프로필/상태 뱃지
// ════════════════════════════════════════════════════════════════════════════

class _SlimHeader extends ConsumerWidget {
  final double totalHeight;
  final double safeTop;
  final double contentHeight;

  const _SlimHeader({
    required this.totalHeight,
    required this.safeTop,
    required this.contentHeight,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final displayName = authState.displayName ?? '사용자';

    return SizedBox(
      height: totalHeight,
      child: Stack(
        children: [
          const Positioned.fill(
            child: ColoredBox(color: SanggamTheme.background),
          ),
          // 3D 베젤 이미지 (배경 텍스처)
          Positioned.fill(
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    SanggamTheme.surface.withValues(alpha: 0.94),
                    SanggamTheme.background.withValues(alpha: 0.86),
                    Colors.transparent,
                  ],
                ),
              ),
              child: Opacity(
                opacity: 0.28,
                child: ColorFiltered(
                  colorFilter: const ColorFilter.mode(
                    Color(0xFF1C2D48),
                    BlendMode.multiply,
                  ),
                  child: Image.asset(
                    'assets/images/header_3d_frame_slim.png',
                    fit: BoxFit.cover,
                    alignment: Alignment.center,
                    filterQuality: FilterQuality.high,
                    errorBuilder: (_, __, ___) => const SizedBox.shrink(),
                  ),
                ),
              ),
            ),
          ),
          // 프로그래매틱 브러시드 메탈 오버레이
          const Positioned.fill(
            child: CustomPaint(painter: _HeaderBezelPainter()),
          ),
          // ── 프로필 + 상태 뱃지 행 ──
          Positioned(
            top: safeTop,
            left: 20,
            right: 20,
            bottom: 8,
            child: Row(
              children: [
                IconButton(
                  tooltip: '전체 메뉴',
                  iconSize: 18,
                  splashRadius: 18,
                  color: SanggamTheme.primary,
                  onPressed: () => _showAppNavigationSheet(context),
                  icon: const Icon(Icons.menu_rounded),
                ),
                // 프로필 아바타
                Container(
                  width: 24,
                  height: 24,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: SanggamTheme.primary.withValues(alpha: 0.8),
                      width: 1,
                    ),
                    gradient: SanggamTheme.goldGradient,
                  ),
                  child: const Center(
                    child: Icon(Icons.person,
                        size: 14, color: SanggamTheme.background),
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  '$displayName님',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.w700,
                    letterSpacing: 0.5,
                  ),
                ),
                const Spacer(),
                // 디바이스 상태 뱃지
                const _HeaderBadge(
                  icon: Icons.bluetooth_connected,
                  label: '연결',
                  color: SanggamTheme.jagaeCyan,
                ),
                const SizedBox(width: 8),
                const _HeaderBadge(
                  icon: Icons.battery_charging_full,
                  label: '95%',
                  color: SanggamTheme.primary,
                ),
              ],
            ),
          ),
          // ── 하단 금선 1px ──
          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            height: 1,
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Colors.transparent,
                    SanggamTheme.primary.withValues(alpha: 0.5),
                    SanggamTheme.primary,
                    SanggamTheme.primary.withValues(alpha: 0.5),
                    Colors.transparent,
                  ],
                  stops: const [0.0, 0.2, 0.5, 0.8, 1.0],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _HeaderBadge extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;

  const _HeaderBadge({
    required this.icon,
    required this.label,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: color.withValues(alpha: 0.3), width: 0.5),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 10, color: color),
          const SizedBox(width: 4),
          Text(
            label,
            style: TextStyle(
              color: color,
              fontSize: _kCaptionFont,
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}

// ── Header Bezel Painter ─────────────────────────────────────────────────

class _HeaderBezelPainter extends CustomPainter {
  const _HeaderBezelPainter();

  @override
  void paint(Canvas canvas, Size size) {
    if (size.isEmpty) return;
    final rect = Offset.zero & size;
    canvas.drawRect(
      rect,
      Paint()
        ..shader = const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xE0060A14),
            Color(0xE00D1525),
            Color(0xE00B1021),
          ],
        ).createShader(rect),
    );
    // 미세 브러시 텍스처 (수평선)
    final brushPaint = Paint()
      ..color = Colors.white.withValues(alpha: 0.015)
      ..strokeWidth = 0.5;
    const step = 16;
    for (int y = 0; y < size.height; y += step) {
      canvas.drawLine(
        Offset(0, y.toDouble()),
        Offset(size.width, y.toDouble()),
        brushPaint,
      );
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

// ════════════════════════════════════════════════════════════════════════════
//  CENTER LAYER — 벤토 박스 그리드 (Hero + Stats + Insight)
// ════════════════════════════════════════════════════════════════════════════

class _BentoGrid extends ConsumerWidget {
  final AnimationController pulseCtrl;

  const _BentoGrid({required this.pulseCtrl});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ResponsiveLayout(
      mobile: _buildMobileList(),
      desktop: _buildDesktopGrid(),
    );
  }

  Widget _buildMobileList() {
    return SingleChildScrollView(
      physics: const BouncingScrollPhysics(),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
      child: Column(
        children: const [
          _HeroBentoCard(),
          SizedBox(height: 16),
          _StatsBentoRow(),
          SizedBox(height: 16),
          _CrossDomainCards(),
          SizedBox(height: 16),
          _ChallengeBanner(),
          SizedBox(height: 16),
          _InsightBentoCard(),
          SizedBox(height: 16),
          _QuickActionBentoRow(),
          SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildDesktopGrid() {
    return SingleChildScrollView(
      physics: const BouncingScrollPhysics(),
      padding: const EdgeInsets.all(32),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Left Side: Bento Grid (Main Content)
          const Expanded(
            flex: 7,
            child: Wrap(
              spacing: 24,
              runSpacing: 24,
              children: [
                // Hero Card - Strategic wide placement
                SizedBox(
                  width: double.infinity,
                  child: _HeroBentoCard(),
                ),

                // Stats & Quick Actions arranged in 2-3 columns within the left area
                SizedBox(
                  width: 500,
                  child: _StatsBentoRow(),
                ),

                SizedBox(
                  width: 400,
                  child: _InsightBentoCard(),
                ),

                SizedBox(
                  width: 400,
                  child: _QuickActionBentoRow(),
                ),

                // Additional Bento components could be added here
              ],
            ),
          ),

          const SizedBox(width: 32),

          // Right Side: Persistent Detailed Analysis Panel
          Expanded(
            flex: 3,
            child: Column(
              children: [
                _buildDetailedAnalysisSidePanel(),
                const SizedBox(height: 24),
                // Could add more vertical widgets here for desktop
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailedAnalysisSidePanel() {
    return SanggamContainer(
      padding: const EdgeInsets.all(24),
      borderRadius: 24,
      jagaeOpacity: 0.1,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Icon(Icons.analytics_outlined, color: SanggamTheme.primary),
              SizedBox(width: 8),
              Expanded(
                child: Text(
                  '심층 분석 리포트',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          _buildDetailItem('코어 안정성', '94%', SanggamTheme.jagaeCyan),
          _buildDetailItem('전력 효율성', '98%', SanggamTheme.primary),
          _buildDetailItem('네트워크 동기화', '99.9%', SanggamTheme.jagaeMagenta),
          const SizedBox(height: 24),
          const Text(
            '최근 24시간 동안 시스템 안정성이 매우 높게 유지되었습니다. 장치 간 동기화 오차가 0.1ms 미만으로 측정되어 최상의 성능을 발휘하고 있습니다.',
            style: TextStyle(color: Colors.white60, fontSize: 13, height: 1.6),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailItem(String label, String value, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Flexible(
                child: Text(label,
                    style: const TextStyle(color: Colors.white38, fontSize: 12),
                    overflow: TextOverflow.ellipsis),
              ),
              Text(value,
                  style: TextStyle(color: color, fontWeight: FontWeight.bold)),
            ],
          ),
          const SizedBox(height: 4),
          LinearProgressIndicator(
            value: double.parse(value.replaceAll('%', '')) / 100,
            backgroundColor: Colors.white10,
            color: color,
            minHeight: 2,
          ),
        ],
      ),
    );
  }
}

// ── 히어로 벤토 카드 (GoldenOrbitGlobe + Health Score) ───────────────────

class _HeroBentoCard extends ConsumerWidget {
  const _HeroBentoCard();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(homeDashboardProvider);
    final score = dashboard.healthScore;
    final latest = dashboard.latestMeasurement;
    final statusText = score >= 80
        ? '조화로운 건강 리듬'
        : score >= 60
            ? '주의가 필요한 상태'
            : '건강 관리 시작하기';

    return SanggamContainer(
      padding: const EdgeInsets.all(24),
      borderRadius: 24,
      jagaeOpacity: 0.08,
      child: LayoutBuilder(
        builder: (context, constraints) {
          final maxW = constraints.maxWidth;
          final heroH = (maxW / _kGoldenRatio).clamp(180.0, 320.0);
          final globeSize = (heroH * 0.75).clamp(120.0, 240.0);

          return SizedBox(
            height: heroH,
            child: Row(
              children: [
                SizedBox(
                  width: globeSize,
                  height: globeSize,
                  child: Hero(
                    tag: 'orbit_globe',
                    child: GoldenOrbitGlobe(size: globeSize),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'HEALTH SCORE',
                        style: TextStyle(
                          color: SanggamTheme.primary.withValues(alpha: 0.6),
                          fontSize: _kCaptionFont,
                          fontWeight: FontWeight.w800,
                          letterSpacing: 3,
                        ),
                      ),
                      const SizedBox(height: 8),
                      ShaderMask(
                        shaderCallback: (bounds) =>
                            SanggamTheme.goldGradient.createShader(bounds),
                        child: Text(
                          '$score',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: _kScoreFont,
                            fontWeight: FontWeight.w900,
                            height: 1.0,
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        statusText,
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.5),
                          fontSize: _kBodyFont * 0.8,
                          letterSpacing: 0.5,
                        ),
                      ),
                      if (latest?.measuredAt != null) ...[
                        const SizedBox(height: 4),
                        Text(
                          '최근 측정: ${_formatTimeAgo(latest!.measuredAt!)}',
                          style: TextStyle(
                            color:
                                SanggamTheme.jagaeCyan.withValues(alpha: 0.6),
                            fontSize: _kCaptionFont,
                          ),
                        ),
                      ],
                      if (latest != null) ...[
                        const SizedBox(height: 8),
                        Align(
                          alignment: Alignment.centerLeft,
                          child: MeasurementEvidenceBadge(
                            evidenceStatus: latest.evidenceStatus,
                            diagnosticReady: latest.diagnosticReady,
                            evidenceGaps: latest.evidenceGaps,
                          ),
                        ),
                      ],
                      if (dashboard.measurementHistoryStale) ...[
                        const SizedBox(height: 4),
                        Text(
                          '측정 기록 동기화 확인 필요',
                          style: TextStyle(
                            color: SanggamTheme.error.withValues(alpha: 0.75),
                            fontSize: _kCaptionFont,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

String _formatTimeAgo(DateTime dt) {
  final diff = DateTime.now().difference(dt);
  if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
  if (diff.inHours < 24) return '${diff.inHours}시간 전';
  return '${diff.inDays}일 전';
}

// ── 3열 통계 벤토 카드 행 ────────────────────────────────────────────────

class _StatsBentoRow extends ConsumerWidget {
  const _StatsBentoRow();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(homeDashboardProvider);
    final latest = dashboard.latestMeasurement;
    final value = latest?.primaryValue;
    final unit = latest?.unit ?? '';
    final type = dashboard.measurementHistoryStale
        ? '기록 상태'
        : latest?.cartridgeType ?? '';

    return Row(
      children: [
        Expanded(
          child: _StatBentoCard(
            icon: Icons.science_rounded,
            label: type.isNotEmpty ? type : '측정값',
            value: dashboard.measurementHistoryStale
                ? '확인'
                : value != null
                    ? value.toStringAsFixed(1)
                    : '--',
            unit: dashboard.measurementHistoryStale ? '' : unit,
            color: dashboard.measurementHistoryStale
                ? SanggamTheme.error
                : SanggamTheme.primary,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _StatBentoCard(
            icon: Icons.notifications_active_rounded,
            label: '알림',
            value: '${dashboard.unreadNotifications}',
            unit: '',
            color: dashboard.unreadNotifications > 0
                ? SanggamTheme.jagaeCyan
                : SanggamTheme.onSurfaceDim,
          ),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: _StatBentoCard(
            icon: Icons.family_restroom_rounded,
            label: '가족',
            value: dashboard.alertMembers.isNotEmpty
                ? '${dashboard.alertMembers.length}'
                : '정상',
            unit: dashboard.alertMembers.isNotEmpty ? '주의' : '',
            color: dashboard.alertMembers.isNotEmpty
                ? SanggamTheme.error
                : SanggamTheme.jagaeCyan,
          ),
        ),
      ],
    );
  }
}

class _StatBentoCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final String unit;
  final Color color;

  const _StatBentoCard({
    required this.icon,
    required this.label,
    required this.value,
    required this.unit,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 8),
      borderRadius: 16,
      borderWidth: 1.0,
      blurSigma: 12,
      jagaeOpacity: 0.05,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 20, color: color),
          const SizedBox(height: 8), // 8px 그리드: 8
          RichText(
            textAlign: TextAlign.center,
            text: TextSpan(
              children: [
                TextSpan(
                  text: value,
                  style: TextStyle(
                    color: Colors.white.withValues(alpha: 0.9),
                    fontSize: _kBodyFont,
                    fontWeight: FontWeight.w800,
                  ),
                ),
                if (unit.isNotEmpty)
                  TextSpan(
                    text: ' $unit',
                    style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: _kCaptionFont,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim.withValues(alpha: 0.7),
              fontSize: _kCaptionFont,
              fontWeight: FontWeight.w600,
              letterSpacing: 1,
            ),
          ),
        ],
      ),
    );
  }
}

// ── AI 인사이트 벤토 카드 (풀 너비) ──────────────────────────────────────

class _InsightBentoCard extends ConsumerWidget {
  const _InsightBentoCard();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(homeDashboardProvider);
    final insightText =
        dashboard.aiInsight ?? '건강 데이터를 분석하여 맞춤 인사이트를 제공합니다. 측정을 시작해 보세요.';

    return SanggamContainer(
      padding: const EdgeInsets.all(16),
      borderRadius: 16,
      jagaeOpacity: 0.06,
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: [
                  SanggamTheme.primary.withValues(alpha: 0.3),
                  SanggamTheme.jagaeCyan.withValues(alpha: 0.15),
                ],
              ),
            ),
            child: const Center(
              child: Icon(
                Icons.auto_awesome_rounded,
                size: 20,
                color: SanggamTheme.primary,
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'AI 인사이트',
                  style: TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 13,
                    fontWeight: FontWeight.w700,
                    letterSpacing: 0.5,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  insightText,
                  style: TextStyle(
                    color: Colors.white.withValues(alpha: 0.6),
                    fontSize: _kBodyFont * 0.8,
                    height: 1.5,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ── 퀵 액션 2열 벤토 카드 ────────────────────────────────────────────────

class _QuickActionBentoRow extends StatelessWidget {
  const _QuickActionBentoRow();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Row(
          children: [
            Expanded(
              child: _QuickActionCard(
                icon: Icons.hub_outlined,
                label: '데이터 허브',
                description: '건강 데이터 통합',
                onTap: () => context.go('/data'),
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _QuickActionCard(
                icon: Icons.psychology_outlined,
                label: 'AI 코치',
                description: '맞춤형 건강 코칭',
                onTap: () => context.go('/coach'),
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _QuickActionCard(
                icon: Icons.hexagon_outlined,
                label: '측정하기',
                description: '바이오마커 측정',
                onTap: () => context.push('/measure'),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: _QuickActionCard(
                icon: Icons.local_hospital_outlined,
                label: '진료예약',
                description: '원격 의료 상담',
                onTap: () => context.go('/medical'),
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _QuickActionCard(
                icon: Icons.forum_outlined,
                label: '커뮤니티',
                description: '건강 정보 공유',
                onTap: () => context.go('/community'),
              ),
            ),
            const SizedBox(width: 8),
            Expanded(
              child: _QuickActionCard(
                icon: Icons.family_restroom_outlined,
                label: '가족관리',
                description: '가족 건강 모니터링',
                onTap: () => context.go('/family'),
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class _QuickActionCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final String description;
  final VoidCallback onTap;

  const _QuickActionCard({
    required this.icon,
    required this.label,
    required this.description,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return SanggamContainer(
      onTap: onTap,
      padding: const EdgeInsets.all(16),
      borderRadius: 16,
      borderWidth: 1.0,
      jagaeOpacity: 0.04,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 24, color: SanggamTheme.primary),
          const SizedBox(height: 8), // 8px 그리드: 8
          Text(
            label,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            description,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim.withValues(alpha: 0.7),
              fontSize: _kCaptionFont + 1, // 11px
              height: 1.4,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }
}

// ── 크로스 도메인 카드 (예약/가족 알림) ──────────────────────────────────

class _CrossDomainCards extends ConsumerWidget {
  const _CrossDomainCards();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboard = ref.watch(homeDashboardProvider);
    final reservations = dashboard.upcomingReservations;
    final alerts = dashboard.alertMembers;

    if (reservations.isEmpty && alerts.isEmpty) {
      return const SizedBox.shrink();
    }

    return Column(
      children: [
        if (reservations.isNotEmpty)
          SanggamContainer(
            onTap: () => context.go('/medical'),
            padding: const EdgeInsets.all(16),
            borderRadius: 16,
            jagaeOpacity: 0.04,
            child: Row(
              children: [
                Container(
                  width: 36,
                  height: 36,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: SanggamTheme.jagaeCyan.withValues(alpha: 0.15),
                  ),
                  child: const Icon(Icons.calendar_today,
                      size: 18, color: SanggamTheme.jagaeCyan),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '예정된 진료 ${reservations.length}건',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 13,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '${reservations.first.doctorName} · ${reservations.first.specialty}',
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.5),
                          fontSize: 11,
                        ),
                      ),
                    ],
                  ),
                ),
                const Icon(Icons.chevron_right,
                    color: SanggamTheme.onSurfaceDim, size: 20),
              ],
            ),
          ),
        if (reservations.isNotEmpty && alerts.isNotEmpty)
          const SizedBox(height: 8),
        if (alerts.isNotEmpty)
          SanggamContainer(
            onTap: () => context.go('/family'),
            padding: const EdgeInsets.all(16),
            borderRadius: 16,
            borderColor: SanggamTheme.error.withValues(alpha: 0.3),
            jagaeOpacity: 0.04,
            child: Row(
              children: [
                Container(
                  width: 36,
                  height: 36,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: SanggamTheme.error.withValues(alpha: 0.15),
                  ),
                  child: const Icon(Icons.warning_amber_rounded,
                      size: 18, color: SanggamTheme.error),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '가족 건강 알림 ${alerts.length}건',
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 13,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '${alerts.first.displayName}님의 건강 상태에 주의가 필요합니다',
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.5),
                          fontSize: 11,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ],
                  ),
                ),
                const Icon(Icons.chevron_right,
                    color: SanggamTheme.onSurfaceDim, size: 20),
              ],
            ),
          ),
      ],
    );
  }
}

// ── 커뮤니티 챌린지 배너 ────────────────────────────────────────────────────

class _ChallengeBanner extends ConsumerWidget {
  const _ChallengeBanner();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final challengesAsync = ref.watch(challengesProvider);

    return challengesAsync.when(
      data: (challenges) {
        if (challenges.isEmpty) return const SizedBox.shrink();
        final active = challenges.first;
        final title = active['title'] as String? ?? '건강 챌린지';
        final participants = active['participant_count'] as int? ?? 0;

        return SanggamContainer(
          onTap: () => context.go('/community'),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          borderRadius: 16,
          borderColor: SanggamTheme.jagaeMagenta.withValues(alpha: 0.3),
          jagaeOpacity: 0.06,
          child: Row(
            children: [
              Container(
                width: 36,
                height: 36,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: SanggamTheme.jagaeMagenta.withValues(alpha: 0.15),
                ),
                child: const Icon(Icons.emoji_events_rounded,
                    size: 18, color: SanggamTheme.jagaeMagenta),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 13,
                        fontWeight: FontWeight.w700,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '$participants명 참여 중',
                      style: TextStyle(
                        color: Colors.white.withValues(alpha: 0.5),
                        fontSize: 11,
                      ),
                    ),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(
                    horizontal: 10, vertical: 4),
                decoration: BoxDecoration(
                  color: SanggamTheme.jagaeMagenta.withValues(alpha: 0.15),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Text(
                  '참여하기',
                  style: TextStyle(
                    color: SanggamTheme.jagaeMagenta,
                    fontSize: 11,
                    fontWeight: FontWeight.w700,
                  ),
                ),
              ),
            ],
          ),
        );
      },
      loading: () => const SizedBox.shrink(),
      error: (_, __) => const SizedBox.shrink(),
    );
  }
}

// ════════════════════════════════════════════════════════════════════════════
//  BOTTOM LAYER — 자개 독 + 플로팅 측정 버튼
// ════════════════════════════════════════════════════════════════════════════

class _SlimDock extends StatelessWidget {
  final double totalHeight;
  final double safeBottom;
  final double contentHeight;
  final AnimationController pulse;
  final VoidCallback onMeasure;

  const _SlimDock({
    required this.totalHeight,
    required this.safeBottom,
    required this.contentHeight,
    required this.pulse,
    required this.onMeasure,
  });

  @override
  Widget build(BuildContext context) {
    final buttonSize = (contentHeight * 0.62).clamp(52.0, 74.0);
    final navTop = math.max(24.0, contentHeight * 0.30);
    final navBottom = safeBottom + math.max(6.0, contentHeight * 0.08);
    final navHeight = math.max(40.0, totalHeight - navTop - navBottom);

    return SizedBox(
      height: totalHeight,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          const Positioned.fill(
            child: ColoredBox(color: SanggamTheme.background),
          ),
          // 독 이미지 배경
          Positioned.fill(
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    SanggamTheme.surface.withValues(alpha: 0.84),
                    SanggamTheme.background.withValues(alpha: 0.94),
                  ],
                ),
              ),
              child: Opacity(
                opacity: 0.26,
                child: ColorFiltered(
                  colorFilter: const ColorFilter.mode(
                    Color(0xFF18273F),
                    BlendMode.multiply,
                  ),
                  child: Image.asset(
                    'assets/images/bottom_3d_dock_slim.png',
                    fit: BoxFit.cover,
                    alignment: Alignment.bottomCenter,
                    filterQuality: FilterQuality.high,
                    errorBuilder: (_, __, ___) => const SizedBox.shrink(),
                  ),
                ),
              ),
            ),
          ),
          // 프로그래매틱 자개 독 오버레이
          const Positioned.fill(
            child: CustomPaint(painter: _DockBezelPainter()),
          ),
          // ── 상단 금선 ──
          Positioned(
            left: 0,
            right: 0,
            top: 0,
            height: 1,
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Colors.transparent,
                    SanggamTheme.primary.withValues(alpha: 0.3),
                    SanggamTheme.primary.withValues(alpha: 0.6),
                    SanggamTheme.primary.withValues(alpha: 0.3),
                    Colors.transparent,
                  ],
                  stops: const [0.0, 0.2, 0.5, 0.8, 1.0],
                ),
              ),
            ),
          ),
          // ── 플로팅 측정 버튼 (독 위로 28px 돌출) ──
          Positioned(
            top: -(buttonSize * 0.42),
            left: 0,
            right: 0,
            child: Center(
              child: _ActionOrb(
                size: buttonSize,
                pulse: pulse,
                onTap: onMeasure,
              ),
            ),
          ),
          // ── 네비게이션 아이템 행 ──
          Positioned(
            left: 16,
            right: 16,
            top: navTop,
            bottom: navBottom,
            child: SizedBox(
              height: navHeight,
              child: _DockNavRow(),
            ),
          ),
        ],
      ),
    );
  }
}

// ── 독 네비게이션 행 ─────────────────────────────────────────────────────

class _DockNavRow extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        _DockNavItem(
          icon: Icons.hub_outlined,
          label: '데이터',
          onTap: () => context.go('/data'),
        ),
        _DockNavItem(
          icon: Icons.shopping_cart_outlined,
          label: '마켓',
          onTap: () => context.go('/market'),
        ),
        // 중앙 공간 (플로팅 버튼 영역)
        const SizedBox(width: 64),
        _DockNavItem(
          icon: Icons.local_hospital_outlined,
          label: '의료',
          onTap: () => context.go('/medical'),
        ),
        _DockNavItem(
          icon: Icons.grid_view_rounded,
          label: '전체',
          onTap: () => _showAppNavigationSheet(context),
        ),
      ],
    );
  }
}

class _DockNavItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _DockNavItem({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 24, color: Colors.white.withValues(alpha: 0.6)),
          const SizedBox(height: 4),
          Text(
            label,
            style: TextStyle(
              color: Colors.white.withValues(alpha: 0.4),
              fontSize: _kCaptionFont,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}

// ── Dock Bezel Painter ───────────────────────────────────────────────────

class _DockBezelPainter extends CustomPainter {
  const _DockBezelPainter();

  @override
  void paint(Canvas canvas, Size size) {
    if (size.isEmpty) return;
    final rect = Offset.zero & size;

    // 자개+심해 배경
    canvas.drawRect(
      rect,
      Paint()
        ..shader = const LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            Color(0xF00D1525),
            Color(0xF00B1021),
            Color(0xF0070C18),
          ],
        ).createShader(rect),
    );

    // 미세 자개 광택 (SweepGradient)
    canvas.drawRect(
      rect,
      Paint()
        ..shader = SweepGradient(
          center: Alignment.topCenter,
          colors: [
            SanggamTheme.jagaeCyan.withValues(alpha: 0.02),
            SanggamTheme.jagaeMagenta.withValues(alpha: 0.015),
            SanggamTheme.primary.withValues(alpha: 0.01),
            SanggamTheme.jagaeCyan.withValues(alpha: 0.02),
          ],
        ).createShader(rect),
    );
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

// ════════════════════════════════════════════════════════════════════════════
//  플로팅 액션 오브 — btn_3d_action_slim.png 통합 버튼 (규칙 4)
// ════════════════════════════════════════════════════════════════════════════

class _ActionOrb extends StatelessWidget {
  final double size;
  final AnimationController pulse;
  final VoidCallback onTap;

  const _ActionOrb({
    required this.size,
    required this.pulse,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: pulse,
      builder: (_, __) {
        final scale = 1.0 + pulse.value * 0.05;
        return GestureDetector(
          onTap: onTap,
          child: Transform.scale(
            scale: scale,
            child: Container(
              width: size,
              height: size,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                boxShadow: [
                  BoxShadow(
                    color: SanggamTheme.primary
                        .withValues(alpha: 0.35 + pulse.value * 0.15),
                    blurRadius: 18,
                    spreadRadius: 2,
                  ),
                ],
              ),
              child: Stack(
                children: [
                  // 버튼 이미지
                  Positioned.fill(
                    child: ClipOval(
                      child: Image.asset(
                        'assets/images/btn_3d_action_slim.png',
                        fit: BoxFit.cover,
                        errorBuilder: (_, __, ___) => const SizedBox.shrink(),
                      ),
                    ),
                  ),
                  // 프로그래매틱 자개 오버레이
                  Positioned.fill(
                    child: CustomPaint(
                      painter: _ActionBtnPainter(phase: pulse.value),
                    ),
                  ),
                  // 아이콘 + 텍스트
                  const Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.favorite,
                            color: SanggamTheme.primary, size: 22),
                        Text(
                          '측정',
                          style: TextStyle(
                            color: SanggamTheme.primary,
                            fontSize: 9,
                            fontWeight: FontWeight.w900,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}

// ── Action Button Painter ────────────────────────────────────────────────

class _ActionBtnPainter extends CustomPainter {
  final double phase;

  const _ActionBtnPainter({required this.phase});

  @override
  void paint(Canvas canvas, Size size) {
    final r = size.width / 2;
    final c = Offset(r, r);

    // 어두운 원 배경
    canvas.drawCircle(
      c,
      r,
      Paint()
        ..shader = const RadialGradient(
          colors: [
            Color(0xFF1A1E30),
            SanggamTheme.background,
          ],
        ).createShader(Rect.fromCircle(center: c, radius: r)),
    );

    // 금선 테두리
    canvas.drawCircle(
      c,
      r - 1,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = 1.5
        ..color = SanggamTheme.primary.withValues(alpha: 0.7 + phase * 0.3),
    );

    // 자개 빛반사 (SweepGradient 회전)
    canvas.drawCircle(
      c,
      r - 2,
      Paint()
        ..shader = SweepGradient(
          transform: GradientRotation(phase * math.pi * 2),
          colors: [
            SanggamTheme.jagaeCyan.withValues(alpha: 0.06),
            SanggamTheme.jagaeMagenta.withValues(alpha: 0.04),
            SanggamTheme.primary.withValues(alpha: 0.03),
            SanggamTheme.jagaeCyan.withValues(alpha: 0.06),
          ],
        ).createShader(Rect.fromCircle(center: c, radius: r)),
    );
  }

  @override
  bool shouldRepaint(covariant _ActionBtnPainter old) => old.phase != phase;
}
