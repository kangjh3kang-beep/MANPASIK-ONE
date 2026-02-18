import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';
import { BASE_URL, HEADERS, randomEmail, randomPassword } from './config.js';

const errorRate = new Rate('errors');

// 스파이크 테스트: 급격한 부하 변동 시 시스템 회복 능력 검증
export const options = {
  scenarios: {
    spike: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 20 },   // 정상 부하
        { duration: '10s', target: 200 },   // 1차 스파이크
        { duration: '30s', target: 20 },    // 회복
        { duration: '10s', target: 200 },   // 2차 스파이크
        { duration: '30s', target: 20 },    // 회복
        { duration: '30s', target: 0 },     // 종료
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<3000'],
    http_req_failed: ['rate<0.15'],
  },
};

export default function () {
  const email = randomEmail();
  const password = randomPassword();

  const regRes = http.post(`${BASE_URL}/api/v1/auth/register`,
    JSON.stringify({ email, password, name: `Spike-${__VU}` }), { headers: HEADERS });
  check(regRes, { '가입': (r) => r.status < 400 }) || errorRate.add(1);

  sleep(0.1);

  const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`,
    JSON.stringify({ email, password }), { headers: HEADERS });
  check(loginRes, { '로그인': (r) => r.status === 200 }) || errorRate.add(1);

  sleep(0.5);
}
