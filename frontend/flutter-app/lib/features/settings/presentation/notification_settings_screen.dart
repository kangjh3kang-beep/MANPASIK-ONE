import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/widgets/cosmic_background.dart';
import 'package:manpasik/shared/widgets/hanji_background.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/shared/widgets/porcelain_container.dart';

class NotificationSettingsScreen extends ConsumerStatefulWidget {
  const NotificationSettingsScreen({super.key});

  @override
  ConsumerState<NotificationSettingsScreen> createState() => _NotificationSettingsScreenState();
}

class _NotificationSettingsScreenState extends ConsumerState<NotificationSettingsScreen> {
  // Temporary local state for UI demonstration
  bool _pushEnabled = true;
  bool _emailEnabled = false;
  bool _marketingEnabled = false;
  bool _healthAlerts = true;
  bool _deviceStatus = true;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: Text('알림 설정', style: TextStyle(color: isDark ? Colors.white : AppTheme.inkBlack, fontWeight: FontWeight.bold)),
        centerTitle: true,
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: isDark ? Colors.white : AppTheme.inkBlack),
          onPressed: () => context.pop(),
        ),
      ),
      body: Stack(
        children: [
          // Background
          if (isDark)
            const CosmicBackground(child: SizedBox.expand())
          else
            const HanjiBackground(child: SizedBox.expand()),

          // Content
          SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(24),
              children: [
                _buildSectionHeader(isDark, '필수 알림'),
                _buildGlassContainer(
                  isDark: isDark,
                  child: Column(
                    children: [
                      _buildSwitchTile(
                        isDark: isDark,
                        title: '건강 측정 알림',
                        subtitle: '비정상 파동 감지 시 즉시 알림',
                        value: _healthAlerts,
                        onChanged: (v) => setState(() => _healthAlerts = v),
                        icon: Icons.monitor_heart_outlined,
                      ),
                      Divider(color: isDark ? Colors.white10 : Colors.black12, height: 1),
                      _buildSwitchTile(
                        isDark: isDark,
                        title: '기기 상태 알림',
                        subtitle: '배터리 부족, 연결 끊김 알림',
                        value: _deviceStatus,
                        onChanged: (v) => setState(() => _deviceStatus = v),
                        icon: Icons.device_thermostat_outlined,
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 32),

                _buildSectionHeader(isDark, '일반 알림'),
                _buildGlassContainer(
                  isDark: isDark,
                  child: Column(
                    children: [
                      _buildSwitchTile(
                        isDark: isDark,
                        title: '푸시 알림',
                        subtitle: '앱 내 활동 및 업데이트 알림',
                        value: _pushEnabled,
                        onChanged: (v) => setState(() => _pushEnabled = v),
                        icon: Icons.notifications_active_outlined,
                      ),
                      Divider(color: isDark ? Colors.white10 : Colors.black12, height: 1),
                      _buildSwitchTile(
                        isDark: isDark,
                        title: '이메일 알림',
                        subtitle: '주간 건강 리포트 수신',
                        value: _emailEnabled,
                        onChanged: (v) => setState(() => _emailEnabled = v),
                        icon: Icons.email_outlined,
                      ),
                      Divider(color: isDark ? Colors.white10 : Colors.black12, height: 1),
                      _buildSwitchTile(
                        isDark: isDark,
                        title: '마케팅 정보 수신',
                        subtitle: '이벤트 및 혜택 알림',
                        value: _marketingEnabled,
                        onChanged: (v) => setState(() => _marketingEnabled = v),
                        icon: Icons.redeem_outlined,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(bool isDark, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(8, 0, 8, 8),
      child: Text(
        title,
        style: TextStyle(
          color: isDark ? AppTheme.sanggamGold : AppTheme.inkBlack,
          fontWeight: FontWeight.bold,
          fontSize: 14,
        ),
      ),
    );
  }

  Widget _buildGlassContainer({required bool isDark, required Widget child}) {
    if (isDark) {
      return JagaeContainer(
        opacity: 0.1,
        showLattice: false,
        decoration: BoxDecoration(
          color: const Color(0xFF1A1F35).withOpacity(0.4),
          borderRadius: BorderRadius.circular(20),
          border: Border.all(color: Colors.white.withOpacity(0.1)),
        ),
        child: child,
      );
    } else {
      return PorcelainContainer(
        child: child,
      );
    }
  }

  Widget _buildSwitchTile({
    required bool isDark,
    required String title,
    required String subtitle,
    required bool value,
    required ValueChanged<bool> onChanged,
    required IconData icon,
  }) {
    return SwitchListTile(
      value: value,
      onChanged: onChanged,
      activeColor: AppTheme.sanggamGold,
      activeTrackColor: AppTheme.sanggamGold.withOpacity(0.3),
      inactiveThumbColor: isDark ? Colors.white54 : Colors.grey,
      inactiveTrackColor: isDark ? Colors.white10 : Colors.black12,
      secondary: Icon(icon, color: isDark ? Colors.white70 : AppTheme.inkBlack),
      title: Text(title, style: TextStyle(color: isDark ? Colors.white : AppTheme.inkBlack)),
      subtitle: Text(subtitle, style: TextStyle(color: isDark ? Colors.white54 : Colors.black54, fontSize: 12)),
    );
  }
}
