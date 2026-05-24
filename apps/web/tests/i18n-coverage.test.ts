/**
 * i18n Coverage Test
 *
 * Reads all 4 message files (ko, en, ja, zh) and verifies
 * that every key present in ko.json also exists in the other
 * three locale files. Reports missing keys per language.
 *
 * Run: npx tsx tests/i18n-coverage.test.ts
 */

import * as fs from 'fs';
import * as path from 'path';

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const MESSAGES_DIR = path.resolve(__dirname, '../messages');
const BASE_LOCALE = 'ko';
const TARGET_LOCALES = ['en', 'ja', 'zh'];
const ALL_LOCALES = [BASE_LOCALE, ...TARGET_LOCALES];

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

/**
 * Recursively collect all dot-separated key paths from a nested object.
 */
function collectKeys(obj: Record<string, unknown>, prefix = ''): string[] {
  const keys: string[] = [];
  for (const [key, value] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefix}.${key}` : key;
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      keys.push(...collectKeys(value as Record<string, unknown>, fullKey));
    } else {
      keys.push(fullKey);
    }
  }
  return keys;
}

/**
 * Check if a dot-separated key exists in a nested object.
 */
function hasKey(obj: Record<string, unknown>, dotPath: string): boolean {
  const parts = dotPath.split('.');
  let current: unknown = obj;
  for (const part of parts) {
    if (current === null || current === undefined || typeof current !== 'object') {
      return false;
    }
    current = (current as Record<string, unknown>)[part];
  }
  return current !== undefined;
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

console.log('=== i18n Coverage Test ===\n');

// Step 1: Load all message files
console.log('--- Step 1: Load Message Files ---');

const messages: Record<string, Record<string, unknown>> = {};

for (const locale of ALL_LOCALES) {
  const filePath = path.join(MESSAGES_DIR, `${locale}.json`);
  if (!fs.existsSync(filePath)) {
    fail(`${locale}.json not found at ${filePath}`);
    continue;
  }
  const content = fs.readFileSync(filePath, 'utf-8');
  messages[locale] = JSON.parse(content);
  pass(`${locale}.json loaded`);
}

if (!(BASE_LOCALE in messages)) {
  console.error(`  FATAL: Base locale '${BASE_LOCALE}' not found. Cannot continue.`);
  process.exit(1);
}

// Step 2: Collect all keys from base locale
console.log('\n--- Step 2: Collect Keys from Base Locale (ko) ---');

const baseKeys = collectKeys(messages[BASE_LOCALE]);
console.log(`  Total keys in ${BASE_LOCALE}.json: ${baseKeys.length}`);

// Group keys by top-level section
const sections = new Map<string, number>();
for (const key of baseKeys) {
  const section = key.split('.')[0];
  sections.set(section, (sections.get(section) || 0) + 1);
}
console.log('  Sections:');
for (const [section, count] of sections) {
  console.log(`    ${section}: ${count} keys`);
}

pass(`Base locale has ${baseKeys.length} keys across ${sections.size} sections`);

// Step 3: Compare each target locale against the base
console.log('\n--- Step 3: Compare Keys ---');

for (const locale of TARGET_LOCALES) {
  if (!(locale in messages)) {
    fail(`Skipping ${locale} (file not loaded)`);
    continue;
  }

  console.log(`\n  Checking ${locale}.json...`);

  const targetKeys = collectKeys(messages[locale]);
  const missingKeys: string[] = [];
  const extraKeys: string[] = [];

  // Find keys in base that are missing in target
  for (const key of baseKeys) {
    if (!hasKey(messages[locale], key)) {
      missingKeys.push(key);
    }
  }

  // Find keys in target that are not in base (informational)
  for (const key of targetKeys) {
    if (!hasKey(messages[BASE_LOCALE], key)) {
      extraKeys.push(key);
    }
  }

  if (missingKeys.length === 0) {
    pass(`${locale}.json: all ${baseKeys.length} keys present`);
  } else {
    fail(`${locale}.json: ${missingKeys.length} keys missing`);
    for (const key of missingKeys) {
      console.log(`    MISSING: ${key}`);
    }
  }

  if (extraKeys.length > 0) {
    console.log(`  [INFO] ${locale}.json has ${extraKeys.length} extra keys not in base:`);
    for (const key of extraKeys) {
      console.log(`    EXTRA: ${key}`);
    }
  }

  // Check key count match
  if (targetKeys.length === baseKeys.length) {
    pass(`${locale}.json: key count matches (${targetKeys.length})`);
  } else {
    const diff = targetKeys.length - baseKeys.length;
    const diffStr = diff > 0 ? `+${diff}` : `${diff}`;
    fail(`${locale}.json: key count mismatch (${targetKeys.length} vs ${baseKeys.length}, ${diffStr})`);
  }
}

// Step 4: Check for empty values in any locale
console.log('\n--- Step 4: Check for Empty Values ---');

for (const locale of ALL_LOCALES) {
  if (!(locale in messages)) continue;

  const keys = collectKeys(messages[locale]);
  const emptyKeys: string[] = [];

  for (const key of keys) {
    const parts = key.split('.');
    let current: unknown = messages[locale];
    for (const part of parts) {
      current = (current as Record<string, unknown>)[part];
    }
    if (current === '' || current === null) {
      emptyKeys.push(key);
    }
  }

  if (emptyKeys.length === 0) {
    pass(`${locale}.json: no empty values`);
  } else {
    fail(`${locale}.json: ${emptyKeys.length} empty values`);
    for (const key of emptyKeys) {
      console.log(`    EMPTY: ${key}`);
    }
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
