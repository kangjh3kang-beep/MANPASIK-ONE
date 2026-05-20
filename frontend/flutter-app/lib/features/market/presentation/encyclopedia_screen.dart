import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// EncyclopediaScreen — Sanggam Orbit 카트리지 백과사전
//
// [Rule 4] app_theme → sanggam_theme + sanggam_container
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData param 1x 제거
// [Rule 4] theme.textTheme ~6x → 직접 TextStyle
// [Rule 4] theme.colorScheme ~2x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 2x → SanggamTheme.primary
// [Rule 4] Card 1x → SanggamContainer
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] 카테고리 Colors → SanggamTheme 상수
// [Rule 4] FilterChip 다크 테마
// [Rule 4] TextField 다크 테마
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// ───────────────────────────────────────────────────

/// 카트리지 백과사전 화면
class EncyclopediaScreen
    extends ConsumerStatefulWidget {
  const EncyclopediaScreen({super.key});

  @override
  ConsumerState<EncyclopediaScreen>
      createState() =>
          _EncyclopediaScreenState();
}

class _EncyclopediaScreenState
    extends ConsumerState<EncyclopediaScreen> {
  String _selectedCategory = 'all';
  String _searchQuery = '';
  final _searchController =
      TextEditingController();

  static const _categories = [
    ('all', '전체', Icons.grid_view, null),
    ('bio', '바이오', Icons.biotech,
        SanggamTheme.jagaeCyan),
    ('env', '환경', Icons.eco,
        SanggamTheme.jagaeCyan),
    ('food', '식품', Icons.restaurant,
        SanggamTheme.primary),
    ('ind', '산업', Icons.factory,
        SanggamTheme.jagaeMagenta),
  ];

  static final _cartridges = [
    _CartridgeInfo(
        id: 'BIO-001',
        name: '혈당 측정',
        category: 'bio',
        icon: Icons.bloodtype,
        description:
            '공복/식후 혈당 수치를 정밀 측정합니다.',
        specs: {
          '측정 범위': '20~600 mg/dL',
          '정확도': '±5%',
          '측정 시간': '5초',
          '샘플량': '0.5μL'
        }),
    _CartridgeInfo(
        id: 'BIO-002',
        name: '콜레스테롤',
        category: 'bio',
        icon: Icons.favorite,
        description:
            'HDL/LDL 콜레스테롤 수치를 분석합니다.',
        specs: {
          '측정 범위': '100~500 mg/dL',
          '정확도': '±8%',
          '측정 시간': '30초',
          '샘플량': '15μL'
        }),
    _CartridgeInfo(
        id: 'BIO-003',
        name: '요산',
        category: 'bio',
        icon: Icons.science,
        description:
            '통풍 위험 지표인 요산 수치를 측정합니다.',
        specs: {
          '측정 범위': '1.5~20 mg/dL',
          '정확도': '±6%',
          '측정 시간': '10초',
          '샘플량': '1μL'
        }),
    _CartridgeInfo(
        id: 'BIO-004',
        name: 'CRP (염증)',
        category: 'bio',
        icon: Icons.local_fire_department,
        description:
            'C반응성 단백질로 체내 염증 수준을 확인합니다.',
        specs: {
          '측정 범위': '0.5~200 mg/L',
          '정확도': '±10%',
          '측정 시간': '15초',
          '샘플량': '5μL'
        }),
    _CartridgeInfo(
        id: 'ENV-001',
        name: '수질 분석',
        category: 'env',
        icon: Icons.water_drop,
        description:
            'pH, 탁도, 잔류염소 등 수질 지표를 분석합니다.',
        specs: {
          'pH 범위': '0~14',
          '탁도': '0~1000 NTU',
          '측정 시간': '60초',
          '샘플량': '5mL'
        }),
    _CartridgeInfo(
        id: 'ENV-002',
        name: '미세먼지',
        category: 'env',
        icon: Icons.cloud,
        description:
            'PM2.5, PM10 농도를 측정합니다.',
        specs: {
          'PM2.5': '0~500 μg/m³',
          'PM10': '0~1000 μg/m³',
          '측정 시간': '30초',
          '정확도': '±15%'
        }),
    _CartridgeInfo(
        id: 'FOOD-001',
        name: '식품 신선도',
        category: 'food',
        icon: Icons.restaurant,
        description:
            '식품의 신선도를 VOC 분석으로 평가합니다.',
        specs: {
          '분석 항목': 'TVB-N, 아민류',
          '측정 시간': '45초',
          '결과': '신선/보통/주의/위험',
          '대상': '육류, 해산물'
        }),
    _CartridgeInfo(
        id: 'FOOD-002',
        name: '잔류 농약',
        category: 'food',
        icon: Icons.eco,
        description:
            '과일/채소의 잔류 농약을 검출합니다.',
        specs: {
          '검출 한계': '0.01 ppm',
          '측정 시간': '120초',
          '대상 농약': '유기인계, 카바메이트',
          '샘플량': '1g'
        }),
    _CartridgeInfo(
        id: 'IND-001',
        name: '윤활유 품질',
        category: 'ind',
        icon: Icons.oil_barrel,
        description:
            '산업용 윤활유의 산화도/오염도를 분석합니다.',
        specs: {
          '측정 항목': 'TAN, TBN, 수분',
          '측정 시간': '90초',
          '정확도': '±5%',
          '샘플량': '3mL'
        }),
  ];

  List<_CartridgeInfo>
      get _filteredCartridges {
    var list = _cartridges;
    if (_selectedCategory != 'all') {
      list = list
          .where((c) =>
              c.category ==
              _selectedCategory)
          .toList();
    }
    if (_searchQuery.isNotEmpty) {
      final q =
          _searchQuery.toLowerCase();
      list = list
          .where((c) =>
              c.name
                  .toLowerCase()
                  .contains(q) ||
              c.description
                  .toLowerCase()
                  .contains(q))
          .toList();
    }
    return list;
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final filtered = _filteredCartridges;

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
                      '카트리지 백과사전',
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

            // 검색바
            Padding(
              padding:
                  const EdgeInsets.fromLTRB(
                      16, 0, 16, 0),
              child: TextField(
                controller:
                    _searchController,
                style: const TextStyle(
                    color: Colors.white),
                decoration: InputDecoration(
                  hintText: '카트리지 검색...',
                  hintStyle: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim),
                  prefixIcon: const Icon(
                      Icons.search,
                      color: SanggamTheme
                          .onSurfaceDim),
                  filled: true,
                  fillColor:
                      SanggamTheme.surface,
                  border:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        BorderSide.none,
                  ),
                  enabledBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .surfaceVariant),
                  ),
                  focusedBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(
                            16),
                    borderSide:
                        const BorderSide(
                            color: SanggamTheme
                                .primary),
                  ),
                  suffixIcon: _searchQuery
                          .isNotEmpty
                      ? IconButton(
                          icon: const Icon(
                              Icons.clear,
                              color: SanggamTheme
                                  .onSurfaceDim),
                          tooltip: '검색 초기화',
                          onPressed: () {
                            _searchController
                                .clear();
                            setState(() =>
                                _searchQuery =
                                    '');
                          },
                        )
                      : null,
                ),
                onChanged: (v) =>
                    setState(() =>
                        _searchQuery = v),
              ),
            ),
            const SizedBox(height: 8),

            // 카테고리 필터
            SizedBox(
              height: 40,
              child: ListView(
                scrollDirection:
                    Axis.horizontal,
                padding:
                    const EdgeInsets.symmetric(
                        horizontal: 16),
                children:
                    _categories.map((cat) {
                  final isSelected =
                      _selectedCategory ==
                          cat.$1;
                  return Padding(
                    padding:
                        const EdgeInsets.only(
                            right: 8),
                    child: FilterChip(
                      selected: isSelected,
                      label: Row(
                        mainAxisSize:
                            MainAxisSize.min,
                        children: [
                          Icon(cat.$3,
                              size: 16,
                              color: isSelected
                                  ? SanggamTheme
                                      .background
                                  : cat.$4 ??
                                      SanggamTheme
                                          .onSurfaceDim),
                          const SizedBox(
                              width: 4),
                          Text(cat.$2,
                              style: TextStyle(
                                  color: isSelected
                                      ? SanggamTheme
                                          .background
                                      : Colors
                                          .white)),
                        ],
                      ),
                      selectedColor:
                          SanggamTheme
                              .primary,
                      backgroundColor:
                          SanggamTheme
                              .surface,
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
                      onSelected: (_) =>
                          setState(() =>
                              _selectedCategory =
                                  cat.$1),
                    ),
                  );
                }).toList(),
              ),
            ),
            const SizedBox(height: 8),

            // 결과 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 16),
              child: Row(
                mainAxisAlignment:
                    MainAxisAlignment
                        .spaceBetween,
                children: [
                  Text(
                      '${filtered.length}개 카트리지',
                      style:
                          const TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      )),
                  const Text(
                      'ManPaSik 공식 카트리지',
                      style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      )),
                ],
              ),
            ),
            const Divider(
                color: SanggamTheme
                    .surfaceVariant),

            // 카트리지 목록
            Expanded(
              child: filtered.isEmpty
                  ? Center(
                      child: Column(
                        mainAxisSize:
                            MainAxisSize
                                .min,
                        children: [
                          const Icon(
                              Icons
                                  .search_off,
                              size: 48,
                              color: SanggamTheme
                                  .onSurfaceDim),
                          const SizedBox(
                              height: 8),
                          const Text(
                              '검색 결과가 없습니다.',
                              style:
                                  TextStyle(
                                color: Colors
                                    .white,
                                fontSize:
                                    14,
                              )),
                        ],
                      ),
                    )
                  : ListView.builder(
                      padding:
                          const EdgeInsets
                              .symmetric(
                              horizontal:
                                  16),
                      itemCount:
                          filtered.length,
                      itemBuilder:
                          (context,
                                  index) =>
                              _buildCartridgeCard(
                                  filtered[
                                      index]),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCartridgeCard(
      _CartridgeInfo info) {
    final catColor = _categories
            .firstWhere(
                (c) => c.$1 == info.category)
            .$4 ??
        SanggamTheme.primary;

    return Padding(
      padding: const EdgeInsets.only(
          bottom: 16),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: Theme(
          data: ThemeData.dark().copyWith(
            dividerColor:
                SanggamTheme.surfaceVariant,
          ),
          child: ExpansionTile(
            leading: CircleAvatar(
              backgroundColor: catColor
                  .withValues(alpha: 0.15),
              child: Icon(info.icon,
                  color: catColor,
                  size: 20),
            ),
            title: Text(info.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight:
                      FontWeight.w600,
                )),
            subtitle: Text(info.id,
                style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
            children: [
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 8),
                child: Text(
                    info.description,
                    style:
                        const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                    )),
              ),
              const Divider(
                  indent: 16,
                  endIndent: 16,
                  color: SanggamTheme
                      .surfaceVariant),
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 16),
                child: Table(
                  columnWidths: const {
                    0: FlexColumnWidth(1),
                    1: FlexColumnWidth(2),
                  },
                  children: info
                      .specs.entries
                      .map((e) => TableRow(
                            children: [
                              Padding(
                                padding: const EdgeInsets
                                    .symmetric(
                                    vertical:
                                        4),
                                child: Text(
                                    e.key,
                                    style:
                                        const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim,
                                      fontSize:
                                          13,
                                      fontWeight:
                                          FontWeight
                                              .w600,
                                    )),
                              ),
                              Padding(
                                padding: const EdgeInsets
                                    .symmetric(
                                    vertical:
                                        4),
                                child: Text(
                                    e.value,
                                    style:
                                        const TextStyle(
                                      color:
                                          Colors.white,
                                      fontSize:
                                          13,
                                    )),
                              ),
                            ],
                          ))
                      .toList(),
                ),
              ),
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 16),
                child: SizedBox(
                  width: double.infinity,
                  child:
                      OutlinedButton.icon(
                    onPressed: () =>
                        context.push(
                            '/market/product/${info.id}'),
                    icon: const Icon(
                        Icons
                            .shopping_cart_outlined,
                        size: 18),
                    label: const Text(
                        '마켓에서 보기'),
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
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _CartridgeInfo {
  final String id, name, category,
      description;
  final IconData icon;
  final Map<String, String> specs;
  const _CartridgeInfo({
    required this.id,
    required this.name,
    required this.category,
    required this.icon,
    required this.description,
    required this.specs,
  });
}
