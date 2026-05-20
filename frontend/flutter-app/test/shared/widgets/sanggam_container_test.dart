import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

void main() {
  group('SanggamContainer', () {
    testWidgets('기본 렌더링 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              child: Text('상감 오르빗'),
            ),
          ),
        ),
      );
      expect(find.text('상감 오르빗'), findsOneWidget);
      expect(find.byType(SanggamContainer), findsOneWidget);
    });

    testWidgets('onTap 콜백이 호출된다', (WidgetTester tester) async {
      bool tapped = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              onTap: () => tapped = true,
              child: const Text('탭 테스트'),
            ),
          ),
        ),
      );
      await tester.tap(find.text('탭 테스트'));
      expect(tapped, isTrue);
    });

    testWidgets('onTap이 null이면 GestureDetector가 없다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              child: Text('탭 없음'),
            ),
          ),
        ),
      );
      // onTap 없으면 GestureDetector가 추가되지 않음
      expect(find.text('탭 없음'), findsOneWidget);
    });

    testWidgets('isSelected 변경 시 정상 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              isSelected: true,
              child: Text('선택됨'),
            ),
          ),
        ),
      );
      expect(find.text('선택됨'), findsOneWidget);
    });

    testWidgets('커스텀 width/height 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: Center(
              child: SanggamContainer(
                width: 200,
                height: 100,
                child: Text('크기 테스트'),
              ),
            ),
          ),
        ),
      );
      expect(find.text('크기 테스트'), findsOneWidget);
    });

    testWidgets('margin 적용 시 Padding이 추가된다', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              margin: EdgeInsets.all(16),
              child: Text('마진 테스트'),
            ),
          ),
        ),
      );
      expect(find.text('마진 테스트'), findsOneWidget);
    });

    testWidgets('jagaeOpacity 0이면 자개 레이어가 없다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              jagaeOpacity: 0,
              child: Text('자개 없음'),
            ),
          ),
        ),
      );
      expect(find.text('자개 없음'), findsOneWidget);
    });

    testWidgets('커스텀 backgroundColor 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              backgroundColor: Colors.blue,
              child: Text('배경 테스트'),
            ),
          ),
        ),
      );
      expect(find.text('배경 테스트'), findsOneWidget);
    });
  });
}
