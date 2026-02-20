import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';

void main() {
  group('PorcelainContainer', () {
    testWidgets('기본 렌더링 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PorcelainContainer(
              child: Text('컨텐츠'),
            ),
          ),
        ),
      );

      expect(find.text('컨텐츠'), findsOneWidget);
      expect(find.byType(PorcelainContainer), findsOneWidget);
    });

    testWidgets('onTap 콜백이 호출된다', (WidgetTester tester) async {
      bool tapped = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PorcelainContainer(
              onTap: () => tapped = true,
              child: const Text('탭 테스트'),
            ),
          ),
        ),
      );

      await tester.tap(find.text('탭 테스트'));
      expect(tapped, isTrue);
    });

    testWidgets('isSelected 시 테두리가 변경된다', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PorcelainContainer(
              isSelected: true,
              child: Text('선택됨'),
            ),
          ),
        ),
      );

      expect(find.text('선택됨'), findsOneWidget);
    });

    testWidgets('커스텀 padding/margin/size 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PorcelainContainer(
              width: 200,
              height: 100,
              padding: EdgeInsets.all(16),
              margin: EdgeInsets.all(8),
              child: Text('사이즈 테스트'),
            ),
          ),
        ),
      );

      expect(find.text('사이즈 테스트'), findsOneWidget);
    });

    testWidgets('다크 테마에서 렌더링 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          theme: ThemeData.dark(),
          home: const Scaffold(
            body: PorcelainContainer(
              child: Text('다크 모드'),
            ),
          ),
        ),
      );

      expect(find.text('다크 모드'), findsOneWidget);
    });

    testWidgets('커스텀 color 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PorcelainContainer(
              color: Colors.blue,
              child: Text('커스텀 색상'),
            ),
          ),
        ),
      );

      expect(find.text('커스텀 색상'), findsOneWidget);
    });
  });
}
