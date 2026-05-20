import './globals.css';

export const metadata = {
  title: 'MMUP 의약품 GxP 준수',
  description: 'cGMP 추적성 확보 및 전자 서명 워크플로우',
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
