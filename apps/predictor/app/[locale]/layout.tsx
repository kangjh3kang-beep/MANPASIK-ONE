/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3
 * @sb SB-1
 * @standard IEC 62304 Class B
 */
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';
import { TrinityIndicator } from '@mmup/ui';

export default async function LocaleLayout({
  children,
  params
}: {
  children: React.ReactNode;
  params: Promise<{locale: string}>;
}) {
  const { locale } = await params;
  const messages = await getMessages();
 
  return (
    

    <NextIntlClientProvider messages={messages}>
      {children}
      <TrinityIndicator />
    </NextIntlClientProvider>


  );
}
