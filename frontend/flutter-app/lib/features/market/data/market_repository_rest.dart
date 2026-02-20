import 'package:dio/dio.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 MarketRepository 구현체
class MarketRepositoryRest implements MarketRepository {
  MarketRepositoryRest(this._client, {required this.userId});

  final ManPaSikRestClient _client;
  final String userId;

  @override
  Future<List<CartridgeProduct>> getProducts({String? tier}) async {
    try {
      final res = await _client.listProducts();
      final products = res['products'] as List<dynamic>? ?? [];
      var list = products
          .map((p) => _mapProduct(p as Map<String, dynamic>))
          .toList();
      if (tier != null) {
        list = list.where((p) => p.tier == tier).toList();
      }
      return list;
    } on DioException {
      return [];
    }
  }

  @override
  Future<CartridgeProduct> getProductDetail(String productId) async {
    final res = await _client.getProduct(productId);
    return _mapProduct(res);
  }

  @override
  Future<List<SubscriptionPlan>> getSubscriptionPlans() async {
    try {
      final res = await _client.listSubscriptionPlans();
      final plans = res['plans'] as List<dynamic>? ?? [];
      return plans
          .map((p) => _mapSubscriptionPlan(p as Map<String, dynamic>))
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<Order> createOrder(List<OrderItem> items) async {
    final res = await _client.createOrder(userId: userId);
    return _mapOrder(res);
  }

  @override
  Future<List<Order>> getOrders() async {
    try {
      final res = await _client.listOrders(userId);
      final orders = res['orders'] as List<dynamic>? ?? [];
      return orders
          .map((o) => _mapOrder(o as Map<String, dynamic>))
          .toList();
    } on DioException {
      return [];
    }
  }

  @override
  Future<Order> getOrderDetail(String orderId) async {
    // Use list and filter since no single-order endpoint
    final orders = await getOrders();
    return orders.firstWhere(
      (o) => o.id == orderId,
      orElse: () => throw Exception('Order not found'),
    );
  }

  @override
  Future<bool> checkCompatibility(String typeCode, String deviceId) async {
    try {
      final res = await _client.validateCartridge(cartridgeUid: typeCode, userId: userId);
      return res['valid'] as bool? ?? true;
    } on DioException {
      return true; // assume compatible on error
    }
  }

  @override
  Future<List<GeneralProduct>> getGeneralProducts({String? category}) async {
    try {
      final res = await _client.listProducts();
      final products = res['general_products'] as List<dynamic>? ?? [];
      var list = products
          .map((p) => _mapGeneralProduct(p as Map<String, dynamic>))
          .toList();
      if (category != null) {
        list = list.where((p) => p.category == category).toList();
      }
      return list;
    } on DioException {
      // 서버에 일반상품 미구현 시 시뮬레이션 데이터
      return _simulatedGeneralProducts(category);
    }
  }

  List<GeneralProduct> _simulatedGeneralProducts(String? category) {
    final all = [
      const GeneralProduct(id: 'gp-1', name: '오메가-3 피쉬오일 (60캡슐)', price: 29900, originalPrice: 39000, rating: 4.5, reviewCount: 342, category: 'supplement', freeShipping: true),
      const GeneralProduct(id: 'gp-2', name: '프로바이오틱스 유산균', price: 24900, rating: 4.3, reviewCount: 218, category: 'supplement', freeShipping: true),
      const GeneralProduct(id: 'gp-3', name: '만파식 리더기 보호 케이스', price: 15000, rating: 4.8, reviewCount: 89, category: 'accessory'),
      const GeneralProduct(id: 'gp-4', name: '휴대용 카트리지 보관함', price: 22000, originalPrice: 28000, rating: 4.6, reviewCount: 156, category: 'accessory', freeShipping: true),
      const GeneralProduct(id: 'gp-5', name: '비타민 D3 + K2 복합', price: 19900, rating: 4.4, reviewCount: 275, category: 'supplement'),
      const GeneralProduct(id: 'gp-6', name: '아로마 테라피 디퓨저', price: 45000, originalPrice: 55000, rating: 4.7, reviewCount: 62, category: 'wellness', freeShipping: true),
      const GeneralProduct(id: 'gp-7', name: '스트레스 릴리프 티 세트', price: 18500, rating: 4.2, reviewCount: 94, category: 'wellness'),
      const GeneralProduct(id: 'gp-8', name: '만파식 프리미엄 선물세트', price: 89000, originalPrice: 110000, rating: 4.9, reviewCount: 45, category: 'giftset', freeShipping: true),
    ];
    if (category != null) {
      return all.where((p) => p.category == category).toList();
    }
    return all;
  }

  GeneralProduct _mapGeneralProduct(Map<String, dynamic> m) {
    return GeneralProduct(
      id: m['id'] as String? ?? '',
      name: m['name'] as String? ?? '',
      price: m['price'] as int? ?? 0,
      originalPrice: m['original_price'] as int?,
      rating: (m['rating'] as num?)?.toDouble() ?? 0.0,
      reviewCount: m['review_count'] as int? ?? 0,
      category: m['category'] as String? ?? 'supplement',
      imageUrl: m['image_url'] as String?,
      freeShipping: m['free_shipping'] as bool? ?? false,
      isWishlisted: m['is_wishlisted'] as bool? ?? false,
    );
  }

  CartridgeProduct _mapProduct(Map<String, dynamic> m) {
    return CartridgeProduct(
      id: m['id'] as String? ?? m['product_id'] as String? ?? '',
      typeCode: m['type_code'] as String? ?? '',
      nameKo: m['name_ko'] as String? ?? m['name'] as String? ?? '',
      nameEn: m['name_en'] as String? ?? '',
      tier: m['tier'] as String? ?? 'Basic',
      price: m['price'] as int? ?? (m['price_krw'] as int? ?? 0),
      unit: m['unit'] as String? ?? '',
      referenceRange: m['reference_range'] as String? ?? '',
      requiredChannels: m['required_channels'] as int? ?? 1,
      measurementSecs: m['measurement_secs'] as int? ?? 60,
      isAvailable: m['is_available'] as bool? ?? true,
    );
  }

  SubscriptionPlan _mapSubscriptionPlan(Map<String, dynamic> m) {
    return SubscriptionPlan(
      id: m['id'] as String? ?? m['plan_id'] as String? ?? '',
      name: m['name'] as String? ?? '',
      monthlyPrice: m['monthly_price'] as int? ?? 0,
      discountPercent: m['discount_percent'] as int? ?? 0,
      includedCartridgeTypes:
          (m['included_cartridge_types'] as List<dynamic>?)
              ?.map((t) => t.toString())
              .toList() ??
              [],
      cartridgesPerMonth: m['cartridges_per_month'] as int? ?? 0,
    );
  }

  Order _mapOrder(Map<String, dynamic> m) {
    return Order(
      id: m['id'] as String? ?? m['order_id'] as String? ?? '',
      items: _parseOrderItems(m['items']),
      totalAmount: m['total_amount'] as int? ?? 0,
      status: _parseOrderStatus(m['status']),
      orderedAt: m['ordered_at'] != null
          ? DateTime.tryParse(m['ordered_at'] as String) ?? DateTime.now()
          : m['created_at'] != null
              ? DateTime.tryParse(m['created_at'] as String) ?? DateTime.now()
              : DateTime.now(),
      trackingNumber: m['tracking_number'] as String?,
    );
  }

  List<OrderItem> _parseOrderItems(dynamic items) {
    if (items is! List) return [];
    return items.map((i) {
      final m = i as Map<String, dynamic>;
      return OrderItem(
        productId: m['product_id'] as String? ?? '',
        productName: m['product_name'] as String? ?? '',
        quantity: m['quantity'] as int? ?? 1,
        unitPrice: m['unit_price'] as int? ?? 0,
      );
    }).toList();
  }

  OrderStatus _parseOrderStatus(dynamic v) {
    if (v is String) {
      switch (v.toLowerCase()) {
        case 'confirmed':
          return OrderStatus.confirmed;
        case 'shipping':
          return OrderStatus.shipping;
        case 'delivered':
          return OrderStatus.delivered;
        case 'cancelled':
          return OrderStatus.cancelled;
      }
    }
    return OrderStatus.pending;
  }
}
