import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/measurement_card.dart';

void main() {
  group('MeasurementCard', () {
    testWidgets('정상 상태 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19, 10, 30),
              value: 95.0,
              unit: 'mg/dL',
              resultType: 'normal',
            ),
          ),
        ),
      );
      expect(find.text('정상'), findsOneWidget);
      expect(find.text('95.0'), findsOneWidget);
      expect(find.text('mg/dL'), findsOneWidget);
    });

    testWidgets('주의 상태 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19),
              value: 130.0,
              unit: 'mg/dL',
              resultType: 'warning',
            ),
          ),
        ),
      );
      expect(find.text('주의'), findsOneWidget);
      expect(find.text('130.0'), findsOneWidget);
    });

    testWidgets('위험 상태 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19),
              value: 250.0,
              unit: 'mg/dL',
              resultType: 'danger',
            ),
          ),
        ),
      );
      expect(find.text('위험'), findsOneWidget);
      expect(find.text('250.0'), findsOneWidget);
    });

    testWidgets('onTap 콜백이 호출된다', (WidgetTester tester) async {
      bool tapped = false;
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19),
              value: 100.0,
              unit: 'mg/dL',
              resultType: 'normal',
              onTap: () => tapped = true,
            ),
          ),
        ),
      );
      await tester.tap(find.byType(MeasurementCard));
      await tester.pump();
      expect(tapped, isTrue);
    });

    testWidgets('날짜 포맷 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19, 14, 30),
              value: 80.0,
              unit: 'mmHg',
              resultType: 'normal',
            ),
          ),
        ),
      );
      expect(find.text('04월 19일 14:30'), findsOneWidget);
    });

    testWidgets('소수점 1자리로 표시된다', (WidgetTester tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: MeasurementCard(
              date: DateTime(2026, 4, 19),
              value: 98.765,
              unit: 'mg/dL',
              resultType: 'normal',
            ),
          ),
        ),
      );
      expect(find.text('98.8'), findsOneWidget);
    });
  });
}
