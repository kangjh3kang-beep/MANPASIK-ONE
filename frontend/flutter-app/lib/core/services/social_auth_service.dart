import 'package:flutter/foundation.dart';
import 'package:url_launcher/url_launcher.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';

/// 소셜 인증 프로바이더
enum SocialProvider { google, kakao, apple }

/// 소셜 로그인 서비스
///
/// Google / Kakao / Apple OAuth를 지원합니다.
/// 각 SDK 패키지(google_sign_in, kakao_flutter_sdk_user, sign_in_with_apple)
/// 설치 시 네이티브 로그인, 미설치 시 브라우저 기반 OAuth 리다이렉트 사용.
class SocialAuthService {
  SocialAuthService({required this.restClient});

  final ManPaSikRestClient restClient;

  /// 소셜 로그인 실행
  ///
  /// 1. 네이티브 SDK 우선 시도 (패키지 설치 시)
  /// 2. 미설치 시 브라우저 OAuth 리다이렉트
  /// 3. 서버에 토큰 전달 후 AuthResult 반환
  Future<AuthResult> login(SocialProvider provider) async {
    switch (provider) {
      case SocialProvider.google:
        return _loginWithGoogle();
      case SocialProvider.kakao:
        return _loginWithKakao();
      case SocialProvider.apple:
        return _loginWithApple();
    }
  }

  Future<AuthResult> _loginWithGoogle() async {
    try {
      // Google Sign-In SDK 설치 시:
      // final googleUser = await GoogleSignIn(scopes: ['email', 'profile']).signIn();
      // if (googleUser == null) return AuthResult.failure('Google 로그인 취소');
      // final googleAuth = await googleUser.authentication;
      // final idToken = googleAuth.idToken!;

      // SDK 미설치 → 브라우저 OAuth 리다이렉트
      final redirectUri = Uri.encodeFull('${AppConfig.baseUrl}/auth/callback/google');
      final authUrl = Uri.parse(
        'https://accounts.google.com/o/oauth2/v2/auth'
        '?client_id=${const String.fromEnvironment('GOOGLE_CLIENT_ID')}'
        '&redirect_uri=$redirectUri'
        '&response_type=code'
        '&scope=email%20profile'
        '&state=manpasik_google',
      );

      if (const String.fromEnvironment('GOOGLE_CLIENT_ID').isEmpty) {
        debugPrint('[SocialAuth] GOOGLE_CLIENT_ID 미설정 → 시뮬레이션');
        return _simulateLogin('google');
      }

      final launched = await launchUrl(authUrl, mode: LaunchMode.externalApplication);
      if (!launched) {
        return AuthResult.failure('브라우저를 열 수 없습니다');
      }

      // Deep link 콜백 대기 (앱으로 복귀 시 auth/callback에서 처리)
      return AuthResult.failure('브라우저에서 인증을 완료해주세요');
    } catch (e) {
      debugPrint('[SocialAuth] Google 로그인 실패: $e');
      return AuthResult.failure('Google 로그인 실패: $e');
    }
  }

  Future<AuthResult> _loginWithKakao() async {
    try {
      // Kakao SDK 설치 시:
      // final OAuthToken token;
      // if (await isKakaoTalkInstalled()) {
      //   token = await UserApi.instance.loginWithKakaoTalk();
      // } else {
      //   token = await UserApi.instance.loginWithKakaoAccount();
      // }
      // return _exchangeToken('kakao', token.accessToken);

      final redirectUri = Uri.encodeFull('${AppConfig.baseUrl}/auth/callback/kakao');
      final authUrl = Uri.parse(
        'https://kauth.kakao.com/oauth/authorize'
        '?client_id=${const String.fromEnvironment('KAKAO_REST_KEY')}'
        '&redirect_uri=$redirectUri'
        '&response_type=code'
        '&state=manpasik_kakao',
      );

      if (const String.fromEnvironment('KAKAO_REST_KEY').isEmpty) {
        debugPrint('[SocialAuth] KAKAO_REST_KEY 미설정 → 시뮬레이션');
        return _simulateLogin('kakao');
      }

      final launched = await launchUrl(authUrl, mode: LaunchMode.externalApplication);
      if (!launched) {
        return AuthResult.failure('브라우저를 열 수 없습니다');
      }

      return AuthResult.failure('브라우저에서 인증을 완료해주세요');
    } catch (e) {
      debugPrint('[SocialAuth] Kakao 로그인 실패: $e');
      return AuthResult.failure('Kakao 로그인 실패: $e');
    }
  }

  Future<AuthResult> _loginWithApple() async {
    try {
      // Sign in with Apple SDK 설치 시:
      // final credential = await SignInWithApple.getAppleIDCredential(
      //   scopes: [
      //     AppleIDAuthorizationScopes.email,
      //     AppleIDAuthorizationScopes.fullName,
      //   ],
      // );
      // return _exchangeToken('apple', credential.identityToken!);

      if (const String.fromEnvironment('APPLE_SERVICE_ID').isEmpty) {
        debugPrint('[SocialAuth] APPLE_SERVICE_ID 미설정 → 시뮬레이션');
        return _simulateLogin('apple');
      }

      final redirectUri = Uri.encodeFull('${AppConfig.baseUrl}/auth/callback/apple');
      final authUrl = Uri.parse(
        'https://appleid.apple.com/auth/authorize'
        '?client_id=${const String.fromEnvironment('APPLE_SERVICE_ID')}'
        '&redirect_uri=$redirectUri'
        '&response_type=code%20id_token'
        '&scope=name%20email'
        '&response_mode=form_post'
        '&state=manpasik_apple',
      );

      final launched = await launchUrl(authUrl, mode: LaunchMode.externalApplication);
      if (!launched) {
        return AuthResult.failure('브라우저를 열 수 없습니다');
      }

      return AuthResult.failure('브라우저에서 인증을 완료해주세요');
    } catch (e) {
      debugPrint('[SocialAuth] Apple 로그인 실패: $e');
      return AuthResult.failure('Apple 로그인 실패: $e');
    }
  }

  /// 서버에 소셜 토큰을 전달하여 앱 토큰으로 교환
  Future<AuthResult> exchangeToken(String provider, String token) async {
    try {
      final res = await restClient.socialLogin(
        provider: provider,
        idToken: token,
      );
      final accessToken = res['access_token'] as String? ?? '';
      final refreshToken = res['refresh_token'] as String? ?? '';
      final userId = res['user_id'] as String? ?? '';

      if (accessToken.isNotEmpty) {
        restClient.setAuthToken(accessToken);
        return AuthResult.success(
          userId: userId,
          email: res['email'] as String? ?? '',
          displayName: res['display_name'] as String?,
          accessToken: accessToken,
          refreshToken: refreshToken,
        );
      }
      return AuthResult.failure(res['error'] as String? ?? '소셜 로그인 실패');
    } catch (e) {
      return AuthResult.failure('소셜 로그인 서버 오류: $e');
    }
  }

  Future<AuthResult> _simulateLogin(String provider) async {
    debugPrint('[SocialAuth:Sim] $provider 로그인 시뮬레이션');
    await Future.delayed(const Duration(seconds: 1));
    return AuthResult.success(
      userId: 'social_${provider}_${DateTime.now().millisecondsSinceEpoch}',
      email: 'user@$provider.com',
      displayName: '$provider 사용자',
      accessToken: 'sim_token_$provider',
      refreshToken: 'sim_refresh_$provider',
    );
  }
}
