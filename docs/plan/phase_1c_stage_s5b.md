# Phase 1C Stage S5b: Rust FFI·BLE/NFC·차트·통합 테스트

> **전제**: S5a 완료 (gRPC 연동, Repository, 화면 고도화, 60+ 테스트)
> **목표**: Flutter에서 Rust 코어 연동, BLE/NFC UI, 측정 결과·트렌드 차트, E2E 후 S5·Phase 1C 완료

---

## 1. 현재 상태

- **S5a**: Go 4서비스 Docker, Dart gRPC 클라이언트, Auth/Device/Measurement/User Repository, Home/Measure/Devices/Settings 실제 데이터 연동, 60+ 테스트
- **Rust**: `rust-core/flutter-bridge` — `#[frb]` API 정의 완료 (differential, fingerprint, ble_scan/ble_connect, nfc_read_cartridge, 유틸)
- **Flutter**: `flutter_rust_bridge` 주석 처리 상태, 측정 완료 후 “결과 확인”은 단순 리셋

---

## 2. S5b 범위

| 순서 | 항목 | 내용 | Gate |
|------|------|------|------|
| 1 | **측정 결과·트렌드 차트** | fl_chart 추가, 측정 결과 전용 화면, GetMeasurementHistory 기반 트렌드 차트 | 화면·라우트·Provider 연동 |
| 2 | **flutter_rust_bridge 활성화** | pubspec 활성화, FFI 빌드 스크립트/CI, 생성 Dart 연동 | 빌드 성공 (플랫폼별 스텁 허용) |
| 3 | **BLE 스캔/연결 UI** | Rust `ble_scan`/`ble_connect` 호출 또는 플랫폼 미지원 시 스텁, MeasurementScreen·DeviceList 연동 | 스캔 목록·연결 시도 UI |
| 4 | **NFC 카트리지 읽기 UI** | Rust `nfc_read_cartridge` 호출 또는 스텁, 측정 시작 전 카트리지 정보 표시 | 카트리지 정보 표시 UI |
| 5 | **통합 테스트·문서** | E2E(측정 플로우) 또는 주요 플로우 위젯 테스트, QUALITY_GATES S5 완료, CHANGELOG/CONTEXT 갱신 | S5 Gate 체크리스트, Phase 1C 완료 선언 |

---

## 3. 구현 순서

### Step 1: fl_chart + 측정 결과·트렌드 화면

- `pubspec.yaml`: `fl_chart: ^0.69.0` 추가
- **측정 결과 화면** (`/measurement/result`): 마지막 측정 값 요약(primary value, unit, sessionId), “트렌드 보기” 버튼
- **트렌드 차트**: GetMeasurementHistory 데이터로 시계열 라인 차트 (날짜·값)
- MeasurementScreen “결과 확인” → `/measurement/result` 로 이동 (sessionId 또는 최근 결과 전달)
- 라우터: `GoRoute(path: '/measurement/result', ...)` 추가

### Step 2: flutter_rust_bridge 활성화 + FFI 빌드

- Flutter: `flutter_rust_bridge: ^2.0.0` (또는 프로젝트 호환 버전) 활성화
- Rust: `flutter-bridge` 빌드 시 codegen으로 Dart 생성 (기존 `flutter_rust_bridge_codegen`)
- Flutter에서 생성된 API import, 초기화(예: `get_engine_version()`) 호출로 연동 검증
- 플랫폼별: Android/iOS 네이티브 빌드 설정, Desktop/Web 미지원 시 스텁 또는 조건부 import

### Step 3: BLE 스캔/연결 UI

- **디바이스 목록**: DeviceListScreen 또는 측정 화면에서 “디바이스 검색” → FFI `ble_scan()` 호출, 목록 표시
- **연결**: 선택 시 `ble_connect(device_id)` 호출, 성공 시 MeasurementScreen에서 해당 deviceId 사용
- 미지원 플랫폼: `ble_scan`/`ble_connect` 스텁(빈 목록/실패)으로 UI만 동작

### Step 4: NFC 카트리지 읽기 UI

- 측정 시작 전 “카트리지 읽기” 버튼 → FFI `nfc_read_cartridge()` 호출
- CartridgeInfoDto(cartridge_id, type, lot_id, expiry, remaining_uses) 표시, 측정 세션 시 cartridgeId 전달
- 미지원: 스텁으로 기본 cartridgeId 사용

### Step 5: 통합 테스트 + S5 Gate

- E2E: 로그인 → 홈 → 측정 시작 → (시뮬레이션) 완료 → 결과 화면 → 트렌드 또는 위젯 테스트로 대체
- QUALITY_GATES.md: S5 Gate 체크리스트 완료, S5 상태 “통과”
- CHANGELOG.md: S5b 작업 로그
- CONTEXT.md: Phase 1C 완료 반영

---

## 4. S5 Gate 통과 기준 (S5b 반영)

- [ ] fl_chart 기반 측정 결과·트렌드 차트 화면 동작
- [ ] flutter_rust_bridge 연동 또는 명시적 스텁(문서화)
- [ ] BLE 스캔/연결 UI (실기기 또는 스텁)
- [ ] NFC 카트리지 읽기 UI (실기기 또는 스텁)
- [ ] 통합/위젯 테스트 유지·보강
- [ ] flutter analyze 0 에러, flutter test 통과
- [ ] CHANGELOG·CONTEXT·QUALITY_GATES 갱신

---

## 5. 핵심 파일 참조

- Rust FFI: `rust-core/flutter-bridge/src/lib.rs`
- Flutter 라우터: `frontend/flutter-app/lib/core/router/app_router.dart`
- 측정 Repository: `frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart`
- 기록 Provider: `frontend/flutter-app/lib/core/providers/grpc_provider.dart` (measurementHistoryProvider)

---

**문서 버전**: 1.0  
**최종 업데이트**: 2026-02-10
