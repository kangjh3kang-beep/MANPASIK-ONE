import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 가족 멤버 역할/모드 편집 화면
class MemberEditScreen extends ConsumerStatefulWidget {
  const MemberEditScreen({super.key, required this.memberId});

  final String memberId;

  @override
  ConsumerState<MemberEditScreen> createState() => _MemberEditScreenState();
}

class _MemberEditScreenState extends ConsumerState<MemberEditScreen> {
  String _role = 'member';
  String _viewMode = 'normal';
  bool _allowDataSharing = true;
  bool _receiveAlerts = true;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('멤버 설정'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
        actions: [
          TextButton(
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('설정이 저장되었습니다.')));
              context.pop();
            },
            child: const Text('저장'),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 프로필 헤더
          Center(
            child: Column(
              children: [
                CircleAvatar(radius: 40, backgroundColor: AppTheme.sanggamGold.withOpacity(0.2), child: const Icon(Icons.person, size: 40, color: AppTheme.sanggamGold)),
                const SizedBox(height: 8),
                Text('멤버 ${widget.memberId}', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // 역할 설정
          Text('역할', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          SegmentedButton<String>(
            segments: const [
              ButtonSegment(value: 'admin', label: Text('관리자'), icon: Icon(Icons.admin_panel_settings)),
              ButtonSegment(value: 'guardian', label: Text('보호자'), icon: Icon(Icons.shield)),
              ButtonSegment(value: 'member', label: Text('일반'), icon: Icon(Icons.person)),
            ],
            selected: {_role},
            onSelectionChanged: (s) => setState(() => _role = s.first),
          ),
          const SizedBox(height: 24),

          // 보기 모드
          Text('보기 모드', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          SegmentedButton<String>(
            segments: const [
              ButtonSegment(value: 'normal', label: Text('일반')),
              ButtonSegment(value: 'senior', label: Text('시니어')),
              ButtonSegment(value: 'child', label: Text('어린이')),
            ],
            selected: {_viewMode},
            onSelectionChanged: (s) => setState(() => _viewMode = s.first),
          ),
          const SizedBox(height: 24),

          // 권한 설정
          Text('권한', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          SwitchListTile(
            title: const Text('데이터 공유 허용'),
            subtitle: const Text('측정 결과를 그룹 멤버에게 공유'),
            value: _allowDataSharing,
            onChanged: (v) => setState(() => _allowDataSharing = v),
          ),
          SwitchListTile(
            title: const Text('이상치 알림 수신'),
            subtitle: const Text('이 멤버의 이상 수치 발생 시 알림'),
            value: _receiveAlerts,
            onChanged: (v) => setState(() => _receiveAlerts = v),
          ),
          const SizedBox(height: 24),

          // 멤버 내보내기
          OutlinedButton.icon(
            onPressed: () {
              showDialog(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('멤버 내보내기'),
                  content: const Text('이 멤버를 가족 그룹에서 내보내시겠습니까?'),
                  actions: [
                    TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
                    FilledButton(
                      onPressed: () { Navigator.pop(ctx); context.pop(); },
                      style: FilledButton.styleFrom(backgroundColor: Colors.red),
                      child: const Text('내보내기'),
                    ),
                  ],
                ),
              );
            },
            icon: const Icon(Icons.person_remove, color: Colors.red),
            label: const Text('멤버 내보내기', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}
