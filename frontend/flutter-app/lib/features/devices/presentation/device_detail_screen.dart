import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/presentation/ble_scan_dialog.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// DeviceDetailScreen — Sanggam Orbit 기기 상세
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 제거
// [Rule 4] theme.textTheme.* 2x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold 5x → SanggamTheme.primary
// [Rule 4] Colors.green 2x → SanggamTheme.jagaeCyan
// [Rule 4] Colors.red 4x → SanggamTheme.error
// [Rule 4] Card 2x → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 4] FilledButton foregroundColor + borderRadius:16
// [Rule 2] horizontal:12→16, vertical:4→8, borderRadius:12→16, h:12→16
// ───────────────────────────────────────────────────

/// 기기 상세 관리 화면
///
/// 기기명 변경, 위치 설정, 펌웨어 정보, 배터리 상태 표시
class DeviceDetailScreen extends ConsumerStatefulWidget {
  const DeviceDetailScreen({super.key, required this.deviceId});

  final String deviceId;

  @override
  ConsumerState<DeviceDetailScreen> createState() =>
      _DeviceDetailScreenState();
}

class _DeviceDetailScreenState
    extends ConsumerState<DeviceDetailScreen> {
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
    final client = ref.watch(restClientProvider);

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
                      '기기 상세',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  TextButton(
                    onPressed: _saving ? null : _saveDevice,
                    child: _saving
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(
                                strokeWidth: 2),
                          )
                        : const Text(
                            '저장',
                            style: TextStyle(
                              color: SanggamTheme.primary,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: FutureBuilder<Map<String, dynamic>>(
                future: client.listDevices(
                    ref.read(authProvider).userId ?? ''),
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
                              backgroundColor: SanggamTheme
                                  .primary
                                  .withValues(alpha: 0.15),
                              child: const Icon(
                                  Icons.bluetooth_connected,
                                  size: 40,
                                  color: SanggamTheme.primary),
                            ),
                            const SizedBox(height: 8),
                            Container(
                              padding:
                                  const EdgeInsets.symmetric(
                                      horizontal: 16,
                                      vertical: 8),
                              decoration: BoxDecoration(
                                color: SanggamTheme.jagaeCyan
                                    .withValues(alpha: 0.1),
                                borderRadius:
                                    BorderRadius.circular(16),
                              ),
                              child: const Text(
                                '연결됨',
                                style: TextStyle(
                                  color:
                                      SanggamTheme.jagaeCyan,
                                  fontSize: 12,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
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
                          prefixIcon: Icon(
                              Icons.location_on_outlined),
                          border: OutlineInputBorder(),
                        ),
                      ),
                      const SizedBox(height: 24),

                      // 기기 정보
                      const Text(
                        '기기 정보',
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
                            _infoTile(Icons.fingerprint,
                                '기기 ID', widget.deviceId),
                            _infoTile(Icons.memory, '펌웨어',
                                'v2.1.3'),
                            _infoTile(
                                Icons.battery_charging_full,
                                '배터리',
                                '85%'),
                            _infoTile(
                                Icons.signal_cellular_alt,
                                '신호 강도',
                                '-42 dBm (우수)'),
                            _infoTile(Icons.access_time,
                                '마지막 동기화', '방금 전'),
                          ],
                        ),
                      ),
                      const SizedBox(height: 24),

                      // 카트리지 상태
                      const Text(
                        '장착된 카트리지',
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
                        child: ListTile(
                          leading: CircleAvatar(
                            backgroundColor: SanggamTheme
                                .primary
                                .withValues(alpha: 0.15),
                            child: const Icon(Icons.science,
                                color: SanggamTheme.primary),
                          ),
                          title: const Text('혈당 측정 카트리지'),
                          subtitle: const Text(
                              '잔여 횟수: 7/10 | 유효기한: 2026-08-15'),
                          trailing: const Icon(
                              Icons.chevron_right,
                              color:
                                  SanggamTheme.onSurfaceDim),
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
                        icon:
                            const Icon(Icons.system_update),
                        label:
                            const Text('펌웨어 업데이트 확인'),
                        style: FilledButton.styleFrom(
                          backgroundColor:
                              SanggamTheme.primary,
                          foregroundColor:
                              SanggamTheme.background,
                          minimumSize:
                              const Size.fromHeight(48),
                          shape: RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius.circular(16),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),

                      // 위험 영역
                      OutlinedButton.icon(
                        onPressed: () =>
                            _showUnpairDialog(context),
                        icon: const Icon(Icons.link_off,
                            color: SanggamTheme.error),
                        label: const Text(
                          '기기 연결 해제',
                          style: TextStyle(
                              color: SanggamTheme.error),
                        ),
                        style: OutlinedButton.styleFrom(
                          side: const BorderSide(
                              color: SanggamTheme.error),
                        ),
                      ),
                    ],
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _infoTile(
      IconData icon, String label, String value) {
    return ListTile(
      dense: true,
      leading: Icon(icon,
          size: 20, color: SanggamTheme.onSurfaceDim),
      title: Text(
        label,
        style: const TextStyle(
          color: SanggamTheme.onSurfaceDim,
          fontSize: 13,
        ),
      ),
      trailing: Text(
        value,
        style: const TextStyle(
          color: Colors.white,
          fontSize: 13,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Future<void> _saveDevice() async {
    setState(() => _saving = true);
    try {
      final client = ref.read(restClientProvider);
      await client.registerDevice(
        deviceId: widget.deviceId,
        userId: ref.read(authProvider).userId ?? '',
      );
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
              content: Text('기기 정보가 저장되었습니다.')),
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
        content: const Text(
            '이 기기와의 연결을 해제하시겠습니까? 이후 다시 페어링이 필요합니다.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('취소'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.error,
            ),
            onPressed: () {
              Navigator.pop(ctx);
              context.pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                    content: Text('기기 연결이 해제되었습니다.')),
              );
            },
            child: const Text('해제'),
          ),
        ],
      ),
    );
  }
}
