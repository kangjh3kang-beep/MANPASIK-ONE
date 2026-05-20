import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// NotificationSettingsScreen — Sanggam Orbit 알림 설정
//
// [Rule 4] app_theme + cosmic/hanji/jagae/porcelain → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + isDark 분기 전체 제거
// [Rule 4] AppTheme.sanggamGold/inkBlack → SanggamTheme 상수
// [Rule 4] CosmicBackground/HanjiBackground → SanggamTheme.background
// [Rule 4] JagaeContainer/PorcelainContainer → SanggamContainer
// [Rule 4] Colors.white10/54/70 → SanggamTheme 상수
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

class NotificationSettingsScreen
    extends ConsumerStatefulWidget {
  const NotificationSettingsScreen(
      {super.key});

  @override
  ConsumerState<NotificationSettingsScreen>
      createState() =>
          _NotificationSettingsScreenState();
}

class _NotificationSettingsScreenState
    extends ConsumerState<
        NotificationSettingsScreen> {
  bool _pushEnabled = true;
  bool _emailEnabled = false;
  bool _marketingEnabled = false;
  bool _healthAlerts = true;
  bool _deviceStatus = true;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadPreferences();
  }

  Future<void> _loadPreferences() async {
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final data = await rest.getNotificationPreferences(userId);
      if (mounted) {
        setState(() {
          _pushEnabled = (data['push_enabled'] as bool?) ?? true;
          _emailEnabled = (data['email_enabled'] as bool?) ?? false;
          _marketingEnabled = (data['marketing_enabled'] as bool?) ?? false;
          _healthAlerts = (data['health_alerts'] as bool?) ?? true;
          _deviceStatus = (data['device_status'] as bool?) ?? true;
          _isLoading = false;
        });
      }
    } on DioException {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Future<void> _savePreference(String key, bool value) async {
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.updateNotificationPreferences(userId, {key: value});
    } on DioException {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('설정 저장에 실패했습니다')),
        );
      }
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
                      '알림 설정',
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
                padding:
                    const EdgeInsets.all(24),
                children: [
                  _buildSectionHeader('필수 알림'),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        _buildSwitchTile(
                          title: '건강 측정 알림',
                          subtitle:
                              '비정상 파동 감지 시 즉시 알림',
                          value: _healthAlerts,
                          onChanged: (v) {
                              setState(() => _healthAlerts = v);
                              _savePreference('health_alerts', v);
                          },
                          icon: Icons
                              .monitor_heart_outlined,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSwitchTile(
                          title: '기기 상태 알림',
                          subtitle:
                              '배터리 부족, 연결 끊김 알림',
                          value: _deviceStatus,
                          onChanged: (v) {
                              setState(() => _deviceStatus = v);
                              _savePreference('device_status', v);
                          },
                          icon: Icons
                              .device_thermostat_outlined,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 32),

                  _buildSectionHeader('일반 알림'),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        _buildSwitchTile(
                          title: '푸시 알림',
                          subtitle:
                              '앱 내 활동 및 업데이트 알림',
                          value: _pushEnabled,
                          onChanged: (v) {
                              setState(() => _pushEnabled = v);
                              _savePreference('push_enabled', v);
                          },
                          icon: Icons
                              .notifications_active_outlined,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSwitchTile(
                          title: '이메일 알림',
                          subtitle:
                              '주간 건강 리포트 수신',
                          value: _emailEnabled,
                          onChanged: (v) {
                              setState(() => _emailEnabled = v);
                              _savePreference('email_enabled', v);
                          },
                          icon: Icons
                              .email_outlined,
                        ),
                        const Divider(
                            height: 1,
                            color: SanggamTheme
                                .surfaceVariant),
                        _buildSwitchTile(
                          title: '마케팅 정보 수신',
                          subtitle:
                              '이벤트 및 혜택 알림',
                          value:
                              _marketingEnabled,
                          onChanged: (v) {
                              setState(() => _marketingEnabled = v);
                              _savePreference('marketing_enabled', v);
                          },
                          icon: Icons
                              .redeem_outlined,
                        ),
                      ],
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

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(
          8, 0, 8, 8),
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

  Widget _buildSwitchTile({
    required String title,
    required String subtitle,
    required bool value,
    required ValueChanged<bool> onChanged,
    required IconData icon,
  }) {
    return SwitchListTile(
      value: value,
      onChanged: onChanged,
      activeColor: SanggamTheme.primary,
      activeTrackColor: SanggamTheme.primary
          .withValues(alpha: 0.3),
      inactiveThumbColor:
          SanggamTheme.onSurfaceDim,
      inactiveTrackColor:
          SanggamTheme.surfaceVariant,
      secondary: Icon(icon,
          color: SanggamTheme.onSurfaceDim),
      title: Text(title,
          style: const TextStyle(
              color: Colors.white)),
      subtitle: Text(subtitle,
          style: const TextStyle(
            color: SanggamTheme.onSurfaceDim,
            fontSize: 12,
          )),
    );
  }
}
