import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/config/feature_flags.dart';
import 'package:manpasik/core/providers/rust_bridge_provider.dart';

void main() {
  setUp(() {
    FeatureFlags.instance.resetOverrides();
  });

  tearDown(() {
    FeatureFlags.instance.resetOverrides();
  });

  group('rustBridgeProvider', () {
    test('인스턴스 반환', () {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final bridge = container.read(rustBridgeProvider);
      expect(bridge, isNotNull);
    });
  });

  group('rustNativeEnabledProvider', () {
    test('기본은 false (디버그/테스트)', () {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      expect(container.read(rustNativeEnabledProvider), false);
    });

    test('Feature Flag 활성 시 true', () {
      FeatureFlags.instance.setOverride(enableRustNative: true);
      final container = ProviderContainer();
      addTearDown(container.dispose);

      expect(container.read(rustNativeEnabledProvider), true);
    });
  });

  group('rustCapabilityProvider', () {
    test('네이티브 비활성 시 모든 capability false', () {
      FeatureFlags.instance.setOverride(enableRustNative: false);
      final container = ProviderContainer();
      addTearDown(container.dispose);

      expect(container.read(rustCapabilityProvider('differential')), false);
      expect(container.read(rustCapabilityProvider('fingerprint')), false);
      expect(container.read(rustCapabilityProvider('nfc')), false);
    });

    test('네이티브 활성 + 지원 capability', () {
      FeatureFlags.instance.setOverride(enableRustNative: true);
      final container = ProviderContainer();
      addTearDown(container.dispose);

      expect(container.read(rustCapabilityProvider('differential')), true);
      expect(container.read(rustCapabilityProvider('fingerprint')), true);
      expect(container.read(rustCapabilityProvider('nfc')), true);
      expect(container.read(rustCapabilityProvider('ble')), true);
    });

    test('지원하지 않는 capability', () {
      FeatureFlags.instance.setOverride(enableRustNative: true);
      final container = ProviderContainer();
      addTearDown(container.dispose);

      expect(container.read(rustCapabilityProvider('unknown')), false);
    });
  });

  group('differentialCorrectionProvider', () {
    test('정상 차분식 계산', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);
      final result = await svc.apply(
        sN: [10.0, 20.0, 30.0],
        rN: [0.1, 0.2, 0.3],
        alpha: 0.98,
      );

      expect(result.success, true);
      expect(result.value!.corrected.length, 3);
      // 10 - 0.98 * 0.1 = 9.902
      expect(result.value!.corrected[0], closeTo(9.902, 0.001));
      // 20 - 0.98 * 0.2 = 19.804
      expect(result.value!.corrected[1], closeTo(19.804, 0.001));
    });

    test('차분식 + 95% CI', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);
      final result = await svc.apply(
        sN: [100.0],
        rN: [10.0],
        alpha: 1.0,
        stdErrRatio: 0.02,
      );

      expect(result.success, true);
      // value = 100 - 1.0 * 10 = 90
      // stdErr = 90 * 0.02 = 1.8
      // CI = 90 ± 1.96 * 1.8 = [86.47, 93.53]
      expect(result.value!.corrected[0], closeTo(90.0, 0.001));
      expect(result.value!.ciLow[0], closeTo(86.472, 0.01));
      expect(result.value!.ciHigh[0], closeTo(93.528, 0.01));
    });

    test('길이 불일치 거부', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);
      final result = await svc.apply(
        sN: [10.0, 20.0],
        rN: [0.1],
        alpha: 0.98,
      );

      expect(result.success, false);
      expect(result.errorMessage, contains('same length'));
    });

    test('alpha 범위 위반 거부', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);

      // alpha 0.85 < 0.90
      final r1 = await svc.apply(sN: [10.0], rN: [1.0], alpha: 0.85);
      expect(r1.success, false);

      // alpha 1.20 > 1.10
      final r2 = await svc.apply(sN: [10.0], rN: [1.0], alpha: 1.20);
      expect(r2.success, false);
    });

    test('빈 입력 처리', () async {
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);
      final result = await svc.apply(sN: [], rN: []);

      expect(result.success, true);
      expect(result.value!.corrected, isEmpty);
    });

    test('Dart 폴백 사용 시 fromNative=false', () async {
      // 네이티브 비활성
      FeatureFlags.instance.setOverride(enableRustNative: false);
      final container = ProviderContainer();
      addTearDown(container.dispose);

      final svc = container.read(differentialCorrectionProvider);
      final result = await svc.apply(sN: [1.0], rN: [0.1]);

      expect(result.success, true);
      expect(result.value, isNotNull);
      // Dart 폴백은 fromNative=false
      expect(result.fromNative, false);
    });
  });

  group('RustCallResult', () {
    test('success factory', () {
      const r = RustCallResult.success(42);
      expect(r.success, true);
      expect(r.value, 42);
      expect(r.errorMessage, null);
      expect(r.isFailure, false);
    });

    test('failure factory', () {
      const r = RustCallResult<int>.failure('error');
      expect(r.success, false);
      expect(r.value, null);
      expect(r.errorMessage, 'error');
      expect(r.isFailure, true);
    });
  });
}
