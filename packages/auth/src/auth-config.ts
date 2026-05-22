/**
 * @mmup-axis 9 전체
 * @mmup-stage 4 인증
 * @sb SB-AUTH
 * @standard IEC 62304 Class C
 *
 * NextAuth v5 (Auth.js) 코어 설정
 * - Credentials 프로바이더 (개발 환경 데모용)
 * - 향후 Auth0 / Azure AD / Google OIDC 프로바이더 추가 가능
 */

import NextAuth from 'next-auth';
import Credentials from 'next-auth/providers/credentials';
import type { Persona, MMUPUser } from './personas';
import { getPersonaRedirectUrl } from './personas';

declare module 'next-auth' {
  interface Session {
    user: {
      id?: string;
      name?: string | null;
      email?: string | null;
      image?: string | null;
      persona?: Persona;
      organization?: string;
      licenseNumber?: string;
      clinicalTrialId?: string;
    };
  }
}

declare module 'next-auth/jwt' {
  interface JWT {
    persona?: Persona;
    organization?: string;
    licenseNumber?: string;
    clinicalTrialId?: string;
  }
}

/**
 * 데모용 사용자 데이터베이스 (개발 환경 전용)
 * - 실제 프로덕션에서는 PostgreSQL + Prisma ORM으로 대체
 */
const DEMO_USERS: Record<string, MMUPUser & { password: string }> = {
  'patient@mmup.io': {
    id: 'usr_patient_001',
    name: '김환자',
    email: 'patient@mmup.io',
    persona: 'patient',
    password: 'demo1234',
  },
  'doctor@mmup.io': {
    id: 'usr_doctor_001',
    name: '박의사',
    email: 'doctor@mmup.io',
    persona: 'doctor',
    licenseNumber: 'KR-MD-2024-00123',
    password: 'demo1234',
  },
  'researcher@mmup.io': {
    id: 'usr_researcher_001',
    name: '이연구원',
    email: 'researcher@mmup.io',
    persona: 'researcher',
    clinicalTrialId: 'CT-2024-MMUP-001',
    password: 'demo1234',
  },
  'pharma@mmup.io': {
    id: 'usr_pharma_001',
    name: '최제약',
    email: 'pharma@mmup.io',
    persona: 'pharma',
    organization: '한국제약 주식회사',
    password: 'demo1234',
  },
  'partner@mmup.io': {
    id: 'usr_partner_001',
    name: '정파트너',
    email: 'partner@mmup.io',
    persona: 'partner',
    organization: '서울대학교병원',
    password: 'demo1234',
  },
  'dev@mmup.io': {
    id: 'usr_dev_001',
    name: '강개발',
    email: 'dev@mmup.io',
    persona: 'developer',
    password: 'demo1234',
  },
  'admin@mmup.io': {
    id: 'usr_admin_001',
    name: '관리자',
    email: 'admin@mmup.io',
    persona: 'admin',
    password: 'demo1234',
  },
};

export const { auth, signIn, signOut, handlers } = NextAuth({
  providers: [
    Credentials({
      name: 'MMUP 계정',
      credentials: {
        email: { label: '이메일', type: 'email', placeholder: 'doctor@mmup.io' },
        password: { label: '비밀번호', type: 'password' },
      },
      async authorize(credentials) {
        const email = credentials?.email as string;
        const password = credentials?.password as string;

        const user = DEMO_USERS[email];
        if (!user || user.password !== password) {
          return null;
        }

        return {
          id: user.id,
          name: user.name,
          email: user.email,
          image: user.image,
        };
      },
    }),
  ],
  callbacks: {
    async jwt({ token, user }) {
      if (user?.email) {
        const demoUser = DEMO_USERS[user.email];
        if (demoUser) {
          token.persona = demoUser.persona;
          token.organization = demoUser.organization;
          token.licenseNumber = demoUser.licenseNumber;
          token.clinicalTrialId = demoUser.clinicalTrialId;
        }
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        session.user.persona = token.persona;
        session.user.organization = token.organization;
        session.user.licenseNumber = token.licenseNumber;
        session.user.clinicalTrialId = token.clinicalTrialId;
      }
      return session;
    },
  },
  pages: {
    signIn: '/login',
  },
  session: {
    strategy: 'jwt',
  },
  secret: process.env.AUTH_SECRET || 'mmup-dev-secret-key-change-in-production',
});
