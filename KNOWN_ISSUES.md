# ManPaSik 알려진 이슈 추적 (Known Issues Tracker)

> **용도**: 프로젝트의 모든 알려진 이슈, 환경 제약, 기술 부채, 우회 방법을 추적하는 문서
> **규칙**: 이슈 발견 시 추가, 해결 시 상태 변경 (삭제 금지 — 해결 이력도 지식이다)
> **업데이트**: 이슈 발견/해결 시 즉시

---

## 📋 이슈 상태 범례

| 상태 | 의미 |
|------|------|
| 🔴 미해결 | 해결되지 않은 활성 이슈 |
| 🟡 우회 중 | 임시 조치로 우회 중 (근본 해결 필요) |
| 🟢 해결됨 | 완전히 해결됨 |
| ⚪ 보류 | 현재 영향 없어 보류 중 |

---

## 🔴 미해결 이슈

### ~~공유 모듈 5개 중 4개 실서비스 미연동 (2026-02-12 식별)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: `shared/cache`, `shared/events`, `shared/search`, `shared/storage`, `shared/vectordb` 어댑터가 모두 구현되어 있으나 실서비스 미연동
- **해결**: Sprint 0 완료 — Redis(device/subscription), Kafka(measurement), Milvus(measurement), Elasticsearch(measurement/community), MinIO(gateway) 전부 연동
- **상태**: 🟢 해결됨

### ~~community/video/translation/telemedicine 서비스 PostgreSQL 미연동 (2026-02-12 식별)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: 4개 서비스가 인메모리 저장소만 사용, 서비스 재시작 시 데이터 유실
- **해결**: 4개 서비스 모두 PostgreSQL 저장소 구현 + 조건부 초기화(DB_HOST) 적용
- **상태**: 🟢 해결됨
- **해결 계획**: Sprint 0에서 PostgreSQL Repository 구현 및 연동
- **우선순위**: P1

### Flutter 단위 테스트 (2026-02-12 검증)
- **이전**: 테스트 파일 부재 또는 0개로 기재됨
- **현재**: `flutter test` 실행 시 exit 0 통과 (2026-02-12 통합 검증). 테스트 파일 5개 존재. 60개+ 목표는 Sprint 1에서 지속 보강
- **우선순위**: P1 (추가 커버리지 확대)

### Rust FFI 브리지 비활성화 (2026-02-12 식별)
- **증상**: `frontend/flutter-app/lib/main.dart` 18행에서 `await RustBridge.init()` 주석 처리
- **영향**: Flutter 앱에서 Rust 코어 엔진(차동측정, 핑거프린트, AI) 미사용
- **해결 계획**: Sprint 1에서 flutter_rust_bridge 빌드 설정 후 활성화
- **우선순위**: P0

### WSL 2: Docker 명령을 찾을 수 없음 (`docker` / `docker-compose`)
- **증상**: `The command 'docker' could not be found in this WSL 2 distro` 또는 `docker-compose` could not be found
- **원인**: Docker Desktop이 설치되어 있어도 **WSL 2 배포판과 연동이 꺼져 있으면** 해당 WSL 터미널에서 `docker` 실행 파일을 찾지 못함
- **해결 (권장)**  
  1. Windows에서 **Docker Desktop** 실행  
  2. **Settings** → **Resources** → **WSL Integration**  
  3. **Enable integration with my default WSL distro** 켜기  
  4. 사용하는 배포판(예: Ubuntu) 옆 **Enable** 켜기  
  5. **Apply & Restart** 후 WSL 터미널 새로 열고 `docker --version` 확인  
- **참고**: [Docker Desktop WSL 2 백엔드](https://docs.docker.com/go/wsl2/)  
- **E2E 테스트**: Docker 없이도 `cd backend && go test -v ./tests/e2e/...` 는 실행 가능하며, 서비스 미기동 시 헬스/플로우 테스트는 스킵되고 `TestDifferentialMeasurement` 등 단위 테스트만 통과함

### E2E TestMeasurementFlow: `grpc: want proto.Message` (marshal / unmarshal)
- **증상**: TestServiceHealth는 통과하지만 TestMeasurementFlow에서 다음 중 하나로 스킵됨  
  - `grpc: error while marshaling: ... *v1.RegisterRequest, want proto.Message` (클라이언트)  
  - `grpc: error unmarshalling request: ... *v1.RegisterRequest, want proto.Message` (서버)
- **원인**: `backend/shared/gen/go/v1/manpasik.pb.go` 가 **수동 스텁**이면 `proto.Message`(ProtoReflect) 미구현. 클라이언트는 marshal, **실행 중인 서비스**는 unmarshal 시 실패함.
- **해결**  
  1. **protoc로 Go 코드 재생성**: `make proto` (프로젝트 루트). 필요 시 `apt install protobuf-compiler`, `go install .../protoc-gen-go@latest`, `go install .../protoc-gen-go-grpc@latest`, `PROTO_GOOGLE_INCLUDE=/usr/include make proto`.  
  2. **의존성**: `make proto` 후 빌드 오류(`SupportPackageIsVersion9` 등) 시 **backend/go.mod** 에서 `google.golang.org/grpc` v1.78.0 이상, `google.golang.org/protobuf` v1.35.2 이상으로 올린 뒤 `cd backend && go mod tidy`.  
  3. **중요**: `make proto` 후 **모든 gRPC 서비스(auth, user, device, measurement)를 재빌드·재기동**해야 함. 기동 중인 바이너리가 예전(스텁) 코드로 빌드되어 있으면 서버 쪽 unmarshal 오류가 난다.  
     - 로컬: `make build-go` 후 서비스 재실행.  
     - Docker: 이미지 재빌드 후 `docker compose up -d` 등으로 재기동.
- **이후**: `cd backend && go test -v -count=1 ./tests/e2e/...` 로 플로우 통과 확인.

### Docker Compose: `pull access denied for manpasik/auth-service` 등 (정상 동작)
- **증상**: `docker compose up` 시 manpasik/auth-service, manpasik/user-service 등에 대해 "pull access denied" 메시지가 반복 출력됨
- **원인**: 해당 이미지는 Docker Hub 등에 푸시되어 있지 않아 pull은 실패함
- **동작**: Compose는 pull 실패 후 **로컬 Dockerfile로 빌드**하여 이미지를 만들고 컨테이너를 기동함. "Building ... FINISHED", "Image manpasik/xxx:dev Built", "Container manpasik-xxx Created" 로그가 나오면 정상임. 별도 `docker login` 불필요

---

## 🟢 해결된 이슈

### ISSUE-001: TFLite 네이티브 빌드 불가 (Bazel 미설치) ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **영향**: `cargo build --features full` 실패
- **증상**: `tflitec` 크레이트 빌드 시 "Cannot find bazel" 에러 (v0.5.2) / spectrogram.cc 컴파일 오류·bindgen Invalid Ident (v0.5 소스 빌드 시)
- **해결**: **tflitec 0.5 → 0.7 업그레이드**. v0.7은 Bazel 없이 빌드 가능(bindgen 0.65 사용). `cargo build -p manpasik-engine --features full` 및 `cargo test -p manpasik-engine --features full` 62테스트 통과.
- **관련 파일**: `rust-core/Cargo.toml` (tflitec = "0.7"), `rust-core/manpasik-engine/Cargo.toml`
- **참고**: Bazel/Bazelisk가 필요한 구버전(v0.5) 소스 빌드를 쓰는 경우, WSL에 Bazelisk 설치(`~/.local/bin/bazelisk`) 및 TensorFlow spectrogram.cc에 `#include <cstdint>` 패치 필요.

---

### ISSUE-002: WSL 셸 명령어 출력 미반환 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `wsl -d Ubuntu -- bash -c "..."` 명령이 exit code 0이나 출력 없음 (0ms 완료)
- **원인**: Cursor IDE의 WSL 셸 세션 상태 이상
- **해결**: Windows 네이티브 명령(`hostname`)을 로컬 working_directory에서 실행하여 셸 리셋
- **교훈**: WSL 명령 출력이 비정상일 때 Windows 명령으로 셸 상태 초기화 시도

### ISSUE-003: Rust 툴체인 미설치 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `which rustc` 및 `cargo --version` 실패
- **해결**: `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y`
- **교훈**: WSL 환경에 Rust가 기본 설치되어 있지 않음, 매번 확인 필요

### ISSUE-004: Cargo 빌드 — 벤치마크 파일 누락 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `error: failed to parse manifest` — `benches/differential_measurement.rs` 파일 없음
- **원인**: `Cargo.toml`에 `[[bench]]` 항목이 선언되었으나 실제 파일 미존재
- **해결**: `rust-core/manpasik-engine/benches/differential_measurement.rs` 더미 벤치마크 파일 생성
- **관련 파일**: `rust-core/manpasik-engine/Cargo.toml`, `rust-core/manpasik-engine/benches/`

### ISSUE-005: OpenSSL 시스템 라이브러리 누락 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `Could not find OpenSSL installation` (`openssl-sys` 크레이트 빌드 실패)
- **원인**: `btleplug` (BLE) → `openssl-sys` 의존, WSL에 `libssl-dev` 미설치
- **해결**: `sudo apt-get install -y libssl-dev libdbus-1-dev pkg-config build-essential libclang-dev cmake`
- **교훈**: BLE 기능 빌드에 시스템 라이브러리 여러 개 필요

### ISSUE-006: `hex` 크레이트 미등록 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `error[E0433]: failed to resolve: use of unresolved module or unlinked crate 'hex'` (`nfc/mod.rs`)
- **원인**: NFC 모듈에서 `hex` 크레이트 사용하나 `Cargo.toml`에 미등록
- **해결**: `Cargo.toml`에 `hex = "0.4"` 추가
- **관련 파일**: `rust-core/manpasik-engine/Cargo.toml`

### ISSUE-007: `futures` 크레이트 미등록 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `error[E0433]: failed to resolve: use of unresolved module or unlinked crate 'futures'` (`ble/mod.rs`)
- **원인**: BLE 모듈에서 `futures` 크레이트의 비동기 스트림 사용하나 `Cargo.toml`에 미등록
- **해결**: `Cargo.toml`에 `futures = "0.3"` 추가
- **관련 파일**: `rust-core/manpasik-engine/Cargo.toml`

### ISSUE-008: `unused_mut` 경고 (DSP) ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `warning: variable does not need to be mutable` (`dsp/mod.rs`의 `sum` 변수)
- **해결**: `let mut sum` → `let sum` 변경
- **관련 파일**: `rust-core/manpasik-engine/src/dsp/mod.rs`

### ISSUE-009: `sudo` 비대화형 셸에서 패스워드 프롬프트 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `sudo apt-get install` 명령이 무한 대기 (exit하지 않음)
- **원인**: WSL 비대화형 셸에서 sudo 패스워드 입력 불가
- **우회**: 사용자에게 수동 실행 요청
- **교훈**: sudo 필요 명령은 사용자에게 직접 실행 요청 필요

### ISSUE-010: Flutter intl 버전 충돌 (flutter_localizations 핀) ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `flutter pub get` 실패 — "Because manpasik depends on flutter_localizations from sdk which depends on intl 0.20.2, intl ^0.19.0 version solving failed."
- **원인**: `flutter_localizations`가 `intl 0.20.2`를 고정(pin)하는데, `pubspec.yaml`에 `intl: ^0.19.0`으로 선언
- **해결**: `pubspec.yaml`에서 `intl: ^0.19.0` → `intl: ^0.20.2`로 변경
- **교훈**: **`intl` 패키지는 Flutter SDK에 의해 버전이 고정됨.** `flutter_localizations`를 사용할 경우 `intl` 버전을 SDK 핀 버전과 일치시켜야 함. `flutter pub add intl:^0.20.2` 명령으로 확인 가능.
- **방지책**: Flutter 프로젝트 초기 설정 시 `flutter pub get`을 먼저 실행하여 SDK 핀 버전 확인 후 `pubspec.yaml` 작성

### ISSUE-011: flutter_gen 합성 패키지 미생성 (gen-l10n 비동작) ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `flutter analyze` 에러 — "Target of URI doesn't exist: 'package:flutter_gen/gen_l10n/app_localizations.dart'"
- **원인**: `pubspec.yaml`에 `generate: true` + `l10n.yaml` 설정했으나, `flutter gen-l10n`이 `.dart_tool/flutter_gen/` 합성 패키지를 제대로 생성하지 못함. `synthetic-package` 옵션도 deprecated.
- **시도한 우회**: (1) `synthetic-package: false` + `output-dir` → "no longer has any effect" 에러 (2) `generate: true` 제거 → "generate flag turned on" 요구 에러
- **최종 해결**: **flutter_gen/gen-l10n 코드 생성 포기 → 수동 AppLocalizations 구현.** `lib/l10n/app_localizations.dart`에 직접 delegate/class 작성, `lib/l10n/translations/{ko,en,ja,zh,fr,hi}.dart`에 Map<String,String> 기반 번역.
- **교훈**: **Flutter의 `generate: true` + `flutter_gen` 합성 패키지는 환경에 따라 불안정.** WSL + Cursor IDE 조합에서 특히 문제 발생. 수동 구현이 100% 안정적이며 코드 생성 의존성 제거. ARB 파일은 참조용으로 보존.
- **방지책**: (1) Flutter l10n은 수동 구현 우선 (2) gen-l10n 사용 시 반드시 생성 파일 존재 확인 (3) `flutter_gen` import 사용 전 `.dart_tool/flutter_gen/` 디렉토리 확인

### ISSUE-012: gen-l10n 잔여 생성 파일 충돌 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `flutter analyze` 에러 — "The argument type 'String' can't be assigned to the parameter type 'Locale'" in `app_localizations_{ko,en,ja,zh,fr,hi}.dart`
- **원인**: 이전 `flutter gen-l10n` 실행이 `lib/l10n/` 디렉토리에 `app_localizations_*.dart` 파일을 생성. 수동 구현으로 전환 후 이 파일들이 잔존하여 충돌.
- **해결**: `lib/l10n/app_localizations_{ko,en,ja,zh,fr,hi}.dart` 6개 파일 삭제
- **교훈**: **gen-l10n → 수동 구현 전환 시 반드시 자동 생성 잔여 파일 확인 및 삭제.** gen-l10n은 ARB 파일당 `app_localizations_{locale}.dart` 파일을 생성하므로 이들이 수동 구현과 충돌.
- **방지책**: l10n 전환 시 `ls lib/l10n/app_localizations_*.dart` 확인 후 삭제

### ISSUE-013: Flutter deprecated API 일괄 수정 ✅
- **발견일**: 2026-02-10
- **해결일**: 2026-02-10
- **발견자**: Claude
- **증상**: `flutter analyze` info 경고 다수 — `withOpacity` deprecated, `debugState` deprecated, `RadioListTile.groupValue/onChanged` deprecated
- **원인**: Flutter 3.32+ 에서 다수 API deprecated
- **해결**:
  - `Color.withOpacity(0.3)` → `Color.withValues(alpha: 0.3)` (6개 파일)
  - `StateNotifier.debugState` → `StateNotifier.state` (테스트 파일)
  - `RadioListTile(groupValue:, onChanged:)` → `ListTile` + 체크 아이콘 (settings_screen)
  - `const` 추가 (analysis hints)
- **교훈**: **Flutter 최신 버전 사용 시 deprecated API 사전 확인 필수.** 특히 `withOpacity`→`withValues`, `debugState`→`state`, `RadioListTile`→`RadioGroup` 전환.
- **방지책**: 코드 작성 시 `withValues(alpha:)` 사용, 테스트에서 `state` 직접 접근, RadioListTile 대신 ListTile+아이콘 패턴 사용

---

## 🟡 경고/주의사항

### WARN-001: 빌드 경고 — 미사용 import (BLE)
- **파일**: `rust-core/manpasik-engine/src/ble/mod.rs`
- **내용**: `unused imports: Adapter and Peripheral`
- **영향**: 없음 (향후 실제 BLE 연결 구현 시 사용 예정)
- **조치**: 보류 (향후 BLE 완전 구현 시 자연 해결)

### WARN-002: 빌드 경고 — 미사용 구조체 (Sync)
- **파일**: `rust-core/manpasik-engine/src/sync/mod.rs`
- **내용**: `struct TaggedElement is never constructed`
- **영향**: 없음 (CRDT 확장 시 사용 예정)
- **조치**: 보류 (향후 Sync 모듈 확장 시 자연 해결)

### WARN-003: AI 모듈 `unused_mut` 경고
- **파일**: `rust-core/manpasik-engine/src/ai/mod.rs`
- **내용**: `variable does not need to be mutable` (anomaly_input)
- **영향**: 없음
- **조치**: 다음 AI 모듈 작업 시 수정

---

## 📊 환경 요구사항 (WSL Ubuntu)

이 프로젝트를 WSL 환경에서 빌드하기 위해 필요한 시스템 패키지:

```bash
# Rust 툴체인
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
source ~/.cargo/env

# 시스템 빌드 의존성 (BLE + OpenSSL + D-Bus)
sudo apt-get update
sudo apt-get install -y \
  libssl-dev \
  libdbus-1-dev \
  pkg-config \
  build-essential \
  libclang-dev \
  cmake

# Python (tflitec 빌드 시 필요)
pip3 install --break-system-packages numpy

# TFLite 전체 빌드 시 (선택)
# Bazel/Bazelisk 설치 필요
```

**빌드 명령어:**
```bash
# AI 제외 빌드 (권장 — Bazel 불필요)
cargo build -p manpasik-engine --no-default-features --features 'std,ble,nfc,fingerprint'

# 전체 빌드 (Bazel 필요)
cargo build -p manpasik-engine --features full

# 테스트
cargo test -p manpasik-engine --no-default-features --features 'std,ble,nfc,fingerprint'
```

---

---

## 📋 Flutter 에러 방지 체크리스트

새 Flutter 코드 작성 시 반드시 확인:

1. **intl 버전**: `flutter_localizations` 사용 시 `intl` 버전을 SDK 핀 버전과 일치 (`^0.20.2`)
2. **l10n 구현**: `flutter_gen` 합성 패키지 대신 수동 `AppLocalizations` 사용 (환경 안정성)
3. **Color API**: `withOpacity()` 사용 금지 → `withValues(alpha:)` 사용
4. **StateNotifier 테스트**: `debugState` 사용 금지 → `state` 직접 접근
5. **RadioListTile**: `groupValue`/`onChanged` deprecated → `ListTile` + 체크 아이콘 패턴
6. **gen-l10n 전환 시**: 잔여 `app_localizations_*.dart` 파일 삭제 필수
7. **import 정리**: 미사용 import는 즉시 제거 (`flutter/material.dart` 등)
8. **const 적극 활용**: `const` 가능한 위젯/리터럴에 항상 적용

---

### ~~admin-service ConfigMetadata/Translation PostgreSQL 미구현 (2026-02-12 식별)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: ConfigMetadataRepository, ConfigTranslationRepository가 인메모리만 구현됨.
- **해결**: `postgres/config_meta.go` 신규 구현 — ConfigMetadataRepository(GetByKey, ListByCategory, ListAll, CountByCategory) + ConfigTranslationRepository(GetByKeyAndLang, ListByKey, ListByLang). main.go에서 DB 연결 시 자동 전환.
- **검증**: `go build` + `go vet` + `go test` 전체 통과. 통합 테스트 9개 작성 (DB 미접속 시 Skip).
- **상태**: 🟢 해결됨

### ~~Proto 생성 코드 수동 추가 (2026-02-12 식별)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: `admin_config_ext.go` + `telemedicine_ext.go`에 수동 타입 정의. protoc 재생성 시 충돌 가능.
- **해결**: `manpasik.proto`에 TelemedicineService 추가 → `protoc` 정식 재생성 → `admin_config_ext.go`, `telemedicine_ext.go` 삭제. `go build/vet/test` 전체 통과.
- **상태**: 🟢 해결됨

### config_manager_test.go 순환 import 수정 (2026-02-12 식별·해결)
- **증상**: `package service` 내부 테스트에서 `memory` 패키지 import → 순환 의존성
- **해결**: `package service_test` (외부 테스트 패키지)로 변경, EventBus 직접 참조
- **상태**: 🟢 해결됨

### notification/payment-service undefined ctx (2026-02-12 식별·해결)
- **증상**: ConfigWatcher.Watch(ctx, ...) 호출 시점에 ctx 미정의 → `go build ./...` 실패
- **해결**: `context.Background()` 사용으로 변경
- **상태**: 🟢 해결됨

### ~~telemedicine-service Proto 타입 미생성 (2026-02-12 식별·우회)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: `v1.UnimplementedTelemedicineServiceServer` 등 미정의 → `go build ./...` 실패
- **해결**: `manpasik.proto`에 TelemedicineService 정의 추가 → `protoc` 정식 재생성 → `telemedicine_ext.go` 삭제. `go build/vet/test` 전체 통과.
- **상태**: 🟢 해결됨

### E2E env.go 포트 매핑 오류 (2026-02-12 식별·해결)
- **증상**: AdminAddr()=50067, NotificationAddr()=50068이었으나 실제 서비스 포트와 불일치
- **원인**: 초기 설정 시 포트 할당 변경이 env.go에 미반영
- **해결**: AdminAddr→50068, NotificationAddr→50062, K8s ConfigMap도 동일 수정
- **상태**: 🟢 해결됨

### ~~screen_widget_test.dart gRPC Provider 의존성 오류 (2026-02-12 식별)~~ → 🟢 해결됨 (2026-02-12)
- **증상**: `screen_widget_test.dart`의 HomeScreen/DeviceListScreen/MeasurementResultScreen 위젯 테스트 4건 FAIL
- **원인**: gRPC Provider 의존성(grpcChannelProvider 등) override 미제공 → 위젯 초기화 시 네트워크 접근 시도
- **해결**: `_baseOverrides()`에 `FakeMeasurementRepository`, `FakeDeviceRepository`, `FakeUserRepository`, `measurementHistoryProvider`, `deviceListProvider` override 추가. 6/6 전체 PASS.
- **상태**: 🟢 해결됨

### Flutter widget_test.dart / screen_widget_test.dart import 오류 (2026-02-12 식별·해결)
- **증상**: `package:manpasik/test/helpers/fake_repositories.dart` — 테스트 헬퍼를 package import로 참조 → URI 미존재
- **원인**: 테스트 파일은 `test/` 디렉토리에 있으므로 상대 경로 import 필요
- **해결**: `import 'helpers/fake_repositories.dart'`로 변경
- **상태**: 🟢 해결됨

### screen_widget_test.dart getter 문법 오류 (2026-02-12 식별·해결)
- **증상**: `List<Override> get _baseOverrides => [...]` — 함수 내부에서 getter 구문 사용 불가
- **해결**: `List<Override> _baseOverrides() => [...]` 일반 함수로 변환, 호출부 `_baseOverrides()` 변경
- **상태**: 🟢 해결됨

### AppTheme 테스트 Google Fonts 네트워크 오류 (2026-02-12 식별·해결)
- **증상**: `app_theme_test.dart` dark 테마 테스트 시 `GoogleFonts.notoSansKr()` HTTP 요청 실패
- **원인**: 테스트 환경에서 네트워크 접근 불가 + Google Fonts 기본 런타임 fetching 활성화
- **해결**: `setUpAll(() { GoogleFonts.config.allowRuntimeFetching = false; })` 추가
- **상태**: 🟢 해결됨

### E2E payment_subscription_flow_test.go Proto 필드 불일치 (2026-02-12 식별·해결)
- **증상**: `Amount`, `Currency`, `OrderName` 필드 미존재, `GetSubscriptionRequest` → `GetSubscriptionDetailRequest`, `UpgradeSubscription` → `UpdateSubscription`
- **원인**: protoc 정식 재생성 후 Proto 필드명이 변경됨 (Amount→AmountKrw 등)
- **해결**: E2E 테스트 코드를 Proto 정의에 맞게 수정
- **상태**: 🟢 해결됨

**마지막 업데이트**: 2026-02-12 (Sprint 2 Phase 2 — D-2 SRS, D-3 SAD, AS-7 LLM, I-5 K8s Overlay, screen_widget_test 수정. 🟡 우회 0건. 🔴 미해결 0건.)
