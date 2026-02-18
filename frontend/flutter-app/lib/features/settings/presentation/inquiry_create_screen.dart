import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 1:1 문의 작성 화면
class InquiryCreateScreen extends ConsumerStatefulWidget {
  const InquiryCreateScreen({super.key});

  @override
  ConsumerState<InquiryCreateScreen> createState() => _InquiryCreateScreenState();
}

class _InquiryCreateScreenState extends ConsumerState<InquiryCreateScreen> {
  String _type = 'device';
  final _titleController = TextEditingController();
  final _contentController = TextEditingController();
  bool _notifyByPush = true;
  bool _notifyByEmail = false;
  bool _isSubmitting = false;

  @override
  void dispose() {
    _titleController.dispose();
    _contentController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('1:1 문의'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            DropdownButtonFormField<String>(
              value: _type,
              decoration: const InputDecoration(labelText: '문의 유형'),
              items: const [
                DropdownMenuItem(value: 'device', child: Text('기기/카트리지')),
                DropdownMenuItem(value: 'subscription', child: Text('구독/결제')),
                DropdownMenuItem(value: 'account', child: Text('계정/인증')),
                DropdownMenuItem(value: 'measurement', child: Text('측정/결과')),
                DropdownMenuItem(value: 'bug', child: Text('오류 신고')),
                DropdownMenuItem(value: 'other', child: Text('기타')),
              ],
              onChanged: (v) => setState(() => _type = v ?? 'device'),
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _titleController,
              decoration: const InputDecoration(labelText: '제목', hintText: '문의 제목을 입력해주세요'),
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _contentController,
              maxLines: 8,
              decoration: const InputDecoration(
                labelText: '문의 내용',
                hintText: '문의하실 내용을 상세히 작성해주세요.',
                alignLabelWithHint: true,
              ),
            ),
            const SizedBox(height: 16),

            // 알림 설정
            Text('답변 알림', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            SwitchListTile(
              title: const Text('푸시 알림'),
              value: _notifyByPush,
              onChanged: (v) => setState(() => _notifyByPush = v),
              dense: true,
            ),
            SwitchListTile(
              title: const Text('이메일 알림'),
              value: _notifyByEmail,
              onChanged: (v) => setState(() => _notifyByEmail = v),
              dense: true,
            ),
            const SizedBox(height: 24),
            FilledButton(
              onPressed: _isSubmitting ? null : () async {
                if (_titleController.text.trim().isEmpty || _contentController.text.trim().isEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('제목과 내용을 입력해주세요.')));
                  return;
                }
                setState(() => _isSubmitting = true);
                await Future.delayed(const Duration(seconds: 1));
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('문의가 접수되었습니다.')));
                  context.pop();
                }
              },
              style: FilledButton.styleFrom(minimumSize: const Size.fromHeight(48), backgroundColor: AppTheme.sanggamGold),
              child: _isSubmitting
                  ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                  : const Text('문의 접수'),
            ),
          ],
        ),
      ),
    );
  }
}
