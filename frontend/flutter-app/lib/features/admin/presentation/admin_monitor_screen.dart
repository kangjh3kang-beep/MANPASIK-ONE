import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminMonitorScreen — Sanggam Orbit 시스템 모니터링
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] +sanggam_container.dart
// [Rule 4] AppBar+TabBar → body 내 커스텀 헤더 + 독립 TabBar
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData 파라미터 4x 제거
// [Rule 4] theme.textTheme.* ~4x → 직접 TextStyle
// [Rule 4] Card ~8x → SanggamContainer
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.purple → SanggamTheme.jagaeMagenta
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 4→8, 12→16
// ───────────────────────────────────────────────────

/// 시스템 모니터링 화면
class AdminMonitorScreen
    extends ConsumerWidget {
  const AdminMonitorScreen({super.key, this.tab});

  final String? tab;

  @override
  Widget build(
      BuildContext context, WidgetRef ref) {
    return DefaultTabController(
      length: 3,
      initialIndex:
          tab == 'emergency' ? 2 : 0,
      child: Scaffold(
        backgroundColor:
            SanggamTheme.background,
        body: SafeArea(
          child: Column(
            children: [
              // 헤더
              Padding(
                padding:
                    const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 8),
                child: Row(
                  children: [
                    IconButton(
                      icon: const Icon(
                          Icons.arrow_back,
                          color: Colors.white),
                      tooltip: '뒤로 가기',
                      onPressed: () =>
                          context.pop(),
                    ),
                    const SizedBox(width: 8),
                    const Expanded(
                      child: Text(
                        '시스템 모니터링',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 20,
                          fontWeight:
                              FontWeight.bold,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              // TabBar
              const TabBar(
                indicatorColor:
                    SanggamTheme.primary,
                labelColor:
                    SanggamTheme.primary,
                unselectedLabelColor:
                    SanggamTheme.onSurfaceDim,
                dividerColor:
                    SanggamTheme.surfaceVariant,
                tabs: [
                  Tab(text: '서비스'),
                  Tab(text: '리소스'),
                  Tab(text: '긴급'),
                ],
              ),
              // TabBarView
              const Expanded(
                child: TabBarView(
                  children: [
                    _ServiceStatusTab(),
                    _ResourceTab(),
                    _EmergencyTab(),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ServiceStatusTab
    extends ConsumerStatefulWidget {
  const _ServiceStatusTab();

  @override
  ConsumerState<_ServiceStatusTab> createState() =>
      _ServiceStatusTabState();
}

class _ServiceStatusTabState
    extends ConsumerState<_ServiceStatusTab> {
  List<_ServiceInfo> _services = _defaultServices;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _loadMetrics();
    _timer = Timer.periodic(
        const Duration(seconds: 30), (_) => _loadMetrics());
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<void> _loadMetrics() async {
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.getSystemMetrics();
      final rawServices = (data['services'] as List?) ?? [];
      if (mounted && rawServices.isNotEmpty) {
        setState(() {
          _services = rawServices
              .map((e) {
                final m = e as Map<String, dynamic>;
                return _ServiceInfo(
                  (m['name'] as String?) ?? '',
                  (m['healthy'] as bool?) ?? true,
                  (m['latency_ms'] as num?)?.toInt() ?? 0,
                );
              })
              .toList();
        });
      }
    } on DioException {
      // API 실패 시 기본값 유지
    }
  }

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        ..._services.map((s) =>
            SanggamContainer(
              borderRadius: 16,
              padding: EdgeInsets.zero,
              margin: const EdgeInsets.only(
                  bottom: 8),
              child: ListTile(
                leading: Icon(Icons.circle,
                    size: 12,
                    color: s.healthy
                        ? SanggamTheme
                            .jagaeCyan
                        : SanggamTheme.error),
                title: Text(s.name,
                    style: const TextStyle(
                        color: Colors.white)),
                subtitle: Text(
                    '응답시간: ${s.latency}ms',
                    style: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim,
                      fontSize: 12,
                    )),
                trailing: Text(
                    s.healthy ? '정상' : '장애',
                    style: TextStyle(
                      color: s.healthy
                          ? SanggamTheme
                              .jagaeCyan
                          : SanggamTheme.error,
                      fontWeight:
                          FontWeight.bold,
                    )),
              ),
            )),
      ],
    );
  }

  static final _defaultServices = [
    _ServiceInfo(
        '게이트웨이 (Gateway)', true, 12),
    _ServiceInfo('인증 (Auth)', true, 8),
    _ServiceInfo(
        '측정 (Measurement)', true, 15),
    _ServiceInfo('사용자 (User)', true, 10),
    _ServiceInfo('디바이스 (Device)', true, 11),
    _ServiceInfo(
        'AI 추론 (AI Inference)', true, 45),
    _ServiceInfo(
        '알림 (Notification)', true, 9),
    _ServiceInfo('가족 (Family)', true, 13),
    _ServiceInfo(
        '커뮤니티 (Community)', true, 14),
    _ServiceInfo(
        '원격진료 (Telemedicine)', true, 20),
  ];
}

class _ResourceTab extends ConsumerStatefulWidget {
  const _ResourceTab();

  @override
  ConsumerState<_ResourceTab> createState() =>
      _ResourceTabState();
}

class _ResourceTabState
    extends ConsumerState<_ResourceTab> {
  double _cpuUsage = 0.35;
  double _memUsage = 0.62;
  double _diskUsage = 0.28;
  double _netUsage = 0.15;
  String _netDisplay = '1.2 Gbps';
  String _totalRequests = '12,450';
  String _avgLatency = '23ms';
  String _errorRate = '0.02%';
  String _activeConns = '342';
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _loadResources();
    _timer = Timer.periodic(
        const Duration(seconds: 30), (_) => _loadResources());
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<void> _loadResources() async {
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.getSystemMetrics();
      final res = (data['resources'] as Map<String, dynamic>?) ?? {};
      final stats = (data['request_stats'] as Map<String, dynamic>?) ?? {};
      if (mounted && res.isNotEmpty) {
        setState(() {
          _cpuUsage = (res['cpu_usage'] as num?)?.toDouble() ?? _cpuUsage;
          _memUsage = (res['memory_usage'] as num?)?.toDouble() ?? _memUsage;
          _diskUsage = (res['disk_usage'] as num?)?.toDouble() ?? _diskUsage;
          _netUsage = (res['network_usage'] as num?)?.toDouble() ?? _netUsage;
          _netDisplay = (res['network_display'] as String?) ?? _netDisplay;
          _totalRequests = (stats['total_requests'] as String?) ?? _totalRequests;
          _avgLatency = (stats['avg_latency'] as String?) ?? _avgLatency;
          _errorRate = (stats['error_rate'] as String?) ?? _errorRate;
          _activeConns = (stats['active_connections'] as String?) ?? _activeConns;
        });
      }
    } on DioException {
      // API 실패 시 기본값 유지
    }
  }

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        _buildMetric('CPU 사용률', _cpuUsage,
            '${(_cpuUsage * 100).toStringAsFixed(0)}%',
            SanggamTheme.jagaeCyan),
        _buildMetric('메모리 사용률', _memUsage,
            '${(_memUsage * 100).toStringAsFixed(0)}%',
            SanggamTheme.primary),
        _buildMetric('디스크 사용률', _diskUsage,
            '${(_diskUsage * 100).toStringAsFixed(0)}%',
            SanggamTheme.jagaeCyan),
        _buildMetric('네트워크 I/O', _netUsage,
            _netDisplay,
            SanggamTheme.jagaeMagenta),
        const SizedBox(height: 16),
        const Text('요청 통계 (최근 1시간)',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            )),
        const SizedBox(height: 8),
        SanggamContainer(
          borderRadius: 16,
          padding: const EdgeInsets.all(16),
          child: Column(
            children: [
              _statRow('총 요청 수', _totalRequests),
              _statRow('평균 응답 시간', _avgLatency),
              _statRow('에러율', _errorRate),
              _statRow('활성 연결', _activeConns),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildMetric(String label,
      double value, String display, Color color) {
    return SanggamContainer(
      borderRadius: 16,
      padding: const EdgeInsets.all(16),
      margin:
          const EdgeInsets.only(bottom: 8),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment:
                MainAxisAlignment.spaceBetween,
            children: [
              Text(label,
                  style: const TextStyle(
                      color: Colors.white)),
              Text(display,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.bold,
                  )),
            ],
          ),
          const SizedBox(height: 8),
          LinearProgressIndicator(
            value: value,
            color: color,
            backgroundColor:
                color.withValues(alpha: 0.1),
          ),
        ],
      ),
    );
  }

  Widget _statRow(
      String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 4),
      child: Row(
        mainAxisAlignment:
            MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim,
                fontSize: 14,
              )),
          Text(value,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontWeight: FontWeight.bold,
              )),
        ],
      ),
    );
  }
}

class _EmergencyTab extends StatelessWidget {
  const _EmergencyTab();

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        SanggamContainer(
          borderRadius: 16,
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Icon(Icons.check_circle,
                  color: SanggamTheme.jagaeCyan,
                  size: 32),
              const SizedBox(width: 16),
              const Expanded(
                  child: Text(
                      '현재 활성화된 긴급 알림이 없습니다.',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight:
                            FontWeight.w600,
                      ))),
            ],
          ),
        ),
        const SizedBox(height: 16),
        const Text('최근 긴급 이벤트',
            style: TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.bold,
            )),
        const SizedBox(height: 8),
        SanggamContainer(
          borderRadius: 16,
          padding: EdgeInsets.zero,
          margin:
              const EdgeInsets.only(bottom: 8),
          child: const ListTile(
            leading: Icon(Icons.warning,
                color: SanggamTheme.primary),
            title: Text('AI 서비스 지연',
                style: TextStyle(
                    color: Colors.white)),
            subtitle: Text(
                '2024-02-14 03:15 — 자동 복구됨',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
          ),
        ),
        SanggamContainer(
          borderRadius: 16,
          padding: EdgeInsets.zero,
          margin:
              const EdgeInsets.only(bottom: 8),
          child: const ListTile(
            leading: Icon(Icons.error,
                color: SanggamTheme.error),
            title: Text('DB 연결 끊김',
                style: TextStyle(
                    color: Colors.white)),
            subtitle: Text(
                '2024-02-10 22:00 — 수동 복구',
                style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12,
                )),
          ),
        ),
      ],
    );
  }
}

class _ServiceInfo {
  final String name;
  final bool healthy;
  final int latency;
  const _ServiceInfo(
      this.name, this.healthy, this.latency);
}
