import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 관리자 사용자 관리 화면
class AdminUsersScreen extends ConsumerStatefulWidget {
  const AdminUsersScreen({super.key});

  @override
  ConsumerState<AdminUsersScreen> createState() => _AdminUsersScreenState();
}

class _AdminUsersScreenState extends ConsumerState<AdminUsersScreen> {
  final _searchController = TextEditingController();
  List<Map<String, dynamic>> _users = [];
  bool _isLoading = false;
  String _roleFilter = 'all';

  @override
  void initState() {
    super.initState();
    _loadUsers();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _loadUsers() async {
    setState(() => _isLoading = true);
    try {
      final client = ref.read(restClientProvider);
      final resp = await client.adminListUsers(
        query: _searchController.text.isEmpty ? null : _searchController.text,
      );
      final items = resp['users'] as List? ?? [];
      setState(() {
        _users = items.cast<Map<String, dynamic>>();
        _isLoading = false;
      });
    } catch (_) {
      setState(() {
        _users = _fallbackUsers;
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('사용자 관리'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // 검색바
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '이름, 이메일로 검색',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: IconButton(
                  icon: const Icon(Icons.search),
                  onPressed: _loadUsers,
                ),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
              onSubmitted: (_) => _loadUsers(),
            ),
          ),

          // 역할 필터
          SizedBox(
            height: 36,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: [
                ('all', '전체'),
                ('user', '일반'),
                ('premium', '프리미엄'),
                ('clinician', '임상의'),
                ('admin', '관리자'),
              ].map((r) {
                final isSelected = _roleFilter == r.$1;
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: FilterChip(
                    selected: isSelected,
                    label: Text(r.$2, style: const TextStyle(fontSize: 12)),
                    selectedColor: AppTheme.sanggamGold,
                    onSelected: (_) {
                      setState(() => _roleFilter = r.$1);
                      _loadUsers();
                    },
                  ),
                );
              }).toList(),
            ),
          ),
          const SizedBox(height: 8),
          const Divider(height: 1),

          // 사용자 카운트
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                Text('총 ${_users.length}명', style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
                const Spacer(),
                Text('최근 가입순', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ],
            ),
          ),

          // 사용자 리스트
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _users.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.person_search, size: 48, color: theme.colorScheme.onSurfaceVariant),
                            const SizedBox(height: 8),
                            Text('사용자를 찾을 수 없습니다.', style: theme.textTheme.bodyMedium),
                          ],
                        ),
                      )
                    : RefreshIndicator(
                        onRefresh: _loadUsers,
                        child: ListView.separated(
                          padding: const EdgeInsets.all(16),
                          itemCount: _users.length,
                          separatorBuilder: (_, __) => const Divider(height: 1),
                          itemBuilder: (context, index) => _buildUserTile(theme, _users[index]),
                        ),
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildUserTile(ThemeData theme, Map<String, dynamic> user) {
    final name = user['display_name'] as String? ?? user['name'] as String? ?? '사용자';
    final email = user['email'] as String? ?? '';
    final userId = user['user_id'] as String? ?? user['id'] as String? ?? email;
    final role = user['role'] as String? ?? 'user';
    final status = user['status'] as String? ?? 'active';
    final createdAt = user['created_at'] as String? ?? '';
    final isActive = status == 'active';

    final roleColor = switch (role) {
      'admin' => Colors.red,
      'clinician' => Colors.purple,
      'premium' => Colors.orange,
      _ => Colors.blue,
    };

    final roleLabel = switch (role) {
      'admin' => '관리자',
      'clinician' => '임상의',
      'premium' => '프리미엄',
      _ => '일반',
    };

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: isActive ? theme.colorScheme.primaryContainer : Colors.grey.withOpacity(0.2),
        child: Text(
          name.isNotEmpty ? name[0] : '?',
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: isActive ? theme.colorScheme.onPrimaryContainer : Colors.grey,
          ),
        ),
      ),
      title: Row(
        children: [
          Flexible(child: Text(name, style: const TextStyle(fontWeight: FontWeight.w600), overflow: TextOverflow.ellipsis)),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
            decoration: BoxDecoration(
              color: roleColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Text(roleLabel, style: TextStyle(fontSize: 10, color: roleColor, fontWeight: FontWeight.w600)),
          ),
          if (!isActive) ...[
            const SizedBox(width: 4),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
              decoration: BoxDecoration(
                color: Colors.red.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: const Text('정지', style: TextStyle(fontSize: 10, color: Colors.red, fontWeight: FontWeight.w600)),
            ),
          ],
        ],
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(email, style: theme.textTheme.bodySmall),
          if (createdAt.isNotEmpty)
            Text(
              '가입일: ${createdAt.length > 10 ? createdAt.substring(0, 10) : createdAt}',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
            ),
        ],
      ),
      trailing: PopupMenuButton<String>(
        itemBuilder: (_) => [
          const PopupMenuItem(value: 'detail', child: Text('상세 보기')),
          const PopupMenuItem(value: 'role', child: Text('역할 변경')),
          if (isActive)
            const PopupMenuItem(value: 'suspend', child: Text('계정 정지'))
          else
            const PopupMenuItem(value: 'activate', child: Text('계정 활성화')),
        ],
        onSelected: (action) {
          switch (action) {
            case 'role':
              _showRoleChangeDialog(context, userId, name);
              break;
            case 'suspend':
            case 'activate':
              _showConfirmAction(context, userId, name, action);
              break;
            default:
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('$name 사용자 상세 보기')),
              );
          }
        },
      ),
    );
  }

  void _showRoleChangeDialog(BuildContext context, String userId, String name) {
    String selectedRole = 'user';
    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setDialogState) => AlertDialog(
          title: Text('$name 역할 변경'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: ['user', 'premium', 'clinician', 'admin'].map((role) =>
              RadioListTile<String>(
                title: Text(role),
                value: role,
                groupValue: selectedRole,
                onChanged: (v) => setDialogState(() => selectedRole = v!),
              ),
            ).toList(),
          ),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
            FilledButton(
              onPressed: () async {
                final client = ref.read(restClientProvider);
                try {
                  await client.adminChangeRole(userId, selectedRole);
                  if (ctx.mounted) Navigator.pop(ctx);
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('$name의 역할이 $selectedRole로 변경되었습니다.')),
                    );
                  }
                } catch (e) {
                  if (ctx.mounted) {
                    ScaffoldMessenger.of(ctx).showSnackBar(
                      SnackBar(content: Text('역할 변경 실패: $e')),
                    );
                  }
                }
              },
              child: const Text('변경'),
            ),
          ],
        ),
      ),
    );
  }

  void _showConfirmAction(BuildContext context, String userId, String name, String action) {
    final isActivate = action == 'activate';
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(isActivate ? '계정 활성화' : '계정 정지'),
        content: Text('$name 사용자를 ${isActivate ? "활성화" : "정지"}하시겠습니까?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            style: isActivate ? null : FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () async {
              final client = ref.read(restClientProvider);
              try {
                await client.adminBulkAction(userIds: [userId], action: action);
                if (ctx.mounted) Navigator.pop(ctx);
                if (context.mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('$name 계정이 ${isActivate ? "활성화" : "정지"}되었습니다.')),
                  );
                }
              } catch (e) {
                if (ctx.mounted) {
                  ScaffoldMessenger.of(ctx).showSnackBar(
                    SnackBar(content: Text('작업 실패: $e')),
                  );
                }
              }
            },
            child: Text(isActivate ? '활성화' : '정지'),
          ),
        ],
      ),
    );
  }

  static final _fallbackUsers = [
    {'display_name': '김건강', 'email': 'kim@test.com', 'role': 'admin', 'status': 'active', 'created_at': '2025-12-01T00:00:00Z'},
    {'display_name': '이안심', 'email': 'lee@test.com', 'role': 'clinician', 'status': 'active', 'created_at': '2026-01-05T00:00:00Z'},
    {'display_name': '박자연', 'email': 'park@test.com', 'role': 'premium', 'status': 'active', 'created_at': '2026-01-15T00:00:00Z'},
    {'display_name': '최희망', 'email': 'choi@test.com', 'role': 'user', 'status': 'active', 'created_at': '2026-02-01T00:00:00Z'},
    {'display_name': '정미래', 'email': 'jung@test.com', 'role': 'user', 'status': 'suspended', 'created_at': '2026-02-10T00:00:00Z'},
    {'display_name': '한솔루션', 'email': 'han@test.com', 'role': 'premium', 'status': 'active', 'created_at': '2026-02-12T00:00:00Z'},
  ];
}
