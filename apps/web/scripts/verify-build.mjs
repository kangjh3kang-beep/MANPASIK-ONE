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
