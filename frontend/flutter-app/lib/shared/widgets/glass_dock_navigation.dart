import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

class GlassDockNavigation extends StatelessWidget {
  final int currentIndex;
  final Function(int) onTap;

  const GlassDockNavigation({
    super.key,
    required this.currentIndex,
    required this.onTap,
  });

  static const _navItems = [
    (icon: Icons.home_outlined, selectedIcon: Icons.home, label: '홈'),
    (icon: Icons.bar_chart_outlined, selectedIcon: Icons.bar_chart, label: '데이터'),
    (icon: Icons.hexagon_outlined, selectedIcon: Icons.hexagon, label: '측정'),
    (icon: Icons.shopping_cart_outlined, selectedIcon: Icons.shopping_cart, label: '마켓'),
    (icon: Icons.settings_outlined, selectedIcon: Icons.settings, label: '설정'),
  ];

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 24),
      child: KoreanEdgeBorder(
        borderRadius: BorderRadius.circular(32),
        borderWidth: 0.5, // Thinner, elegant border
        child: ClipRRect(
          borderRadius: BorderRadius.circular(32),
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 5, sigmaY: 5), // Reduced blur to show background
            child: Container(
              height: 80,
              decoration: BoxDecoration(
                color: const Color(0xFF0A0E21).withOpacity(0.05), // Ultra-Transparent (5% Opacity)
                borderRadius: BorderRadius.circular(32),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.1), // Softer shadow
                    spreadRadius: 0,
                    blurRadius: 20,
                    offset: const Offset(0, 10),
                  ),
                ],
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: List.generate(_navItems.length, (index) {
                  final item = _navItems[index];
                  final isSelected = currentIndex == index;
                  return _GlassNavItem(
                    item: item,
                    isSelected: isSelected,
                    onTap: () => onTap(index),
                  );
                }),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _GlassNavItem extends StatefulWidget {
  final ({IconData icon, IconData selectedIcon, String label}) item;
  final bool isSelected;
  final VoidCallback onTap;

  const _GlassNavItem({
    required this.item,
    required this.isSelected,
    required this.onTap,
  });

  @override
  State<_GlassNavItem> createState() => _GlassNavItemState();
}

class _GlassNavItemState extends State<_GlassNavItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    // Sanggam Gold
    const goldColor = Color(0xFFD4AF37); 
    
    // Dynamic styles based on state
    final isSelected = widget.isSelected;
    final isHovered = _isHovered;

    Color iconColor;
    if (isSelected) {
      iconColor = goldColor;
    } else if (isHovered) {
      iconColor = Colors.white.withOpacity(0.9); // Brighter on hover
    } else {
      iconColor = Colors.white.withOpacity(0.4);
    }

    Color bgColor;
    if (isSelected) {
      bgColor = goldColor.withOpacity(0.15);
    } else if (isHovered) {
      bgColor = Colors.white.withOpacity(0.1); // Subtle light background
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
          duration: const Duration(milliseconds: 200),
          curve: Curves.easeOut,
          width: 60,
          height: 60,
          decoration: BoxDecoration(
            color: bgColor,
            shape: BoxShape.circle,
            boxShadow: isHovered && !isSelected 
                ? [BoxShadow(color: Colors.white.withOpacity(0.1), blurRadius: 8, spreadRadius: 0)] 
                : null,
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                isSelected ? widget.item.selectedIcon : widget.item.icon,
                color: iconColor,
                size: 26,
              ),
              const SizedBox(height: 4),
              AnimatedOpacity(
                duration: const Duration(milliseconds: 200),
                opacity: isSelected ? 1.0 : 0.0,
                child: isSelected 
                  ? Container(
                      width: 4,
                      height: 4,
                      decoration: const BoxDecoration(
                        color: goldColor,
                        shape: BoxShape.circle,
                        boxShadow: [
                          BoxShadow(color: goldColor, blurRadius: 4, spreadRadius: 1),
                        ],
                      ),
                    )
                  : const SizedBox(height: 4),
              ),
            ],
          ),
        ),
      ),
      ),
    );
  }
}
