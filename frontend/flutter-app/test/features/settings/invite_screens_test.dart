import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/settings/presentation/invite_accept_screen.dart';
import 'package:manpasik/features/settings/presentation/invite_create_screen.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient({
    this.memberships = const [],
    this.createResponse = const {'token': 'tok-stub', 'status': 'pending'},
    this.acceptResponse = const {'tenant_id': 'hospA', 'role': 'member'},
    this.failCreate = false,
    this.failAccept = false,
  }) : super(baseUrl: 'http://localhost');

  final List<Map<String, dynamic>> memberships;
  final Map<String, dynamic> createResponse;
  final Map<String, dynamic> acceptResponse;
  final bool failCreate;
  final bool failAccept;

  String? lastCreatedTenantId;
  String? lastAcceptedToken;

  @override
  Future<List<Map<String, dynamic>>> getMyMemberships() async => memberships;

  @override
  Future<Map<String, dynamic>> createTenantInvitation({
    required String tenantId,
    required String role,
    String? inviteeHint,
    int? ttlHours,
  }) async {
    if (failCreate) throw StateError('생성 실패');
    lastCreatedTenantId = tenantId;
    return createResponse;
  }

  @override
  Future<Map<String, dynamic>> acceptTenantInvitation(String token) async {
    if (failAccept) throw StateError('수락 실패');
    lastAcceptedToken = token;
    return acceptResponse;
  }
}

Widget _wrap(Widget child, _StubClient client) => UncontrolledProviderScope(
      container: ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(client),
      ]),
      child: MaterialApp(home: child),
    );

void main() {
  group('InviteCreateScreen', () {
    testWidgets('admin 멤버십 없으면 잠금 메시지', (tester) async {
      final client = _StubClient(memberships: [
        {
          'user_id': 'me',
          'tenant_id': 'famA',
          'role': 'member',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester.pumpWidget(_wrap(const InviteCreateScreen(), client));
      await tester.pumpAndSettle();
      expect(find.text('관리자 권한이 있는 조직이 없습니다'), findsOneWidget);
    });

    testWidgets('admin 멤버십이 있으면 발급 폼 표시', (tester) async {
      final client = _StubClient(memberships: [
        {
          'user_id': 'me',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester.pumpWidget(_wrap(const InviteCreateScreen(), client));
      await tester.pumpAndSettle();
      expect(find.byKey(const Key('tenant-dropdown')), findsOneWidget);
      expect(find.byKey(const Key('role-dropdown')), findsOneWidget);
      expect(find.byKey(const Key('submit-invite')), findsOneWidget);
    });

    testWidgets('조직 미선택 시 에러', (tester) async {
      final client = _StubClient(memberships: [
        {
          'user_id': 'me',
          'tenant_id': 'hospA',
          'role': 'admin',
          'active': true,
          'joined_at': 1,
        },
      ]);
      await tester.pumpWidget(_wrap(const InviteCreateScreen(), client));
      await tester.pumpAndSettle();
      // 드롭다운 미선택 + 발급 버튼 탭
      await tester.tap(find.byKey(const Key('submit-invite')));
      await tester.pumpAndSettle();
      expect(find.text('조직을 선택하세요'), findsOneWidget);
    });

    testWidgets('발급 성공 시 토큰 표시', (tester) async {
      final client = _StubClient(
        memberships: [
          {
            'user_id': 'me',
            'tenant_id': 'hospA',
            'role': 'admin',
            'active': true,
            'joined_at': 1,
          },
        ],
        createResponse: {'token': 'abcdef1234567890', 'status': 'pending'},
      );
      await tester.pumpWidget(_wrap(const InviteCreateScreen(), client));
      await tester.pumpAndSettle();

      // 드롭다운 선택
      await tester.tap(find.byKey(const Key('tenant-dropdown')));
      await tester.pumpAndSettle();
      await tester.tap(find.text('hospA (admin)').last);
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('submit-invite')));
      await tester.pumpAndSettle();

      expect(find.text('초대 토큰 발급 완료'), findsOneWidget);
      expect(find.text('abcdef1234567890'), findsOneWidget);
      expect(client.lastCreatedTenantId, 'hospA');
    });

    testWidgets('발급 실패 시 에러 표시', (tester) async {
      final client = _StubClient(
        memberships: [
          {
            'user_id': 'me',
            'tenant_id': 'hospA',
            'role': 'admin',
            'active': true,
            'joined_at': 1,
          },
        ],
        failCreate: true,
      );
      await tester.pumpWidget(_wrap(const InviteCreateScreen(), client));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('tenant-dropdown')));
      await tester.pumpAndSettle();
      await tester.tap(find.text('hospA (admin)').last);
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(const Key('submit-invite')));
      await tester.pumpAndSettle();

      expect(find.textContaining('생성 실패'), findsOneWidget);
    });
  });

  group('InviteAcceptScreen', () {
    testWidgets('빈 토큰 제출 시 에러', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(_wrap(const InviteAcceptScreen(), client));
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(const Key('submit-accept')));
      await tester.pumpAndSettle();
      expect(find.text('토큰을 입력하세요'), findsOneWidget);
    });

    testWidgets('수락 성공 시 결과 표시', (tester) async {
      final client = _StubClient(
        acceptResponse: {'tenant_id': 'hospA', 'role': 'medical_staff'},
      );
      await tester.pumpWidget(_wrap(const InviteAcceptScreen(), client));
      await tester.pumpAndSettle();

      await tester.enterText(
          find.byKey(const Key('invite-token-field')), 'token-xyz');
      await tester.tap(find.byKey(const Key('submit-accept')));
      await tester.pumpAndSettle();

      expect(find.textContaining('가입 완료: hospA'), findsOneWidget);
      expect(client.lastAcceptedToken, 'token-xyz');
    });

    testWidgets('수락 실패 시 에러 표시', (tester) async {
      final client = _StubClient(failAccept: true);
      await tester.pumpWidget(_wrap(const InviteAcceptScreen(), client));
      await tester.pumpAndSettle();

      await tester.enterText(
          find.byKey(const Key('invite-token-field')), 'bad-token');
      await tester.tap(find.byKey(const Key('submit-accept')));
      await tester.pumpAndSettle();

      expect(find.textContaining('수락 실패'), findsOneWidget);
    });

    testWidgets('initialToken 자동 채움', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(
          _wrap(const InviteAcceptScreen(initialToken: 'preset-token'), client));
      await tester.pumpAndSettle();
      // controller 의 텍스트가 자동 설정되었는지 확인
      final field = tester.widget<TextField>(
          find.byKey(const Key('invite-token-field')));
      expect(field.controller?.text, 'preset-token');
    });
  });
}
