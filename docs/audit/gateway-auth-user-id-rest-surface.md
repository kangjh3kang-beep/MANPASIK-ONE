# Gateway Auth User ID REST Surface 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: Gateway REST 인증 응답의 사용자 식별자 전달 계약

## 상세 구축계획

1. Auth gRPC `LoginResponse.user_id`가 Gateway REST 응답까지 전달되는지 확인한다.
2. `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/social-login` E2E mock 응답에 `UserId`를 넣고 테스트가 값을 검증하게 한다.
3. 카카오 REST 로그인은 auth-service 응답의 `UserId`를 우선 사용하고, 빈 경우에만 기존 카카오 식별자를 fallback으로 쓴다.
4. Flutter REST repository는 Gateway 표준 snake_case와 legacy camelCase를 모두 읽도록 보강한다.
5. 로컬 HTTP 서버 기반 Flutter 테스트로 실제 REST JSON 매핑을 검증한다.
6. Go Gateway/Auth 회귀, Flutter Auth 테스트/analyze, targeted `git diff --check`를 통과시킨다.

## 구현 결과

- Gateway Auth E2E 테스트가 로그인, 갱신, 소셜 로그인 모두에서 `user_id`를 필수 계약으로 본다.
- 카카오 로그인 REST 응답은 `resp.UserId`를 우선 노출한다.
- Flutter `AuthRepositoryRest`는 `user_id`와 `userId` 양쪽을 처리하므로 JSON naming 설정 변화에도 사용자 ID를 유지한다.
- `unknown` fallback은 legacy 서버가 실제로 빈 사용자 ID를 반환하는 경우에만 남는다.

## 변경 파일

- `backend/services/gateway/internal/handler/auth_routes.go`
- `backend/services/gateway/internal/handler/e2e_test.go`
- `frontend/flutter-app/lib/features/auth/data/auth_repository_rest.dart`
- `frontend/flutter-app/test/features/auth/data/auth_repository_rest_test.dart`
- `backend/services/gateway/cmd/main_smoke_test.go`
- `backend/services/auth-service/cmd/main.go`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 \
  ./backend/services/gateway/cmd \
  ./backend/services/gateway/internal/handler \
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

- Go Gateway handler: PASS
- Go Gateway cmd + Gateway handler + Auth service/handler: PASS
- Flutter Auth REST/proto/domain tests: PASS
- Flutter Auth targeted analyze: PASS
- `git diff --check` targeted files: PASS
- 2026-05-02 Gateway auth live binary smoke: PASS

## 잔여 리스크

- 실제 Gateway binary와 auth-service binary를 동시에 띄운 live REST-to-gRPC smoke는 2026-05-02 후속 단계에서 완료됐다.
- Docker Compose 기반 container smoke는 현재 WSL Docker integration 부재로 blocked 상태다.

## 다음 단계

- Docker Desktop WSL integration 활성화 후 measurement compose smoke를 PASS로 확정한다.
- 이후 full dev compose에서 gateway + auth-service container auth smoke를 추가한다.
