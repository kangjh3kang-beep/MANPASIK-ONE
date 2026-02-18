import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 보호자 대시보드 화면
class GuardianDashboardScreen extends ConsumerWidget {
  const GuardianDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('보호자 대시보드'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 알림 요약
          Card(
            color: Colors.orange.withOpacity(0.1),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  const Icon(Icons.warning_amber, color: Colors.orange, size: 32),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('주의 알림 1건', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold, color: Colors.orange)),
                        Text('어머니의 혈압 수치가 정상 범위를 초과했습니다.', style: theme.textTheme.bodySmall),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // 구성원별 건강 현황
          Text('구성원 건강 현황', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          ..._familyMembers.map((m) => _buildMemberCard(context, theme, m)),
        ],
      ),
    );
  }

  Widget _buildMemberCard(BuildContext context, ThemeData theme, _FamilyMemberData member) {
    final statusColor = member.status == 'normal' ? Colors.green : member.status == 'warning' ? Colors.orange : Colors.red;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: InkWell(
        onTap: () => context.push('/family/member/${member.id}/edit'),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  CircleAvatar(
                    backgroundColor: AppTheme.sanggamGold.withOpacity(0.2),
                    child: Text(member.name[0], style: const TextStyle(color: AppTheme.sanggamGold, fontWeight: FontWeight.bold)),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(member.name, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                        Text(member.relation, style: theme.textTheme.bodySmall),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(color: statusColor.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
                    child: Text(
                      member.status == 'normal' ? '정상' : member.status == 'warning' ? '주의' : '위험',
                      style: TextStyle(fontSize: 11, color: statusColor, fontWeight: FontWeight.w600),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              // 7일 트렌드 바 (간략)
              Row(
                children: List.generate(7, (i) {
                  final val = member.weeklyTrend[i];
                  return Expanded(
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 1),
                      child: Column(
                        children: [
                          Container(
                            height: 24,
                            decoration: BoxDecoration(
                              color: val > 0.7 ? Colors.green : val > 0.4 ? Colors.orange : Colors.red,
                              borderRadius: BorderRadius.circular(4),
                            ),
                          ),
                          const SizedBox(height: 2),
                          Text(['월', '화', '수', '목', '금', '토', '일'][i], style: const TextStyle(fontSize: 9)),
                        ],
                      ),
                    ),
                  );
                }),
              ),
              const SizedBox(height: 4),
              Text('최근 측정: ${member.lastMeasured}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ],
          ),
        ),
      ),
    );
  }

  static final _familyMembers = [
    _FamilyMemberData(id: '1', name: '어머니', relation: '부모', status: 'warning', lastMeasured: '2시간 전', weeklyTrend: [0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.5]),
    _FamilyMemberData(id: '2', name: '아버지', relation: '부모', status: 'normal', lastMeasured: '4시간 전', weeklyTrend: [0.7, 0.8, 0.9, 0.8, 0.7, 0.8, 0.9]),
    _FamilyMemberData(id: '3', name: '배우자', relation: '배우자', status: 'normal', lastMeasured: '1일 전', weeklyTrend: [0.9, 0.8, 0.7, 0.8, 0.9, 0.8, 0.7]),
  ];
}

class _FamilyMemberData {
  final String id, name, relation, status, lastMeasured;
  final List<double> weeklyTrend;
  const _FamilyMemberData({required this.id, required this.name, required this.relation, required this.status, required this.lastMeasured, required this.weeklyTrend});
}
