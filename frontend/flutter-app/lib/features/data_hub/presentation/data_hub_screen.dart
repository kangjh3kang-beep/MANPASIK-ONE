import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/holo_globe.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/wave_analysis_panel.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';

class DataHubScreen extends ConsumerStatefulWidget {
  const DataHubScreen({super.key});

  @override
  ConsumerState<DataHubScreen> createState() => _DataHubScreenState();
}

class _DataHubScreenState extends ConsumerState<DataHubScreen> with SingleTickerProviderStateMixin {
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final size = MediaQuery.of(context).size;
    final isCompact = size.width < 600;
    
    // Dynamic Colors for Theme
    final contentColor = isDark ? Colors.white : const Color(0xFF1A1A1A); // Ink Black for Light
    final subContentColor = isDark ? Colors.white70 : const Color(0xFF424242);
    final globeColor = isDark ? AppTheme.waveCyan : const Color(0xFF2C3E50); // Dark Blue/Ink for Globe
    final glowColor = isDark ? AppTheme.waveCyan : Colors.transparent; // No glow in light mode
    final shadowColor = isDark ? AppTheme.waveCyan : Colors.black.withOpacity(0.1); 

    return Scaffold(
      backgroundColor: Colors.transparent, 
      body: SafeArea(
        child: Stack(
          alignment: Alignment.center,
          children: [
              // 2. Dashboard Interface Layer
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                children: [
                  // Title Area (Top)
                  AnimateFadeInUp(
                    child: Column(
                      children: [
                        Text(
                          '만파식 시스템',
                          style: theme.textTheme.headlineMedium?.copyWith(
                            color: contentColor, 
                            fontWeight: FontWeight.bold,
                            letterSpacing: 1.5,
                            shadows: [
                              Shadow(color: shadowColor, blurRadius: 15), 
                            ],
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '글로벌 파동 모니터링',
                            style: theme.textTheme.labelSmall?.copyWith(
                              color: subContentColor,
                              letterSpacing: 1.2,
                              fontWeight: FontWeight.bold, 
                              shadows: isDark ? [const Shadow(color: Colors.black, blurRadius: 4, offset: Offset(1,1))] : null,
                            ),
                        ),
                      ],
                    ),
                  ),

                  const Spacer(),

                  // MIDDLE SECTION: Globe + Panels
                  Expanded(
                    flex: 10,
                    child: Stack(
                      alignment: Alignment.center,
                      children: [
                        // 1. Central Halo (Large Circular Glow behind everything)
                        // 1. Central Halo (Theme Specific Background)
                         Container(
                           width: 500,
                           height: 500,
                           decoration: BoxDecoration(
                             shape: BoxShape.circle,
                             gradient: RadialGradient(
                               center: Alignment.center,
                               radius: 0.7,
                               colors: isDark
                                   ? [
                                       AppTheme.sanggamGold.withOpacity(0.25), // Golden Blur Center
                                       AppTheme.sanggamGold.withOpacity(0.0),  // Fade out
                                     ]
                                   : [
                                       const Color(0xFF00796B).withOpacity(0.15), // Cheongnok (Teal) Transparent Base
                                       const Color(0xFF004D40).withOpacity(0.0),  // Fade out
                                     ],
                               stops: const [0.0, 1.0],
                             ),
                             boxShadow: isDark
                                 ? [
                                     BoxShadow(
                                       color: AppTheme.sanggamGold.withOpacity(0.1),
                                       blurRadius: 60,
                                       spreadRadius: 10,
                                     )
                                   ]
                                 : [
                                     // White Mode Gloss/Glass Effect
                                     BoxShadow(
                                       color: const Color(0xFF00796B).withOpacity(0.1),
                                       blurRadius: 40,
                                       spreadRadius: 5,
                                     ),
                                   ],
                           ),
                         ),
                         
                         // White Mode Gloss Highlight (Jadae/Jade feel)
                         if (!isDark)
                           Positioned(
                             top: 100,
                             right: 120,
                             child: Container(
                               width: 150,
                               height: 80,
                               decoration: BoxDecoration(
                                 shape: BoxShape.circle,
                                 gradient: RadialGradient(
                                   colors: [
                                     Colors.white.withOpacity(0.4),
                                     Colors.transparent,
                                   ],
                                 ),
                               ),
                             ),
                           ),

                        // 2. The Globe (Theme Specific Style)
                        HoloGlobe(
                          size: isCompact ? 280 : 400,
                          color: isDark ? AppTheme.sanggamGold : const Color(0xFF004D40), // Gold vs Teal
                          accentColor: isDark ? Colors.white : const Color(0xFF00796B), // Highlight
                        ),

                        // 3. The Panels (Surrounding in Grid)
                        Column(
                          children: [
                            // Top Row Panels
                            Expanded(
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.end,
                                children: [
                                  // Top Left: Monitoring
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.only(bottom: 20, right: 20),
                                      child: WaveAnalysisPanel(
                                        title: '실시간 파동 모니터링', // REF: "REAL-TIME MONITORING"
                                        isActive: true,
                                        child: _buildRealTimeChart(globeColor),
                                        footer: _buildFooterStat('파동 무결성: 98.7%', '', contentColor, subContentColor),
                                      ),
                                    ),
                                  ),
                                  SizedBox(width: isCompact ? 10 : 350), // Globe Gap
                                  // Top Right: System Health
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.only(bottom: 20, left: 20),
                                      child: WaveAnalysisPanel(
                                        title: '시스템 상태',
                                        child: _buildHealthGauge(globeColor),
                                        footer: _buildFooterStat('코어 안정성: 최적', '', contentColor, subContentColor),
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            
                            // Bottom Row Panels
                            Expanded(
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  // Bottom Left: Predictive
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.only(top: 20, right: 20),
                                      child: WaveAnalysisPanel(
                                        title: 'AI 예측 모델링',
                                        child: _buildHexStructure(globeColor), // Hexagon 3D
                                        footer: _buildFooterStat('미래 예측: 안정적', '', contentColor, subContentColor),
                                      ),
                                    ),
                                  ),
                                  SizedBox(width: isCompact ? 10 : 350),
                                  // Bottom Right: Data Vault
                                  Expanded(
                                    child: Padding(
                                      padding: const EdgeInsets.only(top: 20, left: 20),
                                      child: WaveAnalysisPanel(
                                        title: '보안 데이터 금고',
                                        child: _buildTreasureChest(globeColor),
                                        footer: _buildFooterStat('암호화 수준: 양자', '', contentColor, subContentColor),
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  
                  const Spacer(),
                  const SizedBox(height: 100), // Navigation Bar Space
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildHexStructure(Color color) {
    return Center(
      child: Stack(
        alignment: Alignment.center,
        children: [
          // Background Pulse
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
          Icon(Icons.hexagon_outlined, color: AppTheme.sanggamGold.withOpacity(0.5), size: 80),
          // Inner detail
          Icon(Icons.data_object, color: color, size: 30),
          Positioned(
             top: 20,
             child: Container(width: 40, height: 1, color: AppTheme.sanggamGold),
          ),
           Positioned(
            bottom: 10,
            child: Text('AI 노드 활성', style: TextStyle(color: AppTheme.sanggamGold, fontSize: 6, letterSpacing: 1.0)),
          )
        ],
      ),
    );
  }

  Widget _buildTreasureChest(Color color) {
    return Center(
      child: FittedBox( // Fix overflow
        fit: BoxFit.scaleDown,
        child: Stack(
          alignment: Alignment.center,
          children: [
            Icon(Icons.shield_outlined, color: color.withOpacity(0.2), size: 70),
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.lock_rounded, color: AppTheme.sanggamGold, size: 32),
                const SizedBox(height: 4),
                Text('양자 암호화됨', style: TextStyle(color: color, fontSize: 6, fontWeight: FontWeight.bold)),
              ],
            ),
            // Rotating perimeter
             TweenAnimationBuilder(
              tween: Tween<double>(begin: 0, end: 2 * math.pi),
              duration: const Duration(seconds: 10),
              builder: (context, value, child) {
                 return Transform.rotate(
                   angle: value,
                   child: Container(
                     width: 50, height: 50,
                     decoration: BoxDecoration(
                       border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.3), width: 1, style: BorderStyle.solid),
                       shape: BoxShape.circle, 
                     ),
                   ),
                 );
              },
               onEnd: () {}, // Repeat logic usually requires state, simplifying here with basic Tween
               curve: Curves.linear,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRealTimeChart(Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 8),
      child: Center(
        child: CustomPaint(
          size: const Size(double.infinity, 60),
          painter: _ChartPainter(color: color),
        ),
      ),
    );
  }

  Widget _buildHealthGauge(Color color) {
    return Center(
      child: Stack(
        alignment: Alignment.center,
        children: [
          // Outer Ring
           CircularProgressIndicator(
              value: 1.0,
              color: AppTheme.sanggamGold.withOpacity(0.1),
              strokeWidth: 8,
            ),
          // Actual Value
          TweenAnimationBuilder(
            tween: Tween<double>(begin: 0, end: 0.85),
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
          // Inner Decor
          Container(
             width: 50, height: 50,
             decoration: BoxDecoration(
               shape: BoxShape.circle,
               border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.5), width: 1),
               gradient: RadialGradient(colors: [color.withOpacity(0.2), Colors.transparent]),
             ),
             child: Center(
               child: Text(
                '결', 
                style: TextStyle(
                  fontSize: 18, 
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
  _ChartPainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    // 1. Grid Background
    final gridPaint = Paint()..color = Colors.white.withOpacity(0.05)..strokeWidth = 0.5;
    final double step = size.width / 10;
    for(double x=0; x<=size.width; x+=step) canvas.drawLine(Offset(x,0), Offset(x, size.height), gridPaint);
    for(double y=0; y<=size.height; y+=10) canvas.drawLine(Offset(0,y), Offset(size.width, y), gridPaint);

    // 2. Chart Line with Gradient
    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2
      ..shader = LinearGradient(
        colors: [color.withOpacity(0), color, color, color.withOpacity(0)],
        stops: const [0.0, 0.2, 0.8, 1.0],
      ).createShader(Rect.fromLTWH(0, 0, size.width, size.height));

    final path = Path();
    path.moveTo(0, size.height * 0.5);
    
    // Multi-frequency wave for "tech" look
    for (double x = 0; x <= size.width; x+=2) {
      final y = size.height * 0.5 + 
                math.sin(x * 0.1) * 10 + 
                math.sin(x * 0.5) * 5; // Added harmonics
      path.lineTo(x, y);
    }
    canvas.drawPath(path, paint);

    // 3. Fill below line
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
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
