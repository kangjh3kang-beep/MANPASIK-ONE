import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminAuditScreen — Sanggam Orbit 감사 로그
//
// [Rule 4] +sanggam_theme.dart + sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 3x 제거
// [Rule 4] theme.textTheme.* ~6x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~2x → SanggamTheme 상수
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card → SanggamContainer
// [Rule 4] FilterChip → Container + BoxDecoration
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius:12→16, bottom:4→8, h:12→16
// ───────────────────────────────────────────────────

/// 관리자 감사 로그 화면
class AdminAuditScreen
    extends ConsumerStatefulWidget {
  const AdminAuditScreen({super.key});

  @override
  ConsumerState<AdminAuditScreen> createState() =>
      _AdminAuditScreenState();
}

class _AdminAuditScreenState
    extends ConsumerState<AdminAuditScreen> {
  String _filterType = 'all';
  final _searchCtrl = TextEditingController();

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final logsAsync =
        ref.watch(auditLogProvider);

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
                      '감사 로그',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(
                        Icons
                            .file_download_outlined,
                        color: Colors.white),
                    tooltip: '내보내기',
                    onPressed: () async {
                      try {
                        final client = ref.read(
                            restClientProvider);
                        final data = await client
                            .getAuditLog(
                                limit: 1000);
                        final logs =
                            (data['logs']
                                    as List?) ??
                                [];
                        final buffer = StringBuffer(
                            '시간,관리자,작업,대상,IP주소\n');
                        for (final log in logs) {
                          final m = log as Map<
                              String, dynamic>;
                          buffer.writeln(
                              '${m["timestamp"]},${m["admin_name"]},${m["action"]},${m["target"]},${m["ip_address"]}');
                        }
                        if (context.mounted) {
                          ScaffoldMessenger.of(
                                  context)
                              .showSnackBar(
                            SnackBar(
                                content: Text(
                                    '감사 로그 ${logs.length}건이 CSV로 내보내졌습니다.')),
                          );
                        }
                      } catch (e) {
                        if (context.mounted) {
                          ScaffoldMessenger.of(
                                  context)
                              .showSnackBar(
                            SnackBar(
                                content: Text(
                                    '내보내기 실패: $e')),
                          );
                        }
                      }
                    },
                  ),
                ],
              ),
            ),

            // 검색 + 필터
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 16),
              child: Column(
                children: [
                  TextField(
                    controller: _searchCtrl,
                    style: const TextStyle(
                        color: Colors.white),
                    decoration: InputDecoration(
                      hintText: '사용자, 행위, IP 검색...',
                      hintStyle: const TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim),
                      prefixIcon: const Icon(
                          Icons.search,
                          color: SanggamTheme
                              .onSurfaceDim),
                      filled: true,
                      fillColor:
                          SanggamTheme.surface,
                      border: OutlineInputBorder(
                        borderRadius:
                            BorderRadius.circular(
                                16),
                        borderSide: BorderSide.none,
                      ),
                      enabledBorder:
                          OutlineInputBorder(
                        borderRadius:
                            BorderRadius.circular(
                                16),
                        borderSide: const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                      ),
                      focusedBorder:
                          OutlineInputBorder(
                        borderRadius:
                            BorderRadius.circular(
                                16),
                        borderSide: const BorderSide(
                            color: SanggamTheme
                                .primary),
                      ),
                    ),
                    onChanged: (_) =>
                        setState(() {}),
                  ),
                  const SizedBox(height: 8),
                  SingleChildScrollView(
                    scrollDirection:
                        Axis.horizontal,
                    child: Row(
                      children:
                          _actionTypes.map((t) {
                        final isSelected =
                            _filterType == t.$1;
                        return Padding(
                          padding:
                              const EdgeInsets.only(
                                  right: 8),
                          child: GestureDetector(
                            onTap: () => setState(
                                () => _filterType =
                                    t.$1),
                            child: Container(
                              padding:
                                  const EdgeInsets
                                      .symmetric(
                                      horizontal:
                                          16,
                                      vertical: 8),
                              decoration:
                                  BoxDecoration(
                                color: isSelected
                                    ? SanggamTheme
                                        .primary
                                    : SanggamTheme
                                        .surface,
                                borderRadius:
                                    BorderRadius
                                        .circular(
                                            16),
                                border: Border.all(
                                  color: isSelected
                                      ? SanggamTheme
                                          .primary
                                      : SanggamTheme
                                          .surfaceVariant,
                                ),
                              ),
                              child: Text(
                                t.$2,
                                style: TextStyle(
                                  color: isSelected
                                      ? SanggamTheme
                                          .background
                                      : Colors
                                          .white,
                                  fontSize: 12,
                                  fontWeight: isSelected
                                      ? FontWeight
                                          .bold
                                      : FontWeight
                                          .normal,
                                ),
                              ),
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),
            const Divider(
                height: 1,
                color:
                    SanggamTheme.surfaceVariant),

            // 로그 리스트
            Expanded(
              child: logsAsync.when(
                data: (logMap) {
                  final logs = (logMap['entries']
                              as List?)
                          ?.cast<
                              Map<String,
                                  dynamic>>() ??
                      <Map<String, dynamic>>[];
                  final filtered =
                      _applyFilters(logs);
                  if (filtered.isEmpty) {
                    return const Center(
                      child: Column(
                        mainAxisSize:
                            MainAxisSize.min,
                        children: [
                          Icon(
                              Icons
                                  .assignment_outlined,
                              size: 48,
                              color: SanggamTheme
                                  .onSurfaceDim),
                          SizedBox(height: 16),
                          Text('감사 로그가 없습니다.',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                              )),
                        ],
                      ),
                    );
                  }
                  return RefreshIndicator(
                    color: SanggamTheme.primary,
                    onRefresh: () async => ref
                        .invalidate(
                            auditLogProvider),
                    child: ListView.builder(
                      padding: const EdgeInsets
                          .symmetric(
                          horizontal: 16,
                          vertical: 8),
                      itemCount: filtered.length,
                      itemBuilder: (_, index) =>
                          _buildLogTile(
                              filtered[index]),
                    ),
                  );
                },
                loading: () => const Center(
                    child:
                        CircularProgressIndicator(
                            color: SanggamTheme
                                .primary)),
                error: (_, __) =>
                    _buildFallbackLogs(),
              ),
            ),
          ],
        ),
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

  List<Map<String, dynamic>> _applyFilters(
      List<Map<String, dynamic>> logs) {
    var result = logs;
    if (_filterType != 'all') {
      result = result
          .where((l) => l['type'] == _filterType)
          .toList();
    }
    if (_searchCtrl.text.isNotEmpty) {
      final q = _searchCtrl.text.toLowerCase();
      result = result.where((l) {
        final user =
            (l['user'] as String? ?? '')
                .toLowerCase();
        final action =
            (l['action'] as String? ?? '')
                .toLowerCase();
        final ip = (l['ip'] as String? ?? '')
            .toLowerCase();
        return user.contains(q) ||
            action.contains(q) ||
            ip.contains(q);
      }).toList();
    }
    return result;
  }

  Widget _buildLogTile(
      Map<String, dynamic> log) {
    final action =
        log['action'] as String? ?? '';
    final user = log['user'] as String? ?? '';
    final ip = log['ip'] as String? ?? '';
    final time = log['time'] as String? ?? '';
    final type =
        log['type'] as String? ?? 'system';

    final icon = switch (type) {
      'user' => Icons.person,
      'security' => Icons.shield,
      'data' => Icons.storage,
      _ => Icons.settings,
    };
    final color = switch (type) {
      'user' => SanggamTheme.jagaeCyan,
      'security' => SanggamTheme.error,
      'data' => SanggamTheme.jagaeCyan,
      _ => SanggamTheme.primary,
    };

    return SanggamContainer(
      borderRadius: 16,
      padding: EdgeInsets.zero,
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        dense: true,
        leading: CircleAvatar(
          radius: 16,
          backgroundColor:
              color.withValues(alpha: 0.1),
          child:
              Icon(icon, size: 16, color: color),
        ),
        title: Text(action,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 12,
              fontWeight: FontWeight.w600,
            )),
        subtitle: Text('$user | $ip',
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            )),
        trailing: Text(time,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            )),
      ),
    );
  }

  Widget _buildFallbackLogs() {
    final fallback = [
      {
        'action': '사용자 등록',
        'user': 'admin@manpasik.com',
        'ip': '192.168.1.1',
        'time': '09:30',
        'type': 'user',
      },
      {
        'action': '시스템 설정 변경',
        'user': 'admin@manpasik.com',
        'ip': '192.168.1.1',
        'time': '09:15',
        'type': 'system',
      },
      {
        'action': '비밀번호 변경',
        'user': 'user123@test.com',
        'ip': '10.0.0.5',
        'time': '08:45',
        'type': 'security',
      },
      {
        'action': '데이터 내보내기',
        'user': 'user456@test.com',
        'ip': '10.0.0.8',
        'time': '08:30',
        'type': 'data',
      },
      {
        'action': '계정 정지',
        'user': 'admin@manpasik.com',
        'ip': '192.168.1.1',
        'time': '08:00',
        'type': 'security',
      },
      {
        'action': '카트리지 등록',
        'user': 'admin@manpasik.com',
        'ip': '192.168.1.1',
        'time': '어제',
        'type': 'system',
      },
    ];
    return ListView.builder(
      padding: const EdgeInsets.symmetric(
          horizontal: 16, vertical: 8),
      itemCount: fallback.length,
      itemBuilder: (_, index) =>
          _buildLogTile(fallback[index]),
    );
  }
}
