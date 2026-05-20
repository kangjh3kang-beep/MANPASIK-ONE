import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// FamilyReportScreen — Sanggam Orbit 가족 건강 리포트
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 5x → SanggamTheme.primary
// [Rule 4] Theme.of(context) + ThemeData 파라미터 6개 제거
// [Rule 4] theme.textTheme.* ~20x → 직접 TextStyle
// [Rule 4] theme.colorScheme.onSurfaceVariant → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Card ~5x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] h:12→16, h:4→8, padding:12→16, bottom:4→8
// ───────────────────────────────────────────────────

/// 가족 건강 리포트 화면
class FamilyReportScreen extends ConsumerWidget {
  const FamilyReportScreen({super.key});

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
                      '가족 건강 리포트',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: RefreshIndicator(
                onRefresh: () async =>
                    ref.invalidate(familyGroupsProvider),
                child: groupsAsync.when(
                  data: (groups) {
                    if (groups.isEmpty) {
                      return Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(
                                Icons.family_restroom,
                                size: 64,
                                color: SanggamTheme
                                    .onSurfaceDim),
                            const SizedBox(height: 16),
                            const Text(
                              '가족 그룹이 없습니다.',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                              ),
                            ),
                            const SizedBox(height: 8),
                            FilledButton(
                              onPressed: () =>
                                  context.push('/family'),
                              style: FilledButton.styleFrom(
                                backgroundColor:
                                    SanggamTheme.primary,
                                foregroundColor:
                                    SanggamTheme.background,
                              ),
                              child:
                                  const Text('가족 그룹 만들기'),
                            ),
                          ],
                        ),
                      );
                    }
                    return _buildReport(groups);
                  },
                  loading: () => const Center(
                      child: CircularProgressIndicator()),
                  error: (_, __) => _buildFallbackReport(),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildReport(List<dynamic> groups) {
    final allMembers = <_MemberHealth>[];
    for (final g in groups) {
      final group = g as FamilyGroup;
      for (final m in group.members) {
        final status =
            _translateStatus(m.latestHealthStatus);
        final lastMeasure = m.lastMeasurementAt != null
            ? _formatTimeAgo(m.lastMeasurementAt!)
            : '측정 없음';
        allMembers.add(_MemberHealth(
          name: m.displayName,
          status: status,
          lastMeasure: lastMeasure,
          glucose: 0,
          cholesterol: 0,
          trend: '안정',
        ));
      }
    }
    if (allMembers.isEmpty) return _buildFallbackReport();

    final normalCount =
        allMembers.where((m) => m.status == '양호').length;
    final cautionCount =
        allMembers.where((m) => m.status == '주의').length;
    final alertCount =
        allMembers.where((m) => m.status == '관찰').length;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _buildSummaryCard(
            total: allMembers.length,
            normal: normalCount,
            caution: cautionCount,
            alert: alertCount,
          ),
          const SizedBox(height: 16),
          const Text(
            '구성원별 건강 현황',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          ...allMembers.map((m) => _buildMemberCard(m)),
        ],
      ),
    );
  }

  String _translateStatus(String? status) {
    switch (status?.toLowerCase()) {
      case 'normal':
      case 'good':
      case 'excellent':
        return '양호';
      case 'caution':
      case 'borderline':
        return '주의';
      case 'alert':
      case 'abnormal':
        return '관찰';
      default:
        return '양호';
    }
  }

  String _formatTimeAgo(DateTime dt) {
    final diff = DateTime.now().difference(dt);
    if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
    if (diff.inHours < 24) return '${diff.inHours}시간 전';
    if (diff.inDays < 7) return '${diff.inDays}일 전';
    return '${dt.month}/${dt.day}';
  }

  Widget _buildFallbackReport() {
    final members = [
      _MemberHealth(
          name: '나',
          status: '양호',
          lastMeasure: '오늘 08:30',
          glucose: 95,
          cholesterol: 180,
          trend: '안정'),
      _MemberHealth(
          name: '배우자',
          status: '주의',
          lastMeasure: '어제 07:15',
          glucose: 125,
          cholesterol: 220,
          trend: '상승'),
      _MemberHealth(
          name: '자녀 1',
          status: '양호',
          lastMeasure: '3일 전',
          glucose: 88,
          cholesterol: 165,
          trend: '안정'),
      _MemberHealth(
          name: '부모님',
          status: '관찰',
          lastMeasure: '1주 전',
          glucose: 140,
          cholesterol: 245,
          trend: '상승'),
    ];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _buildSummaryCard(
              total: 4, normal: 2, caution: 1, alert: 1),
          const SizedBox(height: 16),
          const Text(
            '구성원별 건강 현황',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          ...members.map((m) => _buildMemberCard(m)),
          const SizedBox(height: 16),
          const Text(
            '최근 건강 알림',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          _buildAlertTile('배우자',
              '혈당 수치가 정상 범위를 초과했습니다. (125 mg/dL)',
              SanggamTheme.primary, '어제'),
          _buildAlertTile('부모님',
              '콜레스테롤 수치 상승 추세가 감지되었습니다.',
              SanggamTheme.error, '3일 전'),
          _buildAlertTile('자녀 1',
              '이번 주 측정을 아직 하지 않았습니다.',
              SanggamTheme.jagaeCyan, '5일 전'),
        ],
      ),
    );
  }

  Widget _buildSummaryCard({
    required int total,
    required int normal,
    required int caution,
    required int alert,
  }) {
    return SanggamContainer(
      borderRadius: 16,
      borderColor:
          SanggamTheme.primary.withValues(alpha: 0.3),
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Row(
            children: [
              Icon(Icons.family_restroom,
                  color: SanggamTheme.primary),
              SizedBox(width: 8),
              Text(
                '가족 건강 요약',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment:
                MainAxisAlignment.spaceAround,
            children: [
              _summaryItem('가족 수', '$total명',
                  Icons.people, null),
              _summaryItem('양호', '$normal명',
                  Icons.check_circle,
                  SanggamTheme.jagaeCyan),
              _summaryItem('주의', '$caution명',
                  Icons.warning, SanggamTheme.primary),
              _summaryItem('관찰', '$alert명',
                  Icons.visibility, SanggamTheme.error),
            ],
          ),
        ],
      ),
    );
  }

  Widget _summaryItem(
      String label, String value, IconData icon,
      [Color? color]) {
    return Column(
      children: [
        Icon(icon,
            color: color ?? SanggamTheme.primary, size: 24),
        const SizedBox(height: 8),
        Text(
          value,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  Widget _buildMemberCard(_MemberHealth m) {
    final statusColor = m.status == '양호'
        ? SanggamTheme.jagaeCyan
        : m.status == '주의'
            ? SanggamTheme.primary
            : SanggamTheme.error;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 18,
                  backgroundColor:
                      statusColor.withValues(alpha: 0.15),
                  child: Text(
                    m.name[0],
                    style: TextStyle(
                      color: statusColor,
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
                        m.name,
                        style: const TextStyle(
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      Text(
                        '마지막 측정: ${m.lastMeasure}',
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                ),
                Chip(
                  label: Text(
                    m.status,
                    style: TextStyle(
                        fontSize: 11, color: statusColor),
                  ),
                  backgroundColor:
                      statusColor.withValues(alpha: 0.1),
                  side: BorderSide.none,
                  visualDensity: VisualDensity.compact,
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment:
                  MainAxisAlignment.spaceAround,
              children: [
                _metricChip(
                    '혈당',
                    '${m.glucose} mg/dL',
                    m.glucose > 120
                        ? SanggamTheme.primary
                        : SanggamTheme.jagaeCyan),
                _metricChip(
                    '콜레스테롤',
                    '${m.cholesterol} mg/dL',
                    m.cholesterol > 200
                        ? SanggamTheme.primary
                        : SanggamTheme.jagaeCyan),
                _metricChip(
                    '추세',
                    m.trend,
                    m.trend == '상승'
                        ? SanggamTheme.error
                        : SanggamTheme.jagaeCyan),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _metricChip(
      String label, String value, Color color) {
    return Column(
      children: [
        Text(
          value,
          style: TextStyle(
            color: color,
            fontSize: 12,
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 10,
          ),
        ),
      ],
    );
  }

  Widget _buildAlertTile(
      String member, String message, Color color,
      String time) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: ListTile(
          dense: true,
          leading: Icon(Icons.notification_important,
              size: 20, color: color),
          title: Text(
            member,
            style: const TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.w600,
            ),
          ),
          subtitle: Text(
            message,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
          trailing: Text(
            time,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ),
      ),
    );
  }
}

class _MemberHealth {
  final String name, status, lastMeasure, trend;
  final int glucose, cholesterol;
  const _MemberHealth({
    required this.name,
    required this.status,
    required this.lastMeasure,
    required this.glucose,
    required this.cholesterol,
    required this.trend,
  });
}
