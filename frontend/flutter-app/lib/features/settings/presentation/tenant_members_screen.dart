import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// TenantMembersScreen — 조직 멤버 관리 화면 (admin 전용).
///
/// 기능:
///   - 멤버 목록 조회 (`tenantMembersProvider`)
///   - 역할 변경 (드롭다운 → PATCH)
///   - 멤버 제거 (확인 다이얼로그 → DELETE)
class TenantMembersScreen extends ConsumerWidget {
  const TenantMembersScreen({super.key, required this.tenantId});

  final String tenantId;

  static const _availableRoles = <String>[
    'owner',
    'admin',
    'medical_staff',
    'member',
    'viewer',
  ];

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncMembers = ref.watch(tenantMembersProvider(tenantId));

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            _Header(tenantId: tenantId),
            Expanded(
              child: asyncMembers.when(
                data: (list) => _MembersList(
                  tenantId: tenantId,
                  members: list,
                  availableRoles: _availableRoles,
                ),
                loading: () => const Center(
                  child: CircularProgressIndicator(color: SanggamTheme.primary),
                ),
                error: (e, _) => _ErrorState(
                  message: '$e',
                  onRetry: () => ref.invalidate(tenantMembersProvider(tenantId)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Header extends StatelessWidget {
  const _Header({required this.tenantId});
  final String tenantId;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Row(
        children: [
          IconButton(
            icon: const Icon(Icons.arrow_back, color: SanggamTheme.primary),
            onPressed: () => Navigator.of(context).maybePop(),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '조직 멤버 관리',
                  style: TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 20,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  tenantId,
                  style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim, fontSize: 12),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _MembersList extends ConsumerWidget {
  const _MembersList({
    required this.tenantId,
    required this.members,
    required this.availableRoles,
  });

  final String tenantId;
  final List<MembershipDto> members;
  final List<String> availableRoles;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (members.isEmpty) {
      return const Center(
        child: Padding(
          padding: EdgeInsets.all(32),
          child: Text(
            '멤버가 없습니다',
            style: TextStyle(color: SanggamTheme.onSurfaceDim),
          ),
        ),
      );
    }
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: members.length,
      itemBuilder: (context, idx) {
        final m = members[idx];
        return Padding(
          padding: const EdgeInsets.only(bottom: 8),
          child: SanggamContainer(
            child: _MemberTile(
              tenantId: tenantId,
              membership: m,
              availableRoles: availableRoles,
              onRoleChanged: (newRole) =>
                  _changeRole(context, ref, m, newRole),
              onRemove: () => _confirmRemove(context, ref, m),
            ),
          ),
        );
      },
    );
  }

  Future<void> _changeRole(BuildContext context, WidgetRef ref,
      MembershipDto m, String newRole) async {
    if (newRole == m.role) return;
    try {
      final client = ref.read(restClientProvider);
      await client.updateTenantMemberRole(
        tenantId: m.tenantId,
        userId: m.userId,
        newRole: newRole,
      );
      ref.invalidate(tenantMembersProvider(m.tenantId));
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('${m.userId} 역할 → $newRole')),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('역할 변경 실패: $e')),
        );
      }
    }
  }

  Future<void> _confirmRemove(
      BuildContext context, WidgetRef ref, MembershipDto m) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        title: const Text('멤버 제거', style: TextStyle(color: Colors.white)),
        content: Text(
          '${m.userId} 멤버를 정말 제거하시겠습니까?',
          style: const TextStyle(color: SanggamTheme.onSurfaceDim),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('취소', style: TextStyle(color: SanggamTheme.onSurfaceDim)),
          ),
          TextButton(
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('제거', style: TextStyle(color: SanggamTheme.error)),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    try {
      final client = ref.read(restClientProvider);
      await client.removeTenantMember(tenantId: m.tenantId, userId: m.userId);
      ref.invalidate(tenantMembersProvider(m.tenantId));
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('${m.userId} 제거 완료')),
        );
      }
    } catch (e) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('제거 실패: $e')),
        );
      }
    }
  }
}

class _MemberTile extends StatelessWidget {
  const _MemberTile({
    required this.tenantId,
    required this.membership,
    required this.availableRoles,
    required this.onRoleChanged,
    required this.onRemove,
  });

  final String tenantId;
  final MembershipDto membership;
  final List<String> availableRoles;
  final ValueChanged<String> onRoleChanged;
  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(12),
      child: Row(
        children: [
          const Icon(Icons.person, color: SanggamTheme.primary, size: 28),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  membership.userId,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  membership.active ? '활성' : '비활성',
                  style: TextStyle(
                    color: membership.active
                        ? SanggamTheme.primary
                        : SanggamTheme.onSurfaceDim,
                    fontSize: 11,
                  ),
                ),
              ],
            ),
          ),
          DropdownButton<String>(
            key: Key('role-${membership.userId}'),
            value: membership.role,
            dropdownColor: SanggamTheme.surface,
            style: const TextStyle(color: Colors.white, fontSize: 13),
            underline: const SizedBox.shrink(),
            items: availableRoles
                .map((r) => DropdownMenuItem(value: r, child: Text(r)))
                .toList(),
            onChanged: (v) {
              if (v != null) onRoleChanged(v);
            },
          ),
          IconButton(
            key: Key('remove-${membership.userId}'),
            icon: const Icon(Icons.delete_outline, color: SanggamTheme.error),
            onPressed: onRemove,
          ),
        ],
      ),
    );
  }
}

class _ErrorState extends StatelessWidget {
  const _ErrorState({required this.message, required this.onRetry});
  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, color: SanggamTheme.error, size: 48),
            const SizedBox(height: 16),
            const Text(
              '멤버 조회 실패',
              style: TextStyle(
                color: SanggamTheme.error,
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(message,
                textAlign: TextAlign.center,
                style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim, fontSize: 12)),
            const SizedBox(height: 16),
            TextButton(
              onPressed: onRetry,
              child: const Text('다시 시도',
                  style: TextStyle(color: SanggamTheme.primary)),
            ),
          ],
        ),
      ),
    );
  }
}
