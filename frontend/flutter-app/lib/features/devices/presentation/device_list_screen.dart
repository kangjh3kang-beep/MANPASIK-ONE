import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';
import 'package:manpasik/features/devices/presentation/ble_scan_dialog.dart';

/// 디바이스 목록 화면
///
/// device-service ListDevices gRPC로 등록된 디바이스 로드.
class DeviceListScreen extends ConsumerWidget {
  const DeviceListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final devicesAsync = ref.watch(deviceListProvider);

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('디바이스'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add_circle_outline),
            tooltip: '디바이스 검색 (BLE)',
            onPressed: () => showBleScanDialog(context),
          ),
        ],
      ),
      body: devicesAsync.when(
        data: (devices) => devices.isEmpty
            ? _buildEmptyState(context, theme)
            : _buildDeviceListWrapper(theme, devices, ref),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, _) => Center(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Text(
              '디바이스 목록을 불러올 수 없습니다.\n$err',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.error),
            ),
          ),
        ),
      ),
    );
  }

  /// 빈 상태 UI
  Widget _buildEmptyState(BuildContext context, ThemeData theme) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(48),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 120,
              height: 120,
              decoration: BoxDecoration(
                color: theme.colorScheme.surfaceContainerHighest,
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.devices_rounded,
                size: 56,
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              '등록된 디바이스가 없습니다',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              '우측 상단의 + 버튼을 눌러\n새 디바이스를 등록해주세요',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 32),
            FilledButton.icon(
              onPressed: () => showBleScanDialog(context),
              icon: const Icon(Icons.bluetooth_searching_rounded),
              label: const Text('디바이스 검색'),
              style: FilledButton.styleFrom(
                minimumSize: const Size(200, 56),
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

  /// RefreshIndicator 감싸기
  Widget _buildDeviceListWrapper(ThemeData theme, List<DeviceItem> devices, WidgetRef ref) {
    return RefreshIndicator(
      onRefresh: () async => ref.invalidate(deviceListProvider),
      child: _buildDeviceList(theme, devices),
    );
  }

  /// 디바이스 목록
  Widget _buildDeviceList(ThemeData theme, List<DeviceItem> devices) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: devices.length,
      itemBuilder: (context, index) {
        final device = devices[index];
        final isConnected = device.status == 'online';
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: ListTile(
            contentPadding: const EdgeInsets.all(16),
            onTap: () => context.push('/devices/${device.deviceId}'),
            leading: Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: isConnected
                    ? Colors.green.withValues(alpha: 0.1)
                    : theme.colorScheme.surfaceContainerHighest,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(
                isConnected
                    ? Icons.bluetooth_connected_rounded
                    : Icons.bluetooth_disabled_rounded,
                color: isConnected
                    ? Colors.green
                    : theme.colorScheme.onSurfaceVariant,
              ),
            ),
            title: Text(
              device.name,
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.w600,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isConnected
                      ? '연결됨'
                      : (device.status == 'measuring'
                          ? '측정 중'
                          : device.status == 'offline'
                              ? '연결 안됨'
                              : device.status),
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: isConnected ? Colors.green : theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(height: 4),
                _UsageDropdown(deviceId: device.deviceId),
              ],
            ),
            trailing: const Icon(Icons.chevron_right),
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
  State<_UsageDropdown> createState() => _UsageDropdownState();
}

class _UsageDropdownState extends State<_UsageDropdown> {
  String _usage = '개인';

  static const _options = ['개인', '가정', '사무실'];

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(Icons.category_outlined, size: 14, color: Theme.of(context).colorScheme.outline),
        const SizedBox(width: 4),
        DropdownButton<String>(
          value: _usage,
          isDense: true,
          underline: const SizedBox.shrink(),
          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                color: Theme.of(context).colorScheme.primary,
              ),
          items: _options.map((o) => DropdownMenuItem(value: o, child: Text(o))).toList(),
          onChanged: (v) {
            if (v != null) setState(() => _usage = v);
          },
        ),
      ],
    );
  }
}
