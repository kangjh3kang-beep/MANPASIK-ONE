import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

/// 보안 설정 화면
class SecurityScreen extends ConsumerStatefulWidget {
  const SecurityScreen({super.key});

  @override
  ConsumerState<SecurityScreen> createState() => _SecurityScreenState();
}

class _SecurityScreenState extends ConsumerState<SecurityScreen> {
  bool _biometricEnabled = false;
  bool _twoFactorEnabled = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('보안'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: ListView(
        children: [
          // 비밀번호 변경
          _buildSectionHeader(theme, '인증'),
          ListTile(
            leading: const Icon(Icons.lock_outline),
            title: const Text('비밀번호 변경'),
            subtitle: const Text('마지막 변경: 30일 이상'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showChangePasswordDialog(context),
          ),
          SwitchListTile(
            secondary: const Icon(Icons.fingerprint),
            title: const Text('생체인증'),
            subtitle: const Text('지문 또는 얼굴 인식으로 로그인'),
            value: _biometricEnabled,
            onChanged: (v) {
              setState(() => _biometricEnabled = v);
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text(v ? '생체인증이 활성화되었습니다.' : '생체인증이 비활성화되었습니다.')),
              );
            },
          ),
          SwitchListTile(
            secondary: const Icon(Icons.security),
            title: const Text('2단계 인증'),
            subtitle: const Text('SMS 또는 인증 앱으로 추가 보안'),
            value: _twoFactorEnabled,
            onChanged: (v) {
              setState(() => _twoFactorEnabled = v);
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text(v ? '2단계 인증이 활성화되었습니다.' : '2단계 인증이 비활성화되었습니다.')),
              );
            },
          ),
          const Divider(),

          // 로그인 기록
          _buildSectionHeader(theme, '활동'),
          ListTile(
            leading: const Icon(Icons.history),
            title: const Text('로그인 기록'),
            subtitle: const Text('최근 로그인 활동 확인'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showLoginHistorySheet(context, theme),
          ),
          ListTile(
            leading: const Icon(Icons.devices),
            title: const Text('활성 세션'),
            subtitle: const Text('현재 로그인된 기기 관리'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showActiveSessionsSheet(context, theme),
          ),
          const Divider(),

          // 데이터 보안
          _buildSectionHeader(theme, '데이터 보안'),
          ListTile(
            leading: const Icon(Icons.download_outlined),
            title: const Text('내 데이터 다운로드'),
            subtitle: const Text('개인 데이터 전체 백업 (FHIR/CSV)'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('데이터 다운로드 요청이 접수되었습니다. 이메일로 전송됩니다.')),
              );
            },
          ),
          ListTile(
            leading: Icon(Icons.delete_forever, color: theme.colorScheme.error),
            title: Text('계정 삭제', style: TextStyle(color: theme.colorScheme.error)),
            subtitle: const Text('모든 데이터가 영구 삭제됩니다'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showDeleteAccountDialog(context),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(ThemeData theme, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: theme.textTheme.labelLarge?.copyWith(
          color: theme.colorScheme.primary,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  void _showChangePasswordDialog(BuildContext context) {
    final currentCtrl = TextEditingController();
    final newCtrl = TextEditingController();
    final confirmCtrl = TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('비밀번호 변경'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: currentCtrl,
              obscureText: true,
              decoration: const InputDecoration(labelText: '현재 비밀번호'),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: newCtrl,
              obscureText: true,
              decoration: const InputDecoration(labelText: '새 비밀번호'),
            ),
            const SizedBox(height: 8),
            TextField(
              controller: confirmCtrl,
              obscureText: true,
              decoration: const InputDecoration(labelText: '새 비밀번호 확인'),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () {
              if (newCtrl.text == confirmCtrl.text && newCtrl.text.isNotEmpty) {
                Navigator.pop(ctx);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('비밀번호가 변경되었습니다.')),
                );
              }
            },
            child: const Text('변경'),
          ),
        ],
      ),
    );
  }

  void _showLoginHistorySheet(BuildContext context, ThemeData theme) {
    final history = [
      ('오늘 09:30', '서울, 대한민국', 'Chrome / Windows', true),
      ('어제 18:45', '서울, 대한민국', 'ManPaSik App / Android', false),
      ('2일 전 14:20', '부산, 대한민국', 'Safari / iOS', false),
      ('3일 전 08:15', '서울, 대한민국', 'ManPaSik App / iOS', false),
    ];

    showModalBottomSheet(
      context: context,
      builder: (ctx) => Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Text('로그인 기록', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          ),
          ...history.map((h) => ListTile(
                leading: Icon(h.$4 ? Icons.circle : Icons.circle_outlined, size: 12, color: h.$4 ? Colors.green : Colors.grey),
                title: Text(h.$3, style: theme.textTheme.bodyMedium),
                subtitle: Text('${h.$1} | ${h.$2}'),
                trailing: h.$4 ? const Text('현재 세션', style: TextStyle(color: Colors.green, fontSize: 12)) : null,
              )),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  void _showActiveSessionsSheet(BuildContext context, ThemeData theme) {
    showModalBottomSheet(
      context: context,
      builder: (ctx) => Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Text('활성 세션', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          ),
          ListTile(
            leading: const Icon(Icons.phone_android, color: Colors.green),
            title: const Text('이 기기'),
            subtitle: const Text('ManPaSik App | 서울'),
            trailing: const Text('현재', style: TextStyle(color: Colors.green, fontSize: 12)),
          ),
          ListTile(
            leading: const Icon(Icons.computer),
            title: const Text('Chrome / Windows'),
            subtitle: const Text('서울 | 오늘 09:30'),
            trailing: TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('해당 세션이 종료되었습니다.')),
                );
              },
              child: const Text('로그아웃', style: TextStyle(color: Colors.red)),
            ),
          ),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  void _showDeleteAccountDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('계정 삭제'),
        content: const Text(
          '정말 계정을 삭제하시겠습니까?\n\n'
          '이 작업은 되돌릴 수 없으며, 모든 건강 데이터와 측정 기록이 영구적으로 삭제됩니다.',
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Theme.of(ctx).colorScheme.error),
            onPressed: () {
              Navigator.pop(ctx);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('계정 삭제 요청이 접수되었습니다. 14일 후 처리됩니다.')),
              );
            },
            child: const Text('삭제'),
          ),
        ],
      ),
    );
  }
}
