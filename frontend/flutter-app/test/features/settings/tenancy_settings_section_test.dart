import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/network/tenant_interceptor.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/settings/presentation/tenancy_settings_section.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient({this.memberships = const []}) : super(baseUrl: 'http://localhost');
  final List<Map<String, dynamic>> memberships;

  @override
  Future<List<Map<String, dynamic>>> getMyMemberships() async => memberships;
}

GoRouter _buildRouter() => GoRouter(
      initialLocation: '/',
      routes: [
        GoRoute(
          path: '/',
          builder: (_, __) => const Scaffold(body: TenancySettingsSection()),
        ),
        GoRoute(path: '/tenant-switcher', builder: (_, __) => const _Stub('switcher')),
        GoRoute(path: '/invite', builder: (_, __) => const _Stub('accept')),
        GoRoute(path: '/invite-create', builder: (_, __) => const _Stub('create')),
        GoRoute(
          path: '/tenant-members/:tid',
          builder: (_, state) => _Stub('members/${state.pathParameters['tid']}'),
        ),
      ],
    );

class _Stub extends StatelessWidget {
  const _Stub(this.label);
  final String label;
  @override
  Widget build(BuildContext context) =>
      Scaffold(body: Center(child: Text('STUB:$label')));
}

Widget _wrap(_StubClient client, ProviderContainer container, GoRouter router) =>
    UncontrolledProviderScope(
      container: container,
      child: MaterialApp.router(routerConfig: router),
    );

ProviderContainer _container(_StubClient client) => ProviderContainer(
      overrides: [restClientProvider.overrideWithValue(client)],
    );

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  group('TenancySettingsSection', () {
    testWidgets('개인 모드 + 멤버십 없음 → 활성 조직 + 기본 액션만 표시', (tester) async {
      final client = _StubClient();
      final container = _container(client);
      addTearDown(container.dispose);

      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('active-tenant-tile')), findsOneWidget);
      expect(find.text('개인 모드'), findsOneWidget);
      expect(find.byKey(const Key('nav-tenant-switcher')), findsOneWidget);
      expect(find.byKey(const Key('nav-invite-accept')), findsOneWidget);
      // admin 권한 없음 → 발급/멤버 관리 버튼 없음
      expect(find.byKey(const Key('nav-invite-create')), findsNothing);
      expect(find.byKey(const Key('nav-tenant-members')), findsNothing);
    });

    testWidgets('admin 멤버십 있음 → 초대 발급 표시', (tester) async {
      final client = _StubClient(memberships: const [
        {
          'user_id': 'me',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);

      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('nav-invite-create')), findsOneWidget);
    });

    testWidgets('활성 조직이 admin 인 경우 멤버 관리 버튼 표시', (tester) async {
      final client = _StubClient(memberships: const [
        {
          'user_id': 'me',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);

      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'hospA',
            displayName: 'A 병원',
          );

      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('nav-tenant-members')), findsOneWidget);
      expect(find.text('A 병원 멤버 관리'), findsOneWidget);
    });

    testWidgets('활성 조직이 admin 아닌 경우 멤버 관리 미표시', (tester) async {
      final client = _StubClient(memberships: const [
        {
          'user_id': 'me',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
        {
          'user_id': 'me',
          'tenant_id': 'famB',
          'role': 'member',
          'active': true,
          'joined_at': 2,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);

      // famB 활성 (member 역할만)
      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'famB',
            displayName: 'B 가족',
          );

      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      // admin 권한이 다른 조직 (hospA) 에 있어 nav-invite-create 는 표시
      expect(find.byKey(const Key('nav-invite-create')), findsOneWidget);
      // 활성 조직 famB 에서는 admin 아니므로 멤버 관리 미표시
      expect(find.byKey(const Key('nav-tenant-members')), findsNothing);
    });

    testWidgets('조직 전환 버튼 탭 → /tenant-switcher 라우팅', (tester) async {
      final client = _StubClient();
      final container = _container(client);
      addTearDown(container.dispose);
      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('nav-tenant-switcher')));
      await tester.pumpAndSettle();

      expect(find.text('STUB:switcher'), findsOneWidget);
    });

    testWidgets('초대 수락 버튼 탭 → /invite', (tester) async {
      final client = _StubClient();
      final container = _container(client);
      addTearDown(container.dispose);
      await tester.pumpWidget(_wrap(client, container, _buildRouter()));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('nav-invite-accept')));
      await tester.pumpAndSettle();

      expect(find.text('STUB:accept'), findsOneWidget);
    });
  });
}
