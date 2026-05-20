/**
 * @mmup-axis 9 전체
 * @mmup-stage 4 인증
 * @family C
 * @trinity IP3
 * @sb SB-AUTH
 * @standard IEC 62304 Class C
 *
 * MMUP 통합 인증 코어 패키지
 * - 모든 9개 도메인 앱에서 공유하는 SSO(Single Sign-On) 엔진
 * - 페르소나 기반 역할(Role) 분기 로직
 */

import { PERSONA_ROUTES, type Persona, type MMUPUser } from './personas';
export { PERSONA_ROUTES, type Persona, type MMUPUser };
export { auth, signIn, signOut, handlers } from './auth-config';
export { withAuth } from './middleware-helper';


/**
 * 서버 컴포넌트용 현재 유저 세션 헬퍼
 */
export async function currentUser() {
  const { auth } = await import('./auth-config');
  const session = await auth();
  return session?.user as MMUPUser | undefined;
}

