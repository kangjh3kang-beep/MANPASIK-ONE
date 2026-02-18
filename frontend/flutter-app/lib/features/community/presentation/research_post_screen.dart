import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 연구 협업 게시글 화면 (C7)
///
/// 의료 연구 관련 게시글 목록 및 작성 기능을 제공합니다.
/// 연구 데이터 공유, IRB 승인 연구, 사용자 참여 연구 등을 지원합니다.
class ResearchPostScreen extends ConsumerStatefulWidget {
  const ResearchPostScreen({super.key});

  @override
  ConsumerState<ResearchPostScreen> createState() =>
      _ResearchPostScreenState();
}

class _ResearchPostScreenState extends ConsumerState<ResearchPostScreen> {
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
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('연구 게시글 작성'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // 연구 유형 선택
              Text('연구 유형',
                  style: theme.textTheme.titleSmall
                      ?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              SegmentedButton<String>(
                segments: const [
                  ButtonSegment(
                    value: 'participation',
                    label: Text('참여 연구'),
                    icon: Icon(Icons.people_outline, size: 18),
                  ),
                  ButtonSegment(
                    value: 'data_sharing',
                    label: Text('데이터 공유'),
                    icon: Icon(Icons.share_outlined, size: 18),
                  ),
                  ButtonSegment(
                    value: 'discussion',
                    label: Text('연구 토론'),
                    icon: Icon(Icons.forum_outlined, size: 18),
                  ),
                ],
                selected: {_selectedType},
                onSelectionChanged: (s) =>
                    setState(() => _selectedType = s.first),
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
                    (v == null || v.trim().isEmpty) ? '제목을 입력해주세요' : null,
              ),
              const SizedBox(height: 16),

              // 본문
              TextFormField(
                controller: _contentController,
                maxLines: 8,
                decoration: const InputDecoration(
                  labelText: '내용',
                  hintText: '연구 목적, 방법, 참여 조건 등을 상세히 기술해주세요',
                  border: OutlineInputBorder(),
                  alignLabelWithHint: true,
                ),
                validator: (v) =>
                    (v == null || v.trim().isEmpty) ? '내용을 입력해주세요' : null,
              ),
              const SizedBox(height: 16),

              // IRB 안내
              Card(
                color: theme.colorScheme.surfaceContainerHighest
                    .withOpacity(0.5),
                child: Padding(
                  padding: const EdgeInsets.all(12),
                  child: Row(
                    children: [
                      Icon(Icons.info_outline,
                          size: 18,
                          color: theme.colorScheme.onSurfaceVariant),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          '의료 데이터를 활용한 연구는 IRB 승인이 필요할 수 있습니다.',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              FilledButton.icon(
                onPressed: _isSubmitting ? null : _submitPost,
                icon: _isSubmitting
                    ? const SizedBox(
                        width: 18,
                        height: 18,
                        child: CircularProgressIndicator(
                            strokeWidth: 2, color: Colors.white))
                    : const Icon(Icons.publish_rounded),
                label: const Text('게시하기'),
                style: FilledButton.styleFrom(
                  minimumSize: const Size(double.infinity, 48),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
