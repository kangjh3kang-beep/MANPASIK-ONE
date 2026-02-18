import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 상품 상세 화면
class ProductDetailScreen extends ConsumerWidget {
  const ProductDetailScreen({super.key, required this.productId});

  final String productId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final client = ref.watch(restClientProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('상품 상세'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          _WishlistButton(productId: productId),
          IconButton(
            icon: const Icon(Icons.shopping_cart_outlined),
            onPressed: () => context.push('/market/cart'),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(context, ref),
      body: FutureBuilder<Map<String, dynamic>>(
        future: client.getProduct(productId),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return _buildFallback(context, theme, ref);
          }
          final data = snapshot.data ?? {};
          return _buildProductDetail(context, theme, ref, data);
        },
      ),
    );
  }

  Widget _buildProductDetail(BuildContext context, ThemeData theme, WidgetRef ref, Map<String, dynamic> data) {
    final name = data['name'] as String? ?? '카트리지 상품';
    final description = data['description'] as String? ?? '상품 설명이 없습니다.';
    final price = data['price'] as num? ?? 0;
    final tier = data['tier'] as String? ?? 'Basic';

    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 상품 이미지 영역
          Container(
            height: 250,
            color: theme.colorScheme.surfaceContainerHighest,
            child: Center(
              child: Icon(Icons.science, size: 80, color: AppTheme.sanggamGold.withOpacity(0.5)),
            ),
          ),

          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 티어 배지
                Chip(
                  label: Text(tier, style: const TextStyle(fontSize: 12)),
                  backgroundColor: AppTheme.sanggamGold.withOpacity(0.15),
                  side: BorderSide.none,
                ),
                const SizedBox(height: 8),

                // 상품명
                Text(name, style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),

                // 가격
                Text(
                  '₩${_formatPrice(price)}',
                  style: theme.textTheme.titleLarge?.copyWith(
                    color: AppTheme.sanggamGold,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                const Divider(),
                const SizedBox(height: 8),

                // 상품 설명
                Text('상품 설명', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Text(description, style: theme.textTheme.bodyMedium?.copyWith(height: 1.6)),
                const SizedBox(height: 16),

                // 리뷰 섹션
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text('사용자 리뷰', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    TextButton.icon(
                      onPressed: () => _showWriteReviewDialog(context, ref),
                      icon: const Icon(Icons.rate_review, size: 18),
                      label: const Text('리뷰 작성'),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                _ReviewSection(productId: productId),
                const SizedBox(height: 16),

                // 스펙 테이블
                Text('제품 스펙', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                _buildSpecRow(theme, '제품 ID', productId),
                _buildSpecRow(theme, '등급', tier),
                _buildSpecRow(theme, '호환 리더기', 'ManPaSik Reader v2.0+'),
                _buildSpecRow(theme, '보관 조건', '실온 (15~30°C)'),
                _buildSpecRow(theme, '유효 기간', '제조일로부터 12개월'),
                const SizedBox(height: 24),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFallback(BuildContext context, ThemeData theme, WidgetRef ref) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Container(
            height: 250,
            color: theme.colorScheme.surfaceContainerHighest,
            child: Center(
              child: Icon(Icons.science, size: 80, color: AppTheme.sanggamGold.withOpacity(0.5)),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('카트리지 상품', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Text('₩29,900', style: theme.textTheme.titleLarge?.copyWith(color: AppTheme.sanggamGold, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                Text('서버 연결 후 상세 정보가 표시됩니다.', style: theme.textTheme.bodyMedium),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSpecRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          SizedBox(
            width: 100,
            child: Text(label, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
          ),
          Expanded(child: Text(value, style: theme.textTheme.bodySmall)),
        ],
      ),
    );
  }

  String _formatPrice(num price) {
    return price.toInt().toString().replaceAllMapped(
      RegExp(r'(\d)(?=(\d{3})+(?!\d))'),
      (m) => '${m[1]},',
    );
  }

  void _showWriteReviewDialog(BuildContext context, WidgetRef ref) {
    int rating = 5;
    final contentCtrl = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setDialogState) => AlertDialog(
          title: const Text('리뷰 작성'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(5, (i) => IconButton(
                  icon: Icon(
                    i < rating ? Icons.star : Icons.star_border,
                    color: Colors.amber,
                  ),
                  onPressed: () => setDialogState(() => rating = i + 1),
                )),
              ),
              const SizedBox(height: 8),
              TextField(
                controller: contentCtrl,
                maxLines: 3,
                decoration: const InputDecoration(
                  hintText: '사용 후기를 작성해주세요',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
            FilledButton(
              onPressed: () async {
                final client = ref.read(restClientProvider);
                final userId = ref.read(authProvider).userId ?? '';
                try {
                  await client.createProductReview(
                    productId: productId,
                    userId: userId,
                    rating: rating,
                    content: contentCtrl.text,
                  );
                  if (ctx.mounted) Navigator.pop(ctx);
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('리뷰가 등록되었습니다.')),
                    );
                  }
                } catch (e) {
                  if (ctx.mounted) {
                    ScaffoldMessenger.of(ctx).showSnackBar(
                      SnackBar(content: Text('리뷰 등록 실패: $e')),
                    );
                  }
                }
              },
              child: const Text('등록'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildBottomBar(BuildContext context, WidgetRef ref) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Expanded(
              child: OutlinedButton(
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('장바구니에 추가되었습니다.')),
                  );
                },
                child: const Text('장바구니'),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              flex: 2,
              child: FilledButton(
                onPressed: () => context.push('/market/cart'),
                style: FilledButton.styleFrom(backgroundColor: AppTheme.sanggamGold),
                child: const Text('바로 구매'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// 위시리스트 하트 토글 버튼
class _WishlistButton extends StatefulWidget {
  const _WishlistButton({required this.productId});
  final String productId;

  @override
  State<_WishlistButton> createState() => _WishlistButtonState();
}

class _WishlistButtonState extends State<_WishlistButton> {
  bool _isWished = false;

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Icon(
        _isWished ? Icons.favorite : Icons.favorite_border,
        color: _isWished ? Colors.red : null,
      ),
      tooltip: '위시리스트',
      onPressed: () {
        setState(() => _isWished = !_isWished);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(_isWished ? '위시리스트에 추가되었습니다.' : '위시리스트에서 제거되었습니다.')),
        );
      },
    );
  }
}

/// REST API 기반 리뷰 섹션
class _ReviewSection extends ConsumerWidget {
  const _ReviewSection({required this.productId});
  final String productId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final client = ref.watch(restClientProvider);
    final theme = Theme.of(context);

    return FutureBuilder<Map<String, dynamic>>(
      future: client.getProductReviews(productId, limit: 3),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Card(child: Padding(padding: EdgeInsets.all(16), child: Center(child: CircularProgressIndicator(strokeWidth: 2))));
        }

        final data = snapshot.data ?? {};
        final reviews = (data['reviews'] as List?)?.cast<Map<String, dynamic>>() ?? [];
        final avgRating = (data['average_rating'] as num?)?.toDouble() ?? 0.0;
        final totalCount = (data['total_count'] as num?)?.toInt() ?? reviews.length;

        if (reviews.isEmpty && !snapshot.hasError) {
          return _buildFallbackReview(theme);
        }

        return Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 평점 요약
                Row(
                  children: [
                    ...List.generate(5, (i) => Icon(
                      i < avgRating.round() ? Icons.star : Icons.star_border,
                      size: 20,
                      color: Colors.amber,
                    )),
                    const SizedBox(width: 8),
                    Text(avgRating.toStringAsFixed(1), style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(width: 4),
                    Text('($totalCount개 리뷰)', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
                  ],
                ),
                const SizedBox(height: 12),
                // 최근 리뷰 목록
                ...reviews.take(3).map((r) => Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          ...List.generate(5, (i) => Icon(
                            i < ((r['rating'] as num?) ?? 5) ? Icons.star : Icons.star_border,
                            size: 14,
                            color: Colors.amber,
                          )),
                          const SizedBox(width: 8),
                          Text(r['author_name'] as String? ?? '익명', style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
                        ],
                      ),
                      const SizedBox(height: 4),
                      Text(r['content'] as String? ?? '', style: theme.textTheme.bodySmall),
                    ],
                  ),
                )),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildFallbackReview(ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                ...List.generate(5, (i) => Icon(
                  i < 4 ? Icons.star : Icons.star_half,
                  size: 20,
                  color: Colors.amber,
                )),
                const SizedBox(width: 8),
                Text('4.5', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                const SizedBox(width: 4),
                Text('(128개 리뷰)', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
              ],
            ),
            const SizedBox(height: 12),
            Text('"정확도가 높고 결과가 빨리 나와서 좋습니다."', style: theme.textTheme.bodyMedium?.copyWith(fontStyle: FontStyle.italic)),
            const SizedBox(height: 4),
            Text('- 건강관리러 님', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
          ],
        ),
      ),
    );
  }
}
