import './globals.css';

export const metadata = {
  title: 'MMUP 모바일 하이브리드',
  description: 'Flutter Native 통신 연동용 PWA Fallback',
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
