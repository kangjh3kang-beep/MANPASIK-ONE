import 'dart:math' as math;
import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitoring_providers.dart';
import 'package:manpasik/features/data_hub/presentation/providers/bio_ticker_provider.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/device_detail_bottom_sheet.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/device_status_card.dart';
import 'package:manpasik/shared/widgets/golden_orbit_globe.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/medical_stat_card.dart';
import 'package:manpasik/shared/widgets/jagae_card.dart';
import 'package:manpasik/shared/widgets/holo_body.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/orbital_hud.dart';
import 'package:manpasik/shared/widgets/glass_morphism_card.dart';
import 'package:manpasik/shared/widgets/sanggam_orbit_frame.dart';

// ───────────────────────────────────────────────────
// MonitoringDashboardScreen — Sanggam Orbit 모니터링 대시보드
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] isDark 파라미터 전체 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold ~10x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan ~5x → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) ~6x → SanggamTheme.jagaeCyan
// [Rule 4] Colors.redAccent ~4x → SanggamTheme.error
// [Rule 4] Colors.grey ~3x → SanggamTheme.onSurfaceDim
// [Rule 4] Color(0xFF0A1628) → SanggamTheme.background
// [Rule 4] Color(0xFF1A1A2E) → SanggamTheme.surface
// [Rule 4] Color(0xFF1B2640) → SanggamTheme.surfaceVariant
// [Rule 4] Color(0xFF050B14) → SanggamTheme.background
// [Rule 4] Color(0xFF080C17) → SanggamTheme.background
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] 영문 텍스트 → 한국어
// ───────────────────────────────────────────────────

class MonitoringDashboardScreen extends ConsumerStatefulWidget {
  const MonitoringDashboardScreen({super.key});

  @override
  ConsumerState<MonitoringDashboardScreen> createState() =>
      _MonitoringDashboardScreenState();
}

class _MonitoringDashboardScreenState
    extends ConsumerState<MonitoringDashboardScreen>
    with TickerProviderStateMixin {
  late AnimationController _flowController;
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _flowController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat();
    _tabController = TabController(length: 4, vsync: this);
    _tabController.addListener(() {
      final idx = _tabController.index;
      if (ref.read(monitoringFilterTabProvider) != idx) {
        ref.read(monitoringFilterTabProvider.notifier).state = idx;
      }
    });
  }

  @override
  void dispose() {
    _flowController.dispose();
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final filteredAsync = ref.watch(filteredDevicesProvider);
    final summary = ref.watch(monitoringSummaryProvider);
    final alertDevices = ref.watch(alertDevicesProvider);
    final selectedId = ref.watch(selectedDeviceIdProvider);

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: _buildAppBar(),
      body: RoyalCloudBackground(
        child: SanggamOrbitFrame(
          child: filteredAsync.when(
            data: (devices) => KeyedSubtree(
              key: const ValueKey('monitoring_data'),
              child: _buildContent(
                context, devices, summary, alertDevices, selectedId,
              ),
            ),
            loading: () => const KeyedSubtree(
              key: ValueKey('monitoring_loading'),
              child: Center(
                child: CircularProgressIndicator(color: SanggamTheme.primary),
              ),
            ),
            error: (err, _) => KeyedSubtree(
              key: const ValueKey('monitoring_error'),
              child: Center(
                child: Text('Error: $err',
                    style: const TextStyle(color: Colors.white)),
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── AppBar ──────────────────────────────────────────
  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      backgroundColor: Colors.transparent,
      elevation: 0,
      title: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Text(
            '리더기 카트리지 현황',
            style: TextStyle(
              color: SanggamTheme.primary,
              fontWeight: FontWeight.bold,
              letterSpacing: 1.0,
              fontSize: 16,
            ),
          ),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
            decoration: BoxDecoration(
              color: SanggamTheme.primary.withValues(alpha: 0.2),
              borderRadius: BorderRadius.circular(4),
            ),
            child: const Text(
              'v2.1',
              style: TextStyle(fontSize: 8, color: SanggamTheme.primary, fontWeight: FontWeight.bold),
            ),
          ),
        ],
      ),
      leading: IconButton(
        icon: const Icon(Icons.arrow_back, color: Colors.white),
        tooltip: '뒤로 가기',
        onPressed: () => context.pop(),
      ),
      actions: [
        // 프로필 토글 (바이오 탭에서만 표시)
        if (ref.watch(monitoringFilterTabProvider) == 1)
          _buildProfileToggle(),
        IconButton(
          icon: const Icon(Icons.refresh, color: Colors.white70),
          onPressed: () => ref.invalidate(pollingConnectedDevicesProvider),
          tooltip: '새로고침',
        ),
      ],
      bottom: TabBar(
        controller: _tabController,
        indicatorColor: SanggamTheme.primary,
        labelColor: SanggamTheme.primary,
        unselectedLabelColor: Colors.white54,
        tabs: const [
          Tab(text: '전체'),
          Tab(text: '바이오'),
          Tab(text: '가스'),
          Tab(text: '환경'),
        ],
      ),
    );
  }


  String _profileLabel(HoloBodyProfile profile) {
    switch (profile) {
      case HoloBodyProfile.adultMale:
        return '성인남성';
      case HoloBodyProfile.adultFemale:
        return '성인여성';
      case HoloBodyProfile.juniorMale:
        return '남자어린이';
      case HoloBodyProfile.juniorFemale:
        return '여자어린이';
      case HoloBodyProfile.baby:
        return '유아';
    }
  }

  Widget _buildProfileToggle() {
    final selectedProfile = ref.watch(holoBodyProfileProvider);

    return Padding(
      padding: const EdgeInsets.only(right: 4),
      child: PopupMenuButton<HoloBodyProfile>(
        tooltip: '홀로그램 프로필',
        initialValue: selectedProfile,
        onSelected: (HoloBodyProfile next) {
          ref.read(holoBodyProfileProvider.notifier).state = next;
          ref.read(holoGenderProvider.notifier).state = next.fallbackGender;
        },
        color: SanggamTheme.background,
        itemBuilder: (context) {
          return HoloBodyProfile.values
              .map(
                (profile) => PopupMenuItem<HoloBodyProfile>(
                  value: profile,
                  child: Text(
                    _profileLabel(profile),
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              )
              .toList();
        },
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: SanggamTheme.primary.withValues(alpha: 0.6),
            ),
            color: Colors.black.withValues(alpha: 0.3),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              const Icon(
                Icons.view_in_ar,
                size: 14,
                color: SanggamTheme.primary,
              ),
              const SizedBox(width: 4),
              Text(
                _profileLabel(selectedProfile),
                style: const TextStyle(
                  fontSize: 11,
                  color: Colors.white70,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(width: 2),
              const Icon(
                Icons.arrow_drop_down,
                size: 16,
                color: Colors.white70,
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ── Content Layout ──────────────────────────────────
  Widget _buildContent(
    BuildContext context,
    List<ConnectedDevice> devices,
    ({int total, int connected, int alerts, int avgBattery}) summary,
    List<ConnectedDevice> alertDevices,
    String? selectedId,
  ) {
    final double topPadding = MediaQuery.of(context).padding.top;
    const double appBarHeight = kToolbarHeight + 48;

    return Column(
      children: [
        SizedBox(height: topPadding + appBarHeight),
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 4),
          child: _buildSummaryBar(summary),
        ),
        Expanded(
          child: LayoutBuilder(
            builder: (context, constraints) {
              return _buildPortraitLayout(
                  constraints, devices, alertDevices, selectedId);
            },
          ),
        ),
      ],
    );
  }

  // ── Portrait Layout ─────────────────────────────────
  Widget _buildPortraitLayout(
    BoxConstraints constraints,
    List<ConnectedDevice> devices,
    List<ConnectedDevice> alertDevices,
    String? selectedId,
  ) {
    final center = Offset(
      constraints.maxWidth / 2,
      constraints.maxHeight * 0.42,
    );
    final minDim = math.min(constraints.maxWidth, constraints.maxHeight);
    final isCompact = constraints.maxWidth < 600;
    final baseRadius =
        (minDim * (isCompact ? 0.28 : 0.32)).clamp(80.0, 200.0);
    final globeSize = baseRadius * 1.3;

    final positions =
        _computeNodePositions(center, baseRadius, devices, constraints);

    return _buildVisualization(
      center, baseRadius, globeSize,
      positions, devices, selectedId, constraints,
    );
  }

  // ── Visualization (Stack of Positioned.fill layers) ──
  Widget _buildVisualization(
    Offset center,
    double baseRadius,
    double globeSize,
    List<({Offset nodePos, Offset surfacePos, double angle})> positions,
    List<ConnectedDevice> devices,
    String? selectedId,
    BoxConstraints constraints,
  ) {
    final filterTab = ref.watch(monitoringFilterTabProvider);
    final isBody = filterTab == 1;
    final connColor = SanggamTheme.primary.withValues(alpha: 0.6);
    final bodyH = constraints.maxHeight * 0.80;
    final bodyW = math.min(bodyH * 0.42, constraints.maxWidth * 0.60);

    return Stack(
      children: [
        // 1. 배경 궤도 링 (Exclude Semantics - not interactive)
        Positioned.fill(
          child: ExcludeSemantics(
            child: CustomPaint(
              painter: _BackgroundRingsPainter(
                center: center,
                radius: baseRadius,
              ),
            ),
          ),
        ),

        // Consolidated Orbital HUD for Bio Tab
        if (isBody) ...[
          Positioned.fill(
            child: OrbitalHud(
              center: center,
              baseRadius: baseRadius * 1.5,
              nodes: [
                OrbitalNodeData(
                  label: '심박수',
                  value: ref.watch(selectedBioDataProvider)['Pulse']?.toString() ?? '72',
                  unit: 'bpm',
                  icon: Icons.favorite,
                  angle: -135,
                  radiusMultiplier: 1.1,
                  bodyAttachmentOffset: const Offset(0, -100),
                ),
                OrbitalNodeData(
                  label: '산소포화도',
                  value: ref.watch(selectedBioDataProvider)['O2']?.toString() ?? '98',
                  unit: '%',
                  icon: Icons.air,
                  angle: -45,
                  radiusMultiplier: 1.1,
                  bodyAttachmentOffset: const Offset(0, -80),
                ),
                OrbitalNodeData(
                  label: '혈당',
                  value: ref.watch(selectedBioDataProvider)['Glucose']?.toString() ?? '105',
                  unit: 'mg/dL',
                  icon: Icons.bloodtype,
                  angle: 135,
                  radiusMultiplier: 1.2,
                  bodyAttachmentOffset: const Offset(0, 50),
                  isAlert: (double.tryParse(ref.watch(selectedBioDataProvider)['Glucose']?.toString() ?? '') ?? 105) > 140,
                  onTap: () {
                    context.push('/market'); // 1단계 조치: 관련 건강 식품 샵으로 안내
                  },
                ),
                OrbitalNodeData(
                  label: '뇌활동',
                  value: ref.watch(selectedBioDataProvider)['Stress']?.toString() ?? '42',
                  unit: '%',
                  icon: Icons.psychology,
                  angle: 45,
                  radiusMultiplier: 1.0,
                  bodyAttachmentOffset: const Offset(0, -180),
                ),
              ],
              centerWidget: HoloBody(
                key: ValueKey('holo_${ref.watch(holoBodyProfileProvider)}_${ref.watch(holoGenderProvider)}'),
                width: bodyW,
                height: bodyH,
                profile: ref.watch(holoBodyProfileProvider),
                gender: ref.watch(holoGenderProvider),
                bioData: ref.watch(selectedBioDataProvider),
                bioTickerSpeed: ref.watch(bioTickerSpeedProvider),
                showHud: false,
              ),
            ),
          ),
          
          // Zero-Dead-End Ecosystem: AI 코칭 패널 부상 (위험 감지 시)
          if ((double.tryParse(ref.watch(selectedBioDataProvider)['Glucose']?.toString() ?? '') ?? 105) > 140)
            Positioned(
              bottom: 24,
              right: 24,
              child: Semantics(
                label: 'AI 코칭 시작',
                button: true,
                child: MouseRegion(
                  cursor: SystemMouseCursors.click,
                  child: GestureDetector(
                    onTap: () {
                      context.push('/coach'); // Flow A: 위험 감지 → AI 코치 핸드오프
                    },
                    child: GlassmorphismCard(
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
                      baseColor: SanggamTheme.error,
                      opacity: 0.2,
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Container(
                            padding: const EdgeInsets.all(8),
                            decoration: BoxDecoration(
                              shape: BoxShape.circle,
                              color: SanggamTheme.error.withValues(alpha: 0.3),
                            ),
                            child: const Icon(Icons.smart_toy, color: Colors.white, size: 20),
                          ),
                          const SizedBox(width: 12),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: const [
                              Text(
                                '혈당 이상 감지됨',
                                style: TextStyle(
                                  color: SanggamTheme.error,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 10,
                                ),
                              ),
                              Text(
                                'AI 코치와 상담하기',
                                style: TextStyle(
                                  color: Colors.white,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 14,
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ),
        ] else ...[
          // Standard Globe & Node layer for connected devices
          Positioned(
            left: center.dx - globeSize / 2,
            top: center.dy - globeSize / 2,
            width: globeSize,
            height: globeSize,
            child: GoldenOrbitGlobe(
              size: globeSize,
            ),
          ),
          _InteractiveNodesLayer(
            positions: positions,
            devices: devices,
            selectedId: selectedId,
            center: center,
            connectorRadius: globeSize * 0.5,
            connectorColor: connColor,
            connectorAnimation: _flowController,
            viewportSize: constraints.biggest,
            onNodeTap: (device) {
              ref.read(selectedDeviceIdProvider.notifier).state = device.id;
              showDeviceDetailSheet(context, device);
            },
          ),
        ],

        // Empty State
        if (devices.isEmpty && !isBody)
          const Positioned.fill(
            child: Center(
              child: Text(
                '연결된 기기가 없습니다.',
                style: TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 14,
                ),
              ),
            ),
          ),
      ],
    );
  }

  // ── Summary Stats Bar (컴팩트 통계바) ──────────────
  Widget _buildSummaryBar(
    ({int total, int connected, int alerts, int avgBattery}) summary,
  ) {
    return JagaeCard(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      opacity: 0.1,
      child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _CompactStatTile(
                icon: Icons.devices,
                value: '${summary.total}',
                label: '전체',
                color: Colors.white70,
                tooltip: '등록된 전체 기기 수',
              ),
              _statDivider(),
              _CompactStatTile(
                icon: Icons.link,
                value: '${summary.connected}',
                label: '연결됨',
                color: SanggamTheme.jagaeCyan,
                tooltip: '현재 연결된 기기 수',
              ),
              _statDivider(),
              _CompactStatTile(
                icon: Icons.warning_amber_rounded,
                value: '${summary.alerts}',
                label: '경고',
                color: summary.alerts > 0
                    ? SanggamTheme.error
                    : SanggamTheme.onSurfaceDim,
                tooltip: '경고 상태 기기 수',
              ),
              _statDivider(),
              _CompactStatTile(
                icon: Icons.battery_std,
                value: '${summary.avgBattery}%',
                label: '배터리',
                color: summary.avgBattery < 30
                    ? SanggamTheme.error
                    : SanggamTheme.jagaeCyan,
                tooltip: '평균 배터리 잔량',
              ),
            ],
          ),
    );
  }

  Widget _statDivider() {
    return Container(
      width: 1,
      height: 20,
      color: Colors.white.withValues(alpha: 0.1),
    );
  }

  // Phase 6: Deep HCI Optimization - F-Pattern & Progressive Disclosure
  // ignore: unused_element
  Widget _buildBioDataOverlay(double bodyW, double bodyH) {
    final bioData = ref.watch(selectedBioDataProvider);
    const cardColor = SanggamTheme.jagaeCyan;

    Widget buildInteractiveCard(String key, String label, String unitText, IconData icon, {String? tip}) {
      final val = bioData[key]?.toString() ?? '--';
      final num = double.tryParse(RegExp(r'[\d.]+').firstMatch(val)?.group(0) ?? '');
      bool isAlert = false;
      if (num != null) {
        if (key == 'Stress' && num > 70) isAlert = true;
        if (key == 'O2' && num < 95) isAlert = true;
        if (key == 'Pulse' && (num < 50 || num > 100)) isAlert = true;
        if (key == 'Glucose' && (num < 70 || num > 140)) isAlert = true;
      }

      return RepaintBoundary(
        child: MouseRegion(
          cursor: SystemMouseCursors.click,
          child: GestureDetector(
            onTap: () {
               ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$label 세부 데이터 로딩 중...'), duration: const Duration(seconds: 1)));
            },
            child: Tooltip(
              message: tip ?? '',
              child: GlassmorphismCard(
                opacity: isAlert ? 0.3 : 0.15,
                useKoreanBorder: isAlert,
                baseColor: isAlert ? SanggamTheme.error : SanggamTheme.background,
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                child: MedicalStatCard(
                  label: label,
                  value: val,
                  unit: unitText,
                  icon: icon,
                  isAlert: isAlert,
                  color: cardColor,
                  tooltipMessage: null,
                ),
              ),
            ),
          ),
        ),
      );
    }

    return Align(
      alignment: Alignment.topCenter,
      child: Padding(
        padding: const EdgeInsets.only(top: kToolbarHeight + 110),
        child: SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          physics: const BouncingScrollPhysics(),
          padding: const EdgeInsets.symmetric(horizontal: 16),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              buildInteractiveCard('Pulse', '심박수', 'bpm', Icons.favorite, tip: '분당 심박수 (bpm)'),
              const SizedBox(width: 8),
              buildInteractiveCard('O2', '산소포화도', '%', Icons.air, tip: '혈중 산소 농도 (%)'),
              const SizedBox(width: 8),
              buildInteractiveCard('Glucose', '혈당', 'mg/dL', Icons.bloodtype, tip: '혈당 수치 (mg/dL)'),
              const SizedBox(width: 8),
              buildInteractiveCard('Stress', '뇌활동', '%', Icons.psychology, tip: '스트레스 지수 (%)'),
            ],
          ),
        ),
      ),
    );
  }

  // ── Alert Banner ────────────────────────────────────
  Widget _buildAlertBanner(List<ConnectedDevice> alertDevices) {
    final names = alertDevices.take(3).map((d) => d.name).join(', ');
    final extra =
        alertDevices.length > 3 ? ' 외 ${alertDevices.length - 3}대' : '';

    return Tooltip(
      message: '경고 기기 상세 보기',
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        child: GestureDetector(
      onTap: () {
        if (alertDevices.isNotEmpty) {
          ref.read(selectedDeviceIdProvider.notifier).state =
              alertDevices.first.id;
          showDeviceDetailSheet(context, alertDevices.first);
        }
      },
      child: GlassmorphismCard(
        baseColor: SanggamTheme.error,
        opacity: 0.15,
        useKoreanBorder: false,
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        child: Row(
          children: [
            const Icon(Icons.warning_amber_rounded,
                size: 14, color: SanggamTheme.error),
            const SizedBox(width: 6),
            Expanded(
              child: Text(
                '경고 기기 ${alertDevices.length}대: $names$extra',
                style: const TextStyle(
                  color: SanggamTheme.error,
                  fontSize: 11,
                  fontWeight: FontWeight.w600,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            const Icon(Icons.chevron_right,
                size: 14, color: SanggamTheme.error),
          ],
        ),
      ),
    )));
  }

  // ── Device List Panel (DraggableSheet) ──────────────
  // ignore: unused_element
  Widget _buildDeviceListPanel(
    List<ConnectedDevice> devices,
    List<ConnectedDevice> alertDevices,
    String? selectedId,
  ) {
    return DraggableScrollableSheet(
      key: const ValueKey('drg_sheet'),
      initialChildSize: 0.12,
      minChildSize: 0.12,
      maxChildSize: 0.65,
      snap: true,
      snapSizes: const [0.12, 0.35, 0.65],
      builder: (context, scrollController) {
        return ClipRRect(
          borderRadius:
              const BorderRadius.vertical(top: Radius.circular(20)),
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 16, sigmaY: 16),
            child: Container(
              decoration: BoxDecoration(
                borderRadius:
                    const BorderRadius.vertical(top: Radius.circular(20)),
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    SanggamTheme.surfaceVariant.withValues(alpha: 0.85),
                    SanggamTheme.background.withValues(alpha: 0.95),
                  ],
                ),
                border: Border(
                  top: BorderSide(
                    color: SanggamTheme.primary.withValues(alpha: 0.3),
                  ),
                ),
              ),
              child: CustomScrollView(
                controller: scrollController,
                physics: const ClampingScrollPhysics(),
                slivers: [
                  // 헤더: 드래그 핸들 + 알림 배너 + 제목
                  SliverToBoxAdapter(
                    child: Padding(
                      padding: const EdgeInsets.fromLTRB(16, 0, 16, 0),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          Center(
                            child: Container(
                              margin: const EdgeInsets.symmetric(vertical: 10),
                              width: 40,
                              height: 4,
                              decoration: BoxDecoration(
                                color: Colors.white24,
                                borderRadius: BorderRadius.circular(2),
                              ),
                            ),
                          ),
                          if (alertDevices.isNotEmpty)
                            Padding(
                              padding: const EdgeInsets.only(bottom: 8),
                              child: _buildAlertBanner(alertDevices),
                            ),
                          Text(
                            '기기 목록 (${devices.length})',
                            style: const TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              letterSpacing: 1.2,
                            ),
                          ),
                          const SizedBox(height: 12),
                        ],
                      ),
                    ),
                  ),

                  // 기기 목록
                  if (devices.isNotEmpty)
                    SliverList.builder(
                      itemCount: devices.length,
                      itemBuilder: (context, index) {
                        final device = devices[index];
                        return Padding(
                          key: ValueKey('dev_${device.id}'),
                          padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                          child: DeviceStatusCard(
                            device: device,
                            isSelected: device.id == selectedId,
                            onTap: () {
                              ref
                                  .read(selectedDeviceIdProvider.notifier)
                                  .state = device.id;
                              showDeviceDetailSheet(context, device);
                            },
                          ),
                        );
                      },
                    ),

                  // 하단 여백 확보
                  const SliverToBoxAdapter(child: SizedBox(height: 24)),

                  // 비어있을 때 안내 문구
                  if (devices.isEmpty)
                    const SliverToBoxAdapter(
                      child: Padding(
                        padding: EdgeInsets.symmetric(vertical: 32),
                        child: Center(
                          child: Text(
                            '해당 카테고리에 기기가 없습니다.',
                            style: TextStyle(
                              color: SanggamTheme.onSurfaceDim,
                              fontSize: 12,
                            ),
                          ),
                        ),
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


  // HUD panel avoidance rects
  List<Rect> _hudAvoidanceRects(BoxConstraints constraints) {
    final vw = constraints.maxWidth;
    final vh = constraints.maxHeight;
    final cX = vw / 2;
    final cY = vh * 0.42;
    final bH = vh * 0.80;
    final bW = math.min(bH * 0.42, vw * 0.60);
    final cW = (vw * 0.18).clamp(80.0, 130.0);
    final cH = (vh * 0.13).clamp(70.0, 110.0);
    final bL = cX - bW / 2;
    final bR = cX + bW / 2;
    final pG = (vw * 0.012).clamp(4.0, 12.0);
    final lX = (bL - cW - pG).clamp(0.0, bL);
    final rX = (bR + pG).clamp(bR, vw - cW);
    final vG = (vh * 0.02).clamp(8.0, 20.0);
    final tY = cY - vG - cH;
    final bY = cY + vG;
    const double np = 120.0; // 기존 60.0에서 120.0으로 파격적으로 확장하여 간섭 원천 차단
    return [
      Rect.fromLTWH(lX - np, tY - np, cW + np * 2, cH + np * 2),
      Rect.fromLTWH(cX - 210, cY - 340, 420, 680), // 중앙 홀로그램 영역 대폭 확장 (360x560 -> 420x680)
      Rect.fromLTWH(lX - np, bY - np, cW + np * 2, cH + np * 2),
      Rect.fromLTWH(rX - np, tY - np, cW + np * 2, cH + np * 2),
      Rect.fromLTWH(rX - np, bY - np, cW + np * 2, cH + np * 2),
    ];
  }

  // ── Node Position Computation (타원형 각도/간격 기반 + 안전 영역) ──
  List<({Offset nodePos, Offset surfacePos, double angle})>
      _computeNodePositions(
    Offset center,
    double baseRadius,
    List<ConnectedDevice> devices,
    BoxConstraints constraints,
  ) {
    if (devices.isEmpty) return [];

    final List<({Offset nodePos, Offset surfacePos, double angle})> positions =
        [];
    final int count = devices.length;
    final bool useTwoRings = count > 8; // 기존 10에서 8로 조정하여 혼잡도 완화
    final int innerCount = useTwoRings ? (count * 0.4).floor() : count; // 내부 링 비중 축소 (기존 50% -> 40%)

    const double nodeRadius = 23.0;
    const double padding = 12.0;
    const double globeRadius = 60.0;

    // 가용 공간 계산
    final maxRx = center.dx - nodeRadius - padding;
    final maxRyTop = center.dy - nodeRadius - padding;
    final maxRyBot = constraints.maxHeight * 0.75 - center.dy - nodeRadius;
    final maxRy = math.min(maxRyTop, maxRyBot);

    // 타원형 반경 (가용 공간 비율) - 중앙 영역 침범 방지를 위해 최소 반경 확보
    final innerRx = math.max(maxRx * 0.70, 180.0); // 0.62 -> 0.70, min 180
    final innerRy = math.max(maxRy * 0.65, 200.0); // 0.55 -> 0.65, min 200
    final outerRx = math.max(maxRx * 0.92, 260.0); // 0.88 -> 0.92
    final outerRy = math.max(maxRy * 0.88, 300.0); // 0.82 -> 0.88

    // HUD panel avoidance zones
    final hudRects = _hudAvoidanceRects(constraints);

    Offset avoidHud(double x, double y) {
      for (final rect in hudRects) {
        if (rect.contains(Offset(x, y))) {
          final adx = x - center.dx;
          final ady = y - center.dy;
          final dist = math.sqrt(adx * adx + ady * ady);
          if (dist < 1) continue;
          final ux = adx / dist;
          final uy = ady / dist;
          for (double step = 10; step <= 200; step += 10) {
            final nx = x + ux * step;
            final ny = y + uy * step;
            if (!rect.contains(Offset(nx, ny))) {
              return Offset(nx, ny);
            }
          }
        }
      }
      return Offset(x, y);
    }


    void computeRing(
        List<ConnectedDevice> ringDevices, double rx, double ry, int ringIndex) {
      final int ringCount = ringDevices.length;
      final double angleStep = (2 * math.pi) / ringCount;
      final double ringOffset = ringIndex * 0.4;

      for (int i = 0; i < ringCount; i++) {
        final device = ringDevices[i];
        final int hash = device.id.hashCode;
        // jitter를 반경의 8%로 제한
        final double angleJitter =
            ((hash & 0xFF) / 255.0 - 0.5) * 0.175;
        final double rxJitter =
            (((hash >> 8) & 0xFF) / 255.0 - 0.5) * rx * 0.08;
        final double ryJitter =
            (((hash >> 16) & 0xFF) / 255.0 - 0.5) * ry * 0.08;

        final double angle =
            -math.pi / 2 + ringOffset + (i * angleStep) + angleJitter;

        double x = center.dx + (rx + rxJitter) * math.cos(angle);
        double y = center.dy + (ry + ryJitter) * math.sin(angle);

        // Push away from HUD panels
        final avoided = avoidHud(x, y);
        x = avoided.dx;
        y = avoided.dy;

        // Keep nodes inside screen bounds.
        x = x.clamp(nodeRadius + padding, constraints.maxWidth - nodeRadius - padding);
        y = y.clamp(nodeRadius + padding, constraints.maxHeight * 0.75 - nodeRadius);

        final double dx = x - center.dx;
        final double dy = y - center.dy;
        final double dist = math.sqrt(dx * dx + dy * dy);
        final Offset surfacePoint = dist > 0
            ? Offset(center.dx + (dx / dist) * globeRadius,
                center.dy + (dy / dist) * globeRadius)
            : center;

        positions.add(
            (nodePos: Offset(x, y), surfacePos: surfacePoint, angle: angle));
      }
    }

    if (useTwoRings) {
      computeRing(devices.sublist(0, innerCount), innerRx, innerRy, 0);
      computeRing(devices.sublist(innerCount), outerRx, outerRy, 1);
    } else {
      computeRing(devices, innerRx, innerRy, 0);
    }
    return positions;
  }
}

// ═══════════════════════════════════════════════════════
// Interactive Nodes Layer — StatefulWidget + Listener + 드래그 지원
// ═══════════════════════════════════════════════════════
class _InteractiveNodesLayer extends StatefulWidget {
  final List<({Offset nodePos, Offset surfacePos, double angle})> positions;
  final List<ConnectedDevice> devices;
  final String? selectedId;
  final Offset center;
  final double connectorRadius;
  final Color connectorColor;
  final Animation<double> connectorAnimation;
  final Size viewportSize;
  final void Function(ConnectedDevice) onNodeTap;

  const _InteractiveNodesLayer({
    required this.positions,
    required this.devices,
    required this.selectedId,
    required this.center,
    required this.connectorRadius,
    required this.connectorColor,
    required this.connectorAnimation,
    required this.viewportSize,
    required this.onNodeTap,
  });

  @override
  State<_InteractiveNodesLayer> createState() => _InteractiveNodesLayerState();
}

class _InteractiveNodesLayerState extends State<_InteractiveNodesLayer> {
  int? _hoveredIndex;

  // ── Drag state ──────────────────────────────────────
  final Map<String, Offset> _userOffsets = {};
  int? _draggingIndex;
  Offset _dragStartPointer = Offset.zero;
  Offset _dragStartNode = Offset.zero;
  bool _didMove = false;
  static const _prefKeyX = 'node_x_';
  static const _prefKeyY = 'node_y_';

  String _storageKeyX(String deviceId) => '$_prefKeyX$deviceId';
  String _storageKeyY(String deviceId) => '$_prefKeyY$deviceId';

  String? _deviceIdByIndex(int index) {
    if (index < 0 || index >= widget.devices.length) return null;
    return widget.devices[index].id;
  }

  @override
  void initState() {
    super.initState();
    _loadSavedPositions();
  }

  @override
  void didUpdateWidget(covariant _InteractiveNodesLayer oldWidget) {
    super.didUpdateWidget(oldWidget);
    final oldIds = oldWidget.devices.map((d) => d.id).toSet();
    final newIds = widget.devices.map((d) => d.id).toSet();
    if (oldIds.length == newIds.length && oldIds.containsAll(newIds)) {
      return;
    }

    setState(() {
      _userOffsets.removeWhere((deviceId, _) => !newIds.contains(deviceId));
      if (_hoveredIndex != null && _hoveredIndex! >= widget.devices.length) {
        _hoveredIndex = null;
      }
      if (_draggingIndex != null && _draggingIndex! >= widget.devices.length) {
        _draggingIndex = null;
      }
    });

    _loadSavedPositions();
  }

  Future<void> _loadSavedPositions() async {
    final prefs = await SharedPreferences.getInstance();
    final loaded = <String, Offset>{};
    for (final device in widget.devices) {
      final x = prefs.getDouble(_storageKeyX(device.id));
      final y = prefs.getDouble(_storageKeyY(device.id));
      if (x != null && y != null) {
        loaded[device.id] = Offset(x, y);
      }
    }
    if (loaded.isNotEmpty && mounted) {
      setState(() => _userOffsets.addAll(loaded));
    }
  }

  Future<void> _savePosition(String deviceId, Offset pos) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setDouble(_storageKeyX(deviceId), pos.dx);
    await prefs.setDouble(_storageKeyY(deviceId), pos.dy);
  }

  Offset _clampToCanvas(Offset pos) {
    final size = widget.viewportSize;
    if (size.width <= 0 || size.height <= 0) return pos;
    const pad = 26.0;
    return Offset(
      pos.dx.clamp(pad, size.width - pad).toDouble(),
      pos.dy.clamp(pad, size.height - pad).toDouble(),
    );
  }

  Offset _effPos(int i) {
    final deviceId = _deviceIdByIndex(i);
    if (deviceId == null) return _clampToCanvas(widget.positions[i].nodePos);
    return _clampToCanvas(_userOffsets[deviceId] ?? widget.positions[i].nodePos);
  }

  List<({Offset nodePos, Offset surfacePos, double angle})>
      get _effectivePositions => List.generate(
            widget.positions.length,
            (i) {
              final p = widget.positions[i];
              final deviceId = _deviceIdByIndex(i);
              final ov = deviceId == null ? null : _userOffsets[deviceId];
              final nodePos = _clampToCanvas(ov ?? p.nodePos);
              return (nodePos: nodePos, surfacePos: p.surfacePos, angle: p.angle);
            },
          );

  int? _findNodeAt(Offset localPos) {
    for (int i = 0;
        i < widget.positions.length && i < widget.devices.length;
        i++) {
      if ((localPos - _effPos(i)).distance < 28) return i;
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    final effPos = _effectivePositions;

    final cursor = _draggingIndex != null
        ? SystemMouseCursors.grabbing
        : (_hoveredIndex != null ? SystemMouseCursors.grab : MouseCursor.defer);

    return MouseRegion(
      cursor: cursor,
      onExit: (_) {
        if (_hoveredIndex != null) setState(() => _hoveredIndex = null);
      },
      hitTestBehavior: HitTestBehavior.translucent,
      child: Listener(
        behavior: HitTestBehavior.translucent,
        // ── Hover ──
        onPointerHover: (event) {
          if (_draggingIndex != null) return;
          final found = _findNodeAt(event.localPosition);
          if (found != _hoveredIndex) setState(() => _hoveredIndex = found);
        },
        // ── Press ──
        onPointerDown: (event) {
          final found = _findNodeAt(event.localPosition);
          if (found != null) {
            setState(() {
              _draggingIndex    = found;
              _dragStartPointer = event.localPosition;
              _dragStartNode    = _effPos(found);
              _didMove          = false;
            });
          }
        },
        // ── Drag ──
        onPointerMove: (event) {
          if (_draggingIndex == null) return;
          final delta = event.localPosition - _dragStartPointer;
          if (delta.distance > 8) _didMove = true;
          if (_didMove) {
            setState(() {
              final next = _clampToCanvas(_dragStartNode + delta);

              final deviceId = _deviceIdByIndex(_draggingIndex!);
              if (deviceId != null) {
                _userOffsets[deviceId] = next;
              }
            });
          }
        },
        // ── Release ──
        onPointerUp: (event) {
          final dragIndex = _draggingIndex;
          if (dragIndex != null) {
            if (!_didMove) {
              widget.onNodeTap(widget.devices[dragIndex]);
            } else {
              final deviceId = _deviceIdByIndex(dragIndex);
              final pos = deviceId == null ? null : _userOffsets[deviceId];
              if (deviceId != null && pos != null) {
                _savePosition(deviceId, pos);
              }
            }
            setState(() => _draggingIndex = null);
          }
        },
        onPointerCancel: (_) => setState(() => _draggingIndex = null),
        child: Stack(
          children: [
            Positioned.fill(
              child: IgnorePointer(
                child: AnimatedBuilder(
                  animation: widget.connectorAnimation,
                  builder: (_, __) => CustomPaint(
                    painter: _AllConnectorsPainter(
                      connectors: _buildConnectorData(effPos),
                      color: widget.connectorColor,
                      animationValue: widget.connectorAnimation.value,
                    ),
                    size: Size.infinite,
                  ),
                ),
              ),
            ),
            Positioned.fill(
              child: CustomPaint(
                painter: _NodesPainter(
                  positions: effPos,
                  devices: widget.devices,
                  selectedId: widget.selectedId,
                  hoveredIndex: _hoveredIndex,
                  draggingIndex: _draggingIndex,
                ),
                size: Size.infinite,
              ),
            ),
            Positioned.fill(
              child: IgnorePointer(
                ignoring: _hoveredIndex == null,
                child: AnimatedOpacity(
                  opacity: (_hoveredIndex != null &&
                      _hoveredIndex! < widget.devices.length &&
                      _hoveredIndex! < effPos.length)
                      ? 1.0
                      : 0.0,
                  duration: const Duration(milliseconds: 150),
                  child: _buildTooltipLayer(effPos),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // parentDataDirty 규제: Transform.translate 기반 툴팁
  List<({Offset start, Offset end})> _buildConnectorData(
      List<({Offset nodePos, Offset surfacePos, double angle})> effPos) {
    final connectors = <({Offset start, Offset end})>[];
    final radius = widget.connectorRadius.clamp(36.0, 320.0).toDouble();
    final limit = math.min(effPos.length, widget.devices.length);

    for (int i = 0; i < limit; i++) {
      final end = effPos[i].nodePos;
      final dx = end.dx - widget.center.dx;
      final dy = end.dy - widget.center.dy;
      final dist = math.sqrt(dx * dx + dy * dy);
      final start = dist > 0
          ? Offset(
              widget.center.dx + (dx / dist) * radius,
              widget.center.dy + (dy / dist) * radius,
            )
          : widget.center;
      connectors.add((start: start, end: end));
    }

    return connectors;
  }

  Widget _buildTooltipLayer(
      List<({Offset nodePos, Offset surfacePos, double angle})> effPos) {
    final index = _hoveredIndex;
    if (index == null ||
        index >= widget.devices.length ||
        index >= effPos.length) {
      return const SizedBox.shrink();
    }

    final device = widget.devices[index];
    final pos = effPos[index].nodePos;

    const tooltipW = 180.0;
    const tooltipH = 72.0;
    final screenW = MediaQuery.of(context).size.width;
    final left = (pos.dx - tooltipW / 2).clamp(8.0, screenW - tooltipW - 8);
    final top = (pos.dy - tooltipH - 30).clamp(8.0, double.infinity);

    final isConnected = device.status == DeviceConnectionStatus.connected;
    final values = device.currentValues.entries.take(2).toList();

    return Transform.translate(
      offset: Offset(left, top),
      child: Align(
        alignment: Alignment.topLeft,
        child: IgnorePointer(
          child: ClipRRect(
            borderRadius: BorderRadius.circular(10),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 12, sigmaY: 12),
              child: Container(
                width: tooltipW,
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(10),
                  color: SanggamTheme.background.withValues(alpha: 0.85),
                  border: Border.all(
                    color: SanggamTheme.primary.withValues(alpha: 0.5),
                    width: 0.5,
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Row(
                      children: [
                        Container(
                          width: 6, height: 6,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: isConnected
                                ? SanggamTheme.jagaeCyan
                                : SanggamTheme.onSurfaceDim,
                          ),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          isConnected ? '활성' : '비활성',
                          style: TextStyle(
                            fontSize: 8,
                            fontWeight: FontWeight.bold,
                            color: isConnected
                                ? SanggamTheme.jagaeCyan
                                : SanggamTheme.onSurfaceDim,
                          ),
                        ),
                        const SizedBox(width: 6),
                        Expanded(
                          child: Text(
                            device.name,
                            style: const TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                              color: Colors.white,
                            ),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    if (values.isNotEmpty)
                      Text(
                        values.map((e) => '${e.key}: ${e.value}').join('  '),
                        style: const TextStyle(
                          fontSize: 9,
                          color: SanggamTheme.onSurfaceDim,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    const SizedBox(height: 3),
                    Row(
                      children: [
                        Icon(Icons.battery_std, size: 10,
                            color: device.batteryLevel < 20
                                ? SanggamTheme.error
                                : SanggamTheme.jagaeCyan),
                        const SizedBox(width: 2),
                        Text(
                          '${device.batteryLevel}%',
                          style: const TextStyle(
                            fontSize: 9,
                            color: SanggamTheme.onSurfaceDim,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

// ═══════════════════════════════════════════════════════
// Nodes Painter — 각 기기를 원형 노드로 그리기
// ═══════════════════════════════════════════════════════
class _NodesPainter extends CustomPainter {
  final List<({Offset nodePos, Offset surfacePos, double angle})> positions;
  final List<ConnectedDevice> devices;
  final String? selectedId;
  final int? hoveredIndex;
  final int? draggingIndex;

  _NodesPainter({
    required this.positions,
    required this.devices,
    required this.selectedId,
    this.hoveredIndex,
    this.draggingIndex,
  });

  @override
  void paint(Canvas canvas, Size size) {
    for (int i = 0; i < positions.length && i < devices.length; i++) {
      final pos = positions[i].nodePos;
      final device = devices[i];
      final isSelected = device.id == selectedId;
      final isAlert = device.status == DeviceConnectionStatus.disconnected ||
          device.batteryLevel < 20;

      final statusColor = device.status == DeviceConnectionStatus.connected
          ? SanggamTheme.jagaeCyan
          : (isAlert ? SanggamTheme.error : SanggamTheme.onSurfaceDim);

      final isHovered   = i == hoveredIndex;
      final isDragging  = i == draggingIndex;
      final double nodeRadius = isDragging ? 25.0 : (isSelected ? 23.0 : (isHovered ? 21.0 : 19.0));

      // 드래그 글로우 (golden + pulsing ring)
      if (isDragging) {
        canvas.drawCircle(
          pos,
          nodeRadius + 10,
          Paint()
            ..color = SanggamTheme.primary.withValues(alpha: 0.40)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 12),
        );
        canvas.drawCircle(
          pos, nodeRadius + 2,
          Paint()
            ..color = SanggamTheme.primary.withValues(alpha: 0.80)
            ..style = PaintingStyle.stroke
            ..strokeWidth = 1.5,
        );
      }
      // 선택 글로우
      if (isSelected) {
        canvas.drawCircle(
          pos,
          nodeRadius + 8,
          Paint()
            ..color = SanggamTheme.jagaeCyan.withValues(alpha: 0.35)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 10),
        );
      } else if (isHovered) {
        // 호버 글로우
        canvas.drawCircle(
          pos,
          nodeRadius + 6,
          Paint()
            ..color = SanggamTheme.jagaeCyan.withValues(alpha: 0.25)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8),
        );
      }

      // Background fill.
      canvas.drawCircle(
        pos,
        nodeRadius,
        Paint()
          ..color = SanggamTheme.surface.withValues(alpha: 0.92)
          ..style = PaintingStyle.fill,
      );

      // Outer stroke.
      canvas.drawCircle(
        pos,
        nodeRadius,
        Paint()
          ..color = isSelected
              ? SanggamTheme.jagaeCyan
              : (isHovered
                  ? SanggamTheme.jagaeCyan.withValues(alpha: 0.7)
                  : statusColor)
          ..style = PaintingStyle.stroke
          ..strokeWidth = isSelected ? 2.5 : (isHovered ? 2.0 : 1.5),
      );

      // Device icon.
      _paintIcon(
        canvas,
        _deviceTypeIcon(device.type),
        pos,
        isSelected ? 18.0 : 15.0,
        isSelected ? SanggamTheme.primary : Colors.white70,
      );

      // LIVE 표시 (연결 시)
      if (device.status == DeviceConnectionStatus.connected) {
        final dotPos =
            Offset(pos.dx + nodeRadius * 0.65, pos.dy - nodeRadius * 0.65);
        canvas.drawCircle(dotPos, 4, Paint()..color = SanggamTheme.jagaeCyan);
        canvas.drawCircle(
          dotPos,
          4,
          Paint()
            ..color = Colors.white
            ..style = PaintingStyle.stroke
            ..strokeWidth = 1,
        );
      }

      // 이름 라벨 (선택 시)
      if (isSelected) {
        _paintLabel(canvas, device.name,
            Offset(pos.dx, pos.dy + nodeRadius + 6));
      }
    }
  }

  void _paintIcon(
      Canvas canvas, IconData icon, Offset center, double sz, Color color) {
    final tp = TextPainter(
      text: TextSpan(
        text: String.fromCharCode(icon.codePoint),
        style: TextStyle(
          fontFamily: icon.fontFamily,
          package: icon.fontPackage,
          fontSize: sz,
          color: color,
        ),
      ),
      textDirection: TextDirection.ltr,
    );
    tp.layout();
    tp.paint(
        canvas, Offset(center.dx - tp.width / 2, center.dy - tp.height / 2));
  }

  void _paintLabel(Canvas canvas, String text, Offset pos) {
    final tp = TextPainter(
      text: TextSpan(
        text: text,
        style: const TextStyle(
          fontSize: 10,
          color: Colors.white,
          fontWeight: FontWeight.w600,
        ),
      ),
      textDirection: TextDirection.ltr,
    );
    tp.layout(maxWidth: 120);

    final bgRect = RRect.fromRectAndRadius(
      Rect.fromCenter(
        center: Offset(pos.dx, pos.dy + tp.height / 2),
        width: tp.width + 12,
        height: tp.height + 6,
      ),
      const Radius.circular(6),
    );
    canvas.drawRRect(
      bgRect,
      Paint()
        ..color = Colors.black.withValues(alpha: 0.8),
    );
    tp.paint(canvas, Offset(pos.dx - tp.width / 2, pos.dy));
  }

  IconData _deviceTypeIcon(DeviceType type) {
    switch (type) {
      case DeviceType.gasCartridge:
        return Icons.cloud_outlined;
      case DeviceType.envCartridge:
        return Icons.thermostat;
      case DeviceType.bioCartridge:
        return Icons.science;
      case DeviceType.unknown:
        return Icons.device_unknown;
    }
  }

  @override
  bool shouldRepaint(covariant _NodesPainter oldDelegate) =>
      oldDelegate.positions != positions ||
      oldDelegate.selectedId != selectedId ||
      oldDelegate.devices != devices ||
      oldDelegate.hoveredIndex != hoveredIndex ||
      oldDelegate.draggingIndex != draggingIndex;
}

// ═══════════════════════════════════════════════════════
// Background Rings Painter
// ═══════════════════════════════════════════════════════
class _BackgroundRingsPainter extends CustomPainter {
  final Offset center;
  final double radius;

  _BackgroundRingsPainter({
    required this.center,
    required this.radius,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.white.withValues(alpha: 0.05)
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

    canvas.drawCircle(center, radius * 1.5, paint);
    canvas.drawCircle(center, radius * 2.5, paint);
    canvas.drawCircle(center, radius * 4.0, paint..strokeWidth = 0.5);
  }

  @override
  bool shouldRepaint(covariant _BackgroundRingsPainter old) =>
      old.center != center || old.radius != radius;
}

// ═══════════════════════════════════════════════════════
// All Connectors Painter (단일 CustomPaint로 모든 커넥터)
// ═══════════════════════════════════════════════════════
class _AllConnectorsPainter extends CustomPainter {
  final List<({Offset start, Offset end})> connectors;
  final Color color;
  final double animationValue;

  _AllConnectorsPainter({
    required this.connectors,
    required this.color,
    required this.animationValue,
  });

  @override
  void paint(Canvas canvas, Size size) {
    for (final c in connectors) {
      _drawConnector(canvas, c.start, c.end);
    }
  }

  void _drawConnector(Canvas canvas, Offset start, Offset end) {
    final double distance = (end - start).distance;
    if (distance == 0) return;

    final Offset mid = (start + end) / 2;
    final Offset direction = end - start;
    final double curveDir = (end.dx > start.dx) ? 1.0 : -1.0;
    final Offset normal =
        Offset(-direction.dy, direction.dx * curveDir) / distance;
    final Offset control = mid + normal * (distance * 0.25);

    final Path path = Path()
      ..moveTo(start.dx, start.dy)
      ..quadraticBezierTo(control.dx, control.dy, end.dx, end.dy);

    final paint = Paint()
      ..color = color.withValues(alpha: 0.3)
      ..strokeWidth = 1.0
      ..style = PaintingStyle.stroke;

    const double dashWidth = 3;
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

    if (color.a > 0) {
      final double pos = animationValue * length;
      final tangent = metric.getTangentForOffset(pos);
      if (tangent != null) {
        canvas.drawCircle(
          tangent.position,
          4.0,
          Paint()
            ..color = color.withValues(alpha: 0.6)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 3),
        );
        canvas.drawCircle(
            tangent.position, 2.0, Paint()..color = Colors.white);
      }
    }

    canvas.drawCircle(
        start, 2.0, Paint()..color = color.withValues(alpha: 0.5));
    canvas.drawCircle(
        end, 2.0, Paint()..color = color.withValues(alpha: 0.5));
  }

  @override
  bool shouldRepaint(covariant _AllConnectorsPainter oldDelegate) =>
      oldDelegate.animationValue != animationValue ||
      oldDelegate.connectors != connectors ||
      oldDelegate.color != color;
}

// ═══════════════════════════════════════════════════════
// Compact Stat Tile (요약바 타일)
// ═══════════════════════════════════════════════════════
class _CompactStatTile extends StatelessWidget {
  final IconData icon;
  final String value;
  final String label;
  final Color color;
  final String? tooltip;

  const _CompactStatTile({
    required this.icon,
    required this.value,
    required this.label,
    required this.color,
    this.tooltip,
  });

  @override
  Widget build(BuildContext context) {
    final content = Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(icon, size: 13, color: color),
            const SizedBox(width: 4),
            Text(
              value,
              style: TextStyle(
                color: color,
                fontSize: 13,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        const SizedBox(height: 1),
        Text(
          label,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 8,
          ),
        ),
      ],
    );
    if (tooltip != null) {
      return Tooltip(message: tooltip!, child: content);
    }
    return content;
  }
}

// ═══════════════════════════════════════════════════════
// SafeHitTestWrapper — hitTest 중 viewport null check 오류 방어
// ═══════════════════════════════════════════════════════
// ignore: unused_element
class _SafeHitTestWrapper extends SingleChildRenderObjectWidget {
  const _SafeHitTestWrapper({required super.child});

  @override
  RenderObject createRenderObject(BuildContext context) => _RenderSafeHitTest();
}

class _RenderSafeHitTest extends RenderProxyBox {
  @override
  bool hitTest(BoxHitTestResult result, {required Offset position}) {
    try {
      return super.hitTest(result, position: position);
    } catch (_) {
      // DraggableScrollableSheet 리빌드 중 Viewport 리사이즈 도중 발생 가능
      return false;
    }
  }
}
