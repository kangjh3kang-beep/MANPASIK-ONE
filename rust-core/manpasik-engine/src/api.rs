use anyhow::Result;
use image::{ImageBuffer, Rgba};
use ndarray::Array3;
use serde::{Deserialize, Serialize};

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
            format!("{{\"detected\":false,\"biomarker\":\"\",\"value\":0,\"confidence\":0,\"dominant_rgb\":[0,0,0],\"warning\":\"직렬화 실패: {}\"}}", e)
        }),
        Err(e) => {
            let err_result = AnalysisResult {
                detected: false,
                biomarker: String::new(),
                value: 0.0,
                confidence: 0.0,
                dominant_rgb: [0, 0, 0],
                warning: format!("분석 오류: {}", e),
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
