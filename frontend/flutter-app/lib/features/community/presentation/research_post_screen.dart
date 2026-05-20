import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ResearchPostScreen — Sanggam Orbit 연구 게시글 작성
//
// [Rule 4] app_theme.dart(미사용) → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.* 2x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* 3x → SanggamTheme 상수
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] padding:12→16
// ───────────────────────────────────────────────────

/// 연구 협업 게시글 화면 (C7)
class ResearchPostScreen extends ConsumerStatefulWidget {
  const ResearchPostScreen({super.key});

  @override
  ConsumerState<ResearchPostScreen> createState() =>
      _ResearchPostScreenState();
}

class _ResearchPostScreenState
    extends ConsumerState<ResearchPostScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _contentController = TextEditingController();
  String _selectedType = 'participation';
  bool _isSubmitting = false;

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  Future<void> _submitPost() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isSubmitting = true);

    try {
      final userId = ref.read(authProvider).userId ?? '';
      await ref.read(restClientProvider).createPost(
            authorId: userId,
            title: _titleController.text.trim(),
            content: _contentController.text.trim(),
            category: 5, // 5 = research category
            tags: [_selectedType, 'research'],
          );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('연구 게시글이 등록되었습니다')),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('등록 실패: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '연구 게시글 작성',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(24),
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment:
                        CrossAxisAlignment.stretch,
                    children: [
                      // 연구 유형 선택
                      const Text(
                        '연구 유형',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      SegmentedButton<String>(
                        segments: const [
                          ButtonSegment(
                            value: 'participation',
                            label: Text('참여 연구'),
                            icon: Icon(Icons.people_outline,
                                size: 18),
                          ),
                          ButtonSegment(
                            value: 'data_sharing',
                            label: Text('데이터 공유'),
                            icon: Icon(Icons.share_outlined,
                                size: 18),
                          ),
                          ButtonSegment(
                            value: 'discussion',
                            label: Text('연구 토론'),
                            icon: Icon(Icons.forum_outlined,
                                size: 18),
                          ),
                        ],
                        selected: {_selectedType},
                        onSelectionChanged: (s) => setState(
                            () => _selectedType = s.first),
                      ),
                      const SizedBox(height: 24),

                      // 제목
                      TextFormField(
                        controller: _titleController,
                        decoration: const InputDecoration(
                          labelText: '제목',
                          hintText: '연구 주제를 입력하세요',
                          border: OutlineInputBorder(),
                        ),
                        validator: (v) =>
                            (v == null || v.trim().isEmpty)
                                ? '제목을 입력해주세요'
                                : null,
                      ),
                      const SizedBox(height: 16),

                      // 본문
                      TextFormField(
                        controller: _contentController,
                        maxLines: 8,
                        decoration: const InputDecoration(
                          labelText: '내용',
                          hintText:
                              '연구 목적, 방법, 참여 조건 등을 상세히 기술해주세요',
                          border: OutlineInputBorder(),
                          alignLabelWithHint: true,
                        ),
                        validator: (v) =>
                            (v == null || v.trim().isEmpty)
                                ? '내용을 입력해주세요'
                                : null,
                      ),
                      const SizedBox(height: 16),

                      // IRB 안내
                      SanggamContainer(
                        borderRadius: 16,
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          children: [
                            const Icon(Icons.info_outline,
                                size: 18,
                                color:
                                    SanggamTheme.onSurfaceDim),
                            const SizedBox(width: 8),
                            const Expanded(
                              child: Text(
                                '의료 데이터를 활용한 연구는 IRB 승인이 필요할 수 있습니다.',
                                style: TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim,
                                  fontSize: 12,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 24),

                      FilledButton.icon(
                        onPressed:
                            _isSubmitting ? null : _submitPost,
                        icon: _isSubmitting
                            ? const SizedBox(
                                width: 18,
                                height: 18,
                                child:
                                    CircularProgressIndicator(
                                        strokeWidth: 2,
                                        color: Colors.white),
                              )
                            : const Icon(
                                Icons.publish_rounded),
                        label: const Text('게시하기'),
                        style: FilledButton.styleFrom(
                          minimumSize:
                              const Size(double.infinity, 48),
                          backgroundColor:
                              SanggamTheme.primary,
                          foregroundColor:
                              SanggamTheme.background,
                          shape: RoundedRectangleBorder(
                              borderRadius:
                                  BorderRadius.circular(16)),
                        ),
                      ),
                    ],
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
