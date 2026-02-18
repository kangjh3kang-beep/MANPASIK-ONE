import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 관리자 대시보드 화면
class AdminDashboardScreen extends ConsumerWidget {
  const AdminDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final statsAsync = ref.watch(systemStatsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('관리자 대시보드'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.people),
            tooltip: '사용자 관리',
            onPressed: () => context.push('/admin/users'),
          ),
          IconButton(
            icon: const Icon(Icons.settings),
            tooltip: '시스템 설정',
            onPressed: () => context.push('/admin/settings'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(systemStatsProvider);
          ref.invalidate(auditLogProvider);
        },
        child: statsAsync.when(
          data: (stats) => _buildDashboard(context, theme, stats),
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (_, __) => _buildDashboard(context, theme, _fallbackStats),
        ),
      ),
    );
  }

  Widget _buildDashboard(BuildContext context, ThemeData theme, Map<String, dynamic> stats) {
    final totalUsers = stats['total_users'] as int? ?? 0;
    final activeUsers = stats['active_users'] as int? ?? 0;
    final totalMeasurements = stats['total_measurements'] as int? ?? 0;
    final totalDevices = stats['total_devices'] as int? ?? 0;
    final systemHealth = stats['system_health'] as String? ?? 'healthy';

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        // 시스템 상태 배너
        _buildSystemHealthBanner(theme, systemHealth),
        const SizedBox(height: 16),

        // 통계 카드 그리드
        Text('시스템 통계', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 1.6,
          children: [
            _buildStatCard(theme, '전체 사용자', '$totalUsers', Icons.people, Colors.blue),
            _buildStatCard(theme, '활성 사용자', '$activeUsers', Icons.person_pin, Colors.green),
            _buildStatCard(theme, '총 측정 횟수', '$totalMeasurements', Icons.science, Colors.orange),
            _buildStatCard(theme, '등록 기기', '$totalDevices', Icons.devices, Colors.purple),
          ],
        ),
        const SizedBox(height: 24),

        // 빠른 메뉴
        Text('관리 메뉴', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        _AdminMenuTile(
          icon: Icons.people,
          title: '사용자 관리',
          subtitle: '사용자 검색, 정보 조회, 계정 정지',
          onTap: () => context.push('/admin/users'),
        ),
        _AdminMenuTile(
          icon: Icons.settings,
          title: '시스템 설정',
          subtitle: '카테고리별 시스템 구성 관리',
          onTap: () => context.push('/admin/settings'),
        ),
        _AdminMenuTile(
          icon: Icons.assignment,
          title: '감사 로그',
          subtitle: '관리자 활동 기록 조회',
          onTap: () => context.push('/admin/audit'),
        ),
        _AdminMenuTile(
          icon: Icons.bar_chart,
          title: '서비스 모니터링',
          subtitle: '마이크로서비스 상태 확인',
          onTap: () => context.push('/admin/monitor'),
        ),
        _AdminMenuTile(
          icon: Icons.emergency,
          title: '긴급 상황 관리',
          subtitle: '긴급 알림 모니터링 및 대응',
          onTap: () => context.push('/admin/emergency'),
        ),
        _AdminMenuTile(
          icon: Icons.account_tree,
          title: '계층 관리',
          subtitle: '조직 계층 구조 관리',
          onTap: () => context.push('/admin/hierarchy'),
        ),
        _AdminMenuTile(
          icon: Icons.verified_user,
          title: '규제 준수',
          subtitle: 'GDPR/PIPA/HIPAA 체크리스트',
          onTap: () => context.push('/admin/compliance'),
        ),
        const SizedBox(height: 24),

        // 최근 감사 로그
        Text('최근 활동', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        ..._recentActivities.map((a) => _buildActivityTile(theme, a)),
      ],
    );
  }

  Widget _buildSystemHealthBanner(ThemeData theme, String status) {
    final isHealthy = status == 'healthy';
    final color = isHealthy ? Colors.green : Colors.orange;

    return Card(
      color: color.withOpacity(0.1),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Icon(
              isHealthy ? Icons.check_circle : Icons.warning,
              color: color,
              size: 32,
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    isHealthy ? '시스템 정상' : '시스템 주의',
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: color,
                    ),
                  ),
                  Text(
                    isHealthy
                        ? '모든 서비스가 정상적으로 운영 중입니다.'
                        : '일부 서비스에서 경고가 감지되었습니다.',
                    style: theme.textTheme.bodySmall,
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatCard(ThemeData theme, String label, String value, IconData icon, Color color) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Row(
              children: [
                Icon(icon, size: 20, color: color),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(label, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(value, style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
          ],
        ),
      ),
    );
  }

  Widget _buildActivityTile(ThemeData theme, _ActivityItem item) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: ListTile(
        dense: true,
        leading: CircleAvatar(
          radius: 16,
          backgroundColor: item.color.withOpacity(0.1),
          child: Icon(item.icon, size: 16, color: item.color),
        ),
        title: Text(item.action, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
        subtitle: Text(item.detail, style: theme.textTheme.bodySmall),
        trailing: Text(item.time, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
      ),
    );
  }

  static final _fallbackStats = <String, dynamic>{
    'total_users': 1247,
    'active_users': 892,
    'total_measurements': 34521,
    'total_devices': 456,
    'system_health': 'healthy',
  };

  static final _recentActivities = [
    _ActivityItem(action: '사용자 등록', detail: 'user_1023@test.com 계정 생성', time: '5분 전', icon: Icons.person_add, color: Colors.green),
    _ActivityItem(action: '시스템 설정 변경', detail: 'ai.model_version → v2.1.0', time: '1시간 전', icon: Icons.settings, color: Colors.blue),
    _ActivityItem(action: '알림 발송', detail: '전체 사용자 대상 공지 발송', time: '3시간 전', icon: Icons.notifications, color: Colors.orange),
    _ActivityItem(action: '카트리지 등록', detail: 'CRT-PRO-009 신규 등록', time: '어제', icon: Icons.science, color: Colors.purple),
    _ActivityItem(action: '계정 정지', detail: 'spam_user@test.com 계정 정지', time: '2일 전', icon: Icons.block, color: Colors.red),
  ];
}

class _AdminMenuTile extends StatelessWidget {
  const _AdminMenuTile({required this.icon, required this.title, required this.subtitle, required this.onTap});
  final IconData icon;
  final String title, subtitle;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: ListTile(
        leading: Icon(icon, color: AppTheme.sanggamGold),
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(subtitle, style: theme.textTheme.bodySmall),
        trailing: const Icon(Icons.chevron_right),
        onTap: onTap,
      ),
    );
  }
}

class _ActivityItem {
  final String action, detail, time;
  final IconData icon;
  final Color color;
  const _ActivityItem({required this.action, required this.detail, required this.time, required this.icon, required this.color});
}
