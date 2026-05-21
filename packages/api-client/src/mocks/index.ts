export async function initMocks() {
  if (typeof window === 'undefined') {
    // SSR 환경 (Node.js)
    // Edge Runtime 호환성을 위해 Webpack 번들링에서 제외합니다.
    const { server } = await import(/* webpackIgnore: true */ './server');
    server.listen({ onUnhandledRequest: 'bypass' });
  } else {
    // 브라우저 환경
    const { worker } = await import('./browser');
    await worker.start({ onUnhandledRequest: 'bypass' });
  }
}
