import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// FamilyScreen — Sanggam Orbit 가족 건강
//
// [Rule 4] Theme(data:ThemeData.dark().copyWith(...)) 래퍼 제거
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] GlassmorphismCard 4x → SanggamContainer
// [Rule 4] HanjiPanel 2x → SanggamContainer
// [Rule 4] theme.colorScheme.* ~12x → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~7x → 직접 TextStyle
// [Rule 4] Theme.of(context) 제거
// [Rule 4] withOpacity ~7x → withValues(alpha:)
// [Rule 4] Color(0xFF...) 하드코딩 4x → SanggamTheme 상수
// [Rule 4] GlassmorphismCard+Container 이중 래핑 → SanggamContainer 단일
// [Rule 4] _showCreateGroupDialog 미사용 코드 제거
// [Rule 2] h:12→16, padding 20→24, borderRadius 20→24, v:6→8, top:4→8
// ───────────────────────────────────────────────────

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
    final groupsAsync = ref.watch(familyGroupsProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: RefreshIndicator(
          color: SanggamTheme.primary,
          onRefresh: () async =>
              ref.invalidate(familyGroupsProvider),
          child: SingleChildScrollView(
            physics: const AlwaysScrollableScrollPhysics(),
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 헤더: 타이틀 + 액션
                Row(
                  children: [
                    const Expanded(
                      child: Text(
                        '가족 건강',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          letterSpacing: -0.5,
                        ),
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.assessment_outlined,
                          color: Colors.white),
                      tooltip: '가족 건강 리포트',
                      onPressed: () => context.push('/family/report'),
                    ),
                    IconButton(
                      icon: const Icon(Icons.group_add,
                          color: Colors.white),
                      tooltip: '가족 초대',
                      onPressed: () => _showInviteDialog(context),
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                groupsAsync.when(
                  data: (groups) {
                    if (groups.isEmpty) return _buildCreateGroupCard();
                    return Column(
                      children: groups
                          .map((g) => _buildGroupCard(g))
                          .toList(),
                    );
                  },
                  loading: () => const Center(
                      child: Padding(
                          padding: EdgeInsets.all(24),
                          child: CircularProgressIndicator())),
                  error: (_, __) => _buildCreateGroupCard(),
                ),
                const SizedBox(height: 24),

                const Text(
                  '데이터 공유 설정',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                _buildSharingSettings(),
                const SizedBox(height: 24),

                // 보호자 대시보드 바로가기
                SanggamContainer(
                  borderRadius: 16,
                  padding: EdgeInsets.zero,
                  child: ListTile(
                    leading: const Icon(Icons.shield,
                        color: SanggamTheme.primary),
                    title: const Text('보호자 대시보드'),
                    subtitle: const Text('가족 건강 현황 한눈에 확인'),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push('/family/guardian'),
                  ),
                ),
                const SizedBox(height: 8),

                // 긴급 대응 설정 바로가기
                SanggamContainer(
                  borderRadius: 16,
                  borderColor:
                      SanggamTheme.error.withValues(alpha: 0.5),
                  padding: EdgeInsets.zero,
                  child: ListTile(
                    leading: const Icon(Icons.emergency_outlined,
                        color: SanggamTheme.error),
                    title: const Text('긴급 대응 설정'),
                    subtitle: const Text('긴급 연락망, 위험 감지 기준 설정'),
                    trailing: const Icon(Icons.chevron_right),
                    onTap: () => context.push('/settings/emergency'),
                  ),
                ),
                const SizedBox(height: 24),

                const Text(
                  '가족 알림',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 16),
                _buildFamilyAlerts(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildCreateGroupCard() {
    return SanggamContainer(
      borderRadius: 24,
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          const Icon(Icons.family_restroom,
              size: 48, color: SanggamTheme.primary),
          const SizedBox(height: 16),
          const Text(
            '가족 그룹을 만들어보세요',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '가족 구성원의 건강 데이터를 함께 관리하고\n측정 리마인더를 보낼 수 있습니다',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
          const SizedBox(height: 16),
          FilledButton.icon(
            onPressed: () => context.push('/family/create'),
            icon: const Icon(Icons.add),
            label: const Text('가족 그룹 만들기'),
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.primary,
              foregroundColor: SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16)),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildGroupCard(FamilyGroup group) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.family_restroom,
                    color: SanggamTheme.primary),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    group.name,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                Text(
                  '${group.members.length}명',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
            const Divider(height: 24),
            ...group.members.map((m) => _buildMemberTile(m)),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () async {
                final inv = await ref
                    .read(familyRepositoryProvider)
                    .createInvitation(group.id);
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                        content:
                            Text('초대 코드: ${inv.inviteCode}')),
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

  Widget _buildMemberTile(FamilyMember member) {
    final statusColor = switch (member.latestHealthStatus) {
      'normal' => SanggamTheme.jagaeCyan,
      'caution' => SanggamTheme.primary,
      'alert' => SanggamTheme.error,
      _ => SanggamTheme.onSurfaceDim,
    };
    final roleText = switch (member.role) {
      FamilyRole.owner => '소유자',
      FamilyRole.admin => '관리자',
      FamilyRole.member => '구성원',
    };

    return ListTile(
      dense: false,
      minVerticalPadding: 16,
      leading: Container(
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          boxShadow: [
            BoxShadow(
              color: statusColor.withValues(alpha: 0.4),
              blurRadius: 12,
              spreadRadius: 2,
            ),
          ],
        ),
        child: CircleAvatar(
          radius: 24,
          backgroundColor: SanggamTheme.surface,
          child: Text(
            member.displayName.isNotEmpty
                ? member.displayName[0]
                : '?',
            style: const TextStyle(
                color: Colors.white, fontWeight: FontWeight.bold),
          ),
        ),
      ),
      title: Text(member.displayName,
          style: const TextStyle(
              fontWeight: FontWeight.bold, fontSize: 16)),
      subtitle: Padding(
        padding: const EdgeInsets.only(top: 8),
        child: Text(roleText,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 13,
            )),
      ),
      trailing: member.latestHealthStatus != null
          ? Container(
              padding: const EdgeInsets.symmetric(
                  horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: statusColor.withValues(alpha: 0.15),
                borderRadius: BorderRadius.circular(24),
                border: Border.all(
                    color: statusColor.withValues(alpha: 0.5)),
              ),
              child: Text(
                member.latestHealthStatus == 'normal'
                    ? '정상'
                    : member.latestHealthStatus == 'caution'
                        ? '주의'
                        : '위험',
                style: TextStyle(
                    color: statusColor,
                    fontSize: 12,
                    fontWeight: FontWeight.bold),
              ),
            )
          : null,
    );
  }

  Widget _buildSharingSettings() {
    return SanggamContainer(
      borderRadius: 16,
      padding: EdgeInsets.zero,
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
            onChanged: (v) =>
                setState(() => _alertOnAbnormal = v),
            secondary: const Icon(Icons.notifications_active),
          ),
          const Divider(height: 1),
          SwitchListTile(
            title: const Text('측정 리마인더'),
            subtitle: const Text('가족에게 측정 리마인더 전송'),
            value: _sendReminders,
            onChanged: (v) =>
                setState(() => _sendReminders = v),
            secondary: const Icon(Icons.alarm),
          ),
        ],
      ),
    );
  }

  Widget _buildFamilyAlerts() {
    return SanggamContainer(
      borderRadius: 24,
      padding: const EdgeInsets.all(24),
      child: const Center(
        child: Text(
          '가족 알림이 없습니다',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 14,
          ),
        ),
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
          TextButton(
              onPressed: () => Navigator.pop(ctx),
              child: const Text('취소')),
          FilledButton(
            onPressed: () async {
              if (codeCtrl.text.isNotEmpty) {
                await ref
                    .read(familyRepositoryProvider)
                    .acceptInvitation(codeCtrl.text);
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
