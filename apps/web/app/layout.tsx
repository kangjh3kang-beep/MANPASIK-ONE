/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정
 * @sb SB-1
 *
 * 루트 레이아웃 — 모든 페이지에 Tailwind CSS 적용
 */
import React from 'react';
import './globals.css';
import { ApiProvider } from '../components/api-provider';
import { MSWProvider } from '../components/common/msw-provider';
import { AuthProvider } from '../components/common/auth-provider';

export const metadata = {
  title: 'MMUP [VERIFIED_SYNC] 생태계',
  description: 'MMUP 통합 의료 생태계 플랫폼 — 정밀 의학을 일상으로',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ko" suppressHydrationWarning>
      <body suppressHydrationWarning className="antialiased">
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

