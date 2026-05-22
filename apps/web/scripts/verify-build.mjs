#!/usr/bin/env node
/**
 * Phase BW-3 — apps/web (Next.js) 빌드 산출물 자가 검증 (Harness H3 / H5)
 *
 * Next.js 빌드 (`next build`) 이후 다음 계약을 검사:
 *
 *   1. .next/ 디렉토리 존재
 *   2. .next/standalone/ 존재 (output: 'standalone' 모드)
 *   3. .next/static/chunks 에 main bundle 존재 + 0 바이트 아님
 *   4. 소스 + 빌드 산출물 어디에도 placeholder sentinel 없음
 *      ([VERIFIED_*], [TEST_*], [REPLACEMENT_*], [SYNC_*], [PLACEHOLDER_*],
 *       TEST_REPLACEMENT_SUCCESS_LOCAL 등)
 *   5. middleware.ts manifest 가 정상 생성됨
 *
 * 실패 시 종료 코드 1 — CI 차단.
 *
 * 사용:
 *   pnpm build && node scripts/verify-build.mjs
 */

import { existsSync, readFileSync, readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';

const ROOT = new URL('..', import.meta.url).pathname;
const NEXT_DIR = join(ROOT, '.next');

const fails = [];
const warns = [];

function fail(msg) { fails.push(msg); }
function warn(msg) { warns.push(msg); }
function ok(msg) { console.log(`  ✅ ${msg}`); }

console.log('🔍 apps/web/.next/ 빌드 산출물 검증 시작\n');

// (1) .next 디렉토리
if (!existsSync(NEXT_DIR)) {
  fail(`.next/ 가 존재하지 않음 — 'pnpm build' 먼저 실행 필요`);
  printResult();
  process.exit(1);
}
ok('.next/ 디렉토리 존재');

// (2) standalone build
const standalonePath = join(NEXT_DIR, 'standalone');
if (existsSync(standalonePath)) {
  ok('standalone 빌드 활성 (next.config output: standalone)');
} else {
  warn('standalone 빌드 미활성 — next.config.mjs 의 output 설정 확인');
}

// (3) main bundle
const chunksPath = join(NEXT_DIR, 'static', 'chunks');
if (existsSync(chunksPath)) {
  const chunks = readdirSync(chunksPath).filter(f => /\.js$/.test(f));
  const main = chunks.find(f => /main-/.test(f)) ?? chunks[0];
  if (!main) {
    fail('.next/static/chunks 에 JS 파일 없음');
  } else {
    const size = statSync(join(chunksPath, main)).size;
    if (size === 0) {
      fail(`${main} 가 0 바이트`);
    } else {
      ok(`main bundle ${main} (${(size / 1024).toFixed(1)} KB)`);
    }
  }
} else {
  fail('.next/static/chunks/ 누락');
}

// (4) Sentinel 검사 — 소스 + 빌드 산출물 양쪽
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
  try {
    content = readFileSync(filepath, 'utf-8');
  } catch {
    return;
  }
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
      if (entry === 'node_modules' || entry === '.next' || entry.startsWith('.')) continue;
      scanDir(full, prefix, extensions);
    } else if (extensions.some(e => entry.endsWith(e))) {
      scanFile(full, prefix);
    }
  }
}

// 소스 (.tsx/.ts/.mdx)
for (const d of SOURCE_DIRS) {
  scanDir(join(ROOT, d), 'src', ['.tsx', '.ts', '.mdx', '.md']);
}
// 빌드 산출물 .next/server (SSR HTML 안에 sentinel 잔존 검증)
scanDir(join(NEXT_DIR, 'server'), 'build', ['.html', '.json', '.js']);

if (sentinelHits.length > 0) {
  for (const h of sentinelHits) fail(`Sentinel 검출: ${h}`);
} else {
  ok('Sentinel 패턴 미검출');
}

// (5) build-manifest.json — Next.js 빌드 정합
const manifestPath = join(NEXT_DIR, 'build-manifest.json');
if (existsSync(manifestPath)) {
  ok('build-manifest.json 정상');
} else {
  fail('.next/build-manifest.json 누락 — 빌드 비정상 종료');
}

// (6) middleware manifest
const middlewareManifest = join(NEXT_DIR, 'server', 'middleware-manifest.json');
if (existsSync(middlewareManifest)) {
  ok('middleware-manifest.json 정상 (intl + auth 활성)');
} else {
  warn('middleware-manifest.json 누락 — middleware.ts 비활성 가능성');
}

// ============================================================================
// (7) Deploy Contract — Cloudflare Pages 한계 검증 (Phase BX, Harness H7)
// ============================================================================
//
// Cloudflare Pages 공식 한계:
//   - 단일 파일: 25 MiB                       (hard limit)
//   - 총 파일 수: 20,000                       (soft, plan 별 상이)
//   - 총 deploy 크기: 25 GB                   (soft)
//   - Pages Functions 압축 크기: 10 MiB        (Workers)
//
// 빌드 성공 ≠ deploy 성공. 사고 회고 (2026-05-22T01:38): 27 MiB
// webpack cache 가 deploy 단계에서 차단되어 production 미반영.
// 본 검증은 deploy 차단 발생 전에 build 단계에서 사전 감지.

const CLOUDFLARE_MAX_FILE_BYTES = 25 * 1024 * 1024; // 25 MiB
const FILE_COUNT_WARN = 18000;
const FILE_COUNT_FAIL = 20000;
const TOTAL_SIZE_WARN_BYTES = 1024 * 1024 * 1024; // 1 GB
const TOTAL_SIZE_FAIL_BYTES = 25 * 1024 * 1024 * 1024; // 25 GB

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
      // .next/cache 는 deploy 에 포함되면 안 됨 — production webpack cache 잔존 검사
      if (rel === 'cache' && relative === '') {
        // 크기 계산만 하고 경고 처리
        const cacheSize = computeDirSize(full);
        cacheJunkBytes = cacheSize;
        continue; // walk 안 함 — fail 발생 방지
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

walkDeploy(NEXT_DIR);

// 단일 파일 크기 위반 — hard fail
if (oversized.length > 0) {
  for (const f of oversized) {
    fail(`단일 파일 크기 ${(f.bytes / 1024 / 1024).toFixed(1)} MiB > Cloudflare 한계 25 MiB: .next/${f.path}`);
  }
} else {
  ok(`모든 파일 ≤ 25 MiB (총 ${fileCount} 개)`);
}

// 파일 수 검사
if (fileCount > FILE_COUNT_FAIL) {
  fail(`파일 수 ${fileCount} > Cloudflare 한계 ${FILE_COUNT_FAIL}`);
} else if (fileCount > FILE_COUNT_WARN) {
  warn(`파일 수 ${fileCount} (한계 ${FILE_COUNT_FAIL} 에 근접)`);
}

// 총 크기 검사
if (totalBytes > TOTAL_SIZE_FAIL_BYTES) {
  fail(`총 deploy 크기 ${(totalBytes / 1024 / 1024 / 1024).toFixed(2)} GB > 25 GB 한계`);
} else if (totalBytes > TOTAL_SIZE_WARN_BYTES) {
  warn(`총 deploy 크기 ${(totalBytes / 1024 / 1024).toFixed(0)} MiB (1 GB 초과 — 분석 권장)`);
} else {
  ok(`총 deploy 크기 ${(totalBytes / 1024 / 1024).toFixed(1)} MiB`);
}

// .next/cache 잔존 검사 — webpack cache 파일이 27 MiB+ 로 deploy 차단 원인
if (cacheJunkBytes > 0) {
  const mib = (cacheJunkBytes / 1024 / 1024).toFixed(1);
  fail(`.next/cache 잔존 (${mib} MiB) — package.json build 가 rm -rf .next/cache 실행해야 함`);
} else {
  ok('.next/cache 미존재 (deploy 정합)');
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
    console.log('\nNext.js 빌드 산출물이 계약을 만족하지 않습니다. 배포 차단을 권장합니다.');
    process.exit(1);
  }
  console.log(`\n✅ 빌드 산출물 검증 통과 (경고 ${warns.length})`);
  process.exit(0);
}
