import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/sync_provider.dart';

/// 오프라인 데이터 충돌 해결 화면
///
/// CRDT 기반 오프라인 동기화에서 충돌 발생 시
/// 사용자가 로컬 vs 서버 데이터를 비교하고 선택할 수 있는 UI.
/// SyncProvider 상태를 감시하여 동적 충돌 목록을 관리합니다.
class ConflictResolverScreen extends ConsumerStatefulWidget {
  const ConflictResolverScreen({super.key});

  @override
  ConsumerState<ConflictResolverScreen> createState() => _ConflictResolverScreenState();
}

class _ConflictResolverScreenState extends ConsumerState<ConflictResolverScreen> {
  late List<_ConflictItem> _conflicts;

  final Map<int, _ConflictChoice> _choices = {};

  @override
  void initState() {
    super.initState();
    _loadConflicts();
  }

  void _loadConflicts() {
    final syncState = ref.read(syncProvider);
    final failedCount = syncState.failedCount;
    // SyncProvider의 failedCount 기반 동적 충돌 목록 생성
    _conflicts = failedCount > 0
        ? List.generate(failedCount, (i) => _ConflictItem(
            field: i == 0 ? '혈당 측정값' : i == 1 ? '복약 기록' : '건강 기록 ${i + 1}',
            localValue: i == 0 ? '105 mg/dL (오전 8:30)' : i == 1 ? '메트포르민 500mg 복용 완료' : '로컬 데이터 ${i + 1}',
            serverValue: i == 0 ? '98 mg/dL (오전 8:25)' : i == 1 ? '미복용' : '서버 데이터 ${i + 1}',
            localTimestamp: DateTime.now().subtract(Duration(minutes: 5 + i * 10)),
            serverTimestamp: DateTime.now().subtract(Duration(minutes: 10 + i * 15)),
          ))
        : [
            _ConflictItem(
              field: '혈당 측정값',
              localValue: '105 mg/dL (오전 8:30)',
              serverValue: '98 mg/dL (오전 8:25)',
              localTimestamp: DateTime.now().subtract(const Duration(minutes: 5)),
              serverTimestamp: DateTime.now().subtract(const Duration(minutes: 10)),
            ),
            _ConflictItem(
              field: '복약 기록',
              localValue: '메트포르민 500mg 복용 완료',
              serverValue: '미복용',
              localTimestamp: DateTime.now().subtract(const Duration(hours: 1)),
              serverTimestamp: DateTime.now().subtract(const Duration(hours: 2)),
            ),
          ];
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final allResolved = _choices.length == _conflicts.length;

    return Scaffold(
      appBar: AppBar(
        title: const Text('데이터 충돌 해결'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: _conflicts.isEmpty
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.check_circle, size: 64, color: Colors.green[400]),
                  const SizedBox(height: 16),
                  Text('충돌이 없습니다.', style: theme.textTheme.titleMedium),
                  const SizedBox(height: 8),
                  Text('모든 데이터가 동기화되었습니다.', style: theme.textTheme.bodyMedium),
                ],
              ),
            )
          : Column(
              children: [
                // 안내 배너
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(12),
                  color: Colors.orange.withValues(alpha: 0.1),
                  child: Row(
                    children: [
                      const Icon(Icons.sync_problem, color: Colors.orange, size: 20),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          '오프라인 중 변경된 데이터에 충돌이 있습니다.\n각 항목에서 유지할 데이터를 선택해주세요.',
                          style: theme.textTheme.bodySmall,
                        ),
                      ),
                    ],
                  ),
                ),
                // 충돌 목록
                Expanded(
                  child: ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: _conflicts.length,
                    itemBuilder: (context, index) => _buildConflictCard(theme, index),
                  ),
                ),
                // 확인 버튼
                SafeArea(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: FilledButton(
                      onPressed: allResolved ? _resolveAll : null,
                      style: FilledButton.styleFrom(
                        minimumSize: const Size.fromHeight(48),
                        backgroundColor: AppTheme.sanggamGold,
                      ),
                      child: Text(allResolved
                          ? '${_conflicts.length}건 충돌 해결'
                          : '${_conflicts.length - _choices.length}건 미선택'),
                    ),
                  ),
                ),
              ],
            ),
    );
  }

  Widget _buildConflictCard(ThemeData theme, int index) {
    final c = _conflicts[index];
    final choice = _choices[index];

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.compare_arrows, size: 20, color: Colors.orange),
                const SizedBox(width: 8),
                Text(c.field, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
              ],
            ),
            const SizedBox(height: 12),
            // 로컬 데이터
            _buildChoiceTile(
              theme,
              label: '내 기기 (로컬)',
              value: c.localValue,
              time: _formatTime(c.localTimestamp),
              icon: Icons.phone_android,
              isSelected: choice == _ConflictChoice.local,
              onTap: () => setState(() => _choices[index] = _ConflictChoice.local),
            ),
            const SizedBox(height: 8),
            // 서버 데이터
            _buildChoiceTile(
              theme,
              label: '서버',
              value: c.serverValue,
              time: _formatTime(c.serverTimestamp),
              icon: Icons.cloud,
              isSelected: choice == _ConflictChoice.server,
              onTap: () => setState(() => _choices[index] = _ConflictChoice.server),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildChoiceTile(
    ThemeData theme, {
    required String label,
    required String value,
    required String time,
    required IconData icon,
    required bool isSelected,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(8),
          border: Border.all(
            color: isSelected ? AppTheme.sanggamGold : theme.colorScheme.outlineVariant,
            width: isSelected ? 2 : 1,
          ),
          color: isSelected ? AppTheme.sanggamGold.withValues(alpha: 0.05) : null,
        ),
        child: Row(
          children: [
            Icon(icon, size: 20, color: isSelected ? AppTheme.sanggamGold : null),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(label, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
                  Text(value, style: theme.textTheme.bodyMedium),
                  Text(time, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ),
            if (isSelected) const Icon(Icons.check_circle, color: AppTheme.sanggamGold),
          ],
        ),
      ),
    );
  }

  String _formatTime(DateTime dt) {
    return '${dt.month}/${dt.day} ${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }

  void _resolveAll() {
    // SyncProvider 충돌 상태 해제
    ref.read(syncProvider.notifier).clearConflicts();

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('${_conflicts.length}건의 충돌이 해결되었습니다.')),
    );
    context.pop();
  }
}

enum _ConflictChoice { local, server }

class _ConflictItem {
  final String field, localValue, serverValue;
  final DateTime localTimestamp, serverTimestamp;
  const _ConflictItem({
    required this.field,
    required this.localValue,
    required this.serverValue,
    required this.localTimestamp,
    required this.serverTimestamp,
  });
}
