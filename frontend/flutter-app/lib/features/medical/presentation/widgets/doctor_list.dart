import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';
import 'package:manpasik/shared/widgets/scale_button.dart';

/// 추천 의사 횡스크롤 리스트 — Riverpod Provider 연결
class DoctorList extends ConsumerWidget {
  const DoctorList({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final doctorsAsync = ref.watch(recommendedDoctorsProvider);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.symmetric(horizontal: 16),
          child: Text(
            '추천 의료진',
            style: TextStyle(
              color: Colors.white,
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 150,
          child: doctorsAsync.when(
            data: (doctors) {
              if (doctors.isEmpty) {
                return const Center(
                  child: Text(
                    '추천 의료진 정보를 준비 중입니다',
                    style: TextStyle(color: Colors.white38, fontSize: 13),
                  ),
                );
              }
              return ListView.builder(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                itemCount: doctors.length,
                itemBuilder: (context, index) {
                  final doc = doctors[index];
                  return Padding(
                    padding: const EdgeInsets.only(right: 12),
                    child: AnimateFadeInUp(
                      delay: Duration(milliseconds: index * 100),
                      child: _DoctorCard(doctor: doc),
                    ),
                  );
                },
              );
            },
            loading: () => const Center(
              child: CircularProgressIndicator(color: AppTheme.sanggamGold),
            ),
            error: (_, __) => const Center(
              child: Text(
                '의료진 정보를 불러올 수 없습니다',
                style: TextStyle(color: Colors.white38, fontSize: 13),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class _DoctorCard extends StatelessWidget {
  final DoctorInfo doctor;
  const _DoctorCard({required this.doctor});

  @override
  Widget build(BuildContext context) {
    return ScaleButton(
      onPressed: () => context.push('/medical/telemedicine'),
      child: Container(
        width: 120,
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.white10),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // 아바타
            CircleAvatar(
              radius: 24,
              backgroundColor: AppTheme.sanggamGold.withOpacity(0.2),
              child: Text(
                doctor.name.isNotEmpty ? doctor.name[0] : '?',
                style: const TextStyle(
                  color: AppTheme.sanggamGold,
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            const SizedBox(height: 10),
            // 이름
            Text(
              doctor.name,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
                fontSize: 13,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            // 전문분야
            Text(
              doctor.specialty,
              style: TextStyle(
                color: Colors.white.withOpacity(0.6),
                fontSize: 10,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 6),
            // 별점 + 가용 상태
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.star, size: 10, color: Colors.amber),
                const SizedBox(width: 2),
                Text(
                  '${doctor.rating}',
                  style: const TextStyle(fontSize: 10, color: Colors.white70),
                ),
              ],
            ),
            const SizedBox(height: 2),
            if (doctor.isAvailable)
              Text(
                doctor.nextSlot ?? '예약 가능',
                style: const TextStyle(fontSize: 9, color: Color(0xFF00E676), fontWeight: FontWeight.w600),
              )
            else
              const Text(
                '예약 마감',
                style: TextStyle(fontSize: 9, color: Colors.white38),
              ),
          ],
        ),
      ),
    );
  }
}
