import './globals.css';

export const metadata = {
  title: 'MMUP AI 에이전트 허브',
  description: '의료 AI 모델 오케스트레이션 및 상태 관리',
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
