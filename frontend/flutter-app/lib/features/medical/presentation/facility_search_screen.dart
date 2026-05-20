import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/utils/navigation_utils.dart';

// ───────────────────────────────────────────────────
// FacilitySearchScreen — Sanggam Orbit 병원/약국 검색
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* ~4x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold 3x → SanggamTheme.primary
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.amber → SanggamTheme.primary
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] FilterChip → Container + BoxDecoration
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius:12→16, padding:12→16
// ───────────────────────────────────────────────────

/// 병원/약국 검색 화면
class FacilitySearchScreen extends ConsumerStatefulWidget {
  const FacilitySearchScreen({super.key});

  @override
  ConsumerState<FacilitySearchScreen> createState() =>
      _FacilitySearchScreenState();
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
                      '병원/약국 찾기',
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

            // 검색바
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: TextField(
                controller: _searchController,
                style: const TextStyle(color: Colors.white),
                decoration: InputDecoration(
                  hintText: '병원명, 주소, 의사명으로 검색',
                  hintStyle: const TextStyle(color: SanggamTheme.onSurfaceDim),
                  prefixIcon: const Icon(Icons.search,
                      color: SanggamTheme.onSurfaceDim),
                  suffixIcon: IconButton(
                    tooltip: '검색',
                    icon: const Icon(Icons.search, color: SanggamTheme.primary),
                    onPressed: _search,
                  ),
                  filled: true,
                  fillColor: SanggamTheme.surface,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: BorderSide.none,
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide:
                        const BorderSide(color: SanggamTheme.surfaceVariant),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: const BorderSide(color: SanggamTheme.primary),
                  ),
                ),
                onSubmitted: (_) => _search(),
              ),
            ),
            const SizedBox(height: 16),

            // 진료과 필터
            SizedBox(
              height: 40,
              child: ListView(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                children: _specialties.map((s) {
                  final isSelected = _selectedSpecialty == s.$1;
                  return Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: GestureDetector(
                      onTap: () {
                        setState(() => _selectedSpecialty = s.$1);
                        _search();
                      },
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 16, vertical: 8),
                        decoration: BoxDecoration(
                          color: isSelected
                              ? SanggamTheme.primary
                              : SanggamTheme.surface,
                          borderRadius: BorderRadius.circular(16),
                          border: Border.all(
                            color: isSelected
                                ? SanggamTheme.primary
                                : SanggamTheme.surfaceVariant,
                          ),
                        ),
                        child: Text(
                          s.$2,
                          style: TextStyle(
                            color: isSelected
                                ? SanggamTheme.background
                                : Colors.white,
                            fontSize: 12,
                            fontWeight: isSelected
                                ? FontWeight.bold
                                : FontWeight.normal,
                          ),
                        ),
                      ),
                    ),
                  );
                }).toList(),
              ),
            ),
            const SizedBox(height: 8),
            const Divider(height: 1, color: SanggamTheme.surfaceVariant),

            // 결과 리스트
            Expanded(
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : _results.isEmpty
                      ? const Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.local_hospital_outlined,
                                  size: 48, color: SanggamTheme.onSurfaceDim),
                              SizedBox(height: 8),
                              Text(
                                '검색 결과가 없습니다.',
                                style: TextStyle(
                                  color: SanggamTheme.onSurfaceDim,
                                  fontSize: 14,
                                ),
                              ),
                            ],
                          ),
                        )
                      : ListView.separated(
                          padding: const EdgeInsets.all(16),
                          itemCount: _results.length,
                          separatorBuilder: (_, __) =>
                              const Divider(color: SanggamTheme.surfaceVariant),
                          itemBuilder: (context, index) =>
                              _buildFacilityTile(_results[index]),
                        ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFacilityTile(Map<String, dynamic> facility) {
    final name = facility['name'] as String? ?? '의료기관';
    final specialty = facility['specialty'] as String? ?? '';
    final address = facility['address'] as String? ?? '';
    final rating = (facility['rating'] as num?)?.toDouble() ?? 4.0;
    final isPharmacy = specialty.contains('약국');

    return ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: isPharmacy
            ? SanggamTheme.jagaeCyan.withValues(alpha: 0.1)
            : SanggamTheme.jagaeCyan.withValues(alpha: 0.1),
        child: Icon(
          isPharmacy ? Icons.local_pharmacy : Icons.local_hospital,
          color: SanggamTheme.jagaeCyan,
        ),
      ),
      title: Text(name,
          style: const TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.w600,
          )),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (specialty.isNotEmpty)
            Text(specialty,
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                )),
          if (address.isNotEmpty)
            Text(address,
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                )),
          Row(
            children: [
              const Icon(Icons.star, size: 14, color: SanggamTheme.primary),
              const SizedBox(width: 4),
              Text(rating.toStringAsFixed(1),
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                  )),
            ],
          ),
        ],
      ),
      trailing: FilledButton(
        onPressed: () => _showReservationDialog(facility),
        style: FilledButton.styleFrom(
          backgroundColor: SanggamTheme.primary,
          foregroundColor: SanggamTheme.background,
          minimumSize: const Size(60, 32),
          padding: const EdgeInsets.symmetric(horizontal: 16),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
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
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
        title: Text('$name 예약', style: const TextStyle(color: Colors.white)),
        content: const Text('해당 기관에 진료 예약을 요청하시겠습니까?',
            style: TextStyle(color: SanggamTheme.onSurfaceDim)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('취소',
                style: TextStyle(color: SanggamTheme.onSurfaceDim)),
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
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.primary,
              foregroundColor: SanggamTheme.background,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16),
              ),
            ),
            child: const Text('예약하기'),
          ),
        ],
      ),
    );
  }

  static final _fallbackFacilities = [
    {
      'name': '서울대학교병원',
      'specialty': '내과',
      'address': '서울 종로구 대학로 101',
      'rating': 4.8,
    },
    {
      'name': '삼성서울병원',
      'specialty': '심장내과',
      'address': '서울 강남구 일원로 81',
      'rating': 4.7,
    },
    {
      'name': '연세세브란스병원',
      'specialty': '내분비내과',
      'address': '서울 서대문구 연세로 50-1',
      'rating': 4.6,
    },
    {
      'name': '서울아산병원',
      'specialty': '가정의학과',
      'address': '서울 송파구 올림픽로43길 88',
      'rating': 4.7,
    },
    {
      'name': '건강약국',
      'specialty': '약국',
      'address': '서울 강남구 테헤란로 123',
      'rating': 4.5,
    },
  ];
}
