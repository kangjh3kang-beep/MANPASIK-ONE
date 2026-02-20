import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';

class MarketProductCard extends StatelessWidget {
  final CartridgeProduct product;

  const MarketProductCard({super.key, required this.product});

  @override
  Widget build(BuildContext context) {
    // NanoBanana Pro Theme Colors (Local Override for Premium Feel)
    final tierColor = _getTierColor(product.tier);
    final icon = _getProductIcon(product.typeCode);
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return ScaleButton(
      onPressed: () => context.push('/market/product/${product.id}'),
      child: PorcelainContainer(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Icon & Tier Badge
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _buildGlowingIcon(icon, tierColor),
                _buildTierBadge(product.tier, tierColor),
              ],
            ),
            const Spacer(),
            // Product Name
            Text(
              product.nameKo,
              style: TextStyle(
                color: isDark ? Colors.white : const Color(0xFF1A1A1A), // Adaptive Color
                fontSize: 16,
                fontWeight: FontWeight.bold,
                letterSpacing: -0.5,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),
            // English Name
            Text(
              product.nameEn,
              style: TextStyle(
                color: isDark ? Colors.white70 : const Color(0xFF1A1A1A).withOpacity(0.6),
                fontSize: 12,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 12),
            // Price
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '₩${_formatPrice(product.price)}',
                  style: const TextStyle(
                    color: Color(0xFFD4AF37), // Omni Gold
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                    fontFamily: 'Roboto', 
                  ),
                ),
                if (product.unit.isNotEmpty) ...[
                  const SizedBox(width: 4),
                  Text(
                    '/ ${product.unit}',
                    style: TextStyle(
                      color: const Color(0xFF1A1A1A).withOpacity(0.4),
                      fontSize: 11,
                    ),
                  ),
                ],
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildGlowingIcon(IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(10),
      decoration: BoxDecoration(
        color: color.withOpacity(0.15),
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: color.withOpacity(0.4),
            blurRadius: 12,
            spreadRadius: -2,
          ),
        ],
        border: Border.all(color: color.withOpacity(0.3), width: 1.5),
      ),
      child: Icon(icon, color: color, size: 24),
    );
  }

  Widget _buildTierBadge(String tier, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withOpacity(0.3)),
      ),
      child: Text(
        tier.toUpperCase(),
        style: TextStyle(
          color: color,
          fontSize: 10,
          fontWeight: FontWeight.bold,
          letterSpacing: 0.5,
        ),
      ),
    );
  }

  Color _getTierColor(String tier) {
    switch (tier) {
      case 'Premium':
        return const Color(0xFFFFD700); // Gold
      case 'Professional':
        return const Color(0xFFE0FFFF); // Cyan/Diamond
      case 'Standard':
        return const Color(0xFF00BFFF); // Deep Sky Blue
      default:
        return const Color(0xFFAAAAAA); // Silver
    }
  }

  IconData _getProductIcon(String typeCode) {
    switch (typeCode.toLowerCase()) {
      case '0x01':
      case 'glucose':
        return Icons.water_drop_rounded;
      case '0x02':
      case 'hba1c':
        return Icons.timelapse_rounded;
      case '0x05':
      case 'vitamin_d':
        return Icons.wb_sunny_rounded;
      default:
        return Icons.science_rounded;
    }
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (match) => '${match[1]},',
    );
  }
}
