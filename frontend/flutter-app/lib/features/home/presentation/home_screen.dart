import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'dart:async'; // For Timer
import 'dart:math' as math; // For Random

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/measurement_card.dart';
import 'package:manpasik/shared/widgets/holo_glass_card.dart'; // Upgrade!
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/l10n/app_localizations.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';
import 'package:manpasik/shared/widgets/sanggam_decoration.dart';
import 'package:manpasik/shared/widgets/mini_line_chart.dart';


import 'package:manpasik/shared/widgets/holo_globe.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/jarvis_connector.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';

class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> with TickerProviderStateMixin {
  // Animation Controllers for Entrance
  late AnimationController _entranceCtrl;

  // Demo Mode Dynamic Data
  Timer? _demoTimer;
  String _coreStatus = '안정';
  String _networkStatus = '연결됨';
  String _syncRate = '99.8%';
  String _aiAnomaly = '0.00%';
  List<double> _chartData = [10, 15, 8, 20, 12, 18, 14];
  final math.Random _random = math.Random();

  @override
  void initState() {
    super.initState();
    _entranceCtrl = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..forward();

    // Start Mock Data Timer if Demo Mode
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (ref.read(authProvider).isDemo) {
        _startDemoTimer();
      }
    });
  }

  void _startDemoTimer() {
    _demoTimer = Timer.periodic(const Duration(milliseconds: 2000), (timer) {
      if (!mounted) return;
      setState(() {
        // Core Status (Rarely busy)
        _coreStatus = _random.nextDouble() > 0.9 ? '처리중' : '안정';
        
        // Network Status
        _networkStatus = _random.nextDouble() > 0.95 ? '암호화' : '연결됨';
        
        // Sync Rate (98.0 - 99.9)
        _syncRate = '${(98.0 + _random.nextDouble() * 1.9).toStringAsFixed(1)}%';
        
        // AI Anomaly (0.00 - 0.05)
        _aiAnomaly = '${(_random.nextDouble() * 0.05).toStringAsFixed(2)}%';

        // Chart Data Shift
        _chartData.removeAt(0);
        _chartData.add(5.0 + _random.nextInt(25).toDouble()); // Random 5-30
      });
    });
  }

  @override
  void dispose() {
    _entranceCtrl.dispose();
    _demoTimer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    // Theme & Data
    final authState = ref.watch(authProvider);
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    // Fixed HUD Layout
    return Scaffold(
      backgroundColor: Colors.transparent, 
      body: Stack(
        fit: StackFit.expand,
        children: [
          // 1. Background handled by AppRouter (RoyalCloudBackground)
          
          // 2. Central HoloGlobe (The Core)
          // Positioned slightly upwards to leave room for the dock
          Positioned.fill(
             bottom: 100, // Leave room for bottom sheet/dock
             top: 0,
             child: Center(
               child: Transform.scale(
                 scale: 1.1, // Grand Scale
                 child: const SizedBox(
                   width: 350,
                   height: 350,
                   child: HoloGlobe(),
                 ),
               ),
             ),
          ),

          // 3. HUD Modules (Floating panels connected to Center)
          // We use LayoutBuilder to position relative to screen size
          LayoutBuilder(
            builder: (context, constraints) {
              final center = Offset(constraints.maxWidth / 2, constraints.maxHeight / 2 - 50);
              
              // Responsive HUD Logic
              // Responsive HUD Logic (Dynamic Positioning)
              final isCompact = constraints.maxWidth < 800 || constraints.maxHeight < 600;
              final isUltraWide = constraints.maxWidth > 1600;
              
              // Dynamic Spacing based on screen size
              final spreadX = constraints.maxWidth * (isCompact ? 0.35 : 0.25);
              
              // Prevent spreading too far vertically on short screens
              final spreadY = (constraints.maxHeight * 0.25).clamp(80.0, 150.0);
              
              final panelWidth = isCompact ? 140.0 : 180.0;
              // Remove fixed height constraint for flexibility, or use a constrained range
              final panelHeight = 110.0; // Standardized height

              // Define Module Positions relative to Center
              final leftTopPos = Offset(center.dx - spreadX, center.dy - spreadY);
              final rightTopPos = Offset(center.dx + spreadX - panelWidth, center.dy - spreadY);
              
              // Bottom panels position - Ensure they don't overlap with Bottom Dock (approx 100px height + padding)
              double bottomY = center.dy + spreadY;
              final double dockSafeLimit = constraints.maxHeight - 160; // 160px from bottom safety margin
              
              if (bottomY + panelHeight > dockSafeLimit) {
                 bottomY = dockSafeLimit - panelHeight; 
                 // If that pushes it above center, we have a very small screen problem.
                 // In that case, the center needs to move up.
              }

              final leftBotPos = Offset(center.dx - spreadX, bottomY);
              final rightBotPos = Offset(center.dx + spreadX - panelWidth, bottomY);

              return Stack(
                children: [
                   // Connectors (Animated Lines)
                   JarvisConnector(startOffset: center, endOffset: leftTopPos + const Offset(140, 60)),
                   JarvisConnector(startOffset: center, endOffset: rightTopPos + const Offset(0, 60)),
                   JarvisConnector(startOffset: center, endOffset: leftBotPos + const Offset(140, 0)),
                   JarvisConnector(startOffset: center, endOffset: rightBotPos + const Offset(0, 0)),

                   // Top Left: System Status (Health)
                   Positioned(
                     left: leftTopPos.dx, top: leftTopPos.dy,
                     child: _buildHudPanel(
                       width: panelWidth, height: panelHeight,
                       title: '시스템 상태',
                       icon: Icons.monitor_heart_outlined, 
                       // Icon opacity handled in _buildHudPanel or manually here if needed, 
                       // but _buildHudPanel uses the iconData. 
                       // I will update _buildHudPanel to apply 60% opacity to the icon.
                       content: Column(
                         crossAxisAlignment: CrossAxisAlignment.start,
                         children: [
                           _buildStatusRow('코어', _coreStatus, _coreStatus == '안정' ? Colors.greenAccent : Colors.amberAccent),
                           _buildStatusRow('네트워크', _networkStatus, AppTheme.waveCyan),
                           _buildStatusRow('동기화', _syncRate, AppTheme.sanggamGold),
                         ],
                       ),
                       onTap: () => context.push('/measure'),
                     ),
                   ),

                   // Top Right: Environmental / AI Analysis
                   Positioned(
                     left: rightTopPos.dx, top: rightTopPos.dy, // Use calculated pos
                     child: _buildHudPanel(
                       width: panelWidth, height: panelHeight,
                       title: 'AI 분석',
                       icon: Icons.psychology_outlined,
                       content: Text(
                         '실시간 파동 패턴 인식 활성.\n이상 징후: $_aiAnomaly',
                         style: const TextStyle(color: Colors.white70, fontSize: 10),
                       ),
                       onTap: () => context.push('/coach'),
                     ),
                   ),

                   // Bottom Left: Recent Data
                   Positioned(
                     left: leftBotPos.dx, top: leftBotPos.dy,
                     child: _buildHudPanel(
                       width: panelWidth, height: panelHeight,
                       title: '데이터 스트림',
                       icon: Icons.data_usage,
                       content: _buildMiniChart(),
                       onTap: () => context.go('/data'),
                     ),
                   ),

                   // Bottom Right: Security / Device
                   Positioned(
                     left: rightBotPos.dx, top: rightBotPos.dy,
                     child: _buildHudPanel(
                       width: panelWidth, height: panelHeight,
                       title: '보안 금고',
                       icon: Icons.lock_outline,
                       content: Center(
                         child: Icon(Icons.fingerprint, size: 40, color: AppTheme.sanggamGold.withOpacity(0.6)), // 60% Icon
                       ),
                       onTap: () => context.push('/settings/security'),
                     ),
                   ),
                   
                   // Top Center Header (User Ranking/Identity)
                   Positioned(
                     top: 60, left: 0, right: 0,
                     child: Semantics(
                       header: true,
                       label: '만파식 시스템 대시보드',
                       child: Center(
                       child: Column(
                         children: [
                           Text(
                             '만파식 시스템',
                             style: TextStyle(
                               fontFamily: 'Orbitron', // Assuming font exists or fallback
                               color: AppTheme.sanggamGold,
                               fontWeight: FontWeight.bold,
                               fontSize: 16,
                               letterSpacing: 2.0,
                               shadows: [Shadow(color: AppTheme.sanggamGold, blurRadius: 10)]
                             ),
                           ),
                           const SizedBox(height: 4),
                           Text(
                             'CMD: ${authState.displayName ?? "GUEST"}',
                             style: const TextStyle(
                               color: Colors.white54,
                               fontSize: 10,
                               letterSpacing: 1.5,
                             ),
                           ),
                         ],
                       ),
                     ),
                     ),
                   ),
                ],
              );
            },
          ),
          
          // 4. Large Action Button (Bottom Center Overlap)
          Positioned(
            bottom: 110,
            left: 0, right: 0,
            child: Semantics(
              button: true,
              label: '건강 분석 시작 버튼',
              child: Center(
              child: ScaleButton(
                onPressed: () => context.push('/measure'),
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 12),
                  decoration: BoxDecoration(
                    color: AppTheme.sanggamGold.withOpacity(0.2),
                    border: Border.all(color: AppTheme.sanggamGold),
                    borderRadius: BorderRadius.circular(30),
                    boxShadow: [
                      BoxShadow(color: AppTheme.sanggamGold.withOpacity(0.3), blurRadius: 15, spreadRadius: 1)
                    ],
                  ),
                  child: const Text(
                    '분석 시작',
                    style: TextStyle(
                       color: AppTheme.sanggamGold, 
                       fontWeight: FontWeight.bold,
                       letterSpacing: 1.5
                    ),
                  ),
                ),
              ),
            ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusRow(String label, String value, Color valueColor) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.white54, fontSize: 10)),
          AnimatedSwitcher(
            duration: const Duration(milliseconds: 300),
            child: Text(
              value, 
              key: ValueKey(value),
              style: TextStyle(color: valueColor, fontSize: 10, fontWeight: FontWeight.bold)
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMiniChart() {
    return LayoutBuilder(
      builder: (context, constraints) {
        final chartWidth = constraints.maxWidth.isFinite ? constraints.maxWidth : 120.0;
        return SizedBox(
          height: 40,
          width: chartWidth,
          child: MiniLineChart(
            data: _chartData,
            width: chartWidth,
            height: 40,
            lineColor: AppTheme.waveCyan,
            fillColor: AppTheme.waveCyan.withOpacity(0.15),
          ),
        );
      },
    );
  }

  Widget _buildHudPanel({
    required double width,
    required double height,
    required String title,
    required IconData icon,
    required Widget content,
    VoidCallback? onTap,
  }) {
    return ScaleButton(
      onPressed: onTap ?? () {},
      child: Container(
        width: width, height: height,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: const Color(0xFF0A1020).withOpacity(0.6), // Dark Glass
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: AppTheme.sanggamGold.withOpacity(0.3), width: 1),
          boxShadow: [
             BoxShadow(color: Colors.black.withOpacity(0.5), blurRadius: 10),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(icon, size: 14, color: AppTheme.sanggamGold.withOpacity(0.6)), // 60% Icon
                const SizedBox(width: 8),
                Expanded(child: Text(title, overflow: TextOverflow.ellipsis, style: const TextStyle(color: AppTheme.sanggamGold, fontSize: 10, fontWeight: FontWeight.bold, letterSpacing: 1.0))),
              ],
            ),
            const Divider(color: Colors.white10, height: 12),
            Expanded(
              child: FittedBox(
                fit: BoxFit.scaleDown,
                alignment: Alignment.topLeft,
                child: ConstrainedBox(
                   constraints: BoxConstraints(minWidth: width - 24), // Ensure width fill
                   child: content
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// 알림 벨 아이콘
class _NotificationBell extends ConsumerWidget {
  final Color? color;
  const _NotificationBell({this.color});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return IconButton(
      icon: Icon(Icons.notifications_outlined, color: color ?? Colors.white70),
      onPressed: () => context.push('/notifications'),
    );
  }
}

/// Holo Style Quick Action Button
class _HoloQuickAction extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final Color? textColor;

  const _HoloQuickAction({
    required this.icon,
    required this.label,
    required this.onTap,
    this.textColor,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: ScaleButton(
        onPressed: onTap,
        child: Column(
          children: [
            HoloGlassCard(
              padding: const EdgeInsets.symmetric(vertical: 16),
              child: Center(
                child: Icon(icon, size: 28, color: AppTheme.sanggamGold),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: TextStyle(
                color: textColor ?? Colors.white70,
                fontSize: 12,
                fontWeight: FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
