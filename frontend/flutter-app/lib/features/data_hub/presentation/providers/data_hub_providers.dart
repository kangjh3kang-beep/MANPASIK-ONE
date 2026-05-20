import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/config/feature_flags.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';

/// Demo Mode 활성 여부.
///
/// FeatureFlags가 true면 → useRealBackend로 우선 (FeatureFlag 우선순위).
/// FeatureFlags 미설정 시 기존 authState.isDemo 폴백 (점진 마이그레이션).
bool _isDemoMode(AuthState authState) {
  // FeatureFlags가 명시적으로 useRealBackend=true면 Demo 비활성화
  if (FeatureFlags.instance.useRealBackend) {
    return false;
  }
  // 그 외에는 기존 Demo 모드 동작 유지
  return authState.isDemo;
}

// ════════════════════════════════════════════════════════════════════════════
//  DataHub Providers — Riverpod 기반 상태 관리
//
//  구조: StateProvider (UI 선택) → FutureProvider (비동기 데이터)
//        → derived Provider (파생 계산)
// ════════════════════════════════════════════════════════════════════════════

/// 선택된 바이오마커 메트릭 (필터 칩)
final dataHubSelectedMetricProvider = StateProvider<String>((ref) => '혈당');

/// 선택된 기간 (필터 칩)
final dataHubSelectedPeriodProvider = StateProvider<String>((ref) => '주간');

/// 선택된 벤토 카드 (Progressive Disclosure — null이면 비선택)
final dataHubSelectedBentoCardProvider = StateProvider<String?>((ref) => null);

/// 기간 → DateTime 범위 변환
({DateTime from, DateTime to}) _periodToRange(String period) {
  final now = DateTime.now();
  switch (period) {
    case '일간':
      return (from: now.subtract(const Duration(days: 1)), to: now);
    case '월간':
      return (from: now.subtract(const Duration(days: 30)), to: now);
    case '주간':
    default:
      return (from: now.subtract(const Duration(days: 7)), to: now);
  }
}

/// 트렌드 데이터 (metric + period 조합)
final dataHubTrendProvider =
    FutureProvider.family<List<TrendDataPoint>, ({String metric, String period})>(
        (ref, params) async {
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (_isDemoMode(authState)) {
    final range = _periodToRange(params.period);
    final days = range.to.difference(range.from).inDays;
    return List.generate(days.clamp(1, 30), (i) {
      final base = params.metric == '혈당'
          ? 95.0
          : params.metric == '요산'
              ? 5.5
              : params.metric == '케톤'
                  ? 0.8
                  : params.metric == '심박'
                      ? 72.0
                      : 36.5;
      final variance = (i % 3 - 1) * (base * 0.05);
      return TrendDataPoint(
        timestamp: range.from.add(Duration(days: i)),
        value: base + variance,
        unit: params.metric == '혈당'
            ? 'mg/dL'
            : params.metric == '요산'
                ? 'mg/dL'
                : params.metric == '케톤'
                    ? 'mmol/L'
                    : params.metric == '심박'
                        ? 'bpm'
                        : '°C',
        biomarkerType: params.metric,
        isWithinRange: true,
      );
    });
  }

  final range = _periodToRange(params.period);
  try {
    return await ref.read(dataHubRepositoryProvider).getTrendData(
          biomarkerType: params.metric,
          from: range.from,
          to: range.to,
        );
  } catch (_) {
    return [];
  }
});

/// 전체 바이오마커 요약 (30초 폴링)
final dataHubSummariesProvider = StreamProvider<List<BiomarkerSummary>>((ref) {
  final controller = StreamController<List<BiomarkerSummary>>();

  Future<void> fetch() async {
    final authState = ref.read(authProvider);

    if (_isDemoMode(authState)) {
      if (!controller.isClosed) {
        controller.add(const [
          BiomarkerSummary(
              biomarkerType: '혈당', displayName: '혈당', latestValue: 98.2,
              averageValue: 94.7, minValue: 85, maxValue: 124,
              unit: 'mg/dL', referenceMin: 70, referenceMax: 140,
              totalMeasurements: 56, trend: 'stable'),
          BiomarkerSummary(
              biomarkerType: '요산', displayName: '요산', latestValue: 5.8,
              averageValue: 5.5, minValue: 4.2, maxValue: 7.1,
              unit: 'mg/dL', referenceMin: 3.4, referenceMax: 7.0,
              totalMeasurements: 42, trend: 'rising'),
          BiomarkerSummary(
              biomarkerType: '케톤', displayName: '케톤', latestValue: 0.6,
              averageValue: 0.8, minValue: 0.3, maxValue: 1.5,
              unit: 'mmol/L', referenceMin: 0, referenceMax: 3.0,
              totalMeasurements: 38, trend: 'falling'),
          BiomarkerSummary(
              biomarkerType: '심박', displayName: '심박', latestValue: 72,
              averageValue: 74, minValue: 58, maxValue: 95,
              unit: 'bpm', referenceMin: 60, referenceMax: 100,
              totalMeasurements: 120, trend: 'stable'),
          BiomarkerSummary(
              biomarkerType: '체온', displayName: '체온', latestValue: 36.5,
              averageValue: 36.4, minValue: 36.0, maxValue: 37.2,
              unit: '°C', referenceMin: 36.0, referenceMax: 37.5,
              totalMeasurements: 56, trend: 'stable'),
        ]);
      }
      return;
    }

    try {
      final repo = ref.read(dataHubRepositoryProvider);
      final data = await repo.getAllBiomarkerSummaries();
      if (!controller.isClosed) controller.add(data);
    } catch (e) {
      if (!controller.isClosed) controller.addError(e);
    }
  }

  fetch();
  final timer = Timer.periodic(const Duration(seconds: 30), (_) => fetch());

  ref.onDispose(() {
    timer.cancel();
    controller.close();
  });

  return controller.stream;
});

/// 종합 건강 점수 (바이오마커 요약에서 파생)
final dataHubHealthScoreProvider = Provider<int>((ref) {
  final summaries = ref.watch(dataHubSummariesProvider);
  return summaries.when(
    data: (list) {
      if (list.isEmpty) return 0;
      int totalScore = 0;
      for (final s in list) {
        if (s.latestValue == null) continue;
        final range = s.referenceMax - s.referenceMin;
        if (range <= 0) continue;
        final mid = (s.referenceMin + s.referenceMax) / 2;
        final deviation = (s.latestValue! - mid).abs() / (range / 2);
        totalScore += ((1.0 - deviation.clamp(0.0, 1.0)) * 100).round();
      }
      final withValues = list.where((s) => s.latestValue != null).length;
      return withValues > 0 ? totalScore ~/ withValues : 0;
    },
    loading: () => 0,
    error: (_, __) => 0,
  );
});

/// 총 측정 횟수
final dataHubTotalMeasurementsProvider = FutureProvider<int>((ref) async {
  final authState = ref.watch(authProvider);
  if (_isDemoMode(authState)) return 312;

  try {
    return await ref.read(dataHubRepositoryProvider).getTotalMeasurementCount();
  } catch (_) {
    return 0;
  }
});

/// 선택된 메트릭에 대한 요약 (파생)
final dataHubSelectedSummaryProvider = Provider<BiomarkerSummary?>((ref) {
  final metric = ref.watch(dataHubSelectedMetricProvider);
  final summaries = ref.watch(dataHubSummariesProvider);
  return summaries.when(
    data: (list) {
      try {
        return list.firstWhere((s) => s.biomarkerType == metric);
      } catch (_) {
        return null;
      }
    },
    loading: () => null,
    error: (_, __) => null,
  );
});
