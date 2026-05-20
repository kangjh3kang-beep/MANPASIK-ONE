/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

import { TRINITY } from "@mmup/ssot";

export function getTrinityMapping(ip: keyof typeof TRINITY) {
  return TRINITY[ip];
}

export function isProperlyDistributed(ip: keyof typeof TRINITY): boolean {
  const mapping = TRINITY[ip];
  return mapping.primary !== mapping.secondary && mapping.secondary !== mapping.tertiary;
}
