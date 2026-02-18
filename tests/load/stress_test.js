import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';
import { BASE_URL, HEADERS, randomEmail, randomPassword } from './config.js';

const errorRate = new Rate('errors');

// 스트레스 테스트: 시스템 한계점 식별
export const options = {
  scenarios: {
    stress: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },   // 준비 단계
        { duration: '3m', target: 200 },   // 증가 단계
        { duration: '2m', target: 500 },   // 최대 부하
        { duration: '1m', target: 500 },   // 최대 부하 유지
        { duration: '2m', target: 0 },     // 회복 단계
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<2000'],   // 스트레스 시 P95 2초 허용
    http_req_failed: ['rate<0.10'],      // 10% 에러율까지 허용
    errors: ['rate<0.15'],
  },
};

export default function () {
  const email = randomEmail();
  const password = randomPassword();

  // 회원가입 + 로그인 + 측정 전체 플로우
  const regRes = http.post(`${BASE_URL}/api/v1/auth/register`,
    JSON.stringify({ email, password, name: `Stress-${__VU}` }), { headers: HEADERS });
  check(regRes, { '가입 성공': (r) => r.status < 400 }) || errorRate.add(1);

  sleep(0.2);

  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`,
    JSON.stringify({ email, password }), { headers: HEADERS });
  const success = check(loginRes, { '로그인 성공': (r) => r.status === 200 });
  if (!success) { errorRate.add(1); return; }

  let token;
  try { token = JSON.parse(loginRes.body).access_token; } catch { return; }

  sleep(0.2);

  const measureRes = http.post(`${BASE_URL}/api/v1/measurement/sessions`,
    JSON.stringify({ device_id: `MPK-STRESS-${__VU}`, cartridge_id: 'CART-GLU-001', user_id: `user-${__VU}` }),
    { headers: { ...HEADERS, 'Authorization': `Bearer ${token}` } });
  check(measureRes, { '측정 시작 성공': (r) => r.status < 400 }) || errorRate.add(1);

  sleep(1);
}
