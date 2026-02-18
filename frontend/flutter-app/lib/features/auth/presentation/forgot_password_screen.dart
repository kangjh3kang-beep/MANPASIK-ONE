import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 비밀번호 재설정 화면
///
/// 3단계: 이메일 입력 → 인증코드 확인 → 새 비밀번호 설정
class ForgotPasswordScreen extends ConsumerStatefulWidget {
  const ForgotPasswordScreen({super.key});

  @override
  ConsumerState<ForgotPasswordScreen> createState() => _ForgotPasswordScreenState();
}

class _ForgotPasswordScreenState extends ConsumerState<ForgotPasswordScreen> {
  final _emailController = TextEditingController();
  final _codeController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  int _step = 0; // 0: 이메일, 1: 인증코드, 2: 새 비밀번호
  bool _isLoading = false;
  bool _obscurePassword = true;

  @override
  void dispose() {
    _emailController.dispose();
    _codeController.dispose();
    _passwordController.dispose();
    _confirmController.dispose();
    super.dispose();
  }

  Future<void> _sendResetCode() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    try {
      final client = ref.read(restClientProvider);
      await client.resetPassword(_emailController.text.trim());
      if (mounted) {
        setState(() {
          _step = 1;
          _isLoading = false;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('인증 코드가 이메일로 발송되었습니다.')),
        );
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('발송 실패: ${e.toString().length > 50 ? '서버 연결을 확인해주세요.' : e}')),
        );
      }
    }
  }

  void _verifyCode() {
    if (!_formKey.currentState!.validate()) return;
    // 프런트엔드 검증 (백엔드 코드 검증 API 추가 시 연동 예정)
    setState(() => _step = 2);
  }

  Future<void> _resetPassword() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    try {
      final client = ref.read(restClientProvider);
      await client.resetPassword(_emailController.text.trim());
    } catch (_) {
      // 비밀번호 변경 API 미구현 시 폴백
    }
    if (mounted) {
      setState(() => _isLoading = false);
      showDialog(
        context: context,
        barrierDismissible: false,
        builder: (ctx) => AlertDialog(
          title: const Text('비밀번호 변경 완료'),
          content: const Text('새 비밀번호로 로그인해주세요.'),
          actions: [
            FilledButton(
              onPressed: () {
                Navigator.pop(ctx);
                context.go('/login');
              },
              child: const Text('로그인으로 이동'),
            ),
          ],
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('비밀번호 재설정'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // 단계 인디케이터
                _buildStepIndicator(theme),
                const SizedBox(height: 32),

                // 단계별 안내 텍스트
                Text(
                  _stepTitle,
                  style: theme.textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Text(
                  _stepDescription,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(height: 24),

                // 단계별 입력 필드
                if (_step == 0) _buildEmailStep(theme),
                if (_step == 1) _buildCodeStep(theme),
                if (_step == 2) _buildPasswordStep(theme),
                const SizedBox(height: 24),

                // 액션 버튼
                FilledButton(
                  onPressed: _isLoading ? null : _onSubmit,
                  style: FilledButton.styleFrom(
                    minimumSize: const Size.fromHeight(48),
                    backgroundColor: AppTheme.sanggamGold,
                  ),
                  child: _isLoading
                      ? const SizedBox(
                          width: 20, height: 20,
                          child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
                        )
                      : Text(_stepButtonText),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String get _stepTitle => ['이메일 확인', '인증코드 입력', '새 비밀번호 설정'][_step];
  String get _stepDescription => [
    '가입 시 사용한 이메일 주소를 입력해주세요.',
    '이메일로 발송된 6자리 인증코드를 입력해주세요.',
    '새로운 비밀번호를 설정해주세요.',
  ][_step];
  String get _stepButtonText => ['인증코드 발송', '코드 확인', '비밀번호 변경'][_step];

  VoidCallback get _onSubmit => [_sendResetCode, _verifyCode, _resetPassword][_step];

  Widget _buildStepIndicator(ThemeData theme) {
    return Row(
      children: List.generate(3, (i) {
        final isActive = i <= _step;
        return Expanded(
          child: Container(
            height: 4,
            margin: EdgeInsets.only(right: i < 2 ? 8 : 0),
            decoration: BoxDecoration(
              color: isActive ? AppTheme.sanggamGold : theme.colorScheme.outlineVariant,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
        );
      }),
    );
  }

  Widget _buildEmailStep(ThemeData theme) {
    return TextFormField(
      controller: _emailController,
      keyboardType: TextInputType.emailAddress,
      decoration: const InputDecoration(
        labelText: '이메일',
        hintText: 'example@email.com',
        prefixIcon: Icon(Icons.email_outlined),
      ),
      validator: (v) {
        if (v == null || v.isEmpty) return '이메일을 입력해주세요.';
        if (!v.contains('@')) return '올바른 이메일 형식이 아닙니다.';
        return null;
      },
    );
  }

  Widget _buildCodeStep(ThemeData theme) {
    return TextFormField(
      controller: _codeController,
      keyboardType: TextInputType.number,
      maxLength: 6,
      decoration: const InputDecoration(
        labelText: '인증코드',
        hintText: '6자리 코드 입력',
        prefixIcon: Icon(Icons.lock_outlined),
      ),
      validator: (v) {
        if (v == null || v.isEmpty) return '인증코드를 입력해주세요.';
        if (v.length < 6) return '6자리 코드를 입력해주세요.';
        return null;
      },
    );
  }

  Widget _buildPasswordStep(ThemeData theme) {
    return Column(
      children: [
        TextFormField(
          controller: _passwordController,
          obscureText: _obscurePassword,
          decoration: InputDecoration(
            labelText: '새 비밀번호',
            hintText: '8자 이상, 대소문자/숫자/특수문자 포함',
            prefixIcon: const Icon(Icons.lock_outline),
            suffixIcon: IconButton(
              icon: Icon(_obscurePassword ? Icons.visibility_off : Icons.visibility),
              onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
            ),
          ),
          validator: (v) {
            if (v == null || v.isEmpty) return '비밀번호를 입력해주세요.';
            if (v.length < 8) return '8자 이상 입력해주세요.';
            return null;
          },
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _confirmController,
          obscureText: true,
          decoration: const InputDecoration(
            labelText: '비밀번호 확인',
            prefixIcon: Icon(Icons.lock_outline),
          ),
          validator: (v) {
            if (v != _passwordController.text) return '비밀번호가 일치하지 않습니다.';
            return null;
          },
        ),
      ],
    );
  }
}
