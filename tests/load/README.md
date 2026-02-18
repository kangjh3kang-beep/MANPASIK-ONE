# ManPaSik 부하 테스트 (Load Tests)

## 사전 요구사항
- [k6](https://k6.io/) v0.47+ 설치
- ManPaSik 서비스가 실행 중이어야 함

## 실행 방법

```bash
# 인증 부하 테스트 (기본)
k6 run tests/load/auth_load_test.js

# 측정 부하 테스트
k6 run tests/load/measurement_load_test.js

# API Gateway 부하 테스트
k6 run tests/load/api_gateway_load_test.js

# 동시 사용자 시뮬레이션 (100 CCU)
k6 run tests/load/concurrent_users_test.js

# 스트레스 테스트 (500 VU까지)
k6 run tests/load/stress_test.js

# 스파이크 테스트
k6 run tests/load/spike_test.js

# 환경변수로 URL 오버라이드
BASE_URL=https://staging.manpasik.com k6 run tests/load/auth_load_test.js
```

## 성능 목표 (Phase 1)

| 지표 | 목표 |
|------|------|
| 동시 사용자 (CCU) | 100 |
| 인증 P95 | < 200ms |
| 측정 P95 | < 500ms |
| 에러율 | < 1% |

## 결과 해석
- `p(95)`: 95%의 요청이 이 시간 내 완료
- `http_req_failed`: 실패한 요청의 비율
- `vus`: 동시 가상 사용자 수
