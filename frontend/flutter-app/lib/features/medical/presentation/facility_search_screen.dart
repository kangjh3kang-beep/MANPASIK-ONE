import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 병원/약국 검색 화면
class FacilitySearchScreen extends ConsumerStatefulWidget {
  const FacilitySearchScreen({super.key});

  @override
  ConsumerState<FacilitySearchScreen> createState() => _FacilitySearchScreenState();
}

class _FacilitySearchScreenState extends ConsumerState<FacilitySearchScreen> {
  final _searchController = TextEditingController();
  String _selectedSpecialty = 'all';
  List<Map<String, dynamic>> _results = [];
  bool _isLoading = false;

  static const _specialties = [
    ('all', '전체'),
    ('internal', '내과'),
    ('cardiology', '심장내과'),
    ('endocrinology', '내분비내과'),
    ('dermatology', '피부과'),
    ('family', '가정의학과'),
    ('pharmacy', '약국'),
  ];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _search() async {
    setState(() => _isLoading = true);
    try {
      final client = ref.read(restClientProvider);
      final queryText = _selectedSpecialty == 'all'
          ? _searchController.text
          : '${_searchController.text} $_selectedSpecialty'.trim();
      final resp = await client.searchFacilities(
        query: queryText,
      );
      final items = resp['facilities'] as List? ?? resp['items'] as List? ?? [];
      setState(() {
        _results = items.cast<Map<String, dynamic>>();
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _results = _fallbackFacilities;
        _isLoading = false;
      });
    }
  }

  @override
  void initState() {
    super.initState();
    _search();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('병원/약국 찾기'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: Column(
        children: [
          // 검색바
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '병원명, 주소, 의사명으로 검색',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: IconButton(
                  icon: const Icon(Icons.search),
                  onPressed: _search,
                ),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
              onSubmitted: (_) => _search(),
            ),
          ),

          // 진료과 필터
          SizedBox(
            height: 36,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: _specialties.map((s) {
                final isSelected = _selectedSpecialty == s.$1;
                return Padding(
                  padding: const EdgeInsets.only(right: 8),
                  child: FilterChip(
                    selected: isSelected,
                    label: Text(s.$2, style: const TextStyle(fontSize: 12)),
                    selectedColor: AppTheme.sanggamGold,
                    onSelected: (_) {
                      setState(() => _selectedSpecialty = s.$1);
                      _search();
                    },
                  ),
                );
              }).toList(),
            ),
          ),
          const SizedBox(height: 8),
          const Divider(height: 1),

          // 결과 리스트
          Expanded(
            child: _isLoading
                ? const Center(child: CircularProgressIndicator())
                : _results.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(Icons.local_hospital_outlined, size: 48, color: theme.colorScheme.onSurfaceVariant),
                            const SizedBox(height: 8),
                            Text('검색 결과가 없습니다.', style: theme.textTheme.bodyMedium),
                          ],
                        ),
                      )
                    : ListView.separated(
                        padding: const EdgeInsets.all(16),
                        itemCount: _results.length,
                        separatorBuilder: (_, __) => const Divider(),
                        itemBuilder: (context, index) => _buildFacilityTile(theme, _results[index]),
                      ),
          ),
        ],
      ),
    );
  }

  Widget _buildFacilityTile(ThemeData theme, Map<String, dynamic> facility) {
    final name = facility['name'] as String? ?? '의료기관';
    final specialty = facility['specialty'] as String? ?? '';
    final address = facility['address'] as String? ?? '';
    final rating = (facility['rating'] as num?)?.toDouble() ?? 4.0;
    final isPharmacy = specialty.contains('약국');

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: isPharmacy ? Colors.green.withOpacity(0.1) : Colors.blue.withOpacity(0.1),
        child: Icon(
          isPharmacy ? Icons.local_pharmacy : Icons.local_hospital,
          color: isPharmacy ? Colors.green : Colors.blue,
        ),
      ),
      title: Text(name, style: const TextStyle(fontWeight: FontWeight.w600)),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (specialty.isNotEmpty) Text(specialty, style: theme.textTheme.bodySmall),
          if (address.isNotEmpty)
            Text(address, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          Row(
            children: [
              const Icon(Icons.star, size: 14, color: Colors.amber),
              const SizedBox(width: 2),
              Text(rating.toStringAsFixed(1), style: theme.textTheme.bodySmall),
            ],
          ),
        ],
      ),
      trailing: FilledButton(
        onPressed: () => _showReservationDialog(facility),
        style: FilledButton.styleFrom(
          backgroundColor: AppTheme.sanggamGold,
          minimumSize: const Size(60, 32),
          padding: const EdgeInsets.symmetric(horizontal: 12),
        ),
        child: const Text('예약', style: TextStyle(fontSize: 12)),
      ),
    );
  }

  void _showReservationDialog(Map<String, dynamic> facility) {
    final name = facility['name'] as String? ?? '의료기관';
    final facilityId = facility['id'] as String? ?? '';
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('$name 예약'),
        content: const Text('해당 기관에 진료 예약을 요청하시겠습니까?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('취소'),
          ),
          FilledButton(
            onPressed: () async {
              Navigator.pop(ctx);
              try {
                final client = ref.read(restClientProvider);
                await client.createReservation(
                  userId: 'current-user',
                  facilityId: facilityId,
                );
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('$name 예약이 요청되었습니다.')),
                  );
                }
              } catch (e) {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('$name 예약 요청에 실패했습니다.')),
                  );
                }
              }
            },
            style: FilledButton.styleFrom(backgroundColor: AppTheme.sanggamGold),
            child: const Text('예약하기'),
          ),
        ],
      ),
    );
  }

  static final _fallbackFacilities = [
    {'name': '서울대학교병원', 'specialty': '내과', 'address': '서울 종로구 대학로 101', 'rating': 4.8},
    {'name': '삼성서울병원', 'specialty': '심장내과', 'address': '서울 강남구 일원로 81', 'rating': 4.7},
    {'name': '연세세브란스병원', 'specialty': '내분비내과', 'address': '서울 서대문구 연세로 50-1', 'rating': 4.6},
    {'name': '서울아산병원', 'specialty': '가정의학과', 'address': '서울 송파구 올림픽로43길 88', 'rating': 4.7},
    {'name': '건강약국', 'specialty': '약국', 'address': '서울 강남구 테헤란로 123', 'rating': 4.5},
  ];
}
