/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 4 인증
 * @sb SB-AUTH
 * @standard IEC 62304 Class C
 *
 * Clinical 도메인 미들웨어
 * - next-intl 다국어 라우팅 유지
 * - @mmup/auth 접근 제어 통합 (Sprint 5에서 활성화)
 */

import createMiddleware from 'next-intl/middleware';

// Phase 1: 다국어 라우팅 (현재 활성)
const intlMiddleware = createMiddleware({
  locales: ['ko', 'en', 'ja', 'zh'],
  defaultLocale: 'ko'
});

export default intlMiddleware;

export const config = {
  matcher: ['/', '/(ko|en|ja|zh)/:path*']
};

// Phase 2 참고사항 (Sprint 5에서 활성화):
// import { withAuth } from '@mmup/auth';
// allowedPersonas: ['doctor', 'researcher', 'admin']
