import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/utils/validators.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// RegisterScreen — Sanggam Orbit 회원가입
//
// [Rule 4] AppBar 제거 → body 내 뒤로가기 버튼
// [Rule 4] theme.colorScheme.* ~10건 → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~8건 → 직접 TextStyle
// [Rule 4] CosmicBackground + SanggamContainer 적용
// [Rule 2] SizedBox(h:12)→16, borderRadius 12→16
// ───────────────────────────────────────────────────

class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  bool _isLoading = false;
  bool _obscurePassword = true;
  bool _obscureConfirm = true;

  // 약관 동의
  bool _termsAgreed = false;
  bool _privacyAgreed = false;
  bool _healthDataAgreed = false;
  bool _marketingAgreed = false;

  bool get _requiredConsentsChecked =>
      _termsAgreed && _privacyAgreed && _healthDataAgreed;

  bool _socialLoading = false;

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  Future<void> _handleSocialRegister(String provider) async {
    setState(() => _socialLoading = true);
    final success = await ref.read(authProvider.notifier).socialLogin(
          provider,
          'pending-oauth-flow',
        );
    if (!mounted) return;
    setState(() => _socialLoading = false);
    if (success) {
      context.go('/onboarding');
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('$provider 가입에 실패했습니다. 다시 시도해주세요.'),
          behavior: SnackBarBehavior.floating,
        ),
      );
    }
  }

  Future<void> _handleRegister() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final success = await ref.read(authProvider.notifier).register(
          _emailController.text.trim(),
          _passwordController.text,
          _nameController.text.trim(),
        );

    if (!mounted) return;
    setState(() => _isLoading = false);

    if (success) {
      context.go('/onboarding');
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('회원가입에 실패했습니다. 다시 시도해주세요.'),
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
          child: Column(
            children: [
              // 헤더: 뒤로가기 + 타이틀
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8),
                child: SizedBox(
                  height: 48,
                  child: Row(
                    children: [
                      IconButton(
                        icon: const Icon(Icons.arrow_back,
                            color: Colors.white),
                        tooltip: '뒤로 가기',
                        onPressed: () => context.pop(),
                      ),
                      const SizedBox(width: 8),
                      const Text(
                        '회원가입',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ),
              ),

              // 폼 영역
              Expanded(
                child: SingleChildScrollView(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: SanggamContainer(
                    borderRadius: 24,
                    borderWidth: 1.0,
                    borderColor:
                        SanggamTheme.primary.withValues(alpha: 0.3),
                    blurSigma: 16,
                    jagaeOpacity: 0.05,
                    padding: const EdgeInsets.all(24),
                    child: Form(
                      key: _formKey,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.stretch,
                        children: [
                          const SizedBox(height: 8),

                          // 이름
                          TextFormField(
                            controller: _nameController,
                            textInputAction: TextInputAction.next,
                            validator: Validators.validateDisplayName,
                            decoration: InputDecoration(
                              labelText: '이름',
                              hintText: '표시될 이름을 입력해주세요',
                              prefixIcon:
                                  const Icon(Icons.person_outlined),
                              border: OutlineInputBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),

                          // 이메일
                          TextFormField(
                            controller: _emailController,
                            keyboardType: TextInputType.emailAddress,
                            textInputAction: TextInputAction.next,
                            validator: Validators.validateEmail,
                            decoration: InputDecoration(
                              labelText: '이메일',
                              hintText: 'example@manpasik.com',
                              prefixIcon:
                                  const Icon(Icons.email_outlined),
                              border: OutlineInputBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),

                          // 비밀번호
                          TextFormField(
                            controller: _passwordController,
                            obscureText: _obscurePassword,
                            textInputAction: TextInputAction.next,
                            validator: Validators.validatePassword,
                            decoration: InputDecoration(
                              labelText: '비밀번호',
                              hintText: '8자 이상 (영문 + 숫자)',
                              prefixIcon:
                                  const Icon(Icons.lock_outlined),
                              suffixIcon: IconButton(
                                tooltip: '비밀번호 표시/숨김',
                                icon: Icon(
                                  _obscurePassword
                                      ? Icons.visibility_outlined
                                      : Icons.visibility_off_outlined,
                                ),
                                onPressed: () {
                                  setState(() => _obscurePassword =
                                      !_obscurePassword);
                                },
                              ),
                              border: OutlineInputBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),

                          // 비밀번호 확인
                          TextFormField(
                            controller: _confirmPasswordController,
                            obscureText: _obscureConfirm,
                            textInputAction: TextInputAction.done,
                            validator: (value) {
                              if (value != _passwordController.text) {
                                return '비밀번호가 일치하지 않습니다';
                              }
                              return null;
                            },
                            onFieldSubmitted: (_) =>
                                _handleRegister(),
                            decoration: InputDecoration(
                              labelText: '비밀번호 확인',
                              hintText: '비밀번호를 다시 입력해주세요',
                              prefixIcon:
                                  const Icon(Icons.lock_outlined),
                              suffixIcon: IconButton(
                                tooltip: '비밀번호 표시/숨김',
                                icon: Icon(
                                  _obscureConfirm
                                      ? Icons.visibility_outlined
                                      : Icons.visibility_off_outlined,
                                ),
                                onPressed: () {
                                  setState(() => _obscureConfirm =
                                      !_obscureConfirm);
                                },
                              ),
                              border: OutlineInputBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                            ),
                          ),
                          const SizedBox(height: 24),

                          // ── 약관 동의 섹션 ──
                          Container(
                            padding: const EdgeInsets.all(16),
                            decoration: BoxDecoration(
                              color: SanggamTheme.surface
                                  .withValues(alpha: 0.5),
                              borderRadius: BorderRadius.circular(16),
                              border: Border.all(
                                color:
                                    !_requiredConsentsChecked &&
                                            _isLoading
                                        ? SanggamTheme.error
                                        : SanggamTheme.onSurfaceDim
                                            .withValues(alpha: 0.3),
                              ),
                            ),
                            child: Column(
                              crossAxisAlignment:
                                  CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  '약관 동의',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 14,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const SizedBox(height: 8),

                                // 전체 동의
                                CheckboxListTile(
                                  value: _termsAgreed &&
                                      _privacyAgreed &&
                                      _healthDataAgreed &&
                                      _marketingAgreed,
                                  onChanged: (v) {
                                    setState(() {
                                      _termsAgreed = v ?? false;
                                      _privacyAgreed = v ?? false;
                                      _healthDataAgreed = v ?? false;
                                      _marketingAgreed = v ?? false;
                                    });
                                  },
                                  title: const Text(
                                    '전체 동의',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontSize: 14,
                                      fontWeight: FontWeight.w600,
                                    ),
                                  ),
                                  dense: true,
                                  contentPadding: EdgeInsets.zero,
                                  controlAffinity:
                                      ListTileControlAffinity.leading,
                                ),
                                const Divider(height: 1),

                                _consentTile(
                                  '[필수] 서비스 이용약관',
                                  _termsAgreed,
                                  (v) => setState(() =>
                                      _termsAgreed = v ?? false),
                                ),
                                _consentTile(
                                  '[필수] 개인정보 처리방침',
                                  _privacyAgreed,
                                  (v) => setState(() =>
                                      _privacyAgreed = v ?? false),
                                ),
                                _consentTile(
                                  '[필수] 건강정보 수집 및 이용 동의',
                                  _healthDataAgreed,
                                  (v) => setState(() =>
                                      _healthDataAgreed = v ?? false),
                                ),
                                _consentTile(
                                  '[선택] 마케팅 정보 수신 동의',
                                  _marketingAgreed,
                                  (v) => setState(() =>
                                      _marketingAgreed = v ?? false),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 24),

                          // 가입 버튼
                          PrimaryButton(
                            text: '가입하기',
                            isLoading: _isLoading,
                            onPressed: _requiredConsentsChecked
                                ? _handleRegister
                                : null,
                          ),
                          if (!_requiredConsentsChecked)
                            const Padding(
                              padding: EdgeInsets.only(top: 8),
                              child: Text(
                                '필수 약관에 모두 동의해주세요.',
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: SanggamTheme.error,
                                  fontSize: 12,
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
                                padding: EdgeInsets.symmetric(
                                    horizontal: 16),
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

                          // Google 소셜 가입
                          OutlinedButton.icon(
                            onPressed: _socialLoading
                                ? null
                                : () =>
                                    _handleSocialRegister('google'),
                            icon: const Icon(
                                Icons.g_mobiledata, size: 24),
                            label: const Text('Google로 가입하기'),
                            style: OutlinedButton.styleFrom(
                              minimumSize:
                                  const Size(double.infinity, 48),
                              shape: RoundedRectangleBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                              side: BorderSide(
                                  color: SanggamTheme.onSurfaceDim
                                      .withValues(alpha: 0.3)),
                            ),
                          ),
                          const SizedBox(height: 16),

                          // Apple 소셜 가입
                          OutlinedButton.icon(
                            onPressed: _socialLoading
                                ? null
                                : () =>
                                    _handleSocialRegister('apple'),
                            icon:
                                const Icon(Icons.apple, size: 24),
                            label: const Text('Apple로 가입하기'),
                            style: OutlinedButton.styleFrom(
                              minimumSize:
                                  const Size(double.infinity, 48),
                              shape: RoundedRectangleBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                              side: BorderSide(
                                  color: SanggamTheme.onSurfaceDim
                                      .withValues(alpha: 0.3)),
                            ),
                          ),
                          const SizedBox(height: 24),

                          // 로그인 링크
                          Row(
                            mainAxisAlignment:
                                MainAxisAlignment.center,
                            children: [
                              const Text(
                                '이미 계정이 있으신가요?',
                                style: TextStyle(
                                  color: SanggamTheme.onSurfaceDim,
                                  fontSize: 14,
                                ),
                              ),
                              TextButton(
                                onPressed: () => context.pop(),
                                child: const Text('로그인'),
                              ),
                            ],
                          ),
                          const SizedBox(height: 16),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _consentTile(
      String title, bool value, ValueChanged<bool?> onChanged) {
    return CheckboxListTile(
      value: value,
      onChanged: onChanged,
      title: Text(
        title,
        style: const TextStyle(
          color: SanggamTheme.onSurfaceDim,
          fontSize: 12,
        ),
      ),
      dense: true,
      contentPadding: EdgeInsets.zero,
      controlAffinity: ListTileControlAffinity.leading,
    );
  }
}
