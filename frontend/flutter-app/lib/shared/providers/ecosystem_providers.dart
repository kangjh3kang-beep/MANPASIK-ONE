import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/data_hub/presentation/providers/data_hub_providers.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';

// ═══════════════════════════════════════════════════════════════════════════
//  Ecosystem Integration Providers — 크로스 도메인 집계 레이어
//
//  기존 도메인별 Provider를 집계하여 홈 대시보드와 도메인 간 데이터 흐름을
//  연결합니다. 새로운 API 호출 없이 기존 Provider를 조합합니다.
// ═══════════════════════════════════════════════════════════════════════════

/// 측정 완료 이벤트 카운터.
/// 측정 완료 시 increment → 의존하는 Provider들이 자동 갱신됩니다.
final measurementCompletionProvider = StateProvider<int>((ref) => 0);

/// 홈 대시보드 통합 데이터.
/// 기존 Provider들을 집계하여 홈 화면에 필요한 모든 데이터를 제공합니다.
class HomeDashboardData {
  final int healthScore;
  final MeasurementHistoryItem? latestMeasurement;
  final bool measurementHistoryStale;
  final String? measurementHistoryError;
  final List<TelemedicineReservation> upcomingReservations;
  final int unreadNotifications;
  final List<FamilyMember> alertMembers;
  final String? aiInsight;

  const HomeDashboardData({
    required this.healthScore,
    this.latestMeasurement,
    this.measurementHistoryStale = false,
    this.measurementHistoryError,
    this.upcomingReservations = const [],
    this.unreadNotifications = 0,
    this.alertMembers = const [],
    this.aiInsight,
  });
}

/// 홈 대시보드 집계 Provider.
/// 6개 도메인의 기존 Provider를 watch하여 단일 데이터로 통합합니다.
final homeDashboardProvider = Provider<HomeDashboardData>((ref) {
  // measurementCompletion을 watch하여 측정 완료 시 자동 갱신
  ref.watch(measurementCompletionProvider);

  // 1. 건강 점수 (DataHub 바이오마커 기반)
  final healthScore = ref.watch(dataHubHealthScoreProvider);

  // 2. 최근 측정
  final measurementAsync = ref.watch(measurementHistoryProvider);
  final latestMeasurement = measurementAsync.whenOrNull(
    data: (result) => result.items.isNotEmpty ? result.items.first : null,
  );
  final measurementHistoryStale = measurementAsync.whenOrNull(
        data: (result) => result.isStale,
      ) ??
      false;
  final measurementHistoryError = measurementAsync.whenOrNull(
    data: (result) => result.errorMessage,
  );

  // 3. 예약 목록
  final reservationsAsync = ref.watch(reservationsProvider);
  final upcomingReservations = reservationsAsync.whenOrNull(
        data: (list) => list
            .where(
              (r) =>
                  r.status == ReservationStatus.confirmed ||
                  r.status == ReservationStatus.pending,
            )
            .toList(),
      ) ??
      [];

  // 4. 미읽은 알림 수
  final notifAsync = ref.watch(unreadNotificationCountProvider);
  final unreadCount = notifAsync.whenOrNull(data: (c) => c) ?? 0;

  // 5. 가족 경고 멤버
  final familyAsync = ref.watch(familyGroupsProvider);
  final alertMembers = <FamilyMember>[];
  familyAsync.whenOrNull(data: (groups) {
    for (final g in groups) {
      for (final m in g.members) {
        if (m.latestHealthStatus == 'alert' ||
            m.latestHealthStatus == 'caution') {
          alertMembers.add(m);
        }
      }
    }
  });

  // 6. AI 인사이트
  final insightAsync = ref.watch(todayInsightProvider);
  final aiInsight = insightAsync.whenOrNull(
    data: (insight) => insight.summary,
  );

  return HomeDashboardData(
    healthScore: healthScore,
    latestMeasurement: latestMeasurement,
    measurementHistoryStale: measurementHistoryStale,
    measurementHistoryError: measurementHistoryError,
    upcomingReservations: upcomingReservations,
    unreadNotifications: unreadCount,
    alertMembers: alertMembers,
    aiInsight: aiInsight,
  );
});

/// 글로벌 배지 카운트 Provider.
/// 하단 독 네비게이션의 배지에 사용됩니다.
final globalBadgeProvider = Provider<Map<String, int>>((ref) {
  ref.watch(measurementCompletionProvider);

  final notifAsync = ref.watch(unreadNotificationCountProvider);
  final notifCount = notifAsync.whenOrNull(data: (c) => c) ?? 0;

  final familyAsync = ref.watch(familyGroupsProvider);
  int familyAlertCount = 0;
  familyAsync.whenOrNull(data: (groups) {
    for (final g in groups) {
      for (final m in g.members) {
        if (m.latestHealthStatus == 'alert' ||
            m.latestHealthStatus == 'caution') {
          familyAlertCount++;
        }
      }
    }
  });

  return {
    'notification': notifCount,
    'family': familyAlertCount,
  };
});
