/**
 * @mmup-axis 5 분야별 전문의 에이전트
 * @mmup-stage 4 처방
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export const PERSONA = {
  P1: "일반 사용자",
  P2: "만성질환자",
  P3: "부모(영유아)",
  P4: "전문의",
  P5: "검사기사·약사",
  P6: "투자자·파트너",
  P7: "규제·감사"
} as const;

export type PersonaType = keyof typeof PERSONA;
