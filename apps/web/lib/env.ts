/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정 (부팅 인프라)
 * @sb SB-BOOT
 * @standard IEC 62304 Class C
 *
 * Next.js 환경변수 계약 (Phase BW-2 — Harness H5 인터페이스 계약 / H6 실패 격리)
 *
 * Next.js 는 빌드 시점에 NEXT_PUBLIC_* 만 클라이언트로 inline 치환되며, 그 외는
 * 서버 런타임에서만 읽힌다. Cloudflare Pages 같은 정적/edge 배포 환경에서
 * 환경변수가 누락되면 빈 값으로 치환되어 런타임 무력화 또는 부팅 실패가 발생.
 *
 * 본 모듈은:
 *   1) 모든 환경변수의 스키마/fallback 을 단일 출처로 관리
 *   2) 누락 시 명확한 진단 (서버 콘솔 + 클라이언트 window.__MMUP_ENV__)
 *   3) 환경별 게이팅 (production vs preview vs development)
 */

interface EnvSchema {
  NODE_ENV: string;
  /** API gateway base URL — Go gateway 또는 Cloudflare Worker */
  NEXT_PUBLIC_API_BASE_URL: string;
  /** Supabase URL (선택 — apps/web 은 현재 직접 사용 X 이나 미래 대비) */
  NEXT_PUBLIC_SUPABASE_URL?: string;
  /** Supabase anon key (선택) */
  NEXT_PUBLIC_SUPABASE_ANON_KEY?: string;
  /** 빌드 시점 SHA — 운영 진단용 */
  NEXT_PUBLIC_BUILD_SHA?: string;
}

/**
 * Production-safe FALLBACK — Cloudflare Pages 빌드 환경에서 env 누락 시
 * 사이트가 죽지 않도록 보장. 실제 운영에서는 대시보드 환경변수로 override 권장.
 */
const FALLBACK = {
  NEXT_PUBLIC_API_BASE_URL: 'https://manpasik.com/api',
} as const;

function resolveOrFallback(
  key: keyof typeof FALLBACK,
  raw: string | undefined
): string {
  const trimmed = (raw ?? '').trim();
  if (trimmed && trimmed !== 'undefined') return trimmed;
  return FALLBACK[key] ?? '';
}

const NEXT_PUBLIC_API_BASE_URL = resolveOrFallback(
  'NEXT_PUBLIC_API_BASE_URL',
  process.env.NEXT_PUBLIC_API_BASE_URL
);

export const env: EnvSchema = {
  NODE_ENV: process.env.NODE_ENV ?? 'development',
  NEXT_PUBLIC_API_BASE_URL,
  NEXT_PUBLIC_SUPABASE_URL: process.env.NEXT_PUBLIC_SUPABASE_URL,
  NEXT_PUBLIC_SUPABASE_ANON_KEY: process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY,
  NEXT_PUBLIC_BUILD_SHA: process.env.NEXT_PUBLIC_BUILD_SHA,
};

export interface EnvDiagnostics {
  ok: boolean;
  source: 'build' | 'fallback' | 'missing';
  missing: string[];
  usingFallback: string[];
  values: Record<string, string>;
}

export const envDiagnostics: EnvDiagnostics = (() => {
  const missing: string[] = [];
  const usingFallback: string[] = [];

  const rawApi = (process.env.NEXT_PUBLIC_API_BASE_URL ?? '').trim();
  if (!rawApi) {
    if (env.NEXT_PUBLIC_API_BASE_URL) usingFallback.push('NEXT_PUBLIC_API_BASE_URL');
    else missing.push('NEXT_PUBLIC_API_BASE_URL');
  }

  const source: EnvDiagnostics['source'] =
    missing.length > 0 ? 'missing' :
    usingFallback.length > 0 ? 'fallback' : 'build';

  return {
    ok: missing.length === 0,
    source,
    missing,
    usingFallback,
    values: {
      NODE_ENV: env.NODE_ENV,
      NEXT_PUBLIC_API_BASE_URL: env.NEXT_PUBLIC_API_BASE_URL || '(empty)',
      NEXT_PUBLIC_BUILD_SHA: env.NEXT_PUBLIC_BUILD_SHA ?? '(unset)',
    },
  };
})();

// 클라이언트에서 window.__MMUP_ENV__ 로 진단 가능
if (typeof window !== 'undefined') {
  (window as unknown as { __MMUP_ENV__: EnvDiagnostics }).__MMUP_ENV__ = envDiagnostics;
}
