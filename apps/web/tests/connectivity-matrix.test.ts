/**
 * Connectivity Matrix Test
 *
 * Reads all 16 domain page files, extracts outbound links from
 * "관련 도메인" sections, builds a connectivity matrix, and verifies:
 *   1. Every page has >= 3 outbound links
 *   2. No orphan pages (every page has >= 1 inbound link)
 *   3. Patient journey path is complete:
 *      measure -> health-records -> ai-coach -> telemedicine -> store
 *
 * Run: npx tsx tests/connectivity-matrix.test.ts
 */

import * as fs from 'fs';
import * as path from 'path';

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const DOMAINS_DIR = path.resolve(
  __dirname,
  '../app/[locale]/domains',
);

const ALL_DOMAINS = [
  'admin',
  'agents-hub',
  'ai-coach',
  'app',
  'clinical',
  'dev-portal',
  'gxp',
  'hardware-core',
  'health-records',
  'measure',
  'partner',
  'predictor',
  'reward',
  'settings',
  'store',
  'telemedicine',
];

// The required patient journey path (each element must link to the next)
const PATIENT_JOURNEY = [
  'measure',
  'health-records',
  'ai-coach',
  'telemedicine',
  'store',
];

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function extractOutboundLinks(source: string): string[] {
  const links = new Set<string>();

  // Pattern 1: path: '/domains/xxx' (관련 도메인 sections)
  const pathRegex = /path:\s*['"]\/domains\/([a-z-]+)['"]/g;
  let match: RegExpExecArray | null;
  while ((match = pathRegex.exec(source)) !== null) {
    links.add(match[1]);
  }

  // Pattern 2: href={`${localePrefix}/domains/xxx`} (app page style)
  const hrefRegex = /\/domains\/([a-z-]+)/g;
  while ((match = hrefRegex.exec(source)) !== null) {
    links.add(match[1]);
  }

  return [...links];
}

// ---------------------------------------------------------------------------
// Main
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

console.log('=== Connectivity Matrix Test ===\n');

// Step 1: Read all domain page files and extract links
const matrix: Record<string, string[]> = {};
const missingFiles: string[] = [];

for (const domain of ALL_DOMAINS) {
  const pagePath = path.join(DOMAINS_DIR, domain, 'page.tsx');
  if (!fs.existsSync(pagePath)) {
    missingFiles.push(domain);
    matrix[domain] = [];
    continue;
  }
  const source = fs.readFileSync(pagePath, 'utf-8');
  const links = extractOutboundLinks(source).filter(l => l !== domain); // exclude self-links
  matrix[domain] = links;
}

if (missingFiles.length > 0) {
  fail(`Missing page files for: ${missingFiles.join(', ')}`);
} else {
  pass(`All 16 domain page files found`);
}

// Print the matrix
console.log('\n--- Connectivity Matrix ---');
for (const [domain, links] of Object.entries(matrix)) {
  console.log(`  ${domain.padEnd(18)} -> ${links.length > 0 ? links.join(', ') : '(none)'}`);
}

// Step 2: Verify minimum outbound links (>= 3)
console.log('\n--- Check: Minimum 3 Outbound Links ---');
for (const domain of ALL_DOMAINS) {
  const links = matrix[domain] || [];
  if (links.length >= 3) {
    pass(`${domain}: ${links.length} outbound links`);
  } else {
    fail(`${domain}: only ${links.length} outbound links (need >= 3)`);
  }
}

// Step 3: Verify no orphan pages (every page has >= 1 inbound link)
console.log('\n--- Check: No Orphan Pages (>= 1 Inbound) ---');
const inboundCounts: Record<string, number> = {};
for (const domain of ALL_DOMAINS) {
  inboundCounts[domain] = 0;
}
for (const [, links] of Object.entries(matrix)) {
  for (const target of links) {
    if (target in inboundCounts) {
      inboundCounts[target]++;
    }
  }
}

for (const domain of ALL_DOMAINS) {
  if (inboundCounts[domain] >= 1) {
    pass(`${domain}: ${inboundCounts[domain]} inbound links`);
  } else {
    fail(`${domain}: 0 inbound links (orphan page!)`);
  }
}

// Step 4: Verify patient journey path
console.log('\n--- Check: Patient Journey Path ---');
console.log(`  Path: ${PATIENT_JOURNEY.join(' -> ')}`);
let journeyComplete = true;
for (let i = 0; i < PATIENT_JOURNEY.length - 1; i++) {
  const from = PATIENT_JOURNEY[i];
  const to = PATIENT_JOURNEY[i + 1];
  const links = matrix[from] || [];
  if (links.includes(to)) {
    pass(`${from} -> ${to}`);
  } else {
    fail(`${from} -> ${to} (link missing!)`);
    journeyComplete = false;
  }
}

if (journeyComplete) {
  pass('Patient journey path is complete');
}

// Summary
console.log('\n=== Summary ===');
console.log(`  Passed: ${passed}`);
console.log(`  Failed: ${failed}`);
console.log(`  Result: ${failed === 0 ? 'ALL PASS' : 'SOME FAILURES'}`);
process.exit(failed > 0 ? 1 : 0);
