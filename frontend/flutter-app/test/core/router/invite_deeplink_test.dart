import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/settings/presentation/invite_accept_screen.dart';

class _StubClient extends ManPaSikRestClient {
  _StubClient() : super(baseUrl: 'http://localhost');

  @override
  Future<List<Map<String, dynamic>>> getMyMemberships() async => [];

  @override
  Future<Map<String, dynamic>> acceptTenantInvitation(String token) async => {
        'tenant_id': 'hospA',
        'role': 'member',
      };
}

/// 단순 라우터 — 앱 전체 라우터를 로드하지 않고 invite 라우트만 분리하여 검증.
GoRouter _testRouter() => GoRouter(
      initialLocation: '/invite/preset-token-123',
      routes: [
        GoRoute(
          path: '/',
          builder: (_, __) => const Text('home'),
        ),
        GoRoute(
          path: '/invite/:token',
          builder: (context, state) => InviteAcceptScreen(
            initialToken: state.pathParameters['token'],
          ),
        ),
        GoRoute(
          path: '/invite',
          builder: (_, __) => const InviteAcceptScreen(),
        ),
      ],
    );

void main() {
  group('Invite deep link', () {
    testWidgets('/invite/:token 라우트가 InviteAcceptScreen 으로 이동', (tester) async {
      await tester.pumpWidget(
        UncontrolledProviderScope(
          container: ProviderContainer(overrides: [
            restClientProvider.overrideWithValue(_StubClient()),
          ]),
          child: MaterialApp.router(routerConfig: _testRouter()),
        ),
      );
      await tester.pumpAndSettle();

      // 토큰이 자동으로 채워졌는지 확인
      final field = tester.widget<TextField>(
          find.byKey(const Key('invite-token-field')));
      expect(field.controller?.text, 'preset-token-123');
    });

    testWidgets('/invite (path 없음) 도 작동', (tester) async {
      final router = GoRouter(
        initialLocation: '/invite',
        routes: [
          GoRoute(
            path: '/invite',
            builder: (_, __) => const InviteAcceptScreen(),
          ),
          GoRoute(
            path: '/invite/:token',
            builder: (context, state) => InviteAcceptScreen(
              initialToken: state.pathParameters['token'],
            ),
          ),
        ],
      );
      await tester.pumpWidget(
        UncontrolledProviderScope(
          container: ProviderContainer(overrides: [
            restClientProvider.overrideWithValue(_StubClient()),
          ]),
          child: MaterialApp.router(routerConfig: router),
        ),
      );
      await tester.pumpAndSettle();

      // 토큰 미주입
      final field = tester.widget<TextField>(
          find.byKey(const Key('invite-token-field')));
      expect(field.controller?.text, '');
      expect(find.text('조직 초대 수락'), findsOneWidget);
    });

    testWidgets('초대 토큰이 다양한 길이여도 라우팅 통과', (tester) async {
      const longToken =
          'aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899';
      final router = GoRouter(
        initialLocation: '/invite/$longToken',
        routes: [
          GoRoute(
            path: '/invite/:token',
            builder: (context, state) => InviteAcceptScreen(
              initialToken: state.pathParameters['token'],
            ),
          ),
        ],
      );
      await tester.pumpWidget(
        UncontrolledProviderScope(
          container: ProviderContainer(overrides: [
            restClientProvider.overrideWithValue(_StubClient()),
          ]),
          child: MaterialApp.router(routerConfig: router),
        ),
      );
      await tester.pumpAndSettle();

      final field = tester.widget<TextField>(
          find.byKey(const Key('invite-token-field')));
      expect(field.controller?.text, longToken);
    });
  });
}
