import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/features/user/domain/user_repository.dart';

/// 테스트용 Fake AuthRepository (gRPC 없이 로그인/회원가입 성공 시뮬레이션)
class FakeAuthRepository implements AuthRepository {
  @override
  Future<AuthResult> login(String email, String password) async {
    if (email.isEmpty || password.length < 8) {
      return AuthResult.failure('Invalid credentials');
    }
    return AuthResult.success(
      userId: 'test-user-id',
      email: email,
      displayName: email.split('@').first,
      accessToken: 'fake-access-token',
      refreshToken: 'fake-refresh-token',
    );
  }

  @override
  Future<AuthResult> register(
    String email,
    String password,
    String displayName,
  ) async {
    if (email.isEmpty || password.length < 8 || displayName.isEmpty) {
      return AuthResult.failure('Invalid input');
    }
    return AuthResult.success(
      userId: 'test-user-id',
      email: email,
      displayName: displayName,
      accessToken: 'fake-access-token',
      refreshToken: 'fake-refresh-token',
    );
  }

  @override
  Future<void> logout() async {}

  @override
  Future<bool> refreshToken() async => false;

  @override
  Future<bool> isAuthenticated() async => false;
}

/// 테스트용 Fake DeviceRepository
class FakeDeviceRepository implements DeviceRepository {
  @override
  Future<List<DeviceItem>> listDevices(String userId) async {
    if (userId.isEmpty) return [];
    return [
      const DeviceItem(
        deviceId: 'device-1',
        name: '테스트 리더기',
        firmwareVersion: '1.0.0',
        status: 'online',
        batteryPercent: 85,
      ),
    ];
  }
}

/// 테스트용 Fake MeasurementRepository
class FakeMeasurementRepository implements MeasurementRepository {
  @override
  Future<StartSessionResult> startSession({
    required String deviceId,
    required String cartridgeId,
    required String userId,
  }) async {
    return const StartSessionResult(sessionId: 'session-1');
  }

  @override
  Future<EndSessionResult?> endSession(String sessionId) async {
    return const EndSessionResult(
      sessionId: 'session-1',
      totalMeasurements: 5,
    );
  }

  @override
  Future<MeasurementHistoryResult> getHistory({
    required String userId,
    int limit = 20,
    int offset = 0,
  }) async {
    if (userId.isEmpty) {
      return const MeasurementHistoryResult(items: [], totalCount: 0);
    }
    return const MeasurementHistoryResult(
      items: [
        MeasurementHistoryItem(
          sessionId: 's1',
          cartridgeType: 'glucose',
          primaryValue: 98.4,
          unit: 'mg/dL',
          measuredAt: null,
        ),
      ],
      totalCount: 1,
    );
  }
}

/// 테스트용 Fake UserRepository
class FakeUserRepository implements UserRepository {
  @override
  Future<UserProfileInfo?> getProfile(String userId) async {
    if (userId.isEmpty) return null;
    return UserProfileInfo(
      userId: userId,
      email: 'test@manpasik.com',
      displayName: '테스트 사용자',
      subscriptionTier: 0,
    );
  }

  @override
  Future<SubscriptionInfoDto?> getSubscription(String userId) async {
    if (userId.isEmpty) return null;
    return const SubscriptionInfoDto(
      userId: 'test-user-id',
      tier: 0,
      maxDevices: 1,
      maxFamilyMembers: 1,
    );
  }
}
