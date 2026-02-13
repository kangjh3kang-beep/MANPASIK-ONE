import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';

/// 설정 화면
///
/// user-service GetProfile/GetSubscription으로 프로필·구독 표시.
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final authState = ref.watch(authProvider);
    final profileAsync = ref.watch(userProfileProvider);
    final subscriptionAsync = ref.watch(subscriptionInfoProvider);
    final currentThemeMode = ref.watch(themeModeProvider);
    final currentLocale = ref.watch(localeProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('설정'),
      ),
      body: ListView(
        children: [
          // ── 프로필 섹션 (user-service GetProfile) ──
          _buildSectionHeader(theme, '프로필'),
          profileAsync.when(
            data: (profile) => ListTile(
              leading: CircleAvatar(
                radius: 24,
                backgroundColor: theme.colorScheme.primaryContainer,
                child: Icon(
                  Icons.person,
                  color: theme.colorScheme.onPrimaryContainer,
                ),
              ),
              title: Text(profile?.displayName ?? authState.displayName ?? '사용자'),
              subtitle: Text(profile?.email ?? authState.email ?? '로그인이 필요합니다'),
              trailing: const Icon(Icons.chevron_right),
              onTap: () {},
            ),
            loading: () => const ListTile(
              leading: SizedBox(
                width: 48,
                height: 48,
                child: Center(child: CircularProgressIndicator(strokeWidth: 2)),
              ),
              title: Text('프로필 로딩 중…'),
            ),
            error: (_, __) => ListTile(
              leading: const Icon(Icons.person_outline),
              title: Text(authState.displayName ?? '사용자'),
              subtitle: Text(authState.email ?? '로그인이 필요합니다'),
              trailing: const Icon(Icons.chevron_right),
            ),
          ),
          if (subscriptionAsync.hasValue && subscriptionAsync.value != null)
            ListTile(
              leading: const Icon(Icons.card_membership_outlined),
              title: Text(_tierLabel(subscriptionAsync.value!.tier)),
              subtitle: const Text('구독 정보'),
            ),
          const Divider(),

          // ── 일반 설정 섹션 ──
          _buildSectionHeader(theme, '일반'),

          // 테마 설정
          ListTile(
            leading: const Icon(Icons.brightness_6_outlined),
            title: const Text('테마'),
            subtitle: Text(_getThemeModeLabel(currentThemeMode)),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showThemeDialog(context, ref, currentThemeMode),
          ),

          // 언어 설정
          ListTile(
            leading: const Icon(Icons.language_outlined),
            title: const Text('언어'),
            subtitle: Text(SupportedLocales.getLanguageName(currentLocale.languageCode)),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => _showLanguageDialog(context, ref, currentLocale),
          ),
          const Divider(),

          // ── 앱 정보 섹션 ──
          _buildSectionHeader(theme, '앱 정보'),

          const ListTile(
            leading: Icon(Icons.info_outlined),
            title: Text('버전'),
            subtitle: Text('1.0.0'),
          ),

          ListTile(
            leading: const Icon(Icons.description_outlined),
            title: const Text('이용약관'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // TODO: 이용약관 화면
            },
          ),

          ListTile(
            leading: const Icon(Icons.privacy_tip_outlined),
            title: const Text('개인정보처리방침'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () {
              // TODO: 개인정보처리방침 화면
            },
          ),
          const Divider(),

          // ── 계정 섹션 ──
          if (authState.isAuthenticated) ...[
            _buildSectionHeader(theme, '계정'),
            ListTile(
              leading: const Icon(Icons.logout, color: Colors.red),
              title: const Text(
                '로그아웃',
                style: TextStyle(color: Colors.red),
              ),
              onTap: () => _showLogoutDialog(context, ref),
            ),
            const SizedBox(height: 32),
          ],
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

  String _tierLabel(int tier) {
    switch (tier) {
      case 0:
        return 'Free';
      case 1:
        return 'Basic';
      case 2:
        return 'Pro';
      case 3:
        return 'Clinical';
      default:
        return 'Free';
    }
  }

  String _getThemeModeLabel(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.system:
        return '시스템 기본값';
      case ThemeMode.light:
        return '라이트 모드';
      case ThemeMode.dark:
        return '다크 모드';
    }
  }

  /// 테마 선택 다이얼로그
  void _showThemeDialog(BuildContext context, WidgetRef ref, ThemeMode current) {
    final options = [
      (mode: ThemeMode.system, label: '시스템 기본값', icon: Icons.settings_suggest_outlined),
      (mode: ThemeMode.light, label: '라이트 모드', icon: Icons.light_mode_outlined),
      (mode: ThemeMode.dark, label: '다크 모드', icon: Icons.dark_mode_outlined),
    ];

    showDialog(
      context: context,
      builder: (context) => SimpleDialog(
        title: const Text('테마 선택'),
        children: options.map((opt) {
          final isSelected = opt.mode == current;
          return ListTile(
            leading: Icon(opt.icon),
            title: Text(opt.label),
            trailing: isSelected ? const Icon(Icons.check, color: Colors.green) : null,
            onTap: () {
              ref.read(themeModeProvider.notifier).setThemeMode(opt.mode);
              Navigator.pop(context);
            },
          );
        }).toList(),
      ),
    );
  }

  /// 언어 선택 다이얼로그 (6개 언어)
  void _showLanguageDialog(BuildContext context, WidgetRef ref, Locale current) {
    showDialog(
      context: context,
      builder: (context) => SimpleDialog(
        title: const Text('언어 선택'),
        children: SupportedLocales.all.map((locale) {
          final isSelected = locale.languageCode == current.languageCode;
          return ListTile(
            title: Text(SupportedLocales.getLanguageName(locale.languageCode)),
            trailing: isSelected ? const Icon(Icons.check, color: Colors.green) : null,
            onTap: () {
              ref.read(localeProvider.notifier).setLocaleByCode(locale.languageCode);
              Navigator.pop(context);
            },
          );
        }).toList(),
      ),
    );
  }

  /// 로그아웃 확인 다이얼로그
  void _showLogoutDialog(BuildContext context, WidgetRef ref) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('로그아웃'),
        content: const Text('정말 로그아웃하시겠습니까?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('취소'),
          ),
          FilledButton(
            onPressed: () {
              ref.read(authProvider.notifier).logout();
              Navigator.pop(context);
              context.go('/login');
            },
            child: const Text('로그아웃'),
          ),
        ],
      ),
    );
  }
}
