import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// ResultScreen — Sanggam Orbit 분석 결과
//
// [Rule 4] +sanggam_theme.dart, +sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] theme.textTheme.* ~12x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~15x → SanggamTheme 상수
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.amber → SanggamTheme.primary
// [Rule 4] withOpacity 3x → withValues(alpha:)
// [Rule 4] Container 2x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius:20→24, borderRadius:100→16
// ───────────────────────────────────────────────────

class ResultScreen extends StatelessWidget {
  const ResultScreen({
    super.key,
    this.value = 0.0,
    this.unit = 'mg/dL',
    this.status = '',
    this.biomarker = '',
  });

  final double value;
  final String unit;
  final String status;
  final String biomarker;

  String get _statusText {
    if (status.isNotEmpty) return status;
    if (value <= 0) return '측정 대기';
    return '분석 완료';
  }

  Color get _statusColor {
    if (status.contains('정상') || status.contains('normal')) {
      return SanggamTheme.jagaeCyan;
    }
    if (status.contains('주의') || status.contains('warning')) {
      return SanggamTheme.primary;
    }
    if (status.contains('위험') || status.contains('critical')) {
      return SanggamTheme.error;
    }
    return SanggamTheme.primary;
  }

  @override
  Widget build(BuildContext context) {
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
                    icon: const Icon(Icons.close,
                        color: Colors.white),
                    tooltip: '닫기',
                    onPressed: () => context.go('/home'),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '분석 결과',
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
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment.stretch,
                  children: [
                    // 1. 결과 요약 카드
                    SanggamContainer(
                      borderRadius: 24,
                      borderColor: SanggamTheme.primary
                          .withValues(alpha: 0.2),
                      padding: const EdgeInsets.all(32),
                      child: Column(
                        children: [
                          Text(
                            DateFormat('yyyy. MM. dd  HH:mm')
                                .format(DateTime.now()),
                            style: const TextStyle(
                              color: SanggamTheme
                                  .onSurfaceDim,
                              fontSize: 14,
                            ),
                          ),
                          const SizedBox(height: 24),
                          Row(
                            mainAxisAlignment:
                                MainAxisAlignment.center,
                            crossAxisAlignment:
                                CrossAxisAlignment.baseline,
                            textBaseline:
                                TextBaseline.alphabetic,
                            children: [
                              Text(
                                value > 0
                                    ? value.toStringAsFixed(1)
                                    : '--',
                                style: const TextStyle(
                                  color: SanggamTheme
                                      .primary,
                                  fontSize: 40,
                                  fontWeight:
                                      FontWeight.bold,
                                ),
                              ),
                              const SizedBox(width: 8),
                              Text(
                                unit,
                                style: const TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim,
                                  fontSize: 24,
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 16),
                          Container(
                            padding:
                                const EdgeInsets.symmetric(
                                    horizontal: 16,
                                    vertical: 8),
                            decoration: BoxDecoration(
                              color: _statusColor,
                              borderRadius:
                                  BorderRadius.circular(
                                      16),
                            ),
                            child: Text(
                              _statusText,
                              style: const TextStyle(
                                color: SanggamTheme
                                    .background,
                                fontSize: 14,
                                fontWeight:
                                    FontWeight.bold,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),

                    const SizedBox(height: 32),

                    // 2. 상세 분석 그리드
                    const Text(
                      '상세 분석',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),

                    GridView.count(
                      shrinkWrap: true,
                      physics:
                          const NeverScrollableScrollPhysics(),
                      crossAxisCount: 2,
                      mainAxisSpacing: 16,
                      crossAxisSpacing: 16,
                      childAspectRatio: 1.5,
                      children: [
                        _buildDetailCard(
                            '수분량',
                            '적정',
                            Icons.water_drop,
                            SanggamTheme.jagaeCyan),
                        _buildDetailCard(
                            '전해질',
                            '주의',
                            Icons.bolt,
                            SanggamTheme.error),
                        _buildDetailCard(
                            '스트레스',
                            '낮음',
                            Icons.mood,
                            SanggamTheme.jagaeCyan),
                        _buildDetailCard(
                            '피로도',
                            '보통',
                            Icons.battery_charging_full,
                            SanggamTheme.primary),
                      ],
                    ),

                    const SizedBox(height: 32),

                    // 3. AI 코칭 메시지
                    SanggamContainer(
                      borderRadius: 24,
                      padding: const EdgeInsets.all(24),
                      child: Row(
                        crossAxisAlignment:
                            CrossAxisAlignment.start,
                        children: [
                          const Icon(Icons.auto_awesome,
                              color: SanggamTheme.primary),
                          const SizedBox(width: 16),
                          Expanded(
                            child: Column(
                              crossAxisAlignment:
                                  CrossAxisAlignment.start,
                              children: const [
                                Text(
                                  'AI 건강 코칭',
                                  style: TextStyle(
                                    color: SanggamTheme
                                        .primary,
                                    fontSize: 14,
                                    fontWeight:
                                        FontWeight.bold,
                                  ),
                                ),
                                SizedBox(height: 8),
                                Text(
                                  '지난주보다 수치가 안정적입니다. 지금처럼 하루 2L 수분 섭취를 유지하세요. 저녁 영양제 섭취도 잊지 마세요!',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 14,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),

                    const SizedBox(height: 32),

                    PrimaryButton(
                      text: '확인완료',
                      onPressed: () =>
                          context.go('/home'),
                    ),
                    const SizedBox(height: 24),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDetailCard(String label, String value,
      IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: SanggamTheme.surface,
        borderRadius: BorderRadius.circular(24),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, size: 20, color: color),
              const Spacer(),
              Text(
                value,
                style: TextStyle(
                  color: color,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
          const Spacer(),
          Text(
            label,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 14,
            ),
          ),
        ],
      ),
    );
  }
}
