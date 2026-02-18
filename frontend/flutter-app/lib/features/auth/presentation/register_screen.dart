import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/utils/validators.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';

/// 회원가입 화면
///
/// 이메일, 비밀번호, 이름 입력 폼.
/// S5에서 gRPC auth-service Register RPC 연동 예정.
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

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }

  bool _socialLoading = false;

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
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('회원가입'),
      ),
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  const SizedBox(height: 16),

                  // 이름
                  TextFormField(
                    controller: _nameController,
                    textInputAction: TextInputAction.next,
                    validator: Validators.validateDisplayName,
                    decoration: InputDecoration(
                      labelText: '이름',
                      hintText: '표시될 이름을 입력해주세요',
                      prefixIcon: const Icon(Icons.person_outlined),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(16),
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
                      prefixIcon: const Icon(Icons.email_outlined),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(16),
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
                    onFieldSubmitted: (_) => _handleRegister(),
                    decoration: InputDecoration(
                      labelText: '비밀번호 확인',
                      hintText: '비밀번호를 다시 입력해주세요',
                      prefixIcon: const Icon(Icons.lock_outlined),
                      suffixIcon: IconButton(
                        icon: Icon(
                          _obscureConfirm
                              ? Icons.visibility_outlined
                              : Icons.visibility_off_outlined,
                        ),
                        onPressed: () {
                          setState(() => _obscureConfirm = !_obscureConfirm);
                        },
                      ),
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // ── 약관 동의 섹션 ──
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceContainerLow,
                      borderRadius: BorderRadius.circular(12),
                      border: Border.all(
                        color: !_requiredConsentsChecked && _isLoading
                            ? theme.colorScheme.error
                            : theme.colorScheme.outlineVariant,
                      ),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '약관 동의',
                          style: theme.textTheme.titleSmall?.copyWith(
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
                          title: Text(
                            '전체 동의',
                            style: theme.textTheme.bodyMedium?.copyWith(
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          controlAffinity: ListTileControlAffinity.leading,
                        ),
                        const Divider(height: 1),

                        // 서비스 이용약관 (필수)
                        CheckboxListTile(
                          value: _termsAgreed,
                          onChanged: (v) =>
                              setState(() => _termsAgreed = v ?? false),
                          title: Text(
                            '[필수] 서비스 이용약관',
                            style: theme.textTheme.bodySmall,
                          ),
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          controlAffinity: ListTileControlAffinity.leading,
                        ),

                        // 개인정보처리방침 (필수)
                        CheckboxListTile(
                          value: _privacyAgreed,
                          onChanged: (v) =>
                              setState(() => _privacyAgreed = v ?? false),
                          title: Text(
                            '[필수] 개인정보 처리방침',
                            style: theme.textTheme.bodySmall,
                          ),
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          controlAffinity: ListTileControlAffinity.leading,
                        ),

                        // 건강정보 수집 동의 (필수)
                        CheckboxListTile(
                          value: _healthDataAgreed,
                          onChanged: (v) =>
                              setState(() => _healthDataAgreed = v ?? false),
                          title: Text(
                            '[필수] 건강정보 수집 및 이용 동의',
                            style: theme.textTheme.bodySmall,
                          ),
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          controlAffinity: ListTileControlAffinity.leading,
                        ),

                        // 마케팅 동의 (선택)
                        CheckboxListTile(
                          value: _marketingAgreed,
                          onChanged: (v) =>
                              setState(() => _marketingAgreed = v ?? false),
                          title: Text(
                            '[선택] 마케팅 정보 수신 동의',
                            style: theme.textTheme.bodySmall,
                          ),
                          dense: true,
                          contentPadding: EdgeInsets.zero,
                          controlAffinity: ListTileControlAffinity.leading,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // 가입 버튼
                  PrimaryButton(
                    text: '가입하기',
                    isLoading: _isLoading,
                    onPressed: _requiredConsentsChecked ? _handleRegister : null,
                  ),
                  if (!_requiredConsentsChecked)
                    Padding(
                      padding: const EdgeInsets.only(top: 8),
                      child: Text(
                        '필수 약관에 모두 동의해주세요.',
                        textAlign: TextAlign.center,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: theme.colorScheme.error,
                        ),
                      ),
                    ),
                  const SizedBox(height: 16),

                  // 소셜 로그인 구분선
                  Row(
                    children: [
                      Expanded(
                          child: Divider(
                              color: theme.colorScheme.outlineVariant)),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16),
                        child: Text(
                          '또는',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ),
                      Expanded(
                          child: Divider(
                              color: theme.colorScheme.outlineVariant)),
                    ],
                  ),
                  const SizedBox(height: 16),

                  // Google 소셜 가입
                  OutlinedButton.icon(
                    onPressed: _socialLoading ? null : () => _handleSocialRegister('google'),
                    icon: const Icon(Icons.g_mobiledata, size: 24),
                    label: const Text('Google로 가입하기'),
                    style: OutlinedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      side: BorderSide(
                          color: theme.colorScheme.outlineVariant),
                    ),
                  ),
                  const SizedBox(height: 12),

                  // Apple 소셜 가입
                  OutlinedButton.icon(
                    onPressed: _socialLoading ? null : () => _handleSocialRegister('apple'),
                    icon: const Icon(Icons.apple, size: 24),
                    label: const Text('Apple로 가입하기'),
                    style: OutlinedButton.styleFrom(
                      minimumSize: const Size(double.infinity, 48),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      side: BorderSide(
                          color: theme.colorScheme.outlineVariant),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // 로그인 링크
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        '이미 계정이 있으신가요?',
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                      TextButton(
                        onPressed: () => context.pop(),
                        child: const Text('로그인'),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
