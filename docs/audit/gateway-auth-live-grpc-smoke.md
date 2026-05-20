# Gateway Auth Live gRPC Smoke 게이트

**작성일**: 2026-05-02  
**작성자**: Codex  
**범위**: 실제 Gateway binary와 실제 Auth Service binary 사이의 REST-to-gRPC 사용자 식별자 전달 검증

## 상세 구축계획

1. 테스트 안에서 auth-service production `main` binary를 `go build`로 생성한다.
2. 테스트 안에서 gateway production `main` binary를 `go build`로 생성한다.
3. auth-service는 DB/Redis env 없이 인메모리 User/Token repository로 기동한다.
4. auth-service gRPC 포트와 metrics 포트를 동적 loopback 주소로 격리한다.
5. gateway는 동적 HTTP 포트로 기동하고 `AUTH_SERVICE_ADDR`에 실제 auth-service gRPC 주소를 주입한다.
6. REST `POST /api/v1/auth/register`로 사용자 ID를 확보한다.
7. REST `POST /api/v1/auth/login` 응답의 `user_id`가 register의 `user_id`와 같은지 확인한다.
8. REST `POST /api/v1/auth/refresh` 응답의 `user_id`도 같은 사용자 ID인지 확인한다.
9. access/refresh token 존재, gateway `/health`, auth gRPC health `SERVING`, graceful shutdown을 함께 검증한다.

## 구현 결과

- `TestGatewayAuthLiveBinarySmoke`가 실제 OS process 2개를 띄워 REST-to-gRPC 경계를 검증한다.
- Gateway mock E2E가 아니라 실제 gateway binary, 실제 auth-service binary, real TCP HTTP/gRPC stack을 통과한다.
- auth-service metrics/health HTTP server는 `METRICS_PORT` env로 동적 포트 격리가 가능하다.
- `MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_ADDR`를 지정하면 기존 외부 Gateway endpoint에도 같은 auth lifecycle smoke를 실행할 수 있다.

## 변경 파일

- `backend/services/auth-service/cmd/main.go`
- `backend/services/gateway/cmd/main_smoke_test.go`
- `infrastructure/docker/docker-compose.gateway-auth-smoke.yml`
- `scripts/gateway_auth_compose_smoke.sh`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
MANPASIK_GO_BINARY=/home/kangjh3kang/sdk/go-go1.26.2/bin/go \
  /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v -count=1 \
  ./backend/services/gateway/cmd \
  -run TestGatewayAuthLiveBinarySmoke

/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/gateway/cmd \
  ./backend/services/gateway/internal/handler \
  ./backend/services/auth-service/cmd \
  ./backend/services/auth-service/internal/service \
  ./backend/services/auth-service/internal/handler
```

결과:

- Gateway auth live binary smoke: PASS
- Gateway cmd/handler + Auth cmd/service/handler regression: PASS
- Flutter Auth REST/proto/domain regression: PASS
- Flutter Auth targeted analyze: PASS
- `git diff --check` targeted files: PASS
- 2026-05-02 Gateway auth compose container smoke harness: BUILT, current WSL Docker Compose runtime BLOCKED

## 잔여 리스크

- 이번 smoke는 Docker 없이 local process boundary를 검증한다.
- gateway container와 auth-service container 사이 통신을 검증하는 Compose 하네스는 추가됐지만, 현재 WSL Docker Compose runtime 부재로 실제 PASS는 보류다.
- full dev compose의 gateway/auth/measurement 연쇄는 Docker WSL integration 활성화 후 별도 확인이 필요하다.
- social-login live path는 외부 OAuth verifier 의존성이 있어 이번 smoke 범위에서는 register/login/refresh 사용자 식별자 경로를 우선 검증했다.

## 다음 단계

- Docker Desktop WSL integration 활성화 후 measurement compose smoke를 PASS로 확정한다.
- 이후 full dev compose에서 gateway + auth-service container auth smoke를 추가한다.
