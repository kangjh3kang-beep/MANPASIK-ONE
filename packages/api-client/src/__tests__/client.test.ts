/**
 * @mmup/api-client — API 클라이언트 단위 테스트
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MmupApiClient } from '../client';

describe('MmupApiClient', () => {
  let client: MmupApiClient;

  beforeEach(() => {
    client = new MmupApiClient({ baseUrl: 'http://test-api.local' });
    vi.restoreAllMocks();
  });

  it('기본 설정이 올바르게 생성된다', () => {
    const defaultClient = new MmupApiClient();
    expect(defaultClient).toBeDefined();
  });

  it('GET 요청을 올바르게 수행한다', async () => {
    const mockData = { data: [{ id: '1', name: 'Test Patient' }] };
    
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await client.get('/api/v1/patients');
    expect(result.data).toEqual(mockData.data);
    expect(fetch).toHaveBeenCalledWith(
      'http://test-api.local/api/v1/patients',
      expect.objectContaining({ method: 'GET' })
    );
  });

  it('POST 요청을 올바르게 수행한다', async () => {
    const body = { heartRate: 72, spo2: 98 };
    const mockResponse = { data: { id: 'meas-001' } };
    
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse),
    } as Response);

    const result = await client.post('/api/v1/measurements', body);
    expect(result.data).toEqual(mockResponse.data);
    expect(fetch).toHaveBeenCalledWith(
      'http://test-api.local/api/v1/measurements',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify(body),
      })
    );
  });

  it('JWT 토큰이 헤더에 포함된다', async () => {
    client.setToken('test-jwt-token');
    
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ data: {} }),
    } as Response);

    await client.get('/api/v1/protected');
    expect(fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer test-jwt-token',
        }),
      })
    );
  });

  it('HTTP 에러 시 ApiError를 던진다', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 404,
      statusText: 'Not Found',
      json: () => Promise.resolve({ code: 'NOT_FOUND', message: 'Resource not found' }),
    } as Response);

    const noRetryClient = new MmupApiClient({ baseUrl: 'http://test-api.local', retries: 0 });
    
    await expect(noRetryClient.get('/api/v1/unknown')).rejects.toMatchObject({
      status: 404,
      code: 'NOT_FOUND',
    });
  });

  it('쿼리 파라미터가 URL에 포함된다', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ data: [] }),
    } as Response);

    await client.get('/api/v1/patients', { page: '1', limit: '20' });
    
    const calledUrl = (fetch as any).mock.calls[0][0] as string;
    expect(calledUrl).toContain('page=1');
    expect(calledUrl).toContain('limit=20');
  });

  it('인터셉터가 요청을 수정한다', async () => {
    client.addInterceptor((init) => ({
      ...init,
      headers: {
        ...init.headers as Record<string, string>,
        'X-Custom-Header': 'test-value',
      },
    }));

    vi.spyOn(globalThis, 'fetch').mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ data: {} }),
    } as Response);

    await client.get('/api/v1/test');
    expect(fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          'X-Custom-Header': 'test-value',
        }),
      })
    );
  });
});
