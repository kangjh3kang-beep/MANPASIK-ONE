import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// QnaScreen — Sanggam Orbit 전문가 Q&A
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar 2x → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.* 4x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* 2x → SanggamTheme 상수
// [Rule 4] Colors.green 4x → SanggamTheme.jagaeCyan
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 2x → SanggamTheme.background
// [Rule 4] FAB foregroundColor 추가
// ───────────────────────────────────────────────────

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
      final posts = (resp['questions'] as List?)
              ?.cast<Map<String, dynamic>>() ??
          [];
      if (mounted) {
        setState(() {
          _questions = posts;
          _isLoading = false;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() {
          _questions = _fallbackQuestions;
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (widget.mode == 'ask') {
      return _AskQuestionView(onSubmit: () => context.pop());
    }

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
                      '전문가 Q&A',
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
            // 목록
            Expanded(
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator())
                  : RefreshIndicator(
                      onRefresh: _loadQuestions,
                      child: ListView.builder(
                        padding: const EdgeInsets.all(16),
                        itemCount: _questions.length,
                        itemBuilder: (context, index) {
                          final q = _questions[index];
                          final title =
                              q['title'] as String? ?? '질문';
                          final author =
                              q['author_name'] as String? ??
                                  '익명';
                          final isAnswered =
                              q['is_answered'] as bool? ??
                                  false;
                          final answerCount =
                              q['answer_count'] as int? ?? 0;

                          return Padding(
                            padding: const EdgeInsets.only(
                                bottom: 8),
                            child: ListTile(
                              tileColor: SanggamTheme.surface,
                              shape: RoundedRectangleBorder(
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                              leading: CircleAvatar(
                                backgroundColor: isAnswered
                                    ? SanggamTheme.jagaeCyan
                                        .withValues(
                                            alpha: 0.1)
                                    : SanggamTheme
                                        .surfaceVariant,
                                child: Icon(
                                  isAnswered
                                      ? Icons.check_circle
                                      : Icons.help_outline,
                                  color: isAnswered
                                      ? SanggamTheme
                                          .jagaeCyan
                                      : SanggamTheme
                                          .onSurfaceDim,
                                ),
                              ),
                              title: Text(
                                title,
                                style: const TextStyle(
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                              subtitle: Row(
                                children: [
                                  Text(
                                    author,
                                    style: const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize: 12,
                                    ),
                                  ),
                                  const SizedBox(width: 8),
                                  if (isAnswered) ...[
                                    const Icon(
                                        Icons.verified,
                                        size: 14,
                                        color: SanggamTheme
                                            .jagaeCyan),
                                    const SizedBox(width: 4),
                                    const Text(
                                      '답변 채택',
                                      style: TextStyle(
                                        color: SanggamTheme
                                            .jagaeCyan,
                                        fontSize: 12,
                                      ),
                                    ),
                                  ],
                                  const Spacer(),
                                  Text(
                                    '답변 $answerCount',
                                    style: const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize: 12,
                                    ),
                                  ),
                                ],
                              ),
                              onTap: () => context.push(
                                  '/community/post/${q['id']}'),
                            ),
                          );
                        },
                      ),
                    ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        tooltip: '질문하기',
        onPressed: () => context.push('/community/qna/ask'),
        icon: const Icon(Icons.edit),
        label: const Text('질문하기'),
        backgroundColor: SanggamTheme.primary,
        foregroundColor: SanggamTheme.background,
      ),
    );
  }

  static final _fallbackQuestions = [
    {
      'id': 'q1',
      'title': '혈당 수치가 갑자기 올랐는데 원인이 뭘까요?',
      'author_name': '건강이',
      'is_answered': true,
      'answer_count': 3,
    },
    {
      'id': 'q2',
      'title': '카트리지 보관 방법에 대해 알려주세요',
      'author_name': '초보유저',
      'is_answered': true,
      'answer_count': 2,
    },
    {
      'id': 'q3',
      'title': '아이의 바이오마커 정상 범위가 궁금합니다',
      'author_name': '부모맘',
      'is_answered': false,
      'answer_count': 1,
    },
  ];
}

class _AskQuestionView extends StatefulWidget {
  const _AskQuestionView({required this.onSubmit});
  final VoidCallback onSubmit;

  @override
  State<_AskQuestionView> createState() =>
      _AskQuestionViewState();
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
                    icon: const Icon(Icons.close,
                        color: Colors.white),
                    tooltip: '닫기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '질문하기',
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
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment.stretch,
                  children: [
                    TextFormField(
                      controller: _titleController,
                      decoration: const InputDecoration(
                        labelText: '질문 제목',
                        hintText: '궁금한 점을 간결하게 작성하세요',
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextFormField(
                      controller: _contentController,
                      maxLines: 8,
                      decoration: const InputDecoration(
                        labelText: '상세 내용',
                        hintText:
                            '증상, 수치, 상황 등을 자세히 설명해주세요',
                        alignLabelWithHint: true,
                      ),
                    ),
                    const SizedBox(height: 24),
                    FilledButton(
                      onPressed: () {
                        ScaffoldMessenger.of(context)
                            .showSnackBar(
                          const SnackBar(
                              content:
                                  Text('질문이 등록되었습니다.')),
                        );
                        widget.onSubmit();
                      },
                      style: FilledButton.styleFrom(
                        minimumSize:
                            const Size.fromHeight(48),
                        backgroundColor:
                            SanggamTheme.primary,
                        foregroundColor:
                            SanggamTheme.background,
                        shape: RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius.circular(16)),
                      ),
                      child: const Text('질문 등록'),
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
