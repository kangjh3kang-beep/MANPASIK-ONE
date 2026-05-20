import createNextIntlPlugin from 'next-intl/plugin';
const withNextIntl = createNextIntlPlugin('./i18n/request.ts');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  transpilePackages: ['@mmup/ui', '@mmup/design-system', '@mmup/ssot', '@mmup/types'],
  experimental: {
    instrumentationHook: true,
  },
  env: {
    FORCE_RELOAD: Date.now().toString(),
  }
};
export default withNextIntl(nextConfig);
