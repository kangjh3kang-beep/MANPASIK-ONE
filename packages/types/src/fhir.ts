/**
 * @mmup-axis 3 종합 데이터 분석
 * @mmup-stage 2 진단
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */

export interface Observation {
  resourceType: string;
  meta?: Record<string, any>;
  [key: string]: any;
}

export interface MMUPMeasurement extends Observation {
  meta: Observation["meta"] & {
    deviceId: string;
    cartridgeUDI: string;
    skuLayer: "L1" | "L2" | "L3" | "L4" | "L5";
    skuId: string;
    confidence_135P5: number;
    differentialMode: boolean;
    twinModelVersion: string;
    aiCorrectionApplied: boolean;
    hashChain: string;
  };
}
