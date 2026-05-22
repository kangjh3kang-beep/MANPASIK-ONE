import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// -----------------------------------------------
// AdminDeviceFleetScreen -- Device Fleet Management
//
// Lists all registered devices with status indicators,
// firmware version, battery level, and last seen timestamp.
// Connects to /admin/devices REST endpoint.
// -----------------------------------------------

/// 관리자 기기 관리 화면
class AdminDeviceFleetScreen extends ConsumerStatefulWidget {
  const AdminDeviceFleetScreen({super.key});

  @override
  ConsumerState<AdminDeviceFleetScreen> createState() =>
      _AdminDeviceFleetScreenState();
}

class _AdminDeviceFleetScreenState
    extends ConsumerState<AdminDeviceFleetScreen> {
  final _searchController = TextEditingController();
  List<Map<String, dynamic>> _devices = [];
  bool _isLoading = false;
  String? _error;
  String _statusFilter = 'all';
  Timer? _refreshTimer;

  @override
  void initState() {
    super.initState();
    _loadDevices();
    // Auto-refresh every 30 seconds
    _refreshTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) => _loadDevices(),
    );
  }

  @override
  void dispose() {
    _searchController.dispose();
    _refreshTimer?.cancel();
    super.dispose();
  }

  Future<void> _loadDevices() async {
    if (_isLoading) return;
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final client = ref.read(restClientProvider);
      // Try admin devices endpoint first, fall back to general devices
      Map<String, dynamic> resp;
      try {
        resp = await client.listDevices('');
      } on DioException catch (e) {
        setState(() {
          _error = e.message ?? '기기 목록을 불러오는 중 네트워크 오류가 발생했습니다.';
          _devices = _fallbackDevices;
          _isLoading = false;
        });
        return;
      }
      final items = resp['devices'] as List? ?? [];
      setState(() {
        _devices = items.cast<Map<String, dynamic>>();
        _error = null;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = '기기 목록을 불러오는 중 오류가 발생했습니다: $e';
        _devices = _fallbackDevices;
        _isLoading = false;
      });
    }
  }

  List<Map<String, dynamic>> get _filteredDevices {
    var result = _devices;

    // Status filter
    if (_statusFilter != 'all') {
      result = result
          .where((d) => (d['status'] as String? ?? '') == _statusFilter)
          .toList();
    }

    // Search filter
    if (_searchController.text.isNotEmpty) {
      final q = _searchController.text.toLowerCase();
      result = result.where((d) {
        final name = (d['name'] as String? ?? '').toLowerCase();
        final serial = (d['serial_number'] as String? ?? '').toLowerCase();
        final owner = (d['owner_name'] as String? ?? '').toLowerCase();
        return name.contains(q) || serial.contains(q) || owner.contains(q);
      }).toList();
    }

    return result;
  }

  @override
  Widget build(BuildContext context) {
    final filtered = _filteredDevices;

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
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
                      '기기 관리',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.refresh, color: Colors.white),
                    tooltip: '새로고침',
                    onPressed: _loadDevices,
                  ),
                ],
              ),
            ),

            // Summary cards
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Row(
                children: [
                  _buildSummaryChip(
                    '전체',
                    _devices.length,
                    SanggamTheme.jagaeCyan,
                  ),
                  const SizedBox(width: 8),
                  _buildSummaryChip(
                    '온라인',
                    _devices
                        .where((d) => d['status'] == 'online')
                        .length,
                    SanggamTheme.jagaeCyan,
                  ),
                  const SizedBox(width: 8),
                  _buildSummaryChip(
                    '오프라인',
                    _devices
                        .where((d) => d['status'] != 'online')
                        .length,
                    SanggamTheme.error,
                  ),
                ],
              ),
            ),
            const SizedBox(height: 8),

            // Error banner
            if (_error != null)
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(
                      horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    color: SanggamTheme.error.withValues(alpha: 0.15),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: SanggamTheme.error.withValues(alpha: 0.4),
                    ),
                  ),
                  child: Row(
                    children: [
                      const Icon(Icons.error_outline,
                          color: SanggamTheme.error, size: 18),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          _error!,
                          style: const TextStyle(
                            color: SanggamTheme.error,
                            fontSize: 12,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      GestureDetector(
                        onTap: _loadDevices,
                        child: const Text(
                          '재시도',
                          style: TextStyle(
                            color: SanggamTheme.error,
                            fontSize: 12,
                            fontWeight: FontWeight.bold,
                            decoration: TextDecoration.underline,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            if (_error != null) const SizedBox(height: 8),

            // Search bar
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: TextField(
                controller: _searchController,
                style: const TextStyle(color: Colors.white),
                decoration: InputDecoration(
                  hintText: '기기명, 시리얼, 소유자 검색',
                  hintStyle:
                      const TextStyle(color: SanggamTheme.onSurfaceDim),
                  prefixIcon: const Icon(Icons.search,
                      color: SanggamTheme.onSurfaceDim),
                  filled: true,
                  fillColor: SanggamTheme.surface,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: BorderSide.none,
                  ),
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide: const BorderSide(
                        color: SanggamTheme.surfaceVariant),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(16),
                    borderSide:
                        const BorderSide(color: SanggamTheme.primary),
                  ),
                ),
                onChanged: (_) => setState(() {}),
              ),
            ),
            const SizedBox(height: 8),

            // Status filter chips
            SizedBox(
              height: 36,
              child: ListView(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                children: _statusFilters.map((f) {
                  final isSelected = _statusFilter == f.$1;
                  return Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: GestureDetector(
                      onTap: () =>
                          setState(() => _statusFilter = f.$1),
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
                          f.$2,
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
            const Divider(
                height: 1, color: SanggamTheme.surfaceVariant),

            // Device count
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 16, vertical: 8),
              child: Row(
                children: [
                  Text('${filtered.length}개 기기',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      )),
                  const Spacer(),
                  const Text('최근 활성순',
                      style: TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 12,
                      )),
                ],
              ),
            ),

            // Device list
            Expanded(
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : filtered.isEmpty
                      ? const Center(
                          child: Column(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Icon(Icons.devices,
                                  size: 48,
                                  color:
                                      SanggamTheme.onSurfaceDim),
                              SizedBox(height: 8),
                              Text('등록된 기기가 없습니다.',
                                  style: TextStyle(
                                    color: Colors.white,
                                    fontSize: 14,
                                  )),
                            ],
                          ),
                        )
                      : RefreshIndicator(
                          color: SanggamTheme.primary,
                          onRefresh: _loadDevices,
                          child: ListView.builder(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 16),
                            itemCount: filtered.length,
                            itemBuilder: (_, index) =>
                                _buildDeviceTile(filtered[index]),
                          ),
                        ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryChip(String label, int count, Color color) {
    return Expanded(
      child: SanggamContainer(
        borderRadius: 16,
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Text('$count',
                style: TextStyle(
                  color: color,
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                )),
            Text(label,
                style: const TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 11,
                )),
          ],
        ),
      ),
    );
  }

  Widget _buildDeviceTile(Map<String, dynamic> device) {
    final name = device['name'] as String? ?? '만파식 리더기';
    final serial = device['serial_number'] as String? ??
        device['device_id'] as String? ??
        '';
    final status = device['status'] as String? ?? 'offline';
    final firmware = device['firmware_version'] as String? ?? '-';
    final battery = device['battery_percent'] as int? ??
        device['battery'] as int? ??
        0;
    final lastSeen = device['last_seen'] as String? ??
        device['last_seen_at'] as String? ??
        '';
    final ownerName = device['owner_name'] as String? ?? '';

    final isOnline = status == 'online';
    final statusColor =
        isOnline ? SanggamTheme.jagaeCyan : SanggamTheme.error;
    final statusLabel = isOnline ? '온라인' : '오프라인';

    final batteryColor = battery > 50
        ? SanggamTheme.jagaeCyan
        : battery > 20
            ? SanggamTheme.primary
            : SanggamTheme.error;

    return SanggamContainer(
      borderRadius: 16,
      padding: EdgeInsets.zero,
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        leading: Stack(
          children: [
            CircleAvatar(
              radius: 20,
              backgroundColor:
                  SanggamTheme.primary.withValues(alpha: 0.15),
              child: const Icon(Icons.devices,
                  color: SanggamTheme.primary, size: 20),
            ),
            Positioned(
              right: 0,
              bottom: 0,
              child: Container(
                width: 12,
                height: 12,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: statusColor,
                  border: Border.all(
                      color: SanggamTheme.background, width: 2),
                ),
              ),
            ),
          ],
        ),
        title: Row(
          children: [
            Flexible(
              child: Text(name,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                  overflow: TextOverflow.ellipsis),
            ),
            const SizedBox(width: 8),
            Container(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: statusColor.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(statusLabel,
                  style: TextStyle(
                    fontSize: 10,
                    color: statusColor,
                    fontWeight: FontWeight.w600,
                  )),
            ),
          ],
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 4),
            if (serial.isNotEmpty)
              Text('S/N: $serial',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 11,
                  )),
            Row(
              children: [
                Icon(Icons.memory, size: 12, color: SanggamTheme.onSurfaceDim),
                const SizedBox(width: 4),
                Text('FW $firmware',
                    style: const TextStyle(
                      color: SanggamTheme.onSurfaceDim,
                      fontSize: 11,
                    )),
                const SizedBox(width: 12),
                Icon(Icons.battery_std, size: 12, color: batteryColor),
                const SizedBox(width: 4),
                Text('$battery%',
                    style: TextStyle(
                      color: batteryColor,
                      fontSize: 11,
                      fontWeight: FontWeight.w600,
                    )),
              ],
            ),
            if (ownerName.isNotEmpty)
              Text('소유자: $ownerName',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 11,
                  )),
            if (lastSeen.isNotEmpty)
              Text('마지막 연결: ${_formatLastSeen(lastSeen)}',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 11,
                  )),
          ],
        ),
        isThreeLine: true,
      ),
    );
  }

  String _formatLastSeen(String lastSeen) {
    final dt = DateTime.tryParse(lastSeen);
    if (dt == null) return lastSeen;
    final diff = DateTime.now().difference(dt);
    if (diff.inMinutes < 1) return '방금 전';
    if (diff.inMinutes < 60) return '${diff.inMinutes}분 전';
    if (diff.inHours < 24) return '${diff.inHours}시간 전';
    if (diff.inDays < 7) return '${diff.inDays}일 전';
    return '${dt.month}/${dt.day}';
  }

  static const _statusFilters = [
    ('all', '전체'),
    ('online', '온라인'),
    ('offline', '오프라인'),
  ];

  static final _fallbackDevices = [
    {
      'name': '만파식 ONE #001',
      'device_id': 'MPK-001-A1B2',
      'serial_number': 'MPK-2026-001-A1B2',
      'status': 'online',
      'firmware_version': '2.1.0',
      'battery_percent': 92,
      'last_seen': DateTime.now()
          .subtract(const Duration(minutes: 3))
          .toIso8601String(),
      'owner_name': '김건강',
    },
    {
      'name': '만파식 ONE #002',
      'device_id': 'MPK-002-C3D4',
      'serial_number': 'MPK-2026-002-C3D4',
      'status': 'online',
      'firmware_version': '2.1.0',
      'battery_percent': 78,
      'last_seen': DateTime.now()
          .subtract(const Duration(minutes: 15))
          .toIso8601String(),
      'owner_name': '이안심',
    },
    {
      'name': '만파식 ONE #003',
      'device_id': 'MPK-003-E5F6',
      'serial_number': 'MPK-2026-003-E5F6',
      'status': 'offline',
      'firmware_version': '2.0.1',
      'battery_percent': 15,
      'last_seen': DateTime.now()
          .subtract(const Duration(days: 3))
          .toIso8601String(),
      'owner_name': '박자연',
    },
    {
      'name': '만파식 ONE #004',
      'device_id': 'MPK-004-G7H8',
      'serial_number': 'MPK-2026-004-G7H8',
      'status': 'online',
      'firmware_version': '2.1.0',
      'battery_percent': 55,
      'last_seen': DateTime.now()
          .subtract(const Duration(hours: 1))
          .toIso8601String(),
      'owner_name': '최희망',
    },
    {
      'name': '만파식 ONE #005',
      'device_id': 'MPK-005-I9J0',
      'serial_number': 'MPK-2026-005-I9J0',
      'status': 'offline',
      'firmware_version': '1.9.2',
      'battery_percent': 0,
      'last_seen': DateTime.now()
          .subtract(const Duration(days: 14))
          .toIso8601String(),
      'owner_name': '정미래',
    },
    {
      'name': '만파식 ONE #006',
      'device_id': 'MPK-006-K1L2',
      'serial_number': 'MPK-2026-006-K1L2',
      'status': 'online',
      'firmware_version': '2.1.0',
      'battery_percent': 100,
      'last_seen':
          DateTime.now().subtract(const Duration(seconds: 30)).toIso8601String(),
      'owner_name': '한솔루션',
    },
  ];
}
