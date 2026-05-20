import 'package:flutter_test/flutter_test.dart';

import 'package:manpasik/core/config/feature_flags.dart';

void main() {
  setUp(() {
    FeatureFlags.instance.resetOverrides();
  });

  tearDown(() {
    FeatureFlags.instance.resetOverrides();
  });

  group('FeatureFlags', () {
    test('singleton 인스턴스', () {
      final a = FeatureFlags.instance;
      final b = FeatureFlags.instance;
      expect(identical(a, b), true);
    });

    test('기본 디버그/테스트 모드는 demo', () {
      // 테스트 환경(kReleaseMode=false) → demo
      expect(FeatureFlags.instance.isDemoMode, true);
      expect(FeatureFlags.instance.useRealBackend, false);
    });

    test('useRealBackend 오버라이드 적용', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      expect(FeatureFlags.instance.useRealBackend, true);
      expect(FeatureFlags.instance.isDemoMode, false);
    });

    test('useRealBackend=true 시 medical domain 자동 활성', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      expect(FeatureFlags.instance.enableMedicalDomain, true);
      expect(FeatureFlags.instance.enableExternalChannels, true);
      expect(FeatureFlags.instance.enableSLAMonitoring, true);
    });

    test('명시적 medical domain 오버라이드', () {
      FeatureFlags.instance.setOverride(
        useRealBackend: false,
        enableMedicalDomain: true, // 명시적 활성
      );
      expect(FeatureFlags.instance.useRealBackend, false);
      expect(FeatureFlags.instance.enableMedicalDomain, true);
    });

    test('Rust 네이티브 별도 제어', () {
      FeatureFlags.instance.setOverride(
        useRealBackend: true,
        enableRustNative: false,
      );
      expect(FeatureFlags.instance.useRealBackend, true);
      expect(FeatureFlags.instance.enableRustNative, false);
    });

    test('snapshot 반환 형식', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      final snap = FeatureFlags.instance.snapshot();
      expect(snap['useRealBackend'], true);
      expect(snap.containsKey('enableMedicalDomain'), true);
      expect(snap.containsKey('isDemoMode'), true);
    });

    test('resetOverrides 후 기본값 복원', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      expect(FeatureFlags.instance.useRealBackend, true);

      FeatureFlags.instance.resetOverrides();
      expect(FeatureFlags.instance.useRealBackend, false); // 디버그 기본값
    });

    test('isDemoMode = !useRealBackend', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      expect(FeatureFlags.instance.isDemoMode, false);

      FeatureFlags.instance.setOverride(useRealBackend: false);
      expect(FeatureFlags.instance.isDemoMode, true);
    });

    test('각 플래그 독립적 제어', () {
      FeatureFlags.instance.setOverride(
        useRealBackend: false,
        enableRustNative: true,
        enableMedicalDomain: true,
        enableExternalChannels: false,
        enableSLAMonitoring: true,
      );

      expect(FeatureFlags.instance.useRealBackend, false);
      expect(FeatureFlags.instance.enableRustNative, true);
      expect(FeatureFlags.instance.enableMedicalDomain, true);
      expect(FeatureFlags.instance.enableExternalChannels, false);
      expect(FeatureFlags.instance.enableSLAMonitoring, true);
    });

    test('null 오버라이드는 dart-define/기본값 폴백', () {
      FeatureFlags.instance.setOverride(useRealBackend: true);
      FeatureFlags.instance.setOverride(useRealBackend: null);
      // null로 명시적 reset → 기본값
      expect(FeatureFlags.instance.useRealBackend, false);
    });
  });
}
