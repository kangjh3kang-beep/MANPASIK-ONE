import 'package:flutter/material.dart';
import 'package:manpasik/core/services/rust_ffi_stub.dart';

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
      final list = await RustFfiStub.bleScan();
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
    final theme = Theme.of(context);
    return AlertDialog(
      title: const Text('BLE 디바이스 검색'),
      content: SizedBox(
        width: double.maxFinite,
        child: _loading
            ? const Padding(
                padding: EdgeInsets.all(24),
                child: Center(child: CircularProgressIndicator()),
              )
            : _error != null
                ? Text(_error!, style: TextStyle(color: theme.colorScheme.error))
                : _devices.isEmpty
                    ? const Text('검색된 기기가 없습니다.\n(실제 기기 연동 시 Rust FFI ble_scan 사용)')
                    : SingleChildScrollView(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: _devices.map((d) => ListTile(
                            leading: const Icon(Icons.bluetooth),
                            title: Text(d.name.isNotEmpty ? d.name : d.deviceId),
                            subtitle: Text('RSSI: ${d.rssi}'),
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
