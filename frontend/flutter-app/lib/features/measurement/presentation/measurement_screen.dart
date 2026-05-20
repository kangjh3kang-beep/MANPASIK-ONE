import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/wave_ripple_painter.dart';
import 'package:manpasik/shared/widgets/breathing_overlay.dart';

// ───────────────────────────────────────────────────
// MeasurementScreen — Sanggam Orbit 측정 화면
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan 4x → SanggamTheme.jagaeCyan
// [Rule 4] AppTheme.deepSeaBlue → SanggamTheme.surfaceVariant
// [Rule 4] Theme.of(context).textTheme → 직접 TextStyle
// [Rule 4] theme.colorScheme.* → SanggamTheme 상수
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Card → 다크 테마 Container
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] lottie import 제거 (미사용)
// ───────────────────────────────────────────────────

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

  // 실측정 결과 (목업 아닌 실제 파이프라인 결과)
  MeasurementPipelineResult? _pipelineResult;

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
      final info = await RustBridge.nfcReadCartridge();
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

    // 실제 측정 파이프라인 실행 (RustBridge 경유)
    // Phase 1: 차동 계측 데이터 수집
    setState(() {
      _measureProgress = 0.1;
      _phaseText = '차동 신호 수집 중...';
    });
    await Future.delayed(const Duration(milliseconds: 300));
    if (!mounted) return;

    setState(() {
      _measureProgress = 0.3;
      _phaseText = 'S_diff 연산 중...';
    });

    try {
      // Phase 2: RustBridge 파이프라인 (BLE→DSP→AI)
      final result = await RustBridge.runMeasurementPipeline(
        deviceId: 'device-1',
        biomarker: 'glucose',
        unit: 'mg/dL',
      );
      if (!mounted) return;
      setState(() {
        _measureProgress = 0.7;
        _phaseText = 'AI 분석 완료...';
      });

      await Future.delayed(const Duration(milliseconds: 200));
      if (!mounted) return;
      _pipelineResult = result;

      setState(() {
        _measureProgress = 1.0;
        _phaseText = '측정 완료';
      });
    } catch (_) {
      if (!mounted) return;
      _disposeWaveAnimation();
      setState(() => _status = MeasurementStatus.error);
      return;
    }

    if (!mounted) return;
    _disposeWaveAnimation();
    setState(() {
      _status = MeasurementStatus.complete;
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
    _pipelineResult = null;
    setState(() {
      _status = MeasurementStatus.idle;
      _measureProgress = 0.0;
      _phaseText = '파동 안정화 중...';
    });
  }

  static String _riskLevelKo(String? level) {
    return switch (level) {
      'normal' => '정상',
      'caution' => '주의',
      'warning' => '경고',
      'critical' => '위험',
      _ => '--',
    };
  }

  @override
  void dispose() {
    _disposeWaveAnimation();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isMeasuring = _status == MeasurementStatus.measuring;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          tooltip: '뒤로 가기',
          onPressed: () => context.pop(),
        ),
        title: const Text('측정', style: TextStyle(
          color: SanggamTheme.primary,
          fontWeight: FontWeight.bold,
        )),
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
                    child: _buildStatusWidget(),
                  ),
                ),

                // 하단 버튼
                if (_status == MeasurementStatus.idle)
                  Column(
                    children: [
                      Semantics(
                        button: true,
                        label: '건강 측정 시작 버튼',
                        child: PrimaryButton(
                          text: '측정 시작',
                          icon: Icons.play_arrow_rounded,
                          onPressed: _startMeasurement,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Semantics(
                        button: true,
                        label: 'NFC 카트리지 읽기 버튼',
                        child: OutlinedButton.icon(
                          onPressed: _readCartridge,
                          icon: const Icon(Icons.nfc),
                          label: Text(_cartridgeId != null ? '카트리지 읽음' : 'NFC 카트리지 읽기'),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 48),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                        ),
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
                          context.push('/measure/result');
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

  Widget _buildStatusWidget() {
    switch (_status) {
      case MeasurementStatus.idle:
        return const Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.sensors_rounded, size: 72, color: SanggamTheme.primary),
            SizedBox(height: 24),
            Text(
              '디바이스를 준비해주세요',
              style: TextStyle(
                color: Colors.white,
                fontSize: 22,
                fontWeight: FontWeight.bold,
              ),
            ),
            SizedBox(height: 8),
            Text(
              '카트리지를 장착하고\n측정 버튼을 눌러주세요',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
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
                          child: const Center(
                            child: Icon(
                              Icons.bluetooth_searching_rounded,
                              size: 48,
                              color: SanggamTheme.jagaeCyan,
                            ),
                          ),
                        );
                      },
                    )
                  : const CircularProgressIndicator(
                      strokeWidth: 4,
                      color: SanggamTheme.primary,
                    ),
            ),
            const SizedBox(height: 24),
            const Text(
              '디바이스 연결 중...',
              style: TextStyle(
                color: Colors.white,
                fontSize: 22,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              'BLE 연결을 시도하고 있습니다',
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
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
                            waveColor: SanggamTheme.jagaeCyan,
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
                style: const TextStyle(
                  color: SanggamTheme.primary,
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
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
                backgroundColor: SanggamTheme.surfaceVariant.withValues(alpha: 0.3),
                valueColor: AlwaysStoppedAnimation<Color>(
                  Color.lerp(SanggamTheme.jagaeCyan, SanggamTheme.primary, _measureProgress)!,
                ),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${(_measureProgress * 100).toInt()}%',
              style: const TextStyle(
                fontFamily: 'JetBrains Mono',
                color: SanggamTheme.jagaeCyan,
                fontSize: 12,
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
                color: SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.check_circle_rounded,
                size: 72,
                color: SanggamTheme.jagaeCyan,
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              '측정 완료!',
              style: TextStyle(
                color: SanggamTheme.jagaeCyan,
                fontSize: 28,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            // 실측정 결과
            Container(
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.05),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: SanggamTheme.surfaceVariant),
              ),
              padding: const EdgeInsets.all(24),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _buildResultItem(
                    _pipelineResult?.measurement.primaryValue.toStringAsFixed(1) ?? '--',
                    _pipelineResult?.measurement.unit ?? 'mg/dL',
                    '혈당',
                  ),
                  Container(
                    width: 1,
                    height: 40,
                    color: SanggamTheme.surfaceVariant,
                  ),
                  _buildResultItem(
                    _riskLevelKo(_pipelineResult?.analysis.riskLevel),
                    '',
                    '판정',
                  ),
                  Container(
                    width: 1,
                    height: 40,
                    color: SanggamTheme.surfaceVariant,
                  ),
                  _buildResultItem(
                    '${(_pipelineResult?.measurement.confidence ?? 0.0 * 100).toStringAsFixed(0)}%',
                    '',
                    '신뢰도',
                  ),
                ],
              ),
            ),
          ],
        );

      case MeasurementStatus.error:
        return const Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.error_outline, size: 72, color: SanggamTheme.error),
            SizedBox(height: 24),
            Text(
              '측정 실패',
              style: TextStyle(
                color: SanggamTheme.error,
                fontSize: 22,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        );
    }
  }

  Widget _buildResultItem(String value, String unit, String label) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
        if (unit.isNotEmpty)
          Text(
            unit,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 11,
            ),
          ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          ),
        ),
      ],
    );
  }
}
