# Auth Login User Identity Surface 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: Auth 로그인/토큰 갱신 응답의 사용자 식별자 계약 복구

## 상세 구현계획

1. `LoginResponse`에 하위 호환 필드 `user_id = 5`를 추가한다.
2. Go generated proto와 Flutter official generated proto를 같은 원본 proto에서 재생성한다.
3. auth-service 내부 `TokenPair`가 토큰 문자열뿐 아니라 사용자 ID를 함께 운반하게 한다.
4. gRPC handler의 `Login`, `RefreshToken`, `SocialLogin` 응답에 `UserId`를 채운다.
5. Flutter `AuthRepositoryImpl`은 정상 응답에서 `res.userId`를 사용하고, legacy 빈 응답에서만 `unknown` fallback을 유지한다.
6. 서비스, handler, Dart proto 계약 테스트로 회귀를 잠근다.
7. 기존 Measure golden path 테스트를 함께 실행해 proto 재생성이 스트림 계약을 깨지 않는지 확인한다.

## 구현 결과

- `LoginResponse.user_id = 5`가 추가되어 로그인 응답만으로 앱 세션의 사용자 식별자를 확보할 수 있다.
- `generateTokenPair()`가 모든 토큰 발급 경로에서 `TokenPair.UserID`를 채운다.
- `Login`, `RefreshToken`, `SocialLogin` gRPC 응답이 동일한 `UserId` 매핑을 사용한다.
- Flutter AuthRepository는 정상 서버에서 `unknown` 대신 실제 사용자 ID를 세션 결과로 반환한다.
- 공식 Dart generated 산출물은 현재 proto 기준 64,750라인으로 재생성됐다.
- 2026-05-01 후속 Gateway REST 보강으로 `/auth/login`, `/auth/refresh`, `/auth/social-login` E2E도 `user_id` 전달을 검증한다.

## 변경 파일

- `backend/shared/proto/manpasik.proto`
- `backend/shared/gen/go/v1/manpasik.pb.go`
- `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
- `backend/services/auth-service/internal/service/auth.go`
- `backend/services/auth-service/internal/service/auth_test.go`
- `backend/services/auth-service/internal/handler/grpc.go`
- `backend/services/auth-service/internal/handler/grpc_test.go`
- `frontend/flutter-app/lib/generated/manpasik.pb.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`
- `frontend/flutter-app/lib/features/auth/data/auth_repository_impl.dart`
- `frontend/flutter-app/test/features/auth/data/auth_proto_contract_test.dart`
- `backend/services/gateway/internal/handler/auth_routes.go`
- `backend/services/gateway/internal/handler/e2e_test.go`
- `frontend/flutter-app/lib/features/auth/data/auth_repository_rest.dart`
- `frontend/flutter-app/test/features/auth/data/auth_repository_rest_test.dart`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/auth-service/internal/service \
  ./backend/services/auth-service/internal/handler \
  ./backend/services/measurement-service/cmd \
  ./backend/services/measurement-service/internal/handler \
  ./backend/services/measurement-service/internal/service

cd /home/kangjh3kang/Manpasik/frontend/flutter-app
scripts/check_proto_generation_preflight.sh
scripts/check_proto_generation_compile_gate.sh
flutter test --no-pub \
  test/features/auth/data/auth_proto_contract_test.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart \
  test/features/measurement/data/measurement_repository_rest_test.dart \
  test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
flutter analyze --no-pub \
  lib/generated/manpasik.pb.dart \
  lib/generated/manpasik.pbgrpc.dart \
  lib/generated/manpasik.pbenum.dart \
  lib/generated/manpasik.pbjson.dart \
  lib/features/auth/data/auth_repository_impl.dart \
  test/features/auth/data/auth_proto_contract_test.dart \
  test/features/measurement/data/measurement_stream_wire_contract_test.dart \
  test/features/measurement/data/measurement_repository_impl_test.dart
```

결과:

- Go auth service/handler: PASS
- Go auth + measurement cmd/handler/service regression: PASS
- Dart proto preflight: PASS, generated Dart 64,750 lines
- Dart proto compile gate: PASS, `grpc 5.1.0`, `protobuf 6.0.0`
- Flutter auth + measure tests: PASS
- Flutter generated/auth analyze target: PASS
- Gateway REST auth E2E user_id regression: PASS
- Flutter REST auth snake_case/camelCase mapping tests: PASS
- `git diff --check` targeted files: PASS

## 잔여 리스크

- Gateway REST 로그인/갱신/소셜 로그인 응답의 `user_id` 전달은 2026-05-01 후속 단계에서 검증 완료됐다.
- 실제 Gateway binary와 auth-service binary를 함께 띄운 live REST-to-gRPC smoke는 아직 별도 단계가 필요하다.
- legacy 서버나 중간 gateway가 아직 `user_id`를 비워 전달하면 Flutter fallback `unknown`이 유지된다.
- Compose container smoke는 Docker Desktop WSL integration 부재로 아직 runtime blocked 상태다.

## 다음 단계

- Gateway auth live gRPC smoke를 추가해 별도 프로세스 경계에서 REST 로그인 `user_id` 전달을 검증한다.
- Docker Desktop WSL integration 활성화 후 measurement compose smoke를 실제 PASS로 확정한다.
