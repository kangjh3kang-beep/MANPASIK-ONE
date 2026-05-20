import { describe, it, expect } from 'vitest';
import { mockMeasurement } from './index';

/**
 * @mmup-axis 6 예측·예방
 * @mmup-stage 3 예측
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 */
describe('Test Utils - Mock Measurement', () => {
  it('should have required FHIR Observation structure', () => {
    expect(mockMeasurement.resourceType).toBe("Observation");
    expect(mockMeasurement.code.coding[0].system).toBe("http://loinc.org");
  });

  it('should include MMUP specific meta extensions', () => {
    expect(mockMeasurement.meta?.deviceId).toBeDefined();
    expect(mockMeasurement.meta?.differentialMode).toBeDefined();
    expect(mockMeasurement.meta?.hashChain).toBe("mock-hash");
  });
});
