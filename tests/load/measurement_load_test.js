import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { BASE_URL, HEADERS, THRESHOLDS, randomEmail, randomPassword, getAuthHeaders } from './config.js';

const sessionStartDuration = new Trend('session_start_duration');
const sessionEndDuration = new Trend('session_end_duration');
const historyDuration = new Trend('history_query_duration');
const errorRate = new Rate('errors');

export const options = {
  scenarios: {
    constant_rate: {
      executor: 'constant-arrival-rate',
      rate: 50,
      timeUnit: '1s',
      duration: '1m',
      preAllocatedVUs: 100,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_duration: THRESHOLDS.measurement.http_req_duration,
    http_req_failed: THRESHOLDS.measurement.http_req_failed,
    session_start_duration: ['p(95)<500'],
    session_end_duration: ['p(95)<500'],
    history_query_duration: ['p(95)<300'],
  },
};

function login() {
  const email = randomEmail();
  const password = randomPassword();
  http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({ email, password, name: 'Measure Test' }), { headers: HEADERS });
  const res = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({ email, password }), { headers: HEADERS });
  try { return JSON.parse(res.body).access_token; } catch { return null; }
}

export function setup() {
  return { token: login() };
}

export default function (data) {
  const token = data.token;
  if (!token) return;
  const authHeaders = getAuthHeaders(token);

  group('측정 세션 시작', () => {
    const res = http.post(`${BASE_URL}/api/v1/measurement/sessions`, JSON.stringify({
      device_id: `MPK-LOAD-${__VU}`,
      cartridge_id: `CART-GLU-${__ITER % 100}`,
      user_id: `user-load-${__VU}`,
    }), { headers: authHeaders });

    sessionStartDuration.add(res.timings.duration);
    const success = check(res, { '세션 시작 성공': (r) => r.status === 200 || r.status === 201 });
    if (!success) errorRate.add(1);
  });

  sleep(0.5);

  group('측정 세션 종료', () => {
    const res = http.post(`${BASE_URL}/api/v1/measurement/sessions/end`, JSON.stringify({
      session_id: `session-${__VU}-${__ITER}`,
      raw_channels: Array.from({ length: 88 }, () => Math.random() * 2),
      primary_value: 95 + Math.random() * 50,
      unit: 'mg/dL',
      confidence: 0.85 + Math.random() * 0.15,
    }), { headers: authHeaders });

    sessionEndDuration.add(res.timings.duration);
    check(res, { '세션 종료 성공': (r) => r.status === 200 });
  });

  sleep(0.3);

  group('측정 기록 조회', () => {
    const res = http.get(`${BASE_URL}/api/v1/measurement/history?limit=10`, { headers: authHeaders });
    historyDuration.add(res.timings.duration);
    check(res, { '기록 조회 성공': (r) => r.status === 200 });
  });

  sleep(0.5);
}
