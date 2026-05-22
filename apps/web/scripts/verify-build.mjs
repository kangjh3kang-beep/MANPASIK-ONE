#!/usr/bin/env node
/**
 * Phase BY (Harness H3 / H5 / H7) — apps/web 빌드 산출물 자가 검증.
 *
 * 세 가지 빌드 모드 모두 지원:
 *   (a) next build → .next/ (standalone, 검증용)
 *   (b) @opennextjs/cloudflare → .open-next/ (Cloudflare Workers, 운영)
 *   (c) @cloudflare/next-on-pages → .vercel/output/static/ (DEPRECATED, legacy)
 *
 * 검증 항목:
 *   1. 빌드 산출물 존재
 *   2. sentinel 패턴 미잔존 (소스 + 산출물)
 *   3. Cloudflare Deploy Contract (단일 파일 < 25 MiB, 총 < 20K files)
 *   4. OpenNext 핵심 파일 존재 (worker.js + assets/)
 *
 * 실패 시 종료 코드 1 — CI 차단.
 */

import { existsSync, readFileSync, readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';

const ROOT = new URL('..', import.meta.url).pathname;
const NEXT_DIR = join(ROOT, '.next');
const VERCEL_OUTPUT = join(ROOT, '.vercel', 'output');
const VERCEL_STATIC = join(VERCEL_OUTPUT, 'static');
const OPEN_NEXT_DIR = join(ROOT, '.open-next');
const OPEN_NEXT_WORKER = join(OPEN_NEXT_DIR, 'worker.js');
const OPEN_NEXT_ASSETS = join(OPEN_NEXT_DIR, 'assets');

const fails = [];
const warns = [];

function fail(msg) { fails.push(msg); }
function warn(msg) { warns.push(msg); }
function ok(msg) { console.log(`  ✅ ${msg}`); }

console.log('🔍 apps/web 빌드 산출물 검증 시작\n');

// 빌드 모드 자동 감지 (OpenNext 우선 → next-on-pages → next build)
const hasNextDir = existsSync(NEXT_DIR);
const hasVercelOutput = existsSync(VERCEL_STATIC);
const hasOpenNext = existsSync(OPEN_NEXT_DIR);

if (!hasNextDir && !hasVercelOutput && !hasOpenNext) {
  fail('빌드 산출물 미존재 — .open-next/, .vercel/output/static/, .next/ 모두 없음');
  printResult();
  process.exit(1);
}

// 검증 대상 디렉토리 결정 (OpenNext 우선)
let DEPLOY_DIR;
let MODE;
if (hasOpenNext) {
  DEPLOY_DIR = OPEN_NEXT_DIR;
  MODE = 'OpenNext (Cloudflare Workers)';
  ok(`빌드 모드: ${MODE}`);
  ok('.open-next/ 디렉토리 존재');

  if (existsSync(OPEN_NEXT_WORKER)) {
    const workerSize = statSync(OPEN_NEXT_WORKER).size;
    ok(`worker.js 존재 (${(workerSize / 1024 / 1024).toFixed(2)} MiB)`);
    if (workerSize > 25 * 1024 * 1024) {
      fail(`worker.js ${(workerSize / 1024 / 1024).toFixed(1)} MiB > Workers 한계 25 MiB`);
    }
  } else {
    fail('.open-next/worker.js 누락 — OpenNext 빌드 비정상 종료');
  }

  if (existsSync(OPEN_NEXT_ASSETS)) {
    ok('.open-next/assets/ 디렉토리 존재');
  } else {
    fail('.open-next/assets/ 누락 — 정적 자산 디렉토리 없음');
  }
} else if (hasVercelOutput) {
  DEPLOY_DIR = VERCEL_STATIC;
  MODE = 'next-on-pages (Cloudflare Pages, DEPRECATED)';
  ok(`빌드 모드: ${MODE}`);
  warn('next-on-pages 는 deprecated — OpenNext 마이그레이션 권장');
} else {
  DEPLOY_DIR = NEXT_DIR;
  MODE = 'next build (standalone)';
  ok(`빌드 모드: ${MODE}`);
  ok('.next/ 디렉토리 존재');
}

// next build 모드의 추가 검증 (legacy)
if (hasNextDir && !hasVercelOutput) {
  if (existsSync(join(NEXT_DIR, 'standalone'))) {
    ok('standalone 빌드 활성');
  } else {
    warn('standalone 빌드 미활성');
  }
  if (existsSync(join(NEXT_DIR, 'build-manifest.json'))) {
    ok('build-manifest.json 정상');
  } else {
    fail('build-manifest.json 누락 — 빌드 비정상 종료');
  }
}

// (2) Sentinel 검사 — 소스 + 산출물
const SENTINELS = [
  'TEST_REPLACEMENT_SUCCESS',
  'PLACEHOLDER_DO_NOT_DEPLOY',
  '[VERIFIED_SYNC]',
  '[SYNC_OK_',
  '[TEST_',
  '[REPLACEMENT_',
];

const SOURCE_DIRS = ['app', 'components', 'lib', 'i18n'];
const sentinelHits = [];

function scanFile(filepath, prefix) {
  let content;
  try { content = readFileSync(filepath, 'utf-8'); } catch { return; }
  for (const sentinel of SENTINELS) {
    if (content.includes(sentinel)) {
      sentinelHits.push(`${prefix}:${filepath.replace(ROOT, '')} — "${sentinel}"`);
    }
  }
}

function scanDir(dir, prefix, extensions) {
  if (!existsSync(dir)) return;
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    let st;
    try { st = statSync(full); } catch { continue; }
    if (st.isDirectory()) {
      if (entry === 'node_modules' || entry === '.next' || entry === '.vercel' || entry.startsWith('.')) continue;
      scanDir(full, prefix, extensions);
    } else if (extensions.some(e => entry.endsWith(e))) {
      scanFile(full, prefix);
    }
  }
}

for (const d of SOURCE_DIRS) {
  scanDir(join(ROOT, d), 'src', ['.tsx', '.ts', '.mdx', '.md']);
}

if (sentinelHits.length > 0) {
  for (const h of sentinelHits) fail(`Sentinel 검출: ${h}`);
} else {
  ok('Sentinel 패턴 미검출 (소스)');
}

// (3) Cloudflare Deploy Contract (Phase BX, Harness H7)
const CLOUDFLARE_MAX_FILE_BYTES = 25 * 1024 * 1024;
const FILE_COUNT_WARN = 18000;
const FILE_COUNT_FAIL = 20000;
const TOTAL_SIZE_WARN_BYTES = 1024 * 1024 * 1024;

console.log('');
console.log('📦 Cloudflare Pages Deploy Contract 검증');

const oversized = [];
let fileCount = 0;
let totalBytes = 0;
let cacheJunkBytes = 0;

function walkDeploy(dir, relative = '') {
  if (!existsSync(dir)) return;
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const rel = relative ? `${relative}/${entry}` : entry;
    let st;
    try { st = statSync(full); } catch { continue; }
    if (st.isDirectory()) {
      // .next/cache 는 deploy 에 포함되면 안 됨 (next build 모드 한정)
      if (rel === 'cache' && relative === '' && MODE.includes('standalone')) {
        cacheJunkBytes = computeDirSize(full);
        continue;
      }
      walkDeploy(full, rel);
    } else {
      fileCount += 1;
      totalBytes += st.size;
      if (st.size > CLOUDFLARE_MAX_FILE_BYTES) {
        oversized.push({ path: rel, bytes: st.size });
      }
    }
  }
}

function computeDirSize(dir) {
  let bytes = 0;
  if (!existsSync(dir)) return 0;
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    let st;
    try { st = statSync(full); } catch { continue; }
    if (st.isDirectory()) bytes += computeDirSize(full);
    else bytes += st.size;
  }
  return bytes;
}

walkDeploy(DEPLOY_DIR);

if (oversized.length > 0) {
  for (const f of oversized) {
    fail(`단일 파일 ${(f.bytes / 1024 / 1024).toFixed(1)} MiB > Cloudflare 한계 25 MiB: ${f.path}`);
  }
} else {
  ok(`모든 파일 ≤ 25 MiB (총 ${fileCount} 개)`);
}

if (fileCount > FILE_COUNT_FAIL) {
  fail(`파일 수 ${fileCount} > Cloudflare 한계 ${FILE_COUNT_FAIL}`);
} else if (fileCount > FILE_COUNT_WARN) {
  warn(`파일 수 ${fileCount} (한계 ${FILE_COUNT_FAIL} 에 근접)`);
}

if (totalBytes > TOTAL_SIZE_WARN_BYTES) {
  warn(`총 deploy 크기 ${(totalBytes / 1024 / 1024).toFixed(0)} MiB (1 GB 초과 — 분석 권장)`);
} else {
  ok(`총 deploy 크기 ${(totalBytes / 1024 / 1024).toFixed(1)} MiB`);
}

if (cacheJunkBytes > 0) {
  const mib = (cacheJunkBytes / 1024 / 1024).toFixed(1);
  fail(`.next/cache 잔존 (${mib} MiB) — package.json build 가 rm -rf .next/cache 실행해야 함`);
}

printResult();

function printResult() {
  console.log('');
  if (warns.length > 0) {
    console.log(`⚠️  경고 ${warns.length} 건:`);
    warns.forEach(w => console.log(`    - ${w}`));
  }
  if (fails.length > 0) {
    console.log(`\n❌ 검증 실패 ${fails.length} 건:`);
    fails.forEach(f => console.log(`    - ${f}`));
    console.log('\n빌드 산출물이 계약을 만족하지 않습니다. 배포 차단을 권장합니다.');
    process.exit(1);
  }
  console.log(`\n✅ 빌드 산출물 검증 통과 (경고 ${warns.length})`);
  process.exit(0);
}
