import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/chat_provider.dart';
import 'package:manpasik/shared/widgets/streaming_text_bubble.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// AI 건강 어시스턴트 채팅 화면
///
/// gRPC AIInferenceService 연동, 미연결 시 로컬 fallback 응답 제공.
/// 빈 상태: 환영 메시지 + 예시 질문 칩, 메시지 목록: 사용자/AI 구분.
class ChatScreen extends ConsumerStatefulWidget {
  const ChatScreen({super.key});

  @override
  ConsumerState<ChatScreen> createState() => _ChatScreenState();
}

class _ChatScreenState extends ConsumerState<ChatScreen> {
  final _controller = TextEditingController();
  final _scrollController = ScrollController();
  final _focusNode = FocusNode();

  @override
  void dispose() {
    _controller.dispose();
    _scrollController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  void _sendMessage() {
    final text = _controller.text.trim();
    if (text.isEmpty) return;
    _controller.clear();
    ref.read(chatProvider.notifier).sendMessageStream(text);
    _scrollToBottom();
  }

  void _sendExample(String text) {
    ref.read(chatProvider.notifier).sendMessageStream(text);
    _scrollToBottom();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent + 100,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  void _showClearDialog() {
    showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('대화 기록 삭제'),
        content: const Text('모든 대화 기록이 삭제됩니다. 계속하시겠습니까?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('취소'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: Text(
              '삭제',
              style: TextStyle(color: Theme.of(context).colorScheme.error),
            ),
          ),
        ],
      ),
    ).then((confirmed) {
      if (confirmed == true) {
        ref.read(chatProvider.notifier).clearChat();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chatState = ref.watch(chatProvider);

    // 새 메시지 도착 시 자동 스크롤
    ref.listen<ChatState>(chatProvider, (prev, next) {
      if (prev != null && next.messages.length > prev.messages.length) {
        _scrollToBottom();
      }
    });

    return Scaffold(
      appBar: AppBar(
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.smart_toy_rounded,
              color: theme.colorScheme.primary,
              size: 24,
            ),
            const SizedBox(width: 8),
            const Text('AI 건강 코치'),
          ],
        ),
        centerTitle: true,
        actions: [
          if (chatState.messages.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.delete_outline_rounded),
              tooltip: '대화 기록 삭제',
              onPressed: _showClearDialog,
            ),
        ],
      ),
      body: Column(
        children: [
          // 면책 조항
          const _DisclaimerBanner(
            text: 'AI 건강 코치의 응답은 참고용이며, 의료 전문가의 진단을 대체하지 않습니다.',
          ),

          // 메시지 영역
          Expanded(
            child: chatState.messages.isEmpty && !chatState.isStreaming
                ? _EmptyState(onExampleTap: _sendExample)
                : _MessageList(
                    messages: chatState.messages,
                    isLoading: chatState.isLoading,
                    isStreaming: chatState.isStreaming,
                    streamingContent: chatState.streamingContent,
                    scrollController: _scrollController,
                    typingText: 'AI가 응답 중...',
                  ),
          ),

          // 입력 영역
          _ChatInputBar(
            controller: _controller,
            focusNode: _focusNode,
            hintText: '건강에 대해 물어보세요...',
            sendLabel: '전송',
            isLoading: chatState.isLoading || chatState.isStreaming,
            onSend: _sendMessage,
          ),
        ],
      ),
    );
  }
}

// ─── 면책 조항 배너 ───

class _DisclaimerBanner extends StatelessWidget {
  const _DisclaimerBanner({required this.text});
  final String text;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: theme.colorScheme.surfaceContainerHighest.withOpacity(0.5),
      child: Text(
        text,
        textAlign: TextAlign.center,
        style: theme.textTheme.bodySmall?.copyWith(
          color: theme.colorScheme.onSurfaceVariant,
          fontSize: 11,
        ),
      ),
    );
  }
}

// ─── 빈 상태: 환영 메시지 + 예시 질문 ───

class _EmptyState extends StatelessWidget {
  const _EmptyState({required this.onExampleTap});
  final void Function(String) onExampleTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    const examples = [
      (Icons.bloodtype_rounded, '혈당 수치가 높은데 어떻게 해야 하나요?'),
      (Icons.favorite_rounded, '혈압 관리 방법을 알려주세요'),
      (Icons.fitness_center_rounded, '당뇨 환자에게 좋은 운동은?'),
      (Icons.restaurant_rounded, '건강한 식단 추천해주세요'),
      (Icons.bedtime_rounded, '수면 질 개선 방법이 궁금해요'),
    ];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          const SizedBox(height: 32),

          // AI 아이콘
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              gradient: LinearGradient(
                colors: [
                  theme.colorScheme.primary,
                  theme.colorScheme.tertiary,
                ],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color: theme.colorScheme.primary.withOpacity(0.3),
                  blurRadius: 20,
                  offset: const Offset(0, 8),
                ),
              ],
            ),
            child: const Icon(
              Icons.smart_toy_rounded,
              color: Colors.white,
              size: 40,
            ),
          ),

          const SizedBox(height: 24),

          // 환영 메시지
          Text(
            '안녕하세요! MANPASIK AI 건강 코치입니다.\n건강에 관한 질문을 해주세요.',
            textAlign: TextAlign.center,
            style: theme.textTheme.bodyLarge?.copyWith(
              color: theme.colorScheme.onSurface,
              height: 1.6,
            ),
          ),

          const SizedBox(height: 32),

          // 예시 질문 칩
          Wrap(
            spacing: 8,
            runSpacing: 8,
            alignment: WrapAlignment.center,
            children: examples.map((e) {
              return ActionChip(
                avatar: Icon(e.$1, size: 18),
                label: Text(
                  e.$2,
                  style: theme.textTheme.bodySmall,
                ),
                onPressed: () => onExampleTap(e.$2),
                backgroundColor:
                    theme.colorScheme.surfaceContainerHighest.withOpacity(0.6),
                side: BorderSide(
                  color: theme.colorScheme.outlineVariant.withOpacity(0.5),
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(20),
                ),
              );
            }).toList(),
          ),
        ],
      ),
    );
  }
}

// ─── 메시지 목록 ───

class _MessageList extends StatelessWidget {
  const _MessageList({
    required this.messages,
    required this.isLoading,
    required this.scrollController,
    required this.typingText,
    this.isStreaming = false,
    this.streamingContent = '',
  });
  final List<ChatMessage> messages;
  final bool isLoading;
  final bool isStreaming;
  final String streamingContent;
  final ScrollController scrollController;
  final String typingText;

  @override
  Widget build(BuildContext context) {
    final extraItems = isStreaming ? 1 : (isLoading ? 1 : 0);
    return ListView.builder(
      controller: scrollController,
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      itemCount: messages.length + extraItems,
      itemBuilder: (context, index) {
        // 스트리밍 버블
        if (index == messages.length && isStreaming) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                CircleAvatar(
                  radius: 16,
                  backgroundColor: Theme.of(context).colorScheme.primary,
                  child: const Icon(Icons.smart_toy_rounded,
                      color: Colors.white, size: 18),
                ),
                const SizedBox(width: 8),
                Flexible(
                  child: StreamingTextBubble(
                    text: streamingContent,
                    isStreaming: true,
                  ),
                ),
              ],
            ),
          );
        }
        // 타이핑 인디케이터
        if (index == messages.length && isLoading) {
          return _TypingIndicator(text: typingText);
        }
        return _MessageBubble(message: messages[index]);
      },
    );
  }
}

// ─── 메시지 버블 ───

class _MessageBubble extends ConsumerStatefulWidget {
  const _MessageBubble({required this.message});
  final ChatMessage message;

  @override
  ConsumerState<_MessageBubble> createState() => _MessageBubbleState();
}

class _MessageBubbleState extends ConsumerState<_MessageBubble> {
  String? _translatedText;
  bool _translating = false;

  Future<void> _translateMessage() async {
    if (_translatedText != null) {
      setState(() => _translatedText = null);
      return;
    }
    setState(() => _translating = true);
    try {
      final client = ref.read(restClientProvider);
      final res = await client.translateText(
        text: widget.message.content,
        sourceLanguage: 'auto',
        targetLanguage: 'en',
      );
      if (mounted) {
        setState(() {
          _translatedText = res['translated_text'] as String? ?? res['text'] as String? ?? '';
          _translating = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _translating = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isUser = widget.message.role == 'user';

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        mainAxisAlignment:
            isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (!isUser) ...[
            CircleAvatar(
              radius: 16,
              backgroundColor: theme.colorScheme.primary,
              child: const Icon(Icons.smart_toy_rounded, color: Colors.white, size: 18),
            ),
            const SizedBox(width: 8),
          ],
          Flexible(
            child: Column(
              crossAxisAlignment:
                  isUser ? CrossAxisAlignment.end : CrossAxisAlignment.start,
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                  decoration: BoxDecoration(
                    color: isUser
                        ? theme.colorScheme.primary
                        : theme.colorScheme.surfaceContainerHighest,
                    borderRadius: BorderRadius.only(
                      topLeft: const Radius.circular(16),
                      topRight: const Radius.circular(16),
                      bottomLeft: Radius.circular(isUser ? 16 : 4),
                      bottomRight: Radius.circular(isUser ? 4 : 16),
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: (isUser ? theme.colorScheme.primary : Colors.black)
                            .withValues(alpha: 0.08),
                        blurRadius: 8,
                        offset: const Offset(0, 2),
                      ),
                    ],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        widget.message.content,
                        style: theme.textTheme.bodyMedium?.copyWith(
                          color: isUser ? Colors.white : theme.colorScheme.onSurface,
                          height: 1.5,
                        ),
                      ),
                      if (_translatedText != null) ...[
                        const Divider(height: 12),
                        Text(
                          _translatedText!,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: isUser ? Colors.white70 : theme.colorScheme.onSurfaceVariant,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
                const SizedBox(height: 4),
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      _formatTime(widget.message.timestamp),
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant.withValues(alpha: 0.6),
                        fontSize: 11,
                      ),
                    ),
                    const SizedBox(width: 8),
                    InkWell(
                      onTap: _translating ? null : _translateMessage,
                      child: _translating
                          ? const SizedBox(width: 12, height: 12, child: CircularProgressIndicator(strokeWidth: 1.5))
                          : Icon(
                              _translatedText != null ? Icons.translate : Icons.translate,
                              size: 14,
                              color: _translatedText != null
                                  ? theme.colorScheme.primary
                                  : theme.colorScheme.onSurfaceVariant.withValues(alpha: 0.6),
                            ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          if (isUser) ...[
            const SizedBox(width: 8),
            CircleAvatar(
              radius: 16,
              backgroundColor: AppTheme.deepSeaBlue,
              child: const Icon(Icons.person_rounded, color: Colors.white, size: 18),
            ),
          ],
        ],
      ),
    );
  }

  String _formatTime(DateTime time) {
    final h = time.hour.toString().padLeft(2, '0');
    final m = time.minute.toString().padLeft(2, '0');
    return '$h:$m';
  }
}

// ─── 타이핑 인디케이터 ───

class _TypingIndicator extends StatelessWidget {
  const _TypingIndicator({required this.text});
  final String text;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: 16,
            backgroundColor: theme.colorScheme.primary,
            child: const Icon(
              Icons.smart_toy_rounded,
              color: Colors.white,
              size: 18,
            ),
          ),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: theme.colorScheme.surfaceContainerHighest,
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(16),
                topRight: Radius.circular(16),
                bottomLeft: Radius.circular(4),
                bottomRight: Radius.circular(16),
              ),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: theme.colorScheme.primary,
                  ),
                ),
                const SizedBox(width: 10),
                Text(
                  text,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// ─── 입력 바 ───

class _ChatInputBar extends StatelessWidget {
  const _ChatInputBar({
    required this.controller,
    required this.focusNode,
    required this.hintText,
    required this.sendLabel,
    required this.isLoading,
    required this.onSend,
  });

  final TextEditingController controller;
  final FocusNode focusNode;
  final String hintText;
  final String sendLabel;
  final bool isLoading;
  final VoidCallback onSend;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        top: false,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(12, 8, 8, 8),
          child: Row(
            children: [
              Expanded(
                child: TextField(
                  controller: controller,
                  focusNode: focusNode,
                  textInputAction: TextInputAction.send,
                  onSubmitted: (_) => onSend(),
                  enabled: !isLoading,
                  maxLines: 4,
                  minLines: 1,
                  decoration: InputDecoration(
                    hintText: hintText,
                    filled: true,
                    fillColor: theme.colorScheme.surfaceContainerHighest
                        .withOpacity(0.5),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(24),
                      borderSide: BorderSide.none,
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 20,
                      vertical: 10,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 8),
              // 전송 버튼
              Material(
                color: theme.colorScheme.primary,
                shape: const CircleBorder(),
                clipBehavior: Clip.antiAlias,
                child: InkWell(
                  onTap: isLoading ? null : onSend,
                  child: SizedBox(
                    width: 44,
                    height: 44,
                    child: Icon(
                      Icons.send_rounded,
                      color: isLoading
                          ? Colors.white.withOpacity(0.5)
                          : Colors.white,
                      size: 20,
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
}
