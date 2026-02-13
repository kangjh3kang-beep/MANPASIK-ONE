// AppTheme 테스트
// - 라이트/다크 테마 생성, 브랜드 컬러, Material3 설정 확인

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:manpasik/core/theme/app_theme.dart';

void main() {
  // Google Fonts 네트워크 요청 비활성화 (테스트 환경)
  setUpAll(() {
    GoogleFonts.config.allowRuntimeFetching = false;
  });
  group('AppTheme 브랜드 컬러 테스트', () {
    // 브랜드 컬러 상수 값 확인
    test('celadonTeal(Primary)은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.celadonTeal, const Color(0xFF00897B));
    });

    test('dancheongRed(Error)은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.dancheongRed, const Color(0xFFD32F2F));
    });

    test('inkBlack은 올바른 색상 값이어야 한다', () {
      expect(AppTheme.inkBlack, const Color(0xFF121212));
    });

    test('hanjiWhite는 올바른 색상 값이어야 한다', () {
      expect(AppTheme.hanjiWhite, const Color(0xFFFAFAFA));
    });

    test('deepSeaBlue는 올바른 색상 값이어야 한다', () {
      expect(AppTheme.deepSeaBlue, const Color(0xFF1A237E));
    });
  });

  group('AppTheme.light 테스트', () {
    // Material3 사용 확인
    test('라이트 테마는 Material3를 사용해야 한다', () {
      expect(AppTheme.light.useMaterial3, isTrue);
    });

    // Brightness 확인
    test('라이트 테마의 brightness는 light이어야 한다', () {
      expect(AppTheme.light.brightness, Brightness.light);
    });

    // scaffoldBackgroundColor 확인
    test('라이트 테마의 scaffold 배경색은 hanjiWhite이어야 한다', () {
      expect(AppTheme.light.scaffoldBackgroundColor, AppTheme.hanjiWhite);
    });

    // primary 컬러 확인
    test('라이트 테마의 primary 컬러가 celadonTeal이어야 한다', () {
      expect(AppTheme.light.colorScheme.primary, AppTheme.celadonTeal);
    });
  });

  group('AppTheme.dark 테스트', () {
    // Material3 사용 확인
    test('다크 테마는 Material3를 사용해야 한다', () {
      expect(AppTheme.dark.useMaterial3, isTrue);
    });

    // Brightness 확인
    test('다크 테마의 brightness는 dark이어야 한다', () {
      expect(AppTheme.dark.brightness, Brightness.dark);
    });

    // scaffoldBackgroundColor 확인
    test('다크 테마의 scaffold 배경색은 0xFF121212이어야 한다', () {
      expect(
        AppTheme.dark.scaffoldBackgroundColor,
        const Color(0xFF121212),
      );
    });

    // primary 컬러 확인
    test('다크 테마의 primary 컬러가 celadonTeal이어야 한다', () {
      expect(AppTheme.dark.colorScheme.primary, AppTheme.celadonTeal);
    });
  });
}
