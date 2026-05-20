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
      final accessToken =
          _readString(res, 'access_token', alternate: 'accessToken');
      final refreshToken =
          _readString(res, 'refresh_token', alternate: 'refreshToken');
      final userId = _readString(res, 'user_id', alternate: 'userId');

      if (accessToken.isNotEmpty) {
        _client.setAuthToken(accessToken);
        return AuthResult.success(
          userId: userId.isNotEmpty ? userId : 'unknown',
          email: email,
          displayName: _readString(
            res,
            'display_name',
            alternate: 'displayName',
            fallback: email.split('@').first,
          ),
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

  @override
  Future<AuthResult> socialLogin(String provider, String token) async {
    try {
      final res = await _client.socialLogin(
        provider: provider,
        idToken: token,
      );
      final accessToken =
          _readString(res, 'access_token', alternate: 'accessToken');
      final refreshToken =
          _readString(res, 'refresh_token', alternate: 'refreshToken');
      final userId = _readString(res, 'user_id', alternate: 'userId');

      if (accessToken.isNotEmpty) {
        _client.setAuthToken(accessToken);
        return AuthResult.success(
          userId: userId.isNotEmpty ? userId : 'unknown',
          email: _readString(res, 'email'),
          displayName: _readNullableString(
            res,
            'display_name',
            alternate: 'displayName',
          ),
          accessToken: accessToken,
          refreshToken: refreshToken,
        );
      }
      return AuthResult.failure(res['error'] as String? ?? '소셜 로그인 실패');
    } on DioException catch (e) {
      return AuthResult.failure(e.message ?? '네트워크 오류');
    } catch (e) {
      return AuthResult.failure(e.toString());
    }
  }

  String _readString(
    Map<String, dynamic> source,
    String primary, {
    String? alternate,
    String fallback = '',
  }) {
    return _readNullableString(source, primary, alternate: alternate) ??
        fallback;
  }

  String? _readNullableString(
    Map<String, dynamic> source,
    String primary, {
    String? alternate,
  }) {
    final primaryValue = source[primary];
    if (primaryValue is String) {
      return primaryValue;
    }

    if (alternate != null) {
      final alternateValue = source[alternate];
      if (alternateValue is String) {
        return alternateValue;
      }
    }

    return null;
  }
}
