import createNextIntlPlugin from 'next-intl/plugin';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const withNextIntl = createNextIntlPlugin('./i18n/request.ts');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // OpenNext (@opennextjs/cloudflare) 는 .open-next/worker.js 를 자체 생성하므로
  // output: 'standalone' 불필요. OpenNext 가 Next.js 빌드 결과를 변환.
  outputFileTracingRoot: path.join(__dirname, '../../'),
  transpilePackages: ['@mmup/ui', '@mmup/design-system', '@mmup/ssot', '@mmup/types'],
  env: {
    FORCE_RELOAD: Date.now().toString(),
  },
  // Cloudflare Workers 도 단일 파일 25 MiB 제한 유지. webpack persistent cache
  // (`.next/cache/webpack/*.pack`) 가 OpenNext 변환 과정에서 worker.js 로 부풀 위험
  // 있으므로 production 에서 비활성. (Phase BX 사고 정합)
  webpack: (config, { dev }) => {
    if (!dev) {
      config.cache = false;
    }
    return config;
  },
};
export default withNextIntl(nextConfig);
