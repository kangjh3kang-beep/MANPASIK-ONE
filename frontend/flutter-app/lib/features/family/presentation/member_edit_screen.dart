import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

// ───────────────────────────────────────────────────
// MemberEditScreen — Sanggam Orbit 멤버 설정
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.* 4x → 직접 TextStyle
// [Rule 4] Colors.red 4x → SanggamTheme.error
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 가족 멤버 역할/모드 편집 화면
class MemberEditScreen extends ConsumerStatefulWidget {
  const MemberEditScreen({super.key, required this.memberId});

  final String memberId;

  @override
  ConsumerState<MemberEditScreen> createState() =>
      _MemberEditScreenState();
}

class _MemberEditScreenState
    extends ConsumerState<MemberEditScreen> {
  String _role = 'member';
  String _viewMode = 'normal';
  bool _allowDataSharing = true;
  bool _receiveAlerts = true;
  bool _isSaving = false;
  String _memberName = '';

  @override
  void initState() {
    super.initState();
    _loadMember();
  }

  Future<void> _loadMember() async {
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final data = await rest.listFamilyMembers(userId);
      final members = (data['members'] as List?) ?? [];
      final member = members.cast<Map<String, dynamic>>().firstWhere(
        (m) => m['id'] == widget.memberId,
        orElse: () => <String, dynamic>{},
      );
      if (mounted && member.isNotEmpty) {
        setState(() {
          _memberName = (member['name'] as String?) ?? '멤버 ${widget.memberId}';
          _role = (member['role'] as String?) ?? 'member';
          _viewMode = (member['view_mode'] as String?) ?? 'normal';
          _allowDataSharing = (member['allow_data_sharing'] as bool?) ?? true;
          _receiveAlerts = (member['receive_alerts'] as bool?) ?? true;
        });
      }
    } on DioException {
      // 로드 실패 시 기본값 유지
    }
  }

  Future<void> _saveMember() async {
    setState(() => _isSaving = true);
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.updateFamilyMember(
        groupId: userId,
        memberId: widget.memberId,
        role: _role,
        mode: _viewMode,
        permissions: {
          'allow_data_sharing': _allowDataSharing,
          'receive_alerts': _receiveAlerts,
        },
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('설정이 저장되었습니다.')),
        );
        context.pop();
      }
    } on DioException {
      if (mounted) {
        setState(() => _isSaving = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('저장에 실패했습니다. 다시 시도해주세요.')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
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
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '멤버 설정',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  TextButton(
                    onPressed: _isSaving ? null : _saveMember,
                    child: _isSaving
                        ? const SizedBox(
                            width: 16,
                            height: 16,
                            child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: SanggamTheme.primary))
                        : const Text(
                            '저장',
                            style: TextStyle(
                              color: SanggamTheme.primary,
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
                padding: const EdgeInsets.all(16),
                children: [
                  // 프로필 헤더
                  Center(
                    child: Column(
                      children: [
                        CircleAvatar(
                          radius: 40,
                          backgroundColor: SanggamTheme
                              .primary
                              .withValues(alpha: 0.2),
                          child: const Icon(Icons.person,
                              size: 40,
                              color: SanggamTheme.primary),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          _memberName.isEmpty ? '멤버 ${widget.memberId}' : _memberName,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 24),

                  // 역할 설정
                  const Text(
                    '역할',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SegmentedButton<String>(
                    segments: const [
                      ButtonSegment(
                          value: 'admin',
                          label: Text('관리자'),
                          icon: Icon(
                              Icons.admin_panel_settings)),
                      ButtonSegment(
                          value: 'guardian',
                          label: Text('보호자'),
                          icon: Icon(Icons.shield)),
                      ButtonSegment(
                          value: 'member',
                          label: Text('일반'),
                          icon: Icon(Icons.person)),
                    ],
                    selected: {_role},
                    onSelectionChanged: (s) =>
                        setState(() => _role = s.first),
                  ),
                  const SizedBox(height: 24),

                  // 보기 모드
                  const Text(
                    '보기 모드',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SegmentedButton<String>(
                    segments: const [
                      ButtonSegment(
                          value: 'normal',
                          label: Text('일반')),
                      ButtonSegment(
                          value: 'senior',
                          label: Text('시니어')),
                      ButtonSegment(
                          value: 'child',
                          label: Text('어린이')),
                    ],
                    selected: {_viewMode},
                    onSelectionChanged: (s) => setState(
                        () => _viewMode = s.first),
                  ),
                  const SizedBox(height: 24),

                  // 권한 설정
                  const Text(
                    '권한',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SwitchListTile(
                    title: const Text('데이터 공유 허용'),
                    subtitle:
                        const Text('측정 결과를 그룹 멤버에게 공유'),
                    value: _allowDataSharing,
                    onChanged: (v) => setState(
                        () => _allowDataSharing = v),
                  ),
                  SwitchListTile(
                    title: const Text('이상치 알림 수신'),
                    subtitle:
                        const Text('이 멤버의 이상 수치 발생 시 알림'),
                    value: _receiveAlerts,
                    onChanged: (v) => setState(
                        () => _receiveAlerts = v),
                  ),
                  const SizedBox(height: 24),

                  // 멤버 내보내기
                  OutlinedButton.icon(
                    onPressed: () {
                      showDialog(
                        context: context,
                        builder: (ctx) => AlertDialog(
                          title: const Text('멤버 내보내기'),
                          content: const Text(
                              '이 멤버를 가족 그룹에서 내보내시겠습니까?'),
                          actions: [
                            TextButton(
                              onPressed: () =>
                                  Navigator.pop(ctx),
                              child: const Text('취소'),
                            ),
                            FilledButton(
                              onPressed: () {
                                Navigator.pop(ctx);
                                context.pop();
                              },
                              style: FilledButton.styleFrom(
                                backgroundColor:
                                    SanggamTheme.error,
                              ),
                              child: const Text('내보내기'),
                            ),
                          ],
                        ),
                      );
                    },
                    icon: const Icon(Icons.person_remove,
                        color: SanggamTheme.error),
                    label: const Text(
                      '멤버 내보내기',
                      style: TextStyle(
                          color: SanggamTheme.error),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
