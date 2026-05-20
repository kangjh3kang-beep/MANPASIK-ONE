import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/holo_glass_card.dart';

void main() {
  group('HoloGlassCard', () {
    testWidgets('기본 렌더링 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: HoloGlassCard(
              child: Text('홀로그래픽'),
            ),
          ),
        ),
      );
      expect(find.text('홀로그래픽'), findsOneWidget);
      expect(find.byType(HoloGlassCard), findsOneWidget);
    });

    testWidgets('onTap 콜백이 호출된다', (WidgetTester tester) async {
      bool tapped = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: HoloGlassCard(
              onTap: () => tapped = true,
              child: const Text('탭 가능'),
            ),
          ),
        ),
      );
      await tester.tap(find.text('탭 가능'));
      expect(tapped, isTrue);
    });

    testWidgets('isActive=true 시 정상 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: HoloGlassCard(
              isActive: true,
              child: Text('활성 카드'),
            ),
          ),
        ),
      );
      expect(find.text('활성 카드'), findsOneWidget);
      // isActive=true는 추가 네온 글로우 DecoratedBox를 포함
      expect(find.byType(HoloGlassCard), findsOneWidget);
    });

    testWidgets('isActive=false 시 정상 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: HoloGlassCard(
              isActive: false,
              child: Text('비활성 카드'),
            ),
          ),
        ),
      );
      expect(find.text('비활성 카드'), findsOneWidget);
      expect(find.byType(HoloGlassCard), findsOneWidget);
    });

    testWidgets('커스텀 width/height 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: Center(
              child: HoloGlassCard(
                width: 300,
                height: 200,
                child: Text('크기 테스트'),
              ),
            ),
          ),
        ),
      );
      expect(find.text('크기 테스트'), findsOneWidget);
    });

    testWidgets('커스텀 glowColor 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: HoloGlassCard(
              isActive: true,
              glowColor: Colors.purple,
              child: Text('커스텀 글로우'),
            ),
          ),
        ),
      );
      expect(find.text('커스텀 글로우'), findsOneWidget);
    });
  });
}
