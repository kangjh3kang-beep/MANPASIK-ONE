import 'package:flutter/material.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

// ───────────────────────────────────────────────────
// GlassDockNavigation — SanggamContainer 기반 하단 독 내비게이션
//
// [Rule 4] dart:ui, go_router 미사용 import 제거
// [Rule 4] JagaeCard+Container 이중 → SanggamContainer 단일
// [Rule 4] Color(0xFFD4AF37) → SanggamTheme.primary
// [Rule 4] withOpacity 6건 → withValues(alpha:)
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 2] bottom 12→16, navItem height 52→48 (8×6)
// ───────────────────────────────────────────────────

class GlassDockNavigation extends StatelessWidget {
  final int currentIndex;
  final Function(int) onTap;
  final Map<int, int> badgeCounts;

  const GlassDockNavigation({
    super.key,
    required this.currentIndex,
    required this.onTap,
    this.badgeCounts = const {},
  });

  static const _navItems = [
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
    (icon: Icons.settings_outlined, selectedIcon: Icons.settings, label: '설정'),
  ];

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
      child: Stack(
        clipBehavior: Clip.none,
        alignment: Alignment.bottomCenter,
        children: [
          // [Layer 1] Sanggam Container (Glass base)
          SanggamContainer(
            borderRadius: 32,
            borderWidth: 0.8,
            borderColor: SanggamTheme.primary.withValues(alpha: 0.3),
            blurSigma: 8,
            jagaeOpacity: 0.06,
            height: 64,
            padding: EdgeInsets.zero,
            child: Row(
              children: List.generate(_navItems.length, (index) {
                final item = _navItems[index];
                final isSelected = currentIndex == index;
                return Expanded(
                  child: _GlassNavItem(
                    item: item,
                    isSelected: isSelected,
                    badgeCount: badgeCounts[index] ?? 0,
                    onTap: () => onTap(index),
                  ),
                );
              }),
            ),
          ),

          // [Layer 2] Photorealistic Slim Dock Asset (v5)
          Positioned(
            bottom: -4, // Adjust to overlap slightly or align
            left: -12,
            right: -12,
            height: 84, // Slightly taller to cover the container
            child: IgnorePointer(
              child: Opacity(
                opacity: 0.28,
                child: ColorFiltered(
                  colorFilter: const ColorFilter.mode(
                    Color(0xFF15233A),
                    BlendMode.multiply,
                  ),
                  child: Image.asset(
                    'assets/images/bottom_3d_dock_slim.png',
                    fit: BoxFit.cover,
                    alignment: Alignment.bottomCenter,
                    filterQuality: FilterQuality.high,
                    errorBuilder: (context, error, stackTrace) =>
                        const SizedBox.shrink(),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _GlassNavItem extends StatefulWidget {
  final ({IconData icon, IconData selectedIcon, String label}) item;
  final bool isSelected;
  final int badgeCount;
  final VoidCallback onTap;

  const _GlassNavItem({
    required this.item,
    required this.isSelected,
    this.badgeCount = 0,
    required this.onTap,
  });

  @override
  State<_GlassNavItem> createState() => _GlassNavItemState();
}

class _GlassNavItemState extends State<_GlassNavItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final isSelected = widget.isSelected;
    final isHovered = _isHovered;

    Color iconColor;
    if (isSelected) {
      iconColor = SanggamTheme.primary;
    } else if (isHovered) {
      iconColor = Colors.white.withValues(alpha: 0.9);
    } else {
      iconColor = Colors.white.withValues(alpha: 0.5);
    }

    Color bgColor;
    if (isSelected) {
      bgColor = SanggamTheme.primary.withValues(alpha: 0.1);
    } else if (isHovered) {
      bgColor = Colors.white.withValues(alpha: 0.05);
    } else {
      bgColor = Colors.transparent;
    }

    return Semantics(
      button: true,
      label: '${widget.item.label} 탭${isSelected ? ", 선택됨" : ""}',
      child: MouseRegion(
        onEnter: (_) => setState(() => _isHovered = true),
        onExit: (_) => setState(() => _isHovered = false),
        cursor: SystemMouseCursors.click,
        child: ScaleButton(
          onPressed: widget.onTap,
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeOutCubic,
            height: 48,
            margin: const EdgeInsets.symmetric(horizontal: 2),
            decoration: BoxDecoration(
              color: bgColor,
              borderRadius: BorderRadius.circular(16),
              boxShadow: isSelected
                  ? [
                      BoxShadow(
                        color: SanggamTheme.primary.withValues(alpha: 0.3),
                        blurRadius: 10,
                        spreadRadius: -2,
                      ),
                    ]
                  : null,
            ),
            child: Stack(
              alignment: Alignment.center,
              clipBehavior: Clip.none,
              children: [
                Icon(
                  isSelected ? widget.item.selectedIcon : widget.item.icon,
                  color: iconColor,
                  size: 22,
                ),
                if (widget.badgeCount > 0)
                  Positioned(
                    top: 6,
                    right: 4,
                    child: Container(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 4, vertical: 1),
                      constraints:
                          const BoxConstraints(minWidth: 14, minHeight: 14),
                      decoration: BoxDecoration(
                        color: SanggamTheme.error,
                        borderRadius: BorderRadius.circular(7),
                      ),
                      child: Text(
                        widget.badgeCount > 99 ? '99+' : '${widget.badgeCount}',
                        textAlign: TextAlign.center,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 8,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
