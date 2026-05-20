import 'dart:async';
import 'dart:math' as math;
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/measurement/application/measurement_golden_path_orchestrator.dart';
import 'package:manpasik/features/measurement/domain/measurement_evidence_presentation.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/utils/navigation_utils.dart';

const int _kParticleCount = 1792;
const int _kOrbitCount = 7;
const double _kHeaderRatio = 0.08; // 8%
const double _kDockRatio = 0.12; // 12%
const double _kDockVisualScale = 0.70; // 기존 대비 70%
const Color _kDeepSeaBg0 = Color(0xFF031025);
const Color _kDeepSeaBg1 = Color(0xFF072447);
const Color _kDeepSeaBg2 = Color(0xFF0A2F57);
const Color _kGold = Color(0xFFD4AF37);
const Color _kGoldLight = Color(0xFFFFE5A0);
const Color _kGoldDeep = Color(0xFF6F5217);
const Color _kJagaeCyan = Color(0xFF6FE6FF);
const Color _kJagaeMint = Color(0xFF9BF8E4);
const Color _kJagaeBlue = Color(0xFF3DB8FF);
const Color _kSurfaceTone = Color(0xFF0C1D35);
const Color _kSurfaceToneStrong = Color(0xFF081326);
const String _kMeasureDesignModePrefKey = 'measure_design_mode';

enum _SlimDesignMode { royalOrbit, jagaeWave, hanjiGold }

extension _SlimDesignModeX on _SlimDesignMode {
  String get folderName => switch (this) {
        _SlimDesignMode.royalOrbit => 'royal_orbit',
        _SlimDesignMode.jagaeWave => 'jagae_wave',
        _SlimDesignMode.hanjiGold => 'hanji_gold',
      };

  String get label => switch (this) {
        _SlimDesignMode.royalOrbit => '로열 오빗',
        _SlimDesignMode.jagaeWave => '자개 웨이브',
        _SlimDesignMode.hanjiGold => '한지 골드',
      };

  String assetPath(String fileName) =>
      'assets/images/concepts/$folderName/$fileName';
}

_SlimDesignMode _modeFromStorage(String? raw) {
  return switch (raw) {
    'royal_orbit' => _SlimDesignMode.royalOrbit,
    'jagae_wave' => _SlimDesignMode.jagaeWave,
    'hanji_gold' => _SlimDesignMode.hanjiGold,
    _ => _SlimDesignMode.royalOrbit,
  };
}

const List<({IconData icon, String label, String path})> _kMeasureNavActions = [
  (icon: Icons.home_outlined, label: '홈', path: '/home'),
  (icon: Icons.bar_chart_outlined, label: '데이터', path: '/data'),
  // Phase AP-2: 카트리지 비전 분석 진입 (codex VisionAnalyzerScreen 연결)
  (icon: Icons.qr_code_scanner_outlined, label: '카트리지 비전', path: '/measure/vision'),
  (icon: Icons.local_hospital_outlined, label: '의료', path: '/medical'),
  (icon: Icons.shopping_cart_outlined, label: '마켓', path: '/market'),
  (icon: Icons.forum_outlined, label: '커뮤니티', path: '/community'),
  (icon: Icons.settings_outlined, label: '설정', path: '/settings'),
];

Future<void> _showMeasureNavigationSheet(BuildContext context) async {
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
                '이동 메뉴',
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w700,
                  fontSize: 17,
                ),
              ),
            ),
            const Divider(height: 1, color: Color(0x33D4AF37)),
            ..._kMeasureNavActions.map((entry) {
              return ListTile(
                leading: Icon(entry.icon, color: _kGold),
                title: Text(entry.label,
                    style: const TextStyle(color: Colors.white)),
                onTap: () {
                  Navigator.of(sheetContext).pop();
                  context.go(entry.path);
                },
              );
            }),
          ],
        ),
      );
    },
  );
}

class MeasureScreen extends ConsumerStatefulWidget {
  const MeasureScreen({super.key});

  @override
  ConsumerState<MeasureScreen> createState() => _MeasureScreenState();
}

class _MeasureScreenState extends ConsumerState<MeasureScreen>
    with TickerProviderStateMixin {
  late final AnimationController _particleController;
  late final AnimationController _scanController;
  late final List<JagaeParticleNode> _particles;
  _SlimDesignMode _designMode = _SlimDesignMode.royalOrbit;
  bool _isScanning = true;
  bool _isGoldenPathRunning = false;
  MeasurementGoldenPathSnapshot? _goldenPathSnapshot;

  @override
  void initState() {
    super.initState();
    _particles = _buildParticles();
    _restoreSavedDesignMode();

    _particleController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 36),
    )..repeat();

    _scanController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 9),
    )
      ..addListener(() => setState(() {}))
      ..addStatusListener((status) {
        if (status == AnimationStatus.completed) {
          setState(() => _isScanning = false);
        }
      })
      ..forward();
  }

  Future<void> _restoreSavedDesignMode() async {
    final prefs = await SharedPreferences.getInstance();
    final savedMode =
        _modeFromStorage(prefs.getString(_kMeasureDesignModePrefKey));
    if (!mounted || savedMode == _designMode) {
      return;
    }
    setState(() {
      _designMode = savedMode;
    });
  }

  Future<void> _saveDesignMode(_SlimDesignMode mode) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_kMeasureDesignModePrefKey, mode.folderName);
  }

  List<JagaeParticleNode> _buildParticles() {
    final random = math.Random(1792);
    return List.generate(_kParticleCount, (index) {
      return JagaeParticleNode(
        orbit: index % _kOrbitCount,
        angle: random.nextDouble() * math.pi * 2,
        speed: 0.35 + random.nextDouble() * 1.8,
        radialJitter: random.nextDouble(),
        size: 0.65 + random.nextDouble() * 1.8,
        alpha: 0.22 + random.nextDouble() * 0.72,
        phase: random.nextDouble() * math.pi * 2,
        drift: 0.1 + random.nextDouble() * 0.9,
      );
    });
  }

  Future<void> _runGoldenPath() async {
    if (_isScanning || _isGoldenPathRunning) {
      return;
    }
    final authState = ref.read(authProvider);
    final userId = authState.userId;
    if (userId == null || userId.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('로그인 후 측정을 시작할 수 있습니다.')),
      );
      return;
    }

    setState(() {
      _isScanning = true;
      _isGoldenPathRunning = true;
      _goldenPathSnapshot = const MeasurementGoldenPathSnapshot(
        phase: MeasurementGoldenPathPhase.idle,
      );
    });
    _scanController.forward(from: 0);

    final orchestrator = MeasurementGoldenPathOrchestrator(
      repository: ref.read(measurementRepositoryProvider),
      traceSink: ref.read(measurementGoldenPathTraceSinkProvider),
      onSnapshot: (snapshot) {
        if (!mounted) return;
        setState(() {
          _goldenPathSnapshot = snapshot;
        });
      },
    );

    try {
      await orchestrator.run(
        userId: userId,
        deviceId: 'device-1',
      );
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('측정 골든 패스가 완료되었습니다.')),
      );
    } catch (error) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('측정 처리 실패: $error')),
      );
    } finally {
      if (mounted) {
        setState(() {
          _isGoldenPathRunning = false;
          _isScanning = false;
        });
      }
    }
  }

  void _changeDesignMode(_SlimDesignMode mode) {
    if (_designMode == mode) {
      return;
    }
    setState(() {
      _designMode = mode;
    });
    unawaited(_saveDesignMode(mode));
  }

  @override
  void dispose() {
    _particleController.dispose();
    _scanController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final media = MediaQuery.of(context);

    return Scaffold(
      body: LayoutBuilder(
        builder: (context, constraints) {
          final safeTop = media.padding.top;
          final safeBottom = media.padding.bottom;
          final usableHeight = math.max(
            0.0,
            constraints.maxHeight - safeTop - safeBottom,
          );

          final headerBodyHeight =
              (usableHeight * _kHeaderRatio).clamp(50.0, 72.0);
          final dockBodyHeight =
              ((usableHeight * _kDockRatio) * _kDockVisualScale)
                  .clamp(64.0, 96.0);
          final headerHeight = safeTop + headerBodyHeight;
          final dockHeight = safeBottom + dockBodyHeight;
          final progress = _scanController.value;

          return Column(
            children: [
              _SlimTopLayer(
                totalHeight: headerHeight,
                safeTop: safeTop,
                contentHeight: headerBodyHeight,
                isScanning: _isScanning,
                progress: progress,
                designMode: _designMode,
                onDesignModeChanged: _changeDesignMode,
              ),
              Expanded(
                child: _CenterLayer(
                  particles: _particles,
                  particleController: _particleController,
                  isScanning: _isScanning,
                  progress: progress,
                ),
              ),
              _SlimBottomLayer(
                totalHeight: dockHeight,
                safeBottom: safeBottom,
                contentHeight: dockBodyHeight,
                isScanning: _isScanning,
                progress: progress,
                designMode: _designMode,
                statusText: _measureStatusText,
                onActionTap: () => unawaited(_runGoldenPath()),
              ),
            ],
          );
        },
      ),
    );
  }

  String get _measureStatusText {
    final snapshot = _goldenPathSnapshot;
    if (_isGoldenPathRunning && snapshot != null) {
      return switch (snapshot.phase) {
        MeasurementGoldenPathPhase.idle => 'READY',
        MeasurementGoldenPathPhase.readinessChecked =>
          snapshot.engineMode?.toUpperCase() ?? 'ENGINE READY',
        MeasurementGoldenPathPhase.cartridgeRead => 'NFC OK',
        MeasurementGoldenPathPhase.sessionStarted => 'SESSION',
        MeasurementGoldenPathPhase.engineProcessed => 'ENGINE',
        MeasurementGoldenPathPhase.serverProcessed =>
          _evidenceStatusText(snapshot) ?? 'SERVER',
        MeasurementGoldenPathPhase.sessionEnded =>
          _evidenceStatusText(snapshot) ?? 'DONE',
        MeasurementGoldenPathPhase.failed => 'FAILED',
      };
    }
    return _isScanning ? 'SYNC LIVE' : 'SYNC DONE';
  }

  String? _evidenceStatusText(MeasurementGoldenPathSnapshot snapshot) {
    final evidenceStatus = snapshot.evidenceStatus;
    if (evidenceStatus == null || evidenceStatus.isEmpty) {
      return null;
    }
    return MeasurementEvidencePresentation.from(
      evidenceStatus: evidenceStatus,
      diagnosticReady: snapshot.diagnosticReady,
      evidenceGaps: snapshot.evidenceGaps,
    ).badgeLabel;
  }
}

class _SlimTopLayer extends StatelessWidget {
  const _SlimTopLayer({
    required this.totalHeight,
    required this.safeTop,
    required this.contentHeight,
    required this.isScanning,
    required this.progress,
    required this.designMode,
    required this.onDesignModeChanged,
  });

  final double totalHeight;
  final double safeTop;
  final double contentHeight;
  final bool isScanning;
  final double progress;
  final _SlimDesignMode designMode;
  final ValueChanged<_SlimDesignMode> onDesignModeChanged;

  @override
  Widget build(BuildContext context) {
    final progressPct = (progress * 100).clamp(0, 100).toInt();

    return SizedBox(
      height: totalHeight,
      child: Stack(
        children: [
          Positioned.fill(
            child: Opacity(
              opacity: 0.32,
              child: ColorFiltered(
                colorFilter: const ColorFilter.mode(
                  _kSurfaceTone,
                  BlendMode.multiply,
                ),
                child: Image.asset(
                  designMode.assetPath('header_3d_frame_slim.png'),
                  fit: BoxFit.cover,
                  alignment: Alignment.center,
                  filterQuality: FilterQuality.high,
                  errorBuilder: (_, __, ___) => Image.asset(
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
          Positioned.fill(
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    _kSurfaceToneStrong.withValues(alpha: 0.62),
                    _kSurfaceTone.withValues(alpha: 0.22),
                    Colors.transparent,
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            top: safeTop,
            left: 12,
            right: 12,
            height: contentHeight,
            child: Row(
              children: [
                IconButton(
                  onPressed: () => context.popOrHome(),
                  icon: const Icon(Icons.arrow_back_ios_new_rounded),
                  iconSize: 16,
                  color: _kGold,
                  tooltip: '뒤로',
                  splashRadius: 18,
                ),
                IconButton(
                  onPressed: () => _showMeasureNavigationSheet(context),
                  icon: const Icon(Icons.menu_rounded),
                  iconSize: 18,
                  color: _kGold,
                  tooltip: '메뉴',
                  splashRadius: 18,
                ),
                const SizedBox(width: 6),
                _StatusBadge(
                  label: isScanning ? 'SCANNING' : 'COMPLETE',
                  value: isScanning ? '$progressPct%' : '100%',
                  active: isScanning,
                ),
                const Spacer(),
                _DesignModeMenu(
                  mode: designMode,
                  onChanged: onDesignModeChanged,
                ),
                const SizedBox(width: 8),
                _PillBadge(
                  text: isScanning ? 'BIO LIVE' : 'BIO READY',
                  color: isScanning ? _kJagaeCyan : _kGold,
                  icon: Icons.auto_awesome_rounded,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _CenterLayer extends StatelessWidget {
  const _CenterLayer({
    required this.particles,
    required this.particleController,
    required this.isScanning,
    required this.progress,
  });

  final List<JagaeParticleNode> particles;
  final AnimationController particleController;
  final bool isScanning;
  final double progress;

  @override
  Widget build(BuildContext context) {
    final progressPct = (progress * 100).clamp(0, 100).toInt();

    return Stack(
      children: [
        // Center Layer: code-only DeepSea Blue background.
        const Positioned.fill(
          child: DecoratedBox(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [_kDeepSeaBg0, _kDeepSeaBg1, _kDeepSeaBg2],
              ),
            ),
          ),
        ),
        Positioned.fill(
          child: Container(
            decoration: BoxDecoration(
              gradient: RadialGradient(
                center: const Alignment(0.0, -0.2),
                radius: 1.05,
                colors: [
                  _kJagaeBlue.withValues(alpha: 0.16),
                  Colors.transparent,
                ],
              ),
            ),
          ),
        ),
        // Center Layer: 1792-dimensional particle field covering full area.
        Positioned.fill(
          child: AnimatedBuilder(
            animation: particleController,
            builder: (context, _) {
              return SizedBox.expand(
                child: CustomPaint(
                  painter: JagaeParticleSystem(
                    particles: particles,
                    time: particleController.value * math.pi * 2,
                    isScanning: isScanning,
                  ),
                ),
              );
            },
          ),
        ),
        Center(
          child: ClipRRect(
            borderRadius: BorderRadius.circular(20),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 7, sigmaY: 7),
              child: Container(
                width: 280,
                padding:
                    const EdgeInsets.symmetric(horizontal: 18, vertical: 16),
                decoration: BoxDecoration(
                  color: Colors.black.withValues(alpha: 0.2),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(
                    color: _kGold.withValues(alpha: 0.22),
                    width: 0.8,
                  ),
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      isScanning ? 'BODY SCAN ACTIVE' : 'SIGNATURE LOCKED',
                      style: TextStyle(
                        color: _kGold.withValues(alpha: 0.95),
                        fontWeight: FontWeight.w800,
                        letterSpacing: 2.2,
                        fontSize: 11,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      isScanning ? '$progressPct%' : '100%',
                      style: TextStyle(
                        color: isScanning ? _kJagaeCyan : _kJagaeMint,
                        fontSize: 30,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      isScanning
                          ? '1792D 파티클 스캔이 실시간 동기화 중입니다.'
                          : '데이터 시그니처 동기화가 완료되었습니다.',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        color: Colors.white.withValues(alpha: 0.78),
                        fontSize: 12,
                        height: 1.35,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class _SlimBottomLayer extends StatelessWidget {
  const _SlimBottomLayer({
    required this.totalHeight,
    required this.safeBottom,
    required this.contentHeight,
    required this.isScanning,
    required this.progress,
    required this.designMode,
    required this.statusText,
    required this.onActionTap,
  });

  final double totalHeight;
  final double safeBottom;
  final double contentHeight;
  final bool isScanning;
  final double progress;
  final _SlimDesignMode designMode;
  final String statusText;
  final VoidCallback onActionTap;

  @override
  Widget build(BuildContext context) {
    final buttonSize = contentHeight * 0.86;
    final progressPct = (progress * 100).clamp(0, 100).toInt();

    return SizedBox(
      height: totalHeight,
      child: Stack(
        clipBehavior: Clip.none,
        children: [
          Positioned.fill(
            child: Opacity(
              opacity: 0.3,
              child: ColorFiltered(
                colorFilter: const ColorFilter.mode(
                  _kSurfaceTone,
                  BlendMode.multiply,
                ),
                child: Image.asset(
                  designMode.assetPath('bottom_3d_dock_slim.png'),
                  fit: BoxFit.cover,
                  alignment: Alignment.bottomCenter,
                  filterQuality: FilterQuality.high,
                  errorBuilder: (_, __, ___) => Image.asset(
                    'assets/images/bottom_3d_dock_slim.png',
                    fit: BoxFit.cover,
                    alignment: Alignment.bottomCenter,
                    filterQuality: FilterQuality.high,
                    errorBuilder: (_, __, ___) => DecoratedBox(
                      decoration: BoxDecoration(
                        color: _kDeepSeaBg0.withValues(alpha: 0.9),
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ),
          Positioned.fill(
            child: DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    _kSurfaceToneStrong.withValues(alpha: 0.05),
                    _kSurfaceTone.withValues(alpha: 0.18),
                    _kDeepSeaBg0.withValues(alpha: 0.38),
                  ],
                ),
              ),
            ),
          ),
          Positioned.fill(
            child: IgnorePointer(
              child: Opacity(
                opacity: 0.36,
                child: ShaderMask(
                  blendMode: BlendMode.srcIn,
                  shaderCallback: (bounds) {
                    return LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        _kGold.withValues(alpha: 0.22),
                        _kGold.withValues(alpha: 0.52),
                      ],
                    ).createShader(bounds);
                  },
                  child: const CustomPaint(
                    painter: _TraditionalPatternPainter(),
                  ),
                ),
              ),
            ),
          ),
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            child: Container(
              height: 1.1,
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.centerLeft,
                  end: Alignment.centerRight,
                  colors: [
                    Colors.transparent,
                    _kGold.withValues(alpha: 0.58),
                    Colors.transparent,
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            left: 12,
            right: 12,
            bottom: safeBottom + 8,
            height: contentHeight - 12,
            child: Row(
              children: [
                Expanded(
                  child: _PillBadge(
                    text: isScanning ? 'SIGNAL $progressPct%' : 'SIGNAL 100%',
                    color: _kJagaeCyan,
                    icon: Icons.network_check_rounded,
                  ),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: _PillBadge(
                    text: statusText,
                    color: _kJagaeMint,
                    icon: Icons.sync_rounded,
                  ),
                ),
              ],
            ),
          ),
          Positioned(
            top: -(buttonSize * 0.38),
            left: 0,
            right: 0,
            child: Center(
              child: _ActionImageButton(
                size: buttonSize,
                enabled: !isScanning,
                designMode: designMode,
                onTap: onActionTap,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ActionImageButton extends StatelessWidget {
  const _ActionImageButton({
    required this.size,
    required this.enabled,
    required this.designMode,
    required this.onTap,
  });

  final double size;
  final bool enabled;
  final _SlimDesignMode designMode;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: enabled ? onTap : null,
      child: Stack(
        alignment: Alignment.center,
        children: [
          Container(
            width: size,
            height: size,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: _kJagaeCyan.withValues(alpha: enabled ? 0.42 : 0.18),
                  blurRadius: 26,
                  spreadRadius: 4,
                ),
              ],
            ),
          ),
          Image.asset(
            designMode.assetPath('btn_3d_action_slim.png'),
            width: size,
            height: size,
            fit: BoxFit.cover,
            errorBuilder: (_, __, ___) => Image.asset(
              'assets/images/btn_3d_action_slim.png',
              width: size,
              height: size,
              fit: BoxFit.cover,
              errorBuilder: (_, __, ___) => Container(
                width: size,
                height: size,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  gradient: const LinearGradient(
                    begin: Alignment.topLeft,
                    end: Alignment.bottomRight,
                    colors: [_kJagaeBlue, _kJagaeCyan],
                  ),
                  border: Border.all(color: _kGold, width: 2),
                ),
              ),
            ),
          ),
          _EmbossedGoldIcon(
            icon: enabled
                ? Icons.play_arrow_rounded
                : Icons.hourglass_top_rounded,
            size: size * 0.34,
          ),
        ],
      ),
    );
  }
}

class _DesignModeMenu extends StatelessWidget {
  const _DesignModeMenu({
    required this.mode,
    required this.onChanged,
  });

  final _SlimDesignMode mode;
  final ValueChanged<_SlimDesignMode> onChanged;

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<_SlimDesignMode>(
      tooltip: '디자인 모드',
      onSelected: onChanged,
      itemBuilder: (context) {
        return _SlimDesignMode.values.map((entry) {
          return PopupMenuItem<_SlimDesignMode>(
            value: entry,
            child: Row(
              children: [
                Icon(
                  entry == mode
                      ? Icons.radio_button_checked_rounded
                      : Icons.radio_button_unchecked_rounded,
                  size: 14,
                  color: entry == mode ? _kGold : _kJagaeBlue,
                ),
                const SizedBox(width: 8),
                Text(entry.label),
              ],
            ),
          );
        }).toList();
      },
      child: _PillBadge(
        text: mode.label,
        color: _kGold,
        icon: Icons.palette_rounded,
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({
    required this.label,
    required this.value,
    required this.active,
  });

  final String label;
  final String value;
  final bool active;

  @override
  Widget build(BuildContext context) {
    final color = active ? _kJagaeCyan : _kGold;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        color: Colors.black.withValues(alpha: 0.28),
        border: Border.all(color: color.withValues(alpha: 0.38), width: 0.8),
      ),
      child: Row(
        children: [
          Container(
            width: 6,
            height: 6,
            decoration: BoxDecoration(shape: BoxShape.circle, color: color),
          ),
          const SizedBox(width: 6),
          Text(
            label,
            style: TextStyle(
              color: color,
              fontWeight: FontWeight.w700,
              fontSize: 10,
              letterSpacing: 1.1,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            value,
            style: TextStyle(
              color: Colors.white.withValues(alpha: 0.92),
              fontWeight: FontWeight.w700,
              fontSize: 10,
            ),
          ),
        ],
      ),
    );
  }
}

class _PillBadge extends StatelessWidget {
  const _PillBadge({
    required this.text,
    required this.color,
    this.icon,
  });

  final String text;
  final Color color;
  final IconData? icon;

  @override
  Widget build(BuildContext context) {
    final edgeColor = Color.lerp(_kGold, color, 0.35) ?? _kGold;
    return Container(
      height: 34,
      alignment: Alignment.center,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(999),
        gradient: LinearGradient(
          begin: Alignment.topCenter,
          end: Alignment.bottomCenter,
          colors: [
            _kSurfaceToneStrong.withValues(alpha: 0.72),
            _kSurfaceTone.withValues(alpha: 0.56),
          ],
        ),
        border: Border.all(
          color: edgeColor.withValues(alpha: 0.52),
          width: 0.85,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.28),
            blurRadius: 8,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (icon != null) ...[
            _EmbossedGoldIcon(
              icon: icon!,
              size: 15,
            ),
            const SizedBox(width: 7),
          ],
          Flexible(
            child: Text(
              text,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                color: color,
                fontWeight: FontWeight.w700,
                fontSize: 11,
                letterSpacing: 0.8,
                shadows: [
                  Shadow(
                    color: Colors.black.withValues(alpha: 0.55),
                    blurRadius: 3,
                    offset: const Offset(0, 1.2),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _EmbossedGoldIcon extends StatelessWidget {
  const _EmbossedGoldIcon({
    required this.icon,
    required this.size,
  });

  final IconData icon;
  final double size;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: size,
      height: size,
      child: Stack(
        alignment: Alignment.center,
        children: [
          Transform.translate(
            offset: const Offset(0.9, 0.9),
            child: Icon(
              icon,
              size: size,
              color: Colors.black.withValues(alpha: 0.55),
            ),
          ),
          ShaderMask(
            blendMode: BlendMode.srcIn,
            shaderCallback: (bounds) {
              return const LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [_kGoldLight, _kGold, _kGoldDeep],
                stops: [0.0, 0.58, 1.0],
              ).createShader(bounds);
            },
            child: Icon(
              icon,
              size: size,
              color: Colors.white,
            ),
          ),
          Transform.translate(
            offset: const Offset(-0.6, -0.6),
            child: Icon(
              icon,
              size: size * 0.92,
              color: Colors.white.withValues(alpha: 0.32),
            ),
          ),
        ],
      ),
    );
  }
}

class _TraditionalPatternPainter extends CustomPainter {
  const _TraditionalPatternPainter();

  @override
  void paint(Canvas canvas, Size size) {
    final step = math.max(24.0, size.height * 0.48);
    final ring = step * 0.22;
    final strokePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.8
      ..color = Colors.white.withValues(alpha: 0.66);
    final dotPaint = Paint()
      ..style = PaintingStyle.fill
      ..color = Colors.white.withValues(alpha: 0.52);
    final linePaint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 0.65
      ..color = Colors.white.withValues(alpha: 0.58);

    for (double y = -step; y <= size.height + step; y += step) {
      for (double x = -step; x <= size.width + step; x += step) {
        final center = Offset(x + step * 0.5, y + step * 0.5);
        canvas.drawCircle(center, ring, strokePaint);
        canvas.drawLine(
          Offset(center.dx - ring, center.dy),
          Offset(center.dx + ring, center.dy),
          linePaint,
        );
        canvas.drawLine(
          Offset(center.dx, center.dy - ring),
          Offset(center.dx, center.dy + ring),
          linePaint,
        );
        canvas.drawCircle(center, 1.1, dotPaint);
      }
    }
  }

  @override
  bool shouldRepaint(covariant _TraditionalPatternPainter oldDelegate) {
    return false;
  }
}

class JagaeParticleNode {
  const JagaeParticleNode({
    required this.orbit,
    required this.angle,
    required this.speed,
    required this.radialJitter,
    required this.size,
    required this.alpha,
    required this.phase,
    required this.drift,
  });

  final int orbit;
  final double angle;
  final double speed;
  final double radialJitter;
  final double size;
  final double alpha;
  final double phase;
  final double drift;
}

class JagaeParticleSystem extends CustomPainter {
  JagaeParticleSystem({
    required this.particles,
    required this.time,
    required this.isScanning,
  });

  final List<JagaeParticleNode> particles;
  final double time;
  final bool isScanning;

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width * 0.5, size.height * 0.52);
    final baseRadius = math.min(size.width, size.height) * 0.46;

    final pulse = isScanning ? 0.72 + 0.28 * math.sin(time * 1.2) : 0.82;
    final palette = <Color>[_kJagaeCyan, _kJagaeBlue, _kJagaeMint, _kGold];

    for (var i = 0; i < particles.length; i++) {
      final p = particles[i];
      final direction = p.orbit.isEven ? 1.0 : -1.0;
      final orbitRadius =
          baseRadius * (0.18 + p.orbit * 0.105 + p.radialJitter * 0.032);
      final angle = p.angle + (time * p.speed * direction);
      final wobble = math.sin(time * 1.8 + p.phase) * p.drift;

      final x = center.dx + math.cos(angle) * (orbitRadius + wobble * 16);
      final y =
          center.dy + math.sin(angle) * (orbitRadius * 0.64 + wobble * 12);

      final color = palette[(p.orbit + i) % palette.length];
      final alpha = (p.alpha * pulse).clamp(0.05, 0.95);
      final radius = p.size * (isScanning ? 1.0 : 0.82);

      canvas.drawCircle(
        Offset(x, y),
        radius,
        Paint()..color = color.withValues(alpha: alpha),
      );
    }
  }

  @override
  bool shouldRepaint(covariant JagaeParticleSystem oldDelegate) {
    return oldDelegate.time != time || oldDelegate.isScanning != isScanning;
  }
}
