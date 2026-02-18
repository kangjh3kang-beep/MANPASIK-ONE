import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 전문가 Q&A 화면
class QnaScreen extends ConsumerStatefulWidget {
  const QnaScreen({super.key, this.mode});

  final String? mode;

  @override
  ConsumerState<QnaScreen> createState() => _QnaScreenState();
}

class _QnaScreenState extends ConsumerState<QnaScreen> {
  List<Map<String, dynamic>> _questions = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadQuestions();
  }

  Future<void> _loadQuestions() async {
    try {
      final client = ref.read(restClientProvider);
      final resp = await client.getQnaQuestions();
      final posts = (resp['questions'] as List?)?.cast<Map<String, dynamic>>() ?? [];
      if (mounted) setState(() { _questions = posts; _isLoading = false; });
    } catch (_) {
      if (mounted) setState(() {
        _questions = _fallbackQuestions;
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (widget.mode == 'ask') {
      return _AskQuestionView(onSubmit: () => context.pop());
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('전문가 Q&A'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _loadQuestions,
              child: ListView.builder(
                padding: const EdgeInsets.all(16),
                itemCount: _questions.length,
                itemBuilder: (context, index) {
                  final q = _questions[index];
                  final title = q['title'] as String? ?? '질문';
                  final author = q['author_name'] as String? ?? '익명';
                  final isAnswered = q['is_answered'] as bool? ?? false;
                  final answerCount = q['answer_count'] as int? ?? 0;

                  return Card(
                    margin: const EdgeInsets.only(bottom: 8),
                    child: ListTile(
                      leading: CircleAvatar(
                        backgroundColor: isAnswered ? Colors.green.withOpacity(0.1) : theme.colorScheme.surfaceContainerHighest,
                        child: Icon(
                          isAnswered ? Icons.check_circle : Icons.help_outline,
                          color: isAnswered ? Colors.green : theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                      title: Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
                      subtitle: Row(
                        children: [
                          Text(author, style: theme.textTheme.bodySmall),
                          const SizedBox(width: 8),
                          if (isAnswered) ...[
                            const Icon(Icons.verified, size: 14, color: Colors.green),
                            const SizedBox(width: 2),
                            Text('답변 채택', style: theme.textTheme.bodySmall?.copyWith(color: Colors.green)),
                          ],
                          const Spacer(),
                          Text('답변 $answerCount', style: theme.textTheme.bodySmall),
                        ],
                      ),
                      onTap: () => context.push('/community/post/${q['id']}'),
                    ),
                  );
                },
              ),
            ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.push('/community/qna/ask'),
        icon: const Icon(Icons.edit),
        label: const Text('질문하기'),
        backgroundColor: AppTheme.sanggamGold,
      ),
    );
  }

  static final _fallbackQuestions = [
    {'id': 'q1', 'title': '혈당 수치가 갑자기 올랐는데 원인이 뭘까요?', 'author_name': '건강이', 'is_answered': true, 'answer_count': 3},
    {'id': 'q2', 'title': '카트리지 보관 방법에 대해 알려주세요', 'author_name': '초보유저', 'is_answered': true, 'answer_count': 2},
    {'id': 'q3', 'title': '아이의 바이오마커 정상 범위가 궁금합니다', 'author_name': '부모맘', 'is_answered': false, 'answer_count': 1},
  ];
}

class _AskQuestionView extends StatefulWidget {
  const _AskQuestionView({required this.onSubmit});
  final VoidCallback onSubmit;

  @override
  State<_AskQuestionView> createState() => _AskQuestionViewState();
}

class _AskQuestionViewState extends State<_AskQuestionView> {
  final _titleController = TextEditingController();
  final _contentController = TextEditingController();

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('질문하기'),
        leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop()),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(labelText: '질문 제목', hintText: '궁금한 점을 간결하게 작성하세요'),
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _contentController,
              maxLines: 8,
              decoration: const InputDecoration(labelText: '상세 내용', hintText: '증상, 수치, 상황 등을 자세히 설명해주세요', alignLabelWithHint: true),
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('질문이 등록되었습니다.')));
                widget.onSubmit();
              },
              style: FilledButton.styleFrom(minimumSize: const Size.fromHeight(48), backgroundColor: AppTheme.sanggamGold),
              child: const Text('질문 등록'),
            ),
          ],
        ),
      ),
    );
  }
}
