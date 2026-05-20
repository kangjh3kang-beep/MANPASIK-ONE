# Gateway Auth Compose Container Smoke 게이트

**작성일**: 2026-05-02  
**작성자**: Codex  
**범위**: Gateway + Auth Service 컨테이너 경계의 REST-to-gRPC 사용자 식별자 계약 하네스

## 상세 구축계획

1. 전체 dev stack 대신 `auth-service`와 `gateway` 두 컨테이너만 띄운다.
2. auth-service 컨테이너는 DB/Redis 환경변수를 주입하지 않아 인메모리 repository로 실행한다.
3. gateway 컨테이너는 `AUTH_SERVICE_ADDR=auth-service:50051`로 auth-service 컨테이너에 붙는다.
4. gateway host HTTP 포트는 스크립트가 동적으로 예약해 충돌을 줄인다.
5. auth-service `/health`와 gateway `/health` healthcheck를 Compose에 설정한다.
6. Compose stack이 뜬 뒤 기존 `TestGatewayAuthExternalEndpointSmoke`를 실행한다.
7. smoke는 REST `register -> login -> refresh`를 호출하고 `user_id` 일치 및 access/refresh token 존재를 검증한다.
8. 실행 후 compose stack을 cleanup한다.
9. Docker가 없는 환경에서는 명확한 blocked status를 출력한다.

## 구현 결과

- `infrastructure/docker/docker-compose.gateway-auth-smoke.yml`을 추가했다.
- `scripts/gateway_auth_compose_smoke.sh`를 추가했다.
- Docker 연결 환경에서는 다음 명령으로 컨테이너 경계 smoke를 실행할 수 있다.

```bash
cd /home/kangjh3kang/Manpasik
scripts/gateway_auth_compose_smoke.sh
```

현재 환경 결과:

- `GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`

## 변경 파일

- `infrastructure/docker/docker-compose.gateway-auth-smoke.yml`
- `scripts/gateway_auth_compose_smoke.sh`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
bash -n scripts/gateway_auth_compose_smoke.sh
scripts/gateway_auth_compose_smoke.sh
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/gateway/cmd \
  ./backend/services/gateway/internal/handler \
  ./backend/services/auth-service/cmd \
  ./backend/services/auth-service/internal/service \
  ./backend/services/auth-service/internal/handler

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
flutter test --no-pub \
  test/features/auth/data/auth_repository_rest_test.dart \
  test/features/auth/data/auth_proto_contract_test.dart \
  test/features/auth/domain/auth_result_test.dart
flutter analyze --no-pub \
  lib/features/auth/data/auth_repository_rest.dart \
  lib/features/auth/data/auth_repository_impl.dart \
  test/features/auth/data/auth_repository_rest_test.dart \
  test/features/auth/data/auth_proto_contract_test.dart
```

결과:

- Script syntax: PASS
- Compose script runtime: BLOCKED, Docker Compose unavailable in current WSL
- Go Gateway/Auth regression: PASS
- Flutter Auth REST/proto/domain regression: PASS
- Flutter Auth targeted analyze: PASS
- `git diff --check` targeted files: PASS

## 잔여 리스크

- 현재 WSL distro에는 Docker Compose runtime이 연결되어 있지 않아 container smoke를 실제 PASS로 확정하지 못했다.
- Docker Desktop WSL integration 활성화 후 같은 스크립트를 재실행해야 한다.
- full dev compose의 gateway/auth/measurement 연쇄는 다음 단계에서 별도로 확장해야 한다.

## 다음 단계

- Docker Desktop WSL integration을 활성화한 뒤 `scripts/gateway_auth_compose_smoke.sh`를 재실행한다.
- 이후 `scripts/measure_service_compose_smoke.sh`와 full dev compose smoke를 순서대로 확정한다.
