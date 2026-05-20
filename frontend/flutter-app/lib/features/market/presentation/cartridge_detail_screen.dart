import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/cartridge_3d_viewer.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// CartridgeDetailScreen — Sanggam Orbit 카트리지 상세
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] theme.textTheme ~8x → 직접 TextStyle
// [Rule 4] theme.colorScheme 1x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 1x → SanggamTheme.primary
// [Rule 4] Card ~5x → SanggamContainer
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] 데이터 Colors → SanggamTheme 상수
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 카트리지 상세 화면 (백과사전 → 상세)
class CartridgeDetailScreen
    extends ConsumerWidget {
  const CartridgeDetailScreen(
      {super.key,
      required this.cartridgeId});

  final String cartridgeId;

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    final info =
        _getCartridgeById(cartridgeId);

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
                  Expanded(
                    child: Text(
                      info.name,
                      style: const TextStyle(
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
                    const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .start,
                  children: [
                    // 카트리지 3D 뷰어
                    Cartridge3DViewer(
                      height: 220,
                      primaryColor:
                          info.color,
                      label: info.name,
                    ),
                    const SizedBox(
                        height: 24),

                    // 카트리지 헤더
                    Center(
                      child: Column(
                        children: [
                          const SizedBox(
                              height: 16),
                          Text(info.name,
                              style:
                                  const TextStyle(
                                color: Colors
                                    .white,
                                fontSize: 22,
                                fontWeight:
                                    FontWeight
                                        .bold,
                              )),
                          Text(info.code,
                              style:
                                  const TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim,
                                fontSize: 14,
                              )),
                          const SizedBox(
                              height: 8),
                          Container(
                            padding:
                                const EdgeInsets
                                    .symmetric(
                                    horizontal:
                                        16,
                                    vertical:
                                        4),
                            decoration:
                                BoxDecoration(
                              color: info
                                  .color
                                  .withValues(
                                      alpha:
                                          0.1),
                              borderRadius:
                                  BorderRadius
                                      .circular(
                                          16),
                            ),
                            child: Text(
                                info
                                    .category,
                                style:
                                    TextStyle(
                                  color: info
                                      .color,
                                  fontWeight:
                                      FontWeight
                                          .bold,
                                )),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 32),

                    // 설명
                    const Text('측정 원리',
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
                    SanggamContainer(
                      borderRadius: 16,
                      child: Text(
                          info.principle,
                          style:
                              const TextStyle(
                            color: Colors
                                .white,
                            fontSize: 14,
                          )),
                    ),
                    const SizedBox(
                        height: 24),

                    // 스펙 테이블
                    const Text('상세 스펙',
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
                    SanggamContainer(
                      borderRadius: 16,
                      child: Table(
                        columnWidths: const {
                          0: FlexColumnWidth(
                              1),
                          1: FlexColumnWidth(
                              2),
                        },
                        children: info
                            .specs.entries
                            .map((e) =>
                                TableRow(
                                  children: [
                                    Padding(
                                      padding: const EdgeInsets
                                          .symmetric(
                                          vertical:
                                              8),
                                      child: Text(
                                          e.key,
                                          style:
                                              const TextStyle(
                                            color:
                                                SanggamTheme.onSurfaceDim,
                                            fontSize:
                                                14,
                                            fontWeight:
                                                FontWeight.w600,
                                          )),
                                    ),
                                    Padding(
                                      padding: const EdgeInsets
                                          .symmetric(
                                          vertical:
                                              8),
                                      child: Text(
                                          e.value,
                                          style:
                                              const TextStyle(
                                            color:
                                                Colors.white,
                                            fontSize:
                                                14,
                                          )),
                                    ),
                                  ],
                                ))
                            .toList(),
                      ),
                    ),
                    const SizedBox(
                        height: 24),

                    // 활용 사례
                    const Text('활용 사례',
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
                    ...info.useCases.map(
                        (u) => Padding(
                              padding:
                                  const EdgeInsets
                                      .only(
                                      bottom:
                                          8),
                              child:
                                  SanggamContainer(
                                borderRadius:
                                    16,
                                padding:
                                    EdgeInsets
                                        .zero,
                                child:
                                    ListTile(
                                  leading: Icon(
                                      Icons
                                          .check_circle_outline,
                                      color: info
                                          .color),
                                  title: Text(
                                      u,
                                      style: const TextStyle(
                                          color:
                                              Colors.white)),
                                ),
                              ),
                            )),
                    const SizedBox(
                        height: 32),

                    // 마켓 구매 버튼
                    SizedBox(
                      width: double.infinity,
                      height: 56,
                      child:
                          FilledButton.icon(
                        onPressed: () =>
                            context.push(
                                '/market/product/$cartridgeId'),
                        icon: const Icon(Icons
                            .shopping_cart),
                        label: const Text(
                            '마켓에서 구매하기'),
                        style: FilledButton
                            .styleFrom(
                          backgroundColor:
                              SanggamTheme
                                  .primary,
                          foregroundColor:
                              SanggamTheme
                                  .background,
                          shape:
                              RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(
                        height: 16),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  _CartridgeDetail _getCartridgeById(
      String id) {
    return _cartridgeDb[id] ??
        _cartridgeDb.values.first;
  }

  static final _cartridgeDb =
      <String, _CartridgeDetail>{
    'BIO-001': _CartridgeDetail(
      name: '혈당 측정',
      code: 'BIO-001',
      category: '바이오',
      icon: Icons.bloodtype,
      color: SanggamTheme.error,
      principle:
          '전기화학적 방식으로 혈액 내 포도당 농도를 측정합니다. '
          '글루코스 산화효소(GOD)가 포함된 전극에 혈액 샘플이 접촉하면 '
          '전류 변화를 감지하여 정밀한 수치를 산출합니다.',
      specs: {
        '측정 범위': '20~600 mg/dL',
        '정확도': '±5%',
        '측정 시간': '5초',
        '샘플량': '0.5μL',
        '보관 온도': '4~30°C',
        '유효기간': '개봉 후 3개월',
      },
      useCases: [
        '공복 혈당 모니터링',
        '식후 혈당 추적',
        '당뇨 관리 지표',
        '임신성 당뇨 스크리닝',
      ],
    ),
    'BIO-002': _CartridgeDetail(
      name: '콜레스테롤',
      code: 'BIO-002',
      category: '바이오',
      icon: Icons.favorite,
      color: SanggamTheme.jagaeMagenta,
      principle:
          '효소적 비색법을 이용하여 총 콜레스테롤, HDL, LDL 수치를 분석합니다. '
          '콜레스테롤 에스테라아제와 산화효소의 반응으로 생성되는 색소를 광학적으로 측정합니다.',
      specs: {
        '측정 범위': '100~500 mg/dL',
        '정확도': '±8%',
        '측정 시간': '30초',
        '샘플량': '15μL',
        '분석 항목': 'TC/HDL/LDL',
      },
      useCases: [
        '심혈관 질환 위험도 평가',
        '지질 이상증 모니터링',
        '식이요법 효과 확인',
        '약물 치료 모니터링',
      ],
    ),
    'BIO-003': _CartridgeDetail(
      name: '요산',
      code: 'BIO-003',
      category: '바이오',
      icon: Icons.science,
      color: SanggamTheme.primary,
      principle:
          '유리카제(Uricase) 효소법을 사용하여 혈중 요산 농도를 측정합니다. '
          '요산이 효소에 의해 분해될 때 발생하는 과산화수소를 전기화학적으로 검출합니다.',
      specs: {
        '측정 범위': '1.5~20 mg/dL',
        '정확도': '±6%',
        '측정 시간': '10초',
        '샘플량': '1μL',
      },
      useCases: [
        '통풍 위험 평가',
        '신장 기능 모니터링',
        '식이 관리 효과 확인',
      ],
    ),
    'ENV-001': _CartridgeDetail(
      name: '수질 분석',
      code: 'ENV-001',
      category: '환경',
      icon: Icons.water_drop,
      color: SanggamTheme.jagaeCyan,
      principle:
          '다중 전극 시스템으로 pH, 탁도, 잔류염소 등 수질 지표를 동시에 분석합니다. '
          '이온 선택성 전극(ISE)과 광학 센서를 조합하여 정밀한 수질 데이터를 제공합니다.',
      specs: {
        'pH 범위': '0~14',
        '탁도': '0~1000 NTU',
        '잔류염소': '0~10 ppm',
        '측정 시간': '60초',
        '샘플량': '5mL',
      },
      useCases: [
        '가정용 수돗물 검사',
        '수영장 수질 관리',
        '지하수 오염 모니터링',
        '음용수 안전 확인',
      ],
    ),
  };
}

class _CartridgeDetail {
  final String name,
      code,
      category,
      principle;
  final IconData icon;
  final Color color;
  final Map<String, String> specs;
  final List<String> useCases;
  const _CartridgeDetail({
    required this.name,
    required this.code,
    required this.category,
    required this.icon,
    required this.color,
    required this.principle,
    required this.specs,
    required this.useCases,
  });
}
