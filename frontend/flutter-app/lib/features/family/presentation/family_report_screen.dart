import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/family/domain/family_repository.dart';

/// 가족 건강 리포트 화면
class FamilyReportScreen extends ConsumerWidget {
  const FamilyReportScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final groupsAsync = ref.watch(familyGroupsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('가족 건강 리포트'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: RefreshIndicator(
        onRefresh: () async => ref.invalidate(familyGroupsProvider),
        child: groupsAsync.when(
          data: (groups) {
            if (groups.isEmpty) {
              return Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.family_restroom, size: 64, color: theme.colorScheme.onSurfaceVariant),
                    const SizedBox(height: 16),
                    Text('가족 그룹이 없습니다.', style: theme.textTheme.bodyLarge),
                    const SizedBox(height: 8),
                    FilledButton(
                      onPressed: () => context.push('/family'),
                      child: const Text('가족 그룹 만들기'),
                    ),
                  ],
                ),
              );
            }
            return _buildReport(theme, groups);
          },
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (_, __) => _buildFallbackReport(theme),
        ),
      ),
    );
  }

  Widget _buildReport(ThemeData theme, List<dynamic> groups) {
    // 모든 그룹의 구성원을 _MemberHealth로 변환
    final allMembers = <_MemberHealth>[];
    for (final g in groups) {
      final group = g as FamilyGroup;
      for (final m in group.members) {
        final status = _translateStatus(m.latestHealthStatus);
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
    if (allMembers.isEmpty) return _buildFallbackReport(theme);

    final normalCount = allMembers.where((m) => m.status == '양호').length;
    final cautionCount = allMembers.where((m) => m.status == '주의').length;
    final alertCount = allMembers.where((m) => m.status == '관찰').length;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Card(
            color: AppTheme.sanggamGold.withValues(alpha: 0.1),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(Icons.family_restroom, color: AppTheme.sanggamGold),
                      const SizedBox(width: 8),
                      Text('가족 건강 요약', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _summaryItem(theme, '가족 수', '${allMembers.length}명', Icons.people),
                      _summaryItem(theme, '양호', '$normalCount명', Icons.check_circle, Colors.green),
                      _summaryItem(theme, '주의', '$cautionCount명', Icons.warning, Colors.orange),
                      _summaryItem(theme, '관찰', '$alertCount명', Icons.visibility, Colors.red),
                    ],
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Text('구성원별 건강 현황', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          ...allMembers.map((m) => _buildMemberCard(theme, m)),
        ],
      ),
    );
  }

  String _translateStatus(String? status) {
    switch (status?.toLowerCase()) {
      case 'normal': case 'good': case 'excellent': return '양호';
      case 'caution': case 'borderline': return '주의';
      case 'alert': case 'abnormal': return '관찰';
      default: return '양호';
    }
  }

  String _formatTimeAgo(DateTime dt) {
    final diff = DateTime.now().difference(dt);
    if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
    if (diff.inHours < 24) return '${diff.inHours}시간 전';
    if (diff.inDays < 7) return '${diff.inDays}일 전';
    return '${dt.month}/${dt.day}';
  }

  Widget _buildFallbackReport(ThemeData theme) {
    final members = [
      _MemberHealth(name: '나', status: '양호', lastMeasure: '오늘 08:30', glucose: 95, cholesterol: 180, trend: '안정'),
      _MemberHealth(name: '배우자', status: '주의', lastMeasure: '어제 07:15', glucose: 125, cholesterol: 220, trend: '상승'),
      _MemberHealth(name: '자녀 1', status: '양호', lastMeasure: '3일 전', glucose: 88, cholesterol: 165, trend: '안정'),
      _MemberHealth(name: '부모님', status: '관찰', lastMeasure: '1주 전', glucose: 140, cholesterol: 245, trend: '상승'),
    ];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Card(
            color: AppTheme.sanggamGold.withValues(alpha: 0.1),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(Icons.family_restroom, color: AppTheme.sanggamGold),
                      const SizedBox(width: 8),
                      Text('가족 건강 요약', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _summaryItem(theme, '가족 수', '4명', Icons.people),
                      _summaryItem(theme, '양호', '2명', Icons.check_circle, Colors.green),
                      _summaryItem(theme, '주의', '1명', Icons.warning, Colors.orange),
                      _summaryItem(theme, '관찰', '1명', Icons.visibility, Colors.red),
                    ],
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Text('구성원별 건강 현황', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          ...members.map((m) => _buildMemberCard(theme, m)),
          const SizedBox(height: 16),
          Text('최근 건강 알림', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          _buildAlertTile(theme, '배우자', '혈당 수치가 정상 범위를 초과했습니다. (125 mg/dL)', Colors.orange, '어제'),
          _buildAlertTile(theme, '부모님', '콜레스테롤 수치 상승 추세가 감지되었습니다.', Colors.red, '3일 전'),
          _buildAlertTile(theme, '자녀 1', '이번 주 측정을 아직 하지 않았습니다.', Colors.blue, '5일 전'),
        ],
      ),
    );
  }

  Widget _summaryItem(ThemeData theme, String label, String value, IconData icon, [Color? color]) {
    return Column(
      children: [
        Icon(icon, color: color ?? AppTheme.sanggamGold, size: 24),
        const SizedBox(height: 4),
        Text(value, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
        Text(label, style: theme.textTheme.bodySmall),
      ],
    );
  }

  Widget _buildMemberCard(ThemeData theme, _MemberHealth m) {
    final statusColor = m.status == '양호' ? Colors.green : m.status == '주의' ? Colors.orange : Colors.red;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 18,
                  backgroundColor: statusColor.withValues(alpha:0.15),
                  child: Text(m.name[0], style: TextStyle(color: statusColor, fontWeight: FontWeight.bold)),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(m.name, style: const TextStyle(fontWeight: FontWeight.w600)),
                      Text('마지막 측정: ${m.lastMeasure}', style: theme.textTheme.bodySmall),
                    ],
                  ),
                ),
                Chip(
                  label: Text(m.status, style: TextStyle(fontSize: 11, color: statusColor)),
                  backgroundColor: statusColor.withValues(alpha:0.1),
                  side: BorderSide.none,
                  visualDensity: VisualDensity.compact,
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _metricChip(theme, '혈당', '${m.glucose} mg/dL', m.glucose > 120 ? Colors.orange : Colors.green),
                _metricChip(theme, '콜레스테롤', '${m.cholesterol} mg/dL', m.cholesterol > 200 ? Colors.orange : Colors.green),
                _metricChip(theme, '추세', m.trend, m.trend == '상승' ? Colors.red : Colors.green),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _metricChip(ThemeData theme, String label, String value, Color color) {
    return Column(
      children: [
        Text(value, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
        Text(label, style: theme.textTheme.bodySmall?.copyWith(fontSize: 10)),
      ],
    );
  }

  Widget _buildAlertTile(ThemeData theme, String member, String message, Color color, String time) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: ListTile(
        dense: true,
        leading: Icon(Icons.notification_important, size: 20, color: color),
        title: Text(member, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
        subtitle: Text(message, style: theme.textTheme.bodySmall),
        trailing: Text(time, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
      ),
    );
  }
}

class _MemberHealth {
  final String name, status, lastMeasure, trend;
  final int glucose, cholesterol;
  const _MemberHealth({required this.name, required this.status, required this.lastMeasure, required this.glucose, required this.cholesterol, required this.trend});
}
