import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';
import { BASE_URL, HEADERS } from './config.js';

const authLatency = new Trend('auth_endpoint_latency');
const measureLatency = new Trend('measurement_endpoint_latency');
const userLatency = new Trend('user_endpoint_latency');
const deviceLatency = new Trend('device_endpoint_latency');

// API Gateway 부하 테스트: Kong Gateway를 통한 다양한 엔드포인트 테스트
export const options = {
  scenarios: {
    ramping_rate: {
      executor: 'ramping-arrival-rate',
      startRate: 10,
      timeUnit: '1s',
      preAllocatedVUs: 50,
      maxVUs: 300,
      stages: [
        { duration: '1m', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '1m', target: 200 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<300', 'p(99)<800'],
    http_req_failed: ['rate<0.02'],
    auth_endpoint_latency: ['p(95)<200'],
    measurement_endpoint_latency: ['p(95)<500'],
  },
};

export default function () {
  const endpoints = [
    { name: 'health', url: `${BASE_URL}/health`, metric: null },
    { name: 'auth', url: `${BASE_URL}/api/v1/auth/login`, metric: authLatency, method: 'POST', body: JSON.stringify({ email: 'test@test.com', password: 'test' }) },
    { name: 'devices', url: `${BASE_URL}/api/v1/devices`, metric: deviceLatency },
    { name: 'measurement', url: `${BASE_URL}/api/v1/measurement/history`, metric: measureLatency },
  ];

  const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];

  let res;
  if (endpoint.method === 'POST') {
    res = http.post(endpoint.url, endpoint.body || '{}', { headers: HEADERS, tags: { endpoint: endpoint.name } });
  } else {
    res = http.get(endpoint.url, { headers: HEADERS, tags: { endpoint: endpoint.name } });
  }

  if (endpoint.metric) endpoint.metric.add(res.timings.duration);
  check(res, { [`${endpoint.name} 응답`]: (r) => r.status < 500 });

  sleep(0.1);
}
