/**
 * @mmup-axis 9 전체
 * @mmup-stage 4 인증
 * @sb SB-AUTH
 * @standard IEC 62304 Class C
 *
 * 글로벌 미들웨어 헬퍼
 * - 각 도메인 앱의 middleware.ts에서 import하여 사용
 * - 인증되지 않은 유저를 /login으로 리다이렉트
 * - 페르소나 권한이 맞지 않으면 접근 차단
 */

import { type NextRequest, NextResponse } from 'next/server';
import { decode } from 'next-auth/jwt';
import type { Persona } from './personas';

interface AuthMiddlewareConfig {
  /** 이 도메인 앱에 접근 가능한 페르소나 목록 */
  allowedPersonas: Persona[];
  /** 로그인 페이지로의 리다이렉트 URL */
  loginUrl?: string;
  /** 권한 없음 시 리다이렉트 URL */
  unauthorizedUrl?: string;
  /** 인증 체크를 건너뛸 경로 패턴 */
  publicPaths?: string[];
}

const AUTH_SECRET = process.env.AUTH_SECRET || 'mmup-dev-secret-key-change-in-production';

/**
 * 글로벌 미들웨어 팩토리
 * - 각 도메인 앱은 이 함수를 호출해 자기만의 미들웨어를 만듭니다
 *
 * @example
 * // apps/clinical/middleware.ts
 * import { withAuth } from '@mmup/auth';
 * export default withAuth({
 *   allowedPersonas: ['doctor', 'researcher', 'admin'],
 * });
 */
export function withAuth(config: AuthMiddlewareConfig) {
  const {
    allowedPersonas,
    loginUrl = '/login',
    unauthorizedUrl = '/unauthorized',
    publicPaths = ['/api/auth', '/_next', '/favicon.ico'],
  } = config;

  return async function middleware(request: NextRequest) {
    const { pathname } = request.nextUrl;

    // 공개 경로는 체크 스킵
    if (publicPaths.some((p) => pathname.startsWith(p))) {
      return NextResponse.next();
    }

    // 세션 토큰 확인 (next-auth v5 쿠키)
    const sessionToken =
      request.cookies.get('authjs.session-token')?.value ||
      request.cookies.get('__Secure-authjs.session-token')?.value;

    if (!sessionToken) {
      const url = new URL(loginUrl, request.url);
      url.searchParams.set('callbackUrl', request.url);
      return NextResponse.redirect(url);
    }

    // JWT 페이로드에서 persona 추출 및 권한 체크
    try {
      const token = await decode({
        token: sessionToken,
        secret: AUTH_SECRET,
        salt: request.cookies.has('__Secure-authjs.session-token')
          ? '__Secure-authjs.session-token'
          : 'authjs.session-token',
      });

      if (!token) {
        const url = new URL(loginUrl, request.url);
        url.searchParams.set('callbackUrl', request.url);
        return NextResponse.redirect(url);
      }

      const persona = token.persona as Persona | undefined;

      // admin은 모든 도메인 접근 허용
      if (persona !== 'admin' && allowedPersonas.length > 0) {
        if (!persona || !allowedPersonas.includes(persona)) {
          return NextResponse.redirect(new URL(unauthorizedUrl, request.url));
        }
      }
    } catch {
      return NextResponse.redirect(new URL(unauthorizedUrl, request.url));
    }

    return NextResponse.next();
  };
}

