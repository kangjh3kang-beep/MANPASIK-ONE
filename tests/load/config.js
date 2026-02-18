// ManPaSik Load Test Configuration
// k6 공유 설정 파일

export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const GRPC_ADDR = __ENV.GRPC_ADDR || 'localhost:50051';

export const THRESHOLDS = {
  auth: {
    http_req_duration: ['p(95)<200', 'p(99)<500'],
    http_req_failed: ['rate<0.01'],
  },
  measurement: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
  },
  gateway: {
    http_req_duration: ['p(95)<300', 'p(99)<800'],
    http_req_failed: ['rate<0.02'],
  },
};

export const HEADERS = {
  'Content-Type': 'application/json',
  'Accept': 'application/json',
};

export function getAuthHeaders(token) {
  return { ...HEADERS, 'Authorization': `Bearer ${token}` };
}

export function randomEmail() {
  return `loadtest_${Date.now()}_${Math.random().toString(36).substr(2, 9)}@test.manpasik.com`;
}

export function randomPassword() {
  return `Test!${Math.random().toString(36).substr(2, 12)}`;
}
