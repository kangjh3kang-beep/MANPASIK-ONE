import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';

void main() {
  group('PrimaryButton', () {
    testWidgets('기본 텍스트가 렌더링된다', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PrimaryButton(text: '로그인'),
          ),
        ),
      );
      expect(find.text('로그인'), findsOneWidget);
      expect(find.byType(PrimaryButton), findsOneWidget);
    });

    testWidgets('onPressed 콜백이 호출된다', (WidgetTester tester) async {
      bool pressed = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '확인',
              onPressed: () => pressed = true,
            ),
          ),
        ),
      );
      await tester.tap(find.text('확인'));
      expect(pressed, isTrue);
    });

    testWidgets('onPressed가 null이면 비활성 상태', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: PrimaryButton(text: '비활성'),
          ),
        ),
      );
      // onPressed가 null이면 AnimatedOpacity가 0.5가 되어야 한다
      final opacity = tester.widget<AnimatedOpacity>(
        find.byType(AnimatedOpacity),
      );
      expect(opacity.opacity, 0.5);
    });

    testWidgets('isLoading 시 CircularProgressIndicator가 표시된다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '로딩 중',
              isLoading: true,
              onPressed: () {},
            ),
          ),
        ),
      );
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      // 로딩 중에는 텍스트가 보이지 않아야 한다
      expect(find.text('로딩 중'), findsNothing);
    });

    testWidgets('isLoading 시 탭해도 콜백이 호출되지 않는다',
        (WidgetTester tester) async {
      bool pressed = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '로딩',
              isLoading: true,
              onPressed: () => pressed = true,
            ),
          ),
        ),
      );
      await tester.tap(find.byType(PrimaryButton));
      expect(pressed, isFalse);
    });

    testWidgets('아이콘이 표시된다', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '아이콘 버튼',
              icon: Icons.login,
              onPressed: () {},
            ),
          ),
        ),
      );
      expect(find.byIcon(Icons.login), findsOneWidget);
      expect(find.text('아이콘 버튼'), findsOneWidget);
    });

    testWidgets('아이콘 없으면 Icon 위젯이 없다', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '텍스트만',
              onPressed: () {},
            ),
          ),
        ),
      );
      expect(find.byType(Icon), findsNothing);
    });

    testWidgets('PrimaryButton 위젯 트리에 Container가 존재한다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: PrimaryButton(
              text: '높이 테스트',
              onPressed: () {},
            ),
          ),
        ),
      );
      expect(
        find.descendant(
          of: find.byType(PrimaryButton),
          matching: find.byType(Container),
        ),
        findsOneWidget,
      );
    });
  });
}
