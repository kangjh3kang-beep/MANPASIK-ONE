/// Rust FFI Bridge — ManPaSik Core Engine 연동
///
/// flutter_rust_bridge 기반 네이티브-스텁 하이브리드 모드:
/// - 모바일 (Android/iOS): 네이티브 Rust 엔진 사용 (라이브러리 존재 시)
/// - Web/Desktop: 스텁 모드 (시뮬레이션)
/// - 네이티브 로드 실패 시: 자동 스텁 폴백
library;

import 'dart:io' show Platform;
import 'dart:math';

import 'package:flutter/foundation.dart' show kIsWeb;

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

  static String get engineVersion =>
      _useNative ? '1.0.0-native' : '1.0.0-stub';

  /// 플랫폼 감지: 네이티브 가능 여부
  static bool get _isNativePlatform {
    if (kIsWeb) return false;
    return Platform.isAndroid || Platform.isIOS;
  }

  /// Rust Core 엔진 초기화
  ///
  /// AppConfig.enableRustFfi가 true이고 모바일 플랫폼일 때
  /// 네이티브 라이브러리 로드를 시도합니다.
  /// flutter_rust_bridge 패키지 빌드 완료 후 주석을 해제하면 자동 활성화됩니다.
  static Future<void> init() async {
    if (_instance._initialized) return;

    if (_isNativePlatform) {
      try {
        // flutter_rust_bridge 빌드 완료 후 아래 주석 해제:
        // if (AppConfig.enableRustFfi) {
        //   await RustLib.init();
        //   _useNative = true;
        // }
        _useNative = false; // 네이티브 라이브러리 준비 전까지 스텁 유지
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
      // return await api.bleScanDevices(timeoutMs: timeout.inMilliseconds);
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
      // return await api.bleConnectDevice(deviceId: deviceId);
    }
    await Future.delayed(const Duration(milliseconds: 600));
    return true;
  }

  /// BLE 디바이스 연결 해제
  static Future<void> bleDisconnect(String deviceId) async {
    if (_useNative) {
      // await api.bleDisconnectDevice(deviceId: deviceId);
      return;
    }
    await Future.delayed(const Duration(milliseconds: 200));
  }

  /// BLE 배터리 레벨 읽기
  static Future<int> bleReadBattery(String deviceId) async {
    if (_useNative) {
      // return await api.bleReadBattery(deviceId: deviceId);
    }
    await Future.delayed(const Duration(milliseconds: 100));
    return 85 + Random().nextInt(15);
  }

  /// BLE 연결 품질 (RSSI 기반)
  static Future<String> bleConnectionQuality(String deviceId) async {
    if (_useNative) {
      // return await api.bleConnectionQuality(deviceId: deviceId);
    }
    return 'excellent'; // excellent, good, fair, poor
  }

  // ── NFC API ──

  /// NFC 카트리지 읽기
  static Future<CartridgeInfoDto> nfcReadCartridge() async {
    if (_useNative) {
      // return await api.nfcReadCartridge();
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
      // return await api.processDifferentialMeasurement(
      //   sDet: signalData,
      //   sRef: referenceData,
      //   alpha: 0.95,
      //   biomarker: biomarker,
      //   unit: unit,
      // );
    }

    // 스텁: 시뮬레이션된 차동 계측
    await Future.delayed(const Duration(milliseconds: 300));
    final rng = Random(DateTime.now().microsecondsSinceEpoch);

    final signal = signalData.isNotEmpty
        ? signalData.reduce((a, b) => a + b) / signalData.length
        : 85.0 + rng.nextDouble() * 30;
    final reference = referenceData.isNotEmpty
        ? referenceData.reduce((a, b) => a + b) / referenceData.length
        : 2.0 + rng.nextDouble() * 0.5;
    final differential = signal - reference;

    return MeasurementResultDto(
      primaryValue: differential,
      referenceValue: reference,
      differentialValue: differential,
      snr: 35.0 + rng.nextDouble() * 15,
      confidence: 0.92 + rng.nextDouble() * 0.07,
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
      // return await api.aiAnalyzeMeasurement(
      //   value: value,
      //   biomarker: biomarker,
      //   unit: unit,
      //   recentValues: recentValues,
      // );
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
      // 네이티브 파이프라인: Rust에서 BLE→DSP→AI 일괄 처리
      // return await api.runFullPipeline(
      //   deviceId: deviceId,
      //   biomarker: biomarker,
      //   unit: unit,
      // );
    }

    // 스텁: 단계별 시뮬레이션
    // 1단계: BLE 데이터 수집 시뮬레이션
    final rng = Random(DateTime.now().microsecondsSinceEpoch);
    final signalData = List.generate(88, (_) => 0.5 + rng.nextDouble());
    final referenceData = List.generate(88, (_) => 0.01 + rng.nextDouble() * 0.03);
    await Future.delayed(const Duration(milliseconds: 500));

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

  static List<String> _generateRecommendations(
      String biomarker, String risk) {
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
