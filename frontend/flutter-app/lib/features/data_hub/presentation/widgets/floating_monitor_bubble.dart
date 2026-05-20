import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitor_bubble_provider.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitoring_providers.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/mini_device_card.dart';

// ───────────────────────────────────────────────────
// FloatingMonitorBubble — Sanggam Orbit 모니터 버블
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] Theme.of(context) isDark 분기 제거 (항상 다크)
// [Rule 4] AppTheme.sanggamGold ~7x → SanggamTheme.primary
// [Rule 4] Colors.redAccent ~8x → SanggamTheme.error
// [Rule 4] Color(0xFF00E676) → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF1B2640/050B14) → SanggamTheme 상수
// ───────────────────────────────────────────────────

/// Floating Monitor Bubble — 모든 셸 페이지에 표시되는 리더기 현황 오버레이
///
/// 3단계 UX:
///   1. Compact Bubble (56x56) — 연결 수 + 경고 뱃지
///   2. Mini Dashboard (280x340) — 기기 카드 횡스크롤 + 요약
///   3. "전체 보기" → /data/monitoring (기존 HoloGlobe 화면)
class FloatingMonitorBubble
    extends ConsumerStatefulWidget {
  const FloatingMonitorBubble({super.key});

  @override
  ConsumerState<FloatingMonitorBubble>
      createState() =>
          _FloatingMonitorBubbleState();
}

class _FloatingMonitorBubbleState
    extends ConsumerState<FloatingMonitorBubble>
    with SingleTickerProviderStateMixin {
  late AnimationController _breatheController;

  double _bubbleX = -1;
  double _bubbleY = -1;
  bool _dragging = false;

  static const double _bubbleSize = 56;
  static const double _panelWidth = 280;
  static const double _panelHeight = 360;

  @override
  void initState() {
    super.initState();
    _breatheController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _breatheController.dispose();
    super.dispose();
  }

  void _initPosition(
      BoxConstraints constraints) {
    if (_bubbleX < 0) {
      _bubbleX = constraints.maxWidth -
          _bubbleSize -
          16;
      _bubbleY = constraints.maxHeight -
          _bubbleSize -
          100;
    }
  }

  void _snapToEdge(
      BoxConstraints constraints) {
    final centerX =
        _bubbleX + _bubbleSize / 2;
    final halfW = constraints.maxWidth / 2;
    _bubbleX = centerX < halfW
        ? 12
        : constraints.maxWidth -
            _bubbleSize -
            12;
    _bubbleY = _bubbleY.clamp(60,
        constraints.maxHeight - _bubbleSize - 90);
  }

  @override
  Widget build(BuildContext context) {
    final isExpanded = ref
        .watch(monitorBubbleExpandedProvider);
    final counts =
        ref.watch(connectedCountProvider);
    final alerts =
        ref.watch(deviceAlertCountProvider);

    final location =
        GoRouterState.of(context)
            .matchedLocation;
    if (location == '/data/monitoring') {
      return const SizedBox.shrink();
    }

    return LayoutBuilder(
      builder: (context, constraints) {
        _initPosition(constraints);

        double panelX = _bubbleX -
            _panelWidth +
            _bubbleSize;
        double panelY =
            _bubbleY - _panelHeight - 8;
        panelX = panelX.clamp(
            8,
            constraints.maxWidth -
                _panelWidth -
                8);
        panelY = panelY.clamp(
            60,
            constraints.maxHeight -
                _panelHeight -
                90);

        return Stack(
          children: [
            // Scrim when expanded
            if (isExpanded)
              Positioned.fill(
                child: GestureDetector(
                  onTap: () => ref
                      .read(
                          monitorBubbleExpandedProvider
                              .notifier)
                      .state = false,
                  child: Container(
                      color: Colors.black
                          .withValues(
                              alpha: 0.3)),
                ),
              ),

            // Mini Dashboard Panel
            if (isExpanded)
              AnimatedPositioned(
                duration: const Duration(
                    milliseconds: 250),
                curve: Curves.easeOutCubic,
                left: panelX,
                top: panelY,
                child: _buildMiniDashboard(),
              ),

            // Compact Bubble
            AnimatedPositioned(
              duration: _dragging
                  ? Duration.zero
                  : const Duration(
                      milliseconds: 300),
              curve: Curves.easeOutCubic,
              left: _bubbleX,
              top: _bubbleY,
              child: Tooltip(
                message: '디바이스 모니터',
                child: MouseRegion(
                  cursor: _dragging
                      ? SystemMouseCursors
                          .grabbing
                      : SystemMouseCursors
                          .grab,
                  child: GestureDetector(
                    onTap: () {
                      ref
                          .read(
                              monitorBubbleExpandedProvider
                                  .notifier)
                          .state = !isExpanded;
                    },
                    onPanStart: (_) =>
                        setState(() =>
                            _dragging = true),
                    onPanUpdate: (d) {
                      setState(() {
                        _bubbleX +=
                            d.delta.dx;
                        _bubbleY +=
                            d.delta.dy;
                        _bubbleX =
                            _bubbleX.clamp(
                                0,
                                constraints
                                        .maxWidth -
                                    _bubbleSize);
                        _bubbleY =
                            _bubbleY.clamp(
                                60,
                                constraints
                                        .maxHeight -
                                    _bubbleSize -
                                    20);
                      });
                    },
                    onPanEnd: (_) {
                      setState(() {
                        _dragging = false;
                        _snapToEdge(
                            constraints);
                      });
                      if (isExpanded) {
                        ref
                            .read(
                                monitorBubbleExpandedProvider
                                    .notifier)
                            .state = false;
                      }
                    },
                    child:
                        _buildCompactBubble(
                            counts, alerts),
                  ),
                ),
              ),
            ),
          ],
        );
      },
    );
  }

  /// Compact Bubble (56x56 원형)
  Widget _buildCompactBubble(
    ({int connected, int total}) counts,
    int alerts,
  ) {
    return AnimatedBuilder(
      animation: _breatheController,
      builder: (context, child) {
        final glowOpacity = 0.15 +
            _breatheController.value * 0.2;
        final scale = 1.0 +
            _breatheController.value * 0.02;

        return Transform.scale(
          scale: scale,
          child: Container(
            width: _bubbleSize,
            height: _bubbleSize,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: (alerts > 0
                          ? SanggamTheme.error
                          : SanggamTheme
                              .primary)
                      .withValues(
                          alpha:
                              glowOpacity),
                  blurRadius: 20,
                  spreadRadius: 2,
                ),
              ],
            ),
            child: ClipOval(
              child: BackdropFilter(
                filter: ImageFilter.blur(
                    sigmaX: 12, sigmaY: 12),
                child: Container(
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: LinearGradient(
                      begin:
                          Alignment.topLeft,
                      end: Alignment
                          .bottomRight,
                      colors: [
                        SanggamTheme
                            .surfaceVariant
                            .withValues(
                                alpha: 0.6),
                        SanggamTheme
                            .background
                            .withValues(
                                alpha: 0.8),
                      ],
                    ),
                    border: Border.all(
                      color: (alerts > 0
                              ? SanggamTheme
                                  .error
                              : SanggamTheme
                                  .primary)
                          .withValues(
                              alpha: 0.5),
                      width: 1.2,
                    ),
                  ),
                  child: Stack(
                    alignment:
                        Alignment.center,
                    children: [
                      const Icon(
                        Icons.sensors,
                        size: 22,
                        color: SanggamTheme
                            .primary,
                      ),
                      // Connected count badge
                      Positioned(
                        bottom: 6,
                        child: Text(
                          '${counts.connected}/${counts.total}',
                          style:
                              const TextStyle(
                            fontSize: 8,
                            fontWeight:
                                FontWeight
                                    .bold,
                            color: SanggamTheme
                                .onSurfaceDim,
                          ),
                        ),
                      ),
                      // Alert badge
                      if (alerts > 0)
                        Positioned(
                          top: 4,
                          right: 4,
                          child: Container(
                            width: 16,
                            height: 16,
                            decoration:
                                const BoxDecoration(
                              color:
                                  SanggamTheme
                                      .error,
                              shape: BoxShape
                                  .circle,
                            ),
                            child: Center(
                              child: Text(
                                '$alerts',
                                style:
                                    const TextStyle(
                                  color: Colors
                                      .white,
                                  fontSize: 9,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                ),
                              ),
                            ),
                          ),
                        ),
                    ],
                  ),
                ),
              ),
            ),
          ),
        );
      },
    );
  }

  /// Mini Dashboard (280x360 패널)
  Widget _buildMiniDashboard() {
    final devicesAsync = ref
        .watch(pollingConnectedDevicesProvider);
    final counts =
        ref.watch(connectedCountProvider);
    final alerts =
        ref.watch(deviceAlertCountProvider);

    return ClipRRect(
      borderRadius:
          BorderRadius.circular(20),
      child: BackdropFilter(
        filter: ImageFilter.blur(
            sigmaX: 20, sigmaY: 20),
        child: Container(
          width: _panelWidth,
          height: _panelHeight,
          decoration: BoxDecoration(
            borderRadius:
                BorderRadius.circular(20),
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                SanggamTheme.surfaceVariant
                    .withValues(alpha: 0.5),
                SanggamTheme.background
                    .withValues(alpha: 0.7),
              ],
            ),
            border: Border.all(
              color: SanggamTheme.primary
                  .withValues(alpha: 0.3),
              width: 0.8,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black
                    .withValues(alpha: 0.4),
                blurRadius: 24,
                offset: const Offset(0, 8),
              ),
            ],
          ),
          child: Column(
            children: [
              // Header
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 14, 12, 8),
                child: Row(
                  children: [
                    Container(
                      width: 8,
                      height: 8,
                      decoration:
                          BoxDecoration(
                        color: alerts > 0
                            ? SanggamTheme
                                .error
                            : SanggamTheme
                                .jagaeCyan,
                        shape:
                            BoxShape.circle,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '디바이스 상태',
                      style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim
                            .withValues(
                                alpha: 0.8),
                        fontSize: 10,
                        fontWeight:
                            FontWeight.w600,
                        letterSpacing: 1.5,
                      ),
                    ),
                    const Spacer(),
                    Text(
                      '${counts.connected}/${counts.total} 연결',
                      style:
                          const TextStyle(
                        color: SanggamTheme
                            .primary,
                        fontSize: 12,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ],
                ),
              ),

              // Alert Banner
              if (alerts > 0)
                GestureDetector(
                  onTap: () {
                    ref
                        .read(
                            monitorBubbleExpandedProvider
                                .notifier)
                        .state = false;
                    context.push(
                        '/data/monitoring');
                  },
                  child: Container(
                    width: double.infinity,
                    margin: const EdgeInsets
                        .symmetric(
                        horizontal: 12),
                    padding: const EdgeInsets
                        .symmetric(
                        horizontal: 10,
                        vertical: 6),
                    decoration:
                        BoxDecoration(
                      color: SanggamTheme
                          .error
                          .withValues(
                              alpha: 0.1),
                      borderRadius:
                          BorderRadius
                              .circular(8),
                      border: Border.all(
                          color: SanggamTheme
                              .error
                              .withValues(
                                  alpha:
                                      0.3)),
                    ),
                    child: Row(
                      children: [
                        const Icon(
                          Icons
                              .warning_amber_rounded,
                          size: 14,
                          color: SanggamTheme
                              .error,
                        ),
                        const SizedBox(
                            width: 6),
                        Text(
                          '$alerts개 기기 주의 필요',
                          style:
                              const TextStyle(
                            color:
                                SanggamTheme
                                    .error,
                            fontSize: 11,
                            fontWeight:
                                FontWeight
                                    .w600,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              const SizedBox(height: 8),

              // Device Cards Scroll
              Expanded(
                child: devicesAsync.when(
                  data: (devices) {
                    if (devices.isEmpty) {
                      return const Center(
                        child: Text(
                          '등록된 기기가 없습니다',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          ),
                        ),
                      );
                    }
                    return ListView.builder(
                      scrollDirection:
                          Axis.horizontal,
                      padding:
                          const EdgeInsets
                              .symmetric(
                              horizontal:
                                  12),
                      itemCount:
                          devices.length,
                      itemBuilder: (_, i) =>
                          MiniDeviceCard(
                        device: devices[i],
                        onTap: () {
                          ref
                              .read(selectedDeviceIdProvider
                                  .notifier)
                              .state = devices[i]
                              .id;
                          ref
                              .read(
                                  monitorBubbleExpandedProvider
                                      .notifier)
                              .state = false;
                          context.push(
                              '/data/monitoring');
                        },
                      ),
                    );
                  },
                  loading: () =>
                      const Center(
                    child: SizedBox(
                      width: 24,
                      height: 24,
                      child:
                          CircularProgressIndicator(
                        strokeWidth: 2,
                        color: SanggamTheme
                            .primary,
                      ),
                    ),
                  ),
                  error: (_, __) =>
                      const Center(
                    child: Text(
                      '데이터 로드 실패',
                      style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      ),
                    ),
                  ),
                ),
              ),

              // Footer: "전체 보기" button
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        12, 4, 12, 12),
                child: SizedBox(
                  width: double.infinity,
                  height: 36,
                  child: TextButton(
                    onPressed: () {
                      ref
                          .read(
                              monitorBubbleExpandedProvider
                                  .notifier)
                          .state = false;
                      context.push(
                          '/data/monitoring');
                    },
                    style:
                        TextButton.styleFrom(
                      backgroundColor:
                          SanggamTheme.primary
                              .withValues(
                                  alpha:
                                      0.12),
                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius
                                .circular(
                                    10),
                      ),
                    ),
                    child: const Row(
                      mainAxisAlignment:
                          MainAxisAlignment
                              .center,
                      children: [
                        Text(
                          '전체 리더기 현황보기',
                          style: TextStyle(
                            color:
                                SanggamTheme
                                    .primary,
                            fontSize: 12,
                            fontWeight:
                                FontWeight
                                    .w600,
                          ),
                        ),
                        SizedBox(width: 4),
                        Icon(
                          Icons.open_in_new,
                          size: 14,
                          color: SanggamTheme
                              .primary,
                        ),
                      ],
                    ),
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
