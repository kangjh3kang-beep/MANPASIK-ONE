# Flutter Auth Active User Propagation 게이트

**작성일**: 2026-05-02  
**작성자**: Codex  
**범위**: Flutter 로그인 사용자 식별자의 `TenantInterceptor` active user 저장 및 `X-User-ID` 헤더 기반 연결

## 상세 구축계획

1. Auth gRPC/REST에서 확보한 `user_id`가 `AuthState.userId`에만 머무르는지 확인한다.
2. 로그인, 회원가입, 소셜 로그인 성공 후 `TenantInterceptor.setActiveUser(userId)`를 호출한다.
3. 게스트/데모 로그인도 동일하게 active user를 저장해 데모/개발 경로의 헤더 동작을 일관화한다.
4. 로그아웃과 인증 상태 실패 시 `TenantInterceptor.clear()`로 active user와 active tenant를 함께 제거한다.
5. 소셜 로그인 직접 REST 경로는 snake_case와 camelCase 응답을 모두 읽도록 보강한다.
6. `AuthNotifier` 테스트에서 active user 저장/해제를 검증한다.
7. 기존 `TenantInterceptor` Dio 통합 테스트로 `X-User-ID` 헤더 주입 계약을 재확인한다.

## 구현 결과

- `AuthNotifier.login()`과 `register()`가 성공하면 `active_user_id`에 실제 `userId`를 저장한다.
- `socialLogin()`은 `user_id/userId`를 모두 읽고 active user에 저장한다.
- `loginAsGuest()`와 `loginAsDemo()`도 active user를 저장한다.
- `logout()`과 미인증 `checkAuthStatus()`는 active user와 active tenant를 제거한다.
- SharedPreferences 초기화 실패는 인증 상태 자체를 깨지 않도록 격리했다.

## 변경 파일

- `frontend/flutter-app/lib/shared/providers/auth_provider.dart`
- `frontend/flutter-app/test/shared/providers/auth_notifier_test.dart`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
flutter test --no-pub \
  test/shared/providers/auth_notifier_test.dart \
  test/core/network/tenant_interceptor_test.dart \
  test/features/auth/data/auth_repository_rest_test.dart \
  test/features/auth/data/auth_proto_contract_test.dart \
  test/features/auth/domain/auth_result_test.dart

flutter analyze --no-pub \
  lib/shared/providers/auth_provider.dart \
  lib/core/network/tenant_interceptor.dart \
  test/shared/providers/auth_notifier_test.dart \
  test/core/network/tenant_interceptor_test.dart

flutter analyze --no-pub \
  lib/shared/providers/auth_provider.dart \
  lib/features/auth/presentation/login_screen.dart \
  lib/features/settings/presentation/settings_screen.dart \
  test/shared/providers/auth_notifier_test.dart \
  test/widget_test.dart \
  test/core/providers/measurement_history_provider_test.dart
```

결과:

- Flutter Auth/Tenant tests: PASS
- Flutter AuthProvider/Tenant targeted analyze: PASS
- Flutter AuthProvider callsite analyze: PASS
- Go Gateway/Auth regression: PASS
- `git diff --check` targeted files: PASS

## 잔여 리스크

- 이번 단계는 Flutter 저장소와 Dio interceptor 헤더 주입 기반을 검증한다.
- 실제 Gateway live smoke에서 `X-User-ID`가 gRPC metadata로 전파되는지는 별도 backend header propagation smoke가 필요하다.
- Docker Compose container smoke는 현재 WSL Docker integration 부재로 runtime blocked 상태다.

## 다음 단계

- Gateway live smoke에 `X-User-ID`/tenant header propagation 검증을 추가한다.
- Docker Desktop WSL integration 활성화 후 gateway/auth 및 measurement compose smoke를 실제 PASS로 확정한다.
