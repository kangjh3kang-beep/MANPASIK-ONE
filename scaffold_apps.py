import os

base_dir = "."

apps = [
    "app", "clinical", "agents-hub", "predictor", "reward", "partner", "gxp", "dev-portal"
]

files_to_create = {}

for app in apps:
    app_prefix = f"apps/{app}"
    
    files_to_create[f"{app_prefix}/package.json"] = f"""{{
  "name": "@mmup/{app}",
  "version": "1.0.0",
  "private": true,
  "scripts": {{
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  }},
  "dependencies": {{
    "next": "^15.0.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "next-intl": "^3.0.0",
    "@mmup/ui": "workspace:*",
    "@mmup/design-system": "workspace:*",
    "@mmup/ssot": "workspace:*"
  }},
  "devDependencies": {{
    "tailwindcss": "^3.4.1",
    "autoprefixer": "^10.4.17",
    "postcss": "^8.4.35"
  }}
}}
"""
    
    files_to_create[f"{app_prefix}/next.config.mjs"] = """import createNextIntlPlugin from 'next-intl/plugin';
const withNextIntl = createNextIntlPlugin('./i18n/request.ts');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@mmup/ui', '@mmup/design-system', '@mmup/ssot', '@mmup/types'],
};
export default withNextIntl(nextConfig);
"""
    
    files_to_create[f"{app_prefix}/middleware.ts"] = """import createMiddleware from 'next-intl/middleware';
 
export default createMiddleware({
  locales: ['ko', 'en', 'ja', 'zh'],
  defaultLocale: 'ko'
});
 
export const config = {
  matcher: ['/', '/(ko|en|ja|zh)/:path*']
};
"""

    files_to_create[f"{app_prefix}/messages/ko.json"] = """{
  "Index": {
    "title": "MMUP 만파식 생태계"
  }
}"""
    files_to_create[f"{app_prefix}/messages/en.json"] = """{
  "Index": {
    "title": "MMUP Platform Ecosystem"
  }
}"""

    files_to_create[f"{app_prefix}/app/[locale]/layout.tsx"] = """/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';
import { TrinityIndicator } from '@mmup/ui';
import '../globals.css';

export default async function LocaleLayout({
  children,
  params
}: {
  children: React.ReactNode;
  params: Promise<{locale: string}>;
}) {
  const { locale } = await params;
  const messages = await getMessages();
 
  return (
    <html lang={locale} suppressHydrationWarning>
      <body suppressHydrationWarning>
        <NextIntlClientProvider messages={messages}>
          {children}
          <TrinityIndicator />
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
"""

    files_to_create[f"{app_prefix}/app/[locale]/page.tsx"] = f"""/**
 * @mmup-axis 1 측정
 * @mmup-stage 1 측정
 * @sb SB-1
 */
import {{ useTranslations }} from 'next-intl';
import {{ DataIntegrityConsole }} from '@mmup/ui';

export default function Index() {{
  const t = useTranslations('Index');
  return (
    <main className="min-h-screen bg-slate-50 p-8">
      <div className="mx-auto max-w-7xl space-y-6">
        <h1 className="text-3xl font-bold tracking-tight text-slate-900">{{t('title')}} - {app.upper()}</h1>
        <p className="text-slate-500">
          Bolt 저장소에서 이식해 온 실시간 무결성 검증 대시보드입니다. 센서로부터 유입되는 Hash 체인을 시각적으로 검사합니다.
        </p>
        <div className="mt-8">
          <DataIntegrityConsole />
        </div>
      </div>
    </main>
  );
}}
"""

    files_to_create[f"{app_prefix}/app/globals.css"] = """@tailwind base;
@tailwind components;
@tailwind utilities;
"""

    files_to_create[f"{app_prefix}/postcss.config.mjs"] = """export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
"""

    files_to_create[f"{app_prefix}/tailwind.config.ts"] = """import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "../../packages/ui/src/**/*.{js,ts,jsx,tsx}",
    "../../packages/design-system/src/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        primary: "#0A2540",
        accent: "#635BFF",
      }
    },
  },
  plugins: [],
};

export default config;
"""

    files_to_create[f"{app_prefix}/i18n/request.ts"] = """import {getRequestConfig} from 'next-intl/server';
 
export default getRequestConfig(async ({requestLocale}) => {
  let locale = await requestLocale;
  if (!locale) locale = 'ko';
 
  return {
    locale,
    messages: (await import(`../messages/${locale}.json`)).default
  };
});
"""

# Storybook Config
sb_prefix = "packages/ui/.storybook"
files_to_create[f"{sb_prefix}/main.ts"] = """import type { StorybookConfig } from '@storybook/react-vite';
const config: StorybookConfig = {
  stories: ['../src/**/*.stories.@(js|jsx|mjs|ts|tsx)'],
  addons: ['@storybook/addon-essentials'],
  framework: {
    name: '@storybook/react-vite',
    options: {},
  },
};
export default config;
"""

for fp, contents in files_to_create.items():
    full_path = os.path.join(base_dir, fp.replace("/", os.sep))
    os.makedirs(os.path.dirname(full_path), exist_ok=True)
    with open(full_path, "w", encoding="utf-8") as f:
        f.write(contents)

print(f"Created {len(files_to_create)} files successfully for Sprint 2 and 3 foundation.")
