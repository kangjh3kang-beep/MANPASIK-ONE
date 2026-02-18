/// 인증 Repository 인터페이스
///
/// S5에서 gRPC auth-service와 연동 시 구현체 작성.
abstract class AuthRepository {
  Future<AuthResult> login(String email, String password);
  Future<AuthResult> register(String email, String password, String displayName);
  Future<void> logout();
  Future<bool> refreshToken();
  Future<bool> isAuthenticated();
}

/// 인증 결과
class AuthResult {
  final bool success;
  final String? userId;
  final String? email;
  final String? displayName;
  final String? accessToken;
  final String? refreshToken;
  final String? role;
  final String? errorMessage;

  const AuthResult({
    required this.success,
    this.userId,
    this.email,
    this.displayName,
    this.accessToken,
    this.refreshToken,
    this.role,
    this.errorMessage,
  });

  factory AuthResult.success({
    required String userId,
    required String email,
    String? displayName,
    required String accessToken,
    required String refreshToken,
    String? role,
  }) {
    return AuthResult(
      success: true,
      userId: userId,
      email: email,
      displayName: displayName,
      accessToken: accessToken,
      refreshToken: refreshToken,
      role: role,
    );
  }

  factory AuthResult.failure(String message) {
    return AuthResult(success: false, errorMessage: message);
  }
}
