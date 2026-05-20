import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// EscalationScreen — Sanggam Orbit 119 에스컬레이션
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] theme.colorScheme ~5x → SanggamTheme 상수
// [Rule 4] theme.textTheme ~4x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold 1x → SanggamTheme.primary
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius 12→16, spacing 12→16
// ───────────────────────────────────────────────────

/// 119 에스컬레이션 (응급 상황 처리) UI
class EscalationScreen extends StatefulWidget {
  const EscalationScreen({super.key});

  @override
  State<EscalationScreen> createState() =>
      _EscalationScreenState();
}

class _EscalationScreenState
    extends State<EscalationScreen> {
  bool _isEscalating = false;

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
                      '119 에스컬레이션',
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
                          .stretch,
                  children: [
                    // 경고 배너
                    Container(
                      padding:
                          const EdgeInsets.all(
                              16),
                      decoration:
                          BoxDecoration(
                        color: SanggamTheme
                            .error
                            .withValues(
                                alpha: 0.1),
                        borderRadius:
                            BorderRadius
                                .circular(16),
                        border: Border.all(
                            color:
                                SanggamTheme
                                    .error),
                      ),
                      child: Column(
                        children: [
                          const Icon(
                              Icons
                                  .warning_amber_rounded,
                              size: 48,
                              color:
                                  SanggamTheme
                                      .error),
                          const SizedBox(
                              height: 8),
                          const Text(
                            '위험 수치 감지',
                            style: TextStyle(
                              color:
                                  SanggamTheme
                                      .error,
                              fontSize: 20,
                              fontWeight:
                                  FontWeight
                                      .bold,
                            ),
                          ),
                          const SizedBox(
                              height: 4),
                          const Text(
                            '측정값이 설정된 위험 범위를 초과했습니다.',
                            style: TextStyle(
                              color: Colors
                                  .white,
                              fontSize: 14,
                            ),
                            textAlign:
                                TextAlign
                                    .center,
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 측정 정보 카드
                    SanggamContainer(
                      borderRadius: 16,
                      padding:
                          const EdgeInsets
                              .all(16),
                      child: Column(
                        crossAxisAlignment:
                            CrossAxisAlignment
                                .start,
                        children: [
                          const Text('측정 정보',
                              style:
                                  TextStyle(
                                color: Colors
                                    .white,
                                fontSize: 16,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                          const Divider(
                              color: SanggamTheme
                                  .surfaceVariant),
                          _buildInfoRow(
                              '측정 항목',
                              '혈당'),
                          _buildInfoRow(
                              '측정값',
                              '350 mg/dL'),
                          _buildInfoRow(
                              '정상 범위',
                              '70 ~ 180 mg/dL'),
                          _buildInfoRow(
                              '측정 시각',
                              '방금 전'),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 16),

                    // 긴급 연락처 카드
                    SanggamContainer(
                      borderRadius: 16,
                      padding:
                          const EdgeInsets
                              .all(16),
                      child: Column(
                        crossAxisAlignment:
                            CrossAxisAlignment
                                .start,
                        children: [
                          const Text(
                              '긴급 연락처',
                              style:
                                  TextStyle(
                                color: Colors
                                    .white,
                                fontSize: 16,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                          const Divider(
                              color: SanggamTheme
                                  .surfaceVariant),
                          _buildContactTile(
                              '119 구급대',
                              '119',
                              Icons
                                  .local_hospital,
                              true),
                          _buildContactTile(
                              '보호자 1 (배우자)',
                              '010-1234-5678',
                              Icons.person,
                              false),
                          _buildContactTile(
                              '보호자 2 (자녀)',
                              '010-9876-5432',
                              Icons.person,
                              false),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 119 신고 버튼
                    FilledButton.icon(
                      onPressed:
                          _isEscalating
                              ? null
                              : _handleEscalate,
                      icon: _isEscalating
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                  strokeWidth:
                                      2,
                                  color: Colors
                                      .white))
                          : const Icon(
                              Icons.phone),
                      label: Text(
                          _isEscalating
                              ? '신고 중...'
                              : '119 긴급 신고'),
                      style: FilledButton
                          .styleFrom(
                        backgroundColor:
                            SanggamTheme
                                .error,
                        foregroundColor:
                            Colors.white,
                        minimumSize:
                            const Size
                                .fromHeight(
                                56),
                        textStyle:
                            const TextStyle(
                          fontSize: 18,
                          fontWeight:
                              FontWeight.bold,
                        ),
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
                        height: 16),

                    // 보호자 알림 버튼
                    OutlinedButton.icon(
                      onPressed: () =>
                          _notifyGuardians(
                              context),
                      icon: const Icon(Icons
                          .notification_important),
                      label: const Text(
                          '보호자에게 알림 보내기'),
                      style: OutlinedButton
                          .styleFrom(
                        foregroundColor:
                            SanggamTheme
                                .primary,
                        side: const BorderSide(
                            color:
                                SanggamTheme
                                    .primary),
                        minimumSize:
                            const Size
                                .fromHeight(
                                48),
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
                        height: 16),

                    // 상황 종료 버튼
                    TextButton(
                      onPressed: () =>
                          context.pop(),
                      style:
                          TextButton.styleFrom(
                        foregroundColor:
                            SanggamTheme
                                .onSurfaceDim,
                      ),
                      child: const Text(
                          '상황 종료 (위험하지 않음)'),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(
      String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 4),
      child: Row(
        mainAxisAlignment:
            MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim)),
          Text(value,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              )),
        ],
      ),
    );
  }

  Widget _buildContactTile(String name,
      String number, IconData icon,
      bool isPrimary) {
    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: isPrimary
            ? SanggamTheme.error
                .withValues(alpha: 0.2)
            : SanggamTheme.surfaceVariant,
        child: Icon(icon,
            color: isPrimary
                ? SanggamTheme.error
                : SanggamTheme
                    .onSurfaceDim),
      ),
      title: Text(name,
          style: const TextStyle(
              color: Colors.white)),
      subtitle: Text(number,
          style: const TextStyle(
            color:
                SanggamTheme.onSurfaceDim,
            fontSize: 12,
          )),
      trailing: IconButton(
        icon: const Icon(Icons.phone,
            color: SanggamTheme.primary),
        tooltip: '$name 전화',
        onPressed: () {},
      ),
    );
  }

  void _handleEscalate() {
    setState(() => _isEscalating = true);
    Future.delayed(
        const Duration(seconds: 2), () {
      if (mounted) {
        setState(
            () => _isEscalating = false);
        ScaffoldMessenger.of(context)
            .showSnackBar(
          const SnackBar(
            content: Text(
                '119 긴급 신고가 접수되었습니다.'),
            backgroundColor:
                SanggamTheme.error,
          ),
        );
      }
    });
  }

  void _notifyGuardians(
      BuildContext context) {
    ScaffoldMessenger.of(context)
        .showSnackBar(
      const SnackBar(
          content: Text(
              '보호자에게 알림을 보냈습니다.')),
    );
  }
}
