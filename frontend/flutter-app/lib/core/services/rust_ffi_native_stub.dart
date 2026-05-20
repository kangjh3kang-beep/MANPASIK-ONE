/// Rust FFI Native Adapter — Web/Desktop 스텁
///
/// 네이티브 Rust 라이브러리가 불가능한 플랫폼(Web, Desktop)에서 사용됩니다.
/// 모든 함수가 null/false를 반환하여 rust_ffi_stub.dart가 Dart 폴백을 사용하도록 합니다.
library;

Future<bool> tryInitNative() async => false;

String nativeGetEngineVersion() => '0.1.0-stub';

Future<List<dynamic>?> nativeBleScan() async => null;

Future<bool?> nativeBleConnect(String deviceId) async => null;

Future<int?> nativeBleReadBattery(String deviceId) async => null;

String? nativeBleConnectionQuality(int rssi) => null;

Future<dynamic> nativeNfcReadCartridge() async => null;

Future<dynamic> nativeProcessMeasurement({
  required List<double> sDet,
  required List<double> sRef,
  required double alpha,
  required String biomarker,
  required String unit,
}) async => null;

Future<dynamic> nativeAnalyzeMeasurement({
  required double value,
  required String biomarker,
}) async => null;

Future<dynamic> nativeRunMeasurementPipeline({
  required List<double> sDet,
  required List<double> sRef,
  required double alpha,
  required String biomarker,
  required String unit,
}) async => null;
