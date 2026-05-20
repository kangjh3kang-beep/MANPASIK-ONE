import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/settings/presentation/tenant_members_screen.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient({this.members = const [], this.failList = false}) : super(baseUrl: 'http://localhost');

  final List<Map<String, dynamic>> members;
  final bool failList;

  String? lastUpdatedRole;
  String? lastUpdatedUserId;
  String? lastRemovedUserId;

  @override
  Future<List<Map<String, dynamic>>> listTenantMembers(String tenantId) async {
    if (failList) throw StateError('네트워크 에러');
    return members;
  }

  @override
  Future<Map<String, dynamic>> updateTenantMemberRole({
    required String tenantId,
    required String userId,
    required String newRole,
  }) async {
    lastUpdatedUserId = userId;
    lastUpdatedRole = newRole;
    return {'user_id': userId, 'tenant_id': tenantId, 'role': newRole, 'active': true, 'joined_at': 1};
  }

  @override
  Future<void> removeTenantMember({
    required String tenantId,
    required String userId,
  }) async {
    lastRemovedUserId = userId;
  }
}

Widget _wrap(Widget child, _StubClient client) => UncontrolledProviderScope(
      container: ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(client),
      ]),
      child: MaterialApp(home: child),
    );

void main() {
  group('TenantMembersScreen', () {
    testWidgets('빈 목록 시 안내 메시지', (tester) async {
      final client = _StubClient(members: const []);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();
      expect(find.text('멤버가 없습니다'), findsOneWidget);
    });

    testWidgets('멤버 목록 표시', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'u1',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
        {
          'user_id': 'u2',
          'tenant_id': 'hospA',
          'role': 'member',
          'active': false,
          'joined_at': 2,
        },
      ]);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();
      expect(find.text('u1'), findsOneWidget);
      expect(find.text('u2'), findsOneWidget);
      expect(find.text('활성'), findsOneWidget);
      expect(find.text('비활성'), findsOneWidget);
    });

    testWidgets('역할 변경 호출', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'u1',
          'tenant_id': 'hospA',
          'role': 'member',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('role-u1')));
      await tester.pumpAndSettle();
      // 'medical_staff' 옵션 탭
      await tester.tap(find.text('medical_staff').last);
      await tester.pumpAndSettle();

      expect(client.lastUpdatedUserId, 'u1');
      expect(client.lastUpdatedRole, 'medical_staff');
    });

    testWidgets('역할 변경 동일 값 시 호출 안 함', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'u1',
          'tenant_id': 'hospA',
          'role': 'member',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('role-u1')));
      await tester.pumpAndSettle();
      // 동일한 'member' 선택
      await tester.tap(find.text('member').last);
      await tester.pumpAndSettle();

      expect(client.lastUpdatedUserId, isNull);
    });

    testWidgets('멤버 제거 - 확인 다이얼로그 → 취소', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'u1',
          'tenant_id': 'hospA',
          'role': 'member',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('remove-u1')));
      await tester.pumpAndSettle();
      expect(find.text('멤버 제거'), findsOneWidget);

      await tester.tap(find.text('취소'));
      await tester.pumpAndSettle();

      expect(client.lastRemovedUserId, isNull);
    });

    testWidgets('멤버 제거 - 확인 다이얼로그 → 제거', (tester) async {
      final client = _StubClient(members: const [
        {
          'user_id': 'u1',
          'tenant_id': 'hospA',
          'role': 'member',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('remove-u1')));
      await tester.pumpAndSettle();
      await tester.tap(find.text('제거').last);
      await tester.pumpAndSettle();

      expect(client.lastRemovedUserId, 'u1');
    });

    testWidgets('에러 상태 + 재시도 버튼', (tester) async {
      final client = _StubClient(failList: true);
      await tester
          .pumpWidget(_wrap(const TenantMembersScreen(tenantId: 'hospA'), client));
      await tester.pumpAndSettle();
      expect(find.text('멤버 조회 실패'), findsOneWidget);
      expect(find.text('다시 시도'), findsOneWidget);
    });
  });
}
