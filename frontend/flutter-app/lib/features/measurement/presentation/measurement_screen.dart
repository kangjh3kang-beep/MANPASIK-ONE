import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/wave_ripple_painter.dart';
import 'package:manpasik/shared/widgets/breathing_overlay.dart';

/// 측정 상태
enum MeasurementStatus { idle, connecting, measuring, complete, error }

/// 측정 화면
///
/// measurement-service StartSession/EndSession gRPC 연동.
/// BLE/NFC는 S5b에서 연동.
///
/// Wave Ripple/Breathing 애니메이션 적용:
/// - connecting: WaveRippleBackground (동심원 파동)
/// - measuring: WavePainter (사인파 안정화) + BreathingOverlay
/// - 단계 텍스트: "파동 안정화 중..." → "분석 중..." → "측정 완료"
class MeasurementScreen extends ConsumerStatefulWidget {
  const MeasurementScreen({super.key});

  @override
  ConsumerState<MeasurementScreen> createState() => _MeasurementScreenState();
}

class _MeasurementScreenState extends ConsumerState<MeasurementScreen>
    with TickerProviderStateMixin {
  MeasurementStatus _status = MeasurementStatus.idle;
  String? _sessionId;
  String? _cartridgeId;
  double _measureProgress = 0.0;
  String _phaseText = '파동 안정화 중...';

  // Wave 애니메이션 컨트롤러
  AnimationController? _waveController;

  void _initWaveAnimation() {
    _waveController?.dispose();
    _waveController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 3),
    )..repeat();
  }

  void _disposeWaveAnimation() {
    _waveController?.dispose();
    _waveController = null;
  }

  Future<void> _readCartridge() async {
    try {
      final info = await RustFfiStub.nfcReadCartridge();
      if (!mounted) return;
      setState(() => _cartridgeId = info.cartridgeId);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('카트리지: ${info.cartridgeType} (${info.remainingUses}회 남음)'),
          behavior: SnackBarBehavior.floating,
        ),
      );
    } catch (_) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('카트리지 읽기 실패'), behavior: SnackBarBehavior.floating),
      );
    }
  }

  Future<void> _startMeasurement() async {
    final userId = ref.read(authProvider).userId;
    if (userId == null || userId.isEmpty) {
      setState(() => _status = MeasurementStatus.error);
      return;
    }

    setState(() => _status = MeasurementStatus.connecting);
    _initWaveAnimation();

    try {
      final repo = ref.read(measurementRepositoryProvider);
      final result = await repo.startSession(
        deviceId: 'device-1',
        cartridgeId: _cartridgeId ?? 'cartridge-1',
        userId: userId,
      );
      if (!mounted) return;
      _sessionId = result.sessionId;
      setState(() {
        _status = MeasurementStatus.measuring;
        _measureProgress = 0.0;
        _phaseText = '파동 안정화 중...';
      });
    } catch (_) {
      if (!mounted) return;
      _disposeWaveAnimation();
      setState(() => _status = MeasurementStatus.error);
      return;
    }

    // 시뮬레이션: 측정 진행 단계 (S5b에서 실제 스트림 연동)
    // Phase 1: 파동 안정화 중...
    for (int i = 0; i < 10; i++) {
      await Future.delayed(const Duration(milliseconds: 200));
      if (!mounted) return;
      setState(() => _measureProgress = (i + 1) / 20);
    }

    // Phase 2: 분석 중...
    if (!mounted) return;
    setState(() => _phaseText = '분석 중...');
    for (int i = 10; i < 20; i++) {
      await Future.delayed(const Duration(milliseconds: 200));
      if (!mounted) return;
      setState(() => _measureProgress = (i + 1) / 20);
    }

    if (!mounted) return;
    _disposeWaveAnimation();
    setState(() {
      _status = MeasurementStatus.complete;
      _phaseText = '측정 완료';
    });
  }

  Future<void> _endSession() async {
    if (_sessionId == null) return;
    try {
      await ref.read(measurementRepositoryProvider).endSession(_sessionId!);
    } catch (_) {
      // 무시
    }
    if (mounted) _sessionId = null;
  }

  void _reset() {
    _disposeWaveAnimation();
    _sessionId = null;
    setState(() {
      _status = MeasurementStatus.idle;
      _measureProgress = 0.0;
      _phaseText = '파동 안정화 중...';
    });
  }

  @override
  void dispose() {
    _disposeWaveAnimation();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMeasuring = _status == MeasurementStatus.measuring;

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('측정'),
      ),
      body: BreathingOverlay(
        enabled: isMeasuring,
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // 상태 표시 영역
                Expanded(
                  child: Center(
                    child: _buildStatusWidget(theme),
                  ),
                ),

                // 하단 버튼
                if (_status == MeasurementStatus.idle)
                  Column(
                    children: [
                      PrimaryButton(
                        text: '측정 시작',
                        icon: Icons.play_arrow_rounded,
                        onPressed: _startMeasurement,
                      ),
                      const SizedBox(height: 12),
                      OutlinedButton.icon(
                        onPressed: _readCartridge,
                        icon: const Icon(Icons.nfc),
                        label: Text(_cartridgeId != null ? '카트리지 읽음' : 'NFC 카트리지 읽기'),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                        ),
                      ),
                    ],
                  )
                else if (_status == MeasurementStatus.complete)
                  Column(
                    children: [
                      PrimaryButton(
                        text: '결과 확인',
                        icon: Icons.analytics_outlined,
                        onPressed: () async {
                          await _endSession();
                          if (!mounted) return;
                          context.push('/measurement/result');
                          _reset();
                        },
                      ),
                      const SizedBox(height: 12),
                      OutlinedButton(
                        onPressed: _reset,
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 56),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                        ),
                        child: const Text('다시 측정'),
                      ),
                    ],
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildStatusWidget(ThemeData theme) {
    switch (_status) {
      case MeasurementStatus.idle:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 160,
              height: 160,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer,
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.sensors_rounded,
                size: 72,
                color: theme.colorScheme.onPrimaryContainer,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              '디바이스를 준비해주세요',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '카트리지를 장착하고\n측정 버튼을 눌러주세요',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ],
        );

      case MeasurementStatus.connecting:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            SizedBox(
              width: 160,
              height: 160,
              child: _waveController != null
                  ? AnimatedBuilder(
                      animation: _waveController!,
                      builder: (context, _) {
                        return CustomPaint(
                          painter: WaveRipplePainter(
                            animationValue: _waveController!.value,
                            rippleCount: 5,
                          ),
                          child: Center(
                            child: Icon(
                              Icons.bluetooth_searching_rounded,
                              size: 48,
                              color: AppTheme.waveCyan,
                            ),
                          ),
                        );
                      },
                    )
                  : CircularProgressIndicator(
                      strokeWidth: 4,
                      color: theme.colorScheme.primary,
                    ),
            ),
            const SizedBox(height: 24),
            Text(
              '디바이스 연결 중...',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'BLE 연결을 시도하고 있습니다',
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ],
        );

      case MeasurementStatus.measuring:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // MANPASIK Wave 애니메이션 — 파동이 점차 안정화
            SizedBox(
              width: double.infinity,
              height: 120,
              child: _waveController != null
                  ? AnimatedBuilder(
                      animation: _waveController!,
                      builder: (context, _) {
                        return CustomPaint(
                          painter: WavePainter(
                            animationValue: _waveController!.value,
                            progress: _measureProgress,
                            waveColor: AppTheme.waveCyan,
                          ),
                        );
                      },
                    )
                  : const SizedBox.shrink(),
            ),
            const SizedBox(height: 32),
            // 단계별 텍스트 (Morph)
            AnimatedSwitcher(
              duration: const Duration(milliseconds: 500),
              child: Text(
                _phaseText,
                key: ValueKey(_phaseText),
                style: theme.textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppTheme.sanggamGold,
                ),
              ),
            ),
            const SizedBox(height: 16),
            // 진행률 바
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: LinearProgressIndicator(
                value: _measureProgress,
                minHeight: 6,
                backgroundColor: AppTheme.deepSeaBlue.withOpacity(0.3),
                valueColor: AlwaysStoppedAnimation<Color>(
                  Color.lerp(AppTheme.waveCyan, AppTheme.sanggamGold, _measureProgress)!,
                ),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${(_measureProgress * 100).toInt()}%',
              style: theme.textTheme.bodySmall?.copyWith(
                fontFamily: 'JetBrains Mono',
                color: AppTheme.waveCyan,
              ),
            ),
          ],
        );

      case MeasurementStatus.complete:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                color: Colors.green.withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.check_circle_rounded,
                size: 72,
                color: Colors.green,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              '측정 완료!',
              style: theme.textTheme.headlineMedium?.copyWith(
                fontWeight: FontWeight.bold,
                color: Colors.green,
              ),
            ),
            const SizedBox(height: 16),
            // 더미 결과
            Card(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    _buildResultItem(theme, '98.4', 'mg/dL', '혈당'),
                    Container(
                      width: 1,
                      height: 40,
                      color: theme.colorScheme.outlineVariant,
                    ),
                    _buildResultItem(theme, '정상', '', '판정'),
                  ],
                ),
              ),
            ),
          ],
        );

      case MeasurementStatus.error:
        return Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 72, color: Colors.red),
            const SizedBox(height: 24),
            Text(
              '측정 실패',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
                color: Colors.red,
              ),
            ),
          ],
        );
    }
  }

  Widget _buildResultItem(ThemeData theme, String value, String unit, String label) {
    return Column(
      children: [
        Text(
          value,
          style: theme.textTheme.headlineSmall?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        if (unit.isNotEmpty)
          Text(
            unit,
            style: theme.textTheme.labelSmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
        const SizedBox(height: 4),
        Text(
          label,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
      ],
    );
  }
}
