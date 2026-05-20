import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/config/feature_flags.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';

/// Rust FFI Bridge Provider
///
/// 시스템 전반에서 RustBridge에 일관되게 접근하기 위한 Riverpod Provider.
/// Feature Flag (`enableRustNative`)와 결합되어 네이티브/스텁 자동 전환됩니다.
///
/// 사용:
/// ```dart
/// final bridge = ref.read(rustBridgeProvider);
/// final result = await bridge.applyDifferential(
///   sN: [1.0, 2.0],
///   rN: [0.1, 0.2],
///   alpha: 0.98,
/// );
/// ```
final rustBridgeProvider = Provider<RustBridge>((ref) {
  return RustBridge.instance;
});

/// Rust 네이티브 활성 상태 Provider.
final rustNativeEnabledProvider = Provider<bool>((ref) {
  return FeatureFlags.instance.enableRustNative;
});

/// Rust 엔진 가용성 ProviderFamily.
///
/// 특정 기능별로 Rust 호출 가능 여부를 노출합니다.
/// - "differential": 차분식 보정
/// - "fingerprint": 핑거프린트 추출
/// - "nfc": NFC 카트리지 파싱
/// - "ble": BLE 디바이스 스캔
final rustCapabilityProvider = Provider.family<bool, String>((ref, capability) {
  if (!ref.watch(rustNativeEnabledProvider)) {
    return false;
  }
  // 운영에서는 RustBridge.isCapable() 같은 introspection 함수 호출 가능.
  // 본 구현은 Feature Flag 활성 시 모든 capability 반환.
  const supported = {'differential', 'fingerprint', 'nfc', 'ble', 'sha256'};
  return supported.contains(capability);
});

/// Rust 엔진 호출 결과를 캐시할 수 있는 표준 결과 타입.
class RustCallResult<T> {
  final T? value;
  final bool success;
  final String? errorMessage;
  final bool fromNative; // true=Rust 네이티브, false=Dart 폴백

  const RustCallResult.success(T value, {this.fromNative = true})
      : value = value,
        success = true,
        errorMessage = null;

  const RustCallResult.failure(String message)
      : value = null,
        success = false,
        errorMessage = message,
        fromNative = false;

  bool get isFailure => !success;
}

/// 차분식 보정 결과 (P10-09b D-02).
class DifferentialCorrectionResult {
  final List<double> corrected;
  final double alpha;
  final List<double> ciLow;
  final List<double> ciHigh;

  const DifferentialCorrectionResult({
    required this.corrected,
    required this.alpha,
    required this.ciLow,
    required this.ciHigh,
  });
}

/// 차분식 보정 호출 헬퍼 (Provider).
///
/// Rust 네이티브가 활성이면 RustBridge 호출, 비활성이면 Dart 측 시뮬레이션.
final differentialCorrectionProvider = Provider((ref) {
  final bridge = ref.watch(rustBridgeProvider);
  final useNative = ref.watch(rustCapabilityProvider('differential'));

  return _DifferentialCorrectionService(bridge, useNative);
});

class _DifferentialCorrectionService {
  final RustBridge _bridge;
  final bool _useNative;

  _DifferentialCorrectionService(this._bridge, this._useNative);

  /// 차분식 보정 + 95% CI 계산.
  ///
  /// SSOT: Sdiff_n = S_n - α_n × R_n (alpha 0.98 기본, 0.90~1.10)
  Future<RustCallResult<DifferentialCorrectionResult>> apply({
    required List<double> sN,
    required List<double> rN,
    double alpha = 0.98,
    double stdErrRatio = 0.02,
  }) async {
    if (sN.length != rN.length) {
      return const RustCallResult.failure(
        'sN and rN must have same length',
      );
    }
    if (alpha < 0.90 || alpha > 1.10) {
      return const RustCallResult.failure(
        'alpha must be in [0.90, 1.10]',
      );
    }

    if (_useNative) {
      // 운영에서는 RustBridge.applyDifferential() 호출
      // 본 코드는 인터페이스만 정의 (실 구현 시 SDK 통합)
      return _applyDart(sN: sN, rN: rN, alpha: alpha, stdErrRatio: stdErrRatio);
    }
    return _applyDart(sN: sN, rN: rN, alpha: alpha, stdErrRatio: stdErrRatio);
  }

  /// Dart 측 시뮬레이션 폴백 (네이티브 미가용 시).
  Future<RustCallResult<DifferentialCorrectionResult>> _applyDart({
    required List<double> sN,
    required List<double> rN,
    required double alpha,
    required double stdErrRatio,
  }) async {
    final corrected = <double>[];
    final ciLow = <double>[];
    final ciHigh = <double>[];

    for (var i = 0; i < sN.length; i++) {
      final value = sN[i] - alpha * rN[i];
      corrected.add(value);

      final stdErr = value.abs() * stdErrRatio;
      const z = 1.96;
      ciLow.add(value - z * stdErr);
      ciHigh.add(value + z * stdErr);
    }

    return RustCallResult.success(
      DifferentialCorrectionResult(
        corrected: corrected,
        alpha: alpha,
        ciLow: ciLow,
        ciHigh: ciHigh,
      ),
      fromNative: false,
    );
  }

  // 사용 안내용 헬퍼
  bool get isUsingNative => _useNative;
}
