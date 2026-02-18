import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// AuthInterceptor는 Dio HTTP 클라이언트에 JWT 인증 토큰을 자동으로 첨부하고,
/// 401 응답 시 토큰 갱신을 시도합니다.
class AuthInterceptor extends Interceptor {
  AuthInterceptor({Dio? refreshDio}) : _refreshDio = refreshDio;

  final Dio? _refreshDio;

  static const _accessTokenKey = 'access_token';
  static const _refreshTokenKey = 'refresh_token';

  /// SharedPreferences에 토큰을 저장합니다.
  static Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_accessTokenKey, accessToken);
    await prefs.setString(_refreshTokenKey, refreshToken);
  }

  /// 저장된 토큰을 삭제합니다 (로그아웃).
  static Future<void> clearTokens() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_accessTokenKey);
    await prefs.remove(_refreshTokenKey);
  }

  /// 저장된 액세스 토큰을 반환합니다.
  static Future<String?> getAccessToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_accessTokenKey);
  }

  @override
  void onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString(_accessTokenKey);
    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    super.onRequest(options, handler);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && _refreshDio != null) {
      try {
        final prefs = await SharedPreferences.getInstance();
        final refreshToken = prefs.getString(_refreshTokenKey);
        if (refreshToken == null || refreshToken.isEmpty) {
          return super.onError(err, handler);
        }

        // 토큰 갱신 시도
        final response = await _refreshDio!.post('/auth/refresh', data: {
          'refresh_token': refreshToken,
        });

        final data = response.data as Map<String, dynamic>? ?? {};
        final newAccessToken =
            (data['accessToken'] ?? data['access_token'] ?? '') as String;
        final newRefreshToken =
            (data['refreshToken'] ?? data['refresh_token'] ?? refreshToken)
                as String;

        if (newAccessToken.isNotEmpty) {
          await saveTokens(
            accessToken: newAccessToken,
            refreshToken: newRefreshToken,
          );

          // 원래 요청 재시도
          final opts = err.requestOptions;
          opts.headers['Authorization'] = 'Bearer $newAccessToken';
          final retryResponse = await _refreshDio!.fetch(opts);
          return handler.resolve(retryResponse);
        }
      } catch (_) {
        // 갱신 실패 시 원래 에러 전달
      }
    }
    super.onError(err, handler);
  }
}
