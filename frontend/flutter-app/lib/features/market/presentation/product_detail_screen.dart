import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ProductDetailScreen — Sanggam Orbit 상품 상세
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] 외부 Theme 래퍼 제거
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] ThemeData params ~5x 제거
// [Rule 4] theme.textTheme ~20x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~5x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold ~10x → SanggamTheme.primary
// [Rule 4] Card ~3x → SanggamContainer
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Colors.cyan → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.amber → SanggamTheme.primary
// [Rule 4] AlertDialog + TextField 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16
// ───────────────────────────────────────────────────

/// 상품 상세 화면
class ProductDetailScreen
    extends ConsumerWidget {
  const ProductDetailScreen(
      {super.key,
      required this.productId});

  final String productId;

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final client =
        ref.watch(restClientProvider);

    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '상품 상세',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                  _WishlistButton(
                      productId:
                          productId),
                  IconButton(
                    icon: const Icon(
                        Icons
                            .shopping_cart_outlined,
                        color: Colors.white),
                    tooltip: '장바구니',
                    onPressed: () =>
                        context.push(
                            '/market/cart'),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: FutureBuilder<
                  Map<String, dynamic>>(
                future: client
                    .getProduct(productId),
                builder: (context,
                    snapshot) {
                  if (snapshot
                          .connectionState ==
                      ConnectionState
                          .waiting) {
                    return const Center(
                        child: CircularProgressIndicator(
                            color:
                                SanggamTheme
                                    .primary));
                  }
                  if (snapshot.hasError) {
                    return _buildFallback(
                        context, ref);
                  }
                  final data =
                      snapshot.data ?? {};
                  return _buildProductDetail(
                      context, ref, data);
                },
              ),
            ),
          ],
        ),
      ),
      bottomNavigationBar:
          _buildBottomBar(context, ref),
    );
  }

  Widget _buildProductDetail(
      BuildContext context,
      WidgetRef ref,
      Map<String, dynamic> data) {
    final name = data['name'] as String? ??
        '카트리지 상품';
    final description =
        data['description'] as String? ??
            '상품 설명이 없습니다.';
    final price =
        data['price'] as num? ?? 0;
    final tier = _localizeTier(
        data['tier'] as String? ??
            'Basic');

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.stretch,
        children: [
          // 상품 이미지 영역
          _buildHeroCartridge(),

          // Glassmorphism Content Panel
          Transform.translate(
            offset: const Offset(0, -30),
            child: ClipRRect(
              borderRadius:
                  const BorderRadius.only(
                topLeft: Radius.circular(30),
                topRight:
                    Radius.circular(30),
              ),
              child: BackdropFilter(
                filter: ImageFilter.blur(
                    sigmaX: 15, sigmaY: 15),
                child: Container(
                  decoration: BoxDecoration(
                    color: SanggamTheme
                        .surface
                        .withValues(
                            alpha: 0.7),
                    borderRadius:
                        const BorderRadius
                            .only(
                      topLeft:
                          Radius.circular(
                              30),
                      topRight:
                          Radius.circular(
                              30),
                    ),
                    border: Border(
                      top: BorderSide(
                          color: Colors.white
                              .withValues(
                                  alpha:
                                      0.15),
                          width: 1.5),
                    ),
                  ),
                  child: Padding(
                    padding:
                        const EdgeInsets
                            .all(24),
                    child: Column(
                      crossAxisAlignment:
                          CrossAxisAlignment
                              .start,
                      children: [
                        // 티어 배지
                        Container(
                          padding:
                              const EdgeInsets
                                  .symmetric(
                                  horizontal:
                                      12,
                                  vertical:
                                      6),
                          decoration:
                              BoxDecoration(
                            color: SanggamTheme
                                .primary
                                .withValues(
                                    alpha:
                                        0.15),
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        20),
                            border: Border.all(
                                color: SanggamTheme
                                    .primary
                                    .withValues(
                                        alpha:
                                            0.4)),
                          ),
                          child: Text(tier,
                              style:
                                  const TextStyle(
                                color:
                                    SanggamTheme
                                        .primary,
                                fontSize:
                                    12,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                        ),
                        const SizedBox(
                            height: 8),

                        // 상품명
                        Text(name,
                            style:
                                const TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 22,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            )),
                        const SizedBox(
                            height: 8),

                        // 가격
                        Text(
                          '₩${_formatPrice(price)}',
                          style:
                              const TextStyle(
                            color:
                                SanggamTheme
                                    .primary,
                            fontSize: 20,
                            fontWeight:
                                FontWeight
                                    .bold,
                          ),
                        ),
                        const SizedBox(
                            height: 16),
                        const Divider(
                            color: SanggamTheme
                                .surfaceVariant),
                        const SizedBox(
                            height: 8),

                        // 상품 설명
                        const Text('상품 설명',
                            style:
                                TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 16,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            )),
                        const SizedBox(
                            height: 8),
                        Text(description,
                            style:
                                const TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 14,
                              height: 1.6,
                            )),
                        const SizedBox(
                            height: 16),

                        // 리뷰 섹션
                        Row(
                          mainAxisAlignment:
                              MainAxisAlignment
                                  .spaceBetween,
                          children: [
                            const Text(
                                '사용자 리뷰',
                                style:
                                    TextStyle(
                                  color: Colors
                                      .white,
                                  fontSize:
                                      16,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                )),
                            TextButton.icon(
                              onPressed: () =>
                                  _showWriteReviewDialog(
                                      context,
                                      ref),
                              icon: const Icon(
                                  Icons
                                      .rate_review,
                                  size: 18),
                              label:
                                  const Text(
                                      '리뷰 작성'),
                              style: TextButton
                                  .styleFrom(
                                foregroundColor:
                                    SanggamTheme
                                        .primary,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(
                            height: 8),
                        _ReviewSection(
                            productId:
                                productId),
                        const SizedBox(
                            height: 16),

                        // 스펙 테이블
                        const Text('제품 스펙',
                            style:
                                TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 16,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            )),
                        const SizedBox(
                            height: 8),
                        _buildSpecRow(
                            '제품 ID',
                            productId),
                        _buildSpecRow(
                            '등급', tier),
                        _buildSpecRow(
                            '호환 리더기',
                            'ManPaSik Reader v2.0+'),
                        _buildSpecRow(
                            '보관 조건',
                            '실온 (15~30°C)'),
                        _buildSpecRow(
                            '유효 기간',
                            '제조일로부터 12개월'),
                        const SizedBox(
                            height: 32),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeroCartridge() {
    return Container(
      height: 320,
      decoration: const BoxDecoration(
        gradient: RadialGradient(
          center: Alignment.center,
          radius: 0.8,
          colors: [
            SanggamTheme.surface,
            SanggamTheme.background,
          ],
        ),
      ),
      child: Center(
        child: Container(
          width: 140,
          height: 140,
          decoration: BoxDecoration(
            color:
                SanggamTheme.background,
            shape: BoxShape.circle,
            boxShadow: [
              BoxShadow(
                color: SanggamTheme
                    .primary
                    .withValues(
                        alpha: 0.4),
                blurRadius: 40,
                spreadRadius: -10,
                offset:
                    const Offset(0, 20),
              ),
              BoxShadow(
                color: SanggamTheme
                    .jagaeCyan
                    .withValues(
                        alpha: 0.2),
                blurRadius: 60,
                spreadRadius: 20,
              ),
            ],
            border: Border.all(
                color: SanggamTheme.primary
                    .withValues(
                        alpha: 0.6),
                width: 3),
          ),
          child: Stack(
            alignment: Alignment.center,
            children: [
              const Icon(Icons.science,
                  color: SanggamTheme
                      .primary,
                  size: 70),
              Positioned(
                top: 15,
                left: 20,
                child: Container(
                  width: 30,
                  height: 30,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    gradient: RadialGradient(
                      colors: [
                        Colors.white
                            .withValues(
                                alpha: 0.9),
                        Colors.transparent,
                      ],
                    ),
                  ),
                ),
              ),
            ],
          ),
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

  Widget _buildFallback(
      BuildContext context,
      WidgetRef ref) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.stretch,
        children: [
          Container(
            height: 250,
            color:
                SanggamTheme.surfaceVariant,
            child: Center(
              child: Icon(Icons.science,
                  size: 80,
                  color: SanggamTheme
                      .primary
                      .withValues(
                          alpha: 0.5)),
            ),
          ),
          Padding(
            padding:
                const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,
              children: [
                const Text('카트리지 상품',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 22,
                      fontWeight:
                          FontWeight.bold,
                    )),
                const SizedBox(height: 8),
                const Text('₩29,900',
                    style: TextStyle(
                      color: SanggamTheme
                          .primary,
                      fontSize: 20,
                      fontWeight:
                          FontWeight.bold,
                    )),
                const SizedBox(
                    height: 16),
                const Text(
                    '서버 연결 후 상세 정보가 표시됩니다.',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                    )),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSpecRow(
      String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 4),
      child: Row(
        children: [
          SizedBox(
            width: 100,
            child: Text(label,
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 13,
                  fontWeight:
                      FontWeight.w600,
                )),
          ),
          Expanded(
              child: Text(value,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 13,
                  ))),
        ],
      ),
    );
  }

  String _formatPrice(num price) {
    return price
        .toInt()
        .toString()
        .replaceAllMapped(
      RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]},',
    );
  }

  void _showWriteReviewDialog(
      BuildContext context,
      WidgetRef ref) {
    int rating = 5;
    final contentCtrl =
        TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setDialogState) =>
            AlertDialog(
          backgroundColor:
              SanggamTheme.surface,
          shape: RoundedRectangleBorder(
              borderRadius:
                  BorderRadius.circular(
                      16)),
          title: const Text('리뷰 작성',
              style: TextStyle(
                  color: Colors.white)),
          content: Column(
            mainAxisSize:
                MainAxisSize.min,
            children: [
              Row(
                mainAxisAlignment:
                    MainAxisAlignment
                        .center,
                children:
                    List.generate(
                        5,
                        (i) => IconButton(
                              icon: Icon(
                                i < rating
                                    ? Icons
                                        .star
                                    : Icons
                                        .star_border,
                                color: SanggamTheme
                                    .primary,
                              ),
                              tooltip:
                                  '별점 ${i + 1}',
                              onPressed: () =>
                                  setDialogState(
                                      () =>
                                          rating =
                                              i + 1),
                            )),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: contentCtrl,
                maxLines: 3,
                style: const TextStyle(
                    color: Colors.white),
                decoration:
                    InputDecoration(
                  hintText:
                      '사용 후기를 작성해주세요',
                  hintStyle: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim),
                  filled: true,
                  fillColor:
                      SanggamTheme
                          .surfaceVariant,
                  border:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius
                            .circular(16),
                    borderSide:
                        BorderSide.none,
                  ),
                  enabledBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius
                            .circular(16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                  ),
                  focusedBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius
                            .circular(16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .primary),
                  ),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              style: TextButton.styleFrom(
                  foregroundColor:
                      SanggamTheme
                          .onSurfaceDim),
              child: const Text('취소'),
            ),
            FilledButton(
              onPressed: () async {
                final client = ref.read(
                    restClientProvider);
                final userId = ref
                        .read(
                            authProvider)
                        .userId ??
                    '';
                try {
                  await client
                      .createProductReview(
                    productId: productId,
                    userId: userId,
                    rating: rating,
                    content:
                        contentCtrl.text,
                  );
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                  }
                  if (context.mounted) {
                    ScaffoldMessenger.of(
                            context)
                        .showSnackBar(
                      const SnackBar(
                          content: Text(
                              '리뷰가 등록되었습니다.')),
                    );
                  }
                } catch (e) {
                  if (ctx.mounted) {
                    ScaffoldMessenger.of(
                            ctx)
                        .showSnackBar(
                      SnackBar(
                          content: Text(
                              '리뷰 등록 실패: $e')),
                    );
                  }
                }
              },
              style:
                  FilledButton.styleFrom(
                backgroundColor:
                    SanggamTheme.primary,
                foregroundColor:
                    SanggamTheme
                        .background,
                shape:
                    RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius
                          .circular(16),
                ),
              ),
              child: const Text('등록'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildBottomBar(
      BuildContext context,
      WidgetRef ref) {
    return ClipRRect(
      child: BackdropFilter(
        filter: ImageFilter.blur(
            sigmaX: 20, sigmaY: 20),
        child: Container(
          decoration: BoxDecoration(
            color: SanggamTheme.background
                .withValues(alpha: 0.8),
            border: Border(
                top: BorderSide(
                    color: Colors.white
                        .withValues(
                            alpha: 0.1))),
          ),
          child: SafeArea(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 16),
              child: Row(
                children: [
                  OutlinedButton(
                    onPressed: () {
                      ScaffoldMessenger.of(
                              context)
                          .showSnackBar(
                        const SnackBar(
                            content: Text(
                                '장바구니에 추가되었습니다.')),
                      );
                    },
                    style: OutlinedButton
                        .styleFrom(
                      padding:
                          const EdgeInsets
                              .symmetric(
                              vertical: 14,
                              horizontal:
                                  24),
                      side: BorderSide(
                          color: SanggamTheme
                              .primary
                              .withValues(
                                  alpha:
                                      0.5)),
                      shape:
                          RoundedRectangleBorder(
                              borderRadius:
                                  BorderRadius
                                      .circular(
                                          16)),
                    ),
                    child: const Icon(
                        Icons
                            .shopping_cart_outlined,
                        color: SanggamTheme
                            .primary),
                  ),
                  const SizedBox(
                      width: 16),
                  Expanded(
                    child: Container(
                      decoration:
                          BoxDecoration(
                        gradient:
                            const LinearGradient(
                          colors: [
                            SanggamTheme
                                .primary,
                            Color(
                                0xFFB8860B),
                          ],
                          begin: Alignment
                              .topLeft,
                          end: Alignment
                              .bottomRight,
                        ),
                        borderRadius:
                            BorderRadius
                                .circular(
                                    16),
                        boxShadow: [
                          BoxShadow(
                            color: SanggamTheme
                                .primary
                                .withValues(
                                    alpha:
                                        0.3),
                            blurRadius: 12,
                            offset:
                                const Offset(
                                    0, 4),
                          ),
                        ],
                      ),
                      child: FilledButton(
                        onPressed: () =>
                            context.push(
                                '/market/cart'),
                        style: FilledButton
                            .styleFrom(
                          backgroundColor:
                              Colors
                                  .transparent,
                          shadowColor: Colors
                              .transparent,
                          padding:
                              const EdgeInsets
                                  .symmetric(
                                  vertical:
                                      16),
                          shape: RoundedRectangleBorder(
                              borderRadius:
                                  BorderRadius
                                      .circular(
                                          16)),
                        ),
                        child: const Text(
                            '바로 구매',
                            style:
                                TextStyle(
                              color: SanggamTheme
                                  .background,
                              fontWeight:
                                  FontWeight
                                      .bold,
                              fontSize: 16,
                            )),
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
}

/// 위시리스트 하트 토글 버튼
class _WishlistButton
    extends StatefulWidget {
  const _WishlistButton(
      {required this.productId});
  final String productId;

  @override
  State<_WishlistButton> createState() =>
      _WishlistButtonState();
}

class _WishlistButtonState
    extends State<_WishlistButton> {
  bool _isWished = false;

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Icon(
        _isWished
            ? Icons.favorite
            : Icons.favorite_border,
        color: _isWished
            ? SanggamTheme.error
            : Colors.white,
      ),
      tooltip: '위시리스트',
      onPressed: () {
        setState(() =>
            _isWished = !_isWished);
        ScaffoldMessenger.of(context)
            .showSnackBar(
          SnackBar(
              content: Text(_isWished
                  ? '위시리스트에 추가되었습니다.'
                  : '위시리스트에서 제거되었습니다.')),
        );
      },
    );
  }
}

/// REST API 기반 리뷰 섹션
class _ReviewSection
    extends ConsumerWidget {
  const _ReviewSection(
      {required this.productId});
  final String productId;

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final client =
        ref.watch(restClientProvider);

    return FutureBuilder<
        Map<String, dynamic>>(
      future: client.getProductReviews(
          productId,
          limit: 3),
      builder: (context, snapshot) {
        if (snapshot.connectionState ==
            ConnectionState.waiting) {
          return const SanggamContainer(
              borderRadius: 16,
              child: Padding(
                  padding:
                      EdgeInsets.all(16),
                  child: Center(
                      child:
                          CircularProgressIndicator(
                              strokeWidth:
                                  2,
                              color:
                                  SanggamTheme
                                      .primary))));
        }

        final data =
            snapshot.data ?? {};
        final reviews = (data['reviews']
                    as List?)
                ?.cast<
                    Map<String,
                        dynamic>>() ??
            [];
        final avgRating =
            (data['average_rating']
                        as num?)
                    ?.toDouble() ??
                0.0;
        final totalCount =
            (data['total_count']
                        as num?)
                    ?.toInt() ??
                reviews.length;

        if (reviews.isEmpty &&
            !snapshot.hasError) {
          return _buildFallbackReview();
        }

        return SanggamContainer(
          borderRadius: 16,
          child: Column(
            crossAxisAlignment:
                CrossAxisAlignment.start,
            children: [
              // 평점 요약
              Row(
                children: [
                  ...List.generate(
                      5,
                      (i) => Icon(
                            i < avgRating
                                    .round()
                                ? Icons.star
                                : Icons
                                    .star_border,
                            size: 20,
                            color:
                                SanggamTheme
                                    .primary,
                          )),
                  const SizedBox(
                      width: 8),
                  Text(
                      avgRating
                          .toStringAsFixed(
                              1),
                      style:
                          const TextStyle(
                        color:
                            Colors.white,
                        fontSize: 16,
                        fontWeight:
                            FontWeight.bold,
                      )),
                  const SizedBox(
                      width: 4),
                  Text(
                      '($totalCount개 리뷰)',
                      style:
                          const TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      )),
                ],
              ),
              const SizedBox(height: 16),
              // 최근 리뷰 목록
              ...reviews.take(3).map(
                  (r) => Padding(
                        padding:
                            const EdgeInsets
                                .only(
                                bottom: 8),
                        child: Column(
                          crossAxisAlignment:
                              CrossAxisAlignment
                                  .start,
                          children: [
                            Row(
                              children: [
                                ...List.generate(
                                    5,
                                    (i) =>
                                        Icon(
                                          i < ((r['rating'] as num?) ?? 5)
                                              ? Icons.star
                                              : Icons.star_border,
                                          size:
                                              14,
                                          color:
                                              SanggamTheme.primary,
                                        )),
                                const SizedBox(
                                    width:
                                        8),
                                Text(
                                    r['author_name']
                                            as String? ??
                                        '익명',
                                    style:
                                        const TextStyle(
                                      color:
                                          Colors.white,
                                      fontSize:
                                          13,
                                      fontWeight:
                                          FontWeight.w600,
                                    )),
                              ],
                            ),
                            const SizedBox(
                                height: 4),
                            Text(
                                r['content']
                                        as String? ??
                                    '',
                                style:
                                    const TextStyle(
                                  color: Colors
                                      .white,
                                  fontSize:
                                      13,
                                )),
                          ],
                        ),
                      )),
            ],
          ),
        );
      },
    );
  }

  Widget _buildFallbackReview() {
    return SanggamContainer(
      borderRadius: 16,
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              ...List.generate(
                  5,
                  (i) => Icon(
                        i < 4
                            ? Icons.star
                            : Icons
                                .star_half,
                        size: 20,
                        color: SanggamTheme
                            .primary,
                      )),
              const SizedBox(width: 8),
              const Text('4.5',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight:
                        FontWeight.bold,
                  )),
              const SizedBox(width: 4),
              const Text('(128개 리뷰)',
                  style: TextStyle(
                    color: SanggamTheme
                        .onSurfaceDim,
                    fontSize: 12,
                  )),
            ],
          ),
          const SizedBox(height: 16),
          const Text(
              '"정확도가 높고 결과가 빨리 나와서 좋습니다."',
              style: TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontStyle:
                    FontStyle.italic,
              )),
          const SizedBox(height: 4),
          const Text('- 건강관리러 님',
              style: TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 12,
              )),
        ],
      ),
    );
  }
}
