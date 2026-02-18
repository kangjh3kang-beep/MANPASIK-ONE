import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// 동의 관리 화면 (GDPR Art.7, PIPA §15, §23)
///
/// 5단계 계층적 동의:
/// 1차: 필수 동의 (서비스 이용약관, 개인정보 처리)
/// 2차: 건강정보 수집 동의 (GDPR Art.9 명시적)
/// 3차: 위치정보 수집 동의 (선택)
/// 4차: 마케팅/분석 동의 (선택)
/// 5차: AI 학습 데이터 활용 동의 (선택)
class ConsentManagementScreen extends ConsumerStatefulWidget {
  const ConsentManagementScreen({super.key});

  @override
  ConsumerState<ConsentManagementScreen> createState() => _ConsentManagementScreenState();
}

class _ConsentManagementScreenState extends ConsumerState<ConsentManagementScreen> {
  bool _serviceTerms = true;
  bool _privacyPolicy = true;
  bool _healthDataConsent = true;
  bool _locationConsent = false;
  bool _marketingConsent = false;
  bool _aiTrainingConsent = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('동의 관리'), centerTitle: true),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildInfoBanner(theme),
            const SizedBox(height: 24),

            _buildSectionHeader(theme, '필수 동의', '서비스 이용에 반드시 필요합니다'),
            _buildConsentTile(theme, '서비스 이용약관', '만파식 서비스 이용 약관', _serviceTerms, true, 'GDPR Art.6(1)(b)', (v) => _handleRequired(v, '서비스 이용약관')),
            _buildConsentTile(theme, '개인정보 처리방침', '개인정보 수집 및 이용 동의', _privacyPolicy, true, 'PIPA §15', (v) => _handleRequired(v, '개인정보 처리방침')),
            _buildConsentTile(theme, '건강 데이터 수집', '혈액 분석, 바이오마커 수치 수집', _healthDataConsent, true, 'GDPR Art.9(2)(a), PIPA §23', (v) => _handleRequired(v, '건강정보 수집')),
            const SizedBox(height: 16),

            _buildSectionHeader(theme, '선택 동의', '동의하지 않아도 서비스 이용 가능'),
            _buildConsentTile(theme, '위치정보 수집', '환경 보정 및 주변 병원 검색', _locationConsent, false, '위치정보법', (v) => setState(() => _locationConsent = v ?? false)),
            _buildConsentTile(theme, '마케팅 및 분석', '사용 패턴 분석, 맞춤 정보 제공', _marketingConsent, false, 'GDPR Art.6(1)(a)', (v) => setState(() => _marketingConsent = v ?? false)),
            _buildConsentTile(theme, 'AI 학습 데이터 활용', '익명화 데이터로 AI 모델 개선 (Federated Learning)', _aiTrainingConsent, false, 'GDPR Art.6(1)(a)', (v) => setState(() => _aiTrainingConsent = v ?? false)),

            const SizedBox(height: 32),
            OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.history), label: const Text('동의 변경 이력 조회')),
            const SizedBox(height: 8),
            OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.download), label: const Text('내 데이터 내보내기 (FHIR R4 / CSV)')),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () => _showDeleteDialog(theme),
              icon: Icon(Icons.delete_forever, color: theme.colorScheme.error),
              label: Text('계정 및 데이터 삭제', style: TextStyle(color: theme.colorScheme.error)),
              style: OutlinedButton.styleFrom(side: BorderSide(color: theme.colorScheme.error)),
            ),
            const SizedBox(height: 24),
            _buildLegalNotice(theme),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoBanner(ThemeData theme) {
    return Card(
      color: theme.colorScheme.primaryContainer.withOpacity(0.3),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Icon(Icons.info_outline, color: theme.colorScheme.primary),
            const SizedBox(width: 12),
            Expanded(child: Text('각 동의를 확인하고 언제든 변경할 수 있습니다.\n필수 동의 철회 시 서비스 탈퇴로 처리됩니다.', style: theme.textTheme.bodySmall)),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(ThemeData theme, String title, String sub) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(title, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        Text(sub, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
      ]),
    );
  }

  Widget _buildConsentTile(ThemeData theme, String title, String subtitle, bool value, bool required, String legal, ValueChanged<bool?> onChanged) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      child: Column(children: [
        SwitchListTile(
          title: Row(children: [
            Flexible(child: Text(title)),
            if (required) ...[
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                decoration: BoxDecoration(color: theme.colorScheme.error.withOpacity(0.1), borderRadius: BorderRadius.circular(4)),
                child: Text('필수', style: TextStyle(fontSize: 10, color: theme.colorScheme.error)),
              ),
            ],
          ]),
          subtitle: Text(subtitle),
          value: value,
          onChanged: onChanged,
        ),
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
          child: Row(children: [
            Icon(Icons.gavel, size: 12, color: theme.colorScheme.outline),
            const SizedBox(width: 4),
            Flexible(child: Text(legal, style: theme.textTheme.labelSmall?.copyWith(color: theme.colorScheme.outline))),
          ]),
        ),
      ]),
    );
  }

  void _handleRequired(bool? value, String name) {
    if (value == false) {
      showDialog(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('필수 동의 철회'),
          content: Text('$name은(는) 필수 동의입니다.\n철회 시 서비스 탈퇴로 처리됩니다.'),
          actions: [
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
            TextButton(onPressed: () => Navigator.pop(ctx), child: Text('철회', style: TextStyle(color: Theme.of(ctx).colorScheme.error))),
          ],
        ),
      );
    }
  }

  void _showDeleteDialog(ThemeData theme) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('계정 삭제'),
        content: const Text('즉시 비활성화 → 30일 내 PII 삭제\n건강 데이터는 익명화 후 10년 보존 (의료기기법)\n\n되돌릴 수 없습니다.'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          TextButton(onPressed: () => Navigator.pop(ctx), child: Text('삭제 요청', style: TextStyle(color: theme.colorScheme.error))),
        ],
      ),
    );
  }

  Widget _buildLegalNotice(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: theme.colorScheme.surfaceContainerHighest.withOpacity(0.3), borderRadius: BorderRadius.circular(8)),
      child: Text(
        '본 동의 관리는 GDPR(EU), PIPA(한국), HIPAA(미국), PIPL(중국), APPI(일본) 규정을 준수합니다. '
        '동의 기록은 5년간 보존됩니다. 문의: privacy@manpasik.com',
        style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline),
      ),
    );
  }
}
