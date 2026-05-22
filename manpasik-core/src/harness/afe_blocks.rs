// harness/afe_blocks.rs
//
// H4 (Progressive Extension): 9 Universal AFE block implementations
// with stage-based activation (Stage-1 → Stage-3).
//
// H1 (Module Independence): All blocks implement SensorTrait via HAL;
// no concrete hardware types are referenced directly.

use super::cartridge_manifest::CartridgeManifest;
use super::sensor_trait::{CalibrationData, SensorHealth, SensorTrait};

// ────────────────────────────────────────────────────────────────────
// Error types
// ────────────────────────────────────────────────────────────────────

/// Unified AFE error type (H6: failure isolation via Result<T, AfeError>)
#[derive(Debug, Clone, PartialEq)]
pub enum AfeError {
    NotInitialized,
    ChannelOutOfRange { channel: usize, max: usize },
    CalibrationFailed(String),
    DiagnosticFailed(String),
    UnsupportedStage(u8),
    AlphaOutOfRange { alpha: f64 },
    HardwareFault(String),
    ShutdownFailed(String),
}

impl core::fmt::Display for AfeError {
    fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        match self {
            Self::NotInitialized => write!(f, "AFE block not initialized"),
            Self::ChannelOutOfRange { channel, max } => {
                write!(f, "Channel {channel} out of range (max {max})")
            }
            Self::CalibrationFailed(msg) => write!(f, "Calibration failed: {msg}"),
            Self::DiagnosticFailed(msg) => write!(f, "Diagnostic failed: {msg}"),
            Self::UnsupportedStage(s) => write!(f, "Unsupported stage: {s}"),
            Self::AlphaOutOfRange { alpha } => {
                write!(f, "Alpha {alpha} outside valid range [0.90, 1.10]")
            }
            Self::HardwareFault(msg) => write!(f, "Hardware fault: {msg}"),
            Self::ShutdownFailed(msg) => write!(f, "Shutdown failed: {msg}"),
        }
    }
}

// ────────────────────────────────────────────────────────────────────
// Configuration types
// ────────────────────────────────────────────────────────────────────

/// Electrochemistry AFE configuration (LMP91000 x4 + ADS1256)
#[derive(Debug, Clone)]
pub struct ElectrochemConfig {
    /// Number of active channels (1..=4)
    pub channels: usize,
    /// Default alpha for differential correction (SSOT: 0.98, range 0.90..=1.10)
    pub alpha: f64,
    /// Sampling rate in Hz
    pub sample_rate_hz: u32,
}

impl Default for ElectrochemConfig {
    fn default() -> Self {
        Self {
            channels: 4,
            alpha: 0.98,
            sample_rate_hz: 100,
        }
    }
}

/// E-nose AFE configuration (8-channel gas sensor array)
#[derive(Debug, Clone)]
pub struct EnoseConfig {
    pub channels: usize,
    pub alpha: f64,
}

impl Default for EnoseConfig {
    fn default() -> Self {
        Self {
            channels: 8,
            alpha: 0.98,
        }
    }
}

/// EC array expansion configuration
#[derive(Debug, Clone)]
pub struct EcArrayConfig {
    pub channels: usize,
    pub alpha: f64,
}

impl Default for EcArrayConfig {
    fn default() -> Self {
        Self {
            channels: 8,
            alpha: 0.98,
        }
    }
}

/// SiPM-ECL detection configuration (Stage-2)
#[derive(Debug, Clone)]
pub struct SipmEclConfig {
    pub gain: f64,
    pub alpha: f64,
}

impl Default for SipmEclConfig {
    fn default() -> Self {
        Self {
            gain: 1.0,
            alpha: 0.98,
        }
    }
}

/// Colorimetric AFE configuration (Stage-2)
#[derive(Debug, Clone)]
pub struct ColorimetricConfig {
    pub led_wavelength_nm: u16,
    pub alpha: f64,
}

impl Default for ColorimetricConfig {
    fn default() -> Self {
        Self {
            led_wavelength_nm: 630,
            alpha: 0.98,
        }
    }
}

/// TEC (Thermoelectric Cooling) configuration (Stage-2)
#[derive(Debug, Clone)]
pub struct TecConfig {
    pub target_temp_c: f64,
    pub alpha: f64,
}

impl Default for TecConfig {
    fn default() -> Self {
        Self {
            target_temp_c: 37.0,
            alpha: 0.98,
        }
    }
}

/// LAMP/RPA isothermal amplification configuration (Stage-3)
#[derive(Debug, Clone)]
pub struct LampConfig {
    pub target_temp_c: f64,
    pub cycles: u32,
    pub alpha: f64,
}

impl Default for LampConfig {
    fn default() -> Self {
        Self {
            target_temp_c: 65.0,
            cycles: 30,
            alpha: 0.98,
        }
    }
}

// ────────────────────────────────────────────────────────────────────
// Helper: validate alpha is within SSOT range [0.90, 1.10]
// ────────────────────────────────────────────────────────────────────

fn validate_alpha(alpha: f64) -> Result<(), AfeError> {
    if !(0.90..=1.10).contains(&alpha) {
        return Err(AfeError::AlphaOutOfRange { alpha });
    }
    Ok(())
}

/// Apply the SSOT differential formula: Sdiff_n = S_n - alpha * R_n
/// Expects `raw` to contain interleaved [signal, reference] pairs,
/// i.e. length must be even. Returns one differential value per pair.
fn apply_differential_formula(raw: &[f64], alpha: f64) -> Result<Vec<f64>, AfeError> {
    validate_alpha(alpha)?;
    if raw.len() % 2 != 0 {
        return Err(AfeError::ChannelOutOfRange {
            channel: raw.len(),
            max: raw.len() - 1,
        });
    }
    let result = raw
        .chunks_exact(2)
        .map(|pair| pair[0] - alpha * pair[1])
        .collect();
    Ok(result)
}

// ════════════════════════════════════════════════════════════════════
// Block 1 — ElectrochemAfe (Stage-1, primary)
// LMP91000 x4 + ADS1256 — full implementation
// ════════════════════════════════════════════════════════════════════

pub struct ElectrochemAfe {
    channels: usize,
    initialized: bool,
    alpha: f64,
    sample_rate_hz: u32,
}

impl ElectrochemAfe {
    pub fn new() -> Self {
        Self {
            channels: 4,
            initialized: false,
            alpha: 0.98,
            sample_rate_hz: 100,
        }
    }
}

impl SensorTrait for ElectrochemAfe {
    type Config = ElectrochemConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        if config.channels == 0 || config.channels > 4 {
            return Err(AfeError::ChannelOutOfRange {
                channel: config.channels,
                max: 4,
            });
        }
        validate_alpha(config.alpha)?;
        self.channels = config.channels;
        self.alpha = config.alpha;
        self.sample_rate_hz = config.sample_rate_hz;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // In production this would DMA-read from ADS1256.
        // Return placeholder zeros — pairs of (signal, reference) per channel.
        Ok(vec![0.0; self.channels * 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // SSOT: Sdiff_n = S_n - alpha_n * R_n
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // H2: both v1 and v2 cartridges must be supported
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        // In production: read NFC tag calibration data and populate CalibrationData.
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        // Stub: would check LMP91000 register readback and ADS1256 self-test.
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 2 — EnoseAfe (Stage-1)
// 8-channel electronic nose array
// ════════════════════════════════════════════════════════════════════

pub struct EnoseAfe {
    channels: usize,
    initialized: bool,
    alpha: f64,
}

impl EnoseAfe {
    pub fn new() -> Self {
        Self {
            channels: 8,
            initialized: false,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for EnoseAfe {
    type Config = EnoseConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        if config.channels == 0 || config.channels > 8 {
            return Err(AfeError::ChannelOutOfRange {
                channel: config.channels,
                max: 8,
            });
        }
        validate_alpha(config.alpha)?;
        self.channels = config.channels;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        Ok(vec![0.0; self.channels * 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 3 — EcArrayAfe (Stage-1)
// EC array expansion for multi-analyte electrochemistry
// ════════════════════════════════════════════════════════════════════

pub struct EcArrayAfe {
    channels: usize,
    initialized: bool,
    alpha: f64,
}

impl EcArrayAfe {
    pub fn new() -> Self {
        Self {
            channels: 8,
            initialized: false,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for EcArrayAfe {
    type Config = EcArrayConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        if config.channels == 0 || config.channels > 8 {
            return Err(AfeError::ChannelOutOfRange {
                channel: config.channels,
                max: 8,
            });
        }
        validate_alpha(config.alpha)?;
        self.channels = config.channels;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        Ok(vec![0.0; self.channels * 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 4 — SipmEclAfe (Stage-2)
// SiPM-based electrochemiluminescence detection
// ════════════════════════════════════════════════════════════════════

pub struct SipmEclAfe {
    initialized: bool,
    gain: f64,
    alpha: f64,
}

impl SipmEclAfe {
    pub fn new() -> Self {
        Self {
            initialized: false,
            gain: 1.0,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for SipmEclAfe {
    type Config = SipmEclConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        validate_alpha(config.alpha)?;
        self.gain = config.gain;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // Stage-2 stub: SiPM photon count readout placeholder
        Ok(vec![0.0; 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 5 — ColorimetricAfe (Stage-2)
// LED + photodiode colorimetry
// ════════════════════════════════════════════════════════════════════

pub struct ColorimetricAfe {
    initialized: bool,
    led_wavelength_nm: u16,
    alpha: f64,
}

impl ColorimetricAfe {
    pub fn new() -> Self {
        Self {
            initialized: false,
            led_wavelength_nm: 630,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for ColorimetricAfe {
    type Config = ColorimetricConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        validate_alpha(config.alpha)?;
        self.led_wavelength_nm = config.led_wavelength_nm;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // Stage-2 stub: LED/PD absorbance reading placeholder
        Ok(vec![0.0; 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 6 — TecAfe (Stage-2)
// Thermoelectric cooling/heating control
// ════════════════════════════════════════════════════════════════════

pub struct TecAfe {
    initialized: bool,
    target_temp_c: f64,
    alpha: f64,
}

impl TecAfe {
    pub fn new() -> Self {
        Self {
            initialized: false,
            target_temp_c: 37.0,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for TecAfe {
    type Config = TecConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        validate_alpha(config.alpha)?;
        self.target_temp_c = config.target_temp_c;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // Stage-2 stub: temperature sensor reading placeholder
        Ok(vec![0.0; 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 7 — LampAfe (Stage-3)
// LAMP/RPA isothermal nucleic acid amplification
// ════════════════════════════════════════════════════════════════════

pub struct LampAfe {
    initialized: bool,
    target_temp_c: f64,
    cycles: u32,
    alpha: f64,
}

impl LampAfe {
    pub fn new() -> Self {
        Self {
            initialized: false,
            target_temp_c: 65.0,
            cycles: 30,
            alpha: 0.98,
        }
    }
}

impl SensorTrait for LampAfe {
    type Config = LampConfig;
    type RawData = Vec<f64>;
    type ProcessedData = Vec<f64>;
    type Error = AfeError;

    fn init(&mut self, config: &Self::Config) -> Result<(), Self::Error> {
        validate_alpha(config.alpha)?;
        self.target_temp_c = config.target_temp_c;
        self.cycles = config.cycles;
        self.alpha = config.alpha;
        self.initialized = true;
        Ok(())
    }

    fn read_raw(&self) -> Result<Self::RawData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        // Stage-3 stub: fluorescence readout placeholder
        Ok(vec![0.0; 2])
    }

    fn apply_differential(
        &self,
        raw: &Self::RawData,
        alpha: f64,
    ) -> Result<Self::ProcessedData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        apply_differential_formula(raw, alpha)
    }

    fn load_calibration(
        &mut self,
        manifest: &CartridgeManifest,
    ) -> Result<CalibrationData, Self::Error> {
        if !self.initialized {
            return Err(AfeError::NotInitialized);
        }
        if !manifest.is_v1() && !manifest.is_v2() {
            return Err(AfeError::CalibrationFailed(format!(
                "Unsupported CSI version: {}",
                manifest.csi_version
            )));
        }
        Ok(CalibrationData)
    }

    fn self_diagnose(&self) -> SensorHealth {
        SensorHealth
    }

    fn shutdown(&mut self) -> Result<(), Self::Error> {
        self.initialized = false;
        Ok(())
    }
}

// ════════════════════════════════════════════════════════════════════
// Block 8 & 9 — Reserved for future expansion
// ════════════════════════════════════════════════════════════════════

/// Reserved AFE block 8 — placeholder for future sensor modality
pub struct ReservedAfe8;

/// Reserved AFE block 9 — placeholder for future sensor modality
pub struct ReservedAfe9;

// Reserved blocks are not activated by any stage; they exist only to
// maintain the 9-block Universal AFE architecture contract.

// ════════════════════════════════════════════════════════════════════
// Stage-based activation factory (H4: Progressive Extension)
// ════════════════════════════════════════════════════════════════════

/// Describes which AFE blocks are active for a given stage.
/// Because SensorTrait has associated types, we cannot use a single
/// `dyn SensorTrait` trait object. Instead, this enum wraps concrete
/// block instances so they can be stored uniformly.
#[derive(Debug)]
pub enum AfeBlock {
    Electrochem(ElectrochemAfe),
    Enose(EnoseAfe),
    EcArray(EcArrayAfe),
    SipmEcl(SipmEclAfe),
    Colorimetric(ColorimetricAfe),
    Tec(TecAfe),
    Lamp(LampAfe),
}

/// Collection of AFE blocks activated for the current stage.
pub struct AfeBlockSet {
    pub blocks: Vec<AfeBlock>,
    pub stage: u8,
}

/// Initialize AFE blocks appropriate for the cartridge's declared stage.
///
/// - Stage 1 (Electrochemistry): ElectrochemAfe + EnoseAfe + EcArrayAfe
/// - Stage 2 (Optical): Stage-1 blocks + SipmEclAfe + ColorimetricAfe + TecAfe
/// - Stage 3 (NAAT): Stage-1 + Stage-2 blocks + LampAfe
///
/// H4: each stage is a strict superset of the previous stage.
/// H2: v1.0 cartridges are always supported (backward compatible).
pub fn init_afe_for_stage(manifest: &CartridgeManifest) -> Result<AfeBlockSet, AfeError> {
    let stage = manifest.stage;

    let mut blocks: Vec<AfeBlock> = Vec::new();

    // Stage 1 blocks (always included for stage >= 1)
    if stage >= 1 {
        blocks.push(AfeBlock::Electrochem(ElectrochemAfe::new()));
        blocks.push(AfeBlock::Enose(EnoseAfe::new()));
        blocks.push(AfeBlock::EcArray(EcArrayAfe::new()));
    }

    // Stage 2 blocks (optical detection)
    if stage >= 2 {
        blocks.push(AfeBlock::SipmEcl(SipmEclAfe::new()));
        blocks.push(AfeBlock::Colorimetric(ColorimetricAfe::new()));
        blocks.push(AfeBlock::Tec(TecAfe::new()));
    }

    // Stage 3 blocks (nucleic acid amplification)
    if stage >= 3 {
        blocks.push(AfeBlock::Lamp(LampAfe::new()));
    }

    if blocks.is_empty() || stage == 0 || stage > 3 {
        return Err(AfeError::UnsupportedStage(stage));
    }

    Ok(AfeBlockSet { blocks, stage })
}

// ════════════════════════════════════════════════════════════════════
// Tests
// ════════════════════════════════════════════════════════════════════

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn electrochem_init_and_differential() {
        let mut afe = ElectrochemAfe::new();
        let config = ElectrochemConfig::default();
        afe.init(&config).expect("init should succeed");

        // signal=1.0, reference=0.5 => Sdiff = 1.0 - 0.98*0.5 = 0.51
        let raw = vec![1.0, 0.5, 0.8, 0.3];
        let result = afe.apply_differential(&raw, 0.98).expect("diff ok");
        assert_eq!(result.len(), 2);
        assert!((result[0] - 0.51).abs() < 1e-10);
        assert!((result[1] - (0.8 - 0.98 * 0.3)).abs() < 1e-10);
    }

    #[test]
    fn electrochem_not_initialized_error() {
        let afe = ElectrochemAfe::new();
        assert_eq!(afe.read_raw(), Err(AfeError::NotInitialized));
    }

    #[test]
    fn alpha_out_of_range() {
        let mut afe = ElectrochemAfe::new();
        afe.init(&ElectrochemConfig::default()).expect("ok");
        let raw = vec![1.0, 0.5];
        assert!(matches!(
            afe.apply_differential(&raw, 1.15),
            Err(AfeError::AlphaOutOfRange { .. })
        ));
    }

    #[test]
    fn stage_1_factory() {
        let manifest = CartridgeManifest {
            csi_version: 1,
            stage: 1,
        };
        let set = init_afe_for_stage(&manifest).expect("stage 1 ok");
        assert_eq!(set.stage, 1);
        assert_eq!(set.blocks.len(), 3); // Electrochem + Enose + EcArray
    }

    #[test]
    fn stage_2_factory() {
        let manifest = CartridgeManifest {
            csi_version: 1,
            stage: 2,
        };
        let set = init_afe_for_stage(&manifest).expect("stage 2 ok");
        assert_eq!(set.stage, 2);
        assert_eq!(set.blocks.len(), 6); // Stage-1(3) + Stage-2(3)
    }

    #[test]
    fn stage_3_factory() {
        let manifest = CartridgeManifest {
            csi_version: 2,
            stage: 3,
        };
        let set = init_afe_for_stage(&manifest).expect("stage 3 ok");
        assert_eq!(set.stage, 3);
        assert_eq!(set.blocks.len(), 7); // Stage-1(3) + Stage-2(3) + Stage-3(1)
    }

    #[test]
    fn stage_0_unsupported() {
        let manifest = CartridgeManifest {
            csi_version: 1,
            stage: 0,
        };
        assert!(matches!(
            init_afe_for_stage(&manifest),
            Err(AfeError::UnsupportedStage(0))
        ));
    }

    #[test]
    fn stage_4_unsupported() {
        let manifest = CartridgeManifest {
            csi_version: 1,
            stage: 4,
        };
        assert!(matches!(
            init_afe_for_stage(&manifest),
            Err(AfeError::UnsupportedStage(4))
        ));
    }

    #[test]
    fn shutdown_resets_initialized() {
        let mut afe = ElectrochemAfe::new();
        afe.init(&ElectrochemConfig::default()).expect("ok");
        assert!(afe.read_raw().is_ok());
        afe.shutdown().expect("shutdown ok");
        assert_eq!(afe.read_raw(), Err(AfeError::NotInitialized));
    }
}
