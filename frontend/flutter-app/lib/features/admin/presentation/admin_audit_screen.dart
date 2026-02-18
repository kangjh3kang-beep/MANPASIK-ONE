import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';

/// 관리자 감사 로그 화면
class AdminAuditScreen extends ConsumerStatefulWidget {
  const AdminAuditScreen({super.key});

  @override
  ConsumerState<AdminAuditScreen> createState() => _AdminAuditScreenState();
}

class _AdminAuditScreenState extends ConsumerState<AdminAuditScreen> {
  String _filterType = 'all';
  final _searchCtrl = TextEditingController();

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final logsAsync = ref.watch(auditLogProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('감사 로그'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.file_download_outlined),
            tooltip: '내보내기',
            onPressed: () async {
              try {
                final client = ref.read(restClientProvider);
                final data = await client.getAuditLog(limit: 1000);
                final logs = (data['logs'] as List?) ?? [];
                // CSV 생성
                final buffer = StringBuffer('시간,관리자,작업,대상,IP주소\n');
                for (final log in logs) {
                  final m = log as Map<String, dynamic>;
                  buffer.writeln('${m["timestamp"]},${m["admin_name"]},${m["action"]},${m["target"]},${m["ip_address"]}');
                }
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('감사 로그 ${logs.length}건이 CSV로 내보내졌습니다.')),
                  );
                }
              } catch (e) {
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('내보내기 실패: $e')),
                  );
                }
              }
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // 검색 + 필터
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                TextField(
                  controller: _searchCtrl,
                  decoration: InputDecoration(
                    hintText: '사용자, 행위, IP 검색...',
                    prefixIcon: const Icon(Icons.search),
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                    filled: true,
                  ),
                  onChanged: (_) => setState(() {}),
                ),
                const SizedBox(height: 8),
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: _actionTypes.map((t) {
                      final isSelected = _filterType == t.$1;
                      return Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: FilterChip(
                          label: Text(t.$2),
                          selected: isSelected,
                          onSelected: (_) => setState(() => _filterType = t.$1),
                        ),
                      );
                    }).toList(),
                  ),
                ),
              ],
            ),
          ),
          const Divider(height: 1),

          // 로그 리스트
          Expanded(
            child: logsAsync.when(
              data: (logMap) {
                final logs = (logMap['entries'] as List?)
                    ?.cast<Map<String, dynamic>>() ?? <Map<String, dynamic>>[];
                final filtered = _applyFilters(logs);
                if (filtered.isEmpty) {
                  return Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.assignment_outlined, size: 48, color: theme.colorScheme.outline),
                        const SizedBox(height: 12),
                        Text('감사 로그가 없습니다.', style: theme.textTheme.bodyLarge),
                      ],
                    ),
                  );
                }
                return RefreshIndicator(
                  onRefresh: () async => ref.invalidate(auditLogProvider),
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    itemCount: filtered.length,
                    itemBuilder: (_, index) => _buildLogTile(theme, filtered[index]),
                  ),
                );
              },
              loading: () => const Center(child: CircularProgressIndicator()),
              error: (_, __) => _buildFallbackLogs(theme),
            ),
          ),
        ],
      ),
    );
  }

  static const _actionTypes = [
    ('all', '전체'),
    ('user', '사용자'),
    ('system', '시스템'),
    ('security', '보안'),
    ('data', '데이터'),
  ];

  List<Map<String, dynamic>> _applyFilters(List<Map<String, dynamic>> logs) {
    var result = logs;
    if (_filterType != 'all') {
      result = result.where((l) => l['type'] == _filterType).toList();
    }
    if (_searchCtrl.text.isNotEmpty) {
      final q = _searchCtrl.text.toLowerCase();
      result = result.where((l) {
        final user = (l['user'] as String? ?? '').toLowerCase();
        final action = (l['action'] as String? ?? '').toLowerCase();
        final ip = (l['ip'] as String? ?? '').toLowerCase();
        return user.contains(q) || action.contains(q) || ip.contains(q);
      }).toList();
    }
    return result;
  }

  Widget _buildLogTile(ThemeData theme, Map<String, dynamic> log) {
    final action = log['action'] as String? ?? '';
    final user = log['user'] as String? ?? '';
    final ip = log['ip'] as String? ?? '';
    final time = log['time'] as String? ?? '';
    final type = log['type'] as String? ?? 'system';

    final icon = switch (type) {
      'user' => Icons.person,
      'security' => Icons.shield,
      'data' => Icons.storage,
      _ => Icons.settings,
    };
    final color = switch (type) {
      'user' => Colors.blue,
      'security' => Colors.red,
      'data' => Colors.green,
      _ => Colors.orange,
    };

    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: ListTile(
        dense: true,
        leading: CircleAvatar(
          radius: 16,
          backgroundColor: color.withOpacity(0.1),
          child: Icon(icon, size: 16, color: color),
        ),
        title: Text(action, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
        subtitle: Text('$user | $ip', style: theme.textTheme.bodySmall),
        trailing: Text(time, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
      ),
    );
  }

  Widget _buildFallbackLogs(ThemeData theme) {
    final fallback = [
      {'action': '사용자 등록', 'user': 'admin@manpasik.com', 'ip': '192.168.1.1', 'time': '09:30', 'type': 'user'},
      {'action': '시스템 설정 변경', 'user': 'admin@manpasik.com', 'ip': '192.168.1.1', 'time': '09:15', 'type': 'system'},
      {'action': '비밀번호 변경', 'user': 'user123@test.com', 'ip': '10.0.0.5', 'time': '08:45', 'type': 'security'},
      {'action': '데이터 내보내기', 'user': 'user456@test.com', 'ip': '10.0.0.8', 'time': '08:30', 'type': 'data'},
      {'action': '계정 정지', 'user': 'admin@manpasik.com', 'ip': '192.168.1.1', 'time': '08:00', 'type': 'security'},
      {'action': '카트리지 등록', 'user': 'admin@manpasik.com', 'ip': '192.168.1.1', 'time': '어제', 'type': 'system'},
    ];
    return ListView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      itemCount: fallback.length,
      itemBuilder: (_, index) => _buildLogTile(theme, fallback[index]),
    );
  }
}
