/// Rust FFI Bridge — ManPaSik Core Engine 연동
///
/// flutter_rust_bridge 기반 네이티브-스텁 하이브리드 모드:
/// - 모바일 (Android/iOS): 네이티브 Rust 엔진 사용 (라이브러리 존재 시)
/// - Web/Desktop: 스텁 모드 (시뮬레이션)
/// - 네이티브 로드 실패 시: 자동 스텁 폴백
library;

import 'dart:io' show Platform;
import 'dart:math';

import 'package:flutter/foundation.dart' show kIsWeb, kReleaseMode;

import 'rust_ffi_native_stub.dart'
    if (dart.library.io) 'rust_ffi_native_impl.dart' as native_engine;

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

/// 측정 파이프라인 결과 (BLE → DSP → AI 전체 결과)
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

/// ManPaSik Rust Core Engine Bridge
///
/// 싱글톤 패턴. [init]으로 초기화 후 사용.
/// 플랫폼 감지 기반 네이티브/스텁 자동 전환.
class RustBridge {
  RustBridge._();
  static final RustBridge _instance = RustBridge._();
  static RustBridge get instance => _instance;

  bool _initialized = false;
  bool get isInitialized => _initialized;

  /// 네이티브 Rust 엔진 사용 가능 여부
  /// 모바일 플랫폼이고 라이브러리 로드 성공 시 true
  static bool _useNative = false;

  /// 네이티브 엔진 사용 여부 조회
  static bool get isNativeEnabled => _useNative;

  static const bool _releaseStubAllowed = bool.fromEnvironment(
    'MANPASIK_RELEASE_ALLOW_STUB_MEASUREMENT',
    defaultValue: false,
  );

  static RustBridgeDiagnostics get diagnostics => RustBridgeDiagnostics(
        initialized: _instance._initialized,
        nativeEnabled: _useNative,
        nativePlatform: _isNativePlatform,
        releaseBuild: kReleaseMode,
        engineVersion: engineVersion,
        releaseStubAllowed: _releaseStubAllowed,
      );

  /// 엔진 버전 (네이티브: Rust Cargo.toml 기준, 스텁: 시뮬레이션)
  static String get engineVersion =>
      _useNative ? native_engine.nativeGetEngineVersion() : '0.1.0-stub';

  /// 플랫폼 감지: 네이티브 가능 여부
  static bool get _isNativePlatform {
    if (kIsWeb) return false;
    return Platform.isAndroid || Platform.isIOS;
  }

  /// Rust Core 엔진 초기화
  ///
  /// 모바일 플랫폼에서 네이티브 라이브러리 로드를 시도합니다.
  /// Rust 라이브러리가 링크되어 있으면 자동으로 네이티브 모드 활성화.
  ///
  /// 네이티브 활성화 조건:
  ///   1. Android/iOS 플랫폼
  ///   2. `cargo build` 완료된 Rust 라이브러리 (.so/.dylib) 링크
  ///   3. `flutter_rust_bridge_codegen generate` 완료 (frb_generated/ 존재)
  ///
  /// 활성화 방법 (Phase D 완료 시):
  ///   import 'frb_generated/frb_generated.dart' show RustLib;
  ///   await RustLib.init();
  ///   _useNative = true;
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

  // ── BLE API ──

  /// BLE 디바이스 스캔 (주변 ManPaSik 디바이스 검색)
  static Future<List<DeviceInfoDto>> bleScan({
    Duration timeout = const Duration(seconds: 5),
  }) async {
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

  /// BLE 디바이스 연결
  static Future<bool> bleConnect(String deviceId) async {
    if (_useNative) {
      final result = await native_engine.nativeBleConnect(deviceId);
      if (result != null) return result;
    }
    await Future.delayed(const Duration(milliseconds: 600));
    return true;
  }

  /// BLE 디바이스 연결 해제
  static Future<void> bleDisconnect(String deviceId) async {
    // BLE disconnect는 현재 frb_generated API에 미포함 — 항상 스텁
    await Future.delayed(const Duration(milliseconds: 200));
  }

  /// BLE 배터리 레벨 읽기
  static Future<int> bleReadBattery(String deviceId) async {
    if (_useNative) {
      final result = await native_engine.nativeBleReadBattery(deviceId);
      if (result != null) return result;
    }
    await Future.delayed(const Duration(milliseconds: 100));
    return 85 + Random().nextInt(15);
  }

  /// BLE 연결 품질 (RSSI 기반)
  static Future<String> bleConnectionQuality(String deviceId) async {
    if (_useNative) {
      final result = native_engine.nativeBleConnectionQuality(-45);
      if (result != null) return result;
    }
    return 'excellent'; // excellent, good, fair, poor
  }

  // ── NFC API ──

  /// NFC 카트리지 읽기
  static Future<CartridgeInfoDto> nfcReadCartridge() async {
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

  // ── 측정 엔진 API ──

  /// 차동 계측 처리 (시그널 + 레퍼런스 → 보정된 결과)
  static Future<MeasurementResultDto> processMeasurement({
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

  // ── AI 분석 API ──

  /// 측정 결과 AI 분석 (위험도 판정 + 코칭 추천)
  static Future<AiAnalysisDto> analyzeResult({
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

  // ── 측정 파이프라인 (BLE → DSP → AI 통합) ──

  /// 전체 측정 파이프라인 실행
  ///
  /// 1. BLE에서 원시 데이터 수집 (시뮬레이션)
  /// 2. 차동 계측 처리 (S_det - α × S_ref)
  /// 3. AI 분석 (위험도 + 코칭)
  static Future<MeasurementPipelineResult> runMeasurementPipeline({
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

    // Dart-side 파이프라인: BLE 데이터 수집 → 차동 계측 → AI 분석
    // 1단계: BLE 데이터 수집 (네이티브 미사용 시 센서 시뮬레이션)
    final rng = Random(DateTime.now().microsecondsSinceEpoch);
    // 88채널 × 4센서 = 실제 신호/레퍼런스 생성 (랜덤이 아닌 구조적)
    final signalData = List.generate(88, (i) {
      final base = 0.5 + 0.3 * sin(i * 0.15);
      return base + rng.nextDouble() * 0.05;
    });
    final referenceData = List.generate(88, (i) {
      return 0.02 + 0.01 * sin(i * 0.15) + rng.nextDouble() * 0.005;
    });
    await Future.delayed(const Duration(milliseconds: 100));

    // 2단계: 차동 계측 처리
    final measurement = await processMeasurement(
      signalData: signalData,
      referenceData: referenceData,
      biomarker: biomarker,
      unit: unit,
    );

    // 3단계: AI 분석
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

  // ── 내부 스텁 헬퍼 ──

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
