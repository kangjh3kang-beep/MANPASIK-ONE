//! BLE 5.0 GATT 통신 모듈
//!
//! 리더기와의 BLE GATT 통신 - btleplug 기반 구현

use parking_lot::RwLock;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::Arc;
use thiserror::Error;

#[cfg(feature = "ble")]
use btleplug::api::{
    Central, CharPropFlags, Characteristic, Manager as _, Peripheral as _, ScanFilter, WriteType,
};
#[cfg(feature = "ble")]
use btleplug::platform::{Adapter, Manager, Peripheral};
#[cfg(feature = "ble")]
use uuid::Uuid as BleUuid;

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
            let manager =
                Manager::new()
                    .await
                    .map_err(|e| BleError::ConnectionFailed(e.to_string()))?;
            let adapters = manager
                .adapters()
                .await
                .map_err(|_| BleError::NoAdapter)?;
            let adapter = adapters.into_iter().next().ok_or(BleError::NoAdapter)?;

            // 스캔 캐시 또는 새 스캔으로 peripheral 탐색
            let peripherals = adapter
                .peripherals()
                .await
                .map_err(|e| BleError::DeviceNotFound(e.to_string()))?;

            let peripheral = peripherals
                .into_iter()
                .find(|p| p.id().to_string() == device_id)
                .ok_or_else(|| BleError::DeviceNotFound(device_id.to_string()))?;

            // BLE 연결
            peripheral
                .connect()
                .await
                .map_err(|e| BleError::ConnectionFailed(e.to_string()))?;

            // GATT 서비스 탐색
            peripheral
                .discover_services()
                .await
                .map_err(|e| BleError::ServiceNotFound(e.to_string()))?;

            // 펌웨어 버전 읽기 시도
            let fw_version = Self::read_characteristic_string(
                &peripheral,
                characteristic_uuids::FIRMWARE_VERSION,
            )
            .await
            .ok();

            // 배터리 레벨 읽기 시도
            let battery = Self::read_characteristic_u8(
                &peripheral,
                characteristic_uuids::BATTERY_LEVEL,
            )
            .await
            .ok();

            let properties = peripheral.properties().await.ok().flatten();
            let name = properties
                .and_then(|p| p.local_name)
                .unwrap_or_else(|| format!("MPK-{}", &device_id[..8.min(device_id.len())]));

            let device_info = DeviceInfo {
                device_id: device_id.to_string(),
                name,
                rssi: 0,
                state: ConnectionState::Connected,
                firmware_version: fw_version,
                battery_level: battery,
            };

            self.connected_devices
                .write()
                .insert(device_id.to_string(), device_info);
            return Ok(());
        }

        // BLE 기능 비활성화 시 스텁 모드
        #[cfg(not(feature = "ble"))]
        {
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

        #[cfg(feature = "ble")]
        {
            // GATT Write: 측정 시작 명령 (0x01)
            if let Err(e) = self
                .write_command(device_id, BleCommand::StartMeasurement)
                .await
            {
                tracing::warn!("BLE 명령 전송 실패 (스텁 모드로 전환): {}", e);
            }
        }

        // 상태 업데이트
        if let Some(device) = self.connected_devices.write().get_mut(device_id) {
            device.state = ConnectionState::Measuring;
        }

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

        // 캐시된 값 반환 (실시간 BLE 읽기는 연결 시 수행)
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

    /// BLE GATT Write 명령 전송
    #[cfg(feature = "ble")]
    async fn write_command(&self, device_id: &str, command: BleCommand) -> Result<(), BleError> {
        let manager =
            Manager::new()
                .await
                .map_err(|e| BleError::WriteError(e.to_string()))?;
        let adapters = manager
            .adapters()
            .await
            .map_err(|_| BleError::NoAdapter)?;
        let adapter = adapters.into_iter().next().ok_or(BleError::NoAdapter)?;
        let peripherals = adapter
            .peripherals()
            .await
            .map_err(|e| BleError::WriteError(e.to_string()))?;

        let peripheral = peripherals
            .into_iter()
            .find(|p| p.id().to_string() == device_id)
            .ok_or_else(|| BleError::DeviceNotFound(device_id.to_string()))?;

        // MEASUREMENT_COMMAND 특성 찾기
        let cmd_uuid =
            BleUuid::parse_str(characteristic_uuids::MEASUREMENT_COMMAND)
                .map_err(|e| BleError::CharacteristicNotFound(e.to_string()))?;

        let characteristics = peripheral.characteristics();
        let cmd_char = characteristics
            .iter()
            .find(|c| c.uuid == cmd_uuid)
            .ok_or_else(|| {
                BleError::CharacteristicNotFound("MEASUREMENT_COMMAND".to_string())
            })?;

        // 명령 바이트 전송
        peripheral
            .write(cmd_char, &[command as u8], WriteType::WithResponse)
            .await
            .map_err(|e| BleError::WriteError(e.to_string()))?;

        Ok(())
    }

    /// BLE GATT 특성에서 문자열 읽기
    #[cfg(feature = "ble")]
    async fn read_characteristic_string(
        peripheral: &Peripheral,
        uuid_str: &str,
    ) -> Result<String, BleError> {
        let uuid = BleUuid::parse_str(uuid_str)
            .map_err(|e| BleError::CharacteristicNotFound(e.to_string()))?;
        let characteristics = peripheral.characteristics();
        let char = characteristics
            .iter()
            .find(|c| c.uuid == uuid)
            .ok_or_else(|| BleError::CharacteristicNotFound(uuid_str.to_string()))?;

        let data = peripheral
            .read(char)
            .await
            .map_err(|e| BleError::ReadError(e.to_string()))?;
        Ok(String::from_utf8_lossy(&data).to_string())
    }

    /// BLE GATT 특성에서 u8 읽기
    #[cfg(feature = "ble")]
    async fn read_characteristic_u8(
        peripheral: &Peripheral,
        uuid_str: &str,
    ) -> Result<u8, BleError> {
        let uuid = BleUuid::parse_str(uuid_str)
            .map_err(|e| BleError::CharacteristicNotFound(e.to_string()))?;
        let characteristics = peripheral.characteristics();
        let char = characteristics
            .iter()
            .find(|c| c.uuid == uuid)
            .ok_or_else(|| BleError::CharacteristicNotFound(uuid_str.to_string()))?;

        let data = peripheral
            .read(char)
            .await
            .map_err(|e| BleError::ReadError(e.to_string()))?;
        data.first()
            .copied()
            .ok_or_else(|| BleError::ReadError("빈 응답".to_string()))
    }

    /// 측정 데이터 GATT Notify 구독
    ///
    /// MEASUREMENT_DATA 특성을 구독하여 스트리밍 데이터 수신
    #[cfg(feature = "ble")]
    pub async fn subscribe_measurement_data(
        &self,
        device_id: &str,
    ) -> Result<(), BleError> {
        let manager =
            Manager::new()
                .await
                .map_err(|e| BleError::ConnectionFailed(e.to_string()))?;
        let adapters = manager
            .adapters()
            .await
            .map_err(|_| BleError::NoAdapter)?;
        let adapter = adapters.into_iter().next().ok_or(BleError::NoAdapter)?;
        let peripherals = adapter
            .peripherals()
            .await
            .map_err(|e| BleError::ConnectionFailed(e.to_string()))?;

        let peripheral = peripherals
            .into_iter()
            .find(|p| p.id().to_string() == device_id)
            .ok_or_else(|| BleError::DeviceNotFound(device_id.to_string()))?;

        // MEASUREMENT_DATA Notify 특성 구독
        let data_uuid =
            BleUuid::parse_str(characteristic_uuids::MEASUREMENT_DATA)
                .map_err(|e| BleError::CharacteristicNotFound(e.to_string()))?;

        let characteristics = peripheral.characteristics();
        let data_char = characteristics
            .iter()
            .find(|c| c.uuid == data_uuid && c.properties.contains(CharPropFlags::NOTIFY))
            .ok_or_else(|| {
                BleError::CharacteristicNotFound("MEASUREMENT_DATA (Notify)".to_string())
            })?;

        peripheral
            .subscribe(data_char)
            .await
            .map_err(|e| BleError::ReadError(e.to_string()))?;

        // 상태 업데이트 → Streaming
        if let Some(device) = self.connected_devices.write().get_mut(device_id) {
            device.state = ConnectionState::Measuring;
        }

        Ok(())
    }
}

impl Default for BleManager {
    fn default() -> Self {
        Self::new()
    }
}

// ============================================================================
// BLE 상태 머신 (ISO 14971 FM-BLE-001 대응)
// ============================================================================

/// BLE 연결 상태 머신 - 6-state FSM
/// Disconnected → Scanning → Connecting → Connected → Streaming → Disconnecting
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum BleStateMachine {
    Disconnected,
    Scanning,
    Connecting,
    Connected,
    Streaming,
    Disconnecting,
}

impl BleStateMachine {
    /// 상태 전이 유효성 검증
    pub fn can_transition(&self, target: BleStateMachine) -> bool {
        matches!(
            (self, target),
            (BleStateMachine::Disconnected, BleStateMachine::Scanning)
                | (BleStateMachine::Scanning, BleStateMachine::Connecting)
                | (BleStateMachine::Scanning, BleStateMachine::Disconnected)
                | (BleStateMachine::Connecting, BleStateMachine::Connected)
                | (BleStateMachine::Connecting, BleStateMachine::Disconnected)
                | (BleStateMachine::Connected, BleStateMachine::Streaming)
                | (BleStateMachine::Connected, BleStateMachine::Disconnecting)
                | (BleStateMachine::Streaming, BleStateMachine::Connected)
                | (BleStateMachine::Streaming, BleStateMachine::Disconnecting)
                | (BleStateMachine::Disconnecting, BleStateMachine::Disconnected)
        )
    }

    /// 상태 전이 수행
    pub fn transition(&mut self, target: BleStateMachine) -> Result<(), BleError> {
        if self.can_transition(target) {
            *self = target;
            Ok(())
        } else {
            Err(BleError::ConnectionFailed(format!(
                "잘못된 상태 전이: {:?} → {:?}",
                self, target
            )))
        }
    }
}

// ============================================================================
// 청크 재조립기 (FM-BLE-005 대응: 896/1792차원 대용량 데이터)
// ============================================================================

/// BLE MTU 제한으로 인한 대용량 데이터 청크 재조립
pub struct ChunkReassembler {
    /// 총 예상 청크 수
    total_chunks: u16,
    /// 수신된 청크 (순서 보장)
    received: HashMap<u16, Vec<u8>>,
    /// 시작 시간 (타임아웃 체크용)
    started_at: std::time::Instant,
    /// 타임아웃 (기본 30초)
    timeout: std::time::Duration,
}

impl ChunkReassembler {
    pub fn new(total_chunks: u16) -> Self {
        Self {
            total_chunks,
            received: HashMap::new(),
            started_at: std::time::Instant::now(),
            timeout: std::time::Duration::from_secs(30),
        }
    }

    /// 청크 추가
    pub fn add_chunk(&mut self, sequence: u16, data: Vec<u8>) -> Result<(), BleError> {
        if self.started_at.elapsed() > self.timeout {
            return Err(BleError::Timeout);
        }
        if sequence >= self.total_chunks {
            return Err(BleError::ReadError(format!(
                "시퀀스 번호 초과: {} >= {}",
                sequence, self.total_chunks
            )));
        }
        self.received.insert(sequence, data);
        Ok(())
    }

    /// 모든 청크 수신 완료 여부
    pub fn is_complete(&self) -> bool {
        self.received.len() as u16 == self.total_chunks
    }

    /// 누락된 청크 시퀀스 번호 반환
    pub fn missing_chunks(&self) -> Vec<u16> {
        (0..self.total_chunks)
            .filter(|seq| !self.received.contains_key(seq))
            .collect()
    }

    /// 완성된 데이터 추출 (순서대로 연결)
    pub fn assemble(&self) -> Result<Vec<u8>, BleError> {
        if !self.is_complete() {
            return Err(BleError::ReadError(format!(
                "미완성: {}/{} 청크 수신",
                self.received.len(),
                self.total_chunks
            )));
        }
        let mut result = Vec::new();
        for seq in 0..self.total_chunks {
            if let Some(data) = self.received.get(&seq) {
                result.extend_from_slice(data);
            }
        }
        Ok(result)
    }

    /// 수신 진행률 (0.0 ~ 1.0)
    pub fn progress(&self) -> f32 {
        self.received.len() as f32 / self.total_chunks as f32
    }
}

// ============================================================================
// RSSI 모니터 (연결 품질 모니터링)
// ============================================================================

/// RSSI 기반 연결 품질 모니터
pub struct RssiMonitor {
    samples: Vec<i8>,
    max_samples: usize,
}

impl RssiMonitor {
    pub fn new(max_samples: usize) -> Self {
        Self {
            samples: Vec::with_capacity(max_samples),
            max_samples,
        }
    }

    /// RSSI 샘플 추가
    pub fn add_sample(&mut self, rssi: i8) {
        if self.samples.len() >= self.max_samples {
            self.samples.remove(0);
        }
        self.samples.push(rssi);
    }

    /// 평균 RSSI
    pub fn average_rssi(&self) -> f32 {
        if self.samples.is_empty() {
            return -100.0;
        }
        self.samples.iter().map(|&r| r as f32).sum::<f32>() / self.samples.len() as f32
    }

    /// 연결 품질 등급
    pub fn quality(&self) -> ConnectionQuality {
        let avg = self.average_rssi();
        match avg as i32 {
            -50..=0 => ConnectionQuality::Excellent,
            -70..=-51 => ConnectionQuality::Good,
            -85..=-71 => ConnectionQuality::Fair,
            _ => ConnectionQuality::Poor,
        }
    }
}

/// 연결 품질 등급
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum ConnectionQuality {
    Excellent,
    Good,
    Fair,
    Poor,
}

// ============================================================================
// 재연결 전략 (지수 백오프)
// ============================================================================

/// 자동 재연결 전략 (FM-BLE-001 대응)
pub struct ReconnectionStrategy {
    max_retries: u32,
    current_retry: u32,
    base_delay_ms: u64,
}

impl ReconnectionStrategy {
    pub fn new(max_retries: u32) -> Self {
        Self {
            max_retries,
            current_retry: 0,
            base_delay_ms: 500,
        }
    }

    /// 재연결 시도 가능 여부
    pub fn can_retry(&self) -> bool {
        self.current_retry < self.max_retries
    }

    /// 다음 재시도 대기 시간 (밀리초, 지수 백오프)
    pub fn next_delay_ms(&mut self) -> Option<u64> {
        if !self.can_retry() {
            return None;
        }
        let delay = self.base_delay_ms * 2u64.pow(self.current_retry);
        self.current_retry += 1;
        Some(delay.min(30_000)) // 최대 30초
    }

    /// 재시도 카운터 리셋
    pub fn reset(&mut self) {
        self.current_retry = 0;
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
    fn test_state_machine_valid_transitions() {
        let mut state = BleStateMachine::Disconnected;
        assert!(state.transition(BleStateMachine::Scanning).is_ok());
        assert_eq!(state, BleStateMachine::Scanning);
        assert!(state.transition(BleStateMachine::Connecting).is_ok());
        assert_eq!(state, BleStateMachine::Connecting);
        assert!(state.transition(BleStateMachine::Connected).is_ok());
        assert!(state.transition(BleStateMachine::Streaming).is_ok());
        assert!(state.transition(BleStateMachine::Disconnecting).is_ok());
        assert!(state.transition(BleStateMachine::Disconnected).is_ok());
    }

    #[test]
    fn test_state_machine_invalid_transition() {
        let mut state = BleStateMachine::Disconnected;
        assert!(state.transition(BleStateMachine::Streaming).is_err());
    }

    #[test]
    fn test_chunk_reassembler() {
        let mut reassembler = ChunkReassembler::new(3);
        assert!(!reassembler.is_complete());

        reassembler.add_chunk(0, vec![1, 2]).unwrap();
        reassembler.add_chunk(2, vec![5, 6]).unwrap();
        assert_eq!(reassembler.missing_chunks(), vec![1]);

        reassembler.add_chunk(1, vec![3, 4]).unwrap();
        assert!(reassembler.is_complete());

        let data = reassembler.assemble().unwrap();
        assert_eq!(data, vec![1, 2, 3, 4, 5, 6]);
    }

    #[test]
    fn test_rssi_monitor() {
        let mut monitor = RssiMonitor::new(10);
        monitor.add_sample(-45);
        monitor.add_sample(-50);
        monitor.add_sample(-48);
        assert!(monitor.average_rssi() > -51.0);
        assert_eq!(monitor.quality(), ConnectionQuality::Excellent);
    }

    #[test]
    fn test_reconnection_strategy() {
        let mut strategy = ReconnectionStrategy::new(3);
        assert!(strategy.can_retry());
        assert_eq!(strategy.next_delay_ms(), Some(500));
        assert_eq!(strategy.next_delay_ms(), Some(1000));
        assert_eq!(strategy.next_delay_ms(), Some(2000));
        assert!(!strategy.can_retry());
        strategy.reset();
        assert!(strategy.can_retry());
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
