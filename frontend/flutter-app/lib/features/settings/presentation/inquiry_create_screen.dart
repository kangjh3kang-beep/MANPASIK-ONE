import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

// ───────────────────────────────────────────────────
// InquiryCreateScreen — Sanggam Orbit 1:1 문의
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] theme.textTheme 1x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold 1x → SanggamTheme.primary
// [Rule 4] DropdownButtonFormField/TextField/SwitchListTile/FilledButton 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 1:1 문의 작성 화면
class InquiryCreateScreen
    extends ConsumerStatefulWidget {
  const InquiryCreateScreen({super.key});

  @override
  ConsumerState<InquiryCreateScreen>
      createState() =>
          _InquiryCreateScreenState();
}

class _InquiryCreateScreenState
    extends ConsumerState<InquiryCreateScreen> {
  String _type = 'device';
  final _titleController =
      TextEditingController();
  final _contentController =
      TextEditingController();
  bool _notifyByPush = true;
  bool _notifyByEmail = false;
  bool _isSubmitting = false;

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  InputDecoration _inputDecoration({
    required String label,
    String? hint,
    bool alignLabel = false,
  }) {
    return InputDecoration(
      labelText: label,
      labelStyle: const TextStyle(
          color: SanggamTheme.onSurfaceDim),
      hintText: hint,
      hintStyle: const TextStyle(
          color: SanggamTheme.onSurfaceDim),
      alignLabelWithHint: alignLabel,
      filled: true,
      fillColor: SanggamTheme.surface,
      border: OutlineInputBorder(
        borderRadius:
            BorderRadius.circular(16),
        borderSide: BorderSide.none,
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius:
            BorderRadius.circular(16),
        borderSide: const BorderSide(
            color:
                SanggamTheme.surfaceVariant),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius:
            BorderRadius.circular(16),
        borderSide: const BorderSide(
            color: SanggamTheme.primary),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '1:1 문의',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: SingleChildScrollView(
                padding:
                    const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .stretch,
                  children: [
                    DropdownButtonFormField<
                        String>(
                      value: _type,
                      dropdownColor:
                          SanggamTheme
                              .surface,
                      style: const TextStyle(
                          color:
                              Colors.white),
                      decoration:
                          _inputDecoration(
                              label:
                                  '문의 유형'),
                      items: const [
                        DropdownMenuItem(
                            value: 'device',
                            child: Text(
                                '기기/카트리지')),
                        DropdownMenuItem(
                            value:
                                'subscription',
                            child: Text(
                                '구독/결제')),
                        DropdownMenuItem(
                            value: 'account',
                            child: Text(
                                '계정/인증')),
                        DropdownMenuItem(
                            value:
                                'measurement',
                            child: Text(
                                '측정/결과')),
                        DropdownMenuItem(
                            value: 'bug',
                            child:
                                Text('오류 신고')),
                        DropdownMenuItem(
                            value: 'other',
                            child:
                                Text('기타')),
                      ],
                      onChanged: (v) =>
                          setState(() =>
                              _type = v ??
                                  'device'),
                    ),
                    const SizedBox(
                        height: 16),
                    TextFormField(
                      controller:
                          _titleController,
                      style: const TextStyle(
                          color:
                              Colors.white),
                      decoration:
                          _inputDecoration(
                        label: '제목',
                        hint:
                            '문의 제목을 입력해주세요',
                      ),
                    ),
                    const SizedBox(
                        height: 16),
                    TextFormField(
                      controller:
                          _contentController,
                      maxLines: 8,
                      style: const TextStyle(
                          color:
                              Colors.white),
                      decoration:
                          _inputDecoration(
                        label: '문의 내용',
                        hint:
                            '문의하실 내용을 상세히 작성해주세요.',
                        alignLabel: true,
                      ),
                    ),
                    const SizedBox(
                        height: 16),

                    // 알림 설정
                    const Text('답변 알림',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 14,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 8),
                    SwitchListTile(
                      title: const Text(
                          '푸시 알림',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      activeColor:
                          SanggamTheme
                              .primary,
                      value: _notifyByPush,
                      onChanged: (v) =>
                          setState(() =>
                              _notifyByPush =
                                  v),
                      dense: true,
                    ),
                    SwitchListTile(
                      title: const Text(
                          '이메일 알림',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      activeColor:
                          SanggamTheme
                              .primary,
                      value: _notifyByEmail,
                      onChanged: (v) =>
                          setState(() =>
                              _notifyByEmail =
                                  v),
                      dense: true,
                    ),
                    const SizedBox(
                        height: 24),
                    FilledButton(
                      onPressed:
                          _isSubmitting
                              ? null
                              : _submit,
                      style: FilledButton
                          .styleFrom(
                        minimumSize:
                            const Size
                                .fromHeight(
                                48),
                        backgroundColor:
                            SanggamTheme
                                .primary,
                        foregroundColor:
                            SanggamTheme
                                .background,
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                      child: _isSubmitting
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                  strokeWidth:
                                      2,
                                  color:
                                      SanggamTheme
                                          .primary))
                          : const Text(
                              '문의 접수'),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _submit() async {
    if (_titleController.text.trim().isEmpty ||
        _contentController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('제목과 내용을 입력해주세요.')));
      return;
    }
    setState(() => _isSubmitting = true);
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.createInquiry(
        userId: userId,
        type: _type,
        title: _titleController.text.trim(),
        content: _contentController.text.trim(),
        notifyByPush: _notifyByPush,
        notifyByEmail: _notifyByEmail,
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('문의가 접수되었습니다.')));
        context.pop();
      }
    } on DioException {
      if (mounted) {
        setState(() => _isSubmitting = false);
        ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('문의 접수에 실패했습니다. 다시 시도해주세요.')));
      }
    }
  }
}
