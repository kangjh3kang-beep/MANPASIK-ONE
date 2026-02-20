import 'dart:math' as math;
import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitoring_providers.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/device_detail_bottom_sheet.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/device_status_card.dart';
import 'package:manpasik/shared/widgets/holo_globe.dart';
import 'package:manpasik/shared/widgets/holo_body.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/bio_data_panel.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/body_scan_hud.dart';
import 'package:manpasik/shared/widgets/medical_stat_card.dart';

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
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: _buildAppBar(isDark),
      // parentDataDirty 근절: KeyedSubtree로 loading→data 전환 안정화
      body: RoyalCloudBackground(
        child: filteredAsync.when(
          data: (devices) => KeyedSubtree(
            key: const ValueKey('monitoring_data'),
            child: _buildContent(
              context, devices, summary, alertDevices, selectedId, isDark,
            ),
          ),
          loading: () => const KeyedSubtree(
            key: ValueKey('monitoring_loading'),
            child: Center(
              child: CircularProgressIndicator(color: AppTheme.sanggamGold),
            ),
          ),
          error: (err, _) => KeyedSubtree(
            key: const ValueKey('monitoring_error'),
            child: Center(
              child: Text('Error: $err',
                  style: TextStyle(color: Colors.white)),
            ),
          ),
        ),
      ),
    );
  }

  // ── AppBar ────────────────────────────────────────────────────────────
  PreferredSizeWidget _buildAppBar(bool isDark) {
    return AppBar(
      backgroundColor: Colors.transparent,
      elevation: 0,
      title: Text(
        '전체 리더기 현황',
        style: TextStyle(
          color: isDark ? AppTheme.sanggamGold : const Color(0xFF1A1A1A),
          fontWeight: FontWeight.bold,
          letterSpacing: 1.0,
          fontSize: 16,
        ),
      ),
      leading: IconButton(
        icon: Icon(Icons.arrow_back,
            color: isDark ? Colors.white : Colors.black),
        onPressed: () => context.pop(),
      ),
      actions: [
        // 성별 토글 (바이오 탭에서만 표시)
        if (ref.watch(monitoringFilterTabProvider) == 3)
          _buildGenderToggle(isDark),
        IconButton(
          icon: Icon(Icons.refresh,
              color: isDark ? Colors.white70 : Colors.black54),
          onPressed: () => ref.invalidate(pollingConnectedDevicesProvider),
          tooltip: '새로고침',
        ),
      ],
      bottom: TabBar(
        controller: _tabController,
        indicatorColor: AppTheme.sanggamGold,
        labelColor: isDark ? AppTheme.sanggamGold : const Color(0xFF1A1A1A),
        unselectedLabelColor: isDark ? Colors.white54 : Colors.grey,
        tabs: const [
          Tab(text: '전체'),
          Tab(text: '기체'),
          Tab(text: '환경'),
          Tab(text: '바이오'),
        ],
      ),
    );
  }

  // ── Gender Toggle Button ────────────────────────────────────────────────
  Widget _buildGenderToggle(bool isDark) {
    final gender = ref.watch(holoGenderProvider);
    final isMale = gender == HoloGender.male;
    return Padding(
      padding: const EdgeInsets.only(right: 4),
      child: GestureDetector(
        onTap: () {
          ref.read(holoGenderProvider.notifier).state =
              isMale ? HoloGender.female : HoloGender.male;
        },
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: AppTheme.sanggamGold.withValues(alpha: 0.6),
            ),
            color: isDark
                ? Colors.black.withValues(alpha: 0.3)
                : Colors.white.withValues(alpha: 0.3),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                isMale ? '\u2642' : '\u2640',
                style: TextStyle(
                  fontSize: 14,
                  color: AppTheme.sanggamGold,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(width: 4),
              Text(
                isMale ? '남성' : '여성',
                style: TextStyle(
                  fontSize: 11,
                  color: isDark ? Colors.white70 : Colors.black87,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ── Content Layout ────────────────────────────────────────────────────
  Widget _buildContent(
    BuildContext context,
    List<ConnectedDevice> devices,
    ({int total, int connected, int alerts, int avgBattery}) summary,
    List<ConnectedDevice> alertDevices,
    String? selectedId,
    bool isDark,
  ) {
    final double topPadding = MediaQuery.of(context).padding.top;
    const double appBarHeight = kToolbarHeight + 48;

    return Column(
      children: [
        SizedBox(height: topPadding + appBarHeight),
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 8, 16, 4),
          child: _buildSummaryBar(summary, isDark),
        ),
        Expanded(
          child: LayoutBuilder(
            builder: (context, constraints) {
              return _buildPortraitLayout(
                  constraints, devices, alertDevices, isDark, selectedId);
            },
          ),
        ),
      ],
    );
  }

  // ── Portrait Layout ───────────────────────────────────────────────────
  Widget _buildPortraitLayout(
    BoxConstraints constraints,
    List<ConnectedDevice> devices,
    List<ConnectedDevice> alertDevices,
    bool isDark,
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

    final positions = _computeNodePositions(center, baseRadius, devices, constraints);
    final connectorData =
        positions.map((p) => (start: p.surfacePos, end: p.nodePos)).toList();

    return _buildVisualization(
      center, baseRadius, globeSize,
      positions, connectorData, devices, isDark, selectedId, constraints,
    );
  }

  // ── Visualization (Stack of Positioned.fill layers) ────────────────────
  Widget _buildVisualization(
    Offset center,
    double baseRadius,
    double globeSize,
    List<({Offset nodePos, Offset surfacePos, double angle})> positions,
    List<({Offset start, Offset end})> connectorData,
    List<ConnectedDevice> devices,
    bool isDark,
    String? selectedId,
    BoxConstraints constraints,
  ) {
    final filterTab = ref.watch(monitoringFilterTabProvider);
    final isBody = filterTab == 3;
    final connColor =
        (isDark ? AppTheme.sanggamGold : Colors.teal).withValues(alpha: 0.6);
    // v21: 내부 fitScale이 전신 표시용으로 축소 → 외부 컨테이너를 넉넉히
    final bodyH = constraints.maxHeight * 0.80;
    final bodyW = math.min(bodyH * 0.42, constraints.maxWidth * 0.60);

    return Stack(
      children: [
        // 1. 배경 궤도 링
        Positioned.fill(
          child: CustomPaint(
            painter: _BackgroundRingsPainter(
              center: center,
              isDark: isDark,
              radius: baseRadius,
            ),
          ),
        ),

        // 2. 커넥터 (애니메이션)
        Positioned.fill(
          child: AnimatedBuilder(
            animation: _flowController,
            builder: (context, child) {
              return CustomPaint(
                painter: _AllConnectorsPainter(
                  connectors: connectorData,
                  color: connColor,
                  animationValue: _flowController.value,
                ),
              );
            },
          ),
        ),

        // 3. 중앙 홀로그램 — Globe (AnimatedOpacity)
        Positioned(
          left: center.dx - globeSize / 2,
          top: center.dy - globeSize / 2,
          width: globeSize,
          height: globeSize,
          child: IgnorePointer(
            child: AnimatedOpacity(
              opacity: isBody ? 0.0 : 1.0,
              duration: const Duration(milliseconds: 500),
              child: HoloGlobe(
                key: const ValueKey('globe'),
                size: globeSize,
                color: isDark ? AppTheme.sanggamGold : const Color(0xFF004D40),
                accentColor: isDark ? Colors.white : const Color(0xFF00796B),
              ),
            ),
          ),
        ),

        // 4. 중앙 홀로그램 — Body (v6.0 전체 화면 확대)
        Positioned(
          left: center.dx - bodyW / 2,
          top: center.dy - bodyH * 0.48,
          width: bodyW,
          height: bodyH,
          child: IgnorePointer(
            child: AnimatedOpacity(
              opacity: isBody ? 1.0 : 0.0,
              duration: const Duration(milliseconds: 500),
              child: HoloBody(
                key: ValueKey('body_${ref.watch(holoGenderProvider).name}'),
                width: bodyW,
                height: bodyH,
                color: isDark ? AppTheme.waveCyan : const Color(0xFF00ACC1),
                accentColor: isDark ? AppTheme.sanggamGold : const Color(0xFFFF4D4D),
                gender: ref.watch(holoGenderProvider),
                bioData: ref.watch(selectedBioDataProvider),
                showEcg: true,
                showHud: true,
              ),
            ),
          ),
        ),

        // 5. 인터랙티브 노드 레이어
        Positioned.fill(
          child: _InteractiveNodesLayer(
            positions: positions,
            devices: devices,
            selectedId: selectedId,
            isDark: isDark,
            onNodeTap: (device) {
              ref.read(selectedDeviceIdProvider.notifier).state = device.id;
              showDeviceDetailSheet(context, device);
            },
          ),
        ),

        // 7. v6.0 바이오 데이터 패널 오버레이 (바이오탭만)
        Positioned.fill(
          child: IgnorePointer(
            child: AnimatedOpacity(
              opacity: isBody ? 1.0 : 0.0,
              duration: const Duration(milliseconds: 400),
              child: _buildBioDataOverlay(isDark, bodyW, bodyH),
            ),
          ),
        ),

        // 8. v6.0 상단 HUD 오버레이 (바이오탭만)
        Positioned.fill(
          child: IgnorePointer(
            child: AnimatedOpacity(
              opacity: isBody ? 1.0 : 0.0,
              duration: const Duration(milliseconds: 400),
              child: Align(
                alignment: Alignment.topCenter,
                child: AnimatedBuilder(
                  animation: _flowController,
                  builder: (context, _) => BodyScanHud(
                    scanProgress: _flowController.value,
                  ),
                ),
              ),
            ),
          ),
        ),

        // 9. 빈 상태 (항상 배치 + opacity 제어 → parentDataDirty 근절)
        Positioned.fill(
          child: IgnorePointer(
            ignoring: devices.isNotEmpty,
            child: AnimatedOpacity(
              opacity: devices.isEmpty ? 1.0 : 0.0,
              duration: const Duration(milliseconds: 300),
              child: Center(
                child: Text(
                  '연결된 기기가 없습니다.',
                  style: TextStyle(
                    color: isDark ? Colors.white54 : Colors.black38,
                    fontSize: 14,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  // ── Summary Stats Bar (컴팩트) ─────────────────────────────────────────
  Widget _buildSummaryBar(
    ({int total, int connected, int alerts, int avgBattery}) summary,
    bool isDark,
  ) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(12),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: 12, sigmaY: 12),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            color: isDark
                ? Colors.black.withValues(alpha: 0.35)
                : Colors.white.withValues(alpha: 0.65),
            border: Border.all(
              color: isDark
                  ? AppTheme.sanggamGold.withValues(alpha: 0.2)
                  : Colors.black.withValues(alpha: 0.06),
            ),
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _CompactStatTile(
                icon: Icons.devices,
                value: '${summary.total}',
                label: '전체',
                color: isDark ? Colors.white70 : Colors.black87,
              ),
              _statDivider(isDark),
              _CompactStatTile(
                icon: Icons.link,
                value: '${summary.connected}',
                label: '연결',
                color: const Color(0xFF00E676),
              ),
              _statDivider(isDark),
              _CompactStatTile(
                icon: Icons.warning_amber_rounded,
                value: '${summary.alerts}',
                label: '경고',
                color: summary.alerts > 0 ? Colors.redAccent : Colors.grey,
              ),
              _statDivider(isDark),
              _CompactStatTile(
                icon: Icons.battery_std,
                value: '${summary.avgBattery}%',
                label: '배터리',
                color: summary.avgBattery < 30
                    ? Colors.redAccent
                    : const Color(0xFF00E676),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _statDivider(bool isDark) {
    return Container(
      width: 1,
      height: 20,
      color: isDark
          ? Colors.white.withValues(alpha: 0.1)
          : Colors.black.withValues(alpha: 0.08),
    );
  }

  // ── v6.0 Bio Data Overlay (Medical HUD with Stat Cards) ────────────────
  Widget _buildBioDataOverlay(bool isDark, double bodyW, double bodyH) {
    final bioData = ref.watch(selectedBioDataProvider);
    final unit = bodyW * 0.35;
    
    Widget buildCard(String key, String label, String unitText, IconData icon, Alignment alignment, Offset offset) {
       final val = bioData[key]?.toString() ?? '--';
       final num = double.tryParse(RegExp(r'[\d.]+').firstMatch(val)?.group(0) ?? '');
       bool isAlert = false;
       if (num != null) {
          if (key == 'Stress' && num > 70) isAlert = true;
          if (key == 'O2' && num < 95) isAlert = true;
          if (key == 'Pulse' && (num < 50 || num > 100)) isAlert = true;
          if (key == 'Glucose' && (num < 70 || num > 140)) isAlert = true;
       }
       
       return Align(
         alignment: alignment,
         child: Transform.translate(
           offset: offset,
           child: MedicalStatCard(
             label: label,
             value: val,
             unit: unitText,
             icon: icon,
             isAlert: isAlert,
             color: isDark ? AppTheme.waveCyan : const Color(0xFF00ACC1),
           ),
         ),
       );
    }
    
    return Stack(
      children: [
        buildCard('Stress', 'BRAIN ACTIVITY', '%', Icons.psychology, Alignment.center, Offset(-unit * 1.5 - 50, -unit * 0.6)),
        buildCard('O2', 'SpO\u2082 LEVEL', '%', Icons.air, Alignment.center, Offset(-unit * 1.5 - 50, unit * 0.5)),
        buildCard('Pulse', 'HEART RATE', 'bpm', Icons.favorite, Alignment.center, Offset(unit * 1.5 + 50, -unit * 0.2)),
        buildCard('Glucose', 'GLUCOSE', 'mg/dL', Icons.bloodtype, Alignment.center, Offset(unit * 1.5 + 50, unit * 0.9)),
      ],
    );
  }

  // ── Alert Banner ──────────────────────────────────────────────────────
  Widget _buildAlertBanner(List<ConnectedDevice> alertDevices, bool isDark) {
    final names = alertDevices.take(3).map((d) => d.name).join(', ');
    final extra =
        alertDevices.length > 3 ? ' 외 ${alertDevices.length - 3}개' : '';

    return GestureDetector(
      onTap: () {
        if (alertDevices.isNotEmpty) {
          ref.read(selectedDeviceIdProvider.notifier).state =
              alertDevices.first.id;
          showDeviceDetailSheet(context, alertDevices.first);
        }
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
        decoration: BoxDecoration(
          color: Colors.redAccent.withValues(alpha: isDark ? 0.15 : 0.08),
          borderRadius: BorderRadius.circular(8),
          border:
              Border.all(color: Colors.redAccent.withValues(alpha: 0.35)),
        ),
        child: Row(
          children: [
            const Icon(Icons.warning_amber_rounded,
                size: 14, color: Colors.redAccent),
            const SizedBox(width: 6),
            Expanded(
              child: Text(
                '${alertDevices.length}개 기기 주의: $names$extra',
                style: const TextStyle(
                  color: Colors.redAccent,
                  fontSize: 11,
                  fontWeight: FontWeight.w600,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            ),
            const Icon(Icons.chevron_right,
                size: 14, color: Colors.redAccent),
          ],
        ),
      ),
    );
  }

  // ── Device List Panel (DraggableSheet) ─────────────────────────────────
  Widget _buildDeviceListPanel(
    List<ConnectedDevice> devices,
    List<ConnectedDevice> alertDevices,
    bool isDark,
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
                  colors: isDark
                      ? [
                          const Color(0xFF1B2640).withValues(alpha: 0.85),
                          const Color(0xFF050B14).withValues(alpha: 0.95),
                        ]
                      : [
                          Colors.white.withValues(alpha: 0.95),
                          Colors.white.withValues(alpha: 0.85),
                        ],
                ),
                border: Border(
                  top: BorderSide(
                    color: isDark
                        ? AppTheme.sanggamGold.withValues(alpha: 0.3)
                        : Colors.black.withValues(alpha: 0.08),
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
                                color: isDark ? Colors.white24 : Colors.black12,
                                borderRadius: BorderRadius.circular(2),
                              ),
                            ),
                          ),
                          if (alertDevices.isNotEmpty)
                            Padding(
                              padding: const EdgeInsets.only(bottom: 8),
                              child: _buildAlertBanner(alertDevices, isDark),
                            ),
                          Text(
                            '기기 목록 (${devices.length}개)',
                            style: TextStyle(
                              color: isDark ? Colors.white54 : Colors.black45,
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
                  
                  // 기기 목록 (SliverPadding 제거 -> Item Padding으로 대체)
                  if (devices.isNotEmpty)
                    SliverList.builder(
                      itemCount: devices.length,
                      itemBuilder: (context, index) {
                        final device = devices[index];
                        return Padding(
                          key: ValueKey('dev_${device.id}'),
                          // 상하좌우 패딩을 여기서 직접 적용 (SliverPadding 회피)
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

                  // 하단 여백 확보용 (SliverPadding 대체)
                  const SliverToBoxAdapter(child: SizedBox(height: 24)),

                  // 비어있을 때 안내 문구
                  if (devices.isEmpty)
                    SliverToBoxAdapter(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(vertical: 32),
                        child: Center(
                          child: Text(
                            '이 카테고리에 등록된 기기가 없습니다.',
                            style: TextStyle(
                              color: isDark ? Colors.white38 : Colors.black38,
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

  // ── Node Position Computation (v5.0 — 타원형 가용공간 기반 + 안전 클램프) ───
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
    final bool useTwoRings = count > 10;
    final int innerCount = useTwoRings ? (count / 2).ceil() : count;

    const double nodeRadius = 23.0;
    const double padding = 12.0;
    const double globeRadius = 60.0;

    // 가용 공간 계산
    final maxRx = center.dx - nodeRadius - padding;
    final maxRyTop = center.dy - nodeRadius - padding;
    final maxRyBot = constraints.maxHeight * 0.75 - center.dy - nodeRadius;
    final maxRy = math.min(maxRyTop, maxRyBot);

    // 타원형 반경 (가용 공간 비율)
    final innerRx = maxRx * 0.62;
    final innerRy = maxRy * 0.55;
    final outerRx = maxRx * 0.88;
    final outerRy = maxRy * 0.82;

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

        // 안전 클램프
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

// ═══════════════════════════════════════════════════════════════════════════
// Interactive Nodes Layer — StatefulWidget + Listener + 호버 툴팁
// ═══════════════════════════════════════════════════════════════════════════

class _InteractiveNodesLayer extends StatefulWidget {
  final List<({Offset nodePos, Offset surfacePos, double angle})> positions;
  final List<ConnectedDevice> devices;
  final String? selectedId;
  final bool isDark;
  final void Function(ConnectedDevice) onNodeTap;

  const _InteractiveNodesLayer({
    required this.positions,
    required this.devices,
    required this.selectedId,
    required this.isDark,
    required this.onNodeTap,
  });

  @override
  State<_InteractiveNodesLayer> createState() => _InteractiveNodesLayerState();
}

class _InteractiveNodesLayerState extends State<_InteractiveNodesLayer> {
  int? _hoveredIndex;

  int? _findNodeAt(Offset localPos) {
    for (int i = 0; i < widget.positions.length && i < widget.devices.length; i++) {
      if ((localPos - widget.positions[i].nodePos).distance < 28) return i;
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onExit: (_) {
        if (_hoveredIndex != null) setState(() => _hoveredIndex = null);
      },
      hitTestBehavior: HitTestBehavior.translucent,
      child: Listener(
        behavior: HitTestBehavior.translucent,
        onPointerHover: (event) {
          final found = _findNodeAt(event.localPosition);
          if (found != _hoveredIndex) setState(() => _hoveredIndex = found);
        },
        onPointerUp: (event) {
          final found = _findNodeAt(event.localPosition);
          if (found != null) widget.onNodeTap(widget.devices[found]);
        },
        child: Stack(
        children: [
          Positioned.fill(
            child: CustomPaint(
              painter: _NodesPainter(
                positions: widget.positions,
                devices: widget.devices,
                selectedId: widget.selectedId,
                isDark: widget.isDark,
                hoveredIndex: _hoveredIndex,
              ),
              size: Size.infinite,
            ),
          ),
          // parentDataDirty 근절: 항상 2개 자식 고정 (조건부 제거)
          Positioned.fill(
            child: IgnorePointer(
              ignoring: _hoveredIndex == null,
              child: AnimatedOpacity(
                opacity: (_hoveredIndex != null &&
                    _hoveredIndex! < widget.devices.length &&
                    _hoveredIndex! < widget.positions.length)
                    ? 1.0
                    : 0.0,
                duration: const Duration(milliseconds: 150),
                child: _buildTooltipLayer(),
              ),
            ),
          ),
        ],
      ),
      ),
    );
  }

  // parentDataDirty 근절: Transform.translate 기반 툴팁 (Positioned 미사용)
  Widget _buildTooltipLayer() {
    final index = _hoveredIndex;
    if (index == null ||
        index >= widget.devices.length ||
        index >= widget.positions.length) {
      return const SizedBox.shrink();
    }

    final device = widget.devices[index];
    final pos = widget.positions[index].nodePos;
    final isDark = widget.isDark;

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
                  color: isDark
                      ? const Color(0xFF0A1628).withValues(alpha: 0.85)
                      : Colors.white.withValues(alpha: 0.90),
                  border: Border.all(
                    color: AppTheme.sanggamGold.withValues(alpha: 0.5),
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
                            color: isConnected ? const Color(0xFF00E676) : Colors.grey,
                          ),
                        ),
                        const SizedBox(width: 4),
                        Text(
                          isConnected ? 'LIVE' : 'OFF',
                          style: TextStyle(
                            fontSize: 8,
                            fontWeight: FontWeight.bold,
                            color: isConnected ? const Color(0xFF00E676) : Colors.grey,
                          ),
                        ),
                        const SizedBox(width: 6),
                        Expanded(
                          child: Text(
                            device.name,
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                              color: isDark ? Colors.white : Colors.black87,
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
                        style: TextStyle(
                          fontSize: 9,
                          color: isDark ? Colors.white54 : Colors.black54,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                    const SizedBox(height: 3),
                    Row(
                      children: [
                        Icon(Icons.battery_std, size: 10,
                            color: device.batteryLevel < 20 ? Colors.redAccent : const Color(0xFF00E676)),
                        const SizedBox(width: 2),
                        Text(
                          '${device.batteryLevel}%',
                          style: TextStyle(
                            fontSize: 9,
                            color: isDark ? Colors.white54 : Colors.black54,
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

// ═══════════════════════════════════════════════════════════════════════════
// Nodes Painter — 각 기기를 원형 노드로 그리기
// ═══════════════════════════════════════════════════════════════════════════

class _NodesPainter extends CustomPainter {
  final List<({Offset nodePos, Offset surfacePos, double angle})> positions;
  final List<ConnectedDevice> devices;
  final String? selectedId;
  final bool isDark;
  final int? hoveredIndex;

  _NodesPainter({
    required this.positions,
    required this.devices,
    required this.selectedId,
    required this.isDark,
    this.hoveredIndex,
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
          ? const Color(0xFF00E676)
          : (isAlert ? Colors.redAccent : Colors.grey);

      final isHovered = i == hoveredIndex;
      final double nodeRadius = isSelected ? 23.0 : (isHovered ? 21.0 : 19.0);

      // 선택 글로우
      if (isSelected) {
        canvas.drawCircle(
          pos,
          nodeRadius + 8,
          Paint()
            ..color = AppTheme.waveCyan.withValues(alpha: 0.35)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 10),
        );
      } else if (isHovered) {
        // 호버 글로우
        canvas.drawCircle(
          pos,
          nodeRadius + 6,
          Paint()
            ..color = AppTheme.waveCyan.withValues(alpha: 0.25)
            ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 8),
        );
      }

      // 배경 원
      canvas.drawCircle(
        pos,
        nodeRadius,
        Paint()
          ..color = isDark
              ? const Color(0xFF1A1A2E).withValues(alpha: 0.92)
              : Colors.white.withValues(alpha: 0.92)
          ..style = PaintingStyle.fill,
      );

      // 테두리
      canvas.drawCircle(
        pos,
        nodeRadius,
        Paint()
          ..color = isSelected ? AppTheme.waveCyan : (isHovered ? AppTheme.waveCyan.withValues(alpha: 0.7) : statusColor)
          ..style = PaintingStyle.stroke
          ..strokeWidth = isSelected ? 2.5 : (isHovered ? 2.0 : 1.5),
      );

      // 아이콘
      _paintIcon(
        canvas,
        _deviceTypeIcon(device.type),
        pos,
        isSelected ? 18.0 : 15.0,
        isSelected
            ? AppTheme.sanggamGold
            : (isDark ? Colors.white70 : Colors.black87),
      );

      // LIVE 점 (연결 시)
      if (device.status == DeviceConnectionStatus.connected) {
        final dotPos =
            Offset(pos.dx + nodeRadius * 0.65, pos.dy - nodeRadius * 0.65);
        canvas.drawCircle(dotPos, 4, Paint()..color = const Color(0xFF00E676));
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
        style: TextStyle(
          fontSize: 10,
          color: isDark ? Colors.white : Colors.black87,
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
        ..color = (isDark ? Colors.black : Colors.white).withValues(alpha: 0.8),
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
      oldDelegate.selectedId != selectedId ||
      oldDelegate.devices.length != devices.length ||
      oldDelegate.hoveredIndex != hoveredIndex;
}

// ═══════════════════════════════════════════════════════════════════════════
// Background Rings Painter
// ═══════════════════════════════════════════════════════════════════════════

class _BackgroundRingsPainter extends CustomPainter {
  final Offset center;
  final bool isDark;
  final double radius;

  _BackgroundRingsPainter({
    required this.center,
    required this.isDark,
    required this.radius,
  });

  @override
  void paint(Canvas canvas, Size size) {
    final color = isDark
        ? Colors.white.withValues(alpha: 0.05)
        : Colors.black.withValues(alpha: 0.05);
    final paint = Paint()
      ..color = color
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1;

    canvas.drawCircle(center, radius * 1.5, paint);
    canvas.drawCircle(center, radius * 2.5, paint);
    canvas.drawCircle(center, radius * 4.0, paint..strokeWidth = 0.5);
  }

  @override
  bool shouldRepaint(covariant _BackgroundRingsPainter old) =>
      old.center != center || old.isDark != isDark || old.radius != radius;
}

// ═══════════════════════════════════════════════════════════════════════════
// All Connectors Painter (단일 CustomPaint로 모든 커넥터)
// ═══════════════════════════════════════════════════════════════════════════

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
      oldDelegate.animationValue != animationValue;
}

// ═══════════════════════════════════════════════════════════════════════════
// Compact Stat Tile (요약바 내부)
// ═══════════════════════════════════════════════════════════════════════════

class _CompactStatTile extends StatelessWidget {
  final IconData icon;
  final String value;
  final String label;
  final Color color;

  const _CompactStatTile({
    required this.icon,
    required this.value,
    required this.label,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
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
            color: Colors.grey,
            fontSize: 8,
          ),
        ),
      ],
    );
  }
}

// ═══════════════════════════════════════════════════════════════════════════
// SafeHitTestWrapper — hitTest 중 viewport null check 크래시 방어
// ═══════════════════════════════════════════════════════════════════════════

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
      // DraggableScrollableSheet 내부 Viewport 리빌드 타이밍 이슈 방어
      return false;
    }
  }
}
