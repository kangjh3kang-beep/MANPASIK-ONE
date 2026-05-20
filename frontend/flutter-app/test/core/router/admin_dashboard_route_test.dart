import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/admin/presentation/tenancy_admin_dashboard_screen.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient() : super(baseUrl: 'http://localhost');

  @override
  Future<Map<String, dynamic>> getTenancyStats({String? opsBaseUrl}) async => {
        'tenant_count': 1,
        'active_member_count': 1,
        'inactive_member_count': 0,
        'invitations': {'pending': 0, 'accepted': 0, 'expired': 0, 'revoked': 0},
      };

  @override
  Future<Map<String, dynamic>> getAuditStats({
    required String tenantId,
    int hours = 24,
    String? opsBaseUrl,
  }) async =>
      {
        'TenantID': tenantId,
        'TotalCalls': 0,
        'WindowHours': hours,
      };

  @override
  Future<List<Map<String, dynamic>>> getAuditFailures({
    required String tenantId,
    int limit = 10,
    String? opsBaseUrl,
  }) async =>
      const [];

  @override
  Future<List<Map<String, dynamic>>> listWebhookDLQ({
    String? opsBaseUrl,
  }) async =>
      const [];
}

GoRouter _testRouter(String initial) => GoRouter(
      initialLocation: initial,
      routes: [
        GoRoute(
          path: '/admin/tenancy-dashboard',
          builder: (context, state) => TenancyAdminDashboardScreen(
            opsBaseUrl: state.uri.queryParameters['ops_url'],
            targetTenantId: state.uri.queryParameters['tenant'],
          ),
        ),
      ],
    );

Widget _wrap(GoRouter router) => UncontrolledProviderScope(
      container: ProviderContainer(overrides: [
        restClientProvider.overrideWithValue(_StubClient()),
      ]),
      child: MaterialApp.router(routerConfig: router),
    );

void main() {
  group('Admin tenancy dashboard route', () {
    testWidgets('/admin/tenancy-dashboard 기본 진입', (tester) async {
      await tester.pumpWidget(_wrap(_testRouter('/admin/tenancy-dashboard')));
      await tester.pumpAndSettle();
      expect(find.byType(TenancyAdminDashboardScreen), findsOneWidget);
    });

    testWidgets('query 파라미터 tenant 자동 주입', (tester) async {
      await tester.pumpWidget(_wrap(
        _testRouter('/admin/tenancy-dashboard?tenant=hospA'),
      ));
      await tester.pumpAndSettle();
      // tenant 입력란이 hospA 로 자동 채워졌는지
      final field = tester.widget<TextField>(
        find.byKey(const Key('tenant-input')),
      );
      expect(field.controller?.text, 'hospA');
    });

    testWidgets('query 파라미터 ops_url 전달', (tester) async {
      await tester.pumpWidget(_wrap(
        _testRouter('/admin/tenancy-dashboard?ops_url=http%3A%2F%2Fadmin.test%3A9100'),
      ));
      await tester.pumpAndSettle();
      // 위젯이 정상 렌더링되었는지만 확인 (URL 디코딩은 GoRouter 가 처리)
      expect(find.byType(TenancyAdminDashboardScreen), findsOneWidget);
    });
  });
}
