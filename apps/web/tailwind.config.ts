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
        heading: ['var(--font-heading)', 'Figtree', 'system-ui', 'sans-serif'],
        body: ['var(--font-body)', 'Noto Sans KR', 'system-ui', 'sans-serif'],
      },
      colors: {
        primary: "#0A2540",
        accent: "#635BFF",
        medical: { teal: "#0891B2", DEFAULT: "#0891B2", light: "#ECFEFF" },
        trust: { blue: "#1E40AF", DEFAULT: "#1E40AF", light: "#DBEAFE" },
        science: { indigo: "#4F46E5", DEFAULT: "#4F46E5", light: "#EEF2FF" },
        bio: { cyan: "#22D3EE", DEFAULT: "#22D3EE", light: "#CFFAFE" },
      },
    },
  },
  plugins: [],
};

export default config;
