import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/devices/presentation/ble_scan_dialog.dart';

// ───────────────────────────────────────────────────
// DeviceListScreen — Sanggam Orbit 디바이스 목록
//
// [Rule 4] sanggam_theme.dart import 추가
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 4x 제거
// [Rule 4] ThemeData 파라미터 4개 제거
// [Rule 4] theme.textTheme.* ~7x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~9x → SanggamTheme 상수
// [Rule 4] Colors.green 3x → SanggamTheme.jagaeCyan
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] bottom:12→16, h:12→16, h:4→8, borderRadius:12→16
// ───────────────────────────────────────────────────

/// 디바이스 목록 화면
class DeviceListScreen extends ConsumerWidget {
  const DeviceListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final devicesAsync = ref.watch(deviceListProvider);

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
                      '디바이스',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(
                        Icons.add_circle_outline,
                        color: Colors.white),
                    tooltip: '디바이스 검색 (BLE)',
                    onPressed: () =>
                        showBleScanDialog(context),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: devicesAsync.when(
                data: (devices) => devices.isEmpty
                    ? _buildEmptyState(context)
                    : _buildDeviceListWrapper(
                        devices, ref),
                loading: () => const Center(
                    child: CircularProgressIndicator()),
                error: (err, _) => Center(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Text(
                      '디바이스 목록을 불러올 수 없습니다.\n$err',
                      textAlign: TextAlign.center,
                      style: const TextStyle(
                        color: SanggamTheme.error,
                        fontSize: 14,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(48),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120,
              height: 120,
              decoration: const BoxDecoration(
                color: SanggamTheme.surfaceVariant,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.devices_rounded,
                size: 56,
                color: SanggamTheme.onSurfaceDim,
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              '등록된 디바이스가 없습니다',
              style: TextStyle(
                color: Colors.white,
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              '우측 상단의 + 버튼을 눌러\n새 디바이스를 등록해주세요',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
              ),
            ),
            const SizedBox(height: 32),
            FilledButton.icon(
              onPressed: () => showBleScanDialog(context),
              icon: const Icon(
                  Icons.bluetooth_searching_rounded),
              label: const Text('디바이스 검색'),
              style: FilledButton.styleFrom(
                minimumSize: const Size(200, 56),
                backgroundColor: SanggamTheme.primary,
                foregroundColor: SanggamTheme.background,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDeviceListWrapper(
      List<DeviceItem> devices, WidgetRef ref) {
    return RefreshIndicator(
      onRefresh: () async =>
          ref.invalidate(deviceListProvider),
      child: _buildDeviceList(devices),
    );
  }

  Widget _buildDeviceList(List<DeviceItem> devices) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: devices.length,
      itemBuilder: (context, index) {
        final device = devices[index];
        final isConnected = device.status == 'online';
        return Padding(
          padding: const EdgeInsets.only(bottom: 16),
          child: ListTile(
            tileColor: SanggamTheme.surface,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(16),
            ),
            contentPadding: const EdgeInsets.all(16),
            onTap: () => context
                .push('/devices/${device.deviceId}'),
            leading: Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: isConnected
                    ? SanggamTheme.jagaeCyan
                        .withValues(alpha: 0.1)
                    : SanggamTheme.surfaceVariant,
                borderRadius: BorderRadius.circular(16),
              ),
              child: Icon(
                isConnected
                    ? Icons.bluetooth_connected_rounded
                    : Icons.bluetooth_disabled_rounded,
                color: isConnected
                    ? SanggamTheme.jagaeCyan
                    : SanggamTheme.onSurfaceDim,
              ),
            ),
            title: Text(
              device.name,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,
              children: [
                Text(
                  isConnected
                      ? '연결됨'
                      : (device.status == 'measuring'
                          ? '측정 중'
                          : device.status == 'offline'
                              ? '연결 안됨'
                              : device.status),
                  style: TextStyle(
                    color: isConnected
                        ? SanggamTheme.jagaeCyan
                        : SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
                const SizedBox(height: 8),
                _UsageDropdown(
                    deviceId: device.deviceId),
              ],
            ),
            trailing: const Icon(Icons.chevron_right,
                color: SanggamTheme.onSurfaceDim),
          ),
        );
      },
    );
  }
}

/// 기기 용도별 분류 드롭다운
class _UsageDropdown extends StatefulWidget {
  const _UsageDropdown({required this.deviceId});
  final String deviceId;

  @override
  State<_UsageDropdown> createState() =>
      _UsageDropdownState();
}

class _UsageDropdownState extends State<_UsageDropdown> {
  String _usage = '개인';

  static const _options = ['개인', '가정', '사무실'];

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        const Icon(Icons.category_outlined,
            size: 14, color: SanggamTheme.onSurfaceDim),
        const SizedBox(width: 4),
        DropdownButton<String>(
          value: _usage,
          isDense: true,
          underline: const SizedBox.shrink(),
          style: const TextStyle(
            color: SanggamTheme.primary,
            fontSize: 12,
          ),
          items: _options
              .map((o) => DropdownMenuItem(
                  value: o, child: Text(o)))
              .toList(),
          onChanged: (v) {
            if (v != null) setState(() => _usage = v);
          },
        ),
      ],
    );
  }
}
