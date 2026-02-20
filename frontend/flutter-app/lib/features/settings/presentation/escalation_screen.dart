import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 119 에스컬레이션 (응급 상황 처리) UI
class EscalationScreen extends StatefulWidget {
  const EscalationScreen({super.key});

  @override
  State<EscalationScreen> createState() => _EscalationScreenState();
}

class _EscalationScreenState extends State<EscalationScreen> {
  bool _isEscalating = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('119 에스컬레이션'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 경고 배너
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: theme.colorScheme.error.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: theme.colorScheme.error),
              ),
              child: Column(
                children: [
                  Icon(Icons.warning_amber_rounded,
                      size: 48, color: theme.colorScheme.error),
                  const SizedBox(height: 8),
                  Text(
                    '위험 수치 감지',
                    style: theme.textTheme.titleLarge?.copyWith(
                      color: theme.colorScheme.error,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    '측정값이 설정된 위험 범위를 초과했습니다.',
                    style: theme.textTheme.bodyMedium,
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // 측정 정보 카드
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('측정 정보',
                        style: theme.textTheme.titleMedium
                            ?.copyWith(fontWeight: FontWeight.bold)),
                    const Divider(),
                    _buildInfoRow('측정 항목', '혈당'),
                    _buildInfoRow('측정값', '350 mg/dL'),
                    _buildInfoRow('정상 범위', '70 ~ 180 mg/dL'),
                    _buildInfoRow('측정 시각', '방금 전'),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 긴급 연락처 카드
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('긴급 연락처',
                        style: theme.textTheme.titleMedium
                            ?.copyWith(fontWeight: FontWeight.bold)),
                    const Divider(),
                    _buildContactTile(
                        '119 구급대', '119', Icons.local_hospital, true),
                    _buildContactTile(
                        '보호자 1 (배우자)', '010-1234-5678', Icons.person, false),
                    _buildContactTile(
                        '보호자 2 (자녀)', '010-9876-5432', Icons.person, false),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            // 119 신고 버튼
            FilledButton.icon(
              onPressed: _isEscalating ? null : _handleEscalate,
              icon: _isEscalating
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child:
                          CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                  : const Icon(Icons.phone),
              label: Text(_isEscalating ? '신고 중...' : '119 긴급 신고'),
              style: FilledButton.styleFrom(
                backgroundColor: theme.colorScheme.error,
                minimumSize: const Size.fromHeight(56),
                textStyle: const TextStyle(
                    fontSize: 18, fontWeight: FontWeight.bold),
              ),
            ),
            const SizedBox(height: 12),

            // 보호자 알림 버튼
            OutlinedButton.icon(
              onPressed: () => _notifyGuardians(context),
              icon: const Icon(Icons.notification_important),
              label: const Text('보호자에게 알림 보내기'),
              style: OutlinedButton.styleFrom(
                minimumSize: const Size.fromHeight(48),
              ),
            ),
            const SizedBox(height: 12),

            // 상황 종료 버튼
            TextButton(
              onPressed: () => context.pop(),
              child: const Text('상황 종료 (위험하지 않음)'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.grey)),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  Widget _buildContactTile(
      String name, String number, IconData icon, bool isPrimary) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: isPrimary ? Colors.red.shade100 : Colors.grey.shade200,
        child: Icon(icon, color: isPrimary ? Colors.red : Colors.grey),
      ),
      title: Text(name),
      subtitle: Text(number),
      trailing: IconButton(
        icon: const Icon(Icons.phone, color: AppTheme.sanggamGold),
        onPressed: () {},
      ),
    );
  }

  void _handleEscalate() {
    setState(() => _isEscalating = true);
    Future.delayed(const Duration(seconds: 2), () {
      if (mounted) {
        setState(() => _isEscalating = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('119 긴급 신고가 접수되었습니다.'),
            backgroundColor: Colors.red,
          ),
        );
      }
    });
  }

  void _notifyGuardians(BuildContext context) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('보호자에게 알림을 보냈습니다.')),
    );
  }
}
