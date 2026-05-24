/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정
 * @sb SB-1
 *
 * 루트 레이아웃 — 모든 페이지에 Tailwind CSS 적용
 */
import React from 'react';
import { Outfit, Noto_Sans_KR, Gowun_Batang } from 'next/font/google';
import './globals.css';
import { ApiProvider } from '../components/api-provider';
import { MSWProvider } from '../components/common/msw-provider';
import { AuthProvider } from '../components/common/auth-provider';

const outfit = Outfit({
  subsets: ['latin'],
  weight: ['300', '400', '500', '600', '700'],
  variable: '--font-heading',
  display: 'swap',
});

const notoSansKR = Noto_Sans_KR({
  subsets: ['latin'],
  weight: ['400', '700'],
  variable: '--font-body',
  display: 'swap',
});

const gowunBatang = Gowun_Batang({
  subsets: ['latin'],
  weight: ['400', '700'],
  variable: '--font-display',
  display: 'swap',
});

export const metadata = {
  title: 'MPS 만파식 — 정밀 건강 측정 플랫폼',
  description: 'MPS 만파식 다중측정원리 유니버설 POCT 플랫폼 — 정밀 의학을 일상으로',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ko" suppressHydrationWarning className={`${outfit.variable} ${notoSansKR.variable} ${gowunBatang.variable}`}>
      <body suppressHydrationWarning className="antialiased font-body">
        <MSWProvider>
          <AuthProvider>
            <ApiProvider>
              {children}
            </ApiProvider>
          </AuthProvider>
        </MSWProvider>
      </body>
    </html>
  );
}

