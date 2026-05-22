use anyhow::Result;
use image::{ImageBuffer, Rgba};
use ndarray::Array3;
use serde::{Deserialize, Serialize};

use crate::crypto::CryptoEngine;
use crate::differential::DifferentialEngine;
use crate::fingerprint;

/// 카트리지 분석 결과
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AnalysisResult {
    /// 감지 성공 여부
    pub detected: bool,
    /// 바이오마커 이름
    pub biomarker: String,
    /// 측정값 (임의 단위)
    pub value: f64,
    /// 신뢰도 (0.0 ~ 1.0)
    pub confidence: f64,
    /// 색상 영역 RGB 평균
    pub dominant_rgb: [u8; 3],
    /// 경고 메시지 (없으면 빈 문자열)
    pub warning: String,
}

/// 측정 파이프라인 무결성 검증 결과 (IEC 62304 해시체인)
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MeasurementIntegrityResult {
    /// 차동 보정된 신호 배열
    pub corrected_values: Vec<f64>,
    /// 88차원 핑거프린트 벡터 (f32)
    pub fingerprint: Vec<f32>,
    /// 최종 해시체인 해시 (hex)
    pub chain_hash: String,
    /// 체인 무결성 검증 결과
    pub chain_valid: bool,
    /// 체인 엔트리 수
    pub chain_entries: usize,
    /// AI 예측 시뮬레이션 값 (향후 실제 TFLite 연동)
    pub ai_prediction: f64,
}

/// 측정 파이프라인 + IEC 62304 해시체인 무결성 검증
///
/// 차동측정, 핑거프린트 생성, AI 예측의 전 과정을 해시체인으로
/// 기록하여 데이터 변조를 감지합니다 (H3 사전인증 요구사항).
///
/// # Arguments
/// * `s_det` - 검출 전극 신호 배열
/// * `s_ref` - 기준 전극 신호 배열
/// * `alpha` - 보정 계수 (SSOT 범위: 0.90~1.10)
///
/// # Returns
/// 보정된 값, 핑거프린트, 해시체인 검증 결과를 포함하는 `MeasurementIntegrityResult`
///
/// # Errors
/// 채널 수 불일치, 핑거프린트 생성 실패, 해시체인 검증 실패 시 에러 반환
pub fn process_measurement_with_integrity(
    s_det: &[f64],
    s_ref: &[f64],
    alpha: f64,
) -> std::result::Result<MeasurementIntegrityResult, String> {
    if s_det.is_empty() {
        return Err("s_det 신호 배열이 비어있습니다".to_string());
    }
    if s_det.len() != s_ref.len() {
        return Err(format!(
            "s_det({})와 s_ref({})의 길이가 다릅니다",
            s_det.len(),
            s_ref.len()
        ));
    }

    let mut crypto = CryptoEngine::new();

    // Step 1: Record raw input hash
    let raw_bytes = serialize_signal_pair(s_det, s_ref);
    crypto.add_chain_entry("raw_capture", &raw_bytes);

    // Step 2: Differential correction (Sdiff_n = S_n - alpha_n * R_n)
    let clamped_alpha = alpha.clamp(crate::ALPHA_MIN, crate::ALPHA_MAX);
    let diff_engine = DifferentialEngine::with_defaults(s_det.len());
    // with_defaults uses SSOT alpha; we need custom alpha via a params update
    let mut diff_engine = diff_engine;
    diff_engine.set_alpha(clamped_alpha);

    let corrected = diff_engine
        .measure(s_det, s_ref)
        .map_err(|e| e.to_string())?;

    let corrected_bytes = serialize_f64_slice(&corrected);
    crypto.add_chain_entry("differential_correction", &corrected_bytes);

    // Step 3: Fingerprint generation (88-dim)
    let fingerprint = generate_fingerprint_88(&corrected);
    let fp_bytes = serialize_f32_slice(&fingerprint);
    crypto.add_chain_entry("fingerprint_generation", &fp_bytes);

    // Step 4: AI prediction (simulated — real TFLite INT8 inference in Phase 6)
    // Compute a deterministic prediction from the corrected signal mean
    let prediction = if corrected.is_empty() {
        0.0
    } else {
        let mean = corrected.iter().sum::<f64>() / corrected.len() as f64;
        // Sigmoid-like normalization for a 0..1 confidence proxy
        1.0 / (1.0 + (-mean).exp())
    };
    let prediction_bytes = prediction.to_le_bytes();
    crypto.add_chain_entry("ai_prediction", &prediction_bytes);

    // Step 5: Verify chain integrity
    let chain_valid = crypto.verify_chain().map_err(|e| format!("{:?}", e))?;

    let chain_hash = crypto
        .current_chain_hash()
        .unwrap_or("")
        .to_string();

    let chain_entries = crypto.chain_entries().len();

    Ok(MeasurementIntegrityResult {
        corrected_values: corrected,
        fingerprint,
        chain_hash,
        chain_valid,
        chain_entries,
        ai_prediction: prediction,
    })
}

/// 88차원 핑거프린트 생성 (차동 보정 결과로부터)
///
/// 입력 채널 수에 관계없이 88차원 벡터를 생성합니다.
/// 채널이 88개 미만이면 통계 피처(평균, 분산, 기울기 등)로 패딩합니다.
fn generate_fingerprint_88(corrected: &[f64]) -> Vec<f32> {
    let dim = fingerprint::DIM_88;
    let mut fp = Vec::with_capacity(dim);

    // 직접 채널 값 복사 (88개까지)
    for &v in corrected.iter().take(dim) {
        fp.push(v as f32);
    }

    // 남은 슬롯을 통계 피처로 채움
    if fp.len() < dim {
        let n = corrected.len() as f64;
        let mean = if n > 0.0 {
            corrected.iter().sum::<f64>() / n
        } else {
            0.0
        };
        let variance = if n > 1.0 {
            corrected.iter().map(|v| (v - mean).powi(2)).sum::<f64>() / (n - 1.0)
        } else {
            0.0
        };
        let std_dev = variance.sqrt();
        let max_val = corrected
            .iter()
            .copied()
            .fold(f64::NEG_INFINITY, f64::max);
        let min_val = corrected
            .iter()
            .copied()
            .fold(f64::INFINITY, f64::min);
        let range = max_val - min_val;

        // 통계 피처 패딩
        let stats = [mean, variance, std_dev, max_val, min_val, range];
        let mut stat_idx = 0;
        while fp.len() < dim {
            fp.push(stats[stat_idx % stats.len()] as f32);
            stat_idx += 1;
        }
    }

    fp.truncate(dim);
    fp
}

/// f64 슬라이스를 바이트로 직렬화 (해시체인용)
fn serialize_f64_slice(data: &[f64]) -> Vec<u8> {
    data.iter().flat_map(|v| v.to_le_bytes()).collect()
}

/// f32 슬라이스를 바이트로 직렬화 (해시체인용)
fn serialize_f32_slice(data: &[f32]) -> Vec<u8> {
    data.iter().flat_map(|v| v.to_le_bytes()).collect()
}

/// 신호 쌍(s_det + s_ref)을 바이트로 직렬화 (해시체인용)
fn serialize_signal_pair(s_det: &[f64], s_ref: &[f64]) -> Vec<u8> {
    let mut bytes = Vec::with_capacity((s_det.len() + s_ref.len()) * 8);
    for v in s_det {
        bytes.extend_from_slice(&v.to_le_bytes());
    }
    for v in s_ref {
        bytes.extend_from_slice(&v.to_le_bytes());
    }
    bytes
}

/// Flutter에서 호출되는 메인 FFI 엔트리포인트.
///
/// # Arguments
/// * `bytes` — BGRA8888 원시 픽셀 바이트 (카메라 스트림)
/// * `width` — 프레임 너비
/// * `height` — 프레임 높이
///
/// # Returns
/// JSON 문자열 형태의 `AnalysisResult`
pub fn process_camera_frame(bytes: Vec<u8>, width: u32, height: u32) -> String {
    match _process_internal(bytes, width, height) {
        Ok(result) => serde_json::to_string(&result).unwrap_or_else(|e| {
            format!("{{\"detected\":false,\"biomarker\":\"\",\"value\":0,\"confidence\":0,\"dominant_rgb\":[0,0,0],\"warning\":\"직렬화 실패: {e}\"}}")
        }),
        Err(e) => {
            let err_result = AnalysisResult {
                detected: false,
                biomarker: String::new(),
                value: 0.0,
                confidence: 0.0,
                dominant_rgb: [0, 0, 0],
                warning: format!("분석 오류: {e}"),
            };
            serde_json::to_string(&err_result).unwrap_or_default()
        }
    }
}

/// 내부 분석 파이프라인
fn _process_internal(bytes: Vec<u8>, width: u32, height: u32) -> Result<AnalysisResult> {
    let expected_len = (width * height * 4) as usize;
    if bytes.len() < expected_len {
        anyhow::bail!(
            "프레임 크기 불일치: 예상 {}B, 실제 {}B",
            expected_len,
            bytes.len()
        );
    }

    // 1) BGRA → RGBA 변환 및 ImageBuffer 생성
    let mut rgba_bytes = bytes[..expected_len].to_vec();
    for pixel in rgba_bytes.chunks_exact_mut(4) {
        pixel.swap(0, 2); // B <-> R
    }

    let img: ImageBuffer<Rgba<u8>, Vec<u8>> =
        ImageBuffer::from_raw(width, height, rgba_bytes.clone())
            .ok_or_else(|| anyhow::anyhow!("ImageBuffer 생성 실패"))?;

    // 2) ROI(Region of Interest) 추출 — 중앙 40% 영역
    let roi_x = (width as f32 * 0.3) as u32;
    let roi_y = (height as f32 * 0.3) as u32;
    let roi_w = (width as f32 * 0.4) as u32;
    let roi_h = (height as f32 * 0.4) as u32;

    // 3) ndarray 텐서로 변환 (ROI 영역만)
    let mut tensor = Array3::<f64>::zeros((roi_h as usize, roi_w as usize, 3));
    let mut r_sum: u64 = 0;
    let mut g_sum: u64 = 0;
    let mut b_sum: u64 = 0;
    let mut pixel_count: u64 = 0;

    for dy in 0..roi_h {
        for dx in 0..roi_w {
            let px = img.get_pixel(roi_x + dx, roi_y + dy);
            let r = px[0] as f64;
            let g = px[1] as f64;
            let b = px[2] as f64;

            tensor[[dy as usize, dx as usize, 0]] = r / 255.0;
            tensor[[dy as usize, dx as usize, 1]] = g / 255.0;
            tensor[[dy as usize, dx as usize, 2]] = b / 255.0;

            r_sum += px[0] as u64;
            g_sum += px[1] as u64;
            b_sum += px[2] as u64;
            pixel_count += 1;
        }
    }

    if pixel_count == 0 {
        anyhow::bail!("ROI 영역에 유효한 픽셀이 없습니다");
    }

    let avg_r = (r_sum / pixel_count) as u8;
    let avg_g = (g_sum / pixel_count) as u8;
    let avg_b = (b_sum / pixel_count) as u8;

    // 4) 색상 기반 카트리지 반응 감지
    //    - 카트리지 반응 영역은 특정 색상 범위를 가짐
    //    - 여기서는 색상 히스토그램 분석으로 반응 여부 판별
    let (detected, biomarker, value, confidence) =
        analyze_color_response(avg_r, avg_g, avg_b, &tensor);

    // 5) 경고 생성
    let warning = if !detected {
        "카트리지가 감지되지 않았습니다. 카메라를 카트리지 위에 정확히 위치시켜 주세요.".to_string()
    } else if confidence < 0.5 {
        "측정 신뢰도가 낮습니다. 조명을 확인해 주세요.".to_string()
    } else {
        String::new()
    };

    Ok(AnalysisResult {
        detected,
        biomarker,
        value,
        confidence,
        dominant_rgb: [avg_r, avg_g, avg_b],
        warning,
    })
}

/// 색상 응답 분석 — ROI 평균 RGB + 텐서 통계로 바이오마커 판별
fn analyze_color_response(r: u8, g: u8, b: u8, tensor: &Array3<f64>) -> (bool, String, f64, f64) {
    // 채널별 평균/표준편차 계산
    let shape = tensor.shape();
    let total = (shape[0] * shape[1]) as f64;

    let ch_means: Vec<f64> = (0..3)
        .map(|c| tensor.slice(ndarray::s![.., .., c]).iter().sum::<f64>() / total)
        .collect();

    let ch_stds: Vec<f64> = (0..3)
        .map(|c| {
            let mean = ch_means[c];
            let var = tensor
                .slice(ndarray::s![.., .., c])
                .iter()
                .map(|&v| (v - mean).powi(2))
                .sum::<f64>()
                / total;
            var.sqrt()
        })
        .collect();

    // 균일도 지표: 표준편차가 낮을수록 카트리지 반응 영역일 가능성 높음
    let uniformity = 1.0 - (ch_stds.iter().sum::<f64>() / 3.0).min(1.0);

    // 카트리지 반응 감지 기준: 균일도 > 0.7 이고 특정 색상 범위
    let detected = uniformity > 0.7 && (r > 30 || g > 30 || b > 30);

    if !detected {
        return (false, String::new(), 0.0, 0.0);
    }

    // 색상 → 바이오마커 매핑 (카트리지 색변 패턴)
    let (biomarker, value, confidence) = if r > 150 && g < 100 && b < 100 {
        // 적색 → pH (산성)
        let ph = 4.0 + (r as f64 - 150.0) / 105.0 * 3.0;
        ("pH".to_string(), ph, uniformity * 0.95)
    } else if r < 100 && g > 150 && b < 100 {
        // 녹색 → 글루코스
        let glucose = 70.0 + (g as f64 - 150.0) / 105.0 * 130.0;
        ("Glucose".to_string(), glucose, uniformity * 0.92)
    } else if r < 100 && g < 100 && b > 150 {
        // 청색 → 단백질
        let protein = 0.0 + (b as f64 - 150.0) / 105.0 * 30.0;
        ("Protein".to_string(), protein, uniformity * 0.88)
    } else if r > 150 && g > 100 && b < 80 {
        // 황색 → 빌리루빈
        let bilirubin = 0.1 + ((r as f64 + g as f64) / 2.0 - 125.0) / 130.0 * 5.0;
        ("Bilirubin".to_string(), bilirubin, uniformity * 0.85)
    } else {
        // 기타 색상 → 일반 반응도
        let intensity = (r as f64 + g as f64 + b as f64) / (3.0 * 255.0) * 100.0;
        ("General".to_string(), intensity, uniformity * 0.7)
    };

    (true, biomarker, value, confidence)
}
