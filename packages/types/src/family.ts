/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export type FamilyId = "A" | "B" | "C";

export interface FamilyMapping {
  family: FamilyId;
  description: string;
}
