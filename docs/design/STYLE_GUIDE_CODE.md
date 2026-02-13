# MANPASIK World Code-Based Style Guide

**문서번호**: MPK-STYLE-CODE-v2.0
**작성일**: 2026-02-14
**적용대상**: Web, App, OS Developers

---

## 1. Web Module (Next.js + Tailwind CSS v4)

### `globals.css` — CSS Variables
```css
:root {
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
}
```

### Animation Classes
```css
.wave-animate      /* scale(1↔1.02) + opacity pulse */
.wave-ripple-animate /* 동심원 box-shadow 확산 */
.breathing-animate   /* scale(1↔1.015) + opacity(0.85↔1) */
.sanggam-glow-line   /* 금색 그라데이션 라인 이동 */
```

## 2. App Module (Flutter)

### `lib/core/theme/app_theme.dart`
```dart
class AppTheme {
  static const sanggamGold = Color(0xFFD4AF37);
  static const deepSeaBlue = Color(0xFF0A192F);
  static const glassBlue = Color(0x1A64FFDA);
  static const waveCyan = Color(0xFF00E5FF);
  static const inkBlack = Color(0xFF020617);
  static const hanjiWhite = Color(0xFFF8FAFC);
  static const dancheongRed = Color(0xFFFF4D4D);
}
```

### Animation Widgets
```dart
// lib/shared/widgets/wave_ripple_painter.dart
WaveRipplePainter   // 동심원 파동 CustomPainter
WaveRippleBackground // 파동 배경 위젯
WavePainter          // 사인파 안정화 CustomPainter

// lib/shared/widgets/breathing_overlay.dart
BreathingOverlay     // scale+opacity 호흡 래퍼
```

## 3. OS Module (Rust Slint)

### `ui/style.slint`
```slint
export global Style {
    property <color> deep-sea: #0A192F;
    property <color> ink-black: #000000;
    property <color> sanggam-gold: #D4AF37;
    property <color> wave-cyan: #64FFDA;
    property <color> status-ok: #00FF00;
    property <color> status-warn: #FFFF00;
    property <color> status-err: #FF0000;
    
    property <length> text-size-main: 12px;
    property <length> text-size-small: 10px;
    
    property <image-filter> glass-blur: blur(5px);
}
```
