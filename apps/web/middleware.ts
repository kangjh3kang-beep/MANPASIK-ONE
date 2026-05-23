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
  // 도메인 페이지는 데모 모드에서 인증 없이 접근 가능
  // 프로덕션 백엔드 연결 시 인증 활성화: NEXT_PUBLIC_AUTH_REQUIRED=true
  const authRequired = process.env.NEXT_PUBLIC_AUTH_REQUIRED === 'true';

  if (authRequired) {
    const { pathname } = request.nextUrl;
    if (pathname.includes('/domains')) {
      const authRes = await authMiddleware(request);
      if (authRes && authRes.status !== 200 && authRes.headers.get('location')) {
        return authRes;
      }
    }
  }

  return intlMiddleware(request);
}
 
export const config = {
  // 정적 자원, API, _next를 제외한 모든 경로에 미들웨어 적용
  matcher: ['/', '/(ko|en|ja|zh)/:path*', '/domains/:path*', '/login', '/unauthorized']
};
