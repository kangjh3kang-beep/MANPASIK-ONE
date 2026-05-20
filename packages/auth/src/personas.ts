/**
 * @mmup-axis 9 전체
 * @mmup-stage 4 인증
 * @sb SB-AUTH
 * @standard IEC 62304 Class C
 *
 * 페르소나(역할) 정의 및 도메인 라우팅 맵
 * - 로그인 후 사용자의 역할에 따라 적절한 도메인 앱으로 자동 리다이렉트
 */

/** MMUP 플랫폼 페르소나 (역할) 타입 */
export type Persona =
  | 'patient'       // 환자 - 만파식 앱 (모바일/PWA)
  | 'doctor'        // 의사 - 임상 데이터 콘솔
  | 'researcher'    // 연구원 - AI 에이전트 허브 + 예측
  | 'pharma'        // 제약사 관계자 - GxP 준수
  | 'partner'       // 파트너 기관 - 파트너 통합 연동
  | 'developer'     // 개발자 - 개발자 포털
  | 'admin';        // 관리자 - 전체 대시보드

/** MMUP 확장 유저 타입 (NextAuth Session 확장) */
export interface MMUPUser {
  id: string;
  name: string;
  email: string;
  image?: string;
  persona: Persona;
  organization?: string;
  licenseNumber?: string;    // 의사 면허번호 (doctor 전용)
  clinicalTrialId?: string;  // 임상시험 ID (researcher 전용)
}

/**
 * 페르소나별 기본 리다이렉트 경로 매핑
 * - 로그인 성공 후 자동으로 해당 도메인 앱 URL로 리다이렉트
 */
export const PERSONA_ROUTES: Record<Persona, { label: string; path: string; port: number }> = {
  patient: {
    label: '만파식 건강 앱',
    path: '/ko',
    port: 3008,
  },
  doctor: {
    label: '임상 데이터 콘솔',
    path: '/ko',
    port: 3001,
  },
  researcher: {
    label: 'AI 에이전트 허브',
    path: '/ko',
    port: 3002,
  },
  pharma: {
    label: 'GxP 준수 시스템',
    path: '/ko',
    port: 3006,
  },
  partner: {
    label: '파트너 통합 연동',
    path: '/ko',
    port: 3005,
  },
  developer: {
    label: '개발자 포털',
    path: '/ko',
    port: 3007,
  },
  admin: {
    label: '관리자 대시보드',
    path: '/ko',
    port: 3000,
  },
};

/**
 * 페르소나별 리다이렉트 URL 생성
 */
export function getPersonaRedirectUrl(persona: Persona): string {
  const route = PERSONA_ROUTES[persona];
  return `http://localhost:${route.port}${route.path}`;
}
