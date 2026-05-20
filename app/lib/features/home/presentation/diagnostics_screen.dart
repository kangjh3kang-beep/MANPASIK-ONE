import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_theme.dart';
// import '../../../core/ffi/rust_bridge.dart'; // 실 연동 시 해제

/// 실시간 레이어별 건강 상태를 모의로 스트리밍 (Rust 진단 모듈 연동용)
final systemHealthStreamProvider = StreamProvider<List<double>>((ref) async* {
  // [HW, FW, Rust, App, Backend, AI] — 1.0 (최상) 기준
  yield [0.85, 0.92, 0.99, 1.00, 0.88, 0.95];
});

final healingLogsProvider = Provider<List<String>>((ref) {
  return [
    '[06:40] [INFO] FW 센서 C-31 노이즈 보정 완료 (HwAutoRecovery)',
    '[06:38] [WARN] AI 가중치 동기화 지연 (오프라인 상태 캐싱 중)',
    '[06:30] [ERROR] Bluetooth CRC 실패 (재전송 요청됨)',
    '[06:15] [INFO] Rust 차동측정 엔진 Calibration 적용',
  ];
});

class DiagnosticsScreen extends ConsumerWidget {
  const DiagnosticsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final healthAsync = ref.watch(systemHealthStreamProvider);
    final logs = ref.watch(healingLogsProvider);

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text('H-035 자가진단 센터'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        flexibleSpace: ClipRect(
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
            child: Container(color: AppTheme.backgroundDark.withValues(alpha: 0.5)),
          ),
        ),
      ),
      body: Stack(
        children: [
          // 배경 글로우
          Positioned(
            top: 50, right: -50,
            child: Container(
              width: 300, height: 300,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: AppTheme.dangerRed.withValues(alpha: 0.1), // 진단 화면 특유의 경고색 혼합
              ),
              child: BackdropFilter(filter: ImageFilter.blur(sigmaX: 100, sigmaY: 100), child: Container()),
            ),
          ),

          SafeArea(
            child: CustomScrollView(
              physics: const BouncingScrollPhysics(),
              slivers: [
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 20),
                    child: _buildRadarChart(healthAsync),
                  ),
                ),
                
                const SliverPadding(
                  padding: EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                  sliver: SliverToBoxAdapter(
                    child: Text(
                      "실시간 자가치유 (Self-Healing) 로그", 
                      style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white, letterSpacing: 0.5),
                    ),
                  ),
                ),

                SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, index) {
                      return _buildLogCard(logs[index]);
                    },
                    childCount: logs.length,
                  ),
                ),
                const SliverPadding(padding: EdgeInsets.only(bottom: 60)),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRadarChart(AsyncValue<List<double>> healthAsync) {
    return Container(
      height: 350,
      decoration: BoxDecoration(
        color: AppTheme.surfaceDark,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: Colors.white.withValues(alpha: 0.05)),
        boxShadow: [
          BoxShadow(color: Colors.black.withValues(alpha: 0.3), blurRadius: 20),
        ],
      ),
      padding: const EdgeInsets.all(20),
      child: Column(
        children: [
          const Text('6-Layer System Integrity', style: TextStyle(color: Colors.white54, fontSize: 16)),
          const SizedBox(height: 20),
          Expanded(
            child: healthAsync.when(
              data: (values) {
                return RadarChart(
                  RadarChartData(
                    tickCount: 3,
                    ticksTextStyle: const TextStyle(color: Colors.transparent),
                    titlePositionPercentageOffset: 0.15,
                    getTitle: (index, angle) {
                      final layers = ['HW', 'FW', 'Rust', 'App', 'Go', 'AI'];
                      return RadarChartTitle(text: layers[index], angle: angle);
                    },
                    dataSets: [
                      RadarDataSet(
                        fillColor: AppTheme.primaryNeonPurple.withValues(alpha: 0.2),
                        borderColor: AppTheme.primaryNeonPurple,
                        entryRadius: 4,
                        dataEntries: values.map((v) => RadarEntry(value: v)).toList(),
                        borderWidth: 2,
                      )
                    ],
                    radarBorderData: const BorderSide(color: Colors.white12),
                    tickBorderData: const BorderSide(color: Colors.white12),
                    gridBorderData: const BorderSide(color: Colors.white12, width: 2),
                  ),
                  swapAnimationDuration: const Duration(milliseconds: 400),
                );
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (e, st) => Center(child: Text('Error: $e', style: const TextStyle(color: AppTheme.dangerRed))),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLogCard(String logText) {
    bool isError = logText.contains('[ERROR]');
    bool isWarn = logText.contains('[WARN]');
    
    Color accentColor = isError ? AppTheme.dangerRed : (isWarn ? AppTheme.warningOrange : AppTheme.primaryNeonTeal);

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 6),
      decoration: BoxDecoration(
        color: AppTheme.surfaceDark.withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: accentColor.withValues(alpha: 0.3)),
      ),
      child: ListTile(
        leading: Icon(
          isError ? Icons.error_outline : (isWarn ? Icons.warning_amber_rounded : Icons.healing),
          color: accentColor,
        ),
        title: Text(
          logText,
          style: TextStyle(color: Colors.white.withValues(alpha: 0.8), fontSize: 13, height: 1.4),
        ),
      ),
    );
  }
}
