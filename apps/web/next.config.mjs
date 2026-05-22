import createNextIntlPlugin from 'next-intl/plugin';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const withNextIntl = createNextIntlPlugin('./i18n/request.ts');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  outputFileTracingRoot: path.join(__dirname, '../../'),
  transpilePackages: ['@mmup/ui', '@mmup/design-system', '@mmup/ssot', '@mmup/types'],
  env: {
    FORCE_RELOAD: Date.now().toString(),
  },
  // Cloudflare Pages 는 단일 파일 25 MiB 제한이 있다.
  // production 빌드에서 webpack persistent cache (`.next/cache/webpack/*.pack`)
  // 가 27 MiB+ 로 생성되어 deploy 단계에서 "Pages only supports files up to
  // 25 MiB in size" 에러로 차단되는 사고 발생. (Cloudflare 빌드 로그 2026-05-22T01:38)
  //
  // 대응 2단:
  //   (1) production 빌드에서 webpack persistent cache 비활성 — .pack 미생성
  //   (2) build script 가 빌드 후 .next/cache 잔여물 제거 (이중 안전망)
  webpack: (config, { dev }) => {
    if (!dev) {
      config.cache = false;
    }
    return config;
  },
};
export default withNextIntl(nextConfig);
