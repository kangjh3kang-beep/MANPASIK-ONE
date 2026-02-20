import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

/// 일반 마켓 제품 모델 (Mock Data용)
class GeneralProduct {
  final String id;
  final String name;
  final String description;
  final String category;
  final int price;
  final String imageUrl; // Asset path or IconData for now
  final double rating;
  final int reviewCount;
  final bool isBest;

  const GeneralProduct({
    required this.id,
    required this.name,
    required this.description,
    required this.category,
    required this.price,
    required this.imageUrl,
    required this.rating,
    required this.reviewCount,
    this.isBest = false,
  });
}

class GeneralMarketTab extends StatefulWidget {
  const GeneralMarketTab({super.key});

  @override
  State<GeneralMarketTab> createState() => _GeneralMarketTabState();
}

class _GeneralMarketTabState extends State<GeneralMarketTab> with AutomaticKeepAliveClientMixin {
  @override
  bool get wantKeepAlive => true;

  // Mock Data Definition
  final List<GeneralProduct> _products = [
    GeneralProduct(
        id: 'g1',
        name: '프리미엄 6년근 홍삼정',
        description: '면역력 증진 및 피로 개선',
        category: '건강기능식품',
        price: 128000,
        imageUrl: 'assets/images/mock/ginseng.png', 
        rating: 4.9,
        reviewCount: 1240,
        isBest: true),
    GeneralProduct(
        id: 'g2',
        name: '천연 비타민D 2000IU',
        description: '햇빛 에너지 충전, 뼈 건강',
        category: '비타민',
        price: 24000,
        imageUrl: 'assets/images/mock/vitamin.png',
        rating: 4.8,
        reviewCount: 850,
        isBest: true),
    GeneralProduct(
        id: 'g3',
        name: '스마트 경추 베개',
        description: 'C커브 유지, 수면 퀄리티 개선',
        category: '헬스케어 기기',
        price: 89000,
        imageUrl: 'assets/images/mock/pillow.png',
        rating: 4.7,
        reviewCount: 320),
    GeneralProduct(
        id: 'g4',
        name: '유기농 야채수 30포',
        description: '하루 한 팩으로 챙기는 건강',
        category: '건강음료',
        price: 35000,
        imageUrl: 'assets/images/mock/juice.png',
        rating: 4.6,
        reviewCount: 512),
    GeneralProduct(
        id: 'g5',
        name: '루테인 지아잔틴',
        description: '눈 노화 케어, 침침한 눈',
        category: '영양제',
        price: 42000,
        imageUrl: 'assets/images/mock/lutein.png',
        rating: 4.8,
        reviewCount: 930),
    GeneralProduct(
        id: 'g6',
        name: '저주파 마사지기',
        description: '뭉친 어깨와 근육통 완화',
        category: '헬스케어 기기',
        price: 55000,
        imageUrl: 'assets/images/mock/massage.png',
        rating: 4.5,
        reviewCount: 220),
  ];

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return CustomScrollView(
      slivers: [
        // Removed SliverOverlapInjector due to stability issues
        const SliverPadding(padding: EdgeInsets.only(top: 120)), // Safe manual padding
        const SliverPadding(padding: EdgeInsets.only(top: 16)),
        SliverPadding(
          padding: const EdgeInsets.symmetric(horizontal: 16),
          sliver: SliverGrid(
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              childAspectRatio: 0.65, // Taller for description
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            delegate: SliverChildBuilderDelegate(
              (context, index) {
                return AnimateFadeInUp(
                  duration: const Duration(milliseconds: 600),
                  delay: Duration(milliseconds: index * 50),
                  child: _buildGeneralProductCard(_products[index]),
                );
              },
              childCount: _products.length,
            ),
          ),
        ),
        const SliverPadding(padding: EdgeInsets.only(bottom: 24)),
      ],
    );
  }

  Widget _buildGeneralProductCard(GeneralProduct product) {
    return ScaleButton(
      onPressed: () {
        // Navigate to details (To be implemented)
      },
      child: PorcelainContainer(
        padding: const EdgeInsets.all(0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Image Area
            Expanded(
              flex: 6,
              child: Stack(
                children: [
                  Container(
                    width: double.infinity,
                    decoration: BoxDecoration(
                      color: Colors.grey.withOpacity(0.1),
                      borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
                    ),
                    child: Center(
                      child: Icon(Icons.shopping_bag_outlined, 
                        size: 48, 
                        color: AppTheme.sanggamGold.withOpacity(0.5)
                      ), // Placeholder
                    ),
                  ),
                  if (product.isBest)
                    Positioned(
                      top: 12,
                      left: 12,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: const Color(0xFFE53935), // Dancheong Red
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: const Text(
                          'BEST',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
            // Content Area
            Expanded(
              flex: 4,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          product.category,
                          style: TextStyle(
                            color: Colors.grey[500],
                            fontSize: 10,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          product.name,
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                          style: const TextStyle(
                            color: Color(0xFF1A1A1A),
                            fontWeight: FontWeight.bold,
                            fontSize: 14,
                            height: 1.2,
                          ),
                        ),
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Text(
                          '₩${_formatPrice(product.price)}',
                          style: const TextStyle(
                            color: AppTheme.sanggamGold,
                            fontWeight: FontWeight.bold,
                            fontSize: 16,
                          ),
                        ),
                        Row(
                          children: [
                            const Icon(Icons.star, size: 12, color: Colors.amber),
                            const SizedBox(width: 2),
                            Text(
                              product.rating.toString(),
                              style: const TextStyle(fontSize: 10, color: Colors.grey),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (match) => '${match[1]},',
    );
  }
}
