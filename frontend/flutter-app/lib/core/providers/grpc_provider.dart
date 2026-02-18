import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/features/auth/data/auth_repository_impl.dart';
import 'package:manpasik/features/auth/data/auth_repository_rest.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/features/devices/data/device_repository_impl.dart';
import 'package:manpasik/features/devices/data/device_repository_rest.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_impl.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_rest.dart';
import 'package:manpasik/features/user/data/user_repository_impl.dart';
import 'package:manpasik/features/user/data/user_repository_rest.dart';
import 'package:manpasik/features/community/domain/community_repository.dart';
import 'package:manpasik/features/community/data/community_repository_rest.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/features/medical/data/medical_repository_rest.dart';
import 'package:manpasik/features/market/domain/market_repository.dart';
import 'package:manpasik/features/market/data/market_repository_rest.dart';
import 'package:manpasik/features/ai_coach/domain/ai_coach_repository.dart';
import 'package:manpasik/features/ai_coach/data/ai_coach_repository_rest.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/features/family/data/family_repository_rest.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/features/data_hub/data/data_hub_repository_rest.dart';

/// REST Gateway Client Provider (웹/네이티브 공용)
final restClientProvider = Provider<ManPaSikRestClient>((ref) {
  return ManPaSikRestClient();
});

/// gRPC 채널 관리자 Provider (네이티브 전용, 웹에서는 사용하지 않음)
final grpcClientManagerProvider = Provider<GrpcClientManager>((ref) {
  return GrpcClientManager();
});

/// Auth Repository Provider
///
/// 웹: REST Gateway, 네이티브: gRPC 직접 연결
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  if (kIsWeb) {
    return AuthRepositoryRest(ref.watch(restClientProvider));
  }
  final manager = ref.watch(grpcClientManagerProvider);
  return AuthRepositoryImpl(manager);
});

/// Device Repository Provider
final deviceRepositoryProvider = Provider<DeviceRepository>((ref) {
  if (kIsWeb) {
    return DeviceRepositoryRest(ref.watch(restClientProvider));
  }
  final manager = ref.watch(grpcClientManagerProvider);
  return DeviceRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

/// Measurement Repository Provider
final measurementRepositoryProvider = Provider<MeasurementRepository>((ref) {
  if (kIsWeb) {
    return MeasurementRepositoryRest(ref.watch(restClientProvider));
  }
  final manager = ref.watch(grpcClientManagerProvider);
  return MeasurementRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

/// User Repository Provider (프로필/구독)
final userRepositoryProvider = Provider<UserRepository>((ref) {
  if (kIsWeb) {
    return UserRepositoryRest(ref.watch(restClientProvider));
  }
  final manager = ref.watch(grpcClientManagerProvider);
  return UserRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

// ── 화면용 데이터 Provider (gRPC/REST 자동 선택) ──

/// 최근 측정 기록 (HomeScreen). userId 없으면 빈 결과.
final measurementHistoryProvider = FutureProvider<MeasurementHistoryResult>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) {
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
  try {
    return await ref.read(measurementRepositoryProvider).getHistory(userId: userId, limit: 10);
  } catch (_) {
    // 백엔드 미연결 시 빈 결과 반환 (에러 대신)
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
});

/// 디바이스 목록 (DeviceListScreen)
final deviceListProvider = FutureProvider<List<DeviceItem>>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return [];
  try {
    return await ref.read(deviceRepositoryProvider).listDevices(userId);
  } catch (_) {
    return [];
  }
});

/// 사용자 프로필 (SettingsScreen)
final userProfileProvider = FutureProvider<UserProfileInfo?>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return null;
  try {
    return await ref.read(userRepositoryProvider).getProfile(userId);
  } catch (_) {
    return null;
  }
});

/// 구독 정보 (SettingsScreen)
final subscriptionInfoProvider = FutureProvider<SubscriptionInfoDto?>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return null;
  try {
    return await ref.read(userRepositoryProvider).getSubscription(userId);
  } catch (_) {
    return null;
  }
});

// ── 기능별 리포지토리 Provider ──

/// Community Repository Provider
final communityRepositoryProvider = Provider<CommunityRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return CommunityRepositoryRest(client, userId: userId);
});

/// Medical Repository Provider
final medicalRepositoryProvider = Provider<MedicalRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return MedicalRepositoryRest(client, userId: userId);
});

/// Market Repository Provider
final marketRepositoryProvider = Provider<MarketRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return MarketRepositoryRest(client, userId: userId);
});

/// AI Coach Repository Provider
final aiCoachRepositoryProvider = Provider<AiCoachRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return AiCoachRepositoryRest(client, userId: userId);
});

/// Family Repository Provider
final familyRepositoryProvider = Provider<FamilyRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return FamilyRepositoryRest(client, userId: userId);
});

/// DataHub Repository Provider
final dataHubRepositoryProvider = Provider<DataHubRepository>((ref) {
  final client = ref.watch(restClientProvider);
  final userId = ref.watch(authProvider).userId ?? '';
  return DataHubRepositoryRest(client, userId: userId);
});

// ── 화면용 비동기 데이터 Provider ──

/// 커뮤니티 게시글 목록
final communityPostsProvider = FutureProvider.family<List<CommunityPost>, PostCategory?>((ref, category) async {
  try {
    return await ref.read(communityRepositoryProvider).getPosts(category: category);
  } catch (_) {
    return [];
  }
});

/// 커뮤니티 건강 챌린지
final healthChallengesProvider = FutureProvider<List<HealthChallenge>>((ref) async {
  try {
    return await ref.read(communityRepositoryProvider).getChallenges();
  } catch (_) {
    return [];
  }
});

/// 진료 예약 목록
final reservationsProvider = FutureProvider<List<TelemedicineReservation>>((ref) async {
  try {
    return await ref.read(medicalRepositoryProvider).getReservations();
  } catch (_) {
    return [];
  }
});

/// 처방전 목록
final prescriptionsProvider = FutureProvider<List<Prescription>>((ref) async {
  try {
    return await ref.read(medicalRepositoryProvider).getPrescriptions();
  } catch (_) {
    return [];
  }
});

/// 건강 리포트 목록
final healthReportsProvider = FutureProvider<List<HealthReport>>((ref) async {
  try {
    return await ref.read(medicalRepositoryProvider).getHealthReports();
  } catch (_) {
    return [];
  }
});

/// 카트리지 상품 목록
final cartridgeProductsProvider = FutureProvider.family<List<CartridgeProduct>, String?>((ref, tier) async {
  try {
    return await ref.read(marketRepositoryProvider).getProducts(tier: tier);
  } catch (_) {
    return [];
  }
});

/// 구독 플랜 목록
final subscriptionPlansProvider = FutureProvider<List<SubscriptionPlan>>((ref) async {
  try {
    return await ref.read(marketRepositoryProvider).getSubscriptionPlans();
  } catch (_) {
    return [];
  }
});

/// 주문 내역
final ordersProvider = FutureProvider<List<Order>>((ref) async {
  try {
    return await ref.read(marketRepositoryProvider).getOrders();
  } catch (_) {
    return [];
  }
});

/// AI 코치 오늘의 인사이트
final todayInsightProvider = FutureProvider<HealthInsight>((ref) async {
  try {
    return await ref.read(aiCoachRepositoryProvider).getTodayInsight();
  } catch (_) {
    return HealthInsight(
      summary: '데이터를 불러올 수 없습니다.',
      detail: '',
      confidence: 0.0,
      generatedAt: DateTime.now(),
    );
  }
});

/// AI 코치 추천 목록
final aiRecommendationsProvider = FutureProvider.family<List<Recommendation>, String>((ref, category) async {
  try {
    return await ref.read(aiCoachRepositoryProvider).getRecommendations(category);
  } catch (_) {
    return [];
  }
});

/// 가족 그룹 목록
final familyGroupsProvider = FutureProvider<List<FamilyGroup>>((ref) async {
  try {
    return await ref.read(familyRepositoryProvider).getMyGroups();
  } catch (_) {
    return [];
  }
});

/// 바이오마커 요약 목록 (DataHub)
final biomarkerSummariesProvider = FutureProvider<List<BiomarkerSummary>>((ref) async {
  try {
    return await ref.read(dataHubRepositoryProvider).getAllBiomarkerSummaries();
  } catch (_) {
    return [];
  }
});

/// 시스템 통계 (AdminDashboard)
final systemStatsProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  try {
    return await ref.read(restClientProvider).getSystemStats();
  } catch (_) {
    return {};
  }
});

/// 감사 로그 (AdminDashboard)
final auditLogProvider = FutureProvider<Map<String, dynamic>>((ref) async {
  try {
    return await ref.read(restClientProvider).getAuditLog();
  } catch (_) {
    return {};
  }
});

/// 카트리지 종류 목록 (Encyclopedia)
final cartridgeTypesProvider = FutureProvider<List<Map<String, dynamic>>>((ref) async {
  try {
    final resp = await ref.read(restClientProvider).listCartridgeTypes();
    final items = resp['types'] as List? ?? resp['items'] as List? ?? [];
    return items.cast<Map<String, dynamic>>();
  } catch (_) {
    return [];
  }
});

/// 챌린지 리더보드 (ChallengeScreen — C8)
final challengeLeaderboardProvider =
    FutureProvider.family<List<Map<String, dynamic>>, String>((ref, challengeId) async {
  try {
    final resp = await ref.read(restClientProvider).getChallengeLeaderboard(challengeId);
    final entries = resp['entries'] as List? ?? resp['items'] as List? ?? [];
    return entries.cast<Map<String, dynamic>>();
  } catch (_) {
    return [];
  }
});

/// 알림 미읽은 개수
final unreadNotificationCountProvider = FutureProvider<int>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return 0;
  try {
    final client = ref.read(restClientProvider);
    final res = await client.getUnreadCount(userId);
    return res['count'] as int? ?? res['unread_count'] as int? ?? 0;
  } catch (_) {
    return 0;
  }
});
