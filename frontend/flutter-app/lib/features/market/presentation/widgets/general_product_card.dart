import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

/// 일반 상품 카드 (건강식품/웰빙/액세서리)
class GeneralProductCard extends StatelessWidget {
  final GeneralProduct product;

  const GeneralProductCard({super.key, required this.product});

  @override
  Widget build(BuildContext context) {
    return ScaleButton(
      onPressed: () => context.push('/market/product/${product.id}'),
      child: PorcelainContainer(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 이미지 플레이스홀더 + 찜 버튼
            Expanded(
              flex: 3,
              child: Stack(
                children: [
                  Container(
                    width: double.infinity,
                    decoration: BoxDecoration(
                      color: _categoryColor(product.category).withOpacity(0.08),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Icon(
                        _categoryIcon(product.category),
                        size: 40,
                        color: _categoryColor(product.category).withOpacity(0.5),
                      ),
                    ),
                  ),
                  // 할인 뱃지
                  if (product.discountPercent != null)
                    Positioned(
                      top: 4,
                      left: 4,
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                        decoration: BoxDecoration(
                          color: Colors.red,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text(
                          '${product.discountPercent}%',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 10,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  // 찜 버튼
                  Positioned(
                    top: 4,
                    right: 4,
                    child: Icon(
                      product.isWishlisted ? Icons.favorite : Icons.favorite_border,
                      size: 20,
                      color: product.isWishlisted ? Colors.red : Colors.grey,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),
            // 상품명
            Expanded(
              flex: 2,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    product.name,
                    style: const TextStyle(
                      fontSize: 13,
                      fontWeight: FontWeight.w600,
                      color: Color(0xFF1A1A1A),
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  // 별점
                  Row(
                    children: [
                      const Icon(Icons.star, size: 12, color: Color(0xFFD4AF37)),
                      const SizedBox(width: 2),
                      Text(
                        '${product.rating}',
                        style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: Color(0xFF666666)),
                      ),
                      Text(
                        ' (${product.reviewCount})',
                        style: const TextStyle(fontSize: 11, color: Color(0xFF999999)),
                      ),
                    ],
                  ),
                  const Spacer(),
                  // 가격
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      if (product.originalPrice != null) ...[
                        Text(
                          '₩${_formatPrice(product.originalPrice!)}',
                          style: const TextStyle(
                            fontSize: 11,
                            color: Color(0xFFAAAAAA),
                            decoration: TextDecoration.lineThrough,
                          ),
                        ),
                        const SizedBox(width: 4),
                      ],
                      Text(
                        '₩${_formatPrice(product.price)}',
                        style: const TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w700,
                          color: Color(0xFFD4AF37),
                        ),
                      ),
                    ],
                  ),
                  if (product.freeShipping)
                    const Padding(
                      padding: EdgeInsets.only(top: 2),
                      child: Text(
                        '무료 배송',
                        style: TextStyle(fontSize: 10, color: Color(0xFF4CAF50)),
                      ),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Color _categoryColor(String category) {
    switch (category) {
      case 'supplement':
        return Colors.green;
      case 'wellness':
        return Colors.purple;
      case 'accessory':
        return Colors.blue;
      case 'giftset':
        return const Color(0xFFD4AF37);
      default:
        return Colors.grey;
    }
  }

  IconData _categoryIcon(String category) {
    switch (category) {
      case 'supplement':
        return Icons.medication_rounded;
      case 'wellness':
        return Icons.spa_rounded;
      case 'accessory':
        return Icons.devices_rounded;
      case 'giftset':
        return Icons.card_giftcard_rounded;
      default:
        return Icons.shopping_bag_rounded;
    }
  }

  String _formatPrice(int price) {
    return price.toString().replaceAllMapped(
      RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
      (match) => '${match[1]},',
    );
  }
}
