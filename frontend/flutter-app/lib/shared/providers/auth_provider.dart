import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
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
      return true;
    }
    return false;
  }

  /// 회원가입 처리 (gRPC AuthService Register + Login)
  Future<bool> register(String email, String password, String displayName) async {
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
      return true;
    }
    return false;
  }

  /// 소셜 로그인 (Google/Apple OAuth)
  Future<bool> socialLogin(String provider, String idToken) async {
    if (restClient == null) return false;
    try {
      final res = await restClient!.socialLogin(provider: provider, idToken: idToken);
      final accessToken = res['access_token'] as String?;
      final refreshToken = res['refresh_token'] as String?;
      if (accessToken != null && refreshToken != null) {
        state = AuthState(
          isAuthenticated: true,
          userId: res['user_id'] as String? ?? '',
          email: res['email'] as String?,
          displayName: res['display_name'] as String?,
          accessToken: accessToken,
          refreshToken: refreshToken,
          role: res['role'] as String? ?? 'user',
        );
        return true;
      }
    } catch (_) {
      // Social login failed
    }
    return false;
  }

  /// 게스트 로그인 (둘러보기)
  void loginAsGuest() {
    state = const AuthState(
      isAuthenticated: true,
      userId: 'guest-user',
      email: 'guest@example.com',
      displayName: 'Guest',
      accessToken: 'guest-token',
      refreshToken: 'guest-refresh-token',
    );
  }

  /// 데모 모드 (가상 데이터 체험)
  void loginAsDemo() {
    state = const AuthState(
      isAuthenticated: true,
      userId: 'demo-user-id', 
      email: 'demo@manpasik.com',
      displayName: '테스트 계정',
      accessToken: 'demo-token',
      role: 'user',
    );
  }
  
  bool get isDemo => state.userId == 'demo-user-id';

  /// 로그아웃 (로컬 상태 초기화)
  void logout() {
    state = const AuthState();
  }

  /// 초기 인증 상태 확인 (스플래시 화면에서 호출)
  Future<void> checkAuthStatus() async {
    final ok = await _repository.isAuthenticated();
    if (!ok) {
      state = const AuthState();
    }
  }
}

/// 인증 상태 Provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  final client = ref.watch(restClientProvider);
  return AuthNotifier(repository, restClient: client);
});
