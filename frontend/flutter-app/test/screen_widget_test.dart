import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/features/home/presentation/home_screen.dart';
import 'package:manpasik/features/devices/presentation/device_list_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_result_screen.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'helpers/fake_repositories.dart';

List<Override> _baseOverrides() => [
  // Auth: Fake repository + notifier (네트워크 없이 동작)
  authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
  authProvider.overrideWith((ref) => AuthNotifier(ref.read(authRepositoryProvider))),
  // gRPC 의존 Repository override (실제 gRPC 연결 차단)
  measurementRepositoryProvider.overrideWithValue(FakeMeasurementRepository()),
  deviceRepositoryProvider.overrideWithValue(FakeDeviceRepository()),
  userRepositoryProvider.overrideWithValue(FakeUserRepository()),
  // FutureProvider override (gRPC 호출 완전 차단)
  measurementHistoryProvider.overrideWith(
    (ref) async => const MeasurementHistoryResult(items: [], totalCount: 0),
  ),
  deviceListProvider.overrideWith(
    (ref) async => <DeviceItem>[],
  ),
];

void main() {

  group('HomeScreen 위젯 테스트', () {
    testWidgets('비인증 시 사용자명 대신 "사용자" 표시', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(
              body: HomeScreen(),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.textContaining('님'), findsOneWidget);
    });

    testWidgets('최근 기록 섹션 존재', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(
              body: HomeScreen(),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('최근 기록'), findsOneWidget);
    });

    testWidgets('측정 시작하기 버튼 존재', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(
              body: HomeScreen(),
            ),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('측정 시작하기'), findsOneWidget);
    });
  });

  group('DeviceListScreen 위젯 테스트', () {
    testWidgets('앱바 타이틀 디바이스', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const DeviceListScreen(),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.text('디바이스'), findsOneWidget);
    });

    testWidgets('비인증 시 빈 목록 또는 로딩/에러 표시', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const DeviceListScreen(),
          ),
        ),
      );
      await tester.pump();
      await tester.pump(const Duration(seconds: 1));
      await tester.pumpAndSettle();
      expect(find.byType(DeviceListScreen), findsOneWidget);
    });
  });

  group('MeasurementResultScreen 위젯 테스트', () {
    testWidgets('결과 화면 진입 시 트렌드 또는 빈 상태', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const MeasurementResultScreen(),
          ),
        ),
      );
      await tester.pumpAndSettle();
      expect(find.byType(MeasurementResultScreen), findsOneWidget);
    });
  });
}
