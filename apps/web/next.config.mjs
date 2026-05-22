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
  }
};
export default withNextIntl(nextConfig);
