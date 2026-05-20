import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AlertDetailScreen — Sanggam Orbit 긴급 알림 상세
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] Theme.of(context) + ThemeData 파라미터 제거
// [Rule 4] theme.textTheme.* ~12x → 직접 TextStyle
// [Rule 4] theme.colorScheme.onSurfaceVariant 3x → SanggamTheme.onSurfaceDim
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card 4x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] padding:20→24, h:12→16, h:4→8, w:12→16
// ───────────────────────────────────────────────────

/// 긴급 알림 상세 화면
class AlertDetailScreen extends ConsumerStatefulWidget {
  const AlertDetailScreen({super.key, required this.alertId});

  final String alertId;

  @override
  ConsumerState<AlertDetailScreen> createState() =>
      _AlertDetailScreenState();
}

class _AlertDetailScreenState
    extends ConsumerState<AlertDetailScreen> {
  Map<String, dynamic>? _alert;
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadAlert();
  }

  Future<void> _loadAlert() async {
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.getAlertDetail(widget.alertId);
      if (mounted) {
        setState(() {
          _alert = data;
          _isLoading = false;
        });
      }
    } on DioException {
      if (mounted) {
        setState(() {
          _isLoading = false;
          _error = '알림 정보를 불러올 수 없습니다';
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final alertType = (_alert?['type'] as String?) ?? '이상 수치 감지';
    final memberName = (_alert?['member_name'] as String?) ?? '멤버';
    final createdAt = (_alert?['created_at'] as String?) ?? '';
    final analysis = (_alert?['analysis'] as String?) ??
        '수축기 혈압이 정상 범위를 크게 초과했습니다. 최근 7일간의 추세를 볼 때 점진적 상승이 관찰됩니다. '
        '스트레스, 식이 변화, 약물 복용 상태를 확인하시길 권장드립니다.';
    final metrics = (_alert?['metrics'] as List?) ?? [];

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '긴급 알림 상세',
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
            // 본문
            Expanded(
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : _error != null
                      ? Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(_error!,
                                  style: const TextStyle(
                                      color: SanggamTheme
                                          .onSurfaceDim)),
                              const SizedBox(height: 16),
                              FilledButton(
                                onPressed: () {
                                  setState(() {
                                    _isLoading = true;
                                    _error = null;
                                  });
                                  _loadAlert();
                                },
                                child: const Text('다시 시도'),
                              ),
                            ],
                          ),
                        )
                      : ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  // 알림 헤더
                  SanggamContainer(
                    borderRadius: 16,
                    borderColor: SanggamTheme.error
                        .withValues(alpha: 0.3),
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      children: [
                        const Icon(
                            Icons.warning_amber_rounded,
                            size: 48,
                            color: SanggamTheme.error),
                        const SizedBox(height: 16),
                        Text(
                          alertType,
                          style: const TextStyle(
                            color: SanggamTheme.error,
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '$memberName · $createdAt',
                          style: const TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 14,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 측정 수치
                  const Text(
                    '측정 수치',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      children: metrics.isEmpty
                          ? [
                              _buildMetricRow('수축기 혈압',
                                  '-- mmHg', '정상 범위: 90-120',
                                  SanggamTheme.onSurfaceDim),
                            ]
                          : _buildMetricRows(metrics),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // AI 분석
                  const Text(
                    'AI 분석',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment:
                          CrossAxisAlignment.start,
                      children: [
                        const Row(
                          children: [
                            Icon(Icons.auto_awesome,
                                color: SanggamTheme.primary),
                            SizedBox(width: 8),
                            Text(
                              'AI 건강 분석',
                              style: TextStyle(
                                color: Colors.white,
                                fontSize: 14,
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        Text(
                          analysis,
                          style: const TextStyle(
                            color: Colors.white,
                            fontSize: 14,
                          ),
                        ),
                        const SizedBox(height: 8),
                        const Text(
                          '※ 본 분석은 참고용이며, 정확한 진단은 의료 전문가와 상담하세요.',
                          style: TextStyle(
                            color: SanggamTheme.onSurfaceDim,
                            fontSize: 12,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),

                  // 대응 조치
                  const Text(
                    '권장 조치',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  SanggamContainer(
                    borderRadius: 16,
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        ListTile(
                          leading: const Icon(
                              Icons.local_hospital,
                              color: SanggamTheme.jagaeCyan),
                          title: const Text('가까운 병원 찾기'),
                          trailing: const Icon(
                              Icons.chevron_right),
                          onTap: () => context.push(
                              '/medical/facility-search'),
                        ),
                        const Divider(
                            height: 1,
                            color:
                                SanggamTheme.surfaceVariant),
                        ListTile(
                          leading: const Icon(
                              Icons.videocam,
                              color: SanggamTheme.jagaeCyan),
                          title: const Text('화상 진료 예약'),
                          trailing: const Icon(
                              Icons.chevron_right),
                          onTap: () => context.push(
                              '/medical/telemedicine'),
                        ),
                        const Divider(
                            height: 1,
                            color:
                                SanggamTheme.surfaceVariant),
                        ListTile(
                          leading: const Icon(Icons.phone,
                              color: SanggamTheme.error),
                          title: const Text('긴급 연락처 전화'),
                          trailing: const Icon(
                              Icons.chevron_right),
                          onTap: () => context
                              .push('/settings/emergency'),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  List<Widget> _buildMetricRows(List<dynamic> metrics) {
    final widgets = <Widget>[];
    for (int i = 0; i < metrics.length; i++) {
      final m = metrics[i] as Map<String, dynamic>;
      final severity = (m['severity'] as String?) ?? 'normal';
      final color = switch (severity) {
        'critical' => SanggamTheme.error,
        'warning' => SanggamTheme.primary,
        _ => SanggamTheme.jagaeCyan,
      };
      widgets.add(_buildMetricRow(
        (m['label'] as String?) ?? '',
        (m['value'] as String?) ?? '',
        (m['range'] as String?) ?? '',
        color,
      ));
      if (i < metrics.length - 1) {
        widgets.add(const Divider(
            color: SanggamTheme.surfaceVariant));
      }
    }
    return widgets;
  }

  Widget _buildMetricRow(
      String label, String value, String range, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(
            width: 4,
            height: 32,
            decoration: BoxDecoration(
              color: color,
              borderRadius: BorderRadius.circular(2),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
                Text(
                  value,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
          ),
          Text(
            range,
            style: const TextStyle(
              color: SanggamTheme.onSurfaceDim,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }
}
