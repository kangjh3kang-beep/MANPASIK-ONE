/**
 * Dev Portal 도메인 미들웨어
 * - next-intl 다국어 라우팅
 * - @mmup/auth 접근 제어 통합
 */

import createIntlMiddleware from 'next-intl/middleware';
import { withAuth } from '@mmup/auth';
import { NextRequest } from 'next/server';

const intlMiddleware = createIntlMiddleware({
  locales: ['ko', 'en', 'ja', 'zh'],
  defaultLocale: 'ko'
});

const authMiddleware = withAuth({
  allowedPersonas: ['developer', 'admin'],
});

export default async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (!pathname.startsWith('/_next') && !pathname.startsWith('/api')) {
    const authRes = await authMiddleware(request);
    if (authRes && authRes.status !== 200 && authRes.headers.get('location')) {
      return authRes;
    }
  }

  return intlMiddleware(request);
}

export const config = {
  matcher: ['/', '/(ko|en|ja|zh)/:path*']
};
