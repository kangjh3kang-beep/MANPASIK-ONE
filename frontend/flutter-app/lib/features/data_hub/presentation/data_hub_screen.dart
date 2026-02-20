import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/widgets/holo_globe.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/wave_analysis_panel.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/jarvis_connector.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

class DataHubScreen extends ConsumerStatefulWidget {
  const DataHubScreen({super.key});

  @override
  ConsumerState<DataHubScreen> createState() => _DataHubScreenState();
}

class _DataHubScreenState extends ConsumerState<DataHubScreen> with TickerProviderStateMixin {
  bool _showMyZone = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    
    // API 데이터 바인딩
    final historyAsync = ref.watch(measurementHistoryProvider);
    final biomarkersAsync = ref.watch(biomarkerSummariesProvider);
    // Even if we don't display nodes here, we might want to know count or status eventually, 
    // but for static panels we mostly need system stats.
    
    // 측정 통계 계산
    final totalMeasurements = historyAsync.valueOrNull?.totalCount ?? 0;
    final biomarkers = biomarkersAsync.valueOrNull ?? [];
    
    final healthScore = biomarkers.isNotEmpty
        ? (biomarkers.where((b) => b.latestValue != null && b.latestValue! >= b.referenceMin && b.latestValue! <= b.referenceMax).length / biomarkers.length * 100).toInt()
        : 85;

    // Dynamic Colors
    final contentColor = isDark ? Colors.white : const Color(0xFF1A1A1A);
    final subContentColor = isDark ? Colors.white70 : const Color(0xFF424242);
    final globeColor = isDark ? AppTheme.waveCyan : const Color(0xFF2C3E50);

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: SafeArea(
        child: LayoutBuilder(
          builder: (context, constraints) {
            // Adaptive Layout Logic
            final isCompact = constraints.maxWidth < 900 || constraints.maxHeight < 600;
            
            // Fix: Increase minimum height to prevent internal overflow of panels
            // Chrome (44) + Text(35) + Footer(31) = ~110px minimum needed + content
            final double panelWidth = isCompact ? 170.0 : 220.0;
            final double panelHeight = isCompact ? 160.0 : 180.0;
            
            // Ensure minimum canvas size to prevent overlapping/crushing
            final double minCanvasWidth = 600; 
            final double minCanvasHeight = 500;
            
            // Use the larger of actual constraint or minimum size
            // This enables scrolling if the window is too small
            final double canvasWidth = math.max(constraints.maxWidth, minCanvasWidth);
            final double canvasHeight = math.max(constraints.maxHeight, minCanvasHeight);
            
            final center = Offset(canvasWidth / 2, canvasHeight / 2);

            // Dynamic Spread - Scale with available space but clamp to safe limits
            // spreadY determines vertical distance from center. 
            // panelHeight is 160. So center +/- spreadY must leave 80px space? No. 
            // Position is top-left of panel? No, calculated: leftTopPos = center + Offset(..., -spreadY)
            // So spreadY is distance from center Y.
            final double spreadX = isCompact ? canvasWidth * 0.25 : canvasWidth * 0.22;
            final double spreadY = isCompact ? canvasHeight * 0.22 : 160.0; 

            // Calculate positions
            final leftTopPos = center + Offset(-spreadX - panelWidth, -spreadY - panelHeight/2);
            final rightTopPos = center + Offset(spreadX, -spreadY - panelHeight/2); 
            final leftBotPos = center + Offset(-spreadX - panelWidth, spreadY - panelHeight/2);
            final rightBotPos = center + Offset(spreadX, spreadY - panelHeight/2);

            return SingleChildScrollView(
              scrollDirection: Axis.vertical,
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Container(
                  width: canvasWidth,
                  height: canvasHeight,
                  alignment: Alignment.center,
                  child: Stack(
                      alignment: Alignment.center,
                      children: [
                        // 1. Central Decor & Globe
                         Container(
                           width: canvasWidth * 0.6, 
                           height: canvasWidth * 0.6,
                           constraints: const BoxConstraints(maxHeight: 500, maxWidth: 500),
                           decoration: BoxDecoration(
                             shape: BoxShape.circle,
                             gradient: RadialGradient(
                               colors: isDark 
                                 ? [AppTheme.sanggamGold.withOpacity(0.15), Colors.transparent]
                                 : [const Color(0xFF00796B).withOpacity(0.1), Colors.transparent],
                             ),
                           ),
                         ),

                        HoloGlobe(
                          size: isCompact ? 220 : 300,
                          color: isDark ? AppTheme.sanggamGold : const Color(0xFF004D40),
                          accentColor: isDark ? Colors.white : const Color(0xFF00796B),
                        ),

                        // 2. Jarvis Connectors (Static Panels) - Adjusted for new positions
                        JarvisConnector(startOffset: center, endOffset: leftTopPos + Offset(panelWidth, panelHeight/2)),
                        JarvisConnector(startOffset: center, endOffset: rightTopPos + Offset(0, panelHeight/2)),
                        JarvisConnector(startOffset: center, endOffset: leftBotPos + Offset(panelWidth, panelHeight/2)),
                        JarvisConnector(startOffset: center, endOffset: rightBotPos + Offset(0, panelHeight/2)),

                        // 3. Static HUD Panels (4 Panels)
                        
                        // Top Left: Environment Data
                        Positioned(
                          left: leftTopPos.dx, top: leftTopPos.dy,
                          child: SizedBox(
                            width: panelWidth, height: panelHeight,
                            child: WaveAnalysisPanel(
                              title: '수질/환경 데이터',
                              child: _buildRealTimeChart(globeColor, [1,2,3,4,5,4,3,2,1]),
                              footer: _buildFooterStat(
                                '정상 범위',
                                '안전',
                                contentColor, subContentColor,
                              ),
                            ),
                          ),
                        ),

                        // Top Right: System Health
                        Positioned(
                          left: rightTopPos.dx, top: rightTopPos.dy,
                          child: SizedBox(
                            width: panelWidth, height: panelHeight,
                            child: WaveAnalysisPanel(
                              title: '시스템 상태',
                              child: _buildHealthGauge(globeColor, healthScore),
                              footer: _buildFooterStat(
                                 healthScore > 80 ? "최적" : "주의",
                                'Core: ${healthScore}%',
                                contentColor, subContentColor,
                              ),
                            ),
                          ),
                        ),

                        // Bottom Left: AI Model
                        Positioned(
                          left: leftBotPos.dx, top: leftBotPos.dy,
                          child: SizedBox(
                            width: panelWidth, height: panelHeight,
                            child: WaveAnalysisPanel(
                              title: 'AI 예측 모델',
                              child: _buildHexStructure(globeColor, biomarkers.length),
                              footer: _buildFooterStat(
                                '활성 노드',
                                '${biomarkers.length}개',
                                contentColor, subContentColor,
                              ),
                            ),
                          ),
                        ),

                        // Bottom Right: Data Vault
                        Positioned(
                          left: rightBotPos.dx, top: rightBotPos.dy,
                          child: SizedBox(
                            width: panelWidth, height: panelHeight,
                            child: WaveAnalysisPanel(
                              title: '보안 데이터',
                              child: _buildTreasureChest(globeColor),
                              footer: _buildFooterStat(
                                '양자 암호화',
                                '${totalMeasurements}건',
                                contentColor, subContentColor,
                              ),
                            ),
                          ),
                        ),

                         // 5. Header (Floating)
                         Positioned(
                           top: 20,
                           child: AnimateFadeInUp(
                             child: Column(
                               children: [
                                 Text(
                                   '만파식 데이터 허브',
                                   style: theme.textTheme.titleMedium?.copyWith(
                                     color: AppTheme.sanggamGold,
                                     fontWeight: FontWeight.bold,
                                     letterSpacing: 2.0,
                                     shadows: [Shadow(color: AppTheme.sanggamGold, blurRadius: 10)],
                                   ),
                                 ),
                                 Text(
                                   '글로벌 파동 모니터링 시스템',
                                   style: TextStyle(color: subContentColor, fontSize: 10, letterSpacing: 1.0),
                                 ),
                                 if (isCompact)
                                  Padding(
                                    padding: const EdgeInsets.only(top: 8),
                                    child: Text(
                                      '창을 키워 전체 화면으로 확인하세요',
                                      style: TextStyle(color: Colors.redAccent.withOpacity(0.7), fontSize: 10),
                                    ),
                                  )
                               ],
                             ),
                           ),
                         ),
                         
                         // 6. My Zone Toggle
                         Positioned(
                           top: 60, right: 20,
                           child: GestureDetector(
                             onTap: () => setState(() => _showMyZone = !_showMyZone),
                             child: AnimatedContainer(
                               duration: const Duration(milliseconds: 300),
                               padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                               decoration: BoxDecoration(
                                 color: _showMyZone ? AppTheme.sanggamGold.withOpacity(0.2) : Colors.white.withOpacity(0.05),
                                 borderRadius: BorderRadius.circular(16),
                                 border: Border.all(color: _showMyZone ? AppTheme.sanggamGold : Colors.white24),
                               ),
                               child: Row(
                                 mainAxisSize: MainAxisSize.min,
                                 children: [
                                   Icon(Icons.my_location, size: 14, color: _showMyZone ? AppTheme.sanggamGold : Colors.white54),
                                   const SizedBox(width: 4),
                                   Text('My Zone', style: TextStyle(fontSize: 10, color: _showMyZone ? AppTheme.sanggamGold : Colors.white54, fontWeight: _showMyZone ? FontWeight.bold : FontWeight.normal)),
                                 ],
                               ),
                             ),
                           ),
                         ),

                         // 7. My Zone Overlay (개인 기준선 범위)
                         if (_showMyZone)
                           Positioned(
                             bottom: 60, left: 20, right: 20,
                             child: AnimatedOpacity(
                               opacity: _showMyZone ? 1.0 : 0.0,
                               duration: const Duration(milliseconds: 300),
                               child: Container(
                                 padding: const EdgeInsets.all(12),
                                 decoration: BoxDecoration(
                                   color: const Color(0xFF0A1020).withOpacity(0.85),
                                   borderRadius: BorderRadius.circular(16),
                                   border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.3)),
                                 ),
                                 child: Column(
                                   crossAxisAlignment: CrossAxisAlignment.start,
                                   children: [
                                     Text('내 기준선 범위', style: TextStyle(color: AppTheme.sanggamGold, fontSize: 12, fontWeight: FontWeight.bold)),
                                     const SizedBox(height: 8),
                                     _buildMyZoneRow('혈당', '85-110', 'mg/dL', healthScore > 80),
                                     _buildMyZoneRow('SpO2', '96-99', '%', true),
                                     _buildMyZoneRow('심박', '60-80', 'bpm', healthScore > 70),
                                     _buildMyZoneRow('체온', '36.2-36.8', '\u00b0C', true),
                                   ],
                                 ),
                               ),
                             ),
                           ),

                         // 8. Monitoring Shortcut (Center Bottom)
                         Positioned(
                            bottom: 20,
                            child: _buildMonitoringButton(context, isDark),
                         )
                      ],
                    ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildMonitoringButton(BuildContext context, bool isDark) {
    return InkWell(
        onTap: () => context.push('/data/monitoring'),
        borderRadius: BorderRadius.circular(30),
        child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            decoration: BoxDecoration(
                color: isDark ? const Color(0xFF1A1F35) : Colors.white,
                borderRadius: BorderRadius.circular(30),
                border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.5)),
                boxShadow: const [BoxShadow(color: Colors.black26, blurRadius: 10, offset: Offset(0, 4))],
            ),
            child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                    const Icon(Icons.grid_view_rounded, color: AppTheme.sanggamGold, size: 18),
                    const SizedBox(width: 8),
                    Text('전체 리더기 현황 보기', style: TextStyle(color: isDark ? Colors.white : Colors.black87, fontWeight: FontWeight.bold)),
                ],
            ),
        ),
    );
  }

  Widget _buildHexStructure(Color color, int biomarkerCount) {
    return Center(
      child: Stack(
        alignment: Alignment.center,
        children: [
          TweenAnimationBuilder(
            tween: Tween<double>(begin: 0.8, end: 1.2),
            duration: const Duration(seconds: 2),
            builder: (context, value, child) {
               return Transform.scale(
                 scale: value,
                 child: Icon(Icons.hexagon_outlined, color: AppTheme.sanggamGold.withOpacity(0.1), size: 100),
               );
            },
            curve: Curves.easeInOut,
          ),
          Icon(Icons.hexagon_outlined, color: AppTheme.sanggamGold.withOpacity(0.6), size: 80), // 60%
          Icon(Icons.data_object, color: color.withOpacity(0.6), size: 30), // 60%
          Positioned(
             top: 20,
             child: Container(width: 40, height: 1, color: AppTheme.sanggamGold.withOpacity(0.6)), // 60%
          ),
           Positioned(
            bottom: 10,
            child: Text(
              biomarkerCount > 0 ? '$biomarkerCount 노드 활성' : 'AI 노드 활성',
              style: TextStyle(color: AppTheme.sanggamGold, fontSize: 6, letterSpacing: 1.0),
            ),
          )
        ],
      ),
    );
  }

  Widget _buildTreasureChest(Color color) {
    return Center(
      child: FittedBox(
        fit: BoxFit.scaleDown,
        child: Stack(
          alignment: Alignment.center,
          children: [
            Icon(Icons.shield_outlined, color: color.withOpacity(0.2), size: 70),
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.lock_rounded, color: AppTheme.sanggamGold.withOpacity(0.6), size: 32), // 60%
                const SizedBox(height: 4),
                Text('양자 암호화됨', style: TextStyle(color: color, fontSize: 6, fontWeight: FontWeight.bold)),
              ],
            ),
             TweenAnimationBuilder(
              tween: Tween<double>(begin: 0, end: 2 * math.pi),
              duration: const Duration(seconds: 10),
              builder: (context, value, child) {
                 return Transform.rotate(
                   angle: value,
                   child: Container(
                     width: 50, height: 50,
                     decoration: BoxDecoration(
                       border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.6), width: 1, style: BorderStyle.solid), // 60%
                       shape: BoxShape.circle,
                     ),
                   ),
                 );
              },
               onEnd: () {},
               curve: Curves.linear,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRealTimeChart(Color color, List<dynamic> items) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
      child: Center(
        child: CustomPaint(
          size: const Size(double.infinity, 60),
          painter: _ChartPainter(color: color, dataPoints: items.length),
        ),
      ),
    );
  }

  Widget _buildHealthGauge(Color color, int score) {
    final normalizedScore = score.clamp(0, 100) / 100.0;
    return Center(
      child: Stack(
        alignment: Alignment.center,
        children: [
           CircularProgressIndicator(
              value: 1.0,
              color: AppTheme.sanggamGold.withOpacity(0.1),
              strokeWidth: 8,
            ),
          TweenAnimationBuilder(
            tween: Tween<double>(begin: 0, end: normalizedScore),
            duration: const Duration(milliseconds: 1500),
            builder: (context, value, child) {
              return CircularProgressIndicator(
                value: value,
                color: color,
                strokeWidth: 4,
                backgroundColor: Colors.transparent,
              );
            },
          ),
          Container(
             width: 50, height: 50,
             decoration: BoxDecoration(
               shape: BoxShape.circle,
               border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.5), width: 1),
               gradient: RadialGradient(colors: [color.withOpacity(0.2), Colors.transparent]),
             ),
             child: Center(
               child: Text(
                '$score',
                style: TextStyle(
                  fontSize: 16,
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  shadows: [Shadow(color: color, blurRadius: 10)],
                ),
               ),
             ),
          ),
        ],
      ),
    );
  }

  Widget _buildMyZoneRow(String label, String range, String unit, bool inRange) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 3),
      child: Row(
        children: [
          Icon(inRange ? Icons.check_circle : Icons.warning_amber, size: 14, color: inRange ? Colors.greenAccent : Colors.amberAccent),
          const SizedBox(width: 8),
          SizedBox(width: 40, child: Text(label, style: const TextStyle(color: Colors.white70, fontSize: 11))),
          const SizedBox(width: 8),
          Text(range, style: const TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.bold)),
          const SizedBox(width: 4),
          Text(unit, style: const TextStyle(color: Colors.white54, fontSize: 10)),
        ],
      ),
    );
  }

  Widget _buildFooterStat(String label, String value, Color textColor, Color labelColor) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: TextStyle(color: labelColor, fontSize: 10)),
        Text(value, style: TextStyle(color: textColor, fontWeight: FontWeight.bold, fontSize: 12)),
      ],
    );
  }
}

class _ChartPainter extends CustomPainter {
  final Color color;
  final int dataPoints;
  _ChartPainter({required this.color, this.dataPoints = 0});

  @override
  void paint(Canvas canvas, Size size) {
    final gridPaint = Paint()..color = Colors.white.withOpacity(0.05)..strokeWidth = 0.5;
    final double step = size.width / 10;
    for(double x=0; x<=size.width; x+=step) canvas.drawLine(Offset(x,0), Offset(x, size.height), gridPaint);
    for(double y=0; y<=size.height; y+=10) canvas.drawLine(Offset(0,y), Offset(size.width, y), gridPaint);

    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2
      ..shader = LinearGradient(
        colors: [color.withOpacity(0), color, color, color.withOpacity(0)],
        stops: const [0.0, 0.2, 0.8, 1.0],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));

    final path = Path();
    path.moveTo(0, size.height * 0.5);

    // Amplitude varies with data points count
    final amplitude = dataPoints > 0 ? 10.0 + (dataPoints * 0.5).clamp(0, 15) : 10.0;
    for (double x = 0; x <= size.width; x+=2) {
      final y = size.height * 0.5 +
                math.sin(x * 0.1) * amplitude +
                math.sin(x * 0.5) * 5;
      path.lineTo(x, y);
    }
    canvas.drawPath(path, paint);

    final fillPaint = Paint()
      ..style = PaintingStyle.fill
      ..shader = LinearGradient(
        begin: Alignment.topCenter, end: Alignment.bottomCenter,
        colors: [color.withOpacity(0.2), Colors.transparent],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));

    final fillPath = Path.from(path)
      ..lineTo(size.width, size.height)
      ..lineTo(0, size.height)
      ..close();

    canvas.drawPath(fillPath, fillPaint);
  }
  @override
  bool shouldRepaint(covariant _ChartPainter oldDelegate) => oldDelegate.dataPoints != dataPoints;
}
