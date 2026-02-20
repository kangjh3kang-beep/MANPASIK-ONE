import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

void main() {
  group('ScaleButton', () {
    testWidgets('ScaleButton이 올바르게 렌더링된다', (WidgetTester tester) async {
      bool pressed = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ScaleButton(
              onPressed: () => pressed = true,
              child: const Text('테스트 버튼'),
            ),
          ),
        ),
      );

      expect(find.text('테스트 버튼'), findsOneWidget);
      expect(pressed, isFalse);
    });

    testWidgets('ScaleButton 탭 시 onPressed가 호출된다', (WidgetTester tester) async {
      bool pressed = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ScaleButton(
              onPressed: () => pressed = true,
              child: const Text('테스트 버튼'),
            ),
          ),
        ),
      );

      await tester.tap(find.text('테스트 버튼'));
      await tester.pump();
      expect(pressed, isTrue);
    });

    testWidgets('ScaleButton은 MouseRegion을 포함한다', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ScaleButton(
              onPressed: () {},
              child: const Text('Hover 테스트'),
            ),
          ),
        ),
      );

      expect(find.byType(MouseRegion), findsWidgets);
    });

    testWidgets('ScaleButton은 GestureDetector를 포함한다', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ScaleButton(
              onPressed: () {},
              child: const Text('Gesture 테스트'),
            ),
          ),
        ),
      );

      expect(find.byType(GestureDetector), findsWidgets);
    });

    testWidgets('ScaleButton 커스텀 scale 파라미터', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: ScaleButton(
              onPressed: () {},
              scale: 0.9,
              child: const Text('Custom Scale'),
            ),
          ),
        ),
      );
      expect(find.text('Custom Scale'), findsOneWidget);
    });
  });
}
