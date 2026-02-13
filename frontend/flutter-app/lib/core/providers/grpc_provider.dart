import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/features/auth/data/auth_repository_impl.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/features/devices/data/device_repository_impl.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/features/measurement/data/measurement_repository_impl.dart';
import 'package:manpasik/features/user/data/user_repository_impl.dart';

/// gRPC 채널 관리자 Provider (싱글톤)
final grpcClientManagerProvider = Provider<GrpcClientManager>((ref) {
  return GrpcClientManager();
});

/// Auth Repository Provider (gRPC AuthService 연동)
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return AuthRepositoryImpl(manager);
});

/// Device Repository Provider (JWT 인터셉터에 현재 액세스 토큰 전달)
final deviceRepositoryProvider = Provider<DeviceRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return DeviceRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

/// Measurement Repository Provider
final measurementRepositoryProvider = Provider<MeasurementRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return MeasurementRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

/// User Repository Provider (프로필/구독)
final userRepositoryProvider = Provider<UserRepository>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  return UserRepositoryImpl(
    manager,
    accessTokenProvider: () => ref.read(authProvider).accessToken,
  );
});

// ── 화면용 데이터 Provider (실제 gRPC 로드) ──

/// 최근 측정 기록 (HomeScreen). userId 없으면 빈 결과.
final measurementHistoryProvider = FutureProvider<MeasurementHistoryResult>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) {
    return const MeasurementHistoryResult(items: [], totalCount: 0);
  }
  return ref.read(measurementRepositoryProvider).getHistory(userId: userId, limit: 10);
});

/// 디바이스 목록 (DeviceListScreen)
final deviceListProvider = FutureProvider<List<DeviceItem>>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return [];
  return ref.read(deviceRepositoryProvider).listDevices(userId);
});

/// 사용자 프로필 (SettingsScreen)
final userProfileProvider = FutureProvider<UserProfileInfo?>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return null;
  return ref.read(userRepositoryProvider).getProfile(userId);
});

/// 구독 정보 (SettingsScreen)
final subscriptionInfoProvider = FutureProvider<SubscriptionInfoDto?>((ref) async {
  final userId = ref.watch(authProvider).userId;
  if (userId == null || userId.isEmpty) return null;
  return ref.read(userRepositoryProvider).getSubscription(userId);
});
