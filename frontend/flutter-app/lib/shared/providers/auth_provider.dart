import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/auth/domain/auth_repository.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';

/// 인증 상태 모델
class AuthState {
  final bool isAuthenticated;
  final String? userId;
  final String? email;
  final String? displayName;
  final String? accessToken;
  final String? refreshToken;

  const AuthState({
    this.isAuthenticated = false,
    this.userId,
    this.email,
    this.displayName,
    this.accessToken,
    this.refreshToken,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    String? userId,
    String? email,
    String? displayName,
    String? accessToken,
    String? refreshToken,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      userId: userId ?? this.userId,
      email: email ?? this.email,
      displayName: displayName ?? this.displayName,
      accessToken: accessToken ?? this.accessToken,
      refreshToken: refreshToken ?? this.refreshToken,
    );
  }
}

/// 인증 상태 Notifier
///
/// gRPC auth-service와 연동된 AuthRepository 사용.
class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier(this._repository) : super(const AuthState());

  final AuthRepository _repository;

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
      );
      return true;
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
  return AuthNotifier(repository);
});
