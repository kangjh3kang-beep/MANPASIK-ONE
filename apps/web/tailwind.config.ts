import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "../../packages/ui/src/**/*.{js,ts,jsx,tsx}",
    "../../packages/design-system/src/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      fontFamily: {
        display: ['var(--font-display)', 'Gowun Batang', 'Playfair Display', 'serif'],
        heading: ['var(--font-heading)', 'Outfit', 'Figtree', 'system-ui', 'sans-serif'],
        body: ['var(--font-body)', 'Noto Sans KR', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      colors: {
        // 브랜드 기본 (기획서 SSOT)
        'deep-sea': '#0A192F',
        'glass-navy': '#112240',
        'sanggam': {
          DEFAULT: '#D4AF37',
          light: 'rgba(212, 175, 55, 0.1)',
          hover: '#E5C349',
          border: 'rgba(212, 175, 55, 0.3)',
        },
        'wave-cyan': {
          DEFAULT: '#64FFDA',
          light: 'rgba(100, 255, 218, 0.08)',
          hover: '#7BFFE3',
        },
        'hanji': '#FAFAFA',
        'dancheong': '#D32F2F',
        // 호환성 유지
        primary: '#0A192F',
        accent: '#D4AF37',
        medical: { teal: '#0891B2', DEFAULT: '#0891B2', light: '#ECFEFF' },
        trust: { blue: '#1E40AF', DEFAULT: '#1E40AF', light: '#DBEAFE' },
      },
      borderRadius: {
        'sanggam': '12px',
      },
      boxShadow: {
        'sanggam': '0 0 20px rgba(212, 175, 55, 0.15)',
        'glass': '0 4px 30px rgba(0, 0, 0, 0.06)',
        'deep': '0 4px 20px rgba(0, 0, 0, 0.4)',
      },
    },
  },
  plugins: [],
};

export default config;
