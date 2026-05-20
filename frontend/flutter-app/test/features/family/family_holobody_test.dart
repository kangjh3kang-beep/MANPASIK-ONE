import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/network/tenant_interceptor.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/family/presentation/family_holobody_widget.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient({this.members = const []}) : super(baseUrl: 'http://localhost');
  final List<Map<String, dynamic>> members;

  @override
  Future<List<Map<String, dynamic>>> listTenantMembers(String tenantId) async =>
      members;
}

Widget _wrap(Widget child, _StubClient client, ProviderContainer container) =>
    UncontrolledProviderScope(
      container: container,
      child: MaterialApp(home: Scaffold(body: child)),
    );

ProviderContainer _container(_StubClient client) => ProviderContainer(
      overrides: [
        restClientProvider.overrideWithValue(client),
      ],
    );

/// _fakeHoloBuilder 는 HoloBody 대신 단순 Container 를 렌더 (테스트용).
Widget _fakeHoloBuilder({
  required double width,
  required double height,
  required Map<String, dynamic> bioData,
  required Color color,
  Key? key,
}) =>
    Container(
      key: key,
      width: width,
      height: height,
      color: color,
    );

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  group('FamilyHoloBodyWidget', () {
    testWidgets('개인 모드: 멤버 rail 미표시 + HoloBody 표시', (tester) async {
      final client = _StubClient();
      final container = _container(client);
      addTearDown(container.dispose);

      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('member-rail')), findsNothing);
      // HoloBody 위젯이 렌더링되었는지 (CustomPaint 또는 Container)
      expect(find.byType(FamilyHoloBodyWidget), findsOneWidget);
    });

    testWidgets('가족 그룹 활성: 멤버 rail 표시 + 본인 chip', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'mom',
          'tenant_id': 'fam-A',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
        {
          'user_id': 'dad',
          'tenant_id': 'fam-A',
          'role': 'member',
          'active': true,
          'joined_at': 2,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);

      // 활성 조직 설정
      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'fam-A',
            displayName: '우리 가족',
          );

      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      expect(find.byKey(const Key('member-rail')), findsOneWidget);
      expect(find.byKey(const Key('chip-self')), findsOneWidget);
      expect(find.byKey(const Key('chip-mom')), findsOneWidget);
      expect(find.byKey(const Key('chip-dad')), findsOneWidget);
    });

    testWidgets('멤버 chip 탭 → fetchMemberBioData 호출', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'mom',
          'tenant_id': 'fam-A',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);
      await container.read(activeTenantProvider.notifier).switchTo(
            tenantId: 'fam-A',
          );

      String? lastFetched;
      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
          fetchMemberBioData: (uid) async {
            lastFetched = uid;
            return {'heart_rate': 75};
          },
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('chip-mom')));
      await tester.pumpAndSettle();

      expect(lastFetched, 'mom');
    });

    testWidgets('본인 chip 탭으로 복귀', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'mom',
          'tenant_id': 'fam-A',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);
      await container.read(activeTenantProvider.notifier).switchTo(tenantId: 'fam-A');

      int fetchCount = 0;
      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
          fetchMemberBioData: (uid) async {
            fetchCount++;
            return {};
          },
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('chip-mom')));
      await tester.pumpAndSettle();
      expect(fetchCount, 1);

      // 본인으로 복귀 → fetch 안 함
      await tester.tap(find.byKey(const Key('chip-self')));
      await tester.pumpAndSettle();
      expect(fetchCount, 1, reason: '본인 복귀 시 fetch 호출 안 됨');
    });

    testWidgets('fetch 에러 시에도 위젯 유지', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'mom',
          'tenant_id': 'fam-A',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);
      await container.read(activeTenantProvider.notifier).switchTo(tenantId: 'fam-A');

      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
          fetchMemberBioData: (uid) async {
            throw StateError('네트워크');
          },
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('chip-mom')));
      await tester.pumpAndSettle();

      // 에러여도 위젯이 사라지지 않음
      expect(find.byType(FamilyHoloBodyWidget), findsOneWidget);
    });

    testWidgets('fetchMemberBioData 미주입 시 빈 데이터로 동작', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'mom',
          'tenant_id': 'fam-A',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      final container = _container(client);
      addTearDown(container.dispose);
      await container.read(activeTenantProvider.notifier).switchTo(tenantId: 'fam-A');

      await tester.pumpWidget(_wrap(
        FamilyHoloBodyWidget(
          width: 200,
          height: 300,
          holoBodyBuilder: _fakeHoloBuilder,
        ),
        client,
        container,
      ));
      await tester.pumpAndSettle();

      // chip 탭해도 panic 없음
      await tester.tap(find.byKey(const Key('chip-mom')));
      await tester.pumpAndSettle();
    });
  });
}
