/// 카트리지 마켓 도메인 모델 및 리포지토리
///
/// 카트리지 상품, 구독 플랜, 주문/배송 관리

/// 카트리지 상품
class CartridgeProduct {
  final String id;
  final String typeCode;
  final String nameKo;
  final String nameEn;
  final String tier; // Basic, Standard, Premium, Professional
  final int price;
  final String unit;
  final String referenceRange;
  final int requiredChannels;
  final int measurementSecs;
  final bool isAvailable;

  const CartridgeProduct({
    required this.id,
    required this.typeCode,
    required this.nameKo,
    required this.nameEn,
    required this.tier,
    required this.price,
    required this.unit,
    required this.referenceRange,
    required this.requiredChannels,
    required this.measurementSecs,
    required this.isAvailable,
  });
}

/// 구독 플랜
class SubscriptionPlan {
  final String id;
  final String name;
  final int monthlyPrice;
  final int discountPercent;
  final List<String> includedCartridgeTypes;
  final int cartridgesPerMonth;

  const SubscriptionPlan({
    required this.id,
    required this.name,
    required this.monthlyPrice,
    required this.discountPercent,
    required this.includedCartridgeTypes,
    required this.cartridgesPerMonth,
  });
}

/// 주문
class Order {
  final String id;
  final List<OrderItem> items;
  final int totalAmount;
  final OrderStatus status;
  final DateTime orderedAt;
  final String? trackingNumber;

  const Order({
    required this.id,
    required this.items,
    required this.totalAmount,
    required this.status,
    required this.orderedAt,
    this.trackingNumber,
  });
}

/// 주문 항목
class OrderItem {
  final String productId;
  final String productName;
  final int quantity;
  final int unitPrice;

  const OrderItem({
    required this.productId,
    required this.productName,
    required this.quantity,
    required this.unitPrice,
  });
}

/// 주문 상태
enum OrderStatus { pending, confirmed, shipping, delivered, cancelled }

/// 마켓 리포지토리 인터페이스
abstract class MarketRepository {
  /// 카트리지 상품 목록
  Future<List<CartridgeProduct>> getProducts({String? tier});

  /// 상품 상세
  Future<CartridgeProduct> getProductDetail(String productId);

  /// 구독 플랜 목록
  Future<List<SubscriptionPlan>> getSubscriptionPlans();

  /// 주문 생성
  Future<Order> createOrder(List<OrderItem> items);

  /// 주문 내역
  Future<List<Order>> getOrders();

  /// 주문 상세
  Future<Order> getOrderDetail(String orderId);

  /// 카트리지 호환성 확인
  Future<bool> checkCompatibility(String typeCode, String deviceId);
}
