import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 화상진료 예약 화면 (storyboard-telemedicine.md)
///
/// 진료과 선택 → 의사 목록 → 예약 확인 → 대기실
/// 실제 WebRTC는 Phase 4에서 활성화 예정; 현재는 REST API 기반 UI 플로우
class TelemedicineScreen extends ConsumerStatefulWidget {
  const TelemedicineScreen({super.key});

  @override
  ConsumerState<TelemedicineScreen> createState() => _TelemedicineScreenState();
}

class _TelemedicineScreenState extends ConsumerState<TelemedicineScreen> {
  int _step = 0; // 0: 진료과, 1: 의사, 2: 예약확인, 3: 대기실
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
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_stepTitle),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            if (_step > 0) {
              setState(() => _step--);
            } else {
              context.pop();
            }
          },
        ),
      ),
      body: _buildStep(theme),
    );
  }

  String get _stepTitle => ['비대면 진료', '의사 선택', '예약 확인', '대기실'][_step];

  Widget _buildStep(ThemeData theme) {
    switch (_step) {
      case 0: return _buildSpecialtyStep(theme);
      case 1: return _buildDoctorStep(theme);
      case 2: return _buildConfirmStep(theme);
      case 3: return _buildWaitingStep(theme);
      default: return const SizedBox.shrink();
    }
  }

  Widget _buildSpecialtyStep(ThemeData theme) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('어떤 진료가 필요하신가요?', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 4),
          Text('진료과를 선택해주세요.', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          const SizedBox(height: 16),
          Expanded(
            child: GridView.builder(
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 1.5,
              ),
              itemCount: _specialties.length,
              itemBuilder: (context, index) {
                final s = _specialties[index];
                return Card(
                  child: InkWell(
                    borderRadius: BorderRadius.circular(12),
                    onTap: () => _selectSpecialty(s.$1),
                    child: Padding(
                      padding: const EdgeInsets.all(12),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(s.$2, size: 28, color: AppTheme.sanggamGold),
                          const SizedBox(height: 8),
                          Text(s.$1, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                          Text(s.$3, style: theme.textTheme.bodySmall?.copyWith(fontSize: 10, color: theme.colorScheme.onSurfaceVariant)),
                        ],
                      ),
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
              hospital: m['hospital'] as String? ?? m['facility_name'] as String? ?? '',
              rating: (m['rating'] as num?)?.toDouble() ?? 0.0,
              experience: m['experience'] as String? ?? '',
              available: m['available'] as bool? ?? true,
            );
          }).toList();
          _loadingDoctors = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _loadingDoctors = false);
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
      if (mounted) setState(() { _step = 3; _reserving = false; });
    } catch (e) {
      if (mounted) {
        setState(() => _reserving = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('예약 실패: $e')),
        );
      }
    }
  }

  Widget _buildDoctorStep(ThemeData theme) {
    if (_loadingDoctors) {
      return const Center(child: CircularProgressIndicator());
    }
    if (_doctors.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.person_search, size: 48, color: theme.colorScheme.onSurfaceVariant),
            const SizedBox(height: 12),
            Text('현재 $_selectedSpecialty 전문의를 찾을 수 없습니다.', style: theme.textTheme.bodyMedium),
            const SizedBox(height: 8),
            FilledButton(
              onPressed: () => _selectSpecialty(_selectedSpecialty!),
              child: const Text('다시 검색'),
            ),
          ],
        ),
      );
    }
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Text('$_selectedSpecialty 전문의', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 12),
        ..._doctors.map((d) => Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: AppTheme.sanggamGold.withValues(alpha:0.15),
              child: Text(d.name[0], style: const TextStyle(fontWeight: FontWeight.bold)),
            ),
            title: Row(
              children: [
                Text('${d.name} 전문의', style: const TextStyle(fontWeight: FontWeight.w600)),
                const SizedBox(width: 8),
                if (!d.available)
                  Chip(
                    label: const Text('마감', style: TextStyle(fontSize: 10, color: Colors.white)),
                    backgroundColor: Colors.grey,
                    side: BorderSide.none,
                    visualDensity: VisualDensity.compact,
                  ),
              ],
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('${d.hospital} | 경력 ${d.experience}', style: theme.textTheme.bodySmall),
                Row(
                  children: [
                    const Icon(Icons.star, size: 14, color: Colors.amber),
                    Text(' ${d.rating}', style: theme.textTheme.bodySmall),
                  ],
                ),
              ],
            ),
            trailing: FilledButton(
              onPressed: d.available ? () {
                setState(() {
                  _selectedDoctor = d;
                  _step = 2;
                });
              } : null,
              style: FilledButton.styleFrom(
                backgroundColor: AppTheme.sanggamGold,
                minimumSize: const Size(60, 32),
                padding: const EdgeInsets.symmetric(horizontal: 12),
              ),
              child: const Text('선택', style: TextStyle(fontSize: 12)),
            ),
          ),
        )),
      ],
    );
  }

  Widget _buildConfirmStep(ThemeData theme) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('예약 정보 확인', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                  const Divider(),
                  _infoRow(theme, '진료과', _selectedSpecialty ?? ''),
                  _infoRow(theme, '담당의', '${_selectedDoctor?.name ?? ''} 전문의'),
                  _infoRow(theme, '소속', _selectedDoctor?.hospital ?? ''),
                  _infoRow(theme, '진료 방식', '비대면 화상 진료'),
                  _infoRow(theme, '예약 시간', '오늘 14:00~14:30'),
                  _infoRow(theme, '진료비', '₩30,000 (보험 적용 별도)'),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Card(
            color: Colors.blue.withValues(alpha:0.05),
            child: const Padding(
              padding: EdgeInsets.all(12),
              child: Row(
                children: [
                  Icon(Icons.info_outline, size: 20, color: Colors.blue),
                  SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      '화상 진료 시 카메라와 마이크 권한이 필요합니다.\n안정적인 Wi-Fi 환경에서 진행해주세요.',
                      style: TextStyle(fontSize: 12),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const Spacer(),
          FilledButton(
            onPressed: _reserving ? null : _confirmReservation,
            style: FilledButton.styleFrom(
              minimumSize: const Size.fromHeight(48),
              backgroundColor: AppTheme.sanggamGold,
            ),
            child: _reserving
                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                : const Text('예약 확정'),
          ),
        ],
      ),
    );
  }

  Widget _buildWaitingStep(ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120, height: 120,
              decoration: BoxDecoration(
                color: AppTheme.sanggamGold.withValues(alpha:0.1),
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.videocam, size: 48, color: AppTheme.sanggamGold),
            ),
            const SizedBox(height: 24),
            Text('대기실', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Text(
              '${_selectedDoctor?.name ?? ''} 전문의와의 진료를 준비하고 있습니다.',
              style: theme.textTheme.bodyMedium,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            const CircularProgressIndicator(),
            const SizedBox(height: 16),
            Text(
              '잠시만 기다려주세요. 곧 연결됩니다.',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () {
                // 화상 통화 화면으로 이동 (sessionId = 예약 기반 roomId)
                final roomId = 'room-${DateTime.now().millisecondsSinceEpoch}';
                context.push('/medical/video-call/$roomId');
              },
              icon: const Icon(Icons.videocam),
              label: const Text('진료실 입장'),
              style: FilledButton.styleFrom(
                backgroundColor: AppTheme.sanggamGold,
                minimumSize: const Size.fromHeight(48),
              ),
            ),
            const SizedBox(height: 12),
            OutlinedButton(
              onPressed: () => context.pop(),
              child: const Text('대기실 나가기'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _infoRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          SizedBox(width: 80, child: Text(label, style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600))),
          Expanded(child: Text(value, style: theme.textTheme.bodyMedium)),
        ],
      ),
    );
  }
}

class _DoctorInfo {
  final String id, name, specialty, hospital, experience;
  final double rating;
  final bool available;
  const _DoctorInfo({this.id = '', required this.name, required this.specialty, required this.hospital, required this.rating, required this.experience, required this.available});
}
