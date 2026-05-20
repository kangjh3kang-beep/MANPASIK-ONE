import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/providers/theme_provider.dart';
import 'package:manpasik/shared/providers/locale_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';
import 'tenancy_settings_section.dart';

// ───────────────────────────────────────────────────
// SettingsScreen — Sanggam Orbit 설정
//
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppTheme.sanggamGold ~8x → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] JagaeContainer → SanggamContainer
// [Rule 4] jagae_pattern.dart import 제거
// [Rule 4] Theme.of(context) 3x 제거
// [Rule 4] theme.colorScheme.* ~5x → SanggamTheme 직접
// [Rule 4] theme.textTheme.* ~4x → 직접 TextStyle
// [Rule 4] withOpacity ~12x → withValues(alpha:)
// [Rule 4] Color(0xFF...) 하드코딩 → SanggamTheme 상수
// [Rule 4] Colors.white70/54/60 → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.redAccent/red → SanggamTheme.error
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] _showThemeDialog + _getThemeModeLabel 미사용 코드 제거
// [Rule 2] h:12→16, borderRadius 20→24
// ───────────────────────────────────────────────────

/// 설정 화면
///
/// user-service GetProfile/GetSubscription으로 프로필·구독 표시.
class SettingsScreen extends ConsumerWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final profileAsync = ref.watch(userProfileProvider);
    final subscriptionAsync = ref.watch(subscriptionInfoProvider);
    final currentThemeMode = ref.watch(themeModeProvider);
    final currentLocale = ref.watch(localeProvider);

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 0, 16, 120),
          physics: const BouncingScrollPhysics(),
          children: [
            // 헤더: 뒤로가기 + 타이틀
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Text(
                    '설정',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
            ),

            // ── 프로필 섹션 ──
            _buildSectionHeader('프로필'),
            _buildGlassSection([
              profileAsync.when(
                data: (profile) => ListTile(
                  leading: CircleAvatar(
                    radius: 20,
                    backgroundColor:
                        SanggamTheme.primary.withValues(alpha: 0.2),
                    child: const Icon(Icons.person,
                        color: SanggamTheme.primary),
                  ),
                  title: Text(
                      profile?.displayName ??
                          authState.displayName ??
                          '사용자',
                      style: const TextStyle(color: Colors.white)),
                  subtitle: Text(
                      profile?.email ??
                          authState.email ??
                          '로그인이 필요합니다',
                      style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim)),
                  trailing: const Icon(Icons.chevron_right,
                      color: SanggamTheme.onSurfaceDim),
                  onTap: () => context.push('/settings/profile'),
                ),
                loading: () => const Center(
                    child: Padding(
                        padding: EdgeInsets.all(16),
                        child: CircularProgressIndicator())),
                error: (_, __) => ListTile(
                  leading: const Icon(Icons.person,
                      color: SanggamTheme.onSurfaceDim),
                  title: Text(authState.displayName ?? '사용자',
                      style: const TextStyle(color: Colors.white)),
                  subtitle: const Text('로그인 필요',
                      style: TextStyle(
                          color: SanggamTheme.onSurfaceDim)),
                ),
              ),
              if (subscriptionAsync.hasValue &&
                  subscriptionAsync.value != null)
                ListTile(
                  leading: const Icon(Icons.card_membership,
                      color: SanggamTheme.jagaeCyan),
                  title: Text(
                      _tierLabel(subscriptionAsync.value!.tier),
                      style: const TextStyle(
                          color: SanggamTheme.jagaeCyan,
                          fontWeight: FontWeight.bold)),
                  subtitle: const Text('구독 정보',
                      style: TextStyle(
                          color: SanggamTheme.onSurfaceDim)),
                  trailing: const Icon(Icons.chevron_right,
                      color: SanggamTheme.onSurfaceDim),
                  onTap: () => context.push('/market/subscription'),
                ),
            ]),

            const SizedBox(height: 24),

            // ── 서비스 설정 ──
            _buildSectionHeader('서비스'),
            _buildGlassSection([
              _buildGlassTile(Icons.notifications_outlined, '알림 설정',
                  onTap: () =>
                      context.push('/settings/notifications')),
              _buildGlassTile(Icons.shield_outlined, '보안',
                  subtitle: '비밀번호, 생체인증',
                  onTap: () => context.push('/settings/security')),
              _buildGlassTile(
                  Icons.accessibility_new_outlined, '접근성',
                  subtitle: '화면 읽기, 글꼴',
                  onTap: () =>
                      context.push('/settings/accessibility')),
              _buildGlassTile(Icons.emergency_outlined, '긴급 대응',
                  subtitle: '119 자동 신고',
                  iconColor: SanggamTheme.error,
                  onTap: () =>
                      context.push('/settings/emergency')),
              _buildGlassTile(
                  Icons.warning_amber_outlined, '119 에스컬레이션',
                  subtitle: '응급 상황 처리',
                  iconColor: SanggamTheme.error,
                  onTap: () =>
                      context.push('/settings/escalation')),
              _buildGlassTile(Icons.privacy_tip_outlined, '동의 관리',
                  onTap: () =>
                      context.push('/settings/consent')),
              _buildGlassTile(
                  Icons.cloud_download_outlined, '데이터 내보내기',
                  subtitle: 'FHIR/CSV 형식 지원',
                  onTap: () =>
                      context.push('/settings/data-export')),
            ]),

            const SizedBox(height: 24),

            // ── 조직 관리 (Phase AE-3) ──
            const TenancySettingsSection(),
            const SizedBox(height: 24),

            // ── 바로가기 ──
            _buildSectionHeader('바로가기'),
            _buildGlassSection([
              _buildGlassTile(
                  Icons.local_hospital_outlined, '의료 서비스',
                  subtitle: '화상진료, 처방전, 병원 검색',
                  iconColor: SanggamTheme.jagaeCyan,
                  onTap: () => context.push('/medical')),
              _buildGlassTile(
                  Icons.family_restroom_outlined, '가족 건강',
                  subtitle: '가족 관리, 보호자 모니터링',
                  iconColor: SanggamTheme.primary,
                  onTap: () => context.push('/family')),
              _buildGlassTile(Icons.forum_outlined, '커뮤니티',
                  subtitle: '건강 포럼, 챌린지, Q&A',
                  onTap: () => context.push('/community')),
              _buildGlassTile(Icons.devices_outlined, '기기 관리',
                  subtitle: '리더기 목록, 펌웨어 업데이트',
                  onTap: () => context.push('/devices')),
            ]),

            const SizedBox(height: 24),

            // ── 일반 설정 ──
            _buildSectionHeader('일반'),
            _buildThemeSelector(ref, currentThemeMode),
            const SizedBox(height: 16),
            _buildGlassSection([
              _buildGlassTile(Icons.language_outlined, '언어',
                  subtitle: SupportedLocales.getLanguageName(
                      currentLocale.languageCode),
                  onTap: () => _showLanguageDialog(
                      context, ref, currentLocale)),
            ]),

            const SizedBox(height: 24),

            // ── 앱 정보 ──
            _buildSectionHeader('앱 정보'),
            _buildGlassSection([
              const ListTile(
                leading: Icon(Icons.info_outline,
                    color: SanggamTheme.onSurfaceDim),
                title: Text('버전',
                    style: TextStyle(color: Colors.white)),
                subtitle: Text('1.0.0',
                    style: TextStyle(
                        color: SanggamTheme.onSurfaceDim)),
              ),
              _buildGlassTile(Icons.description_outlined, '이용약관',
                  onTap: () => context.push('/settings/terms')),
              _buildGlassTile(
                  Icons.privacy_tip_outlined, '개인정보처리방침',
                  onTap: () => context.push('/settings/privacy')),
              _buildGlassTile(Icons.history, '약관 변경 이력',
                  onTap: () => _showTermsChangeHistory(context)),
            ]),

            const SizedBox(height: 24),

            // ── 고객 지원 ──
            _buildSectionHeader('고객 지원'),
            _buildGlassSection([
              _buildGlassTile(
                  Icons.help_outline, '자주 묻는 질문 (FAQ)',
                  onTap: () => context.push('/support')),
              _buildGlassTile(Icons.headset_mic_outlined, '1:1 문의',
                  onTap: () =>
                      context.push('/settings/inquiry/create')),
              _buildGlassTile(Icons.campaign_outlined, '공지사항',
                  onTap: () => context.push('/support/notices')),
            ]),

            // ── 관리자 ──
            if (authState.isAdmin) ...[
              const SizedBox(height: 24),
              _buildSectionHeader('관리자'),
              _buildGlassSection([
                _buildGlassTile(
                    Icons.admin_panel_settings_outlined, '관리자 대시보드',
                    iconColor: SanggamTheme.primary,
                    onTap: () =>
                        context.push('/admin/dashboard')),
              ]),
            ],

            if (authState.isAuthenticated) ...[
              const SizedBox(height: 32),
              _buildGlassSection([
                ListTile(
                  leading: const Icon(Icons.logout,
                      color: SanggamTheme.error),
                  title: const Text('로그아웃',
                      style: TextStyle(
                          color: SanggamTheme.error,
                          fontWeight: FontWeight.bold)),
                  onTap: () => _showLogoutDialog(context, ref),
                ),
              ]),
              const SizedBox(height: 32),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildGlassSection(List<Widget> children) {
    return SanggamContainer(
      borderRadius: 24,
      jagaeOpacity: 0.1,
      borderWidth: 0.5,
      borderColor: Colors.white.withValues(alpha: 0.1),
      padding: EdgeInsets.zero,
      child: Column(children: children),
    );
  }

  Widget _buildGlassTile(IconData icon, String title,
      {String? subtitle,
      VoidCallback? onTap,
      Color iconColor = SanggamTheme.onSurfaceDim}) {
    return ListTile(
      leading: Icon(icon, color: iconColor),
      title: Text(title,
          style: const TextStyle(color: Colors.white)),
      subtitle: subtitle != null
          ? Text(subtitle,
              style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim))
          : null,
      trailing: const Icon(Icons.chevron_right,
          color: SanggamTheme.onSurfaceDim),
      onTap: onTap,
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: const TextStyle(
          color: SanggamTheme.primary,
          fontSize: 14,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  String _tierLabel(int tier) {
    switch (tier) {
      case 0:
        return '무료';
      case 1:
        return '베이직';
      case 2:
        return '프로';
      case 3:
        return '클리닉';
      default:
        return '무료';
    }
  }

  Widget _buildThemeSelector(WidgetRef ref, ThemeMode current) {
    final options = [
      (
        mode: ThemeMode.system,
        label: '시스템',
        icon: Icons.settings_suggest_rounded
      ),
      (
        mode: ThemeMode.light,
        label: '백자',
        icon: Icons.wb_sunny_rounded
      ),
      (
        mode: ThemeMode.dark,
        label: '코스믹',
        icon: Icons.nightlight_round
      ),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.symmetric(horizontal: 4, vertical: 8),
          child: Text('테마 설정',
              style: TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12)),
        ),
        Row(
          children: options.map((opt) {
            final isSelected = current == opt.mode;
            return Expanded(
              child: GestureDetector(
                onTap: () => ref
                    .read(themeModeProvider.notifier)
                    .setThemeMode(opt.mode),
                child: Container(
                  margin:
                      const EdgeInsets.symmetric(horizontal: 4),
                  padding:
                      const EdgeInsets.symmetric(vertical: 16),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? SanggamTheme.primary
                            .withValues(alpha: 0.2)
                        : Colors.white
                            .withValues(alpha: 0.05),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: isSelected
                          ? SanggamTheme.primary
                          : Colors.white
                              .withValues(alpha: 0.1),
                      width: isSelected ? 1.5 : 0.5,
                    ),
                    boxShadow: isSelected
                        ? [
                            BoxShadow(
                              color: SanggamTheme.primary
                                  .withValues(alpha: 0.3),
                              blurRadius: 12,
                              spreadRadius: -2,
                            )
                          ]
                        : null,
                  ),
                  child: Column(
                    children: [
                      Icon(opt.icon,
                          color: isSelected
                              ? SanggamTheme.primary
                              : SanggamTheme.onSurfaceDim,
                          size: 28),
                      const SizedBox(height: 8),
                      Text(
                        opt.label,
                        style: TextStyle(
                          color: isSelected
                              ? Colors.white
                              : SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                          fontWeight: isSelected
                              ? FontWeight.bold
                              : FontWeight.normal,
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

  /// 언어 선택 다이얼로그 (6개 언어)
  void _showLanguageDialog(
      BuildContext context, WidgetRef ref, Locale current) {
    showDialog(
      context: context,
      builder: (context) => SimpleDialog(
        title: const Text('언어 선택'),
        children: SupportedLocales.all.map((locale) {
          final isSelected =
              locale.languageCode == current.languageCode;
          return ListTile(
            title: Text(SupportedLocales.getLanguageName(
                locale.languageCode)),
            trailing: isSelected
                ? const Icon(Icons.check,
                    color: SanggamTheme.jagaeCyan)
                : null,
            onTap: () {
              ref
                  .read(localeProvider.notifier)
                  .setLocaleByCode(locale.languageCode);
              Navigator.pop(context);
            },
          );
        }).toList(),
      ),
    );
  }

  /// 약관 변경 이력 바텀시트
  void _showTermsChangeHistory(BuildContext context) {
    final history = [
      {
        'date': '2026-02-01',
        'title': '개인정보처리방침 개정',
        'summary': '건강 데이터 제3자 제공 조항 추가'
      },
      {
        'date': '2026-01-15',
        'title': '이용약관 개정',
        'summary': '원격 진료 서비스 이용 조건 신설'
      },
      {
        'date': '2025-12-01',
        'title': '이용약관 제정',
        'summary': '만파식 서비스 최초 약관 제정'
      },
    ];
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: SanggamTheme.background,
      builder: (ctx) => DraggableScrollableSheet(
        initialChildSize: 0.5,
        maxChildSize: 0.85,
        expand: false,
        builder: (_, scrollCtrl) => Column(
          children: [
            Container(
              margin: const EdgeInsets.only(top: 8),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: SanggamTheme.onSurfaceDim
                    .withValues(alpha: 0.3),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const Padding(
              padding: EdgeInsets.all(16),
              child: Text('약관 변경 이력',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                  )),
            ),
            Expanded(
              child: ListView.separated(
                controller: scrollCtrl,
                padding:
                    const EdgeInsets.symmetric(horizontal: 16),
                itemCount: history.length,
                separatorBuilder: (_, __) =>
                    const Divider(height: 1),
                itemBuilder: (_, i) {
                  final h = history[i];
                  return ListTile(
                    leading: const Icon(
                        Icons.article_outlined,
                        color: SanggamTheme.primary),
                    title: Text(h['title']!),
                    subtitle: Text(h['summary']!),
                    trailing: Text(h['date']!,
                        style: const TextStyle(
                          color: SanggamTheme.onSurfaceDim,
                          fontSize: 12,
                        )),
                  );
                },
              ),
            ),
          ],
        ),
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
