import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/services/rust_ffi_stub.dart';

/// H1: Rust 엔진 인터페이스 Provider
///
/// [RustEngineInterface] 타입으로 노출하여 소비 코드가
/// 구체 구현(RustBridge, NativeRustEngine, StubRustEngine)에 의존하지 않도록 합니다.
///
/// 사용:
/// ```dart
/// final engine = ref.read(rustEngineProvider);
/// final devices = await engine.bleScan();
/// final result = await engine.differentialMeasureSingle(1.0, 0.1, 0.98);
/// ```
final rustEngineProvider = Provider<RustEngineInterface>((ref) {
  return RustBridge.instance; // Returns the singleton implementing RustEngineInterface
});

/// Rust 엔진 진단 정보 Provider
///
/// 디버그 화면 등에서 엔진 상태를 확인할 때 사용합니다.
final rustEngineDiagnosticsProvider = Provider<RustBridgeDiagnostics>((ref) {
  final engine = ref.watch(rustEngineProvider);
  return engine.diagnostics;
});

/// Rust 엔진 네이티브 활성 여부 Provider
final rustEngineNativeEnabledProvider = Provider<bool>((ref) {
  final engine = ref.watch(rustEngineProvider);
  return engine.isNativeEnabled;
});
