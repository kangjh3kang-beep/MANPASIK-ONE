import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/medical/domain/medical_repository.dart';

/// 의료 서비스 화면
class MedicalScreen extends ConsumerWidget {
  const MedicalScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final reservationsAsync = ref.watch(reservationsProvider);
    final prescriptionsAsync = ref.watch(prescriptionsProvider);
    final reportsAsync = ref.watch(healthReportsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('의료 서비스'),
        centerTitle: true,
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(reservationsProvider);
          ref.invalidate(prescriptionsProvider);
          ref.invalidate(healthReportsProvider);
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildEmergencyBanner(theme),
              const SizedBox(height: 24),
              _buildServiceGrid(theme, context),
              const SizedBox(height: 24),

              Text('예약 현황', style: theme.textTheme.titleLarge),
              const SizedBox(height: 12),
              reservationsAsync.when(
                data: (reservations) => reservations.isEmpty
                    ? _buildEmptyCard(context, theme, '예약된 진료가 없습니다', '진료 예약하기')
                    : Column(
                        children: reservations.map((r) => _buildReservationCard(theme, r)).toList(),
                      ),
                loading: () => const Center(child: Padding(padding: EdgeInsets.all(16), child: CircularProgressIndicator())),
                error: (_, __) => _buildEmptyCard(context, theme, '예약 정보를 불러올 수 없습니다', null),
              ),
              const SizedBox(height: 24),

              Text('최근 처방전', style: theme.textTheme.titleLarge),
              const SizedBox(height: 12),
              prescriptionsAsync.when(
                data: (prescriptions) => prescriptions.isEmpty
                    ? Card(
                        child: Padding(
                          padding: const EdgeInsets.all(24),
                          child: Center(
                            child: Text('처방전 내역이 없습니다',
                                style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.outline)),
                          ),
                        ),
                      )
                    : Column(
                        children: prescriptions.map((p) => _buildPrescriptionCard(context, theme, p)).toList(),
                      ),
                loading: () => const Center(child: Padding(padding: EdgeInsets.all(16), child: CircularProgressIndicator())),
                error: (_, __) => Card(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Center(child: Text('처방전 정보를 불러올 수 없습니다', style: theme.textTheme.bodyMedium)),
                  ),
                ),
              ),
              const SizedBox(height: 24),

              Text('건강 리포트', style: theme.textTheme.titleLarge),
              const SizedBox(height: 12),
              reportsAsync.when(
                data: (reports) => reports.isEmpty
                    ? Card(
                        child: ListTile(
                          leading: Icon(Icons.analytics, color: theme.colorScheme.primary),
                          title: const Text('건강 분석 리포트'),
                          subtitle: const Text('충분한 측정 데이터가 쌓이면 AI가 분석합니다'),
                          trailing: const Icon(Icons.chevron_right),
                          onTap: () async {
                            try {
                              await ref.read(medicalRepositoryProvider).generateHealthReport();
                              ref.invalidate(healthReportsProvider);
                            } catch (_) {}
                          },
                        ),
                      )
                    : Column(
                        children: reports.map((r) => _buildReportCard(theme, r)).toList(),
                      ),
                loading: () => const Center(child: Padding(padding: EdgeInsets.all(16), child: CircularProgressIndicator())),
                error: (_, __) => const SizedBox.shrink(),
              ),

              const SizedBox(height: 32),
              _buildDisclaimer(theme),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildEmergencyBanner(ThemeData theme) {
    return Card(
      color: theme.colorScheme.errorContainer,
      child: ListTile(
        leading: Icon(Icons.emergency, color: theme.colorScheme.error),
        title: Text(
          '응급 상황 시 119에 전화하세요',
          style: TextStyle(color: theme.colorScheme.onErrorContainer, fontWeight: FontWeight.bold),
        ),
        subtitle: Text(
          '본 서비스는 응급 의료 서비스를 대체하지 않습니다',
          style: TextStyle(color: theme.colorScheme.onErrorContainer),
        ),
      ),
    );
  }

  Widget _buildServiceGrid(ThemeData theme, BuildContext context) {
    final services = [
      _ServiceItem(Icons.videocam, '비대면 진료', '화상 진료 예약', theme.colorScheme.primary, '/medical/telemedicine'),
      _ServiceItem(Icons.medication, '처방전 조회', '처방 내역 확인', Colors.green, null),
      _ServiceItem(Icons.description, '건강 리포트', 'AI 분석 결과', Colors.orange, null),
      _ServiceItem(Icons.local_hospital, '병원 찾기', '주변 병원 검색', Colors.blue, '/medical/facility-search'),
    ];
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.3,
      children: services.map((s) {
        return Card(
          child: InkWell(
            onTap: s.route != null ? () => context.push(s.route!) : null,
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(s.icon, size: 32, color: s.color),
                  const SizedBox(height: 8),
                  Text(s.title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                  Text(s.subtitle, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
                ],
              ),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildEmptyCard(BuildContext context, ThemeData theme, String message, String? buttonLabel) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: Column(
            children: [
              Icon(Icons.calendar_today, size: 40, color: theme.colorScheme.outline),
              const SizedBox(height: 12),
              Text(message, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.outline)),
              if (buttonLabel != null) ...[
                const SizedBox(height: 12),
                OutlinedButton(
                  onPressed: () => context.push('/medical/telemedicine'),
                  child: Text(buttonLabel),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildReservationCard(ThemeData theme, TelemedicineReservation r) {
    final statusColor = switch (r.status) {
      ReservationStatus.confirmed => Colors.green,
      ReservationStatus.inProgress => Colors.orange,
      ReservationStatus.completed => Colors.grey,
      ReservationStatus.cancelled => Colors.red,
      _ => theme.colorScheme.primary,
    };
    final statusText = switch (r.status) {
      ReservationStatus.pending => '대기 중',
      ReservationStatus.confirmed => '확정',
      ReservationStatus.inProgress => '진행 중',
      ReservationStatus.completed => '완료',
      ReservationStatus.cancelled => '취소됨',
    };
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Icon(Icons.videocam, color: statusColor),
        title: Text(r.doctorName.isNotEmpty ? r.doctorName : '비대면 진료'),
        subtitle: Text('${r.scheduledAt.month}/${r.scheduledAt.day} ${r.scheduledAt.hour}:${r.scheduledAt.minute.toString().padLeft(2, '0')}'),
        trailing: Chip(
          label: Text(statusText, style: TextStyle(color: statusColor, fontSize: 12)),
          backgroundColor: statusColor.withOpacity(0.1),
          side: BorderSide.none,
        ),
      ),
    );
  }

  Widget _buildPrescriptionCard(BuildContext context, ThemeData theme, Prescription p) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Icon(Icons.medication, color: Colors.green.shade600),
        title: Text(p.doctorName.isNotEmpty ? '${p.doctorName} 처방' : '처방전'),
        subtitle: Text('${p.items.length}개 의약품 | ${p.issuedAt.month}/${p.issuedAt.day} 발행'),
        trailing: const Icon(Icons.chevron_right),
        onTap: () => context.push('/medical/prescription/${p.id}'),
      ),
    );
  }

  Widget _buildReportCard(ThemeData theme, HealthReport r) {
    final statusIcon = switch (r.overallStatus) {
      'excellent' => Icons.sentiment_very_satisfied,
      'good' => Icons.sentiment_satisfied,
      'caution' => Icons.sentiment_neutral,
      'alert' => Icons.sentiment_very_dissatisfied,
      _ => Icons.analytics,
    };
    final statusColor = switch (r.overallStatus) {
      'excellent' => Colors.green,
      'good' => Colors.blue,
      'caution' => Colors.orange,
      'alert' => Colors.red,
      _ => theme.colorScheme.primary,
    };
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Icon(statusIcon, color: statusColor),
        title: Text(r.periodDescription.isNotEmpty ? r.periodDescription : '건강 리포트'),
        subtitle: Text('${r.generatedAt.month}/${r.generatedAt.day} | ${r.analyses.length}개 바이오마커 분석'),
        trailing: const Icon(Icons.chevron_right),
        onTap: () {},
      ),
    );
  }

  Widget _buildDisclaimer(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest.withOpacity(0.3),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, size: 16, color: theme.colorScheme.outline),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '본 서비스의 측정 결과 및 AI 분석은 의료 진단을 대체하지 않습니다. '
              '정확한 진단은 반드시 의료 전문가와 상담하시기 바랍니다.',
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline),
            ),
          ),
        ],
      ),
    );
  }
}

class _ServiceItem {
  final IconData icon;
  final String title;
  final String subtitle;
  final Color color;
  final String? route;
  const _ServiceItem(this.icon, this.title, this.subtitle, this.color, this.route);
}
