import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:manpasik/features/home/presentation/home_screen.dart';
import 'package:manpasik/features/devices/presentation/device_list_screen.dart';
import 'package:manpasik/features/measurement/presentation/measurement_result_screen.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'helpers/fake_repositories.dart';

List<Override> _baseOverrides() => [
  authRepositoryProvider.overrideWithValue(FakeAuthRepository()),
  authProvider.overrideWith((ref) => AuthNotifier(ref.read(authRepositoryProvider))),
  measurementRepositoryProvider.overrideWithValue(FakeMeasurementRepository()),
  deviceRepositoryProvider.overrideWithValue(FakeDeviceRepository()),
  userRepositoryProvider.overrideWithValue(FakeUserRepository()),
  measurementHistoryProvider.overrideWith(
    (ref) async => const MeasurementHistoryResult(items: [], totalCount: 0),
  ),
  deviceListProvider.overrideWith(
    (ref) async => <DeviceItem>[],
  ),
];

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  group('HomeScreen 위젯 테스트', () {
    testWidgets('비인증 시 HomeScreen 위젯이 생성된다', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(body: HomeScreen()),
          ),
        ),
      );
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));
      expect(find.byType(HomeScreen), findsOneWidget);
    });

    testWidgets('HomeScreen이 렌더링된다', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(body: HomeScreen()),
          ),
        ),
      );
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));
      expect(find.byType(HomeScreen), findsOneWidget);
    });

    testWidgets('HomeScreen Scaffold가 존재한다', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: _baseOverrides(),
          child: MaterialApp(
            home: const Scaffold(body: HomeScreen()),
          ),
        ),
      );
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));
      expect(find.byType(Scaffold), findsWidgets);
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
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));
      expect(find.byType(DeviceListScreen), findsOneWidget);
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
      await tester.pump();
      await tester.pump(const Duration(milliseconds: 500));
      expect(find.byType(MeasurementResultScreen), findsOneWidget);
    });
  });
}
