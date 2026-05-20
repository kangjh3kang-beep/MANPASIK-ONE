import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/utils/navigation_utils.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// TelemedicineScreen — Sanggam Orbit 화상진료 예약
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) + ThemeData 파라미터 7x 제거
// [Rule 4] theme.textTheme.* ~15x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~4x → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold ~8x → SanggamTheme.primary
// [Rule 4] Colors.blue 2x → SanggamTheme.jagaeCyan
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.amber → SanggamTheme.primary
// [Rule 4] Card ~5x → SanggamContainer / ListTile
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing:12→16, borderRadius:12→16, vertical:4→8,
//          padding:12→16, h:12→16
// ───────────────────────────────────────────────────

/// 화상진료 예약 화면
///
/// 진료과 선택 → 의사 목록 → 예약 확인 → 대기실
class TelemedicineScreen extends ConsumerStatefulWidget {
  const TelemedicineScreen({super.key});

  @override
  ConsumerState<TelemedicineScreen> createState() => _TelemedicineScreenState();
}

class _TelemedicineScreenState extends ConsumerState<TelemedicineScreen> {
  int _step = 0;
  String? _selectedSpecialty;
  _DoctorInfo? _selectedDoctor;
  List<_DoctorInfo> _doctors = [];
  bool _loadingDoctors = false;
  bool _reserving = false;

  static const _specialties = [
    ('내과', Icons.medical_information, '일반 내과 진료'),
    ('심장내과', Icons.favorite, '심혈관 전문 진료'),
    ('내분비내과', Icons.science, '당뇨/갑상선 전문'),
    ('피부과', Icons.face, '피부 질환 상담'),
    ('가정의학과', Icons.home, '종합 건강 상담'),
    ('정신건강의학과', Icons.psychology, '정신건강 상담'),
  ];

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
                    onPressed: () {
                      if (_step > 0) {
                        setState(() => _step--);
                      } else {
                        context.popOrHome();
                      }
                    },
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      _stepTitle,
                      style: const TextStyle(
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
            Expanded(child: _buildStep()),
          ],
        ),
      ),
    );
  }

  String get _stepTitle => ['비대면 진료', '의사 선택', '예약 확인', '대기실'][_step];

  Widget _buildStep() {
    switch (_step) {
      case 0:
        return _buildSpecialtyStep();
      case 1:
        return _buildDoctorStep();
      case 2:
        return _buildConfirmStep();
      case 3:
        return _buildWaitingStep();
      default:
        return const SizedBox.shrink();
    }
  }

  Widget _buildSpecialtyStep() {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '어떤 진료가 필요하신가요?',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '진료과를 선택해주세요.',
            style: TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
          const SizedBox(height: 16),
          Expanded(
            child: GridView.builder(
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                mainAxisSpacing: 16,
                crossAxisSpacing: 16,
                childAspectRatio: 1.5,
              ),
              itemCount: _specialties.length,
              itemBuilder: (context, index) {
                final s = _specialties[index];
                return SanggamContainer(
                  borderRadius: 16,
                  padding: EdgeInsets.zero,
                  onTap: () => _selectSpecialty(s.$1),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(s.$2, size: 28, color: SanggamTheme.primary),
                        const SizedBox(height: 8),
                        Text(
                          s.$1,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 14,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        Text(
                          s.$3,
                          style: const TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 10,
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _selectSpecialty(String specialty) async {
    setState(() {
      _selectedSpecialty = specialty;
      _step = 1;
      _loadingDoctors = true;
      _doctors = [];
    });
    try {
      final client = ref.read(restClientProvider);
      final res = await client.searchDoctors(specialty: specialty);
      final list = res['doctors'] as List<dynamic>? ?? [];
      if (mounted) {
        setState(() {
          _doctors = list.map((d) {
            final m = d as Map<String, dynamic>;
            return _DoctorInfo(
              id: m['doctor_id'] as String? ?? m['id'] as String? ?? '',
              name: m['name'] as String? ?? m['doctor_name'] as String? ?? '',
              specialty: m['specialty'] as String? ?? specialty,
              hospital: m['hospital'] as String? ??
                  m['facility_name'] as String? ??
                  '',
              rating: (m['rating'] as num?)?.toDouble() ?? 0.0,
              experience: m['experience'] as String? ?? '',
              available: m['available'] as bool? ?? true,
            );
          }).toList();
          _loadingDoctors = false;
        });
      }
    } catch (_) {
      if (mounted) {
        setState(() => _loadingDoctors = false);
      }
    }
  }

  Future<void> _confirmReservation() async {
    if (_selectedDoctor == null) return;
    setState(() => _reserving = true);
    try {
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await client.createConsultation(
        userId: userId,
        doctorId: _selectedDoctor!.id,
        specialty: _selectedSpecialty ?? '',
        reason: '비대면 화상 진료',
      );
      if (mounted) {
        setState(() {
          _step = 3;
          _reserving = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() => _reserving = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('예약 실패: $e')),
        );
      }
    }
  }

  Widget _buildDoctorStep() {
    if (_loadingDoctors) {
      return const Center(child: CircularProgressIndicator());
    }
    if (_doctors.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.person_search,
                size: 48, color: SanggamTheme.onSurfaceDim),
            const SizedBox(height: 16),
            Text(
              '현재 $_selectedSpecialty 전문의를 찾을 수 없습니다.',
              style: const TextStyle(
                color: Colors.white,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 8),
            FilledButton(
              onPressed: () => _selectSpecialty(_selectedSpecialty!),
              style: FilledButton.styleFrom(
                backgroundColor: SanggamTheme.primary,
                foregroundColor: SanggamTheme.background,
              ),
              child: const Text('다시 검색'),
            ),
          ],
        ),
      );
    }
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Text(
          '$_selectedSpecialty 전문의',
          style: const TextStyle(
            color: Colors.white,
            fontSize: 16,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 16),
        ..._doctors.map((d) => Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: ListTile(
                tileColor: SanggamTheme.surface,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
                contentPadding: const EdgeInsets.all(16),
                leading: CircleAvatar(
                  backgroundColor: SanggamTheme.primary.withValues(alpha: 0.15),
                  child: Text(
                    d.name.isNotEmpty ? d.name[0] : '?',
                    style: const TextStyle(
                      color: SanggamTheme.primary,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                title: Row(
                  children: [
                    Text(
                      '${d.name} 전문의',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(width: 8),
                    if (!d.available)
                      Container(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 8, vertical: 2),
                        decoration: BoxDecoration(
                          color: SanggamTheme.onSurfaceDim,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: const Text(
                          '마감',
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.white,
                          ),
                        ),
                      ),
                  ],
                ),
                subtitle: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${d.hospital} | 경력 ${d.experience}',
                      style: const TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 12,
                      ),
                    ),
                    Row(
                      children: [
                        const Icon(Icons.star,
                            size: 14, color: SanggamTheme.primary),
                        Text(
                          ' ${d.rating}',
                          style: const TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
                trailing: FilledButton(
                  onPressed: d.available
                      ? () {
                          setState(() {
                            _selectedDoctor = d;
                            _step = 2;
                          });
                        }
                      : null,
                  style: FilledButton.styleFrom(
                    backgroundColor: SanggamTheme.primary,
                    foregroundColor: SanggamTheme.background,
                    minimumSize: const Size(60, 32),
                    padding: const EdgeInsets.symmetric(horizontal: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                  ),
                  child: const Text('선택', style: TextStyle(fontSize: 12)),
                ),
              ),
            )),
      ],
    );
  }

  Widget _buildConfirmStep() {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          SanggamContainer(
            borderRadius: 16,
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '예약 정보 확인',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const Divider(color: SanggamTheme.surfaceVariant),
                _infoRow('진료과', _selectedSpecialty ?? ''),
                _infoRow('담당의', '${_selectedDoctor?.name ?? ''} 전문의'),
                _infoRow('소속', _selectedDoctor?.hospital ?? ''),
                _infoRow('진료 방식', '비대면 화상 진료'),
                _infoRow('예약 시간', '오늘 14:00~14:30'),
                _infoRow('진료비', '₩30,000 (보험 적용 별도)'),
              ],
            ),
          ),
          const SizedBox(height: 16),
          SanggamContainer(
            borderRadius: 16,
            borderColor: SanggamTheme.jagaeCyan.withValues(alpha: 0.3),
            padding: const EdgeInsets.all(16),
            child: const Row(
              children: [
                Icon(Icons.info_outline,
                    size: 20, color: SanggamTheme.jagaeCyan),
                SizedBox(width: 8),
                Expanded(
                  child: Text(
                    '화상 진료 시 카메라와 마이크 권한이 필요합니다.\n안정적인 Wi-Fi 환경에서 진행해주세요.',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 12,
                    ),
                  ),
                ),
              ],
            ),
          ),
          const Spacer(),
          FilledButton(
            onPressed: _reserving ? null : _confirmReservation,
            style: FilledButton.styleFrom(
              minimumSize: const Size.fromHeight(48),
              backgroundColor: SanggamTheme.primary,
              foregroundColor: SanggamTheme.background,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(16),
              ),
            ),
            child: _reserving
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                        strokeWidth: 2, color: Colors.white),
                  )
                : const Text('예약 확정'),
          ),
        ],
      ),
    );
  }

  Widget _buildWaitingStep() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                color: SanggamTheme.primary.withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.videocam,
                  size: 48, color: SanggamTheme.primary),
            ),
            const SizedBox(height: 24),
            const Text(
              '대기실',
              style: TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '${_selectedDoctor?.name ?? ''} 전문의와의 진료를 준비하고 있습니다.',
              style: const TextStyle(
                color: Colors.white,
                fontSize: 14,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            const CircularProgressIndicator(),
            const SizedBox(height: 16),
            const Text(
              '잠시만 기다려주세요. 곧 연결됩니다.',
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 12,
              ),
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () {
                final roomId = 'room-${DateTime.now().millisecondsSinceEpoch}';
                context.push('/medical/video-call/$roomId');
              },
              icon: const Icon(Icons.videocam),
              label: const Text('진료실 입장'),
              style: FilledButton.styleFrom(
                backgroundColor: SanggamTheme.primary,
                foregroundColor: SanggamTheme.background,
                minimumSize: const Size.fromHeight(48),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
            ),
            const SizedBox(height: 16),
            OutlinedButton(
              onPressed: () => context.popOrHome(),
              child: const Text('대기실 나가기'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _infoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
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
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 14,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _DoctorInfo {
  final String id, name, specialty, hospital, experience;
  final double rating;
  final bool available;
  const _DoctorInfo({
    this.id = '',
    required this.name,
    required this.specialty,
    required this.hospital,
    required this.rating,
    required this.experience,
    required this.available,
  });
}
