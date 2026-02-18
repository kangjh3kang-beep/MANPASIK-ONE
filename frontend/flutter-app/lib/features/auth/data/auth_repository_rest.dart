import 'package:dio/dio.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// REST Gateway를 사용하는 AuthRepository 구현체
///
/// 웹 플랫폼에서 gRPC 대신 REST API를 통해 인증 처리.
class AuthRepositoryRest implements AuthRepository {
  AuthRepositoryRest(this._client);

  final ManPaSikRestClient _client;

  @override
  Future<AuthResult> login(String email, String password) async {
    try {
      final res = await _client.login(email, password);
      final accessToken = res['access_token'] as String? ?? '';
      final refreshToken = res['refresh_token'] as String? ?? '';
      final userId = res['user_id'] as String? ?? '';

      if (accessToken.isNotEmpty) {
        _client.setAuthToken(accessToken);
        return AuthResult.success(
          userId: userId.isNotEmpty ? userId : 'unknown',
          email: email,
          displayName: res['display_name'] as String? ?? email.split('@').first,
          accessToken: accessToken,
          refreshToken: refreshToken,
        );
      }
      return AuthResult.failure(res['error'] as String? ?? 'Login failed');
    } on DioException catch (e) {
      return AuthResult.failure(e.message ?? 'Network error');
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
      await _client.register(email, password, displayName);
      return login(email, password);
    } on DioException catch (e) {
      return AuthResult.failure(e.message ?? 'Register failed');
    } catch (e) {
      return AuthResult.failure(e.toString());
    }
  }

  @override
  Future<void> logout() async {
    _client.clearAuthToken();
  }

  @override
  Future<bool> refreshToken() async {
    return false;
  }

  @override
  Future<bool> isAuthenticated() async {
    return false;
  }
}
