import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/utils/validators.dart';
import 'package:manpasik/features/auth/kakao_login_service.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// LoginScreen — Sanggam Orbit 로그인
//
// [Rule 4] isDark 분기 제거 (항상 dark) → CosmicBackground 단일
// [Rule 4] KoreanEdgeBorder+JagaeContainer → SanggamContainer
// [Rule 4] hanji_background, jagae_pattern, porcelain_container import 제거
// [Rule 4] theme.colorScheme.* ~13건 → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~8건 → 직접 TextStyle
// [Rule 4] withOpacity 1건 → withValues(alpha:)
// [Rule 2] SizedBox(h:12)→16, button height 52→48
// ───────────────────────────────────────────────────

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _isLoading = false;
  bool _obscurePassword = true;
  bool _socialLoading = false;
  bool _kakaoLoading = false;

  // Phase AP-1: KakaoLoginService 인스턴스 (DI 가능하도록 final)
  KakaoLoginService _kakaoService = KakaoLoginService();

  @visibleForTesting
  set kakaoService(KakaoLoginService service) => _kakaoService = service;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleSocialLogin(String provider) async {
    setState(() => _socialLoading = true);

    try {
      final idToken = await _acquireOAuthToken(provider);
      if (idToken == null) {
        if (mounted) {
          setState(() => _socialLoading = false);
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('$provider 인증이 취소되었습니다.'),
              behavior: SnackBarBehavior.floating,
            ),
          );
        }
        return;
      }

      final success = await ref.read(authProvider.notifier).socialLogin(
            provider,
            idToken,
          );
      if (!mounted) return;
      setState(() => _socialLoading = false);
      if (success) {
        context.go('/home');
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('$provider 로그인에 실패했습니다. 다시 시도해주세요.'),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        setState(() => _socialLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('$provider 로그인 오류: $e'),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    }
  }

  Future<String?> _acquireOAuthToken(String provider) async {
    return 'pending-oauth-$provider-${DateTime.now().millisecondsSinceEpoch}';
  }

  /// 카카오 로그인 핸들러 (Phase AP-1).
  ///
  /// `KakaoLoginService.login()` 호출 → 성공 시 `applyKakaoLogin` 으로
  /// authProvider 상태 갱신 → /home 으로 이동. 실패 시 SnackBar 노출.
  Future<void> _handleKakaoLogin() async {
    setState(() => _kakaoLoading = true);
    try {
      final result = await _kakaoService.login();
      if (!mounted) return;
      if (result.success && result.accessToken != null) {
        await ref.read(authProvider.notifier).applyKakaoLogin(
              accessToken: result.accessToken!,
              refreshToken: result.refreshToken ?? '',
              userId: result.userId ?? '',
              email: result.email,
              displayName: result.displayName,
            );
        if (!mounted) return;
        setState(() => _kakaoLoading = false);
        context.go('/home');
      } else {
        setState(() => _kakaoLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(result.errorMessage ?? '카카오 로그인 실패'),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    } catch (e) {
      if (!mounted) return;
      setState(() => _kakaoLoading = false);
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('카카오 로그인 오류: $e'),
          behavior: SnackBarBehavior.floating,
        ),
      );
    }
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final success = await ref.read(authProvider.notifier).login(
          _emailController.text.trim(),
          _passwordController.text,
        );

    if (!mounted) return;
    setState(() => _isLoading = false);

    if (success) {
      context.go('/home');
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('로그인에 실패했습니다. 이메일과 비밀번호를 확인해주세요.'),
          behavior: SnackBarBehavior.floating,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: CosmicBackground(
        child: SafeArea(
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: SanggamContainer(
                borderRadius: 24,
                borderWidth: 1.0,
                borderColor:
                    SanggamTheme.primary.withValues(alpha: 0.3),
                blurSigma: 16,
                jagaeOpacity: 0.05,
                padding: const EdgeInsets.all(32),
                child: _buildLoginForm(),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildLoginForm() {
    return FocusTraversalGroup(
      child: Form(
        key: _formKey,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 로고 영역
            Semantics(
              label: '만파식 AI 헬스케어 로고',
              child: const Icon(
                Icons.biotech_rounded,
                size: 64,
                color: SanggamTheme.primary,
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              '만파식',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: SanggamTheme.primary,
                fontSize: 28,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '건강한 일상을 위한 AI 헬스케어',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 48),

            // 이메일 입력
            TextFormField(
              controller: _emailController,
              keyboardType: TextInputType.emailAddress,
              textInputAction: TextInputAction.next,
              validator: Validators.validateEmail,
              decoration: InputDecoration(
                labelText: '이메일',
                hintText: 'example@manpasik.com',
                prefixIcon: const Icon(Icons.email_outlined),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 비밀번호 입력
            TextFormField(
              controller: _passwordController,
              obscureText: _obscurePassword,
              textInputAction: TextInputAction.done,
              validator: Validators.validatePassword,
              onFieldSubmitted: (_) => _handleLogin(),
              decoration: InputDecoration(
                labelText: '비밀번호',
                hintText: '8자 이상 (영문 + 숫자)',
                prefixIcon: const Icon(Icons.lock_outlined),
                suffixIcon: IconButton(
                  tooltip: '비밀번호 표시/숨김',
                  icon: Icon(
                    _obscurePassword
                        ? Icons.visibility_outlined
                        : Icons.visibility_off_outlined,
                  ),
                  onPressed: () {
                    setState(
                        () => _obscurePassword = !_obscurePassword);
                  },
                ),
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
            ),
            const SizedBox(height: 24),

            // 로그인 버튼
            PrimaryButton(
              text: '로그인',
              isLoading: _isLoading,
              onPressed: _handleLogin,
            ),
            const SizedBox(height: 8),

            // 비밀번호 찾기
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: () => context.push('/forgot-password'),
                child: const Text(
                  '비밀번호를 잊으셨나요?',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 소셜 로그인 구분선
            Row(
              children: [
                Expanded(
                    child: Divider(
                        color: SanggamTheme.onSurfaceDim
                            .withValues(alpha: 0.3))),
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 16),
                  child: Text(
                    '또는',
                    style: TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: 12,
                    ),
                  ),
                ),
                Expanded(
                    child: Divider(
                        color: SanggamTheme.onSurfaceDim
                            .withValues(alpha: 0.3))),
              ],
            ),
            const SizedBox(height: 16),

            // Google 소셜 로그인
            OutlinedButton.icon(
              onPressed: _socialLoading
                  ? null
                  : () => _handleSocialLogin('google'),
              icon: const Icon(Icons.g_mobiledata, size: 24),
              label: const Text('Google로 계속하기'),
              style: OutlinedButton.styleFrom(
                minimumSize: const Size(double.infinity, 48),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                side: BorderSide(
                    color: SanggamTheme.onSurfaceDim
                        .withValues(alpha: 0.3)),
              ),
            ),
            const SizedBox(height: 16),

            // Apple 소셜 로그인
            OutlinedButton.icon(
              onPressed: _socialLoading
                  ? null
                  : () => _handleSocialLogin('apple'),
              icon: const Icon(Icons.apple, size: 24),
              label: const Text('Apple로 계속하기'),
              style: OutlinedButton.styleFrom(
                minimumSize: const Size(double.infinity, 48),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                side: BorderSide(
                    color: SanggamTheme.onSurfaceDim
                        .withValues(alpha: 0.3)),
              ),
            ),
            const SizedBox(height: 24),

            // 회원가입 링크
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Text(
                  '계정이 없으신가요?',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 14,
                  ),
                ),
                TextButton(
                  onPressed: () => context.push('/register'),
                  child: const Text('회원가입'),
                ),
              ],
            ),
            const SizedBox(height: 16),

            // 카카오 로그인 버튼 (Phase AP-1 — codex G1 통합)
            Semantics(
              button: true,
              label: '카카오로 시작 버튼',
              child: ElevatedButton.icon(
                key: const Key('kakao-login-button'),
                onPressed: _kakaoLoading ? null : _handleKakaoLogin,
                icon: _kakaoLoading
                    ? const SizedBox(
                        width: 16,
                        height: 16,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Icon(Icons.chat_bubble, color: Color(0xFF3C1E1E)),
                label: const Text(
                  '카카오로 시작',
                  style: TextStyle(
                    color: Color(0xFF3C1E1E),
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFFFEE500),
                  minimumSize: const Size.fromHeight(48),
                ),
              ),
            ),
            const SizedBox(height: 12),

            // 가상 체험 버튼 (Demo Mode)
            Semantics(
              button: true,
              label: '가상 데이터 체험 시작 버튼',
              child: TextButton.icon(
                onPressed: () {
                  ref.read(authProvider.notifier).loginAsDemo();
                  context.go('/home');
                },
                icon: const Icon(Icons.science_outlined,
                    color: SanggamTheme.jagaeCyan),
                label: const Text(
                  '가상 데이터 체험 시작',
                  style: TextStyle(
                    color: SanggamTheme.jagaeCyan,
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
            const SizedBox(height: 8),

            // 둘러보기 버튼
            Semantics(
              button: true,
              label: '비회원 둘러보기 버튼',
              child: TextButton(
                onPressed: () {
                  ref.read(authProvider.notifier).loginAsGuest();
                  context.go('/home');
                },
                child: const Text(
                  '비회원 둘러보기',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 14,
                    decoration: TextDecoration.underline,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
