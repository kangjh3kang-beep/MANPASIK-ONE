// harness/sensor_trait.rs

use crate::harness::cartridge_manifest::CartridgeManifest;

pub struct CalibrationData;
pub struct SensorHealth;

/// 통합 AFE 블록의 하네스 엔지니어링 추상화 레이어 트레이트
pub trait SensorTrait {
    type Config;
    type RawData;
    type ProcessedData;
    type Error;

    /// 센서 초기화 (카트리지 매니페스트 기반)
    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error>;

    /// 원시 데이터 읽기 (DMA 기반 고속 수집)
    fn read_raw(&self) -> Result<Self::RawData, Self::Error>;

    /// 차동측정 적용: S_det - alpha * S_ref
    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error>;

    /// 보정 데이터 로드 (NFC 태그 → 폴백 체인)
    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error>;

    /// 센서 상태 진단 (디지털 트윈 잔차 기반)
    fn self_diagnose(&self) -> SensorHealth;

    /// 안전 종료
    fn shutdown(&mut self) -> Result<(), Self::Error>;
}
