import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { BASE_URL, HEADERS, THRESHOLDS, randomEmail, randomPassword } from './config.js';

const loginDuration = new Trend('login_duration');
const registerDuration = new Trend('register_duration');
const errorRate = new Rate('errors');

export const options = {
  scenarios: {
    constant_load: {
      executor: 'constant-vus',
      vus: 10,
      duration: '30s',
      tags: { scenario: 'constant' },
    },
    ramping_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 50 },
        { duration: '30s', target: 100 },
        { duration: '30s', target: 50 },
        { duration: '30s', target: 0 },
      ],
      tags: { scenario: 'ramping' },
    },
  },
  thresholds: {
    http_req_duration: THRESHOLDS.auth.http_req_duration,
    http_req_failed: THRESHOLDS.auth.http_req_failed,
    errors: ['rate<0.05'],
  },
};

export default function () {
  const email = randomEmail();
  const password = randomPassword();

  group('회원가입 플로우', () => {
    const registerRes = http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({
      email: email,
      password: password,
      name: 'Load Test User',
    }), { headers: HEADERS });

    registerDuration.add(registerRes.timings.duration);
    check(registerRes, {
      '회원가입 성공 (200/201)': (r) => r.status === 200 || r.status === 201,
    }) || errorRate.add(1);
  });

  sleep(0.5);

  group('로그인 플로우', () => {
    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
      email: email,
      password: password,
    }), { headers: HEADERS });

    loginDuration.add(loginRes.timings.duration);
    const success = check(loginRes, {
      '로그인 성공 (200)': (r) => r.status === 200,
      '토큰 반환': (r) => {
        try { return JSON.parse(r.body).access_token !== undefined; }
        catch { return false; }
      },
    });
    if (!success) errorRate.add(1);

    if (loginRes.status === 200) {
      let token;
      try { token = JSON.parse(loginRes.body).access_token; } catch { return; }

      // 토큰 갱신 테스트
      sleep(0.3);
      const refreshRes = http.post(`${BASE_URL}/api/v1/auth/refresh`, JSON.stringify({
        refresh_token: JSON.parse(loginRes.body).refresh_token || '',
      }), { headers: HEADERS });

      check(refreshRes, {
        '토큰 갱신 성공': (r) => r.status === 200,
      });
    }
  });

  sleep(1);
}
