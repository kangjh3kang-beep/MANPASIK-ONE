import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/network/tenant_interceptor.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// 인증 상태 모델
class AuthState {
  final bool isAuthenticated;
  final String? userId;
  final String? email;
  final String? displayName;
  final String? accessToken;
  final String? refreshToken;
  final String role;

  const AuthState({
    this.isAuthenticated = false,
    this.userId,
    this.email,
    this.displayName,
    this.accessToken,
    this.refreshToken,
    this.role = 'user',
  });

  bool get isAdmin => role == 'admin' || role == 'super_admin';
  bool get isDemo => userId == 'demo-user-id';

  AuthState copyWith({
    bool? isAuthenticated,
    String? userId,
    String? email,
    String? displayName,
    String? accessToken,
    String? refreshToken,
    String? role,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      userId: userId ?? this.userId,
      email: email ?? this.email,
      displayName: displayName ?? this.displayName,
      accessToken: accessToken ?? this.accessToken,
      refreshToken: refreshToken ?? this.refreshToken,
      role: role ?? this.role,
    );
  }
}

/// 인증 상태 Notifier
///
/// gRPC auth-service와 연동된 AuthRepository 사용.
class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier(this._repository, {this.restClient}) : super(const AuthState());

  final AuthRepository _repository;
  final ManPaSikRestClient? restClient;

  /// 로그인 처리 (gRPC AuthService Login)
  Future<bool> login(String email, String password) async {
    final result = await _repository.login(email, password);
    if (result.success &&
        result.accessToken != null &&
        result.refreshToken != null) {
      state = AuthState(
        isAuthenticated: true,
        userId: result.userId,
        email: result.email ?? email,
        displayName: result.displayName ?? email.split('@').first,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
        role: result.role ?? 'user',
      );
      await _setActiveUser(state.userId);
      return true;
    }
    return false;
  }

  /// 회원가입 처리 (gRPC AuthService Register + Login)
  Future<bool> register(
      String email, String password, String displayName) async {
    final result = await _repository.register(email, password, displayName);
    if (result.success &&
        result.accessToken != null &&
        result.refreshToken != null) {
      state = AuthState(
        isAuthenticated: true,
        userId: result.userId,
        email: result.email ?? email,
        displayName: result.displayName ?? displayName,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
        role: result.role ?? 'user',
      );
      await _setActiveUser(state.userId);
      return true;
    }
    return false;
  }

  /// 소셜 로그인 (Google/Apple OAuth)
  Future<bool> socialLogin(String provider, String idToken) async {
    if (restClient == null) return false;
    try {
      final res =
          await restClient!.socialLogin(provider: provider, idToken: idToken);
      final accessToken = _readAuthField(res, 'access_token', 'accessToken');
      final refreshToken = _readAuthField(res, 'refresh_token', 'refreshToken');
      if (accessToken != null && refreshToken != null) {
        state = AuthState(
          isAuthenticated: true,
          userId: _readAuthField(res, 'user_id', 'userId') ?? '',
          email: _readAuthField(res, 'email'),
          displayName: _readAuthField(res, 'display_name', 'displayName'),
          accessToken: accessToken,
          refreshToken: refreshToken,
          role: _readAuthField(res, 'role') ?? 'user',
        );
        await _setActiveUser(state.userId);
        return true;
      }
    } catch (_) {
      // Social login failed
    }
    return false;
  }

  /// 카카오 소셜 로그인 결과를 상태에 반영 (Phase AP-1).
  ///
  /// LoginScreen 에서 KakaoLoginService.login() 으로 받은 결과를 그대로 전달.
  /// 성공 시 인증 상태가 갱신되고 활성 사용자 ID 가 TenantInterceptor 에 저장됨.
  Future<bool> applyKakaoLogin({
    required String accessToken,
    required String refreshToken,
    required String userId,
    String? email,
    String? displayName,
  }) async {
    state = AuthState(
      isAuthenticated: true,
      userId: userId,
      email: email ?? '',
      displayName: displayName ?? (email?.split('@').first ?? 'kakao-user'),
      accessToken: accessToken,
      refreshToken: refreshToken,
      role: 'user',
    );
    await _setActiveUser(state.userId);
    return true;
  }

  /// 게스트 로그인 (둘러보기)
  Future<void> loginAsGuest() async {
    state = const AuthState(
      isAuthenticated: true,
      userId: 'guest-user',
      email: 'guest@example.com',
      displayName: 'Guest',
      accessToken: 'guest-token',
      refreshToken: 'guest-refresh-token',
    );
    await _setActiveUser(state.userId);
  }

  /// 데모 모드 (가상 데이터 체험)
  Future<void> loginAsDemo() async {
    state = const AuthState(
      isAuthenticated: true,
      userId: 'demo-user-id',
      email: 'demo@manpasik.com',
      displayName: '테스트 계정',
      accessToken: 'demo-token',
      role: 'user',
    );
    await _setActiveUser(state.userId);
  }

  bool get isDemo => state.userId == 'demo-user-id';

  /// 로그아웃 (로컬 상태 초기화)
  Future<void> logout() async {
    state = const AuthState();
    await _clearTenantContext();
  }

  /// 초기 인증 상태 확인 (스플래시 화면에서 호출)
  Future<void> checkAuthStatus() async {
    final ok = await _repository.isAuthenticated();
    if (!ok) {
      state = const AuthState();
      await _clearTenantContext();
    }
  }

  Future<void> _setActiveUser(String? userId) async {
    try {
      await TenantInterceptor.setActiveUser(userId);
    } catch (_) {
      // SharedPreferences 초기화 실패가 로그인 상태 자체를 깨지 않도록 격리한다.
    }
  }

  Future<void> _clearTenantContext() async {
    try {
      await TenantInterceptor.clear();
    } catch (_) {
      // SharedPreferences 초기화 실패가 로그아웃 상태 자체를 깨지 않도록 격리한다.
    }
  }

  String? _readAuthField(
    Map<String, dynamic> source,
    String primary, [
    String? alternate,
  ]) {
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

/// 인증 상태 Provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  final client = ref.watch(restClientProvider);
  return AuthNotifier(repository, restClient: client);
});
