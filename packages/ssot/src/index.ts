/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 * 
 * MMUP SSOT v1.6 — 영구 단일출처 (변경 금지·Board 결의 시만 갱신)
 * 만파식: 모든 축 / 단계: 모든 단계의 기준
 */

export const MMUP = {
  fullName_KO: "MPS 만파식 다중측정원리 유니버설 POCT 플랫폼(MMUP)",
  fullName_EN: "MPS Manpasik Multi-Measurement Universal POCT Platform (MMUP)",
  abbr: "MMUP",
} as const;

export const FAMILY = {
  A: { id: "A", appNo: "APP2026-0022KR", title: "디지털 트윈 기반의 적응형 유체장 제어 장치 및 방법" },
  B: { id: "B", appNo: "APP2025-0967KR", title: "차동 측정 기반 범용 분석 장치 및 방법" },
  C: { id: "C", appNo: "APP2025-0968KR", title: "플랫폼 공통층과 결합된 모듈러 측정 장치 및 방법" },
} as const;

export const TRINITY = {
  IP1: { name: "Physical Matrix Removal", primary: "A", secondary: "B", tertiary: "C" },
  IP2: { name: "Differential Measurement", primary: "B", secondary: "A", tertiary: "C" },
  IP3: { name: "AI Correction (distributed)", primary: "C", secondary: "A", tertiary: "B" },
} as const;

export const BUSINESS_SSOT = {
  fxKRWperUSD: 1480,
  skuTotal: 44,
  layerSkus: { L1: 16, L2: 12, L3: 10, L4: 8, L5: 6 },
  clinicalAssays: 169,
  tamUSD_B: 34.47,
  samUSD_B: 8.0,
  somUSD_B: 0.13,
  fundingKRW_B: { seriesA: 6.5, seriesB: 5.0, seriesC: 15.0 },
} as const;

export const MFG_SSOT = {
  connector: "Samtec MECF-08-01-L-DV (16핀, 1.27mm)",
  parallelPlane_mm: { length: 49.7, width: 30, height: 4.3 },
  differentialFormula: "Sdiff_n = S_n - alpha_n * R_n",
} as const;

export const FAMILY_C_PARTS = {
  dockingInterface: 110,
  dataPacketEngine: 130,
  uncertainty: "135-P5",
  updateRollback: 140,
  safetyGuard: 150,
  multimodalScheduler: 160,
  quietWindow: 164,
  sensorFusion: 170,
  hmiSafety: 175,
  securityModule: 190,
  tpm: "191-T",
  pcr: 194,
  measurementEngines: [200, 300, 400],
} as const;

export const DEPRECATED_FOREVER = [
  "9블록", "9-block", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9",
  "Family L",
  "100% 완성", "Zero 오류",
  "소형 정량 면역분석기", "혈액 전용 분석기", "유니버설 POCT 분석장치 (단축)",
] as const;

export const PRIORITY = [
  "환자 안전",
  "규제 준수",
  "SSOT 무결성",
  "사실성",
  "절차",
  "형식·간결성",
] as const;
