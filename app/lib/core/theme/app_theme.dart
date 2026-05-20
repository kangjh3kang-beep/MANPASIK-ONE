import 'package:flutter/material.dart';

// =========================================================================
// ManPaSik (萬波息) - Premium Design System (앱 테마)
// Concept: "미래 지향적 바이오 헬스케어 (Neon Bio-Tech)"
// =========================================================================

class AppTheme {
  // --- 컬러 팰레트 (Dark Mode Base) ---
  static const Color backgroundDark = Color(0xFF0F172A); // Slate 900
  static const Color surfaceDark = Color(0xFF1E293B);    // Slate 800
  static const Color primaryNeonTeal = Color(0xFF2DD4BF);   // Teal 400 (정상/연결)
  static const Color primaryNeonPurple = Color(0xFFA855F7); // Purple 500 (AI 분석/핑거프린트)
  static const Color warningOrange = Color(0xFFFB923C);     // Orange 400 (주의)
  static const Color dangerRed = Color(0xFFEF4444);         // Red 500 (위험)
  
  // --- 그라데이션 ---
  static const LinearGradient premiumGradient = LinearGradient(
    colors: [primaryNeonPurple, primaryNeonTeal],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  static const LinearGradient scoreGoodGradient = LinearGradient(
    colors: [Color(0xFF10B981), Color(0xFF047857)], // Emerald
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  static ThemeData get darkTheme {
    return ThemeData(
      brightness: Brightness.dark,
      scaffoldBackgroundColor: backgroundDark,
      primaryColor: primaryNeonTeal,
      fontFamily: 'Inter', // 모던 타이포그래피 (pubspec.yaml에 폰트 애셋 추가 가정)
      
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: true,
        titleTextStyle: TextStyle(
          color: Colors.white,
          fontSize: 20,
          fontWeight: FontWeight.w600,
          letterSpacing: 1.2,
        ),
        iconTheme: IconThemeData(color: Colors.white),
      ),

      cardTheme: CardThemeData(
        color: surfaceDark.withValues(alpha: 0.7),
        elevation: 8,
        shadowColor: Colors.black45,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      ),

      floatingActionButtonTheme: const FloatingActionButtonThemeData(
        backgroundColor: primaryNeonTeal,
        foregroundColor: backgroundDark,
        elevation: 12,
      ),

      textTheme: const TextTheme(
        headlineLarge: TextStyle(color: Colors.white, fontSize: 32, fontWeight: FontWeight.bold),
        headlineMedium: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.w700),
        titleLarge: TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.w600),
        bodyLarge: TextStyle(color: Color(0xFFCBD5E1), fontSize: 16), // Slate 300
        bodyMedium: TextStyle(color: Color(0xFF94A3B8), fontSize: 14), // Slate 400
      ),
    );
  }
}
