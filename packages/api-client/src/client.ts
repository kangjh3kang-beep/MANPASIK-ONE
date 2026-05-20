/**
 * @mmup/api-client — 통합 API 클라이언트
 * 
 * fetch wrapper + 에러 핸들링 + 재시도 로직 + 인터셉터
 * IEC 62304 Class B 준수
 */

export interface ApiConfig {
  baseUrl: string;
  token?: string;
  timeout?: number;
  retries?: number;
}

export interface ApiError {
  status: number;
  code: string;
  message: string;
  details?: Record<string, unknown>;
}

export interface ApiResponse<T> {
  data: T;
  meta?: {
    total?: number;
    page?: number;
    pageSize?: number;
  };
}

const DEFAULT_CONFIG: ApiConfig = {
  baseUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  timeout: 15000,
  retries: 2,
};

class MmupApiClient {
  private config: ApiConfig;
  private interceptors: Array<(init: RequestInit) => RequestInit> = [];

  constructor(config?: Partial<ApiConfig>) {
    this.config = { ...DEFAULT_CONFIG, ...config };
  }

  setToken(token: string) {
    this.config.token = token;
  }

  addInterceptor(fn: (init: RequestInit) => RequestInit) {
    this.interceptors.push(fn);
  }

  private async request<T>(
    method: string,
    path: string,
    body?: unknown,
    params?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    // 기본 베이스 URL 처리 (상대 경로인 경우 처리)
    let baseUrl = this.config.baseUrl;
    if (baseUrl.startsWith('/')) {
      if (typeof window !== 'undefined') {
        baseUrl = window.location.origin + baseUrl;
      } else {
        baseUrl = 'http://localhost:3000' + baseUrl; // SSR 대응
      }
    }

    const url = new URL(path.startsWith('/') ? path.slice(1) : path, baseUrl.endsWith('/') ? baseUrl : baseUrl + '/');
    if (params) {
      Object.entries(params).forEach(([k, v]) => url.searchParams.set(k, v));
    }

    let init: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...(this.config.token ? { Authorization: `Bearer ${this.config.token}` } : {}),
      },
      ...(body ? { body: JSON.stringify(body) } : {}),
      signal: AbortSignal.timeout(this.config.timeout || 15000),
    };

    // Apply interceptors
    for (const interceptor of this.interceptors) {
      init = interceptor(init);
    }

    let lastError: Error | null = null;
    const maxRetries = this.config.retries || 0;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        const res = await fetch(url.toString(), init);

        if (!res.ok) {
          const errorBody = await res.json().catch(() => ({}));
          const apiError: ApiError = {
            status: res.status,
            code: errorBody.code || `HTTP_${res.status}`,
            message: errorBody.message || res.statusText,
            details: errorBody.details,
          };
          throw apiError;
        }

        const data = await res.json();
        return data as ApiResponse<T>;
      } catch (err) {
        lastError = err as Error;
        if (attempt < maxRetries) {
          await new Promise(r => setTimeout(r, Math.pow(2, attempt) * 500));
        }
      }
    }

    throw lastError;
  }

  get<T>(path: string, params?: Record<string, string>) {
    return this.request<T>('GET', path, undefined, params);
  }

  post<T>(path: string, body?: unknown) {
    return this.request<T>('POST', path, body);
  }

  put<T>(path: string, body?: unknown) {
    return this.request<T>('PUT', path, body);
  }

  delete<T>(path: string) {
    return this.request<T>('DELETE', path);
  }
}

// Singleton instance
export const apiClient = new MmupApiClient();
export { MmupApiClient };
