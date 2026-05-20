export async function register() {
  // 서버 사이드 모킹 활성화
  if (process.env.NEXT_RUNTIME === 'nodejs' && process.env.NEXT_PUBLIC_API_MOCKING === 'enabled') {
    const { initMocks } = await import('@mmup/api-client');
    await initMocks();
  }
}
