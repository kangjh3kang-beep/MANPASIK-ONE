export async function initMocks() {
  if (typeof window === 'undefined') {
    // SSR 환경 (Node.js) - 추후 필요 시 server.ts 추가
    const { server } = await import('./server');
    server.listen({ onUnhandledRequest: 'bypass' });
  } else {
    // 브라우저 환경
    const { worker } = await import('./browser');
    await worker.start({ onUnhandledRequest: 'bypass' });
  }
}
