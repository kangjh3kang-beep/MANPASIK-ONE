import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 긴급 알림 상세 화면
class AlertDetailScreen extends ConsumerWidget {
  const AlertDetailScreen({super.key, required this.alertId});

  final String alertId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('긴급 알림 상세'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 알림 헤더
          Card(
            color: Colors.red.withOpacity(0.05),
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  const Icon(Icons.warning_amber_rounded, size: 48, color: Colors.red),
                  const SizedBox(height: 12),
                  Text('이상 수치 감지', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold, color: Colors.red)),
                  const SizedBox(height: 4),
                  Text('어머니 · 2024-02-15 14:32', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // 측정 수치
          Text('측정 수치', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  _buildMetricRow(theme, '수축기 혈압', '158 mmHg', '정상 범위: 90-120', Colors.red),
                  const Divider(),
                  _buildMetricRow(theme, '이완기 혈압', '95 mmHg', '정상 범위: 60-80', Colors.orange),
                  const Divider(),
                  _buildMetricRow(theme, '심박수', '92 bpm', '정상 범위: 60-100', Colors.green),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // AI 분석
          Text('AI 분석', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(Icons.auto_awesome, color: AppTheme.sanggamGold),
                      const SizedBox(width: 8),
                      Text('AI 건강 분석', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    '수축기 혈압이 정상 범위를 크게 초과했습니다. 최근 7일간의 추세를 볼 때 점진적 상승이 관찰됩니다. '
                    '스트레스, 식이 변화, 약물 복용 상태를 확인하시길 권장드립니다.',
                    style: theme.textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 8),
                  Text('※ 본 분석은 참고용이며, 정확한 진단은 의료 전문가와 상담하세요.',
                      style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // 대응 조치
          Text('권장 조치', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Card(
            child: Column(
              children: [
                ListTile(
                  leading: const Icon(Icons.local_hospital, color: Colors.blue),
                  title: const Text('가까운 병원 찾기'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/medical/facility-search'),
                ),
                const Divider(height: 1),
                ListTile(
                  leading: const Icon(Icons.videocam, color: Colors.green),
                  title: const Text('화상 진료 예약'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/medical/telemedicine'),
                ),
                const Divider(height: 1),
                ListTile(
                  leading: const Icon(Icons.phone, color: Colors.red),
                  title: const Text('긴급 연락처 전화'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/settings/emergency'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMetricRow(ThemeData theme, String label, String value, String range, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(width: 4, height: 32, decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(2))),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(label, style: theme.textTheme.bodySmall),
                Text(value, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              ],
            ),
          ),
          Text(range, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
      ),
    );
  }
}
