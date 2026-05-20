import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// CreatePostScreen — Sanggam Orbit 게시글 작성
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.bodySmall → 직접 TextStyle
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 게시글 작성 화면
class CreatePostScreen extends ConsumerStatefulWidget {
  const CreatePostScreen({super.key});

  @override
  ConsumerState<CreatePostScreen> createState() =>
      _CreatePostScreenState();
}

class _CreatePostScreenState extends ConsumerState<CreatePostScreen> {
  final _titleController = TextEditingController();
  final _contentController = TextEditingController();
  String _category = 'general';
  bool _isAnonymous = false;
  bool _isSubmitting = false;

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (_titleController.text.trim().isEmpty ||
        _contentController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('제목과 내용을 입력해주세요.')),
      );
      return;
    }
    setState(() => _isSubmitting = true);
    try {
      final client = ref.read(restClientProvider);
      await client.createPost(
        authorId: 'current-user',
        title: _titleController.text.trim(),
        content: _contentController.text.trim(),
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('게시글이 작성되었습니다.')),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isSubmitting = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('게시글 작성에 실패했습니다.')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더: 닫기 + 타이틀 + 등록 버튼
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.close,
                        color: Colors.white),
                    tooltip: '닫기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '글쓰기',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  FilledButton(
                    onPressed: _isSubmitting ? null : _submit,
                    style: FilledButton.styleFrom(
                      backgroundColor: SanggamTheme.primary,
                      foregroundColor: SanggamTheme.background,
                      shape: RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius.circular(16)),
                    ),
                    child: _isSubmitting
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                                strokeWidth: 2),
                          )
                        : const Text('등록'),
                  ),
                  const SizedBox(width: 8),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    DropdownButtonFormField<String>(
                      value: _category,
                      decoration:
                          const InputDecoration(labelText: '카테고리'),
                      items: const [
                        DropdownMenuItem(
                            value: 'general', child: Text('전체')),
                        DropdownMenuItem(
                            value: 'health_tips',
                            child: Text('건강 팁')),
                        DropdownMenuItem(
                            value: 'reviews',
                            child: Text('측정 후기')),
                        DropdownMenuItem(
                            value: 'qna', child: Text('Q&A')),
                      ],
                      onChanged: (v) =>
                          setState(() => _category = v ?? 'general'),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _titleController,
                      decoration: const InputDecoration(
                        labelText: '제목',
                        hintText: '게시글 제목을 입력하세요',
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _contentController,
                      maxLines: 12,
                      decoration: const InputDecoration(
                        labelText: '내용',
                        hintText: '내용을 작성하세요...',
                        alignLabelWithHint: true,
                      ),
                    ),
                    const SizedBox(height: 16),
                    SwitchListTile(
                      title: const Text('익명으로 작성'),
                      subtitle: const Text(
                        '닉네임이 "익명"으로 표시됩니다.',
                        style: TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                      value: _isAnonymous,
                      onChanged: (v) =>
                          setState(() => _isAnonymous = v),
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
}
