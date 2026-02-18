import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 진료 결과 화면
class ConsultationResultScreen extends ConsumerWidget {
  const ConsultationResultScreen({super.key, required this.consultationId});

  final String consultationId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('진료 결과'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
        actions: [
          IconButton(icon: const Icon(Icons.share), tooltip: 'PDF 내보내기', onPressed: () {
            ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('진료 결과 PDF가 생성되었습니다.')));
          }),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 진료 정보
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(Icons.medical_services, color: AppTheme.sanggamGold),
                      const SizedBox(width: 8),
                      Text('진료 정보', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 12),
                  _infoRow(theme, '진료일시', '2024-02-15 14:00'),
                  _infoRow(theme, '담당의', '김건강 전문의'),
                  _infoRow(theme, '진료과', '내과'),
                  _infoRow(theme, '진료유형', '화상진료'),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),

          // 진료 소견
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('진료 소견', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  Text(
                    '바이오마커 측정 결과 혈당 수치가 경계 범위에 있습니다. 식이 조절과 규칙적인 운동을 권장합니다. '
                    '2주 후 재측정하여 추이를 확인하시기 바랍니다.',
                    style: theme.textTheme.bodyMedium,
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 12),

          // 처방전
          Card(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                  child: Text('처방전', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                ),
                ListTile(
                  leading: const Icon(Icons.medication, color: Colors.blue),
                  title: const Text('메트포르민 500mg'),
                  subtitle: const Text('1일 2회, 식후 30분'),
                  trailing: const Text('14일분'),
                ),
                const Divider(height: 1),
                ListTile(
                  leading: const Icon(Icons.medication, color: Colors.green),
                  title: const Text('비타민 D 1000IU'),
                  subtitle: const Text('1일 1회, 아침'),
                  trailing: const Text('30일분'),
                ),
                Padding(
                  padding: const EdgeInsets.all(16),
                  child: OutlinedButton.icon(
                    onPressed: () => context.push('/medical/facility-search'),
                    icon: const Icon(Icons.local_pharmacy),
                    label: const Text('근처 약국 찾기'),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),

          // 다음 단계
          Card(
            child: Column(
              children: [
                ListTile(
                  leading: const Icon(Icons.calendar_today, color: AppTheme.sanggamGold),
                  title: const Text('다음 진료 예약'),
                  subtitle: const Text('2주 후 재진 권장'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.push('/medical/telemedicine'),
                ),
                const Divider(height: 1),
                ListTile(
                  leading: const Icon(Icons.science, color: AppTheme.sanggamGold),
                  title: const Text('바이오마커 재측정'),
                  subtitle: const Text('2주 후 추이 확인'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () => context.go('/measure'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _infoRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        children: [
          SizedBox(width: 80, child: Text(label, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant))),
          Text(value, style: theme.textTheme.bodyMedium),
        ],
      ),
    );
  }
}
