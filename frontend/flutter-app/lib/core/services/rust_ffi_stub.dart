/// Rust FFI 스텁 (S5b)
///
/// flutter_rust_bridge 활성화 전까지 BLE/NFC/엔진 API 호출을 스텁으로 처리.
/// 실제 FFI 연동 시 이 파일을 생성된 API 래퍼로 교체.

/// BLE 디바이스 정보 (Rust DeviceInfoDto 대응)
class DeviceInfoDto {
  final String deviceId;
  final String name;
  final int rssi;
  final String state;

  const DeviceInfoDto({
    required this.deviceId,
    required this.name,
    required this.rssi,
    required this.state,
  });
}

/// 카트리지 정보 (Rust CartridgeInfoDto 대응)
class CartridgeInfoDto {
  final String cartridgeId;
  final String cartridgeType;
  final String lotId;
  final String expiryDate;
  final int remainingUses;

  const CartridgeInfoDto({
    required this.cartridgeId,
    required this.cartridgeType,
    required this.lotId,
    required this.expiryDate,
    required this.remainingUses,
  });
}

/// Rust FFI API 스텁 (실제 연동 시 native 바인딩으로 교체)
class RustFfiStub {
  RustFfiStub._();

  static String get engineVersion => '0.1.0-stub';

  static Future<List<DeviceInfoDto>> bleScan() async {
    await Future.delayed(const Duration(milliseconds: 500));
    return [];
  }

  static Future<bool> bleConnect(String deviceId) async {
    await Future.delayed(const Duration(milliseconds: 300));
    return false;
  }

  static Future<CartridgeInfoDto> nfcReadCartridge() async {
    await Future.delayed(const Duration(milliseconds: 300));
    return const CartridgeInfoDto(
      cartridgeId: 'stub-cartridge-1',
      cartridgeType: 'Glucose',
      lotId: 'LOT-STUB',
      expiryDate: '20261231',
      remainingUses: 10,
    );
  }
}
