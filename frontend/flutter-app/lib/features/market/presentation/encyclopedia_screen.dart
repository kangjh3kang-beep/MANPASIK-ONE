import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 카트리지 백과사전 화면 (storyboard-encyclopedia.md)
///
/// 카테고리별 카트리지 목록, 검색, 상세 스펙 표시
class EncyclopediaScreen extends ConsumerStatefulWidget {
  const EncyclopediaScreen({super.key});

  @override
  ConsumerState<EncyclopediaScreen> createState() => _EncyclopediaScreenState();
}

class _EncyclopediaScreenState extends ConsumerState<EncyclopediaScreen> {
  String _selectedCategory = 'all';
  String _searchQuery = '';
  final _searchController = TextEditingController();

  static const _categories = [
    ('all', '전체', Icons.grid_view, null),
    ('bio', '바이오', Icons.biotech, Color(0xFF4CAF50)),
    ('env', '환경', Icons.eco, Color(0xFF2196F3)),
    ('food', '식품', Icons.restaurant, Color(0xFFFF9800)),
    ('ind', '산업', Icons.factory, Color(0xFF9C27B0)),
  ];

  // 하드코딩 카트리지 데이터 (REST 미연결 시 fallback)
  static final _cartridges = [
    _CartridgeInfo(id: 'BIO-001', name: '혈당 측정', category: 'bio', icon: Icons.bloodtype, description: '공복/식후 혈당 수치를 정밀 측정합니다.', specs: {'측정 범위': '20~600 mg/dL', '정확도': '±5%', '측정 시간': '5초', '샘플량': '0.5μL'}),
    _CartridgeInfo(id: 'BIO-002', name: '콜레스테롤', category: 'bio', icon: Icons.favorite, description: 'HDL/LDL 콜레스테롤 수치를 분석합니다.', specs: {'측정 범위': '100~500 mg/dL', '정확도': '±8%', '측정 시간': '30초', '샘플량': '15μL'}),
    _CartridgeInfo(id: 'BIO-003', name: '요산', category: 'bio', icon: Icons.science, description: '통풍 위험 지표인 요산 수치를 측정합니다.', specs: {'측정 범위': '1.5~20 mg/dL', '정확도': '±6%', '측정 시간': '10초', '샘플량': '1μL'}),
    _CartridgeInfo(id: 'BIO-004', name: 'CRP (염증)', category: 'bio', icon: Icons.local_fire_department, description: 'C반응성 단백질로 체내 염증 수준을 확인합니다.', specs: {'측정 범위': '0.5~200 mg/L', '정확도': '±10%', '측정 시간': '15초', '샘플량': '5μL'}),
    _CartridgeInfo(id: 'ENV-001', name: '수질 분석', category: 'env', icon: Icons.water_drop, description: 'pH, 탁도, 잔류염소 등 수질 지표를 분석합니다.', specs: {'pH 범위': '0~14', '탁도': '0~1000 NTU', '측정 시간': '60초', '샘플량': '5mL'}),
    _CartridgeInfo(id: 'ENV-002', name: '미세먼지', category: 'env', icon: Icons.cloud, description: 'PM2.5, PM10 농도를 측정합니다.', specs: {'PM2.5': '0~500 μg/m³', 'PM10': '0~1000 μg/m³', '측정 시간': '30초', '정확도': '±15%'}),
    _CartridgeInfo(id: 'FOOD-001', name: '식품 신선도', category: 'food', icon: Icons.restaurant, description: '식품의 신선도를 VOC 분석으로 평가합니다.', specs: {'분석 항목': 'TVB-N, 아민류', '측정 시간': '45초', '결과': '신선/보통/주의/위험', '대상': '육류, 해산물'}),
    _CartridgeInfo(id: 'FOOD-002', name: '잔류 농약', category: 'food', icon: Icons.eco, description: '과일/채소의 잔류 농약을 검출합니다.', specs: {'검출 한계': '0.01 ppm', '측정 시간': '120초', '대상 농약': '유기인계, 카바메이트', '샘플량': '1g'}),
    _CartridgeInfo(id: 'IND-001', name: '윤활유 품질', category: 'ind', icon: Icons.oil_barrel, description: '산업용 윤활유의 산화도/오염도를 분석합니다.', specs: {'측정 항목': 'TAN, TBN, 수분', '측정 시간': '90초', '정확도': '±5%', '샘플량': '3mL'}),
  ];

  List<_CartridgeInfo> get _filteredCartridges {
    var list = _cartridges;
    if (_selectedCategory != 'all') {
      list = list.where((c) => c.category == _selectedCategory).toList();
    }
    if (_searchQuery.isNotEmpty) {
      final q = _searchQuery.toLowerCase();
      list = list.where((c) => c.name.toLowerCase().contains(q) || c.description.toLowerCase().contains(q)).toList();
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
    final theme = Theme.of(context);
    final filtered = _filteredCartridges;

    return Scaffold(
      appBar: AppBar(
        title: const Text('카트리지 백과사전'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // 검색바
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '카트리지 검색...',
                prefixIcon: const Icon(Icons.search),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                filled: true,
                fillColor: theme.colorScheme.surfaceContainerHighest.withOpacity(0.3),
                suffixIcon: _searchQuery.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() => _searchQuery = '');
                        },
                      )
                    : null,
              ),
              onChanged: (v) => setState(() => _searchQuery = v),
            ),
          ),
          const SizedBox(height: 8),

          // 카테고리 필터
          SizedBox(
            height: 40,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: _categories.map((cat) {
                final isSelected = _selectedCategory == cat.$1;
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: FilterChip(
                    selected: isSelected,
                    label: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(cat.$3, size: 16, color: isSelected ? Colors.white : cat.$4),
                        const SizedBox(width: 4),
                        Text(cat.$2),
                      ],
                    ),
                    selectedColor: AppTheme.sanggamGold,
                    onSelected: (_) => setState(() => _selectedCategory = cat.$1),
                  ),
                );
              }).toList(),
            ),
          ),
          const SizedBox(height: 8),

          // 결과 헤더
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('${filtered.length}개 카트리지', style: theme.textTheme.bodySmall),
                Text('ManPaSik 공식 카트리지', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ],
            ),
          ),
          const Divider(),

          // 카트리지 목록
          Expanded(
            child: filtered.isEmpty
                ? Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.search_off, size: 48, color: theme.colorScheme.onSurfaceVariant),
                        const SizedBox(height: 8),
                        Text('검색 결과가 없습니다.', style: theme.textTheme.bodyMedium),
                      ],
                    ),
                  )
                : ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    itemCount: filtered.length,
                    itemBuilder: (context, index) => _buildCartridgeCard(theme, filtered[index]),
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildCartridgeCard(ThemeData theme, _CartridgeInfo info) {
    final catColor = _categories.firstWhere((c) => c.$1 == info.category).$4 ?? AppTheme.sanggamGold;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ExpansionTile(
        leading: CircleAvatar(
          backgroundColor: catColor.withOpacity(0.15),
          child: Icon(info.icon, color: catColor, size: 20),
        ),
        title: Text(info.name, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(info.id, style: theme.textTheme.bodySmall),
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
            child: Text(info.description, style: theme.textTheme.bodyMedium),
          ),
          const Divider(indent: 16, endIndent: 16),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
            child: Table(
              columnWidths: const {0: FlexColumnWidth(1), 1: FlexColumnWidth(2)},
              children: info.specs.entries.map((e) => TableRow(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    child: Text(e.key, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600)),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    child: Text(e.value, style: theme.textTheme.bodySmall),
                  ),
                ],
              )).toList(),
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
            child: SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: () => context.push('/market/product/${info.id}'),
                icon: const Icon(Icons.shopping_cart_outlined, size: 18),
                label: const Text('마켓에서 보기'),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _CartridgeInfo {
  final String id, name, category, description;
  final IconData icon;
  final Map<String, String> specs;
  const _CartridgeInfo({required this.id, required this.name, required this.category, required this.icon, required this.description, required this.specs});
}
