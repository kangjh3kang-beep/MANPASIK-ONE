import { describe, it, expect } from 'vitest';
import { MMUP, FAMILY, TRINITY, BUSINESS_SSOT, DEPRECATED_FOREVER } from './index';

/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
describe('SSOT Constants', () => {
  it('should have correct MMUP definitions', () => {
    expect(MMUP.abbr).toBe("MMUP");
    expect(MMUP.fullName_KO).toContain("만파식");
  });

  it('should properly define FAMILY', () => {
    expect(FAMILY.A.title).toBeDefined();
    expect(FAMILY.B.id).toBe("B");
    expect(FAMILY.C.appNo).toBeDefined();
  });

  it('should not contain deprecated terms in BUSINESS_SSOT', () => {
    const serialized = JSON.stringify(BUSINESS_SSOT).toLowerCase();
    DEPRECATED_FOREVER.forEach(term => {
      expect(serialized).not.toContain(term.toLowerCase());
    });
  });
});
