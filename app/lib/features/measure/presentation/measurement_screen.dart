import 'dart:async';
import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';
import '../../../core/ffi/rust_bridge.dart';
import '../../../core/state/global_providers.dart';
import '../../../core/theme/app_theme.dart';
import '../../../core/database/local_measurement.dart';
import '../data/measurement_local_repository.dart';

class MeasurementScreen extends ConsumerStatefulWidget {
  const MeasurementScreen({super.key});

  @override
  ConsumerState<MeasurementScreen> createState() => _MeasurementScreenState();
}

class _MeasurementScreenState extends ConsumerState<MeasurementScreen> with SingleTickerProviderStateMixin {
  final List<FlSpot> _diffDataPoints = [];
  Timer? _bleStreamTimer;
  double _currentTimeMs = 0;
  
  late AnimationController _pulseController;
  late Animation<double> _pulseAnimation;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(vsync: this, duration: const Duration(seconds: 2))..repeat(reverse: true);
    _pulseAnimation = Tween<double>(begin: 0.8, end: 1.2).animate(CurvedAnimation(parent: _pulseController, curve: Curves.easeInOut));
    
    _startSimulatedBleStream();
  }

  void _startSimulatedBleStream() {
    _bleStreamTimer = Timer.periodic(const Duration(milliseconds: 50), (timer) async {
      _currentTimeMs += 50;
      double sDet = 1.0 + sin(_currentTimeMs / 300) * 0.7 + (Random().nextDouble() * 0.1);
      double sRef = 0.5 + sin(_currentTimeMs / 300) * 0.3;
      double sDiff = await rustEngine.applyDifferential(sDet, sRef);
      
      if (mounted) {
        setState(() {
          _diffDataPoints.add(FlSpot(_currentTimeMs / 1000, sDiff));
          if (_diffDataPoints.length > 150) _diffDataPoints.removeAt(0);
        });
      }

      if (timer.tick >= 150) { 
        timer.cancel();
        _runRustInference();
      }
    });
  }

  Future<void> _runRustInference() async {
    ref.read(measurementStateProvider.notifier).state = MeasurementState(phase: MeasurementPhase.INFERENCE);
    
    // Rust 코어 통합(FFI)에 차동 시그널 전송 및 분석 
    final resultText = await rustEngine.runInference(_diffDataPoints.map((e) => e.y).toList(), 25.0);
    
    // [Phase 3] 분석 완료 직후, Isar 로컬 DB에 CRDT 동기화를 위한 오프라인 캐싱
    final isarRepo = await ref.read(measurementLocalRepoProvider.future);
    final localItem = LocalMeasurement()
      ..clientLocalId = 'meas_${DateTime.now().millisecondsSinceEpoch}'
      ..deviceMac = 'CSI-V1-MAC-0001'
      ..measuredAt = DateTime.now()
      ..diffSignal = _diffDataPoints.map((e) => e.y).toList()
      ..fingerprint = List.filled(896, 0.0) // TODO: 실제 FFI 병합 피처
      ..healthScore = 92
      ..riskLabel = 'Normal'
      ..isSynced = false
      ..updatedAt = DateTime.now();

    await isarRepo.saveMeasurement(localItem);

    if (mounted) {
      ref.read(measurementStateProvider.notifier).state = MeasurementState(
        phase: MeasurementPhase.COMPLETED, 
        errorMessage: '로컬 DB 적재 & 예측 완료\n$resultText'
      );
    }
  }

  @override
  void dispose() {
    _bleStreamTimer?.cancel();
    _pulseController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final curState = ref.watch(measurementStateProvider);

    return PopScope(
      canPop: false, 
      child: Scaffold(
        backgroundColor: AppTheme.backgroundDark,
        appBar: AppBar(
          title: const Text('M-040 분석 진행중', style: TextStyle(fontFamily: 'Inter')),
          automaticallyImplyLeading: false,
          actions: [
            Container(
              margin: const EdgeInsets.only(right: 16),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(color: AppTheme.surfaceDark, borderRadius: BorderRadius.circular(20), border: Border.all(color: AppTheme.primaryNeonTeal.withValues(alpha: 0.5))),
              child: const Row(
                children: [
                  Icon(Icons.lock_outline, size: 14, color: AppTheme.primaryNeonTeal),
                  SizedBox(width: 6),
                  Text('보호됨', style: TextStyle(color: AppTheme.primaryNeonTeal, fontSize: 12, fontWeight: FontWeight.bold)),
                ],
              ),
            ),
          ],
        ),
        body: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              const SizedBox(height: 40),
              // 중앙 AI 분석 상태 (Pulsing)
              ScaleTransition(
                scale: curState.phase == MeasurementPhase.COMPLETED ? const AlwaysStoppedAnimation(1.0) : _pulseAnimation,
                child: Container(
                  width: 120, height: 120,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: AppTheme.premiumGradient,
                    boxShadow: [
                      BoxShadow(color: AppTheme.primaryNeonPurple.withValues(alpha: 0.6), blurRadius: 30, spreadRadius: 10)
                    ]
                  ),
                  child: Center(
                    child: Icon(
                      curState.phase == MeasurementPhase.INFERENCE ? Icons.auto_awesome 
                        : (curState.phase == MeasurementPhase.COMPLETED ? Icons.check : Icons.graphic_eq),
                      color: Colors.white, size: 50,
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 40),

              Text(
                curState.phase == MeasurementPhase.INFERENCE 
                  ? 'Edge Neural Network 분석 중...' 
                  : curState.phase == MeasurementPhase.COMPLETED
                      ? '분석 완료'
                      : 'CSI v1.0 16-pin 실시간 차동 측정 연결됨',
                style: const TextStyle(fontSize: 18, color: Colors.white, fontWeight: FontWeight.w600, letterSpacing: 0.5),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 10),
              
              if (curState.phase == MeasurementPhase.COMPLETED && curState.errorMessage != null)
                Column(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(color: AppTheme.primaryNeonTeal.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12), border: Border.all(color: AppTheme.primaryNeonTeal)),
                      child: Text(curState.errorMessage!, style: const TextStyle(color: AppTheme.primaryNeonTeal, fontSize: 16, fontWeight: FontWeight.bold), textAlign: TextAlign.center),
                    ),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppTheme.primaryNeonTeal,
                        foregroundColor: AppTheme.backgroundDark,
                        padding: const EdgeInsets.symmetric(horizontal: 40, vertical: 16),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(30)),
                        elevation: 8,
                      ),
                      onPressed: () => context.go('/home'), 
                      child: const Text('종합 스코어 보기 (홈으로)', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    )
                  ],
                ),

              const SizedBox(height: 40),

              // 고품질 차트 (그라데이션 및 곡선)
              Expanded(
                child: Container(
                  decoration: BoxDecoration(
                    color: AppTheme.surfaceDark,
                    borderRadius: BorderRadius.circular(24),
                    border: Border.all(color: Colors.white.withValues(alpha: 0.05)),
                    boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.3), blurRadius: 20)],
                  ),
                  padding: const EdgeInsets.fromLTRB(20, 30, 20, 10),
                  child: LineChart(
                    LineChartData(
                      lineBarsData: [
                        LineChartBarData(
                          spots: _diffDataPoints,
                          isCurved: true,
                          curveSmoothness: 0.35,
                          gradient: const LinearGradient(colors: [AppTheme.primaryNeonTeal, AppTheme.primaryNeonPurple]),
                          barWidth: 4,
                          isStrokeCapRound: true,
                          dotData: const FlDotData(show: false),
                          belowBarData: BarAreaData(
                            show: true,
                            gradient: LinearGradient(
                              colors: [AppTheme.primaryNeonTeal.withValues(alpha: 0.3), AppTheme.primaryNeonPurple.withValues(alpha: 0.0)],
                              begin: Alignment.topCenter, end: Alignment.bottomCenter,
                            ),
                          ),
                        )
                      ],
                      titlesData: const FlTitlesData(show: false),
                      gridData: FlGridData(
                        show: true, 
                        drawVerticalLine: false, 
                        horizontalInterval: 0.5,
                        getDrawingHorizontalLine: (value) => FlLine(color: Colors.white.withValues(alpha: 0.05), strokeWidth: 1),
                      ),
                      borderData: FlBorderData(show: false),
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
