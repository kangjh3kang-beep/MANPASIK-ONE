import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  // MANPASIK R&D Lab - Premium Color Palette (Korean Futuristic)
  static const sanggamGold = Color(0xFFD4AF37); // Traditional Gold Inlay
  static const deepSeaBlue = Color(0xFF050B14); // Deep Space Background
  static const glassBlue = Color(0x1A64FFDA); // Glassmorphism Tint
  static const waveCyan = Color(0xFF00E5FF); // Analysis Energy
  static const inkBlack = Color(0xFF020617); // Extra Dark
  static const hanjiWhite = Color(0xFFF8FAFC);
  static const dancheongRed = Color(0xFFFF4D4D); // Critical Alerts
  static const celadonTeal = Color(0xFF00897B); // Primary Brand Color
  
  // Light Theme (Clean Professional)
  static final ThemeData light = ThemeData(
    useMaterial3: true,
    brightness: Brightness.light,
    fontFamily: GoogleFonts.notoSansKr().fontFamily,
    
    colorScheme: ColorScheme.fromSeed(
      seedColor: deepSeaBlue,
      primary: deepSeaBlue,
      secondary: sanggamGold,
      error: dancheongRed,
      surface: hanjiWhite,
      onSurface: inkBlack,
      brightness: Brightness.light,
    ),
    
    scaffoldBackgroundColor: hanjiWhite,
    
    textTheme: TextTheme(
      headlineLarge: GoogleFonts.gowunBatang(fontWeight: FontWeight.bold, color: deepSeaBlue),
      headlineMedium: GoogleFonts.gowunBatang(fontWeight: FontWeight.bold, color: deepSeaBlue),
      titleLarge: GoogleFonts.notoSansKr(fontWeight: FontWeight.w700),
      bodyLarge: GoogleFonts.notoSansKr(),
    ),
    
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Colors.grey.shade50,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Colors.grey, width: 0.5),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: sanggamGold, width: 1.5),
      ),
      labelStyle: GoogleFonts.notoSansKr(color: inkBlack.withOpacity(0.6)),
    ),

    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: deepSeaBlue,
        foregroundColor: Colors.white,
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        padding: const EdgeInsets.symmetric(vertical: 16),
        textStyle: GoogleFonts.notoSansKr(fontSize: 16, fontWeight: FontWeight.bold),
      ),
    ),
    cardTheme: CardThemeData(
      elevation: 3,
      shadowColor: Colors.black12,
      surfaceTintColor: Colors.white,
      color: Colors.white,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
    ),
  );

  // Dark Theme (MANPASIK R&D Lab - Primary Mode)
  static final ThemeData dark = ThemeData(
    useMaterial3: true,
    brightness: Brightness.dark,
    fontFamily: GoogleFonts.notoSansKr().fontFamily,
    
    colorScheme: ColorScheme.fromSeed(
      seedColor: sanggamGold,
      primary: sanggamGold,
      secondary: waveCyan,
      surface: const Color(0xFF112240), // Darker Navy
      onSurface: Colors.white,
      brightness: Brightness.dark,
    ),
    
    scaffoldBackgroundColor: const Color(0xFF050B14), // Deep Space Black/Blue
    
    textTheme: TextTheme(
      headlineLarge: GoogleFonts.gowunBatang(
        fontWeight: FontWeight.bold, 
        color: sanggamGold,
        letterSpacing: -0.5,
      ),
      headlineMedium: GoogleFonts.gowunBatang(
        fontWeight: FontWeight.bold, 
        color: sanggamGold,
      ),
      titleLarge: GoogleFonts.notoSansKr(
        fontWeight: FontWeight.w700, 
        color: Colors.white,
      ),
      bodyLarge: GoogleFonts.notoSansKr(color: Colors.white),
    ),

    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: const Color(0xFF1D2D50),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide.none,
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: sanggamGold, width: 1.5),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      labelStyle: GoogleFonts.notoSansKr(color: Colors.white70),
      hintStyle: GoogleFonts.notoSansKr(color: Colors.white30),
    ),

    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: sanggamGold,
        foregroundColor: deepSeaBlue,
        elevation: 4,
        shadowColor: sanggamGold.withOpacity(0.3),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        padding: const EdgeInsets.symmetric(vertical: 16),
        textStyle: GoogleFonts.notoSansKr(fontSize: 16, fontWeight: FontWeight.bold),
      ),
    ),

    cardTheme: CardThemeData(
      elevation: 0,
      color: const Color(0xFF112240).withOpacity(0.8),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: BorderSide(color: sanggamGold.withOpacity(0.2), width: 0.5),
      ),
    ),
  );
  // Ultra-Premium Korean White Theme (Baekja/Hanji)
  static final ThemeData koreanWhite = ThemeData(
    useMaterial3: true,
    brightness: Brightness.light,
    fontFamily: GoogleFonts.gowunBatang().fontFamily, // Elegant Serif for Traditional Feel
    
    colorScheme: ColorScheme.fromSeed(
      seedColor: const Color(0xFF1A1A1A), // Ink Black
      primary: const Color(0xFF1A1A1A), // Ink Black for Logo/Headers
      secondary: sanggamGold, // Gold Accents
      surface: hanjiWhite, 
      onSurface: const Color(0xFF1A1A1A), // Clear Ink Black
      onPrimary: Colors.white,
      brightness: Brightness.light,
    ),
    
    scaffoldBackgroundColor: const Color(0xFFF9F9F7), // Warm Baekja White
    
    textTheme: TextTheme(
      headlineLarge: GoogleFonts.gowunBatang(
        fontWeight: FontWeight.bold, 
        color: const Color(0xFF1A1A1A), // High Contrast Ink Black
        letterSpacing: -0.5,
      ),
      headlineMedium: GoogleFonts.gowunBatang(
        fontWeight: FontWeight.bold, 
        color: const Color(0xFF1A1A1A),
      ),
      titleLarge: GoogleFonts.notoSansKr(
        fontWeight: FontWeight.w700, 
        color: const Color(0xFF1A1A1A),
      ),
      bodyLarge: GoogleFonts.notoSansKr(color: const Color(0xFF222222)), // Deep Grey/Black
      bodyMedium: GoogleFonts.notoSansKr(color: const Color(0xFF424242)), // Readable Grey
      labelSmall: GoogleFonts.notoSansKr(color: const Color(0xFF616161), fontWeight: FontWeight.w600), // Captions
    ),
    
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: const Color(0xFF1A1A1A).withOpacity(0.03), // Light Ink Wash
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: BorderSide(color: const Color(0xFF1A1A1A).withOpacity(0.15), width: 1.0),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: sanggamGold, width: 1.5), // Gold Focus
      ),
      labelStyle: GoogleFonts.notoSansKr(color: const Color(0xFF1A1A1A).withOpacity(0.7)),
    ),

    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: sanggamGold, // Gold Background
        foregroundColor: const Color(0xFF1A1A1A), // Ink Text
        elevation: 2,
        shadowColor: sanggamGold.withOpacity(0.4),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        padding: const EdgeInsets.symmetric(vertical: 16),
        textStyle: GoogleFonts.notoSansKr(fontSize: 16, fontWeight: FontWeight.bold),
      ),
    ),
    
    cardTheme: CardThemeData(
      elevation: 4,
      shadowColor: Colors.black.withOpacity(0.05), // Soft diffused shadow
      color: Colors.white, // Porcelain White
      surfaceTintColor: Colors.white,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: BorderSide(color: Colors.black.withOpacity(0.05), width: 0.5),
      ),
    ),
  );
}
