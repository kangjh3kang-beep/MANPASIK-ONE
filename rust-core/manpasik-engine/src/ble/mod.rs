//! BLE 5.0 GATT 통신 모듈
//!
//! 리더기와의 BLE GATT 통신 - btleplug 기반 구현

use parking_lot::RwLock;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use thiserror::Error;

#[cfg(feature = "ble")]
use btleplug::api::{Central, Manager as _, Peripheral as _, ScanFilter};
#[cfg(feature = "ble")]
use btleplug::platform::{Adapter, Manager, Peripheral};

/// BLE 연결 상태
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ConnectionState {
    Disconnected,
    Connecting,
    Connected,
    Measuring,
    Error,
}

/// BLE 디바이스 정보
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DeviceInfo {
    pub device_id: String,
    pub name: String,
    pub rssi: i8,
    pub state: ConnectionState,
    pub firmware_version: Option<String>,
    pub battery_level: Option<u8>,
}

/// ManPaSik 리더기 GATT 서비스 UUID
pub mod service_uuids {
    /// 메인 측정 서비스
    pub const MEASUREMENT_SERVICE: &str = "0000fff0-0000-1000-8000-00805f9b34fb";
    /// 디바이스 정보 서비스
    pub const DEVICE_INFO_SERVICE: &str = "0000180a-0000-1000-8000-00805f9b34fb";
    /// 배터리 서비스
    pub const BATTERY_SERVICE: &str = "0000180f-0000-1000-8000-00805f9b34fb";
}

/// ManPaSik 리더기 GATT 특성 UUID
pub mod characteristic_uuids {
    /// 측정 명령 특성 (Write)
    pub const MEASUREMENT_COMMAND: &str = "0000fff1-0000-1000-8000-00805f9b34fb";
    /// 측정 데이터 특성 (Notify)
    pub const MEASUREMENT_DATA: &str = "0000fff2-0000-1000-8000-00805f9b34fb";
    /// 측정 상태 특성 (Read/Notify)
    pub const MEASUREMENT_STATUS: &str = "0000fff3-0000-1000-8000-00805f9b34fb";
    /// 보정 데이터 특성 (Read/Write)
    pub const CALIBRATION_DATA: &str = "0000fff4-0000-1000-8000-00805f9b34fb";
    /// 펌웨어 버전 특성
    pub const FIRMWARE_VERSION: &str = "00002a26-0000-1000-8000-00805f9b34fb";
    /// 배터리 레벨 특성
    pub const BATTERY_LEVEL: &str = "00002a19-0000-1000-8000-00805f9b34fb";
}

/// BLE 명령 코드
#[derive(Debug, Clone, Copy)]
pub enum BleCommand {
    StartMeasurement = 0x01,
    StopMeasurement = 0x02,
    GetStatus = 0x03,
    StartCalibration = 0x04,
    SetParameters = 0x05,
    Reset = 0xFF,
}

/// 측정 데이터 패킷 (BLE로 수신)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MeasurementDataPacket {
    /// 패킷 시퀀스 번호
    pub sequence: u16,
    /// 채널 데이터 (88 또는 896 채널)
    pub channels: Vec<f32>,
    /// 온도 (섭씨)
    pub temperature: f32,
    /// 습도 (%)
    pub humidity: f32,
    /// 배터리 레벨 (%)
    pub battery: u8,
    /// 타임스탬프 (밀리초)
    pub timestamp_ms: u64,
}

#[derive(Debug, Error)]
pub enum BleError {
    #[error("BLE 어댑터를 찾을 수 없습니다")]
    NoAdapter,

    #[error("디바이스를 찾을 수 없습니다: {0}")]
    DeviceNotFound(String),

    #[error("연결 실패: {0}")]
    ConnectionFailed(String),

    #[error("서비스를 찾을 수 없습니다: {0}")]
    ServiceNotFound(String),

    #[error("특성을 찾을 수 없습니다: {0}")]
    CharacteristicNotFound(String),

    #[error("읽기 실패: {0}")]
    ReadError(String),

    #[error("쓰기 실패: {0}")]
    WriteError(String),

    #[error("타임아웃")]
    Timeout,
}

/// BLE 매니저 - 여러 리더기 동시 관리 (무제한 확장)
pub struct BleManager {
    /// 연결된 디바이스 맵 (device_id -> DeviceInfo)
    connected_devices: Arc<RwLock<HashMap<String, DeviceInfo>>>,
    /// 스캔 결과 캐시
    scan_cache: Arc<RwLock<Vec<DeviceInfo>>>,
}

impl BleManager {
    pub fn new() -> Self {
        Self {
            connected_devices: Arc::new(RwLock::new(HashMap::new())),
            scan_cache: Arc::new(RwLock::new(Vec::new())),
        }
    }

    /// BLE 디바이스 스캔
    pub async fn scan(&self) -> Vec<DeviceInfo> {
        self.scan_with_timeout(std::time::Duration::from_secs(5))
            .await
    }

    /// 타임아웃 지정 스캔
    pub async fn scan_with_timeout(&self, _timeout: std::time::Duration) -> Vec<DeviceInfo> {
        #[cfg(feature = "ble")]
        {
            // btleplug 기반 실제 스캔 구현
            match Manager::new().await {
                Ok(manager) => {
                    if let Ok(adapters) = manager.adapters().await {
                        if let Some(adapter) = adapters.into_iter().next() {
                            let _ = adapter.start_scan(ScanFilter::default()).await;
                            tokio::time::sleep(_timeout).await;
                            let _ = adapter.stop_scan().await;

                            if let Ok(peripherals) = adapter.peripherals().await {
                                let devices: Vec<DeviceInfo> = peripherals
                                    .iter()
                                    .filter_map(|p| {
                                        let properties =
                                            futures::executor::block_on(p.properties()).ok()??;
                                        let name = properties.local_name.unwrap_or_default();

                                        // ManPaSik 리더기 필터링 (이름에 "MPK" 포함)
                                        if name.contains("MPK") || name.contains("ManPaSik") {
                                            Some(DeviceInfo {
                                                device_id: p.id().to_string(),
                                                name,
                                                rssi: properties.rssi.unwrap_or(0) as i8,
                                                state: ConnectionState::Disconnected,
                                                firmware_version: None,
                                                battery_level: None,
                                            })
                                        } else {
                                            None
                                        }
                                    })
                                    .collect();

                                *self.scan_cache.write() = devices.clone();
                                return devices;
                            }
                        }
                    }
                }
                Err(_) => {}
            }
        }

        // BLE 기능이 비활성화된 경우 빈 목록 반환
        Vec::new()
    }

    /// 디바이스 연결
    pub async fn connect(&mut self, device_id: &str) -> Result<(), BleError> {
        // 이미 연결된 경우 스킵
        if self.is_connected(device_id) {
            return Ok(());
        }

        #[cfg(feature = "ble")]
        {
            // btleplug 기반 실제 연결 구현
            // TODO: 실제 구현
        }

        // 연결 상태 업데이트
        let device_info = DeviceInfo {
            device_id: device_id.to_string(),
            name: format!("MPK-{}", &device_id[..8.min(device_id.len())]),
            rssi: 0,
            state: ConnectionState::Connected,
            firmware_version: Some("1.0.0".to_string()),
            battery_level: Some(100),
        };

        self.connected_devices
            .write()
            .insert(device_id.to_string(), device_info);
        Ok(())
    }

    /// 디바이스 연결 해제
    pub async fn disconnect(&mut self, device_id: &str) -> Result<(), BleError> {
        self.connected_devices.write().remove(device_id);
        Ok(())
    }

    /// 연결 상태 확인
    pub fn is_connected(&self, device_id: &str) -> bool {
        self.connected_devices.read().contains_key(device_id)
    }

    /// 연결된 디바이스 수
    pub fn connected_count(&self) -> usize {
        self.connected_devices.read().len()
    }

    /// 연결된 모든 디바이스 목록
    pub fn connected_devices(&self) -> Vec<DeviceInfo> {
        self.connected_devices.read().values().cloned().collect()
    }

    /// 측정 시작 명령 전송
    pub async fn start_measurement(&self, device_id: &str) -> Result<(), BleError> {
        if !self.is_connected(device_id) {
            return Err(BleError::DeviceNotFound(device_id.to_string()));
        }

        // 상태 업데이트
        if let Some(device) = self.connected_devices.write().get_mut(device_id) {
            device.state = ConnectionState::Measuring;
        }

        // TODO: 실제 BLE 명령 전송
        Ok(())
    }

    /// 측정 중지 명령 전송
    pub async fn stop_measurement(&self, device_id: &str) -> Result<(), BleError> {
        if !self.is_connected(device_id) {
            return Err(BleError::DeviceNotFound(device_id.to_string()));
        }

        // 상태 업데이트
        if let Some(device) = self.connected_devices.write().get_mut(device_id) {
            device.state = ConnectionState::Connected;
        }

        Ok(())
    }

    /// 배터리 레벨 읽기
    pub async fn read_battery_level(&self, device_id: &str) -> Result<u8, BleError> {
        if !self.is_connected(device_id) {
            return Err(BleError::DeviceNotFound(device_id.to_string()));
        }

        // TODO: 실제 BLE 읽기
        Ok(self
            .connected_devices
            .read()
            .get(device_id)
            .and_then(|d| d.battery_level)
            .unwrap_or(100))
    }

    /// 펌웨어 버전 읽기
    pub async fn read_firmware_version(&self, device_id: &str) -> Result<String, BleError> {
        if !self.is_connected(device_id) {
            return Err(BleError::DeviceNotFound(device_id.to_string()));
        }

        Ok(self
            .connected_devices
            .read()
            .get(device_id)
            .and_then(|d| d.firmware_version.clone())
            .unwrap_or_else(|| "Unknown".to_string()))
    }
}

impl Default for BleManager {
    fn default() -> Self {
        Self::new()
    }
}

/// 측정 데이터 파서
pub fn parse_measurement_packet(data: &[u8]) -> Result<MeasurementDataPacket, BleError> {
    if data.len() < 16 {
        return Err(BleError::ReadError("패킷이 너무 짧습니다".to_string()));
    }

    // 패킷 구조:
    // [0-1]: 시퀀스 번호 (u16 LE)
    // [2-3]: 채널 수 (u16 LE)
    // [4-7]: 온도 (f32 LE)
    // [8-11]: 습도 (f32 LE)
    // [12]: 배터리 (u8)
    // [13-16]: 타임스탬프 하위 4바이트
    // [17+]: 채널 데이터

    let sequence = u16::from_le_bytes([data[0], data[1]]);
    let channel_count = u16::from_le_bytes([data[2], data[3]]) as usize;
    let temperature = f32::from_le_bytes([data[4], data[5], data[6], data[7]]);
    let humidity = f32::from_le_bytes([data[8], data[9], data[10], data[11]]);
    let battery = data[12];
    let timestamp_ms = u32::from_le_bytes([data[13], data[14], data[15], data[16]]) as u64;

    // 채널 데이터 파싱
    let mut channels = Vec::with_capacity(channel_count);
    let channel_data_start = 17;
    for i in 0..channel_count {
        let offset = channel_data_start + i * 4;
        if offset + 4 <= data.len() {
            let value = f32::from_le_bytes([
                data[offset],
                data[offset + 1],
                data[offset + 2],
                data[offset + 3],
            ]);
            channels.push(value);
        }
    }

    Ok(MeasurementDataPacket {
        sequence,
        channels,
        temperature,
        humidity,
        battery,
        timestamp_ms,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ble_manager_creation() {
        let manager = BleManager::new();
        assert_eq!(manager.connected_count(), 0);
    }

    #[tokio::test]
    async fn test_connect_disconnect() {
        let mut manager = BleManager::new();

        manager.connect("test-device-001").await.unwrap();
        assert!(manager.is_connected("test-device-001"));
        assert_eq!(manager.connected_count(), 1);

        manager.disconnect("test-device-001").await.unwrap();
        assert!(!manager.is_connected("test-device-001"));
        assert_eq!(manager.connected_count(), 0);
    }

    #[test]
    fn test_parse_measurement_packet() {
        // 최소 패킷 생성
        let mut data = vec![0u8; 17 + 4 * 4]; // 4채널 데이터

        // 시퀀스: 1
        data[0] = 1;
        data[1] = 0;

        // 채널 수: 4
        data[2] = 4;
        data[3] = 0;

        // 온도: 23.5
        let temp_bytes = 23.5f32.to_le_bytes();
        data[4..8].copy_from_slice(&temp_bytes);

        // 습도: 45.0
        let humidity_bytes = 45.0f32.to_le_bytes();
        data[8..12].copy_from_slice(&humidity_bytes);

        // 배터리: 85%
        data[12] = 85;

        let packet = parse_measurement_packet(&data).unwrap();
        assert_eq!(packet.sequence, 1);
        assert!((packet.temperature - 23.5).abs() < 0.01);
        assert_eq!(packet.battery, 85);
    }
}
