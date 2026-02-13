import 'package:manpasik/features/user/domain/user_repository.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:grpc/grpc.dart';

/// gRPC UserService를 사용하는 UserRepository 구현체
class UserRepositoryImpl implements UserRepository {
  UserRepositoryImpl(
    this._grpcManager, {
    required String? Function() accessTokenProvider,
  }) : _authInterceptor = AuthInterceptor(accessTokenProvider);

  final GrpcClientManager _grpcManager;
  final AuthInterceptor _authInterceptor;

  UserServiceClient? _client;

  UserServiceClient get _userClient {
    _client ??= UserServiceClient(
      _grpcManager.userChannel,
      interceptors: [_authInterceptor],
    );
    return _client!;
  }

  @override
  Future<UserProfileInfo?> getProfile(String userId) async {
    try {
      final res = await _userClient.getProfile(
        GetProfileRequest()..userId = userId,
      );
      return UserProfileInfo(
        userId: res.userId,
        email: res.email,
        displayName: res.displayName,
        avatarUrl: res.avatarUrl.isNotEmpty ? res.avatarUrl : null,
        language: res.language.isNotEmpty ? res.language : null,
        timezone: res.timezone.isNotEmpty ? res.timezone : null,
        subscriptionTier: res.subscriptionTier,
      );
    } on GrpcError {
      rethrow;
    }
  }

  @override
  Future<SubscriptionInfoDto?> getSubscription(String userId) async {
    try {
      final res = await _userClient.getSubscription(
        GetSubscriptionRequest()..userId = userId,
      );
      return SubscriptionInfoDto(
        userId: res.userId,
        tier: res.tier,
        maxDevices: res.maxDevices,
        maxFamilyMembers: res.maxFamilyMembers,
        aiCoachingEnabled: res.aiCoachingEnabled,
        telemedicineEnabled: res.telemedicineEnabled,
      );
    } on GrpcError {
      rethrow;
    }
  }
}
