import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminHierarchyScreen — Sanggam Orbit 계층 관리
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* → 직접 TextStyle
// [Rule 4] theme.colorScheme.* → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] bottom:4→8
// ───────────────────────────────────────────────────

/// 계층형 관리 화면
class AdminHierarchyScreen extends ConsumerWidget {
  const AdminHierarchyScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '조직 계층 관리',
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
              child: ListView(
                padding:
                    const EdgeInsets.all(16),
                children: [
                  const Text('조직 구조',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      )),
                  const SizedBox(height: 8),
                  _buildNode(0, 'ManPaSik 본사',
                      'organization', [
                    _buildNode(
                        1, '운영팀', 'team', [
                      _buildNode(2, '시스템 관리자',
                          'super_admin', []),
                      _buildNode(2, '고객 지원',
                          'admin', []),
                    ]),
                    _buildNode(
                        1, '의료팀', 'team', [
                      _buildNode(2, '의료 감독',
                          'medical_admin', []),
                      _buildNode(2, '데이터 분석',
                          'analyst', []),
                    ]),
                    _buildNode(
                        1, '파트너 기관', 'team', [
                      _buildNode(2, '서울대병원',
                          'partner', [
                        _buildNode(
                            3,
                            '김의사 (내과)',
                            'doctor',
                            []),
                        _buildNode(
                            3,
                            '이의사 (심장내과)',
                            'doctor',
                            []),
                      ]),
                      _buildNode(2, '연세세브란스',
                          'partner', [
                        _buildNode(
                            3,
                            '박의사 (내분비)',
                            'doctor',
                            []),
                      ]),
                    ]),
                  ]),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNode(int depth, String label,
      String role, List<Widget> children) {
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
      'organization' => SanggamTheme.primary,
      'team' => SanggamTheme.jagaeCyan,
      'super_admin' => SanggamTheme.error,
      'admin' => SanggamTheme.primary,
      'partner' => SanggamTheme.jagaeCyan,
      _ => SanggamTheme.onSurfaceDim,
    };

    return Column(
      crossAxisAlignment:
          CrossAxisAlignment.start,
      children: [
        Padding(
          padding:
              EdgeInsets.only(left: depth * 24.0),
          child: SanggamContainer(
            borderRadius: 16,
            padding: EdgeInsets.zero,
            margin:
                const EdgeInsets.only(bottom: 8),
            child: ListTile(
              dense: true,
              leading:
                  Icon(icon, color: color, size: 20),
              title: Text(label,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w500,
                  )),
              trailing: Text(
                  role.replaceAll('_', ' '),
                  style: const TextStyle(
                    color:
                        SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  )),
            ),
          ),
        ),
        ...children,
      ],
    );
  }
}
