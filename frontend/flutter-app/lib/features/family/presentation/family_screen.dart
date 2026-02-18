import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';

/// 가족 건강 관리 화면
class FamilyScreen extends ConsumerStatefulWidget {
  const FamilyScreen({super.key});

  @override
  ConsumerState<FamilyScreen> createState() => _FamilyScreenState();
}

class _FamilyScreenState extends ConsumerState<FamilyScreen> {
  bool _shareResults = false;
  bool _alertOnAbnormal = true;
  bool _sendReminders = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final groupsAsync = ref.watch(familyGroupsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('가족 건강'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.assessment_outlined),
            tooltip: '가족 건강 리포트',
            onPressed: () => context.push('/family/report'),
          ),
          IconButton(
            icon: const Icon(Icons.group_add),
            tooltip: '가족 초대',
            onPressed: () => _showInviteDialog(context),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(familyGroupsProvider),
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              groupsAsync.when(
                data: (groups) {
                  if (groups.isEmpty) return _buildCreateGroupCard(theme);
                  return Column(
                    children: groups.map((g) => _buildGroupCard(theme, g)).toList(),
                  );
                },
                loading: () => const Center(child: Padding(padding: EdgeInsets.all(24), child: CircularProgressIndicator())),
                error: (_, __) => _buildCreateGroupCard(theme),
              ),
              const SizedBox(height: 24),

              Text('데이터 공유 설정', style: theme.textTheme.titleLarge),
              const SizedBox(height: 12),
              _buildSharingSettings(theme),
              const SizedBox(height: 24),

              // 보호자 대시보드 바로가기
              Card(
                color: theme.colorScheme.primaryContainer.withOpacity(0.3),
                child: ListTile(
                  leading: Icon(Icons.shield, color: theme.colorScheme.primary),
                  title: const Text('보호자 대시보드'),
                  subtitle: const Text('가족 건강 현황 한눈에 확인'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/family/guardian'),
                ),
              ),
              const SizedBox(height: 8),

              // 긴급 대응 설정 바로가기
              Card(
                color: theme.colorScheme.errorContainer.withOpacity(0.3),
                child: ListTile(
                  leading: Icon(Icons.emergency_outlined, color: theme.colorScheme.error),
                  title: const Text('긴급 대응 설정'),
                  subtitle: const Text('긴급 연락망, 위험 감지 기준 설정'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/settings/emergency'),
                ),
              ),
              const SizedBox(height: 24),

              Text('가족 알림', style: theme.textTheme.titleLarge),
              const SizedBox(height: 12),
              _buildFamilyAlerts(theme),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildCreateGroupCard(ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            Icon(Icons.family_restroom, size: 48, color: theme.colorScheme.primary),
            const SizedBox(height: 12),
            Text('가족 그룹을 만들어보세요', style: theme.textTheme.titleMedium),
            const SizedBox(height: 8),
            Text(
              '가족 구성원의 건강 데이터를 함께 관리하고\n측정 리마인더를 보낼 수 있습니다',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
            ),
            const SizedBox(height: 16),
            FilledButton.icon(
              onPressed: () => context.push('/family/create'),
              icon: const Icon(Icons.add),
              label: const Text('가족 그룹 만들기'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildGroupCard(ThemeData theme, FamilyGroup group) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(Icons.family_restroom, color: theme.colorScheme.primary),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(group.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                ),
                Text('${group.members.length}명', style: theme.textTheme.bodySmall),
              ],
            ),
            const Divider(height: 24),
            ...group.members.map((m) => _buildMemberTile(theme, m)),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () async {
                final inv = await ref.read(familyRepositoryProvider).createInvitation(group.id);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('초대 코드: ${inv.inviteCode}')),
                  );
                }
              },
              icon: const Icon(Icons.link, size: 18),
              label: const Text('초대 링크 생성'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMemberTile(ThemeData theme, FamilyMember member) {
    final statusColor = switch (member.latestHealthStatus) {
      'normal' => Colors.green,
      'caution' => Colors.orange,
      'alert' => Colors.red,
      _ => theme.colorScheme.outline,
    };
    final roleText = switch (member.role) {
      FamilyRole.owner => '소유자',
      FamilyRole.admin => '관리자',
      FamilyRole.member => '구성원',
    };

    return ListTile(
      dense: true,
      leading: CircleAvatar(
        radius: 18,
        child: Text(member.displayName.isNotEmpty ? member.displayName[0] : '?'),
      ),
      title: Text(member.displayName),
      subtitle: Text(roleText),
      trailing: member.latestHealthStatus != null
          ? Container(
              width: 10,
              height: 10,
              decoration: BoxDecoration(color: statusColor, shape: BoxShape.circle),
            )
          : null,
    );
  }

  Widget _buildSharingSettings(ThemeData theme) {
    return Card(
      child: Column(
        children: [
          SwitchListTile(
            title: const Text('측정 결과 공유'),
            subtitle: const Text('가족에게 측정 결과를 자동으로 공유'),
            value: _shareResults,
            onChanged: (v) => setState(() => _shareResults = v),
            secondary: const Icon(Icons.share),
          ),
          const Divider(height: 1),
          SwitchListTile(
            title: const Text('이상 수치 알림'),
            subtitle: const Text('가족의 이상 수치 발생 시 알림'),
            value: _alertOnAbnormal,
            onChanged: (v) => setState(() => _alertOnAbnormal = v),
            secondary: const Icon(Icons.notifications_active),
          ),
          const Divider(height: 1),
          SwitchListTile(
            title: const Text('측정 리마인더'),
            subtitle: const Text('가족에게 측정 리마인더 전송'),
            value: _sendReminders,
            onChanged: (v) => setState(() => _sendReminders = v),
            secondary: const Icon(Icons.alarm),
          ),
        ],
      ),
    );
  }

  Widget _buildFamilyAlerts(ThemeData theme) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: Text(
            '가족 알림이 없습니다',
            style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.outline),
          ),
        ),
      ),
    );
  }

  void _showCreateGroupDialog(BuildContext context) {
    final nameCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('가족 그룹 만들기'),
        content: TextField(
          controller: nameCtrl,
          decoration: const InputDecoration(labelText: '그룹 이름'),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () async {
              if (nameCtrl.text.isNotEmpty) {
                await ref.read(familyRepositoryProvider).createGroup(nameCtrl.text);
                ref.invalidate(familyGroupsProvider);
                if (ctx.mounted) Navigator.pop(ctx);
              }
            },
            child: const Text('만들기'),
          ),
        ],
      ),
    );
  }

  void _showInviteDialog(BuildContext context) {
    final codeCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('초대 코드 입력'),
        content: TextField(
          controller: codeCtrl,
          decoration: const InputDecoration(labelText: '초대 코드'),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () async {
              if (codeCtrl.text.isNotEmpty) {
                await ref.read(familyRepositoryProvider).acceptInvitation(codeCtrl.text);
                ref.invalidate(familyGroupsProvider);
                if (ctx.mounted) Navigator.pop(ctx);
              }
            },
            child: const Text('참가'),
          ),
        ],
      ),
    );
  }
}
