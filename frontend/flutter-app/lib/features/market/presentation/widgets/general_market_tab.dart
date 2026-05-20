import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

// ───────────────────────────────────────────────────
// GeneralMarketTab — Sanggam Orbit 일반 마켓 탭
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppTheme.sanggamGold ~4x → SanggamTheme.primary
// [Rule 4] Color(0xFF1A1A1A) → Colors.white (다크 테마)
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.amber → SanggamTheme.primary
// [Rule 4] Color(0xFFE53935) → SanggamTheme.error
// [Rule 4] Color(0xFF1A1F35/0C101F) → SanggamTheme 상수
// ───────────────────────────────────────────────────

/// 일반 마켓 제품 모델
class GeneralProduct {
  final String id;
  final String name;
  final String description;
  final String category;
  final int price;
  final String imageUrl;
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

class GeneralMarketTab
    extends StatefulWidget {
  const GeneralMarketTab({super.key});

  @override
  State<GeneralMarketTab> createState() =>
      _GeneralMarketTabState();
}

class _GeneralMarketTabState
    extends State<GeneralMarketTab>
    with AutomaticKeepAliveClientMixin {
  @override
  bool get wantKeepAlive => true;

  final List<GeneralProduct> _products = [
    GeneralProduct(
        id: 'g1',
        name: '프리미엄 6년근 홍삼정',
        description: '면역력 증진 및 피로 개선',
        category: '건강기능식품',
        price: 128000,
        imageUrl:
            'assets/images/mock/ginseng.png',
        rating: 4.9,
        reviewCount: 1240,
        isBest: true),
    GeneralProduct(
        id: 'g2',
        name: '천연 비타민D 2000IU',
        description: '햇빛 에너지 충전, 뼈 건강',
        category: '비타민',
        price: 24000,
        imageUrl:
            'assets/images/mock/vitamin.png',
        rating: 4.8,
        reviewCount: 850,
        isBest: true),
    GeneralProduct(
        id: 'g3',
        name: '스마트 경추 베개',
        description: 'C커브 유지, 수면 퀄리티 개선',
        category: '헬스케어 기기',
        price: 89000,
        imageUrl:
            'assets/images/mock/pillow.png',
        rating: 4.7,
        reviewCount: 320),
    GeneralProduct(
        id: 'g4',
        name: '유기농 야채수 30포',
        description: '하루 한 팩으로 챙기는 건강',
        category: '건강음료',
        price: 35000,
        imageUrl:
            'assets/images/mock/juice.png',
        rating: 4.6,
        reviewCount: 512),
    GeneralProduct(
        id: 'g5',
        name: '루테인 지아잔틴',
        description: '눈 노화 케어, 침침한 눈',
        category: '영양제',
        price: 42000,
        imageUrl:
            'assets/images/mock/lutein.png',
        rating: 4.8,
        reviewCount: 930),
    GeneralProduct(
        id: 'g6',
        name: '저주파 마사지기',
        description: '뭉친 어깨와 근육통 완화',
        category: '헬스케어 기기',
        price: 55000,
        imageUrl:
            'assets/images/mock/massage.png',
        rating: 4.5,
        reviewCount: 220),
  ];

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return CustomScrollView(
      slivers: [
        const SliverPadding(
            padding:
                EdgeInsets.only(top: 120)),
        const SliverPadding(
            padding:
                EdgeInsets.only(top: 16)),
        SliverPadding(
          padding:
              const EdgeInsets.symmetric(
                  horizontal: 16),
          sliver: SliverGrid(
            gridDelegate:
                const SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 280,
              childAspectRatio: 0.65,
              crossAxisSpacing: 16,
              mainAxisSpacing: 16,
            ),
            delegate:
                SliverChildBuilderDelegate(
              (context, index) {
                return AnimateFadeInUp(
                  duration: const Duration(
                      milliseconds: 600),
                  delay: Duration(
                      milliseconds:
                          index * 50),
                  child:
                      _buildGeneralProductCard(
                          _products[
                              index]),
                );
              },
              childCount:
                  _products.length,
            ),
          ),
        ),
        const SliverPadding(
            padding:
                EdgeInsets.only(
                    bottom: 24)),
      ],
    );
  }

  Widget _buildGeneralProductCard(
      GeneralProduct product) {
    return ScaleButton(
      onPressed: () {},
      child: KoreanEdgeBorder(
        borderColor:
            SanggamTheme.primary,
        borderWidth: 1.5,
        child: ClipRRect(
          borderRadius:
              BorderRadius.circular(16),
          child: BackdropFilter(
            filter: ImageFilter.blur(
                sigmaX: 10, sigmaY: 10),
            child: JagaeContainer(
              padding:
                  const EdgeInsets.all(0),
              opacity: 0.25,
              showLattice: true,
              baseColor:
                  SanggamTheme.primary,
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    SanggamTheme.surface
                        .withValues(
                            alpha: 0.8),
                    SanggamTheme
                        .background
                        .withValues(
                            alpha: 0.9),
                  ],
                  begin:
                      Alignment.topLeft,
                  end: Alignment
                      .bottomRight,
                ),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black
                        .withValues(
                            alpha: 0.4),
                    blurRadius: 15,
                    offset:
                        const Offset(
                            0, 8),
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment:
                    CrossAxisAlignment
                        .start,
                children: [
                  // Image Area
                  Expanded(
                    flex: 6,
                    child: Stack(
                      children: [
                        Container(
                          width: double
                              .infinity,
                          decoration:
                              const BoxDecoration(
                            borderRadius:
                                BorderRadius
                                    .vertical(
                                        top: Radius.circular(
                                            16)),
                          ),
                          child: Center(
                            child: Icon(
                                Icons
                                    .shopping_bag_outlined,
                                size: 48,
                                color: SanggamTheme
                                    .primary
                                    .withValues(
                                        alpha:
                                            0.5)),
                          ),
                        ),
                        if (product
                            .isBest)
                          Positioned(
                            top: 12,
                            left: 12,
                            child:
                                Container(
                              padding: const EdgeInsets
                                  .symmetric(
                                  horizontal:
                                      8,
                                  vertical:
                                      4),
                              decoration:
                                  BoxDecoration(
                                color:
                                    SanggamTheme
                                        .error,
                                borderRadius:
                                    BorderRadius
                                        .circular(
                                            4),
                              ),
                              child:
                                  const Text(
                                '인기',
                                style:
                                    TextStyle(
                                  color: Colors
                                      .white,
                                  fontSize:
                                      10,
                                  fontWeight:
                                      FontWeight
                                          .bold,
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
                      padding:
                          const EdgeInsets
                              .all(12),
                      child: Column(
                        crossAxisAlignment:
                            CrossAxisAlignment
                                .start,
                        mainAxisAlignment:
                            MainAxisAlignment
                                .spaceBetween,
                        children: [
                          Column(
                            crossAxisAlignment:
                                CrossAxisAlignment
                                    .start,
                            children: [
                              Text(
                                product
                                    .category,
                                style:
                                    const TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim,
                                  fontSize:
                                      10,
                                ),
                              ),
                              const SizedBox(
                                  height:
                                      4),
                              Text(
                                product
                                    .name,
                                maxLines:
                                    2,
                                overflow:
                                    TextOverflow
                                        .ellipsis,
                                style:
                                    const TextStyle(
                                  color: Colors
                                      .white,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                  fontSize:
                                      14,
                                  height:
                                      1.2,
                                ),
                              ),
                            ],
                          ),
                          Row(
                            mainAxisAlignment:
                                MainAxisAlignment
                                    .spaceBetween,
                            crossAxisAlignment:
                                CrossAxisAlignment
                                    .end,
                            children: [
                              Text(
                                '₩${_formatPrice(product.price)}',
                                style:
                                    const TextStyle(
                                  color: SanggamTheme
                                      .primary,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                  fontSize:
                                      16,
                                ),
                              ),
                              Row(
                                children: [
                                  const Icon(
                                      Icons
                                          .star,
                                      size:
                                          12,
                                      color:
                                          SanggamTheme.primary),
                                  const SizedBox(
                                      width:
                                          2),
                                  Text(
                                    product
                                        .rating
                                        .toString(),
                                    style:
                                        const TextStyle(
                                      fontSize:
                                          10,
                                      color:
                                          SanggamTheme.onSurfaceDim,
                                    ),
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
          ),
        ),
      ),
    );
  }

  String _formatPrice(int price) {
    return price
        .toString()
        .replaceAllMapped(
      RegExp(
          r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (match) => '${match[1]},',
    );
  }
}
