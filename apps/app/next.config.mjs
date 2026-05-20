import createNextIntlPlugin from 'next-intl/plugin';
const withNextIntl = createNextIntlPlugin('./i18n/request.ts');

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@mmup/ui', '@mmup/design-system', '@mmup/ssot', '@mmup/types'],
};
export default withNextIntl(nextConfig);
