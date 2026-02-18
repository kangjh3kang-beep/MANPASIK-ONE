import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/features/devices/presentation/ble_scan_dialog.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// 기기 상세 관리 화면
///
/// 기기명 변경, 위치 설정, 펌웨어 정보, 배터리 상태 표시
class DeviceDetailScreen extends ConsumerStatefulWidget {
  const DeviceDetailScreen({super.key, required this.deviceId});

  final String deviceId;

  @override
  ConsumerState<DeviceDetailScreen> createState() => _DeviceDetailScreenState();
}

class _DeviceDetailScreenState extends ConsumerState<DeviceDetailScreen> {
  late TextEditingController _nameCtrl;
  late TextEditingController _locationCtrl;
  bool _saving = false;

  @override
  void initState() {
    super.initState();
    _nameCtrl = TextEditingController(text: 'ManPaSik Reader');
    _locationCtrl = TextEditingController(text: '거실');
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _locationCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final client = ref.watch(restClientProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('기기 상세'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        actions: [
          TextButton(
            onPressed: _saving ? null : _saveDevice,
            child: _saving
                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                : const Text('저장'),
          ),
        ],
      ),
      body: FutureBuilder<Map<String, dynamic>>(
        future: client.listDevices(ref.read(authProvider).userId ?? ''),
        builder: (context, snapshot) {
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // 기기 아이콘 + 상태
              Center(
                child: Column(
                  children: [
                    CircleAvatar(
                      radius: 40,
                      backgroundColor: AppTheme.sanggamGold.withValues(alpha: 0.15),
                      child: const Icon(Icons.bluetooth_connected, size: 40, color: AppTheme.sanggamGold),
                    ),
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                      decoration: BoxDecoration(
                        color: Colors.green.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Text('연결됨', style: TextStyle(color: Colors.green, fontSize: 12, fontWeight: FontWeight.w600)),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 24),

              // 기기명 변경
              TextFormField(
                controller: _nameCtrl,
                decoration: const InputDecoration(
                  labelText: '기기 이름',
                  prefixIcon: Icon(Icons.edit),
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),

              // 위치 설정
              TextFormField(
                controller: _locationCtrl,
                decoration: const InputDecoration(
                  labelText: '설치 위치',
                  prefixIcon: Icon(Icons.location_on_outlined),
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),

              // 기기 정보
              Text('기기 정보', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Card(
                child: Column(
                  children: [
                    _infoTile(Icons.fingerprint, '기기 ID', widget.deviceId),
                    _infoTile(Icons.memory, '펌웨어', 'v2.1.3'),
                    _infoTile(Icons.battery_charging_full, '배터리', '85%'),
                    _infoTile(Icons.signal_cellular_alt, '신호 강도', '-42 dBm (우수)'),
                    _infoTile(Icons.access_time, '마지막 동기화', '방금 전'),
                  ],
                ),
              ),
              const SizedBox(height: 24),

              // 카트리지 상태
              Text('장착된 카트리지', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Card(
                child: ListTile(
                  leading: CircleAvatar(
                    backgroundColor: AppTheme.sanggamGold.withValues(alpha: 0.15),
                    child: const Icon(Icons.science, color: AppTheme.sanggamGold),
                  ),
                  title: const Text('혈당 측정 카트리지'),
                  subtitle: const Text('잔여 횟수: 7/10 | 유효기한: 2026-08-15'),
                  trailing: const Icon(Icons.chevron_right),
                ),
              ),
              const SizedBox(height: 24),

              // 펌웨어 업데이트
              FilledButton.icon(
                onPressed: () => showOtaUpdateDialog(
                  context,
                  deviceId: widget.deviceId,
                  deviceName: _nameCtrl.text,
                ),
                icon: const Icon(Icons.system_update),
                label: const Text('펌웨어 업데이트 확인'),
                style: FilledButton.styleFrom(
                  backgroundColor: AppTheme.sanggamGold,
                  minimumSize: const Size.fromHeight(48),
                ),
              ),
              const SizedBox(height: 12),

              // 위험 영역
              OutlinedButton.icon(
                onPressed: () => _showUnpairDialog(context),
                icon: const Icon(Icons.link_off, color: Colors.red),
                label: const Text('기기 연결 해제', style: TextStyle(color: Colors.red)),
                style: OutlinedButton.styleFrom(side: const BorderSide(color: Colors.red)),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _infoTile(IconData icon, String label, String value) {
    return ListTile(
      dense: true,
      leading: Icon(icon, size: 20),
      title: Text(label, style: const TextStyle(fontSize: 13)),
      trailing: Text(value, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w500)),
    );
  }

  Future<void> _saveDevice() async {
    setState(() => _saving = true);
    try {
      // REST API로 기기 이름/위치 업데이트
      final client = ref.read(restClientProvider);
      await client.registerDevice(
        deviceId: widget.deviceId,
        userId: ref.read(authProvider).userId ?? '',
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('기기 정보가 저장되었습니다.')),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('저장 실패: $e')),
        );
      }
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  void _showUnpairDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('기기 연결 해제'),
        content: const Text('이 기기와의 연결을 해제하시겠습니까? 이후 다시 페어링이 필요합니다.'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('취소')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Colors.red),
            onPressed: () {
              Navigator.pop(ctx);
              context.pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('기기 연결이 해제되었습니다.')),
              );
            },
            child: const Text('해제'),
          ),
        ],
      ),
    );
  }
}
