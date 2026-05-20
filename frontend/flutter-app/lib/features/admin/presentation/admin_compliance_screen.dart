import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminComplianceScreen — Sanggam Orbit 규제 준수
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* ~4x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 규제 준수 체크리스트 화면
class AdminComplianceScreen extends ConsumerWidget {
  const AdminComplianceScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '규제 준수 현황',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
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
                    const EdgeInsets.all(16),
                children: [
                  // 종합 준수율
                  SanggamContainer(
                    borderRadius: 16,
                    padding:
                        const EdgeInsets.all(16),
                    child: Row(
                      children: [
                        Stack(
                          alignment:
                              Alignment.center,
                          children: [
                            SizedBox(
                              width: 60,
                              height: 60,
                              child:
                                  CircularProgressIndicator(
                                value: 0.92,
                                strokeWidth: 6,
                                color: SanggamTheme
                                    .jagaeCyan,
                                backgroundColor:
                                    SanggamTheme
                                        .jagaeCyan
                                        .withValues(
                                            alpha:
                                                0.1),
                              ),
                            ),
                            const Text('92%',
                                style: TextStyle(
                                  color:
                                      Colors.white,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                  fontSize: 16,
                                )),
                          ],
                        ),
                        const SizedBox(width: 16),
                        const Expanded(
                          child: Column(
                            crossAxisAlignment:
                                CrossAxisAlignment
                                    .start,
                            children: [
                              Text('종합 준수율',
                                  style:
                                      TextStyle(
                                    color: Colors
                                        .white,
                                    fontSize: 16,
                                    fontWeight:
                                        FontWeight
                                            .bold,
                                  )),
                              Text('46/50 항목 충족',
                                  style:
                                      TextStyle(
                                    color: SanggamTheme
                                        .onSurfaceDim,
                                    fontSize: 12,
                                  )),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 규제별 상세
                  ..._regulations.map((reg) =>
                      _buildRegulationCard(reg)),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRegulationCard(
      _RegulationData reg) {
    final progress = reg.passed / reg.total;
    final color = progress >= 0.9
        ? SanggamTheme.jagaeCyan
        : progress >= 0.7
            ? SanggamTheme.primary
            : SanggamTheme.error;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: Theme(
          data: ThemeData.dark().copyWith(
            dividerColor: SanggamTheme.surfaceVariant,
          ),
          child: ExpansionTile(
            leading: Icon(reg.icon,
                color: SanggamTheme.primary),
            title: Text(reg.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w600,
                )),
            subtitle: Row(
              children: [
                Expanded(
                  child:
                      LinearProgressIndicator(
                    value: progress,
                    color: color,
                    backgroundColor: color
                        .withValues(alpha: 0.1),
                  ),
                ),
                const SizedBox(width: 8),
                Text('${reg.passed}/${reg.total}',
                    style: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim,
                      fontSize: 12,
                    )),
              ],
            ),
            children: reg.items
                .map((item) => ListTile(
                      dense: true,
                      leading: Icon(
                        item.passed
                            ? Icons.check_circle
                            : Icons.cancel,
                        size: 18,
                        color: item.passed
                            ? SanggamTheme
                                .jagaeCyan
                            : SanggamTheme.error,
                      ),
                      title: Text(item.label,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 12,
                          )),
                    ))
                .toList(),
          ),
        ),
      ),
    );
  }

  static final _regulations = [
    _RegulationData(
        'GDPR (EU)', Icons.public, 12, 12, [
      _CheckItem('데이터 수집 동의', true),
      _CheckItem('잊힐 권리 구현', true),
      _CheckItem('데이터 이동권', true),
      _CheckItem('DPO 지정', true),
      _CheckItem('데이터 침해 통지 절차', true),
      _CheckItem('프라이버시 영향 평가', true),
    ]),
    _RegulationData(
        'PIPA (한국)', Icons.flag, 10, 10, [
      _CheckItem('개인정보처리방침 공개', true),
      _CheckItem('수집·이용 동의', true),
      _CheckItem('제3자 제공 동의', true),
      _CheckItem('파기 절차', true),
      _CheckItem('안전성 확보 조치', true),
    ]),
    _RegulationData('HIPAA (미국)',
        Icons.health_and_safety, 14, 12, [
      _CheckItem('PHI 암호화 (전송)', true),
      _CheckItem('PHI 암호화 (저장)', true),
      _CheckItem('접근 제어', true),
      _CheckItem('감사 로그', true),
      _CheckItem('비상 접근 절차', false),
      _CheckItem('비즈니스 연관 계약', false),
    ]),
    _RegulationData(
        'ISO 13485', Icons.verified, 10, 8, [
      _CheckItem('품질 관리 시스템', true),
      _CheckItem('설계 및 개발 제어', true),
      _CheckItem('위험 관리 (ISO 14971)', true),
      _CheckItem('추적성', true),
      _CheckItem('CAPA 프로세스', false),
      _CheckItem('내부 감사 일정', false),
    ]),
    _RegulationData(
        'IEC 62304', Icons.code, 8, 8, [
      _CheckItem('소프트웨어 안전 분류', true),
      _CheckItem('개발 계획', true),
      _CheckItem('요구사항 분석', true),
      _CheckItem('검증/확인(V&V)', true),
    ]),
  ];
}

class _RegulationData {
  final String name;
  final IconData icon;
  final int total, passed;
  final List<_CheckItem> items;
  const _RegulationData(this.name, this.icon,
      this.total, this.passed, this.items);
}

class _CheckItem {
  final String label;
  final bool passed;
  const _CheckItem(this.label, this.passed);
}
