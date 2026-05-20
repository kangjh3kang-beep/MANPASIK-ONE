import 'package:flutter/material.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';

/// BLE 스캔 다이얼로그 (Rust FFI 스텁 연동)
void showBleScanDialog(BuildContext context) {
  showDialog<void>(
    context: context,
    builder: (context) => const _BleScanDialog(),
  );
}

class _BleScanDialog extends StatefulWidget {
  const _BleScanDialog();

  @override
  State<_BleScanDialog> createState() => _BleScanDialogState();
}

class _BleScanDialogState extends State<_BleScanDialog> {
  List<DeviceInfoDto> _devices = [];
  bool _loading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _scan();
  }

  Future<void> _scan() async {
    setState(() {
      _loading = true;
      _error = null;
      _devices = [];
    });
    try {
      final list = await RustBridge.bleScan();
      if (!mounted) return;
      setState(() {
        _devices = list;
        _loading = false;
      });
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: SanggamTheme.surface,
      title: const Text('BLE 디바이스 검색', style: TextStyle(color: Colors.white)),
      content: SizedBox(
        width: double.maxFinite,
        child: _loading
            ? const Padding(
                padding: EdgeInsets.all(24),
                child: Center(child: CircularProgressIndicator(color: SanggamTheme.primary)),
              )
            : _error != null
                ? Text(_error!, style: const TextStyle(color: SanggamTheme.error))
                : _devices.isEmpty
                    ? const Text('검색된 기기가 없습니다.\n(실제 기기 연동 시 Rust FFI ble_scan 사용)',
                        style: TextStyle(color: SanggamTheme.onSurfaceDim))
                    : SingleChildScrollView(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: _devices.map((d) => ListTile(
                            leading: const Icon(Icons.bluetooth, color: SanggamTheme.jagaeCyan),
                            title: Text(d.name.isNotEmpty ? d.name : d.deviceId,
                                style: const TextStyle(color: Colors.white)),
                            subtitle: Text('RSSI: ${d.rssi}',
                                style: const TextStyle(color: SanggamTheme.onSurfaceDim)),
                            onTap: () => Navigator.of(context).pop(d.deviceId),
                          )).toList(),
                        ),
                      ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text('닫기'),
        ),
        if (!_loading)
          TextButton(
            onPressed: _scan,
            child: const Text('다시 검색'),
          ),
      ],
    );
  }
}

/// 펌웨어 OTA 업데이트 다이얼로그
void showOtaUpdateDialog(BuildContext context, {required String deviceId, required String deviceName}) {
  showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (context) => _OtaUpdateDialog(deviceId: deviceId, deviceName: deviceName),
  );
}

class _OtaUpdateDialog extends StatefulWidget {
  const _OtaUpdateDialog({required this.deviceId, required this.deviceName});
  final String deviceId;
  final String deviceName;

  @override
  State<_OtaUpdateDialog> createState() => _OtaUpdateDialogState();
}

enum _OtaStage { checking, downloading, installing, complete, error }

class _OtaUpdateDialogState extends State<_OtaUpdateDialog> {
  _OtaStage _stage = _OtaStage.checking;
  double _progress = 0.0;
  String? _errorMessage;
  String _currentVersion = '1.2.0';
  String _newVersion = '1.3.0';

  @override
  void initState() {
    super.initState();
    _checkForUpdate();
  }

  Future<void> _checkForUpdate() async {
    setState(() => _stage = _OtaStage.checking);

    try {
      await Future.delayed(const Duration(seconds: 2));
    } catch (_) {
      // 버전 조회 실패 → 기본값 사용
    }

    if (!mounted) return;
    setState(() {
      _currentVersion = '1.2.0';
      _newVersion = '1.3.0';
      _stage = _OtaStage.downloading;
    });
    _downloadFirmware();
  }

  Future<void> _downloadFirmware() async {
    for (var i = 0; i <= 100; i += 5) {
      await Future.delayed(const Duration(milliseconds: 150));
      if (!mounted) return;
      setState(() => _progress = i / 100);
    }
    setState(() {
      _stage = _OtaStage.installing;
      _progress = 0.0;
    });
    _installFirmware();
  }

  Future<void> _installFirmware() async {
    for (var i = 0; i <= 100; i += 2) {
      await Future.delayed(const Duration(milliseconds: 200));
      if (!mounted) return;
      setState(() => _progress = i / 100);
    }
    setState(() => _stage = _OtaStage.complete);
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: SanggamTheme.surface,
      title: const Row(
        children: [
          Icon(Icons.system_update, size: 24, color: SanggamTheme.primary),
          SizedBox(width: 8),
          Expanded(child: Text('펌웨어 업데이트', style: TextStyle(color: Colors.white))),
        ],
      ),
      content: SizedBox(
        width: double.maxFinite,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 기기 정보
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: SanggamTheme.surfaceVariant,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  const Icon(Icons.bluetooth_connected, size: 20, color: SanggamTheme.jagaeCyan),
                  const SizedBox(width: 8),
                  Expanded(child: Text(widget.deviceName, style: const TextStyle(color: Colors.white, fontSize: 14))),
                ],
              ),
            ),
            const SizedBox(height: 16),

            // 버전 정보
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                _buildVersionChip('현재', _currentVersion, false),
                const Padding(
                  padding: EdgeInsets.symmetric(horizontal: 8),
                  child: Icon(Icons.arrow_forward, size: 16, color: SanggamTheme.onSurfaceDim),
                ),
                _buildVersionChip('최신', _newVersion, true),
              ],
            ),
            const SizedBox(height: 20),

            // 진행 상태
            _buildStageContent(),
          ],
        ),
      ),
      actions: [
        if (_stage == _OtaStage.error)
          TextButton(
            onPressed: _checkForUpdate,
            child: const Text('다시 시도'),
          ),
        if (_stage == _OtaStage.complete || _stage == _OtaStage.error)
          FilledButton(
            onPressed: () => Navigator.of(context).pop(),
            style: FilledButton.styleFrom(backgroundColor: SanggamTheme.primary),
            child: Text(_stage == _OtaStage.complete ? '완료' : '닫기'),
          ),
        if (_stage == _OtaStage.checking || _stage == _OtaStage.downloading || _stage == _OtaStage.installing)
          TextButton(
            onPressed: null,
            child: Text(_stageLabel()),
          ),
      ],
    );
  }

  Widget _buildVersionChip(String label, String version, bool isNew) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: isNew ? SanggamTheme.primary.withValues(alpha: 0.1) : SanggamTheme.surfaceVariant,
        borderRadius: BorderRadius.circular(16),
        border: isNew ? Border.all(color: SanggamTheme.primary, width: 1) : null,
      ),
      child: Column(
        children: [
          Text(label, style: const TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12)),
          Text('v$version', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 14)),
        ],
      ),
    );
  }

  Widget _buildStageContent() {
    switch (_stage) {
      case _OtaStage.checking:
        return const Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: SanggamTheme.primary)),
            SizedBox(width: 12),
            Text('업데이트 확인 중...', style: TextStyle(color: Colors.white)),
          ],
        );
      case _OtaStage.downloading:
        return Column(
          children: [
            const Text('펌웨어 다운로드 중...', style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12)),
            const SizedBox(height: 8),
            LinearProgressIndicator(value: _progress, color: SanggamTheme.primary),
            const SizedBox(height: 4),
            Text('${(_progress * 100).toInt()}%', style: const TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12)),
          ],
        );
      case _OtaStage.installing:
        return Column(
          children: [
            const Text('기기에 설치 중... (전원을 끄지 마세요)', style: TextStyle(color: SanggamTheme.primary, fontSize: 12)),
            const SizedBox(height: 8),
            LinearProgressIndicator(value: _progress, color: SanggamTheme.primary),
            const SizedBox(height: 4),
            Text('${(_progress * 100).toInt()}%', style: const TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12)),
          ],
        );
      case _OtaStage.complete:
        return Column(
          children: [
            Icon(Icons.check_circle, size: 48, color: SanggamTheme.jagaeCyan.withValues(alpha: 0.8)),
            const SizedBox(height: 8),
            const Text('업데이트 완료!', style: TextStyle(color: SanggamTheme.jagaeCyan, fontSize: 16, fontWeight: FontWeight.bold)),
            Text('v$_newVersion이 설치되었습니다.', style: const TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12)),
          ],
        );
      case _OtaStage.error:
        return Column(
          children: [
            const Icon(Icons.error_outline, size: 48, color: SanggamTheme.error),
            const SizedBox(height: 8),
            Text(_errorMessage ?? '업데이트에 실패했습니다.', style: const TextStyle(color: SanggamTheme.error, fontSize: 12)),
          ],
        );
    }
  }

  String _stageLabel() {
    switch (_stage) {
      case _OtaStage.checking:
        return '확인 중...';
      case _OtaStage.downloading:
        return '다운로드 중...';
      case _OtaStage.installing:
        return '설치 중...';
      default:
        return '';
    }
  }
}
