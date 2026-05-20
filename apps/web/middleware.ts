import { NextRequest } from 'next/server';
import createIntlMiddleware from 'next-intl/middleware';
import { withAuth } from '@mmup/auth';

const intlMiddleware = createIntlMiddleware({
  locales: ['ko', 'en', 'ja', 'zh'],
  defaultLocale: 'ko'
});

const authMiddleware = withAuth({
  allowedPersonas: ['patient', 'doctor', 'researcher', 'admin', 'pharma', 'partner', 'developer'],
  publicPaths: ['/api', '/_next', '/favicon.ico', '/login'],
});

export default async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  
  // 1. 도메인 경로는 인증 먼저 체크
  if (pathname.includes('/domains')) {
    const authRes = await authMiddleware(request);
    // 리다이렉트가 발생했다면(로그인 필요 등) 바로 반환
    if (authRes && authRes.status !== 200 && authRes.headers.get('location')) {
      return authRes;
    }
  }

  // 2. 다국어 처리 진행
  return intlMiddleware(request);
}
 
export const config = {
  matcher: ['/', '/(ko|en|ja|zh)/:path*', '/domains/:path*']
};

