import 'package:flutter/foundation.dart';
import 'package:kakao_flutter_sdk_user/kakao_flutter_sdk_user.dart';
import 'package:dio/dio.dart';

import 'package:manpasik/core/config/app_config.dart';

/// 카카오 로그인 결과
class KakaoLoginResult {
  final bool success;
  final String? accessToken;
  final String? refreshToken;
  final String? userId;
  final String? email;
  final String? displayName;
  final String? errorMessage;

  const KakaoLoginResult._({
    required this.success,
    this.accessToken,
    this.refreshToken,
    this.userId,
    this.email,
    this.displayName,
    this.errorMessage,
  });

  factory KakaoLoginResult.success({
    required String accessToken,
    required String refreshToken,
    required String userId,
    String? email,
    String? displayName,
  }) {
    return KakaoLoginResult._(
      success: true,
      accessToken: accessToken,
      refreshToken: refreshToken,
      userId: userId,
      email: email,
      displayName: displayName,
    );
  }

  factory KakaoLoginResult.failure(String message) {
    return KakaoLoginResult._(success: false, errorMessage: message);
  }
}

/// 카카오 소셜 로그인 서비스
///
/// 카카오톡 앱 설치 여부에 따라 앱/웹 로그인을 분기하고,
/// 로그인 성공 시 카카오 Access Token을 백엔드 API Gateway에
/// 전송하여 만파식 JWT(Access/Refresh Token)를 발급받습니다.
class KakaoLoginService {
  KakaoLoginService({Dio? dio})
      : _dio = dio ??
            Dio(BaseOptions(
              baseUrl: AppConfig.baseUrl,
              connectTimeout: const Duration(seconds: 10),
              receiveTimeout: const Duration(seconds: 10),
              headers: {'Content-Type': 'application/json'},
            ));

  final Dio _dio;

  /// 카카오 로그인 실행
  ///
  /// 전체 흐름:
  /// 1. 카카오톡 앱 설치 여부 확인
  /// 2. 앱 설치 → [loginWithKakaoTalk], 미설치 → [loginWithKakaoAccount]
  /// 3. 카카오 Access Token 획득
  /// 4. 백엔드 `/api/v1/auth/kakao` 엔드포인트로 토큰 전송
  /// 5. 만파식 JWT(Access/Refresh Token) 수신 → [KakaoLoginResult] 반환
  Future<KakaoLoginResult> login() async {
    try {
      // 1단계: 카카오 OAuth 토큰 획득
      final OAuthToken kakaoToken = await _acquireKakaoToken();
      final String kakaoAccessToken = kakaoToken.accessToken;

      debugPrint('[KakaoLogin] 카카오 Access Token 획득 완료');

      // 2단계: 카카오 사용자 정보 조회 (닉네임/이메일)
      final kakaoUser = await UserApi.instance.me();
      final kakaoEmail = kakaoUser.kakaoAccount?.email;
      final kakaoNickname = kakaoUser.kakaoAccount?.profile?.nickname;

      debugPrint('[KakaoLogin] 카카오 사용자: ${kakaoUser.id}, $kakaoNickname');

      // 3단계: 백엔드에 카카오 Access Token 전송 → JWT 발급
      final jwtResult = await _exchangeTokenWithBackend(
        kakaoAccessToken: kakaoAccessToken,
        kakaoEmail: kakaoEmail,
        kakaoNickname: kakaoNickname,
      );

      return jwtResult;
    } on KakaoAuthException catch (e) {
      debugPrint('[KakaoLogin] 카카오 인증 오류: ${e.error}');
      return KakaoLoginResult.failure('카카오 인증 실패: ${e.error.name}');
    } on KakaoClientException catch (e) {
      debugPrint('[KakaoLogin] 카카오 클라이언트 오류: ${e.msg}');
      return KakaoLoginResult.failure('카카오 로그인 취소 또는 오류: ${e.msg}');
    } on DioException catch (e) {
      debugPrint('[KakaoLogin] 서버 통신 오류: ${e.message}');
      return KakaoLoginResult.failure(
        '서버 통신 오류: ${e.response?.statusCode ?? "네트워크 연결 실패"}',
      );
    } catch (e) {
      debugPrint('[KakaoLogin] 예상치 못한 오류: $e');
      return KakaoLoginResult.failure('카카오 로그인 실패: $e');
    }
  }

  /// 카카오 OAuth 토큰 획득 (앱/웹 자동 분기)
  ///
  /// - 카카오톡 앱 설치됨 → [UserApi.instance.loginWithKakaoTalk]
  /// - 카카오톡 미설치 → [UserApi.instance.loginWithKakaoAccount] (웹뷰)
  Future<OAuthToken> _acquireKakaoToken() async {
    // 카카오톡 앱 설치 여부 확인
    final bool isKakaoTalkAvailable = await isKakaoTalkInstalled();
    debugPrint('[KakaoLogin] 카카오톡 앱 설치 여부: $isKakaoTalkAvailable');

    if (isKakaoTalkAvailable) {
      // 카카오톡 앱으로 로그인
      try {
        final token = await UserApi.instance.loginWithKakaoTalk();
        debugPrint('[KakaoLogin] 카카오톡 앱 로그인 성공');
        return token;
      } catch (e) {
        debugPrint('[KakaoLogin] 카카오톡 앱 로그인 실패, 웹뷰로 폴백: $e');
        // 앱 로그인 실패 시 웹뷰로 폴백
        final token = await UserApi.instance.loginWithKakaoAccount();
        debugPrint('[KakaoLogin] 카카오 웹뷰 로그인 성공 (폴백)');
        return token;
      }
    } else {
      // 카카오톡 미설치 → 웹뷰로 로그인
      final token = await UserApi.instance.loginWithKakaoAccount();
      debugPrint('[KakaoLogin] 카카오 웹뷰 로그인 성공');
      return token;
    }
  }

  /// 백엔드 API Gateway에 카카오 Access Token 전송 → 만파식 JWT 발급
  ///
  /// POST /api/v1/auth/kakao
  /// Body: { "kakao_access_token": "...", "email": "...", "nickname": "..." }
  /// Response: { "access_token": "...", "refresh_token": "...", "user_id": "...", ... }
  Future<KakaoLoginResult> _exchangeTokenWithBackend({
    required String kakaoAccessToken,
    String? kakaoEmail,
    String? kakaoNickname,
  }) async {
    final response = await _dio.post(
      '/api/v1/auth/kakao',
      data: {
        'kakao_access_token': kakaoAccessToken,
        'email': kakaoEmail ?? '',
        'nickname': kakaoNickname ?? '',
      },
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final data = response.data as Map<String, dynamic>;

      final accessToken = data['access_token'] as String?;
      final refreshToken = data['refresh_token'] as String?;
      final userId = data['user_id'] as String?;
      final email = data['email'] as String?;
      final displayName = data['display_name'] as String?;

      if (accessToken == null || accessToken.isEmpty) {
        return KakaoLoginResult.failure('서버에서 유효한 토큰을 반환하지 않았습니다.');
      }

      debugPrint('[KakaoLogin] 만파식 JWT 발급 완료 — userId: $userId');

      return KakaoLoginResult.success(
        accessToken: accessToken,
        refreshToken: refreshToken ?? '',
        userId: userId ?? '',
        email: email,
        displayName: displayName,
      );
    } else {
      final errorMsg = response.data is Map
          ? (response.data['error'] as String? ?? '서버 오류')
          : '서버 응답 오류 (${response.statusCode})';
      return KakaoLoginResult.failure(errorMsg);
    }
  }

  /// 카카오 로그아웃
  ///
  /// 카카오 SDK 세션을 종료합니다.
  /// 만파식 JWT 무효화는 별도의 백엔드 logout API를 호출해야 합니다.
  Future<void> logout() async {
    try {
      await UserApi.instance.logout();
      debugPrint('[KakaoLogin] 카카오 로그아웃 완료');
    } catch (e) {
      debugPrint('[KakaoLogin] 카카오 로그아웃 실패: $e');
    }
  }

  /// 카카오 연결 해제 (탈퇴)
  ///
  /// 카카오와의 앱 연결을 완전히 해제합니다.
  Future<void> unlink() async {
    try {
      await UserApi.instance.unlink();
      debugPrint('[KakaoLogin] 카카오 연결 해제 완료');
    } catch (e) {
      debugPrint('[KakaoLogin] 카카오 연결 해제 실패: $e');
    }
  }

  /// 현재 카카오 토큰 유효성 확인
  Future<bool> hasValidKakaoToken() async {
    try {
      final tokenInfo = await UserApi.instance.accessTokenInfo();
      return tokenInfo.id != null;
    } catch (_) {
      return false;
    }
  }
}
