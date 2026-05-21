export async function initMocks() {
  if (typeof window !== 'undefined') {
    // 브라우저 환경
    const { worker } = await import('./browser');
    await worker.start({ onUnhandledRequest: 'bypass' });
  } else {
    console.log('[MSW] Server-side mocking is disabled in this environment.');
  }
}
