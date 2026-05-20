import createMiddleware from 'next-intl/middleware';
 
export default createMiddleware({
  locales: ['ko', 'en', 'ja', 'zh'],
  defaultLocale: 'ko'
});
 
export const config = {
  matcher: ['/', '/(ko|en|ja|zh)/:path*']
};
