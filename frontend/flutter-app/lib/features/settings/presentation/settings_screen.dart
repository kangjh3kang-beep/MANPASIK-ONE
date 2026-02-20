import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/widgets/jagae_pattern.dart';
import 'package:manpasik/core/theme/app_theme.dart';

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
      backgroundColor: Colors.transparent, // Global Cosmic Background
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => context.pop(),
        ),
        title: const Text('설정', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
        centerTitle: true,
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 120), // Bottom padding for Glass Dock
        physics: const BouncingScrollPhysics(),
        children: [
          // ── 프로필 섹션 ──
          _buildSectionHeader(theme, '프로필'),
          _buildGlassSection([
            profileAsync.when(
              data: (profile) => ListTile(
                leading: CircleAvatar(
                  radius: 20,
                  backgroundColor: AppTheme.sanggamGold.withOpacity(0.2),
                  child: const Icon(Icons.person, color: AppTheme.sanggamGold),
                ),
                title: Text(profile?.displayName ?? authState.displayName ?? '사용자', style: const TextStyle(color: Colors.white)),
                subtitle: Text(profile?.email ?? authState.email ?? '로그인이 필요합니다', style: TextStyle(color: Colors.white.withOpacity(0.6))),
                trailing: const Icon(Icons.chevron_right, color: Colors.white54),
                onTap: () => context.push('/settings/profile'),
              ),
              loading: () => const Center(child: Padding(padding: EdgeInsets.all(16), child: CircularProgressIndicator())),
              error: (_, __) => ListTile(
                leading: const Icon(Icons.person, color: Colors.white54),
                title: Text(authState.displayName ?? '사용자', style: const TextStyle(color: Colors.white)),
                subtitle: Text('로그인 필요', style: TextStyle(color: Colors.white.withOpacity(0.6))),
              ),
            ),
            if (subscriptionAsync.hasValue && subscriptionAsync.value != null)
              ListTile(
                leading: const Icon(Icons.card_membership, color: AppTheme.waveCyan),
                title: Text(_tierLabel(subscriptionAsync.value!.tier), style: const TextStyle(color: AppTheme.waveCyan, fontWeight: FontWeight.bold)),
                subtitle: Text('구독 정보', style: TextStyle(color: Colors.white.withOpacity(0.6))),
                trailing: const Icon(Icons.chevron_right, color: Colors.white54),
                onTap: () => context.push('/market/subscription'),
              ),
          ]),

          const SizedBox(height: 24),

          // ── 서비스 설정 ──
          _buildSectionHeader(theme, '서비스'),
          _buildGlassSection([
            _buildGlassTile(Icons.notifications_outlined, '알림 설정', 
              onTap: () => context.push('/settings/notifications')),
            _buildGlassTile(Icons.shield_outlined, '보안', subtitle: '비밀번호, 생체인증',
              onTap: () => context.push('/settings/security')),
            _buildGlassTile(Icons.accessibility_new_outlined, '접근성', subtitle: '화면 읽기, 글꼴',
              onTap: () => context.push('/settings/accessibility')),
            _buildGlassTile(Icons.emergency_outlined, '긴급 대응', subtitle: '119 자동 신고', iconColor: Colors.redAccent,
              onTap: () => context.push('/settings/emergency')),
            _buildGlassTile(Icons.privacy_tip_outlined, '동의 관리',
              onTap: () => context.push('/settings/consent')),
          ]),

          const SizedBox(height: 24),

          // ── 일반 설정 ──
          _buildSectionHeader(theme, '일반'),
            // 테마 선택 섹션 (Horizontal Cards)
            _buildThemeSelector(context, ref, currentThemeMode),
            const SizedBox(height: 12),
            
            // 언어 설정
            _buildGlassSection([
              _buildGlassTile(Icons.language_outlined, '언어', 
                subtitle: SupportedLocales.getLanguageName(currentLocale.languageCode),
                onTap: () => _showLanguageDialog(context, ref, currentLocale)),
            ]),

          const SizedBox(height: 24),

          // ── 앱 정보 ──
          _buildSectionHeader(theme, '앱 정보'),
          _buildGlassSection([
            const ListTile(
              leading: Icon(Icons.info_outline, color: Colors.white70),
              title: Text('버전', style: TextStyle(color: Colors.white)),
              subtitle: Text('1.0.0', style: TextStyle(color: Colors.white54)),
            ),
            _buildGlassTile(Icons.description_outlined, '이용약관',
              onTap: () => context.push('/settings/terms')),
            _buildGlassTile(Icons.privacy_tip_outlined, '개인정보처리방침',
              onTap: () => context.push('/settings/privacy')),
            _buildGlassTile(Icons.history, '약관 변경 이력',
              onTap: () => _showTermsChangeHistory(context)),
          ]),

          const SizedBox(height: 24),

           // ── 고객 지원 ──
          _buildSectionHeader(theme, '고객 지원'),
          _buildGlassSection([
            _buildGlassTile(Icons.help_outline, '자주 묻는 질문 (FAQ)',
              onTap: () => context.push('/support')),
            _buildGlassTile(Icons.headset_mic_outlined, '1:1 문의',
              onTap: () => context.push('/support')),
          ]),

          if (authState.isAuthenticated) ...[
             const SizedBox(height: 32),
             _buildGlassSection([
               ListTile(
                leading: const Icon(Icons.logout, color: Colors.redAccent),
                title: const Text('로그아웃', style: TextStyle(color: Colors.redAccent, fontWeight: FontWeight.bold)),
                onTap: () => _showLogoutDialog(context, ref),
              ),
             ]),
             const SizedBox(height: 32),
          ],
        ],
      ),
    );
  }

  Widget _buildGlassSection(List<Widget> children) {
    return JagaeContainer(
      opacity: 0.1,
      showLattice: true,
      decoration: BoxDecoration(
        color: const Color(0xFF1A1F35).withOpacity(0.4),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.white.withOpacity(0.1)),
      ),
      child: Column(
        children: children,
      ),
    );
  }

  Widget _buildGlassTile(IconData icon, String title, {String? subtitle, VoidCallback? onTap, Color iconColor = Colors.white70}) {
    return ListTile(
      leading: Icon(icon, color: iconColor),
      title: Text(title, style: const TextStyle(color: Colors.white)),
      subtitle: subtitle != null ? Text(subtitle, style: TextStyle(color: Colors.white.withOpacity(0.6))) : null,
      trailing: const Icon(Icons.chevron_right, color: Colors.white54),
      onTap: onTap,
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
      (mode: ThemeMode.light, label: '한국적 화이트 모드 (백자)', icon: Icons.light_mode_outlined), // Updated Label
      (mode: ThemeMode.dark, label: '만파식 다크 모드 (우주)', icon: Icons.dark_mode_outlined),
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
  void _showTermsChangeHistory(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      builder: (ctx) {
        final theme = Theme.of(ctx);
        final history = [
          {'date': '2026-02-01', 'title': '개인정보처리방침 개정', 'summary': '건강 데이터 제3자 제공 조항 추가'},
          {'date': '2026-01-15', 'title': '이용약관 개정', 'summary': '원격 진료 서비스 이용 조건 신설'},
          {'date': '2025-12-01', 'title': '이용약관 제정', 'summary': '만파식 서비스 최초 약관 제정'},
        ];
        return DraggableScrollableSheet(
          initialChildSize: 0.5,
          maxChildSize: 0.85,
          expand: false,
          builder: (_, scrollCtrl) => Column(
            children: [
              Container(
                margin: const EdgeInsets.only(top: 8),
                width: 40, height: 4,
                decoration: BoxDecoration(
                  color: theme.colorScheme.outline.withOpacity(0.3),
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(16),
                child: Text('약관 변경 이력', style: theme.textTheme.titleMedium),
              ),
              Expanded(
                child: ListView.separated(
                  controller: scrollCtrl,
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  itemCount: history.length,
                  separatorBuilder: (_, __) => const Divider(height: 1),
                  itemBuilder: (_, i) {
                    final h = history[i];
                    return ListTile(
                      leading: Icon(Icons.article_outlined, color: theme.colorScheme.primary),
                      title: Text(h['title']!),
                      subtitle: Text(h['summary']!),
                      trailing: Text(h['date']!,
                          style: theme.textTheme.bodySmall?.copyWith(
                              color: theme.colorScheme.outline)),
                    );
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }

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
  Widget _buildThemeSelector(BuildContext context, WidgetRef ref, ThemeMode current) {
    final options = [
      (mode: ThemeMode.system, label: 'System', icon: Icons.settings_suggest_rounded),
      (mode: ThemeMode.light, label: 'Baekja', icon: Icons.wb_sunny_rounded),
      (mode: ThemeMode.dark, label: 'Cosmic', icon: Icons.nightlight_round),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 8),
          child: Text('테마 설정', style: TextStyle(color: Colors.white.withOpacity(0.8), fontSize: 12)),
        ),
        Row(
          children: options.map((opt) {
            final isSelected = current == opt.mode;
            return Expanded(
              child: GestureDetector(
                onTap: () => ref.read(themeModeProvider.notifier).setThemeMode(opt.mode),
                child: Container(
                  margin: const EdgeInsets.symmetric(horizontal: 4),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  decoration: BoxDecoration(
                    color: isSelected 
                        ? AppTheme.sanggamGold.withOpacity(0.2) 
                        : Colors.white.withOpacity(0.05),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: isSelected ? AppTheme.sanggamGold : Colors.white10,
                      width: isSelected ? 1.5 : 0.5,
                    ),
                    boxShadow: isSelected ? [
                      BoxShadow(
                        color: AppTheme.sanggamGold.withOpacity(0.3),
                        blurRadius: 12,
                        spreadRadius: -2,
                      )
                    ] : null,
                  ),
                  child: Column(
                    children: [
                      Icon(opt.icon, 
                        color: isSelected ? AppTheme.sanggamGold : Colors.white70,
                        size: 28,
                      ),
                      const SizedBox(height: 8),
                      Text(
                        opt.label,
                        style: TextStyle(
                          color: isSelected ? Colors.white : Colors.white60,
                          fontSize: 12,
                          fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            );
          }).toList(),
        ),
      ],
    );
  }
}
