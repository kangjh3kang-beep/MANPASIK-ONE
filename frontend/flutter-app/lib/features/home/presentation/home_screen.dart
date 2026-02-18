import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

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


import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final textColor = isDark ? Colors.white : AppTheme.inkBlack;
    final subTextColor = isDark ? Colors.white70 : AppTheme.inkBlack.withOpacity(0.7);

    final authState = ref.watch(authProvider);
    final historyAsync = ref.watch(measurementHistoryProvider);

    return Scaffold(
      backgroundColor: Colors.transparent, // Background handled by AppShell
      body: SafeArea(
        child: CustomScrollView(
          physics: const BouncingScrollPhysics(),
          slivers: [
            // Apple-Style Collapsing Header (Transparent Glass)
            SliverAppBar(
              expandedHeight: 100.0,
              floating: false,
              pinned: true,
              backgroundColor: Colors.transparent, // Glass effect below handles this
              surfaceTintColor: Colors.transparent,
              elevation: 0,
              flexibleSpace: FlexibleSpaceBar(
                titlePadding: const EdgeInsets.only(left: 24, bottom: 16),
                centerTitle: false,
                title: Text(
                  authState.displayName ?? 'Guest',
                  style: theme.textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: textColor,
                    shadows: isDark ? [
                      Shadow(color: AppTheme.sanggamGold, blurRadius: 10),
                    ] : null,
                  ),
                ),
                background: Padding(
                  padding: const EdgeInsets.only(left: 24, top: 20),
                  child: Align(
                    alignment: Alignment.topLeft,
                    child: Text(
                      AppLocalizations.of(context)!.greeting,
                      style: theme.textTheme.bodyLarge?.copyWith(
                        color: subTextColor,
                      ),
                    ),
                  ),
                ),
              ),
              actions: [
                _NotificationBell(color: subTextColor),
                IconButton(
                  icon: Icon(Icons.settings_rounded, color: subTextColor),
                  tooltip: AppLocalizations.of(context)!.settings,
                  onPressed: () => context.push('/settings'),
                ),
                const SizedBox(width: 16),
              ],
            ),

            // Main Content Body
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: Column(
                  children: [
                    const SizedBox(height: 16),
                    
                    // 1. Hero Section (Holographic Glass + Breathing Glow)
                    AnimateFadeInUp(
                      duration: const Duration(milliseconds: 700),
                      child: OrnateGoldFrame( // Upgraded to OrnateGoldFrame
                        width: double.infinity,
                        isActive: true,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Icon(
                                  Icons.query_stats_rounded,
                                  color: isDark ? AppTheme.waveCyan : AppTheme.celadonTeal, 
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  AppLocalizations.of(context)!.newMeasurement,
                                  style: theme.textTheme.titleMedium?.copyWith(
                                    color: isDark ? AppTheme.waveCyan : AppTheme.celadonTeal,
                                    fontWeight: FontWeight.bold,
                                    shadows: [
                                       Shadow(color: isDark ? AppTheme.waveCyan : Colors.transparent, blurRadius: 8),
                                    ],
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 16),
                            Text(
                              AppLocalizations.of(context)!.checkHealth,
                              style: theme.textTheme.headlineSmall?.copyWith(
                                color: textColor,
                                fontWeight: FontWeight.bold,
                                height: 1.2,
                              ),
                            ),
                            const SizedBox(height: 24),
                            
                            // Premium Gradient Button
                            ScaleButton(
                              onPressed: () => context.push('/measure'),
                              child: Container(
                                width: double.infinity,
                                padding: const EdgeInsets.symmetric(vertical: 16),
                                decoration: BoxDecoration(
                                  borderRadius: BorderRadius.circular(12),
                                  gradient: const LinearGradient(
                                    colors: [Color(0xFFE6C15D), Color(0xFFB38B24)], // Polished Gold
                                  ),
                                  boxShadow: [
                                    BoxShadow(
                                      color: AppTheme.sanggamGold.withOpacity(0.5),
                                      blurRadius: 16,
                                      spreadRadius: -2,
                                      offset: const Offset(0, 4),
                                    ),
                                  ],
                                ),
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    const Icon(Icons.auto_awesome_rounded, color: Color(0xFF050B14)),
                                    const SizedBox(width: 8),
                                    Text(
                                      AppLocalizations.of(context)!.startMeasurementAction,
                                      style: theme.textTheme.titleMedium?.copyWith(
                                        color: const Color(0xFF050B14),
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),

                    // 2. Quick Actions (Holo Glass Tiles)
                    AnimateFadeInUp(
                      delay: const Duration(milliseconds: 100),
                      child: Row(
                        children: [
                          _HoloQuickAction(
                            icon: Icons.bar_chart_rounded,
                            label: '데이터',
                            textColor: subTextColor,
                            onTap: () => context.go('/data'),
                          ),
                          const SizedBox(width: 12),
                          _HoloQuickAction(
                            icon: Icons.smart_toy_rounded,
                            label: 'AI 코치',
                            textColor: subTextColor,
                            onTap: () => context.push('/coach'),
                          ),
                          const SizedBox(width: 12),
                          _HoloQuickAction(
                            icon: Icons.family_restroom_rounded,
                            label: '가족',
                            textColor: subTextColor,
                            onTap: () => context.push('/family'),
                          ),
                          const SizedBox(width: 12),
                          _HoloQuickAction(
                            icon: Icons.local_hospital_rounded,
                            label: '의료',
                            textColor: subTextColor,
                            onTap: () => context.push('/medical'),
                          ),
                        ],
                      ),
                    ),

                    const SizedBox(height: 16),

                    // Device Quick Access Banner
                    AnimateFadeInUp(
                      delay: const Duration(milliseconds: 150),
                      child: GestureDetector(
                        onTap: () => context.push('/devices'),
                        child: HoloGlassCard(
                          child: Row(
                            children: [
                              Icon(Icons.bluetooth_connected, color: isDark ? AppTheme.waveCyan : AppTheme.celadonTeal, size: 20),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Text(
                                  '기기 관리',
                                  style: theme.textTheme.bodyMedium?.copyWith(
                                    color: textColor,
                                    fontWeight: FontWeight.w600,
                                  ),
                                ),
                              ),
                              Text(
                                '연결 상태 확인',
                                style: TextStyle(color: subTextColor.withOpacity(0.5), fontSize: 12),
                              ),
                              const SizedBox(width: 4),
                              Icon(Icons.chevron_right, color: subTextColor.withOpacity(0.5), size: 16),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),

                    // Recent History Header
                    AnimateFadeInUp(
                      delay: const Duration(milliseconds: 200),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            AppLocalizations.of(context)!.recentHistory,
                            style: theme.textTheme.titleLarge?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: textColor,
                            ),
                          ),
                          TextButton(
                            onPressed: () => context.go('/data'),
                            child: Text(
                              AppLocalizations.of(context)!.viewAll,
                              style: TextStyle(color: subTextColor.withOpacity(0.6)),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 8),
                  ],
                ),
              ),
            ),

            // 3. History List (Updated Card Style needed later, keeping logic for now)
            historyAsync.when(
              data: (result) {
                if (result.items.isEmpty) {
                  return SliverFillRemaining(
                    hasScrollBody: false,
                    child: Center(
                      child: AnimateFadeInUp(
                        delay: const Duration(milliseconds: 300),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(Icons.analytics_outlined, size: 48, color: Colors.white24),
                            const SizedBox(height: 12),
                            Text(
                              '아직 측정 기록이 없습니다.\n첫 측정을 시작해보세요!',
                              textAlign: TextAlign.center,
                              style: theme.textTheme.bodyLarge?.copyWith(color: Colors.white54),
                            ),
                          ],
                        ),
                      ),
                    ),
                  );
                }
                return SliverPadding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  sliver: SliverList(
                    delegate: SliverChildBuilderDelegate(
                      (context, index) {
                        final item = result.items[index];
                        final date = item.measuredAt ?? DateTime.now();
                        final type = item.primaryValue <= 100
                            ? 'normal'
                            : (item.primaryValue <= 125 ? 'warning' : 'high');
                        
                        return AnimateFadeInUp(
                          delay: Duration(milliseconds: 300 + (index * 50)),
                          offset: 20,
                          child: Padding(
                            padding: const EdgeInsets.only(bottom: 12),
                            child: MeasurementCard( // Needs upgrade to Holo style later
                              date: date,
                              value: item.primaryValue,
                              unit: item.unit.isNotEmpty ? item.unit : 'mg/dL',
                              resultType: type,
                              onTap: () => context.push('/measure/result'),
                            ),
                          ),
                        );
                      },
                      childCount: result.items.length,
                    ),
                  ),
                );
              },
              loading: () => const SliverFillRemaining(
                child: Center(child: CircularProgressIndicator(color: AppTheme.sanggamGold)),
              ),
              error: (err, _) => SliverFillRemaining(
                child: Center(child: Text('Error', style: TextStyle(color: Colors.white))),
              ),
            ),
            
            const SliverToBoxAdapter(child: SizedBox(height: 100)), // Space for Glass Dock
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
