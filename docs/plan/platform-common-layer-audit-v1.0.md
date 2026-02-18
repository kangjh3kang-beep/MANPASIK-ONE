# ManPaSik 플랫폼 공통층 구현 감사 보고서 v1.0

**작성일**: 2026-02-18
**대상**: `frontend/flutter-app/lib/core/`, `lib/shared/`, `lib/l10n/`, `lib/main.dart`

---

## 1. 전체 평가

| 구분 | 평가 |
|------|------|
| **코드 레벨 완성도** | 92% — 모든 기능 모듈 존재, 빌드/분석 통과 |
| **프로덕션 준비도** | 75% — 운영 필수 인프라(에러 핸들링, 로깅, 크래시 리포팅) 미비 |
| **외부 SDK 연동** | 40% — PG/본인인증/FCM/HealthKit 모두 시뮬레이션 모드 |

---

## 2. 완벽 구현 (100%) — 8개 영역

| 영역 | 파일 | 설명 |
|------|------|------|
| 환경 설정 | `core/config/app_config.dart` | dev/staging/prod 3환경 분리, feature flags, SSL pins |
| SSL Pinning | `core/network/ssl_pinning.dart` | X509 SHA-256 검증, release-only 동작 |
| 인증 인터셉터 | `core/network/auth_interceptor.dart` | JWT 자동갱신(401), refresh token 로직 |
| 오프라인 큐 | `core/network/offline_queue.dart` | Hive FIFO, 재시도 추적, max 3회 |
| 라우터 | `core/router/app_router.dart` | 72+ 라우트, ShellRoute 5탭, Admin RBAC guard |
| REST 클라이언트 | `core/services/rest_client.dart` | 90+ API 메서드, multipart upload |
| 테마 | `core/theme/app_theme.dart` | light/dark/koreanWhite 3종, 전통 색상 팔레트 |
| 다국어 | `l10n/` | 6개 언어(ko/en/ja/zh/fr/hi), ~90키 |

## 3. 실질 완성 (95~98%) — 2개 영역

| 영역 | 파일 | 완성도 | 미비점 |
|------|------|--------|--------|
| gRPC Provider | `core/providers/grpc_provider.dart` | 98% | 12 Repository + 21 FutureProvider. 일부 서비스 연결 미사용 |
| 공유 Provider | `shared/providers/` | 95% | 7개 프로바이더. chat_provider 스트리밍 시뮬레이션, sync_provider backoff 없음 |

## 4. 부분 구현 (35~70%) — 4개 영역

| 영역 | 파일 | 완성도 | 설명 |
|------|------|--------|------|
| Rust FFI | `core/services/rust_ffi_stub.dart` | 70% | 전체 시뮬레이션 동작, flutter_rust_bridge 연결점 주석 처리 |
| 결제(PG) | `core/services/payment_service.dart` | 40% | 추상 인터페이스 + SimulatedPaymentService |
| 본인인증 | `core/services/identity_verification_service.dart` | 35% | 추상 인터페이스 + 시뮬레이션 |
| HealthConnect | `core/services/health_connect_service.dart` | 60% | 플랫폼 채널 + 시뮬레이션 fallback |

## 5. 미구현 (0%) — 프로덕션 필수 6건

| # | 항목 | 위치 | 영향도 |
|---|------|------|--------|
| 1 | 글로벌 에러 핸들러 | `main.dart` — FlutterError.onError 미설정 | **높음** |
| 2 | 로깅/크래시 리포팅 | Crashlytics/Sentry 미연동 | **높음** |
| 3 | 애널리틱스 | Firebase Analytics/Mixpanel 미설정 | **중간** |
| 4 | FCM 푸시 | Firebase 초기화 주석 처리, 폴링으로 대체 | **중간** |
| 5 | 딥링크 핸들링 | GoRouter initialLocation 고정, 딥링크 스킴 미등록 | **중간** |
| 6 | 앱 라이프사이클 관리 | WidgetsBindingObserver 미구현 | **낮음** |

## 6. 구조적 이슈

### Feature 계층 불균일
- **Clean Architecture 준수 (9/17)**: auth, community, data_hub, devices, family, market, measurement, medical, ai_coach
- **domain 계층 부재 (8/17)**: admin, chat, home, notification, onboarding, settings, user

### 미사용 위젯 (6개)
- `cached_image.dart`, `cartridge_3d_viewer.dart`, `glass_dock_navigation.dart`
- `leaderboard_widget.dart`, `network_indicator.dart`, `sanggam_decoration.dart`

### Provider 개선 필요
- `chat_provider.dart`: sendMessageStream이 char-by-char 시뮬레이션 (실 SSE 아님)
- `sync_provider.dart`: 5분 주기 동기화, exponential backoff 없음

---

## 7. Sprint 10 구현 계획 요약

### Phase 1: 글로벌 에러 핸들러 + Zone.runGuarded
- `main.dart` — runZonedGuarded 래핑, FlutterError.onError 설정
- `core/services/crash_reporter.dart` (신규) — 추상 CrashReporter + ConsoleCrashReporter

### Phase 2: 로깅/크래시 리포팅 인프라
- `core/services/app_logger.dart` (신규) — 레벨별 로깅(debug/info/warning/error)
- CrashReporter와 AppLogger 통합, main.dart에서 초기화

### Phase 3: 앱 라이프사이클 관리
- `core/services/app_lifecycle_observer.dart` (신규) — WidgetsBindingObserver
- 백그라운드 전환 시 WebSocket 해제, 포그라운드 복귀 시 토큰 갱신

### Phase 4: 딥링크 핸들링
- `app_router.dart` 수정 — initialLocation을 딥링크 파라미터에서 읽기
- AndroidManifest.xml + Info.plist 딥링크 스킴 등록

### Phase 5: Feature domain 계층 보완 (8개)
- admin, chat, home, notification, onboarding, settings, user — 각 domain/ 디렉토리 + repository 인터페이스

### Phase 6: 미사용 위젯 정리 + Provider 개선
- 미사용 위젯 6개를 실 화면에 연결하거나 삭제
- chat_provider SSE 스트리밍 구조 개선
- sync_provider exponential backoff 추가

### 검증
- `flutter analyze` — 에러 0건 유지
- `go build` — 11/11 서비스 통과
