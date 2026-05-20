# Phase D — 네이티브 통합 계획서

**문서 ID**: MPS-PLAN-PHASE-D-v1.0
**작성일**: 2026-04-25
**상태**: Phase D-1~D-3 구현 완료

---

## 1. 실사 결과 요약

### 1.1 네이티브 통합 현황 (Phase D 이전)

초기 MEMORY.md에서 네이티브 완성도를 25%로 추정했으나, 정밀 실사 결과 **~45%** 이미 구현 완료 확인:

| 영역 | 구현체 | 파일 | 상태 |
|------|--------|------|------|
| Rust 엔진 9모듈 | differential, fingerprint, ble, nfc, ai, dsp, crypto, sync, api | `manpasik-engine/src/` | **100%** |
| flutter_rust_bridge codegen | RustLib + 20 FFI API | `frb_generated/` | **100%** |
| flutter-bridge 래퍼 | 16 pub fn (BLE/NFC/DSP/AI/Fingerprint) | `flutter-bridge/src/lib.rs` | **100%** |
| Flutter RustBridge 싱글톤 | DeviceInfoDto, CartridgeInfoDto 등 DTO | `rust_ffi_stub.dart` | **90%** |
| `_useNative` 호출 경로 | if 분기 존재하나 **모두 TODO 빈칸** | `rust_ffi_stub.dart` | **0%** |
| Conditional Import (Web/Native) | 미구현 — frb_generated 직접 import 시 Web 빌드 깨짐 | - | **0%** |
| CRDT 오프라인 동기화 | GCounter, LWWRegister, ORSet, SyncManager | `manpasik-engine/src/sync/` | **100%** |
| Rust 테스트 | engine 85 + bridge 8 = 93 | `rust-core/` | **85%** |

### 1.2 실제 GAP 3가지

1. **네이티브 호출 경로 0%**: `_useNative = true`일 때 실제 Rust 호출이 모두 비어있음
2. **Conditional Import 미구현**: `frb_generated.web.dart` 미존재로 Web 빌드 호환 불가
3. **NFC 파싱 FFI 부재**: NFC 바이트 파싱이 flutter-bridge에 노출되지 않음

---

## 2. Phase D 구현 완료 내역

### D-1: RustBridge 네이티브 호출 경로 (완료)

**아키텍처**: Conditional Import 패턴으로 Web/Native 분리

```
rust_ffi_stub.dart
  ├── import 'rust_ffi_native_stub.dart'     ← Web/Desktop (no-op)
  │   if (dart.library.io)
  └── import 'rust_ffi_native_impl.dart'     ← Android/iOS (실제 FFI)
```

**신규 파일**:
- `lib/core/services/rust_ffi_native_stub.dart` — Web/Desktop 스텁 (모든 함수 null/false 반환)
- `lib/core/services/rust_ffi_native_impl.dart` — Native 어댑터 (frb_generated → try/catch → null 폴백)

**수정 파일**: `lib/core/services/rust_ffi_stub.dart`
- Conditional import 추가: `import '...native_stub.dart' if (dart.library.io) '...native_impl.dart'`
- `init()`: `native_engine.tryInitNative()` 호출로 네이티브 초기화 시도
- `engineVersion`: `native_engine.nativeGetEngineVersion()` 호출
- **9개 메서드의 `_useNative` 분기 전부 구현**:
  - `bleScan()` — dynamic → DeviceInfoDto 변환
  - `bleConnect()` — native_engine.nativeBleConnect()
  - `bleReadBattery()` — native_engine.nativeBleReadBattery()
  - `bleConnectionQuality()` — native_engine.nativeBleConnectionQuality()
  - `nfcReadCartridge()` — dynamic → CartridgeInfoDto 변환
  - `processMeasurement()` — dynamic → MeasurementResultDto 변환
  - `analyzeResult()` — dynamic → AiAnalysisDto 변환
  - `runMeasurementPipeline()` — dynamic → MeasurementPipelineResult 변환
  - `bleDisconnect()` — always-stub (Rust API에 미존재)

**설계 원칙**:
- `dynamic` 반환 타입으로 순환 의존 방지 (frb_generated ↔ rust_ffi_stub)
- H6 실패 격리: native 호출 실패 시 자동으로 Dart 스텁 폴백
- null 반환 = "네이티브 불가" → 스텁 경로로 투명 전환

### D-2: NFC Parse Tag FFI (완료)

**수정 파일**: `rust-core/flutter-bridge/src/lib.rs`
- `nfc_parse_tag(data: Vec<u8>) -> Result<CartridgeInfoDto, String>` 추가
  - 플랫폼 독립적 NFC 바이트 파싱 (하드웨어 의존 없음)
  - v1.0 (64바이트) 및 v2.0 (80바이트) 포맷 지원
  - Flutter에서 platform plugin으로 raw NFC 바이트 읽기 → Rust로 파싱 위임 가능

### D-3: 테스트 보강 (완료)

**수정 파일**: `rust-core/flutter-bridge/src/lib.rs`
- 기존 8개 → **13개 테스트** (5개 추가):
  - `test_nfc_parse_tag_v1` — v1.0 64바이트 태그 데이터 파싱 검증
  - `test_nfc_parse_tag_too_short` — 짧은 데이터 에러 처리 검증
  - `test_sha256_hash` — SHA-256 해시 출력 길이 검증
  - `test_differential_measure_vec` — 벡터 차동 측정 (S_corrected ≈ 0.9902)
  - `test_cosine_similarity_identical` — 동일 벡터 코사인 유사도 = 1.0

---

## 3. 검증 결과

| 항목 | 결과 |
|------|------|
| Rust manpasik-engine 빌드 | **PASS** |
| Rust flutter-bridge 빌드 | **PASS** |
| Rust manpasik-engine 테스트 | **85/85 PASS** |
| Rust flutter-bridge 테스트 | **13/13 PASS** |
| Flutter analyze | **에러 0** (info/warning만) |
| Go 빌드 | **ALL PASS** |
| Go 테스트 | **46 패키지 PASS** |

---

## 4. 네이티브 완성도 변화

| 시점 | 완성도 | 근거 |
|------|--------|------|
| Phase D 이전 (추정) | 25% | MEMORY.md 초기 추정 |
| Phase D 이전 (실사) | ~45% | Rust 엔진/브릿지/codegen 구현 완료, 호출 경로 0% |
| **Phase D 이후** | **~60%** | 네이티브 호출 경로 100% + NFC FFI + 테스트 13개 |

### 남은 네이티브 작업 (40%)
- Android NDK / iOS 빌드 설정 + 실 디바이스 테스트
- BLE 실 하드웨어 연동 (btleplug → Android/iOS BLE 권한)
- NFC 실 하드웨어 연동 (PN7150 → nfc_manager plugin)
- CRDT SyncManager → Flutter SyncProvider 연결
- flutter_rust_bridge codegen 재실행 (nfc_parse_tag 노출)
- 오프라인 큐 영속화 (SharedPreferences/Hive)

---

*자동 생성: 2026-04-25 | 만파식(萬波息) Phase D 네이티브 통합 계획서*
