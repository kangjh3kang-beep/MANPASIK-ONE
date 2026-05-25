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
        display: ['var(--font-display)', 'Gowun Batang', 'serif'],
        heading: ['var(--font-heading)', 'Outfit', 'system-ui', 'sans-serif'],
        body: ['var(--font-body)', 'Noto Sans KR', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      colors: {
        // ── Tailwind 기본색을 CSS 변수로 연결 ──
        // 페이지에서 bg-white, bg-slate-50 등을 쓰면 테마 변수를 따름
        'white': 'var(--mps-bg-card)',
        'background': 'var(--mps-bg)',

        // 기존 호환 + 브랜드
        primary: '#0A192F',
        accent: '#D4AF37',
        'deep-sea': '#0A192F',
        'glass-navy': '#112240',
        sanggam: {
          DEFAULT: '#D4AF37',
          light: 'rgba(212, 175, 55, 0.1)',
          hover: '#E5C349',
          border: 'rgba(212, 175, 55, 0.3)',
        },
        'wave-cyan': {
          DEFAULT: '#64FFDA',
          light: 'rgba(100, 255, 218, 0.08)',
        },
        hanji: '#FAFAFA',
        dancheong: '#D32F2F',
        medical: { teal: '#0891B2', DEFAULT: '#0891B2' },
        trust: { blue: '#1E40AF', DEFAULT: '#1E40AF' },
      },
      backgroundColor: {
        // bg-slate-50 → CSS 변수 참조 (테마 자동 반영)
        'slate-50': 'var(--mps-bg)',
        'slate-100': 'var(--mps-bg-elevated)',
        // bg-*-50 (라이트 틴트) → 다크 모드 대응
        'sky-50': 'var(--mps-tint-sky)',
        'emerald-50': 'var(--mps-tint-emerald)',
        'purple-50': 'var(--mps-tint-purple)',
        'amber-50': 'var(--mps-tint-amber)',
        'rose-50': 'var(--mps-tint-rose)',
        'cyan-50': 'var(--mps-tint-cyan)',
        'indigo-50': 'var(--mps-tint-indigo)',
        'orange-50': 'var(--mps-tint-orange)',
        'teal-50': 'var(--mps-tint-teal)',
        'violet-50': 'var(--mps-tint-violet)',
        'blue-50': 'var(--mps-tint-blue)',
        'pink-50': 'var(--mps-tint-pink)',
        'green-50': 'var(--mps-tint-emerald)',
        'red-50': 'var(--mps-tint-rose)',
        'gray-50': 'var(--mps-bg)',
        'stone-50': 'var(--mps-bg)',
      },
      textColor: {
        // text-slate-900 → CSS 변수 참조
        'slate-900': 'var(--mps-text)',
        'slate-700': 'var(--mps-text)',
        'slate-600': 'var(--mps-text-secondary)',
        'slate-500': 'var(--mps-text-secondary)',
        'slate-400': 'var(--mps-text-muted)',
      },
      borderColor: {
        // border-slate-200 → CSS 변수 참조
        'slate-200': 'var(--mps-border)',
        'slate-100': 'var(--mps-border-subtle)',
      },
      boxShadow: {
        'sanggam': '0 0 20px rgba(212, 175, 55, 0.15)',
        'glass': '0 4px 30px rgba(0, 0, 0, 0.06)',
        'deep': '0 4px 20px rgba(0, 0, 0, 0.4)',
      },
      borderRadius: {
        'sanggam': '12px',
      },
    },
  },
  plugins: [],
};

export default config;
