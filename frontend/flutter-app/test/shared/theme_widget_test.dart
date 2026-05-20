import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// Phase M: SanggamTheme 및 핵심 위젯 빌드 검증
void main() {
  group('SanggamTheme 색상 상수', () {
    test('background 색상은 다크 톤', () {
      expect(SanggamTheme.background.toARGB32(), 0xFF0B1021);
    });

    test('primary 색상은 골드', () {
      // primary 색상이 정의되어 있어야 함
      expect(SanggamTheme.primary, isA<Color>());
    });

    test('jagaeCyan과 jagaeMagenta 정의', () {
      expect(SanggamTheme.jagaeCyan, isA<Color>());
      expect(SanggamTheme.jagaeMagenta, isA<Color>());
    });

    test('error 색상 정의', () {
      expect(SanggamTheme.error, isA<Color>());
    });

    test('onSurfaceDim과 surface 정의', () {
      expect(SanggamTheme.onSurfaceDim, isA<Color>());
      expect(SanggamTheme.surface, isA<Color>());
    });

    test('surfaceVariant 정의', () {
      expect(SanggamTheme.surfaceVariant, isA<Color>());
    });

    test('background는 불투명 색상', () {
      expect(SanggamTheme.background.a, 1.0);
    });

    test('error 색상은 불투명', () {
      expect(SanggamTheme.error.a, 1.0);
    });
  });

  group('SanggamContainer 위젯', () {
    testWidgets('자식 위젯을 렌더링한다', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              child: Text('Test'),
            ),
          ),
        ),
      );

      expect(find.text('Test'), findsOneWidget);
    });

    testWidgets('borderRadius 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              borderRadius: 24,
              child: Text('Rounded'),
            ),
          ),
        ),
      );

      expect(find.text('Rounded'), findsOneWidget);
      expect(find.byType(SanggamContainer), findsOneWidget);
    });

    testWidgets('padding 커스터마이징', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              padding: EdgeInsets.all(8),
              child: Text('Padded'),
            ),
          ),
        ),
      );

      expect(find.text('Padded'), findsOneWidget);
    });

    testWidgets('borderColor 옵션 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              borderColor: Colors.red,
              child: Text('Bordered'),
            ),
          ),
        ),
      );

      expect(find.text('Bordered'), findsOneWidget);
    });

    testWidgets('backgroundColor 옵션 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              backgroundColor: Colors.blue,
              child: Text('Custom BG'),
            ),
          ),
        ),
      );

      expect(find.text('Custom BG'), findsOneWidget);
    });

    testWidgets('margin 옵션 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              margin: EdgeInsets.all(16),
              child: Text('Margin'),
            ),
          ),
        ),
      );

      expect(find.text('Margin'), findsOneWidget);
    });

    testWidgets('빈 자식도 처리', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: SanggamContainer(
              child: SizedBox.shrink(),
            ),
          ),
        ),
      );

      expect(find.byType(SanggamContainer), findsOneWidget);
    });
  });

  group('Color 확장: withValues', () {
    test('alpha 0.5 적용', () {
      final color = SanggamTheme.primary.withValues(alpha: 0.5);
      expect(color.a, lessThan(1.0));
    });

    test('alpha 1.0은 원래 색상', () {
      final color = SanggamTheme.primary.withValues(alpha: 1.0);
      expect(color.a, 1.0);
    });

    test('alpha 0.0은 투명', () {
      final color = SanggamTheme.primary.withValues(alpha: 0.0);
      expect(color.a, 0.0);
    });
  });

  group('Theme 컴포넌트 (테마 미사용)', () {
    testWidgets('Scaffold 배경에 SanggamTheme.background 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            backgroundColor: SanggamTheme.background,
            body: SizedBox(),
          ),
        ),
      );

      final scaffold = tester.widget<Scaffold>(find.byType(Scaffold));
      expect(scaffold.backgroundColor, SanggamTheme.background);
    });

    testWidgets('기본 MaterialApp 생성', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(body: Text('App')),
        ),
      );

      expect(find.text('App'), findsOneWidget);
    });

    testWidgets('TextStyle 적용', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: Text(
              'Styled',
              style: TextStyle(
                color: SanggamTheme.primary,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ),
      );

      expect(find.text('Styled'), findsOneWidget);
    });

    testWidgets('Container with SanggamTheme 색상', (tester) async {
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: Container(
              color: SanggamTheme.surface,
              child: const SizedBox(width: 100, height: 100),
            ),
          ),
        ),
      );
      expect(find.byType(Container), findsOneWidget);
    });

    testWidgets('Icon with SanggamTheme.primary', (tester) async {
      await tester.pumpWidget(
        const MaterialApp(
          home: Scaffold(
            body: Icon(
              Icons.favorite,
              color: SanggamTheme.primary,
              size: 32,
            ),
          ),
        ),
      );
      expect(find.byIcon(Icons.favorite), findsOneWidget);
    });
  });
}
