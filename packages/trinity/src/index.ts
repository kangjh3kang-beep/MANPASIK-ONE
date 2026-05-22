/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family A
 * @trinity IP1
 * @sb SB-1
 * @standard IEC 62304 Class B
 *
 * MMUP Trinity IP classification utilities
 */

import { TRINITY } from "@mmup/ssot";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------
export interface TrinityClassification {
  ip: string;
  primary: string;
  secondary: string;
  tertiary: string;
}

export interface TrinityMapping {
  domainId: string;
  ip: string;
  axis: number;
  stage: number;
}

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

/** Maps all 9 MMUP axes (domains) to their Trinity IP classification */
export const TRINITY_DOMAINS: Record<string, { ip: string; axis: number; stage: number }> = {
  "universal-measurement":  { ip: "IP1", axis: 1, stage: 1 },
  "home-workplace-checkup": { ip: "IP3", axis: 2, stage: 1 },
  "data-analysis":          { ip: "IP3", axis: 3, stage: 2 },
  "ai-diagnosis":           { ip: "IP2", axis: 4, stage: 2 },
  "specialist-agent":       { ip: "IP3", axis: 5, stage: 4 },
  "prediction-prevention":  { ip: "IP3", axis: 6, stage: 3 },
  "prescription-coaching":  { ip: "IP3", axis: 7, stage: 4 },
  "compliance-verification":{ ip: "IP3", axis: 8, stage: 5 },
  "platform-common":        { ip: "IP3", axis: 9, stage: 1 },
} as const;

const TRINITY_LABELS: Record<string, Record<string, string>> = {
  IP1: { en: "Physical Matrix Removal", ko: "물리적 매트릭스 제거" },
  IP2: { en: "Differential Measurement", ko: "차동 측정" },
  IP3: { en: "AI Correction (distributed)", ko: "AI 보정 (분산)" },
};

const TRINITY_COLORS: Record<string, string> = {
  IP1: "#0ea5e9", // sky-500 — primary
  IP2: "#a855f7", // purple-500 — secondary
  IP3: "#10b981", // emerald-500 — tertiary
};

// ---------------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------------

/** Returns the TRINITY mapping for a given IP key */
export function getTrinityMapping(ip: keyof typeof TRINITY) {
  return TRINITY[ip];
}

/** Checks whether the given IP has properly distributed family assignments */
export function isProperlyDistributed(ip: keyof typeof TRINITY): boolean {
  const mapping = TRINITY[ip];
  return mapping.primary !== mapping.secondary && mapping.secondary !== mapping.tertiary;
}

/** Maps a domain ID to its Trinity IP classification */
export function classifyTrinityAxis(domainId: string): TrinityClassification {
  const domain = TRINITY_DOMAINS[domainId];
  if (!domain) {
    return { ip: "unknown", primary: "", secondary: "", tertiary: "" };
  }
  const trinityEntry = TRINITY[domain.ip as keyof typeof TRINITY];
  if (!trinityEntry) {
    return { ip: domain.ip, primary: "", secondary: "", tertiary: "" };
  }
  return {
    ip: domain.ip,
    primary: trinityEntry.primary,
    secondary: trinityEntry.secondary,
    tertiary: trinityEntry.tertiary,
  };
}

/** Returns a localized label for a Trinity IP key */
export function getTrinityLabel(ip: string, locale: string = "en"): string {
  const labels = TRINITY_LABELS[ip];
  if (!labels) {
    return `Unknown IP: ${ip}`;
  }
  return labels[locale] ?? labels["en"] ?? `Unknown IP: ${ip}`;
}

/** Validates a Trinity mapping object */
export function validateTrinityMapping(mapping: TrinityMapping): { valid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (!mapping.domainId || mapping.domainId.trim() === "") {
    errors.push("domainId is required");
  }

  const validIPs = ["IP1", "IP2", "IP3"];
  if (!validIPs.includes(mapping.ip)) {
    errors.push(`ip must be one of ${validIPs.join(", ")}, got "${mapping.ip}"`);
  }

  if (mapping.axis < 1 || mapping.axis > 9 || !Number.isInteger(mapping.axis)) {
    errors.push(`axis must be an integer between 1 and 9, got ${mapping.axis}`);
  }

  if (mapping.stage < 1 || mapping.stage > 8 || !Number.isInteger(mapping.stage)) {
    errors.push(`stage must be an integer between 1 and 8, got ${mapping.stage}`);
  }

  // Cross-validate: if domainId exists in TRINITY_DOMAINS, the ip should match
  const known = TRINITY_DOMAINS[mapping.domainId];
  if (known && known.ip !== mapping.ip) {
    errors.push(`domainId "${mapping.domainId}" is assigned to ${known.ip}, but mapping specifies ${mapping.ip}`);
  }

  return { valid: errors.length === 0, errors };
}

/** Returns the brand color hex for a Trinity IP axis */
export function getTrinityColor(ip: string): string {
  return TRINITY_COLORS[ip] ?? "#64748b"; // fallback to neutral-500
}
