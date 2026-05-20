import { describe, it, expect } from 'vitest';
import { getTrinityMapping, isProperlyDistributed } from './index';

/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
describe('Trinity Mapping', () => {
  it('should return correct mapping for IP1', () => {
    const mapping = getTrinityMapping('IP1');
    expect(mapping.name).toBe("Physical Matrix Removal");
    expect(mapping.primary).toBe("A");
  });

  it('should correctly ensure IP distribution is proper', () => {
    expect(isProperlyDistributed('IP1')).toBe(true);
    expect(isProperlyDistributed('IP2')).toBe(true);
    expect(isProperlyDistributed('IP3')).toBe(true);
  });
});
