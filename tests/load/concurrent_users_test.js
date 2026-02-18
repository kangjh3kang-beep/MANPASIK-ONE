import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend } from 'k6/metrics';
import { BASE_URL, HEADERS, randomEmail, randomPassword, getAuthHeaders } from './config.js';

const sessionDuration = new Trend('user_session_duration');

// 동시 사용자 시뮬레이션: Phase 1 목표 100 CCU
export const options = {
  scenarios: {
    concurrent_users: {
      executor: 'constant-vus',
      vus: 100,
      duration: '5m',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    user_session_duration: ['p(95)<10000'],
  },
};

export default function () {
  const startTime = Date.now();
  const email = randomEmail();
  const password = randomPassword();

  group('1. 회원가입/로그인', () => {
    http.post(`${BASE_URL}/api/v1/auth/register`,
      JSON.stringify({ email, password, name: `CCU-${__VU}` }), { headers: HEADERS });
    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`,
      JSON.stringify({ email, password }), { headers: HEADERS });
    check(loginRes, { '로그인 성공': (r) => r.status === 200 });
  });

  sleep(1);

  group('2. 측정 시작', () => {
    http.post(`${BASE_URL}/api/v1/measurement/sessions`,
      JSON.stringify({ device_id: `MPK-CCU-${__VU}`, cartridge_id: 'CART-GLU-001', user_id: `ccu-${__VU}` }),
      { headers: HEADERS });
  });

  sleep(2);

  group('3. 결과 조회', () => {
    http.get(`${BASE_URL}/api/v1/measurement/history?limit=5`, { headers: HEADERS });
  });

  sleep(1);

  group('4. 로그아웃', () => {
    http.post(`${BASE_URL}/api/v1/auth/logout`, null, { headers: HEADERS });
  });

  sessionDuration.add(Date.now() - startTime);
  sleep(2);
}
