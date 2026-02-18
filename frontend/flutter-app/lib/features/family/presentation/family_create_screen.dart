import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 가족 그룹 생성 화면
class FamilyCreateScreen extends ConsumerStatefulWidget {
  const FamilyCreateScreen({super.key, this.mode});

  final String? mode;

  @override
  ConsumerState<FamilyCreateScreen> createState() => _FamilyCreateScreenState();
}

class _FamilyCreateScreenState extends ConsumerState<FamilyCreateScreen> {
  final _nameController = TextEditingController();
  final _inviteController = TextEditingController();
  String _inviteMethod = 'link';
  bool _isSubmitting = false;

  bool get _isInviteMode => widget.mode == 'invite';

  @override
  void dispose() {
    _nameController.dispose();
    _inviteController.dispose();
    super.dispose();
  }

  Future<void> _createGroup() async {
    if (_nameController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('그룹 이름을 입력해주세요.')));
      return;
    }
    setState(() => _isSubmitting = true);
    try {
      final client = ref.read(restClientProvider);
      await client.createFamilyGroup(userId: 'current-user', name: _nameController.text.trim());
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('가족 그룹이 생성되었습니다.')));
        context.pop();
      }
    } catch (_) {
      if (mounted) {
        setState(() => _isSubmitting = false);
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('그룹 생성에 실패했습니다.')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_isInviteMode ? '가족 초대' : '가족 그룹 만들기'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            if (!_isInviteMode) ...[
              Icon(Icons.family_restroom, size: 64, color: AppTheme.sanggamGold),
              const SizedBox(height: 16),
              Text('가족 그룹을 만들어\n건강 데이터를 함께 관리하세요.',
                  textAlign: TextAlign.center, style: theme.textTheme.bodyLarge),
              const SizedBox(height: 32),
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(labelText: '그룹 이름', hintText: '예: 우리 가족', prefixIcon: Icon(Icons.group)),
              ),
              const SizedBox(height: 24),
            ],
            Text('초대 방법', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(value: 'link', label: Text('초대 링크'), icon: Icon(Icons.link)),
                ButtonSegment(value: 'email', label: Text('이메일'), icon: Icon(Icons.email)),
                ButtonSegment(value: 'qr', label: Text('QR 코드'), icon: Icon(Icons.qr_code)),
              ],
              selected: {_inviteMethod},
              onSelectionChanged: (s) => setState(() => _inviteMethod = s.first),
            ),
            const SizedBox(height: 16),
            if (_inviteMethod == 'email')
              TextFormField(
                controller: _inviteController,
                keyboardType: TextInputType.emailAddress,
                decoration: const InputDecoration(labelText: '초대할 이메일', hintText: 'family@email.com', prefixIcon: Icon(Icons.email)),
              ),
            if (_inviteMethod == 'link')
              Card(
                child: ListTile(
                  leading: const Icon(Icons.link, color: AppTheme.sanggamGold),
                  title: const Text('초대 링크 생성'),
                  subtitle: const Text('링크를 복사하여 가족에게 공유하세요.'),
                  trailing: IconButton(icon: const Icon(Icons.copy), onPressed: () {
                    ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('초대 링크가 복사되었습니다.')));
                  }),
                ),
              ),
            if (_inviteMethod == 'qr')
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(32),
                  child: Column(
                    children: [
                      Icon(Icons.qr_code_2, size: 120, color: theme.colorScheme.onSurfaceVariant),
                      const SizedBox(height: 8),
                      Text('QR 코드를 스캔하여 참여하세요', style: theme.textTheme.bodySmall),
                    ],
                  ),
                ),
              ),
            const SizedBox(height: 32),
            if (!_isInviteMode)
              FilledButton(
                onPressed: _isSubmitting ? null : _createGroup,
                style: FilledButton.styleFrom(minimumSize: const Size.fromHeight(48), backgroundColor: AppTheme.sanggamGold),
                child: _isSubmitting
                    ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                    : const Text('그룹 만들기'),
              ),
          ],
        ),
      ),
    );
  }
}
