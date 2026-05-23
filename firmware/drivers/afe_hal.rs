//! AFE 9블록 HAL 추상화 — embedded-hal trait 기반
//!
//! H1: 구체 IC(LMP91000, ADS1256 등)를 직접 참조하지 않음
//! H4: Stage-1(전기화학) → Stage-2(광학) → Stage-3(NAAT) 순차 활성화

#![no_std]

/// AFE 블록 식별자 (9블록 Universal AFE)
#[derive(Debug, Clone, Copy, PartialEq)]
#[repr(u8)]
pub enum AfeBlockId {
    Electrochem = 1,    // Stage-1: 전기화학 (LMP91000×4 + ADS1256)
    Enose = 2,          // Stage-1: 8채널 전자코
    EcArray = 3,        // Stage-1: EC 어레이 확장
    SipmEcl = 4,        // Stage-2: SiPM-ECL
    Colorimetric = 5,   // Stage-2: LED/PD 색도
    Tec = 6,            // Stage-2: 열전 냉각
    Lamp = 7,           // Stage-3: LAMP/RPA 등온 증폭
    Reserved8 = 8,      // 예약
    Reserved9 = 9,      // 예약
}

/// 측정 Stage
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum MeasurementStage {
    Stage1Electrochem,  // 전기화학 (CSI v1.0 기본)
    Stage2Optical,      // + 광학 (핀 15-16 활성화)
    Stage3Naat,         // + NAAT (CSI v2.0+)
}

/// AFE 블록 HAL trait (H1: 구체 구현 격리)
pub trait AfeBlock {
    type Error;
    type RawReading;

    /// 블록 초기화
    fn init(&mut self) -> Result<(), Self::Error>;

    /// 원시 읽기 (채널별)
    fn read_channel(&self, channel: u8) -> Result<Self::RawReading, Self::Error>;

    /// 전체 채널 읽기
    fn read_all_channels(&self) -> Result<[Self::RawReading; 16], Self::Error>;

    /// 자가진단 (H6: 실패 시 해당 블록만 비활성화)
    fn self_test(&self) -> Result<bool, Self::Error>;

    /// 블록 셧다운
    fn shutdown(&mut self) -> Result<(), Self::Error>;

    /// 블록 ID
    fn block_id(&self) -> AfeBlockId;
}

/// RAFE 스위치 제어 (Stage별 블록 라우팅)
pub trait RafeSwitch {
    type Error;

    /// 측정 Stage에 따라 AFE 블록 활성화
    fn activate_stage(&mut self, stage: MeasurementStage) -> Result<(), Self::Error>;

    /// 특정 블록 활성화/비활성화
    fn set_block_enabled(&mut self, block: AfeBlockId, enabled: bool) -> Result<(), Self::Error>;

    /// 현재 활성 블록 목록
    fn active_blocks(&self) -> &[AfeBlockId];
}

/// Stage별 필요 AFE 블록 매핑
pub fn blocks_for_stage(stage: MeasurementStage) -> &'static [AfeBlockId] {
    match stage {
        MeasurementStage::Stage1Electrochem => &[
            AfeBlockId::Electrochem,
            AfeBlockId::Enose,
            AfeBlockId::EcArray,
        ],
        MeasurementStage::Stage2Optical => &[
            AfeBlockId::Electrochem,
            AfeBlockId::Enose,
            AfeBlockId::EcArray,
            AfeBlockId::SipmEcl,
            AfeBlockId::Colorimetric,
            AfeBlockId::Tec,
        ],
        MeasurementStage::Stage3Naat => &[
            AfeBlockId::Electrochem,
            AfeBlockId::Enose,
            AfeBlockId::EcArray,
            AfeBlockId::SipmEcl,
            AfeBlockId::Colorimetric,
            AfeBlockId::Tec,
            AfeBlockId::Lamp,
        ],
    }
}

/// CSI v1.0 핀 할당 (14 신호핀 + 2 마운트)
pub const CSI_V1_SIGNAL_PINS: u8 = 14;
pub const CSI_V1_TOTAL_PINS: u8 = 16;
pub const CSI_V1_PITCH_MM: f32 = 1.27;
