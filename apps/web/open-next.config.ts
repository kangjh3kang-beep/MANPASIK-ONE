// @ts-expect-error - @opennextjs/cloudflare provides this at build time
import { defineCloudflareConfig } from '@opennextjs/cloudflare';

/**
 * Phase BY — OpenNext + Cloudflare Workers 구성.
 *
 * Vercel CLI 의존을 완전 제거하고 Cloudflare 가 공식 권장하는 OpenNext 어댑터로
 * Next.js 15 빌드 결과를 Workers 진입점(.open-next/worker.js)으로 변환한다.
 *
 * 기본 설정만 사용 — incremental cache, queue, tag cache 는 추후 KV/R2/DO 와 연동
 * 필요 시 단계적 활성화 (Hyperdrive, Smart Placement 도 검토 대상).
 */
export default defineCloudflareConfig({
  // 추후 활성화 옵션 (Phase BZ+):
  // incrementalCache: 'kv',   // KV namespace 바인딩 필요
  // queue: 'durable-objects', // DO 바인딩 필요
  // tagCache: 'd1',           // D1 데이터베이스 필요
});
