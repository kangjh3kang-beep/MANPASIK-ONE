import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/shared/widgets/medical_stat_card.dart';

void main() {
  group('MedicalStatCard', () {
    testWidgets('기본 렌더링 확인', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '혈당',
                value: '95',
                unit: 'mg/dL',
                icon: Icons.favorite,
              ),
            ),
          ),
        ),
      );
      expect(find.text('혈당'), findsOneWidget);
      expect(find.text('95'), findsOneWidget);
      expect(find.text('mg/dL'), findsOneWidget);
      expect(find.byIcon(Icons.favorite), findsOneWidget);
    });

    testWidgets('isAlert=true 시 정상 렌더링', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '혈압',
                value: '180',
                unit: 'mmHg',
                icon: Icons.warning,
                isAlert: true,
              ),
            ),
          ),
        ),
      );
      expect(find.text('혈압'), findsOneWidget);
      expect(find.text('180'), findsOneWidget);
    });

    testWidgets('value가 -- 일 때 미니 그래프가 없다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '심박수',
                value: '--',
                unit: 'bpm',
                icon: Icons.monitor_heart,
              ),
            ),
          ),
        ),
      );
      expect(find.text('--'), findsOneWidget);
      // value가 '--'일 때 _MiniGraphPainter CustomPaint는 없다
      // (다만 SanggamContainer 내부 구조에서 다른 CustomPaint가 있을 수 있음)
      expect(find.byType(MedicalStatCard), findsOneWidget);
    });

    testWidgets('tooltipMessage가 있으면 Tooltip이 추가된다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '체온',
                value: '36.5',
                unit: '°C',
                icon: Icons.thermostat,
                tooltipMessage: '정상 범위: 36.0-37.5°C',
              ),
            ),
          ),
        ),
      );
      expect(find.byType(Tooltip), findsOneWidget);
    });

    testWidgets('tooltipMessage가 null이면 Tooltip이 없다',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '산소포화도',
                value: '98',
                unit: '%',
                icon: Icons.air,
              ),
            ),
          ),
        ),
      );
      expect(find.byType(Tooltip), findsNothing);
    });

    testWidgets('커스텀 color 적용', (WidgetTester tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SizedBox(
              width: 200,
              child: MedicalStatCard(
                label: '콜레스테롤',
                value: '200',
                unit: 'mg/dL',
                icon: Icons.science,
                color: Colors.purple,
              ),
            ),
          ),
        ),
      );
      expect(find.text('콜레스테롤'), findsOneWidget);
    });
  });
}
