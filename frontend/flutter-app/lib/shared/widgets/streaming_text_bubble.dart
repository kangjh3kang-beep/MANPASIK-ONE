import 'package:flutter/material.dart';

/// AI 스트리밍 응답을 실시간으로 표시하는 텍스트 버블 (C1)
///
/// 텍스트가 점진적으로 표시되며, 커서 깜빡임 애니메이션으로
/// 실시간 응답 생성 중임을 시각적으로 표현합니다.
class StreamingTextBubble extends StatefulWidget {
  const StreamingTextBubble({
    super.key,
    required this.text,
    required this.isStreaming,
    this.style,
  });

  final String text;
  final bool isStreaming;
  final TextStyle? style;

  @override
  State<StreamingTextBubble> createState() => _StreamingTextBubbleState();
}

class _StreamingTextBubbleState extends State<StreamingTextBubble>
    with SingleTickerProviderStateMixin {
  late AnimationController _cursorController;

  @override
  void initState() {
    super.initState();
    _cursorController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _cursorController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final textStyle = widget.style ??
        theme.textTheme.bodyMedium?.copyWith(
          color: theme.colorScheme.onSurface,
          height: 1.5,
        );

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest,
        borderRadius: const BorderRadius.only(
          topLeft: Radius.circular(16),
          topRight: Radius.circular(16),
          bottomLeft: Radius.circular(4),
          bottomRight: Radius.circular(16),
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.08),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: widget.text.isEmpty && widget.isStreaming
          ? _buildThinkingDots(theme)
          : RichText(
              text: TextSpan(
                children: [
                  TextSpan(text: widget.text, style: textStyle),
                  if (widget.isStreaming)
                    WidgetSpan(
                      child: AnimatedBuilder(
                        animation: _cursorController,
                        builder: (context, child) => Opacity(
                          opacity: _cursorController.value,
                          child: Text(
                            '\u258C',
                            style: textStyle?.copyWith(
                              color: theme.colorScheme.primary,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
    );
  }

  Widget _buildThinkingDots(ThemeData theme) {
    return Row(
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
          'AI가 생각 중...',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
      ],
    );
  }
}
