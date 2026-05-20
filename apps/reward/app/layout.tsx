import './globals.css';

export const metadata = {
  title: 'MMUP 환자 리워드 풀',
  description: '개인 데이터 주권 및 토큰 보상',
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
