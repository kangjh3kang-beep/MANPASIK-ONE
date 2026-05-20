import './globals.css';

export const metadata = {
  title: 'MMUP 개발자 포털',
  description: '오픈 API 명세서 및 실시간 샌드박스',
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
