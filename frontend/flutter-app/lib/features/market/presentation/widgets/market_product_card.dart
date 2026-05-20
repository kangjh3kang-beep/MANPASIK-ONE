import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/jagae_card.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

class MarketProductCard extends StatelessWidget {
  final CartridgeProduct product;

  const MarketProductCard({super.key, required this.product});

  @override
  Widget build(BuildContext context) {
    final tierColor = _getTierColor(product.tier);
    final icon = _getProductIcon(product.typeCode);

    return ScaleButton(
      onPressed: () => context.push('/market/product/${product.id}'),
      child: KoreanEdgeBorder(
        borderColor: tierColor,
        borderWidth: 1.5,
      child: JagaeCard(
        padding: const EdgeInsets.all(16),
        opacity: 0.2,
        borderRadius: BorderRadius.circular(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _buildTierBadge(product.tier, tierColor),
                const Icon(
                  Icons.favorite_border,
                  color: Colors.white38,
                  size: 20,
                ),
              ],
            ),
            const Spacer(),
            Center(child: _buildFloatingCartridgeObject(icon, tierColor)),
            const Spacer(),
            // Phase 6 Parallax Typography: Bold & tight letter spacing (-0.5)
            Text(
              product.nameKo,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 22,
                fontWeight: FontWeight.w900,
                letterSpacing: -0.5,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 4),
            Text(
              product.nameEn,
              style: const TextStyle(
                color: Colors.white70,
                fontSize: 14,
                fontWeight: FontWeight.w600,
                letterSpacing: -0.3,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 16),
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text(
                  '₩${_formatPrice(product.price)}',
                  style: const TextStyle(
                    color: Color(0xFFD4AF37),
                    fontSize: 22,
                    fontWeight: FontWeight.w800,
                    fontFamily: 'Roboto',
                  ),
                ),
                if (product.unit.isNotEmpty) ...[
                  const SizedBox(width: 4),
                  Text(
                    '/ ${product.unit}',
                    style: const TextStyle(
                      color: Colors.white60,
                      fontSize: 13,
                    ),
                  ),
                ],
              ],
            ),
          ],
        ),
      ),
      ),
    );
  }

  Widget _buildFloatingCartridgeObject(IconData icon, Color color) {
    return Container(
      width: 130,
      height: 130,
      decoration: BoxDecoration(
        color: const Color(0xFF080C17),
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: color.withValues(alpha: 0.5),
            blurRadius: 30,
            spreadRadius: -5,
            offset: const Offset(0, 15),
          ),
          BoxShadow(
            color: color.withValues(alpha: 0.2),
            blurRadius: 50,
            spreadRadius: 15,
          ),
        ],
        border: Border.all(color: color.withValues(alpha: 0.8), width: 2.5),
      ),
      child: Stack(
        alignment: Alignment.center,
        clipBehavior: Clip.none, // Phase 6 Parallax: Allow content to overflow slightly
        children: [
          // The icon itself is shifted slightly upwards to simulate 3D popping out of the card
          Positioned(
            top: -10, // Parallax shift
            child: Icon(icon, color: color, size: 70), // Slightly larger
          ),
          Positioned(
            top: 10,
            left: 20,
            child: Container(
              width: 25,
              height: 25,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [Colors.white.withValues(alpha: 0.8), Colors.transparent],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTierBadge(String tier, Color color) {
    final localizedTier = _localizeTier(tier);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: color.withValues(alpha: 0.4)),
      ),
      child: Text(
        localizedTier,
        style: TextStyle(
          color: color,
          fontSize: 10,
          fontWeight: FontWeight.bold,
          letterSpacing: 0.5,
        ),
      ),
    );
  }

  String _localizeTier(String tier) {
    switch (tier) {
      case 'Premium':
        return '프리미엄';
      case 'Professional':
        return '프로페셔널';
      case 'Standard':
        return '스탠다드';
      case 'Basic':
        return '베이직';
      default:
        return tier;
    }
  }

  Color _getTierColor(String tier) {
    switch (tier) {
      case 'Premium':
        return const Color(0xFFFFD700);
      case 'Professional':
        return const Color(0xFFE0FFFF);
      case 'Standard':
        return const Color(0xFF00BFFF);
      default:
        return const Color(0xFFAAAAAA);
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
