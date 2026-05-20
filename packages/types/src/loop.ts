/**
 * @mmup-axis 8 이행·검증·보상
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class B
 */

export const LOOP_STAGES = {
  1: "measure",
  2: "diagnose",
  3: "predict",
  4: "prescribe",
  5: "execute",
  6: "verify",
  7: "reward",
  8: "remeasure",
} as const;

export type LoopStage = keyof typeof LOOP_STAGES;

export const PRESCRIPTION_TYPES = [
  "clinic_referral",
  "diet_adjustment",
  "exercise_plan",
  "environment_setup",
  "supplement_recommendation",
] as const;
