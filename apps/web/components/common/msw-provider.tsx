'use client';

import React, { useEffect, useState } from 'react';

export function MSWProvider(props: React.PropsWithChildren<{}>) {
  const { children } = props;
  // 초기 상태: 모킹이 필요 없으면 즉시 ready 상태로 시작
  const isMockingEnabled =
    process.env.NODE_ENV === 'development' ||
    process.env.NEXT_PUBLIC_API_MOCKING === 'enabled';
    
  const [mswReady, setMswReady] = useState(!isMockingEnabled);

  useEffect(() => {
    if (!isMockingEnabled) return;

    const init = async () => {
      try {
        console.log('[MSW] Initializing mocks...');
        const { initMocks } = await import('@mmup/api-client');
        await initMocks();
        console.log('[MSW] Mocks initialized successfully');
        setMswReady(true);
      } catch (err) {
        console.error('[MSW] Mock initialization failed:', err);
        setMswReady(true); // 에러가 나도 앱은 구동되어야 하므로 transition
      }
    };

    if (!mswReady) {
      init();
    }
  }, [mswReady]);

  if (!mswReady) {
    // MSW가 서비스 워커를 가로챌 준비가 될 때까지 렌더링을 지연시키거나 
    // 최소한의 로딩 상태를 보여줄 수 있습니다.
    return null; 
  }

  return <>{children}</>;
}
