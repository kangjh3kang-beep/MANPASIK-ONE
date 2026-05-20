import 'package:flutter/foundation.dart';

/// Feature Flag 시스템
///
/// 기존 Demo Mode 자동 폴백을 환경 기반 명시적 플래그로 교체합니다.
/// 모세혈관 워크플로우 검증을 위해 실제 API 호출 경로를 활성화할 수 있습니다.
///
/// 우선순위:
///   1. 컴파일 타임 dart-define (--dart-define=USE_REAL_BACKEND=true)
///   2. 런타임 환경변수 (테스트 시 setOverride)
///   3. 기본값 (kReleaseMode면 true, 디버그/테스트면 false)
class FeatureFlags {
  FeatureFlags._();

  static final FeatureFlags _instance = FeatureFlags._();

  /// 인스턴스 접근자.
  static FeatureFlags get instance => _instance;

  // 컴파일 타임 dart-define 값들
  static const bool _useRealBackendCompile = bool.fromEnvironment(
    'USE_REAL_BACKEND',
    defaultValue: false,
  );
  static const bool _enableRustNativeCompile = bool.fromEnvironment(
    'ENABLE_RUST_NATIVE',
    defaultValue: false,
  );
  static const bool _enableMedicalDomainCompile = bool.fromEnvironment(
    'ENABLE_MEDICAL_DOMAIN',
    defaultValue: false,
  );

  // 런타임 오버라이드 (테스트/QA 용)
  bool? _useRealBackendOverride;
  bool? _enableRustNativeOverride;
  bool? _enableMedicalDomainOverride;
  bool? _enableExternalChannelsOverride;
  bool? _enableSLAMonitoringOverride;

  /// 실제 백엔드 호출 활성화.
  ///
  /// false면 기존 Demo Mode 폴백 동작 유지 (Mock 데이터).
  /// true면 모든 REST 호출이 실제 서버로 전송됩니다.
  bool get useRealBackend {
    if (_useRealBackendOverride != null) return _useRealBackendOverride!;
    if (_useRealBackendCompile) return true;
    return kReleaseMode;
  }

  /// Rust 네이티브 FFI 활성화.
  ///
  /// false면 Dart 측 시뮬레이션 폴백.
  /// true면 rust-core/flutter-bridge 직접 호출.
  bool get enableRustNative {
    if (_enableRustNativeOverride != null) return _enableRustNativeOverride!;
    if (_enableRustNativeCompile) return true;
    return false;
  }

  /// 의료 도메인 모듈 (POCT/FHIR/Polaris/CoT/Scribe) 호출 활성화.
  ///
  /// false면 화상진료 시 의료 분석 생략.
  /// true면 telemedicine-service의 MedicalIntegrator 호출.
  bool get enableMedicalDomain {
    if (_enableMedicalDomainOverride != null) {
      return _enableMedicalDomainOverride!;
    }
    if (_enableMedicalDomainCompile) return true;
    return useRealBackend; // 실 백엔드 활성화 시 자동 활성
  }

  /// 외부 채널 (SMS/Email/FCM/Kakao) 알림 활성화.
  bool get enableExternalChannels {
    if (_enableExternalChannelsOverride != null) {
      return _enableExternalChannelsOverride!;
    }
    return useRealBackend;
  }

  /// STAT 30초 SLA 모니터링 활성화.
  bool get enableSLAMonitoring {
    if (_enableSLAMonitoringOverride != null) {
      return _enableSLAMonitoringOverride!;
    }
    return useRealBackend;
  }

  /// 데모 모드 활성화 여부 (legacy).
  ///
  /// useRealBackend의 반대값. 점진적 마이그레이션을 위해 유지.
  bool get isDemoMode => !useRealBackend;

  /// 테스트 전용: 런타임 오버라이드 설정.
  @visibleForTesting
  void setOverride({
    bool? useRealBackend,
    bool? enableRustNative,
    bool? enableMedicalDomain,
    bool? enableExternalChannels,
    bool? enableSLAMonitoring,
  }) {
    _useRealBackendOverride = useRealBackend;
    _enableRustNativeOverride = enableRustNative;
    _enableMedicalDomainOverride = enableMedicalDomain;
    _enableExternalChannelsOverride = enableExternalChannels;
    _enableSLAMonitoringOverride = enableSLAMonitoring;
  }

  /// 모든 오버라이드 초기화 (테스트 cleanup).
  @visibleForTesting
  void resetOverrides() {
    _useRealBackendOverride = null;
    _enableRustNativeOverride = null;
    _enableMedicalDomainOverride = null;
    _enableExternalChannelsOverride = null;
    _enableSLAMonitoringOverride = null;
  }

  /// 현재 활성 플래그 요약 (디버그용).
  Map<String, bool> snapshot() {
    return {
      'useRealBackend': useRealBackend,
      'enableRustNative': enableRustNative,
      'enableMedicalDomain': enableMedicalDomain,
      'enableExternalChannels': enableExternalChannels,
      'enableSLAMonitoring': enableSLAMonitoring,
      'isDemoMode': isDemoMode,
    };
  }
}
