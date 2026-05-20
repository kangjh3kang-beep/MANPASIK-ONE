import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// GuardianDashboardScreen — Sanggam Orbit 보호자 대시보드
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* ~6x → 직접 TextStyle
// [Rule 4] theme.colorScheme.onSurfaceVariant → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] withOpacity 3x → withValues(alpha:)
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] w:12→16, h:12→16, h:4→8, borderRadius:12→16
// ───────────────────────────────────────────────────

/// 보호자 대시보드 화면
class GuardianDashboardScreen extends ConsumerWidget {
  const GuardianDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final groupsAsync = ref.watch(familyGroupsProvider);

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
                      '보호자 대시보드',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.refresh,
                        color: Colors.white),
                    tooltip: '새로고침',
                    onPressed: () =>
                        ref.invalidate(familyGroupsProvider),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: groupsAsync.when(
                loading: () => const Center(
                    child: CircularProgressIndicator()),
                error: (e, _) => Center(
                  child: Text(
                    '데이터를 불러올 수 없습니다: $e',
                    style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: 14,
                    ),
                  ),
                ),
                data: (groups) {
                  final allMembers = <_FamilyMemberData>[];
                  for (final group in groups) {
                    for (final member in group.members) {
                      allMembers.add(_FamilyMemberData(
                        id: member.userId,
                        name: member.displayName,
                        relation: member.role.name,
                        status:
                            member.latestHealthStatus ??
                                'normal',
                        lastMeasured:
                            member.lastMeasurementAt != null
                                ? '${DateTime.now().difference(member.lastMeasurementAt!).inHours}시간 전'
                                : '미측정',
                        weeklyTrend: [
                          0.7, 0.8, 0.7, 0.8, 0.9, 0.8, 0.7
                        ],
                      ));
                    }
                  }

                  if (allMembers.isEmpty) {
                    allMembers.addAll(_fallbackMembers);
                  }

                  final warningCount = allMembers
                      .where((m) =>
                          m.status == 'warning' ||
                          m.status == 'danger')
                      .length;

                  return RefreshIndicator(
                    onRefresh: () async => ref
                        .invalidate(familyGroupsProvider),
                    child: ListView(
                      padding: const EdgeInsets.all(16),
                      children: [
                        if (warningCount > 0) ...[
                          SanggamContainer(
                            borderRadius: 16,
                            borderColor: SanggamTheme.primary
                                .withValues(alpha: 0.3),
                            padding:
                                const EdgeInsets.all(16),
                            child: Row(
                              children: [
                                const Icon(
                                    Icons.warning_amber,
                                    color: SanggamTheme
                                        .primary,
                                    size: 32),
                                const SizedBox(width: 16),
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment
                                            .start,
                                    children: [
                                      Text(
                                        '주의 알림 $warningCount건',
                                        style:
                                            const TextStyle(
                                          color:
                                              SanggamTheme
                                                  .primary,
                                          fontSize: 16,
                                          fontWeight:
                                              FontWeight
                                                  .bold,
                                        ),
                                      ),
                                      const Text(
                                        '일부 구성원의 건강 수치에 주의가 필요합니다.',
                                        style: TextStyle(
                                          color: SanggamTheme
                                              .onSurfaceDim,
                                          fontSize: 12,
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                              ],
                            ),
                          ),
                          const SizedBox(height: 8),
                          // 크로스 도메인 CTA: 긴급 연락 / 원격 진료
                          Row(
                            children: [
                              Expanded(
                                child: OutlinedButton.icon(
                                  onPressed: () =>
                                      context.push('/settings/emergency'),
                                  icon: const Icon(Icons.emergency,
                                      size: 16),
                                  label: const Text('긴급 연락',
                                      style: TextStyle(fontSize: 12)),
                                  style: OutlinedButton.styleFrom(
                                    foregroundColor:
                                        SanggamTheme.error,
                                    side: BorderSide(
                                      color: SanggamTheme.error
                                          .withValues(alpha: 0.5),
                                    ),
                                  ),
                                ),
                              ),
                              const SizedBox(width: 8),
                              Expanded(
                                child: OutlinedButton.icon(
                                  onPressed: () =>
                                      context.push('/medical/telemedicine'),
                                  icon: const Icon(Icons.videocam,
                                      size: 16),
                                  label: const Text('원격 진료',
                                      style: TextStyle(fontSize: 12)),
                                  style: OutlinedButton.styleFrom(
                                    foregroundColor:
                                        SanggamTheme.jagaeCyan,
                                    side: BorderSide(
                                      color: SanggamTheme.jagaeCyan
                                          .withValues(alpha: 0.5),
                                    ),
                                  ),
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 16),
                        ],
                        Text(
                          '구성원 건강 현황 (${allMembers.length}명)',
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 8),
                        ...allMembers.map((m) =>
                            _buildMemberCard(context, m)),
                      ],
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMemberCard(
      BuildContext context, _FamilyMemberData member) {
    final statusColor = member.status == 'normal'
        ? SanggamTheme.jagaeCyan
        : member.status == 'warning'
            ? SanggamTheme.primary
            : SanggamTheme.error;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: const EdgeInsets.all(16),
        onTap: () =>
            context.push('/family/member/${member.id}/edit'),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                CircleAvatar(
                  backgroundColor: SanggamTheme.primary
                      .withValues(alpha: 0.2),
                  child: Text(
                    member.name.isNotEmpty
                        ? member.name[0]
                        : '?',
                    style: const TextStyle(
                      color: SanggamTheme.primary,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment:
                        CrossAxisAlignment.start,
                    children: [
                      Text(
                        member.name,
                        style: const TextStyle(
                          color: Colors.white,
                          fontSize: 14,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        member.relation,
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(
                      horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: statusColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Text(
                    member.status == 'normal'
                        ? '정상'
                        : member.status == 'warning'
                            ? '주의'
                            : '위험',
                    style: TextStyle(
                      fontSize: 11,
                      color: statusColor,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: List.generate(7, (i) {
                final val =
                    i < member.weeklyTrend.length
                        ? member.weeklyTrend[i]
                        : 0.5;
                return Expanded(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 1),
                    child: Column(
                      children: [
                        Container(
                          height: 24,
                          decoration: BoxDecoration(
                            color: val > 0.7
                                ? SanggamTheme.jagaeCyan
                                : val > 0.4
                                    ? SanggamTheme.primary
                                    : SanggamTheme.error,
                            borderRadius:
                                BorderRadius.circular(4),
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          ['월', '화', '수', '목', '금', '토', '일'][i],
                          style: const TextStyle(
                            fontSize: 9,
                            color:
                                SanggamTheme.onSurfaceDim,
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              }),
            ),
            const SizedBox(height: 8),
            Text(
              '최근 측정: ${member.lastMeasured}',
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              ),
            ),
          ],
        ),
      ),
    );
  }

  static final _fallbackMembers = [
    _FamilyMemberData(
      id: '1',
      name: '어머니',
      relation: '부모',
      status: 'warning',
      lastMeasured: '2시간 전',
      weeklyTrend: [0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.5],
    ),
    _FamilyMemberData(
      id: '2',
      name: '아버지',
      relation: '부모',
      status: 'normal',
      lastMeasured: '4시간 전',
      weeklyTrend: [0.7, 0.8, 0.9, 0.8, 0.7, 0.8, 0.9],
    ),
    _FamilyMemberData(
      id: '3',
      name: '배우자',
      relation: '배우자',
      status: 'normal',
      lastMeasured: '1일 전',
      weeklyTrend: [0.9, 0.8, 0.7, 0.8, 0.9, 0.8, 0.7],
    ),
  ];
}

class _FamilyMemberData {
  final String id, name, relation, status, lastMeasured;
  final List<double> weeklyTrend;
  const _FamilyMemberData({
    required this.id,
    required this.name,
    required this.relation,
    required this.status,
    required this.lastMeasured,
    required this.weeklyTrend,
  });
}
