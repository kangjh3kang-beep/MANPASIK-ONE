import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/utils/navigation_utils.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ConsultationResultScreen — Sanggam Orbit 진료 결과
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* ~8x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 4x → SanggamTheme.primary
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Card 5x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] h:12→16, bottom:4→8
// ───────────────────────────────────────────────────

/// 진료 결과 화면
class ConsultationResultScreen extends ConsumerWidget {
  const ConsultationResultScreen({super.key, required this.consultationId});

  final String consultationId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back, color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.popOrHome(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '진료 결과',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.share, color: Colors.white),
                    tooltip: 'PDF 내보내기',
                    onPressed: () {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('진료 결과 PDF가 생성되었습니다.')),
                      );
                    },
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // 진료 정보
                  SanggamContainer(
                    borderRadius: 16,
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Row(
                          children: [
                            Icon(Icons.medical_services,
                                color: SanggamTheme.primary),
                            SizedBox(width: 8),
                            Text(
                              '진료 정보',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),
                        _infoRow('진료일시', '2024-02-15 14:00'),
                        _infoRow('담당의', '김건강 전문의'),
                        _infoRow('진료과', '내과'),
                        _infoRow('진료유형', '화상진료'),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 진료 소견
                  const SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '진료 소견',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        SizedBox(height: 8),
                        Text(
                          '바이오마커 측정 결과 혈당 수치가 경계 범위에 있습니다. 식이 조절과 규칙적인 운동을 권장합니다. '
                          '2주 후 재측정하여 추이를 확인하시기 바랍니다.',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 처방전
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Padding(
                          padding: EdgeInsets.fromLTRB(16, 16, 16, 8),
                          child: Text(
                            '처방전',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        const ListTile(
                          leading: Icon(Icons.medication,
                              color: SanggamTheme.jagaeCyan),
                          title: Text('메트포르민 500mg'),
                          subtitle: Text('1일 2회, 식후 30분'),
                          trailing: Text('14일분'),
                        ),
                        const Divider(
                            height: 1, color: SanggamTheme.surfaceVariant),
                        const ListTile(
                          leading: Icon(Icons.medication,
                              color: SanggamTheme.jagaeCyan),
                          title: Text('비타민 D 1000IU'),
                          subtitle: Text('1일 1회, 아침'),
                          trailing: Text('30일분'),
                        ),
                        Padding(
                          padding: const EdgeInsets.all(16),
                          child: OutlinedButton.icon(
                            onPressed: () =>
                                context.push('/medical/facility-search'),
                            icon: const Icon(Icons.local_pharmacy),
                            label: const Text('근처 약국 찾기'),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 다음 단계 — 크로스 도메인 연계
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        ListTile(
                          leading: const Icon(Icons.folder_shared,
                              color: SanggamTheme.jagaeCyan),
                          title: const Text('건강기록에서 확인'),
                          subtitle: const Text('진료 결과가 건강기록에 저장되었습니다'),
                          trailing: const Icon(Icons.chevron_right,
                              color: SanggamTheme.onSurfaceDim),
                          onTap: () => context.go('/data'),
                        ),
                        const Divider(
                            height: 1, color: SanggamTheme.surfaceVariant),
                        ListTile(
                          leading: const Icon(Icons.receipt_long,
                              color: SanggamTheme.primary),
                          title: const Text('처방전 상세 보기'),
                          subtitle: const Text('처방 약품 및 복약 안내'),
                          trailing: const Icon(Icons.chevron_right,
                              color: SanggamTheme.onSurfaceDim),
                          onTap: () => context.push(
                              '/medical/prescription/$consultationId'),
                        ),
                        const Divider(
                            height: 1, color: SanggamTheme.surfaceVariant),
                        ListTile(
                          leading: const Icon(Icons.calendar_today,
                              color: SanggamTheme.primary),
                          title: const Text('다음 진료 예약'),
                          subtitle: const Text('2주 후 재진 권장'),
                          trailing: const Icon(Icons.chevron_right,
                              color: SanggamTheme.onSurfaceDim),
                          onTap: () => context.push('/medical/telemedicine'),
                        ),
                        const Divider(
                            height: 1, color: SanggamTheme.surfaceVariant),
                        ListTile(
                          leading: const Icon(Icons.science,
                              color: SanggamTheme.primary),
                          title: const Text('바이오마커 재측정'),
                          subtitle: const Text('2주 후 추이 확인'),
                          trailing: const Icon(Icons.chevron_right,
                              color: SanggamTheme.onSurfaceDim),
                          onTap: () => context.go('/measure'),
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

  Widget _infoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          SizedBox(
            width: 80,
            child: Text(
              label,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 14,
            ),
          ),
        ],
      ),
    );
  }
}
