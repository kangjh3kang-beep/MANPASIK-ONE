import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 계층형 관리 화면
class AdminHierarchyScreen extends ConsumerWidget {
  const AdminHierarchyScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('조직 계층 관리'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Text('조직 구조', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          _buildNode(theme, 0, 'ManPaSik 본사', 'organization', [
            _buildNode(theme, 1, '운영팀', 'team', [
              _buildNode(theme, 2, '시스템 관리자', 'super_admin', []),
              _buildNode(theme, 2, '고객 지원', 'admin', []),
            ]),
            _buildNode(theme, 1, '의료팀', 'team', [
              _buildNode(theme, 2, '의료 감독', 'medical_admin', []),
              _buildNode(theme, 2, '데이터 분석', 'analyst', []),
            ]),
            _buildNode(theme, 1, '파트너 기관', 'team', [
              _buildNode(theme, 2, '서울대병원', 'partner', [
                _buildNode(theme, 3, '김의사 (내과)', 'doctor', []),
                _buildNode(theme, 3, '이의사 (심장내과)', 'doctor', []),
              ]),
              _buildNode(theme, 2, '연세세브란스', 'partner', [
                _buildNode(theme, 3, '박의사 (내분비)', 'doctor', []),
              ]),
            ]),
          ]),
        ],
      ),
    );
  }

  Widget _buildNode(ThemeData theme, int depth, String label, String role, List<Widget> children) {
    final icon = switch (role) {
      'organization' => Icons.business,
      'team' => Icons.group,
      'super_admin' => Icons.admin_panel_settings,
      'admin' => Icons.manage_accounts,
      'medical_admin' => Icons.medical_services,
      'analyst' => Icons.analytics,
      'partner' => Icons.local_hospital,
      'doctor' => Icons.person,
      _ => Icons.person,
    };

    final color = switch (role) {
      'organization' => AppTheme.sanggamGold,
      'team' => Colors.blue,
      'super_admin' => Colors.red,
      'admin' => Colors.orange,
      'partner' => Colors.green,
      _ => Colors.grey,
    };

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: EdgeInsets.only(left: depth * 24.0),
          child: Card(
            margin: const EdgeInsets.only(bottom: 4),
            child: ListTile(
              dense: true,
              leading: Icon(icon, color: color, size: 20),
              title: Text(label, style: const TextStyle(fontWeight: FontWeight.w500)),
              trailing: Text(role.replaceAll('_', ' '), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ),
          ),
        ),
        ...children,
      ],
    );
  }
}
