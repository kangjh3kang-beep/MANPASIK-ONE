import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { THRESHOLDS } from './config.js';

// 커스텀 메트릭
const grpcErrorRate = new Rate('grpc_error_rate');
const measurementDuration = new Trend('grpc_measurement_duration');
const deviceQueryDuration = new Trend('grpc_device_query_duration');
const healthCheckDuration = new Trend('grpc_health_check_duration');

const GRPC_ADDR = __ENV.GRPC_ADDR || 'localhost:50051';

const client = new grpc.Client();
client.load(
  ['../../backend/shared/proto'],
  'manpasik.proto',
  'health.proto'
);

export const options = {
  scenarios: {
    // 시나리오 1: gRPC 헬스체크 부하
    grpc_health: {
      executor: 'constant-arrival-rate',
      rate: 100,
      timeUnit: '1s',
      duration: '2m',
      preAllocatedVUs: 50,
      maxVUs: 200,
      exec: 'healthCheck',
    },
    // 시나리오 2: 측정 서비스 gRPC 부하
    grpc_measurement: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '1m', target: 200 },
        { duration: '30s', target: 0 },
      ],
      exec: 'measurementFlow',
    },
    // 시나리오 3: 디바이스 서비스 gRPC 부하
    grpc_device: {
      executor: 'constant-vus',
      vus: 30,
      duration: '2m',
      exec: 'deviceQuery',
    },
    // 시나리오 4: 혼합 gRPC 트래픽
    grpc_mixed: {
      executor: 'ramping-arrival-rate',
      startRate: 10,
      timeUnit: '1s',
      stages: [
        { duration: '30s', target: 50 },
        { duration: '1m', target: 100 },
        { duration: '30s', target: 50 },
      ],
      preAllocatedVUs: 100,
      maxVUs: 300,
      exec: 'mixedTraffic',
    },
  },
  thresholds: {
    grpc_error_rate: ['rate<0.01'],                // gRPC 에러율 < 1%
    grpc_measurement_duration: ['p(95)<500'],       // 측정 p95 < 500ms
    grpc_device_query_duration: ['p(95)<200'],      // 디바이스 조회 p95 < 200ms
    grpc_health_check_duration: ['p(99)<100'],      // 헬스체크 p99 < 100ms
    'grpc_req_duration{scenario:grpc_health}': ['p(95)<50'],
  },
};

// 시나리오 1: 헬스체크
export function healthCheck() {
  client.connect(GRPC_ADDR, { plaintext: true });

  const start = Date.now();
  const response = client.invoke('grpc.health.v1.Health/Check', {
    service: 'manpasik.MeasurementService',
  });
  healthCheckDuration.add(Date.now() - start);

  const success = check(response, {
    'health check status OK': (r) => r && r.status === grpc.StatusOK,
    'serving status': (r) => r && r.message && r.message.status === 'SERVING',
  });
  grpcErrorRate.add(!success);

  client.close();
}

// 시나리오 2: 측정 플로우
export function measurementFlow() {
  client.connect(GRPC_ADDR, { plaintext: true });

  // 세션 시작
  const startTime = Date.now();
  const sessionResp = client.invoke('manpasik.v1.MeasurementService/StartSession', {
    device_id: `device-load-${__VU}`,
    cartridge_type: '0x01',
    channel_count: 88,
  });

  const sessionOk = check(sessionResp, {
    'session started': (r) => r && r.status === grpc.StatusOK,
  });
  grpcErrorRate.add(!sessionOk);

  if (sessionOk && sessionResp.message) {
    // 측정 데이터 전송 (88채널 시뮬레이션)
    const channels = Array.from({ length: 88 }, (_, i) => ({
      channel_id: i,
      detector_value: Math.random() * 1000,
      reference_value: Math.random() * 100,
    }));

    const measureResp = client.invoke('manpasik.v1.MeasurementService/EndSession', {
      session_id: sessionResp.message.session_id,
      channels: channels,
    });

    measurementDuration.add(Date.now() - startTime);

    const measureOk = check(measureResp, {
      'measurement completed': (r) => r && r.status === grpc.StatusOK,
      'has corrected values': (r) => r && r.message && r.message.corrected_values,
    });
    grpcErrorRate.add(!measureOk);
  }

  client.close();
  sleep(0.5);
}

// 시나리오 3: 디바이스 조회
export function deviceQuery() {
  client.connect(GRPC_ADDR, { plaintext: true });

  const start = Date.now();
  const resp = client.invoke('manpasik.v1.DeviceService/GetDevices', {
    user_id: `user-load-${__VU}`,
    page: 1,
    page_size: 10,
  });
  deviceQueryDuration.add(Date.now() - start);

  const success = check(resp, {
    'device query success': (r) => r && r.status === grpc.StatusOK,
  });
  grpcErrorRate.add(!success);

  client.close();
  sleep(1);
}

// 시나리오 4: 혼합 트래픽
export function mixedTraffic() {
  const rand = Math.random();

  if (rand < 0.4) {
    healthCheck();
  } else if (rand < 0.7) {
    deviceQuery();
  } else {
    measurementFlow();
  }
}
