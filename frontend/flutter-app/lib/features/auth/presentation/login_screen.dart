
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/utils/validators.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:manpasik/shared/widgets/hanji_background.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';

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

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  bool _socialLoading = false;

  Future<void> _handleSocialLogin(String provider) async {
    setState(() => _socialLoading = true);

    try {
      // 1단계: 플랫폼별 OAuth SDK로 idToken 획득
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

      // 2단계: idToken → 백엔드 socialLogin 엔드포인트로 교환
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

  /// 플랫폼별 OAuth SDK에서 ID 토큰을 획득합니다.
  ///
  /// google_sign_in / sign_in_with_apple 패키지 설치 시 실제 SDK 호출,
  /// 미설치 또는 Web 환경에서는 REST 기반 시뮬레이션 폴백.
  Future<String?> _acquireOAuthToken(String provider) async {
    // google_sign_in / sign_in_with_apple 패키지 연동 지점
    // 패키지 설치 후 아래 주석 해제:
    //
    // if (provider == 'google') {
    //   final googleUser = await GoogleSignIn(scopes: ['email', 'profile']).signIn();
    //   return googleUser?.authentication.then((auth) => auth.idToken);
    // }
    // if (provider == 'apple') {
    //   final credential = await SignInWithApple.getAppleIDCredential(
    //     scopes: [AppleIDAuthorizationScopes.email, AppleIDAuthorizationScopes.fullName],
    //   );
    //   return credential.identityToken;
    // }

    // 시뮬레이션 폴백: OAuth 미연동 시 pending 토큰으로 서버 호출
    return 'pending-oauth-$provider-${DateTime.now().millisecondsSinceEpoch}';
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
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: isDark 
          ? CosmicBackground(
              child: _buildLoginBody(context, theme, isDark),
            )
          : HanjiBackground(
              child: _buildLoginBody(context, theme, isDark),
            ),
    );
  }

  Widget _buildLoginBody(BuildContext context, ThemeData theme, bool isDark) {
    return SafeArea(
      child: Center(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(horizontal: 24),
          child: KoreanEdgeBorder(
            borderRadius: BorderRadius.circular(24),
            child: isDark
                ? JagaeContainer(
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surface.withOpacity(0.1), // Glassmorphic transparency
                      borderRadius: BorderRadius.circular(24),
                      boxShadow: [
                        BoxShadow(
                          color: theme.shadowColor.withValues(alpha: 0.1),
                          blurRadius: 20,
                          offset: const Offset(0, 10),
                        ),
                      ],
                    ),
                    padding: const EdgeInsets.all(32),
                    child: _buildLoginForm(context, theme),
                  )
                : PorcelainContainer(
                    padding: const EdgeInsets.all(32),
                    child: _buildLoginForm(context, theme),
                  ),
          ),
        ),
      ),
    );
  }

  Widget _buildLoginForm(BuildContext context, ThemeData theme) {
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
                        child: Icon(
                          Icons.biotech_rounded,
                          size: 64,
                          color: theme.colorScheme.primary,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        '만파식',
                        textAlign: TextAlign.center,
                        style: theme.textTheme.headlineMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: theme.colorScheme.primary,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '건강한 일상을 위한 AI 헬스케어',
                        textAlign: TextAlign.center,
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
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
                            icon: Icon(
                              _obscurePassword
                                  ? Icons.visibility_outlined
                                  : Icons.visibility_off_outlined,
                            ),
                            onPressed: () {
                              setState(() => _obscurePassword = !_obscurePassword);
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
                          child: Text(
                            '비밀번호를 잊으셨나요?',
                            style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),

                      // 소셜 로그인 구분선
                      Row(
                        children: [
                          Expanded(child: Divider(color: theme.colorScheme.outlineVariant)),
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 16),
                            child: Text(
                              '또는',
                              style: theme.textTheme.bodySmall?.copyWith(
                                color: theme.colorScheme.onSurfaceVariant,
                              ),
                            ),
                          ),
                          Expanded(child: Divider(color: theme.colorScheme.outlineVariant)),
                        ],
                      ),
                      const SizedBox(height: 16),

                      // Google 소셜 로그인
                      OutlinedButton.icon(
                        onPressed: _socialLoading ? null : () => _handleSocialLogin('google'),
                        icon: const Icon(Icons.g_mobiledata, size: 24),
                        label: const Text('Google로 계속하기'),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 52),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                          side: BorderSide(color: theme.colorScheme.outlineVariant),
                        ),
                      ),
                      const SizedBox(height: 12),

                      // Apple 소셜 로그인
                      OutlinedButton.icon(
                        onPressed: _socialLoading ? null : () => _handleSocialLogin('apple'),
                        icon: const Icon(Icons.apple, size: 24),
                        label: const Text('Apple로 계속하기'),
                        style: OutlinedButton.styleFrom(
                          minimumSize: const Size(double.infinity, 52),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                          side: BorderSide(color: theme.colorScheme.outlineVariant),
                        ),
                      ),
                      const SizedBox(height: 24),

                      // 회원가입 링크
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            '계정이 없으신가요?',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
                            ),
                          ),
                          TextButton(
                            onPressed: () => context.push('/register'),
                            child: const Text('회원가입'),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),

                      // 가상 체험 버튼 (Demo Mode)
                      Semantics(
                        button: true,
                        label: '가상 데이터 체험 시작 버튼',
                        child: TextButton.icon(
                          onPressed: () {
                            ref.read(authProvider.notifier).loginAsDemo();
                            context.go('/home');
                          },
                          icon: Icon(Icons.science_outlined, color: theme.colorScheme.secondary),
                          label: Text(
                            '가상 데이터 체험 시작',
                            style: theme.textTheme.labelLarge?.copyWith(
                              color: theme.colorScheme.secondary,
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
                          child: Text(
                            '비회원 둘러보기',
                            style: theme.textTheme.labelLarge?.copyWith(
                              color: theme.colorScheme.onSurfaceVariant,
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
} // End of LoginScreen
