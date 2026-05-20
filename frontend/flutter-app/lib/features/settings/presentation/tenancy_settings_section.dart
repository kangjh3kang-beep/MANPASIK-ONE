import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/active_tenant_provider.dart';
import 'package:manpasik/shared/providers/memberships_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// TenancySettingsSection — 설정 화면 내 조직 관리 섹션.
///
/// 구성:
///   1. 활성 조직 표시 (개인 모드 / 가족 / 병원 등)
///   2. "조직 전환" 버튼 → /tenant-switcher
///   3. "초대 수락" 버튼 → /invite (토큰 입력 폼)
///   4. admin/owner 멤버십이 있으면 "초대 발급" + "멤버 관리" 표시
class TenancySettingsSection extends ConsumerWidget {
  const TenancySettingsSection({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final active = ref.watch(activeTenantProvider);
    final asyncMems = ref.watch(myMembershipsProvider);

    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _SectionTitle(title: '조직 관리'),
            const SizedBox(height: 12),
            _ActiveTenantTile(active: active),
            const SizedBox(height: 8),
            _NavTile(
              key: const Key('nav-tenant-switcher'),
              icon: Icons.swap_horiz,
              label: '조직 전환',
              onTap: () => context.push('/tenant-switcher'),
            ),
            _NavTile(
              key: const Key('nav-invite-accept'),
              icon: Icons.qr_code_scanner,
              label: '초대 수락',
              onTap: () => context.push('/invite'),
            ),
            asyncMems.when(
              data: (list) => _AdminActions(memberships: list, active: active),
              loading: () => const SizedBox(),
              error: (_, __) => const SizedBox(),
            ),
          ],
        ),
      ),
    );
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({required this.title});
  final String title;

  @override
  Widget build(BuildContext context) {
    return Text(
      title,
      style: const TextStyle(
        color: SanggamTheme.primary,
        fontSize: 16,
        fontWeight: FontWeight.w600,
      ),
    );
  }
}

class _ActiveTenantTile extends StatelessWidget {
  const _ActiveTenantTile({required this.active});
  final ActiveTenant active;

  @override
  Widget build(BuildContext context) {
    final isPersonal = active.isPersonal;
    final label = isPersonal
        ? '개인 모드'
        : (active.displayName ?? active.tenantId ?? '');
    return Container(
      key: const Key('active-tenant-tile'),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: SanggamTheme.surfaceVariant,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Icon(
            isPersonal ? Icons.person : Icons.business,
            color: SanggamTheme.primary,
            size: 20,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('활성 조직',
                    style: TextStyle(
                        color: SanggamTheme.onSurfaceDim, fontSize: 11)),
                Text(
                  label,
                  style: const TextStyle(color: Colors.white, fontSize: 14),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _NavTile extends StatelessWidget {
  const _NavTile({
    super.key,
    required this.icon,
    required this.label,
    required this.onTap,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10),
        child: Row(
          children: [
            Icon(icon, color: SanggamTheme.primary, size: 20),
            const SizedBox(width: 12),
            Expanded(
              child: Text(label,
                  style: const TextStyle(color: Colors.white, fontSize: 14)),
            ),
            const Icon(Icons.chevron_right,
                color: SanggamTheme.onSurfaceDim, size: 18),
          ],
        ),
      ),
    );
  }
}

class _AdminActions extends StatelessWidget {
  const _AdminActions({required this.memberships, required this.active});

  final List<MembershipDto> memberships;
  final ActiveTenant active;

  @override
  Widget build(BuildContext context) {
    final adminMemberships = memberships
        .where((m) => m.role == 'admin' || m.role == 'owner')
        .toList();

    if (adminMemberships.isEmpty) {
      return const SizedBox(); // admin 권한 없음
    }

    final activeIsAdmin = !active.isPersonal &&
        adminMemberships.any((m) => m.tenantId == active.tenantId);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Divider(color: SanggamTheme.surfaceVariant, height: 16),
        const Padding(
          padding: EdgeInsets.only(bottom: 4),
          child: Text(
            '관리자 권한',
            style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 11),
          ),
        ),
        _NavTile(
          key: const Key('nav-invite-create'),
          icon: Icons.send,
          label: '초대 발급',
          onTap: () => context.push('/invite-create'),
        ),
        if (activeIsAdmin)
          _NavTile(
            key: const Key('nav-tenant-members'),
            icon: Icons.group,
            label: '${active.displayName ?? active.tenantId} 멤버 관리',
            onTap: () => context.push('/tenant-members/${active.tenantId}'),
          ),
      ],
    );
  }
}

