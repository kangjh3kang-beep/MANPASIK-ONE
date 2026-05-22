/**
 * @mmup-axis 2 가정·직장 종합검진
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard WCAG 2.2 AA
 *
 * MMUP Design System — comprehensive design tokens for medical POCT UI
 */

// ---------------------------------------------------------------------------
// Color palette — medical POCT context (Tailwind sky/slate/purple/emerald/red)
// ---------------------------------------------------------------------------
export const colors = {
  primary: {
    50: "#f0f9ff",
    100: "#e0f2fe",
    200: "#bae6fd",
    300: "#7dd3fc",
    400: "#38bdf8",
    500: "#0ea5e9",
    600: "#0284c7",
    700: "#0369a1",
    800: "#075985",
    900: "#0c4a6e",
  },
  secondary: {
    50: "#faf5ff",
    100: "#f3e8ff",
    200: "#e9d5ff",
    300: "#d8b4fe",
    400: "#c084fc",
    500: "#a855f7",
    600: "#9333ea",
    700: "#7e22ce",
    800: "#6b21a8",
    900: "#581c87",
  },
  neutral: {
    50: "#f8fafc",
    100: "#f1f5f9",
    200: "#e2e8f0",
    300: "#cbd5e1",
    400: "#94a3b8",
    500: "#64748b",
    600: "#475569",
    700: "#334155",
    800: "#1e293b",
    900: "#0f172a",
  },
  success: {
    light: "#d1fae5",
    DEFAULT: "#10b981",
    dark: "#065f46",
  },
  warning: {
    light: "#fef3c7",
    DEFAULT: "#f59e0b",
    dark: "#92400e",
  },
  error: {
    light: "#fee2e2",
    DEFAULT: "#ef4444",
    dark: "#991b1b",
  },
  info: {
    light: "#dbeafe",
    DEFAULT: "#3b82f6",
    dark: "#1e40af",
  },
} as const;

// ---------------------------------------------------------------------------
// Spacing
// ---------------------------------------------------------------------------
export const spacing = {
  0: "0",
  1: "4px",
  2: "8px",
  3: "12px",
  4: "16px",
  5: "20px",
  6: "24px",
  8: "32px",
  10: "40px",
  12: "48px",
  16: "64px",
} as const;

// ---------------------------------------------------------------------------
// Typography
// ---------------------------------------------------------------------------
export const typography = {
  fontFamily: {
    sans: "Inter, system-ui, sans-serif",
    mono: "JetBrains Mono, monospace",
  },
  fontSize: {
    xs: "0.75rem",
    sm: "0.875rem",
    base: "1rem",
    lg: "1.125rem",
    xl: "1.25rem",
    "2xl": "1.5rem",
    "3xl": "1.875rem",
    "4xl": "2.25rem",
  },
  fontWeight: {
    normal: "400",
    medium: "500",
    semibold: "600",
    bold: "700",
  },
  lineHeight: {
    tight: "1.25",
    normal: "1.5",
    relaxed: "1.75",
  },
} as const;

// ---------------------------------------------------------------------------
// Border radius
// ---------------------------------------------------------------------------
export const borderRadius = {
  none: "0",
  sm: "4px",
  md: "8px",
  lg: "12px",
  xl: "16px",
  "2xl": "24px",
  full: "9999px",
} as const;

// ---------------------------------------------------------------------------
// Shadows
// ---------------------------------------------------------------------------
export const shadows = {
  sm: "0 1px 2px rgba(0,0,0,0.05)",
  md: "0 4px 6px rgba(0,0,0,0.07)",
  lg: "0 10px 15px rgba(0,0,0,0.1)",
  xl: "0 20px 25px rgba(0,0,0,0.1)",
} as const;

// ---------------------------------------------------------------------------
// Breakpoints
// ---------------------------------------------------------------------------
export const breakpoints = {
  sm: "640px",
  md: "768px",
  lg: "1024px",
  xl: "1280px",
  "2xl": "1536px",
} as const;

// ---------------------------------------------------------------------------
// Z-index
// ---------------------------------------------------------------------------
export const zIndex = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modal: 1040,
  popover: 1050,
  tooltip: 1060,
} as const;

// ---------------------------------------------------------------------------
// Component tokens
// ---------------------------------------------------------------------------
export const components = {
  button: {
    height: { sm: "32px", md: "40px", lg: "48px" },
    borderRadius: "8px",
  },
  input: {
    height: { sm: "32px", md: "40px", lg: "48px" },
    borderRadius: "8px",
  },
  card: {
    borderRadius: "16px",
    padding: "24px",
  },
  badge: {
    borderRadius: "9999px",
    paddingX: "8px",
    paddingY: "2px",
  },
} as const;
