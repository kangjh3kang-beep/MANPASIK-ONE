import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/market/data/market_repository_rest.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

void main() {
  group('MarketRepositoryRest', () {
    test('MarketRepositoryRest는 MarketRepository를 구현한다', () {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      expect(repo, isA<MarketRepository>());
    });

    test('getProducts는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final products = await repo.getProducts();
      expect(products, isEmpty);
    });

    test('getProducts tier 필터 적용 시에도 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final products = await repo.getProducts(tier: 'Premium');
      expect(products, isEmpty);
    });

    test('getSubscriptionPlans는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final plans = await repo.getSubscriptionPlans();
      expect(plans, isEmpty);
    });

    test('getOrders는 DioException 시 빈 리스트를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final orders = await repo.getOrders();
      expect(orders, isEmpty);
    });

    test('checkCompatibility는 DioException 시 true를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final compatible = await repo.checkCompatibility('BIO-01', 'dev-1');
      expect(compatible, isTrue);
    });

    test('getGeneralProducts는 DioException 시 시뮬레이션 데이터를 반환한다', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final products = await repo.getGeneralProducts();
      expect(products, isNotEmpty);
      expect(products.length, 8);
    });

    test('getGeneralProducts category 필터 적용', () async {
      final client = ManPaSikRestClient(baseUrl: 'http://localhost:99999/api/v1');
      final repo = MarketRepositoryRest(client, userId: 'user-1');
      final supplements = await repo.getGeneralProducts(category: 'supplement');
      expect(supplements.every((p) => p.category == 'supplement'), isTrue);
    });
  });

  group('도메인 모델 테스트', () {
    test('CartridgeProduct 생성 확인', () {
      const product = CartridgeProduct(
        id: 'p-1', typeCode: 'BIO-GLU', nameKo: '혈당',
        nameEn: 'Glucose', tier: 'Basic', price: 5000,
        unit: 'mg/dL', referenceRange: '70-110',
        requiredChannels: 1, measurementSecs: 60, isAvailable: true,
      );
      expect(product.id, 'p-1');
      expect(product.tier, 'Basic');
      expect(product.isAvailable, isTrue);
    });

    test('SubscriptionPlan 생성 확인', () {
      const plan = SubscriptionPlan(
        id: 'sub-1', name: '프리미엄', monthlyPrice: 39900,
        discountPercent: 20, includedCartridgeTypes: ['BIO-GLU', 'BIO-CHO'],
        cartridgesPerMonth: 30,
      );
      expect(plan.includedCartridgeTypes, hasLength(2));
      expect(plan.discountPercent, 20);
    });

    test('Order 및 OrderItem 생성 확인', () {
      final order = Order(
        id: 'ord-1',
        items: const [
          OrderItem(productId: 'p-1', productName: '혈당', quantity: 10, unitPrice: 5000),
          OrderItem(productId: 'p-2', productName: '콜레스테롤', quantity: 5, unitPrice: 8000),
        ],
        totalAmount: 90000,
        status: OrderStatus.confirmed,
        orderedAt: DateTime(2026, 2, 19),
      );
      expect(order.items, hasLength(2));
      expect(order.status, OrderStatus.confirmed);
      expect(order.trackingNumber, isNull);
    });

    test('GeneralProduct discountPercent 계산', () {
      const product = GeneralProduct(
        id: 'gp-1', name: '오메가-3', price: 29900,
        originalPrice: 39000, rating: 4.5, reviewCount: 100,
        category: 'supplement',
      );
      expect(product.discountPercent, 23); // (39000-29900)/39000*100 = 23.3 → 23
    });

    test('GeneralProduct originalPrice 없으면 discountPercent null', () {
      const product = GeneralProduct(
        id: 'gp-2', name: '유산균', price: 24900,
        rating: 4.3, reviewCount: 50, category: 'supplement',
      );
      expect(product.discountPercent, isNull);
    });

    test('OrderStatus enum 값 확인', () {
      expect(OrderStatus.values, hasLength(5));
      expect(OrderStatus.values, contains(OrderStatus.pending));
      expect(OrderStatus.values, contains(OrderStatus.delivered));
    });
  });
}
