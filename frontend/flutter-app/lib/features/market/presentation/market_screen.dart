import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/features/market/presentation/widgets/general_market_tab.dart';
import 'package:manpasik/features/market/presentation/widgets/market_product_card.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/breathing_glow.dart';
import 'package:manpasik/features/data_hub/presentation/widgets/ornate_gold_frame.dart';

// ───────────────────────────────────────────────────
// MarketScreen — Sanggam Orbit 카트리지 마켓
//
// [Rule 4] _backgroundColor/_goldColor 필드 → SanggamTheme.*
// [Rule 4] Theme(data:...) 래핑 제거 (항상 dark)
// [Rule 4] AppTheme.sanggamGold 6건 → SanggamTheme.primary
// [Rule 4] withOpacity ~15건 → withValues(alpha:)
// [Rule 4] jagae_pattern, app_theme 미사용 import 제거
// [Rule 4] 하드코딩 Color(0xFF...) 4건 → SanggamTheme 상수
// [Rule 2] spacing 20→16, 12→16, 4→8, borderRadius 12→16, h:44→48
// ───────────────────────────────────────────────────

class MarketScreen extends ConsumerStatefulWidget {
  const MarketScreen({super.key});

  @override
  ConsumerState<MarketScreen> createState() => _MarketScreenState();
}

class _MarketScreenState extends ConsumerState<MarketScreen> {
  String _selectedTier = 'all';
  final _couponController = TextEditingController();
  bool _isAnnual = false;

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        backgroundColor: Colors.transparent,
        body: Stack(
          children: [
            // 1. Deep Cosmic Background Gradient
            Positioned.fill(
              child: Container(
                decoration: BoxDecoration(
                  gradient: RadialGradient(
                    center: const Alignment(-0.8, -0.6),
                    radius: 1.5,
                    colors: [
                      SanggamTheme.surfaceVariant,
                      SanggamTheme.background,
                    ],
                  ),
                ),
              ),
            ),
            // 2. Ambient Gold Glow
            Positioned(
              top: 150,
              right: -100,
              child: Container(
                width: 300,
                height: 300,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  boxShadow: [
                    BoxShadow(
                      color:
                          SanggamTheme.primary.withValues(alpha: 0.15),
                      blurRadius: 100,
                      spreadRadius: 50,
                    ),
                  ],
                ),
              ),
            ),
            // Main Scrollable Content
            NestedScrollView(
              headerSliverBuilder: (context, innerBoxIsScrolled) {
                return [
                  _buildSliverAppBar(innerBoxIsScrolled),
                ];
              },
              body: Builder(
                builder: (context) {
                  return TabBarView(
                    children: [
                      // Tab 1: Cartridge Market
                      CustomScrollView(
                        key: const PageStorageKey('cartridge_market'),
                        physics: const BouncingScrollPhysics(
                            parent: AlwaysScrollableScrollPhysics()),
                        slivers: [
                          const SliverPadding(
                              padding: EdgeInsets.only(top: 120)),
                          SliverToBoxAdapter(
                              child: _buildTierFilter()),
                          SliverToBoxAdapter(
                              child: _buildSubscriptionBanner()),
                          SliverToBoxAdapter(
                              child: _buildCouponSection()),
                          _buildProductGrid(),
                          const SliverPadding(
                              padding: EdgeInsets.only(bottom: 24)),
                        ],
                      ),
                      // Tab 2: General Market
                      const GeneralMarketTab(),
                    ],
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSliverAppBar(bool innerBoxIsScrolled) {
    return SliverAppBar(
      expandedHeight: 120.0,
      floating: true,
      pinned: true,
      forceElevated: innerBoxIsScrolled,
      backgroundColor: SanggamTheme.background,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 48),
        title: const Text(
          '마켓',
          style: TextStyle(
            color: SanggamTheme.primary,
            fontWeight: FontWeight.bold,
            letterSpacing: 2.0,
            fontSize: 28,
          ),
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [
                SanggamTheme.background.withValues(alpha: 0.8),
                SanggamTheme.background,
              ],
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Opacity(
              opacity: 0.2,
              child: Icon(Icons.shopping_bag_outlined,
                  size: 100, color: SanggamTheme.primary),
            ),
          ),
        ),
      ),
      bottom: const TabBar(
        indicatorColor: SanggamTheme.primary,
        labelColor: SanggamTheme.primary,
        unselectedLabelColor: SanggamTheme.onSurfaceDim,
        indicatorWeight: 3,
        tabs: [
          Tab(text: '바이오 카트리지'),
          Tab(text: '헬스 리빙'),
        ],
      ),
      actions: [
        IconButton(
          icon: Icon(Icons.menu_book_outlined,
              color: Colors.white.withValues(alpha: 0.7)),
          tooltip: '카트리지 도감',
          onPressed: () => context.push('/market/encyclopedia'),
        ),
        IconButton(
          icon: Icon(Icons.receipt_long_outlined,
              color: Colors.white.withValues(alpha: 0.7)),
          tooltip: '주문 내역',
          onPressed: () => context.push('/market/orders'),
        ),
        IconButton(
          icon: Icon(Icons.search,
              color: Colors.white.withValues(alpha: 0.7)),
          tooltip: '검색',
          onPressed: () {},
        ),
        IconButton(
          icon: Icon(Icons.shopping_cart_outlined,
              color: Colors.white.withValues(alpha: 0.7)),
          tooltip: '장바구니',
          onPressed: () => context.push('/market/cart'),
        ),
      ],
    );
  }

  Widget _buildTierFilter() {
    final tiers = {
      'all': '전체',
      'Basic': '베이직',
      'Standard': '스탠다드',
      'Premium': '프리미엄',
      'Professional': '프로페셔널',
    };
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      padding: const EdgeInsets.all(16),
      child: Row(
        children: tiers.entries.map((entry) {
          final isSelected = _selectedTier == entry.key;
          return Padding(
            padding: const EdgeInsets.only(right: 8),
            child: ChoiceChip(
              label: Text(entry.value),
              selected: isSelected,
              onSelected: (selected) =>
                  setState(() => _selectedTier = entry.key),
              selectedColor:
                  SanggamTheme.primary.withValues(alpha: 0.2),
              backgroundColor:
                  Colors.white.withValues(alpha: 0.05),
              labelStyle: TextStyle(
                color: isSelected
                    ? SanggamTheme.primary
                    : SanggamTheme.onSurfaceDim,
                fontWeight: isSelected
                    ? FontWeight.bold
                    : FontWeight.normal,
              ),
              side: BorderSide(
                color: isSelected
                    ? SanggamTheme.primary
                    : Colors.white.withValues(alpha: 0.12),
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
        child: OrnateGoldFrame(
          width: double.infinity,
          isActive: true,
          child: Column(
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.1),
                      shape: BoxShape.circle,
                      border: Border.all(
                          color: SanggamTheme.primary, width: 1),
                    ),
                    child: const Icon(Icons.percent,
                        color: SanggamTheme.primary, size: 24),
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
                              fontSize: 16),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _isAnnual
                              ? '연간 구독 시 30% 할인 + 무료 배송'
                              : '최대 20% 할인 + 무료 배송',
                          style: TextStyle(
                              color:
                                  Colors.white.withValues(alpha: 0.8),
                              fontSize: 12),
                        ),
                      ],
                    ),
                  ),
                  const Icon(Icons.arrow_forward_ios,
                      color: SanggamTheme.primary, size: 16),
                ],
              ),
              const SizedBox(height: 16),
              // 연간/월간 토글
              Row(
                children: [
                  Expanded(
                    child: _buildToggle(
                      label: '월간',
                      isActive: !_isAnnual,
                      onTap: () => setState(() => _isAnnual = false),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: _buildToggle(
                      label: '연간 (-30%)',
                      isActive: _isAnnual,
                      onTap: () => setState(() => _isAnnual = true),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildToggle({
    required String label,
    required bool isActive,
    required VoidCallback onTap,
  }) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 8),
          decoration: BoxDecoration(
            color: isActive
                ? SanggamTheme.primary.withValues(alpha: 0.2)
                : Colors.transparent,
            borderRadius: BorderRadius.circular(8),
            border: Border.all(
              color: isActive
                  ? SanggamTheme.primary
                  : Colors.white.withValues(alpha: 0.24),
            ),
          ),
          child: Text(
            label,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: isActive
                  ? SanggamTheme.primary
                  : SanggamTheme.onSurfaceDim,
              fontWeight:
                  isActive ? FontWeight.bold : FontWeight.normal,
              fontSize: 13,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildCouponSection() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _couponController,
              style: const TextStyle(
                  color: Colors.white, fontSize: 14),
              decoration: InputDecoration(
                hintText: '쿠폰 코드 입력',
                hintStyle: TextStyle(
                    color: Colors.white.withValues(alpha: 0.4)),
                prefixIcon: Icon(Icons.confirmation_number_outlined,
                    color: SanggamTheme.primary
                        .withValues(alpha: 0.6),
                    size: 20),
                filled: true,
                fillColor: Colors.white.withValues(alpha: 0.05),
                contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16, vertical: 8),
                border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: BorderSide(
                        color:
                            Colors.white.withValues(alpha: 0.1))),
                enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: BorderSide(
                        color:
                            Colors.white.withValues(alpha: 0.1))),
                focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: const BorderSide(
                        color: SanggamTheme.primary)),
              ),
            ),
          ),
          const SizedBox(width: 8),
          SizedBox(
            height: 48,
            child: FilledButton(
              onPressed: () {
                if (_couponController.text.trim().isNotEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                          content: Text(
                              '쿠폰 "${_couponController.text.trim()}" 적용됨')));
                  _couponController.clear();
                }
              },
              style: FilledButton.styleFrom(
                  backgroundColor: SanggamTheme.primary,
                  foregroundColor: SanggamTheme.background,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16))),
              child: const Text('적용'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildProductGrid() {
    final tier = _selectedTier == 'all' ? null : _selectedTier;
    final productsAsync = ref.watch(cartridgeProductsProvider(tier));

    return productsAsync.when(
      data: (products) {
        if (products.isEmpty) {
          return _buildFallbackGrid();
        }
        return SliverPadding(
          padding: const EdgeInsets.all(16),
          sliver: SliverGrid(
            gridDelegate:
                const SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 280,
              childAspectRatio: 0.8,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                return AnimateFadeInUp(
                  duration: const Duration(milliseconds: 600),
                  delay: Duration(milliseconds: index * 50),
                  child:
                      MarketProductCard(product: products[index]),
                );
              },
              childCount: products.length,
            ),
          ),
        );
      },
      loading: () => const SliverToBoxAdapter(
          child: Center(child: CircularProgressIndicator())),
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
        gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
          maxCrossAxisExtent: 280,
          childAspectRatio: 0.72,
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
