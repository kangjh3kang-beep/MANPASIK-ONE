/**
 * API Endpoint Verification Test
 *
 * Verifies consistency between:
 *   1. Mock endpoint paths in the catch-all route (app/api/v1/[...path]/route.ts)
 *   2. API hooks in packages/api-client/src/hooks.ts
 *   3. MSW handlers in packages/api-client/src/mocks/handlers.ts
 *
 * Run: npx tsx tests/api-endpoints.test.ts
 */

import * as fs from 'fs';
import * as path from 'path';

// ---------------------------------------------------------------------------
// File paths
// ---------------------------------------------------------------------------

const ROUTE_FILE = path.resolve(
  __dirname,
  '../app/api/v1/[...path]/route.ts',
);

const HOOKS_FILE = path.resolve(
  __dirname,
  '../../../packages/api-client/src/hooks.ts',
);

const HANDLERS_FILE = path.resolve(
  __dirname,
  '../../../packages/api-client/src/mocks/handlers.ts',
);

const ENDPOINTS_FILE = path.resolve(
  __dirname,
  '../../../packages/api-client/src/endpoints.ts',
);

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

let passed = 0;
let failed = 0;

function pass(msg: string) {
  console.log(`  [PASS] ${msg}`);
  passed++;
}

function fail(msg: string) {
  console.log(`  [FAIL] ${msg}`);
  failed++;
}

function readFileChecked(filePath: string): string {
  if (!fs.existsSync(filePath)) {
    throw new Error(`File not found: ${filePath}`);
  }
  return fs.readFileSync(filePath, 'utf-8');
}

// ---------------------------------------------------------------------------
// Step 1: Extract mock endpoint paths from catch-all route
// ---------------------------------------------------------------------------

console.log('=== API Endpoint Verification Test ===\n');
console.log('--- Step 1: Extract Mock Endpoints from Route ---');

const routeSource = readFileChecked(ROUTE_FILE);

// Match MOCK_DATA keys like 'patients', 'patients/*/vitals', etc.
const mockPathRegex = /^\s*'([^']+)'/gm;
const mockPaths: string[] = [];
let match: RegExpExecArray | null;

// Only extract from MOCK_DATA block
const mockDataBlock = routeSource.match(/const MOCK_DATA[^}]*\{([\s\S]*?)^};/m);
if (mockDataBlock) {
  const block = mockDataBlock[0];
  while ((match = mockPathRegex.exec(block)) !== null) {
    mockPaths.push(match[1]);
  }
}

console.log(`  Found ${mockPaths.length} mock endpoint paths:`);
for (const p of mockPaths) {
  console.log(`    /api/v1/${p}`);
}

if (mockPaths.length > 0) {
  pass(`Route file contains ${mockPaths.length} mock endpoints`);
} else {
  fail('No mock endpoints found in route file');
}

// ---------------------------------------------------------------------------
// Step 2: Extract hooks from hooks.ts
// ---------------------------------------------------------------------------

console.log('\n--- Step 2: Extract Hooks ---');

const hooksSource = readFileChecked(HOOKS_FILE);

// Match exported function names like usePatients, useSensors, etc.
const hookRegex = /export function (use\w+)/g;
const hooks: string[] = [];
while ((match = hookRegex.exec(hooksSource)) !== null) {
  hooks.push(match[1]);
}

console.log(`  Found ${hooks.length} hooks:`);
for (const h of hooks) {
  console.log(`    ${h}`);
}

if (hooks.length > 0) {
  pass(`hooks.ts exports ${hooks.length} hooks`);
} else {
  fail('No hooks found in hooks.ts');
}

// ---------------------------------------------------------------------------
// Step 3: Extract endpoint paths from endpoints.ts
// ---------------------------------------------------------------------------

console.log('\n--- Step 3: Extract Endpoint Paths ---');

const endpointsSource = readFileChecked(ENDPOINTS_FILE);

// Match API paths like '/api/v1/patients', '/api/v1/sensors', etc.
const endpointPathRegex = /['"]\/api\/v1\/([^'"$]+)['"]/g;
const endpointPaths: string[] = [];
while ((match = endpointPathRegex.exec(endpointsSource)) !== null) {
  // Normalize template expressions like `patients/${patientId}/vitals` -> `patients/*/vitals`
  const normalized = match[1].replace(/\$\{[^}]+\}/g, '*');
  if (!endpointPaths.includes(normalized)) {
    endpointPaths.push(normalized);
  }
}

console.log(`  Found ${endpointPaths.length} unique endpoint paths in endpoints.ts`);

// ---------------------------------------------------------------------------
// Step 4: Verify each endpoint path has a mock
// ---------------------------------------------------------------------------

console.log('\n--- Step 4: Verify Endpoints Have Mock Data ---');

// Normalize mock paths for comparison (wildcards)
function normalizePath(p: string): string {
  return p.replace(/\*/g, '*').replace(/:[a-zA-Z]+/g, '*');
}

const normalizedMockPaths = mockPaths.map(normalizePath);

for (const ep of endpointPaths) {
  const normalized = normalizePath(ep);
  const hasMock = normalizedMockPaths.some(mp => {
    // Direct match
    if (mp === normalized) return true;
    // Wildcard matching: mock pattern matches endpoint
    const regex = new RegExp('^' + mp.replace(/\*/g, '[^/]+') + '$');
    return regex.test(normalized);
  });

  if (hasMock) {
    pass(`/api/v1/${ep} has a mock endpoint`);
  } else {
    fail(`/api/v1/${ep} has NO mock endpoint in route.ts`);
  }
}

// ---------------------------------------------------------------------------
// Step 5: Extract MSW handlers
// ---------------------------------------------------------------------------

console.log('\n--- Step 5: Extract MSW Handlers ---');

const handlersSource = readFileChecked(HANDLERS_FILE);

// Match patterns like http.get('*/api/v1/patients', ...)
const mswRegex = /http\.(get|post|put|delete|patch)\(\s*['"].*?\/api\/v1\/([^'"]+)['"]/g;
const mswPaths: { method: string; path: string }[] = [];
while ((match = mswRegex.exec(handlersSource)) !== null) {
  const method = match[1].toUpperCase();
  const rawPath = match[2].replace(/:[a-zA-Z]+/g, '*');
  mswPaths.push({ method, path: rawPath });
}

console.log(`  Found ${mswPaths.length} MSW handlers:`);
for (const h of mswPaths) {
  console.log(`    ${h.method} /api/v1/${h.path}`);
}

if (mswPaths.length > 0) {
  pass(`handlers.ts contains ${mswPaths.length} MSW handlers`);
} else {
  fail('No MSW handlers found');
}

// ---------------------------------------------------------------------------
// Step 6: Verify each endpoint has a corresponding MSW handler
// ---------------------------------------------------------------------------

console.log('\n--- Step 6: Verify Endpoints Have MSW Handlers ---');

const normalizedMswPaths = mswPaths.map(h => normalizePath(h.path));

for (const ep of endpointPaths) {
  const normalized = normalizePath(ep);
  const hasHandler = normalizedMswPaths.some(mp => {
    if (mp === normalized) return true;
    const regex = new RegExp('^' + mp.replace(/\*/g, '[^/]+') + '$');
    return regex.test(normalized);
  });

  if (hasHandler) {
    pass(`/api/v1/${ep} has an MSW handler`);
  } else {
    fail(`/api/v1/${ep} has NO MSW handler in handlers.ts`);
  }
}

// ---------------------------------------------------------------------------
// Step 7: Hook-to-Endpoint mapping verification
// ---------------------------------------------------------------------------

console.log('\n--- Step 7: Hook-to-Endpoint Mapping ---');

// Map hook names to expected API category prefixes
const HOOK_TO_CATEGORY: Record<string, string[]> = {
  usePatients: ['patients'],
  usePatient: ['patients'],
  usePatientVitals: ['patients/*/vitals'],
  usePatientBiomarkers: ['patients/*/biomarkers'],
  useSensors: ['sensors'],
  useAgents: ['agents'],
  useRunPrediction: ['agents/*/predict'],
  usePredictions: ['predictor/risk-scores', 'predictions/patient/*'],
  useRewardBalance: ['rewards/*/balance'],
  useRewardContributions: ['rewards/*/contributions'],
  usePartners: ['partners'],
  useAuditLogs: ['gxp/audit-logs'],
  useCompliance: ['gxp/compliance'],
  useAgentPipeline: ['agents/pipeline'],
  useDevPortalKeys: ['dev-portal/keys'],
  useDevPortalEndpoints: ['dev-portal/endpoints'],
  useHardwareDevices: ['hardware/devices'],
  useHardwareDiagnostics: ['hardware/diagnostics'],
  useAppMetrics: ['app/metrics'],
  useMeasurementHistory: ['measurements/history'],
  useTelemedicineDoctors: ['telemedicine/doctors'],
  useStoreProducts: ['store/products'],
  useAdminKPIs: ['admin/kpis'],
  useAdminEvents: ['admin/events'],
};

for (const hook of hooks) {
  const expectedPaths = HOOK_TO_CATEGORY[hook];
  if (!expectedPaths) {
    fail(`${hook}: no expected endpoint mapping defined (update test)`);
    continue;
  }

  // Check if at least one expected path has a mock
  const hasMock = expectedPaths.some(ep => {
    const normalized = normalizePath(ep);
    return normalizedMockPaths.some(mp => {
      if (mp === normalized) return true;
      const regex = new RegExp('^' + mp.replace(/\*/g, '[^/]+') + '$');
      return regex.test(normalized);
    });
  });

  // Check if at least one expected path has an MSW handler
  const hasHandler = expectedPaths.some(ep => {
    const normalized = normalizePath(ep);
    return normalizedMswPaths.some(mp => {
      if (mp === normalized) return true;
      const regex = new RegExp('^' + mp.replace(/\*/g, '[^/]+') + '$');
      return regex.test(normalized);
    });
  });

  if (hasMock && hasHandler) {
    pass(`${hook} -> mock + MSW handler`);
  } else if (hasMock) {
    fail(`${hook} -> has mock but NO MSW handler`);
  } else if (hasHandler) {
    fail(`${hook} -> has MSW handler but NO mock`);
  } else {
    fail(`${hook} -> missing BOTH mock and MSW handler`);
  }
}

// ---------------------------------------------------------------------------
// Summary
// ---------------------------------------------------------------------------

console.log('\n=== Summary ===');
console.log(`  Passed: ${passed}`);
console.log(`  Failed: ${failed}`);
console.log(`  Result: ${failed === 0 ? 'ALL PASS' : 'SOME FAILURES'}`);
process.exit(failed > 0 ? 1 : 0);
