import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

// ───────────────────────────────────────────────────
// AdminUsersScreen — Sanggam Orbit 사용자 관리
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] +sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData 파라미터 1x 제거
// [Rule 4] theme.textTheme.* ~5x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~5x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] FilterChip 5x → Container + BoxDecoration
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.purple → SanggamTheme.jagaeMagenta
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 4] AlertDialog/TextField/FilledButton 다크 테마
// [Rule 2] borderRadius:12→16, spacing:6→8
// ───────────────────────────────────────────────────

/// 관리자 사용자 관리 화면
class AdminUsersScreen
    extends ConsumerStatefulWidget {
  const AdminUsersScreen({super.key});

  @override
  ConsumerState<AdminUsersScreen> createState() =>
      _AdminUsersScreenState();
}

class _AdminUsersScreenState
    extends ConsumerState<AdminUsersScreen> {
  final _searchController =
      TextEditingController();
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
      final client =
          ref.read(restClientProvider);
      final resp = await client.adminListUsers(
        query: _searchController.text.isEmpty
            ? null
            : _searchController.text,
      );
      final items =
          resp['users'] as List? ?? [];
      setState(() {
        _users =
            items.cast<Map<String, dynamic>>();
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
    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '사용자 관리',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // 검색바
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 16),
              child: TextField(
                controller: _searchController,
                style: const TextStyle(
                    color: Colors.white),
                decoration: InputDecoration(
                  hintText: '이름, 이메일로 검색',
                  hintStyle: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim),
                  prefixIcon: const Icon(
                      Icons.search,
                      color: SanggamTheme
                          .onSurfaceDim),
                  suffixIcon: IconButton(
                    tooltip: '검색',
                    icon: const Icon(
                        Icons.search,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onPressed: _loadUsers,
                  ),
                  filled: true,
                  fillColor:
                      SanggamTheme.surface,
                  border: OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        BorderSide.none,
                  ),
                  enabledBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                  ),
                  focusedBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .primary),
                  ),
                ),
                onSubmitted: (_) =>
                    _loadUsers(),
              ),
            ),
            const SizedBox(height: 8),

            // 역할 필터
            SizedBox(
              height: 36,
              child: ListView(
                scrollDirection:
                    Axis.horizontal,
                padding:
                    const EdgeInsets.symmetric(
                        horizontal: 16),
                children: _roleFilters.map((r) {
                  final isSelected =
                      _roleFilter == r.$1;
                  return Padding(
                    padding:
                        const EdgeInsets.only(
                            right: 8),
                    child: GestureDetector(
                      onTap: () {
                        setState(() =>
                            _roleFilter =
                                r.$1);
                        _loadUsers();
                      },
                      child: Container(
                        padding: const EdgeInsets
                            .symmetric(
                            horizontal: 16,
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
                          r.$2,
                          style: TextStyle(
                            color: isSelected
                                ? SanggamTheme
                                    .background
                                : Colors.white,
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
            const SizedBox(height: 8),
            const Divider(
                height: 1,
                color: SanggamTheme
                    .surfaceVariant),

            // 사용자 카운트
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8),
              child: Row(
                children: [
                  Text(
                      '총 ${_users.length}명',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight:
                            FontWeight.w600,
                      )),
                  const Spacer(),
                  const Text('최근 가입순',
                      style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      )),
                ],
              ),
            ),

            // 사용자 리스트
            Expanded(
              child: _isLoading
                  ? const Center(
                      child:
                          CircularProgressIndicator(
                              color:
                                  SanggamTheme
                                      .primary))
                  : _users.isEmpty
                      ? const Center(
                          child: Column(
                            mainAxisSize:
                                MainAxisSize
                                    .min,
                            children: [
                              Icon(
                                  Icons
                                      .person_search,
                                  size: 48,
                                  color: SanggamTheme
                                      .onSurfaceDim),
                              SizedBox(
                                  height: 8),
                              Text(
                                  '사용자를 찾을 수 없습니다.',
                                  style:
                                      TextStyle(
                                    color: Colors
                                        .white,
                                    fontSize: 14,
                                  )),
                            ],
                          ),
                        )
                      : RefreshIndicator(
                          color: SanggamTheme
                              .primary,
                          onRefresh:
                              _loadUsers,
                          child: ListView
                              .separated(
                            padding:
                                const EdgeInsets
                                    .all(16),
                            itemCount:
                                _users.length,
                            separatorBuilder:
                                (_, __) =>
                                    const Divider(
                                        height:
                                            1,
                                        color: SanggamTheme
                                            .surfaceVariant),
                            itemBuilder: (_,
                                    index) =>
                                _buildUserTile(
                                    _users[
                                        index]),
                          ),
                        ),
            ),
          ],
        ),
      ),
    );
  }

  static const _roleFilters = [
    ('all', '전체'),
    ('user', '일반'),
    ('premium', '프리미엄'),
    ('clinician', '임상의'),
    ('admin', '관리자'),
  ];

  Widget _buildUserTile(
      Map<String, dynamic> user) {
    final name = user['display_name']
            as String? ??
        user['name'] as String? ??
        '사용자';
    final email =
        user['email'] as String? ?? '';
    final userId = user['user_id']
            as String? ??
        user['id'] as String? ??
        email;
    final role =
        user['role'] as String? ?? 'user';
    final status =
        user['status'] as String? ?? 'active';
    final createdAt =
        user['created_at'] as String? ?? '';
    final isActive = status == 'active';

    final roleColor = switch (role) {
      'admin' => SanggamTheme.error,
      'clinician' => SanggamTheme.jagaeMagenta,
      'premium' => SanggamTheme.primary,
      _ => SanggamTheme.jagaeCyan,
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
        backgroundColor: isActive
            ? SanggamTheme.primary
                .withValues(alpha: 0.15)
            : SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.2),
        child: Text(
          name.isNotEmpty ? name[0] : '?',
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: isActive
                ? SanggamTheme.primary
                : SanggamTheme.onSurfaceDim,
          ),
        ),
      ),
      title: Row(
        children: [
          Flexible(
              child: Text(name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                  ),
                  overflow: TextOverflow
                      .ellipsis)),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(
                horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: roleColor.withValues(
                  alpha: 0.1),
              borderRadius:
                  BorderRadius.circular(8),
            ),
            child: Text(roleLabel,
                style: TextStyle(
                  fontSize: 10,
                  color: roleColor,
                  fontWeight: FontWeight.w600,
                )),
          ),
          if (!isActive) ...[
            const SizedBox(width: 4),
            Container(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 2),
              decoration: BoxDecoration(
                color: SanggamTheme.error
                    .withValues(alpha: 0.1),
                borderRadius:
                    BorderRadius.circular(8),
              ),
              child: const Text('정지',
                  style: TextStyle(
                    fontSize: 10,
                    color: SanggamTheme.error,
                    fontWeight: FontWeight.w600,
                  )),
            ),
          ],
        ],
      ),
      subtitle: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Text(email,
              style: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim,
                fontSize: 12,
              )),
          if (createdAt.isNotEmpty)
            Text(
              '가입일: ${createdAt.length > 10 ? createdAt.substring(0, 10) : createdAt}',
              style: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim,
                fontSize: 12,
              ),
            ),
        ],
      ),
      trailing: PopupMenuButton<String>(
        tooltip: '사용자 작업',
        color: SanggamTheme.surface,
        itemBuilder: (_) => [
          const PopupMenuItem(
              value: 'detail',
              child: Text('상세 보기',
                  style: TextStyle(
                      color: Colors.white))),
          const PopupMenuItem(
              value: 'role',
              child: Text('역할 변경',
                  style: TextStyle(
                      color: Colors.white))),
          if (isActive)
            const PopupMenuItem(
                value: 'suspend',
                child: Text('계정 정지',
                    style: TextStyle(
                        color: Colors.white)))
          else
            const PopupMenuItem(
                value: 'activate',
                child: Text('계정 활성화',
                    style: TextStyle(
                        color: Colors.white))),
        ],
        onSelected: (action) {
          switch (action) {
            case 'role':
              _showRoleChangeDialog(
                  context, userId, name);
              break;
            case 'suspend':
            case 'activate':
              _showConfirmAction(
                  context,
                  userId,
                  name,
                  action);
              break;
            default:
              ScaffoldMessenger.of(context)
                  .showSnackBar(
                SnackBar(
                    content: Text(
                        '$name 사용자 상세 보기')),
              );
          }
        },
      ),
    );
  }

  void _showRoleChangeDialog(
      BuildContext context,
      String userId,
      String name) {
    String selectedRole = 'user';
    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setDialogState) =>
            AlertDialog(
          backgroundColor:
              SanggamTheme.surface,
          shape: RoundedRectangleBorder(
              borderRadius:
                  BorderRadius.circular(16)),
          title: Text('$name 역할 변경',
              style: const TextStyle(
                  color: Colors.white)),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              'user',
              'premium',
              'clinician',
              'admin'
            ]
                .map((role) =>
                    RadioListTile<String>(
                      title: Text(role,
                          style:
                              const TextStyle(
                                  color: Colors
                                      .white)),
                      value: role,
                      groupValue: selectedRole,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (v) =>
                          setDialogState(() =>
                              selectedRole =
                                  v!),
                    ))
                .toList(),
          ),
          actions: [
            TextButton(
                onPressed: () =>
                    Navigator.pop(ctx),
                child: const Text('취소',
                    style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim))),
            FilledButton(
              style: FilledButton.styleFrom(
                backgroundColor:
                    SanggamTheme.primary,
                foregroundColor:
                    SanggamTheme.background,
                shape: RoundedRectangleBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16)),
              ),
              onPressed: () async {
                final client = ref
                    .read(restClientProvider);
                try {
                  await client.adminChangeRole(
                      userId, selectedRole);
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                  }
                  if (context.mounted) {
                    ScaffoldMessenger.of(
                            context)
                        .showSnackBar(
                      SnackBar(
                          content: Text(
                              '$name의 역할이 $selectedRole로 변경되었습니다.')),
                    );
                  }
                } catch (e) {
                  if (ctx.mounted) {
                    ScaffoldMessenger.of(ctx)
                        .showSnackBar(
                      SnackBar(
                          content: Text(
                              '역할 변경 실패: $e')),
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

  void _showConfirmAction(
      BuildContext context,
      String userId,
      String name,
      String action) {
    final isActivate = action == 'activate';
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: Text(
            isActivate ? '계정 활성화' : '계정 정지',
            style: const TextStyle(
                color: Colors.white)),
        content: Text(
            '$name 사용자를 ${isActivate ? "활성화" : "정지"}하시겠습니까?',
            style: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim)),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              child: const Text('취소',
                  style: TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim))),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor: isActivate
                  ? SanggamTheme.primary
                  : SanggamTheme.error,
              foregroundColor:
                  SanggamTheme.background,
              shape: RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16)),
            ),
            onPressed: () async {
              final client =
                  ref.read(restClientProvider);
              try {
                await client.adminBulkAction(
                    userIds: [userId],
                    action: action);
                if (ctx.mounted) {
                  Navigator.pop(ctx);
                }
                if (context.mounted) {
                  ScaffoldMessenger.of(context)
                      .showSnackBar(
                    SnackBar(
                        content: Text(
                            '$name 계정이 ${isActivate ? "활성화" : "정지"}되었습니다.')),
                  );
                }
              } catch (e) {
                if (ctx.mounted) {
                  ScaffoldMessenger.of(ctx)
                      .showSnackBar(
                    SnackBar(
                        content: Text(
                            '작업 실패: $e')),
                  );
                }
              }
            },
            child: Text(
                isActivate ? '활성화' : '정지'),
          ),
        ],
      ),
    );
  }

  static final _fallbackUsers = [
    {
      'display_name': '김건강',
      'email': 'kim@test.com',
      'role': 'admin',
      'status': 'active',
      'created_at': '2025-12-01T00:00:00Z',
    },
    {
      'display_name': '이안심',
      'email': 'lee@test.com',
      'role': 'clinician',
      'status': 'active',
      'created_at': '2026-01-05T00:00:00Z',
    },
    {
      'display_name': '박자연',
      'email': 'park@test.com',
      'role': 'premium',
      'status': 'active',
      'created_at': '2026-01-15T00:00:00Z',
    },
    {
      'display_name': '최희망',
      'email': 'choi@test.com',
      'role': 'user',
      'status': 'active',
      'created_at': '2026-02-01T00:00:00Z',
    },
    {
      'display_name': '정미래',
      'email': 'jung@test.com',
      'role': 'user',
      'status': 'suspended',
      'created_at': '2026-02-10T00:00:00Z',
    },
    {
      'display_name': '한솔루션',
      'email': 'han@test.com',
      'role': 'premium',
      'status': 'active',
      'created_at': '2026-02-12T00:00:00Z',
    },
  ];
}
