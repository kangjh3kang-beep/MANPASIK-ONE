/// Rust FFI Bridge — ManPaSik Core Engine 연동
///
/// flutter_rust_bridge 기반 네이티브-스텁 하이브리드 모드:
/// - 모바일 (Android/iOS): 네이티브 Rust 엔진 사용 (라이브러리 존재 시)
/// - Web/Desktop: 스텁 모드 (시뮬레이션)
/// - 네이티브 로드 실패 시: 자동 스텁 폴백
///
/// H1: [RustEngineInterface]를 구현하여 구체 타입(네이티브/스텁)을
/// 인터페이스 뒤에 격리합니다. 소비 코드는 [RustEngineInterface]만 참조합니다.
library;

import 'dart:io' show Platform;
import 'dart:math';

import 'package:flutter/foundation.dart' show kIsWeb, kReleaseMode;

import 'rust_engine_interface.dart';
import 'rust_ffi_native_stub.dart'
    if (dart.library.io) 'rust_ffi_native_impl.dart' as native_engine;

// Re-export interface types so consumers can import from one place.
export 'rust_engine_interface.dart';

/// BLE 디바이스 정보 (Rust DeviceInfoDto 대응)
class DeviceInfoDto {
  final String deviceId;
  final String name;
  final int rssi;
  final String state;

  const DeviceInfoDto({
    required this.deviceId,
    required this.name,
    required this.rssi,
    required this.state,
  });
}

/// 카트리지 정보 (Rust CartridgeInfoDto 대응)
class CartridgeInfoDto {
  final String cartridgeId;
  final String cartridgeType;
  final String lotId;
  final String expiryDate;
  final int remainingUses;

  const CartridgeInfoDto({
    required this.cartridgeId,
    required this.cartridgeType,
    required this.lotId,
    required this.expiryDate,
    required this.remainingUses,
  });
}

/// 차동 측정 결과 (Rust MeasurementResult 대응)
class MeasurementResultDto {
  final double primaryValue;
  final double referenceValue;
  final double differentialValue;
  final double snr;
  final double confidence;
  final String unit;
  final String biomarker;
  final DateTime timestamp;

  const MeasurementResultDto({
    required this.primaryValue,
    required this.referenceValue,
    required this.differentialValue,
    required this.snr,
    required this.confidence,
    required this.unit,
    required this.biomarker,
    required this.timestamp,
  });
}

/// AI 분석 결과 (코칭 추천 포함)
class AiAnalysisDto {
  final String riskLevel; // normal, caution, warning, critical
  final double healthScore; // 0-100
  final String summary;
  final List<String> recommendations;
  final String trend; // improving, stable, declining

  const AiAnalysisDto({
    required this.riskLevel,
    required this.healthScore,
    required this.summary,
    required this.recommendations,
    required this.trend,
  });
}

/// 측정 파이프라인 결과 (BLE -> DSP -> AI 전체 결과)
class MeasurementPipelineResult {
  final MeasurementResultDto measurement;
  final AiAnalysisDto analysis;
  final double pipelineDurationMs;

  const MeasurementPipelineResult({
    required this.measurement,
    required this.analysis,
    required this.pipelineDurationMs,
  });
}

class RustBridgeDiagnostics {
  const RustBridgeDiagnostics({
    required this.initialized,
    required this.nativeEnabled,
    required this.nativePlatform,
    required this.releaseBuild,
    required this.engineVersion,
    required this.releaseStubAllowed,
  });

  final bool initialized;
  final bool nativeEnabled;
  final bool nativePlatform;
  final bool releaseBuild;
  final String engineVersion;
  final bool releaseStubAllowed;

  bool get isStubFallback => !nativeEnabled;

  bool get canRunMeasurement =>
      nativeEnabled || !releaseBuild || releaseStubAllowed;

  String get modeLabel => nativeEnabled ? 'native' : 'stub-fallback';

  String? get blockingReason {
    if (canRunMeasurement) return null;
    return 'Rust native measurement engine is unavailable in release mode.';
  }
}

// ============================================================================
// RustBridge — 싱글톤 팩토리 (H1 인터페이스 제공)
// ============================================================================

/// ManPaSik Rust Core Engine Bridge
///
/// 싱글톤 패턴. [RustBridge.init]으로 초기화 후 사용.
/// 플랫폼 감지 기반 네이티브/스텁 자동 전환.
///
/// H1: [RustEngineInterface]를 구현하여 구체 타입(네이티브/스텁)을
/// 인터페이스 뒤에 격리합니다. 소비 코드는 항상 [RustEngineInterface]를 통해 접근합니다.
class RustBridge implements RustEngineInterface {
  RustBridge._();
  static final RustBridge _instance = RustBridge._();

  /// 싱글톤 인스턴스 (RustEngineInterface 타입으로 노출)
  static RustBridge get instance => _instance;

  bool _initialized = false;
  bool get isInitialized => _initialized;

  /// 네이티브 Rust 엔진 사용 가능 여부
  static bool _useNative = false;

  static const bool _releaseStubAllowed = bool.fromEnvironment(
    'MANPASIK_RELEASE_ALLOW_STUB_MEASUREMENT',
    defaultValue: false,
  );

  /// 플랫폼 감지: 네이티브 가능 여부
  static bool get _isNativePlatform {
    if (kIsWeb) return false;
    return Platform.isAndroid || Platform.isIOS;
  }

  // ── RustEngineInterface 구현 ──

  @override
  bool get isNativeEnabled => _useNative;

  @override
  String get engineVersion =>
      _useNative ? native_engine.nativeGetEngineVersion() : '0.1.0-stub';

  @override
  RustBridgeDiagnostics get diagnostics => RustBridgeDiagnostics(
        initialized: _initialized,
        nativeEnabled: _useNative,
        nativePlatform: _isNativePlatform,
        releaseBuild: kReleaseMode,
        engineVersion: engineVersion,
        releaseStubAllowed: _releaseStubAllowed,
      );

  @override
  Future<void> initialize() => RustBridge.init();

  /// 정적 초기화 (기존 호출 사이트 호환)
  static Future<void> init() async {
    if (_instance._initialized) return;

    if (_isNativePlatform) {
      try {
        _useNative = await native_engine.tryInitNative();
      } catch (_) {
        _useNative = false;
      }
    }

    _instance._initialized = true;
  }

  // ── BLE API (인스턴스 메서드 — RustEngineInterface 계약) ──

  @override
  Future<List<DeviceInfoDto>> bleScan() => _bleScanImpl();

  @override
  Future<bool> bleConnect(String deviceId) => _bleConnectImpl(deviceId);

  @override
  Future<void> bleDisconnect(String deviceId) =>
      _bleDisconnectImpl(deviceId);

  @override
  Future<int> bleReadBattery(String deviceId) =>
      _bleReadBatteryImpl(deviceId);

  // ── NFC API ──

  @override
  Future<CartridgeInfoDto> nfcReadCartridge() => _nfcReadCartridgeImpl();

  // ── 차분 측정 API (SSOT: Sdiff_n = S_n - alpha_n * R_n) ──

  @override
  Future<DifferentialCorrectionDto> differentialMeasureSingle(
    double sDet,
    double sRef,
    double alpha,
  ) async {
    final corrected = sDet - alpha * sRef;
    return DifferentialCorrectionDto(
      sDet: sDet,
      sRef: sRef,
      alpha: alpha,
      corrected: corrected,
    );
  }

  @override
  Future<List<double>> differentialMeasure(
    List<double> sDet,
    List<double> sRef,
    double alpha,
  ) async {
    final n = sDet.length < sRef.length ? sDet.length : sRef.length;
    return List.generate(n, (i) => sDet[i] - alpha * sRef[i]);
  }

  // ── 핑거프린트 API ──

  @override
  Future<FingerprintDto> createFingerprint88(List<double> data) async {
    // 88차원 기본 핑거프린트 — 입력을 그대로 사용하거나 88차원으로 자름/패딩
    final values = _normalizeToLength(data, 88);
    return FingerprintDto(
      values: values,
      dimensions: 88,
      algorithm: 'electrochemical-v1',
    );
  }

  @override
  Future<FingerprintDto> createFingerprint896(List<double> data) async {
    // 896차원 통합 핑거프린트 — 88차원 기본 + 확장 특징 추출
    final base88 = _normalizeToLength(data, 88);
    final extended = <double>[];
    extended.addAll(base88);
    // 확장: 10개 윈도우별 통계 특징으로 896차원 생성
    for (int w = 0; w < 10; w++) {
      final windowStart = (w * base88.length / 10).floor();
      final windowEnd = ((w + 1) * base88.length / 10).floor();
      final window = base88.sublist(windowStart, windowEnd);
      if (window.isEmpty) continue;
      final mean = window.reduce((a, b) => a + b) / window.length;
      final variance =
          window.map((v) => (v - mean) * (v - mean)).reduce((a, b) => a + b) /
              window.length;
      extended.addAll([
        mean,
        sqrt(variance),
        window.reduce(max),
        window.reduce(min),
      ]);
    }
    final values = _normalizeToLength(extended, 896);
    return FingerprintDto(
      values: values,
      dimensions: 896,
      algorithm: 'multi-modal-v1',
    );
  }

  // ── 측정 파이프라인 (RustEngineInterface 계약) ──

  @override
  Future<MeasurementPipelineDto> runMeasurementPipeline({
    required List<double> sDet,
    required List<double> sRef,
    required double alpha,
    required String biomarker,
    required String unit,
  }) async {
    final stopwatch = Stopwatch()..start();

    final measurement = await _processMeasurementImpl(
      signalData: sDet,
      referenceData: sRef,
      biomarker: biomarker,
      unit: unit,
    );

    final analysis = await _analyzeResultImpl(
      value: measurement.primaryValue,
      biomarker: biomarker,
      unit: unit,
    );

    stopwatch.stop();

    return MeasurementPipelineDto(
      measurement: measurement,
      analysis: analysis,
      pipelineDurationMs: stopwatch.elapsedMilliseconds.toDouble(),
    );
  }

  // ── AI 분석 (RustEngineInterface 계약) ──

  @override
  Future<AiAnalysisResultDto> analyzeMeasurement(
    double value,
    String biomarker,
  ) async {
    final dto = await _analyzeResultImpl(
      value: value,
      biomarker: biomarker,
      unit: 'mg/dL',
    );
    return AiAnalysisResultDto(
      riskLevel: dto.riskLevel,
      healthScore: dto.healthScore,
      summary: dto.summary,
      recommendations: dto.recommendations,
      trend: dto.trend,
    );
  }

  // ══════════════════════════════════════════════════════════════════════════
  // 레거시 정적 API (기존 호출 사이트 후방 호환 — H2)
  //
  // 신규 코드는 RustEngineInterface 인스턴스 메서드를 사용해주세요.
  // ══════════════════════════════════════════════════════════════════════════

  /// BLE 디바이스 스캔 (정적 — 레거시 호환)
  static Future<List<DeviceInfoDto>> bleScanStatic({
    Duration timeout = const Duration(seconds: 5),
  }) =>
      _bleScanImpl();

  /// BLE 디바이스 연결 (정적 — 레거시 호환)
  static Future<bool> bleConnectStatic(String deviceId) =>
      _bleConnectImpl(deviceId);

  /// BLE 디바이스 연결 해제 (정적 — 레거시 호환)
  static Future<void> bleDisconnectStatic(String deviceId) =>
      _bleDisconnectImpl(deviceId);

  /// BLE 배터리 읽기 (정적 — 레거시 호환)
  static Future<int> bleReadBatteryStatic(String deviceId) =>
      _bleReadBatteryImpl(deviceId);

  /// BLE 연결 품질 (RSSI 기반)
  static Future<String> bleConnectionQuality(String deviceId) async {
    if (_useNative) {
      final result = native_engine.nativeBleConnectionQuality(-45);
      if (result != null) return result;
    }
    return 'excellent'; // excellent, good, fair, poor
  }

  /// NFC 카트리지 읽기 (정적 — 레거시 호환)
  static Future<CartridgeInfoDto> nfcReadCartridgeStatic() =>
      _nfcReadCartridgeImpl();

  /// 차동 계측 처리 (정적 — 레거시 호환)
  static Future<MeasurementResultDto> processMeasurement({
    required List<double> signalData,
    required List<double> referenceData,
    required String biomarker,
    String unit = 'mg/dL',
  }) =>
      _processMeasurementImpl(
        signalData: signalData,
        referenceData: referenceData,
        biomarker: biomarker,
        unit: unit,
      );

  /// AI 분석 (정적 — 레거시 호환)
  static Future<AiAnalysisDto> analyzeResult({
    required double value,
    required String biomarker,
    required String unit,
    List<double>? recentValues,
  }) =>
      _analyzeResultImpl(
        value: value,
        biomarker: biomarker,
        unit: unit,
        recentValues: recentValues,
      );

  /// 전체 측정 파이프라인 (정적 — 레거시 호환)
  static Future<MeasurementPipelineResult> runMeasurementPipelineStatic({
    required String deviceId,
    required String biomarker,
    String unit = 'mg/dL',
    List<double>? recentValues,
  }) async {
    final stopwatch = Stopwatch()..start();

    if (_useNative) {
      final rng = Random(DateTime.now().microsecondsSinceEpoch);
      final nativeSDet = List.generate(88, (i) {
        final base = 0.5 + 0.3 * sin(i * 0.15);
        return base + rng.nextDouble() * 0.05;
      });
      final nativeSRef = List.generate(88, (i) {
        return 0.02 + 0.01 * sin(i * 0.15) + rng.nextDouble() * 0.005;
      });
      final result = await native_engine.nativeRunMeasurementPipeline(
        sDet: nativeSDet,
        sRef: nativeSRef,
        alpha: 0.98,
        biomarker: biomarker,
        unit: unit,
      );
      if (result != null) {
        return MeasurementPipelineResult(
          measurement: MeasurementResultDto(
            primaryValue: result.primaryValue as double,
            referenceValue: result.referenceValue as double,
            differentialValue: result.differentialValue as double,
            snr: result.snr as double,
            confidence: result.confidence as double,
            unit: result.unit as String,
            biomarker: result.biomarker as String,
            timestamp: DateTime.now(),
          ),
          analysis: AiAnalysisDto(
            riskLevel: result.riskLevel as String,
            healthScore: result.healthScore as double,
            summary: '측정 완료',
            recommendations: (result.recommendations as List).cast<String>(),
            trend: 'stable',
          ),
          pipelineDurationMs: result.pipelineDurationMs as double,
        );
      }
    }

    final rng = Random(DateTime.now().microsecondsSinceEpoch);
    final signalData = List.generate(88, (i) {
      final base = 0.5 + 0.3 * sin(i * 0.15);
      return base + rng.nextDouble() * 0.05;
    });
    final referenceData = List.generate(88, (i) {
      return 0.02 + 0.01 * sin(i * 0.15) + rng.nextDouble() * 0.005;
    });
    await Future.delayed(const Duration(milliseconds: 100));

    final measurement = await processMeasurement(
      signalData: signalData,
      referenceData: referenceData,
      biomarker: biomarker,
      unit: unit,
    );

    final analysis = await analyzeResult(
      value: measurement.primaryValue,
      biomarker: biomarker,
      unit: unit,
      recentValues: recentValues,
    );

    stopwatch.stop();

    return MeasurementPipelineResult(
      measurement: measurement,
      analysis: analysis,
      pipelineDurationMs: stopwatch.elapsedMilliseconds.toDouble(),
    );
  }

  // ══════════════════════════════════════════════════════════════════════════
  // 내부 구현 (정적 — 인스턴스/레거시 양쪽에서 공유)
  // ══════════════════════════════════════════════════════════════════════════

  static Future<List<DeviceInfoDto>> _bleScanImpl() async {
    if (_useNative) {
      final nativeDevices = await native_engine.nativeBleScan();
      if (nativeDevices != null) {
        return nativeDevices
            .map((d) => DeviceInfoDto(
                  deviceId: d.deviceId as String,
                  name: d.name as String,
                  rssi: d.rssi as int,
                  state: d.state as String,
                ))
            .toList();
      }
    }
    await Future.delayed(const Duration(milliseconds: 800));
    return const [
      DeviceInfoDto(
        deviceId: 'MPK-DEMO-001',
        name: 'ManPaSik Pro X1',
        rssi: -45,
        state: 'discovered',
      ),
      DeviceInfoDto(
        deviceId: 'MPK-DEMO-002',
        name: 'ManPaSik Lite S1',
        rssi: -62,
        state: 'discovered',
      ),
    ];
  }

  static Future<bool> _bleConnectImpl(String deviceId) async {
    if (_useNative) {
      final result = await native_engine.nativeBleConnect(deviceId);
      if (result != null) return result;
    }
    await Future.delayed(const Duration(milliseconds: 600));
    return true;
  }

  static Future<void> _bleDisconnectImpl(String deviceId) async {
    // BLE disconnect는 현재 frb_generated API에 미포함 — 항상 스텁
    await Future.delayed(const Duration(milliseconds: 200));
  }

  static Future<int> _bleReadBatteryImpl(String deviceId) async {
    if (_useNative) {
      final result = await native_engine.nativeBleReadBattery(deviceId);
      if (result != null) return result;
    }
    await Future.delayed(const Duration(milliseconds: 100));
    return 85 + Random().nextInt(15);
  }

  static Future<CartridgeInfoDto> _nfcReadCartridgeImpl() async {
    if (_useNative) {
      final result = await native_engine.nativeNfcReadCartridge();
      if (result != null) {
        return CartridgeInfoDto(
          cartridgeId: result.cartridgeId as String,
          cartridgeType: result.cartridgeType as String,
          lotId: result.lotId as String,
          expiryDate: result.expiryDate as String,
          remainingUses: result.remainingUses as int,
        );
      }
    }
    await Future.delayed(const Duration(milliseconds: 400));
    return const CartridgeInfoDto(
      cartridgeId: 'CART-2026-001',
      cartridgeType: 'Glucose',
      lotId: 'LOT-2026A',
      expiryDate: '20270630',
      remainingUses: 8,
    );
  }

  static Future<MeasurementResultDto> _processMeasurementImpl({
    required List<double> signalData,
    required List<double> referenceData,
    required String biomarker,
    String unit = 'mg/dL',
  }) async {
    if (_useNative) {
      final result = await native_engine.nativeProcessMeasurement(
        sDet: signalData,
        sRef: referenceData,
        alpha: 0.98,
        biomarker: biomarker,
        unit: unit,
      );
      if (result != null) {
        return MeasurementResultDto(
          primaryValue: result.primaryValue as double,
          referenceValue: result.referenceValue as double,
          differentialValue: result.differentialValue as double,
          snr: result.snr as double,
          confidence: result.confidence as double,
          unit: result.unit as String,
          biomarker: result.biomarker as String,
          timestamp: DateTime.now(),
        );
      }
    }

    // Dart-side 차동 계측: S_diff = S_det - alpha * S_ref (SSOT 공식)
    const double alpha = 0.98; // SSOT 기본 alpha
    await Future.delayed(const Duration(milliseconds: 50));

    final n = signalData.length.clamp(1, referenceData.length.clamp(1, 99999));
    double sumDiff = 0.0;
    double sumSq = 0.0;
    for (int i = 0; i < n; i++) {
      final sDet = signalData[i];
      final sRef = i < referenceData.length ? referenceData[i] : 0.0;
      final diff = sDet - alpha * sRef;
      sumDiff += diff;
      sumSq += diff * diff;
    }
    final meanDiff = sumDiff / n;
    final variance = sumSq / n - meanDiff * meanDiff;
    final noiseFloor = variance > 1e-12 ? sqrt(variance) : 0.001;
    final snr = 20.0 *
        log(meanDiff.abs().clamp(1e-12, double.infinity) / noiseFloor) /
        ln10;
    final confidence = (1.0 / (1.0 + exp(-0.1 * (snr - 20)))).clamp(0.0, 1.0);

    return MeasurementResultDto(
      primaryValue: meanDiff,
      referenceValue: referenceData.isNotEmpty
          ? referenceData.reduce((a, b) => a + b) / referenceData.length
          : 0.0,
      differentialValue: meanDiff,
      snr: snr.clamp(0, 100),
      confidence: confidence,
      unit: unit,
      biomarker: biomarker,
      timestamp: DateTime.now(),
    );
  }

  static Future<AiAnalysisDto> _analyzeResultImpl({
    required double value,
    required String biomarker,
    required String unit,
    List<double>? recentValues,
  }) async {
    if (_useNative) {
      final result = await native_engine.nativeAnalyzeMeasurement(
        value: value,
        biomarker: biomarker,
      );
      if (result != null) {
        return AiAnalysisDto(
          riskLevel: result.riskLevel as String,
          healthScore: result.healthScore as double,
          summary: result.summary as String,
          recommendations: (result.recommendations as List).cast<String>(),
          trend: result.trend as String,
        );
      }
    }

    await Future.delayed(const Duration(milliseconds: 200));

    final (riskLevel, healthScore) = _classifyRisk(value, biomarker);
    final trend = _calculateTrend(value, recentValues);

    return AiAnalysisDto(
      riskLevel: riskLevel,
      healthScore: healthScore,
      summary: _generateSummary(value, biomarker, unit, riskLevel),
      recommendations: _generateRecommendations(biomarker, riskLevel),
      trend: trend,
    );
  }

  // ── 내부 헬퍼 ──

  static List<double> _normalizeToLength(List<double> data, int targetLen) {
    if (data.length == targetLen) return List.of(data);
    if (data.length > targetLen) return data.sublist(0, targetLen);
    return [...data, ...List.filled(targetLen - data.length, 0.0)];
  }

  static (String, double) _classifyRisk(double value, String biomarker) {
    switch (biomarker.toLowerCase()) {
      case 'glucose':
        if (value < 70) return ('caution', 65.0);
        if (value <= 100) return ('normal', 90.0);
        if (value <= 125) return ('caution', 70.0);
        return ('warning', 45.0);
      case 'hba1c':
        if (value <= 5.6) return ('normal', 92.0);
        if (value <= 6.4) return ('caution', 68.0);
        return ('warning', 40.0);
      case 'uric_acid':
        if (value >= 3.5 && value <= 7.2) return ('normal', 88.0);
        return ('caution', 60.0);
      default:
        return ('normal', 85.0);
    }
  }

  static String _calculateTrend(double current, List<double>? recent) {
    if (recent == null || recent.length < 2) return 'stable';
    final avg = recent.reduce((a, b) => a + b) / recent.length;
    final diff = current - avg;
    if (diff.abs() < avg * 0.05) return 'stable';
    return diff < 0 ? 'improving' : 'declining';
  }

  static String _generateSummary(
      double value, String biomarker, String unit, String risk) {
    final name = {
          'glucose': '혈당',
          'hba1c': '당화혈색소',
          'uric_acid': '요산',
          'creatinine': '크레아티닌',
          'vitamin_d': '비타민D',
        }[biomarker.toLowerCase()] ??
        biomarker;

    final statusText = {
          'normal': '정상 범위',
          'caution': '주의 범위',
          'warning': '경고 범위',
          'critical': '위험 범위',
        }[risk] ??
        '측정 완료';

    return '$name ${value.toStringAsFixed(1)} $unit — $statusText입니다.';
  }

  static List<String> _generateRecommendations(String biomarker, String risk) {
    if (risk == 'normal') {
      return ['현재 건강 상태가 양호합니다.', '규칙적인 측정을 유지해주세요.'];
    }
    switch (biomarker.toLowerCase()) {
      case 'glucose':
        return [
          '식사 후 30분 가벼운 산책을 추천합니다.',
          '정제 탄수화물 섭취를 줄여보세요.',
          '다음 측정은 공복 상태에서 해주세요.',
        ];
      case 'uric_acid':
        return [
          '수분 섭취를 충분히 해주세요 (하루 2L 이상).',
          '퓨린 함량이 높은 음식을 자제해주세요.',
        ];
      default:
        return ['의료 전문가와 상담을 권장합니다.', '정기적인 추적 측정이 필요합니다.'];
    }
  }
}

// ============================================================================
// NativeRustEngine — 네이티브 FFI 구현 (H1 모듈 독립성)
// ============================================================================

/// 네이티브 Rust 엔진 구현.
///
/// flutter_rust_bridge를 통해 실제 Rust 라이브러리를 호출합니다.
/// 이 클래스는 [RustBridge] 싱글톤에 위임하지만,
/// DI(Dependency Injection) 시나리오에서 독립적으로 사용할 수 있습니다.
class NativeRustEngine implements RustEngineInterface {
  NativeRustEngine._();
  static final NativeRustEngine _instance = NativeRustEngine._();
  static NativeRustEngine get instance => _instance;

  @override
  Future<void> initialize() => RustBridge.init();

  @override
  String get engineVersion => RustBridge.instance.engineVersion;

  @override
  bool get isNativeEnabled => RustBridge.instance.isNativeEnabled;

  @override
  RustBridgeDiagnostics get diagnostics => RustBridge.instance.diagnostics;

  @override
  Future<List<DeviceInfoDto>> bleScan() => RustBridge.instance.bleScan();

  @override
  Future<bool> bleConnect(String deviceId) =>
      RustBridge.instance.bleConnect(deviceId);

  @override
  Future<void> bleDisconnect(String deviceId) =>
      RustBridge.instance.bleDisconnect(deviceId);

  @override
  Future<int> bleReadBattery(String deviceId) =>
      RustBridge.instance.bleReadBattery(deviceId);

  @override
  Future<CartridgeInfoDto> nfcReadCartridge() =>
      RustBridge.instance.nfcReadCartridge();

  @override
  Future<DifferentialCorrectionDto> differentialMeasureSingle(
    double sDet,
    double sRef,
    double alpha,
  ) =>
      RustBridge.instance.differentialMeasureSingle(sDet, sRef, alpha);

  @override
  Future<List<double>> differentialMeasure(
    List<double> sDet,
    List<double> sRef,
    double alpha,
  ) =>
      RustBridge.instance.differentialMeasure(sDet, sRef, alpha);

  @override
  Future<FingerprintDto> createFingerprint88(List<double> data) =>
      RustBridge.instance.createFingerprint88(data);

  @override
  Future<FingerprintDto> createFingerprint896(List<double> data) =>
      RustBridge.instance.createFingerprint896(data);

  @override
  Future<MeasurementPipelineDto> runMeasurementPipeline({
    required List<double> sDet,
    required List<double> sRef,
    required double alpha,
    required String biomarker,
    required String unit,
  }) =>
      RustBridge.instance.runMeasurementPipeline(
        sDet: sDet,
        sRef: sRef,
        alpha: alpha,
        biomarker: biomarker,
        unit: unit,
      );

  @override
  Future<AiAnalysisResultDto> analyzeMeasurement(
    double value,
    String biomarker,
  ) =>
      RustBridge.instance.analyzeMeasurement(value, biomarker);
}

// ============================================================================
// StubRustEngine — 순수 Dart 스텁 구현 (H1 모듈 독립성)
// ============================================================================

/// 스텁 Rust 엔진 구현.
///
/// 네이티브 Rust 라이브러리 없이 Dart-only로 동작합니다.
/// 테스트, Web, Desktop 환경에서 사용됩니다.
class StubRustEngine implements RustEngineInterface {
  StubRustEngine._();
  static final StubRustEngine _instance = StubRustEngine._();
  static StubRustEngine get instance => _instance;

  bool _initialized = false;

  @override
  Future<void> initialize() async {
    _initialized = true;
  }

  @override
  String get engineVersion => '0.1.0-stub';

  @override
  bool get isNativeEnabled => false;

  @override
  RustBridgeDiagnostics get diagnostics => RustBridgeDiagnostics(
        initialized: _initialized,
        nativeEnabled: false,
        nativePlatform: false,
        releaseBuild: kReleaseMode,
        engineVersion: engineVersion,
        releaseStubAllowed: true,
      );

  @override
  Future<List<DeviceInfoDto>> bleScan() async {
    await Future.delayed(const Duration(milliseconds: 800));
    return const [
      DeviceInfoDto(
        deviceId: 'MPK-STUB-001',
        name: 'ManPaSik Stub X1',
        rssi: -45,
        state: 'discovered',
      ),
    ];
  }

  @override
  Future<bool> bleConnect(String deviceId) async {
    await Future.delayed(const Duration(milliseconds: 300));
    return true;
  }

  @override
  Future<void> bleDisconnect(String deviceId) async {
    await Future.delayed(const Duration(milliseconds: 100));
  }

  @override
  Future<int> bleReadBattery(String deviceId) async => 95;

  @override
  Future<CartridgeInfoDto> nfcReadCartridge() async {
    return const CartridgeInfoDto(
      cartridgeId: 'CART-STUB-001',
      cartridgeType: 'Glucose',
      lotId: 'LOT-STUB',
      expiryDate: '20270630',
      remainingUses: 10,
    );
  }

  @override
  Future<DifferentialCorrectionDto> differentialMeasureSingle(
    double sDet,
    double sRef,
    double alpha,
  ) async {
    return DifferentialCorrectionDto(
      sDet: sDet,
      sRef: sRef,
      alpha: alpha,
      corrected: sDet - alpha * sRef,
    );
  }

  @override
  Future<List<double>> differentialMeasure(
    List<double> sDet,
    List<double> sRef,
    double alpha,
  ) async {
    final n = sDet.length < sRef.length ? sDet.length : sRef.length;
    return List.generate(n, (i) => sDet[i] - alpha * sRef[i]);
  }

  @override
  Future<FingerprintDto> createFingerprint88(List<double> data) async {
    final values = data.length >= 88
        ? data.sublist(0, 88)
        : [...data, ...List.filled(88 - data.length, 0.0)];
    return FingerprintDto(
      values: values,
      dimensions: 88,
      algorithm: 'stub-v1',
    );
  }

  @override
  Future<FingerprintDto> createFingerprint896(List<double> data) async {
    final values = data.length >= 896
        ? data.sublist(0, 896)
        : [...data, ...List.filled(896 - data.length, 0.0)];
    return FingerprintDto(
      values: values,
      dimensions: 896,
      algorithm: 'stub-v1',
    );
  }

  @override
  Future<MeasurementPipelineDto> runMeasurementPipeline({
    required List<double> sDet,
    required List<double> sRef,
    required double alpha,
    required String biomarker,
    required String unit,
  }) async {
    final measurement = MeasurementResultDto(
      primaryValue: 95.0,
      referenceValue: 0.5,
      differentialValue: 94.5,
      snr: 35.0,
      confidence: 0.92,
      unit: unit,
      biomarker: biomarker,
      timestamp: DateTime.now(),
    );
    final analysis = AiAnalysisDto(
      riskLevel: 'normal',
      healthScore: 90.0,
      summary: 'Stub 측정 완료',
      recommendations: const ['스텁 모드 — 참고용 결과입니다.'],
      trend: 'stable',
    );
    return MeasurementPipelineDto(
      measurement: measurement,
      analysis: analysis,
      pipelineDurationMs: 50.0,
    );
  }

  @override
  Future<AiAnalysisResultDto> analyzeMeasurement(
    double value,
    String biomarker,
  ) async {
    return AiAnalysisResultDto(
      riskLevel: 'normal',
      healthScore: 85.0,
      summary: 'Stub AI 분석 결과',
      recommendations: const ['스텁 모드 — 참고용 결과입니다.'],
      trend: 'stable',
    );
  }
}
