import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

/// Sanggam Orbit 디자인 철학을 적용한 전역 테마.
///
/// 한국 전통 상감(象嵌) 기법에서 영감을 받아,
/// 깊은 바다빛 배경 위에 금(Gold)과 자개(Jagae) 빛이
/// 궤도(Orbit)를 그리며 빛나는 디자인 언어를 구현한다.
class SanggamTheme {
  SanggamTheme._();

  // ──────────────────────────────────────────────
  //  1. Color Palette
  // ──────────────────────────────────────────────

  /// 배경색 — 심해(深海) 블루
  static const Color background = Color(0xFF0B1021);

  /// 주 색상 — 상감 금(金)
  static const Color primary = Color(0xFFD4AF37);

  /// 자개 시안 — 자개 빛의 한 축
  static const Color jagaeCyan = Color(0xFF00FFFF);

  /// 자개 마젠타 — 자개 빛의 한 축
  static const Color jagaeMagenta = Color(0xFFFF00FF);

  /// 어두운 표면 — 카드/다이얼로그 배경
  static const Color surface = Color(0xFF111A30);

  /// 보조 표면 — 입력 필드, 약간 밝은 패널
  static const Color surfaceVariant = Color(0xFF162040);

  /// 오류 색상 — 단청 붉은빛
  static const Color error = Color(0xFFFF4D4D);

  /// 비활성 텍스트/아이콘
  static const Color onSurfaceDim = Color(0xFF8899BB);

  /// 상감 금 위의 텍스트 (대비용)
  static const Color onPrimary = Color(0xFF0B1021);

  // ──────────────────────────────────────────────
  //  2. Jagae(자개) Accent Gradient
  // ──────────────────────────────────────────────

  /// 자개 오로라 그래디언트 — 수평 방향.
  /// 시안(#00FFFF) → 마젠타(#FF00FF) 혼합으로
  /// 전통 자개의 무지개빛 광택을 표현한다.
  static const LinearGradient jagaeGradient = LinearGradient(
    colors: [jagaeCyan, jagaeMagenta],
    begin: Alignment.centerLeft,
    end: Alignment.centerRight,
  );

  /// 자개 오로라 그래디언트 — 대각선 방향.
  /// 보다 역동적인 자개 빛 표현에 사용한다.
  static const LinearGradient jagaeGradientDiagonal = LinearGradient(
    colors: [jagaeCyan, Color(0xFF8844FF), jagaeMagenta],
    stops: [0.0, 0.5, 1.0],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  /// 자개 그래디언트 — 수직 방향 (Shimmer/Glow 용).
  static const LinearGradient jagaeGradientVertical = LinearGradient(
    colors: [jagaeCyan, jagaeMagenta],
    begin: Alignment.topCenter,
    end: Alignment.bottomCenter,
  );

  /// 상감 금 그래디언트 — 버튼/배지 하이라이트.
  static const LinearGradient goldGradient = LinearGradient(
    colors: [Color(0xFFE8C84A), primary, Color(0xFFC49B26)],
    stops: [0.0, 0.5, 1.0],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );

  // ──────────────────────────────────────────────
  //  3. Typography
  // ──────────────────────────────────────────────

  /// 영문 기본 폰트: Outfit (Geometric Sans-Serif).
  static String? get _outfitFamily => GoogleFonts.outfit().fontFamily;

  /// 한글 기본 폰트: Noto Sans KR.
  static String? get _notoSansKrFamily => GoogleFonts.notoSansKr().fontFamily;

  /// 영문 TextStyle 헬퍼 — Outfit 기반.
  static TextStyle _outfit({
    double fontSize = 14,
    FontWeight fontWeight = FontWeight.w400,
    Color color = Colors.white,
    double letterSpacing = 0,
    double? height,
  }) {
    return GoogleFonts.outfit(
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color,
      letterSpacing: letterSpacing,
      height: height,
    );
  }

  /// 한글 TextStyle 헬퍼 — Noto Sans KR 기반.
  static TextStyle _notoSansKr({
    double fontSize = 14,
    FontWeight fontWeight = FontWeight.w400,
    Color color = Colors.white,
    double letterSpacing = 0,
    double? height,
  }) {
    return GoogleFonts.notoSansKr(
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: color,
      letterSpacing: letterSpacing,
      height: height,
    );
  }

  /// 전체 TextTheme 구성.
  ///
  /// Display/Headline: Outfit + 상감 금(Gold) 색상
  /// Title/Body/Label: Noto Sans KR + 화이트/그레이 계열
  static TextTheme get _textTheme => TextTheme(
        // ── Display (가장 큰 제목) ──
        displayLarge: _outfit(
          fontSize: 57,
          fontWeight: FontWeight.w700,
          color: primary,
          letterSpacing: -1.5,
          height: 1.12,
        ),
        displayMedium: _outfit(
          fontSize: 45,
          fontWeight: FontWeight.w600,
          color: primary,
          letterSpacing: -0.5,
          height: 1.16,
        ),
        displaySmall: _outfit(
          fontSize: 36,
          fontWeight: FontWeight.w600,
          color: primary,
          height: 1.22,
        ),

        // ── Headline (섹션 제목) ──
        headlineLarge: _outfit(
          fontSize: 32,
          fontWeight: FontWeight.w700,
          color: primary,
          height: 1.25,
        ),
        headlineMedium: _outfit(
          fontSize: 28,
          fontWeight: FontWeight.w600,
          color: primary,
          height: 1.29,
        ),
        headlineSmall: _outfit(
          fontSize: 24,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          height: 1.33,
        ),

        // ── Title (카드/리스트 제목) ──
        titleLarge: _notoSansKr(
          fontSize: 22,
          fontWeight: FontWeight.w700,
          color: Colors.white,
          height: 1.27,
        ),
        titleMedium: _notoSansKr(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          letterSpacing: 0.15,
          height: 1.50,
        ),
        titleSmall: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          letterSpacing: 0.1,
          height: 1.43,
        ),

        // ── Body (본문) ──
        bodyLarge: _notoSansKr(
          fontSize: 16,
          fontWeight: FontWeight.w400,
          color: Colors.white,
          letterSpacing: 0.5,
          height: 1.50,
        ),
        bodyMedium: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w400,
          color: const Color(0xFFCCDDEE),
          letterSpacing: 0.25,
          height: 1.43,
        ),
        bodySmall: _notoSansKr(
          fontSize: 12,
          fontWeight: FontWeight.w400,
          color: onSurfaceDim,
          letterSpacing: 0.4,
          height: 1.33,
        ),

        // ── Label (버튼/탭/칩) ──
        labelLarge: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          letterSpacing: 0.1,
          height: 1.43,
        ),
        labelMedium: _notoSansKr(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          letterSpacing: 0.5,
          height: 1.33,
        ),
        labelSmall: _notoSansKr(
          fontSize: 11,
          fontWeight: FontWeight.w500,
          color: onSurfaceDim,
          letterSpacing: 0.5,
          height: 1.45,
        ),
      );

  // ──────────────────────────────────────────────
  //  4. getTheme() — 완성된 ThemeData
  // ──────────────────────────────────────────────

  /// Sanggam Orbit 다크 테마를 반환한다.
  ///
  /// 앱 전역에서 사용하며, 다음을 포함한다:
  /// - Material 3 기반 ColorScheme (Gold primary, DeepSea background)
  /// - Scaffold / AppBar / Card / Dialog / BottomSheet / NavigationBar 테마
  /// - ElevatedButton (다크 배경 + 골드 텍스트)
  /// - Input / Chip / Divider 등 전역 위젯 테마
  static ThemeData getTheme() {
    final colorScheme = ColorScheme(
      brightness: Brightness.dark,
      primary: primary,
      onPrimary: onPrimary,
      primaryContainer: const Color(0xFF2A2210),
      onPrimaryContainer: primary,
      secondary: jagaeCyan,
      onSecondary: background,
      secondaryContainer: const Color(0xFF0A2A2A),
      onSecondaryContainer: jagaeCyan,
      tertiary: jagaeMagenta,
      onTertiary: background,
      tertiaryContainer: const Color(0xFF2A0A2A),
      onTertiaryContainer: jagaeMagenta,
      error: error,
      onError: Colors.white,
      errorContainer: const Color(0xFF3B1010),
      onErrorContainer: const Color(0xFFFFB4AB),
      surface: surface,
      onSurface: Colors.white,
      surfaceContainerHighest: surfaceVariant,
      outline: const Color(0xFF334466),
      outlineVariant: const Color(0xFF1E2E4A),
      shadow: Colors.black,
      scrim: Colors.black,
      inverseSurface: const Color(0xFFE2E2E6),
      onInverseSurface: const Color(0xFF1A1C1E),
      inversePrimary: const Color(0xFF8A7520),
    );

    final textTheme = _textTheme;

    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: colorScheme,
      fontFamily: _notoSansKrFamily,

      // ── Scaffold ──
      scaffoldBackgroundColor: background,

      // ── AppBar ──
      appBarTheme: AppBarTheme(
        backgroundColor: background,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        scrolledUnderElevation: 0.5,
        shadowColor: primary.withOpacity(0.15),
        centerTitle: true,
        iconTheme: const IconThemeData(color: Colors.white, size: 22),
        actionsIconTheme: IconThemeData(color: primary, size: 22),
        titleTextStyle: _outfit(
          fontSize: 20,
          fontWeight: FontWeight.w600,
          color: Colors.white,
          letterSpacing: 0.5,
        ),
        toolbarTextStyle: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w500,
          color: Colors.white,
        ),
      ),

      // ── TextTheme ──
      textTheme: textTheme,

      // ── ElevatedButton: 다크 배경 + 골드 텍스트 ──
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ButtonStyle(
          backgroundColor: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) {
              return surfaceVariant;
            }
            return surface;
          }),
          foregroundColor: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) {
              return onSurfaceDim;
            }
            return primary;
          }),
          overlayColor: WidgetStateProperty.all(primary.withOpacity(0.08)),
          elevation: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.pressed)) return 0;
            if (states.contains(WidgetState.hovered)) return 6;
            return 3;
          }),
          shadowColor: WidgetStateProperty.all(primary.withOpacity(0.25)),
          shape: WidgetStateProperty.all(
            RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
          ),
          side: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) {
              return BorderSide.none;
            }
            return BorderSide(color: primary.withOpacity(0.35), width: 1);
          }),
          padding: WidgetStateProperty.all(
            const EdgeInsets.symmetric(horizontal: 28, vertical: 16),
          ),
          textStyle: WidgetStateProperty.all(
            _notoSansKr(
              fontSize: 15,
              fontWeight: FontWeight.w700,
              color: primary,
            ),
          ),
          minimumSize: WidgetStateProperty.all(const Size(88, 52)),
          animationDuration: const Duration(milliseconds: 200),
        ),
      ),

      // ── OutlinedButton ──
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: ButtonStyle(
          foregroundColor: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) return onSurfaceDim;
            return primary;
          }),
          side: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) {
              return BorderSide(color: onSurfaceDim.withOpacity(0.3));
            }
            return BorderSide(color: primary.withOpacity(0.5));
          }),
          shape: WidgetStateProperty.all(
            RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
          ),
          padding: WidgetStateProperty.all(
            const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
          ),
          textStyle: WidgetStateProperty.all(
            _notoSansKr(fontSize: 14, fontWeight: FontWeight.w600),
          ),
          overlayColor: WidgetStateProperty.all(primary.withOpacity(0.06)),
        ),
      ),

      // ── TextButton ──
      textButtonTheme: TextButtonThemeData(
        style: ButtonStyle(
          foregroundColor: WidgetStateProperty.resolveWith((states) {
            if (states.contains(WidgetState.disabled)) return onSurfaceDim;
            return primary;
          }),
          textStyle: WidgetStateProperty.all(
            _notoSansKr(fontSize: 14, fontWeight: FontWeight.w600),
          ),
          overlayColor: WidgetStateProperty.all(primary.withOpacity(0.06)),
          shape: WidgetStateProperty.all(
            RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
          ),
        ),
      ),

      // ── FloatingActionButton ──
      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: primary,
        foregroundColor: onPrimary,
        elevation: 6,
        focusElevation: 8,
        hoverElevation: 10,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
        extendedTextStyle: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w700,
          color: onPrimary,
        ),
      ),

      // ── IconButton ──
      iconButtonTheme: IconButtonThemeData(
        style: ButtonStyle(
          foregroundColor: WidgetStateProperty.all(Colors.white),
          overlayColor: WidgetStateProperty.all(primary.withOpacity(0.08)),
        ),
      ),

      // ── Card ──
      cardTheme: CardThemeData(
        elevation: 0,
        color: surface,
        surfaceTintColor: Colors.transparent,
        shadowColor: Colors.black45,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
          side: BorderSide(color: primary.withOpacity(0.12), width: 0.5),
        ),
        margin: const EdgeInsets.symmetric(horizontal: 4, vertical: 4),
      ),

      // ── Dialog ──
      dialogTheme: DialogThemeData(
        backgroundColor: surface,
        surfaceTintColor: Colors.transparent,
        elevation: 16,
        shadowColor: Colors.black54,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
          side: BorderSide(color: primary.withOpacity(0.15)),
        ),
        titleTextStyle: _outfit(
          fontSize: 22,
          fontWeight: FontWeight.w700,
          color: primary,
        ),
        contentTextStyle: _notoSansKr(
          fontSize: 14,
          color: const Color(0xFFCCDDEE),
          height: 1.5,
        ),
      ),

      // ── BottomSheet ──
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: surface,
        surfaceTintColor: Colors.transparent,
        elevation: 12,
        shadowColor: Colors.black54,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        dragHandleColor: primary.withOpacity(0.4),
        dragHandleSize: const Size(36, 4),
        modalBarrierColor: Colors.black54,
      ),

      // ── InputDecoration ──
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: surfaceVariant,
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: BorderSide(
            color: const Color(0xFF334466).withOpacity(0.5),
            width: 0.5,
          ),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: primary, width: 1.5),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: error, width: 1),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: error, width: 1.5),
        ),
        labelStyle: _notoSansKr(fontSize: 14, color: onSurfaceDim),
        hintStyle: _notoSansKr(
          fontSize: 14,
          color: onSurfaceDim.withOpacity(0.6),
        ),
        errorStyle: _notoSansKr(fontSize: 12, color: error),
        prefixIconColor: onSurfaceDim,
        suffixIconColor: onSurfaceDim,
        floatingLabelStyle: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w500,
          color: primary,
        ),
      ),

      // ── Chip ──
      chipTheme: ChipThemeData(
        backgroundColor: surfaceVariant,
        selectedColor: primary.withOpacity(0.15),
        disabledColor: surfaceVariant.withOpacity(0.5),
        labelStyle: _notoSansKr(
          fontSize: 13,
          fontWeight: FontWeight.w500,
          color: Colors.white,
        ),
        secondaryLabelStyle: _notoSansKr(fontSize: 13, color: primary),
        side: BorderSide(color: primary.withOpacity(0.2)),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
        showCheckmark: true,
        checkmarkColor: primary,
      ),

      // ── NavigationBar (Bottom) ──
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: background,
        surfaceTintColor: Colors.transparent,
        indicatorColor: primary.withOpacity(0.12),
        elevation: 0,
        height: 68,
        labelTextStyle: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return _notoSansKr(
              fontSize: 12,
              fontWeight: FontWeight.w700,
              color: primary,
            );
          }
          return _notoSansKr(
            fontSize: 12,
            fontWeight: FontWeight.w400,
            color: onSurfaceDim,
          );
        }),
        iconTheme: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return IconThemeData(color: primary, size: 24);
          }
          return IconThemeData(color: onSurfaceDim, size: 24);
        }),
      ),

      // ── NavigationRail (사이드) ──
      navigationRailTheme: NavigationRailThemeData(
        backgroundColor: background,
        indicatorColor: primary.withOpacity(0.12),
        selectedIconTheme: IconThemeData(color: primary, size: 24),
        unselectedIconTheme: IconThemeData(color: onSurfaceDim, size: 22),
        selectedLabelTextStyle: _notoSansKr(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          color: primary,
        ),
        unselectedLabelTextStyle: _notoSansKr(
          fontSize: 12,
          color: onSurfaceDim,
        ),
      ),

      // ── TabBar ──
      tabBarTheme: TabBarThemeData(
        labelColor: primary,
        unselectedLabelColor: onSurfaceDim,
        indicatorColor: primary,
        indicatorSize: TabBarIndicatorSize.label,
        labelStyle: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w700,
        ),
        unselectedLabelStyle: _notoSansKr(
          fontSize: 14,
          fontWeight: FontWeight.w400,
        ),
        dividerColor: Colors.transparent,
      ),

      // ── Divider ──
      dividerTheme: DividerThemeData(
        color: const Color(0xFF1E2E4A),
        thickness: 0.5,
        space: 1,
      ),

      // ── SnackBar ──
      snackBarTheme: SnackBarThemeData(
        backgroundColor: surface,
        contentTextStyle: _notoSansKr(fontSize: 14, color: Colors.white),
        actionTextColor: primary,
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: BorderSide(color: primary.withOpacity(0.2)),
        ),
        elevation: 6,
      ),

      // ── Tooltip ──
      tooltipTheme: TooltipThemeData(
        decoration: BoxDecoration(
          color: surfaceVariant,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: primary.withOpacity(0.2)),
        ),
        textStyle: _notoSansKr(fontSize: 12, color: Colors.white),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        waitDuration: const Duration(milliseconds: 500),
      ),

      // ── ListTile ──
      listTileTheme: ListTileThemeData(
        tileColor: Colors.transparent,
        selectedTileColor: primary.withOpacity(0.06),
        iconColor: onSurfaceDim,
        selectedColor: primary,
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
        titleTextStyle: _notoSansKr(
          fontSize: 15,
          fontWeight: FontWeight.w500,
          color: Colors.white,
        ),
        subtitleTextStyle: _notoSansKr(
          fontSize: 13,
          color: onSurfaceDim,
        ),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),

      // ── Switch ──
      switchTheme: SwitchThemeData(
        thumbColor: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) return primary;
          return onSurfaceDim;
        }),
        trackColor: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return primary.withOpacity(0.3);
          }
          return surfaceVariant;
        }),
        trackOutlineColor: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return primary.withOpacity(0.5);
          }
          return const Color(0xFF334466);
        }),
      ),

      // ── Slider ──
      sliderTheme: SliderThemeData(
        activeTrackColor: primary,
        inactiveTrackColor: surfaceVariant,
        thumbColor: primary,
        overlayColor: primary.withOpacity(0.12),
        valueIndicatorColor: surface,
        valueIndicatorTextStyle: _notoSansKr(
          fontSize: 13,
          fontWeight: FontWeight.w600,
          color: primary,
        ),
      ),

      // ── ProgressIndicator ──
      progressIndicatorTheme: const ProgressIndicatorThemeData(
        color: primary,
        linearTrackColor: surfaceVariant,
        circularTrackColor: surfaceVariant,
      ),

      // ── PopupMenu ──
      popupMenuTheme: PopupMenuThemeData(
        color: surface,
        surfaceTintColor: Colors.transparent,
        elevation: 8,
        shadowColor: Colors.black54,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(14),
          side: BorderSide(color: primary.withOpacity(0.12)),
        ),
        textStyle: _notoSansKr(fontSize: 14, color: Colors.white),
      ),

      // ── Drawer ──
      drawerTheme: DrawerThemeData(
        backgroundColor: background,
        surfaceTintColor: Colors.transparent,
        elevation: 16,
        shadowColor: Colors.black54,
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.horizontal(right: Radius.circular(24)),
        ),
      ),

      // ── ScrollBar ──
      scrollbarTheme: ScrollbarThemeData(
        thumbColor: WidgetStateProperty.all(primary.withOpacity(0.3)),
        radius: const Radius.circular(4),
        thickness: WidgetStateProperty.all(4),
      ),
    );
  }
}
