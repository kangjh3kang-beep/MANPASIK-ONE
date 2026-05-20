# Measure gRPC Compose Container Smoke 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: containerized `measurement-service` HTTP/gRPC/lifecycle smoke 하네스

## 상세 구현계획

1. 실제 service binary smoke 다음 단계로, Docker Compose가 빌드한 `measurement-service` 컨테이너를 검증한다.
2. 전체 dev dependency stack을 모두 띄우기 전에, narrow compose 파일로 container image, runtime env, port publishing, health endpoint, gRPC endpoint를 먼저 고정한다.
3. `infrastructure/docker/docker-compose.measurement-smoke.yml`은 `backend/services/measurement-service/Dockerfile`을 그대로 사용해 `manpasik/measurement-service:smoke` 이미지를 빌드한다.
4. 컨테이너 내부는 `GRPC_PORT=:50054`, `HTTP_PORT=:8080`, `VERSION=compose-smoke`, `TENANCY_ENFORCED=false`로 기동한다.
5. host port는 스크립트가 `127.0.0.1` 빈 포트를 동적으로 예약해 `MANPASIK_MEASUREMENT_SMOKE_GRPC_PORT`, `MANPASIK_MEASUREMENT_SMOKE_HTTP_PORT`로 주입한다.
6. `scripts/measure_service_compose_smoke.sh`가 Docker/Compose/Go preflight를 수행하고, compose up/down cleanup을 담당한다.
7. 컨테이너가 뜨면 `TestMeasurementServiceExternalEndpointSmoke`가 host-published endpoint에 접속해 HTTP `/health`, gRPC health `SERVING`, `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` lifecycle을 검증한다.
8. Docker가 연결되지 않은 환경에서는 스크립트가 명시적인 blocked status를 출력하고 비정상 성공으로 숨기지 않는다.

## 구현 결과

- container smoke 전용 compose 파일을 추가했다.
- compose 실행 스크립트를 추가했다.
- 실제 외부 endpoint에 붙는 Go smoke test를 추가했다.
- 현재 WSL 환경에서는 Docker Desktop WSL integration이 연결되어 있지 않아 실제 compose up은 실행하지 못했다.
- 따라서 이번 gate의 상태는 `하네스 구축 완료 / 현재 환경 실행 blocked`이다.

## 변경 파일

- `infrastructure/docker/docker-compose.measurement-smoke.yml`
  - measurement-service 단일 컨테이너 smoke compose 정의 추가.
  - gRPC `50054`, HTTP `8080` container port를 host loopback 동적 포트에 publish.
  - HTTP `/health` 기반 healthcheck 추가.
- `scripts/measure_service_compose_smoke.sh`
  - Docker/Compose/Go preflight 추가.
  - 동적 host port 예약, compose up/down cleanup, external endpoint Go smoke 실행.
  - Docker Compose 미사용 가능 환경에서는 `MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable` 출력.
- `backend/services/measurement-service/cmd/main_smoke_test.go`
  - `TestMeasurementServiceExternalEndpointSmoke` 추가.
  - 기존 binary smoke lifecycle 검증을 공통 helper로 정리.
  - 외부 HTTP health, gRPC health, measurement lifecycle 검증 helper 추가.

## 실행 명령

```bash
cd /home/kangjh3kang/Manpasik
scripts/measure_service_compose_smoke.sh
```

Docker가 연결된 환경에서는 다음 순서로 실행된다.

1. `docker compose -p <project> -f infrastructure/docker/docker-compose.measurement-smoke.yml up -d --build measurement-service`
2. `go test -v ./backend/services/measurement-service/cmd -run TestMeasurementServiceExternalEndpointSmoke -count=1`
3. `docker compose ... down --remove-orphans -v`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
bash -n scripts/measure_service_compose_smoke.sh
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v \
  ./backend/services/measurement-service/cmd \
  -run 'TestMeasurementService(BinarySmoke|ExternalEndpointSmoke)'
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/measurement-service/cmd \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
export PATH=/mnt/d/우리집/flutter_cache/flutter/bin:/usr/local/bin:/usr/bin:/bin
flutter test --no-pub \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
flutter analyze --no-pub \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/features/measurement/data/measurement_repository_impl.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
```

결과:

- Script syntax: PASS
- Go binary smoke: PASS
- Go external endpoint smoke: SKIP when `MANPASIK_MEASUREMENT_SERVICE_SMOKE_ADDR` is unset
- Go cmd/handler/service regression: PASS
- Proto generation preflight: PASS, full replacement blocked by locked `protobuf 3.1.0` Timestamp import compatibility
- Flutter stream subset test/analyze: PASS
- Compose execution in current WSL: BLOCKED, `MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`

## 다음 단계

- Docker Desktop WSL integration을 활성화한 뒤 `scripts/measure_service_compose_smoke.sh`를 재실행한다.
- narrow container smoke가 PASS하면 full dev compose의 `measurement-service + postgres + gateway` 경로로 확장한다.
- Dart proto full generated compile gate는 `docs/audit/dart-proto-full-generated-compile-gate.md`에서 완료했다.
- Flutter 앱 의존성 lock 정렬과 checked-in official generated file 교체는 `docs/audit/dart-proto-official-generated-replacement.md`에서 완료했다.
