# E2E 테스트

E2E 테스트는 **backend 모듈**에서 실행됩니다 (gRPC 클라이언트·shared gen 사용).

## 실행 방법

```bash
cd backend
go test -v ./tests/e2e/...
```

또는 프로젝트 루트에서:

```bash
make test-integration
```

## 서비스 기동 후 전체 플로우 검증

Postgres + Go 4서비스 기동 후 동일 명령으로 실행하면 `TestServiceHealth`(4서비스 헬스체크)와 `TestMeasurementFlow`(Login → StartSession → EndSession → GetHistory)가 실제 연결로 동작합니다.

**반드시 프로젝트 루트(Manpasik)에서 실행하세요.** 다른 디렉터리(예: `frontend/flutter-app`)에서는 `infrastructure/docker` 경로를 찾지 못합니다.

- **Docker**: 최신 Docker Desktop은 `docker compose`(공백, V2)를 사용합니다. WSL 2에서 `docker-compose`를 찾을 수 없다면 `docker compose`로 실행하거나, [Docker Desktop WSL 통합](https://docs.docker.com/go/wsl2/)을 활성화하세요.
- **실행 위치**: `docker-compose.dev.yml`은 **infrastructure/docker/** 안에 있으므로, 반드시 **프로젝트 루트**에서 `cd infrastructure/docker` 후 실행하거나, 루트에서 `-f infrastructure/docker/docker-compose.dev.yml` 처럼 경로를 지정하세요. `backend/` 디렉터리에서는 해당 파일이 없어 오류가 납니다.

**주의:** 여러 줄을 한꺼번에 붙여넣으면 줄바꿈이 사라져 `./tests/e2e/...` 뒤에 다음 명령이 붙어 `./tests/e2e/...cd` 처럼 해석될 수 있습니다. **한 줄씩 실행**하거나, 아래처럼 **서비스 기동**과 **E2E 실행**을 나눠서 하세요.

**1단계 — 서비스 기동 (프로젝트 루트에서 한 줄):**
```bash
cd ~/Manpasik && cd infrastructure/docker && docker compose -f docker-compose.dev.yml up -d postgres auth-service user-service device-service measurement-service
```

**2단계 — E2E 테스트 (아래 한 줄만 복사해서 실행):**
```bash
cd ~/Manpasik/backend && go test -v ./tests/e2e/...
```

또는 1단계 후 터미널에서 `cd backend` 입력한 다음, 새 줄에 `go test -v ./tests/e2e/...` 만 실행해도 됩니다.

## 서비스 연결이 스킵될 때 (context deadline exceeded)

헬스체크·플로우 테스트가 모두 "연결 불가", "연결 실패"로 **스킵**되면, Go 서비스 컨테이너가 떠 있지 않거나 `localhost:50051` 등에 연결되지 않는 상태입니다.

1. **컨테이너 확인** (프로젝트 루트에서):
   ```bash
   cd ~/Manpasik/infrastructure/docker && docker compose -f docker-compose.dev.yml ps
   ```
   `manpasik-auth-service`, `manpasik-user-service`, `manpasik-device-service`, `manpasik-measurement-service`가 **Up** 이어야 합니다.

2. **중지됐다면 다시 기동**:
   ```bash
   cd ~/Manpasik/infrastructure/docker && docker compose -f docker-compose.dev.yml up -d postgres auth-service user-service device-service measurement-service
   ```
   기동 후 수 초 뒤에 E2E를 다시 실행하세요.

3. **WSL 2 + Docker Desktop**: Docker Desktop에서 WSL 연동이 켜져 있어야 WSL 터미널의 `localhost`로 컨테이너 포트에 접근할 수 있습니다. Settings → Resources → WSL Integration에서 사용 중인 배포판을 활성화하세요.

4. **컨테이너는 Up인데 여전히 스킵되면**: E2E 기본 주소가 **127.0.0.1** 로 설정되어 있음 (WSL2에서 `localhost`가 IPv6(::1)로 해석되어 Docker 포트에 안 붙는 경우 대비). 캐시 없이 다시 실행: `cd ~/Manpasik/backend && go test -v -count=1 ./tests/e2e/...`. 포트 확인: `nc -zv 127.0.0.1 50051`.

## TestMeasurementFlow가 "want proto.Message" 로 스킵될 때

헬스체크는 통과하는데 플로우만 스킵되고, 로그에 `grpc: error while marshaling: want proto.Message` 가 보이면, shared gen이 **수동 작성**된 상태라 `proto.Message`(ProtoReflect)가 없기 때문입니다. **protoc로 Go 코드를 다시 생성**하면 해결됩니다.

1. `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` 설치 (필요 시 `apt install protobuf-compiler`, `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`, `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`)
2. 프로젝트 루트에서: `make proto` (Linux/WSL에서 Google 타입 경로가 다르면 `PROTO_GOOGLE_INCLUDE=/usr/include make proto`)
3. 다시 E2E 실행: `cd ~/Manpasik/backend && go test -v -count=1 ./tests/e2e/...`

자세한 내용은 `KNOWN_ISSUES.md` 의 "E2E TestMeasurementFlow: want proto.Message" 항목을 참고하세요.

## 참조

- 계획: `docs/plan/phase_1d_integration_mvp.md`
- 테스트 코드: `backend/tests/e2e/`
- 이슈: `KNOWN_ISSUES.md` (want proto.Message, Docker WSL 등)
