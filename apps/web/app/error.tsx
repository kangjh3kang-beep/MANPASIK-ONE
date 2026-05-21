'use client';

/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정 (부팅 인프라)
 * @sb SB-BOOT
 * @standard IEC 62304 Class B
 *
 * 페이지 단위 Error Boundary (Phase BW-2 — Harness H6 실패 격리)
 *
 * Next.js App Router 에서 페이지 렌더링 또는 데이터 페칭 중 에러 발생 시
 * 본 컴포넌트가 자동 호출되어 placeholder 빈 화면 대신 진단 가능한 UI 표시.
 *
 * - reset() 으로 사용자가 재시도 가능
 * - 진단 정보 (에러 메시지 + digest) 콘솔 + window.__MMUP_LAST_ERROR__ 노출
 * - production 에서는 user-friendly 메시지, development 에서는 상세 노출
 */

import { useEffect } from 'react';

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function PageError({ error, reset }: ErrorProps) {
  useEffect(() => {
    console.error('[MMUP error.tsx]', error);
    if (typeof window !== 'undefined') {
      (window as unknown as Record<string, unknown>).__MMUP_LAST_ERROR__ = {
        message: error.message,
        digest: error.digest,
        stack: error.stack,
        at: new Date().toISOString(),
      };
    }
  }, [error]);

  const isProduction = process.env.NODE_ENV === 'production';

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-8">
      <div className="max-w-md w-full text-center">
        <div className="mb-6 inline-flex h-20 w-20 items-center justify-center rounded-full bg-red-500/10 ring-2 ring-red-500/40">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="h-10 w-10 text-red-400"
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
          </svg>
        </div>

        <h2 className="mb-3 text-2xl font-bold text-white">화면 로딩 오류</h2>

        <p className="mb-6 text-sm text-slate-400">
          {isProduction
            ? '예상치 못한 오류가 발생했습니다. 다시 시도해 주세요.'
            : (error.message || 'Unknown error')}
        </p>

        {error.digest && (
          <p className="mb-4 text-xs text-slate-500 font-mono">
            오류 ID: {error.digest}
          </p>
        )}

        <div className="flex flex-col gap-2 sm:flex-row sm:justify-center">
          <button
            type="button"
            onClick={reset}
            className="rounded-lg bg-sky-600 px-5 py-2.5 text-sm font-semibold text-white hover:bg-sky-500 transition"
          >
            다시 시도
          </button>
          <button
            type="button"
            onClick={() => {
              if (typeof window !== 'undefined') {
                window.location.replace(window.location.pathname + '?nosw=1&_=' + Date.now());
              }
            }}
            className="rounded-lg border border-slate-600 bg-slate-800 px-5 py-2.5 text-sm font-semibold text-slate-200 hover:bg-slate-700 transition"
          >
            캐시 비우고 새로고침
          </button>
        </div>

        <p className="mt-8 text-xs text-slate-500">
          문제가 계속되면 브라우저 콘솔(F12) 의 <code className="text-sky-400">__MMUP_LAST_ERROR__</code> 정보로 문의해 주세요.
        </p>
      </div>
    </div>
  );
}
