/// Rust FFI Native Adapter — 네이티브 플랫폼 (Android/iOS)
///
/// flutter_rust_bridge를 통해 Rust 코어 엔진과 실제 FFI 통신을 수행합니다.
/// frb_generated/ 코드젠 결과를 사용하여 Rust 함수를 호출합니다.
///
/// 반환 타입은 dynamic으로 통일하여 rust_ffi_stub.dart와의 순환 의존을
/// 방지합니다. 호출측에서 필드 접근 후 로컬 DTO로 변환합니다.
library;

import 'frb_generated/frb_generated.dart';
import 'frb_generated/lib.dart' as frb;

bool _nativeReady = false;

/// Rust 네이티브 라이브러리 초기화 시도
///
/// 성공 시 true 반환, 라이브러리 미링크 시 false 반환 (스텁 폴백)
Future<bool> tryInitNative() async {
  if (_nativeReady) return true;
  try {
    await RustLib.init();
    _nativeReady = true;
    return true;
  } catch (_) {
    _nativeReady = false;
    return false;
  }
}

/// 엔진 버전 조회 (Cargo.toml 기준)
String nativeGetEngineVersion() {
  if (!_nativeReady) return '0.1.0-stub';
  try {
    return frb.getEngineVersion();
  } catch (_) {
    return '0.1.0-error';
  }
}

/// BLE 디바이스 스캔 — frb DeviceInfoDto 리스트 반환
Future<List<dynamic>?> nativeBleScan() async {
  if (!_nativeReady) return null;
  try {
    return await frb.bleScan();
  } catch (_) {
    return null;
  }
}

/// BLE 디바이스 연결
Future<bool?> nativeBleConnect(String deviceId) async {
  if (!_nativeReady) return null;
  try {
    return await frb.bleConnect(deviceId: deviceId);
  } catch (_) {
    return null;
  }
}

/// BLE 배터리 레벨 읽기
Future<int?> nativeBleReadBattery(String deviceId) async {
  if (!_nativeReady) return null;
  try {
    return await frb.bleReadBattery(deviceId: deviceId);
  } catch (_) {
    return null;
  }
}

/// BLE 연결 품질 조회
String? nativeBleConnectionQuality(int rssi) {
  if (!_nativeReady) return null;
  try {
    return frb.bleConnectionQuality(rssi: rssi);
  } catch (_) {
    return null;
  }
}

/// NFC 카트리지 읽기 — frb CartridgeInfoDto 반환
Future<dynamic> nativeNfcReadCartridge() async {
  if (!_nativeReady) return null;
  try {
    return await frb.nfcReadCartridge();
  } catch (_) {
    return null;
  }
}

/// 차동 계측 + AI 분석 파이프라인 — frb MeasurementPipelineDto 반환
Future<dynamic> nativeProcessMeasurement({
  required List<double> sDet,
  required List<double> sRef,
  required double alpha,
  required String biomarker,
  required String unit,
}) async {
  if (!_nativeReady) return null;
  try {
    return await frb.runMeasurementPipeline(
      sDet: sDet,
      sRef: sRef,
      alpha: alpha,
      biomarker: biomarker,
      unit: unit,
    );
  } catch (_) {
    return null;
  }
}

/// AI 분석만 수행 — frb AiAnalysisResultDto 반환
Future<dynamic> nativeAnalyzeMeasurement({
  required double value,
  required String biomarker,
}) async {
  if (!_nativeReady) return null;
  try {
    return await frb.analyzeMeasurement(value: value, biomarker: biomarker);
  } catch (_) {
    return null;
  }
}

/// 전체 측정 파이프라인 (BLE→DSP→AI) — frb MeasurementPipelineDto 반환
Future<dynamic> nativeRunMeasurementPipeline({
  required List<double> sDet,
  required List<double> sRef,
  required double alpha,
  required String biomarker,
  required String unit,
}) async {
  if (!_nativeReady) return null;
  try {
    return await frb.runMeasurementPipeline(
      sDet: sDet,
      sRef: sRef,
      alpha: alpha,
      biomarker: biomarker,
      unit: unit,
    );
  } catch (_) {
    return null;
  }
}
