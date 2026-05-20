/**
 * @mmup-axis 6 예측·예방
 * @mmup-stage 3 예측
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */

import type { MMUPMeasurement } from "@mmup/types";

export const mockMeasurement: MMUPMeasurement = {
  resourceType: "Observation",
  status: "final",
  code: {
    coding: [{ system: "http://loinc.org", code: "1234-5" }]
  },
  meta: {
    deviceId: "DEV-12345",
    cartridgeUDI: "00123456789012",
    skuLayer: "L1",
    skuId: "SKU-BASIC-01",
    confidence_135P5: 0.95,
    differentialMode: true,
    twinModelVersion: "1.0",
    aiCorrectionApplied: true,
    hashChain: "mock-hash"
  }
};
