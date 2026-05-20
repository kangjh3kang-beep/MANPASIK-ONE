import './globals.css';

export const metadata = {
  title: 'MMUP 파트너 통합 연동',
  description: 'HL7/FHIR 국제 표준 의료 데이터 파이프라인',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ko" suppressHydrationWarning>
      <body suppressHydrationWarning className="antialiased">
        {children}
      </body>
    </html>
  );
}
