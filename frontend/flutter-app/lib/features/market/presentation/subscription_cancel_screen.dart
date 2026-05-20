import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// SubscriptionCancelScreen — Sanggam Orbit 구독 해지
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] theme.textTheme ~3x → 직접 TextStyle
// [Rule 4] theme.colorScheme 1x → SanggamTheme.error
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] Card 1x → SanggamContainer
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] RadioListTile 6x 다크 테마
// [Rule 4] TextField 다크 테마
// [Rule 4] AlertDialog 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 구독 해지 화면
class SubscriptionCancelScreen
    extends StatefulWidget {
  const SubscriptionCancelScreen(
      {super.key});

  @override
  State<SubscriptionCancelScreen>
      createState() =>
          _SubscriptionCancelScreenState();
}

class _SubscriptionCancelScreenState
    extends State<SubscriptionCancelScreen> {
  String? _selectedReason;
  final _detailController =
      TextEditingController();
  bool _isCancelling = false;

  static const _reasons = [
    '서비스를 더 이상 사용하지 않음',
    '가격이 너무 비쌈',
    '다른 서비스로 전환',
    '필요한 기능이 부족',
    '기기를 분실/교체함',
    '기타',
  ];

  @override
  void dispose() {
    _detailController.dispose();
    super.dispose();
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
                      '구독 해지',
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
                    // 현재 구독 정보
                    SanggamContainer(
                      borderRadius: 16,
                      child: Column(
                        crossAxisAlignment:
                            CrossAxisAlignment
                                .start,
                        children: [
                          const Text(
                              '현재 구독 정보',
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
                              '플랜',
                              'Pro 구독'),
                          _buildInfoRow(
                              '월 요금',
                              '₩29,900'),
                          _buildInfoRow(
                              '다음 결제일',
                              '2026-03-18'),
                          _buildInfoRow(
                              '구독 시작일',
                              '2025-12-18'),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 16),

                    // 혜택 안내
                    SanggamContainer(
                      borderRadius: 16,
                      borderColor:
                          SanggamTheme
                              .primary,
                      jagaeOpacity: 0.15,
                      child: Column(
                        children: [
                          const Icon(
                              Icons
                                  .info_outline,
                              color:
                                  SanggamTheme
                                      .primary,
                              size: 32),
                          const SizedBox(
                              height: 8),
                          const Text(
                              '해지 시 잃게 되는 혜택',
                              style:
                                  TextStyle(
                                color: Colors
                                    .white,
                                fontSize: 14,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                          const SizedBox(
                              height: 8),
                          ...[
                            'AI 건강 코칭 이용',
                            '원격 진료 서비스',
                            '가족 그룹 (최대 6명)',
                            'FHIR 데이터 내보내기',
                            '고급 카트리지 접근',
                          ].map((feature) =>
                              Padding(
                                padding: const EdgeInsets
                                    .symmetric(
                                    vertical:
                                        2),
                                child: Row(
                                  children: [
                                    const Icon(
                                        Icons
                                            .remove_circle_outline,
                                        size:
                                            16,
                                        color: SanggamTheme
                                            .error),
                                    const SizedBox(
                                        width:
                                            8),
                                    Text(
                                        feature,
                                        style: const TextStyle(
                                            color:
                                                Colors.white)),
                                  ],
                                ),
                              )),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 해지 사유 선택
                    const Text(
                        '해지 사유를 선택해주세요',
                        style: TextStyle(
                          color:
                              Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight
                                  .bold,
                        )),
                    const SizedBox(
                        height: 8),
                    ...(_reasons.map(
                        (reason) =>
                            RadioListTile<
                                String>(
                              title: Text(
                                  reason,
                                  style: const TextStyle(
                                      color: Colors
                                          .white)),
                              activeColor:
                                  SanggamTheme
                                      .primary,
                              value: reason,
                              groupValue:
                                  _selectedReason,
                              onChanged: (val) =>
                                  setState(() =>
                                      _selectedReason =
                                          val),
                              contentPadding:
                                  EdgeInsets
                                      .zero,
                            ))),
                    const SizedBox(
                        height: 8),

                    // 상세 의견
                    TextField(
                      controller:
                          _detailController,
                      maxLines: 3,
                      style:
                          const TextStyle(
                              color: Colors
                                  .white),
                      decoration:
                          InputDecoration(
                        labelText:
                            '추가 의견 (선택)',
                        labelStyle: const TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim),
                        hintText:
                            '서비스 개선을 위해 의견을 남겨주세요',
                        hintStyle: const TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim),
                        filled: true,
                        fillColor:
                            SanggamTheme
                                .surface,
                        border: OutlineInputBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                            borderSide:
                                BorderSide
                                    .none),
                        enabledBorder: OutlineInputBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                            borderSide: const BorderSide(
                                color: SanggamTheme
                                    .surfaceVariant)),
                        focusedBorder: OutlineInputBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                            borderSide: const BorderSide(
                                color: SanggamTheme
                                    .primary)),
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 해지 버튼
                    FilledButton(
                      onPressed:
                          _selectedReason !=
                                      null &&
                                  !_isCancelling
                              ? _handleCancel
                              : null,
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
                                48),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                      child: _isCancelling
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child: CircularProgressIndicator(
                                  strokeWidth:
                                      2,
                                  color: Colors
                                      .white))
                          : const Text(
                              '구독 해지하기'),
                    ),
                    const SizedBox(
                        height: 8),

                    // 유지 버튼
                    OutlinedButton(
                      onPressed: () =>
                          context.pop(),
                      style: OutlinedButton
                          .styleFrom(
                        minimumSize:
                            const Size
                                .fromHeight(
                                48),
                        foregroundColor:
                            SanggamTheme
                                .primary,
                        side: const BorderSide(
                            color:
                                SanggamTheme
                                    .primary),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                      child: const Text(
                          '구독 유지하기'),
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
                    .onSurfaceDim,
                fontSize: 13,
              )),
          Text(value,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              )),
        ],
      ),
    );
  }

  void _handleCancel() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor:
            SanggamTheme.surface,
        shape: RoundedRectangleBorder(
            borderRadius:
                BorderRadius.circular(16)),
        title: const Text(
            '정말 해지하시겠습니까?',
            style: TextStyle(
                color: Colors.white)),
        content: const Text(
            '현재 결제 기간이 끝난 후 Pro 혜택이 중단됩니다.\n무료 플랜으로 전환됩니다.',
            style: TextStyle(
                color: SanggamTheme
                    .onSurfaceDim)),
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
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              setState(() =>
                  _isCancelling = true);
              Future.delayed(
                  const Duration(
                      seconds: 2), () {
                if (mounted) {
                  setState(() =>
                      _isCancelling =
                          false);
                  ScaffoldMessenger.of(
                          context)
                      .showSnackBar(
                    const SnackBar(
                        content: Text(
                            '구독이 해지되었습니다. 현재 결제 기간까지 이용 가능합니다.')),
                  );
                  context.pop();
                }
              });
            },
            style:
                FilledButton.styleFrom(
              backgroundColor:
                  SanggamTheme.error,
              foregroundColor:
                  Colors.white,
              shape:
                  RoundedRectangleBorder(
                borderRadius:
                    BorderRadius.circular(
                        16),
              ),
            ),
            child:
                const Text('해지 확인'),
          ),
        ],
      ),
    );
  }
}
