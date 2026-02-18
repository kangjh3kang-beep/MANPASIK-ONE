# MANPASIK World Code-Based Style Guide

**문서번호**: MPK-STYLE-CODE-v3.0
**작성일**: 2026-02-14
**적용대상**: Web, App, OS Developers

---

## 1. Web Module (Next.js + Tailwind CSS v4)

### 1.1 CSS Variables
```css
:root {
  /* Primary */
  --background: #0A192F;
  --foreground: #FAFAFA;
  --deep-sea: #0A192F;
  --ink-black: #000000;
  --sanggam-gold: #D4AF37;
  --wave-cyan: #64FFDA;
  --celadon-teal: #00897B;
  --dancheong-red: #D32F2F;
  --hanji-white: #FAFAFA;
  --glass-navy: rgba(26, 35, 126, 0.6);

  /* Semantic */
  --success: #4CAF50;
  --warning: #FFC107;
  --info: #2196F3;

  /* Spacing (4px base unit) */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-6: 24px;
  --space-8: 32px;
  --space-12: 48px;
  --space-16: 64px;

  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-2xl: 24px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.3);
  --shadow-sanggam: 0 0 20px rgba(212, 175, 55, 0.3);

  /* Transitions */
  --transition-micro: 150ms ease-out;
  --transition-small: 250ms ease-in-out;
  --transition-medium: 400ms cubic-bezier(0.4, 0, 0.2, 1);
  --transition-large: 600ms cubic-bezier(0.0, 0, 0.2, 1);
}
```

### 1.2 Responsive Breakpoints
```css
/* Tailwind v4 breakpoints */
--breakpoint-sm: 640px;   /* Mobile landscape */
--breakpoint-md: 768px;   /* Tablet portrait */
--breakpoint-lg: 1024px;  /* Tablet landscape / Small desktop */
--breakpoint-xl: 1280px;  /* Desktop */
--breakpoint-2xl: 1536px; /* Wide desktop */
```

### 1.3 Grid System
```css
/* 12-column grid */
.grid-layout {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: var(--space-6);
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 var(--space-4);
}

/* Mobile: single column */
@media (max-width: 640px) {
  .grid-layout { grid-template-columns: 1fr; gap: var(--space-4); }
}
/* Tablet: 2 columns */
@media (min-width: 641px) and (max-width: 1023px) {
  .grid-layout { grid-template-columns: repeat(6, 1fr); }
}
```

### 1.4 Animation Keyframes
```css
@keyframes wave-animate {
  0%, 100% { transform: scale(1); opacity: 0.9; }
  50% { transform: scale(1.02); opacity: 1; }
}

@keyframes wave-ripple {
  0% { box-shadow: 0 0 0 0 rgba(100, 255, 218, 0.4); }
  70% { box-shadow: 0 0 0 20px rgba(100, 255, 218, 0); }
  100% { box-shadow: 0 0 0 0 rgba(100, 255, 218, 0); }
}

@keyframes breathing {
  0%, 100% { transform: scale(1); opacity: 0.85; }
  50% { transform: scale(1.015); opacity: 1; }
}

@keyframes sanggam-glow {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

@keyframes data-flow {
  0% { transform: translateY(0) scale(1); opacity: 0; }
  50% { opacity: 1; }
  100% { transform: translateY(-100px) scale(0.5); opacity: 0; }
}

.wave-animate { animation: wave-animate 3s ease-in-out infinite; }
.wave-ripple-animate { animation: wave-ripple 2s ease-out infinite; }
.breathing-animate { animation: breathing 3s ease-in-out infinite; }
.sanggam-glow-line {
  background: linear-gradient(90deg, transparent, var(--sanggam-gold), transparent);
  background-size: 200% 100%;
  animation: sanggam-glow 3s linear infinite;
}
```

### 1.5 Dark Mode Overrides
```css
[data-theme="light"] {
  --background: #FAFAFA;
  --foreground: #0A192F;
  --glass-navy: rgba(200, 210, 230, 0.6);
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.1);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```

### 1.6 Accessibility: Focus Indicators
```css
*:focus-visible {
  outline: 2px solid var(--wave-cyan);
  outline-offset: 2px;
  border-radius: var(--radius-sm);
}

/* High contrast mode */
@media (prefers-contrast: high) {
  :root {
    --sanggam-gold: #FFD700;
    --wave-cyan: #00FFFF;
    --dancheong-red: #FF0000;
  }
}
```

---

## 2. App Module (Flutter)

### 2.1 Theme Configuration
```dart
// lib/core/theme/app_theme.dart
class AppTheme {
  // Primary Colors
  static const sanggamGold = Color(0xFFD4AF37);
  static const deepSeaBlue = Color(0xFF0A192F);
  static const glassBlue = Color(0x1A64FFDA);
  static const waveCyan = Color(0xFF00E5FF);  // App-optimized cyan
  static const inkBlack = Color(0xFF020617);
  static const hanjiWhite = Color(0xFFF8FAFC);
  static const dancheongRed = Color(0xFFFF4D4D);

  // Semantic Colors
  static const success = Color(0xFF4CAF50);
  static const warning = Color(0xFFFFC107);
  static const info = Color(0xFF2196F3);
}
```

### 2.2 Spacing System
```dart
// lib/core/theme/spacing.dart
class AppSpacing {
  static const double xs = 4.0;    // space-1
  static const double sm = 8.0;    // space-2
  static const double md = 12.0;   // space-3
  static const double base = 16.0; // space-4
  static const double lg = 24.0;   // space-6
  static const double xl = 32.0;   // space-8
  static const double xxl = 48.0;  // space-12
  static const double xxxl = 64.0; // space-16
}
```

### 2.3 Border Radius
```dart
class AppRadius {
  static const double sm = 4.0;
  static const double md = 8.0;
  static const double lg = 12.0;
  static const double xl = 16.0;
  static const double xxl = 24.0;
  static final BorderRadius cardRadius = BorderRadius.circular(lg);
  static final BorderRadius buttonRadius = BorderRadius.circular(md);
  static final BorderRadius chipRadius = BorderRadius.circular(xxl);
}
```

### 2.4 Shadow/Elevation System
```dart
class AppShadows {
  static final List<BoxShadow> sm = [
    BoxShadow(color: Colors.black.withOpacity(0.3), blurRadius: 2, offset: Offset(0, 1)),
  ];
  static final List<BoxShadow> md = [
    BoxShadow(color: Colors.black.withOpacity(0.3), blurRadius: 6, offset: Offset(0, 4)),
  ];
  static final List<BoxShadow> lg = [
    BoxShadow(color: Colors.black.withOpacity(0.3), blurRadius: 15, offset: Offset(0, 10)),
  ];
  static final List<BoxShadow> sanggamGlow = [
    BoxShadow(color: AppTheme.sanggamGold.withOpacity(0.3), blurRadius: 20),
  ];
}
```

### 2.5 Animation Widgets
```dart
// lib/shared/widgets/wave_ripple_painter.dart
class WaveRipplePainter extends CustomPainter {
  final double animationValue; // 0.0 ~ 1.0
  final Color color;

  // Draws concentric circles expanding outward
  // Used for: measurement button tap effect
  // Duration: 2000ms, Curve: Curves.easeOut
}

class WaveRippleBackground extends StatefulWidget {
  // Full-screen wave background with multiple ripple origins
  // Duration: 3000ms per cycle, infinite loop
}

class WavePainter extends CustomPainter {
  // Sine wave animation for calm/loading states
  // Duration: 2000ms, Curve: Curves.linear
}

// lib/shared/widgets/breathing_overlay.dart
class BreathingOverlay extends StatefulWidget {
  final Widget child;
  final Duration duration; // default: 3000ms
  final double minScale;  // default: 1.0
  final double maxScale;  // default: 1.015

  // Scale + opacity breathing wrapper
  // Curve: Curves.easeInOut
}

// lib/shared/widgets/sanggam_decoration.dart
class SanggamDecoration extends BoxDecoration {
  // Gold gradient border + glass background
  // Border: 1px linear-gradient(Sanggam Gold → Wave Cyan)
  // Background: Glass Navy with 60% opacity
  // Shadow: sanggamGlow
}
```

### 2.6 Touch Targets & Accessibility
```dart
// Minimum touch target: 48x48dp (Material Design 3)
class AppTouchTarget {
  static const double minimum = 48.0;
  static const double comfortable = 56.0;
  static const double emergency = 64.0; // 119, delete buttons
}

// Senior mode scaling
class SeniorMode {
  static const double fontScale = 1.5;
  static const double touchScale = 1.3;
  static const double spacingScale = 1.2;
}
```

### 2.7 Dark Mode Color Mapping
```dart
// lib/core/theme/dark_theme.dart
final darkTheme = ThemeData(
  brightness: Brightness.dark,
  scaffoldBackgroundColor: AppTheme.deepSeaBlue,
  colorScheme: ColorScheme.dark(
    primary: AppTheme.sanggamGold,
    secondary: AppTheme.waveCyan,
    surface: Color(0xFF112240),
    error: AppTheme.dancheongRed,
    onPrimary: AppTheme.deepSeaBlue,
    onSecondary: AppTheme.deepSeaBlue,
    onSurface: AppTheme.hanjiWhite,
  ),
);

// lib/core/theme/light_theme.dart
final lightTheme = ThemeData(
  brightness: Brightness.light,
  scaffoldBackgroundColor: AppTheme.hanjiWhite,
  colorScheme: ColorScheme.light(
    primary: AppTheme.deepSeaBlue,
    secondary: AppTheme.sanggamGold,
    surface: Colors.white,
    error: AppTheme.dancheongRed,
    onPrimary: AppTheme.hanjiWhite,
    onSecondary: AppTheme.deepSeaBlue,
    onSurface: AppTheme.inkBlack,
  ),
);
```

---

## 3. OS Module (Rust Slint)

### 3.1 Global Style
```slint
export global Style {
    // Colors
    property <color> deep-sea: #0A192F;
    property <color> ink-black: #000000;
    property <color> sanggam-gold: #D4AF37;
    property <color> wave-cyan: #64FFDA;
    property <color> hanji-white: #FAFAFA;
    property <color> status-ok: #4CAF50;
    property <color> status-warn: #FFC107;
    property <color> status-err: #D32F2F;

    // Typography (160x80 OLED optimized)
    property <length> text-size-main: 12px;
    property <length> text-size-small: 10px;
    property <length> text-size-tiny: 8px;

    // Spacing
    property <length> space-xs: 2px;
    property <length> space-sm: 4px;
    property <length> space-md: 8px;
    property <length> space-lg: 12px;

    // Effects
    property <image-filter> glass-blur: blur(5px);
}
```

### 3.2 Component Examples
```slint
// Status bar component for OLED display
component StatusBar inherits Rectangle {
    height: 12px;
    background: Style.deep-sea;

    HorizontalLayout {
        spacing: Style.space-sm;

        Text {
            text: "●";
            color: Style.status-ok;
            font-size: Style.text-size-tiny;
        }
        Text {
            text: "BLE";
            color: Style.hanji-white;
            font-size: Style.text-size-tiny;
        }
        Rectangle { /* spacer */ }
        Text {
            text: "85%";
            color: Style.hanji-white;
            font-size: Style.text-size-tiny;
        }
    }
}

// Measurement progress indicator
component MeasureProgress inherits Rectangle {
    in property <float> progress: 0.0;
    height: 2px;
    background: Style.deep-sea;

    Rectangle {
        width: parent.width * root.progress;
        height: parent.height;
        background: Style.sanggam-gold;

        animate width {
            duration: 100ms;
            easing: ease-in-out;
        }
    }
}
```

---

## 4. 플랫폼별 Wave Cyan 차이 설명

> **참고**: Wave Cyan은 플랫폼별 최적화된 값을 사용합니다.
> - Web: `#64FFDA` (Material Teal A200) — 대형 디스플레이 최적화
> - App: `#00E5FF` (Material Cyan A400) — 모바일 디스플레이 채도 보정
> - OS: `#64FFDA` — OLED 저전력 최적화
>
> 이는 각 플랫폼의 디스플레이 특성과 접근성 기준에 맞춘 의도적 차이입니다.

---

**참조**: `docs/design/ECOSYSTEM_DESIGN_MASTER_PLAN.md`, `docs/design/BRAND_GUIDELINE.md`, `docs/DESIGN_SYSTEM.md`
