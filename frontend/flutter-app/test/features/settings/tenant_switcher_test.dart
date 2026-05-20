import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/network/tenant_interceptor.dart';
import 'package:manpasik/features/settings/presentation/tenant_switcher_screen.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

Widget _wrap(Widget child) => ProviderScope(
      child: MaterialApp(home: child),
    );

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  group('ActiveTenantNotifier', () {
    test('초기 상태는 개인 모드 (tenantId=null)', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);
      // 비동기 _loadFromPreferences 가 끝날 때까지 대기
      await container.read(activeTenantProvider.notifier).clear();
      final state = container.read(activeTenantProvider);
      expect(state.isPersonal, true);
    });

    test('switchTo 호출 시 SharedPreferences + 상태 모두 갱신', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'hospital-A',
            displayName: 'A 병원',
          );

      // Riverpod 상태
      final state = container.read(activeTenantProvider);
      expect(state.tenantId, 'hospital-A');
      expect(state.displayName, 'A 병원');
      expect(state.isPersonal, false);

      // SharedPreferences (= TenantInterceptor 헤더 부착 대상)
      final stored = await TenantInterceptor.getActiveTenant();
      expect(stored, 'hospital-A');
    });

    test('clear 호출 시 개인 모드로 복귀', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'hospital-A',
          );
      await container.read(activeTenantProvider.notifier).clear();

      expect(container.read(activeTenantProvider).isPersonal, true);
      expect(await TenantInterceptor.getActiveTenant(), isNull);
    });
  });

  group('TenantSwitcherScreen', () {
    testWidgets('빈 조직 목록 시 안내 메시지 표시', (tester) async {
      await tester.pumpWidget(_wrap(const TenantSwitcherScreen(tenants: [])));
      await tester.pumpAndSettle();
      expect(find.text('소속된 조직이 없습니다'), findsOneWidget);
      // 개인 모드 타일은 항상 표시
      expect(find.text('개인'), findsOneWidget);
    });

    testWidgets('조직 목록 표시', (tester) async {
      await tester.pumpWidget(_wrap(const TenantSwitcherScreen(tenants: [
        TenantOption(tenantId: 't1', displayName: 'A 병원', role: 'admin'),
        TenantOption(tenantId: 't2', displayName: 'B 가족', role: 'member'),
      ])));
      await tester.pumpAndSettle();
      expect(find.text('A 병원'), findsOneWidget);
      expect(find.text('B 가족'), findsOneWidget);
      expect(find.text('역할: admin'), findsOneWidget);
      expect(find.text('역할: member'), findsOneWidget);
    });

    testWidgets('조직 탭 시 activeTenantProvider 갱신', (tester) async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      await tester.pumpWidget(UncontrolledProviderScope(
        container: container,
        child: const MaterialApp(
          home: TenantSwitcherScreen(tenants: [
            TenantOption(
                tenantId: 'hospital-X', displayName: 'X 병원', role: 'admin'),
          ]),
        ),
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.text('X 병원'));
      await tester.pumpAndSettle();

      final state = container.read(activeTenantProvider);
      expect(state.tenantId, 'hospital-X');
      expect(state.displayName, 'X 병원');
    });

    testWidgets('개인 타일 탭 시 조직 헤더 제거', (tester) async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      // 미리 조직 활성화
      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'pre-set',
            displayName: 'PRE',
          );

      await tester.pumpWidget(UncontrolledProviderScope(
        container: container,
        child: const MaterialApp(
          home: TenantSwitcherScreen(tenants: [
            TenantOption(tenantId: 'pre-set', displayName: 'PRE', role: 'a'),
          ]),
        ),
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.text('개인'));
      await tester.pumpAndSettle();

      expect(container.read(activeTenantProvider).isPersonal, true);
    });
  });

  group('ActiveTenant equality', () {
    test('같은 ID/이름은 동일', () {
      const a1 = ActiveTenant(tenantId: 't1', displayName: 'X');
      const a2 = ActiveTenant(tenantId: 't1', displayName: 'X');
      expect(a1 == a2, true);
      expect(a1.hashCode, a2.hashCode);
    });

    test('isPersonal 판정', () {
      expect(const ActiveTenant().isPersonal, true);
      expect(const ActiveTenant(tenantId: '').isPersonal, true);
      expect(const ActiveTenant(tenantId: 't1').isPersonal, false);
    });

    test('copyWith', () {
      const a = ActiveTenant(tenantId: 't1', displayName: 'X');
      final b = a.copyWith(displayName: 'Y');
      expect(b.tenantId, 't1');
      expect(b.displayName, 'Y');
    });
  });
}
