import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ConsentManagementScreen — Sanggam Orbit 동의 관리
//
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] ThemeData params 5x 제거
// [Rule 4] theme.textTheme ~7x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~12x → SanggamTheme 상수
// [Rule 4] Card ~5x → SanggamContainer
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] AlertDialog 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 12→16, borderRadius 4→8/8→16
// ───────────────────────────────────────────────────

/// 동의 관리 화면 (GDPR Art.7, PIPA §15, §23)
///
/// 5단계 계층적 동의:
/// 1차: 필수 동의 (서비스 이용약관, 개인정보 처리)
/// 2차: 건강정보 수집 동의 (GDPR Art.9 명시적)
/// 3차: 위치정보 수집 동의 (선택)
/// 4차: 마케팅/분석 동의 (선택)
/// 5차: AI 학습 데이터 활용 동의 (선택)
class ConsentManagementScreen
    extends ConsumerStatefulWidget {
  const ConsentManagementScreen(
      {super.key});

  @override
  ConsumerState<ConsentManagementScreen>
      createState() =>
          _ConsentManagementScreenState();
}

class _ConsentManagementScreenState
    extends ConsumerState<
        ConsentManagementScreen> {
  bool _serviceTerms = true;
  bool _privacyPolicy = true;
  bool _healthDataConsent = true;
  bool _locationConsent = false;
  bool _marketingConsent = false;
  bool _aiTrainingConsent = false;

  @override
  void initState() {
    super.initState();
    _loadConsents();
  }

  Future<void> _loadConsents() async {
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final data = await rest.getNotificationPreferences(userId);
      if (mounted) {
        setState(() {
          _locationConsent = (data['location_consent'] as bool?) ?? false;
          _marketingConsent = (data['marketing_consent'] as bool?) ?? false;
          _aiTrainingConsent = (data['ai_training_consent'] as bool?) ?? false;
        });
      }
    } on DioException {
      // 로드 실패 시 기본값 유지
    }
  }

  Future<void> _saveConsent(String key, bool value) async {
    try {
      final rest = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await rest.updateNotificationPreferences(userId, {key: value});
    } on DioException {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('동의 설정 저장에 실패했습니다')),
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
                      '동의 관리',
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
              child: SingleChildScrollView(
                padding:
                    const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .start,
                  children: [
                    _buildInfoBanner(),
                    const SizedBox(
                        height: 24),

                    _buildSectionHeader(
                        '필수 동의',
                        '서비스 이용에 반드시 필요합니다'),
                    _buildConsentTile(
                      '서비스 이용약관',
                      '만파식 서비스 이용 약관',
                      _serviceTerms,
                      true,
                      'GDPR Art.6(1)(b)',
                      (v) => _handleRequired(
                          v, '서비스 이용약관'),
                    ),
                    _buildConsentTile(
                      '개인정보 처리방침',
                      '개인정보 수집 및 이용 동의',
                      _privacyPolicy,
                      true,
                      'PIPA §15',
                      (v) => _handleRequired(
                          v, '개인정보 처리방침'),
                    ),
                    _buildConsentTile(
                      '건강 데이터 수집',
                      '혈액 분석, 바이오마커 수치 수집',
                      _healthDataConsent,
                      true,
                      'GDPR Art.9(2)(a), PIPA §23',
                      (v) => _handleRequired(
                          v, '건강정보 수집'),
                    ),
                    const SizedBox(
                        height: 16),

                    _buildSectionHeader(
                        '선택 동의',
                        '동의하지 않아도 서비스 이용 가능'),
                    _buildConsentTile(
                      '위치정보 수집',
                      '환경 보정 및 주변 병원 검색',
                      _locationConsent,
                      false,
                      '위치정보법',
                      (v) {
                        setState(() =>
                            _locationConsent =
                                v ?? false);
                        _saveConsent(
                            'location_consent',
                            v ?? false);
                      },
                    ),
                    _buildConsentTile(
                      '마케팅 및 분석',
                      '사용 패턴 분석, 맞춤 정보 제공',
                      _marketingConsent,
                      false,
                      'GDPR Art.6(1)(a)',
                      (v) {
                        setState(() =>
                            _marketingConsent =
                                v ?? false);
                        _saveConsent(
                            'marketing_consent',
                            v ?? false);
                      },
                    ),
                    _buildConsentTile(
                      'AI 학습 데이터 활용',
                      '익명화 데이터로 AI 모델 개선 (Federated Learning)',
                      _aiTrainingConsent,
                      false,
                      'GDPR Art.6(1)(a)',
                      (v) {
                        setState(() =>
                            _aiTrainingConsent =
                                v ?? false);
                        _saveConsent(
                            'ai_training_consent',
                            v ?? false);
                      },
                    ),

                    const SizedBox(
                        height: 32),
                    OutlinedButton.icon(
                      onPressed: () {},
                      icon: const Icon(
                          Icons.history),
                      label: const Text(
                          '동의 변경 이력 조회'),
                      style: OutlinedButton
                          .styleFrom(
                        foregroundColor:
                            SanggamTheme
                                .onSurfaceDim,
                        side: const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                    ),
                    const SizedBox(
                        height: 8),
                    OutlinedButton.icon(
                      onPressed: () {},
                      icon: const Icon(
                          Icons.download),
                      label: const Text(
                          '내 데이터 내보내기 (FHIR R4 / CSV)'),
                      style: OutlinedButton
                          .styleFrom(
                        foregroundColor:
                            SanggamTheme
                                .onSurfaceDim,
                        side: const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                    ),
                    const SizedBox(
                        height: 8),
                    OutlinedButton.icon(
                      onPressed:
                          _showDeleteDialog,
                      icon: const Icon(
                          Icons
                              .delete_forever,
                          color: SanggamTheme
                              .error),
                      label: const Text(
                          '계정 및 데이터 삭제',
                          style: TextStyle(
                              color:
                                  SanggamTheme
                                      .error)),
                      style: OutlinedButton
                          .styleFrom(
                        side: const BorderSide(
                            color:
                                SanggamTheme
                                    .error),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                    ),
                    const SizedBox(
                        height: 24),
                    _buildLegalNotice(),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoBanner() {
    return SanggamContainer(
      borderRadius: 16,
      backgroundColor: SanggamTheme.primary
          .withValues(alpha: 0.1),
      child: Row(
        children: [
          const Icon(Icons.info_outline,
              color: SanggamTheme.primary),
          const SizedBox(width: 16),
          const Expanded(
            child: Text(
              '각 동의를 확인하고 언제든 변경할 수 있습니다.\n'
              '필수 동의 철회 시 서비스 탈퇴로 처리됩니다.',
              style: TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 13,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(
      String title, String sub) {
    return Padding(
      padding: const EdgeInsets.only(
          bottom: 8),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Text(title,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              )),
          Text(sub,
              style: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 12,
              )),
        ],
      ),
    );
  }

  Widget _buildConsentTile(
    String title,
    String subtitle,
    bool value,
    bool required,
    String legal,
    ValueChanged<bool?> onChanged,
  ) {
    return Padding(
      padding: const EdgeInsets.only(
          bottom: 4),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: Column(
          children: [
            SwitchListTile(
              activeColor:
                  SanggamTheme.primary,
              title: Row(
                children: [
                  Flexible(
                    child: Text(title,
                        style:
                            const TextStyle(
                                color: Colors
                                    .white)),
                  ),
                  if (required) ...[
                    const SizedBox(width: 8),
                    Container(
                      padding:
                          const EdgeInsets
                              .symmetric(
                              horizontal: 8,
                              vertical: 2),
                      decoration:
                          BoxDecoration(
                        color: SanggamTheme
                            .error
                            .withValues(
                                alpha: 0.1),
                        borderRadius:
                            BorderRadius
                                .circular(
                                    8),
                      ),
                      child: const Text(
                          '필수',
                          style: TextStyle(
                            fontSize: 10,
                            color:
                                SanggamTheme
                                    .error,
                          )),
                    ),
                  ],
                ],
              ),
              subtitle: Text(subtitle,
                  style: const TextStyle(
                    color: SanggamTheme
                        .onSurfaceDim,
                    fontSize: 12,
                  )),
              value: value,
              onChanged: onChanged,
            ),
            Padding(
              padding:
                  const EdgeInsets.fromLTRB(
                      16, 0, 16, 8),
              child: Row(
                children: [
                  const Icon(Icons.gavel,
                      size: 12,
                      color: SanggamTheme
                          .onSurfaceDim),
                  const SizedBox(width: 4),
                  Flexible(
                    child: Text(legal,
                        style:
                            const TextStyle(
                          color: SanggamTheme
                              .onSurfaceDim,
                          fontSize: 11,
                        )),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _handleRequired(
      bool? value, String name) {
    if (value == false) {
      showDialog(
        context: context,
        builder: (ctx) => AlertDialog(
          backgroundColor:
              SanggamTheme.surface,
          shape: RoundedRectangleBorder(
              borderRadius:
                  BorderRadius.circular(16)),
          title: const Text('필수 동의 철회',
              style: TextStyle(
                  color: Colors.white)),
          content: Text(
            '$name은(는) 필수 동의입니다.\n철회 시 서비스 탈퇴로 처리됩니다.',
            style: const TextStyle(
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
              child: const Text('취소'),
            ),
            TextButton(
              onPressed: () =>
                  Navigator.pop(ctx),
              style: TextButton.styleFrom(
                  foregroundColor:
                      SanggamTheme.error),
              child: const Text('철회'),
            ),
          ],
        ),
      );
    }
  }

  void _showDeleteDialog() {
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
          '즉시 비활성화 → 30일 내 PII 삭제\n'
          '건강 데이터는 익명화 후 10년 보존 (의료기기법)\n\n'
          '되돌릴 수 없습니다.',
          style: TextStyle(
              color: SanggamTheme
                  .onSurfaceDim),
        ),
        actions: [
          TextButton(
            onPressed: () =>
                Navigator.pop(ctx),
            style: TextButton.styleFrom(
                foregroundColor: SanggamTheme
                    .onSurfaceDim),
            child: const Text('취소'),
          ),
          TextButton(
            onPressed: () =>
                Navigator.pop(ctx),
            style: TextButton.styleFrom(
                foregroundColor:
                    SanggamTheme.error),
            child: const Text('삭제 요청'),
          ),
        ],
      ),
    );
  }

  Widget _buildLegalNotice() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: SanggamTheme.surfaceVariant
            .withValues(alpha: 0.3),
        borderRadius:
            BorderRadius.circular(16),
      ),
      child: const Text(
        '본 동의 관리는 GDPR(EU), PIPA(한국), HIPAA(미국), PIPL(중국), APPI(일본) 규정을 준수합니다. '
        '동의 기록은 5년간 보존됩니다. 문의: privacy@manpasik.com',
        style: TextStyle(
          color:
              SanggamTheme.onSurfaceDim,
          fontSize: 12,
        ),
      ),
    );
  }
}
