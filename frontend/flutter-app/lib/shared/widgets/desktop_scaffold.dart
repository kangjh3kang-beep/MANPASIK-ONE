import 'package:flutter/material.dart';
import 'package:manpasik/shared/widgets/sanggam_orbit_frame.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/hanji_background.dart';
import 'package:manpasik/shared/widgets/network_indicator.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

class DesktopScaffold extends StatelessWidget {
  final Widget body;
  final int currentIndex;
  final Function(int) onNavTap;
  final bool isDark;
  final bool hideGlobalTopFrame;

  const DesktopScaffold({
    super.key,
    required this.body,
    required this.currentIndex,
    required this.onNavTap,
    required this.isDark,
    this.hideGlobalTopFrame = false,
  });

  @override
  Widget build(BuildContext context) {
    final useDarkShellBackground = isDark || hideGlobalTopFrame;
    final shellContent = Row(
      children: [
        // [Side Navigation Rail] - Jagae & Sanggam style
        _buildSideRail(context),

        // Main Content Area
        Expanded(
          child: Column(
            children: [
              const NetworkIndicator(),
              Expanded(child: body),
            ],
          ),
        ),
      ],
    );

    final routedBody = hideGlobalTopFrame
        ? shellContent
        : SanggamOrbitFrame(child: shellContent);

    return Scaffold(
      backgroundColor: const Color(0xFF0B1021),
      body: useDarkShellBackground
          ? RoyalCloudBackground(child: routedBody)
          : HanjiBackground(child: routedBody),
    );
  }

  Widget _buildSideRail(BuildContext context) {
    final navItems = [
      (icon: Icons.home_outlined, selectedIcon: Icons.home, label: '홈'),
      (
        icon: Icons.bar_chart_outlined,
        selectedIcon: Icons.bar_chart,
        label: '데이터'
      ),
      (icon: Icons.hexagon_outlined, selectedIcon: Icons.hexagon, label: '측정'),
      (
        icon: Icons.local_hospital_outlined,
        selectedIcon: Icons.local_hospital,
        label: '의료'
      ),
      (
        icon: Icons.shopping_cart_outlined,
        selectedIcon: Icons.shopping_cart,
        label: '마켓'
      ),
      (icon: Icons.forum_outlined, selectedIcon: Icons.forum, label: '커뮤니티'),
      (
        icon: Icons.family_restroom_outlined,
        selectedIcon: Icons.family_restroom,
        label: '가족'
      ),
      (
        icon: Icons.settings_outlined,
        selectedIcon: Icons.settings,
        label: '설정'
      ),
    ];

    return Container(
      width: 120, // Slightly wider for premium feel
      margin: const EdgeInsets.all(16),
      child: Stack(
        children: [
          // [Layer 1] Side Rail Container with Jagae texture
          SanggamContainer(
            borderRadius: 32,
            borderWidth: 1.2,
            borderColor: SanggamTheme.primary.withValues(alpha: 0.4),
            blurSigma: 16,
            jagaeOpacity: 0.12, // Increased for desktop presence
            padding: const EdgeInsets.symmetric(vertical: 32),
            child: Column(
              children: [
                // Brand Logo with Sanggam Orbit style
                Padding(
                  padding: const EdgeInsets.only(bottom: 40),
                  child: Hero(
                    tag: 'brand_logo',
                    child: Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(
                          color: SanggamTheme.primary.withValues(alpha: 0.5),
                          width: 0.8,
                        ),
                        boxShadow: [
                          BoxShadow(
                            color: SanggamTheme.primary.withValues(alpha: 0.2),
                            blurRadius: 12,
                          ),
                        ],
                      ),
                      child: const Icon(
                        Icons.auto_awesome,
                        color: SanggamTheme.primary,
                        size: 32,
                      ),
                    ),
                  ),
                ),

                Expanded(
                  child: ListView.separated(
                    itemCount: navItems.length,
                    separatorBuilder: (context, index) =>
                        const SizedBox(height: 24),
                    itemBuilder: (context, index) {
                      final isSelected = currentIndex == index;
                      final item = navItems[index];

                      return GestureDetector(
                        onTap: () => onNavTap(index),
                        child: MouseRegion(
                          cursor: SystemMouseCursors.click,
                          child: Column(
                            children: [
                              AnimatedContainer(
                                duration: const Duration(milliseconds: 400),
                                curve: Curves.easeOutCubic,
                                padding: const EdgeInsets.all(14),
                                decoration: BoxDecoration(
                                  color: isSelected
                                      ? SanggamTheme.primary
                                          .withValues(alpha: 0.18)
                                      : Colors.transparent,
                                  borderRadius: BorderRadius.circular(20),
                                  border: isSelected
                                      ? Border.all(
                                          color: SanggamTheme.primary
                                              .withValues(alpha: 0.6),
                                          width: 0.5)
                                      : null,
                                  boxShadow: isSelected
                                      ? [
                                          BoxShadow(
                                            color: SanggamTheme.primary
                                                .withValues(alpha: 0.15),
                                            blurRadius: 10,
                                            spreadRadius: -2,
                                          )
                                        ]
                                      : null,
                                ),
                                child: Icon(
                                  isSelected ? item.selectedIcon : item.icon,
                                  color: isSelected
                                      ? SanggamTheme.primary
                                      : Colors.white60,
                                  size: 28,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Text(
                                item.label,
                                style: TextStyle(
                                  color: isSelected
                                      ? SanggamTheme.primary
                                      : Colors.white38,
                                  fontSize: 11,
                                  fontWeight: isSelected
                                      ? FontWeight.bold
                                      : FontWeight.normal,
                                  letterSpacing: 0.5,
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),

          // [Layer 2] Vertical Sanggam Gold Line (Subtle accent)
          Positioned(
            top: 100,
            bottom: 40,
            right: 0,
            width: 1,
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [
                    Colors.transparent,
                    SanggamTheme.primary.withValues(alpha: 0.5),
                    Colors.transparent,
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
