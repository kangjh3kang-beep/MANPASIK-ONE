/// H1: 하네스 추상화 계층 — Rust 엔진 인터페이스 계약
///
/// 구체 구현(네이티브/스텁)은 이 인터페이스 뒤에 격리됩니다.
/// 모든 AFE 블록은 이 HAL(SensorTrait) 기반 인터페이스를 통해 접근하며,
/// 구체 타입을 직접 참조하지 않습니다.
///
/// H2: CartridgeManifest v1.0 카트리지는 v2.0 리더기/앱에서 반드시 동작합니다.
library;

import 'rust_ffi_stub.dart';

/// Rust 엔진에 대한 추상 인터페이스 계약.
///
/// [NativeRustEngine] (네이티브 Rust FFI) 또는 [StubRustEngine] (Dart 폴백)이
/// 이 인터페이스를 구현합니다. 소비 코드는 구체 구현을 알 필요가 없습니다.
abstract class RustEngineInterface {
  /// 엔진 초기화
  Future<void> initialize();

  /// 엔진 버전 (semver)
  String get engineVersion;

  /// 네이티브 엔진 활성화 여부
  bool get isNativeEnabled;

  /// BLE 디바이스 스캔
  Future<List<DeviceInfoDto>> bleScan();

  /// BLE 연결
  Future<bool> bleConnect(String deviceId);

  /// BLE 연결 해제
  Future<void> bleDisconnect(String deviceId);

  /// BLE 배터리 읽기
  Future<int> bleReadBattery(String deviceId);

  /// NFC 카트리지 읽기
  Future<CartridgeInfoDto> nfcReadCartridge();

  /// 차분 측정 (단일) — SSOT: Sdiff_n = S_n - alpha_n * R_n
  Future<DifferentialCorrectionDto> differentialMeasureSingle(
    double sDet,
    double sRef,
    double alpha,
  );

  /// 차분 측정 (배열) — SSOT: Sdiff_n = S_n - alpha_n * R_n
  Future<List<double>> differentialMeasure(
    List<double> sDet,
    List<double> sRef,
    double alpha,
  );

  /// 핑거프린트 생성 (88차원 — 기본)
  Future<FingerprintDto> createFingerprint88(List<double> data);

  /// 핑거프린트 생성 (896차원 — 통합)
  Future<FingerprintDto> createFingerprint896(List<double> data);

  /// 전체 측정 파이프라인 실행 (BLE -> DSP -> AI)
  Future<MeasurementPipelineDto> runMeasurementPipeline({
    required List<double> sDet,
    required List<double> sRef,
    required double alpha,
    required String biomarker,
    required String unit,
  });

  /// AI 분석
  Future<AiAnalysisResultDto> analyzeMeasurement(
    double value,
    String biomarker,
  );

  /// 진단 정보
  RustBridgeDiagnostics get diagnostics;
}

// ============================================================================
// DTO 정의 — 인터페이스 계약에 필요한 데이터 전송 객체
// ============================================================================

/// 차분 보정 결과 DTO
class DifferentialCorrectionDto {
  final double sDet;
  final double sRef;
  final double alpha;
  final double corrected;

  const DifferentialCorrectionDto({
    required this.sDet,
    required this.sRef,
    required this.alpha,
    required this.corrected,
  });
}

/// 핑거프린트 DTO (88/448/896/1792차원)
class FingerprintDto {
  final List<double> values;
  final int dimensions;
  final String algorithm;

  const FingerprintDto({
    required this.values,
    required this.dimensions,
    required this.algorithm,
  });
}

/// 측정 파이프라인 결과 DTO
class MeasurementPipelineDto {
  final MeasurementResultDto measurement;
  final AiAnalysisDto analysis;
  final double pipelineDurationMs;

  const MeasurementPipelineDto({
    required this.measurement,
    required this.analysis,
    required this.pipelineDurationMs,
  });
}

/// AI 분석 결과 DTO
class AiAnalysisResultDto {
  final String riskLevel;
  final double healthScore;
  final String summary;
  final List<String> recommendations;
  final String trend;

  const AiAnalysisResultDto({
    required this.riskLevel,
    required this.healthScore,
    required this.summary,
    required this.recommendations,
    required this.trend,
  });
}
