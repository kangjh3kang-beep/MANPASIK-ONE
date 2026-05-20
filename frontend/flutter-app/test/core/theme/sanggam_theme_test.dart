import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

void main() {
  group('SanggamTheme 색상 상수', () {
    test('background 색상 확인', () {
      expect(SanggamTheme.background, const Color(0xFF0B1021));
    });

    test('primary 색상 확인 (상감 금)', () {
      expect(SanggamTheme.primary, const Color(0xFFD4AF37));
    });

    test('jagaeCyan 색상 확인', () {
      expect(SanggamTheme.jagaeCyan, const Color(0xFF00FFFF));
    });

    test('jagaeMagenta 색상 확인', () {
      expect(SanggamTheme.jagaeMagenta, const Color(0xFFFF00FF));
    });

    test('surface 색상 확인', () {
      expect(SanggamTheme.surface, const Color(0xFF111A30));
    });

    test('error 색상 확인', () {
      expect(SanggamTheme.error, const Color(0xFFFF4D4D));
    });

    test('onSurfaceDim 색상 확인', () {
      expect(SanggamTheme.onSurfaceDim, const Color(0xFF8899BB));
    });

    test('onPrimary 색상 확인', () {
      expect(SanggamTheme.onPrimary, const Color(0xFF0B1021));
    });
  });

  group('SanggamTheme 그래디언트', () {
    test('jagaeGradient는 시안→마젠타', () {
      expect(SanggamTheme.jagaeGradient.colors, hasLength(2));
      expect(SanggamTheme.jagaeGradient.colors[0], SanggamTheme.jagaeCyan);
      expect(SanggamTheme.jagaeGradient.colors[1], SanggamTheme.jagaeMagenta);
    });

    test('goldGradient 존재 확인', () {
      expect(SanggamTheme.goldGradient, isA<LinearGradient>());
      expect(SanggamTheme.goldGradient.colors, isNotEmpty);
    });
  });

  group('SanggamTheme 상수 관계', () {
    test('surfaceVariant는 surface와 다르다', () {
      expect(SanggamTheme.surfaceVariant, isNot(equals(SanggamTheme.surface)));
    });

    test('onPrimary는 background와 같다', () {
      expect(SanggamTheme.onPrimary, SanggamTheme.background);
    });

    test('jagaeGradientDiagonal는 3색 그래디언트', () {
      expect(SanggamTheme.jagaeGradientDiagonal.colors, hasLength(3));
    });

    test('색상 대비: primary는 background와 다르다', () {
      expect(SanggamTheme.primary, isNot(equals(SanggamTheme.background)));
    });

    test('error 색상은 빨간 계열', () {
      expect(SanggamTheme.error.red, greaterThan(200));
    });
  });
}
