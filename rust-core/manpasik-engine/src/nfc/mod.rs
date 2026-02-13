//! NFC 카트리지 인식 모듈
//!
//! ISO 14443A 기반 카트리지 인식 및 보정 데이터 관리
//! v2.0: 무한확장 레지스트리 구조 (2-byte 계층형 코드)

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::OnceLock;
use thiserror::Error;

// ============================================================================
// 카트리지 카테고리 (무한확장 기반)
// ============================================================================

/// 카트리지 카테고리 코드 (0x01~0xFF)
/// 256개 카테고리 × 256개 타입/카테고리 = 65,536종 수용
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq, Hash)]
#[repr(u8)]
pub enum CartridgeCategory {
    /// 건강 바이오마커 (혈액/타액/체액)
    HealthBiomarker = 0x01,
    /// 환경 모니터링
    Environmental = 0x02,
    /// 식품 안전
    FoodSafety = 0x03,
    /// 전자코/전자혀 센서
    ElectronicSensor = 0x04,
    /// 고급 분석 (비표적/다중패널)
    AdvancedAnalysis = 0x05,
    /// 산업용 분석
    Industrial = 0x06,
    /// 수의학
    Veterinary = 0x07,
    /// 제약
    Pharmaceutical = 0x08,
    /// 농업
    Agricultural = 0x09,
    /// 화장품
    Cosmetic = 0x0A,
    /// 법의학
    Forensic = 0x0B,
    /// 해양
    Marine = 0x0C,
    /// 베타/실험용
    Beta = 0xFE,
    /// 맞춤형 연구용
    CustomResearch = 0xFF,
    /// 알 수 없는 카테고리
    Unknown = 0x00,
}

impl CartridgeCategory {
    /// 코드에서 카테고리 파싱
    pub fn from_code(code: u8) -> Self {
        match code {
            0x01 => CartridgeCategory::HealthBiomarker,
            0x02 => CartridgeCategory::Environmental,
            0x03 => CartridgeCategory::FoodSafety,
            0x04 => CartridgeCategory::ElectronicSensor,
            0x05 => CartridgeCategory::AdvancedAnalysis,
            0x06 => CartridgeCategory::Industrial,
            0x07 => CartridgeCategory::Veterinary,
            0x08 => CartridgeCategory::Pharmaceutical,
            0x09 => CartridgeCategory::Agricultural,
            0x0A => CartridgeCategory::Cosmetic,
            0x0B => CartridgeCategory::Forensic,
            0x0C => CartridgeCategory::Marine,
            0xF0..=0xFD => CartridgeCategory::Unknown, // ThirdParty → 런타임 레지스트리
            0xFE => CartridgeCategory::Beta,
            0xFF => CartridgeCategory::CustomResearch,
            _ => CartridgeCategory::Unknown,
        }
    }

    /// 카테고리를 코드로 변환
    pub fn to_code(&self) -> u8 {
        *self as u8
    }

    /// 카테고리 한국어 이름
    pub fn name_ko(&self) -> &'static str {
        match self {
            CartridgeCategory::HealthBiomarker => "건강 바이오마커",
            CartridgeCategory::Environmental => "환경 모니터링",
            CartridgeCategory::FoodSafety => "식품 안전",
            CartridgeCategory::ElectronicSensor => "전자코/전자혀",
            CartridgeCategory::AdvancedAnalysis => "고급 분석",
            CartridgeCategory::Industrial => "산업용",
            CartridgeCategory::Veterinary => "수의학",
            CartridgeCategory::Pharmaceutical => "제약",
            CartridgeCategory::Agricultural => "농업",
            CartridgeCategory::Cosmetic => "화장품",
            CartridgeCategory::Forensic => "법의학",
            CartridgeCategory::Marine => "해양",
            CartridgeCategory::Beta => "베타/실험용",
            CartridgeCategory::CustomResearch => "맞춤형 연구",
            CartridgeCategory::Unknown => "알 수 없음",
        }
    }
}

// ============================================================================
// 카트리지 풀 코드 (2-Byte 계층형)
// ============================================================================

/// 카트리지 2-byte 풀 코드: (category_code, type_index)
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct CartridgeFullCode {
    /// 카테고리 코드 (0x01~0xFF)
    pub category: u8,
    /// 타입 인덱스 (카테고리 내 순번, 0x01~0xFF)
    pub type_index: u8,
}

impl CartridgeFullCode {
    pub fn new(category: u8, type_index: u8) -> Self {
        Self {
            category,
            type_index,
        }
    }

    /// 2-byte 값으로 변환 (상위=category, 하위=type_index)
    pub fn to_u16(&self) -> u16 {
        (self.category as u16) << 8 | self.type_index as u16
    }

    /// 2-byte 값에서 파싱
    pub fn from_u16(code: u16) -> Self {
        Self {
            category: (code >> 8) as u8,
            type_index: (code & 0xFF) as u8,
        }
    }

    /// 레거시 1-byte 코드에서 변환
    pub fn from_legacy(legacy_code: u8) -> Self {
        legacy_to_full_code(legacy_code)
    }
}

// ============================================================================
// 카트리지 정보 (v2.0 확장)
// ============================================================================

/// 카트리지 정보
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CartridgeInfo {
    /// 카트리지 고유 ID (UID from NFC tag)
    pub cartridge_id: String,
    /// 카트리지 타입 (레거시 호환)
    pub cartridge_type: CartridgeType,
    /// 카트리지 풀 코드 (v2.0 무한확장)
    pub full_code: CartridgeFullCode,
    /// 카테고리
    pub category: CartridgeCategory,
    /// 태그 포맷 버전 (1=v1.0, 2=v2.0)
    pub tag_version: u8,
    /// 제조 로트 ID
    pub lot_id: String,
    /// 유효 기간 (YYYY-MM-DD)
    pub expiry_date: String,
    /// 잔여 사용 횟수
    pub remaining_uses: u32,
    /// 최대 사용 횟수
    pub max_uses: u32,
    /// 팩토리 보정 데이터
    pub calibration_data: Vec<u8>,
    /// 보정 계수 (alpha 등)
    pub calibration_coefficients: CalibrationCoefficients,
}

/// 보정 계수
#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct CalibrationCoefficients {
    /// 알파 계수 (차동측정용)
    pub alpha: f64,
    /// 채널별 오프셋
    pub offsets: Vec<f64>,
    /// 채널별 게인
    pub gains: Vec<f64>,
    /// 온도 보정 계수
    pub temp_coefficient: f64,
    /// 습도 보정 계수
    pub humidity_coefficient: f64,
}

// ============================================================================
// 레거시 CartridgeType enum (하위호환 유지)
// ============================================================================

/// 카트리지 타입 (레거시 29종 + 동적 확장)
/// 기존 코드 하위호환을 위해 enum 유지. 신규 타입은 CartridgeFullCode + 레지스트리 사용.
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub enum CartridgeType {
    // === 건강 바이오마커 (혈액/타액) ===
    /// 혈당 (Glucose)
    Glucose,
    /// 지질 패널 (총콜레스테롤, LDL, HDL, 중성지방)
    LipidPanel,
    /// HbA1c (당화혈색소)
    HbA1c,
    /// 요산 (Uric Acid)
    UricAcid,
    /// 크레아티닌 (Creatinine)
    Creatinine,
    /// 비타민 D
    VitaminD,
    /// 비타민 B12
    VitaminB12,
    /// 철분 (Ferritin)
    Ferritin,
    /// 갑상선 (TSH)
    Tsh,
    /// 코르티솔 (스트레스 호르몬)
    Cortisol,
    /// 테스토스테론
    Testosterone,
    /// 에스트로겐
    Estrogen,
    /// CRP (염증 마커)
    Crp,
    /// 인슐린
    Insulin,

    // === 환경 모니터링 ===
    /// 수질 검사 (pH, 잔류염소, 중금속)
    WaterQuality,
    /// 실내 공기질 (PM2.5, VOC, CO2)
    IndoorAirQuality,
    /// 라돈 측정
    Radon,
    /// 방사능 측정
    Radiation,

    // === 식품 안전 ===
    /// 농약 잔류
    PesticideResidue,
    /// 식품 신선도
    FoodFreshness,
    /// 알레르겐 검출 (글루텐, 땅콩 등)
    Allergen,
    /// 데이트레이프약물 검출
    DateDrug,

    // === 전자코/전자혀 ===
    /// 전자코 (E-Nose) 8채널
    ENose,
    /// 전자혀 (E-Tongue) 8채널
    ETongue,
    /// EHD 기체 분석
    EhdGas,

    // === 고급 분석 ===
    /// 비표적 448차원 분석
    NonTarget448,
    /// 비표적 896차원 분석
    NonTarget896,
    /// 비표적 1792차원 궁극 분석 (Phase 5: E12-IF + 시간축 확장)
    NonTarget1792,
    /// 다중 바이오마커 패널
    MultiBiomarker,
    /// 맞춤형 연구용
    CustomResearch,

    /// 알 수 없는 타입
    Unknown,
}

impl CartridgeType {
    /// 카트리지 이름 (한국어)
    pub fn name_ko(&self) -> &'static str {
        match self {
            CartridgeType::Glucose => "혈당",
            CartridgeType::LipidPanel => "지질 패널",
            CartridgeType::HbA1c => "당화혈색소",
            CartridgeType::UricAcid => "요산",
            CartridgeType::Creatinine => "크레아티닌",
            CartridgeType::VitaminD => "비타민 D",
            CartridgeType::VitaminB12 => "비타민 B12",
            CartridgeType::Ferritin => "철분(페리틴)",
            CartridgeType::Tsh => "갑상선(TSH)",
            CartridgeType::Cortisol => "코르티솔",
            CartridgeType::Testosterone => "테스토스테론",
            CartridgeType::Estrogen => "에스트로겐",
            CartridgeType::Crp => "C-반응성단백",
            CartridgeType::Insulin => "인슐린",
            CartridgeType::WaterQuality => "수질 검사",
            CartridgeType::IndoorAirQuality => "실내 공기질",
            CartridgeType::Radon => "라돈",
            CartridgeType::Radiation => "방사능",
            CartridgeType::PesticideResidue => "농약 잔류",
            CartridgeType::FoodFreshness => "식품 신선도",
            CartridgeType::Allergen => "알레르겐",
            CartridgeType::DateDrug => "데이트약물",
            CartridgeType::ENose => "전자코",
            CartridgeType::ETongue => "전자혀",
            CartridgeType::EhdGas => "EHD 기체",
            CartridgeType::NonTarget448 => "비표적 448차원",
            CartridgeType::NonTarget896 => "비표적 896차원",
            CartridgeType::NonTarget1792 => "비표적 1792차원(궁극)",
            CartridgeType::MultiBiomarker => "다중 바이오마커",
            CartridgeType::CustomResearch => "맞춤형 연구용",
            CartridgeType::Unknown => "알 수 없음",
        }
    }

    /// 필요한 채널 수
    pub fn required_channels(&self) -> usize {
        match self {
            CartridgeType::NonTarget1792 => 1792,
            CartridgeType::NonTarget896 => 896,
            CartridgeType::NonTarget448 => 448,
            CartridgeType::ENose | CartridgeType::ETongue => 8,
            CartridgeType::MultiBiomarker => 88,
            _ => 88,
        }
    }

    /// 측정 시간 (초)
    pub fn measurement_duration_secs(&self) -> u32 {
        match self {
            CartridgeType::NonTarget1792 => 180, // 2× 시간 윈도우
            CartridgeType::NonTarget896 => 90,
            CartridgeType::NonTarget448 => 60,
            CartridgeType::MultiBiomarker => 45,
            CartridgeType::ENose | CartridgeType::ETongue => 30,
            _ => 15,
        }
    }

    /// 코드에서 타입 파싱
    pub fn from_code(code: u8) -> Self {
        match code {
            0x01 => CartridgeType::Glucose,
            0x02 => CartridgeType::LipidPanel,
            0x03 => CartridgeType::HbA1c,
            0x04 => CartridgeType::UricAcid,
            0x05 => CartridgeType::Creatinine,
            0x06 => CartridgeType::VitaminD,
            0x07 => CartridgeType::VitaminB12,
            0x08 => CartridgeType::Ferritin,
            0x09 => CartridgeType::Tsh,
            0x0A => CartridgeType::Cortisol,
            0x0B => CartridgeType::Testosterone,
            0x0C => CartridgeType::Estrogen,
            0x0D => CartridgeType::Crp,
            0x0E => CartridgeType::Insulin,
            0x20 => CartridgeType::WaterQuality,
            0x21 => CartridgeType::IndoorAirQuality,
            0x22 => CartridgeType::Radon,
            0x23 => CartridgeType::Radiation,
            0x30 => CartridgeType::PesticideResidue,
            0x31 => CartridgeType::FoodFreshness,
            0x32 => CartridgeType::Allergen,
            0x33 => CartridgeType::DateDrug,
            0x40 => CartridgeType::ENose,
            0x41 => CartridgeType::ETongue,
            0x42 => CartridgeType::EhdGas,
            0x50 => CartridgeType::NonTarget448,
            0x51 => CartridgeType::NonTarget896,
            0x52 => CartridgeType::NonTarget1792,
            0x53 => CartridgeType::MultiBiomarker,
            0xFF => CartridgeType::CustomResearch,
            _ => CartridgeType::Unknown,
        }
    }

    /// 타입을 코드로 변환
    pub fn to_code(&self) -> u8 {
        match self {
            CartridgeType::Glucose => 0x01,
            CartridgeType::LipidPanel => 0x02,
            CartridgeType::HbA1c => 0x03,
            CartridgeType::UricAcid => 0x04,
            CartridgeType::Creatinine => 0x05,
            CartridgeType::VitaminD => 0x06,
            CartridgeType::VitaminB12 => 0x07,
            CartridgeType::Ferritin => 0x08,
            CartridgeType::Tsh => 0x09,
            CartridgeType::Cortisol => 0x0A,
            CartridgeType::Testosterone => 0x0B,
            CartridgeType::Estrogen => 0x0C,
            CartridgeType::Crp => 0x0D,
            CartridgeType::Insulin => 0x0E,
            CartridgeType::WaterQuality => 0x20,
            CartridgeType::IndoorAirQuality => 0x21,
            CartridgeType::Radon => 0x22,
            CartridgeType::Radiation => 0x23,
            CartridgeType::PesticideResidue => 0x30,
            CartridgeType::FoodFreshness => 0x31,
            CartridgeType::Allergen => 0x32,
            CartridgeType::DateDrug => 0x33,
            CartridgeType::ENose => 0x40,
            CartridgeType::ETongue => 0x41,
            CartridgeType::EhdGas => 0x42,
            CartridgeType::NonTarget448 => 0x50,
            CartridgeType::NonTarget896 => 0x51,
            CartridgeType::NonTarget1792 => 0x52,
            CartridgeType::MultiBiomarker => 0x53,
            CartridgeType::CustomResearch => 0xFF,
            CartridgeType::Unknown => 0x00,
        }
    }

    /// 레거시 타입을 풀 코드로 변환
    pub fn to_full_code(&self) -> CartridgeFullCode {
        legacy_to_full_code(self.to_code())
    }

    /// 카테고리 반환
    pub fn category(&self) -> CartridgeCategory {
        CartridgeCategory::from_code(self.to_full_code().category)
    }
}

// ============================================================================
// 레거시 → 풀 코드 변환 매핑
// ============================================================================

/// 레거시 1-byte 코드 → 2-byte 풀 코드 변환
fn legacy_to_full_code(legacy: u8) -> CartridgeFullCode {
    match legacy {
        // HealthBiomarker (0x01~0x0E → 0x01:0x01~0x01:0x0E)
        0x01..=0x0E => CartridgeFullCode::new(0x01, legacy),
        // Environmental (0x20~0x23 → 0x02:0x01~0x02:0x04)
        0x20..=0x23 => CartridgeFullCode::new(0x02, legacy - 0x1F),
        // FoodSafety (0x30~0x33 → 0x03:0x01~0x03:0x04)
        0x30..=0x33 => CartridgeFullCode::new(0x03, legacy - 0x2F),
        // ElectronicSensor (0x40~0x42 → 0x04:0x01~0x04:0x03)
        0x40..=0x42 => CartridgeFullCode::new(0x04, legacy - 0x3F),
        // AdvancedAnalysis (0x50~0x53 → 0x05:0x01~0x05:0x04)
        0x50..=0x53 => CartridgeFullCode::new(0x05, legacy - 0x4F),
        // CustomResearch (0xFF → 0xFF:0x01)
        0xFF => CartridgeFullCode::new(0xFF, 0x01),
        // Unknown
        _ => CartridgeFullCode::new(0x00, 0x00),
    }
}

// ============================================================================
// 카트리지 레지스트리 (동적 타입 관리)
// ============================================================================

/// 레지스트리 항목: 서버에서 동기화되거나 OTA로 추가된 카트리지 타입 정보
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CartridgeRegistryEntry {
    /// 풀 코드
    pub full_code: CartridgeFullCode,
    /// 한국어 이름
    pub name_ko: String,
    /// 영어 이름
    pub name_en: String,
    /// 필요 채널 수
    pub required_channels: usize,
    /// 측정 시간 (초)
    pub measurement_secs: u32,
    /// 측정 단위
    pub unit: String,
    /// 정상 범위
    pub reference_range: String,
    /// 베타 여부
    pub is_beta: bool,
    /// 제조사
    pub manufacturer: String,
}

/// 카트리지 레지스트리 (런타임 확장 가능)
/// 기본 29종은 내장, 추가 타입은 서버 동기화 또는 OTA로 등록
#[derive(Debug, Clone)]
pub struct CartridgeRegistry {
    entries: HashMap<CartridgeFullCode, CartridgeRegistryEntry>,
}

/// 기본 내장 레지스트리 (29종)
static DEFAULT_REGISTRY: OnceLock<CartridgeRegistry> = OnceLock::new();

impl CartridgeRegistry {
    /// 기본 29종으로 초기화된 레지스트리 생성
    pub fn new_with_defaults() -> Self {
        let mut entries = HashMap::new();

        // 건강 바이오마커 14종
        let health_types = [
            (0x01, "혈당", "Glucose", 88, 15, "mg/dL", "70-100"),
            (0x02, "지질 패널", "LipidPanel", 88, 15, "mg/dL", ""),
            (0x03, "당화혈색소", "HbA1c", 88, 15, "%", "4.0-5.6"),
            (0x04, "요산", "UricAcid", 88, 15, "mg/dL", "3.5-7.2"),
            (0x05, "크레아티닌", "Creatinine", 88, 15, "mg/dL", "0.7-1.3"),
            (0x06, "비타민 D", "VitaminD", 88, 15, "ng/mL", "30-100"),
            (0x07, "비타민 B12", "VitaminB12", 88, 15, "pg/mL", "200-900"),
            (0x08, "철분(페리틴)", "Ferritin", 88, 15, "ng/mL", "12-300"),
            (0x09, "갑상선(TSH)", "Tsh", 88, 15, "mIU/L", "0.4-4.0"),
            (0x0A, "코르티솔", "Cortisol", 88, 15, "μg/dL", "6-23"),
            (
                0x0B,
                "테스토스테론",
                "Testosterone",
                88,
                15,
                "ng/dL",
                "300-1000",
            ),
            (0x0C, "에스트로겐", "Estrogen", 88, 15, "pg/mL", "15-350"),
            (0x0D, "C-반응성단백", "Crp", 88, 15, "mg/L", "0-3"),
            (0x0E, "인슐린", "Insulin", 88, 15, "μIU/mL", "2.6-24.9"),
        ];
        for (idx, name_ko, name_en, ch, secs, unit, range) in health_types {
            entries.insert(
                CartridgeFullCode::new(0x01, idx),
                CartridgeRegistryEntry {
                    full_code: CartridgeFullCode::new(0x01, idx),
                    name_ko: name_ko.to_string(),
                    name_en: name_en.to_string(),
                    required_channels: ch,
                    measurement_secs: secs,
                    unit: unit.to_string(),
                    reference_range: range.to_string(),
                    is_beta: false,
                    manufacturer: "ManPaSik".to_string(),
                },
            );
        }

        // 환경 4종
        let env_types = [
            (0x01, "수질 검사", "WaterQuality", 88, 15),
            (0x02, "실내 공기질", "IndoorAirQuality", 88, 15),
            (0x03, "라돈", "Radon", 88, 15),
            (0x04, "방사능", "Radiation", 88, 15),
        ];
        for (idx, name_ko, name_en, ch, secs) in env_types {
            entries.insert(
                CartridgeFullCode::new(0x02, idx),
                CartridgeRegistryEntry {
                    full_code: CartridgeFullCode::new(0x02, idx),
                    name_ko: name_ko.to_string(),
                    name_en: name_en.to_string(),
                    required_channels: ch,
                    measurement_secs: secs,
                    unit: String::new(),
                    reference_range: String::new(),
                    is_beta: false,
                    manufacturer: "ManPaSik".to_string(),
                },
            );
        }

        // 식품 4종
        let food_types = [
            (0x01, "농약 잔류", "PesticideResidue"),
            (0x02, "식품 신선도", "FoodFreshness"),
            (0x03, "알레르겐", "Allergen"),
            (0x04, "데이트약물", "DateDrug"),
        ];
        for (idx, name_ko, name_en) in food_types {
            entries.insert(
                CartridgeFullCode::new(0x03, idx),
                CartridgeRegistryEntry {
                    full_code: CartridgeFullCode::new(0x03, idx),
                    name_ko: name_ko.to_string(),
                    name_en: name_en.to_string(),
                    required_channels: 88,
                    measurement_secs: 15,
                    unit: String::new(),
                    reference_range: String::new(),
                    is_beta: false,
                    manufacturer: "ManPaSik".to_string(),
                },
            );
        }

        // 전자코/전자혀 3종
        let sensor_types = [
            (0x01, "전자코", "ENose", 8, 30),
            (0x02, "전자혀", "ETongue", 8, 30),
            (0x03, "EHD 기체", "EhdGas", 8, 30),
        ];
        for (idx, name_ko, name_en, ch, secs) in sensor_types {
            entries.insert(
                CartridgeFullCode::new(0x04, idx),
                CartridgeRegistryEntry {
                    full_code: CartridgeFullCode::new(0x04, idx),
                    name_ko: name_ko.to_string(),
                    name_en: name_en.to_string(),
                    required_channels: ch,
                    measurement_secs: secs,
                    unit: String::new(),
                    reference_range: String::new(),
                    is_beta: false,
                    manufacturer: "ManPaSik".to_string(),
                },
            );
        }

        // 고급 분석 4종 (1792차원 궁극 확장 포함)
        let adv_types = [
            (0x01, "비표적 448차원", "NonTarget448", 448, 60),
            (0x02, "비표적 896차원", "NonTarget896", 896, 90),
            (0x03, "비표적 1792차원(궁극)", "NonTarget1792", 1792, 180),
            (0x04, "다중 바이오마커", "MultiBiomarker", 88, 45),
        ];
        for (idx, name_ko, name_en, ch, secs) in adv_types {
            entries.insert(
                CartridgeFullCode::new(0x05, idx),
                CartridgeRegistryEntry {
                    full_code: CartridgeFullCode::new(0x05, idx),
                    name_ko: name_ko.to_string(),
                    name_en: name_en.to_string(),
                    required_channels: ch,
                    measurement_secs: secs,
                    unit: String::new(),
                    reference_range: String::new(),
                    is_beta: false,
                    manufacturer: "ManPaSik".to_string(),
                },
            );
        }

        // CustomResearch 1종
        entries.insert(
            CartridgeFullCode::new(0xFF, 0x01),
            CartridgeRegistryEntry {
                full_code: CartridgeFullCode::new(0xFF, 0x01),
                name_ko: "맞춤형 연구용".to_string(),
                name_en: "CustomResearch".to_string(),
                required_channels: 896,
                measurement_secs: 90,
                unit: String::new(),
                reference_range: String::new(),
                is_beta: false,
                manufacturer: "ManPaSik".to_string(),
            },
        );

        Self { entries }
    }

    /// 글로벌 기본 레지스트리 참조
    pub fn global() -> &'static CartridgeRegistry {
        DEFAULT_REGISTRY.get_or_init(Self::new_with_defaults)
    }

    /// 풀 코드로 레지스트리 항목 조회
    pub fn get(&self, code: &CartridgeFullCode) -> Option<&CartridgeRegistryEntry> {
        self.entries.get(code)
    }

    /// 레거시 코드로 레지스트리 항목 조회
    pub fn get_by_legacy(&self, legacy_code: u8) -> Option<&CartridgeRegistryEntry> {
        let full = legacy_to_full_code(legacy_code);
        self.entries.get(&full)
    }

    /// 카테고리별 전체 타입 목록
    pub fn list_by_category(&self, category_code: u8) -> Vec<&CartridgeRegistryEntry> {
        self.entries
            .values()
            .filter(|e| e.full_code.category == category_code)
            .collect()
    }

    /// 신규 카트리지 타입 등록 (서버 동기화/OTA)
    pub fn register(&mut self, entry: CartridgeRegistryEntry) {
        self.entries.insert(entry.full_code, entry);
    }

    /// 등록된 전체 카트리지 수
    pub fn count(&self) -> usize {
        self.entries.len()
    }

    /// 풀 코드로 필요 채널 수 조회 (레지스트리 기반)
    pub fn required_channels(&self, code: &CartridgeFullCode) -> usize {
        self.entries
            .get(code)
            .map(|e| e.required_channels)
            .unwrap_or(88) // 기본값 88
    }

    /// 풀 코드로 측정 시간 조회 (레지스트리 기반)
    pub fn measurement_duration_secs(&self, code: &CartridgeFullCode) -> u32 {
        self.entries
            .get(code)
            .map(|e| e.measurement_secs)
            .unwrap_or(15) // 기본값 15초
    }
}

#[derive(Debug, Error)]
pub enum NfcError {
    #[error("NFC 리더를 찾을 수 없습니다")]
    NoReader,

    #[error("카트리지가 감지되지 않습니다")]
    NoCartridge,

    #[error("읽기 실패: {0}")]
    ReadError(String),

    #[error("쓰기 실패: {0}")]
    WriteError(String),

    #[error("유효하지 않은 카트리지: {0}")]
    InvalidCartridge(String),

    #[error("카트리지 만료됨")]
    Expired,

    #[error("사용 횟수 초과")]
    UsesExceeded,

    #[error("보정 데이터 오류")]
    CalibrationError,
}

/// NFC 리더
pub struct NfcReader {
    /// 마지막으로 읽은 카트리지 정보 캐시
    last_cartridge: Option<CartridgeInfo>,
}

impl NfcReader {
    pub fn new() -> Self {
        Self {
            last_cartridge: None,
        }
    }

    /// 카트리지 읽기
    pub async fn read_cartridge(&mut self) -> Result<CartridgeInfo, NfcError> {
        // TODO: 실제 NFC 읽기 구현 (nfc_manager 또는 플랫폼별 구현)

        // 현재는 시뮬레이션 모드로 동작
        Err(NfcError::NoCartridge)
    }

    /// NFC 태그 데이터 파싱 (v1.0 + v2.0 자동 감지)
    pub fn parse_tag_data(&self, data: &[u8]) -> Result<CartridgeInfo, NfcError> {
        if data.len() < 64 {
            return Err(NfcError::InvalidCartridge(
                "데이터가 너무 짧습니다".to_string(),
            ));
        }

        // 태그 버전 감지: v2.0은 byte[11] == 0x02
        let tag_version = if data.len() >= 80 && data[11] == 0x02 {
            2u8
        } else {
            1u8
        };

        if tag_version == 2 {
            self.parse_tag_v2(data)
        } else {
            self.parse_tag_v1(data)
        }
    }

    /// v1.0 태그 파싱 (레거시 호환)
    fn parse_tag_v1(&self, data: &[u8]) -> Result<CartridgeInfo, NfcError> {
        // v1.0 레이아웃 (64+ 바이트):
        // [0-7]: 카트리지 ID (8바이트 UID)
        // [8]: 카트리지 타입 코드 (레거시 1-byte)
        // [9-16]: 로트 ID (8바이트)
        // [17-24]: 유효 기간 (YYYYMMDD 문자열)
        // [25-26]: 잔여 사용 횟수 (u16 LE)
        // [27-28]: 최대 사용 횟수 (u16 LE)
        // [29-36]: 알파 계수 (f64 LE)
        // [37-44]: 온도 보정 계수 (f64 LE)
        // [45-52]: 습도 보정 계수 (f64 LE)
        // [53+]: 추가 보정 데이터

        let cartridge_id = hex::encode(&data[0..8]);
        let legacy_code = data[8];
        let cartridge_type = CartridgeType::from_code(legacy_code);
        let full_code = CartridgeFullCode::from_legacy(legacy_code);
        let category = CartridgeCategory::from_code(full_code.category);
        let lot_id = String::from_utf8_lossy(&data[9..17])
            .trim_end_matches('\0')
            .to_string();

        let expiry_str = String::from_utf8_lossy(&data[17..25]).to_string();
        let expiry_date = format!(
            "{}-{}-{}",
            &expiry_str[0..4],
            &expiry_str[4..6],
            &expiry_str[6..8]
        );

        let remaining_uses = u16::from_le_bytes([data[25], data[26]]) as u32;
        let max_uses = u16::from_le_bytes([data[27], data[28]]) as u32;

        let alpha = f64::from_le_bytes([
            data[29], data[30], data[31], data[32], data[33], data[34], data[35], data[36],
        ]);

        let temp_coefficient = f64::from_le_bytes([
            data[37], data[38], data[39], data[40], data[41], data[42], data[43], data[44],
        ]);

        let humidity_coefficient = f64::from_le_bytes([
            data[45], data[46], data[47], data[48], data[49], data[50], data[51], data[52],
        ]);

        let calibration_data = data[53..].to_vec();

        Ok(CartridgeInfo {
            cartridge_id,
            cartridge_type,
            full_code,
            category,
            tag_version: 1,
            lot_id,
            expiry_date,
            remaining_uses,
            max_uses,
            calibration_data,
            calibration_coefficients: CalibrationCoefficients {
                alpha,
                offsets: Vec::new(),
                gains: Vec::new(),
                temp_coefficient,
                humidity_coefficient,
            },
        })
    }

    /// v2.0 태그 파싱 (무한확장 지원)
    fn parse_tag_v2(&self, data: &[u8]) -> Result<CartridgeInfo, NfcError> {
        if data.len() < 80 {
            return Err(NfcError::InvalidCartridge(
                "v2.0 태그: 최소 80바이트 필요".to_string(),
            ));
        }

        // v2.0 레이아웃 (80+ 바이트):
        // [0-7]: 카트리지 UID
        // [8]: 카테고리 코드
        // [9]: 타입 인덱스
        // [10]: 레거시 호환 코드 (0x00이면 없음)
        // [11]: 태그 포맷 버전 (0x02)
        // [12-19]: 로트 ID
        // [20-27]: 유효 기간
        // [28-29]: 잔여 사용 횟수
        // [30-31]: 최대 사용 횟수
        // [32-33]: 필요 채널 수 (u16 BE)
        // [34]: 측정 시간 (초)
        // [35]: 플래그
        // [36-43]: α 계수
        // [44-51]: 온도 보정 계수
        // [52-59]: 습도 보정 계수
        // [60-63]: CRC-32
        // [64+]: 확장 보정 데이터

        let cartridge_id = hex::encode(&data[0..8]);
        let category_code = data[8];
        let type_index = data[9];
        let legacy_code = data[10];
        let full_code = CartridgeFullCode::new(category_code, type_index);
        let category = CartridgeCategory::from_code(category_code);
        let cartridge_type = if legacy_code != 0 {
            CartridgeType::from_code(legacy_code)
        } else {
            CartridgeType::Unknown // v2.0 전용 타입은 레거시 enum에 없음
        };

        let lot_id = String::from_utf8_lossy(&data[12..20])
            .trim_end_matches('\0')
            .to_string();
        let expiry_str = String::from_utf8_lossy(&data[20..28]).to_string();
        let expiry_date = format!(
            "{}-{}-{}",
            &expiry_str[0..4],
            &expiry_str[4..6],
            &expiry_str[6..8]
        );

        let remaining_uses = u16::from_le_bytes([data[28], data[29]]) as u32;
        let max_uses = u16::from_le_bytes([data[30], data[31]]) as u32;

        let alpha = f64::from_le_bytes([
            data[36], data[37], data[38], data[39], data[40], data[41], data[42], data[43],
        ]);
        let temp_coefficient = f64::from_le_bytes([
            data[44], data[45], data[46], data[47], data[48], data[49], data[50], data[51],
        ]);
        let humidity_coefficient = f64::from_le_bytes([
            data[52], data[53], data[54], data[55], data[56], data[57], data[58], data[59],
        ]);

        let calibration_data = if data.len() > 64 {
            data[64..].to_vec()
        } else {
            Vec::new()
        };

        Ok(CartridgeInfo {
            cartridge_id,
            cartridge_type,
            full_code,
            category,
            tag_version: 2,
            lot_id,
            expiry_date,
            remaining_uses,
            max_uses,
            calibration_data,
            calibration_coefficients: CalibrationCoefficients {
                alpha,
                offsets: Vec::new(),
                gains: Vec::new(),
                temp_coefficient,
                humidity_coefficient,
            },
        })
    }

    /// 사용 횟수 업데이트 (사용 후)
    pub async fn decrement_usage(&mut self, cartridge_id: &str) -> Result<u32, NfcError> {
        // TODO: 실제 NFC 쓰기 구현

        if let Some(ref mut cartridge) = self.last_cartridge {
            if cartridge.cartridge_id == cartridge_id {
                if cartridge.remaining_uses > 0 {
                    cartridge.remaining_uses -= 1;
                    return Ok(cartridge.remaining_uses);
                } else {
                    return Err(NfcError::UsesExceeded);
                }
            }
        }

        Err(NfcError::NoCartridge)
    }

    /// 마지막으로 읽은 카트리지 정보
    pub fn last_cartridge(&self) -> Option<&CartridgeInfo> {
        self.last_cartridge.as_ref()
    }

    /// 카트리지 유효성 검증
    pub fn validate_cartridge(&self, cartridge: &CartridgeInfo) -> Result<(), NfcError> {
        // 만료 확인
        let today = chrono::Utc::now().format("%Y-%m-%d").to_string();
        if cartridge.expiry_date < today {
            return Err(NfcError::Expired);
        }

        // 사용 횟수 확인
        if cartridge.remaining_uses == 0 {
            return Err(NfcError::UsesExceeded);
        }

        // 알 수 없는 타입 확인
        if cartridge.cartridge_type == CartridgeType::Unknown {
            return Err(NfcError::InvalidCartridge(
                "알 수 없는 카트리지 타입".to_string(),
            ));
        }

        Ok(())
    }
}

impl Default for NfcReader {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cartridge_type_code() {
        assert_eq!(CartridgeType::Glucose.to_code(), 0x01);
        assert_eq!(CartridgeType::from_code(0x01), CartridgeType::Glucose);
        assert_eq!(CartridgeType::from_code(0x51), CartridgeType::NonTarget896);
        assert_eq!(CartridgeType::from_code(0x52), CartridgeType::NonTarget1792);
        assert_eq!(CartridgeType::NonTarget1792.to_code(), 0x52);
    }

    #[test]
    fn test_cartridge_channels() {
        assert_eq!(CartridgeType::Glucose.required_channels(), 88);
        assert_eq!(CartridgeType::NonTarget896.required_channels(), 896);
        assert_eq!(CartridgeType::NonTarget1792.required_channels(), 1792);
        assert_eq!(CartridgeType::ENose.required_channels(), 8);
    }

    #[test]
    fn test_cartridge_duration() {
        assert_eq!(CartridgeType::Glucose.measurement_duration_secs(), 15);
        assert_eq!(CartridgeType::NonTarget896.measurement_duration_secs(), 90);
        assert_eq!(
            CartridgeType::NonTarget1792.measurement_duration_secs(),
            180
        );
    }

    #[test]
    fn test_cartridge_category() {
        assert_eq!(
            CartridgeCategory::from_code(0x01),
            CartridgeCategory::HealthBiomarker
        );
        assert_eq!(
            CartridgeCategory::from_code(0x04),
            CartridgeCategory::ElectronicSensor
        );
        assert_eq!(CartridgeCategory::HealthBiomarker.to_code(), 0x01);
        assert_eq!(
            CartridgeCategory::HealthBiomarker.name_ko(),
            "건강 바이오마커"
        );
    }

    #[test]
    fn test_legacy_to_full_code() {
        let glucose = CartridgeFullCode::from_legacy(0x01);
        assert_eq!(glucose.category, 0x01);
        assert_eq!(glucose.type_index, 0x01);

        let water = CartridgeFullCode::from_legacy(0x20);
        assert_eq!(water.category, 0x02);
        assert_eq!(water.type_index, 0x01);

        let enose = CartridgeFullCode::from_legacy(0x40);
        assert_eq!(enose.category, 0x04);
        assert_eq!(enose.type_index, 0x01);

        let custom = CartridgeFullCode::from_legacy(0xFF);
        assert_eq!(custom.category, 0xFF);
        assert_eq!(custom.type_index, 0x01);
    }

    #[test]
    fn test_full_code_u16() {
        let code = CartridgeFullCode::new(0x01, 0x03);
        assert_eq!(code.to_u16(), 0x0103);
        let parsed = CartridgeFullCode::from_u16(0x0103);
        assert_eq!(parsed.category, 0x01);
        assert_eq!(parsed.type_index, 0x03);
    }

    #[test]
    fn test_cartridge_type_to_full_code() {
        assert_eq!(
            CartridgeType::Glucose.to_full_code(),
            CartridgeFullCode::new(0x01, 0x01)
        );
        assert_eq!(
            CartridgeType::Glucose.category(),
            CartridgeCategory::HealthBiomarker
        );
        assert_eq!(
            CartridgeType::ENose.category(),
            CartridgeCategory::ElectronicSensor
        );
    }

    #[test]
    fn test_registry_defaults() {
        let registry = CartridgeRegistry::new_with_defaults();
        // 14(건강) + 4(환경) + 4(식품) + 3(전자센서) + 4(고급분석: 448+896+1792+MultiBio) + 1(연구용) = 30
        assert_eq!(registry.count(), 30);

        // 레거시 코드로 조회
        let glucose = registry.get_by_legacy(0x01).unwrap();
        assert_eq!(glucose.name_ko, "혈당");
        assert_eq!(glucose.required_channels, 88);

        // 풀 코드로 조회
        let enose = registry.get(&CartridgeFullCode::new(0x04, 0x01)).unwrap();
        assert_eq!(enose.name_en, "ENose");
        assert_eq!(enose.required_channels, 8);

        // 카테고리별 목록
        let health = registry.list_by_category(0x01);
        assert_eq!(health.len(), 14);
    }

    #[test]
    fn test_registry_dynamic_register() {
        let mut registry = CartridgeRegistry::new_with_defaults();
        assert_eq!(registry.count(), 30); // 기본 30종

        // 신규 카트리지 타입 동적 등록
        registry.register(CartridgeRegistryEntry {
            full_code: CartridgeFullCode::new(0x06, 0x01),
            name_ko: "화학물질 분석".to_string(),
            name_en: "ChemicalAgent".to_string(),
            required_channels: 88,
            measurement_secs: 30,
            unit: "ppm".to_string(),
            reference_range: String::new(),
            is_beta: false,
            manufacturer: "ThirdParty Inc.".to_string(),
        });

        assert_eq!(registry.count(), 31); // 30 + 1 동적 추가
        let chem = registry.get(&CartridgeFullCode::new(0x06, 0x01)).unwrap();
        assert_eq!(chem.name_ko, "화학물질 분석");
        assert_eq!(chem.manufacturer, "ThirdParty Inc.");
    }

    #[test]
    fn test_parse_tag_v1() {
        let mut data = vec![0u8; 64];

        data[0..8].copy_from_slice(&[0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08]);
        data[8] = 0x01; // Glucose (v1.0 legacy code)
        data[9..17].copy_from_slice(b"LOT12345");
        data[17..25].copy_from_slice(b"20271231");
        data[25] = 50;
        data[26] = 0;
        data[27] = 100;
        data[28] = 0;
        let alpha_bytes = 0.95f64.to_le_bytes();
        data[29..37].copy_from_slice(&alpha_bytes);

        let reader = NfcReader::new();
        let cartridge = reader.parse_tag_data(&data).unwrap();

        assert_eq!(cartridge.cartridge_type, CartridgeType::Glucose);
        assert_eq!(cartridge.tag_version, 1);
        assert_eq!(cartridge.full_code, CartridgeFullCode::new(0x01, 0x01));
        assert_eq!(cartridge.category, CartridgeCategory::HealthBiomarker);
        assert_eq!(cartridge.remaining_uses, 50);
        assert!((cartridge.calibration_coefficients.alpha - 0.95).abs() < 0.001);
    }

    #[test]
    fn test_parse_tag_v2() {
        let mut data = vec![0u8; 80];

        data[0..8].copy_from_slice(&[0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07, 0x18]);
        data[8] = 0x06; // Industrial category (v2.0 전용)
        data[9] = 0x01; // type_index = 1
        data[10] = 0x00; // legacy_code = 없음
        data[11] = 0x02; // v2.0 태그
        data[12..20].copy_from_slice(b"LOT2XXXX");
        data[20..28].copy_from_slice(b"20281231");
        data[28] = 30;
        data[29] = 0;
        data[30] = 50;
        data[31] = 0;
        let alpha_bytes = 0.93f64.to_le_bytes();
        data[36..44].copy_from_slice(&alpha_bytes);

        let reader = NfcReader::new();
        let cartridge = reader.parse_tag_data(&data).unwrap();

        assert_eq!(cartridge.tag_version, 2);
        assert_eq!(cartridge.full_code, CartridgeFullCode::new(0x06, 0x01));
        assert_eq!(cartridge.category, CartridgeCategory::Industrial);
        assert_eq!(cartridge.cartridge_type, CartridgeType::Unknown); // 레거시에 없는 타입
        assert_eq!(cartridge.remaining_uses, 30);
        assert!((cartridge.calibration_coefficients.alpha - 0.93).abs() < 0.001);
    }
}
