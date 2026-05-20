import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/providers/chat_provider.dart';
import 'package:manpasik/shared/widgets/streaming_text_bubble.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';

// ───────────────────────────────────────────────────
// ChatScreen — Sanggam Orbit AI 채팅
//
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] AppTheme.celadonTeal → SanggamTheme.jagaeCyan (alpha)
// [Rule 4] AppTheme.deepSeaBlue → SanggamTheme.background
// [Rule 4] AppTheme.aiPersonaAsset → 직접 경로 상수
// [Rule 4] jagae_pattern.dart + KoreanEdgeBorder → Container border
// [Rule 4] Theme.of(context) ~6x 제거
// [Rule 4] theme.colorScheme.* ~15x → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~8x → 직접 TextStyle
// [Rule 4] withOpacity ~10x → withValues(alpha:)
// [Rule 4] Color(0xFF1A1A1A) → SanggamTheme.surface
// [Rule 2] 12→16, 4→8, 10→8, 20→16, 28→32, 44→48, borderRadius 20→24
// ───────────────────────────────────────────────────

const _aiPersonaAsset =
    'assets/images/premium/premium_ai_persona_hologram_orb_1771750122274.png';

/// AI 건강 어시스턴트 채팅 화면
///
/// gRPC AIInferenceService 연동, 미연결 시 로컬 fallback 응답 제공.
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
        content:
            const Text('모든 대화 기록이 삭제됩니다. 계속하시겠습니까?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('취소'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text(
              '삭제',
              style: TextStyle(color: SanggamTheme.error),
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
    final chatState = ref.watch(chatProvider);

    ref.listen<ChatState>(chatProvider, (prev, next) {
      if (prev != null &&
          next.messages.length > prev.messages.length) {
        _scrollToBottom();
      }
    });

    return Scaffold(
      backgroundColor: Colors.transparent,
      body: CosmicBackground(
        child: SafeArea(
          child: Column(
            children: [
              // 헤더: 뒤로가기 + AI 코치 타이틀 + 삭제
              Padding(
                padding: const EdgeInsets.symmetric(
                    horizontal: 8, vertical: 8),
                child: Row(
                  children: [
                    const BackButton(color: Colors.white),
                    Image.asset(_aiPersonaAsset,
                        height: 32, width: 32, fit: BoxFit.contain),
                    const SizedBox(width: 8),
                    const Expanded(
                      child: Text(
                        '만파식 AI 코치',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                    if (chatState.messages.isNotEmpty)
                      IconButton(
                        icon: const Icon(
                            Icons.delete_outline_rounded,
                            color: Colors.white),
                        tooltip: '대화 기록 삭제',
                        onPressed: _showClearDialog,
                      ),
                  ],
                ),
              ),

              // 면책 조항
              const _DisclaimerBanner(
                text:
                    'AI 건강 코치의 응답은 참고용이며, 의료 전문가의 진단을 대체하지 않습니다.',
              ),

              // 메시지 영역
              Expanded(
                child: chatState.messages.isEmpty &&
                        !chatState.isStreaming
                    ? _EmptyState(onExampleTap: _sendExample)
                    : _MessageList(
                        messages: chatState.messages,
                        isLoading: chatState.isLoading,
                        isStreaming: chatState.isStreaming,
                        streamingContent:
                            chatState.streamingContent,
                        scrollController: _scrollController,
                        typingText: 'AI가 분석 중...',
                      ),
              ),

              // 입력 영역
              _ChatInputBar(
                controller: _controller,
                focusNode: _focusNode,
                hintText: '건강에 대해 물어보세요...',
                sendLabel: '전송',
                isLoading:
                    chatState.isLoading || chatState.isStreaming,
                onSend: _sendMessage,
              ),
            ],
          ),
        ),
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
    return Container(
      width: double.infinity,
      padding:
          const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: SanggamTheme.surfaceVariant.withValues(alpha: 0.5),
      child: Text(
        text,
        textAlign: TextAlign.center,
        style: const TextStyle(
          color: SanggamTheme.onSurfaceDim,
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
            width: 96,
            height: 96,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              boxShadow: [
                BoxShadow(
                  color:
                      SanggamTheme.jagaeCyan.withValues(alpha: 0.3),
                  blurRadius: 32,
                  spreadRadius: 8,
                ),
              ],
            ),
            child: Image.asset(_aiPersonaAsset,
                fit: BoxFit.contain),
          ),

          const SizedBox(height: 24),

          // 환영 메시지
          const Text(
            '안녕하세요! MANPASIK AI 건강 코치입니다.\n건강에 관한 질문을 해주세요.',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
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
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
                onPressed: () => onExampleTap(e.$2),
                backgroundColor: SanggamTheme.surfaceVariant
                    .withValues(alpha: 0.6),
                side: BorderSide(
                  color: SanggamTheme.onSurfaceDim
                      .withValues(alpha: 0.3),
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(24),
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
      padding:
          const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
      itemCount: messages.length + extraItems,
      itemBuilder: (context, index) {
        // 스트리밍 버블
        if (index == messages.length && isStreaming) {
          return Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const CircleAvatar(
                  radius: 16,
                  backgroundColor: SanggamTheme.primary,
                  child: Icon(Icons.smart_toy_rounded,
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
  ConsumerState<_MessageBubble> createState() =>
      _MessageBubbleState();
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
          _translatedText = res['translated_text'] as String? ??
              res['text'] as String? ??
              '';
          _translating = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _translating = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isUser = widget.message.role == 'user';
    final bubbleRadius = BorderRadius.only(
      topLeft: const Radius.circular(16),
      topRight: const Radius.circular(16),
      bottomLeft: Radius.circular(isUser ? 16 : 4),
      bottomRight: Radius.circular(isUser ? 4 : 16),
    );

    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Row(
        mainAxisAlignment:
            isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (!isUser) ...[
            const CircleAvatar(
              radius: 16,
              backgroundColor: SanggamTheme.primary,
              child: Icon(Icons.smart_toy_rounded,
                  color: Colors.white, size: 18),
            ),
            const SizedBox(width: 8),
          ],
          Flexible(
            child: Column(
              crossAxisAlignment: isUser
                  ? CrossAxisAlignment.end
                  : CrossAxisAlignment.start,
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 16, vertical: 16),
                  decoration: BoxDecoration(
                    color: isUser
                        ? SanggamTheme.jagaeCyan
                            .withValues(alpha: 0.2)
                        : SanggamTheme.surface
                            .withValues(alpha: 0.8),
                    borderRadius: bubbleRadius,
                    border: isUser
                        ? null
                        : Border.all(
                            color: SanggamTheme.primary
                                .withValues(alpha: 0.5),
                            width: 1.0,
                          ),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        widget.message.content,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          height: 1.5,
                        ),
                      ),
                      if (_translatedText != null) ...[
                        Divider(
                            height: 16,
                            color: SanggamTheme.onSurfaceDim
                                .withValues(alpha: 0.3)),
                        Text(
                          _translatedText!,
                          style: const TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 12,
                            fontStyle: FontStyle.italic,
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
                const SizedBox(height: 8),
                Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      _formatTime(widget.message.timestamp),
                      style: TextStyle(
                        color: SanggamTheme.onSurfaceDim
                            .withValues(alpha: 0.6),
                        fontSize: 11,
                      ),
                    ),
                    const SizedBox(width: 8),
                    InkWell(
                      onTap:
                          _translating ? null : _translateMessage,
                      child: _translating
                          ? const SizedBox(
                              width: 12,
                              height: 12,
                              child: CircularProgressIndicator(
                                  strokeWidth: 1.5))
                          : Icon(
                              Icons.translate,
                              size: 14,
                              color: _translatedText != null
                                  ? SanggamTheme.primary
                                  : SanggamTheme.onSurfaceDim
                                      .withValues(alpha: 0.6),
                            ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          if (isUser) ...[
            const SizedBox(width: 8),
            const CircleAvatar(
              radius: 16,
              backgroundColor: SanggamTheme.background,
              child: Icon(Icons.person_rounded,
                  color: Colors.white, size: 18),
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
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: 18,
            backgroundColor: Colors.transparent,
            child:
                Image.asset(_aiPersonaAsset, fit: BoxFit.contain),
          ),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(
                horizontal: 16, vertical: 16),
            decoration: BoxDecoration(
              color: SanggamTheme.surfaceVariant,
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
                const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: SanggamTheme.primary,
                  ),
                ),
                const SizedBox(width: 8),
                const Text(
                  'AI가 분석 중...',
                  style: TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
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
    return Container(
      decoration: BoxDecoration(
        color: SanggamTheme.surface,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.3),
            blurRadius: 8,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        top: false,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(16, 8, 8, 8),
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
                    fillColor: SanggamTheme.surfaceVariant
                        .withValues(alpha: 0.5),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(24),
                      borderSide: BorderSide.none,
                    ),
                    contentPadding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Material(
                color: SanggamTheme.primary,
                shape: const CircleBorder(),
                clipBehavior: Clip.antiAlias,
                child: InkWell(
                  onTap: isLoading ? null : onSend,
                  child: SizedBox(
                    width: 48,
                    height: 48,
                    child: Icon(
                      Icons.send_rounded,
                      color: isLoading
                          ? Colors.white.withValues(alpha: 0.5)
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
