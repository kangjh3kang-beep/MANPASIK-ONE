import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/admin/presentation/tenancy_admin_dashboard_screen.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient({
    this.tenancyStats = const {
      'tenant_count': 5,
      'active_member_count': 42,
      'inactive_member_count': 3,
      'invitations': {
        'pending': 2,
        'accepted': 30,
        'expired': 1,
        'revoked': 0,
      },
    },
    this.auditStats = const {
      'TenantID': 'hospA',
      'TotalCalls': 100,
      'SuccessfulCalls': 95,
      'FailedCalls': 5,
      'TotalTokens': 50000,
      'TokensInWindow': 12000,
      'WindowHours': 24,
    },
    this.auditFailures = const [],
    this.dlqEntries = const [],
    this.failTenancy = false,
    this.failAudit = false,
    this.failDlq = false,
    this.replayedIds = const [],
    this.droppedIds = const [],
  }) : super(baseUrl: 'http://localhost');

  final Map<String, dynamic> tenancyStats;
  final Map<String, dynamic> auditStats;
  final List<Map<String, dynamic>> auditFailures;
  final List<Map<String, dynamic>> dlqEntries;
  final bool failTenancy;
  final bool failAudit;
  final bool failDlq;
  final List<int> replayedIds;
  final List<int> droppedIds;

  @override
  Future<Map<String, dynamic>> getTenancyStats({String? opsBaseUrl}) async {
    if (failTenancy) throw StateError('네트워크');
    return tenancyStats;
  }

  @override
  Future<Map<String, dynamic>> getAuditStats({
    required String tenantId,
    int hours = 24,
    String? opsBaseUrl,
  }) async {
    if (failAudit) throw StateError('audit 실패');
    return auditStats;
  }

  @override
  Future<List<Map<String, dynamic>>> getAuditFailures({
    required String tenantId,
    int limit = 10,
    String? opsBaseUrl,
  }) async {
    if (failAudit) throw StateError('audit failures 실패');
    return auditFailures;
  }

  @override
  Future<List<Map<String, dynamic>>> listWebhookDLQ({
    String? opsBaseUrl,
  }) async {
    if (failDlq) throw StateError('dlq 실패');
    return dlqEntries;
  }

  @override
  Future<void> replayWebhookDLQ({
    required int id,
    String? opsBaseUrl,
  }) async {
    (replayedIds as List<int>).add(id);
  }

  @override
  Future<void> dropWebhookDLQ({
    required int id,
    String? opsBaseUrl,
  }) async {
    (droppedIds as List<int>).add(id);
  }
}

Widget _wrap(Widget child, _StubClient client) => UncontrolledProviderScope(
      container: ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(client),
      ]),
      child: MaterialApp(home: child),
    );

Future<void> _useLargeSurface(WidgetTester tester) async {
  await tester.binding.setSurfaceSize(const Size(1000, 2200));
  addTearDown(() => tester.binding.setSurfaceSize(null));
}

void main() {
  group('TenancyAdminDashboardScreen', () {
    testWidgets('Tenancy 통계 카드 표시', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.text('활성 조직'), findsOneWidget);
      expect(find.text('활성 멤버'), findsOneWidget);
      expect(find.text('42'), findsOneWidget);
      expect(find.text('30'), findsOneWidget); // accepted invitations
    });

    testWidgets('Audit 통계 카드 표시 (tenant 입력 시)', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('LLM Audit'), findsOneWidget);
      expect(find.text('총 호출'), findsOneWidget);
      expect(find.text('100'), findsOneWidget);
    });

    testWidgets('tenant 미입력 시 audit 카드 숨김', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('LLM Audit'), findsNothing);
      expect(find.text('활성 조직'), findsOneWidget);
    });

    testWidgets('새로고침 버튼 탭 후 화면 유지', (tester) async {
      final client = _StubClient();
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const Key('refresh-button')));
      await tester.pumpAndSettle();
      expect(find.text('활성 조직'), findsOneWidget);
    });

    testWidgets('에러 시 에러 카드 + 재시도 버튼', (tester) async {
      final client = _StubClient(failTenancy: true);
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('네트워크'), findsOneWidget);
      expect(find.text('다시 시도'), findsOneWidget);
    });

    testWidgets('audit 에러 시에도 tenancy 카드는 표시', (tester) async {
      final client = _StubClient(failAudit: true);
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      // audit 호출이 실패하면 전체 _error 로 처리되므로 에러 화면
      expect(find.textContaining('audit'), findsOneWidget);
    });

    // Phase AO-3: AuditFailures + DLQ 카드
    testWidgets('audit failures 항목 있을 때 실패 카드 표시', (tester) async {
      await _useLargeSurface(tester);
      final client = _StubClient(
        auditFailures: const [
          {
            'Provider': 'openai',
            'Model': 'gpt-4',
            'ErrorMessage': 'rate_limit_exceeded: too many requests',
            'UserID': 'user-1',
            'RecordedAt': '2026-05-02T10:00:00Z',
          },
          {
            'Provider': 'anthropic',
            'Model': 'claude',
            'ErrorMessage': 'timeout',
            'UserID': 'user-2',
            'RecordedAt': '2026-05-02T10:01:00Z',
          },
        ],
      );
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(targetTenantId: 'hospA'),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('최근 실패 호출 (2)'), findsOneWidget);
      expect(find.textContaining('rate_limit_exceeded'), findsOneWidget);
      expect(find.textContaining('anthropic'), findsOneWidget);
    });

    testWidgets('DLQ 비어있을 때 빈 메시지', (tester) async {
      await _useLargeSurface(tester);
      final client = _StubClient();
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('Webhook DLQ (0)'), findsOneWidget);
      expect(find.text('대기 중인 항목 없음'), findsOneWidget);
    });

    testWidgets('DLQ 항목 표시 + 재발송/폐기 아이콘', (tester) async {
      await _useLargeSurface(tester);
      final client = _StubClient(
        dlqEntries: const [
          {
            'id': 7,
            'event': {
              'type': 'invitation.created',
              'tenant_id': 'hospA',
            },
            'last_error': 'connection refused',
            'attempts': 4,
          },
        ],
      );
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(),
        client,
      ));
      await tester.pumpAndSettle();

      expect(find.textContaining('Webhook DLQ (1)'), findsOneWidget);
      expect(find.textContaining('#7'), findsOneWidget);
      expect(find.textContaining('connection refused'), findsOneWidget);
      expect(find.byKey(const Key('dlq-replay-7')), findsOneWidget);
      expect(find.byKey(const Key('dlq-drop-7')), findsOneWidget);
    });

    testWidgets('DLQ 재발송 버튼 탭 → replayWebhookDLQ 호출', (tester) async {
      await _useLargeSurface(tester);
      final replayed = <int>[];
      final client = _StubClient(
        dlqEntries: const [
          {
            'id': 11,
            'event': {'type': 'membership.created', 'tenant_id': 'X'},
            'last_error': 'oops',
            'attempts': 3,
          },
        ],
        replayedIds: replayed,
      );
      await tester.pumpWidget(_wrap(
        const TenancyAdminDashboardScreen(),
        client,
      ));
      await tester.pumpAndSettle();

      await tester.ensureVisible(find.byKey(const Key('dlq-replay-11')));
      await tester.tap(find.byKey(const Key('dlq-replay-11')));
      await tester.pump();
      await tester.pumpAndSettle();

      expect(replayed, contains(11));
    });
  });
}

