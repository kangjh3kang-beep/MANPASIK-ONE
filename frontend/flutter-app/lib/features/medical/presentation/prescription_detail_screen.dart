import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 처방전 상세 화면
class PrescriptionDetailScreen extends ConsumerWidget {
  const PrescriptionDetailScreen({super.key, required this.prescriptionId});

  final String prescriptionId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('처방전 상세'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 처방전 헤더
            Card(
              color: AppTheme.sanggamGold.withValues(alpha:0.1),
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.receipt_long, color: AppTheme.sanggamGold),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text('처방전 #${prescriptionId.length > 8 ? prescriptionId.substring(0, 8) : prescriptionId}',
                              style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                        ),
                        Chip(
                          label: const Text('유효', style: TextStyle(fontSize: 11, color: Colors.white)),
                          backgroundColor: Colors.green,
                          side: BorderSide.none,
                          visualDensity: VisualDensity.compact,
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Text('발급일: 2026-02-15', style: theme.textTheme.bodySmall),
                    Text('유효기간: 2026-02-22까지', style: theme.textTheme.bodySmall),
                    Text('담당의: 김건강 전문의 (내과)', style: theme.textTheme.bodySmall),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 처방 약품 리스트
            Text('처방 약품', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            _buildMedicineCard(theme, '메트포르민 500mg', '1일 2회, 식후 30분', '30일분', Icons.medication),
            _buildMedicineCard(theme, '아토르바스타틴 10mg', '1일 1회, 취침 전', '30일분', Icons.medication_liquid),
            _buildMedicineCard(theme, '오메가3 1000mg', '1일 1회, 식사와 함께', '30일분', Icons.local_pharmacy),
            const SizedBox(height: 16),

            // 주의사항
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.warning_amber, size: 20, color: Colors.orange),
                        const SizedBox(width: 8),
                        Text('복약 주의사항', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
                      ],
                    ),
                    const SizedBox(height: 8),
                    _buildWarningItem(theme, '메트포르민 복용 중 과음을 삼가세요.'),
                    _buildWarningItem(theme, '근육통이 지속되면 아토르바스타틴 복용을 중단하고 의사에게 알리세요.'),
                    _buildWarningItem(theme, '처방 약품을 임의로 중단하지 마세요.'),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            // 약국 전송 버튼
            FilledButton.icon(
              onPressed: () => _showPharmacySearchSheet(context, theme, ref),
              icon: const Icon(Icons.local_pharmacy),
              label: const Text('약국으로 처방전 전송'),
              style: FilledButton.styleFrom(
                minimumSize: const Size.fromHeight(48),
                backgroundColor: AppTheme.sanggamGold,
              ),
            ),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () => _showReminderDialog(context),
              icon: const Icon(Icons.alarm),
              label: const Text('복약 리마인더 설정'),
              style: OutlinedButton.styleFrom(minimumSize: const Size.fromHeight(48)),
            ),
            const SizedBox(height: 8),
            OutlinedButton.icon(
              onPressed: () {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('처방전 PDF를 다운로드 중입니다...')),
                );
                // PDF 생성: printing + pdf 패키지 설치 후 실제 구현
                Future.delayed(const Duration(seconds: 1), () {
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('처방전 PDF가 다운로드되었습니다.')),
                    );
                  }
                });
              },
              icon: const Icon(Icons.picture_as_pdf),
              label: const Text('처방전 PDF 다운로드'),
              style: OutlinedButton.styleFrom(minimumSize: const Size.fromHeight(48)),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMedicineCard(ThemeData theme, String name, String dosage, String duration, IconData icon) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: CircleAvatar(
          backgroundColor: Colors.blue.withValues(alpha:0.1),
          child: Icon(icon, color: Colors.blue, size: 20),
        ),
        title: Text(name, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(dosage, style: theme.textTheme.bodySmall),
            Text(duration, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ),
      ),
    );
  }

  void _showPharmacySearchSheet(BuildContext context, ThemeData theme, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => _PharmacySearchSheet(
        prescriptionId: prescriptionId,
        theme: theme,
        parentContext: context,
      ),
    );
  }

  void _showReminderDialog(BuildContext context) {
    TimeOfDay morningTime = const TimeOfDay(hour: 8, minute: 0);
    TimeOfDay eveningTime = const TimeOfDay(hour: 20, minute: 0);

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('복약 리마인더 설정'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.wb_sunny_outlined),
              title: const Text('오전 알림'),
              subtitle: Text('${morningTime.hour}:${morningTime.minute.toString().padLeft(2, '0')}'),
              trailing: const Icon(Icons.chevron_right),
            ),
            ListTile(
              contentPadding: EdgeInsets.zero,
              leading: const Icon(Icons.nights_stay_outlined),
              title: const Text('오후 알림'),
              subtitle: Text('${eveningTime.hour}:${eveningTime.minute.toString().padLeft(2, '0')}'),
              trailing: const Icon(Icons.chevron_right),
            ),
            const SizedBox(height: 8),
            Text(
              '알림은 처방 기간(30일) 동안 매일 설정한 시간에 전송됩니다.',
              style: Theme.of(ctx).textTheme.bodySmall?.copyWith(color: Theme.of(ctx).colorScheme.outline),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('복약 리마인더가 설정되었습니다.')),
              );
            },
            child: const Text('설정'),
          ),
        ],
      ),
    );
  }

  Widget _buildWarningItem(ThemeData theme, String text) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('• ', style: TextStyle(color: Colors.orange)),
          Expanded(child: Text(text, style: theme.textTheme.bodySmall)),
        ],
      ),
    );
  }
}

/// 약국 검색 바텀시트 (REST API 연동)
class _PharmacySearchSheet extends ConsumerStatefulWidget {
  const _PharmacySearchSheet({
    required this.prescriptionId,
    required this.theme,
    required this.parentContext,
  });

  final String prescriptionId;
  final ThemeData theme;
  final BuildContext parentContext;

  @override
  ConsumerState<_PharmacySearchSheet> createState() => _PharmacySearchSheetState();
}

class _PharmacySearchSheetState extends ConsumerState<_PharmacySearchSheet> {
  List<Map<String, dynamic>> _pharmacies = [];
  bool _loading = true;
  String? _sendingTo;

  @override
  void initState() {
    super.initState();
    _loadPharmacies();
  }

  Future<void> _loadPharmacies() async {
    try {
      final client = ref.read(restClientProvider);
      final res = await client.searchFacilities(query: '약국', limit: 10);
      final list = res['facilities'] as List<dynamic>? ?? [];
      if (mounted) {
        setState(() {
          _pharmacies = list.cast<Map<String, dynamic>>();
          _loading = false;
        });
      }
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _sendToPharmacy(Map<String, dynamic> pharmacy) async {
    final pharmacyId = pharmacy['facility_id'] as String? ?? pharmacy['id'] as String? ?? '';
    final pharmacyName = pharmacy['name'] as String? ?? '약국';
    setState(() => _sendingTo = pharmacyId);
    try {
      final client = ref.read(restClientProvider);
      await client.selectPharmacy(
        widget.prescriptionId,
        pharmacyId: pharmacyId,
        pharmacyName: pharmacyName,
      );
      await client.sendToPharmacy(widget.prescriptionId);
      if (mounted) Navigator.pop(context);
      if (widget.parentContext.mounted) {
        ScaffoldMessenger.of(widget.parentContext).showSnackBar(
          SnackBar(content: Text('$pharmacyName에 처방전이 전송되었습니다.')),
        );
      }
    } catch (e) {
      if (mounted) {
        setState(() => _sendingTo = null);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('전송 실패: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = widget.theme;
    return DraggableScrollableSheet(
      initialChildSize: 0.5,
      minChildSize: 0.3,
      maxChildSize: 0.8,
      expand: false,
      builder: (_, scrollController) => Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Text('주변 약국', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                const Spacer(),
                TextButton(
                  onPressed: () {
                    Navigator.pop(context);
                    widget.parentContext.push('/medical/pharmacy');
                  },
                  child: const Text('더보기'),
                ),
              ],
            ),
          ),
          if (_loading)
            const Expanded(child: Center(child: CircularProgressIndicator()))
          else if (_pharmacies.isEmpty)
            Expanded(
              child: Center(
                child: Text('주변 약국을 찾을 수 없습니다.', style: theme.textTheme.bodyMedium),
              ),
            )
          else
            Expanded(
              child: ListView.builder(
                controller: scrollController,
                itemCount: _pharmacies.length,
                itemBuilder: (_, index) {
                  final p = _pharmacies[index];
                  final name = p['name'] as String? ?? '약국';
                  final address = p['address'] as String? ?? '';
                  final distance = p['distance_km'] as num?;
                  final pharmacyId = p['facility_id'] as String? ?? p['id'] as String? ?? '';
                  return ListTile(
                    leading: const Icon(Icons.local_pharmacy, color: Colors.green),
                    title: Text(name),
                    subtitle: Text('$address${distance != null ? '\n${distance.toStringAsFixed(1)}km' : ''}'),
                    isThreeLine: distance != null,
                    trailing: FilledButton(
                      onPressed: _sendingTo == pharmacyId ? null : () => _sendToPharmacy(p),
                      child: _sendingTo == pharmacyId
                          ? const SizedBox(width: 16, height: 16, child: CircularProgressIndicator(strokeWidth: 2))
                          : const Text('전송'),
                    ),
                  );
                },
              ),
            ),
        ],
      ),
    );
  }
}
