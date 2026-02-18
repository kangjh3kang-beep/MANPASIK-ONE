import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 시스템 모니터링 화면
class AdminMonitorScreen extends ConsumerWidget {
  const AdminMonitorScreen({super.key, this.tab});

  final String? tab;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return DefaultTabController(
      length: 3,
      initialIndex: tab == 'emergency' ? 2 : 0,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('시스템 모니터링'),
          leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
          bottom: const TabBar(
            tabs: [Tab(text: '서비스'), Tab(text: '리소스'), Tab(text: '긴급')],
          ),
        ),
        body: TabBarView(
          children: [
            _ServiceStatusTab(theme: theme),
            _ResourceTab(theme: theme),
            _EmergencyTab(theme: theme),
          ],
        ),
      ),
    );
  }
}

class _ServiceStatusTab extends StatelessWidget {
  const _ServiceStatusTab({required this.theme});
  final ThemeData theme;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        ..._services.map((s) => Card(
          margin: const EdgeInsets.only(bottom: 4),
          child: ListTile(
            leading: Icon(Icons.circle, size: 12, color: s.healthy ? Colors.green : Colors.red),
            title: Text(s.name),
            subtitle: Text('응답시간: ${s.latency}ms'),
            trailing: Text(s.healthy ? 'UP' : 'DOWN', style: TextStyle(color: s.healthy ? Colors.green : Colors.red, fontWeight: FontWeight.bold)),
          ),
        )),
      ],
    );
  }

  static final _services = [
    _ServiceInfo('Gateway', true, 12), _ServiceInfo('Auth', true, 8),
    _ServiceInfo('Measurement', true, 15), _ServiceInfo('User', true, 10),
    _ServiceInfo('Device', true, 11), _ServiceInfo('AI Inference', true, 45),
    _ServiceInfo('Notification', true, 9), _ServiceInfo('Family', true, 13),
    _ServiceInfo('Community', true, 14), _ServiceInfo('Telemedicine', true, 20),
  ];
}

class _ResourceTab extends StatelessWidget {
  const _ResourceTab({required this.theme});
  final ThemeData theme;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildMetric('CPU 사용률', 0.35, '35%', Colors.blue),
        _buildMetric('메모리 사용률', 0.62, '62%', Colors.orange),
        _buildMetric('디스크 사용률', 0.28, '28%', Colors.green),
        _buildMetric('네트워크 I/O', 0.15, '1.2 Gbps', Colors.purple),
        const SizedBox(height: 16),
        Text('요청 통계 (최근 1시간)', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                _statRow(theme, '총 요청 수', '12,450'),
                _statRow(theme, '평균 응답 시간', '23ms'),
                _statRow(theme, '에러율', '0.02%'),
                _statRow(theme, '활성 연결', '342'),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildMetric(String label, double value, String display, Color color) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [Text(label), Text(display, style: const TextStyle(fontWeight: FontWeight.bold))],
            ),
            const SizedBox(height: 8),
            LinearProgressIndicator(value: value, color: color, backgroundColor: color.withOpacity(0.1)),
          ],
        ),
      ),
    );
  }

  Widget _statRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [Text(label, style: theme.textTheme.bodyMedium), Text(value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold))],
      ),
    );
  }
}

class _EmergencyTab extends StatelessWidget {
  const _EmergencyTab({required this.theme});
  final ThemeData theme;

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          color: Colors.green.withOpacity(0.05),
          child: const Padding(
            padding: EdgeInsets.all(16),
            child: Row(
              children: [
                Icon(Icons.check_circle, color: Colors.green, size: 32),
                SizedBox(width: 12),
                Expanded(child: Text('현재 활성화된 긴급 알림이 없습니다.', style: TextStyle(fontWeight: FontWeight.w600))),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Text('최근 긴급 이벤트', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        Card(child: ListTile(
          leading: const Icon(Icons.warning, color: Colors.orange),
          title: const Text('AI 서비스 지연'),
          subtitle: const Text('2024-02-14 03:15 — 자동 복구됨'),
        )),
        Card(child: ListTile(
          leading: const Icon(Icons.error, color: Colors.red),
          title: const Text('DB 연결 끊김'),
          subtitle: const Text('2024-02-10 22:00 — 수동 복구'),
        )),
      ],
    );
  }
}

class _ServiceInfo {
  final String name;
  final bool healthy;
  final int latency;
  const _ServiceInfo(this.name, this.healthy, this.latency);
}
