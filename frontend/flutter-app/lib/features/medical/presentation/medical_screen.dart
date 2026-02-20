import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';
import 'package:manpasik/features/medical/presentation/widgets/doctor_list.dart';
import 'package:manpasik/features/medical/presentation/widgets/hero_telemedicine_card.dart';
import 'package:manpasik/shared/widgets/animate_fade_in_up.dart';

/// 의료 서비스 메인 화면 (Digital Hospital Lobby)
class MedicalScreen extends ConsumerStatefulWidget {
  const MedicalScreen({super.key});

  @override
  ConsumerState<MedicalScreen> createState() => _MedicalScreenState();
}

class _MedicalScreenState extends ConsumerState<MedicalScreen> {
  @override
  Widget build(BuildContext context) {
    final reservationsAsync = ref.watch(reservationsProvider);
    final prescriptionsAsync = ref.watch(prescriptionsProvider);
    final reportsAsync = ref.watch(healthReportsProvider);

    // Force Dark Theme for Premium Medical Feel
    return Theme(
      data: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: const Color(0xFF050B14),
        colorScheme: ColorScheme.dark(
          primary: const Color(0xFF00E5FF),
          surface: const Color(0xFF1A1F35),
        ),
      ),
      child: Scaffold(
        backgroundColor: Colors.transparent,
        body: RefreshIndicator(
          color: AppTheme.sanggamGold,
          backgroundColor: const Color(0xFF1A1F35),
          onRefresh: () async {
            ref.invalidate(reservationsProvider);
            ref.invalidate(prescriptionsProvider);
            ref.invalidate(healthReportsProvider);
            ref.invalidate(recommendedDoctorsProvider);
          },
          child: CustomScrollView(
            physics: const BouncingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
            slivers: [
              SliverAppBar(
                expandedHeight: 80.0,
                floating: true,
                pinned: false,
                backgroundColor: Colors.transparent,
                centerTitle: false,
                title: const Text(
                  'Digital Hospital',
                  style: TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                    fontSize: 24,
                    letterSpacing: -0.5,
                  ),
                ),
                actions: [
                  IconButton(
                    icon: const Icon(Icons.notifications_outlined, color: Colors.white),
                    onPressed: () {},
                  ),
                ],
              ),

              // 1. Hero Telemedicine Section
              const SliverToBoxAdapter(
                child: Padding(
                  padding: EdgeInsets.only(bottom: 24),
                  child: AnimateFadeInUp(child: HeroTelemedicineCard()),
                ),
              ),

              // 2. Quick Services (2x2 Grid)
              SliverPadding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                sliver: SliverGrid.count(
                  crossAxisCount: 2,
                  mainAxisSpacing: 12,
                  crossAxisSpacing: 12,
                  childAspectRatio: 1.4,
                  children: [
                    _buildQuickServiceItem(Icons.videocam_outlined, '화상 진료', '/medical/telemedicine'), // Added Video Call
                    _buildQuickServiceItem(Icons.medication_outlined, '처방전 관리', '/medical/prescription/1'),
                    _buildQuickServiceItem(Icons.local_pharmacy_outlined, '약국 찾기', '/medical/pharmacy'),
                    _buildQuickServiceItem(Icons.local_hospital_outlined, '주변 병원 찾기', '/medical/facility-search'),
                  ],
                ),
              ),

              const SliverToBoxAdapter(child: SizedBox(height: 32)),

              // 3. Recommended Doctors
              const SliverToBoxAdapter(
                child: DoctorList(),
              ),

              const SliverToBoxAdapter(child: SizedBox(height: 32)),

              // 4. 나의 진료 현황 (예약)
              SliverToBoxAdapter(
                child: _buildSectionWithData(
                  title: '나의 진료 현황',
                  asyncValue: reservationsAsync,
                  emptyMessage: '예정된 진료 내역이 없습니다.',
                  emptyAction: '진료 예약하기',
                  onEmptyAction: () => context.push('/medical/telemedicine'),
                  builder: (reservations) => Column(
                    children: reservations.map((r) => _buildReservationTile(r)).toList(),
                  ),
                ),
              ),

              const SliverToBoxAdapter(child: SizedBox(height: 24)),

              // 5. 최근 처방전
              SliverToBoxAdapter(
                child: _buildSectionWithData(
                  title: '최근 처방전',
                  asyncValue: prescriptionsAsync,
                  emptyMessage: '처방전 내역이 없습니다.',
                  builder: (prescriptions) => Column(
                    children: prescriptions.map((p) => _buildPrescriptionTile(p)).toList(),
                  ),
                ),
              ),

              const SliverToBoxAdapter(child: SizedBox(height: 24)),

              // 6. 건강 리포트
              SliverToBoxAdapter(
                child: _buildSectionWithData(
                  title: '건강 리포트',
                  asyncValue: reportsAsync,
                  emptyMessage: '충분한 측정 데이터가 쌓이면 AI가 분석합니다.',
                  builder: (reports) => Column(
                    children: reports.map((r) => _buildReportTile(r)).toList(),
                  ),
                ),
              ),

              // 7. 응급 안내 + 면책
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(16, 24, 16, 0),
                  child: _buildEmergencyBanner(),
                ),
              ),

              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(16, 12, 16, 0),
                  child: _buildDisclaimer(),
                ),
              ),

              const SliverPadding(padding: EdgeInsets.only(bottom: 100)),
            ],
          ),
        ),
      ),
    );
  }

  // ── Quick Service Item ──
  Widget _buildQuickServiceItem(IconData icon, String label, String? route) {
    return AnimateFadeInUp(
      delay: const Duration(milliseconds: 200),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: route != null ? () => context.push(route) : () {},
          borderRadius: BorderRadius.circular(16),
          child: Container(
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.05),
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: Colors.white10),
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(icon, color: AppTheme.sanggamGold, size: 28),
                const SizedBox(height: 8),
                Text(
                  label,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // ── 범용 섹션 빌더 ──
  Widget _buildSectionWithData<T>({
    required String title,
    required AsyncValue<List<T>> asyncValue,
    required String emptyMessage,
    String? emptyAction,
    VoidCallback? onEmptyAction,
    required Widget Function(List<T> data) builder,
  }) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 12),
          asyncValue.when(
            data: (list) {
              if (list.isEmpty) return _buildEmptyTile(emptyMessage, emptyAction, onEmptyAction);
              return builder(list);
            },
            loading: () => Container(
              width: double.infinity,
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: const Color(0xFF1A1F35),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: Colors.white10),
              ),
              child: const Center(
                child: SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: AppTheme.sanggamGold)),
              ),
            ),
            error: (_, __) => _buildEmptyTile('정보를 불러올 수 없습니다', null, null),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyTile(String message, String? actionLabel, VoidCallback? onAction) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1F35),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.white10),
      ),
      child: Column(
        children: [
          Text(message, style: const TextStyle(color: Colors.white38, fontSize: 13)),
          if (actionLabel != null && onAction != null) ...[
            const SizedBox(height: 12),
            OutlinedButton(
              onPressed: onAction,
              style: OutlinedButton.styleFrom(
                side: BorderSide(color: AppTheme.sanggamGold.withOpacity(0.5)),
                foregroundColor: AppTheme.sanggamGold,
              ),
              child: Text(actionLabel, style: const TextStyle(fontSize: 12)),
            ),
          ],
        ],
      ),
    );
  }

  // ── 예약 타일 ──
  Widget _buildReservationTile(TelemedicineReservation r) {
    final statusColor = switch (r.status) {
      ReservationStatus.confirmed => const Color(0xFF00E676),
      ReservationStatus.inProgress => Colors.orange,
      ReservationStatus.completed => Colors.grey,
      ReservationStatus.cancelled => Colors.red,
      _ => const Color(0xFF00E5FF),
    };
    final statusText = switch (r.status) {
      ReservationStatus.pending => '대기 중',
      ReservationStatus.confirmed => '확정',
      ReservationStatus.inProgress => '진행 중',
      ReservationStatus.completed => '완료',
      ReservationStatus.cancelled => '취소됨',
    };

    return Container(
      width: double.infinity,
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1F35),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: Colors.white10),
      ),
      child: Row(
        children: [
          Icon(Icons.videocam, color: statusColor, size: 22),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  r.doctorName.isNotEmpty ? r.doctorName : '비대면 진료',
                  style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w600),
                ),
                Text(
                  '${r.scheduledAt.month}/${r.scheduledAt.day} ${r.scheduledAt.hour}:${r.scheduledAt.minute.toString().padLeft(2, '0')}',
                  style: const TextStyle(color: Colors.white38, fontSize: 12),
                ),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              color: statusColor.withOpacity(0.15),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Text(statusText, style: TextStyle(color: statusColor, fontSize: 11, fontWeight: FontWeight.w600)),
          ),
        ],
      ),
    );
  }

  // ── 처방전 타일 ──
  Widget _buildPrescriptionTile(Prescription p) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: () => context.push('/medical/prescription/${p.id}'),
        borderRadius: BorderRadius.circular(14),
        child: Container(
          width: double.infinity,
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(
            color: const Color(0xFF1A1F35),
            borderRadius: BorderRadius.circular(14),
            border: Border.all(color: Colors.white10),
          ),
          child: Row(
            children: [
              const Icon(Icons.medication, color: Color(0xFF66BB6A), size: 22),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      p.doctorName.isNotEmpty ? '${p.doctorName} 처방' : '처방전',
                      style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w600),
                    ),
                    Text(
                      '${p.items.length}개 의약품 | ${p.issuedAt.month}/${p.issuedAt.day} 발행',
                      style: const TextStyle(color: Colors.white38, fontSize: 12),
                    ),
                  ],
                ),
              ),
              const Icon(Icons.chevron_right, color: Colors.white24, size: 20),
            ],
          ),
        ),
      ),
    );
  }

  // ── 리포트 타일 ──
  Widget _buildReportTile(HealthReport r) {
    final statusColor = switch (r.overallStatus) {
      'excellent' => const Color(0xFF00E676),
      'good' => const Color(0xFF29B6F6),
      'caution' => Colors.orange,
      'alert' => Colors.red,
      _ => const Color(0xFF00E5FF),
    };
    final statusIcon = switch (r.overallStatus) {
      'excellent' => Icons.sentiment_very_satisfied,
      'good' => Icons.sentiment_satisfied,
      'caution' => Icons.sentiment_neutral,
      'alert' => Icons.sentiment_very_dissatisfied,
      _ => Icons.analytics,
    };

    return Container(
      width: double.infinity,
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1F35),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: Colors.white10),
      ),
      child: Row(
        children: [
          Icon(statusIcon, color: statusColor, size: 22),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  r.periodDescription.isNotEmpty ? r.periodDescription : '건강 리포트',
                  style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w600),
                ),
                Text(
                  '${r.generatedAt.month}/${r.generatedAt.day} | ${r.analyses.length}개 바이오마커',
                  style: const TextStyle(color: Colors.white38, fontSize: 12),
                ),
              ],
            ),
          ),
          const Icon(Icons.chevron_right, color: Colors.white24, size: 20),
        ],
      ),
    );
  }

  // ── 응급 배너 ──
  Widget _buildEmergencyBanner() {
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: Colors.red.withOpacity(0.1),
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: Colors.red.withOpacity(0.2)),
      ),
      child: Row(
        children: [
          const Icon(Icons.emergency, color: Colors.redAccent, size: 22),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '응급 상황 시 119에 전화하세요',
                  style: TextStyle(color: Colors.redAccent, fontWeight: FontWeight.bold, fontSize: 13),
                ),
                Text(
                  '본 서비스는 응급 의료 서비스를 대체하지 않습니다',
                  style: TextStyle(color: Colors.redAccent.withOpacity(0.6), fontSize: 11),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  // ── 면책 조항 ──
  Widget _buildDisclaimer() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.03),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.info_outline, size: 14, color: Colors.white24),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '본 서비스의 측정 결과 및 AI 분석은 의료 진단을 대체하지 않습니다. '
              '정확한 진단은 반드시 의료 전문가와 상담하시기 바랍니다.',
              style: TextStyle(fontSize: 11, color: Colors.white.withOpacity(0.3)),
            ),
          ),
        ],
      ),
    );
  }
}
