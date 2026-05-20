import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

// ───────────────────────────────────────────────────
// SecurityScreen — Sanggam Orbit 보안 설정
//
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] ThemeData params 3x 제거
// [Rule 4] theme.textTheme ~5x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~3x → SanggamTheme 상수
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] AlertDialog/BottomSheet/TextField 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 보안 설정 화면
class SecurityScreen
    extends ConsumerStatefulWidget {
  const SecurityScreen({super.key});

  @override
  ConsumerState<SecurityScreen>
      createState() =>
          _SecurityScreenState();
}

class _SecurityScreenState
    extends ConsumerState<SecurityScreen> {
  bool _biometricEnabled = false;
  bool _twoFactorEnabled = false;

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
                      '보안',
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
            // 본문
            Expanded(
              child: ListView(
                children: [
                  // 비밀번호 변경
                  _buildSectionHeader('인증'),
                  ListTile(
                    leading: const Icon(
                        Icons.lock_outline,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '비밀번호 변경',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        '마지막 변경: 30일 이상',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () =>
                        _showChangePasswordDialog(
                            context),
                  ),
                  SwitchListTile(
                    secondary: const Icon(
                        Icons.fingerprint,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '생체인증',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        '지문 또는 얼굴 인식으로 로그인',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    activeColor:
                        SanggamTheme.primary,
                    value: _biometricEnabled,
                    onChanged: (v) {
                      setState(() =>
                          _biometricEnabled =
                              v);
                      ScaffoldMessenger.of(
                              context)
                          .showSnackBar(
                        SnackBar(
                            content: Text(v
                                ? '생체인증이 활성화되었습니다.'
                                : '생체인증이 비활성화되었습니다.')),
                      );
                    },
                  ),
                  SwitchListTile(
                    secondary: const Icon(
                        Icons.security,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '2단계 인증',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        'SMS 또는 인증 앱으로 추가 보안',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    activeColor:
                        SanggamTheme.primary,
                    value: _twoFactorEnabled,
                    onChanged: (v) {
                      setState(() =>
                          _twoFactorEnabled =
                              v);
                      ScaffoldMessenger.of(
                              context)
                          .showSnackBar(
                        SnackBar(
                            content: Text(v
                                ? '2단계 인증이 활성화되었습니다.'
                                : '2단계 인증이 비활성화되었습니다.')),
                      );
                    },
                  ),
                  const Divider(
                      color: SanggamTheme
                          .surfaceVariant),

                  // 로그인 기록
                  _buildSectionHeader('활동'),
                  ListTile(
                    leading: const Icon(
                        Icons.history,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '로그인 기록',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        '최근 로그인 활동 확인',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () =>
                        _showLoginHistorySheet(
                            context),
                  ),
                  ListTile(
                    leading: const Icon(
                        Icons.devices,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '활성 세션',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        '현재 로그인된 기기 관리',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () =>
                        _showActiveSessionsSheet(
                            context),
                  ),
                  const Divider(
                      color: SanggamTheme
                          .surfaceVariant),

                  // 데이터 보안
                  _buildSectionHeader(
                      '데이터 보안'),
                  ListTile(
                    leading: const Icon(
                        Icons
                            .download_outlined,
                        color: SanggamTheme
                            .onSurfaceDim),
                    title: const Text(
                        '내 데이터 다운로드',
                        style: TextStyle(
                            color: Colors
                                .white)),
                    subtitle: const Text(
                        '개인 데이터 전체 백업 (FHIR/CSV)',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () {
                      ScaffoldMessenger.of(
                              context)
                          .showSnackBar(
                        const SnackBar(
                            content: Text(
                                '데이터 다운로드 요청이 접수되었습니다. 이메일로 전송됩니다.')),
                      );
                    },
                  ),
                  ListTile(
                    leading: const Icon(
                        Icons.delete_forever,
                        color: SanggamTheme
                            .error),
                    title: const Text(
                        '계정 삭제',
                        style: TextStyle(
                            color:
                                SanggamTheme
                                    .error)),
                    subtitle: const Text(
                        '모든 데이터가 영구 삭제됩니다',
                        style: TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 12,
                        )),
                    trailing: const Icon(
                        Icons.chevron_right,
                        color: SanggamTheme
                            .onSurfaceDim),
                    onTap: () =>
                        _showDeleteAccountDialog(
                            context),
                  ),
                  const SizedBox(height: 32),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(
          16, 16, 16, 8),
      child: Text(
        title,
        style: const TextStyle(
          color: SanggamTheme.primary,
          fontWeight: FontWeight.bold,
          fontSize: 14,
        ),
      ),
    );
  }

  void _showChangePasswordDialog(
      BuildContext context) {
    final currentCtrl =
        TextEditingController();
    final newCtrl = TextEditingController();
    final confirmCtrl =
        TextEditingController();

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor:
            SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text('비밀번호 변경',
            style: TextStyle(
                color: Colors.white)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildDialogTextField(
                currentCtrl, '현재 비밀번호',
                obscure: true),
            const SizedBox(height: 8),
            _buildDialogTextField(
                newCtrl, '새 비밀번호',
                obscure: true),
            const SizedBox(height: 8),
            _buildDialogTextField(
                confirmCtrl, '새 비밀번호 확인',
                obscure: true),
          ],
        ),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              style: TextButton.styleFrom(
                  foregroundColor:
                      SanggamTheme
                          .onSurfaceDim),
              child: const Text('취소')),
          FilledButton(
            onPressed: () async {
              if (newCtrl.text ==
                      confirmCtrl.text &&
                  newCtrl.text.isNotEmpty &&
                  currentCtrl.text.isNotEmpty) {
                try {
                  final rest = ref.read(restClientProvider);
                  final userId = ref.read(authProvider).userId ?? '';
                  await rest.changePassword(
                    userId: userId,
                    currentPassword: currentCtrl.text,
                    newPassword: newCtrl.text,
                  );
                  if (ctx.mounted) Navigator.pop(ctx);
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('비밀번호가 변경되었습니다.')),
                    );
                  }
                } on DioException {
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('비밀번호 변경에 실패했습니다.')),
                    );
                  }
                }
              }
            },
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
            child: const Text('변경'),
          ),
        ],
      ),
    );
  }

  Widget _buildDialogTextField(
      TextEditingController ctrl,
      String label,
      {bool obscure = false}) {
    return TextField(
      controller: ctrl,
      obscureText: obscure,
      style: const TextStyle(
          color: Colors.white),
      decoration: InputDecoration(
        labelText: label,
        labelStyle: const TextStyle(
            color:
                SanggamTheme.onSurfaceDim),
        filled: true,
        fillColor:
            SanggamTheme.surfaceVariant,
        border: OutlineInputBorder(
          borderRadius:
              BorderRadius.circular(16),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius:
              BorderRadius.circular(16),
          borderSide: const BorderSide(
              color: SanggamTheme.primary),
        ),
      ),
    );
  }

  void _showLoginHistorySheet(
      BuildContext context) {
    final history = [
      ('오늘 09:30', '서울, 대한민국',
          'Chrome / Windows', true),
      ('어제 18:45', '서울, 대한민국',
          'ManPaSik App / Android', false),
      ('2일 전 14:20', '부산, 대한민국',
          'Safari / iOS', false),
      ('3일 전 08:15', '서울, 대한민국',
          'ManPaSik App / iOS', false),
    ];

    showModalBottomSheet(
      context: context,
      backgroundColor:
          SanggamTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)),
      ),
      builder: (ctx) => Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.all(16),
            child: Text('로그인 기록',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                )),
          ),
          ...history.map((h) => ListTile(
                leading: Icon(
                    h.$4
                        ? Icons.circle
                        : Icons
                            .circle_outlined,
                    size: 12,
                    color: h.$4
                        ? SanggamTheme
                            .jagaeCyan
                        : SanggamTheme
                            .onSurfaceDim),
                title: Text(h.$3,
                    style: const TextStyle(
                        color:
                            Colors.white,
                        fontSize: 14)),
                subtitle: Text(
                    '${h.$1} | ${h.$2}',
                    style: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim,
                      fontSize: 12,
                    )),
                trailing: h.$4
                    ? const Text('현재 세션',
                        style: TextStyle(
                          color: SanggamTheme
                              .jagaeCyan,
                          fontSize: 12,
                        ))
                    : null,
              )),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  void _showActiveSessionsSheet(
      BuildContext context) {
    showModalBottomSheet(
      context: context,
      backgroundColor:
          SanggamTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)),
      ),
      builder: (ctx) => Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          const Padding(
            padding: EdgeInsets.all(16),
            child: Text('활성 세션',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                )),
          ),
          ListTile(
            leading: const Icon(
                Icons.phone_android,
                color:
                    SanggamTheme.jagaeCyan),
            title: const Text('이 기기',
                style: TextStyle(
                    color: Colors.white)),
            subtitle: const Text(
                'ManPaSik App | 서울',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
            trailing: const Text('현재',
                style: TextStyle(
                  color:
                      SanggamTheme.jagaeCyan,
                  fontSize: 12,
                )),
          ),
          ListTile(
            leading: const Icon(
                Icons.computer,
                color: SanggamTheme
                    .onSurfaceDim),
            title: const Text(
                'Chrome / Windows',
                style: TextStyle(
                    color: Colors.white)),
            subtitle: const Text(
                '서울 | 오늘 09:30',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
            trailing: TextButton(
              onPressed: () {
                Navigator.pop(ctx);
                ScaffoldMessenger.of(context)
                    .showSnackBar(
                  const SnackBar(
                      content: Text(
                          '해당 세션이 종료되었습니다.')),
                );
              },
              style: TextButton.styleFrom(
                  foregroundColor:
                      SanggamTheme.error),
              child: const Text('로그아웃'),
            ),
          ),
          const SizedBox(height: 16),
        ],
      ),
    );
  }

  void _showDeleteAccountDialog(
      BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor:
            SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text('계정 삭제',
            style: TextStyle(
                color: Colors.white)),
        content: const Text(
          '정말 계정을 삭제하시겠습니까?\n\n'
          '이 작업은 되돌릴 수 없으며, 모든 건강 데이터와 측정 기록이 영구적으로 삭제됩니다.',
          style: TextStyle(
              color: SanggamTheme
                  .onSurfaceDim),
        ),
        actions: [
          TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              style: TextButton.styleFrom(
                  foregroundColor:
                      SanggamTheme
                          .onSurfaceDim),
              child: const Text('취소')),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.error,
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                  borderRadius:
                      BorderRadius.circular(
                          16)),
            ),
            onPressed: () {
              Navigator.pop(ctx);
              ScaffoldMessenger.of(context)
                  .showSnackBar(
                const SnackBar(
                    content: Text(
                        '계정 삭제 요청이 접수되었습니다. 14일 후 처리됩니다.')),
              );
            },
            child: const Text('삭제'),
          ),
        ],
      ),
    );
  }
}
