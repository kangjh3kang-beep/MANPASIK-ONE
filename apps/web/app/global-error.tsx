'use client';

/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정 (부팅 인프라)
 * @sb SB-BOOT
 * @standard IEC 62304 Class C
 *
 * Root Layout Error Boundary (Phase BW-2 — Harness H6 최외곽 격리)
 *
 * RootLayout 자체가 에러를 throw 하는 극단적 상황 대비. 이 fallback 은 자체
 * <html><body> 를 가져야 하며, 모든 Provider/SSR 의존성 외부에서 동작.
 *
 * 호출되는 경우:
 *   - RootLayout/Provider 자체 throw (drastic: 의존성 미설치, env throw 등)
 *   - 모든 Provider 트리 실패
 *
 * 일반 페이지 에러는 app/error.tsx 가 처리.
 */

interface GlobalErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function GlobalError({ error, reset }: GlobalErrorProps) {
  return (
    <html lang="ko">
      <body
        style={{
          margin: 0,
          fontFamily: "'Noto Sans KR', system-ui, sans-serif",
          background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
          color: '#e2e8f0',
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          padding: '32px',
        }}
      >
        <div style={{ textAlign: 'center', maxWidth: '480px' }}>
          <div
            style={{
              fontSize: '32px',
              fontWeight: 700,
              background: 'linear-gradient(135deg, #06b6d4, #3b82f6)',
              WebkitBackgroundClip: 'text',
              backgroundClip: 'text',
              color: 'transparent',
              marginBottom: '12px',
            }}
          >
            MMUP
          </div>
          <h2 style={{ fontSize: '20px', marginBottom: '12px' }}>
            애플리케이션을 시작할 수 없습니다
          </h2>
          <p
            style={{
              fontSize: '13px',
              color: '#94a3b8',
              marginBottom: '24px',
              wordBreak: 'keep-all',
            }}
          >
            루트 컴포넌트가 초기화에 실패했습니다. 캐시를 비우고 다시 시도해
            주세요.
          </p>
          {error.digest && (
            <p
              style={{
                fontSize: '11px',
                color: '#64748b',
                fontFamily: 'monospace',
                marginBottom: '20px',
              }}
            >
              오류 ID: {error.digest}
            </p>
          )}
          <div style={{ display: 'flex', gap: '8px', justifyContent: 'center', flexWrap: 'wrap' }}>
            <button
              type="button"
              onClick={reset}
              style={{
                background: '#0284c7',
                color: '#fff',
                border: 'none',
                borderRadius: '8px',
                padding: '10px 20px',
                fontSize: '14px',
                fontWeight: 600,
                cursor: 'pointer',
              }}
            >
              다시 시도
            </button>
            <button
              type="button"
              onClick={() => {
                if (typeof window !== 'undefined') {
                  try {
                    if ('serviceWorker' in navigator) {
                      navigator.serviceWorker.getRegistrations().then((regs) => {
                        regs.forEach((r) => r.unregister());
                      });
                    }
                    if ('caches' in window) {
                      caches.keys().then((keys) => {
                        keys.forEach((k) => caches.delete(k));
                      });
                    }
                    localStorage.clear();
                    sessionStorage.clear();
                  } catch {
                    /* noop */
                  }
                  window.location.replace(
                    window.location.pathname + '?nosw=1&_=' + Date.now()
                  );
                }
              }}
              style={{
                background: '#1e293b',
                color: '#e2e8f0',
                border: '1px solid #475569',
                borderRadius: '8px',
                padding: '10px 20px',
                fontSize: '14px',
                fontWeight: 600,
                cursor: 'pointer',
              }}
            >
              캐시 비우고 새로고침
            </button>
          </div>
        </div>
      </body>
    </html>
  );
}
