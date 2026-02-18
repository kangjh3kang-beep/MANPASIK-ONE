import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/features/market/presentation/widgets/market_product_card.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 카트리지 마켓 화면 (Korean Futuristic Ver.)
class MarketScreen extends ConsumerStatefulWidget {
  const MarketScreen({super.key});

  @override
  ConsumerState<MarketScreen> createState() => _MarketScreenState();
}

class _MarketScreenState extends ConsumerState<MarketScreen> {
  String _selectedTier = 'all';

  // Premium Dark Color Palette
  final Color _backgroundColor = const Color(0xFF0A0E21); // Midnight Blue
  final Color _goldColor = const Color(0xFFD4AF37); // Sanggam Gold

  @override
  Widget build(BuildContext context) {
    // Override Theme to Dark Mode for this screen only to match "NanoBanana Pro" aesthetic
    return Theme(
      data: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: _backgroundColor,
        colorScheme: ColorScheme.dark(
          primary: _goldColor,
          surface: const Color(0xFF1A1F35),
        ),
      ),
      child: Scaffold(
        backgroundColor: Colors.transparent,
        body: CustomScrollView(
          physics: const BouncingScrollPhysics(),
          slivers: [
            _buildSliverAppBar(),
            SliverToBoxAdapter(child: _buildTierFilter()),
            SliverToBoxAdapter(child: _buildSubscriptionBanner()),
            _buildProductGrid(),
            const SliverPadding(padding: EdgeInsets.only(bottom: 24)),
          ],
        ),
      ),
    );
  }

  Widget _buildSliverAppBar() {
    return SliverAppBar(
      expandedHeight: 120.0,
      floating: false,
      pinned: true,
      backgroundColor: _backgroundColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Text(
          'Cartridge Market',
          style: TextStyle(
            color: _goldColor,
            fontWeight: FontWeight.bold,
            letterSpacing: 1.2,
          ),
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [
                _backgroundColor.withOpacity(0.8),
                _backgroundColor,
              ],
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Opacity(
              opacity: 0.2,
              child: Icon(Icons.shopping_bag_outlined, size: 100, color: _goldColor),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.menu_book_outlined, color: Colors.white70),
          tooltip: '카트리지 도감',
          onPressed: () => context.push('/market/encyclopedia'),
        ),
        IconButton(
          icon: const Icon(Icons.receipt_long_outlined, color: Colors.white70),
          tooltip: '주문 내역',
          onPressed: () => context.push('/market/orders'),
        ),
        IconButton(
          icon: const Icon(Icons.search, color: Colors.white70),
          onPressed: () {},
        ),
        IconButton(
          icon: const Icon(Icons.shopping_cart_outlined, color: Colors.white70),
          onPressed: () => context.push('/market/cart'),
        ),
      ],
    );
  }

  Widget _buildTierFilter() {
    final tiers = {
      'all': '전체',
      'Basic': 'Basic',
      'Standard': 'Standard',
      'Premium': 'Premium',
      'Professional': 'Professional',
    };
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 16),
      child: Row(
        children: tiers.entries.map((entry) {
          final isSelected = _selectedTier == entry.key;
          return Padding(
            padding: const EdgeInsets.only(right: 8),
            child: ChoiceChip(
              label: Text(entry.value),
              selected: isSelected,
              onSelected: (selected) => setState(() => _selectedTier = entry.key),
              selectedColor: _goldColor.withOpacity(0.2),
              backgroundColor: Colors.white.withOpacity(0.05),
              labelStyle: TextStyle(
                color: isSelected ? _goldColor : Colors.white60,
                fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              ),
              side: BorderSide(
                color: isSelected ? _goldColor : Colors.white12,
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildSubscriptionBanner() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: BreathingGlow(
        child: OrnateGoldFrame( // Upgraded to OrnateGoldFrame
          width: double.infinity,
          isActive: true,
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.1),
                  shape: BoxShape.circle,
                  border: Border.all(color: AppTheme.sanggamGold, width: 1),
                ),
                child: const Icon(Icons.percent, color: AppTheme.sanggamGold, size: 24),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      '정기 구독 멤버십',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '최대 20% 할인 + 무료 배송',
                      style: TextStyle(color: Colors.white.withOpacity(0.8), fontSize: 12),
                    ),
                  ],
                ),
              ),
              const Icon(Icons.arrow_forward_ios, color: AppTheme.sanggamGold, size: 16),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProductGrid() {
    final tier = _selectedTier == 'all' ? null : _selectedTier;
    final productsAsync = ref.watch(cartridgeProductsProvider(tier));

    return productsAsync.when(
      data: (products) {
        if (products.isEmpty) {
          return const SliverToBoxAdapter(
            child: SizedBox(
              height: 200,
              child: Center(child: Text('등록된 상품이 없습니다.', style: TextStyle(color: Colors.white54))),
            ),
          );
        }
        return SliverPadding(
          padding: const EdgeInsets.all(16),
          sliver: SliverGrid(
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 0.7,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                return AnimateFadeInUp(
                  duration: const Duration(milliseconds: 600),
                  delay: Duration(milliseconds: index * 50), // Staggered Effect
                  child: MarketProductCard(product: products[index]),
                );
              },
              childCount: products.length,
            ),
          ),
        );
      },
      loading: () => const SliverToBoxAdapter(child: Center(child: CircularProgressIndicator())),
      error: (_, __) => _buildFallbackGrid(),
    );
  }

  Widget _buildFallbackGrid() {
    final fallback = [
      CartridgeProduct(id: '1', typeCode: '0x01', nameKo: '혈당 카트리지', nameEn: 'Glucose', tier: 'Basic', price: 15000, unit: 'mg/dL', referenceRange: '70-100', requiredChannels: 1, measurementSecs: 60, isAvailable: true),
      CartridgeProduct(id: '2', typeCode: '0x02', nameKo: '당화혈색소', nameEn: 'HbA1c', tier: 'Standard', price: 25000, unit: '%', referenceRange: '4.0-5.6', requiredChannels: 2, measurementSecs: 90, isAvailable: true),
      CartridgeProduct(id: '3', typeCode: '0x03', nameKo: '요산', nameEn: 'Uric Acid', tier: 'Basic', price: 18000, unit: 'mg/dL', referenceRange: '3.0-7.0', requiredChannels: 1, measurementSecs: 60, isAvailable: true),
      CartridgeProduct(id: '4', typeCode: '0x05', nameKo: '비타민D', nameEn: 'Vitamin D', tier: 'Premium', price: 35000, unit: 'ng/mL', referenceRange: '30-100', requiredChannels: 3, measurementSecs: 120, isAvailable: true),
    ];

    return SliverPadding(
      padding: const EdgeInsets.all(16),
      sliver: SliverGrid(
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          childAspectRatio: 0.7,
          crossAxisSpacing: 16,
          mainAxisSpacing: 16,
        ),
        delegate: SliverChildBuilderDelegate(
          (context, index) {
            return AnimateFadeInUp(
              child: MarketProductCard(product: fallback[index]),
            );
          },
          childCount: fallback.length,
        ),
      ),
    );
  }
}
