// AppTheme 테스트
// - 라이트/다크 테마 생성, 브랜드 컬러, Material3 설정 확인

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:manpasik/core/theme/app_theme.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  // Google Fonts 네트워크 요청 비활성화 (테스트 환경)
  setUpAll(() {
    GoogleFonts.config.allowRuntimeFetching = false;
  });

  group('AppTheme 브랜드 컬러 테스트', () {
    test('celadonTeal은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.celadonTeal, const Color(0xFF00897B));
    });

    test('dancheongRed은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.dancheongRed, const Color(0xFFFF4D4D));
    });

    test('inkBlack은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.inkBlack, const Color(0xFF020617));
    });

    test('hanjiWhite는 올바른 색상 값이어야 한다', () {
      expect(AppTheme.hanjiWhite, const Color(0xFFF8FAFC));
    });

    test('deepSeaBlue는 올바른 색상 값이어야 한다', () {
      expect(AppTheme.deepSeaBlue, const Color(0xFF050B14));
    });
  });

  group('AppTheme.light 테스트', () {
    testWidgets('라이트 테마는 Material3를 사용해야 한다', (tester) async {
      expect(AppTheme.light.useMaterial3, isTrue);
    });

    testWidgets('라이트 테마의 brightness는 light이어야 한다', (tester) async {
      expect(AppTheme.light.brightness, Brightness.light);
    });

    testWidgets('라이트 테마의 scaffold 배경색은 hanjiWhite이어야 한다', (tester) async {
      expect(AppTheme.light.scaffoldBackgroundColor, AppTheme.hanjiWhite);
    });

    testWidgets('라이트 테마의 colorScheme이 존재해야 한다', (tester) async {
      expect(AppTheme.light.colorScheme, isNotNull);
      expect(AppTheme.light.colorScheme.brightness, Brightness.light);
    });
  });

  group('AppTheme.dark 테스트', () {
    testWidgets('다크 테마는 Material3를 사용해야 한다', (tester) async {
      expect(AppTheme.dark.useMaterial3, isTrue);
    });

    testWidgets('다크 테마의 brightness는 dark이어야 한다', (tester) async {
      expect(AppTheme.dark.brightness, Brightness.dark);
    });

    testWidgets('다크 테마의 scaffold 배경색은 deepSeaBlue이어야 한다', (tester) async {
      expect(AppTheme.dark.scaffoldBackgroundColor, AppTheme.deepSeaBlue);
    });

    testWidgets('다크 테마의 colorScheme이 존재해야 한다', (tester) async {
      expect(AppTheme.dark.colorScheme, isNotNull);
      expect(AppTheme.dark.colorScheme.brightness, Brightness.dark);
    });
  });
}
