import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';

/// 건강 점수 계산 로직 (data_hub_providers.dart의 dataHubHealthScoreProvider와 동일)
///
/// 공식: score = (1 - |latestValue - mid| / (range/2)) * 100, clamped to [0, 100]
/// 여러 바이오마커의 평균
int calculateHealthScore(List<BiomarkerSummary> summaries) {
  if (summaries.isEmpty) return 0;
  int totalScore = 0;
  for (final s in summaries) {
    if (s.latestValue == null) continue;
    final range = s.referenceMax - s.referenceMin;
    if (range <= 0) continue;
    final mid = (s.referenceMin + s.referenceMax) / 2;
    final deviation = (s.latestValue! - mid).abs() / (range / 2);
    totalScore += ((1.0 - deviation.clamp(0.0, 1.0)) * 100).round();
  }
  final withValues = summaries.where((s) => s.latestValue != null).length;
  return withValues > 0 ? totalScore ~/ withValues : 0;
}

void main() {
  group('건강 점수 계산', () {
    test('빈 목록 → 0점', () {
      expect(calculateHealthScore([]), 0);
    });

    test('정중앙 값 → 100점', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 90.0, // mid = (70+110)/2 = 90
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 100);
    });

    test('참조 범위 경계 값 → 0점', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 70.0, // 하한 경계 (deviation = 1.0)
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 0);
    });

    test('참조 범위 초과 → 0점 (클램프)', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 200.0, // 범위 초과
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 0);
    });

    test('여러 바이오마커 평균 점수', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 90.0, // mid → 100점
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
        BiomarkerSummary(
          biomarkerType: '체온', displayName: '체온', unit: '°C',
          latestValue: 36.75, // mid = (36+37.5)/2 = 36.75 → 100점
          referenceMin: 36.0, referenceMax: 37.5,
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 100);
    });

    test('latestValue가 null인 항목은 제외', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 90.0, // 100점
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
        BiomarkerSummary(
          biomarkerType: '콜레스테롤', displayName: '콜레스테롤', unit: 'mg/dL',
          // latestValue is null → 제외
          referenceMin: 0, referenceMax: 200,
          totalMeasurements: 0, trend: 'insufficient',
        ),
      ];
      expect(calculateHealthScore(summaries), 100);
    });

    test('모든 latestValue가 null → 0점', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 0, trend: 'insufficient',
        ),
      ];
      expect(calculateHealthScore(summaries), 0);
    });

    test('range가 0인 항목 건너뛰기', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: 'test', displayName: 'test', unit: 'unit',
          latestValue: 50.0,
          referenceMin: 50, referenceMax: 50, // range = 0
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 0);
    });

    test('50% 편차 → 50점', () {
      const summaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', unit: 'mg/dL',
          latestValue: 100.0, // mid=90, range/2=20, deviation=10/20=0.5 → 50점
          referenceMin: 70, referenceMax: 110,
          totalMeasurements: 10, trend: 'stable',
        ),
      ];
      expect(calculateHealthScore(summaries), 50);
    });

    test('데모 데이터 건강 점수는 50~100 사이', () {
      const demoSummaries = [
        BiomarkerSummary(
          biomarkerType: '혈당', displayName: '혈당', latestValue: 98.2,
          averageValue: 94.7, minValue: 85, maxValue: 124,
          unit: 'mg/dL', referenceMin: 70, referenceMax: 140,
          totalMeasurements: 56, trend: 'stable',
        ),
        BiomarkerSummary(
          biomarkerType: '요산', displayName: '요산', latestValue: 5.8,
          averageValue: 5.5, minValue: 4.2, maxValue: 7.1,
          unit: 'mg/dL', referenceMin: 3.4, referenceMax: 7.0,
          totalMeasurements: 42, trend: 'rising',
        ),
        BiomarkerSummary(
          biomarkerType: '케톤', displayName: '케톤', latestValue: 0.6,
          averageValue: 0.8, minValue: 0.3, maxValue: 1.5,
          unit: 'mmol/L', referenceMin: 0, referenceMax: 3.0,
          totalMeasurements: 38, trend: 'falling',
        ),
        BiomarkerSummary(
          biomarkerType: '심박', displayName: '심박', latestValue: 72,
          averageValue: 74, minValue: 58, maxValue: 95,
          unit: 'bpm', referenceMin: 60, referenceMax: 100,
          totalMeasurements: 120, trend: 'stable',
        ),
        BiomarkerSummary(
          biomarkerType: '체온', displayName: '체온', latestValue: 36.5,
          averageValue: 36.4, minValue: 36.0, maxValue: 37.2,
          unit: '°C', referenceMin: 36.0, referenceMax: 37.5,
          totalMeasurements: 56, trend: 'stable',
        ),
      ];
      final score = calculateHealthScore(demoSummaries);
      expect(score, greaterThan(50));
      expect(score, lessThanOrEqualTo(100));
    });
  });

  group('모니터링 요약 로직', () {
    test('빈 기기 목록 → 모든 통계 0', () {
      final stats = computeSummary([]);
      expect(stats.total, 0);
      expect(stats.connected, 0);
      expect(stats.alerts, 0);
      expect(stats.avgBattery, 0);
    });

    test('혼합 상태 기기 통계 계산', () {
      final devices = [
        _Dev('d1', true, 85),
        _Dev('d2', true, 90),
        _Dev('d3', false, 50), // disconnected → alert
        _Dev('d4', true, 15), // battery < 20 → alert
        _Dev('d5', true, 70),
      ];
      final stats = computeSummary(devices);
      expect(stats.total, 5);
      expect(stats.connected, 4); // d1, d2, d4, d5
      expect(stats.alerts, 2);   // d3 (disconnected) + d4 (battery<20)
      expect(stats.avgBattery, 62); // (85+90+50+15+70)/5
    });

    test('모든 기기 연결 + 배터리 정상', () {
      final devices = [
        _Dev('d1', true, 100),
        _Dev('d2', true, 80),
      ];
      final stats = computeSummary(devices);
      expect(stats.alerts, 0);
      expect(stats.avgBattery, 90);
    });

    test('배터리 20%는 경고 아님', () {
      final devices = [_Dev('d1', true, 20)];
      final stats = computeSummary(devices);
      expect(stats.alerts, 0);
    });

    test('배터리 19%는 경고', () {
      final devices = [_Dev('d1', true, 19)];
      final stats = computeSummary(devices);
      expect(stats.alerts, 1);
    });
  });

  group('기기 타입 필터링', () {
    test('전체 탭 → 모든 기기', () {
      final devices = _mixedTypeDevices();
      expect(filterByTab(devices, 0).length, 5);
    });

    test('Bio 탭 필터링', () {
      final filtered = filterByTab(_mixedTypeDevices(), 1);
      expect(filtered.length, 2);
      expect(filtered.every((d) => d.type == 'bio'), isTrue);
    });

    test('Gas 탭 필터링', () {
      final filtered = filterByTab(_mixedTypeDevices(), 2);
      expect(filtered.length, 2);
      expect(filtered.every((d) => d.type == 'gas'), isTrue);
    });

    test('Env 탭 필터링', () {
      final filtered = filterByTab(_mixedTypeDevices(), 3);
      expect(filtered.length, 1);
      expect(filtered.every((d) => d.type == 'env'), isTrue);
    });
  });

  group('경고 기기 필터링', () {
    test('disconnected만 경고', () {
      final alerts = alertDevices([
        _Dev('d1', true, 85),
        _Dev('d2', false, 50),
      ]);
      expect(alerts.length, 1);
      expect(alerts.first.id, 'd2');
    });

    test('배터리 부족만 경고', () {
      final alerts = alertDevices([
        _Dev('d1', true, 10),
        _Dev('d2', true, 50),
      ]);
      expect(alerts.length, 1);
      expect(alerts.first.id, 'd1');
    });

    test('경고 없음', () {
      final alerts = alertDevices([
        _Dev('d1', true, 80),
        _Dev('d2', true, 90),
      ]);
      expect(alerts, isEmpty);
    });
  });
}

// ════════════════════════════════════════════════════════════════════════════
// 테스트용 헬퍼 (Provider 로직 추출한 순수 함수)
// ════════════════════════════════════════════════════════════════════════════

class _Dev {
  final String id;
  final bool isConnected;
  final int battery;
  final String type;
  const _Dev(this.id, this.isConnected, this.battery, [this.type = 'bio']);
}

({int total, int connected, int alerts, int avgBattery}) computeSummary(
    List<_Dev> devices) {
  if (devices.isEmpty) {
    return (total: 0, connected: 0, alerts: 0, avgBattery: 0);
  }
  final connected = devices.where((d) => d.isConnected).length;
  final alerts = devices
      .where((d) => !d.isConnected || d.battery < 20)
      .length;
  final avgBattery =
      (devices.fold<int>(0, (s, d) => s + d.battery) / devices.length).round();
  return (
    total: devices.length,
    connected: connected,
    alerts: alerts,
    avgBattery: avgBattery,
  );
}

List<_Dev> filterByTab(List<_Dev> devices, int tab) {
  switch (tab) {
    case 1: return devices.where((d) => d.type == 'bio').toList();
    case 2: return devices.where((d) => d.type == 'gas').toList();
    case 3: return devices.where((d) => d.type == 'env').toList();
    default: return devices;
  }
}

List<_Dev> alertDevices(List<_Dev> devices) =>
    devices.where((d) => !d.isConnected || d.battery < 20).toList();

List<_Dev> _mixedTypeDevices() => const [
      _Dev('d1', true, 80, 'bio'),
      _Dev('d2', true, 90, 'bio'),
      _Dev('d3', true, 70, 'gas'),
      _Dev('d4', false, 50, 'env'),
      _Dev('d5', true, 95, 'gas'),
    ];
