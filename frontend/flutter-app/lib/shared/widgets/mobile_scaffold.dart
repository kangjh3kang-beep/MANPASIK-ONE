import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/shared/providers/ecosystem_providers.dart';
import 'package:manpasik/shared/widgets/glass_dock_navigation.dart';
import 'package:manpasik/shared/widgets/sanggam_orbit_frame.dart';
import 'package:manpasik/shared/widgets/royal_cloud_background.dart';
import 'package:manpasik/shared/widgets/hanji_background.dart';
import 'package:manpasik/shared/widgets/network_indicator.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/floating_monitor_bubble.dart';

class MobileScaffold extends ConsumerWidget {
  final Widget body;
  final int currentIndex;
  final Function(int) onNavTap;
  final bool isDark;
  final bool hideGlobalDock;
  final bool hideGlobalTopFrame;

  const MobileScaffold({
    super.key,
    required this.body,
    required this.currentIndex,
    required this.onNavTap,
    required this.isDark,
    this.hideGlobalDock = false,
    this.hideGlobalTopFrame = false,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scaffoldKey = GlobalKey<ScaffoldState>();
    final useDarkShellBackground =
        isDark || hideGlobalTopFrame || hideGlobalDock;
    final shellContent = Stack(
      children: [
        Column(
          children: [
            const NetworkIndicator(),
            Expanded(child: body),
          ],
        ),
        const FloatingMonitorBubble(),
      ],
    );

    final routedBody = hideGlobalTopFrame
        ? shellContent
        : SanggamOrbitFrame(child: shellContent);

    // 글로벌 배지 카운트 연동: 가족(6번 탭) 알림 표시
    final badges = ref.watch(globalBadgeProvider);
    final badgeCounts = <int, int>{};
    final familyCount = badges['family'] ?? 0;
    final notifCount = badges['notification'] ?? 0;
    if (familyCount > 0) badgeCounts[6] = familyCount; // 가족 탭
    if (notifCount > 0) badgeCounts[7] = notifCount; // 설정 탭 (알림 대용)

    return Scaffold(
      key: scaffoldKey,
      backgroundColor: const Color(0xFF0B1021),
      drawer: _MobileNavDrawer(
        currentIndex: currentIndex,
        onNavTap: onNavTap,
      ),
      extendBody: true,
      body: useDarkShellBackground
          ? RoyalCloudBackground(child: routedBody)
          : HanjiBackground(child: routedBody),
      bottomNavigationBar: hideGlobalDock
          ? null
          : GlassDockNavigation(
              currentIndex: currentIndex,
              onTap: onNavTap,
              badgeCounts: badgeCounts,
            ),
    );
  }
}

class _MobileNavDrawer extends StatelessWidget {
  const _MobileNavDrawer({
    required this.currentIndex,
    required this.onNavTap,
  });

  final int currentIndex;
  final Function(int) onNavTap;

  static const _items = [
    (icon: Icons.home_outlined, label: '홈'),
    (icon: Icons.bar_chart_outlined, label: '데이터'),
    (icon: Icons.hexagon_outlined, label: '측정'),
    (icon: Icons.local_hospital_outlined, label: '의료'),
    (icon: Icons.shopping_cart_outlined, label: '마켓'),
    (icon: Icons.forum_outlined, label: '커뮤니티'),
    (icon: Icons.family_restroom_outlined, label: '가족'),
    (icon: Icons.settings_outlined, label: '설정'),
  ];

  @override
  Widget build(BuildContext context) {
    return Drawer(
      backgroundColor: const Color(0xFF091325),
      child: SafeArea(
        child: Column(
          children: [
            const ListTile(
              title: Text(
                '전체 메뉴',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 18,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
            const Divider(height: 1, color: Color(0x33D4AF37)),
            Expanded(
              child: ListView.builder(
                itemCount: _items.length,
                itemBuilder: (context, index) {
                  final item = _items[index];
                  final selected = index == currentIndex;
                  return ListTile(
                    selected: selected,
                    selectedTileColor: const Color(0x22D4AF37),
                    leading: Icon(
                      item.icon,
                      color:
                          selected ? const Color(0xFFD4AF37) : Colors.white70,
                    ),
                    title: Text(
                      item.label,
                      style: TextStyle(
                        color:
                            selected ? const Color(0xFFD4AF37) : Colors.white,
                        fontWeight:
                            selected ? FontWeight.w700 : FontWeight.w500,
                      ),
                    ),
                    onTap: () {
                      Navigator.of(context).pop();
                      onNavTap(index);
                    },
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
