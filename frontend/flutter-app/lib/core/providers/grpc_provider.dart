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
  final authState = ref.watch(authProvider);
  final userId = authState.userId;
  
  // [DEMO MODE]
  if (authState.isDemo) {
    return MeasurementHistoryResult(
      items: List.generate(10, (index) => MeasurementHistoryItem(
        sessionId: 'mock-measure-$index',
        measuredAt: DateTime.now().subtract(Duration(days: index)),
        primaryValue: (85 + (index % 10)).toDouble(), // 85~94
        cartridgeType: 'Focus',
        unit: '점',
      )),
      totalCount: 56,
    );
  }

  if (userId == null || userId.isEmpty) {
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
  try {
    return await ref.read(measurementRepositoryProvider).getHistory(userId: userId, limit: 10);
  } catch (_) {
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
});

/// 디바이스 목록 (DeviceListScreen)
final deviceListProvider = FutureProvider<List<DeviceItem>>((ref) async {
  final authState = ref.watch(authProvider);
  final userId = authState.userId;

  // [DEMO MODE]
  if (authState.isDemo) {
    return [
      DeviceItem(
        deviceId: 'mock-device-1',
        name: '만파식 ONE (가상)',
        firmwareVersion: '1.2.0',
        status: 'online',
        lastSeen: DateTime.now(),
        batteryPercent: 98,
      ),
    ];
  }

  if (userId == null || userId.isEmpty) return [];
  try {
    return await ref.read(deviceRepositoryProvider).listDevices(userId);
  } catch (_) {
    return [];
  }
});

/// 연결된 디바이스 (모니터링용 - DataHub)
final connectedDevicesProvider = FutureProvider<List<ConnectedDevice>>((ref) async {
  final authState = ref.watch(authProvider);

  // [DEMO MODE] - 10 Simulated Devices
  if (authState.isDemo) {
    await Future.delayed(const Duration(milliseconds: 500));
    return [
      ConnectedDevice(
        id: 'gas-001',
        name: '거실 공기질 측정기',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 85,
        signalStrength: 92,
        currentValues: {'CO2': '450 ppm', 'VOC': '0.05 ppm', 'Radon': 'Safe'},
        latestReadings: [420, 430, 450, 440, 460, 450, 455, 450, 448, 452],
      ),
      ConnectedDevice(
        id: 'env-002',
        name: '안방 환경 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 90,
        signalStrength: 88,
        currentValues: {'Temp': '24.5°C', 'Humidity': '45%', 'Light': '300 lux'},
        latestReadings: [24.0, 24.1, 24.2, 24.5, 24.5, 24.4, 24.5, 24.6, 24.5, 24.5],
      ),
      ConnectedDevice(
        id: 'gas-002',
        name: '주방 가스 감지기',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 72,
        signalStrength: 75,
        currentValues: {'CO': '0 ppm', 'LNG': '0%', 'Smoke': 'None'},
        latestReadings: [0, 0, 0, 0, 0, 0, 1, 0, 0, 0],
      ),
      ConnectedDevice(
        id: 'bio-001',
        name: '바이오 카트리지 #1',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.disconnected,
        batteryLevel: 0,
        signalStrength: 0,
      ),
      ConnectedDevice(
        id: 'env-003',
        name: '아이방 온습도계',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 95,
        signalStrength: 98,
        currentValues: {'Temp': '23.0°C', 'Humidity': '50%'},
        latestReadings: [23, 23, 23, 23, 23],
      ),
      ConnectedDevice(
        id: 'gas-003',
        name: '베란다 환기 센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 60,
        signalStrength: 82,
        currentValues: {'Dust': '15 ug/m3'},
        latestReadings: [10, 12, 15, 14, 15],
      ),
      ConnectedDevice(
        id: 'bio-002',
        name: '웨어러블 밴드 Left',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 45,
        signalStrength: 90,
        currentValues: {'Pulse': '72 bpm', 'O2': '98%'},
        latestReadings: [70, 72, 71, 72, 75],
      ),
      ConnectedDevice(
        id: 'bio-003',
        name: '웨어러블 밴드 Right',
        type: DeviceType.bioCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 42,
        signalStrength: 88,
        currentValues: {'Pulse': '73 bpm', 'Stress': 'Low'},
        latestReadings: [72, 73, 73, 74, 73],
      ),
      ConnectedDevice(
        id: 'env-004',
        name: '서재 조명 센서',
        type: DeviceType.envCartridge,
        status: DeviceConnectionStatus.disconnected,
        batteryLevel: 10,
        signalStrength: 20,
      ),
      ConnectedDevice(
        id: 'gas-004',
        name: '차고 배기 센서',
        type: DeviceType.gasCartridge,
        status: DeviceConnectionStatus.connected,
        batteryLevel: 88,
        signalStrength: 65,
        currentValues: {'CO': '2 ppm'},
        latestReadings: [1, 1, 2, 2, 1],
      ),
    ];
  }

  try {
    // Repository Mock 구현 사용
    return await ref.read(deviceRepositoryProvider).getConnectedDevices();
  } catch (_) {
    return [];
  }
});

/// 사용자 프로필 (SettingsScreen)
final userProfileProvider = FutureProvider<UserProfileInfo?>((ref) async {
  final authState = ref.watch(authProvider);
  final userId = authState.userId;

  // [DEMO MODE]
  if (authState.isDemo) {
    return const UserProfileInfo(
      userId: 'demo-user-id',
      email: 'demo@manpasik.com',
      displayName: '테스트 계정',
      avatarUrl: null,
      subscriptionTier: 1, // Pro
    );
  }

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
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (authState.isDemo) {
    return List.generate(5, (index) => CommunityPost(
      id: 'mock-post-$index',
      authorId: 'user-$index',
      authorName: '사용자 ${index + 1}',
      title: '만파식 체험 후기 $index',
      content: '오늘 만파식으로 측정한 결과가 아주 좋네요! 다들 건강 챙기세요. #건강 #만파식',
      likeCount: 10 + index * 5,
      commentCount: index,
      isLikedByMe: index % 2 == 0,
      isBookmarkedByMe: false,
      createdAt: DateTime.now().subtract(Duration(hours: index)),
      category: category ?? PostCategory.reviews,
    ));
  }

  try {
    return await ref.read(communityRepositoryProvider).getPosts(category: category);
  } catch (_) {
    return [];
  }
});

/// 커뮤니티 건강 챌린지
final healthChallengesProvider = FutureProvider<List<HealthChallenge>>((ref) async {
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (authState.isDemo) {
    return [
      HealthChallenge(
        id: 'mock-challenge-1',
        title: '30일 꾸준한 측정',
        description: '30일 동안 매일 아침 스트레스를 측정하고 기록해보세요.',
        startDate: DateTime.now().subtract(const Duration(days: 5)),
        endDate: DateTime.now().add(const Duration(days: 25)),
        participantCount: 1240,
        isJoined: false,
      ),
      HealthChallenge(
        id: 'mock-challenge-2',
        title: '수면 패턴 개선하기',
        description: '잠자기 1시간 전 스마트폰 멀리하기 챌린지!',
        startDate: DateTime.now(),
        endDate: DateTime.now().add(const Duration(days: 7)),
        participantCount: 856,
        isJoined: true,
      ),
    ];
  }

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
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (authState.isDemo) {
    return [
      const CartridgeProduct(
        id: 'mock-prod-1',
        nameKo: '집중력 강화 카트리지',
        nameEn: 'Focus Booster',
        typeCode: 'FOCUS',
        tier: 'Standard',
        price: 9900,
        unit: 'ea',
        referenceRange: '0-100',
        requiredChannels: 1,
        measurementSecs: 60,
        isAvailable: true,
      ),
      const CartridgeProduct(
        id: 'mock-prod-2',
        nameKo: '수면 유도 카트리지',
        nameEn: 'Sleep Aid',
        typeCode: 'SLEEP',
        tier: 'Premium',
        price: 12900,
        unit: 'ea',
        referenceRange: '0-100',
        requiredChannels: 1,
        measurementSecs: 120,
        isAvailable: true,
      ),
      const CartridgeProduct(
        id: 'mock-prod-3',
        nameKo: '스트레스 해소 팩',
        nameEn: 'Stress Relief',
        typeCode: 'RELAX',
        tier: 'Standard',
        price: 15000,
        unit: 'pack',
        referenceRange: '0-100',
        requiredChannels: 1,
        measurementSecs: 60,
        isAvailable: true,
      ),
    ];
  }

  try {
    return await ref.read(marketRepositoryProvider).getProducts(tier: tier);
  } catch (_) {
    return [];
  }
});

/// 구독 플랜 목록
final subscriptionPlansProvider = FutureProvider<List<SubscriptionPlan>>((ref) async {
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (authState.isDemo) {
    return [
      const SubscriptionPlan(
        id: 'plan-basic',
        name: '베이직 플랜',
        monthlyPrice: 0,
        discountPercent: 0,
        includedCartridgeTypes: ['BASIC'],
        cartridgesPerMonth: 5,
      ),
      const SubscriptionPlan(
        id: 'plan-pro',
        name: '프로 플랜',
        monthlyPrice: 9900,
        discountPercent: 10,
        includedCartridgeTypes: ['BASIC', 'FOCUS', 'SLEEP'],
        cartridgesPerMonth: 999,
      ),
    ];
  }

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

/// 추천 의사 목록
final recommendedDoctorsProvider = FutureProvider<List<DoctorInfo>>((ref) async {
  try {
    return await ref.read(medicalRepositoryProvider).getRecommendedDoctors();
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
  final authState = ref.watch(authProvider);

  // [DEMO MODE]
  if (authState.isDemo) {
    return [
      const BiomarkerSummary(biomarkerType: 'stress', displayName: '스트레스', latestValue: 45, unit: '점', referenceMin: 0, referenceMax: 60, totalMeasurements: 10, trend: 'stable'),
      const BiomarkerSummary(biomarkerType: 'energy', displayName: '에너지', latestValue: 85, unit: '점', referenceMin: 60, referenceMax: 100, totalMeasurements: 10, trend: 'rising'),
      const BiomarkerSummary(biomarkerType: 'hydration', displayName: '수분 균형', latestValue: 72, unit: '%', referenceMin: 70, referenceMax: 100, totalMeasurements: 10, trend: 'stable'),
      const BiomarkerSummary(biomarkerType: 'sleep', displayName: '수면 질', latestValue: 65, unit: '점', referenceMin: 60, referenceMax: 90, totalMeasurements: 10, trend: 'falling'),
      const BiomarkerSummary(biomarkerType: 'hrv', displayName: '심박 변이도', latestValue: 42, unit: 'ms', referenceMin: 30, referenceMax: 100, totalMeasurements: 10, trend: 'stable'),
      const BiomarkerSummary(biomarkerType: 'focus', displayName: '집중력', latestValue: 88, unit: '점', referenceMin: 70, referenceMax: 100, totalMeasurements: 10, trend: 'rising'),
    ];
  }

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
