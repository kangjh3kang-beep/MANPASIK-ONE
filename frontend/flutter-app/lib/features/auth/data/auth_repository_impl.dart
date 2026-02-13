import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/services/grpc_client.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:grpc/grpc.dart';

/// gRPC AuthService를 사용하는 AuthRepository 구현체
class AuthRepositoryImpl implements AuthRepository {
  AuthRepositoryImpl(this._grpcManager);

  final GrpcClientManager _grpcManager;

  AuthServiceClient? _client;

  AuthServiceClient get _authClient {
    _client ??= AuthServiceClient(_grpcManager.authChannel);
    return _client!;
  }

  @override
  Future<AuthResult> login(String email, String password) async {
    try {
      final res = await _authClient.login(
        LoginRequest()
          ..email = email
          ..password = password,
      );
      return AuthResult.success(
        userId: res.userId.isNotEmpty ? res.userId : 'unknown',
        email: email,
        displayName: email.split('@').first,
        accessToken: res.accessToken,
        refreshToken: res.refreshToken,
      );
    } on GrpcError catch (e) {
      return AuthResult.failure(e.message ?? 'Login failed');
    } catch (e) {
      return AuthResult.failure(e.toString());
    }
  }

  @override
  Future<AuthResult> register(
    String email,
    String password,
    String displayName,
  ) async {
    try {
      await _authClient.register(
        RegisterRequest()
          ..email = email
          ..password = password
          ..displayName = displayName,
      );
      // Register 응답에 토큰이 없을 수 있으므로 로그인 한 번 더 호출
      return login(email, password);
    } on GrpcError catch (e) {
      return AuthResult.failure(e.message ?? 'Register failed');
    } catch (e) {
      return AuthResult.failure(e.toString());
    }
  }

  @override
  Future<void> logout() async {
    // 로컬 상태만 초기화. 서버 Logout RPC는 호출하지 않음 (토큰 없을 수 있음)
  }

  @override
  Future<bool> refreshToken() async {
    // 호출 측에서 refresh token을 넘겨줘야 함. 현재는 상태 기반이므로 AuthNotifier에서 처리
    return false;
  }

  @override
  Future<bool> isAuthenticated() async {
    // 토큰 유효성은 ValidateToken RPC로 검사 가능. 현재는 로컬 상태만 확인
    return false;
  }
}
