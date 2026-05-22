/**
 * @mmup-axis 6 예측·예방
 * @mmup-stage 3 예측
 * @family C
 * @trinity IP3 (분산)
 * @sb SB-1
 * @standard IEC 62304 Class C
 *
 * MMUP Test Utilities — mock factories, provider wrappers, MSW helpers, matchers
 */

import React from "react";
import type { MMUPMeasurement } from "@mmup/types";

// ---------------------------------------------------------------------------
// Types (local — avoid hard dependency on third-party type packages)
// ---------------------------------------------------------------------------
export interface Patient {
  id: string;
  resourceType: "Patient";
  name: string;
  birthDate: string;
  gender: "male" | "female" | "other" | "unknown";
  identifier: string;
}

export interface Device {
  id: string;
  resourceType: "Device";
  manufacturer: string;
  modelNumber: string;
  serialNumber: string;
  status: "active" | "inactive" | "entered-in-error";
  udiCarrier: string;
}

export interface VitalSign {
  id: string;
  resourceType: "Observation";
  status: "final" | "preliminary";
  code: string;
  value: number;
  unit: string;
  effectiveDateTime: string;
}

export interface AuditLog {
  id: string;
  timestamp: string;
  actor: string;
  action: string;
  resource: string;
  outcome: "success" | "failure";
  detail: string;
}

// ---------------------------------------------------------------------------
// Mock measurement (kept from v1 for backward compat)
// ---------------------------------------------------------------------------
export const mockMeasurement: MMUPMeasurement = {
  resourceType: "Observation",
  status: "final",
  code: {
    coding: [{ system: "http://loinc.org", code: "1234-5" }],
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
    hashChain: "mock-hash",
  },
};

// ---------------------------------------------------------------------------
// Factory functions
// ---------------------------------------------------------------------------

let _counter = 0;
function nextId(prefix: string): string {
  _counter += 1;
  return `${prefix}-${_counter.toString().padStart(5, "0")}`;
}

/** Creates a mock MMUPMeasurement with sensible defaults */
export function createMockMeasurement(overrides?: Partial<MMUPMeasurement>): MMUPMeasurement {
  return {
    resourceType: "Observation",
    status: "final",
    code: {
      coding: [{ system: "http://loinc.org", code: "2345-7" }],
    },
    meta: {
      deviceId: nextId("DEV"),
      cartridgeUDI: "00123456789012",
      skuLayer: "L1",
      skuId: "SKU-BASIC-01",
      confidence_135P5: 0.92,
      differentialMode: true,
      twinModelVersion: "1.0",
      aiCorrectionApplied: false,
      hashChain: `hash-${Date.now()}`,
    },
    ...overrides,
    meta: {
      deviceId: nextId("DEV"),
      cartridgeUDI: "00123456789012",
      skuLayer: "L1",
      skuId: "SKU-BASIC-01",
      confidence_135P5: 0.92,
      differentialMode: true,
      twinModelVersion: "1.0",
      aiCorrectionApplied: false,
      hashChain: `hash-${Date.now()}`,
      ...(overrides?.meta ?? {}),
    },
  };
}

/** Creates a mock Patient resource */
export function createMockPatient(overrides?: Partial<Patient>): Patient {
  return {
    id: nextId("PAT"),
    resourceType: "Patient",
    name: "Test Patient",
    birthDate: "1990-01-15",
    gender: "unknown",
    identifier: nextId("MRN"),
    ...overrides,
  };
}

/** Creates a mock Device resource */
export function createMockDevice(overrides?: Partial<Device>): Device {
  return {
    id: nextId("DEV"),
    resourceType: "Device",
    manufacturer: "MPS ManPaSik",
    modelNumber: "MMUP-R1",
    serialNumber: nextId("SN"),
    status: "active",
    udiCarrier: "00123456789012",
    ...overrides,
  };
}

/** Creates an array of mock VitalSign observations */
export function createMockVitalSigns(count: number = 3): VitalSign[] {
  const codes = [
    { code: "8867-4", unit: "bpm", label: "heart-rate", range: [60, 100] },
    { code: "8310-5", unit: "°C", label: "body-temp", range: [36.0, 37.5] },
    { code: "2708-6", unit: "%", label: "spo2", range: [95, 100] },
    { code: "8480-6", unit: "mmHg", label: "systolic-bp", range: [90, 140] },
    { code: "8462-4", unit: "mmHg", label: "diastolic-bp", range: [60, 90] },
    { code: "8302-2", unit: "cm", label: "height", range: [150, 200] },
    { code: "29463-7", unit: "kg", label: "weight", range: [40, 120] },
    { code: "39156-5", unit: "kg/m2", label: "bmi", range: [18, 35] },
  ];

  return Array.from({ length: count }, (_, i) => {
    const spec = codes[i % codes.length];
    const value = spec.range[0] + Math.round((spec.range[1] - spec.range[0]) * 0.5 * 100) / 100;
    return {
      id: nextId("VS"),
      resourceType: "Observation" as const,
      status: "final" as const,
      code: spec.code,
      value,
      unit: spec.unit,
      effectiveDateTime: new Date(Date.now() - i * 3600_000).toISOString(),
    };
  });
}

/** Creates a mock AuditLog entry */
export function createMockAuditLog(overrides?: Partial<AuditLog>): AuditLog {
  return {
    id: nextId("AUDIT"),
    timestamp: new Date().toISOString(),
    actor: "system",
    action: "read",
    resource: "Observation/12345",
    outcome: "success",
    detail: "Mock audit log entry",
    ...overrides,
  };
}

// ---------------------------------------------------------------------------
// Provider wrapper for testing (React / Next.js)
// ---------------------------------------------------------------------------

/**
 * Wraps a React element with common providers for testing.
 *
 * Usage (with @testing-library/react):
 *   const result = renderWithProviders(<MyComponent />);
 *
 * This is a lightweight wrapper — it returns a render() call with providers.
 * Consumers must have @testing-library/react installed.
 */
export function renderWithProviders(
  ui: React.ReactElement,
  options?: { session?: Record<string, unknown>; queryClient?: unknown },
) {
  // Lazy-require so the package works even when @testing-library/react is absent
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  let render: (ui: React.ReactElement) => unknown;
  try {
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const rtl = require("@testing-library/react");
    render = rtl.render;
  } catch {
    throw new Error(
      "renderWithProviders requires @testing-library/react to be installed as a devDependency.",
    );
  }

  // Build a wrapper that injects providers
  const Wrapper = ({ children }: { children: React.ReactNode }) => {
    // If the consumer supplies a queryClient (e.g. from @tanstack/react-query),
    // we'd wrap it here. For now, we just pass children through.
    return React.createElement(React.Fragment, null, children);
  };

  return render(React.createElement(Wrapper, null, ui));
}

// ---------------------------------------------------------------------------
// MSW helpers (mock-service-worker)
// ---------------------------------------------------------------------------

/**
 * Sets up a mock service worker server for API testing.
 * Requires `msw` to be installed.
 */
export function setupMockServer(handlers?: unknown[]) {
  try {
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const { setupServer } = require("msw/node");
    const server = setupServer(...(handlers ?? []));
    return server;
  } catch {
    throw new Error("setupMockServer requires msw to be installed as a devDependency.");
  }
}

/**
 * Creates a single MSW request handler.
 * Requires `msw` to be installed.
 */
export function createMockHandler(method: string, path: string, response: unknown, status: number = 200) {
  try {
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    const { http, HttpResponse } = require("msw");
    const methodFn = http[method.toLowerCase() as keyof typeof http];
    if (typeof methodFn !== "function") {
      throw new Error(`Unsupported HTTP method: ${method}`);
    }
    return methodFn(path, () => {
      return HttpResponse.json(response, { status });
    });
  } catch (err) {
    if (err instanceof Error && err.message.startsWith("Unsupported")) throw err;
    throw new Error("createMockHandler requires msw to be installed as a devDependency.");
  }
}

// ---------------------------------------------------------------------------
// Custom matchers / assertion helpers
// ---------------------------------------------------------------------------

/** Asserts that a resource object is FHIR-compliant (basic structural checks) */
export function expectFhirCompliant(resource: Record<string, unknown>): void {
  if (!resource) {
    throw new Error("Expected a FHIR resource but received null/undefined");
  }
  if (typeof resource.resourceType !== "string" || resource.resourceType.trim() === "") {
    throw new Error(
      `FHIR resource must have a non-empty "resourceType" string, got: ${JSON.stringify(resource.resourceType)}`,
    );
  }
  const validResourceTypes = [
    "Observation", "Patient", "Device", "DiagnosticReport",
    "Specimen", "Bundle", "Practitioner", "Organization",
  ];
  if (!validResourceTypes.includes(resource.resourceType as string)) {
    throw new Error(
      `Unknown FHIR resourceType "${resource.resourceType}". Expected one of: ${validResourceTypes.join(", ")}`,
    );
  }
}

/** Biomarker normal ranges for basic validation */
const BIOMARKER_RANGES: Record<string, { min: number; max: number; unit: string }> = {
  glucose: { min: 70, max: 200, unit: "mg/dL" },
  cholesterol: { min: 100, max: 300, unit: "mg/dL" },
  hemoglobin: { min: 7, max: 20, unit: "g/dL" },
  creatinine: { min: 0.4, max: 5.0, unit: "mg/dL" },
  potassium: { min: 2.5, max: 7.0, unit: "mmol/L" },
  sodium: { min: 120, max: 160, unit: "mmol/L" },
  hba1c: { min: 3.0, max: 15.0, unit: "%" },
  troponin: { min: 0, max: 100, unit: "ng/L" },
};

/**
 * Asserts that a measurement value falls within the plausible range for a given biomarker.
 * Does NOT assert clinical normality — only that the reading is physically possible.
 */
export function expectMeasurementInRange(
  measurement: { value?: number; [key: string]: unknown },
  biomarker: string,
): void {
  const range = BIOMARKER_RANGES[biomarker.toLowerCase()];
  if (!range) {
    throw new Error(
      `Unknown biomarker "${biomarker}". Supported: ${Object.keys(BIOMARKER_RANGES).join(", ")}`,
    );
  }
  if (measurement.value === undefined || measurement.value === null) {
    throw new Error(`Measurement is missing a "value" property`);
  }
  if (typeof measurement.value !== "number" || Number.isNaN(measurement.value)) {
    throw new Error(`Measurement value must be a number, got: ${measurement.value}`);
  }
  if (measurement.value < range.min || measurement.value > range.max) {
    throw new Error(
      `Measurement value ${measurement.value} for "${biomarker}" is outside plausible range [${range.min}, ${range.max}] ${range.unit}`,
    );
  }
}
